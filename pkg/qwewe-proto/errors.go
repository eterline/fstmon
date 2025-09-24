// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package qweweproto

type QweweProtoError string

func (e QweweProtoError) Error() string {
	return "qwewe protocol error: " + string(e)
}

const (
	ErrInvalidPacket  QweweProtoError = "invalid data structure"
	ErrInvalidLength  QweweProtoError = "invalid packet length information"
	ErrInvalidPayload QweweProtoError = "invalid packet data payload"
	ErrInvalidCRC     QweweProtoError = "invalid control sum CRC8"
)
