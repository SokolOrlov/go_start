package kafkaConsumer

import (
	"encoding/json"
	"log"

	common "qwe/pkg/models"

	"github.com/IBM/sarama"
)

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
	ch    chan common.UpdateTripStatus
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}
			m := common.UpdateTripStatus{}
			json.Unmarshal(message.Value, &m)

			consumer.ch <- m
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
