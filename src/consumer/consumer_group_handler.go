package consumer

import (
	"fmt"
	"go-ka/logic"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

type ConsumerGroupHandlerImpl struct {
	logic   logic.Logic[any]
	topic   string
	session *sarama.ConsumerGroupSession
	wg      sync.WaitGroup
}

func NewConsumerGroupHandler(logic logic.Logic[any], topic string) ConsumerGroupHandlerImpl {
	one := sync.WaitGroup{}
	one.Add(1)
	return ConsumerGroupHandlerImpl{
		logic: logic,
		topic: topic,
		wg:    one,
	}
}

func (c ConsumerGroupHandlerImpl) GetPartitons() []int32 {
	if c.session != nil && c.topic != "" && len((*c.session).Claims()) > 0 {
		retv := (*c.session).Claims()[c.topic]
		return retv
	}

	return []int32{}
}
func (c *ConsumerGroupHandlerImpl) Setup(session sarama.ConsumerGroupSession) error {
	c.session = &session

	c.wg.Done()
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
// but before the offsets are committed for the very last time.
func (c ConsumerGroupHandlerImpl) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c ConsumerGroupHandlerImpl) ConsumeOffset(offset time.Time) map[string]string {

	partition := c.GetPartitons()

	retv := make(map[string]string)
	epoch := offset.Unix()

	for i := 0; i < len(partition); i++ {
		(*c.session).ResetOffset(c.topic, partition[i], epoch, "")
		retv[fmt.Sprintf("%v", partition[i])] = fmt.Sprintf("%v", epoch)
	}
	return retv
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (c ConsumerGroupHandlerImpl) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for message := range claim.Messages() {

		value := c.logic.Deserialize(message.Value)
		err := c.logic.DoAction(value)
		if err != nil {
			return err
		}
		session.Commit()
	}

	return nil
}
