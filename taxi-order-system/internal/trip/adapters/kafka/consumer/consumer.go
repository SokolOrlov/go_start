package kafkaConsumer

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"qwe/internal/trip/config"
	common "qwe/pkg/models"
	"strings"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
)

func Run(consumeCh chan common.KafkaMessage, log *slog.Logger, cfg *config.KAFKA) error {
	keepRunning := true
	log.Info("Starting a new Sarama consumer")

	version, err := sarama.ParseKafkaVersion(sarama.DefaultVersion.String())
	if err != nil {
		return err
	}

	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	config := sarama.NewConfig()
	config.Version = version

	switch cfg.CONSUMER.ASSIGNOR {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	case "roundrobin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	default:
		return errors.New("Unrecognized consumer group partition assignor: " + cfg.CONSUMER.ASSIGNOR)
	}

	if cfg.CONSUMER.OLDEST {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	/**
	 * Setup a new Sarama consumer group
	 */
	consumer := Consumer{
		ready: make(chan bool),
		ch:    consumeCh,
		log:   log,
	}

	ctx, cancel := context.WithCancel(context.Background())
	var brokers = []string{strings.Join([]string{cfg.BROKER.URL, cfg.BROKER.PORT}, ":")}
	client, err := sarama.NewConsumerGroup(brokers, cfg.CONSUMER.GROUP, config)
	if err != nil {
		return errors.New("Error creating consumer group client: " + err.Error())
	}

	// consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, []string{cfg.CONSUMER.TOPIC}, &consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Error("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	log.Info("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Info("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Info("terminating: via signal")
			keepRunning = false
			os.Exit(0)
		}
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Error("Error closing client: %v", err)
	}

	return nil
}
