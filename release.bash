#!/bin/bash

mkdir -p ./release
env CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags=-d -o ./release/rmsgnogpsd-linux-i386 ./cmd/rmsgnogpsd/rmsgnogpsd.go
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags=-d -o ./release/rmsgnogpsd-linux-amd64 ./cmd/rmsgnogpsd/rmsgnogpsd.go
env CGO_ENABLED=0 GOOS=linux GOARCH=ARM GOARM=6 go build -ldflags=-d -o ./release/rmsgnogpsd-linux-ARM6 ./cmd/rmsgnogpsd/rmsgnogpsd.go
env GOOS=windows GOARCH=386 go build -ldflags=-d -o ./release/rmsgnogpsd-windows-i386.exe ./cmd/rmsgnogpsd/rmsgnogpsd.go
env GOOS=windows GOARCH=amd64 go build -ldflags=-d -o ./release/rmsgnogpsd-windows-amd64.exe ./cmd/rmsgnogpsd/rmsgnogpsd.go
