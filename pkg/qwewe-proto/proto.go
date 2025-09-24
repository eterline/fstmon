// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package qweweproto

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
)

// NewPacket - creates new qwewe protocol packet with wide area ids and zero payload buffer
func NewPacket() QwewePacket {
	return QwewePacket{
		SourceID: PKT_QW_WIDE_ID,
		TargetID: PKT_QW_WIDE_ID,
		Data:     &bytes.Buffer{},
		isValid:  false,
	}
}

// =============================================================================

// Length - get lenght of packet without CRC and start byte
func (pkt QwewePacket) Length() uint {

	if !pkt.isValid {
		return 0
	}

	if pkt.Data != nil {
		return 2 + uint(pkt.Data.Len())
	}

	return 2
}

// CRC - calculates CRC8 from packet data
func (pkt QwewePacket) CRC() byte {

	if !pkt.isValid {
		return 0x00
	}

	crcByte := InitCRCByte(CRC8)
	crcByte.Wite([]byte{pkt.SourceID, pkt.TargetID})
	crcByte.Wite(pkt.Data.Bytes())

	return crcByte.Complete()
}

// =============================================================================

// Marshal - packet marshaling to byte array fro uploading with separate byte, length and crc8
func (pkt QwewePacket) Marshal() ([]byte, error) {

	if !pkt.isValid {
		return nil, ErrInvalidPacket
	}

	// [\n:(1 byte)] + [len:(2 byte)] + [ids and payload:(N bytes)] + [CRC8:(1 byte)]
	length := pkt.Length()
	totalLength := length + 4
	data := make([]byte, totalLength, totalLength)

	data[0] = PKT_QW_START
	data[1] = byte(length >> 8) // upload hight byte for len
	data[2] = byte(length)      // upload low byte for len
	data[3] = pkt.SourceID
	data[4] = pkt.TargetID
	if totalLength > 4 {
		copy(data[5:], pkt.Data.Bytes()) // copy payload to data
	}
	data[totalLength-1] = pkt.CRC()

	return data, nil
}

// WriteTo - marsalling packet to bytes and write it io.Writer
func (pkt QwewePacket) WriteTo(w io.Writer) (int, error) {
	data, err := pkt.Marshal()
	if err != nil {
		return 0, err
	}
	return w.Write(data)
}

// Write - write bytes to packet payload buffer
func (pkt *QwewePacket) Write(p []byte) (int, error) {
	if pkt == nil || pkt.Data == nil {
		return 0, ErrInvalidPacket
	}

	size, err := pkt.Data.Write(p)
	if err != nil {
		return 0, err
	}

	pkt.isValid = true
	return size, nil
}

// Read - read bytes from packet payload buffer
func (pkt *QwewePacket) Read(p []byte) (int, error) {

	if pkt == nil || pkt.Data == nil || !pkt.isValid {
		return 0, ErrInvalidPacket
	}

	return pkt.Data.Read(p)
}

// Upload - upload structure as byte array to packet payload buffer
func (pkt *QwewePacket) Upload(v any) error {
	return binary.Write(pkt.Data, binary.LittleEndian, v)
}

// Extract - parse structure from payload packet byte array buffer
func (pkt *QwewePacket) Extract(v any) error {
	return binary.Read(pkt.Data, binary.LittleEndian, v)
}

// =============================================================================

func (pkt QwewePacket) WideSource() bool {
	return pkt.SourceID == PKT_QW_WIDE_ID
}

func (pkt *QwewePacket) MustWideSource() {
	pkt.SourceID = PKT_QW_WIDE_ID
}

func (pkt QwewePacket) WideTarget() bool {
	return pkt.TargetID == PKT_QW_WIDE_ID
}

func (pkt *QwewePacket) MustWideTarget() {
	pkt.TargetID = PKT_QW_WIDE_ID
}

// =============================================================================

// ReadQwewe - reading packet structure from byte array in  io.Reader
func ReadQwewe(r io.Reader) (QwewePacket, error) {

	stream := bufio.NewReaderSize(r, 128)
	packet := NewPacket()

	// searching for packet start byte
	for {
		data, err := stream.ReadByte()
		if err != nil {
			return packet, err
		}

		if data == '\n' {
			break
		}
	}

	totalLength := uint16(0) // packet byte array len
	// reading packet len from header bytes
	for i := 1; i >= 0; i-- {
		data, err := stream.ReadByte()
		if err != nil {
			return packet, err
		}

		totalLength |= uint16(data) << (8 * i)
	}

	// check size value overflow
	if MAX_PACKET_SIZE < int(totalLength) {
		return packet, ErrInvalidLength
	}

	// init crc byte
	expectedCRC := InitCRCByte(CRC8)
	// reading bytes to buffer
	packetBuffer, err := stream.Peek(int(totalLength))
	if err != nil {
		return packet, err
	}

	// getting crc from buffer
	crc := packetBuffer[totalLength-1]
	// calculating crc from buffer
	if _, err := expectedCRC.
		Wite(packetBuffer[:totalLength-1]); err != nil {
		return packet, err
	}
	// crc sum test
	if expectedCRC.Complete() != crc {
		return packet, ErrInvalidCRC
	}

	if err := QweweFromBytes(&packet, packetBuffer[:totalLength-1]); err != nil {
		return packet, err
	}

	return packet, nil
}

func QweweFromBytes(pkt *QwewePacket, data []byte) error {
	if len(data) < 3 {
		return ErrInvalidLength
	}

	pkt.isValid = true
	pkt.SourceID = data[0]
	pkt.TargetID = data[1]
	pkt.Data = bytes.NewBuffer(data[2:])

	return nil
}
