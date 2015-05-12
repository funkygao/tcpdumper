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
		max   int
	}

	startedAt time.Time

	lines []string = make([]string, 0, 1<<20)
)

func init() {
	flag.StringVar(&options.ifdev, "i", "eth0", "device")
	flag.StringVar(&options.port, "p", "", "port")
	flag.IntVar(&options.max, "max", 1<<20, "max num of tcpdump output lines to collect")

	flag.Parse()

	if options.port == "" {
		flag.Usage()
		os.Exit(0)
	}
}
