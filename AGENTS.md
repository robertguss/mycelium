# Agent Rules — mycelium

This repository is the **mycelium product repo** (the CLI), not an idea instance.
Idea-instance rules live in `program/skeleton/AGENTS.md`.

Human-facing docs: [`README.md`](README.md), [`docs/user-guide.md`](docs/user-guide.md),
[`docs/install.md`](docs/install.md).

## Verify

```bash
CGO_ENABLED=0 go test ./...
```

If Go or embed is touched, product packages must stay ≥85% coverage (exclude
generated, vendor, and data-only fixtures). Stamp **85**.

## Never emit framework/

Scaffold must never emit `framework/`. `ForbiddenPaths` stay:

```text
framework
Justfile
scripts
research-program.toml
```

Deleting leftover ADRP files from *this* repo does not change `ForbiddenPaths`.

## Git and commits

The CLI never git-commits instance work product. Agents do not push to `main`.

## Cursor Cloud

- Cloud env name: exactly `robertguss/mycelium`
- Go **1.26** at `/usr/local/go` or `/usr/local/bin/go`
- Distro `/usr/bin/go` 1.22 is wrong — use the 1.26 toolchain
- `CGO_ENABLED=0`
- No long-running services, databases, Docker Compose, Node, or web servers

## Operator surface

Use `mycelium …`. The only Just recipe is:

```bash
just release version=…
```

That wraps `scripts/release.sh`. Do not restore ADRP bootstrap recipes
(`init` / `status` / `check`).

## Authority

Shipped idea rules live in `program/contracts/` (emitted into instances).
Operator walkthrough: `docs/user-guide.md`. After scaffold, start at
`program/operator/getting-started.md`.

## Build

```bash
CGO_ENABLED=0 go build -o mycelium ./cmd/mycelium
```
