#!/bin/bash

mkdir -p ./release
cd ./cmd/rmsgnogpsd
env CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags=-d -o ../../release/rmsgnogpsd-linux-i386 ./rmsgnogpsd.go
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags=-d -o ../../release/rmsgnogpsd-linux-amd64 ./rmsgnogpsd.go
env CGO_ENABLED=0 GOOS=linux GOARCH=arm ARM=6 go build -ldflags=-d -o ../../release/rmsgnogpsd-linux-ARM6 ./rmsgnogpsd.go
env GOOS=windows GOARCH=386 go build -o ../../release/rmsgnogpsd-windows-i386.exe ./rmsgnogpsd.go
env GOOS=windows GOARCH=amd64 go build -o ../../release/rmsgnogpsd-windows-amd64.exe ./rmsgnogpsd.go
