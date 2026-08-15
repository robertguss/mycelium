# mycelium

Single-binary Go CLI that scaffolds and operates idea repos. Not an ADRP GitHub template.

## Install

See [`docs/install.md`](docs/install.md). Then:

```bash
mycelium version
```

## First commands

| Command | Purpose |
| --- | --- |
| `mycelium new idea` | Scaffold a new idea repo |
| `mycelium check` | Verify convention conformance |
| `mycelium state` | Lifecycle transition |
| `mycelium wake` | Simmering → exploring ritual (writes re-entry brief) |
| `mycelium status` | Single-instance status |
| `mycelium status --all` | Portfolio status across local (and optional GitHub) ideas |
| `mycelium supersede` | Supersede one artifact with another (does not change idea state) |
| `mycelium handoff` | Write `handoff/PACKET.md` then `state=handed-off` (legal only from `clarified`) |
| `mycelium index` | Rebuild `index.md` |
| `mycelium tier` | Raise or lower rigor tier |
| `mycelium publish` | Create or update the GitHub repo for this idea |

## Lifecycle

Exact DEC-006:

```text
spark → exploring ⇄ simmering → clarified → handed-off
any → archived
```

Reach `handed-off` via `mycelium handoff` (or `mycelium state handed-off` only if the packet already exists).

## This repo vs an idea repo

This repository builds the `mycelium` CLI. An idea repo has `mycelium.toml` (DEC-012). Idea-repo agent rules are `program/skeleton/AGENTS.md`, not this root `AGENTS.md`.

## program/ and framework/

- `program/` is shipped methodology and is emitted into new idea repos.
- `framework/` is this repo's build record and is **never** emitted.

## What to read next

- [`docs/install.md`](docs/install.md) — install the binary
- [`program/README.md`](program/README.md) — methodology library
- [`framework/blueprint.md`](framework/blueprint.md) — product blueprint

CLI version: `0.1.0-dev`.
