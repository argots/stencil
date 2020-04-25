#!/bin/bash

function build() {
    export GOARCH=$1
    export GOOS=$2
    
    mkdir -p ./bin/"$2"_"$1"
    GOARCH=$1 GOOS=$2 go build -o ./bin/"$2"_"$1"/. ./...

    pushd ./bin/"$2"_"$1"
    tar cvf ../"$2"_"$1".zip *
    popd
}

set -ex
build 386 linux
build amd64 linux
build 386 darwin
build amd64 darwin
build 386 windows
build amd64 windows

