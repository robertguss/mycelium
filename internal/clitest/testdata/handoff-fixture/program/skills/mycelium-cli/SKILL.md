---
name: mycelium-cli
description: >
  Operate a Mycelium idea instance with the mycelium CLI: scaffold, generate
  registered artifacts, check conformance, change tier, lifecycle transitions,
  wake, handoff, index, status, publish, and supersede. Use when working inside
  an idea repo that has mycelium.toml, or when the human asks for mycelium
  commands. Does not commit to git unless the human asks.
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
| `mycelium state <target> [--dir PATH] [--revisit …]` | Lifecycle transition |
| `mycelium wake [--dir PATH]` | Simmering → exploring ritual (writes re-entry brief) |
| `mycelium handoff [--dir PATH]` | Write `handoff/PACKET.md` then set `state = handed-off` (from `clarified` only) |
| `mycelium status [--dir PATH]` | Single-instance status |
| `mycelium status --all [--offline]` | Portfolio listing (local + GitHub when available) |
| `mycelium index [--dir PATH]` | Rebuild `index.md` |
| `mycelium publish [--dir PATH]` | GitHub create/topic when authenticated |
| `mycelium supersede <OLD-ID> --by <NEW-ID> [--dir PATH]` | Artifact supersede (bidirectional cross-links + log line) |

Exit `0` on success, `1` on failure.

## Handoff

`mycelium handoff [--dir PATH]` writes the implementation packet at
`handoff/PACKET.md`, then sets `state = handed-off`. Legal only from
`clarified`. `handed-off` is terminal except `archived`.

Refuse cases (teaching error, exit 1, no writes):

- not `clarified` — legal only from `clarified`
- `handoff/PACKET.md` already exists — use `mycelium state handed-off` instead
- `mycelium state handed-off` without a passing packet — run `mycelium handoff` first

The packet is `handoff/PACKET.md`. Do not treat v1 session-attachment manifests
as the packet.

## Supersede

`mycelium supersede <OLD-ID> --by <NEW-ID> [--dir PATH]` is **artifact-level**.
It sets OLD `status = "Superseded"` + `superseded_by`, and NEW `supersedes`.
It does **not** change idea `state` and does **not** implement `handed-off`.

Eligible: DEC, ASM, EVD, SPK. Refuse (teaching error, exit 1, no writes):

- ineligible type (including OQ — open a new question instead)
- idea-state token as an ID (`spark`, `exploring`, …)
- OLD already Superseded
- NEW already has `supersedes` (one-to-one this phase)
- different namespace, missing IDs, OLD == NEW

Idea lifecycle stays `mycelium state` / `wake`.

## Sparring

Sparring surface is `mycelium new question`, `mycelium new assumption`,
`mycelium new decision`, `mycelium new spike`, `mycelium check`, plus the
`thinking` skill. No `think` / `spar` verb.

## Perspective ladder

When `program/packs/council/` is present, ladder surface is
`mycelium new commissioning`, `mycelium new model-report`,
`mycelium new reconciliation`, `mycelium check`, plus pack skills `council` and
`second-opinion`. **No portable council CLI.** No `council` / `ladder` /
`replicate` verb.

Runtimes that cannot fan out skip rungs 2–3; sparring still applies.

## Skills on scaffold

New `mycelium new idea` scaffolds emit
`.agents/skills/{mycelium-cli,spark,wake,portfolio,thinking,council,second-opinion}/SKILL.md`
when the pack is present. `tier`, `index`, `state`, and `wake` do not retrofit
skills. Existing instances: re-scaffold or copy the pack **manually**.

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
