#!/bin/sh

# build windows
GOOS=windows GOARCH=amd64 go build -o bin/q2-windows-amd64.exe

# build linux
GOOS=linux GOARCH=amd64 go build -o bin/q2-linux-amd64

# build mac
GOOS=darwin GOARCH=amd64 go build -o bin/q2-darwin-amd64