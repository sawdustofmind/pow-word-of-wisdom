package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/sawdustofmind/pow-word-of-wisdom/internal/client"
	"github.com/sawdustofmind/pow-word-of-wisdom/internal/config"
	"github.com/sawdustofmind/pow-word-of-wisdom/internal/log"
)

func main() {
	log.Error("start client")

	// loading config from file and env
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal("error load config:", zap.Error(err))
		return
	}

	cl, err := client.NewClient(cfg)
	if err != nil {
		log.Error("error create client:", zap.Error(err))
		return
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	cl.Start(ctx)

	exitSignal := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(exitSignal, os.Interrupt, syscall.SIGTERM)

	<-exitSignal
	log.Info("Shutting down")
	cancelFn()

	_ = cl.Close()
}
