#!/usr/bin/env sh
set -eu

if [ "$#" -ne 1 ]; then
  echo "usage: $0 <version-tag>" >&2
  exit 1
fi

VERSION="$1"
ROOT_DIR="$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)"
RELEASE_DIR="${ROOT_DIR}/release/${VERSION}"
GOCACHE_DIR="${GOCACHE:-/tmp/crona-go-release-cache}"

TARGETS="
darwin arm64
darwin amd64
linux amd64
linux arm64
"

rm -rf "${RELEASE_DIR}"
mkdir -p "${RELEASE_DIR}"

for target in ${TARGETS}; do
  :
done

echo "${TARGETS}" | while read -r GOOS GOARCH; do
  [ -n "${GOOS}" ] || continue

  echo "Building ${GOOS}/${GOARCH}"
  env CGO_ENABLED=0 GOOS="${GOOS}" GOARCH="${GOARCH}" GOCACHE="${GOCACHE_DIR}" \
    go build -o "${RELEASE_DIR}/crona-kernel-${VERSION}-${GOOS}-${GOARCH}" ./kernel/cmd/crona-kernel
  env CGO_ENABLED=0 GOOS="${GOOS}" GOARCH="${GOARCH}" GOCACHE="${GOCACHE_DIR}" \
    go build -o "${RELEASE_DIR}/crona-tui-${VERSION}-${GOOS}-${GOARCH}" ./tui
done

sed "s/__VERSION__/${VERSION}/g" "${ROOT_DIR}/scripts/install_tui.sh.tmpl" > "${RELEASE_DIR}/install-crona-tui.sh"
chmod +x "${RELEASE_DIR}/install-crona-tui.sh"

(
  cd "${RELEASE_DIR}"
  shasum -a 256 ./* > checksums.txt
)

echo "Release artifacts written to ${RELEASE_DIR}"
