# DEC-002 — The product of a thinking session is a durable, tiered, handoff-ready record

- **Status:** Accepted
- **Date:** 2026-08-14
- **Owner:** Robert Guss
- **Supersedes:** None
- **Related recommendations:** None
- **Related evidence:** karpathy llm-wiki gist (LLM-as-bookkeeper);
  mattpocock `handoff` skill (compaction ethos)

## Context

Ideation sessions previously ended in chat transcripts. Chat is not
authoritative (v1 rule), gets compacted or lost over weeks-long horizons, and
leaves nothing for a future session (human or agent) to re-enter through.

## Decision

Every session's output lands as git-tracked artifacts: decisions with
assumptions, disagreements with cruxes, glossary updates, log entries, open
questions with agreement states. Capture is tiered: a spark captures almost
nothing; rigor artifacts bind only at the tier that demands them. The
terminal artifact of a clarified idea is a handoff packet sufficient for a
fresh implementation agent with no chat history.

## Rationale

Sessions span hours to weeks; both parties must get back up to speed from
artifacts alone. "I feel clearer" with nothing written down was explicitly
rejected. Tiering keeps thirty-minute brainstorms from paying program-scale
ceremony.

## Consequences

The framework needs log/index conventions, a wake ritual, and tier-aware
conformance. Sessions carry a small always-on capture overhead.

## Alternatives Considered

Thinking-as-the-product with optional capture (rejected: violates
auditability requirement). Uniform full-rigor capture (rejected: kills
sparks).

## Risks

Capture discipline decays when it costs effort. Mitigation: generators make
the conforming path the lazy path; conformance checks catch omissions.

## Revisit Triggers

If session overhead measurably suppresses idea capture (fewer sparks written
down than pre-Mycelium), thin the tier-one requirements further.

## Approval

Approved by Robert Guss in the Mycelium discovery interview, 2026-08-14
(Q1 answer plus "everything should also be captured and auditable" message).
