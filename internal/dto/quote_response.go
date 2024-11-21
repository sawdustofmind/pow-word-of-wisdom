package dto

const (
	QuoteResponseType = "quote"
)

type QuoteResponse struct {
	Quote string `json:"quote"`
}
