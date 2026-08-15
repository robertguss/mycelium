---
name: mycelium-cli
description: >
  Operate a Mycelium idea instance with the mycelium CLI: scaffold, generate
  registered artifacts, check conformance, change tier, and publish. Use when
  working inside an idea repo that has mycelium.toml, or when the human asks
  for mycelium commands. Does not commit to git unless the human asks.
---

# Mycelium CLI

## Commands

| Command | Purpose |
| --- | --- |
| `mycelium version` | Print stamped version string |
| `mycelium new idea <name>` with flags `--dir`, `--offline`, `--publish`, `--tier` | Scaffold idea; always `git init`; never commit |
| `mycelium new <type> <title> [--dir PATH]` | Generate next ID for a registered type |
| `mycelium check [--dir PATH]` | Convention conformance |
| `mycelium tier <tier> [--dir PATH]` | Set tier; emit newly required dirs only; never delete |
| `mycelium publish [--dir PATH]` | GitHub create/topic when authenticated |

Exit `0` on success, `1` on failure.

## Manual floor

Hand-edit Markdown is allowed. Keep front matter, required H2s, and ID-to-path filenames. Run `mycelium check` after edits.

## Teaching-error shape

```text
mycelium: <one-line failure>
convention: <name>
contract: program/contracts/<file>.md
fix: <command or rename>
```

## Git

Do not `git commit` unless the human asks. The CLI never commits.
