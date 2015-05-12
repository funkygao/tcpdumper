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

var (
	lines []string = make([]string, 0, 1<<20)
)

func main() {
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

	startedAt = time.Now()

	signal.RegisterSignalHandler(syscall.SIGINT, func(sig os.Signal) {
		td.Close()
		showReport()
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
			showReport()
		}
	}

	select {}
}

func showReport() {
	fmt.Printf("elapsed: %s\n", time.Since(startedAt))

	fmt.Println(lines)

	os.Exit(0)

}
