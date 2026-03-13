#!/usr/bin/env sh
set -eu

cd "$(dirname "$0")/.."
go run ./cli/cmd/crona-dev clear
