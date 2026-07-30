# Stable Identifier System

Every substantive item receives a stable identifier.

## Required namespaces

| ID | Meaning |
| -- | ------- |
| `DEC-###` | Accepted decision records |
| `REC-###` | Research recommendations |
| `REQ-###` | Normative specification requirements |
| `FND-###` | Adversarial review findings |
| `RSK-###` | Risks |
| `OQ-###` | Open questions |
| `SPK-###` | Evidence spikes |
| `PHASE-##` | Implementation phases |
| `MS-###` | Implementation milestones |

Optional: `EVD-###` (Evidence Ledger entries), `ASM-###` (assumptions).

## Allocation

The Program Blueprint and manifest allocate non-overlapping ranges by stage
before use. Example:

```text
Focused Research 1: REC-001..REC-099
Focused Research 2: REC-100..REC-199
Specification: REQ-001..REQ-399
Specification Review: FND-001..FND-199
Plan Review: FND-200..FND-399
```

## Stability rules

- Never reuse an identifier for a different subject.
- Preserve identifiers when a subject is modified.
- Mark deleted items superseded or rejected; do not silently remove history.
- Later stages must disposition every material upstream identifier in scope.
- Findings remain traceability items; they do not automatically become
  requirements.
