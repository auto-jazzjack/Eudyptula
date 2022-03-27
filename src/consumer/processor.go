package consumer

import (
	"config"
	"errors"
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"os"
)

type Process struct {
	config    config.ProcessorConfig
	consumers []*kafka.Consumer
	DeadCount int
}

type ProcessImpl interface {
	Consume()
}

func NewProcess(cfg config.ProcessorConfig) *Process {
	consumer := newConsumer(cfg)
	return &Process{
		config:    cfg,
		consumers: consumer,
		/*For the init status, all processor just created and not executed*/
		DeadCount: len(consumer),
	}
}

func newConsumer(cfg config.ProcessorConfig) []*kafka.Consumer {
	var consumers []*kafka.Consumer

	for i := 0; i < cfg.Concurrency; i++ {
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
		consumers = append(consumers, c)
	}
	return consumers
}

func (p *Process) Consume() error {

	//P의 concurrency 만큼 실행해야함.

	for {
		ev := p.consumer.Poll(p.config.PollTimeout)
		switch e := ev.(type) {
		case *kafka.Message:
			fmt.Printf("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)

			/**Manager will get error from this when kafka dead*/
			return errors.New("kafka dead")
		}
	}
}
