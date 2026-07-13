@echo off

echo Running Knowledge Base Integration Tests...
echo.

.\astramind.exe --script tests\integration\commands\kb.txt

exit /b %ERRORLEVEL%

