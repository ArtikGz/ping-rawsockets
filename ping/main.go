package ping

import (
	"flag"
	"fmt"
	"sync"
)

var (
	missedPackets   uint32
	recievedPackets uint32
	times           = flag.Int("c", 4, "Number of packets you want to send")
	tcpflag         = flag.Bool("tcpscan", false, "Enable TCP port scan")
	flagip          = flag.String("h", "127.0.0.1", "Host you want to ping")
	ip              []byte
)

func init() {
	flag.Parse()
	ip = validateHostFlag()
}

func GoPing() {
	// Channels and WaitGroup for communicating between goroutines
	var wg sync.WaitGroup
	wg.Add(2)
	confirmListener := make(chan uint8)
	communicationChan := make(chan uint8)
	TTLChan := make(chan uint8)

	go Listener(&wg, confirmListener, communicationChan, TTLChan)

	// Block until listener is up
	<-confirmListener
	go Ping(&wg, communicationChan, TTLChan)
	wg.Wait()
}

func Exit(exitStatus int) {
	totalSent := recievedPackets + missedPackets
	percentSent := (float32(totalSent) / float32(recievedPackets)) * 100
	fmt.Printf("Total sent: %d\tTotal recived: %d\tTotal lost: %d\t%.2f%% Recieved\n", totalSent, recievedPackets, missedPackets, percentSent)
}
