// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package logger

import (
	"strings"
	"sync"
)

type StateStorage interface {
	CleanValue(key any)
	PushValue(key, value any)
	GetValue(key any) any
}

type LoggerCollector struct {
	size  int
	stack []string
	sync.RWMutex
}

func NewLoggerCollector() *LoggerCollector {
	return &LoggerCollector{
		size:  0,
		stack: make([]string, 0),
	}
}

func (lc *LoggerCollector) Write(p []byte) (int, error) {

	lc.Lock()
	defer lc.Unlock()

	lc.stack = append(lc.stack, strings.Trim(string(p), "\n"))

	lc.size++

	return len(p), nil
}

func (lc *LoggerCollector) GetStack() ([]string, error) {

	lc.RLock()
	defer lc.RUnlock()

	return lc.stack, nil
}

func (lc *LoggerCollector) Clean() {
	lc.Lock()
	defer lc.Unlock()

	lc.stack = (lc.stack)[:0]
}
