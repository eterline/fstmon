package metricstore

import (
	"sync"
	"sync/atomic"
)

type ValueStore[T any] struct {
	ptr   atomic.Pointer[T]
	pool  sync.Pool
	cycle *uint64
}

// NewValueStore - creates value container (parallel safe)
func NewValueStore[T any]() *ValueStore[T] {
	return &ValueStore[T]{
		pool: sync.Pool{
			New: func() any {
				var zero T
				return &zero
			},
		},
		cycle: new(uint64),
	}
}

// Save - put object to container
func (s *ValueStore[T]) Save(v T) {
	newVal := s.pool.Get().(*T)
	*newVal = v

	old := s.ptr.Swap(newVal)
	if old != nil {
		*old = *new(T)
		s.pool.Put(old)
		atomic.AddUint64(s.cycle, 1)
	}
}

// Get - return object from container
func (s *ValueStore[T]) Get() (T, bool) {
	val := s.ptr.Load()
	return *val, (val != nil)
}

// Clear - free container
func (s *ValueStore[T]) Clear() {
	old := s.ptr.Swap(nil)
	if old != nil {
		*old = *new(T)
		s.pool.Put(old)
	}
}

// Updated - update counter
func (s *ValueStore[T]) Updated() uint64 {
	return atomic.LoadUint64(s.cycle)
}
