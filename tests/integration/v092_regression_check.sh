#!/usr/bin/env bash
#
# v092_regression_check.sh
#
# Runs the manual verification checklist for the three fixes made
# just before the v0.9.2 release (plural "what is the" routing,
# relevance gate for out-of-scope questions, and the short-word
# substring collision found while fixing the relevance gate), plus
# regression checks confirming the already-shipped enumeration and
# single-fact precision behavior still works correctly.
#
# Usage:
#   bash v092_regression_check.sh                # from repo root, CLI checks only
#   bash v092_regression_check.sh --web           # also run the web UI smoke test
#   bash v092_regression_check.sh ./astramind.exe # explicit binary path
#   bash v092_regression_check.sh --web ./astramind.exe   # both, any order
#
set -e

# Resolve repo root relative to this script's own location, matching
# the pattern already used by manual_testing.sh - works whether
# invoked directly, or from scripts/regression.sh via a relative path.
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
cd "$REPO_ROOT"

RUN_WEB=0
BIN=""

for arg in "$@"; do
    if [ "$arg" = "--web" ]; then
        RUN_WEB=1
    else
        BIN="$arg"
    fi
done

if [ -z "$BIN" ]; then
    if [ -f "./astramind.exe" ]; then
        BIN="./astramind.exe"
    elif [ -f "./astramind" ]; then
        BIN="./astramind"
    else
        echo "Could not find astramind.exe or astramind in $REPO_ROOT."
        echo "Build first: go build -o astramind.exe ./cmd/astramind"
        exit 1
    fi
fi

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
LOGFILE="test_output_${TIMESTAMP}.log"
CMDFILE=$(mktemp)

cat > "$CMDFILE" << 'EOF'
# Test 1 - Bug A fix: "what is the [plural noun]" must route to
# enumeration, not single-fact extraction. Previously this exact
# phrasing returned only 4 of 9 real entries from a single chunk.
/kb clear
/kb import Sanskrit1.txt
/kb ask what is the timings of sanskrit classes?

# Test 2 - confirm the already-working enumeration phrasing still
# returns all 9 entries (no regression from the Test 1 fix).
/kb ask what are the sanskrit class timings

# Test 3 - confirm single-fact precision still works (no regression):
# should return just the Meeting ID/Zoom section, one source, not
# the whole 3-chunk dump.
/kb ask what is the zoom meeting id

# Test 4 - Bug B fix: a genuinely out-of-scope question must say "no
# relevant knowledge found", not confidently return unrelated content.
/kb ask what is the capital of delhi

# Test 5 - the short-word substring collision found and fixed while
# building the Test 4 fix: "id" was matching inside "confidentiality"
# as a false positive, letting an unrelated contract answer a
# Zoom-related question. Import a document with zero Zoom content and
# confirm this specific false positive is gone.
/kb clear
/kb import internal/features/kb/testdata/sample_contract.docx
/kb ask what is the zoom meeting id

# Test 6 - confirm docx single-fact precision is still correct: exact
# $45,000 paragraph, one source.
/kb ask how much should client pay consultant

# Test 7 - mixed KB: both documents loaded together, confirm each
# question still routes to the correct document, and the out-of-scope
# question is still correctly rejected even with more content present.
/kb import Sanskrit1.txt
/kb ask what is the timings of sanskrit classes?
/kb ask how much should client pay consultant
/kb ask what is the capital of delhi

/kb clear
exit
EOF

echo "=========================================="
echo "v0.9.2 Regression Check"
echo "Binary : $BIN"
echo "Log    : $LOGFILE"
echo "=========================================="
echo

"$BIN" --script "$CMDFILE" 2>&1 | tee "$LOGFILE"

rm -f "$CMDFILE"

echo
echo "=========================================="
echo "Done. Full transcript: $LOGFILE"
echo "=========================================="
echo
echo "Expected results, check the log above against these:"
echo "  Test 1: all 9 entries, 3 sources (was 4 entries/1 source before the fix)"
echo "  Test 2: all 9 entries, 3 sources"
echo "  Test 3: Meeting ID / Zoom section only, 1 source"
echo "  Test 4: 'No relevant knowledge found to answer this question.'"
echo "  Test 5: 'No relevant knowledge found to answer this question.' (was returning the contract before the fix)"
echo "  Test 6: exact \$45,000 paragraph, 1 source"
echo "  Test 7: each question correctly answered from its own document; capital-of-delhi still rejected"

