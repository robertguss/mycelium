# DEC-008 — Perspective ladder; councils opt-in, engine-agnostic, Cursor adapter first

- **Status:** Accepted
- **Date:** 2026-08-14
- **Owner:** Robert Guss
- **Supersedes:** None
- **Related recommendations:** None
- **Related evidence:** `program/contracts/replication-reconciliation.md`
  ("optional but first-class"; no majority vote; dissent retained);
  karpathy/llm-council (three-stage shape); pstack (second-opinion doctrine,
  model-role config)

## Context

Multi-model perspective on ideas is a core requirement, but councils are
expensive, not every idea needs them, Robert's council-capable runtime today
is Cursor (subscription-covered fan-out), and tooling churns weekly to
monthly. Every council member's full output must be captured for later
synthesis.

## Decision

1. A three-rung perspective ladder: **sparring** (free, always on),
   **second opinion** (one different model, same commissioning prompt,
   cheap), **council** (full multi-model replication + reconciliation,
   expensive, rare).
2. Councils are opt-in, never a required stage, never auto-run. The agent
   suggests one only when the replication contract's triggers fire, stating
   panel size and cost class first.
3. Councils reuse v1's replication and reconciliation contracts verbatim,
   with replicas on different models. Reports land per-model as committed
   artifacts; reconciliation retains dissent; majority vote, prose
   confidence, and model reputation remain banned selectors.
4. Contracts are engine-agnostic: they define commissioning-prompt and
   report file shapes only. Producers are swappable adapters. Cursor
   parallel subagents are the first adapter; the manual floor (paste into N
   chat UIs, save N files) satisfies the contract with zero tooling.
5. No portable council CLI in v1 (an OpenRouter-based CLI was proposed and
   withdrawn as hypothetical). `AGENTS.md` carries a capability note so
   non-fan-out runtimes skip rungs 2–3 gracefully.
6. Panel presets (quick / standard / high-stakes) live in user-level config.

## Rationale

Opt-in was already v1 doctrine. The stable half of a churning tool landscape
is the paper trail, so standardize files, not machinery. Anonymized
cross-review with retained dissent beats llm-council's chairman-smoothing,
so v1's reconciliation rules govern.

## Consequences

Council value depends on runtime capability; costs stay visible; the ladder's
middle rung likely becomes the workhorse.

## Alternatives Considered

Councils as a default stage (expensive, rejected). Portable CLI engine in v1
(no current need outside Cursor; rejected on laziness grounds).

## Risks

Adapter rot as runtimes change. Mitigation: adapters are thin; contracts are
the ground truth; the manual floor always works.

## Revisit Triggers

Build the portable CLI when a real non-Cursor council need occurs (the
recorded crux for the withdrawn proposal).

## Approval

Approved by Robert Guss in the Mycelium discovery interview, 2026-08-14
("councils should be opt-in and not standard nor required... it's just a
nice thing to have when I need it").
