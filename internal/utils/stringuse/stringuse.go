package stringuse

import (
	"regexp"
	"strings"
)

var (
	splitRegex = regexp.MustCompile(`[,\s;]+`)
)

func SplitBySpaces(s string) ([]string, bool) {
	slice := splitRegex.Split(strings.TrimSpace(s), -1)
	if len(slice) < 1 || slice == nil {
		return nil, false
	}

	return slice, true
}
