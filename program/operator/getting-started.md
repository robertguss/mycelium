# Getting Started

## Purpose of this template

A reusable operating system for deep, multi-stage research with LLMs and
repository-aware agents. Ordinary brainstorming or a one-pass plan is not
enough for the projects this workflow targets.

The program ends with:

- A validated, adversarially reviewed **definitive specification**
- A revised **implementation plan** organized into **phases and milestones**

It deliberately **stops before** a granular coding backlog or agent task
packets.

## First-time flow (recommended)

1. **Create** a new repository from this GitHub template (or clone and rename).
2. **`just init name="your-working-title"`** — fills project name placeholders
   only. Does not run git. Does not invent scope or research tracks.
3. **Discovery** — Research Program Architect interviews you **one question at
   a time** (with a recommendation each time) until ≥95% confidence.
4. **Approve framing** — problem, outcome, locked scope, uncertainties, tracks,
   rigor tier.
5. **Accept Program Blueprint** (`docs/00-program-blueprint.md`).
6. **Accept Research Charter** (`docs/01-research-charter.md`).
7. **Execute the adaptive research graph** with just-in-time prompts, fresh
   sessions, validation, and human approval gates.
8. **Synthesis → adversarial review → revised spec → plan → plan review →
   final plan.**

## What `just init` does and does not do

**Does:** set `program_name`, derive `program_id`, set dates, replace
`{{PROJECT_NAME}}` / obvious title placeholders in root and fixed spine
metadata.

**Does not:** run git; invent decisions, tracks, or conclusions; accept any
stage; skip discovery.

Rigor tier defaults to `standard` as **proposed** until Blueprint approval.

## Required reading order for a new contributor

1. Root [`README.md`](../../README.md)
2. [`AGENTS.md`](../../AGENTS.md)
3. This file and [`resume-protocol.md`](resume-protocol.md)
4. Accepted Blueprint and Charter (once they exist)
5. Current implementation authority (revised spec / final plan when accepted)

## Tooling

| Command              | Purpose                                   |
| -------------------- | ----------------------------------------- |
| `just init name="…"` | Bootstrap project name into placeholders  |
| `just status`        | Manifest stage overview / eligible stages |
| `just check`         | Tree + placeholder-vs-accepted sanity     |

## Skills

Use `.agents/skills/` entry points when your agent supports skills; otherwise
follow `AGENTS.md` and this `program/` tree directly.
