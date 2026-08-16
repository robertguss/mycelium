---
name: second-opinion
description: >
  Cheap second-opinion rung for a Mycelium idea: one different model, identical
  commissioning prompt, exactly one model-report, no reconciliation. Use when
  program/packs/council/ is present. Does not invent council/ladder/replicate
  verbs. Does not commit unless the human asks.
---

# Second opinion

Perspective-ladder rung 2 (DEC-008). Credits: pstack second-opinion doctrine;
DEC-008.

## When

One word to invoke: get exactly one different model's independent read of the
same commissioning prompt. Agreement is high-signal; disagreement surfaces a
fork. No RCL.

## Procedure

1. Read `index.md` and `CONTEXT.md` (not the whole tree).
2. `mycelium new commissioning "<title>"`. Keep or set `rung = "second-opinion"`,
   `opt_in = true`, `cost_class = "cheap"`, `adapter` to `cursor` or `manual`.
3. Fill `## Prompt` (identical for the one model), `## Attachments`, `## Cost`.
4. Follow `program/packs/council/adapters/cursor.md` or
   `program/packs/council/adapters/manual.md` for a single model. Do not ask the
   CLI to call a model.
5. `mycelium new model-report` exactly once. Set `commissioning`,
   `rung = "second-opinion"`, `adapter`, `model`, and `prompt_sha256` to the
   sha256 hex of the trimmed `## Prompt` body. Fill `## Position`,
   `## Findings`, `## Dissent` (body may be `none`).
6. Do **not** create a reconciliation (no RCL for this rung).
7. Run `mycelium check` before handing back.
8. Do not `git commit` unless the human asks. The CLI never commits.
9. If the runtime cannot fan out, skip this rung; sparring still applies.

## Surface

Ladder surface when `program/packs/council/` is present:
`mycelium new commissioning`, `mycelium new model-report`, `mycelium check`,
plus this skill. **No portable council CLI.** No `mycelium new reconciliation`
for this rung.

## Do not

- Run `mycelium council`, `mycelium ladder`, or `mycelium replicate`
  (those verbs do not exist).
- Create an RCL for a second-opinion commissioning.
- Treat CI / an Actions job as the done bar.
- Flip `handed-off`.
- `git commit` unless the human asks.
