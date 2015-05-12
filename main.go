package main

import (
	"bufio"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/funkygao/golib/pipestream"
	"github.com/funkygao/golib/signal"
)

var (
	lines []string = make([]string, 0, 1<<20)
)

func main() {
	tcpdumpFlag := fmt.Sprintf("-i %s -nnN port %d", options.ifdev, options.port)
	td := pipestream.New("tcpdump", tcpdumpFlag)
	if err := td.Open(); err != nil {
		panic(err)
	}

	startedAt = time.Now()

	signal.RegisterSignalHandler(syscall.SIGINT, func(sig os.Signal) {
		td.Close()
		showReport()
		os.Exit(0)
	})

	fmt.Printf("running tcpdump %s\n", tcpdumpFlag)
	fmt.Println("Ctrl-C to stop")

	scanner := bufio.NewScanner(td.Reader())
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	select {}
}

func showReport() {
	fmt.Printf("elapsed: %s\n", time.Since(startedAt))

	fmt.Println(lines)

}
