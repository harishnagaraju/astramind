@echo off
setlocal enabledelayedexpansion

rem check_rag_behavior.bat
rem
rem Runs the manual verification checklist for the three fixes made
rem just before the v0.9.2 release (plural "what is the" routing,
rem relevance gate for out-of-scope questions, and the short-word
rem substring collision found while fixing the relevance gate), plus
rem regression checks confirming the already-shipped enumeration and
rem single-fact precision behavior still works correctly.
rem
rem Usage:
rem   scripts\check_rag_behavior.bat          - CLI checks only
rem   scripts\check_rag_behavior.bat --web     - also runs the web UI smoke test
rem
rem Requires curl.exe (bundled with Windows 10 1803+ and Windows 11 by
rem default) for the --web section only.
rem
rem Log capture: batch has no direct equivalent of a live "tee" (show
rem on screen AND write to file at the same time). Every line below
rem that should end up in the saved log is therefore written twice -
rem once as a normal "echo" (screen) and once appended to %LOGFILE%
rem with ">>" (file). This fixes a real, confirmed bug: an earlier
rem version only captured the astramind.exe invocation into the log
rem file, silently dropping the "Expected results" checklist and the
rem "Web UI Smoke Test" section from the saved file even though both
rem appeared on screen - confirmed by directly diffing a saved .sh
rem run's log against its own terminal output.

cd /d "%~dp0.."

set RUN_WEB=0
if "%1"=="--web" set RUN_WEB=1

if not exist astramind.exe (
    echo Could not find astramind.exe in %CD%.
    echo Build first: go build -o astramind.exe ./cmd/astramind
    exit /b 1
)

for /f "delims=" %%T in ('powershell -NoProfile -Command "Get-Date -Format yyyyMMdd_HHmmss"') do set TIMESTAMP=%%T
set LOGFILE=check_rag_behavior_win_%TIMESTAMP%.log
set CMDFILE=%TEMP%\rag_cmdfile_%RANDOM%.txt

rem Batch has no heredoc - build the command file line by line,
rem matching the exact question set the .sh version uses.
> "%CMDFILE%" (
    echo /kb clear
    echo /kb import tests/fixtures/kb_documents/Sanskrit1.txt
    echo /kb ask what is the timings of sanskrit classes?
    echo /kb ask what are the sanskrit class timings
    echo /kb ask what is the zoom meeting id
    echo /kb ask what is the capital of delhi
    echo /kb clear
    echo /kb import tests/fixtures/kb_documents/sample_contract.docx
    echo /kb ask what is the zoom meeting id
    echo /kb ask how much should client pay consultant
    echo /kb import tests/fixtures/kb_documents/Sanskrit1.txt
    echo /kb ask what is the timings of sanskrit classes?
    echo /kb ask how much should client pay consultant
    echo /kb ask what is the capital of delhi
    echo /kb clear
    echo exit
)

echo ========================================== > "%LOGFILE%"
echo v0.9.2 Regression Check >> "%LOGFILE%"
echo Platform : Windows ^(cmd.exe / check_rag_behavior.bat^) >> "%LOGFILE%"
echo Binary : astramind.exe >> "%LOGFILE%"
echo Log    : %LOGFILE% >> "%LOGFILE%"
echo ========================================== >> "%LOGFILE%"
echo. >> "%LOGFILE%"

astramind.exe --script "%CMDFILE%" >> "%LOGFILE%" 2>&1
type "%LOGFILE%"

del "%CMDFILE%" >nul 2>&1

echo.
echo.>>"%LOGFILE%"
echo ==========================================
echo ==========================================>>"%LOGFILE%"
echo Done. Full transcript: %LOGFILE%
echo Done. Full transcript: %LOGFILE%>>"%LOGFILE%"
echo ==========================================
echo ==========================================>>"%LOGFILE%"
echo.
echo.>>"%LOGFILE%"
echo Expected results, check the log above against these:
echo Expected results, check the log above against these:>>"%LOGFILE%"
echo   Test 1: all 9 entries, 3 sources ^(was 4 entries/1 source before the fix^)
echo   Test 1: all 9 entries, 3 sources ^(was 4 entries/1 source before the fix^)>>"%LOGFILE%"
echo   Test 2: all 9 entries, 3 sources
echo   Test 2: all 9 entries, 3 sources>>"%LOGFILE%"
echo   Test 3: Meeting ID / Zoom section only, 1 source
echo   Test 3: Meeting ID / Zoom section only, 1 source>>"%LOGFILE%"
echo   Test 4: 'No relevant knowledge found to answer this question.'
echo   Test 4: 'No relevant knowledge found to answer this question.'>>"%LOGFILE%"
echo   Test 5: 'No relevant knowledge found to answer this question.' ^(was returning the contract before the fix^)
echo   Test 5: 'No relevant knowledge found to answer this question.' ^(was returning the contract before the fix^)>>"%LOGFILE%"
echo   Test 6: exact $45,000 paragraph, 1 source
echo   Test 6: exact $45,000 paragraph, 1 source>>"%LOGFILE%"
echo   Test 7: each question correctly answered from its own document; capital-of-delhi still rejected
echo   Test 7: each question correctly answered from its own document; capital-of-delhi still rejected>>"%LOGFILE%"

