@echo off

echo ==========================================
echo AstraMind Integration Test Suite
echo ==========================================
echo.

echo [1/5] Formatting...
go fmt ./...
if errorlevel 1 goto :failed

echo.
echo [2/5] Vet...
go vet ./...
if errorlevel 1 goto :failed

echo.
echo [3/5] Build...
go build -o astramind.exe ./cmd/astramind
if errorlevel 1 goto :failed

echo.
echo [4/5] Unit Tests...
go test ./...
if errorlevel 1 goto :failed

echo.
echo [5/5] Knowledge Base Integration...
call tests\integration\run_kb.bat
if errorlevel 1 goto :failed

echo.
echo ==========================================
echo ALL TESTS PASSED
echo ==========================================
exit /b 0

:failed
echo.
echo ==========================================
echo TESTS FAILED
echo ==========================================
exit /b 1