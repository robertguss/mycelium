# DEC-001 — Evolve ADRP in place into Mycelium, one unified master

- **Status:** Accepted
- **Date:** 2026-08-14
- **Owner:** Robert Guss
- **Supersedes:** None
- **Related recommendations:** None
- **Related evidence:** `program/reference/rigor-tiers.md` ("the governance
  spine remains constant; only intensity changes")

## Context

The ideation/thinking framework needed a home. Candidates: a sibling template
repo alongside ADRP, a fresh master repo superseding ADRP, or evolving ADRP
itself. Separately, ideation and research risked living in two repo shapes,
forcing a migration when an idea grows into a program.

## Decision

Evolve this repository, in place, into the single unified master (Mycelium).
The idea lifecycle wraps around the existing research spine. An idea starts
thin and runs discovery/blueprint/charter/tracks in the same repo when it
earns rigor. Git history of the master is retained.

## Rationale

One repo to maintain (Robert's stated goal). No migration seam when an idea
graduates — changing houses mid-life is the ceremony that kills simmering
ideas. Rigor-tiers doc already promises spine-constant/intensity-varies. The
methodology's own evolution history stays auditable in the same git log.

## Consequences

ADRP's identity changes; v1 remains reachable via tags/history. The master
needs a home for its own artifacts (`framework/`, stripped at instance init).

## Alternatives Considered

Sibling thinking-template vendoring ADRP contracts (two masters, 70% shared
DNA, guaranteed drift). Fresh master with ADRP frozen as fallback (safe but
two repos during transition; rejected by Robert in favor of direct evolution).

## Risks

Churn lands in a template Robert relies on for real research programs.
Mitigation: instances are snapshots; nothing breaks retroactively; tag v1
before PHASE-01 begins.

## Revisit Triggers

If v2 churn corrupts the template's usefulness for pure research programs
before PHASE-05 stabilizes, split the sibling repo after all (the rejected
alternative becomes the escape hatch).

## Approval

Approved by Robert Guss in the Mycelium discovery interview, 2026-08-14
("yes I agree with A" — unified master; "yes I agree with A" — in place).
Recorded by the interviewing agent; the human's commit constitutes recording.
