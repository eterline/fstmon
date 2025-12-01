package fastparse

import "errors"

func StringParseInt[I int | int8 | int16 | int32 | int64](s string) (I, error) {
	return ParseInt[I]([]byte(s))
}

func ParseInt[I int | int8 | int16 | int32 | int64](b []byte) (I, error) {
	if len(b) == 0 {
		return 0, errors.New("empty input")
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
			return 0, errors.New("invalid int format")
		}
		val = val*10 + int64(c-'0')
	}

	return I(sign * val), nil
}
