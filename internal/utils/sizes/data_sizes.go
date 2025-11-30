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
	Bit SizeBitBase = 1 << (10 * iota) * 8 // 1 bit, 1024 bits / 8 = 128 bytes, ...
	Kb                                     // Kilobit
	Mb                                     // Megabit
	Gb                                     // Gigabit
	Tb                                     // Terabit
	Pb                                     // Petabit
	Eb                                     // Exabit
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

func DetermByte2BitBase(bitsSize uint64) (float64, SizeBitBase) {
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

	return value, SizeBitBase(unit * 8)
}

// SizeISOMetric — 1000-based SI units (Extended)
type SizeISOMetric uint64

// 1000^n units starting from K (10^3) up to Y (10^24)
const (
	ZERO SizeISOMetric = iota
	K                  // kilo  (10^3)
	M                  // mega  (10^6)
	G                  // giga  (10^9)
	T                  // tera  (10^12)
	P                  // peta  (10^15)
	E                  // exa   (10^18)
	Z                  // zetta (10^21)
	Y                  // yotta (10^24)
)

func (u SizeISOMetric) String() string {
	switch u {
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
	case E:
		return "E"
	case Z:
		return "Z"
	case Y:
		return "Y"
	default:
		return ""
	}
}

// DetermMetricBase — chooses the best SI unit (K..Y) for the given value.
// Returns normalized value + the matched SI unit.
func DetermMetricBase(size uint64) (float64, SizeISOMetric) {
	if size < 1000 {
		return float64(size), ZERO
	}

	msb := bits.Len64(size) - 1
	unitExp := msb / 10
	if unitExp > 8 {
		unitExp = 8
	}

	// calculate 1000^unitExp
	var unit uint64 = 1
	for i := 0; i < unitExp; i++ {
		unit *= 1000
	}

	value := float64(size) / float64(unit)

	var metric SizeISOMetric
	switch unitExp {
	case 1:
		metric = K
	case 2:
		metric = M
	case 3:
		metric = G
	case 4:
		metric = T
	case 5:
		metric = P
	case 6:
		metric = E
	case 7:
		metric = Z
	case 8:
		metric = Y
	default:
		metric = K
	}

	return value, metric
}
