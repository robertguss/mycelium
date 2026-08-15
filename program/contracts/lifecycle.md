# Idea Lifecycle (Mycelium 2.0)

Exact vocabulary (DEC-006):

```text
spark → exploring ⇄ simmering → clarified → handed-off
any (except archived) → archived
archived → (none)
```

Birth state is always `spark`.

`wake` is **not** a seventh state. It is the simmering → exploring ritual.

## Legal edges (PHASE-02)

| From | Allowed next | Command | Extra rule |
| --- | --- | --- | --- |
| `spark` | `exploring` | `mycelium state exploring` | No `--revisit`. Clear `revisit` to `""` if somehow set. |
| `spark` | `archived` | `mycelium state archived` | Does not delete files. Clear `revisit`. |
| `exploring` | `simmering` | `mycelium state simmering --revisit VALUE` | **Require** `--revisit`. Refuse if missing or empty. VALUE must match `program/contracts/revisit.md`. |
| `exploring` | `clarified` | `mycelium state clarified` | No packet. `clarified` is a **legal stored state**. |
| `exploring` | `archived` | `mycelium state archived` | No deletion. |
| `simmering` | `exploring` | `mycelium wake` (preferred) or `mycelium state exploring` | **Must write the re-entry brief.** Silent wake forbidden. Log op is `wake`. Clear `revisit` to `""`. |
| `simmering` | `archived` | `mycelium state archived` | No deletion. Clear `revisit`. No brief required. |
| `clarified` | `archived` | `mycelium state archived` | No deletion. |
| `clarified` | `handed-off` | *refused* | Teaching error names the PHASE-06 packet. |
| `handed-off` | anything | *refused* | Same PHASE-06 teaching error. Check still FAILs stored `handed-off`. |
| `archived` | anything | *refused* | Terminal. Teaching error: archived is terminal. |

## Allowed PHASE-02 argv targets

`exploring` | `simmering` | `clarified` | `archived`

- `handed-off` is **never** a legal argv target this phase.
- `spark` is not a legal *target* (nothing transitions back to spark).

## Commands

```text
mycelium state <target> [--dir PATH] [--revisit VALUE]
mycelium wake [--dir PATH]
```

- `state` is the one transition command. Legal edges come from this contract.
- `wake` is the ritual-bearing alias: **ONLY** legal from `simmering` → `exploring`. It writes the re-entry brief **then** flips state.
- `mycelium state exploring` from `simmering` is **also legal** and **MUST write the same brief** as `wake` (shared function). Silent wake is forbidden.
- Log op: `state` for generic transitions; `wake` when the transition is simmering → exploring (whether invoked as `wake` or as `state exploring`).

## PHASE-02 storage rules

Replace PHASE-01 storage rules.

| Stored state | Check |
| --- | --- |
| `spark` | Legal. |
| `exploring` | Legal. Do not fail leftover `revisit` (empty or set). Only `state`/`wake` clear it. |
| `simmering` | Legal iff `revisit` matches `program/contracts/revisit.md`. Empty or malformed → FAIL. |
| `clarified` | **Legal** (PHASE-01 fail rule lifted for clarified only). |
| `handed-off` | **FAIL**. Teaching error names PHASE-06 packet. |
| `archived` | Legal. Terminal. |
| unknown | FAIL. |

`exploring` / `clarified` / `archived` / `spark` do not fail on `revisit` value.

## Protocol

`state` and `wake` run under the existing operation protocol (lock / journal /
stage / commit / rollback). See `program/contracts/operation-protocol.md`.

Journal `op` values this phase: `scaffold` | `new` | `tier` | `publish` |
`state` | `wake` | `index`.

Commit order for `state` / `wake`:

1. brief + `briefs/LATEST.md` (if this transition is a wake)
2. `index.md`
3. `log.md`
4. `mycelium.toml` **last** (`state` updated, `revisit` cleared or set, `updated_date` bumped)

Never `git add`. Never `git commit`.

`mycelium index` uses the same protocol with journal `op=index`. `index` does
**not** append a log line. Check log-ops stay
`scaffold|new|tier|publish|check|state|wake|supersede`.

## Illegal edge

Teaching error, name `program/contracts/lifecycle.md`, print allowed next
states from the table. Exit 1. No writes.

## `--revisit`

`--revisit` is legal **only** when `<target>` is `simmering`. Refuse on any
other target.

## Same-state

- Already `<target>` (and not a simmering revisit update): no-op, exit 0,
  stdout `already <target>`.
- Same-state simmering with a new legal `--revisit`: allowed. Update `revisit`,
  bump `updated_date`, append a `state` log line, regen `index.md`.

## Teaching-error examples

```text
mycelium: illegal transition spark → clarified
convention: lifecycle
contract: program/contracts/lifecycle.md
fix: allowed next states: exploring, archived

mycelium: state=handed-off requires a PHASE-06 handoff packet
convention: lifecycle
contract: program/contracts/lifecycle.md
fix: stay in clarified, or mycelium state archived; packet command is not shipped
```
