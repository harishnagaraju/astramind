#!/usr/bin/env bash

set -e

echo
echo "[Test] Running Unit & Integration Tests..."

go test -v ./...

echo
echo "[Test] SUCCESS"