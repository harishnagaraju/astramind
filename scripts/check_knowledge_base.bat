@echo off

echo Running Knowledge Base Integration Tests...
echo.

cd /d "%~dp0.."

if not exist astramind.exe (
    echo Could not find astramind.exe in %CD%.
    echo Build first: go build -o astramind.exe ./cmd/astramind
    exit /b 1
)

astramind.exe --script tests\fixtures\commands\kb.txt

exit /b %ERRORLEVEL%