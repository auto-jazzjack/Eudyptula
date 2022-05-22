package consumer

import (
	"go-ka/logic"

	"github.com/Shopify/sarama"
)

type ConsumerGroupHandlerImpl struct {
	ready chan bool
	logic logic.Logic[any]
}

func NewConsumerGroupHandler() ConsumerGroupHandlerImpl {
	return ConsumerGroupHandlerImpl{}
}

func (c ConsumerGroupHandlerImpl) Setup(sarama.ConsumerGroupSession) error {
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
		//log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		//session.MarkMessage(message, "")
		value := c.logic.Deserialize(message.Value)
		c.logic.DoAction(value)
	}

	return nil
}
