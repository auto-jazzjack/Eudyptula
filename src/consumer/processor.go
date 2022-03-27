package consumer

import (
	"config"
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"sync"
)

type Process struct {
	config    config.ProcessorConfig
	consumers []*kafka.Consumer /**Only running cousumer*/
	DeadCount int               /**TODO should be atomic*/
	mutex     *sync.Mutex       /**guarantee atomic for consumer*/
}

type ProcessImpl interface {
	Consume()
}

func NewProcess(cfg config.ProcessorConfig) *Process {
	//consumer := newConsumer(cfg)
	return &Process{
		config:    cfg,
		consumers: []*kafka.Consumer{},
		/*For the init status, all processor just created and not executed*/
		DeadCount: cfg.Concurrency,
		mutex:     &sync.Mutex{},
	}
}

func newConsumer(cfg config.ProcessorConfig) *kafka.Consumer {
	//var consumers []*kafka.Consumer

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
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

func (p *Process) Consume() {
	/**Count to execute*/
	cnt := p.config.Concurrency - p.DeadCount

	for i := 0; i < cnt; i++ {

		//Create one consumer
		c := newConsumer(p.config)

		//register to consumer
		p.mutex.Lock()
		p.consumers = append(p.consumers, c)
		p.mutex.Unlock()

		go func() {
			for {
				//start to consume
				ev := c.Poll(p.config.PollTimeout)
				switch e := ev.(type) {
				case *kafka.Message:
					fmt.Printf("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))

					/**
					TODO place for some action
					*/
				case kafka.Error:
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
}

func (p *Process) removeObject(target *kafka.Consumer) {

	p.mutex.Lock()
	var newValue []*kafka.Consumer

	for _, v := range p.consumers {
		if v != target {
			newValue = append(newValue, v)
		}
	}
	p.consumers = newValue
	p.mutex.Unlock()
}
