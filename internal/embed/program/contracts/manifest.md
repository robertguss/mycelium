# Instance Manifest — mycelium.toml (Mycelium 2.0)

Scaffolded instances use `mycelium.toml` as the sole manifest filename
(DEC-012).
The master repository is the mycelium product repo and does not carry
`research-program.toml`; idea instances use `mycelium.toml` (DEC-012).
Runtime commands detect an instance by the presence of `mycelium.toml`.

`schema_version = 1` is required.

## Required fields

| Field | Rule |
| --- | --- |
| `schema_version` | Must be `1` |
| `idea_name` | Human-readable title |
| `slug` | Must equal `slugify(idea_name)` |
| `state` | One of: `spark`, `exploring`, `simmering`, `clarified`, `handed-off`, `archived`. Birth = `spark` |
| `tier` | One of: `focused`, `standard`, `high-assurance`. Birth = `focused` |
| `methodology_version` | Methodology pin (scaffold writes `2.0.0`) |
| `generated_by_cli_version` | CLI version that scaffolded the instance |
| `created_date` | Set once at scaffold; never rewritten |
| `updated_date` | Bumped on every mutating operation |
| `revisit` | Empty string except **required non-empty** when `state = simmering` |
| `github_repo` | Empty until publish; then `owner/name` |

State and tier are orthogonal. `methodology_version` and
`generated_by_cli_version` are orthogonal (DEC-011).

## Optional `[identifiers]`

Allowed keys only:

| Key | Grammar |
| --- | --- |
| `findings` | `^([A-Z]+)-([0-9]+)\.\.([A-Z]+)-([0-9]+)$` |
| `recommendations` | same |
| `requirements` | same |

Range rules (brief §5):

- Range grammar: `^([A-Z]+)-([0-9]+)\.\.([A-Z]+)-([0-9]+)$`
- Both namespaces must match the key's NS (`findings` → FND, `recommendations` → REC, `requirements` → REQ)
- Start integer <= end integer
- One range per key this phase
- Unknown identifier keys: refuse

Stage-scoped allocation (FND / REC / REQ) requires a declared range
(DEC-013).

## Optional `[[deviations]]`

Each entry requires exactly:

- `convention` — name of the convention being waived
- `reason` — why

Extra deviation keys → refuse.

Example deviation for an allowed extra top-level path:
`convention = "extra-top-level:<path>"`.

## Refuse on unknown keys

- Unknown top-level keys → refuse.
- Unknown `[identifiers]` keys → refuse.
- Extra keys inside a `[[deviations]]` entry → refuse.
