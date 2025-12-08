#!/usr/bin/env bash
set -e

echo "Building heic2jpeg..."

GOOS=linux   GOARCH=amd64 go build -o build/heic2jpeg-linux-amd64   ./cmd/heic2jpeg
GOOS=darwin  GOARCH=arm64 go build -o build/heic2jpeg-macos-arm64    ./cmd/heic2jpeg
GOOS=windows GOARCH=amd64 go build -o build/heic2jpeg-windows-amd64.exe ./cmd/heic2jpeg

echo "Done. Binaries are in ./build"
