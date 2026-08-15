# Naming and ID-to-Path (Mycelium 2.0)

Naming, IDs, and layout are derivable, not searchable. The ID-to-path rule is a
pure function both directions.

## Identifier grammar

- ParseID = `^([A-Z]+)-([0-9]+)$`
- Format zero-pads to the namespace digit width (`PHASE` = 2, all others = 3).
- Allocation allows gaps. Next ID = `max(N)+1` over existing files in the home.
- No tombstones. Never overwrite an existing path.

## Path rule

```text
path = <home>/<NS>-<zero-padded-digits>-<kebab-slug>.md
```

Example: `DEC-014` with slug `manifest-name` → `decisions/DEC-014-manifest-name.md`.

## Link resolution

References matching

```text
\b(DEC|ASM|EVD|SPK|FND|REC|REQ|OQ|RSK|PHASE|MS)-[0-9]+\b
```

must resolve to exactly one matching file in that namespace's home directory.

## Registered types

| Type key | NS | Home | Filename pattern | Digits | Stage-scoped |
| --- | --- | --- | --- | --- | --- |
| decision | `DEC` | `decisions/` | `DEC-###-slug.md` | 3 | no |
| assumption | `ASM` | `assumptions/` | `ASM-###-slug.md` | 3 | no |
| evidence | `EVD` | `evidence/` | `EVD-###-slug.md` | 3 | no |
| spike | `SPK` | `spikes/` | `SPK-###-slug.md` | 3 | no |
| finding | `FND` | `findings/` | `FND-###-slug.md` | 3 | yes |
| recommendation | `REC` | `recommendations/` | `REC-###-slug.md` | 3 | yes |
| requirement | `REQ` | `requirements/` | `REQ-###-slug.md` | 3 | yes |
| question | `OQ` | `questions/` | `OQ-###-slug.md` | 3 | no |
| risk | `RSK` | `risks/` | `RSK-###-slug.md` | 3 | no |
| phase | `PHASE` | `phases/` | `PHASE-##-slug.md` | 2 | no |
| milestone | `MS` | `milestones/` | `MS-###-slug.md` | 3 | no |

`PHASE` is **not** stage-scoped.

Directories are plural lowercase nouns, one artifact type each. IDs are
`UPPER-###` (or `UPPER-##` for `PHASE`). Slugs are kebab-case.

## Conformance

Check enforces the mapping both directions: files must match their pattern;
references must resolve. See `program/contracts/conformance.md` and
`program/contracts/identifiers.md`.
