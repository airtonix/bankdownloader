BINARY_NAME := "bankdownloader"

build: 
	go build -o bin/

test:
	for PACKAGE in $(go list ./...); do go test -v ${PACKAGE}; done;

setup:
	go install github.com/go-delve/delve/cmd/dlv@latest
	go get download