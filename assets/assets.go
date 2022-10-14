package assets

import (
	_ "embed"
	"math/rand"
	"strings"
	"time"
)

var (
	//go:embed data/quotes.txt
	quotes string

	// QuotesList holds list of all quotes stored within a service.
	QuotesList []string
)

func init() {
	QuotesList = strings.Split(quotes, "\n")
}

// RandomQuote returns a random quote string from the QuotesList.
func RandomQuote() string {
	rand.Seed(time.Now().UnixNano())

	return QuotesList[rand.Int63n(int64(len(QuotesList)))]
}
