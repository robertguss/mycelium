# DEC-007 — Sparring stance: mandatory positions, agreement states, crux-bearing disagreement records

- **Status:** Accepted
- **Date:** 2026-08-14
- **Owner:** Robert Guss
- **Supersedes:** None
- **Related recommendations:** None
- **Related evidence:** mattpocock `grilling` skill (recommendation per
  question; decisions are the user's); pstack poteto-mode candor clause
  ("no is an acceptable answer"); `program/reference/discovery-protocol.md`
  (a clear recommendation on every substantive question)

## Context

Robert requires a thinking partner, not an order-taker: the model actively
proposes ideas and hypotheses, challenges terms, and disagrees when
warranted. Shared language matters; hidden assumptions and presuppositions
must be surfaced. Disagreement is where the good stuff lives and must be
captured, including agree-to-disagree outcomes.

## Decision

The thinking-mode skill enforces a sparring stance:

1. Every substantive question carries the agent's position or
   recommendation. Bare questions are a smell.
2. Every substantive question carries an agreement state: `open`, `aligned`,
   or `agree-to-disagree`. The third is terminal and honorable.
3. Disagreement records capture both positions, both sets of reasons, and
   the crux: what evidence would change each mind. Cruxes are eligible to
   become research stages or spikes.
4. Dissent is retained forever, including the agent's after being overruled.
5. The assumption audit is a standing move: both parties periodically dump
   what they're taking for granted, and the agent infers and reads back the
   human's unstated presuppositions.
6. Glossary discipline per domain-modeling conventions: terms challenged on
   drift, sharpened on vagueness, recorded on resolution.

## Rationale

Mandatory positions convert an interviewer into a partner. Cruxes convert
disagreement from a museum piece into the next experiment. Retained dissent
is the audit trail's highest-signal content (per the replication contract's
own ethos).

## Consequences

Session transcripts and artifacts gain structure conformance can check
(crux presence, agreement-state fields). The mode skill must encode the
stance so it survives model and runtime changes.

## Alternatives Considered

Pure grilling (extracts but does not propose). Council-only perspective
(episodic, expensive; sparring must be free and continuous).

## Risks

Performative disagreement (manufactured dissent to satisfy the stance).
Mitigation: positions require reasons; the containers-not-contents boundary
keeps enforcement structural while humans judge substance.

## Revisit Triggers

If crux capture degenerates into boilerplate, tighten the disagreement-record
contract rather than dropping the field.

## Approval

Approved by Robert Guss in the Mycelium discovery interview, 2026-08-14
("I want to capture disagreements and their reasons... this is actually
gold").
