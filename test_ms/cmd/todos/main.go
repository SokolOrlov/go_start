package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"test_ms/internal/app/todos"
	"test_ms/internal/app/todos/config"
	"time"

	logBase "log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		logBase.Fatalln(err)
	}

	ctxWithCancel, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctxSys := newSystemContext(ctxWithCancel, 5*time.Second, newLogSystemSignalCallback())

	app := todos.NewApp()

	if err := app.Init(ctxSys, cfg); err != nil {
		logBase.Fatal("start app failed")
	}

	if err := app.Start(ctxSys); err != nil {
		logBase.Fatal("start app failed")
	}

	if err := app.Stop(ctxSys); err != nil {
		logBase.Fatalf("stop app failed")
	}
}

type Callback func(signal os.Signal)

// NewSystemContext returns new Context, which will be cancelled on receiving SIGTERM and SIGINT signals after supplied delay.
// Additionally multiple Callback functions can be passed, they will be called immediately after receiving signals, before delay.
func newSystemContext(ctx context.Context, delay time.Duration, callbacks ...Callback) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

		sig := <-sigint
		for _, cb := range callbacks {
			go cb(sig)
		}

		time.Sleep(delay)

		cancel()
	}()

	return ctx
}

func newLogSystemSignalCallback() Callback {
	return func(signal os.Signal) {
		logBase.Printf("system signal %d (%s) received, context will be canceled shortly\n", signal, signal.String())
	}
}
