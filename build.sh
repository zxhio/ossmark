#!/bin/bash

set -e

build_os_arch() {
    os=$1
    arch=$2

    echo "+build ossmark $os/$arch"

    rm -rf ossmark && cp -r release ossmark && mkdir -p ossmark/bin
    CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -ldflags "-w -s" -o ossmark/bin/ossmark-article cmd/ossmark-article/main.go
    CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -ldflags "-w -s" -o ossmark/bin/ossmark-sync cmd/ossmark-sync/main.go
    tar -czf ossmark."$os"-"$arch".tar.gz ossmark && rm -rf ossmark
}

case "$1" in
"all")
    build_os_arch "linux" "amd64"
    build_os_arch "linux" "arm64"
    build_os_arch "darwin" "amd64"
    build_os_arch "darwin" "arm64"
    build_os_arch "windows" "amd64"
    ;;
*)
    build_os_arch "linux" "amd64"
    ;;
esac
