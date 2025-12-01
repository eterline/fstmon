package fastparse

import "errors"

func StringParseUint[U uint | uint8 | uint16 | uint32 | uint64](s string) (U, error) {
	return ParseUint[U]([]byte(s))
}

func ParseUint[U uint | uint8 | uint16 | uint32 | uint64](b []byte) (U, error) {
	if len(b) == 0 {
		return 0, errors.New("empty input")
	}

	var val uint64 = 0

	for i := 0; i < len(b); i++ {
		c := b[i]
		if c < '0' || c > '9' {
			return 0, errors.New("invalid uint format")
		}
		val = val*10 + uint64(c-'0')
	}

	return U(val), nil
}
