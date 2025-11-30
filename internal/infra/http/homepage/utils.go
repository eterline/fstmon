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

type queryByteBaseSizesBuilder struct {
	delim byte
	q     int
	sB    *strings.Builder
}

func NewQBBSBuilder(delim byte) *queryByteBaseSizesBuilder {
	b := &strings.Builder{}
	b.Grow(32)
	return &queryByteBaseSizesBuilder{
		delim: delim,
		q:     0,
		sB:    b,
	}
}

func (q *queryByteBaseSizesBuilder) Add(v uint64) *queryByteBaseSizesBuilder {
	if q.q > 0 {
		q.sB.WriteByte(q.delim)
	}
	q.q++

	vf, unit := sizes.DetermByte2ByteBase(v)
	fmt.Fprintf(q.sB, "%.2f", vf)
	q.sB.WriteString(unit.String())

	return q
}

func (q *queryByteBaseSizesBuilder) Build() string {
	if q.q == 0 {
		return "..."
	}
	return q.sB.String()
}
