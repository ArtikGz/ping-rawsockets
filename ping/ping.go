package ping

import (
	"fmt"
	"net"
	"reutility/headers"
	"sync"
	"syscall"
	"time"
)

func Ping(wg *sync.WaitGroup, comunicationChan chan uint8, TTLChan chan uint8) {
	defer wg.Done()
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		panic(fmt.Sprintf("Could't create fd on ping(): %v", err))
	}

	// Setup the headers
	headers := headers.Ip4Headers{
		Version:  4,
		IHL:      20,
		TLen:     20 + 8, // IPv4Header len + ICMPHeader len
		TTL:      64,
		Protocol: 1, // 1 -> ICMP
		Dst:      net.IPv4(ip[0], ip[1], ip[2], ip[3]),
	}

	// TODO: Move this to a struct
	ICMPHeader := []byte{
		8, // Echo request
		0, // Code
		0, // Chksum
		0, // Chksum
		0,
		0,
		0,
		0,
	}

	cs := csum(ICMPHeader)
	ICMPHeader[2] = byte(cs)
	ICMPHeader[3] = byte(cs >> 8)

	ipv4headers := headers.Marshall()
	payload := append(ipv4headers, ICMPHeader...)

	addr := syscall.SockaddrInet4{
		Port: 0,
		Addr: [4]byte{ip[0], ip[1], ip[2], ip[3]},
	}

	for i := 0; i < *times; i++ {
		t1 := time.Now()
		err = syscall.Sendto(fd, payload, 0, &addr)

		// Block until listener response
		// 0 -> Correct
		// 1 -> Timeout
		res := <-comunicationChan
		TTL := <-TTLChan

		if res == 0 {
			elapsed := time.Since(t1)
			fmt.Printf("Time: %.3fms\tTTL: %d (Probably OS: %s)\n", (float32(elapsed) / float32(time.Millisecond)), TTL, getOsFromTTL(TTL))
			recievedPackets++
			if i != *times-1 {
				time.Sleep(1*time.Second - elapsed)
			}
		} else {
			missedPackets++
			fmt.Printf("Resource unavailable\n")
		}
	}

	// Print analitics and exit
	Exit(0)
}