rem ==========================================
rem PART 2: Web UI (/api/ask) smoke test
rem ==========================================
rem
rem Skipped by default - starting a background --web server is the
rem kind of thing that should never silently hang an automated
rem pre-release gate. Pass --web to include it.
if "%RUN_WEB%"=="0" (
    echo.
    echo.>>"%LOGFILE%"
    echo ==========================================
    echo ==========================================>>"%LOGFILE%"
    echo Web UI ^(/api/ask^) Smoke Test - SKIPPED
    echo Web UI ^(/api/ask^) Smoke Test - SKIPPED>>"%LOGFILE%"
    echo ==========================================
    echo ==========================================>>"%LOGFILE%"
    echo Pass --web to include this section, e.g.:
    echo Pass --web to include this section, e.g.:>>"%LOGFILE%"
    echo   scripts\check_rag_behavior.bat --web
    echo   scripts\check_rag_behavior.bat --web>>"%LOGFILE%"
    echo.
    echo.>>"%LOGFILE%"
    exit /b 0
)

echo.
echo.>>"%LOGFILE%"
echo ==========================================
echo ==========================================>>"%LOGFILE%"
echo Web UI ^(/api/ask^) Smoke Test
echo Web UI ^(/api/ask^) Smoke Test>>"%LOGFILE%"
echo ==========================================
echo ==========================================>>"%LOGFILE%"
echo.
echo.>>"%LOGFILE%"

set WEB_ADDR=localhost:8420
set WEB_LOG=check_rag_behavior_web_win_%TIMESTAMP%.log

start "" /b astramind.exe --web > "%WEB_LOG%" 2>&1

echo Waiting for server to start...
echo Waiting for server to start...>>"%LOGFILE%"
timeout /t 2 /nobreak >nul

echo Clearing knowledge base via CLI first, for a clean starting state...
echo Clearing knowledge base via CLI first, for a clean starting state...>>"%LOGFILE%"
set CLEAR_CMDFILE=%TEMP%\rag_clear_%RANDOM%.txt
> "%CLEAR_CMDFILE%" (
    echo /kb clear
    echo exit
)
astramind.exe --script "%CLEAR_CMDFILE%" >nul 2>&1
del "%CLEAR_CMDFILE%" >nul 2>&1

echo.
echo.>>"%LOGFILE%"
echo Importing Sanskrit1.txt ^(tests/fixtures/kb_documents^) via /api/documents...
echo Importing Sanskrit1.txt ^(tests/fixtures/kb_documents^) via /api/documents...>>"%LOGFILE%"
curl -s -o nul -F "file=@tests/fixtures/kb_documents/Sanskrit1.txt" "http://%WEB_ADDR%/api/documents"

echo.
echo.>>"%LOGFILE%"
echo --- Web Test 1: plural 'what is the' routing ^(same as CLI Test 1^) ---
echo --- Web Test 1: plural 'what is the' routing ^(same as CLI Test 1^) --->>"%LOGFILE%"
curl -s -X POST -H "Content-Type: application/json" -d "{\"question\":\"what is the timings of sanskrit classes?\"}" "http://%WEB_ADDR%/api/ask" > "%TEMP%\web1.json"
type "%TEMP%\web1.json"
type "%TEMP%\web1.json" >> "%LOGFILE%"
echo.
echo.>>"%LOGFILE%"

echo --- Web Test 2: single-fact precision ^(same as CLI Test 3^) ---
echo --- Web Test 2: single-fact precision ^(same as CLI Test 3^) --->>"%LOGFILE%"
curl -s -X POST -H "Content-Type: application/json" -d "{\"question\":\"what is the zoom meeting id\"}" "http://%WEB_ADDR%/api/ask" > "%TEMP%\web2.json"
type "%TEMP%\web2.json"
type "%TEMP%\web2.json" >> "%LOGFILE%"
echo.
echo.>>"%LOGFILE%"

