package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"reutility/ping"
	"reutility/tcpscan"
	"runtime"
	"syscall"
)

var (
	tcpflag = flag.Bool("-tcpscan", true, "Enable tcp port scan")
)

func init() {
	flag.Parse()
}

func main() {
	if runtime.GOOS != "linux" {
		fmt.Printf("[ERROR] This only works for linux!\n")
		os.Exit(1)
	}

	// End when ctrl+c
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT)
	go func() {
		<-sigint
		fmt.Printf("Exiting...\n")
		syscall.Exit(1)
	}()

	ping.GoPing()

	// After ping perform the TCP Port Scan if enabled
	tcpscan.Tcpscan()
}
