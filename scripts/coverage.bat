@echo off

echo ================================
echo AstraMind Coverage Report
echo ================================

if not exist tests\coverage mkdir tests\coverage

go test -coverprofile=tests\coverage\coverage.out ./...

go tool cover ^
-html=tests\coverage\coverage.out ^
-o tests\coverage\coverage.html

go tool cover ^
-func=tests\coverage\coverage.out ^
> tests\coverage\coverage.txt

type tests\coverage\coverage.txt

echo.
echo Coverage report generated.