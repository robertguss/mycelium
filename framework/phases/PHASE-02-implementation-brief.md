# PHASE-02 Implementation Brief — Lifecycle, Wake, Portfolio

- **Status:** Binding
- **Date:** 2026-08-15
- **Audience:** Engineering (pstack / poteto-mode)
- **Authority:** `framework/blueprint.md` (Accepted 2026-08-14). DEC-001 through DEC-014 stand. This brief commissions PHASE-02 only. It does not record a new DEC (see Appendix A).
- **Repo:** https://github.com/robertguss/mycelium
- **Pin:** Engineering starts from `main` @ `75645c3c2b48cd485a590cb0f0158d7cb29da1df` (PHASE-01 slices 0–10 on main). Do not implement from a later SHA unless Arvo re-pins in writing.
- **Product:** single-binary Go CLI `mycelium`. Master builds the CLI. `program/` is `go:embed`'d and emitted on scaffold. PHASE-02 adds lifecycle transitions, index + wake brief, and portfolio status.
- **Phase gate:** MS-201 via hermetic `go test ./...` (fixture wake). GitHub Actions is **not** a gate (Robert waived CI). One dogfood wake after seven real days is human evidence for Arvo, not the gate.
- **How to use this file:** Engineering executes from THIS FILE ONLY. No "see chat". No TBD. Open items are decided here and labeled **Architect default**.

Headings: §§1–19 then Appendices A–F (no DEC-015, wake brief template, index.md example, `status --all` format, MS-201 fixture recipe, target file tree).

Cloud env name is exactly `robertguss/mycelium`. Go 1.26 at `/usr/local/go` or `/usr/local/bin/go`. `CGO_ENABLED=0`.

## 1. Scope / out of scope

Tonight is PHASE-02 only. Later phases get their own briefs after MS-201 is accepted. Do not implement, stub-ship, or "leave a hook command" for later phases. Do not implement PHASE-01 leftovers.

### In scope

- Lifecycle command `mycelium state <target>` and ritual alias `mycelium wake`.
- PHASE-02 storage rules that **replace** PHASE-01 storage rules in `program/contracts/lifecycle.md` and `mycelium check`.
- Revisit trigger grammar (date + event) and due/overdue math against `internal/clock`.
- `index.md` convention, generator, skeleton emit, `mycelium index` rebuild.
- Re-entry brief at `briefs/WAKE-YYYY-MM-DD.md` plus `briefs/LATEST.md`.
- First usable `mycelium status` and `mycelium status --all` (live GitHub enumeration with documented local fallback).
- Skills: `spark`, `wake`, `portfolio`, plus `mycelium-cli` command update.
- Check updates listed in §11.
- MS-201 hermetic fixture test in `go test ./...`. No Actions job.
- Commissioning artifacts: this brief, rewritten lifecycle contract, new contracts (index, wake, status, revisit), acceptance stub.

### Out of scope (Quality refuses PRs that add them)

- PHASE-03 sparring / thinking-mode / cruxes (DEC-007).
- PHASE-04 perspective ladder / council / second-opinion (DEC-008).
- PHASE-05 `mycelium supersede` and release/install work.
- PHASE-06 handoff packet. `handed-off` stays unreachable.
- Implementing MS-101(b). Do not commission a `GH_TOKEN` job. Do not reopen publish.
- Growing `latinFold`, adding NFKD, or adding `golang.org/x/text` (DEC-014).
- Converting master (`research-program.toml`, `just init`, deleting Justfile/scripts).
- Emitting `framework/`.
- CLI `git add` / `git commit` of instance work product.
- Content grading of wake-brief prose (DEC-005).
- Requiring a 7-real-day dogfood wake as the phase gate.
- `destroy`, `range`, `explore` / `simmer` as separate verbs, `council`, `handoff`.
- Retrofitting skills into existing PHASE-01 instances via `state` / `wake` / `index`. Pushing to `main`.

### Master vs instance (unchanged)

Master remains an ADRP v1 instance for its own evolution. Do not convert master's `research-program.toml` to `mycelium.toml`. Do not run `just init` on master. `framework/` stays master-only and is NEVER emitted. Justfile/scripts stay on master. PHASE-02 changes the operational surface for *idea instances* only.

## 2. Authority and do-not-reopen

### Governing documents

| Document | Rule |
| --- | --- |
| `framework/blueprint.md` (Accepted 2026-08-14) | Do not rewrite vision. DEC-001–014 stand. |
| DEC-003 | One repo per idea. `status --all` enumerates by `idea` topic. |
| DEC-005 | Checks validate containers, never contents. |
| DEC-006 | spark → exploring ⇄ simmering → clarified → handed-off; any → archived. |
| DEC-007 / DEC-008 | Sparring = PHASE-03. Ladder = PHASE-04. Not this brief. |
| DEC-010 / DEC-011 | CLI never commits. No migration. Runtime reads instance files. |
| DEC-012 / DEC-013 / DEC-014 | Do not reopen (`mycelium.toml`; refuse out-of-range; `latinFold` only, no NFKD, no `x/text`, do not grow the map). |
| This brief | Binding 2026-08-15. PHASE-02 only. Architect defaults are binding. |

### Process override (unchanged)

Blueprint "humans-own-git" is overridden for the *master* repo's engineering workflow: Arvo merges Quality-green PRs and accepts the phase. Engineering opens PRs. Engineering does NOT push to main. The CLI still never git-commits *instance* work product.

### Do not reopen

Do not reopen the product shape, the language, the dependency floor, the state vocabulary, the manifest filename, the refuse-vs-warn range rule, the no-commit rule, the instance-files-are-truth rule, slugify/DEC-014, publish, or MS-101(b). If a later discovery contradicts a locked item, stop and write a DEC; do not silently patch this brief in a code PR.

Do not reopen:

- DEC-012 (`mycelium.toml`)
- DEC-013 (refuse out-of-range)
- DEC-014 (slugify = existing `latinFold`; no NFKD; no `x/text`; do not grow the map)

### MS-101 status (planning accepted; leftovers stay leftovers)

- MS-101(a) hermetic is green on the pin.
- MS-101(b) has **not** passed (`GH_TOKEN` missing; skip ≠ pass).
- Arvo accepted PHASE-02 planning anyway.
- Do **not** implement PHASE-01 leftovers.
- Do **not** commission a `GH_TOKEN` job.
- Do **not** reopen publish.

### Phases you must not commission

| Phase | Why it is not this brief |
| --- | --- |
| PHASE-03 | DEC-007 sparring |
| PHASE-04 | DEC-008 ladder |
| PHASE-05 | `mycelium supersede` |
| PHASE-06 | Handoff packet; `handed-off` becomes legal only then |

## 3. What PHASE-01 left on main (floor; do not reimplement)

Pin: `75645c3c2b48cd485a590cb0f0158d7cb29da1df`. Treat this SHA as the floor. Reuse packages. Do not rewrite working PHASE-01 commands.

### Already shipped (do not rebuild)

Reuse: `cmd/mycelium`, `internal/{cli,version,embed,clock,execrun,metadata,idpath,manifest,schema,slug,logfmt,teach,lock,journal,op,scaffold,generate,check,tiercmd,publish,clitest}`.

| Touch | Fate |
| --- | --- |
| `internal/cli` | Add `state`, `wake`, `status`, `index`. |
| `internal/clock` + `MYCELIUM_NOW` | Reuse. MS-201 injects this clock. |
| `internal/execrun` | `status --all` calls `gh` only through this. |
| `internal/slug` | Do not touch (DEC-014). |
| `internal/logfmt` | Keep. Extend allowed ops in check. |
| `internal/lock` / `journal` / `op` | Extend journal `op` enum. |
| `internal/scaffold` / `generate` / `tiercmd` / `publish` | Emit or regenerate `index.md`. Do not emit skills from `tier`. Do not reopen publish. |
| `internal/check` | Replace PHASE-01 storage-rule bullets (§11). |
| `internal/clitest` | Add MS-201 fixture test. |
| `program/contracts/lifecycle.md` | **Rewrite** (Slice 0). |
| `program/skeleton/` + `program/skills/mycelium-cli` | Add `index.md`; update command table. |
| `phase-01-hermetic.yml` / `phase-01-github.yml` | Leave alone. Do **not** add a PHASE-02 workflow. Actions is not a gate. |
| Justfile / scripts / `research-program.toml` | Keep. Do not delete. Do not `just init`. |
| `framework/` | Master-only. NEVER emitted. |
| `internal/version` | Stay `0.1.0-dev`. Re-run embed generate after `program/` edits. |

### PHASE-01 behaviors that stay true

