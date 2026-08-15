# Attachment Manifest — blueprint-revision

## Required Full Artifacts

1. `framework/prompts/01-blueprint-revision-prompt.md` — the commissioning
   prompt; defines the task, allowed scope, and binding disposition
   resolutions
2. `framework/blueprint.md` — the proposed artifact under revision
3. `framework/reviews/01-blueprint-adversarial-review.md` — all thirteen
   accepted findings with Required Corrections and Robert's Dispositions
4. `AGENTS.md` — repository operating rules (fresh sessions, file scope,
   authority)

## Required Decision Records

- `framework/decisions/DEC-001-evolve-adrp-in-place-into-mycelium.md`
- `framework/decisions/DEC-002-durable-tiered-record-as-the-product.md`
- `framework/decisions/DEC-003-one-repository-per-idea.md`
- `framework/decisions/DEC-004-template-owned-self-contained-multi-runtime.md`
- `framework/decisions/DEC-005-convention-over-configuration.md`
- `framework/decisions/DEC-006-idea-lifecycle-with-simmer.md`
- `framework/decisions/DEC-007-sparring-stance-agreement-states-cruxes.md`
- `framework/decisions/DEC-008-perspective-ladder-opt-in-councils.md`
- `framework/decisions/DEC-009-name-mycelium.md`
- `framework/decisions/DEC-010-mycelium-is-a-cli.md`
- `framework/decisions/DEC-011-defer-migrations.md`

## Required Handoff Digests

- None. All prerequisites attach in full.

## Explicitly Excluded Artifacts

- `program/` methodology internals — this stage revises blueprint text only;
  the blueprint's references to `program/` paths are unchanged by the
  commissioned corrections.
- `research-program.toml` — unfilled v1 template; the framework evolution
  loop does not run through it.

## Authority Notes

Precedence within this stage: accepted `DEC-001`–`DEC-011`, then the review's
Dispositions section, then each finding's Required Correction, then the
current blueprint text. Chat history and model memory are not authoritative.

## Expected Output

`framework/blueprint.md` — revised in place; Status remains
`Proposed — awaiting acceptance by Robert Guss`.
