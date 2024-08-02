#!/bin/bash
cd "$(pwd)" || exit

export NAME="gotools"

#x64
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
go build -o $NAME-$GOOS-$GOARCH

#aarch64
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=arm64
go build -o $NAME-$GOOS-$GOARCH

#windows
export CGO_ENABLED=0
export GOOS=windows
export GOARCH=amd64
go build -o $NAME-$GOOS-$GOARCH.exe