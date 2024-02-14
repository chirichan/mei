#!/bin/bash

release_dir="release"

if [ ! -d "$folder" ]; then
    mkdir -p "$folder" > /dev/null 2>&1
fi

rm -f ./release/pwdgen-darwin-arm64
rm -f ./release/pwdgen-darwin-x86_64
rm -f ./release/pwdgen-linux-x86_64
rm -f ./release/pwdgen-windows-x86_64.exe

rm -f ./release/bing15-darwin-arm64
rm -f ./release/bing15-darwin-x86_64
rm -f ./release/bing15-linux-x86_64
rm -f ./release/bing15-windows-x86_64.exe

CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ./release/pwdgen-darwin-arm64 ./cmd/pwdgen/
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./release/pwdgen-darwin-x86_64 ./cmd/pwdgen/
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./release/pwdgen-linux-x86_64 ./cmd/pwdgen/
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./release/pwdgen-windows-x86_64.exe ./cmd/pwdgen/

CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ./release/bing15-darwin-arm64 ./cmd/bing15/
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./release/bing15-darwin-x86_64 ./cmd/bing15/
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./release/bing15-linux-x86_64 ./cmd/bing15/
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./release/bing15-windows-x86_64.exe ./cmd/bing15/

echo "build success"