---
name: portfolio
description: >
  Cross-idea Mycelium portfolio view: run status --all, interpret local-only
  partials honestly, and do not mutate. Use when surveying many ideas or the
  human asks for portfolio status. Pass --offline when hermetic or gh is absent.
---

# Portfolio

Cross-idea survey.

## Steps

1. `mycelium status --all` (pass `--offline` when hermetic / no `gh`)
2. Interpret `partial: local-only (...)` as incomplete — do not invent
   remote ideas
3. `status --all` may print `partial: legacy-manifest` and still list
   other ideas
4. Report what is simmering, what is due, what is unpublished

Done: the listing is shown and partials are named.

Do not create repos, publish, or mutate instances from this skill.
Do not `git commit` unless the human asks.
