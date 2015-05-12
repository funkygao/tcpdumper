package main

import (
	"flag"
	"os"
	"time"
)

var (
	options struct {
		ifdev string
		port  string
	}

	startedAt time.Time
)

func init() {
	flag.StringVar(&options.ifdev, "i", "eth0", "device")
	flag.StringVar(&options.port, "p", "", "port")

	flag.Parse()

	if options.port == "" {
		flag.Usage()
		os.Exit(0)
	}
}
