BINARY_NAME := "bankdownloader"
export REGISTRY := "ghcr.io"
export IMAGE_NAME := "airtonix/bankdownloader"

default:
  @just --choose

help:
  @just --list

dev *ARGS:
  go run . \
    --config ./store/config.example.json \
    --history ./store/history.example.json \
    {{ARGS}}

build: 
  goreleaser build \
    --snapshot \
    --clean

release:
  goreleaser release \
    --clean \
    --skip-publish \
    --snapshot \
    --clean

publish:
  goreleaser release --clean

test:
  gotest -v \
    -failfast \
    -race \
    -coverpkg=./... \
    -covermode=atomic \
    -coverprofile=coverage.txt \
    ./...


lint:
  go vet ./...
  staticcheck ./...

setup:
  go install golang.org/x/tools/cmd/godoc@latest
  go install golang.org/x/tools/cmd/goimports@latest
  go install golang.org/x/tools/gopls@latest
  go install github.com/go-delve/delve/cmd/dlv@latest
  go install github.com/ramya-rao-a/go-outline@latest
  go install github.com/stamblerre/gocode@v1.0.0
  go install github.com/rogpeppe/godef@v1.1.2
  go install honnef.co/go/tools/cmd/staticcheck@latest
  
  go get .

workflow:="release"
job:="Build"
event:="push"
test_ci_build :
  act {{event}} \
    -s GITHUB_TOKEN="$(gh auth token)" \
    --platform ubuntu-22.04=catthehacker/ubuntu:act-22.04 \
    --eventpath .actevent.json \
    --workflows .github/workflows/{{workflow}}.yml \
    --job {{job}} 

docs:
  godocs -http=:6060