@echo off

echo ================================
echo AstraMind Coverage Report
echo ================================

if not exist tests\output\coverage mkdir tests\output\coverage

go test -coverprofile=tests\output\coverage\coverage.out ./...

go tool cover ^
-html=tests\output\coverage\coverage.out ^
-o tests\output\coverage\coverage.html

go tool cover ^
-func=tests\output\coverage\coverage.out ^
> tests\output\coverage\coverage.txt

type tests\output\coverage\coverage.txt

echo.
echo Coverage report generated.