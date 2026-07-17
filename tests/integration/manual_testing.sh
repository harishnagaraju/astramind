#!/usr/bin/env bash
#
# manual_test.sh
#
# Drives AstraMind's interactive CLI end-to-end by piping the same
# command sequence you've been typing by hand into stdin, and saves
# the full transcript for review.
#
# Can be run from anywhere - it resolves the repo root from its own
# location (tests/integration/) and always executes the binary from
# there, matching how run_all.sh / run_kb.sh already behave.
#
# Usage:
#   bash tests/integration/manual_test.sh     # from repo root
#   ./manual_test.sh                          # from tests/integration/
#   ./manual_test.sh ./path/to/bin.exe        # explicit binary override
#
set -e

# Resolve repo root relative to this script's own location, so it works
# whether invoked from the repo root (bash tests/integration/manual_test.sh)
# or from inside tests/integration itself (./manual_test.sh).
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
cd "$REPO_ROOT"

BIN="${1:-}"

if [ -z "$BIN" ]; then
    if [ -f "./astramind.exe" ]; then
        BIN="./astramind.exe"
    elif [ -f "./astramind" ]; then
        BIN="./astramind"
    else
        echo "Could not find astramind.exe or astramind in repo root ($REPO_ROOT)."
        echo "Build first: go build -o astramind.exe ./cmd/astramind"
        exit 1
    fi
fi

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
LOGFILE="manual_test_${TIMESTAMP}.log"
TMP_SESSION="manualtest_tmp_$$"
KB_FILE_1="manualtest_kb1.md"
KB_FILE_2="manualtest_kb2.md"

echo "This is testing /kb import <filename> (file 1)" > "$KB_FILE_1"
echo "This is testing /kb import <filename> (file 2)" > "$KB_FILE_2"

echo "=========================================="
echo "AstraMind Manual Command Walkthrough"
echo "Binary : $BIN"
echo "Log    : $LOGFILE"
echo "=========================================="
echo

"$BIN" <<EOF | tee "$LOGFILE"
/help
/about
/history
/stats
/config
/export
/export md
/sessions
/search golang
/searchall golang
/kb import $KB_FILE_1
/kb list
/kb search import
/kb import $KB_FILE_2
/kb list
/kb stats
/kb clear
/kb list
/kb stats
/sessions
/new $TMP_SESSION
/sessions
/load default
/delete $TMP_SESSION
/sessions
/provider
exit
EOF

rm -f "$KB_FILE_1" "$KB_FILE_2"

echo
echo "=========================================="
echo "Walkthrough complete."
echo "Full transcript saved to: $LOGFILE"
echo "=========================================="
echo
echo "Quick sanity scan (look for these — 'MISSING' means check the log):"

check() {
    if grep -q "$1" "$LOGFILE"; then
        echo "  [OK]      $2"
    else
        echo "  [MISSING] $2"
    fi
}

check "Available Commands"                 "/help output"
check "AstraMind"                           "/about output"
check "Session Statistics"                  "/stats output"
check "Current Configuration"               "/config output"
check "Session exported to exports"         "/export (txt + md)"
check "Available Sessions"                  "/sessions output"
check "Knowledge Search Results"            "/kb search hit"
check "Knowledge base cleared"              "/kb clear"
check "Knowledge base is empty"             "/kb list after clear"
check "Created and switched to session: ${TMP_SESSION}" "/new session"
check "Loaded session: default"             "/load default"
check "Deleted session: ${TMP_SESSION}"     "/delete session"
check "Current AI Provider"                 "/provider output"
check "Goodbye!"                            "clean exit"