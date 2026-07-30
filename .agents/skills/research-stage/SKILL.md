---
name: research-stage
description: >
  Produce the just-in-time five-item stage package for an artifact-driven
  research program: canonical stage prompt, repository installation task,
  attachment manifest, fresh-session launch message, and post-stage validation
  task with recommended commit message. Use when commissioning a stage,
  generating a research prompt, building an attachment list, launching a fresh
  session, or running /research-stage. Does not execute the substantive stage
  itself in the packaging session.
---

# Research Stage Package

Generate a complete **just-in-time stage package**. Do not execute the
substantive research, synthesis, review, or revision in this session unless the
human explicitly commissions that stage **and** provides the full attachment
set for a fresh-session role—and even then, prefer packaging first.

## Read first

1. `AGENTS.md`
2. `research-program.toml` (stage id, depends_on, outputs, identifier ranges)
3. Accepted Blueprint + Charter
4. `program/contracts/handoffs.md`
5. `program/templates/stage-package.md`
6. Relevant artifact contract under `program/contracts/`
7. Relevant template under `program/templates/`

## Five required deliverables

1. **Canonical stage prompt** — self-contained; all sections from the matching
   contract (e.g. focused research prompt contract).
2. **Repository installation task** — destination path under `docs/prompts/`;
   do not mix prompt install with report execution.
3. **Attachment manifest** — template `program/templates/attachment-manifest.md`;
   store under `docs/handoffs/<stage-id>-attachment-manifest.md` when installing.
4. **Fresh-session launch message** — `program/templates/launch-message.md`.
5. **Validation task + recommended commit message** —
   `program/templates/validation-task.md`;
   commit patterns in `program/reference/commit-boundaries.md`.

## Attachment selection rules

Include:

1. Governing artifacts in full (Blueprint, Charter as applicable)
2. Current stage prompt in full
3. Direct prerequisite reports in full
4. Accepted Decision Records
5. Handoff Digests for indirect reports
6. Full indirect reports when nuance, weak evidence, or conflict is material

Synthesis and adversarial review should get all materially relevant **full**
reports unless reliable repo retrieval is available.

## Just-in-time rules

- Independent research prompts may be generated together after Blueprint +
  Charter acceptance.
- Dependent, reconciliation, synthesis, review, revision, and implementation
  planning prompts only after prerequisites are **accepted**.
- Inherit actual upstream recommendations, IDs, weak evidence, contradictions,
  open questions, risks, spike results, handoff requirements.
- Allocate non-overlapping identifier ranges declared in Blueprint/manifest.

## Focused research tracks

This template does **not** pre-create track files. Create:

```text
docs/prompts/NN-<focus>-research-prompt.md
docs/reports/NN-<focus>-research-report.md
```

and add `[[stages]]` to `research-program.toml`. Use
`program/templates/focused-research-prompt.md`.

## Output behavior

- Prefer writing package artifacts to reserved paths when asked to install.
- Do not mark stages `accepted`.
- Do not begin a downstream substantive stage.
- Do not invent project decisions absent from approved discovery/Blueprint.
