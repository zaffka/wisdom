package assets_test

import (
	"testing"

	"github.com/zaffka/wisdom/assets"
)

func TestQuotesListLen(t *testing.T) {
	awaitedQuotesLen := 18
	gotQuotesLen := len(assets.QuoteList)
	if gotQuotesLen != awaitedQuotesLen {
		t.Logf("wrong quotes len, %d awaited, %d found", awaitedQuotesLen, gotQuotesLen)
		t.Fail()
	}
}
