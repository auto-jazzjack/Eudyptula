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

func NewProcess(cfgs *config.ProcessorConfigs) *Process {
	//consumer := newConsumer(cfg)
	return &Process{
		configs:   *cfgs,
		consumers: newConsumers(cfgs),
	}
}

func newConsumers(cfgs *config.ProcessorConfigs) map[string]*Consumer {
	//var consumers []*sarama.Consumer

	var retv = make(map[string]*Consumer)

	for k, v := range cfgs.Processors {

		config := sarama.NewConfig()
		config.Consumer.Return.Errors = true
		config.Consumer.Fetch.Max = 50
		config.Consumer.MaxProcessingTime = time.Duration(v.PollTimeout * 1000 * 1000) //milli to nao

		c, err := sarama.NewConsumer([]string{v.BoostrapServer}, config)

		if err != nil {
			panic(err)
		}

		partitions, err2 := c.Partitions(v.Topic)

		if err2 != nil {
			panic(err2)
		}

		csm := &Consumer{
			config:        v,
			client:        &c,
			worker:        toMap(partitions), //all dead(nil) for init
			topic:         v.Topic,
			allPartition:  partitions,
			livePartition: []int32{}, ///[]int32 empty live
			deadPartition: partitions,
			mutex:         &sync.Mutex{}, /**guarantee atomic for consumer*/

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

func (p *Process) Consume() map[int32]int {
	retv := make(map[int32]int)
	for _, v := range p.consumers {

		cpy := v.deadPartition
		for idx, partitionNum := range cpy {
			v.mutex.Lock()
			c, err := (*v.client).ConsumePartition(v.topic, partitionNum, sarama.OffsetOldest)

			if err != nil {
				fmt.Println(err)
				continue
			} else {
				//success to revive dead Partition
				v.deadPartition = append(v.deadPartition[:idx], v.deadPartition[idx:]...)
				v.livePartition = append(v.livePartition, partitionNum)
				v.worker[partitionNum] = &c
				num := partitionNum

				//retv[partitionNum] = 1

				//fmt.Println(*v.worker[partitionNum])
				go func() {

					fmt.Println(num)
					for {
						select {
						case msg1 := <-(*v.worker[num]).Messages():
							fmt.Println("received", msg1.Value)
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

/*func (p *Process) removeObject(target *sarama.Consumer) {

	fmt.Println(target)
	var newValue []*sarama.Consumer

	for _, v := range p.consumers {
		if v != target {
			newValue = append(newValue, v)
		}
	}
	p.consumers = newValue
}*/
