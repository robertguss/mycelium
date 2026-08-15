# Definitive Specification Contract

The structure adapts to the project. A software-first specification should cover:

- Artifact metadata
- Executive decision summary
- Authority and intended use
- Problem and product definition
- Goals and non-goals
- Locked decisions and invariants
- Final technology stack
- System context
- Architecture
- Components and boundaries
- Data model
- Interfaces and integrations
- User workflows
- Security and privacy
- Reliability and operations
- Testing and verification
- CI and release
- Migration where applicable
- Performance expectations
- Internal contracts
- Dependency bill of materials
- Normative requirements
- Traceability
- Risk register
- Open questions
- Deferred work
- Rejected work
- Definition of done
- Handoff to adversarial review

The specification must be **standalone and implementation-ready**.

## Revised definitive specification

### Finding disposition

Every `FND-###` receives exactly one: Accepted; Accepted with modification;
Rejected; Deferred to a bounded evidence spike; Not applicable because another
correction removes the cause. No finding may disappear silently.

### Revision rules

- Integrate accepted corrections throughout affected sections.
- Remove superseded or contradictory language.
- Reconcile overlapping diffs.
- Preserve important strengths.
- Prefer simplification over new machinery.
- Retain stable requirement identifiers where the subject remains the same.
- Add new requirements only from unused identifiers.
- Remain standalone.

### Status

Use `Accepted — implementation authority` only when all Critical and
implementation-blocking findings are resolved or validly rejected and no known
blocking contradiction remains. Otherwise `Proposed — implementation blocked`
with blockers explicit.

### Revision-specific sections

- Revision Summary
- Finding Disposition Ledger
- Integrated Correction Ledger
- Preserved Strengths
- Updated implementation handoff
