# DEC-005 — Convention over configuration; containers-not-contents conformance

- **Status:** Accepted
- **Date:** 2026-08-14
- **Owner:** Robert Guss
- **Supersedes:** None
- **Related recommendations:** None
- **Related evidence:** Rails' convention-over-configuration doctrine;
  existing `just check`, ID namespaces, and `program/contracts/` as prior art

## Context

Robert named the organizing identity: a coding framework, but for thinking.
Components and packs may come and go; the core contract never changes.
Standardized structure serves the human (everything has a place) and the
agent (fewer degrees of freedom in generation means less drift and slop).

## Decision

Mycelium adopts convention over configuration as its identity:

1. Repo structure, file naming, and ID allocation are standardized in
   contracts and produced by deterministic generators (`just new-decision`,
   `just new-report`, `just new-council`), which allocate IDs and wire
   cross-references mechanically.
2. Instance evolution uses Rails-style migrations: versioned, ordered,
   deterministic transformations in `program/migrations/`, applied by
   `just upgrade` against the pinned methodology version.
3. A conformance suite validates structure: ID uniqueness/sequence, link
   integrity, front-matter schemas, legal lifecycle transitions,
   tier-appropriate artifact presence, crux presence on disagreement records,
   revisit triggers on simmering ideas. It runs locally and in CI.
4. Failures teach: every error names the violated convention and links its
   contract.
5. Hard boundary: checks validate containers, never contents. Thinking
   quality is judged by adversarial review, councils, and the human.

## Rationale

Bookkeeping is only automatable when structure is predictable (the llm-wiki
thesis). Generators kill the ID/cross-reference error class agents fumble.
Migrations solve live-instance evolution with a proven pattern. The
containers/contents boundary guards against Goodhart: agents writing to
please a linter.

## Consequences

Conventions must be documented once and enforced by tooling; convention
changes become migrations plus DEC records, not casual edits.

## Alternatives Considered

Prose-only conventions (agents skip prose; unenforceable). Content-aware
linting (Goodhart; rejected outright).

## Risks

Ceremony creep taxing sparks. Mitigation: tier-aware checks; a spark passes
nearly empty. Framework hobbyism displacing thinking. Mitigation: the
Laziness guard in blueprint success criteria.

## Revisit Triggers

If generators/migrations prove heavier to maintain than the errors they
prevent, subtract to the smallest enforcing set.

## Approval

Approved by Robert Guss in the Mycelium discovery interview, 2026-08-14
("conventions over configurations... think of it like a coding framework but
for thinking").
