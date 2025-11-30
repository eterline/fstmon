package domain

import "time"

// ==========================

/*
Numerable - type constraint matching all numeric types in Go.

	Used for generic IO structures to support arithmetic operations.
*/
type Numerable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

/*
NewIO - creates a new IO instance with given RX and TX values.

	Initializes Summary as RX + TX.
*/
func NewIO[T Numerable](rx, tx T) IO[T] {
	return IO[T]{
		RX:      rx,
		TX:      tx,
		Summary: rx + tx,
	}
}

/*
IO - generic structure representing input/output metrics.

	Holds RX (received), TX (transmitted) and a Summary field (RX+TX).
*/
type IO[T Numerable | time.Duration] struct {
	Summary T `json:"summary"` // Total of RX and TX
	RX      T `json:"rx"`      // Received value
	TX      T `json:"tx"`      // Transmitted value
}

// ====

// IncRX - increments the RX field by given value and updates Summary.
func (io *IO[T]) IncRX(v T) {
	io.RX += v
	io.Summary += v
}

// IncTX - increments the TX field by given value and updates Summary.
func (io *IO[T]) IncTX(v T) {
	io.TX += v
	io.Summary += v
}

// ====

// DecRX - decrements the RX field by given value and updates Summary.
func (io *IO[T]) DecRX(v T) {
	io.RX -= v
	io.Summary -= v
}

// DecTX - decrements the TX field by given value and updates Summary.
func (io *IO[T]) DecTX(v T) {
	io.TX -= v
	io.Summary -= v
}
