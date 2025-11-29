package sizes

import "math/bits"

type (
	sizeByteBase uint64
	sizeBitBase  uint64
)

// Byte-based sizes
const (
	Byte sizeByteBase = 1 << (10 * iota) // 1, 1024, 1024^2, ...
	KB                                   // Kilobyte
	MB                                   // Megabyte
	GB                                   // Gigabyte
	TB                                   // Terabyte
	PB                                   // Petabyte
	EB                                   // Exabyte
)

func (d sizeByteBase) String() string {
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

func DetermByte2ByteBase(bytesSize uint64) (float64, sizeByteBase) {
	if bytesSize == 0 {
		return 0, Byte
	}

	msb := bits.Len64(bytesSize) - 1
	unitExp := msb / 10
	if unitExp > 6 {
		unitExp = 6
	}

	unit := sizeByteBase(1 << (10 * unitExp))
	value := float64(bytesSize) / float64(unit)

	return value, unit
}

// Bit-based sizes
const (
	Bit sizeBitBase = 1 << (10 * iota) * 8 // 1 bit, 1024 bits / 8 = 128 bytes, ...
	Kb                                     // Kilobit
	Mb                                     // Megabit
	Gb                                     // Gigabit
	Tb                                     // Terabit
	Pb                                     // Petabit
	Eb                                     // Exabit
)

func (d sizeBitBase) String() string {
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

func DetermByte2BitBase(bitsSize uint64) (float64, sizeBitBase) {
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

	return value, sizeBitBase(unit * 8)
}
