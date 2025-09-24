// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package procf

import (
	"bytes"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	errUint8  uint8  = 255
	errUint16 uint16 = 65535
	errUint32 uint32 = 4294967295
	errUint64 uint64 = 18446744073709551615

	errInt8  int8  = -128
	errInt16 int16 = -32768
	errInt32 int32 = -2147483648
	errInt64 int64 = -9223372036854775808

	errFloat32 float32 = 1.401298464e-45
	errFloat64 float64 = 4.9406564584124654e-324
)

type setLine []byte

func (l setLine) isNil() bool {
	return len(l) == 0 || l == nil
}

func (l setLine) String() string {
	if l.isNil() {
		return ""
	}
	return bytesToString(l)
}

func (l setLine) StringList(delim string) []string {
	return strings.Split(l.String(), delim)
}

func (l setLine) Int32() int32 {
	return bytesToInt32(l)
}

func (l setLine) Int64() int64 {
	return bytesToInt64(l)
}

func (l setLine) Uint32() uint32 {
	value, err := strconv.ParseUint(l.String(), 10, 32)
	if err != nil {
		return errUint32
	}

	return uint32(value)
}

func (l setLine) Uint64() uint64 {
	value, err := strconv.ParseUint(l.String(), 10, 64)
	if err != nil {
		return errUint64
	}

	return value
}

func (l setLine) Float32() float32 {
	value, err := strconv.ParseFloat(l.String(), 32)
	if err != nil {
		return errFloat32
	}

	return float32(value)
}

func (l setLine) Float64() float64 {
	value, err := strconv.ParseFloat(l.String(), 64)
	if err != nil {
		return errFloat64
	}

	return value
}

func (l setLine) ParseFunc(f func(v string) string) setLine {
	res := f(l.String())
	return setLine(res)
}

func (l setLine) ParseBool(yes string) bool {
	return l.String() == yes
}

type FileDataSetParser struct {
	dataSrc []byte
	dataMap map[string]setLine
}

func NewFileDataSetParser(data, lineEnd, paramDelim []byte, keyPos, paramPos int) *FileDataSetParser {

	var (
		setLines      = bytes.Split(data, lineEnd)
		setLinesCount = len(setLines)
		resMap        = make(map[string]setLine, setLinesCount)
	)

	for _, setLine := range setLines {
		lineFields := bytes.Split(setLine, paramDelim)
		if len(lineFields) < (paramPos + 1) {
			continue
		}

		key := clrKeyBytes(lineFields[keyPos])

		resMap[string(key)] = clrSideSpaces(lineFields[paramPos])
	}

	return &FileDataSetParser{
		dataSrc: data,
		dataMap: resMap,
	}
}

func clrKeyBytes(p []byte) []byte {
	result := p[:0]

	for _, b := range p {
		if !unicode.IsSpace(rune(b)) && b >= 0x20 {
			result = append(result, b)
		}
	}

	return result
}

func clrSideSpaces(p []byte) []byte {
	lPtr, rPtr := 0, len(p)-1

	for lPtr <= rPtr && p[lPtr] == ' ' {
		lPtr++
	}

	for rPtr >= lPtr && p[rPtr] == ' ' {
		rPtr--
	}

	if lPtr > rPtr {
		return []byte{}
	}

	return p[lPtr : rPtr+1]
}

func (prs *FileDataSetParser) Param(key string) setLine {
	return prs.dataMap[key]
}

func (prs *FileDataSetParser) Count() int {
	return len(prs.dataMap)
}

func (prs *FileDataSetParser) Data() map[string]setLine {
	return prs.dataMap
}

func bytesToInt32(p []byte) int32 {

	if len(p) == 0 {
		return errInt32
	}

	start := 0
	negative := false

	if p[0] == '-' {
		negative = true
		start = 1
	} else if p[0] == '+' {
		start = 1
	}

	if start >= len(p) {
		return errInt32
	}

	var result int32 = 0
	for i := start; i < len(p); i++ {
		b := p[i]
		if b < '0' || b > '9' {
			return errInt32
		}
		digit := int32(b - '0')

		if result > (2147483647-digit)/10 {
			return errInt32
		}

		result = result*10 + digit
	}

	if negative {
		result = -result
	}

	return result
}

func bytesToInt64(p []byte) int64 {

	if len(p) == 0 {
		return errInt64
	}

	start := 0
	negative := false

	if p[0] == '-' {
		negative = true
		start = 1
	} else if p[0] == '+' {
		start = 1
	}

	if start >= len(p) {
		return errInt64
	}

	var result int64 = 0
	for i := start; i < len(p); i++ {
		b := p[i]
		if b < '0' || b > '9' {
			return errInt64
		}
		digit := int64(b - '0')

		if result > (2147483647-digit)/10 {
			return errInt64
		}

		result = result*10 + digit
	}

	if negative {
		result = -result
	}

	return result
}

func bytesToUint32(p []byte) uint32 {
	if len(p) == 0 {
		return errUint32
	}

	var result uint32 = 0
	for i := 0; i < len(p); i++ {
		b := p[i]
		if b < '0' || b > '9' {
			return errUint32
		}
		digit := uint32(b - '0')

		if result > (4294967295-digit)/10 {
			return errUint32
		}

		result = result*10 + digit
	}

	return result
}

func bytesToUint64(p []byte) uint64 {
	if len(p) == 0 {
		return errUint64
	}

	var result uint64 = 0
	for i := 0; i < len(p); i++ {
		b := p[i]
		if b < '0' || b > '9' {
			return errUint64
		}
		digit := uint64(b - '0')

		if result > (18446744073709551615-digit)/10 {
			return errUint64
		}

		result = result*10 + digit
	}

	return result
}

func bytesToFloat32(p []byte) float32 {
	v, err := strconv.ParseFloat(string(p), 10)
	if err != nil {
		return errFloat32
	}
	return float32(v)
}

func bytesToFloat64(p []byte) float64 {
	v, err := strconv.ParseFloat(string(p), 10)
	if err != nil {
		return errFloat64
	}
	return v
}

func bytesToString(p []byte) string {

	for len(p) > 0 {
		r, size := utf8.DecodeRune(p)
		if !unicode.IsSpace(r) {
			break
		}
		p = p[size:]
	}

	for len(p) > 0 {
		r, size := utf8.DecodeLastRune(p)
		if !unicode.IsSpace(r) {
			break
		}
		p = p[:len(p)-size]
	}

	return string(p)
}

func belongsWithPrefix(prefixList []string, s string) bool {
	for _, prefix := range prefixList {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}
