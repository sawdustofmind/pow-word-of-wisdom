package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/sawdustofmind/pow-word-of-wisdom/internal/config"
	"github.com/sawdustofmind/pow-word-of-wisdom/internal/dto"
	"github.com/sawdustofmind/pow-word-of-wisdom/internal/log"
	"github.com/sawdustofmind/pow-word-of-wisdom/internal/pow"
)

type Client struct {
	id   string
	conn net.Conn

	cfg *config.Config

	logger *zap.Logger
}

func NewClient(cfg *config.Config) (*Client, error) {
	address := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	id := uuid.New().String()
	logger := log.With(
		zap.String("service", "client"),
		zap.String("id", id),
		zap.String("address", address),
	)

	logger.Info("connected")
	return &Client{
		id:     id,
		conn:   conn,
		cfg:    cfg,
		logger: logger,
	}, nil
}

func (c *Client) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				message, err := c.RequestQuote()
				if err != nil {
					c.logger.Error("request failed", zap.Error(err))
				}
				c.logger.Info("quote result", zap.String("message", message.Quote))
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (c *Client) RequestQuote() (dto.QuoteResponse, error) {
	challengeAnswer, err := c.ReqAndSolveChallenge()
	if err != nil {
		return dto.QuoteResponse{}, err
	}

	req := dto.ClientMessage{
		Type:             dto.GetQuoteClientMessageType,
		ChallengeCounter: challengeAnswer,
	}
	reqPayload, err := json.Marshal(req)
	if err != nil {
		return dto.QuoteResponse{}, err
	}

	_, err = c.conn.Write(reqPayload)
	if err != nil {
		return dto.QuoteResponse{}, err
	}

	resp, err := c.MustReadMessage()
	if err != nil {
		return dto.QuoteResponse{}, err
	}

	if resp.Type != dto.QuoteResponseType {
		return dto.QuoteResponse{}, fmt.Errorf("unexpeted response type: %s", resp.Type)
	}

	quoteResponse := dto.QuoteResponse{}
	err = json.Unmarshal(resp.Data, &quoteResponse)
	if err != nil {
		return dto.QuoteResponse{}, err
	}
	return quoteResponse, nil
}

func (c *Client) ReqAndSolveChallenge() (int, error) {
	start := time.Now()

	req := dto.ClientMessage{}
	req.Type = dto.ChallengeClientMessageType

	reqPayload, err := json.Marshal(req)
	if err != nil {
		return 0, err
	}

	_, err = c.conn.Write(reqPayload)
	if err != nil {
		return 0, err
	}

	resp, err := c.MustReadMessage()
	if err != nil {
		return 0, err
	}

	if resp.Type != dto.ChallengeResponseType {
		return 0, fmt.Errorf("unexpeted response type: %s", resp.Type)
	}

	challengeResponse := dto.ChallengeResponse{}
	err = json.Unmarshal(resp.Data, &challengeResponse)
	if err != nil {
		return 0, err
	}

	_, answer, err := pow.ComputeHashcash(challengeResponse.Challenge, c.cfg.Hashcash.MaxIterations, c.cfg.Hashcash.ZerosCount)
	if err != nil {
		return 0, err
	}

	c.logger.Info("Hashcash computed", zap.Duration("elapsed", time.Since(start)))
	fmt.Println(challengeResponse.Challenge, answer, c.cfg.Hashcash.ZerosCount)
	return answer, nil
}

func (c *Client) MustReadMessage() (dto.ServerResponse, error) {
	err := c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return dto.ServerResponse{}, err
	}

	resp := dto.ServerResponse{}

	buff := make([]byte, 1024)
	n, err := c.conn.Read(buff)
	if err != nil {
		return dto.ServerResponse{}, err
	}

	err = c.conn.SetReadDeadline(time.Time{})
	if err != nil {
		return dto.ServerResponse{}, err
	}

	err = json.Unmarshal(buff[:n], &resp)
	if err != nil {
		return dto.ServerResponse{}, err
	}

	return resp, nil
}
func (c *Client) Close() error {
	return c.conn.Close()
}
