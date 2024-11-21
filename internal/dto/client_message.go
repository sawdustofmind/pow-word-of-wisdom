package dto

const (
	ChallengeClientMessageType = "challenge"
	GetQuoteClientMessageType  = "quote"
)

type ClientMessage struct {
	Type             string `json:"type"`
	ChallengeCounter int    `json:"challenge_counter"`
}
