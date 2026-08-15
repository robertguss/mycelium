# Revisit Trigger Grammar (Mycelium 2.0)

Manifest field `revisit` is a **string**. Two shapes only.

## Shapes

| Shape | Grammar | Overdue? | `status` `due:` |
| --- | --- | --- | --- |
| Date | `YYYY-MM-DD` (UTC, real calendar date) | Yes, iff clock UTC **date** is **strictly after** that date | `yes` if clock date ≥ that date; `no` if clock date < that date |
| Event | `event:<kebab>` | Never auto-overdue | `event` |

Due on that date (`status` lists it under Due; single-instance `due: yes`).
Overdue strictly after (next UTC day). On the due date it is not yet overdue.

## Date parse

`time.Parse("2006-01-02", value)` in UTC.

- `2026-02-30` → refuse / fail.
- No time component.
- No `2026-8-8`.
- No `2026/08/08`.

## Event parse

```text
^event:[a-z0-9]+(?:-[a-z0-9]+)*$
```

Prefix is lowercase `event:`.

Refuse: `EVENT:foo`, `event:After-Launch`, `event:`, `event:after_iphone`.

Legal examples: `event:after-iphone-launch`, `event:budget-review`, `event:q4`.

## Due / overdue examples

| `revisit` | Clock UTC date | due | overdue | `--all` bucket |
| --- | --- | --- | --- | --- |
| `2026-08-08` | `2026-08-07` | no | no | the rest (future-dated simmering) |
| `2026-08-08` | `2026-08-08` | yes | no | due today |
| `2026-08-08` | `2026-08-09` | yes | yes | overdue simmering first |
| `event:after-iphone-launch` | any | event | no | event-simmering |

## Verdict table

| Value | Verdict |
| --- | --- |
| `""` | refuse / fail when simmering |
| `2026-08-08` | date |
| `event:after-iphone-launch` | event |
| `in two weeks` | refuse / fail |
| `2026-08-08T00:00:00Z` | refuse / fail |
| `after-iphone-launch` | refuse / fail (missing `event:`) |

Anything else: refuse on `mycelium state simmering` (no writes). Fail `check`
if `state=simmering`.

## Clear on leave

Clear `revisit` to `""` when leaving simmering (wake, archive from simmering,
or any successful transition whose target is not `simmering`).

## Trigger-date extract (containers)

Used by wake for evidence / assumption sections.

Regex: first line matching `^\s*(\d{4}-\d{2}-\d{2})\b` in the named section.
First date wins. No NLP. Missing date ≠ due.

| File | Section H2 (exact, case-sensitive) | What is extracted |
| --- | --- | --- |
| `evidence/EVD-###-*.md` | `## Revalidation Trigger` | first `YYYY-MM-DD` |
| `assumptions/ASM-###-*.md` | `## Revisit Triggers` | first `YYYY-MM-DD` |

Section body = bytes after the exact H2 until the next H2 or EOF.

`2026-08-06T00:00:00Z` does **not** match (`\b` after the date; next char must
be a non-word or end). `2026-08-06` and `2026-08-06 something` match.

Do not parse Claim / Statement / other sections for dates.

## Package

`internal/revisit` lands in Slice 1 (spec only in this contract):

```text
Parse(s string) (Kind, dateOrEvent, error)    # Kind = Date | Event
Due(kind, date, now time.Time) bool           # date shape && nowDate >= date
Overdue(kind, date, now time.Time) bool       # date shape && nowDate > date
ExtractTriggerDate(sectionBody string) (date, ok)
```

No filesystem. No CLI.
