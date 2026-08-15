# Conformance Suite (Mycelium 2.0)

Structure only (DEC-005).
Checks never grade prose or thinking quality.

## Must-implement checks (17)

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
   `scaffold|new|tier|publish|check|state|wake`. Regex:
   `^\d{4}-\d{2}-\d{2}\t(scaffold|new|tier|publish|check|state|wake)\t(\S+)\t`
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

## Lift timing

Items 1–14 stay as written (PHASE-02). Items 15–17 land by slice:

| Slice | Check behavior |
| --- | --- |
| 1 | Schema `required_sections` drops `Crux`. Schema-driven check stops requiring Crux on every OQ. **No** IFF bind yet. |
| 2 | Check calls `internal/sparring` for each `questions/OQ-*.md`. IFF rules bind (item 15). |
| 3 | `CONTEXT.md` glossary rules bind (item 16). |
| 4 | Optional DEC Dissent rule binds (item 17). |

Do not require a wake brief on instances that never simmered. Do not grade
brief prose. Do not require N artifacts.

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
