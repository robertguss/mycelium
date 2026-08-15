# Operation Protocol (Mycelium 2.0)

Applies to mutating multi-file commands:

- `mycelium new idea` (local emit)
- `mycelium new <type>`
- `mycelium tier`
- `mycelium publish`
- `mycelium state`
- `mycelium wake`
- `mycelium index`

Does **not** apply to `mycelium version`, `mycelium status` (read-only; never
mutates), or `mycelium check`, except `mycelium check --abort-journal` which
clears an interrupted journal (no confirmation step; see
`program/contracts/conformance.md`).

## Steps

1. **Preflight** — validate that the manifest and log parse, the schema
   resolves (when relevant), and the target path is free, before anything is
   written.
2. **Lock** — take exclusive repository lock at `.mycelium/lock` (flock).
   Record `pid` and `started`. Non-blocking acquire. Stale lock = recorded PID
   is dead.
3. **Stage** — write outputs under `.mycelium/stage/<op-id>/` and record intent
   in the journal.
4. **Journal** — `.mycelium/journal.json` with `schema_version = 1` and
   `op` ∈ `scaffold` | `new` | `tier` | `publish` | `state` | `wake` | `index`.
5. **Commit** — atomic renames in fixed order:
   - Artifact generation: artifact file(s), then log, then manifest **last**.
   - Scaffold: skeleton files, then `program/`, then log, then manifest **last**.
   - `state` / `wake`: brief + `briefs/LATEST.md` (if wake), then `index.md`,
     then `log.md`, then `mycelium.toml` **last**.
   - `index`: regenerate `index.md` only (no log line, no skill emit, no state
     change).
6. **Rollback** — failure before the first rename removes staged files and
   changes nothing. After a partial commit the journal survives.
7. **Resume** — re-running the same argv resumes under `original_id`; do not
   allocate a new ID.

## Git

- `git init -b main` runs **after** file commit.
- If `git init` fails, do **not** roll back emitted files.
- The CLI never `git add` or `git commit` instance work product.

## Publish

- Invoke `gh` after lock and before the manifest rename that writes
  `github_repo`.
- Publish is idempotent (re-run must not create a second repo or corrupt
  state).

## Filesystem floor

Local filesystem with atomic rename. Network filesystems are outside the floor;
lock + journal still bound damage.
