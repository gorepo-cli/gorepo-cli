VERSION := $(shell git describe --tags)

build:
	go build -ldflags="-X 'gorepo-cli/internal/commands.version=$(VERSION)'" -o bin/gorepo ./cmd/cli
