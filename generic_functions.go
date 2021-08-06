package main

import (
	"fmt"
	"os"
)

// param1: TTL
// Returns: Normaly by default Windows have 128 TTL and Linux 64 TTL
// It's not 100% accurate but it's usually like that
func getOsFromTTL(TTL uint8) string {
	if TTL <= 64 && TTL > 0 {
		return "Linux"
	} else if TTL <= 128 {
		return "Windows"
	} else {
		return "Unknown"
	}
}

func csum(b []byte) uint16 {
	var s uint32
	for i := 0; i < len(b); i += 2 {
		s += uint32(b[i+1])<<8 | uint32(b[i])
	}
	// add back the carry
	s = s>>16 + s&0xffff
	s = s + s>>16
	return uint16(^s)
}

func exit(exitStatus int) {
	totalSent := recievedPackets + missedPackets
	percentSent := (float32(totalSent) / float32(recievedPackets)) * 100
	fmt.Printf("Total sent: %d\tTotal recived: %d\tTotal lost: %d\t%.2f%% Recieved\n", totalSent, recievedPackets, missedPackets, percentSent)
	os.Exit(exitStatus)
}
