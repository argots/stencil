#!/bin/bash

set -ex

unset GOPATH
ROOT=$(dirname "$0")/..

cd $ROOT
if [[ ! -f "./bin/golangci-lint" ]]
then
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.24.0
fi

./bin/golangci-lint run --enable-all -D wsl ./...
