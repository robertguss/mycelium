# Council pack — Commissioning (CMP)

Pack contract for `CMP-###` files under `reviews/commissioning/`.
Authority: this file and `program/contracts/replication-reconciliation.md`.
Structure only (DEC-005). Check does not grade prose.

## Required front matter

| Field | Binding |
| --- | --- |
| `id` | `CMP-###` |
| `title` | string |
| `date` | `YYYY-MM-DD` |
| `rung` | `second-opinion` \| `council` |
| `opt_in` | TOML boolean. **Must be `true`.** Missing or `false` → FAIL. String `"true"` → FAIL. |
| `cost_class` | IFF `rung == "second-opinion"`: must be `cheap`. IFF `rung == "council"`: must be `quick` \| `standard` \| `high-stakes`. Cross-rung values FAIL. |
| `adapter` | `cursor` \| `manual` |

`prompt_sha256` is **not** required on the CMP. Check computes the hash at
check time from `## Prompt`; every matching RPT must repeat the same hex
(see `program/packs/council/contracts/report.md`).

## Required H2s

Exact, case-sensitive:

```text
Prompt
Attachments
Cost
```

Bodies are not graded. `Attachments` may be `none`. `Cost` restates the class
in prose; check does not parse it. Extra H2s are allowed.

## Cardinality (WIP rule)

Cardinality binds IFF any matching RPT or RCL exists for that CMP. A CMP alone
(commissioned, not yet run) passes, provided `opt_in`, `cost_class`, `rung`,
and `adapter` are legal.

| CMP `rung` | Matching files | `mycelium check` |
| --- | --- | --- |
| `second-opinion` | no RPT, no RCL | PASS |
| `second-opinion` | exactly one RPT, no RCL | PASS |
| `second-opinion` | two or more RPTs | FAIL |
| `second-opinion` | any RCL | FAIL |
| `council` | no RPT, no RCL | PASS |
| `council` | ≥1 RPT or any RCL, but not (≥2 RPTs **and** exactly one RCL) | FAIL |
| `council` | ≥2 RPTs and exactly one RCL | PASS |
| `council` | ≥2 RPTs and two or more RCLs | FAIL |

"Matching" means the file's `commissioning` field equals the CMP id.

## Teaching errors (paths that point here)

Four lines, stderr, exit 1. Examples that name this contract:

```text
mycelium: CMP-001 opt_in must be true (got false)
convention: council-opt-in
contract: program/packs/council/contracts/commissioning.md
fix: set opt_in = true, or delete the commissioning file

mycelium: CMP-001 cost_class "cheap" is not quick|standard|high-stakes when rung=council
convention: council-cost-class
contract: program/packs/council/contracts/commissioning.md
fix: set cost_class to quick, standard, or high-stakes

mycelium: CMP-001 council requires >=2 model reports and exactly one reconciliation
convention: council-cardinality
contract: program/packs/council/contracts/commissioning.md
fix: add matching RPT-### files (>=2) and exactly one RCL-###, or remove started reports
```
