# Conformance Suite (Mycelium 2.0)

Structure only ([DEC-005](../../framework/decisions/DEC-005-convention-over-configuration.md)).
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
    ([DEC-013](../../framework/decisions/DEC-013-stage-range-refuse.md)).

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

## Allowed top-level paths (focused spark baseline)

`mycelium.toml`, `README.md`, `AGENTS.md`, `CONTEXT.md`, `log.md`, `program/`,
`.agents/`, `.mycelium/`, `.git/`, plus homes that the active tier newly binds.

Deviation key for extras: `extra-top-level:<path>`.

## check --abort-journal

When a journal remains after an interrupted operation, `mycelium check` teaches
recovery. `mycelium check --abort-journal` discards the journal and staged files
after confirming abort intent (operator-driven; does not invent a new ID).
