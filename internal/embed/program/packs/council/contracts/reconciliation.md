# Council pack — Reconciliation (RCL)

Thin pointer. Obey
[`program/contracts/replication-reconciliation.md`](../../../contracts/replication-reconciliation.md)
**verbatim**. Do not fork that file. Do not copy the v1 body here.

Authority: `framework/phases/PHASE-04-implementation-brief.md` §6.4, §6.6, §6.7.
Pack path for `RCL-###` files: `reviews/reconciliations/`.

## Required front matter

| Field | Binding |
| --- | --- |
| `id` | `RCL-###` |
| `title` | string |
| `date` | `YYYY-MM-DD` |
| `commissioning` | `CMP-###` that resolves |
| `rung` | `council` **only**. `second-opinion` on an RCL → FAIL. |

Second-opinion does **not** require an RCL. An RCL whose CMP is
`rung = "second-opinion"` → FAIL.

## Required H2s (v1 list)

Exact strings from the v1 contract, including **`Retained dissent`**:

```text
Convergence
Material disagreement
Evidence unique to one report
Contradictory evidence
Different assumptions
Different scope interpretations
Recommendations independently supported
Questions requiring another spike
Final reconciled recommendation
Retained dissent
```

No majority vote / reputation / prose-confidence selectors — stated in the v1
contract. Check does **not** content-score the method (DEC-005).

## SEED-DISSENT substring rule

If **any** matching RPT's `## Dissent` section body contains the exact token
`SEED-DISSENT`, the matching RCL's `## Retained dissent` section body **must**
contain `SEED-DISSENT`. Substring presence check, not a quality score.

| Situation | Check |
| --- | --- |
| No RPT contains `SEED-DISSENT` | Rule does not fire |
| RPT contains `SEED-DISSENT`, matching RCL retains it | PASS |
| RPT contains `SEED-DISSENT`, matching RCL lacks it | FAIL |
| `second-opinion` RPT contains `SEED-DISSENT` (no RCL required) | Rule does not require an RCL |

## Teaching errors

```text
mycelium: RCL-001 ## Retained dissent missing SEED-DISSENT
convention: seeded-dissent
contract: program/packs/council/contracts/reconciliation.md
fix: retain the SEED-DISSENT token in ## Retained dissent
```
