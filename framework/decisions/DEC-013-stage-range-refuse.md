# DEC-013 — Refuse allocation outside a declared stage-scoped range

- **Status:** Accepted
- **Date:** 2026-08-14
- **Owner:** Robert Guss
- **Supersedes:** None (settles blueprint OQ-007)
- **Related recommendations:** None
- **Related evidence:** program/contracts/identifiers.md (v1 ranges)

## Context

FND, REC, and REQ are stage-scoped. The blueprint allocated ranges per
stage so identifiers do not collide across reviews. OQ-007 asked whether
the generator should warn or refuse when the next ID falls outside a
declared range, or when no range is declared.

## Decision

1. Stage-scoped types (finding, recommendation, requirement) require a
   declared range in the instance manifest `[identifiers]` table before
   allocation.
2. If no range is declared, the generator REFUSES.
3. If the next ID would fall outside all declared ranges, the generator
   REFUSES.
4. A warning is not sufficient. Check also fails existing files whose
   IDs sit outside every declared range for that key.
5. Non-stage-scoped types (DEC, ASM, EVD, SPK, OQ, RSK, PHASE, MS) do
   not require a range.

## Rationale

A warning is a log line agents ignore. A refuse is a teaching error with
a contract link. Ranges exist to keep stage traces disjoint; silent
overflow defeats them.

## Consequences

Fixture CI must write ranges before generating FND/REC/REQ. There is no
`mycelium range` command in PHASE-01; tests edit the manifest.

## Alternatives Considered

Warn and allocate (rejected: unenforceable). Auto-extend the range
(rejected: hides the stage boundary).

## Risks

Sparks that want a finding before declaring a range are blocked.
Mitigation: focused sparks do not bind findings; declaring a range is
one manifest edit.

## Revisit Triggers

If a real idea needs multiple disjoint ranges per type in PHASE-01,
reopen to allow an array of ranges. This brief allows one range per key.

## Approval

Settled by Arvo 2026-08-14 (OQ-007 → REFUSE). Recorded by Architect in the
PHASE-01 implementation brief.
