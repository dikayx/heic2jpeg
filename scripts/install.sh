#!/usr/bin/env bash
set -e

go build -o /usr/local/bin/heic2jpeg ./cmd/heic2jpeg
echo "Installed to /usr/local/bin/heic2jpeg"
