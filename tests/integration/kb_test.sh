#!/usr/bin/env bash

set -e

echo "=========================================="
echo "Knowledge Base Integration Test"
echo "=========================================="
echo

./astramind --script tests/integration/commands/kb.txt

echo
echo "=========================================="
echo "Knowledge Base Integration Passed"
echo "=========================================="