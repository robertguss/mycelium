# Research Program Methodology

This directory is the **methodology library** for the artifact-driven research
program. It redistributes the v1 master specification into three layers so
agents and humans can attach only what each stage needs—without losing rigor.

## Layers

| Layer         | Path                       | Use when                                                      |
| ------------- | -------------------------- | ------------------------------------------------------------- |
| **Operator**  | [`operator/`](operator/)   | Starting, resuming, bootstrapping, approval gates, completion |
| **Contracts** | [`contracts/`](contracts/) | Writing or validating a specific artifact type                |
| **Templates** | [`templates/`](templates/) | Copy-paste structures for prompts, records, ledgers, tasks    |
| **Reference** | [`reference/`](reference/) | Deep rules: tiers, stage library, anti-patterns, amendments   |

## Authority

1. Accepted `DEC-###` records that supersede earlier authority
2. Locked decisions in `docs/00-program-blueprint.md`
3. Normative rules in `docs/01-research-charter.md`
4. The commissioning prompt for the current stage
5. Current accepted revised definitive specification
6. Accepted focused research reports (evidence and recommendations)
7. Adversarial reviews (proposed corrections)
8. Current accepted revised implementation plan
9. `research-program.toml` (operational index only)
10. Community convention
11. Model or reviewer preference

**Chat history, model memory, and unstored reasoning are never authoritative.**

See [`contracts/authority-and-precedence.md`](contracts/authority-and-precedence.md).

## Fixed governance spine

```text
Init (just init) → Discovery → Blueprint → Charter
  → Adaptive focused-research graph
  → Optional replication / reconciliation
  → Chief Architect synthesis → Spec adversarial review → Revised spec
  → Implementation plan → Plan adversarial review → Final revised plan
  → Program closure and implementation handoff
```

Details: [`reference/governance-spine.md`](reference/governance-spine.md).

## Project-facing tree

Project artifacts live at the repository root (`README.md`, `AGENTS.md`,
`research-program.toml`, `decisions/`, `docs/`). This `program/` tree is
methodology; do not put project conclusions here.

## Fresh-session rule

Every substantive stage runs in a **fresh** LLM/agent session with a
self-contained context packet (attachment manifest). A session may prepare
prompts, manifests, and mechanical fixes—but must not execute multiple
substantive stages in one context.

## Skills

Agent entry points (thin; they point here):

- `.agents/skills/research-program/` — discovery, resume, next stage
- `.agents/skills/research-stage/` — just-in-time stage package
- `.agents/skills/research-validate/` — independent validation gate

## Version

- **Methodology version:** 1.0
- **Source:** Redistributed from Artifact-Driven Research Program master spec v1.0
- **Status:** Accepted reusable workflow (template form)
