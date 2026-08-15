# Conformance Suite (Mycelium 2.0)

Structure only (DEC-005).
Checks never grade prose or thinking quality.

## Must-implement checks (23)

1. Manifest present (`mycelium.toml`) and parses; required fields valid per
   `program/contracts/manifest.md`.
2. Lifecycle state legality per `program/contracts/lifecycle.md` **PHASE-02
   storage rules** (`clarified` legal; `handed-off` still FAIL; `simmering`
   requires revisit grammar).
3. Tier legality and tier-aware artifact binds.
4. ID uniqueness within each namespace home.
5. ID-to-path integrity both directions (`program/contracts/naming.md`).
6. Link resolution for ID references (also scans `index.md` +
   `briefs/*.md`).
7. Required front matter and sections per sidecar schema.
8. Log line prefixes parseable. Ops:
   `scaffold|new|tier|publish|check|state|wake|supersede`. Regex:
   `^\d{4}-\d{2}-\d{2}\t(scaffold|new|tier|publish|check|state|wake|supersede)\t(\S+)\t`
9. Interrupted operation: leftover journal or stale lock → teaching recovery
   (complete or `--abort-journal`).
10. Undeclared extra top-level paths unless deviation
    `extra-top-level:<path>` is declared.
11. Stage-scoped IDs outside every declared range → FAIL
    (DEC-013).
12. `index.md` present + required H2s (`State`, `Artifacts`, `Log tail`,
    `Wake`) — all tiers.
13. `briefs/` is an allowed top-level path (not an ID-to-path home).
14. If the log contains a `wake` op, `briefs/LATEST.md` must exist and pass
    the five wake H2s (`program/contracts/wake.md`).
15. Agreement-conditional OQ headings per `program/contracts/sparring.md`.
    Invalid `agreement` → FAIL.
16. If `CONTEXT.md` exists: H1 `# Glossary`; any H2 ⇒ H3 `Definition`. Empty
    glossary legal.
17. If a DEC contains `## Dissent`, the section must contain at least one
    resolvable `OQ-###` or `ASM-###`. Heading absent → pass.
18. Pack presence-is-registration: load `program/packs/<name>/` when the
    directory exists. Namespace, type-key, or home collision between packs
    (or pack vs core) → FAIL.
19. `reviews/` is an allowed top-level path only when
    `program/packs/council/` is present. Without the pack and without
    `reviews/`, core check still passes. Without the pack but with
    `reviews/` leftover: FAIL extra-top-level unless deviation
    `extra-top-level:reviews/`.
20. When the pack is present and CMP files exist: required front matter +
    H2s per `program/packs/council/contracts/commissioning.md`; `opt_in`
    must be `true`; `cost_class` IFF `rung`; `adapter` enum.
21. When the pack is present and RPT files exist: required front matter +
    H2s including `Dissent`; `commissioning` resolves; `rung`/`adapter`
    match the CMP; `prompt_sha256` equals the check-computed CMP hash.
22. Rung cardinality per pack commissioning contract. RCL required H2s
    including `Retained dissent`. `SEED-DISSENT` substring rule per pack
    reconciliation contract. RCL `rung` is `council` only.
23. Bidirectional IFF + one-to-one. If `status = "Superseded"`:
    `superseded_by` present, same namespace, resolves, and peer `supersedes`
    equals this ID. If `supersedes` is set: peer exists, peer
    `status = "Superseded"`, peer `superseded_by` equals this ID. At most
    one inbound `superseded_by` per NEW. Binds in Slice 2.

## Lift timing

Items 1–14 stay as written (PHASE-02). Items 15–17 landed by PHASE-03 slice:

| Slice | Check behavior |
| --- | --- |
| PHASE-03 Slice 1 | Schema `required_sections` drops `Crux`. Schema-driven check stops requiring Crux on every OQ. **No** IFF bind yet. |
| PHASE-03 Slice 2 | Check calls `internal/sparring` for each `questions/OQ-*.md`. IFF rules bind (item 15). |
| PHASE-03 Slice 3 | `CONTEXT.md` glossary rules bind (item 16). |
| PHASE-03 Slice 4 | Optional DEC Dissent rule binds (item 17). |

Items 18–22 land by PHASE-04 slice:

| Slice | Check behavior |
| --- | --- |
| PHASE-04 Slice 1 | Pack presence + collision + `reviews/` extra-top-level (items 18–19). **No** CMP/RPT/RCL content checks yet. |
| PHASE-04 Slice 2 | Pack schemas registered. Item 7 applies to pack types. `opt_in` must be *present* but "must be `true`" and cost-class IFF and hash and cardinality do **not** bind yet. `mycelium new` discovers pack types. |
| PHASE-04 Slice 3 | Check calls `internal/ladder`. Items 20–22 IFF / hash / cardinality / `SEED-DISSENT` bind. |
| PHASE-04 Slice 4 | Skills + adapters; no new check rule. |
| PHASE-04 Slice 5 | MS-401 matrix fixtures in `go test ./...` **are** the gate. |

PHASE-05 slices 1–5 (item 23 + related check/status deltas):

