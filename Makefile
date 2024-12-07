VERSION := $(shell git describe --tags)

test:
	go test ./...

build:
	go build -ldflags="-X 'main.version=$(VERSION)'" -o bin/gorepo .
