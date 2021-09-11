package headers

import "encoding/binary"

// DOCS: https://datatracker.ietf.org/doc/html/rfc793#section-3.1
type TCPHeader struct {
	SrcPort       uint16
	DstPort       uint16
	SequenceNum   uint32
	Ack           uint32
	DataOffset    int
	Reserved      int
	Flags         int
	Window        uint16
	Checksum      uint16
	UrgentPointer uint16
	Options       uint32
	Padding       uint32
}

func (header TCPHeader) Marshall() []byte {
	payload := make([]byte, 24)

	// Convert each port into BigEndian
	binary.BigEndian.PutUint16(payload[0:2], header.SrcPort)
	binary.BigEndian.PutUint16(payload[2:4], header.DstPort)

	// Sequence number
	binary.BigEndian.PutUint32(payload[4:8], header.SequenceNum)

	// ACK
	binary.BigEndian.PutUint32(payload[8:12], header.Ack)

	payload[12] = byte((header.DataOffset << 4) | ((header.Reserved >> 2) & 0xF))
	payload[13] = byte(((header.Reserved & 0xF) << 6) | header.Flags&0b111111)

	binary.BigEndian.PutUint16(payload[14:16], header.Window)
	binary.BigEndian.PutUint16(payload[16:18], header.Checksum)
	binary.BigEndian.PutUint16(payload[18:20], header.UrgentPointer)

	// Options + Padding
	optPadding := (header.Options << 8) | header.Padding
	binary.BigEndian.PutUint32(payload[20:24], optPadding)

	return payload
}
