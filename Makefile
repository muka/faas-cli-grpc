.PHONY: all build clean deps docker/push docker/build

VERSION := $(shell git describe --tags --always)
GOOS ?= linux
GOARCH ?= amd64

deps:
	@go generate api/api.go

clean:
    @rm -f swagger/*.json
	@rm -f api/faas.pb.go

build: deps

docker/build: build

docker/push: docker/build

all: build
