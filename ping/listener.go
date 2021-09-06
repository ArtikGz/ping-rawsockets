package ping

import (
	"fmt"
	"sync"
	"syscall"
)

func Listener(wg *sync.WaitGroup, confirm chan uint8, comunicationChan chan uint8, TTLChan chan uint8) {
	defer wg.Done()
	// Send the confirmation that the listener is on
	confirm <- 1
	for i := 0; i < *times; i++{
		fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
		if err != nil {
			panic(fmt.Sprintf("An error happened while creating the socket: %v", err))
		}

		// Set timeout to 1 second
		timeval := syscall.Timeval{
			Sec:  1,
			Usec: 1,
		}
		syscall.SetsockoptTimeval(fd, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &timeval)

		// Read the echo reply
		data := make([]byte, 128) // Recolect packet data to extract TTL from IPv4 header
		_, _, _, _, err = syscall.Recvmsg(fd, data, nil, 0)

		// Timeout causes err != nil
		if err != nil {
			// STATUS
			comunicationChan <- 1
			// TTL (Avoid blocking)
			TTLChan <- 0
		} else {
			// STATUS
			comunicationChan <- 0
			// TTL
			TTLChan <- uint8(data[8])
		}
	}
}
