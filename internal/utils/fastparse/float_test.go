package fastparse_test

import (
	"testing"

	"github.com/eterline/fstmon/internal/utils/fastparse"
)

func BenchmarkParseFloat(b *testing.B) {

	str := []byte("11.32")

	for n := 0; n < b.N; n++ {
		_, _ = fastparse.ParseFloat[float64](str)
	}
}
