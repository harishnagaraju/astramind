#!/usr/bin/env bash
#
# Usage:
#   bash scripts/regression.sh          # fast, no web server involved
#   bash scripts/regression.sh --web    # also runs the web UI smoke test
#
# Deliberately does NOT use `set -e`. Each step's real exit code is
# captured into its own status variable and the summary at the end
# reports what actually happened - PASS, FAIL, or SKIPPED (if an
# earlier step failed and this one never ran). An earlier version
# printed a hardcoded "PASS" for every step unconditionally,
# regardless of whether it actually passed - technically harmless
# only because `set -e` guaranteed an early exit before that summary
# could ever be reached on a real failure, but genuinely dishonest
# reporting: the text never reflected a real, checked outcome.

RUN_WEB=""
for arg in "$@"; do
    if [ "$arg" = "--web" ]; then
        RUN_WEB="--web"
    fi
done

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
OSTAG=$(uname -s | tr '[:upper:]' '[:lower:]')
LOGFILE="regression_${OSTAG}_${TIMESTAMP}.log"

exec > >(tee "$LOGFILE") 2>&1

START_TIME=$(date +%s)

echo "===================================="
echo " AstraMind Regression Test Suite"
echo "===================================="
echo "Platform : $(uname -s) (bash / regression.sh)"
echo "Log: $LOGFILE"

BUILD_STATUS="SKIPPED"
TEST_STATUS="SKIPPED"
COVERAGE_STATUS="SKIPPED"
KBRAG_STATUS="SKIPPED"

echo
echo "[1/4] Build..."
./scripts/build.sh
if [ $? -eq 0 ]; then BUILD_STATUS="PASS"; else BUILD_STATUS="FAIL"; fi

if [ "$BUILD_STATUS" = "PASS" ]; then
    echo
    echo "[2/4] Tests..."
    ./scripts/test.sh
    if [ $? -eq 0 ]; then TEST_STATUS="PASS"; else TEST_STATUS="FAIL"; fi
fi

if [ "$TEST_STATUS" = "PASS" ]; then
    echo
    echo "[3/4] Coverage..."
    ./scripts/coverage.sh
    if [ $? -eq 0 ]; then COVERAGE_STATUS="PASS"; else COVERAGE_STATUS="FAIL"; fi
fi

if [ "$COVERAGE_STATUS" = "PASS" ]; then
    echo
    echo "[4/4] Knowledge Base & RAG Regression..."
    ./scripts/check_knowledge_base.sh
    KB_EXIT=$?
    bash ./scripts/check_rag_behavior.sh $RUN_WEB
    RAG_EXIT=$?
    if [ $KB_EXIT -eq 0 ] && [ $RAG_EXIT -eq 0 ]; then
        KBRAG_STATUS="PASS"
    else
        KBRAG_STATUS="FAIL"
    fi
fi

END_TIME=$(date +%s)
ELAPSED=$((END_TIME-START_TIME))

echo
echo "===================================="
echo " Regression Summary"
echo "===================================="
echo "Build      : $BUILD_STATUS"
echo "Tests      : $TEST_STATUS"
echo "Coverage   : $COVERAGE_STATUS"
echo "KB & RAG   : $KBRAG_STATUS"
echo "Elapsed    : ${ELAPSED} sec"
echo "Log saved  : ${LOGFILE}"
echo "===================================="

if [ "$BUILD_STATUS" = "PASS" ] && [ "$TEST_STATUS" = "PASS" ] && [ "$COVERAGE_STATUS" = "PASS" ] && [ "$KBRAG_STATUS" = "PASS" ]; then
    echo " AstraMind is READY"
    echo "===================================="
    exit 0
else
    echo " AstraMind is NOT READY - see above"
    echo "===================================="
    exit 1
fi