#!/bin/bash

echo "===================================="
echo "AstraMind Coverage Report"
echo "===================================="

mkdir -p tests/coverage

go test -coverprofile=tests/coverage/coverage.out ./...

go tool cover \
-html=tests/coverage/coverage.out \
-o tests/coverage/coverage.html

go tool cover \
-func=tests/coverage/coverage.out \
> tests/coverage/coverage.txt

cat tests/coverage/coverage.txt

echo
echo "Coverage report generated."