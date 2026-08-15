# Independent Validation Gate

Every substantive artifact must pass a separate repository-level validation
before acceptance.

## Validator independence

Validation should be performed by a separate repository-aware agent or fresh
session. The validator reads: commissioning prompt, governing artifacts,
produced artifact, relevant upstream inputs, repository instructions.

## Validator scope

- Required sections
- Artifact metadata
- Identifier ranges and uniqueness
- Citation quality and portability
- Evidence Ledger completeness
- Required tables
- Completion checklist truthfulness
- Scope compliance
- Authority compliance
- Internal contradictions
- Placeholder remnants
- Allowed file scope
- Git diff
- Manifest status
- Downstream handoff completeness

## Mechanical vs substantive corrections

Validator may fix only mechanical issues (whitespace, heading hierarchy,
malformed fences, incorrect internal links, mechanical metadata typos).

Validator must **not** invent missing research, citations, findings,
recommendations, spike results, or architectural decisions.

Substantive defects → stage `requires-revision`.

## Standard transition

```text
Fresh stage session
  → Copy artifact to reserved path
  → awaiting-validation
  → Independent validation
  → Correct substantive defects if required
  → Human approval
  → Commit one coherent artifact
  → Record accepting commit in manifest
  → Unlock dependent stages
```

Templates:

- [`../templates/validation-report.md`](../templates/validation-report.md)
- [`../templates/validation-task.md`](../templates/validation-task.md)
