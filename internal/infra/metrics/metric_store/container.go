// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package metricstore

import (
	"sync"
	"sync/atomic"
)

type ValueStore[T any] struct {
	ptr  atomic.Pointer[T]
	pool sync.Pool
}

// NewValueStore – creates value container (parallel safe)
func NewValueStore[T any]() *ValueStore[T] {
	return &ValueStore[T]{
		pool: sync.Pool{
			New: func() any {
				var zero T
				return &zero
			},
		},
	}
}

// Save – put object to container
func (s *ValueStore[T]) Save(v T) {
	newVal := s.pool.Get().(*T)
	*newVal = v

	old := s.ptr.Swap(newVal)
	if old != nil {
		*old = *new(T)
		s.pool.Put(old)
	}
}

// Get – return object from container
func (s *ValueStore[T]) Get() T {
	val := s.ptr.Load()
	return *val
}

// Clear – free container
func (s *ValueStore[T]) Clear() {
	old := s.ptr.Swap(nil)
	if old != nil {
		*old = *new(T)
		s.pool.Put(old)
	}
}
