# Instance Manifest — mycelium.toml (frozen G1 / pre-github_repo)

Scaffolded instances use `mycelium.toml` as the sole manifest filename
(DEC-012). This frozen copy omits `github_repo` from required fields
(DEC-011 / PHASE-05 G1).

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

State and tier are orthogonal. `methodology_version` and
`generated_by_cli_version` are orthogonal (DEC-011).

## Optional `[identifiers]`

Allowed keys only:

| Key | Grammar |
| --- | --- |
| `findings` | `^([A-Z]+)-([0-9]+)\.\.([A-Z]+)-([0-9]+)$` |
| `recommendations` | same |
| `requirements` | same |

## Optional `[[deviations]]`

Each entry requires exactly:

- `convention` — name of the convention being waived
- `reason` — why

Extra deviation keys → refuse.

## Refuse on unknown keys

- Unknown top-level keys → refuse.
- Unknown `[identifiers]` keys → refuse.
- Extra keys inside a `[[deviations]]` entry → refuse.
