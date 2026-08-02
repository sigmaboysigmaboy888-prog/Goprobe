#!/usr/bin/env bash
# build.sh – compile goprobe for the current platform, or cross-compile for common targets.
set -euo pipefail

VERSION="1.0.0"
LDFLAGS="-s -w -X main.version=${VERSION}"

echo ">> downloading dependencies..."
go mod download

if [[ "${1:-}" == "all" ]]; then
    echo ">> cross-compiling for all targets..."
    mkdir -p dist

    GOOS=linux   GOARCH=amd64  go build -ldflags "${LDFLAGS}" -o dist/goprobe-linux-amd64  .
    GOOS=linux   GOARCH=arm64  go build -ldflags "${LDFLAGS}" -o dist/goprobe-linux-arm64  .
    GOOS=darwin  GOARCH=amd64  go build -ldflags "${LDFLAGS}" -o dist/goprobe-darwin-amd64 .
    GOOS=darwin  GOARCH=arm64  go build -ldflags "${LDFLAGS}" -o dist/goprobe-darwin-arm64 .
    GOOS=windows GOARCH=amd64  go build -ldflags "${LDFLAGS}" -o dist/goprobe-windows-amd64.exe .

    echo ">> binaries written to ./dist/"
    ls -lh dist/
else
    echo ">> building for current platform..."
    go build -ldflags "${LDFLAGS}" -o goprobe .
    echo ">> built ./goprobe"
    echo ""
    echo "   usage: ./goprobe --help"
fi
