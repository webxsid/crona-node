#!/usr/bin/env sh
set -eu

cd "$(dirname "$0")/.."
CRONA_ENV=Dev go run ./cli/cmd/crona dev clear
