package server

import (
	"errors"
	"net"
	"sync"
	"time"

	"go.uber.org/atomic"
	"go.uber.org/zap"

	"github.com/sawdustofmind/pow-word-of-wisdom/internal/config"
	"github.com/sawdustofmind/pow-word-of-wisdom/internal/log"
	"github.com/sawdustofmind/pow-word-of-wisdom/internal/pow"
	"github.com/sawdustofmind/pow-word-of-wisdom/internal/queue"
	"github.com/sawdustofmind/pow-word-of-wisdom/internal/quotes"
)

type Client struct {
	cfg           *config.Config
	quotesService quotes.QuoteService

	conn net.Conn

	messages                 queue.Queue
	closedOnce               sync.Once
	queueSizeLimitExceedOnce sync.Once

	challenge *atomic.String
	logger    *zap.Logger
}

func newClient(
	id string,
	cfg *config.Config,
	conn net.Conn,
	quotesService quotes.QuoteService,
) Client {
	return Client{
		cfg:           cfg,
		conn:          conn,
		quotesService: quotesService,

		messages:  queue.New(),
		challenge: atomic.NewString(""),

		logger: log.GetLogger().With(
			zap.String("component", "client"),
			zap.String("id", id),
			zap.String("address", conn.RemoteAddr().String()),
		),
	}
}

func (c *Client) start() {
	// start read - write loops
	c.refreshChallenge()
	go c.writeLoop()
	go c.readLoop()
}

func (c *Client) readLoop() {
	// if read loop finished, it is necessary to close that all
	defer c.close()

	for {
		buffer := make([]byte, 1024) // TODO: pool, local field to avoid allocations

		nReadBytes, err := c.conn.Read(buffer)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				c.logger.Info("client disconnected")
				return
			}
			c.logger.Info("failed to read message", zap.Error(err))
			return
		}

		rawMessage := buffer[:nReadBytes]

		if c.cfg.Client.DebugRead {
			c.logger.Info("received message", zap.String("message", string(rawMessage)))
		}

		err = c.processMessage(rawMessage)
		if err != nil {
			c.logger.Info("failed to process message", zap.Error(err))
		}
	}
}

func (c *Client) Enqueue(data []byte) {
	if c.cfg.Client.DebugWrite {
		c.logger.Info("write message to client", zap.String("message", string(data)))
	}

	if c.closed() {
		return
	}

	// client's buffer limit exceed
	if c.cfg.Client.QueueLenLimit > 0 && c.messages.Len() >= c.cfg.Client.QueueLenLimit {
		c.queueSizeLimitExceedOnce.Do(func() {
			// stop write, ping-pong routines
			c.messages.Close()

			c.logger.Info("client is too slow, close connection")
			// disconnect client in another routine
			go func() {
				c.close()
			}()
		})
		return
	}

	c.messages.Add(data)
}

func (c *Client) writeLoop() {
	for {
		// Wait for message from queue.
		msg, ok := c.messages.Wait()
		if !ok {
			if c.closed() {
				return
			}
			continue
		}

		if writeErr := c.write(msg); writeErr != nil {
			c.logger.Warn("failed to write to client", zap.Error(writeErr))
			c.close()
			return
		}
	}
}

func (c *Client) write(message []byte) error {
	if c.cfg.Client.WriteTimeout > 0 {
		_ = c.conn.SetWriteDeadline(time.Now().In(time.UTC).Add(c.cfg.Client.WriteTimeout))
	}

	if _, err := c.conn.Write(message); err != nil {
		return err
	}

	if c.cfg.Client.WriteTimeout > 0 {
		_ = c.conn.SetWriteDeadline(time.Time{})
	}

	return nil
}

func (c *Client) close() {
	c.closedOnce.Do(func() {
		c.messages.Close()
		_ = c.conn.Close()
	})
}

func (c *Client) closed() bool {
	return c.messages.Closed()
}

func (c *Client) refreshChallenge() {
	c.challenge.Store(pow.GenerateChallenge(c.cfg.Client.ChallengeLen))
}
