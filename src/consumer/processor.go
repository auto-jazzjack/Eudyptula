package consumer

import (
	"fmt"
	cluster "github.com/bsm/sarama-cluster"
	"go-ka/config"
	"go-ka/logic"
	"sync"
	"sync/atomic"
	"time"
)

type Consumer[V any] struct {
	config      config.ProcessorConfig[V]
	worker      *cluster.Consumer
	groupId     string
	topic       string
	live        int32
	concurrency int32
	mutex       *sync.Mutex /**guarantee atomic for consumer*/
	logic       logic.Logic[any]
}

type Process[V any] struct {
	configs   config.ProcessorConfigs[V]
	consumers map[string]*Consumer[V]
}

type ProcessImpl interface {
	Consume() int
}

func NewProcess[V any](cfgs *config.ProcessorConfigs[V]) *Process[V] {
	return &Process[V]{
		configs:   *cfgs,
		consumers: newConsumers[V](cfgs, cfgs.Zookeeper),
	}
}

func newConsumers[V any](cfgs *config.ProcessorConfigs[V], zkper []string) map[string]*Consumer[V] {

	var retv = make(map[string]*Consumer[V])

	for k, v := range cfgs.Processors {

		newConfig := &cluster.Config{}
		newConfig.Consumer.Return.Errors = true
		newConfig.Consumer.Fetch.Max = v.FetchSize
		//cfg..Return.Notifications = true
		newConfig.Consumer.MaxProcessingTime = time.Duration(v.PollTimeout * 1000 * 1000) //milli to nao

		c, err := cluster.NewConsumer([]string{v.BoostrapServer}, zkper, v.GroupId, []string{v.Topic}, newConfig)

		if err != nil {
			panic(err)
		}

		/*if err2 != nil {
			panic(err2)
		}*/

		csm := &Consumer[V]{
			config: v,
			//client:        &c,
			worker:  c,
			groupId: v.GroupId,
			//worker:      toMap(partitions), //all dead(nil) for init
			topic:       v.Topic,
			live:        0,
			concurrency: v.Concurrency,
			mutex:       &sync.Mutex{}, /**guarantee atomic for consumer*/
			logic:       logic.Logic[any](v.LogicContainer.Logic),
		}
		retv[k] = csm
	}

	return retv
}

func toMap(nums []int32) map[int32]*cluster.Consumer {
	var retv = make(map[int32]*cluster.Consumer)

	for _, v := range nums {
		retv[v] = nil
	}
	return retv
}

func (p *Process[V]) Consume() map[int32]int {
	retv := make(map[int32]int)
	for _, v := range p.consumers {

		numToRevive := v.concurrency - v.live
		for i := int32(0); i < numToRevive; i++ {
			{

				go func() {

					for {
						select {
						case msg1 := <-(*v.worker).Messages():
							res := (v.logic.Deserialize)(msg1.Value)
							atomic.AddInt32(&v.live, 1)
							v.logic.DoAction(res)
						case msg1 := <-(*v.worker).Errors():
							fmt.Println("error", msg1)
							err := (*v.worker).Close()
							if err != nil {
								fmt.Printf("%s", err)
							}
							atomic.AddInt32(&v.live, -1)
							break

						}
					}
				}()

			}
		}

	}
	return retv

}