- Birth state is `spark`. There was no state-transition command; PHASE-02 adds one.
- Stored `clarified` / `handed-off` currently FAIL check. PHASE-02 lifts **clarified only**.
- `simmering` already requires `revisit` non-empty (manifest parse). PHASE-02 adds grammar.
- Stored `archived` is already legal.
- Clock is injectable (`internal/clock` + `MYCELIUM_NOW` RFC3339).
- `git` / `gh` go through `internal/execrun`. Hermetic tests assert `gh` was not called under `--offline` / `MYCELIUM_OFFLINE=1`.
- Log lines are tab-separated (`internal/logfmt`).
- Exit 0 success, exit 1 every failure. Teaching errors: four lines on stderr (`mycelium` / `convention` / `contract` / `fix`). Cap 20.
- `--dir PATH` is the instance root. Absent: walk upward for `mycelium.toml`, stop at `.git` or filesystem root.
- CLI never git-commits. `git init` only on scaffold, already shipped.
- Runtime check/generate read instance `program/`, never embed (DEC-011). Dependency floor unchanged (stdlib + `pelletier/go-toml/v2` only). Module `github.com/robertguss/mycelium`. `CGO_ENABLED=0`. Go 1.26.

### What must not be broken

`just check` on master; hermetic `go test ./...`; no `framework/` emit; no master conversion. Existing spark instances without `index.md` are *not* silently migrated by `state`/`wake` — they fail `check` once Slice 2 binds `index.md`, and are repaired by `mycelium index`.

If a PHASE-02 PR is bad: revert that PR. Floor is the pin SHA.

## 4. Lifecycle commands + transition table

PHASE-02 storage rules **replace** PHASE-01 storage rules. Rewrite `program/contracts/lifecycle.md` in Slice 0. Implement the command in Slice 3. The table is encoded now (Slice 1 tests it as a pure function).

### Vocabulary (DEC-006; do not bikeshed)

```text
spark → exploring ⇄ simmering → clarified → handed-off
any (except archived) → archived
archived → (none)
```

`wake` is **not** a seventh state. It is the simmering → exploring ritual.

### Legal edges (commanded this phase)

| From | Allowed next | Command | Extra rule |
| --- | --- | --- | --- |
| `spark` | `exploring` | `mycelium state exploring` | No `--revisit`. Clear `revisit` to `""` if somehow set. |
| `spark` | `archived` | `mycelium state archived` | Does not delete files. Clear `revisit`. |
| `exploring` | `simmering` | `mycelium state simmering --revisit VALUE` | **Require** `--revisit`. Refuse if missing or empty. VALUE must match §5. |
| `exploring` | `clarified` | `mycelium state clarified` | No packet. `clarified` is now a **legal stored state**. |
| `exploring` | `archived` | `mycelium state archived` | No deletion. |
| `simmering` | `exploring` | `mycelium wake` (preferred) or `mycelium state exploring` | **Must write the re-entry brief.** Silent wake is a Quality refuse. Log op is `wake`. Clear `revisit` to `""`. |
| `simmering` | `archived` | `mycelium state archived` | No deletion. Clear `revisit`. No brief required (not a wake). |
| `clarified` | `archived` | `mycelium state archived` | No deletion. |
| `clarified` | `handed-off` | *refused* | Teaching error names the PHASE-06 packet. |
| `handed-off` | anything | *refused* | Same PHASE-06 teaching error. Check still FAILs stored `handed-off`. |
| `archived` | anything | *refused* | Terminal. Teaching error: archived is terminal. |

### Allowed PHASE-02 targets

`exploring` | `simmering` | `clarified` | `archived`

`handed-off` is **never** a legal argv target this phase.

`spark` is not a legal *target* (nothing transitions back to spark).

### Singular command + ritual alias

```text
mycelium state <target> [--dir PATH] [--revisit VALUE]
mycelium wake [--dir PATH]
```

- `state` is the one transition command. Legal edges come from `program/contracts/lifecycle.md`.
- `wake` is the ritual-bearing alias: **ONLY** legal from `simmering` → `exploring`. It writes the re-entry brief **then** flips state. It is not a second state machine.
- `mycelium state exploring` from `simmering` is **also legal** and **MUST write the same brief**. Do not allow a silent wake.
- Log op: `state` for generic transitions; `wake` when the transition is simmering → exploring (whether invoked as `wake` or as `state exploring`).

### PHASE-02 storage rules (replace PHASE-01)

| Stored state | Check |
| --- | --- |
| `spark` | Legal. |
| `exploring` | Legal. `revisit` should be `""` after a wake; check does not fail a leftover empty-or-set revisit on exploring unless you want a warning — **Architect default:** exploring/clarified/archived/spark do not fail on `revisit` value; only `state`/`wake` clear it. |
| `simmering` | Legal iff `revisit` matches §5 grammar. Empty or malformed → FAIL. |
| `clarified` | **Legal** (PHASE-01 fail rule is lifted for clarified only). |
| `handed-off` | **FAIL**. Teaching error names PHASE-06 packet. No packet contract yet. |
| `archived` | Legal. Terminal. |
| unknown | FAIL. |

### Protocol

`state` and `wake` run under the existing operation protocol (lock / journal / stage / commit / rollback).

Journal `op` values this phase: `scaffold` | `new` | `tier` | `publish` | `state` | `wake` | `index`.

Commit order for `state` / `wake`:

1. brief + `briefs/LATEST.md` (if this transition is a wake)
2. `index.md`
3. `log.md`
4. `mycelium.toml` **last** (`state` updated, `revisit` cleared or set, `updated_date` bumped)

Never `git add`. Never `git commit`.

`mycelium index` uses the same protocol with journal `op=index`. **Architect default:** `index` does **not** append a log line (derived view). Check log-ops stay `scaffold|new|tier|publish|check|state|wake`.

### Illegal edge

Teaching error, name `program/contracts/lifecycle.md`, print allowed next states from the table. Exit 1. No writes.

### `--revisit` on a non-simmering target

**Architect default:** refuse. `--revisit` is legal only when `<target>` is `simmering`.

## 5. Revisit trigger grammar

Manifest field `revisit` stays a **string**. Two shapes only.

### Shapes

| Shape | Grammar | Overdue? | `status` `due:` |
| --- | --- | --- | --- |
| Date | `YYYY-MM-DD` (UTC, real calendar date) | Yes, iff `clock.Now().UTC()` **date** is **strictly after** that date | `yes` if clock date ≥ that date; `no` if clock date < that date |
| Event | `event:<kebab>` | Never auto-overdue | `event` |

**Architect default (due vs overdue):** the idea is **due on that date** (`status` lists it under Due; single-instance `due: yes`). It is **overdue after that date** (next UTC day). On the due date it is not yet overdue.

Examples:

| `revisit` | Clock UTC date | due | overdue | `--all` bucket |
| --- | --- | --- | --- | --- |
| `2026-08-08` | `2026-08-07` | no | no | the rest (future-dated simmering; after event-simmering, before archived) |
| `2026-08-08` | `2026-08-08` | yes | no | due today |
| `2026-08-08` | `2026-08-09` | yes | yes | overdue simmering first |
| `event:after-iphone-launch` | any | event | no | event-simmering |

### Date parse

**Architect default:** `time.Parse("2006-01-02", value)` in UTC. `2026-02-30` is refuse/fail. No time component. No `2026-8-8`. No `2026/08/08`.

### Event parse

**Architect default:** `^event:[a-z0-9]+(?:-[a-z0-9]+)*$`. Prefix is lowercase `event:`. `EVENT:foo`, `event:After-Launch`, `event:`, `event:after_iphone` all refuse.

Examples of legal events: `event:after-iphone-launch`, `event:budget-review`, `event:q4`.

### Anything else

Refuse on `mycelium state simmering` (no writes). Fail `check` if `state=simmering`.

| Value | Verdict |
| --- | --- |
| `""` | refuse / fail when simmering |
| `2026-08-08` | date |
| `event:after-iphone-launch` | event |
| `in two weeks` | refuse / fail |
| `2026-08-08T00:00:00Z` | refuse / fail |
| `after-iphone-launch` | refuse / fail (missing `event:`) |

### Clear on leave

Clear `revisit` to `""` when leaving simmering (wake, archive from simmering, or any successful transition whose target is not `simmering`).

### Package

New pure package `internal/revisit`:

```text
Parse(s string) (Kind, dateOrEvent, error)    # Kind = Date | Event
Due(kind, date, now time.Time) bool           # date shape && nowDate >= date
Overdue(kind, date, now time.Time) bool       # date shape && nowDate > date
ExtractTriggerDate(sectionBody string) (date, ok)  # §7 containers
```

No filesystem. No CLI. Slice 1 lands this.

## 6. log + index conventions

