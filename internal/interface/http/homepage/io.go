// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package httphomepage

import (
	"fmt"
	"strings"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/internal/utils/sizes"
)

type IO[T domain.Numerable] struct {
	RawRX T      `json:"raw_rx"`
	RawTX T      `json:"raw_tx"`
	RX    string `json:"rx"`
	TX    string `json:"tx"`
}

type Stringer interface {
	String() string
}

type IOBuilder[T domain.Numerable] struct {
	rx       T
	tx       T
	unitFunc func(v T) (float64, Stringer, bool) // bool: "/s" mode
	custom   func(v T) (float64, string)         // for fully custom formatting
	postfix  string
}

// NewIOBuilder — start builder
func NewIOBuilder[T domain.Numerable](rx, tx T) *IOBuilder[T] {
	return &IOBuilder[T]{rx: rx, tx: tx}
}

// AutoUnits — auto-detect byte units
func (b *IOBuilder[T]) AutoUnits() *IOBuilder[T] {
	b.unitFunc = func(v T) (float64, Stringer, bool) {
		fv, unit := sizes.DetermByte2ByteBase(uint64(v))
		return fv, unit, false
	}
	return b
}

// AutoUnitsPerSec — auto-detect byte units + add `/s`
func (b *IOBuilder[T]) AutoUnitsPerSec() *IOBuilder[T] {
	b.unitFunc = func(v T) (float64, Stringer, bool) {
		fv, unit := sizes.DetermByte2ByteBase(uint64(v))
		return fv, unit, true
	}
	return b
}

// CustomUnit — always use specific unit (e.g., MB, KiB)
func (b *IOBuilder[T]) CustomUnit(unit sizes.SizeByteBase) *IOBuilder[T] {
	b.custom = func(v T) (float64, string) {
		fv := float64(v) / float64(unit)
		return fv, unit.String()
	}
	return b
}

// CustomUnitPerSec — same but append `/s`
func (b *IOBuilder[T]) CustomUnitPerSec(unit sizes.SizeByteBase) *IOBuilder[T] {
	b.custom = func(v T) (float64, string) {
		fv := float64(v) / float64(unit)
		return fv, unit.String() + "/s"
	}
	return b
}

// AutoMetricUnits — auto-detect K / M / G / T / P / E / Z / Y
func (b *IOBuilder[T]) AutoMetricUnits() *IOBuilder[T] {
	b.unitFunc = func(v T) (float64, Stringer, bool) {
		fv, u := sizes.DetermMetricBase(uint64(v))
		return fv, u, false
	}
	return b
}

func (b *IOBuilder[T]) WithPostfix(s string) *IOBuilder[T] {
	b.postfix = s
	return b
}

// Format — user supplies ready string values
func (b *IOBuilder[T]) Format(s string) IO[T] {
	return IO[T]{
		RawRX: b.rx,
		RawTX: b.tx,
		RX:    fmt.Sprintf(s, b.rx),
		TX:    fmt.Sprintf(s, b.tx),
	}
}

func (b *IOBuilder[T]) Build() IO[T] {
	io := IO[T]{
		RawRX: b.rx,
		RawTX: b.tx,
	}

	addPostfix := func(rx, tx *string) {
		if b.postfix != "" {
			*rx += b.postfix
			*tx += b.postfix
		}
	}

	if b.custom != nil {
		// RX
		{
			fv, unit := b.custom(b.rx)
			var sb strings.Builder
			sb.Grow(32)

			// "%.2f"
			fmt.Fprintf(&sb, "%.2f", fv)
			sb.WriteString(unit)

			io.RX = sb.String()
		}

		// TX
		{
			fv, unit := b.custom(b.tx)
			var sb strings.Builder
			sb.Grow(32)

			fmt.Fprintf(&sb, "%.2f", fv)
			sb.WriteString(unit)

			io.TX = sb.String()
		}

		addPostfix(&io.RX, &io.TX)
		return io
	}

	if b.unitFunc != nil {
		// RX
		{
			fv, unit, perSec := b.unitFunc(b.rx)

			var sb strings.Builder
			sb.Grow(32)

			fmt.Fprintf(&sb, "%.2f", fv)
			sb.WriteString(unit.String())
			if perSec {
				sb.WriteString("/s")
			}

			io.RX = sb.String()
		}

		// TX
		{
			fv, unit, perSec := b.unitFunc(b.tx)

			var sb strings.Builder
			sb.Grow(32)

			fmt.Fprintf(&sb, "%.2f", fv)
			sb.WriteString(unit.String())
			if perSec {
				sb.WriteString("/s")
			}

			io.TX = sb.String()
		}

		addPostfix(&io.RX, &io.TX)
		return io
	}

	{
		io.RX = fmt.Sprintf("%v", b.rx)
		io.TX = fmt.Sprintf("%v", b.tx)
	}

	addPostfix(&io.RX, &io.TX)
	return io
}
