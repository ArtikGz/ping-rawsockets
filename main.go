package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

var (
	times           = flag.Int("c", 4, "Number of packets you want to send")
	flagip          = flag.String("h", "127.0.0.1", "Host you want to ping")
	ip              []byte
	missedPackets   uint32
	recievedPackets uint32
)

func init() {
	flag.Parse()
	ip = validateHostFlag()
}

func main() {
	if runtime.GOOS != "linux" {
		fmt.Printf("[ERROR] This only works for linux!\n")
		os.Exit(1)
	}

	// Channels and WaitGroup for communicating between goroutines
	var wg sync.WaitGroup
	wg.Add(2)
	confirmListener := make(chan uint8)
	communicationChan := make(chan uint8)
	TTLChan := make(chan uint8)
	sigint := make(chan os.Signal, 1)

	// End when ctrl+c
	signal.Notify(sigint, syscall.SIGINT)
	go func() {
		<-sigint
		exit(1)
		syscall.Exit(1)
	}()

	go listener(&wg, confirmListener, communicationChan, TTLChan)

	// Block until listener is up
	<-confirmListener
	go ping(&wg, communicationChan, TTLChan)
	wg.Wait()
}

func ping(wg *sync.WaitGroup, comunicationChan chan uint8, TTLChan chan uint8) {
	defer wg.Done()
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		panic(fmt.Sprintf("Could't create fd on ping(): %v", err))
	}

	// Setup the headers
	headers := Ip4Headers{
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
	exit(0)
}

func listener(wg *sync.WaitGroup, confirm chan uint8, comunicationChan chan uint8, TTLChan chan uint8) {
	defer wg.Done()
	// Send the confirmation that the listener is on
	confirm <- 1
	for {
		fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
		if err != nil {
			panic(fmt.Sprintf("An error happened when trying to create the socket: %v", err))
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
