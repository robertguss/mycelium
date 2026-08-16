# Agent Rules — Idea Instance

Instance-facing rules for humans and agents operating this idea.

Open this folder in an agent runtime. State the human's goal. Follow
`.agents/skills/session/SKILL.md`. The CLI is the ledger; the agent is
the sparring partner.

## CLI

Use the `mycelium` CLI for scaffolded operations:

- `mycelium version` — print CLI version
- `mycelium new idea <name>` — scaffold a new idea (flags: `--dir`, `--offline`, `--publish`, `--tier`)
- `mycelium new <type> <title>` — generate a registered artifact
- `mycelium check` — verify convention conformance
- `mycelium tier <tier>` — raise or lower rigor tier (`focused` | `standard` | `high-assurance`)
- `mycelium state <target>` — lifecycle transition (flags: `--dir`, `--revisit`)
- `mycelium wake` — simmering → exploring ritual (writes re-entry brief)
- `mycelium status` — single-instance status (`--all` for portfolio; `--offline` when hermetic)
- `mycelium index` — rebuild `index.md`
- `mycelium publish` — create/update the GitHub repo for this idea
- `mycelium supersede <OLD-ID> --by <NEW-ID>` — supersede one artifact with another
- `mycelium handoff` — write `handoff/PACKET.md` then set `state = handed-off` (legal only from `clarified`)

Idea handoff is `mycelium handoff`. Use `mycelium state handed-off` only if the packet already exists.

Artifact supersede is `mycelium supersede` (front-matter cross-links + log line).
It does not change idea `state`. Idea lifecycle stays `mycelium state` / `wake`.

Global exit codes: `0` success, `1` failure. Teaching errors go to stderr.

## Implementation systems

The isolated implementer receives **ONLY** `handoff/`. No chat history. No instance source beyond the packet. Default system is **pstack/poteto-mode**. `manual` (read the packet, implement, run `handoff/acceptance/`) is the floor and always legal. Do not fetch the rest of the instance. Do not reopen locked decisions. Do not treat v1 attachment manifests as the packet.

See `program/reference/implementation-systems.md` for the packet → pstack/poteto mapping table.

## Skills

New scaffolds emit `.agents/skills/{session,mycelium-cli,spark,thinking,simmer,wake,clarify,handoff,portfolio,council,second-opinion}/SKILL.md` when `program/packs/council/` is present. `tier`, `index`, `state`, and `wake` do not retrofit skills into older instances. Existing instances: re-scaffold or copy the pack **manually**.

Runtimes that cannot fan out skip rungs 2–3; sparring still applies.

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
