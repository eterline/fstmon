package output

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// SpeedString - formatted output of bytes: 32KB
func SizeString(v uint64) string {
	const (
		_          = iota
		KB float64 = 1 << (10 * iota)
		MB
		GB
		TB
	)

	fv := float64(v)

	switch {
	case fv >= TB:
		return fmt.Sprintf("%.2fTB", fv/TB)
	case fv >= GB:
		return fmt.Sprintf("%.2fGB", fv/GB)
	case fv >= MB:
		return fmt.Sprintf("%.2fMB", fv/MB)
	case fv >= KB:
		return fmt.Sprintf("%.2fKB", fv/KB)
	default:
		return fmt.Sprintf("%dB", v)
	}
}

func UsageSizes(use, total uint64) string {
	str := strings.Builder{}
	str.Grow(3)

	str.WriteString(SizeString(use))
	str.WriteString("/")
	str.WriteString(SizeString(total))

	return str.String()
}

type Float interface {
	~float32 | ~float64
}

func AverageFloat[F Float](l []F) F {
	len := len(l)
	if len == 0 {
		return 0.0
	}

	sum := F(0)
	for _, v := range l {
		sum += v
	}
	return sum / F(len)
}

// SpeedString - formatted output of bytes speed: 32KB/s
func SpeedString(v uint64) string {
	return SizeString(v) + "/s"
}

// FmtTime - formatted output of duration: 1h22m33s
func FmtTime(t time.Duration) string {
	seconds := int(t.Seconds())

	h := seconds / 3600
	m := (seconds % 3600) / 60
	s := seconds % 60

	res := strings.Builder{}

	if h > 0 {
		res.WriteString(fmt.Sprintf("%dh", h))
	}
	if m > 0 {
		res.WriteString(fmt.Sprintf("%dm", m))
	}
	if s > 0 || res.String() == "" {
		res.WriteString(fmt.Sprintf("%ds", s))
	}

	return res.String()
}

func CelsiusString(temp float64) string {
	if temp < (-273.15) {
		return fmt.Sprint("-273.15°C")
	}

	return strconv.FormatFloat(temp, 'f', 1, 64) + "°C"
}

func PercentString(v float64) string {
	return strconv.Itoa(int(v)) + "%"
}
