@echo off
setlocal

rem Usage:
rem   scripts\regression.bat          - fast, no web server involved
rem   scripts\regression.bat --web    - also runs the web UI smoke test

set RUN_WEB=
if "%1"=="--web" set RUN_WEB=--web

echo ====================================
echo  AstraMind Regression Test Suite
echo ====================================
echo.

echo [1/4] Build...
call scripts\build.bat
if errorlevel 1 exit /b %ERRORLEVEL%

echo.
echo [2/4] Tests...
call scripts\test.bat
if errorlevel 1 exit /b %ERRORLEVEL%

echo.
echo [3/4] Coverage...
call scripts\coverage.bat
if errorlevel 1 exit /b %ERRORLEVEL%

echo.
echo [4/4] Knowledge Base ^& RAG Regression...
call scripts\check_knowledge_base.bat
if errorlevel 1 exit /b %ERRORLEVEL%

call scripts\check_rag_behavior.bat %RUN_WEB%
if errorlevel 1 exit /b %ERRORLEVEL%

echo.
echo ====================================
echo  Regression Summary
echo ====================================
echo Build      : PASS
echo Tests      : PASS
echo Coverage   : PASS
echo KB ^& RAG   : PASS
echo ====================================
echo  AstraMind is READY
echo ====================================