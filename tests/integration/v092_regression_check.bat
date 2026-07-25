@echo off
setlocal enabledelayedexpansion

rem v092_regression_check.bat
rem
rem Runs the manual verification checklist for the three fixes made
rem just before the v0.9.2 release (plural "what is the" routing,
rem relevance gate for out-of-scope questions, and the short-word
rem substring collision found while fixing the relevance gate), plus
rem regression checks confirming the already-shipped enumeration and
rem single-fact precision behavior still works correctly.
rem
rem Usage:
rem   tests\integration\v092_regression_check.bat          - CLI checks only
rem   tests\integration\v092_regression_check.bat --web     - also runs the web UI smoke test
rem
rem Requires curl.exe (bundled with Windows 10 1803+ and Windows 11 by
rem default) for the --web section only.

cd /d "%~dp0..\.."

set RUN_WEB=0
if "%1"=="--web" set RUN_WEB=1

if not exist astramind.exe (
    echo Could not find astramind.exe in %CD%.
    echo Build first: go build -o astramind.exe ./cmd/astramind
    exit /b 1
)

set TIMESTAMP=%RANDOM%%RANDOM%
set LOGFILE=test_output_%TIMESTAMP%.log
set CMDFILE=%TEMP%\v092_cmdfile_%RANDOM%.txt

rem Batch has no heredoc - build the command file line by line,
rem matching the exact question set the .sh version uses.
> "%CMDFILE%" (
    echo /kb clear
    echo /kb import Sanskrit1.txt
    echo /kb ask what is the timings of sanskrit classes?
    echo /kb ask what are the sanskrit class timings
    echo /kb ask what is the zoom meeting id
    echo /kb ask what is the capital of delhi
    echo /kb clear
    echo /kb import internal/features/kb/testdata/sample_contract.docx
    echo /kb ask what is the zoom meeting id
    echo /kb ask how much should client pay consultant
    echo /kb import Sanskrit1.txt
    echo /kb ask what is the timings of sanskrit classes?
    echo /kb ask how much should client pay consultant
    echo /kb ask what is the capital of delhi
    echo /kb clear
    echo exit
)

echo ==========================================
echo v0.9.2 Regression Check
echo Binary : astramind.exe
echo Log    : %LOGFILE%
echo ==========================================
echo.

astramind.exe --script "%CMDFILE%" > "%LOGFILE%" 2>&1
type "%LOGFILE%"

del "%CMDFILE%" >nul 2>&1

echo.
echo ==========================================
echo Done. Full transcript: %LOGFILE%
echo ==========================================
echo.
echo Expected results, check the log above against these:
echo   Test 1: all 9 entries, 3 sources (was 4 entries/1 source before the fix)
echo   Test 2: all 9 entries, 3 sources
echo   Test 3: Meeting ID / Zoom section only, 1 source
echo   Test 4: 'No relevant knowledge found to answer this question.'
echo   Test 5: 'No relevant knowledge found to answer this question.' (was returning the contract before the fix)
echo   Test 6: exact $45,000 paragraph, 1 source
echo   Test 7: each question correctly answered from its own document; capital-of-delhi still rejected

rem ==========================================
rem PART 2: Web UI (/api/ask) smoke test
rem ==========================================
rem
rem Skipped by default - starting a background --web server is the
rem kind of thing that should never silently hang an automated
rem pre-release gate. Pass --web to include it.
if "%RUN_WEB%"=="0" (
    echo.
    echo ==========================================
    echo Web UI (/api/ask) Smoke Test - SKIPPED
    echo ==========================================
    echo Pass --web to include this section, e.g.:
    echo   tests\integration\v092_regression_check.bat --web
    echo.
    exit /b 0
)

echo.
echo ==========================================
echo Web UI (/api/ask) Smoke Test
echo ==========================================
echo.

set WEB_ADDR=localhost:8420
set WEB_LOG=test_output_web_%TIMESTAMP%.log

start "" /b astramind.exe --web > "%WEB_LOG%" 2>&1

echo Waiting for server to start...
timeout /t 2 /nobreak >nul

echo Clearing knowledge base via CLI first, for a clean starting state...
set CLEAR_CMDFILE=%TEMP%\v092_clear_%RANDOM%.txt
> "%CLEAR_CMDFILE%" (
    echo /kb clear
    echo exit
)
astramind.exe --script "%CLEAR_CMDFILE%" >nul 2>&1
del "%CLEAR_CMDFILE%" >nul 2>&1

echo.
echo Importing Sanskrit1.txt via /api/documents...
curl -s -o nul -F "file=@Sanskrit1.txt" "http://%WEB_ADDR%/api/documents"

echo.
echo --- Web Test 1: plural 'what is the' routing (same as CLI Test 1) ---
curl -s -X POST -H "Content-Type: application/json" -d "{\"question\":\"what is the timings of sanskrit classes?\"}" "http://%WEB_ADDR%/api/ask" > "%TEMP%\web1.json"
type "%TEMP%\web1.json"
echo.

echo --- Web Test 2: single-fact precision (same as CLI Test 3) ---
curl -s -X POST -H "Content-Type: application/json" -d "{\"question\":\"what is the zoom meeting id\"}" "http://%WEB_ADDR%/api/ask" > "%TEMP%\web2.json"
type "%TEMP%\web2.json"
echo.

echo --- Web Test 3: out-of-scope question (same as CLI Test 4) ---
curl -s -X POST -H "Content-Type: application/json" -d "{\"question\":\"what is the capital of delhi\"}" "http://%WEB_ADDR%/api/ask" > "%TEMP%\web3.json"
type "%TEMP%\web3.json"
echo.

echo Web sanity scan:

findstr /c:"Sanskrit Term 14" "%TEMP%\web1.json" >nul && (echo   [OK]      Web Test 1: enumeration includes a real entry - Monday) || (echo   [MISSING] Web Test 1: enumeration includes a real entry - Monday)
findstr /c:"Thursday Senior Sanskrit" "%TEMP%\web1.json" >nul && (echo   [OK]      Web Test 1: enumeration includes a real entry - Thursday - the exact fix) || (echo   [MISSING] Web Test 1: enumeration includes a real entry - Thursday - the exact fix)
findstr /c:"Meeting ID 795 777 3585" "%TEMP%\web2.json" >nul && (echo   [OK]      Web Test 2: precise single-fact match present) || (echo   [MISSING] Web Test 2: precise single-fact match present)
findstr /c:"No relevant knowledge found" "%TEMP%\web3.json" >nul && (echo   [OK]      Web Test 3: out-of-scope question correctly rejected) || (echo   [MISSING] Web Test 3: out-of-scope question correctly rejected)

del "%TEMP%\web1.json" "%TEMP%\web2.json" "%TEMP%\web3.json" >nul 2>&1

echo.
echo Cleaning up test document from the knowledge base...
taskkill /IM astramind.exe /F >nul 2>&1
timeout /t 1 /nobreak >nul

set CLEANUP_CMDFILE=%TEMP%\v092_cleanup_%RANDOM%.txt
> "%CLEANUP_CMDFILE%" (
    echo /kb clear
    echo exit
)
astramind.exe --script "%CLEANUP_CMDFILE%" >nul 2>&1
del "%CLEANUP_CMDFILE%" >nul 2>&1

echo.
echo ==========================================
echo All done. Logs:
echo   CLI transcript : %LOGFILE%
echo   Web server log : %WEB_LOG%
echo ==========================================