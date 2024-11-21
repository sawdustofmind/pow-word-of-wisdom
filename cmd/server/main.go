package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/sawdustofmind/pow-word-of-wisdom/internal/config"
	"github.com/sawdustofmind/pow-word-of-wisdom/internal/log"
	"github.com/sawdustofmind/pow-word-of-wisdom/internal/quotes"
	"github.com/sawdustofmind/pow-word-of-wisdom/internal/server"
)

func main() {

	// loading config from file and env
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal("failed to read config", zap.Error(err))
	}

	// run server
	qts := quotes.NewInMemoryQuoteStore()
	srv := server.NewServer(cfg, qts)

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	err = srv.Start(ctx)
	if err != nil {
		log.Fatal("start server", zap.Error(err))
	}

	exitSignal := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(exitSignal, os.Interrupt, syscall.SIGTERM)

	<-exitSignal
	log.Info("Shutting down")
	cancelFn()

	time.Sleep(5 * time.Second) // waiting till clients reads its messages
	srv.Close(ctx)
}
