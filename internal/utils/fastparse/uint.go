// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package fastparse

type Uint interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func StringParseUint[U Uint](s string) (U, error) {
	return ParseUint[U]([]byte(s))
}

func ParseUint[U Uint](b []byte) (U, error) {
	if len(b) == 0 {
		return 0, ErrEmptyInput
	}

	var val uint64 = 0

	for i := 0; i < len(b); i++ {
		c := b[i]
		if c < '0' || c > '9' {
			return 0, ErrInvalidUint.value(b)
		}
		val = val*10 + uint64(c-'0')
	}

	return U(val), nil
}
