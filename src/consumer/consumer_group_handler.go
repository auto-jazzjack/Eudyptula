package consumer

import (
	"go-ka/logic"

	"github.com/Shopify/sarama"
)

type ConsumerGroupHandlerImpl struct {
	logic   logic.Logic[any]
	topic   string
	session *sarama.ConsumerGroupSession
}

func NewConsumerGroupHandler() ConsumerGroupHandlerImpl {
	return ConsumerGroupHandlerImpl{}
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

	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
// but before the offsets are committed for the very last time.
func (c ConsumerGroupHandlerImpl) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
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
