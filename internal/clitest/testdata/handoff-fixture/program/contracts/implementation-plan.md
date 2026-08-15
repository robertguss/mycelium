# Implementation Plan Contract

## Purpose

Translate the accepted revised specification into a safe delivery sequence. The
plan defines **how to build**, not what the architecture should become.

## Plan boundary

Stop at: phases, milestones, dependencies, integration points, evidence spikes,
dogfooding, rollback/reconsideration boundaries, executable phase acceptance
criteria.

**Do not** create a granular execution backlog or coding-agent task packets.

## Required content

- Artifact metadata
- Implementation authority
- Objectives and non-goals
- Assumptions
- Dependency graph
- Phase overview and one section per phase
- Milestones
- Cross-phase integration
- Data or migration sequencing where applicable
- Testing strategy by phase
- Security activities by phase
- Operations and release readiness
- Dogfooding
- Risk register
- Open questions
- Rollback and reconsideration triggers
- Requirement-to-phase traceability
- Definition of plan completion

Templates:

- [`../templates/phase.md`](../templates/phase.md)
- [`../templates/milestone.md`](../templates/milestone.md)

## Sequencing principles

- Resolve unknowns before they harden into architecture.
- Produce thin end-to-end capability early.
- Integrate continuously; dogfood before broad feature expansion.
- Keep risky decisions reversible.
- Avoid phases that are only horizontal infrastructure with no usable outcome.
- Avoid circular entry/exit criteria and acceptance criteria that depend on
  later phases.

## Final revised plan

Disposition every plan-review finding; preserve accepted revised specification;
integrate corrections; remove circular sequencing; make entry/exit executable;
preserve early integration and dogfooding; state blockers honestly; remain at
phase/milestone granularity.

Status `Accepted — delivery authority` only when no implementation-blocking plan
finding remains.
