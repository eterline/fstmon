// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package fastparse

type Float interface {
	~float32 | ~float64
}

func StringParseFloat[F Float](s string) (F, error) {
	return ParseFloat[F]([]byte(s))
}

func ParseFloat[F Float](b []byte) (F, error) {
	if len(b) == 0 {
		return 0, ErrEmptyInput
	}

	if len(b) == 1 && b[0] == '0' {
		return 0, nil
	}

	var sign float64 = 1
	var intPart float64 = 0
	var fracPart float64 = 0
	var fracDiv float64 = 1
	i := 0

	if b[i] == '-' {
		sign = -1
		i++
	} else if b[i] == '+' {
		i++
	}

	for ; i < len(b); i++ {
		c := b[i]

		if c >= '0' && c <= '9' {
			intPart = intPart*10 + float64(c-'0')
			continue
		}

		if c == '.' {
			i++
			break
		}

		return 0, ErrInvalidFloat.value(b)
	}

	for ; i < len(b); i++ {
		c := b[i]
		if c < '0' || c > '9' {
			return 0, ErrInvalidFloat.value(b)
		}
		fracPart = fracPart*10 + float64(c-'0')
		fracDiv *= 10
	}

	if r := F(sign * (intPart + fracPart/fracDiv)); r != 0 {
		return r, nil
	}

	return 0, ErrInvalidFloat.value(b)
}
