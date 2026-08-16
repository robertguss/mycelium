# Council pack — Model report (RPT)

Pack contract for `RPT-###` files under `reviews/reports/`.
Authority: this file and `program/contracts/replication-reconciliation.md`.
Structure only (DEC-005). Check does not grade prose.

## Required front matter

| Field | Binding |
| --- | --- |
| `id` | `RPT-###` |
| `title` | string |
| `date` | `YYYY-MM-DD` |
| `model` | string; **not graded**; not an enum; empty FAIL |
| `commissioning` | `CMP-###` that resolves to an existing file |
| `rung` | must equal the CMP's `rung` |
| `adapter` | must equal the CMP's `adapter` |
| `prompt_sha256` | 64 lowercase hex; must equal the CMP-computed hash |

Check does **not** require distinct `model` strings across RPTs.

## Required H2s

Exact, case-sensitive:

```text
Position
Findings
Dissent
```

`## Dissent` must exist. Body may be `none`. Extra H2s are allowed.

## Prompt identity

```text
prompt_sha256 = hex(sha256(TrimSpace(SectionBody(body, "Prompt"))))
```

- `SectionBody` = bytes after the exact H2 `## Prompt` until the next H2 or EOF.
- `TrimSpace` = surrounding whitespace only.
- Hex = 64 **lowercase** `[0-9a-f]` characters. No `sha256:` prefix. Uppercase FAIL.
- First `## Prompt` wins. Extra `## Prompt` headings do not change the hash.
- Empty-after-trim is structurally legal (hash of empty string).

RPT front matter `prompt_sha256` must equal the CMP-computed hex. Mismatch →
FAIL.

## Teaching errors (hash mismatch)

```text
mycelium: RPT-001 prompt_sha256 mismatch (want ec87bfc2afd545807ca87b5c29cae8e77262cb3c746fc63e4539d8daeb2a77de)
convention: prompt-identity
contract: program/packs/council/contracts/report.md
fix: set prompt_sha256 to the sha256 hex of CMP-001 ## Prompt (trim surrounding whitespace)
```
