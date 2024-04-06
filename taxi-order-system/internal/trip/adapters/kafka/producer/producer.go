package kafkaProducer

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"qwe/internal/trip/config"
	"qwe/internal/trip/models"
	"strings"

	"github.com/IBM/sarama"
)

func Run(produceCh chan models.Message, log *slog.Logger, cfg *config.KAFKA) error {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	var Brokers = []string{strings.Join([]string{cfg.BROKER.URL, cfg.BROKER.PORT}, ":")}
	producer, err := sarama.NewSyncProducer(Brokers, config)

	go func(producer *sarama.SyncProducer, log *slog.Logger) {
		for {
			trip := <-produceCh
			mess := PrepareMessage(trip.Topic, trip.Model)
			SendMessage(mess, producer, log)
		}

	}(&producer, log)
	return err
}

func SendMessage(msg *sarama.ProducerMessage, producer *sarama.SyncProducer, log *slog.Logger) {
	partition, offset, err := (*producer).SendMessage(msg)

	fmt.Println(partition, offset, err)

	if err != nil {
		log.Error("error occured %v", err)
	} else {
		log.Info("Message was saved to partion: %d.\nMessage offset is: %d.\n", slog.Int("partition", int(partition)), slog.Int64("offset", offset))
	}
}

func PrepareMessage(topic string, message interface{}) *sarama.ProducerMessage {
	bytes, _ := json.Marshal(message)
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.StringEncoder(bytes),
	}

	return msg
}
