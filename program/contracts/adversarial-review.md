# Adversarial Review Contract

Adversarial review is a **separate discipline**. The author of a specification
or plan should not be assumed to have found its most dangerous flaws.

## Purpose

> What is wrong, contradictory, unsafe, non-total, under-specified,
> over-engineered, unprovable, difficult to implement, difficult to test, or
> likely to fail in real use?

## Reviewer posture

- Attack polished sections.
- Trace workflows end to end.
- Check failure, cancellation, rollback, and cleanup.
- Check consistency across prose, requirements, tables, examples, appendices.
- Attempt to delete unnecessary machinery.
- Avoid adding features.
- Avoid treating preference as defect.
- Produce a small number of strong findings rather than meeting a quota.

## Severity

| Level    | Meaning                                                     |
| -------- | ----------------------------------------------------------- |
| Critical | Blocks all implementation or risks catastrophic harm        |
| High     | Blocks the affected phase or creates major invalid behavior |
| Medium   | Must be fixed before the affected phase completes           |
| Low      | Should be corrected in revision; does not block early work  |

Finding template: [`../templates/finding.md`](../templates/finding.md).

## Software-first review scope

Product scope, requirements, architecture, data and interfaces, user workflows,
security, filesystem behavior, determinism, dependencies, testing, CI,
operations, migration, implementation phases, agent legibility, framework creep,
acceptance criteria.

## Implementation-plan review attacks

Circular phase dependencies; missing prerequisites; overlarge phases;
milestones without integration evidence; acceptance criteria that do not prove
outcomes; late risk discovery; delayed dogfooding; deferred integration;
migration/rollback/security gaps; test environments not ready; phase boundaries
that conflict with the specification; plan steps that reinterpret architecture;
excessive parallel work; order that hardens wrong decisions before evidence.

Use a **separate `FND-###` range** from specification review.

## Risk-triggered additional rounds

Permitted when: major restructuring; multiple Critical/High findings accepted;
new machinery introduced; blocking contradictions remain; reviewer recommends
another pass; spike materially changes architecture. Not automatic; not an
endless loop. Second review focuses on machinery introduced or changed.
