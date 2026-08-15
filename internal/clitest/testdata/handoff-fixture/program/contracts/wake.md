# Wake Ritual (Mycelium 2.0)

Wake = simmering → exploring plus a written brief.

The CLI writes the brief; the human/agent reads it. The CLI does not grade
Suggested next.

## Paths

| Path | Rule |
| --- | --- |
| `briefs/WAKE-YYYY-MM-DD.md` | Date = `clock.Now().UTC()` date. If a file for that UTC date exists, **overwrite**. |
| `briefs/LATEST.md` | Identical bytes of that day's file. Always overwrite. |

Print the dated path on success.

## Required H2s

```text
## Parked
## Log since simmer
## Evidence triggers
## Assumptions
## Suggested next
```

Check validates these H2s on `briefs/LATEST.md` when the log contains a `wake`
op. Do not require a wake brief on instances that never simmered (no `wake` op
in the log). Do not grade Suggested next prose.

## Must cite (IDs, not URLs)

| Must cite | How |
| --- | --- |
| The simmer log line | Date + note (the line that recorded simmer). |
| Each due/overdue EVD | Cite `EVD-###` whose `## Revalidation Trigger` date is due or overdue vs clock (`date <= clock date`). |
| Each qualifying ASM | Cite `ASM-###` whose status is `Open` or `Held`, **plus** any whose `## Revisit Triggers` has a `YYYY-MM-DD` that is due or overdue. |
| Must not be required | Retired assumptions with no due date (fixture `ASM-002`). |

Deterministic citation IDs, not prose. Same fixture + same clock → same
citation IDs.

Link scan extends to `index.md` and `briefs/*.md` for
`\b(DEC|ASM|EVD|...)-[0-9]+\b`. Do not crawl the web.

## Wake algorithm (binding)

1. **Preflight.** Instance root found. Manifest+log parse. `state` must be
   `simmering`. `revisit` must be set and match `program/contracts/revisit.md`.
   Else teaching error, no writes. Early wake (clock date ≤ revisit date) is
   **legal**.
2. **Collect log since simmer.** Most recent parseable line with
   `op == "state"` and `note` starting with `simmering`. Collect subsequent
   parseable lines **in file order** after that line. If none, use all
   parseable log lines on or after `created_date`. The brief must still cite
   the simmer line itself (date + note).
3. **Scan `evidence/*.md`.** For each EVD file, parse `## Revalidation Trigger`.
   Extract trigger date (`program/contracts/revisit.md`). Include if date
   present and `date <= clock.Now().UTC() date`. Missing date ≠ due.
4. **Scan `assumptions/*.md`.** Include if `status ∈ {Open, Held}` **OR** a
   Revisit Triggers date is present and `date <= clock date`.
5. **Stage** `briefs/WAKE-YYYY-MM-DD.md`, `briefs/LATEST.md` (identical bytes),
   regenerated `index.md`, log (`wake` line), manifest (`state=exploring`,
   `revisit=""`, `updated_date` bumped). Commit order per
   `program/contracts/lifecycle.md`.
6. **Print** the dated brief path. Exit 0.

`mycelium state exploring` from simmering runs this same algorithm. Shared
function. No second implementation.

## Commands

```text
mycelium wake [--dir PATH]
mycelium state exploring [--dir PATH]   # from simmering: same algorithm
```

Log op is `wake`, ID `-`, note `exploring`.

Success stdout:

```text
woke briefs/WAKE-YYYY-MM-DD.md
state: exploring
```

## Refuse from non-simmering

```text
mycelium: wake is only legal from simmering
convention: lifecycle
contract: program/contracts/lifecycle.md
fix: mycelium state simmering --revisit YYYY-MM-DD   # if you meant to park first
```

## Template shape

```text
# Wake — 2026-08-09

## Parked

Simmered on 2026-08-01 with revisit 2026-08-08.

## Log since simmer

2026-08-01	state	-	simmering revisit=2026-08-08

## Evidence triggers

EVD-001 (revalidation 2026-08-06)

## Assumptions

ASM-001 (Held; revisit 2026-08-05)

## Suggested next

<!-- fill -->
```

The dated H1 is optional extra (not a required H2). Required H2s are the five
`##` names. Citation IDs in the body must resolve.

## Packages

`internal/wakebrief` + `internal/statecmd` land in later slices. This contract
is the spec only.