### Log (extend, do not replace)

Format unchanged:

```text
YYYY-MM-DD\t<op>\t<ID-or-->\t<title-or-note>
```

PHASE-01 ops: `scaffold` | `new` | `tier` | `publish` | `check` (`check` itself still does not append a line).

PHASE-02 ops added: `state` | `wake`.

| Transition | op | ID | note |
| --- | --- | --- | --- |
| spark → exploring | `state` | `-` | `exploring` |
| exploring → simmering | `state` | `-` | `simmering revisit=<VALUE>` |
| simmering → exploring (wake or `state exploring`) | `wake` | `-` | `exploring` |
| any legal → archived | `state` | `-` | `archived` |
| exploring → clarified | `state` | `-` | `clarified` |

Check log-line regex becomes:

```text
^\d{4}-\d{2}-\d{2}\t(scaffold|new|tier|publish|check|state|wake)\t(\S+)\t
```

Blank lines and ATX headings remain allowed. `mycelium index` does not append a log line.

### `index.md` (new, required)

Always emitted at spark from PHASE-02 onward. PHASE-02 updates `program/skeleton/` to include `index.md`.

| Rule | Binding |
| --- | --- |
| New scaffolds | Generator writes concrete `index.md` (not tokens). Skeleton may hold a structural stub that scaffold overwrites. |
| Existing PHASE-01 instances | May lack `index.md`. `mycelium index [--dir PATH]` rebuilds. `check` **FAILS** if missing (PHASE-02 bind, all tiers). |
| Regeneration | Every mutating command regenerates it: `state`, `wake`, `new <type>`, `tier`, `publish`, `index`. Scaffold writes it at birth. |
| Skills | `mycelium index` does **not** emit skills. |

### Required structure (check validates presence only)

```text
# <idea_name>
## State
## Artifacts
## Log tail
## Wake
```

| H2 | Generator writes | Check |
| --- | --- | --- |
| `# <idea_name>` | Concrete `idea_name` from the manifest | File has an H1. **Architect default:** H1 text need not equal `idea_name` (containers, not contents). |
| `## State` | `state`, `tier`, `revisit`, `github_repo` as concrete values | H2 present |
| `## Artifacts` | Counts per registered NS (files matching the filename pattern). Missing home = 0. No prose grading. | H2 present |
| `## Log tail` | Last 20 parseable log lines (not headings). Fewer than 20 → all of them. | H2 present |
| `## Wake` | `briefs/LATEST.md` if that file exists, else `none` | H2 present |

Tokens are not required. Extra H2s are allowed. Body prose is not graded.

### Allowed top-level paths (additions)

Add to the PHASE-01 always-allowed set:

```text
index.md
briefs/
```

`briefs/` is not an ID-to-path home. Extra files inside `briefs/` do not fail ID-to-path. If the log contains a `wake` op, `briefs/LATEST.md` must exist and pass the wake H2s (§7).

### Tier binds

**Architect default:** add `index.md` to `binds` in all three `program/tiers/*.toml` files (focused, standard, high-assurance). Do not add `briefs/` as a bind (only allowed; required only after a wake op).

### Shared renderer

New package `internal/indexmd` with `Render(instance) []byte`. Called by scaffold, `index`, `state`, `wake`, `new <type>`, `tier`, `publish`. No network.

## 7. Wake ritual + re-entry brief

Wake is simmering → exploring plus a written brief. Blueprint: reread index and log tail, check evidence revalidation triggers and assumption records against what changed, then brief the human. The CLI writes the brief; the human/agent reads it. The CLI does not grade Suggested next.

### Paths

| Path | Rule |
| --- | --- |
| `briefs/WAKE-YYYY-MM-DD.md` | Date = `clock.Now().UTC()` date. If a file for that UTC date exists, **overwrite** (one wake per UTC day is enough; a second wake the same day replaces). |
| `briefs/LATEST.md` | Identical copy of that day's file. Always overwrite. |

Print the dated path on success.

### Required H2s (structure only)

```text
## Parked
## Log since simmer
## Evidence triggers
## Assumptions
## Suggested next
```

Check validates these H2s on `briefs/LATEST.md` when the log contains a `wake` op. Do not grade the prose of Suggested next. Do not require a wake brief on instances that never simmered (no `wake` op in the log).

### What the body MUST cite (MS-201 / fixture)

Sources = artifact IDs and log dates, **not** web URLs.

| Must cite | How |
| --- | --- |
| The simmer log line | Date + note (the line that recorded simmer). |
| Each due/overdue EVD | Cite `EVD-###` whose Revalidation Trigger date is due or overdue vs clock (`date <= clock date`). |
| Each qualifying ASM | Cite `ASM-###`. **Architect default:** cite every assumption whose status is `Open` or `Held`, **plus** any whose Revisit Triggers section contains a `YYYY-MM-DD` that is due or overdue (`date <= clock date`). |
| Must not be required | Retired assumptions with no due date (fixture `ASM-002`). |

Deterministic: same fixture + same clock → same **citation IDs**, not the same prose. Check validates H2s + that cited IDs resolve (existing link resolution, extended to `briefs/*.md` and `index.md`).

**Architect default — link scan extension:** add `index.md` and `briefs/*.md` to the files scanned for `\b(DEC|ASM|EVD|...)-[0-9]+\b`. Do not crawl the web.

### Wake algorithm (binding)

1. **Preflight.** Instance root found. Manifest+log parse. `state` must be `simmering`. `revisit` must be set and match §5. Else teaching error, no writes. Early wake (clock date ≤ revisit date) is **legal** — humans can wake early.
2. **Collect log since simmer.** Find the most recent log line that recorded the simmer. **Architect default:** the most recent parseable line with `op == "state"` and `note` starting with `simmering`, or any op that recorded the simmer (same note prefix). Collect subsequent parseable lines **in file order** after that line. If none, use all parseable log lines on or after `created_date`. The brief must still cite the simmer line itself (date + note).
3. **Scan `evidence/*.md`.** For each file matching the EVD pattern, parse the `## Revalidation Trigger` section. Extract a trigger date (§7 containers). Include the EVD if a date is present and `date <= clock.Now().UTC() date`. Missing date ≠ due.
4. **Scan `assumptions/*.md`.** Parse front-matter `status` via the existing metadata reader. Include if `status ∈ {Open, Held}` **OR** a Revisit Triggers date is present and `date <= clock date`.
5. **Stage** `briefs/WAKE-YYYY-MM-DD.md`, `briefs/LATEST.md` (identical bytes), regenerated `index.md`, log (`wake` line), manifest (`state=exploring`, `revisit=""`, `updated_date` bumped). Protocol commit order §4.
6. **Print** the dated brief path. Exit 0.

`mycelium state exploring` from simmering runs this same algorithm. Shared function. No second implementation.

### Evidence / assumption trigger parse (containers)

**Architect default:** a line matching `^\s*(\d{4}-\d{2}-\d{2})\b` anywhere in the named section is a trigger date. First date wins. No NLP. Missing date ≠ due.

Section body = bytes after the exact H2 (`## Revalidation Trigger` or `## Revisit Triggers`) until the next H2 or EOF. Heading match is exact, case-sensitive (already the PHASE-01 H2 rule).

| File | Section H2 | What is extracted |
| --- | --- | --- |
| `evidence/EVD-###-*.md` | `Revalidation Trigger` | first `YYYY-MM-DD` |
| `assumptions/ASM-###-*.md` | `Revisit Triggers` | first `YYYY-MM-DD` |

Do not parse Claim / Statement / other sections for dates. Do not treat ISO datetimes as dates unless the line matches the regex (a leading `2026-08-06T…` still matches the date prefix — **Architect default:** the regex is `\b` after the date, so `2026-08-06T00:00:00Z` does **not** match; the next char must be a non-word or end. `2026-08-06` and `2026-08-06 something` match).

`internal/revisit.ExtractTriggerDate` owns this. Wake calls it. Tests table-drive it.

### Package

`internal/wakebrief` — collect citations + write markdown. `internal/statecmd` — `state` and `wake` CLI. Do not invent a second writer.

## 8. Commands

Exact CLI. Flags. Exit codes. Teaching errors.

Global: exit 0 success, exit 1 failure. Teaching errors on stderr (Appendix E of PHASE-01 format, restated in examples below). Success text on stdout.

Env (unchanged + one addition):

| Env | Effect |
| --- | --- |
| `MYCELIUM_OFFLINE=1` | Treat as `--offline` on every command that can touch GitHub. Never exec `gh`. Never open network. |
| `MYCELIUM_NOW` | RFC3339 clock override. |
| `MYCELIUM_IDEAS_ROOT` | Default ideas root for `status --all` when `--root` is absent. |

