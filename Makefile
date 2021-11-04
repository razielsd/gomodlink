BIN := "./bin/abf"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X 'github.com/razielsd/antibruteforce/app/cmd.version=develop' -X 'github.com/razielsd/antibruteforce/app/cmd.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S)' -X 'github.com/razielsd/antibruteforce/app/cmd.gitHash=$(GIT_HASH)'
LOCAL_PATH := $(shell pwd)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)"

test:
	go test ./internal/...

lint:
	golangci-lint run ./...

.PHONY: build test lint