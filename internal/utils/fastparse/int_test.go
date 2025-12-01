package fastparse_test

import (
	"testing"

	"github.com/eterline/fstmon/internal/utils/fastparse"
)

func BenchmarkParseInt(b *testing.B) {

	// str := []byte("-2553")
	strs := "-2553"

	for n := 0; n < b.N; n++ {
		_, _ = fastparse.StringParseInt[int64](strs)
	}
}
