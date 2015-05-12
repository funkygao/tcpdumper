package main

import (
	"fmt"
	"os"
	"time"
)

func showReportAndExit() {
	fmt.Printf("%d lines, elapsed: %s\n", len(lines), time.Since(startedAt))

	fmt.Println(lines)

	os.Exit(0)

}
