BINARY_NAME := "bankdownloader"

build: 
	go build -o bin/

test:
	for PACKAGE in $(go list ./...); do go test -v ${PACKAGE}; done;

lint:
	#!/bin/sh
	if [ "$(gofmt -s -l . | wc -l)" -gt 0 ];
		then exit 1;
	fi

setup:
	go install github.com/go-delve/delve/cmd/dlv@latest
	go get .