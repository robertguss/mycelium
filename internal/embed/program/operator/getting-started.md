# Getting Started

You are inside an **idea instance** (`mycelium.toml` is in this tree). The
product CLI is already installed. This file is the in-repo start page for
you and the agent.

Mycelium is meant to be used **with an LLM**. Open this folder in your
agent runtime. State your goal. Follow `.agents/skills/session/SKILL.md`.

The CLI is the ledger. It does not invent scope, accept stages, or commit.

## What this idea is for

A durable place to think. The program ends, when you take it that far, with:

- A validated, adversarially reviewed **definitive specification**
- A revised **implementation plan** organized into **phases and milestones**

It **stops before** a granular coding backlog. The last artifact is
`handoff/PACKET.md`.

You do not owe that spine to every idea. A thirty-minute spark is a legal
use.

## Pick a session

Say the goal in the first message. The `session` skill maps it:

| You want | Skill |
| --- | --- |
| Brand-new idea, first thought | `spark` |
| Work a question or decide something | `thinking` |
| Park it until a date or event | `simmer` |
| Return after a gap | `wake` |
| Destination reached | `clarify` |
| Ready to implement | `handoff` |
| Survey many ideas | `portfolio` |
| Multi-day research program | this file + `program/reference/governance-spine.md` |
| Another model’s read | `second-opinion` |
| Multi-model replication | `council` |

`thinking` stays on in every non-archived working session.

## First sitting on a new idea

1. Read `index.md` and `mycelium.toml` (not the whole tree).
2. Interview the human **one question at a time**, with a recommendation
   each time, until the first thought is stateable.
3. `mycelium state exploring`
4. `mycelium new question "…"` or `mycelium new decision "…"`
5. `mycelium check`

Do not invent tracks, a blueprint, or conclusions. Do not commit unless
the human asks.

## Longer program (when the idea earns it)

```text
Discovery → Blueprint → Charter
  → Adaptive focused-research graph
  → Optional replication / reconciliation
  → Synthesis → spec review → revised spec
  → Implementation plan → plan review → final plan
  → mycelium handoff
```

Raise the tier first (`mycelium tier standard` or `high-assurance`). Run
each substantive stage in a **fresh** agent session with a self-contained
packet. Human approval gates: [`approval-gates.md`](approval-gates.md).

## What `mycelium new idea` already did

**Did:** scaffold layout, write `mycelium.toml`, emit `program/` and
skills, set `state = spark` and default `tier = focused`.

**Did not:** invent decisions, skip discovery, or git-commit.

## Required reading order

1. This idea’s `README.md` and `AGENTS.md`
2. This file and [`resume-protocol.md`](resume-protocol.md)
3. `.agents/skills/session/SKILL.md`
4. Accepted Blueprint and Charter (once they exist)
5. Current implementation authority (revised spec / final plan when accepted)

## Tooling

| Command | Purpose |
| --- | --- |
| `mycelium check` | Convention conformance |
| `mycelium status` | This idea (`--all` for the portfolio) |
| `mycelium state` / `wake` | Lifecycle |
| `mycelium new <type> "<Title>"` | Next artifact |
| `mycelium handoff` | Packet + `handed-off` from `clarified` |

Full command list: `.agents/skills/mycelium-cli/SKILL.md`.
