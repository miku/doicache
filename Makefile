SHELL = /bin/bash
VERSION := $(shell git describe --always --long --dirty)

doicache: cmd/doicache/main.go fetch.go home.go
	go build -v -ldflags="-X main.version=${VERSION}" -o $@ $<

clean:
	rm -f doicache

install: doicache
	mkdir -p $(HOME)/bin
	cp $< $(HOME)/bin

