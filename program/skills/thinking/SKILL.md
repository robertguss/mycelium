---
name: thinking
description: >
  Always-on sparring for a Mycelium idea: take positions on substantive
  questions, set agreement, record disagreement with cruxes when disputed,
  challenge glossary drift, and audit assumptions. Use in any non-archived
  working session. Does not invent think/spar verbs.
---

# Thinking

Always-on sparring (DEC-007). Credits: mattpocock grilling (recommendation per
question; decisions stay the user's); domain-modeling (glossary + ADR);
pstack poteto candor ("no is an acceptable answer").

## When

- Any non-archived state: apply this skill.
- On `archived`: do not open new OQs; existing records stay.

## Procedure

1. Read `index.md` and `CONTEXT.md` (not the whole tree).
2. On every substantive question: take a position; capture it with
   `mycelium new question` (or edit the OQ). Bare questions are a smell —
   the human or an adversarial reviewer judges substance, never an
   automated content score.
3. Set `agreement` to `open`, `aligned`, or `agree-to-disagree`.
4. If `agree-to-disagree`: fill `## Positions`, `## Reasons`, and
   `## Crux`, each with `### Human` and `### Agent`. The record retains
   both positions, both reasons, and cruxes.
5. If `aligned`: `## Positions` H2 present; no disagreement record
   required (no Crux/Reasons required).
6. If `open`: keep working; Positions container present; Crux/Reasons not
   required.
7. Do not edit `aligned` or `agree-to-disagree` back to `open`. Open a
   new OQ instead.
8. Glossary challenge: on drift or vagueness, sharpen the term; record
   under `CONTEXT.md` as `## <Term>` + `### Definition`.
9. Assumption audit: periodically dump presuppositions via
   `mycelium new assumption`. No AUDIT file.
10. A crux is eligible to become `mycelium new spike` or a research stage.
    Do not auto-promote. No new command.
11. Recommendation per question; decisions stay the user's (grilling).
    "no is an acceptable answer" (poteto candor).
12. Run `mycelium check` before handing back.
13. Do not `git commit` unless the human asks. The CLI never commits.

## Surface

Sparring uses existing commands only: `mycelium new question`,
`mycelium new assumption`, `mycelium new decision`, `mycelium new spike`,
`mycelium check`, plus this skill.

## Do not

- Run `mycelium think`, `mycelium spar`, `mycelium session`, or
  `mycelium council` (those verbs do not exist).
- Create a `DSG-###` file or a `session.md`.
- Flip `handed-off`.
- Treat CI / an Actions job as the done bar.
- `git commit` unless the human asks.