### Commands that exist after PHASE-02

| Command | New? | Protocol? |
| --- | --- | --- |
| `mycelium version` | no | no |
| `mycelium new idea <name> [--dir] [--offline] [--publish] [--tier]` | no (emit `index.md`) | yes |
| `mycelium new <type> <title> [--dir]` | no (regen `index.md`) | yes |
| `mycelium check [--dir] [--abort-journal]` | no (storage-rule update) | abort only |
| `mycelium tier <tier> [--dir]` | no (regen `index.md`; do not emit skills) | yes |
| `mycelium publish [--dir]` | no (regen `index.md`; do not reopen) | yes |
| `mycelium state <target> [--dir PATH] [--revisit VALUE]` | **yes** | yes |
| `mycelium wake [--dir PATH]` | **yes** | yes |
| `mycelium status [--dir PATH] [--all] [--root PATH] [--archived] [--offline]` | **yes** | no (read-only; never mutates) |
| `mycelium index [--dir PATH]` | **yes** | yes |

### Commands that do not exist this phase

`supersede`, `council`, `handoff`, `explore` / `simmer` as separate verbs, `destroy`, `range`. Quality refuses PRs that add them.

### 8.1 `mycelium state <target>`

```text
mycelium state <target> [--dir PATH] [--revisit VALUE]
```

| Flag / arg | Rule |
| --- | --- |
| `<target>` | One of `exploring`, `simmering`, `clarified`, `archived`. Else teaching error listing those four plus the note that `handed-off` is PHASE-06. |
| `--dir PATH` | Instance root. Walk-up if absent. |
| `--revisit VALUE` | Required when target is `simmering`. Forbidden on any other target. |

Algorithm:

1. Resolve root. Preflight: manifest+log parse, no leftover journal for a different op, lock free.
2. If target == current state and (target != simmering or revisit unchanged): **Architect default:** no-op, exit 0, stdout `already <target>`. If target is simmering and `--revisit` differs, treat as an update: rewrite `revisit`, bump `updated_date`, append a `state` log line, regen index. Still require the edge to be legal (it is; same state). **Architect default:** same-state simmering with a new legal revisit is allowed (not an illegal edge).
3. If target is `handed-off` (or any non-allowed argv): refuse PHASE-06 / unknown-target teaching error.
4. If current is `archived`: refuse (terminal).
5. If edge is illegal: refuse, print allowed next states.
6. If target is `simmering`: parse `--revisit` (§5). Refuse on miss/empty/bad grammar.
7. If current is `simmering` and target is `exploring`: run the wake algorithm (§7). Stop. Do not take the generic path.
8. Else: set `state`, set or clear `revisit`, bump `updated_date`, regen `index.md`, append `state` log line, commit under the protocol.
9. Print a short summary. Exit 0.

Success stdout (**Architect default**):

```text
state: <target>
revisit: <value-or-empty>
```

When the generic path was a wake, use the wake stdout (§8.2) instead.

### 8.2 `mycelium wake`

```text
mycelium wake [--dir PATH]
```

No `--revisit`. Extra unknown flags → teaching error.

Preflight: `state` must be `simmering` and `revisit` must parse. Else:

```text
mycelium: wake is only legal from simmering
convention: lifecycle
contract: program/contracts/lifecycle.md
fix: mycelium state simmering --revisit YYYY-MM-DD   # if you meant to park first
```

Then run §7. Success stdout:

```text
woke briefs/WAKE-YYYY-MM-DD.md
state: exploring
```

The first line is the path (tests may match the `briefs/WAKE-` prefix plus the clock date).

### 8.3 `mycelium index`

```text
mycelium index [--dir PATH]
```

Rebuild `index.md` from the current instance. Create the file if missing (the PHASE-01 repair). Protocol `op=index`. No log line. No skill emit. No state change.

Success stdout:

```text
wrote index.md
```

### 8.4 Teaching errors (PHASE-02 additions)

Four lines, stderr, exit 1.

```text
mycelium: illegal transition spark → clarified
convention: lifecycle
contract: program/contracts/lifecycle.md
fix: allowed next states: exploring, archived

mycelium: revisit "in two weeks" is not a date or event:<kebab>
convention: revisit
contract: program/contracts/revisit.md
fix: use YYYY-MM-DD (UTC) or event:after-iphone-launch

mycelium: state=handed-off requires a PHASE-06 handoff packet
convention: lifecycle
contract: program/contracts/lifecycle.md
fix: stay in clarified, or mycelium state archived; packet command is not shipped

```

The public `state exploring` path from simmering must write the brief (shared function). A silent internal flip is a Quality refuse, not a user-facing happy path.

### 8.5 Instance root resolution

Unchanged from PHASE-01 for single-instance commands (`state`, `wake`, `index`, `status` without `--all`, `check`, `new`, `tier`, `publish`).

`status --all` uses `--root` / `MYCELIUM_IDEAS_ROOT` / `~/ideas` (§9). It does not require cwd to be an instance.

## 9. `status` and `status --all`

