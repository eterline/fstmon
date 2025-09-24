// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package calc

func AverageFloat64(f []float64) float64 {
	ln := len(f)
	if ln == 0 {
		return 0.0
	}

	sum := float64(0)
	for _, v := range f {
		sum += v
	}

	return sum / float64(ln)
}

func AverageFloat32(f []float32) float32 {
	ln := len(f)
	if ln == 0 {
		return 0.0
	}

	sum := float32(0)
	for _, v := range f {
		sum += v
	}

	return sum / float32(ln)
}
