// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package sizes_test

import (
	"math"
	"testing"

	"github.com/eterline/fstmon/internal/utils/sizes"
)

func TestSizeByteBaseString(t *testing.T) {
	tests := []struct {
		input    sizes.SizeByteBase
		expected string
	}{
		{sizes.Byte, "B"},
		{sizes.KB, "KB"},
		{sizes.MB, "MB"},
		{sizes.GB, "GB"},
		{sizes.TB, "TB"},
		{sizes.PB, "PB"},
		{sizes.EB, "EB"},
		{sizes.SizeByteBase(999999999999999), "UNKNOWN"},
	}

	for _, tt := range tests {
		if got := tt.input.String(); got != tt.expected {
			t.Errorf("SizeByteBase.String() = %q; want %q", got, tt.expected)
		}
	}
}

func TestDetermByte2ByteBase(t *testing.T) {
	tests := []struct {
		input    uint64
		expValue float64
		expUnit  sizes.SizeByteBase
	}{
		{0, 0, sizes.Byte},
		{500, 500, sizes.Byte},
		{1024, 1, sizes.KB},
		{5 * 1024 * 1024, 5, sizes.MB},
		{7 * 1024 * 1024 * 1024, 7, sizes.GB},
	}

	for _, tt := range tests {
		val, unit := sizes.DetermByte2ByteBase(tt.input)
		if val != tt.expValue || unit != tt.expUnit {
			t.Errorf("DetermByte2ByteBase(%d) = (%f, %v); want (%f, %v)", tt.input, val, unit, tt.expValue, tt.expUnit)
		}
	}
}

func TestSizeBitBaseString(t *testing.T) {
	tests := []struct {
		input    sizes.SizeBitBase
		expected string
	}{
		{sizes.Bit, "b"},
		{sizes.Kb, "Kb"},
		{sizes.Mb, "Mb"},
		{sizes.Gb, "Gb"},
		{sizes.Tb, "Tb"},
		{sizes.Pb, "Pb"},
		{sizes.Eb, "Eb"},
		{sizes.SizeBitBase(999999999999999), "UNKNOWN"},
	}

	for _, tt := range tests {
		if got := tt.input.String(); got != tt.expected {
			t.Errorf("SizeBitBase.String() = %q; want %q", got, tt.expected)
		}
	}
}

func TestDetermByte2BitBase(t *testing.T) {
	tests := []struct {
		name     string
		input    uint64
		expValue float64
		expUnit  sizes.SizeBitBase
	}{
		{
			name:     "zero",
			input:    0,
			expValue: 0,
			expUnit:  sizes.Bit,
		},
		{
			name:     "1 byte -> 8 bits",
			input:    1,
			expValue: 8,
			expUnit:  sizes.Bit,
		},
		{
			name:     "128 bytes -> 1024 bits → 1 Kb",
			input:    128,
			expValue: 1,
			expUnit:  sizes.Kb,
		},
		{
			name:     "1 MiB -> 8,388,608 bits → ~8 Mb",
			input:    1024 * 1024,
			expValue: 8,
			expUnit:  sizes.Mb,
		},
		{
			name:     "1 GiB -> 8,589,934,592 bits → ~8 Gb",
			input:    1024 * 1024 * 1024,
			expValue: 8,
			expUnit:  sizes.Gb,
		},
	}

	const eps = 1e-9

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, unit := sizes.DetermByte2BitBase(tt.input)

			if math.Abs(val-tt.expValue) > eps {
				t.Errorf("value mismatch: got %.12f want %.12f", val, tt.expValue)
			}

			if unit != tt.expUnit {
				t.Errorf("unit mismatch: got %v want %v", unit, tt.expUnit)
			}
		})
	}
}

func TestSizeISOMetricString(t *testing.T) {
	tests := []struct {
		input    sizes.SizeISOMetric
		expected string
	}{
		{sizes.ZERO, ""},
		{sizes.K, "K"},
		{sizes.M, "M"},
		{sizes.G, "G"},
		{sizes.T, "T"},
		{sizes.P, "P"},
		{sizes.E, "E"},
	}

	for _, tt := range tests {
		if got := tt.input.String(); got != tt.expected {
			t.Errorf("SizeISOMetric.String() = %q; want %q", got, tt.expected)
		}
	}
}

func TestDetermMetricBase(t *testing.T) {
	tests := []struct {
		input    uint64
		expValue float64
		expUnit  sizes.SizeISOMetric
	}{
		{500, 500, sizes.ZERO},
		{1000, 1, sizes.K},
		{1_500_000, 1.5, sizes.M},
		{3_200_000_000, 3.2, sizes.G},
		{7_000_000_000_000, 7, sizes.T},
	}

	for _, tt := range tests {
		val, unit := sizes.DetermMetricBase(tt.input)
		if val != tt.expValue || unit != tt.expUnit {
			t.Errorf("DetermMetricBase(%d) = (%f, %v); want (%f, %v)", tt.input, val, unit, tt.expValue, tt.expUnit)
		}
	}
}
