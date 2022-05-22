package consumer

import (
	"context"
	"go-ka/config"
	"go-ka/logic"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Shopify/sarama"
)

type Consumer[V any] struct {
	config        config.ProcessorConfig[V]
	groupId       string
	topic         string
	live          int32
	concurrency   int32
	client        sarama.ConsumerGroup
	logic         logic.Logic[any]
	livePartition map[string]bool
}

type Process[V any] struct {
	configs config.ProcessorConfigs[V]

	consumers map[string]*Consumer[V]
}

type ProcessImpl interface {
	Consume() int
	Rewind(time.Time) map[string][]string
}

func NewProcess[V any](cfgs *config.ProcessorConfigs[V]) *Process[V] {

	return &Process[V]{
		configs:   *cfgs,
		consumers: newConsumers(cfgs, cfgs.Zookeeper),
	}
}

func newConsumers[V any](cfgs *config.ProcessorConfigs[V], zkper []string) map[string]*Consumer[V] {

	var retv = make(map[string]*Consumer[V])
	for k, v := range cfgs.Processors {

		newConfig := sarama.NewConfig()

		newConfig.Consumer.Return.Errors = false
		newConfig.Consumer.Fetch.Max = 15
		newConfig.Consumer.Offsets.Initial = -1
		newConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
		newConfig.Consumer.MaxProcessingTime = time.Duration(v.PollTimeout * 1000 * 1000) //milli to nao

		//If userName is not empty we can suppose that sasl is enabled
		if v.UserName != "" {
			newConfig.Net.SASL.Password = v.Password
			newConfig.Net.SASL.Enable = true
			newConfig.Net.SASL.User = v.UserName
			newConfig.Net.SASL.Handshake = true

			if v.Algorithm == "SCRAM-SHA-256" {
				newConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
				newConfig.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
					return &XDGSCRAMClient{
						HashGeneratorFcn: SHA256,
					}

				}
			} else if v.Algorithm == "SCRAM-SHA-512" {
				newConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
				newConfig.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
					return &XDGSCRAMClient{
						HashGeneratorFcn: SHA512,
					}

				}
			}

		}

		client, err := sarama.NewConsumerGroup([]string{v.BoostrapServer}, v.GroupId, newConfig)
		if err != nil {
			panic(err)
		}

		csm := &Consumer[V]{
			config:      v,
			groupId:     v.GroupId,
			client:      client,
			topic:       v.Topic,
			live:        0,
			concurrency: v.Concurrency,
			logic:       logic.Logic[any](v.LogicContainer.Logic),
		}
		retv[k] = csm
	}

	return retv
}

/**
Reqeust : target time stamp
Return : Key : consumerName, Value : partition
*/
func (p *Process[V]) Rewind(date time.Time) map[string][]string {
	return nil

}

func (p *Process[V]) Consume() map[string]int32 {
	retv := make(map[string]int32)
	for _, v := range p.consumers {

		numToRevive := v.concurrency - v.live
		if numToRevive > 0 {
			retv["topic:"+v.topic+" group id : "+v.groupId] = numToRevive
		}

		for i := int32(0); i < numToRevive; i++ {

			wg := &sync.WaitGroup{}
			wg.Add(1)
			consumer := ConsumerGroupHandlerImpl{
				ready: make(chan bool),
				logic: v.logic,
			}
			go func() {
				atomic.AddInt32(&numToRevive, 1)

				defer func() {
					wg.Done()
					//if go func finished, let's assume consumer dead.
					atomic.AddInt32(&numToRevive, -1)
				}()

				for {

					ctx, err := context.WithCancel(context.Background())
					if err != nil {
						break

					}

					if err := v.client.Consume(ctx, strings.Split(v.topic, ","), &consumer); err != nil {
						log.Panicf("Error from consumer: %v", err)
					}
					if ctx.Err() != nil {
						return
					}
				}
			}()

			<-consumer.ready // Await till the consumer has been set up

		}

	}
	return retv

}