Read-only. Do not create repos. Do not publish. Do not mutate instances. Do not take the instance lock unless a live lock is noticed (do not fail solely on a live lock; **Architect default:** ignore live locks, same spirit as check's notice).

```text
mycelium status [--dir PATH] [--all] [--root PATH] [--archived] [--offline]
```

| Flag | Default | Rule |
| --- | --- | --- |
| `--dir PATH` | walk-up instance | Single-instance root. With `--all`, also include this instance even if it sits outside `--root`. |
| `--all` | off | Portfolio mode. |
| `--root PATH` | `$MYCELIUM_IDEAS_ROOT` or `~/ideas` | Ideas root for local scan. Expand `~`. Do not mkdir. `--root` without `--all` → teaching error. |
| `--archived` | off | Include ideas whose state is `archived` **or** whose GitHub `isArchived` is true. Hidden otherwise. |
| `--offline` | off (or on if `MYCELIUM_OFFLINE=1`) | Local scan only. Never exec `gh`. |

`--offline` and a successful GitHub scan are opposites. `--offline` always wins.

### 9.1 `status` (no `--all`) — this instance

Requires an instance (walk-up or `--dir`). Exit 1 with a teaching error if none.

Stdout exactly six lines:

```text
slug: garden-lighting
state: simmering
tier: focused
revisit: 2026-08-22
due: yes|no|event
github: owner/name|unpublished
```

| Field | Rule |
| --- | --- |
| `slug` | manifest `slug` |
| `state` | manifest `state` |
| `tier` | manifest `tier` |
| `revisit` | manifest `revisit` (may be empty) |
| `due` | `event` if revisit is `event:<kebab>`; `yes` if date shape and clock date ≥ that date; `no` otherwise (including non-simmering / empty revisit) |
| `github` | `github_repo` if non-empty, else `unpublished` |

MS-201 unit clock `2026-08-07` + revisit `2026-08-08` → `due: no`.

### 9.2 `status --all` — local scan

Scan root = `--root` else `$MYCELIUM_IDEAS_ROOT` else `~/ideas` (home + `/ideas`).

Rules:

- Immediate children only (not recursive) that contain `mycelium.toml`.
- If cwd (or `--dir`) is an instance, include it even if it sits outside the root. Dedup by slug.
- Missing root directory → empty local set, not an error.
- Parse each `mycelium.toml`. Unreadable child → skip that child with a teaching line on **stderr** (do not fail the whole command). **Architect default:** continue; do not put a broken local dir in stdout.
- Master (`research-program.toml` only) is not an instance.

### 9.3 `status --all` — GitHub

`gh` via `internal/execrun` only.

**Architect default — prefer one search call:**

1. `gh api user --jq .login` → `<login>`
2. `gh search repos topic:idea user:<login> --json name,url,isArchived,owner`

If search is insufficient (command error): fall back to `gh search repos --owner <auth-user> --topic idea --json name,url,isArchived,owner`, then if still insufficient `gh repo list <auth-user> --limit 1000 --json name,url,isArchived` filtered by `gh api GET /repos/{owner}/{name}/topics`. Fake `Runner` in tests; never call real `gh` in hermetic tests.

Remote manifest:

```text
gh api repos/{owner}/{name}/contents/mycelium.toml --jq .content
```

Content is base64. Decode and parse. If missing or unreadable, **do not skip silently**. List the repo as `manifest: unread`:

```text
<name>	unread	-	-	remote	owner/name
```

`state` column token `unread` means `manifest: unread`. Count it in `n ideas`.

### 9.4 Merge

Merge key = local `slug` vs remote repo `name`.

| Case | `flag` | `github` column |
| --- | --- | --- |
| Local only | `unpublished` | `unpublished` (or local `github_repo` if set but not in the remote set — still `unpublished` if no remote match) |
| Remote only | `remote` | `owner/name` |
| Both | `ok` | `owner/name` |

State/tier/revisit: prefer **local** manifest when both exist. Remote-only uses the remote manifest (or `unread`).

Archived (local `state=archived` OR remote `isArchived`): omit unless `--archived`.

### 9.5 Offline / partial

Any of these → local scan only, **never pretend the list is complete**:

- `--offline`
- `MYCELIUM_OFFLINE=1`
- `gh` missing (`LookPath` fails)
- `gh` failing (non-zero, timeout, unauthenticated)

Stdout **first line**:

```text
partial: local-only (<reason>)
```

| Situation | `<reason>` |
| --- | --- |
| `--offline` or env | `offline` |
| `gh` not on PATH | `gh missing` |
| `gh` non-zero | `gh failed` |

Then the idea lines, then the summary with `partial`.

Exit 0. Partial is not a failure.

Hermetic tests MUST use `--offline` or `MYCELIUM_OFFLINE=1` and a fake execrun that records `gh` was not called.

### 9.6 Sort

Default sort buckets, then slug ASC inside a bucket:

1. overdue simmering (date shape, clock date > revisit)
2. due today (date shape, clock date == revisit)
3. exploring
4. spark
5. clarified
6. event-simmering
7. the rest (including future-dated simmering)
8. archived last, if shown

### 9.7 Output (parseable)

One line per idea:

```text
slug\tstate\ttier\trevisit\tflag\tgithub
```

Then one summary line:

```text
n ideas (k overdue, partial|complete)
```

`k overdue` counts date-shape simmering ideas with clock date > revisit (visible rows only).

`complete` only when the GitHub half ran and succeeded. Otherwise `partial`.

Zero ideas is exit 0:

```text
0 ideas (0 overdue, complete)
```

(or `partial` if offline). If offline, the `partial: local-only (...)` line still comes first.

See Appendix D.

### 9.8 Tests for GitHub half

**Architect default:** GitHub merge is unit-tested with a fake runner. One optional `//go:build github_integration` test may exist. Do **not** block MS-201 on `GH_TOKEN`. Do not add a PHASE-02 credentialed workflow as a gate.

## 10. Skills (spark, wake, portfolio) + mycelium-cli update

Emit path: `.agents/skills/<name>/SKILL.md` like `mycelium-cli`.

Source of truth: `program/skills/<name>/SKILL.md`. Scaffold copies into the instance. Re-run `go generate` in `internal/embed`.

### When skills are emitted

| Event | Emit spark/wake/portfolio? |
| --- | --- |
| New `mycelium new idea` scaffold (PHASE-02+) | **Yes** |
| `mycelium tier` | **No** |
| `mycelium index` | **No** |
| `mycelium state` / `mycelium wake` | **No** (do not retrofit) |
| One-shot `program/skills/` copy command | **Not a PHASE-02 command** |

**Architect default:** mutating `state`/`wake` does not retrofit skills. New scaffolds get them. Existing PHASE-01 instances: re-scaffold or copy skills **manually**. Document that sentence in `program/skeleton/AGENTS.md` and in `mycelium-cli` SKILL.md.

### `program/skills/spark/SKILL.md`

Front matter `name: spark`. First session inside a new instance:

1. Read `index.md` and `log.md` (not the whole tree).
2. `mycelium state exploring`
3. `mycelium new decision "…"` (first thought)
4. `mycelium check`

Do not invent a hub. Do not commit unless the human asks.

### `program/skills/wake/SKILL.md`

Front matter `name: wake`. On re-entry to a simmering idea:

1. Run `mycelium wake`
2. Read `briefs/LATEST.md`
3. Do **not** reread raw logs first
4. Then work. Run `mycelium check` before handing back.

### `program/skills/portfolio/SKILL.md`

Front matter `name: portfolio`. Cross-idea:

1. `mycelium status --all` (pass `--offline` when hermetic / no `gh`)
2. Interpret `partial: local-only (...)` as incomplete — do not invent remote ideas
3. Do not create repos. Do not publish. Do not mutate.

### `program/skills/mycelium-cli/SKILL.md` update

Add rows for `state`, `wake`, `status`, `status --all`, `index`. Keep the manual floor, teaching-error shape, and "do not git commit unless the human asks".

### Slice

Slice 7 lands the three skills + mycelium-cli update + embed generate. Update `program/skeleton/AGENTS.md` in the same slice to name the new commands and the manual skill-copy note for old instances.

## 11. Check updates (what changes from PHASE-01 conformance)

Replace the PHASE-01 stored-state bullets. Structure only (DEC-005). Runtime still reads instance files, never embed.

### Storage-rule bullets (new)

| # | Rule |
| --- | --- |
| 1 | `clarified` is **LEGAL**. |
| 2 | `handed-off` still **FAIL** (PHASE-06). Teaching error names the packet. |
| 3 | `simmering` ⇒ `revisit` matches §5 grammar (not merely non-empty). |
| 4 | `index.md` required (all tiers). Missing → FAIL, fix `mycelium index`. |
| 5 | `briefs/` is an allowed top-level path. |
| 6 | If the log contains a `wake` op, `briefs/LATEST.md` must exist and pass the five H2s. |
| 7 | `index.md` required H2s: State, Artifacts, Log tail, Wake. |
| 8 | Log ops extended: `state`, `wake` (keep `scaffold\|new\|tier\|publish\|check`). |
| 9 | Do not require a wake brief on instances that never simmered. |
| 10 | Do not grade brief prose. Do not require N artifacts. |

### Unchanged PHASE-01 checks

ID uniqueness, ID-to-path both ways, link resolution (now also `index.md` + `briefs/*.md`), front matter + H2s per schema, legal tier, tier binds, leftover journal/stale lock, declared deviations, stage-scoped ranges (DEC-013), teaching-error shape, `--abort-journal`.

### Lift timing

**Architect default:** Slice 1 is pure parsers only. Check storage-rule changes land in **Slice 3** with `state` / `wake` (so clarified becomes legal in the same PR that can reach it). Slice 2 may bind `index.md` required + H2s + `briefs/` allowed, because `mycelium index` and new scaffolds exist in that slice.

### `program/contracts/conformance.md`

Rewrite the lifecycle bullet and the allowed-top-level list. Keep the 11 must-implement checks; extend #2 (lifecycle) and #8 (log ops) and add index/wake bullets as additional must-implement items 12–14, or fold them into #2/#8/#11. **Architect default:** add items 12–14 (index present+H2s; briefs allowed; wake brief if wake op) so Quality can thermos a numbered list.

## 12. Vertical slices with build order, each checkable

PR-per-slice, sequential, rebase on main. Arvo merges Quality-green PRs. Engineering opens PRs. Engineering does NOT push to main. Prefer one live PR at a time. Do not stack unpublished slices on one branch unless Quality is backed up.

Each PR title: `PHASE-02 Slice N: <done-bar noun>`. Each PR body links this brief and the slice done bar. No drive-by refactors. No v1 deletions. No PHASE-03+ commands.

### Slice 0 — Commissioning (docs only)

This brief + lifecycle rewrite + new contracts + acceptance stub. No Go.

Land: `framework/phases/PHASE-02-implementation-brief.md`, `framework/phases/PHASE-02-acceptance.md` (rows = §15), rewritten `program/contracts/lifecycle.md`, new `program/contracts/{revisit,index,wake,status}.md`, updates to `conformance.md` (may finish in Slice 2/3) and `operation-protocol.md` (`state`/`wake`/`index` on the mutating list + journal `op` enum).

Done: files exist on a PR. Quality reads them against this brief. No product code.

### Slice 1 — Pure parsers (no check-rule change yet)

`internal/revisit` (grammar, due/overdue, trigger-date extract). `internal/lifecycle` (transition table; move `check.LegalNext` here or thin-alias). Table-driven tests only.

**Architect default:** Slice 1 does **not** lift the clarified fail and does **not** bind revisit grammar in check. Check rules change in Slice 3 with the command (index bind may land in Slice 2).

Done: `go test` for parse / due / overdue / extract / legal edges. No CLI yet.

### Slice 2 — `index.md` generate + skeleton + check bind

- `internal/indexmd` renderer.
- `program/skeleton/index.md` stub + scaffold writes concrete index.
- `mycelium index`.
- Existing mutating commands (`new <type>`, `tier`, `publish`) regenerate index.
- `index.md` added to all three tier `binds`.
- `index.md` and `briefs/` added to allowed top-level.
- Check: missing `index.md` FAIL; required H2s; `briefs/` allowed.
- New scaffolds include `index.md`.

Done (hermetic): `new idea --offline` has `index.md` with the four H2s; `check` 0; delete `index.md` → `check` 1 → `mycelium index` → `check` 0; `new decision` regenerates counts.

### Slice 3 — `state` + `wake` + brief writer + storage-rule update

- `internal/wakebrief`, `internal/statecmd`.
- `mycelium state`, `mycelium wake`.
- Check: clarified legal; handed-off still fail; simmering ⇒ grammar; wake op ⇒ `briefs/LATEST.md` H2s; log ops include `state`/`wake`; link scan includes briefs + index.
- Silent wake impossible (shared function).

Done (hermetic): legal edges; simmer without `--revisit` refuse; `handed-off` refuse; `state exploring` from simmering writes the same brief as `wake`; archived refuse; `check` 0 after wake.

### Slice 4 — `mycelium status` (single instance) + due/overdue

Done (hermetic): fixture at clock 2026-08-07 + revisit 2026-08-08 prints `due: no`; clock 2026-08-08 prints `due: yes`; event prints `due: event`; unpublished github line.

### Slice 5 — `status --all` local scan + partial/offline

Temp ideas root with 2 instances. `--offline`. First line `partial: local-only (...)`. `gh` never called. Sort order covered by unit tests on the merge/sort function.

Done (hermetic): two-instance root; cwd-outside-root instance included; missing root → 0 ideas partial; no mutation.

### Slice 6 — `status --all` GitHub merge

Fake execrun. Programmed `gh api user`, `gh search repos`, `gh api …/contents/mycelium.toml`. Merge flags `ok` / `unpublished` / `remote`. Unreadable remote → `unread`. Archived hidden unless `--archived`.

**Architect default:** GitHub half is unit-tested with a fake runner; one optional build-tag test; do not block MS-201 on `GH_TOKEN`. No new credentialed workflow.

Done: fake-runner tests green. Hermetic `go test ./...` never calls real `gh`.

### Slice 7 — Skills + mycelium-cli + embed generate

Three skills + mycelium-cli update + `program/skeleton/AGENTS.md` + `go generate` embed.

Done: new scaffold emits `.agents/skills/{spark,wake,portfolio,mycelium-cli}/SKILL.md`; `tier` / `index` / `state` do not retrofit a fixture that lacks them.

### Slice 8 — MS-201 fixture test (`go test` only)

Hermetic fixture (§13, Appendix E) lives in `go test ./...` with `MYCELIUM_OFFLINE=1`. Do **not** add `.github/workflows/phase-02-*.yml`. Do not extend `phase-01-hermetic.yml` as a phase gate. Actions is not a gate.

Done: `go test ./...` runs the MS-201 fixture green. `gh` never invoked.

## 13. Done / verified mapped onto MS-201

MS-201 is the hermetic phase gate (injectable clock + known fixture). **Dogfood 7 real days is out of the gate** — human evidence for Arvo, not Engineering. Quality refuses any PR that makes the real-day dogfood a gate, or that adds an Actions job as the MS-201 gate.

### MS-201 expected (authoritative; recipe in Appendix E)

Clock timeline UTC:

| When | Clock | Actions |
| --- | --- | --- |
| T0 | 2026-08-01 | scaffold `--offline` `"Wake Fixture"` → slug `wake-fixture`, state `spark` |
| T0 | 2026-08-01 | `state exploring` |
| T0 | 2026-08-01 | `new decision "Park the idea"` |
| T0 | 2026-08-01 | `new assumption "API stays stable"` then edit: status `Held`; Revisit Triggers contains `2026-08-05` |
| T0 | 2026-08-01 | `new evidence "Vendor changelog"` then edit: Revalidation Trigger contains `2026-08-06` |
| T0 | 2026-08-01 | `new assumption "Budget is unlimited"` then edit: status `Retired` (must NOT be cited) |
| T0 | 2026-08-01 | `state simmering --revisit 2026-08-08` |
| T1 | 2026-08-09 | `mycelium wake` (≥7 days after T0; revisit 2026-08-08 is overdue) |

Expected after T1:

| Check | Value |
| --- | --- |
| `state` | `exploring` |
| `revisit` | `""` |
| `briefs/WAKE-2026-08-09.md` | exists |
| `briefs/LATEST.md` | exists, identical bytes to the dated file |
| brief H2s | Parked; Log since simmer; Evidence triggers; Assumptions; Suggested next |
| body contains | `ASM-001` and `EVD-001` |
| body does not require | `ASM-002` (Retired) |
| body cites | simmer date `2026-08-01` **or** revisit `2026-08-08` (at least one) |
| `mycelium check` | exit 0 |
| `gh` | never invoked |

Also a unit / hermetic case on the **same fixture** with clock `2026-08-07` (before the revisit date):

| Check | Value |
| --- | --- |
| `mycelium status` | `due: no` |
| `mycelium wake` | still legal (human can wake early) |
| brief cites | `EVD-001` (trigger 2026-08-06 due) and `ASM-001` (Held + 2026-08-05 due) |

### Slice → MS-201 map

| Slice | MS-201 clause |
| --- | --- |
| 0 | commissioning; not a runtime clause |
| 1 | due/overdue + trigger extract used by the fixture |
| 2 | index present so check can pass after wake |
| 3 | `wake` writes the brief; storage rules allow the post-wake instance |
| 4 | `status due:no` on the early-clock case |
| 5–6 | not required for MS-201 (portfolio); must not break hermetic `go test` |
| 7 | not required for MS-201; new scaffolds only |
| 8 | the fixture test in `go test ./...` **is** the gate |

PHASE-02 is accepted when MS-201 is green in `go test ./...` on main. Arvo accepts the phase. Engineering does not self-accept. The dogfood wake is Arvo/Robert evidence after acceptance, not a blocker.

## 14. Automated test plan

Engineering MUST write these tests. Quality thermos against this list. Do NOT require Playwright, Docker, live GitHub, or `GH_TOKEN` in default `go test ./...`.

### Unit (no network, no gh)

| Area | Cases |
| --- | --- |
| revisit parse | date, event, empty, garbage, `2026-02-30`, `EVENT:foo`, `event:After`, `2026-08-08T00:00:00Z` |
| due / overdue | clock 07/08/09 vs date 08; event never overdue |
| trigger extract | first date in section; leading whitespace; date in a later section ignored; missing date ≠ due; `2026-08-06T…` does not match |
| transition table | every legal edge; every illegal edge; archived terminal; handed-off refused as a target |
| index render | H1 + four H2s; artifact counts; last 20 log lines; Wake `none` vs `briefs/LATEST.md` |
| brief citation set | IDs only: Held+date ASM, Open ASM, Retired no-date excluded, EVD date due included, EVD future excluded |
| status sort | overdue, due-today, exploring, spark, clarified, event, rest, archived |
| merge flags | local-only `unpublished`, remote-only `remote`, both `ok`, unread remote |
| partial reason | `offline`, `gh missing`, `gh failed` |

### Hermetic CLI (built binary, temp dirs, `MYCELIUM_OFFLINE=1` or `--offline`)

| Case | Expect |
| --- | --- |
| state edges | legal succeed; illegal teaching error names `lifecycle.md` and lists next states |
| simmer without `--revisit` | refuse |
| `--revisit` on `clarified` | refuse |
| `state handed-off` | refuse, names PHASE-06 packet |
| stored `handed-off` | `check` FAIL |
| stored `clarified` after `state clarified` | `check` 0 |
| `wake` from spark/exploring | refuse |
| `wake` from simmering | writes dated brief + LATEST; state exploring; revisit `""` |
| `state exploring` from simmering | **same brief** (not silent) |
| second wake same UTC day | overwrites that day's file and LATEST |
| `index` rebuild | restores missing `index.md`; `check` 0 |
| `status` | six-line format |
| `status --all --offline` | temp ideas root with 2 instances; first line `partial: local-only (...)`; `gh` not called |
| MS-201 fixture | Appendix E; clock T0 then T1; citations; `check` 0 |
| MS-201 early wake | clock 2026-08-07; `due: no`; wake legal; cites EVD-001 + ASM-001 |
| fake execrun | `gh` not called under offline |

Hermetic tests MUST use `--offline` or `MYCELIUM_OFFLINE=1` and a fake `internal/execrun.Runner` that records `(name, args, dir)` and asserts `gh` was not called. Slice 6 unit tests program fake `gh` results and still do not open a network.

Do not require live GitHub for `go test ./...`. Optional `//go:build github_integration` only. Not a phase gate. Do not commission a `GH_TOKEN` job.

## 15. Acceptance matrix / in-repo contract paths

Slice 0 lands the paths listed in §12 Slice 0. Later slices also land: `program/skeleton/index.md`; `program/skills/{spark,wake,portfolio}/SKILL.md`; updated `program/skills/mycelium-cli/SKILL.md`; `index.md` bind on `program/tiers/*.toml`. No workflow file. No DEC-015 file.

### Acceptance matrix rows (copy into `PHASE-02-acceptance.md`)

Each row: id, check, evidence, owner (Engineering | CI | Arvo).

| id | check | evidence |
| --- | --- | --- |
| A-S0 | Slice 0 files exist and match this brief | PR diff; Quality read |
| A-S1 | revisit / due / overdue / trigger / table tests green | `go test` |
| A-S2 | new scaffold emits `index.md`; `mycelium index` repairs; check binds | hermetic CLI |
| A-S3 | `state` / `wake` edges; brief written; clarified legal; handed-off fails | hermetic CLI |
| A-S4 | single `status` + due/overdue | hermetic CLI |
| A-S5 | `status --all --offline` two-instance root; partial line; no gh | hermetic CLI |
| A-S6 | merge flags + unread + archived filter via fake runner | `go test` |
| A-S7 | three skills emitted on new scaffold; no retrofit on state/index | hermetic CLI |
| A-S8 | MS-201 fixture green | `go test ./...` |
| MS-201 | all §13 expected bullets | A-S8 (uses 1–4) |
| DOGFOOD-7d | one real 7-day wake | Arvo / Robert; **not the gate** |

## 16. Decided / Architect defaults

No TBD. Open items are decided inline and labeled **Architect default**. Do not bikeshed them in a code PR. Write a DEC to change one. No DEC-015 is required for these.

Index of defaults that are easy to miss:

- `state` is the one transition command. `wake` is a ritual alias, not a second machine.
- `state exploring` from simmering writes the same brief as `wake`. Log op is `wake`.
- Allowed targets: exploring, simmering, clarified, archived. `handed-off` always refused (PHASE-06).
- `clarified` is a legal stored state. PHASE-01 fail rule lifted in Slice 3, not Slice 1.
- Slice 1 = pure parsers only. Check storage-rule change in Slice 3. Index bind in Slice 2.
- Due on the revisit date; overdue strictly after (next UTC day).
- Event shape never auto-overdue. `due: event`.
- `--revisit` required on `state simmering`; forbidden on other targets. Same-state simmering may update revisit.
- Trigger date: first `^\s*(\d{4}-\d{2}-\d{2})\b` in the named section. No NLP. Missing date ≠ due.
- Cite every Open or Held ASM plus any ASM/EVD with a due/overdue date in the named section.
- Simmer line = most recent `state` log line whose note starts with `simmering`. Collect later lines in file order.
- Wake paths: `briefs/WAKE-YYYY-MM-DD.md` (clock UTC date, overwrite same day) + `briefs/LATEST.md`.
- `index.md` required all tiers. `briefs/` allowed. `mycelium index` does not emit skills and does not log.
- `state`/`wake` do not retrofit skills. No one-shot skills-copy command.
- `status --all` search: `gh api user --jq .login` then `gh search repos topic:idea user:<login> --json name,url,isArchived,owner`. Fallback only if search is insufficient.
- Offline / missing / failing gh → first line `partial: local-only (<reason>)`. Never pretend complete.
- Merge by slug vs remote name. Flags: `unpublished` / `remote` / `ok`. Unreadable remote: state token `unread` (`manifest: unread`).
- Sort: overdue simmering, due today, exploring, spark, clarified, event-simmering, rest, archived last. Slug ASC inside bucket.
- Local scan: `--root` or `$MYCELIUM_IDEAS_ROOT` or `~/ideas`. Immediate children with `mycelium.toml`. Include cwd/`--dir` instance even if outside root.
- GitHub half: fake runner; optional build-tag; do not block MS-201 on `GH_TOKEN`.
- Do **not** add a PHASE-02 Actions workflow. Do not make phase-01-github.yml a dependency. Actions is not a gate.
- methodology_version `2.0.0`; CLI `0.1.0-dev`. Journal ops add `state`\|`wake`\|`index`. Commit: brief(s), index, log, manifest last. Link scan adds `index.md` + `briefs/*.md`. `2026-08-06T…` is not a trigger date.
- Future-dated simmering sorts in "the rest". Check does not fail leftover `revisit` on non-simmering states. No `mycelium range`. No DEC-015.

## 17. Risks, rollback, what Quality should refuse

### Risks and mitigations

| Risk | Mitigation |
| --- | --- |
| Silent wake (`state exploring` from simmering skips the brief) | Shared function. Hermetic test. Quality refuse. |
| `status --all` pretends complete when `gh` failed | Mandatory `partial: local-only (...)` first line. Test the reasons. |
| Hermetic tests call network / real `gh` | Fake runner + `--offline` / `MYCELIUM_OFFLINE=1`. Slice 8 asserts no `gh`. |
| MS-201 blocked on `GH_TOKEN` | MS-201 is hermetic. Slice 6 is fake-runner. No new credentialed job. |
| Clarified fail lifted before `state` exists | Lift in Slice 3 with the command. |
| Old instances fail check (no `index.md`) | `mycelium index` repair. Documented. Do not silently migrate in `state`. |
| Skills missing on old instances | Do not retrofit. Manual copy / re-scaffold. |
| Content grading of Suggested next | DEC-005. Check H2s + ID resolution only. |
| Dogfood 7 real days treated as the gate | §13 says it is not. Quality refuse. |
| Reopening DEC-014 / publish / MS-101(b) | §2. Quality refuse. |
| `handed-off` succeeds without a packet | Command refuse + check FAIL. |
| Growing `latinFold` or adding `x/text` | DEC-014. Do not touch `internal/slug`. |
| Emitting `framework/` or deleting Justfile | Absence tests stay. Master v1 stays. |

### Rollback

Revert the offending PR on master. Do not `git push --force` to main. Floor is `75645c3c2b48cd485a590cb0f0158d7cb29da1df`. Do not delete Justfile/scripts as a "cleanup" rollback.

### Quality should refuse

Refuse to approve if:

- any PHASE-03+ skill or command is shipped (`thinking-mode`, council, `supersede`, handoff packet)
- `handed-off` succeeds without a packet, or check stops failing stored `handed-off`
- silent wake (simmering → exploring without a brief)
- `status --all` pretends complete when `gh` failed or `--offline` is set
- hermetic tests call network or real `gh`
- cobra / viper / yaml / testify / go-github / `golang.org/x/text` appears
- `latinFold` grows or NFKD is implemented (DEC-014)
- MS-101(b) is implemented or a new `GH_TOKEN` job is commissioned as a gate
- an Actions job is added as the MS-201 gate
- Justfile/scripts deleted from master
- `framework/` is emitted into an instance
- CLI git-commits instance work product
- content grading of wake briefs
- dogfood 7 real days is required as the gate
- PR pushed straight to main
- `just init` was run on master, or `research-program.toml` was renamed
- DEC-012 / DEC-013 / DEC-014 reopened in a code PR
- `explore` / `simmer` added as separate verbs
- a destroy command appears

## 18. Overnight execution order

Same order as §12 (slices 0→8). PR-per-slice, sequential, rebase on main. One live PR at a time. Slice 3 is the largest — do not combine it with 4–8. Slice 8 must be green in `go test ./...` on its PR (not Actions).

Title: `PHASE-02 Slice N: <done-bar noun>`. Body links this brief and the slice done bar. No drive-by refactors, v1 deletions, or PHASE-01 leftover work. Engineering opens PRs; Arvo merges Quality-green PRs; Engineering does NOT push to main.

Cursor cloud env name is exactly `robertguss/mycelium`.

## 19. Handoff

### What Engineering starts with

This file. Only this file. Clone `https://github.com/robertguss/mycelium` at `75645c3c2b48cd485a590cb0f0158d7cb29da1df` (current `main` at pin time). Read `framework/blueprint.md` and DEC-001–014 for authority, not for a second plan. Execute Slice 0 first.

Cursor cloud: env `robertguss/mycelium`. Go 1.26. `CGO_ENABLED=0`.

### What Engineering must not do

See §17 (Quality should refuse) and §16 (Architect defaults). Do not open a design debate in the PR. Do not write a second brief. Do not write DEC-015 unless Arvo orders a reopen.

### What Quality reads

This brief, the acceptance matrix, the rewritten lifecycle contract, the new revisit/index/wake/status contracts, and the PR diff. Thermos: §14 tests exist and match; §17 refuse list is clean; MS-201 hermetic; no `GH_TOKEN` gate.

### What Arvo does

Merges Quality-green PRs. Accepts PHASE-02 when MS-201 is green on main. Optionally records a 7-real-day dogfood wake as human evidence after acceptance.

## Appendix A — No new DEC

No DEC-015. PHASE-02 does not reopen DEC-012, DEC-013, or DEC-014. Remaining choices are Architect defaults in §16. Engineering lands **zero** new files under `framework/decisions/`.

If a later discovery contradicts a locked item, stop and write a DEC; do not silently patch this brief in a code PR.

## Appendix B — Wake brief template

Generator output shape. Tokens are **not** used; the writer fills concrete IDs and dates. Prose under Suggested next is a stub; check does not grade it.

`briefs/WAKE-2026-08-09.md` (and identical `briefs/LATEST.md`):

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

The dated heading is optional extra (not a required H2). Required H2s are the five `##` names. Citation IDs in the body must resolve.

A second wake on 2026-08-09 overwrites both files.

## Appendix C — `index.md` example

Post-wake `wake-fixture` instance (concrete values, no tokens):

```text
# Wake Fixture

## State

state: exploring
tier: focused
revisit:
github_repo:

## Artifacts

DEC: 1
ASM: 2
EVD: 1
(other NS: 0)

## Log tail

2026-08-01	scaffold	-	Wake Fixture
2026-08-01	state	-	exploring
2026-08-01	new	DEC-001	Park the idea
2026-08-01	new	ASM-001	API stays stable
2026-08-01	new	EVD-001	Vendor changelog
2026-08-01	new	ASM-002	Budget is unlimited
2026-08-01	state	-	simmering revisit=2026-08-08
2026-08-09	wake	-	exploring

## Wake

briefs/LATEST.md
```

Exact log-line notes may differ if a command writes a slightly different note string; the **op** column and simmer note prefix `simmering` are binding. Check validates H2s, not these exact counts, except that the renderer must be deterministic for a given instance.

At spark, `## Wake` is `none`.

## Appendix D — `status --all` output format

### Single instance (`mycelium status`)

```text
slug: garden-lighting
state: simmering
tier: focused
revisit: 2026-08-22
due: yes
github: unpublished
```

### Portfolio, offline

```text
partial: local-only (offline)
garden-lighting	simmering	focused	2026-08-08	unpublished	unpublished
wake-fixture	exploring	focused		unpublished	unpublished
2 ideas (1 overdue, partial)
```

(Sort depends on clock. Example assumes clock `2026-08-09`, first idea overdue simmering, second exploring.)

### Portfolio, GitHub succeeded

No `partial:` line.

```text
garden-lighting	simmering	focused	2026-08-08	ok	robertguss/garden-lighting
other-idea	unread	-	-	remote	robertguss/other-idea
local-only	spark	focused		unpublished	unpublished
3 ideas (1 overdue, complete)
```

`unread` in the state column means `manifest: unread`.

### Portfolio, `gh` failed

```text
partial: local-only (gh failed)
garden-lighting	simmering	focused	2026-08-08	unpublished	unpublished
1 ideas (1 overdue, partial)
```

Reasons: `offline` | `gh missing` | `gh failed`.

Columns: `slug`, `state`, `tier`, `revisit`, `flag`, `github` separated by one tab. `flag` ∈ `ok` | `unpublished` | `remote`.

## Appendix E — MS-201 fixture recipe

Hermetic. No network. `MYCELIUM_OFFLINE=1`. Fake execrun. Clock via `MYCELIUM_NOW` (RFC3339) or injected `clock.Fixed`.

Work in a temp dir. Binary = freshly built `mycelium`.

### T0 — `MYCELIUM_NOW=2026-08-01T00:00:00Z`

```text
mycelium new idea "Wake Fixture" --offline --dir PATH/wake-fixture
# slug = wake-fixture, state = spark, created_date = 2026-08-01

mycelium state exploring --dir PATH/wake-fixture

mycelium new decision "Park the idea" --dir PATH/wake-fixture
# decisions/DEC-001-park-the-idea.md

mycelium new assumption "API stays stable" --dir PATH/wake-fixture
# assumptions/ASM-001-api-stays-stable.md
# EDIT front matter: status = "Held"
# EDIT ## Revisit Triggers body to contain a line: 2026-08-05

mycelium new evidence "Vendor changelog" --dir PATH/wake-fixture
# evidence/EVD-001-vendor-changelog.md
# EDIT ## Revalidation Trigger body to contain a line: 2026-08-06

mycelium new assumption "Budget is unlimited" --dir PATH/wake-fixture
# assumptions/ASM-002-budget-is-unlimited.md
# EDIT front matter: status = "Retired"
# Do not put a YYYY-MM-DD in Revisit Triggers

mycelium state simmering --revisit 2026-08-08 --dir PATH/wake-fixture
```

Edits may use the test helper / stdlib. Do not add a `mycelium edit` command.

ASM-001 after edit: front matter `status = "Held"`; `## Revisit Triggers` body contains a line `2026-08-05` (other required H2s may stay as template placeholders).

EVD-001 after edit: `## Revalidation Trigger` body contains a line `2026-08-06`.

ASM-002 after edit: `status = "Retired"`; no `YYYY-MM-DD` in Revisit Triggers.

### Early-clock case — `MYCELIUM_NOW=2026-08-07T00:00:00Z`

On a **copy** of the T0-after-simmer fixture (do not consume the T1 fixture):

```text
mycelium status --dir PATH/wake-fixture
# due: no

mycelium wake --dir PATH/wake-fixture
# legal; writes briefs/WAKE-2026-08-07.md
# body contains ASM-001 and EVD-001
```

### T1 — `MYCELIUM_NOW=2026-08-09T00:00:00Z`

On the T0-after-simmer fixture (not the early-wake copy):

```text
mycelium wake --dir PATH/wake-fixture
```

Assert:

1. manifest `state == exploring`, `revisit == ""`
2. `briefs/WAKE-2026-08-09.md` exists
3. `briefs/LATEST.md` exists and `bytes.Equal` the dated file
4. both files contain H2s Parked, Log since simmer, Evidence triggers, Assumptions, Suggested next
5. body contains `ASM-001` and `EVD-001`
6. body need not contain `ASM-002` (assert absence **or** simply do not require it — **Architect default:** assert `ASM-002` is absent from the brief)
7. body contains `2026-08-01` or `2026-08-08` (at least one)
8. `mycelium check --dir PATH/wake-fixture` exits 0
9. fake runner: `gh` never invoked

**Architect default:** test lives at `internal/clitest/ms201_hermetic_test.go` (execs the binary). Unit citation-set tests live in `internal/wakebrief`.

## Appendix F — Target file tree additions

### Master (additions on top of the PHASE-01 tree; v1 files retained)

```text
internal/revisit/          # parse, due/overdue, trigger extract
internal/lifecycle/        # transition table
internal/indexmd/          # render index.md
internal/wakebrief/        # citation collect + brief write
internal/statecmd/         # mycelium state + wake
internal/statuscmd/        # mycelium status / status --all
internal/indexcmd/         # mycelium index
internal/clitest/ms201_hermetic_test.go
program/contracts/lifecycle.md          # rewritten
program/contracts/revisit.md            # new
program/contracts/index.md              # new
program/contracts/wake.md               # new
program/contracts/status.md             # new
program/contracts/conformance.md        # updated
program/contracts/operation-protocol.md # updated
program/skeleton/index.md
program/skills/spark/SKILL.md
program/skills/wake/SKILL.md
program/skills/portfolio/SKILL.md
program/skills/mycelium-cli/SKILL.md    # updated
program/tiers/{focused,standard,high-assurance}.toml  # index.md bind
framework/phases/PHASE-02-implementation-brief.md
framework/phases/PHASE-02-acceptance.md
internal/embed/program/                 # regenerate after program/ edits
```

Do **not** add `framework/decisions/DEC-015-*.md`, a PHASE-02 workflow, a GitHub-credential job, or delete Justfile/scripts/`research-program.toml`/PHASE-01 workflows.

### Emitted instance (spark / focused, local-only, PHASE-02 scaffold)

```text
README.md  mycelium.toml  log.md  index.md  CONTEXT.md  AGENTS.md  .gitignore
.agents/skills/{mycelium-cli,spark,wake,portfolio}/SKILL.md
program/ …
.git/          # init only; no commit
```

Absent: `framework/`, `cmd/`, `internal/`, `go.mod`, `Justfile`, `scripts/`, `research-program.toml`, v1 `research-*` skills.

After a wake:

```text
briefs/WAKE-YYYY-MM-DD.md
briefs/LATEST.md
```

`index.md` `## Wake` becomes `briefs/LATEST.md`.
`mycelium.toml` `state=exploring`, `revisit=""`.

Unexported helpers may live next to their tests. No `pkg/`. No extra public command packages. Do not touch `internal/slug`.

End of PHASE-02 implementation brief. Engineering executes from this file only.
