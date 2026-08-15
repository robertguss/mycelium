# Status / Portfolio (Mycelium 2.0)

Read-only. Never mutates. Never creates repos. Never publishes. Ignore live
locks (do not fail solely on a live lock). Does **not** use the operation
protocol.

```text
mycelium status [--dir PATH] [--all] [--root PATH] [--archived] [--offline]
```

## Flags

| Flag | Default | Rule |
| --- | --- | --- |
| `--dir PATH` | walk-up instance | Single-instance root. With `--all`, also include this instance even if it sits outside `--root`. |
| `--all` | off | Portfolio mode. |
| `--root PATH` | `$MYCELIUM_IDEAS_ROOT` or `~/ideas` | Ideas root for local scan. Expand `~`. Do not mkdir. `--root` without `--all` → teaching error. |
| `--archived` | off | Include ideas whose state is `archived` **or** whose GitHub `isArchived` is true. Hidden otherwise. |
| `--offline` | off (or on if `MYCELIUM_OFFLINE=1`) | Local scan only. Never exec `gh`. |

`--offline` always wins over a successful GitHub scan.

## Env

| Env | Effect |
| --- | --- |
| `MYCELIUM_OFFLINE=1` | Treat as `--offline`. Never exec `gh`. Never open network. |
| `MYCELIUM_NOW` | RFC3339 clock override. |
| `MYCELIUM_IDEAS_ROOT` | Default ideas root for `status --all` when `--root` is absent. |

## Single instance (no `--all`)

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

MS-201: clock `2026-08-07` + revisit `2026-08-08` → `due: no`.

## `status --all` — local scan

Scan root = `--root` else `$MYCELIUM_IDEAS_ROOT` else `~/ideas`.

- Immediate children only (not recursive) that contain `mycelium.toml`.
- If cwd (or `--dir`) is an instance, include it even if outside the root.
  Dedup by slug.
- Missing root directory → empty local set, not an error.
- Unreadable child → teaching line on **stderr**, skip, continue. Do not put a
  broken local dir in stdout.
- Master (`research-program.toml` only) is not an instance.

## `status --all` — GitHub

`gh` via `internal/execrun` only.

Prefer:

1. `gh api user --jq .login` → `<login>`
2. `gh search repos topic:idea user:<login> --json name,url,isArchived,owner`

Fallbacks if search is insufficient (command error):
`gh search repos --owner <auth-user> --topic idea --json name,url,isArchived,owner`,
then if still insufficient
`gh repo list <auth-user> --limit 1000 --json name,url,isArchived` filtered by
`gh api GET /repos/{owner}/{name}/topics`.

Remote manifest:

```text
gh api repos/{owner}/{name}/contents/mycelium.toml --jq .content
```

Content is base64. Decode and parse. If missing or unreadable, **do not skip
silently**. List as unread:

```text
<name>	unread	-	-	remote	owner/name
```

`state` column token `unread` means `manifest: unread`. Count it in `n ideas`.

## Merge

Merge key = local `slug` vs remote repo `name`.

| Case | `flag` | `github` column |
| --- | --- | --- |
| Local only | `unpublished` | `unpublished` (or local `github_repo` if set but not in the remote set — still `unpublished` if no remote match) |
| Remote only | `remote` | `owner/name` |
| Both | `ok` | `owner/name` |

State/tier/revisit: prefer **local** manifest when both exist. Remote-only uses
the remote manifest (or `unread`).

Archived (local `state=archived` OR remote `isArchived`): omit unless
`--archived`.

## Offline / partial

Any of these → local scan only, **never pretend complete**:

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

Then the idea lines, then the summary with `partial`. Exit 0. Partial is not a
failure.

## Sort

Buckets, then slug ASC inside a bucket:

1. overdue simmering (date shape, clock date > revisit)
2. due today (date shape, clock date == revisit)
3. exploring
4. spark
5. clarified
6. event-simmering
7. the rest (including future-dated simmering)
8. archived last, if shown

## Output

One line per idea:

```text
slug\tstate\ttier\trevisit\tflag\tgithub
```

Then one summary line:

```text
n ideas (k overdue, partial|complete)
```

- `k overdue` counts date-shape simmering ideas with clock date > revisit
  (visible rows only).
- `complete` only when the GitHub half ran and succeeded. Otherwise `partial`.
- Zero ideas is exit 0: `0 ideas (0 overdue, complete)` (or `partial` if
  offline). If offline, the `partial: local-only (...)` line still comes first.

`flag` ∈ `ok` | `unpublished` | `remote`.

## Examples

### Single instance

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

(Sort depends on clock. Example assumes clock `2026-08-09`.)

### Portfolio, GitHub succeeded

No `partial:` line.

```text
garden-lighting	simmering	focused	2026-08-08	ok	robertguss/garden-lighting
other-idea	unread	-	-	remote	robertguss/other-idea
local-only	spark	focused		unpublished	unpublished
3 ideas (1 overdue, complete)
```

### Portfolio, `gh` failed

```text
partial: local-only (gh failed)
garden-lighting	simmering	focused	2026-08-08	unpublished	unpublished
1 ideas (1 overdue, partial)
```

## Tests

GitHub merge is unit-tested with a fake runner. One optional
`//go:build github_integration` test may exist. Do **not** block MS-201 on
`GH_TOKEN`. Do not add a PHASE-02 credentialed workflow as a gate.
