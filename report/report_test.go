package report

import (
	"testing"
)

func TestLineInfo(t *testing.T) {
	var lines []string = []string{
		"08:13:45.138216 IP 10.1.169.236.3298 > 10.77.144.193.10080: Flags [S], seq 2713400699, win 65535, options [mss 1460,nop,wscale 5,nop,nop,TS val 628283432 ecr 0,sackOK,eol], length 0",
		"08:13:45.138325 IP 10.77.144.193.10080 > 10.1.169.236.3298: Flags [S.], seq 3740416998, ack 2713400700, win 14480, options [mss 1460,sackOK,TS val 3595656504 ecr 628283432,nop,wscale 7], length 0",
		"08:13:45.143712 IP 10.1.169.236.3298 > 10.77.144.193.10080: Flags [.], ack 1, win 4117, options [nop,nop,TS val 628283437 ecr 3595656504], length 0",
		"08:13:45.144665 IP 10.1.169.236.3298 > 10.77.144.193.10080: Flags [P.], seq 1:84, ack 1, win 4117, options [nop,nop,TS val 628283437 ecr 3595656504], length 83",
	}
	for _, line := range lines {
		src, dst, flag, err := lineInfo(line)
		if err != nil {
			t.Error(err)
		}

		t.Logf("%s %s %s", src, dst, flag)
	}
}
