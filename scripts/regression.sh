#!/usr/bin/env bash

set -e

START_TIME=$(date +%s)

echo "===================================="
echo " AstraMind Regression Test Suite"
echo "===================================="

echo
echo "[1/3] Build..."
./scripts/build.sh

echo
echo "[2/3] Tests..."
./scripts/test.sh

echo
echo "[3/3] Coverage..."
./scripts/coverage.sh

END_TIME=$(date +%s)
ELAPSED=$((END_TIME-START_TIME))

echo
echo "===================================="
echo " Regression Summary"
echo "===================================="
echo "Build      : PASS"
echo "Tests      : PASS"
echo "Coverage   : PASS"
echo "Elapsed    : ${ELAPSED} sec"
echo "===================================="
echo " AstraMind is READY"
echo "===================================="