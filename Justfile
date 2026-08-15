# mycelium product repo — release helper only (not an ADRP operator file).

set shell := ["bash", "-euo", "pipefail", "-c"]

# Default: list available recipes
default:
    @just --list

# Build checksummed release binaries (linux-amd64 + darwin-arm64). No tag/upload.
# Usage: just release version=0.1.0
release version:
    bash scripts/release.sh {{quote(version)}}
