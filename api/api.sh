#!/usr/bin/env bash

# generate the gRPC code
protoc  -I/usr/local/include -I. --go_out=plugins=grpc:. ./faas.proto

# generate the JSON interface code
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --go_out=plugins=grpc:. \
  ./faas.proto

# generate the swagger definitions
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --swagger_out=logtostderr=true:../swagger \
  ./faas.proto
