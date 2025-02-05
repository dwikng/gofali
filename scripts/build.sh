#!/bin/bash

MAIN_GO_PATH="./"
BIN_DIR="scripts/bin"

# Make sure the directory exists
mkdir -p $BIN_DIR

build_binary() {
    local os=$1
    local arch=$2
    local output_name=$3

    export GOOS=$os
    export GOARCH=$arch
    echo "Building for $os/$arch..."

    go build -o ./$BIN_DIR/$output_name $MAIN_GO_PATH
    if [ $? -ne 0 ]; then
        echo "Error building $os/$arch binary!"
        exit 1
    fi
}

# Windows 64-bit
build_binary "windows" "amd64" "gofali-amd64.exe"

# Linux 64-bit
build_binary "linux" "amd64" "gofali-linux-amd64"

# MacOS 64-bit
build_binary "darwin" "amd64" "gofali-darwin-amd64"

echo "All 64-bit binaries are built."
