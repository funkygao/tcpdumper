package report

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/funkygao/golib/color"
)

var errBadLine = errors.New("bad tcpdump output line")

const (
	SYNC_SEND = ">S"
	FIN       = "F"
)

type trip struct {
	src, dst, flag string
}

type report map[string][]trip

func ShowReportAndExit(startedAt time.Time, lines []string) {
	fmt.Printf("%d lines, elapsed: %s\n", len(lines), time.Since(startedAt))

	var rp = make(report, 1<<10)

	for _, line := range lines {
		src, dst, flag, err := lineInfo(line)
		if err != nil {
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

	retransmitSync := 0
	for endpoint, trips := range rp {
		fmt.Printf("%21s", endpoint)
		if len(trips) > 1000 {
			fmt.Printf(" skipped\n")
			continue
		}

		syncSentN := 0
		finN := 0
		for _, t := range trips {
			if t.flag == SYNC_SEND {
				syncSentN++
				if syncSentN > 1 {
					// retransmit
					retransmitSync++
					t.flag = color.Colorize([]string{color.FgBlue, color.Blink},
						t.flag)
				} else {
					t.flag = color.Blue(t.flag)
				}
			} else if strings.Contains(t.flag, FIN) {
				finN++
				if finN == 1 {
					t.flag = color.Red(t.flag)
				}

			}
			fmt.Printf(" %-2s", t.flag)
		}
		fmt.Println()
	}

	fmt.Printf("sync retrans: %d\n", retransmitSync)

	os.Exit(0)
}

func lineInfo(line string) (src, dst, flag string, err error) {
	parts := strings.Split(line, " ")
	if len(parts) < 7 {
		err = errBadLine
		return
	}

	src = parts[2]
	dst = parts[4][:len(parts[4])-1]
	flag = parts[6][1 : len(parts[6])-2]
	return
}
