package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"qwe/internal/trip"
	"qwe/internal/trip/config"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("app error: %s", err)
		os.Exit(1)
	}

	log := initLogger("dev")

	app, err := trip.New(cfg, log)

	if err != nil {
		log.Error("app error", err)
		os.Exit(1)
	}
	app.Start()

	<-ctx.Done()

	app.Stop()
}

func initLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "dev":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	}

	return log
}
