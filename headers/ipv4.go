package headers

import (
	"encoding/binary"
	"net"
)

// DOCS: https://en.wikipedia.org/wiki/IPv4#Header
type Ip4Headers struct {
	Version        int // Version (4 for IPv4)
	IHL            int // Header len
	DSCP           int // Differentiated Services Code Point OR Type of service
	TLen           int // Total length
	ID             int // Identification
	Flags          int
	FragmentOffset int
	TTL            int // Time to live
	Protocol       int // DOCS: https://en.wikipedia.org/wiki/List_of_IP_protocol_numbers
	Checksum       int
	Src            net.IP // Source IP
	Dst            net.IP // Destination IP
	Options        []byte // Only if IHL > 5
}

// Parse headers
func (headers Ip4Headers) Marshall() []byte {
	hdrlen := headers.IHL
	ipv4headers := make([]byte, hdrlen)

	ipv4headers[0] = byte(headers.Version<<4 | (headers.IHL >> 2 & 0x0f))
	// Unmarshall headers.Version = (ipv4headers[0] >> 4)
	// Unmarshall headers.IHL = ((ipv4headers[0] & 0x0F) << 2)
	ipv4headers[1] = byte(headers.DSCP)
	binary.BigEndian.PutUint16(ipv4headers[2:4], uint16(headers.TLen))
	binary.BigEndian.PutUint16(ipv4headers[4:6], uint16(headers.ID))

	flagsAndFragmentOffset := (headers.FragmentOffset & 0x1FFF) | (headers.Flags << 13)
	binary.BigEndian.PutUint16(ipv4headers[6:8], uint16(flagsAndFragmentOffset))
	ipv4headers[8] = byte(headers.TTL)
	ipv4headers[9] = byte(headers.Protocol)
	binary.BigEndian.PutUint16(ipv4headers[10:12], uint16(headers.Checksum))
	if i := headers.Src.To4(); i != nil {
		copy(ipv4headers[12:16], headers.Src.To4())
	}
	if i := headers.Dst.To4(); i != nil {
		copy(ipv4headers[16:20], headers.Dst.To4())
	}

	if len(headers.Options) > 0 {
		copy(ipv4headers[hdrlen:], headers.Options)
	}

	return ipv4headers
}
