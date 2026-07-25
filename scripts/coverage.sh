#!/bin/bash

set -e

echo "===================================="
echo "AstraMind Coverage Report"
echo "===================================="

mkdir -p tests/output/coverage

rm -f tests/output/coverage/coverage.out
rm -f tests/output/coverage/coverage.txt
rm -f tests/output/coverage/coverage.html
rm -f tests/output/coverage/package_coverage.txt

go test \
    -coverprofile=tests/output/coverage/coverage.out \
    ./... \
    > tests/output/coverage/package_coverage.txt

echo
echo "Test Coverage Summary"

grep "coverage:" tests/output/coverage/package_coverage.txt \
| grep -v "cmd/astramind" \
| sed 's/^ok[[:space:]]*//' \
| sed 's/[[:space:]]*[0-9.]\+s//' \
| sed 's/(cached)//g'

go tool cover \
-html=tests/output/coverage/coverage.out \
-o tests/output/coverage/coverage.html

go tool cover \
    -func=tests/output/coverage/coverage.out \
    > tests/output/coverage/coverage.txt

TOTAL=$(grep "^total:" tests/output/coverage/coverage.txt | awk '{print $3}')

echo
echo "Overall Coverage : $TOTAL"

echo
echo "Coverage report generated."