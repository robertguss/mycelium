---
name: council
description: >
  Opt-in council rung for a Mycelium idea: commission identical prompts across
  N models, collect per-model reports, and reconcile with retained dissent.
  Suggest only on v1 replication triggers. Use when program/packs/council/ is
  present. Does not invent council/ladder/replicate verbs. Does not commit
  unless the human asks.
---

# Council

Opt-in perspective-ladder rung 3 (DEC-008). Credits: DEC-008; v1
`program/contracts/replication-reconciliation.md`; karpathy/llm-council
(three-stage shape, subordinated to v1 reconciliation); pstack second-opinion
doctrine (neighboring rung 2 move).

## Panel defaults

Documented defaults when `~/.config/mycelium/panels.toml` is absent (missing
file is legal; check does not read it; `$MYCELIUM_CONFIG` is a directory
override for that file):

| cost_class | panel size |
| --- | --- |
| quick | 2 |
| standard | 3 |
| high-stakes | 4 |

State panel size and `cost_class` before running.

## When

Suggest a council only when a v1 replication trigger fires (security-critical,
safety-critical, legally or financially consequential, difficult to reverse,
architecturally foundational, weak or conflicting evidence, ecosystem/vendor
bias risk, or still low-confidence after a spike). Source:
`program/contracts/replication-reconciliation.md`. Never auto-run. Never
require a council to leave `spark`.

## Procedure

1. Read `index.md` and `CONTEXT.md` (not the whole tree).
2. Suggest a council only when a v1 replication trigger fires. State panel size
   and `cost_class` first.
3. `mycelium new commissioning "<title>"`. Set `rung = "council"`,
   `opt_in = true`, `cost_class` to `quick` / `standard` / `high-stakes`,
   `adapter` to `cursor` or `manual`.
4. Fill `## Prompt` (identical for every model), `## Attachments`, `## Cost`.
5. Follow `program/packs/council/adapters/cursor.md` or
   `program/packs/council/adapters/manual.md`. Do not ask the CLI to call a
   model.
6. `mycelium new model-report` once per model. Set `commissioning`, `rung`,
   `adapter`, `model`, and `prompt_sha256` to the sha256 hex of the trimmed
   `## Prompt` body. Fill `## Position`, `## Findings`, `## Dissent` (body may
   be `none`).
7. `mycelium new reconciliation`. Set `commissioning` and `rung = "council"`.
   Fill every v1 H2, including `## Retained dissent`. Do **not** choose by
   majority vote, length, confidence of prose, or model reputation.
8. Run `mycelium check` before handing back. On `prompt_sha256` mismatch the
   teaching error names the expected hex.
9. Do not `git commit` unless the human asks. The CLI never commits.
10. If the runtime cannot fan out, skip this rung; sparring still applies.

## Surface

Ladder surface when `program/packs/council/` is present:
`mycelium new commissioning`, `mycelium new model-report`,
`mycelium new reconciliation`, `mycelium check`, plus this skill. **No portable
council CLI.**

## Do not

- Run `mycelium council`, `mycelium ladder`, `mycelium replicate`, or
  `mycelium handoff` (those verbs do not exist).
- Treat CI / an Actions job as the done bar.
- Flip `handed-off`.
- `git commit` unless the human asks.
