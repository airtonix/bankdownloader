BINARY_NAME := "bankdownloader"

dev *ARGS:
  go run . {{ARGS}}

build: 
  go build -o dist/

test:
  for PACKAGE in $(go list ./...); do go test -v ${PACKAGE}; done;

lint:
  #!/bin/sh
  if [ "$(gofmt -s -l . | wc -l)" -gt 0 ];
    then exit 1;
  fi

setup:
  go install github.com/go-delve/delve/cmd/dlv@latest
  go install golang.org/x/tools/cmd/godoc@latest
  go get .

ci:
  act push

docs:
  godocs -http=:6060