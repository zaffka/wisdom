package assets_test

import (
	_ "embed"
	"reflect"
	"testing"

	"github.com/zaffka/wisdom/assets"
)

func TestQuotesListLen(t *testing.T) {
	awaitedQuotesLen := 18
	gotQuotesLen := len(assets.QuotesList)
	if gotQuotesLen != awaitedQuotesLen {
		t.Logf("wrong quotes len, %d awaited, %d found", awaitedQuotesLen, gotQuotesLen)
		t.Fail()
	}
}

func TestRandomQuote(t *testing.T) {
	set1 := make([]string, 3)
	for i := 0; i < 3; i++ {
		set1[i] = assets.RandomQuote()
	}

	set2 := make([]string, 3)
	for i := 0; i < 3; i++ {
		set2[i] = assets.RandomQuote()
	}

	if reflect.DeepEqual(set1, set2) {
		t.Fatal("random quotes set match")
	}
}
