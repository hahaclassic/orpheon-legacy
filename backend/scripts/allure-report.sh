#!/usr/bin/env bash
set -euo pipefail

# Папки
RESULTS=report/allure-results
REPORT=report/allure-report
HISTORY_STORE=report

mkdir -p "$RESULTS" "$REPORT" "$HISTORY_STORE"

if [ -d "$HISTORY_STORE/history" ]; then
  rm -rf "$RESULTS/history"
  cp -r "$HISTORY_STORE/history" "$RESULTS/history"
fi

go clean -testcache

go test ./internal/domain/... ./internal/repository/... -json \
  | golurectl -l -e -s -a -o "$RESULTS" --allure-suite Unit --allure-tags UNIT || true

go test ./tests/integration/... -json \
  | golurectl -l -e -s -a -o "$RESULTS" --allure-suite Integration --allure-tags INTEGRATION || true

go test ./tests/e2e/... -json \
  | golurectl -l -e -s -a -o "$RESULTS" --allure-suite E2E --allure-tags E2E || true

allure generate "$RESULTS" -o "$REPORT" --clean --name "Run $(date +%s)"

rm -rf "$HISTORY_STORE/history"
cp -r "$REPORT/history" "$HISTORY_STORE"

allure open "$REPORT"
