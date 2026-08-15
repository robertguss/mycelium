#!/usr/bin/env bash
# Build checksummed release binaries for linux-amd64 and darwin-arm64.
# Does not tag, commit, upload, or invoke gh.
set -euo pipefail

usage() {
  echo "usage: scripts/release.sh <version>" >&2
  exit 2
}

if [[ $# -ne 1 ]]; then
  usage
fi

VERSION="$1"
if [[ -z "$VERSION" ]]; then
  usage
fi

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

CHANGELOG="$ROOT/CHANGELOG.md"
DIST="$ROOT/dist"
HEADING_PLAIN="## [${VERSION}]"
HEADING_DATED_PREFIX="## [${VERSION}] - "

has_heading() {
  local line
  while IFS= read -r line || [[ -n "$line" ]]; do
    if [[ "$line" == "$HEADING_PLAIN" || "$line" == "$HEADING_DATED_PREFIX"* ]]; then
      return 0
    fi
  done <"$CHANGELOG"
  return 1
}

if [[ ! -f "$CHANGELOG" ]] || ! has_heading; then
  echo "release refused: CHANGELOG.md missing heading ${HEADING_PLAIN} (optional ' - YYYY-MM-DD' suffix allowed)" >&2
  exit 1
fi

mkdir -p "$DIST"
rm -f "$DIST/mycelium-linux-amd64" "$DIST/mycelium-darwin-arm64" "$DIST/SHA256SUMS"

LDFLAGS="-X github.com/robertguss/mycelium/internal/version.Version=${VERSION}"

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "$DIST/mycelium-linux-amd64" ./cmd/mycelium
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "$LDFLAGS" -o "$DIST/mycelium-darwin-arm64" ./cmd/mycelium

(
  cd "$DIST"
  sha256sum mycelium-linux-amd64 mycelium-darwin-arm64 >SHA256SUMS
)

echo "release ${VERSION}: wrote dist/mycelium-linux-amd64 dist/mycelium-darwin-arm64 dist/SHA256SUMS"
