# Getting Started

## Purpose

A reusable operating system for deep, multi-stage research with LLMs and
repository-aware agents. Ordinary brainstorming or a one-pass plan is not
enough for the projects this workflow targets.

The program ends with:

- A validated, adversarially reviewed **definitive specification**
- A revised **implementation plan** organized into **phases and milestones**

It deliberately **stops before** a granular coding backlog or agent task
packets.

## First-time flow (recommended)

1. **Install** the `mycelium` CLI (see the product repo `docs/install.md`).
2. **`mycelium new idea "your-working-title"`** — scaffolds an idea repo with
   `mycelium.toml` and the methodology copy. Does not invent scope or research
   tracks. Does not git-commit instance work product.
3. **Discovery** — Research Program Architect interviews you **one question at
   a time** (with a recommendation each time) until ≥95% confidence.
4. **Approve framing** — problem, outcome, locked scope, uncertainties, tracks,
   rigor tier.
5. **Accept Program Blueprint** (when used for this idea).
6. **Accept Research Charter** (when used for this idea).
7. **Execute the adaptive research graph** with just-in-time prompts, fresh
   sessions, validation, and human approval gates.
8. **Synthesis → adversarial review → revised spec → plan → plan review →
   final plan.**

## What `mycelium new idea` does and does not do

**Does:** scaffold layout, write `mycelium.toml`, emit `program/` and skills,
set birth `state = spark` and default tier.

**Does not:** invent decisions, tracks, or conclusions; accept any stage; skip
discovery; git-commit instance work product.

Rigor tier defaults to `focused` at birth unless `--tier` is set; raise with
`mycelium tier` when appropriate.

## Required reading order for a new contributor

1. Instance [`README.md`](../../README.md) (idea root)
2. Instance [`AGENTS.md`](../../AGENTS.md)
3. This file and [`resume-protocol.md`](resume-protocol.md)
4. Accepted Blueprint and Charter (once they exist)
5. Current implementation authority (revised spec / final plan when accepted)

## Tooling

| Command | Purpose |
| --- | --- |
| `mycelium new idea "…"` | Scaffold a new idea repo |
| `mycelium check` | Convention conformance |
| `mycelium status` | Instance status (`--all` for portfolio) |
| `mycelium state` / `wake` | Lifecycle transitions |
| `mycelium handoff` | Packet + `handed-off` from `clarified` |

## Skills

Use `.agents/skills/` entry points when your agent supports skills; otherwise
follow `AGENTS.md` and this `program/` tree directly.
