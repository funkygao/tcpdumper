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

	totalRetransmitSynN := 0
	totalPushSessionN := 0
	totalAckN := 0
	totalPushN := 0
	totalFinSessionN := 0
	totalRstSessionN := 0
	totalIncompleteHandshakeN := 0
	for endpoint, trips := range rp {
		// the lifespan of an endpoint
		fmt.Printf("%21s", endpoint)

		if strings.HasSuffix(endpoint, "."+port) {
			fmt.Printf(" skipped\n")
			continue
		}

		syncSentN := 0
		fin := false
		rst := false
		push := false
		sentSync := false
		recvSync := false
		for _, t := range trips {
			if strings.Contains(t.flag, ".") {
				totalAckN++
			}

			if t.flag == SYN_SEND {
				sentSync = true
				syncSentN++
				if syncSentN > 1 {
					// retransmit SYN is blue
					totalRetransmitSynN++
					t.flag = color.Blue(t.flag)
				}
			} else if strings.Contains(t.flag, SYN_RECV) {
				recvSync = true
			} else if strings.Contains(t.flag, FIN) {
				if !fin {
					// the first FIN is red
					t.flag = color.Red(t.flag)
					totalFinSessionN++
				}
				fin = true
			} else if strings.Contains(t.flag, RST) {
				if !rst {
					totalRstSessionN++
					if fin {
						totalFinSessionN--
					}
				}
				rst = true
				t.flag = color.Yellow(t.flag)
			} else if strings.Contains(t.flag, PUSH) {
				push = true
				totalPushN++
			}

			fmt.Printf(" %-3s", t.flag)
		}

		if sentSync && !recvSync {
			totalIncompleteHandshakeN++
		}
		if push {
			totalPushSessionN++
		}

		fmt.Println()
	}

	totalSessionN := len(rp)
	if totalSessionN > 1 {
		totalSessionN-- // the skipped endpoint excluded
	}

	for _, err := range errorLines {
		fmt.Println(err)
	}

	fmt.Println(strings.Repeat("=", 78))
	fmt.Printf("%d lines, elapsed: %s processed: %s\n", len(lines),
		time.Since(startedAt),
		time.Since(t1))
	fmt.Println(strings.Repeat("=", 78))
	fmt.Printf("%25s:%12d\n", "sessions", totalSessionN)
	fmt.Printf("%25s:%12d\n", "PUSH", totalPushSessionN)
	fmt.Printf("%25s:%12d\n", "incomplete handshakes",
		totalIncompleteHandshakeN)
	fmt.Printf("%25s:%12d\n", "SYN retry", totalRetransmitSynN)
	fmt.Printf("%25s:%12d\n", "FIN", totalFinSessionN)
	fmt.Printf("%25s:%12d\n", "RST", totalRstSessionN)
	fmt.Println()
	fmt.Printf("%25s:%12d\n", "data sent", totalPushN)
	fmt.Printf("%25s:%12d\n", "ack sent", totalAckN)
	fmt.Printf("%25s:%12.2f%%\n", "data/ack",
		100*float32(totalPushN)/float32(totalAckN))

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
