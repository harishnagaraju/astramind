@echo off

echo.
echo [Test] Running Unit ^& Integration Tests...

go test -v ./...

echo.
echo [Test] SUCCESS