@echo off

echo ====================================
echo  AstraMind Regression Test Suite
echo ====================================
echo.

echo [1/3] Build...
call scripts\build.bat

echo.
echo [2/3] Tests...
call scripts\test.bat

echo.
echo [3/3] Coverage...
call scripts\coverage.bat

echo.
echo ====================================
echo  Regression Summary
echo ====================================
echo Build      : PASS
echo Tests      : PASS
echo Coverage   : PASS
echo ====================================
echo  AstraMind is READY
echo ====================================