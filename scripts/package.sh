#!/usr/bin/env bash
set -e

mkdir -p dist
rm -rf dist/*

for file in build/*; do
    name=$(basename "$file")
    zip "dist/$name.zip" "$file"
done

shasum -a 256 dist/* > dist/checksums.txt
