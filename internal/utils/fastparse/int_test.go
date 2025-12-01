// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
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
