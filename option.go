package main

import (
	"flag"
	"os"
	"time"
)

var (
	options struct {
		ifdev string
		port  int
	}

	startedAt time.Time
)

func init() {
	flag.StringVar(&options.ifdev, "i", "eth0", "device")
	flag.IntVar(&options.port, "p", 0, "port")

	flag.Parse()

	if options.port == 0 {
		flag.Usage()
		os.Exit(0)
	}
}
