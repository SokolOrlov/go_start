package kafkaConsumer

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"qwe/internal/driver/config"
	common "qwe/pkg/models"
	"strings"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
)

func Run(consumeCh chan common.Trip, log *slog.Logger, cfg *config.KAFKA) error {
	keepRunning := true
	log.Info("Starting a new Sarama consumer")

	// if verbose {
	// 	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	// }

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

	// sigusr1 := make(chan os.Signal, 1)
	// signal.Notify(sigusr1, syscall.SIGUSR1)

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
			// case <-sigusr1:
			// 	toggleConsumptionFlow(client, &consumptionIsPaused)
		}
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Error("Error closing client: %v", err)
	}

	return nil
}

// func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
// 	if *isPaused {
// 		client.ResumeAll()
// 		log.Println("Resuming consumption")
// 	} else {
// 		client.PauseAll()
// 		log.Println("Pausing consumption")
// 	}

// 	*isPaused = !*isPaused
// }

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
	ch    chan common.Trip
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
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}
			m := common.Trip{}
			json.Unmarshal(message.Value, &m)

			consumer.ch <- m
			session.MarkMessage(message, "")
		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}
