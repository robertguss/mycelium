# Agent Rules — {{PROJECT_NAME}}

Repository-local rules for humans and agents operating this research program.

## Authority

1. Git-tracked artifacts are authoritative. Chat history and model memory are not.
2. Precedence: accepted `DEC-###` → Blueprint → Charter → current stage prompt →
   revised specification → research reports (evidence) → reviews (proposals) →
   revised plan → `research-program.toml` (index only).
3. Details: `program/contracts/authority-and-precedence.md`.

## Fresh sessions

Every **substantive** stage runs in a fresh session with a self-contained
attachment manifest. Do not execute multiple substantive stages in one context.
Preparing prompts, manifests, and mechanical fixes is allowed in the current
session.

## Allowed file scope

Each stage may modify only the paths declared in its commissioning prompt and
manifest outputs. Do **not** silently edit governing artifacts (Blueprint,
Charter, accepted specs/plans) outside a commissioned revision stage.

## Validation and acceptance

- Placeholders (`Status: Placeholder — not accepted`) never unlock work.
- Independent validation before acceptance (`program/contracts/validation.md`).
- Validators fix mechanical issues only; no invented research.
- Human approval gates: `program/operator/approval-gates.md`.
- Humans own git; do not mark stages accepted without human approval and commit
  recording in the manifest.

## Identifiers and citations

- Stable IDs: `DEC`, `REC`, `REQ`, `FND`, `RSK`, `OQ`, `SPK`, `PHASE`, `MS`.
- Never reuse IDs. Disposition upstream IDs explicitly.
- Portable citations only (Markdown links, footnotes, source ledgers).

## Evidence

- Evidence before confidence.
- Evidence Ledgers on focused reports.
- Bounded spikes when load-bearing claims are testable.
- No popularity-as-proof; no silent recommendation loss.

## Skills

Portable skills live under `.agents/skills/`:

| Skill               | Use for                                                                |
| ------------------- | ---------------------------------------------------------------------- |
| `research-program`  | Discovery, resume, next stage, program orchestration                   |
| `research-stage`    | Just-in-time stage package (prompt, install, attach, launch, validate) |
| `research-validate` | Independent validation gate                                            |

Methodology library: `program/`. Operator start: `program/operator/getting-started.md`.

## Commands

```text
just init name="…"   # name bootstrap only; no git
just status          # stages and eligible next work
just check           # tree + placeholder/acceptance sanity
```

## Anti-patterns

See `program/reference/anti-patterns.md`. Especially: chat-history authority,
placeholder completion, plan-as-backlog, implementation before authority.

## Cursor Cloud specific instructions

This master repo is **Mycelium** (`framework/blueprint.md`, DEC-010): a
convention-over-configuration thinking framework whose product is a single
static Go binary (`mycelium`). PHASE-01 has not landed yet — there is no
`go.mod`, `cmd/`, or `internal/` on `main`. Until the CLI exists, the
operator surface is `just` + the stdlib Python scripts in `scripts/`.

There are **no long-running services**. Do not start databases, Docker
Compose, Node, or a web server. `gh` is optional and is only required for
the authenticated GitHub half of MS-101 (`gh repo create` + `idea` topic);
hermetic local work must succeed offline.

### Tooling already expected on the image

- `just` (operator recipes in `Justfile`)
- Go **1.26.x** on `PATH` as `go` (install to `/usr/local/go` so it wins
  over the older distro `/usr/bin/go`)
- `python3` (scripts are stdlib-only; no `requirements.txt`)
- `git`, `gh`

Canonical commands live in the root `Justfile` / `README.md`. After
`go.mod` exists, use `go test ./...` and `go build -o mycelium ./cmd/mycelium`
instead of inventing a parallel task runner.

### Gotchas

- Do **not** run `just init` on this master unless a human is naming a
  program. It rewrites `{{PROJECT_NAME}}` / `{{PROGRAM_ID}}` /
  `{{CREATED_DATE}}` across root docs and the manifest. Prove `init` in a
  temp copy.
- Framework evolution is governed by accepted `framework/decisions/DEC-###`
  and `framework/blueprint.md`. The `docs/00-program-blueprint.md` spine
  files are still placeholders and do not unlock PHASE-01.
- `just check` returning OK means tree shape and acceptance consistency,
  not that research or the CLI is done.
- Distro `apt` ships `/usr/bin/go` 1.22.x. The Cloud Agent install
  places Go 1.26 at `/usr/local/go` and symlinks `go`/`gofmt` into
  `/usr/local/bin`, which precedes `/usr/bin` on the default `PATH`.
  If `go version` is 1.22, you are on the distro binary — use
  `/usr/local/bin/go` or repair the symlink. Do not reinstall Go on
  every boot.
- Startup dependency refresh is a no-op until `go.mod` exists; then it is
  `go mod download` only. Do not reinstall `just` or Go on every boot.
