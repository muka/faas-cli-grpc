.PHONY: build clean deps docker/push docker/build

VERSION := $(shell git describe --tags --always)
GOOS ?= linux
GOARCH ?= amd64

deps:
	@go generate api/api.go

clean:

build: deps

docker/build: build

docker/push: docker/build
