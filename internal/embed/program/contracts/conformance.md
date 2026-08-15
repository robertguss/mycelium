# Conformance Suite (Mycelium 2.0)

Structure only (DEC-005).
Checks never grade prose or thinking quality.

## Must-implement checks (11)

1. Manifest present (`mycelium.toml`) and parses; required fields valid per
   `program/contracts/manifest.md`.
2. Lifecycle state legality per `program/contracts/lifecycle.md` (PHASE-01
   storage rules).
3. Tier legality and tier-aware artifact binds.
4. ID uniqueness within each namespace home.
5. ID-to-path integrity both directions (`program/contracts/naming.md`).
6. Link resolution for ID references.
7. Required front matter and sections per sidecar schema.
8. Log line prefixes parseable.
9. Interrupted operation: leftover journal or stale lock → teaching recovery
   (complete or `--abort-journal`).
10. Undeclared extra top-level paths unless deviation
    `extra-top-level:<path>` is declared.
11. Stage-scoped IDs outside every declared range → FAIL
    (DEC-013).

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
.agents/
.mycelium/
.git/
.github/
program/
```

Type homes are allowed when they exist: `decisions/`, `assumptions/`, `evidence/`, `questions/`, `risks/`, `spikes/`, `findings/`, `recommendations/`, `requirements/`, `phases/`, `milestones/`.

Inside a type home, files must match the filename pattern or be `README.md`. Other files fail ID-to-path.

Inside `program/` and `.agents/`, extra files are allowed (methodology copy / runtime adapters).

**Architect default** deviation key for a human scratch file: `extra-top-level:<path>` e.g. `extra-top-level:notes.md`.

Missing bound files/dirs fail. Extra valid artifacts at a low tier pass. Homes are allowed when they exist at ANY tier — do not require them to be in the active tier's binds.

## check --abort-journal

`--abort-journal`: delete staged temps listed in `.mycelium/journal.json`, delete the journal, delete a stale lock file. Do **not** delete already-renamed artifacts (no-deletion). Print the surviving paths. Exit 0 if the journal/lock are gone; exit 1 if there was nothing to abort (teaching error: no journal).

No confirmation step.
