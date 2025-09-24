package process

import (
	"sync/atomic"
)

type ValueHolder[T interface{}] struct {
	value atomic.Pointer[T]
}

func (h *ValueHolder[T]) Set(v T) {
	h.value.Store(&v)
}

func (h *ValueHolder[T]) Clear() {
	h.value.Store(nil)
}

func (h *ValueHolder[T]) Value() (T, bool) {
	ptr := h.value.Load()
	if ptr == nil {
		var zero T
		return zero, false
	}
	return *ptr, true
}

// ErrorHolder - structure with multithread error provide in some services
type ErrorHolder struct {
	internal ValueHolder[error]
}

// SetError - puts error into holder
func (e *ErrorHolder) SetError(err error) {
	e.internal.Set(err)
}

// SetError - clears error from holder
func (e *ErrorHolder) ClearError() {
	e.internal.Clear()
}

// SetError - return error value from holder
func (e *ErrorHolder) Err() (error, bool) {
	return e.internal.Value()
}
