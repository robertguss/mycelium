# Evidence Spike Protocol

Evidence spikes are first-class research artifacts.

## When to use a spike

- Documentary evidence is weak or contradictory.
- A claim is economically testable.
- A technical behavior is load-bearing.
- Platform or filesystem semantics matter.
- Agent behavior is uncertain.
- A benchmark could change architecture.
- An API, migration, or tool assumption needs direct verification.
- Library ergonomics or failure behavior cannot be assessed from docs alone.

## Constraints

A spike must be:

- Bounded
- Decision-oriented
- Disposable
- Reproducible where practical
- Explicit about environment and limitations
- Kept outside the research repository unless the spike **report** is committed

**Prototype code must not silently become production architecture.**

## Mandatory spike review

The consuming report must not overgeneralize from one OS, dataset, hardware
profile, agent run, dependency version, or network condition.

## Record format

Use [`../templates/evidence-spike.md`](../templates/evidence-spike.md).
