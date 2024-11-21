package dto

const (
	ChallengeResponseType = "challenge"
)

type ChallengeResponse struct {
	Challenge string `json:"challenge"`
}
