package server

import (
	"context"
	"fmt"
	"net"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/sawdustofmind/pow-word-of-wisdom/internal/config"
	"github.com/sawdustofmind/pow-word-of-wisdom/internal/log"
	"github.com/sawdustofmind/pow-word-of-wisdom/internal/quotes"
)

type Server struct {
	cfg           *config.Config
	quotesService quotes.QuoteService
	logger        *zap.Logger

	listener net.Listener
}

func NewServer(config *config.Config, quotesService quotes.QuoteService) *Server {
	return &Server{
		cfg:           config,
		quotesService: quotesService,
		logger:        log.GetLogger().With(zap.String("component", "server")),
	}
}

func (s *Server) Start(ctx context.Context) error {
	address := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	s.listener = listener
	s.logger.Info("start listening", zap.String("address", listener.Addr().String()))

	go func() {
		defer listener.Close()

		for {
			// Listen for an incoming connection.
			conn, err := listener.Accept()
			if err != nil {
				s.logger.Error("accept connection", zap.Error(err))
				return
			}

			// Handle connections in a new goroutine.
			client := newClient(uuid.New().String(), s.cfg, conn, s.quotesService)
			client.start()
		}
	}()

	return nil
}

func (s *Server) Close(ctx context.Context) {
	err := s.listener.Close()
	if err != nil {
		s.logger.Error("close listener", zap.Error(err))
	}
}
