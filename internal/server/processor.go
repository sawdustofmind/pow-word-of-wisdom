package server

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sawdustofmind/pow-word-of-wisdom/internal/dto"
	"github.com/sawdustofmind/pow-word-of-wisdom/internal/pow"
)

// TODO: move to service
func (c *Client) processMessage(rawMessage []byte) error {
	req := dto.ClientMessage{}
	err := json.Unmarshal(rawMessage, &req)
	if err != nil {
		return err
	}
	switch req.Type {
	case dto.ChallengeClientMessageType:
		return c.processRequestChallenge()
	case dto.GetQuoteClientMessageType:
		return c.processRequestQuote(req)
	default:
		return fmt.Errorf("unknown message type: %s", req.Type)
	}
}

func (c *Client) processRequestQuote(req dto.ClientMessage) error {
	err := c.checkChallenge(req)
	if err != nil {
		return err
	}

	quote, err := c.quotesService.GetRandomQuote()
	if err != nil {
		return err
	}

	quoteResp, err := json.Marshal(dto.QuoteResponse{Quote: quote})
	if err != nil {
		return err
	}

	err = c.writeResponse(dto.QuoteResponseType, quoteResp)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) processRequestChallenge() error {
	c.refreshChallenge()

	challengeResp, err := json.Marshal(dto.ChallengeResponse{
		Challenge: c.challenge.Load(),
	})
	if err != nil {
		return err
	}

	err = c.writeResponse(dto.ChallengeResponseType, challengeResp)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) checkChallenge(req dto.ClientMessage) error {
	fmt.Println(c.challenge.Load(), req.ChallengeCounter, c.cfg.Hashcash.ZerosCount)
	if pow.IsChallengeSolved(c.challenge.Load(), req.ChallengeCounter, c.cfg.Hashcash.ZerosCount) {
		return nil
	}
	return fmt.Errorf("invalid challenge answer")
}

func (c *Client) writeResponse(dataType string, data json.RawMessage) error {
	// todo: sync pool for replies

	response, err := json.Marshal(dto.ServerResponse{
		Type: dataType,
		Data: data,
		Ts:   time.Now().UnixNano(),
	})
	if err != nil {
		return err
	}

	response = append(response, byte('\n'))
	c.Enqueue(response)
	return nil
}
