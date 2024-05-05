package trip

import (
	"context"
	"fmt"
	"log/slog"
	kafkaConsumer "qwe/internal/trip/adapters/kafka/consumer"
	kafkaProducer "qwe/internal/trip/adapters/kafka/producer"
	repository "qwe/internal/trip/adapters/postgres"
	"qwe/internal/trip/config"
	"qwe/internal/trip/models"
	"qwe/internal/trip/service"
	common "qwe/pkg/models"
	"time"
)

type App struct {
	log       *slog.Logger
	cfg       *config.Config
	produceCh chan models.Message
	consumeCh chan common.KafkaMessage
	service   *service.Service
}

func (a *App) Start() error {
	const op = "Trip.Start"

	log := a.log.With(slog.String("op", op))

	log.Info("starting Trip service...")

	err := kafkaProducer.Run(a.produceCh, a.log, &a.cfg.KAFKA)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = kafkaConsumer.Run(a.consumeCh, a.log, &a.cfg.KAFKA)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.service.Run()

	return nil
}

func (a *App) Stop() {
	const op = "client.Stop"

	log := a.log.With(slog.String("op", op))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	a.service.Stop()

	<-ctx.Done()
	log.Info("bye.")
}

func New(cfg *config.Config, log *slog.Logger) (*App, error) {
	db, err := repository.New(&cfg.DB)

	if err != nil {
		return nil, err
	}

	produceCh := make(chan models.Message)
	consumeCh := make(chan common.KafkaMessage)

	svc := service.New(cfg, log, db, produceCh, consumeCh)

	return &App{
		log:       log,
		cfg:       cfg,
		produceCh: produceCh,
		consumeCh: consumeCh,
		service:   svc,
	}, nil
}
