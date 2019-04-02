#! /bin/bash

goos="linux"

case "$1" in
  windows)
    goos="windows"
    ;;
  macos)
    goos="darwin"
    ;;
esac

CGO_ENABLED=0 GOOS=$goos \
go build -a -installsuffix cgo -o build/app ./src
