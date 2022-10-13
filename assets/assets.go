package assets

import (
	_ "embed"

	jsoniter "github.com/json-iterator/go"
)

//go:embed data/quotes.json
var quotes []byte

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary

	// QuoteList holds list of all quotes stored within a service.
	QuoteList []Quote
)

// Quote represents a single citation of Wisdom.
type Quote struct {
	Text   string `json:"text"`
	Author string `json:"author"`
}

func init() {
	if err := json.Unmarshal(quotes, &QuoteList); err != nil {
		panic(err)
	}
}
