package quotes

import (
	"math/rand/v2"
	"regexp"
)

type InMemoryQuoteStore struct {
}

var _ QuoteService = (*InMemoryQuoteStore)(nil)

func NewInMemoryQuoteStore() *InMemoryQuoteStore {
	return &InMemoryQuoteStore{}
}

func (i InMemoryQuoteStore) GetRandomQuote() (string, error) {
	return imMemoryQuotes[rand.IntN(len(imMemoryQuotes))], nil
}

var imMemoryQuotes = []string{
	`Heavy rains were approaching.`,
	`A man had a huge stock of dry fruits because of a 
		wedding in his family.`,
	`One of his neighbors approached him and offered thrice the cost of dry fruits.
		He happily traded them.`,
	`In a week, the entire village was flooded, and people found shelter
		on the roofs of their houses. Soon most people were out of food and the 
		leftover became stale except for edibles like dry fruits.`,
	`The same man now repurchased some of the dry fruits from the neighbor at five times
		the cost to survive.`,
	`Most people take decisions keeping in mind immediate or short-run results and returns.`,
	`They become blinded by the shine of bumper returns and the resulting spike in lifestyle.`,
	`A decision has value when it foresees its impact on the future and keeps in mind the
		learnings from the past and touches major subjects or matters of life.`,
	`Life must be seen with a wholesome vision. Keeping in mind materialistic and 
		spiritual returns, short-term and long-term benefits and losses, learning 
		and experience, values and commitments.`,
	`Even optimism has to be balanced with respectful thought given to the pessimistic
		approach.I have to live the entire story well and support my present and
		future self with a balance.Wisdom is simply not being blinded by one thing and 
		thinking broadly and following a balanced approach.`,
	`The wood from a tree would definitely give you handsome returns and solve your present 
		issues but selling its fruits would give you much higher, stable, and consistent 
		returns in the long run.`,
}

func init() {
	spaceRe := regexp.MustCompile(`\s+`)
	for i := range imMemoryQuotes {
		imMemoryQuotes[i] = spaceRe.ReplaceAllString(imMemoryQuotes[i], " ")
	}
}
