package main

import (
	"fmt"
	_io "io"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/funkygao/golib/io"
	"github.com/funkygao/golib/pipestream"
	"github.com/funkygao/golib/signal"
)

func main() {
	startedAt = time.Now()

	tcpdumpFlag := []string{
		"-i",
		options.ifdev,
		"-nnN",
		"port",
		options.port,
	}
	td := pipestream.New("/usr/sbin/tcpdump", tcpdumpFlag...)
	if err := td.Open(); err != nil {
		panic(err)
	}

	signal.RegisterSignalHandler(syscall.SIGINT, func(sig os.Signal) {
		td.Close()
		showReportAndExit()
	})

	fmt.Printf("running /usr/sbin/tcpdump %s ...\n", strings.Join(tcpdumpFlag, " "))
	fmt.Println("Ctrl-C to stop")

	for {
		line, err := io.ReadLine(td.Reader())
		if err != nil {
			if err != _io.EOF {
				panic(err)
			}

			break
		}

		lines = append(lines, string(line))
		if len(lines) == options.max {
			td.Close()
			showReportAndExit()
		}
	}

	select {}
}
