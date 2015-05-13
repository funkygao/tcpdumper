package report

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

var errBadLine = errors.New("bad tcpdump output line")

func ShowReportAndExit(startedAt time.Time, lines []string) {
	fmt.Printf("%d lines, elapsed: %s\n", len(lines), time.Since(startedAt))

	for _, line := range lines {
		fmt.Println(line)
		src, dst, flag, err := lineInfo(line)
		if err != nil {
			continue
		}
		fmt.Println(src, dst, flag)
	}

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
