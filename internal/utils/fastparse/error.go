// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package fastparse

import "fmt"

type FastparseError string

func (er FastparseError) value(b []byte) FastparseError {
	return FastparseError(fmt.Sprintf("%s: '%v'", er, string(b)))
}

func (er FastparseError) Error() string {
	return "fastparse error: " + string(er)
}

const (
	ErrEmptyInput   FastparseError = "empty input"
	ErrInvalidFloat FastparseError = "invalid float format"
	ErrInvalidInt   FastparseError = "invalid int format"
	ErrInvalidUint  FastparseError = "invalid uint format"
)
