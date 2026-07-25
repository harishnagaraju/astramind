#!/usr/bin/env bash

set -e

echo "Running Knowledge Base Integration Tests..."
echo

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
cd "$REPO_ROOT"

if [ -f "./astramind.exe" ]; then
    BIN="./astramind.exe"
elif [ -f "./astramind" ]; then
    BIN="./astramind"
else
    echo "Could not find astramind.exe or astramind in $REPO_ROOT."
    echo "Build first: go build -o astramind.exe ./cmd/astramind"
    exit 1
fi

"$BIN" --script tests/integration/commands/kb.txt