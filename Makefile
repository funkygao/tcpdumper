all:linux mac

linux:
	GOOS=linux GOARCH=amd64 go build -o bin/tcpdumper.linux

mac:
	GOOS=darwin GOARCH=amd64 go build -o bin/tcpdumper.mac

