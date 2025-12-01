// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package system

type numerable interface {
	~float32 | ~float64 |
		~uint | ~uint32 | ~uint64 |
		~int | ~int32 | ~int64
}

// uwPtr – unwrap T numerable pointer without panic
func uwPtr[T numerable](ptr *T) T {
	if ptr == nil {
		return T(0)
	}
	return *ptr
}

func usedPercent[T, N numerable](frac, full T) N {
	if full > 0 {
		return N(frac) / N(full) * N(100)
	}
	return 0
}
