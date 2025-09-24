// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

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
