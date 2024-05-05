package driver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	gohttp "net/http"
	"qwe/internal/driver/adapters/http"
	kafkaConsumer "qwe/internal/driver/adapters/kafka/consumer"
	kafkaProducer "qwe/internal/driver/adapters/kafka/producer"
	repository "qwe/internal/driver/adapters/postgres"
	"qwe/internal/driver/config"
	"qwe/internal/driver/service"
	common "qwe/pkg/models"
	"time"
)

type App struct {
	http      http.Adapter
	log       *slog.Logger
	cfg       *config.Config
	produceCh chan common.UpdateTripStatus
	consumeCh chan common.Trip
	service   *service.Service
}

func (a *App) Start() error {
	const op = "client.Start"

	log := a.log.With(slog.String("op", op))

	log.Info("starting Client service...")

	err := a.startHttp(log)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = kafkaProducer.Run(a.produceCh, a.log, &a.cfg.KAFKA)
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

	log.Info("stopping http server", slog.String("port", a.cfg.HTTP.PORT))

	a.http.Shutdown(ctx)

	a.service.Stop()

	<-ctx.Done()
	log.Info("bye.")
}

func New(cfg *config.Config, log *slog.Logger) (*App, error) {
	db, err := repository.New(&cfg.DB)

	if err != nil {
		return nil, err
	}

	produceCh := make(chan common.UpdateTripStatus)
	consumeCh := make(chan common.Trip)

	svc := service.New(db, log, produceCh, consumeCh)
	httpserver := http.New(svc, &cfg.HTTP, log)

	return &App{
		http:      httpserver,
		log:       log,
		cfg:       cfg,
		produceCh: produceCh,
		consumeCh: consumeCh,
		service:   svc,
	}, nil
}

func (app *App) startHttp(log *slog.Logger) error {
	go func() {
		if err := app.http.Serve(); err != nil && !errors.Is(err, gohttp.ErrServerClosed) {
			log.Error("Could not listen on port :%s: %s\n", app.cfg.HTTP.PORT, err.Error())
		}
	}()

	log.Info("http server is running", slog.String("port", app.cfg.HTTP.PORT))
	return nil
}
