BINARY_NAME := "bankdownloader"

build: 
    go build \
        -o bin/ \
        -v

test:
    go test -v \
        ./...

setup:
    go run \
        github.com/playwright-community/playwright-go/cmd/playwright@latest \
        install \
        --with-deps