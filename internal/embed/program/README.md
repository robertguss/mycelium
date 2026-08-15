# Research Program Methodology

This directory is the **methodology library** shipped into idea repos by the
`mycelium` CLI. It redistributes the artifact-driven research program into
layers so agents and humans can attach only what each stage needs—without
losing rigor.

## Layers

| Layer         | Path                       | Use when                                                      |
| ------------- | -------------------------- | ------------------------------------------------------------- |
| **Operator**  | [`operator/`](operator/)   | Starting, resuming, bootstrapping, approval gates, completion |
| **Contracts** | [`contracts/`](contracts/) | Writing or validating a specific artifact type                |
| **Templates** | [`templates/`](templates/) | Copy-paste structures for prompts, records, ledgers, tasks    |
| **Reference** | [`reference/`](reference/) | Deep rules: tiers, stage library, anti-patterns, amendments   |

## Authority

1. Accepted `DEC-###` records that supersede earlier authority
2. Locked decisions in the Program Blueprint (when the instance uses one)
3. Normative rules in the Research Charter (when the instance uses one)
4. The commissioning prompt for the current stage
5. Current accepted revised definitive specification
6. Accepted focused research reports (evidence and recommendations)
7. Adversarial reviews (proposed corrections)
8. Current accepted revised implementation plan
9. `mycelium.toml` (operational index only; DEC-012)
10. Community convention
11. Model or reviewer preference

**Chat history, model memory, and unstored reasoning are never authoritative.**

See [`contracts/authority-and-precedence.md`](contracts/authority-and-precedence.md).

## Fixed governance spine

```text
mycelium new idea → Discovery → Blueprint → Charter
  → Adaptive focused-research graph
  → Optional replication / reconciliation
  → Chief Architect synthesis → Spec adversarial review → Revised spec
  → Implementation plan → Plan adversarial review → Final revised plan
  → Program closure and implementation handoff
```

Details: [`reference/governance-spine.md`](reference/governance-spine.md).

## Idea-repo tree

Idea artifacts live at the instance root (`README.md`, `AGENTS.md`,
`mycelium.toml`, type homes, `docs/` as needed). This `program/` tree is
methodology; do not put project conclusions here.

Create an idea with `mycelium new idea`. The instance operational index is
`mycelium.toml`. Operate with the `mycelium` CLI (`check`, `status`, `state`,
`wake`, `handoff`, …).

## Fresh-session rule

Every substantive stage runs in a **fresh** LLM/agent session with a
self-contained context packet (attachment manifest). A session may prepare
prompts, manifests, and mechanical fixes—but must not execute multiple
substantive stages in one context.

## Skills

New scaffolds emit mycelium skills under `.agents/skills/` (see
`program/skeleton/AGENTS.md`). Prefer those entry points when your agent
supports skills; otherwise follow the instance `AGENTS.md` and this `program/`
tree.

## Version

- **Methodology version:** 2.0.0
- **Status:** Shipped with the mycelium CLI into idea instances
