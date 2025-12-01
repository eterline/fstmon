// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package fastparse

type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

func StringParseInt[I Int](s string) (I, error) {
	return ParseInt[I]([]byte(s))
}

func ParseInt[I Int](b []byte) (I, error) {
	if len(b) == 0 {
		return 0, ErrEmptyInput
	}

	sign := int64(1)
	val := int64(0)
	i := 0

	if b[i] == '-' {
		sign = -1
		i++
	} else if b[i] == '+' {
		i++
	}

	for ; i < len(b); i++ {
		c := b[i]
		if c < '0' || c > '9' {
			return 0, ErrInvalidInt.value(b)
		}
		val = val*10 + int64(c-'0')
	}

	return I(sign * val), nil
}
