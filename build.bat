@echo off
cd /d %~dp0

set NAME=gotools

::x64
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build -o %NAME%-%GOOS%-%GOARCH%

::aarch64
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=arm64
go build -o %NAME%-%GOOS%-%GOARCH%



::x64
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
go build -o %NAME%-%GOOS%-%GOARCH%.exe



