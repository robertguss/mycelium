# Program Blueprint — {{PROJECT_NAME}}

- **Artifact type:** Program Blueprint
- **Program:** {{PROJECT_NAME}}
- **Status:** Placeholder — not accepted
- **Version:** 0.0
- **Created:** {{CREATED_DATE}}
- **Last updated:** {{CREATED_DATE}}
- **Rigor tier (proposed):** standard — approve or change during discovery

> This file is a **section skeleton**. It does not prove stage completion.
> Fill after discovery framing is approved. Accept only after human approval.

## 1. Artifact Metadata

| Field | Value |
| ----- | ----- |
| Program ID | {{PROGRAM_ID}} |
| Owner | |
| Repository | |

## 2. Product or Project Vision

_To be filled after discovery._

## 3. Problem Statement

## 4. Intended Users and Stakeholders

## 5. Goals

## 6. Non-Goals

## 7. Locked Constraints

## 8. Success Criteria

## 9. Rigor Tier

- **Selected:** standard (proposed)
- **Rationale:**
- **Approval:** Pending

## 10. Research Graph

Focused research tracks are **not** pre-declared in this template. After
discovery, list stages here and add corresponding `[[stages]]` entries in
`research-program.toml`. Create prompt/report paths just-in-time from
`program/templates/`.

| Stage ID | Name | Kind | Depends on | Output | Parallel group |
| -------- | ---- | ---- | ---------- | ------ | -------------- |
| discovery | Project Discovery | discovery | — | this file | — |
| charter | Research Charter | research-charter | discovery | docs/01-research-charter.md | — |
| synthesis | Definitive Specification Synthesis | chief-architect-synthesis | _(all accepted research)_ | docs/specifications/01-definitive-specification.md | — |
| spec-review | Specification Adversarial Review | adversarial-review | synthesis | docs/reviews/01-specification-adversarial-review.md | — |
| spec-revision | Revised Definitive Specification | artifact-revision | spec-review | docs/specifications/02-definitive-specification-revised.md | — |
| implementation-plan | Implementation Plan | implementation-plan | spec-revision | docs/plans/01-implementation-plan.md | — |
| plan-review | Implementation Plan Adversarial Review | adversarial-review | implementation-plan | docs/reviews/02-implementation-plan-adversarial-review.md | — |
| plan-revision | Final Revised Implementation Plan | artifact-revision | plan-review | docs/plans/02-implementation-plan-revised.md | — |

## 11. Stage Descriptions and Dependencies

### Fixed spine

Document completion criteria for each fixed spine stage. Add focused research
stage sections as tracks are selected.

## 12. Parallelism

Sequential default. Justify any parallel groups explicitly.

## 13. Optional Replication Points

## 14. Artifact Inventory

See repository layout in root `README.md` and `program/operator/bootstrap.md`.

## 15. Identifier Allocations

| Namespace | Range | Notes |
| --------- | ----- | ----- |
| DEC | DEC-001..DEC-999 | Decision records |
| REC | _(allocate per research track)_ | |
| REQ | REQ-001..REQ-299 | Specification |
| FND (spec) | FND-001..FND-199 | Spec review |
| FND (plan) | FND-200..FND-399 | Plan review |
| RSK | RSK-001..RSK-999 | |
| OQ | OQ-001..OQ-999 | |
| SPK | SPK-001..SPK-999 | |
| PHASE | PHASE-01..PHASE-99 | |
| MS | MS-001..MS-999 | |

## 16. Authority and Precedence

See `program/contracts/authority-and-precedence.md`. Project-specific
overrides:

_None yet._

## 17. Human Approval Gates

See `program/operator/approval-gates.md`.

## 18. Fresh-Session Policy

Every substantive stage runs in a fresh session with an attachment manifest.

## 19. Validation and Commit Gates

Independent validation before acceptance. Human owns git commits.

## 20. Amendment Protocol

See `program/reference/amendment-protocol.md`.

## 21. Completion Criteria

See `program/operator/completion-criteria.md`.

## 22. Implementation Handoff Expectation

Final revised plan is delivery authority subordinate to revised specification.

## Completion Checklist

- [ ] Discovery framing approved by human
- [ ] All required sections filled (not placeholder prose)
- [ ] Research tracks justified; omitted tracks justified
- [ ] Identifier ranges allocated
- [ ] Rigor tier approved
- [ ] Human accepts Blueprint
- [ ] Manifest updated; accepting commit recorded
