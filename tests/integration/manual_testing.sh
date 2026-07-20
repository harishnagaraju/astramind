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

# Topically distinct content (not near-duplicates) so semantic search
# has something meaningful to distinguish, and so a paraphrase query
# has zero substring overlap with either file - proving the gap
# between keyword and semantic search rather than just smoke-testing
# that both commands run without error.
echo "The Eiffel Tower is a wrought-iron lattice tower located in Paris, France." > "$KB_FILE_1"
echo "Photosynthesis allows plants to convert sunlight into chemical energy." > "$KB_FILE_2"

# A paraphrase of file 2's content, sharing no contiguous substring
# with either file. /kb search does a plain substring match on the
# full phrase, so this should find nothing. /kb ssearch should still
# surface manualtest_kb2.md by meaning.
SEMANTIC_QUERY="how do plants get energy from the sun"

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
/kb search Paris
/kb import $KB_FILE_2
/kb list
/kb stats
/kb search $SEMANTIC_QUERY
/kb ssearch $SEMANTIC_QUERY
/kb ask $SEMANTIC_QUERY
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
check "chunks embedded"                     "/kb import generates embeddings"
check "Semantic Search Results"             "/kb ssearch hit"
check "Sources:"                            "/kb ask completed the RAG loop"
check "Knowledge base cleared"              "/kb clear"
check "Knowledge base is empty"             "/kb list after clear"
check "Created and switched to session: ${TMP_SESSION}" "/new session"
check "Loaded session: default"             "/load default"
check "Deleted session: ${TMP_SESSION}"     "/delete session"
check "Current AI Provider"                 "/provider output"
check "Goodbye!"                            "clean exit"

echo
echo "Manual comparison (not auto-checked - eyeball the log above):"
echo "  /kb search \"$SEMANTIC_QUERY\"  -> should say 'No matching knowledge found.'"
echo "  /kb ssearch \"$SEMANTIC_QUERY\" -> should surface manualtest_kb2.md by meaning"
echo "  /kb ask \"$SEMANTIC_QUERY\"     -> should give an actual answer, citing manualtest_kb2.md as the source"
echo "  If both search modes found it, or neither did, the semantic path isn't adding real value yet."

# ==========================================
# PART 2: Web UI API smoke test (--web mode)
# ==========================================
#
# This drives the same backend through the local web server's JSON
# API instead of stdin - the same code path the browser UI uses.
# Also confirms --web launches correctly and stays up.
echo
echo "=========================================="
echo "Web UI API Smoke Test (--web mode)"
echo "=========================================="
echo

WEB_ADDR="localhost:8420"
WEB_FILE_1="manualtest_web1.md"
WEB_FILE_2="manualtest_web2.md"
WEB_LOG="manual_test_web_${TIMESTAMP}.log"

echo "The Eiffel Tower is a wrought-iron lattice tower located in Paris, France." > "$WEB_FILE_1"
echo "Photosynthesis allows plants to convert sunlight into chemical energy." > "$WEB_FILE_2"

"$BIN" --web > "$WEB_LOG" 2>&1 &
WEB_PID=$!

# Always clean up the background server and fixture files, even if a
# later step in this section fails.
trap 'kill "$WEB_PID" 2>/dev/null || true; rm -f "$WEB_FILE_1" "$WEB_FILE_2"' EXIT

echo "Waiting for server to start (pid $WEB_PID)..."
sleep 2

STATUS_RESPONSE=$(curl -s "http://${WEB_ADDR}/api/status" || echo "CURL_FAILED")
echo "GET /api/status -> $STATUS_RESPONSE"

curl -s -o /dev/null -w "" -F "file=@${WEB_FILE_1}" "http://${WEB_ADDR}/api/documents" || true
curl -s -o /dev/null -w "" -F "file=@${WEB_FILE_2}" "http://${WEB_ADDR}/api/documents" || true

DOCS_RESPONSE=$(curl -s "http://${WEB_ADDR}/api/documents" || echo "CURL_FAILED")
echo "GET /api/documents -> $DOCS_RESPONSE"

# Extract manualtest_web2.md's document ID from the /api/documents
# response, so the /api/ask citation check compares against the
# actual ID the API returns - not the filename, which never appears
# in /api/ask's response at all.
PHOTOSYNTHESIS_DOC_ID=$(echo "$DOCS_RESPONSE" | grep -o '"id":"[^"]*","name":"manualtest_web2.md"' | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

ASK_RESPONSE=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d "{\"question\":\"$SEMANTIC_QUERY\"}" \
    "http://${WEB_ADDR}/api/ask" || echo "CURL_FAILED")
echo "POST /api/ask -> $ASK_RESPONSE"

echo
echo "Web API sanity scan:"

web_check() {
    if echo "$2" | grep -q "$1"; then
        echo "  [OK]      $3"
    else
        echo "  [MISSING] $3"
    fi
}

web_check "provider"                "$STATUS_RESPONSE" "/api/status returned provider info"
web_check "manualtest_web1.md"      "$DOCS_RESPONSE"    "/api/documents lists imported file 1"
web_check "manualtest_web2.md"      "$DOCS_RESPONSE"    "/api/documents lists imported file 2"
web_check "\"sources\""             "$ASK_RESPONSE"     "/api/ask returned cited sources"
web_check "$PHOTOSYNTHESIS_DOC_ID"  "$ASK_RESPONSE"     "/api/ask correctly cited the photosynthesis doc"

kill "$WEB_PID" 2>/dev/null || true
trap - EXIT
rm -f "$WEB_FILE_1" "$WEB_FILE_2"

# The web smoke test imports real documents into the actual
# data/ knowledge base (the same one used interactively and via
# the browser) - clean them up now, or they silently pollute every
# real /kb ask afterward, the same class of bug fixed earlier for
# the search tests writing into the real data/sessions folder.
echo
echo "Cleaning up test documents from the knowledge base..."
CLEANUP_SCRIPT=$(mktemp)
echo "/kb clear" > "$CLEANUP_SCRIPT"
echo "exit" >> "$CLEANUP_SCRIPT"
"$BIN" --script "$CLEANUP_SCRIPT" > /dev/null 2>&1 || true
rm -f "$CLEANUP_SCRIPT"

echo
echo "=========================================="
echo "All done. Logs:"
echo "  Interactive walkthrough : $LOGFILE"
echo "  Web server stdout       : $WEB_LOG"
echo "=========================================="