# ==========================================
# PART 2: Web UI (/api/ask) smoke test
# ==========================================
#
# Skipped by default: this section starts a background --web server,
# which is the kind of thing that should never silently hang an
# automated pre-release gate (confirmed firsthand - backgrounding
# --web hung a sandboxed shell during testing of this exact script).
# Pass --web to include it, e.g.:
#   bash tests/integration/v092_regression_check.sh --web
#
# scripts/regression.sh does NOT pass --web by default - the
# automated gate stays fast and non-interactive. Add --web (see that
# script's own usage) when you specifically want this deeper check.
if [ "$RUN_WEB" != "1" ]; then
    echo
    echo "=========================================="
    echo "Web UI (/api/ask) Smoke Test - SKIPPED"
    echo "=========================================="
    echo "Pass --web to include this section, e.g.:"
    echo "  bash tests/integration/v092_regression_check.sh --web"
    echo
    exit 0
fi

echo
echo "=========================================="
echo "Web UI (/api/ask) Smoke Test"
echo "=========================================="
echo

WEB_ADDR="localhost:8420"
WEB_LOG="test_output_web_${TIMESTAMP}.log"

"$BIN" --web > "$WEB_LOG" 2>&1 &
WEB_PID=$!

trap 'kill "$WEB_PID" 2>/dev/null || true' EXIT

echo "Waiting for server to start (pid $WEB_PID)..."
sleep 2

echo "Clearing knowledge base via CLI first, for a clean starting state..."
CLEAR_SCRIPT=$(mktemp)
echo "/kb clear" > "$CLEAR_SCRIPT"
echo "exit" >> "$CLEAR_SCRIPT"
"$BIN" --script "$CLEAR_SCRIPT" > /dev/null 2>&1 || true
rm -f "$CLEAR_SCRIPT"

echo
echo "Importing Sanskrit1.txt via /api/documents..."
curl -s -o /dev/null -w "" -F "file=@Sanskrit1.txt" "http://${WEB_ADDR}/api/documents" || true

echo
echo "--- Web Test 1: plural 'what is the' routing (same as CLI Test 1) ---"
ask_web() {
    curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "{\"question\":\"$1\"}" \
        "http://${WEB_ADDR}/api/ask"
}

WEB1=$(ask_web "what is the timings of sanskrit classes?")
echo "$WEB1"
echo

echo "--- Web Test 2: single-fact precision (same as CLI Test 3) ---"
WEB2=$(ask_web "what is the zoom meeting id")
echo "$WEB2"
echo

echo "--- Web Test 3: out-of-scope question (same as CLI Test 4) ---"
WEB3=$(ask_web "what is the capital of delhi")
echo "$WEB3"
echo

echo "Web sanity scan:"

web_check() {
    if echo "$2" | grep -q "$1"; then
        echo "  [OK]      $3"
    else
        echo "  [MISSING] $3"
    fi
}

web_check "Sanskrit Term 14"                                  "$WEB1" "Web Test 1: enumeration includes a real entry (Monday)"
web_check "Thursday Senior Sanskrit"                           "$WEB1" "Web Test 1: enumeration includes a real entry (Thursday) - the exact fix"
web_check "Meeting ID 795 777 3585"                            "$WEB2" "Web Test 2: precise single-fact match present"
web_check "No relevant knowledge found to answer this question" "$WEB3" "Web Test 3: out-of-scope question correctly rejected"

echo
echo "Cleaning up test document from the knowledge base..."
CLEANUP_SCRIPT=$(mktemp)
echo "/kb clear" > "$CLEANUP_SCRIPT"
echo "exit" >> "$CLEANUP_SCRIPT"
kill "$WEB_PID" 2>/dev/null || true
sleep 1
"$BIN" --script "$CLEANUP_SCRIPT" > /dev/null 2>&1 || true
rm -f "$CLEANUP_SCRIPT"
trap - EXIT

echo
echo "=========================================="
echo "All done. Logs:"
echo "  CLI transcript : $LOGFILE"
echo "  Web server log : $WEB_LOG"
echo "=========================================="