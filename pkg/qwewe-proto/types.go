// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package qweweproto

import (
	"bytes"
	"fmt"
)

// =============================================================================

type QweweMeanByte byte

const (
	PKT_QW_START   byte = '\n'
	PKT_QW_WIDE_ID byte = 0xFF
	PKT_QW_ZERO_ID byte = 0x00

	TYPE_QW_ACK     QweweMeanByte = 0x01
	TYPE_QW_DATA    QweweMeanByte = 0x02
	TYPE_QW_COMMAND QweweMeanByte = 0x03
)

const (
	MAX_PACKET_SIZE int = 1024
)

// =============================================================================

type QwewePacket struct {
	SourceID byte
	TargetID byte
	Data     *bytes.Buffer

	isValid bool
}

func (p QwewePacket) String() string {
	return fmt.Sprintf(
		"Qwewe packet unit:\nLength = %d\nsrc = 0x%X\ntrg = 0x%X\nCRC8 = 0x%X\npayload:\n{%v}",
		p.Length(), p.SourceID, p.TargetID, p.CRC(), p.Data.Bytes(),
	)
}
