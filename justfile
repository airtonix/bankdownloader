BINARY_NAME := "bankdownloader"
export REGISTRY := "ghcr.io"
export IMAGE_NAME := "airtonix/bankdownloader"


dev *ARGS:
  go run . \
    --config ./store/config.example.json \
    --history ./store/history.example.json \
    {{ARGS}}

build: 
  goreleaser build --snapshot --clean

release:
  goreleaser release --clean --skip-publish --snapshot --clean

test:
  for PACKAGE in $(go list ./...); do go test -v ${PACKAGE}; done;

lint:
  #!/bin/sh
  if [ "$(gofmt -s -l . | wc -l)" -gt 0 ];
    then exit 1;
  fi

setup:
  go install golang.org/x/tools/cmd/godoc@latest
  go install golang.org/x/tools/cmd/goimports@latest
  go install golang.org/x/tools/gopls@latest
  go install github.com/go-delve/delve/cmd/dlv@latest
  go install github.com/ramya-rao-a/go-outline@latest
  go install github.com/stamblerre/gocode@v1.0.0
  go install github.com/rogpeppe/godef@v1.1.2
  go get .

ci:
  act push

docs:
  godocs -http=:6060