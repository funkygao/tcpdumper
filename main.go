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
	"github.com/funkygao/tcpdumper/report"
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
		report.ShowReportAndExit(startedAt, lines, options.port, options.mute)
	})

	fmt.Printf("running /usr/sbin/tcpdump %s ...\n", strings.Join(tcpdumpFlag, " "))
	fmt.Println("Ctrl-C to stop")

	outfile, err := os.OpenFile(options.rawfile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer outfile.Close()

	n := 0
	for {
		line, err := io.ReadLine(td.Reader())
		if err != nil {
			if err != _io.EOF {
				panic(err)
			}

			break
		}

		if len(line) == 0 {
			// empty line
			continue
		}

		n++
		if n == options.max {
			td.Close()
			report.ShowReportAndExit(startedAt, lines, options.port, options.mute)
			return
		}

		linestr := string(line)
		lines = append(lines, linestr)
		outfile.WriteString(linestr + "\n")
		if n%500 == 0 {
			outfile.Sync()
		}

	}

	select {}
}
