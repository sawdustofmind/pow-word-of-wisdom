package quotes

type QuoteService interface {
	GetRandomQuote() (string, error)
}
