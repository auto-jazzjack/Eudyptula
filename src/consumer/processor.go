package consumer

import (
	"config"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type Process struct {
	config   config.ProcessorConfig
	consumer *kafka.Consumer
}

type ProcessImpl interface {
}

func NewProcess(cfg config.ProcessorConfig) *Process {
	return &Process{
		config:   cfg,
		consumer: newConsumer(cfg),
	}
}

func newConsumer(cfg config.ProcessorConfig) *kafka.Consumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.BoostrapServer,
		"group.id":          cfg.GroupId,
		"auto.offset.reset": cfg.Offset,
	})

	if err != nil {
		panic(err)
	}
	return consumer
}

func Consume() {

}