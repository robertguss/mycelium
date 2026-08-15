# DEC-012 — New instances use mycelium.toml as the manifest filename

- **Status:** Accepted
- **Date:** 2026-08-14
- **Owner:** Robert Guss
- **Supersedes:** None (settles blueprint OQ-002)
- **Related recommendations:** None
- **Related evidence:** DEC-010 (CLI scaffold); DEC-011 (two version fields)

## Context

Blueprint OQ-002 asked whether new instances keep `research-program.toml`
or rename as a 2.0 migration. The master repo remains an ADRP v1 instance
and already has `research-program.toml`. New idea repos are Mycelium 2.0
from birth.

## Decision

1. Scaffolded instances use `mycelium.toml` as the sole manifest filename.
2. The master repository keeps `research-program.toml` and is not converted.
3. Runtime commands detect an instance by the presence of `mycelium.toml`.
   They do not read `research-program.toml` as a 2.0 manifest.
4. No migration machinery is added to rename existing files (DEC-011).

## Rationale

The 2.0 name matches the product. Keeping master's v1 filename avoids
pretending the master is a scaffolded idea and keeps `just check` working.

## Consequences

PHASE-01 contracts, checker, and scaffold all say `mycelium.toml`.
Documentation that still says `research-program.toml` refers to master v1.

## Alternatives Considered

Keep `research-program.toml` everywhere (rejects the 2.0 name).
Rename master too (converts master into an instance; rejected).

## Risks

Agents open master and look for `mycelium.toml`. Mitigation: AGENTS.md on
master stays v1; instance AGENTS.md is emitted and names `mycelium.toml`.

## Revisit Triggers

If master itself is ever re-scaffolded as a Mycelium instance, reopen.

## Approval

Settled by Arvo 2026-08-14 (OQ-002). Recorded by Architect in the
PHASE-01 implementation brief.
