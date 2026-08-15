# DEC-004 — Template-owned, self-contained, multi-runtime methodology; version-pinned instances

- **Status:** Accepted
- **Date:** 2026-08-14
- **Owner:** Robert Guss
- **Supersedes:** None (supersedes an unrecorded in-session proposal)
- **Related recommendations:** None
- **Related evidence:** `.agents/skills/` convention loading across runtimes;
  `program/reference/amendment-protocol.md`

## Context

The agent initially recommended housing the methodology in a Cursor plugin
(single source of truth, thin scaffolded instances). Robert then surfaced the
deciding constraint: he works across many models, agents, and cloud VMs, and
a self-contained git repo is critical to how he works. Tooling in this space
also changes weekly to monthly.

## Decision

The master template owns the methodology. Instances are self-contained:
clone = everything an agent needs (`AGENTS.md`, `.agents/skills/`,
`program/`, `scripts/`). Each instance pins the methodology version at
creation; upgrades are explicit (`just upgrade` applying
`program/migrations/`, per DEC-005), never silent. No Cursor plugin is the
system of record; thin per-runtime adapters remain permitted as
conveniences.

## Rationale

Multi-runtime portability dominates. A git repo is the universal artifact;
plugins are runtime-specific caches that die on update. Pinning is
epistemically correct for in-flight ideas: a program whose methodology
changes under it mid-run is an audit hazard. The plugin-owned proposal was
withdrawn on redesign-from-first-principles grounds once the constraint
surfaced.

## Consequences

Methodology fixes reach live instances only through opt-in migrations.
Runtime-specific optimizations (e.g., Cursor parallel subagents) must remain
optional layers over portable contracts.

## Alternatives Considered

Plugin-owned methodology (agent's original recommendation; fails
multi-runtime). Full-snapshot vendoring with no master evolution (frozen
spec, no living framework).

## Risks

Master/instance drift accumulating. Mitigation: migrations + amendment
protocol + version field in the manifest.

## Revisit Triggers

Agent's retained caveat, recorded as the crux: if maintaining live instances
across versions becomes painful in practice, revisit central distribution
mechanisms.

## Approval

Approved by Robert Guss in the Mycelium discovery interview, 2026-08-14
("I don't want this to be solely a cursor plugin... having a self-contained
git repo with everything I need is critical").
