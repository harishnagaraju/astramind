#!/usr/bin/env bash

set -e

echo
echo "[Build] Formatting..."
go fmt ./...

echo
echo "[Build] Static Analysis..."
go vet ./...

echo
echo "[Build] Building..."
go build -o astramind.exe ./cmd/astramind

echo
echo "[Build] SUCCESS"