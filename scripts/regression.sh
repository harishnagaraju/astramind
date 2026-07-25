#!/usr/bin/env bash
#
# Usage:
#   bash scripts/regression.sh          # fast, no web server involved
#   bash scripts/regression.sh --web    # also runs the web UI smoke test
#
set -e

RUN_WEB=""
for arg in "$@"; do
    if [ "$arg" = "--web" ]; then
        RUN_WEB="--web"
    fi
done

START_TIME=$(date +%s)

echo "===================================="
echo " AstraMind Regression Test Suite"
echo "===================================="

echo
echo "[1/4] Build..."
./scripts/build.sh

echo
echo "[2/4] Tests..."
./scripts/test.sh

echo
echo "[3/4] Coverage..."
./scripts/coverage.sh

echo
echo "[4/4] Knowledge Base & RAG Regression..."
./scripts/check_knowledge_base.sh
bash ./scripts/check_rag_behavior.sh $RUN_WEB

END_TIME=$(date +%s)
ELAPSED=$((END_TIME-START_TIME))

echo
echo "===================================="
echo " Regression Summary"
echo "===================================="
echo "Build      : PASS"
echo "Tests      : PASS"
echo "Coverage   : PASS"
echo "KB & RAG   : PASS"
echo "Elapsed    : ${ELAPSED} sec"
echo "===================================="
echo " AstraMind is READY"
echo "===================================="