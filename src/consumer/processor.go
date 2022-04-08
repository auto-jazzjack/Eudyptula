package consumer

import (
	"fmt"
	"github.com/Shopify/sarama"
	"go-ka/config"
	"go-ka/logic"
	"go-ka/util"
	"sync"
	"time"
)

type Consumer[V any] struct {
	config        config.ProcessorConfig[V]
	client        *sarama.Consumer
	worker        map[int32]*sarama.PartitionConsumer
	topic         string
	allPartition  []int32
	livePartition []int32 ///[]int32
	deadPartition []int32
	mutex         *sync.Mutex /**guarantee atomic for consumer*/
	logic         logic.Logic[any]
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
		consumers: newConsumers[V](cfgs),
	}
}

func newConsumers[V any](cfgs *config.ProcessorConfigs[V]) map[string]*Consumer[V] {

	var retv = make(map[string]*Consumer[V])

	for k, v := range cfgs.Processors {

		newConfig := sarama.NewConfig()
		newConfig.Consumer.Return.Errors = true
		newConfig.Consumer.Fetch.Max = 50
		newConfig.Consumer.MaxProcessingTime = time.Duration(v.PollTimeout * 1000 * 1000) //milli to nao

		c, err := sarama.NewConsumer([]string{v.BoostrapServer}, newConfig)

		if err != nil {
			panic(err)
		}

		partitions, err2 := c.Partitions(v.Topic)

		if err2 != nil {
			panic(err2)
		}

		csm := &Consumer[V]{
			config:        v,
			client:        &c,
			worker:        toMap(partitions), //all dead(nil) for init
			topic:         v.Topic,
			allPartition:  partitions,
			livePartition: []int32{}, ///[]int32 empty live
			deadPartition: partitions,
			mutex:         &sync.Mutex{}, /**guarantee atomic for consumer*/
			logic:         logic.Logic[any](v.LogicContainer.Logic),
		}
		retv[k] = csm
	}

	//consumers = append(consumers, c)

	return retv
}

func toMap(nums []int32) map[int32]*sarama.PartitionConsumer {
	var retv = make(map[int32]*sarama.PartitionConsumer)

	for _, v := range nums {
		retv[v] = nil
	}
	return retv
}

func (p *Process[V]) Consume() map[int32]int {
	retv := make(map[int32]int)
	for _, v := range p.consumers {

		cpy := v.deadPartition
		for _, partitionNum := range cpy {
			v.mutex.Lock()
			c, err := (*v.client).ConsumePartition(v.topic, partitionNum, sarama.OffsetOldest)

			if err != nil {
				fmt.Println(err)
				continue
			} else {
				//success to revive dead Partition
				v.deadPartition = util.FilterExactValue(v.deadPartition, partitionNum)
				v.livePartition = append(v.livePartition, partitionNum)
				v.worker[partitionNum] = &c
				num := partitionNum

				go func() {

					fmt.Println(num)
					for {
						select {
						case msg1 := <-(*v.worker[num]).Messages():
							fmt.Println("received", msg1.Value)
							res := (v.logic.Deserialize)(msg1.Value)
							//fmt.Println(res)
							v.logic.DoAction(res)
						case msg1 := <-(*v.worker[num]).Errors():
							fmt.Println("error", msg1)
							*v.worker[num] = nil
							v.livePartition = append(v.livePartition[:1], v.livePartition[2:]...)
							break

						}
					}
				}()

			}
			v.mutex.Unlock()
		}

	}
	return retv

}
