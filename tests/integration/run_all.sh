#!/usr/bin/env bash

set -e

echo "=========================================="
echo "AstraMind Integration Test Suite"
echo "=========================================="
echo

echo "[1/5] Formatting..."
go fmt ./...

echo
echo "[2/5] Vet..."
go vet ./...

echo
echo "[3/5] Build..."
go build -o astramind ./cmd/astramind

echo
echo "[4/5] Unit Tests..."
go test ./...

echo
echo "[5/5] Knowledge Base Integration..."
bash tests/integration/run_kb.sh

echo
echo "=========================================="
echo "ALL TESTS PASSED"
echo "=========================================="