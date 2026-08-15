# DEC-006 — Idea lifecycle with simmer as a first-class state and a wake ritual

- **Status:** Accepted
- **Date:** 2026-08-14
- **Owner:** Robert Guss
- **Supersedes:** None
- **Related recommendations:** None
- **Related evidence:** `program/contracts/evidence-model.md` (revalidation
  triggers); karpathy llm-wiki gist (decisions as living objects / Vigil)

## Context

Some ideas take a long time to simmer; that space is part of the process and
must be accounted for structurally. Sessions span hours, days, and weeks.
Nothing in the surveyed landscape modeled deliberate parking: wayfinder's fog
means "can't specify yet," and blocked means "waiting on a decision."

## Decision

The idea lifecycle is an explicit state machine stored in the instance
manifest:

```text
spark → exploring ⇄ simmering → clarified → handed-off
any state → archived
```

Simmering means "could decide now, choosing not to" and requires a revisit
trigger (date or event). Waking is a ritual, not a reload: reread index and
log tail, check evidence revalidation triggers and assumption records against
what changed while parked, then brief the human. Conformance validates legal
transitions.

## Rationale

Scattered "we'll come back to this" notes are unenforceable. An explicit
state with a required trigger makes simmering queryable (portfolio), safe
(nothing parked is lost), and warm on return (the wake brief restores context
for human and agent).

## Consequences

Manifest gains lifecycle fields. The portfolio skill and conformance suite
consume them. Log and index conventions must support the wake brief.

## Alternatives Considered

Simmer as a prose convention (unenforceable). Simmer as tracker state
(rejected with hub/tracker-first storage in DEC-003).

## Risks

Revisit triggers pile up unactioned. Mitigation: portfolio surfaces overdue
wakes by default.

## Revisit Triggers

If wake briefs prove low-value in practice (Robert rereads raw logs anyway),
redesign the brief's contract rather than abandoning the state.

## Approval

Approved by Robert Guss in the Mycelium discovery interview, 2026-08-14
("some ideas take a long time to simmer and baking in that... is critical").
