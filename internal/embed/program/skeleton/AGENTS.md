# Agent Rules — Idea Instance

Instance-facing rules for humans and agents operating this idea.

## CLI

Use the `mycelium` CLI for scaffolded operations:

- `mycelium version` — print CLI version
- `mycelium new idea <name>` — scaffold a new idea (flags: `--dir`, `--offline`, `--publish`, `--tier`)
- `mycelium new <type> <title>` — generate a registered artifact
- `mycelium check` — verify convention conformance
- `mycelium tier <tier>` — raise or lower rigor tier (`focused` | `standard` | `high-assurance`)
- `mycelium publish` — create/update the GitHub repo for this idea

Global exit codes: `0` success, `1` failure. Teaching errors go to stderr.

## Manual floor

You may edit Markdown artifacts by hand. Keep `+++` TOML front matter valid, preserve required H2 headings, and keep filenames on the ID-to-path rule (`<home>/<NS>-<digits>-<slug>.md`). After manual edits, run `mycelium check`.

Do not invent new identifier namespaces. Do not emit `framework/`. Do not git commit unless the human asks.

## Teaching errors

Failed commands print four stderr lines:

```text
mycelium: <one-line failure>
convention: <name>
contract: program/contracts/<file>.md
fix: <command or rename>
```

Fix the named convention, then re-run the suggested command.
