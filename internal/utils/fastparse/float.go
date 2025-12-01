package fastparse

import "errors"

func StringParseFloat[F float32 | float64](s string) (F, error) {
	return ParseFloat[F]([]byte(s))
}

func ParseFloat[F float32 | float64](b []byte) (F, error) {
	if len(b) == 0 {
		return 0, errors.New("empty input")
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

		return 0, errors.New("invalid float format")
	}

	for ; i < len(b); i++ {
		c := b[i]
		if c < '0' || c > '9' {
			return 0, errors.New("invalid float format")
		}
		fracPart = fracPart*10 + float64(c-'0')
		fracDiv *= 10
	}

	return F(sign * (intPart + fracPart/fracDiv)), nil
}