| Slice | Check / status behavior |
| --- | --- |
| PHASE-05 Slice 1 | Schema enum + optional keys land. Parsers exist. **No** CLI. **No** item 23 bind. Item 8 regex **not** yet changed (a hand-written `supersede` log line would still fail check — do not write one in Slice 1 fixtures). |
| PHASE-05 Slice 2 | Command + item 8 + item 23 + item 6 link bind. Happy DEC pair + refuse table. **Bound.** |
| PHASE-05 Slice 3 | Item 1 G1 rule. `status` / `status --all` tolerance. G0–G3. |
| PHASE-05 Slice 4 | No new check item. CHANGELOG + release script + checksum tests. |
| PHASE-05 Slice 5 | MS-501 matrix harness in `internal/clitest` runs every PHASE-05 acceptance row. |

Do not require a wake brief on instances that never simmered. Do not grade
brief prose. Do not require N artifacts. Do not add `reviews/` to the
always-allowed top-level list this slice (Slice 1 binds that).

## What check must not do

| Temptation | Verdict |
| --- | --- |
| Require any OQ on spark / exploring / any state | **No.** Spark with zero questions still passes. |
| Require `## Crux` or `## Reasons` on `aligned` or `open` | **No.** |
| Require H3 Human/Agent on `aligned` or `open` | **No.** |
| Grade Positions / Reasons / Crux / Definition prose | **No.** DEC-005. |
| Score substantive vs bare | **No.** Human or adversarial reviewer. |
| Require a periodic assumption-audit file | **No.** |
| Require `## Dissent` on existing or new DECs | **No.** |
| Add a required `index.md` H2 | **No.** |
| Fail extra Crux/Reasons on `open`/`aligned` | **No.** |
| Keep agreement history / fail a flip back to `open` | **No.** |
| Add `think` / `spar` to the log-op regex | **No.** |
| Call network / `gh` / read `GH_TOKEN` | **No.** |
| Require the council pack | **No.** Core items 1–17 pass without it. |
| Require any CMP / RPT / RCL | **No.** Spark with zero reviews still passes. |
| Require a council to leave `spark` | **No.** |
| Require `## Crux` changes or reopen DEC-007 | **No.** |
| Grade Position / Findings / Dissent / Retained dissent / Prompt prose | **No.** DEC-005. |
| Content-score reports or reconciliation method | **No.** |
| Require distinct `model` strings | **No.** |
| Read `~/.config/mycelium` or `$MYCELIUM_CONFIG` | **No.** |
| Enforce panel size beyond council ≥2 | **No.** |
| Require `panels.toml` to exist | **No.** |
| Call a model / Cursor / network / `gh` / read `GH_TOKEN` | **No.** |
| Add `council` / `replicate` / `ladder` to the log-op regex | **No.** |
| Fail a lone CMP (no RPT yet) | **No.** WIP is legal. |
| Treat OQ-006 as a council | **No.** |
| Change `state` because an artifact was superseded | **No.** |
| Require a new H2 on OLD or NEW | **No.** |
| Grade OLD/NEW prose | **No.** DEC-005. |
| Fail G3 as a *status* fixture | **No.** G3 is status-only. |
| Fail G1 because the *binary* contract requires `github_repo` | **No.** Instance contract wins. |
| Add `handoff` / `upgrade` / `council` to the log-op regex | **No.** |
| Call network / `gh` / read `GH_TOKEN` | **No.** |
| Treat Install SLO or a GitHub Release as a check | **No.** |

## Teaching errors

Four lines on stderr, exit 1:

```text
mycelium: <what failed>
convention: <violated convention name>
contract: <path or DEC id>
fix: <suggested command or edit>
```

Cap 20 errors, then:

```text
mycelium: further errors omitted
```

## Success stdout

```text
mycelium check: ok
instance: <slug>
state: spark
tier: focused
artifacts: <n>
```

## Allowed top-level paths (undeclared extras fail)

Always allowed:

```text
README.md
mycelium.toml
log.md
CONTEXT.md
AGENTS.md
.gitignore
LICENSE
CHANGELOG.md
index.md
briefs/
.agents/
.mycelium/
.git/
.github/
program/
```

Type homes are allowed when they exist: `decisions/`, `assumptions/`, `evidence/`, `questions/`, `risks/`, `spikes/`, `findings/`, `recommendations/`, `requirements/`, `phases/`, `milestones/`.

Inside a type home, files must match the filename pattern or be `README.md`. Other files fail ID-to-path.

Inside `program/` and `.agents/`, extra files are allowed (methodology copy / runtime adapters).

Inside `briefs/`, extra files do not fail ID-to-path (`briefs/` is not an
ID-to-path home).

**Architect default** deviation key for a human scratch file: `extra-top-level:<path>` e.g. `extra-top-level:notes.md`.

Missing bound files/dirs fail. Extra valid artifacts at a low tier pass. Homes are allowed when they exist at ANY tier — do not require them to be in the active tier's binds.

## check --abort-journal

`--abort-journal`: delete staged temps listed in `.mycelium/journal.json`, delete the journal, delete a stale lock file. Do **not** delete already-renamed artifacts (no-deletion). Print the surviving paths. Exit 0 if the journal/lock are gone; exit 1 if there was nothing to abort (teaching error: no journal).

No confirmation step.
