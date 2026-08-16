# Methodology library

This directory is the methodology shipped into every idea by the
`mycelium` CLI. Humans and agents attach only what the current sitting
needs.

You are probably in an idea repo. Start at
[`operator/getting-started.md`](operator/getting-started.md) and
`.agents/skills/session/SKILL.md`. Do not begin in this README.

## Layers

| Layer | Path | Use when |
| --- | --- | --- |
| **Operator** | [`operator/`](operator/) | Starting, resuming, bootstrapping, approval gates, completion |
| **Skills** | [`skills/`](skills/) (emitted to `.agents/skills/`) | Session rituals the agent should follow |
| **Contracts** | [`contracts/`](contracts/) | Writing or validating a specific artifact type |
| **Templates** | [`templates/`](templates/) | Copy-paste structures for prompts, records, tasks |
| **Reference** | [`reference/`](reference/) | Deep rules: tiers, stage library, anti-patterns, amendments |

## Authority

1. Accepted `DEC-###` records in **this idea** that supersede earlier authority
2. Locked decisions in the Program Blueprint (when the instance uses one)
3. Normative rules in the Research Charter (when the instance uses one)
4. The commissioning prompt for the current stage
5. Current accepted revised definitive specification
6. Accepted focused research reports (evidence and recommendations)
7. Adversarial reviews (proposed corrections)
8. Current accepted revised implementation plan
9. `mycelium.toml` (operational index only)
10. Community convention
11. Model or reviewer preference

**Chat history, model memory, and unstored reasoning are never authoritative.**

See [`contracts/authority-and-precedence.md`](contracts/authority-and-precedence.md).

## Two ways to use this library

**Short sitting.** Spark → explore → maybe simmer. Skills under
`.agents/skills/` are enough. Do not open the governance spine.

**Long program.** Discovery through handoff. Spine:

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
`mycelium.toml`, type homes). This `program/` tree is methodology; do not
put project conclusions here.

## Fresh-session rule

Every substantive stage of the long program runs in a **fresh** agent
session with a self-contained context packet. A session may prepare
prompts and mechanical fixes — it must not execute multiple substantive
stages in one context.

## Version

- **Methodology version:** 2.0.0
- **Status:** Shipped with the mycelium CLI into idea instances
