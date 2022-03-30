package consumer

import (
	"fmt"
	"github.com/Shopify/sarama"
	"go-ka/config"
	"sync"
)

type Process[V interface{}] struct {
	config    config.ProcessorConfig
	consumers []*sarama.Consumer /**Only running cousumer*/
	DeadCount int                /**TODO should be atomic*/
	mutex     *sync.Mutex        /**guarantee atomic for consumer*/
	executor  Executor[V]
}

type ProcessImpl interface {
	Consume() int
}

func NewProcess[V any](cfg config.ProcessorConfig, executor Executor[V]) *Process[V] {
	//consumer := newConsumer(cfg)
	return &Process[V]{
		config:    cfg,
		consumers: []*sarama.Consumer{},
		/*For the init status, all processor just created and not executed*/
		DeadCount: cfg.Concurrency,
		mutex:     &sync.Mutex{},
		executor:  executor,
	}
}

func newConsumer(cfg config.ProcessorConfig) *sarama.Consumer {
	//var consumers []*kafka.Consumer

	c, err := sarama.NewConsumer(&sarama.Config{
		"bootstrap.servers": cfg.BoostrapServer,
		"group.id":          cfg.GroupId,
		"auto.offset.reset": cfg.Offset,
	})

	if err != nil {
		panic(err)
	}

	err2 := c.SubscribeTopics([]string{cfg.Topic}, nil)
	if err2 != nil {
		panic(err2)
	}
	//consumers = append(consumers, c)

	return c
}

func (p *Process[V]) Consume() int {
	/**Count to execute*/
	cnt := p.DeadCount
	p.mutex.Lock()

	for i := 0; i < cnt; i++ {

		//Create one consumer
		c := newConsumer(p.config)

		//register to consumer
		p.consumers = append(p.consumers, c)
		fmt.Println(p.consumers)

		go func() {
			for {
				//start to consume
				ev := c.Pool(p.config.PollTimeout)
				switch e := ev.(type) {
				case *sarama.Message:
					fmt.Printf("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
					defaultValue := p.executor.DefaultValue()
					res := p.executor.DeSerialize(e.Value, defaultValue)
					p.executor.DoAction(res)

					/**
					TODO place for some action
					*/
				case sarama.ConsumerError:
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
	}

}

func (p *Process[V]) removeObject(target *sarama.Consumer) {

	fmt.Println(target)
	var newValue []*sarama.Consumer

	for _, v := range p.consumers {
		if v != target {
			newValue = append(newValue, v)
		}
	}
	p.consumers = newValue
}
