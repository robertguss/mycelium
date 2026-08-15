# Stable Identifier System

Every substantive item receives a stable identifier.

## Required namespaces

| ID | Meaning |
| --- | --- |
| `DEC-###` | Accepted decision records |
| `ASM-###` | Assumptions (first-class; falsifiable) |
| `EVD-###` | Evidence Ledger entries (first-class) |
| `SPK-###` | Evidence spikes |
| `FND-###` | Adversarial review findings (stage-scoped) |
| `REC-###` | Research recommendations (stage-scoped) |
| `REQ-###` | Normative specification requirements (stage-scoped) |
| `OQ-###` | Open questions |
| `RSK-###` | Risks |
| `PHASE-##` | Implementation phases |
| `MS-###` | Implementation milestones |
| `CMP-###` | Commissioning prompts (perspective ladder; not stage-scoped) |
| `RPT-###` | Model reports (perspective ladder; not stage-scoped) |
| `RCL-###` | Reconciliations (perspective ladder; not stage-scoped) |

`ASM` and `EVD` are first-class required namespaces (not optional).
`CMP`, `RPT`, and `RCL` are pack-owned (council pack); not stage-scoped.

## Allocation

The Program Blueprint and manifest allocate non-overlapping ranges by stage
before use. Example (v1 ranges; retained as examples):

```text
Focused Research 1: REC-001..REC-099
Focused Research 2: REC-100..REC-199
Specification: REQ-001..REQ-399
Specification Review: FND-001..FND-199
Plan Review: FND-200..FND-399
```

### Stage-scoped ranges (Mycelium 2.0)

`FND`, `REC`, and `REQ` are stage-scoped. Before allocation, the instance
manifest must declare a range under `[identifiers]`
(`findings` / `recommendations` / `requirements`).

- If no range is declared, the generator **REFUSES**.
- If the next ID would fall outside all declared ranges, the generator
  **REFUSES**.
- A warning is not sufficient.
- Check fails existing files whose IDs sit outside every declared range for
  that key.

See DEC-013.
Non-stage-scoped types do not require a range.

## Stability rules

- Never reuse an identifier for a different subject.
- Preserve identifiers when a subject is modified.
- Mark deleted items superseded or rejected; do not silently remove history.
- Later stages must disposition every material upstream identifier in scope.
- Findings remain traceability items; they do not automatically become
  requirements.
- Allow gaps; next = max(N)+1; no tombstones; never overwrite.
