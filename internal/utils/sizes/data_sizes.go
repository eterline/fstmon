// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package sizes

import "math/bits"

type (
	SizeByteBase uint64
	SizeBitBase  uint64
)

// Byte-based sizes
const (
	Byte SizeByteBase = 1 << (10 * iota) // 1, 1024, 1024^2, ...
	KB                                   // Kilobyte
	MB                                   // Megabyte
	GB                                   // Gigabyte
	TB                                   // Terabyte
	PB                                   // Petabyte
	EB                                   // Exabyte
)

func (d SizeByteBase) In(v uint64) uint64 {
	return v * uint64(d)
}

func (d SizeByteBase) String() string {
	switch d {
	case Byte:
		return "B"
	case KB:
		return "KB"
	case MB:
		return "MB"
	case GB:
		return "GB"
	case TB:
		return "TB"
	case PB:
		return "PB"
	case EB:
		return "EB"
	default:
		return "UNKNOWN"
	}
}

func DetermByte2ByteBase(bytesSize uint64) (float64, SizeByteBase) {
	if bytesSize == 0 {
		return 0, Byte
	}

	msb := bits.Len64(bytesSize) - 1
	unitExp := msb / 10
	if unitExp > 6 {
		unitExp = 6
	}

	unit := SizeByteBase(1 << (10 * unitExp))
	value := float64(bytesSize) / float64(unit)

	return value, unit
}

// Bit-based sizes
const (
	Bit SizeBitBase = 1 << (10 * iota) // 1, 1024, 1024^2, ...
	Kb
	Mb
	Gb
	Tb
	Pb
	Eb
)

func (d SizeBitBase) String() string {
	switch d {
	case Bit:
		return "b"
	case Kb:
		return "Kb"
	case Mb:
		return "Mb"
	case Gb:
		return "Gb"
	case Tb:
		return "Tb"
	case Pb:
		return "Pb"
	case Eb:
		return "Eb"
	default:
		return "UNKNOWN"
	}
}

func DetermByte2BitBase(bytesSize uint64) (float64, SizeBitBase) {
	return DetermBitBase(bytesSize * 8)
}

func DetermBitBase(bitsSize uint64) (float64, SizeBitBase) {
	if bitsSize == 0 {
		return 0, Bit
	}

	msb := bits.Len64(bitsSize) - 1

	unitExp := msb / 10
	if unitExp > 6 {
		unitExp = 6
	}

	unit := uint64(1 << (10 * unitExp))

	value := float64(bitsSize) / float64(unit)

	return value, SizeBitBase(unit)
}

// SizeISOMetric — 1000-based SI units (Extended)
type SizeISOMetric uint64

// 1000^n units starting from K (10^3) up to Y (10^24)
const (
	ZERO SizeISOMetric = 0

	K SizeISOMetric = 1_000
	M SizeISOMetric = 1_000_000             // 1000^2
	G SizeISOMetric = 1_000_000_000         // 1000^3
	T SizeISOMetric = 1_000_000_000_000     // 1000^4
	P SizeISOMetric = 1_000_000_000_000_000 // 1000^5
	E SizeISOMetric = 1_000_000_000_000_000_000
)

func (u SizeISOMetric) String() string {
	switch u {
	case ZERO:
		return ""
	case K:
		return "K"
	case M:
		return "M"
	case G:
		return "G"
	case T:
		return "T"
	case P:
		return "P"
	default:
		return "E"
	}
}

// DetermMetricBase — chooses the best SI unit (K..Y) for the given value.
// Returns normalized value + the matched SI unit.
func DetermMetricBase(size uint64) (float64, SizeISOMetric) {
	if size < 1_000 {
		return float64(size), ZERO
	}

	// exact thresholds for units
	switch {
	case size < 1_000_000:
		return float64(size) / float64(K), K
	case size < 1_000_000_000:
		return float64(size) / float64(M), M
	case size < 1_000_000_000_000:
		return float64(size) / float64(G), G
	case size < 1_000_000_000_000_000:
		return float64(size) / float64(T), T
	case size < 1_000_000_000_000_000_000:
		return float64(size) / float64(P), P
	default:
		// 1000^6 still fits uint64
		return float64(size) / float64(E), E
	}
}
