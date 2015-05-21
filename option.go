package main

import (
	"flag"
	"os"
	"time"
)

var (
	options struct {
		ifdev   string
		port    string
		rawfile string
		mute    bool
		max     int
	}

	startedAt time.Time

	lines []string = make([]string, 0, 1<<20)
)

func init() {
	flag.StringVar(&options.ifdev, "i", "eth0", "device")
	flag.StringVar(&options.port, "p", "", "port")
	flag.BoolVar(&options.mute, "mute", true, "dont show detailed session info")
	flag.StringVar(&options.rawfile, "f", "tcpdump.out", "output file of tcpdump")
	flag.IntVar(&options.max, "max", 1<<20, "max num of tcpdump output lines to collect")

	flag.Parse()

	if options.port == "" {
		flag.Usage()
		os.Exit(0)
	}
}
