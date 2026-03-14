#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)"
GOCACHE_DIR="${GOCACHE:-/tmp/crona-go-cache}"
LINT_CACHE_DIR="${GOLANGCI_LINT_CACHE:-/tmp/crona-golangci-lint-cache}"
MASON_LINTER="/Users/sm2101/.local/share/nvim/mason/bin/golangci-lint"
LINTER_BIN="${GOLANGCI_LINT_BIN:-}"

if [ -z "${LINTER_BIN}" ]; then
  if [ -x "${MASON_LINTER}" ]; then
    LINTER_BIN="${MASON_LINTER}"
  elif command -v golangci-lint >/dev/null 2>&1; then
    LINTER_BIN="$(command -v golangci-lint)"
  else
    echo "golangci-lint is not installed. Run: make install-lint" >&2
    exit 1
  fi
fi

for module in shared kernel tui cli; do
  (
    cd "${ROOT_DIR}/${module}"
    GOCACHE="${GOCACHE_DIR}" GOLANGCI_LINT_CACHE="${LINT_CACHE_DIR}" "${LINTER_BIN}" run ./...
  )
done
