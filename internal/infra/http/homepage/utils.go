package httphomepage

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/eterline/fstmon/internal/utils/sizes"
)

func formatDuration(d time.Duration, delim string, upperSet bool) string {

	sB := &strings.Builder{}
	sB.Grow(32)

	delimWrite := func() {
		if delim != "" {
			sB.WriteString(delim)
		}
	}

	if v := d.Hours(); v > 0 {
		writeUint64(sB, uint64(v))
		sB.Write(updownSet([]byte("h"), upperSet))
	}

	delimWrite()

	if v := d.Minutes(); v > 0 {
		writeUint64(sB, uint64(v)%60)
		sB.Write(updownSet([]byte("m"), upperSet))
	}

	delimWrite()

	if v := d.Seconds(); v > 0 {
		writeUint64(sB, uint64(v)%60)
		sB.Write(updownSet([]byte("s"), upperSet))
	}

	return sB.String()
}

func writeUint64(w io.Writer, v uint64) error {
	var buf [32]byte
	b := buf[:0]

	b = strconv.AppendUint(b, v, 10)

	_, err := w.Write(b)
	return err
}

func updownSet(p []byte, upperSet bool) []byte {
	if upperSet {
		return bytes.ToUpper(p)
	}
	return bytes.ToLower(p)
}

func queryByteBaseSizes(delim byte, ss ...uint64) string {
	if len(ss) == 0 {
		return "..."
	}

	sB := &strings.Builder{}
	sB.Grow(32)

	for i, s := range ss {
		v, unit := sizes.DetermByte2ByteBase(s)
		fmt.Fprintf(sB, "%.2f", v)
		sB.WriteString(unit.String())

		if i >= 0 && i < (len(ss)-1) {
			sB.WriteByte(delim)
		}
	}

	return sB.String()
}
