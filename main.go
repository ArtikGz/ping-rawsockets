package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"ping/pingutils"
	"runtime"
	"syscall"
)

var (
	tcpflag = flag.Bool("-tcpscan", false, "Enable tcp port scan")
)

func main() {
	if runtime.GOOS != "linux" {
		fmt.Printf("[ERROR] This only works for linux!\n")
		os.Exit(1)
	}
	sigint := make(chan os.Signal, 1)

	// End when ctrl+c
	signal.Notify(sigint, syscall.SIGINT)
	go func() {
		<-sigint
		fmt.Printf("Exiting...\n")
		syscall.Exit(1)
	}()

	pingutils.GoPing()

	// After ping perform the TCP Port Scan if enabled
	if *tcpflag {
	}
}
