package report

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/funkygao/golib/color"
)

const (
	SYN_SEND = ">S"
	SYN_RECV = "<S"
	FIN      = "F"
	RST      = "R"
	PUSH     = "P"
)

// a single tcp packet trip
type trip struct {
	src, dst, flag string
}

type report map[string][]trip // key is endpoint

func ShowReportAndExit(startedAt time.Time, lines []string, port string) {
	var rp = make(report, 1<<16)

	t1 := time.Now()
	errorLines := make([]string, 0)
	for _, line := range lines {
		src, dst, flag, err := lineInfo(line)
		if err != nil {
			errorLines = append(errorLines, err.Error())
			continue
		}

		if _, present := rp[src]; present {
			rp[src] = append(rp[src], trip{src, dst, ">" + flag})
		} else {
			rp[src] = []trip{trip{src, dst, ">" + flag}}
		}

		if _, present := rp[dst]; present {
			rp[dst] = append(rp[dst], trip{src, dst, "<" + flag})
		} else {
			rp[dst] = []trip{trip{src, dst, "<" + flag}}
		}

	}

	retransmitSynN := 0
	pushN := 0
	rstN := 0
	incompleteHandshakeN := 0
	port = "." + port
	for endpoint, trips := range rp {
		fmt.Printf("%21s", endpoint)
		if strings.HasSuffix(endpoint, port) {
			fmt.Printf(" skipped\n")
			continue
		}

		syncSentN := 0
		fin := false
		rst := false
		sentSync := false
		recvSync := false
		push := false
		for _, t := range trips {
			if t.flag == SYN_SEND {
				sentSync = true
				syncSentN++
				if syncSentN > 1 {
					// retransmit
					retransmitSynN++
					t.flag = color.Blue(t.flag)
				}
			} else if strings.Contains(t.flag, SYN_RECV) {
				recvSync = true
			} else if strings.Contains(t.flag, FIN) {
				if !fin {
					t.flag = color.Red(t.flag)
				}
				fin = true
			} else if strings.Contains(t.flag, RST) {
				if !rst {
					rstN++
				}
				rst = true
				t.flag = color.Yellow(t.flag)
			} else if strings.Contains(t.flag, PUSH) {
				push = true
			}
			fmt.Printf(" %-3s", t.flag)
		}

		if sentSync && !recvSync {
			incompleteHandshakeN++
		}
		if push {
			pushN++
		}

		fmt.Println()
	}

	endpointN := len(rp)
	if endpointN > 1 {
		endpointN-- // the skipped endpoint excluded
	}

	for _, err := range errorLines {
		fmt.Println(err)
	}

	fmt.Println(strings.Repeat("=", 78))
	fmt.Printf("%d lines, elapsed: %s processed: %s\n", len(lines),
		time.Since(startedAt),
		time.Since(t1))
	fmt.Println(strings.Repeat("=", 78))
	fmt.Printf("%25s:%8d\n", "endpoint", endpointN)
	fmt.Printf("%25s:%8d\n", "incomplete handshakes", incompleteHandshakeN)
	fmt.Printf("%25s:%8d\n", "no PUSH", endpointN-pushN)
	fmt.Printf("%25s:%8d\n", "PUSH", pushN)
	fmt.Printf("%25s:%8d\n", "SYN retry", retransmitSynN)
	fmt.Printf("%25s:%8d\n", "RST", rstN)

	os.Exit(0)
}

func lineInfo(line string) (src, dst, flag string, err error) {
	parts := strings.Split(line, " ")
	if len(parts) < 7 {
		err = errors.New(fmt.Sprintf("bad line: %s", line))
		return
	}

	src = parts[2]
	dst = parts[4][:len(parts[4])-1]
	flag = parts[6][1 : len(parts[6])-2]
	return
}
