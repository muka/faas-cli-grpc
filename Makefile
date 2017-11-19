.PHONY: build clean deps docker/push docker/build

VERSION := $(shell git describe --tags --always)
GOOS ?= linux
GOARCH ?= amd64

deps:
	@go generate api/api.go

clean:

test:
	cd ./test && zip -r test1.zip test1.yml ./test1
	go test -v ./api

build: deps

docker/build: build

docker/push: docker/build