echo --- Web Test 3: out-of-scope question ^(same as CLI Test 4^) ---
echo --- Web Test 3: out-of-scope question ^(same as CLI Test 4^) --->>"%LOGFILE%"
curl -s -X POST -H "Content-Type: application/json" -d "{\"question\":\"what is the capital of delhi\"}" "http://%WEB_ADDR%/api/ask" > "%TEMP%\web3.json"
type "%TEMP%\web3.json"
type "%TEMP%\web3.json" >> "%LOGFILE%"
echo.
echo.>>"%LOGFILE%"

echo Web sanity scan:
echo Web sanity scan:>>"%LOGFILE%"

findstr /c:"Sanskrit Term 14" "%TEMP%\web1.json" >nul && (echo   [OK]      Web Test 1: enumeration includes a real entry - Monday) || (echo   [MISSING] Web Test 1: enumeration includes a real entry - Monday)
findstr /c:"Sanskrit Term 14" "%TEMP%\web1.json" >nul && (echo   [OK]      Web Test 1: enumeration includes a real entry - Monday>>"%LOGFILE%") || (echo   [MISSING] Web Test 1: enumeration includes a real entry - Monday>>"%LOGFILE%")
findstr /c:"Thursday Senior Sanskrit" "%TEMP%\web1.json" >nul && (echo   [OK]      Web Test 1: enumeration includes a real entry - Thursday - the exact fix) || (echo   [MISSING] Web Test 1: enumeration includes a real entry - Thursday - the exact fix)
findstr /c:"Thursday Senior Sanskrit" "%TEMP%\web1.json" >nul && (echo   [OK]      Web Test 1: enumeration includes a real entry - Thursday - the exact fix>>"%LOGFILE%") || (echo   [MISSING] Web Test 1: enumeration includes a real entry - Thursday - the exact fix>>"%LOGFILE%")
findstr /c:"Meeting ID 795 777 3585" "%TEMP%\web2.json" >nul && (echo   [OK]      Web Test 2: precise single-fact match present) || (echo   [MISSING] Web Test 2: precise single-fact match present)
findstr /c:"Meeting ID 795 777 3585" "%TEMP%\web2.json" >nul && (echo   [OK]      Web Test 2: precise single-fact match present>>"%LOGFILE%") || (echo   [MISSING] Web Test 2: precise single-fact match present>>"%LOGFILE%")
findstr /c:"No relevant knowledge found" "%TEMP%\web3.json" >nul && (echo   [OK]      Web Test 3: out-of-scope question correctly rejected) || (echo   [MISSING] Web Test 3: out-of-scope question correctly rejected)
findstr /c:"No relevant knowledge found" "%TEMP%\web3.json" >nul && (echo   [OK]      Web Test 3: out-of-scope question correctly rejected>>"%LOGFILE%") || (echo   [MISSING] Web Test 3: out-of-scope question correctly rejected>>"%LOGFILE%")

del "%TEMP%\web1.json" "%TEMP%\web2.json" "%TEMP%\web3.json" >nul 2>&1

echo.
echo.>>"%LOGFILE%"
echo Cleaning up test document from the knowledge base...
echo Cleaning up test document from the knowledge base...>>"%LOGFILE%"
taskkill /IM astramind.exe /F >nul 2>&1
timeout /t 1 /nobreak >nul

set CLEANUP_CMDFILE=%TEMP%\rag_cleanup_%RANDOM%.txt
> "%CLEANUP_CMDFILE%" (
    echo /kb clear
    echo exit
)
astramind.exe --script "%CLEANUP_CMDFILE%" >nul 2>&1
del "%CLEANUP_CMDFILE%" >nul 2>&1

echo.
echo.>>"%LOGFILE%"
echo ==========================================
echo ==========================================>>"%LOGFILE%"
echo All done. Logs:
echo All done. Logs:>>"%LOGFILE%"
echo   CLI transcript : %LOGFILE%
echo   CLI transcript : %LOGFILE%>>"%LOGFILE%"
echo   Web server log : %WEB_LOG%
echo   Web server log : %WEB_LOG%>>"%LOGFILE%"
echo ==========================================
echo ==========================================>>"%LOGFILE%"