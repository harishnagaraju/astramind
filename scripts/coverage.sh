#!/bin/bash

set -e

echo "===================================="
echo "AstraMind Coverage Report"
echo "===================================="

mkdir -p tests/coverage

rm -f tests/coverage/coverage.out
rm -f tests/coverage/coverage.txt
rm -f tests/coverage/coverage.html
rm -f tests/coverage/package_coverage.txt

go test \
    -coverprofile=tests/coverage/coverage.out \
    ./... \
    > tests/coverage/package_coverage.txt

echo
echo "Test Coverage Summary"

grep "coverage:" tests/coverage/package_coverage.txt \
| grep -v "cmd/astramind" \
| sed 's/^ok[[:space:]]*//' \
| sed 's/[[:space:]]*[0-9.]\+s//' \
| sed 's/(cached)//g'

go tool cover \
-html=tests/coverage/coverage.out \
-o tests/coverage/coverage.html

go tool cover \
    -func=tests/coverage/coverage.out \
    > tests/coverage/coverage.txt

TOTAL=$(grep "^total:" tests/coverage/coverage.txt | awk '{print $3}')

echo
echo "Overall Coverage : $TOTAL"

echo
echo "Coverage report generated."