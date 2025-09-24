// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package hostdata

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
)

type RecordDMI struct {
	Handle      uint16
	DMItype     uint8
	Size        uint8
	Description string
	Values      map[string]string
}

func ReadDMI() ([]RecordDMI, error) {
	data, err := exec.Command("dmidecode").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to read dmidecode: %v", err)
	}

	DMIs := getRecordsDMI(data)
	recs := make([]RecordDMI, len(DMIs))

	for i, rec := range DMIs {
		recs[i] = parseRecordDMI(rec)
	}

	return recs, nil
}

func parseRecordDMI(p []byte) RecordDMI {
	rec := RecordDMI{}

	lines := bytes.Split(p, []byte{'\n'})
	if len(lines) < 4 {
		return rec
	}

	rec.Handle = parseHandle(lines[0])
	rec.DMItype = parseDmiType(lines[1])
	rec.Size = parseSize(lines[2])
	rec.Description = string(lines[3])

	if len(lines) < 5 {
		return rec
	}

	lines = lines[3:]
	rec.Values = map[string]string{}

	for _, line := range lines {
		if !bytes.Contains(line, []byte{':', ' '}) {
			continue
		}

		v := bytes.Split(line, []byte{':', ' '})
		rec.Values[string(v[0])] = string(v[1])
	}

	return rec
}

func getRecordsDMI(p []byte) [][]byte {
	p = formatDmiOutputBytes(p)
	return splitSegments(p)
}

func parseDmiType(p []byte) uint8 {
	if !bytes.HasPrefix(p, []byte("DMI type")) {
		return 0
	}

	v := bytes.Fields(p)[2]
	value, _ := strconv.ParseUint(string(v), 10, 8)
	return uint8(value)
}

func parseSize(p []byte) uint8 {
	if !bytes.HasSuffix(p, []byte("bytes")) {
		return 0
	}

	v := bytes.Fields(p)[0]
	value, _ := strconv.ParseUint(string(v), 10, 8)
	return uint8(value)
}

func parseHandle(p []byte) uint16 {
	if !bytes.HasPrefix(p, []byte("Handle")) {
		return 0
	}

	v := bytes.Fields(p)[1]
	value, _ := strconv.ParseUint(string(v), 0, 16)
	return uint16(value)
}

func formatDmiOutputBytes(p []byte) []byte {

	for i, pByte := range p {
		if pByte == ',' {
			p[i] = '\n'
		}
		if pByte == '\t' {
			p[i] = ' '
		}
	}

	buf := bytes.NewBuffer(make([]byte, len(p)))
	lines := bytes.Split(p, []byte("\n"))

	for i, line := range lines {
		lines[i] = clrSideSpaces(line)
		buf.Write(lines[i])
		buf.WriteByte('\n')
	}

	return buf.Bytes()
}

func splitSegments(p []byte) [][]byte {
	segments := bytes.Split(p, []byte("\n\n"))
	startBytes := []byte("Handle")
	res := [][]byte{}

	for _, seg := range segments {
		if bytes.HasPrefix(seg, startBytes) {
			res = append(res, seg)
		}
	}

	return res
}
