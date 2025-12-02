// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package config

import (
	"io"
	"os"
)

// stdout proxy without close implementation
type stdoutWrap struct {
	stdout *os.File
}

func newStdoutWrap(stdout *os.File) io.WriteCloser {
	return &stdoutWrap{
		stdout: stdout,
	}
}

func (sw *stdoutWrap) Write(p []byte) (n int, err error) {
	return sw.stdout.Write(p)
}

func (sw *stdoutWrap) Close() error {
	return nil
}
