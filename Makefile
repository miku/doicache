SHELL = /bin/bash
VERSION := $(shell git describe --always --long --dirty)

doicache: cmd/doicache/main.go
	go build -i -v -ldflags="-X main.version=${VERSION}" -o $@ $<

clean:
	rm -f doicache
