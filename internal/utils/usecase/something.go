// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package usecase

import "time"

func MapSlicesLen[K comparable, V any](m map[K]V) (keys []K, values []V, n int) {
	mapLen := len(m)

	if mapLen == 0 {
		return nil, nil, 0
	}

	keys = make([]K, 0, mapLen)
	values = make([]V, 0, mapLen)

	for key, value := range m {
		keys = append(keys, key)
		values = append(values, value)
	}

	return keys, values, mapLen
}

func MapSlices[K comparable, V any](m map[K]V) (keys []K, values []V) {
	keys, values, _ = MapSlicesLen(m)
	return keys, values
}

/*
Numerable – type constraint matching all numeric types in Go.

	Used for generic IO structures to support arithmetic operations.
*/
type Numerable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

func MsToDuration[T Numerable](ms T) time.Duration {
	return time.Millisecond * time.Duration(ms)
}
