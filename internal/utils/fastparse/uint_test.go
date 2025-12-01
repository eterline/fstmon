// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package fastparse_test

import (
	"testing"

	"github.com/eterline/fstmon/internal/utils/fastparse"
)

func BenchmarkParseUint(b *testing.B) {

	str := []byte("34567")

	for n := 0; n < b.N; n++ {
		_, _ = fastparse.ParseUint[uint64](str)
	}
}
