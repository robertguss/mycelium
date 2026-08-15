# DEC-003 — One repository per idea; no hub repo; portfolio queried live

- **Status:** Accepted
- **Date:** 2026-08-14
- **Owner:** Robert Guss
- **Supersedes:** None
- **Related recommendations:** None
- **Related evidence:** None

## Context

Ideas need homes. A single garden/wiki repo (one directory per idea) was the
agent's initial recommendation, prioritizing cross-idea trails and one-stop
re-entry. Robert rejected it: each idea requires its own research and
artifacts, and a shared repo goes noisy, messy, and disorganized.

## Decision

Every idea gets its own repository, instantiated thin from the master
template. No hub or garden repository exists. Cross-idea state lives in
GitHub metadata: an `idea` topic per repo plus lifecycle fields in each
instance manifest. A portfolio skill answers "what's simmering / what's due
to wake" by querying live (`gh` + manifest reads). Instances clone under one
predictable root (`~/ideas/<slug>`).

## Rationale

Per-idea isolation matches how research artifacts accumulate, keeps agent
sessions context-isolated by directory boundary (the fresh-session rule
enforced structurally), and fits the existing template-instantiation
workflow Robert already trusts.

## Consequences

Repo count grows with ideas. Cross-idea discovery depends on the portfolio
skill and topic hygiene (scaffold applies the topic automatically).

## Alternatives Considered

Single garden repo (agent's original recommendation; rejected for noise).
Ideas inside related project repos (orphans pre-project sparks). Tracker-first
storage (weeks-long prose thinking reads terribly as issue threads).

## Risks

Dead-spark repo litter. Mitigation: `archived` state plus portfolio filtering;
archiving is a state change, not a deletion.

## Revisit Triggers

Agent's retained dissent, recorded as the crux: if live portfolio querying
proves flaky or slow enough that Robert stops trusting the cross-idea view,
revisit a thin generated (never hand-maintained) registry artifact.

## Approval

Approved by Robert Guss in the Mycelium discovery interview, 2026-08-14
("B only because each idea requires its own research and artifacts").
