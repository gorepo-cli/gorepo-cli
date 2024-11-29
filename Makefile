VERSION := $(shell git describe --tags)

build:
	go build -ldflags="-X 'main.version=$(VERSION)'" -o bin/gorepo ./src
