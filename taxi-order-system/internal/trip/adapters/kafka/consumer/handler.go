package kafkaConsumer

import (
	"encoding/json"
	"log/slog"

	common "qwe/pkg/models"

	"github.com/IBM/sarama"
)

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
	ch    chan common.KafkaMessage
	log   *slog.Logger
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

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				consumer.log.Error("message channel was closed")
				return nil
			}
			consumer.log.Debug("Message claimed", slog.String("model", string(message.Value)), slog.String("topic", message.Topic))

			var m common.KafkaMessage
			json.Unmarshal(message.Value, &m)
			consumer.ch <- m
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}
