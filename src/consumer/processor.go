package consumer

import (
	"fmt"
	"github.com/Shopify/sarama"
	"go-ka/config"
	"sync"
	"time"
)

type Consumer struct {
	config        config.ProcessorConfig
	client        *sarama.Consumer
	worker        map[int32]*sarama.PartitionConsumer
	topic         string
	allPartition  []int32
	livePartition []int32 ///[]int32
	deadPartition []int32
	mutex         *sync.Mutex /**guarantee atomic for consumer*/

}

type Process struct {
	configs   config.ProcessorConfigs
	consumers map[string]*Consumer
}

type ProcessImpl interface {
	Consume() int
}

func NewProcess(cfgs config.ProcessorConfigs) *Process {
	//consumer := newConsumer(cfg)
	return &Process{
		configs:   cfgs,
		consumers: newConsumers(cfgs),
	}
}

func newConsumers(cfgs config.ProcessorConfigs) map[string]*Consumer {
	//var consumers []*sarama.Consumer

	var retv map[string]*Consumer

	for k, v := range cfgs.Processors {

		config := sarama.NewConfig()
		config.Consumer.Return.Errors = true
		config.Consumer.Fetch.Max = 50
		config.Consumer.MaxProcessingTime = time.Duration(v.PollTimeout / 1000 / 1000) //milli to nao
		c, err := sarama.NewConsumer([]string{v.BoostrapServer}, config)
		partitions, err2 := c.Partitions(v.Topic)

		if err != nil {
			panic(err)
		}

		if err2 != nil {
			panic(err2)
		}

		csm := &Consumer{
			config:        v,
			client:        &c,
			worker:        toMap(partitions),
			topic:         v.Topic,
			allPartition:  partitions,
			livePartition: []int32{}, ///[]int32
			deadPartition: partitions,
			mutex:         &sync.Mutex{}, /**guarantee atomic for consumer*/

		}
		retv[k] = csm
	}

	//consumers = append(consumers, c)

	return retv
}

func toMap(nums []int32) map[int32]*sarama.PartitionConsumer {
	var retv map[int32]*sarama.PartitionConsumer

	for _, v := range nums {
		retv[v] = nil
	}
	return retv
}

func (p *Process) Consume() int {
	for _, v := range p.consumers {
		v.mutex.Lock()
		for _, itr := range v.deadPartition {

			(*v.worker[itr]).Messages().
		}
		v.mutex.Unlock()
	}
	/*Count to execute*/
	/*cnt := p.DeadCount
	p.consumers.mutex.Lock()

	for i := 0; i < cnt; i++ {

		//Create one consumer
		c := newConsumer(p.config)

		//register to consumer
		p.consumers = append(p.consumers, c)
		fmt.Println(p.consumers)

		(*c).Partitions()
		go func() {
			for {
				//start to consume

				ev := (*c).ConsumePartition()
				v1, _ := (*c).ConsumePartition()
				v1.Messages()
				switch e := ev.(type) {
				case *sarama.Message:
					fmt.Printf("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))


					//  TODO place for some action

				case sarama.Error:
					//fmt.Println(p.consumers)
					p.DeadCount += 1
					p.removeObject(c) //Remove consumer from list

					err2 := c.Close()
					if err2 != nil {
						fmt.Println(err2)
					}
					fmt.Println(e)
				}
			}
		}()

	}
	p.mutex.Unlock()
	p.DeadCount = p.config.Concurrency - len(p.consumers)

	if cnt > 0 {

		return cnt
	} else {
		return 0
	}*/

}

func (p *Process) removeObject(target *sarama.Consumer) {

	fmt.Println(target)
	var newValue []*sarama.Consumer

	for _, v := range p.consumers {
		if v != target {
			newValue = append(newValue, v)
		}
	}
	p.consumers = newValue
}
