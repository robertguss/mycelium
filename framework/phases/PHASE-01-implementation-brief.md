# PHASE-01 Implementation Brief — Mycelium CLI Foundation

- **Status:** Binding
- **Date:** 2026-08-14
- **Audience:** Engineering (pstack / poteto-mode)
- **Authority:** `framework/blueprint.md` (Accepted 2026-08-14). DEC-001 through DEC-011 stand. This brief records DEC-012 and DEC-013 and commissions PHASE-01 only.
- **Repo:** https://github.com/robertguss/mycelium
- **Product:** single-binary Go CLI `mycelium` that scaffolds per-idea repos. Master builds the CLI. `program/` is `go:embed`'d and emitted on scaffold.
- **Phase gate:** MS-101 (hermetic local + separately credentialed GitHub integration). Five-minute spark-to-first-thought is a user SLO, not the phase gate.
- **How to use this file:** Engineering executes from THIS FILE ONLY. No "see chat". No TBD. Open items are decided here and labeled **Architect default**.

Headings: §§1–20 then Appendices A–F (DEC-012, DEC-013, manifest example, schema example, teaching-error format, file tree).

## 1. Scope / out of scope

Tonight is PHASE-01 only. Later phases get their own briefs after MS-101 is accepted. Do not implement, stub-ship, or "leave a hook command" for later phases.

### In scope

- Go 1.26 module `github.com/robertguss/mycelium`, binary `mycelium`, `CGO_ENABLED=0`.
- `go:embed` of `program/` and emit on scaffold.
- Commands: `version`, `new idea`, `new <type>`, `check`, `tier`, `publish`.
- `mycelium.toml` as the instance manifest (DEC-012).
- Machine-readable tiers. Default scaffold tier: `focused`.
- Data-driven generator and checker (template + sidecar schema; zero Go changes to add a type).
- Operation protocol: lock, journal, stage, atomic rename, resume, abort.
- Teaching errors.
- Fixture-instance CI (hermetic) and a separately credentialed GitHub job.
- Commissioning artifacts: this brief, DEC-012, DEC-013, acceptance matrix, 2.0 contracts, registered templates+schemas, skeleton, mycelium-cli skill.
- Rewrite of registered `program/templates/` to add `+++` TOML front matter. That is program/ authorship, in scope.

### Out of scope (Quality refuses PRs that add them)

- `mycelium status` / `status --all`
- `mycelium supersede`
- council / second-opinion / perspective ladder
- handoff packet generator
- PHASE-02+ skills (`spark`, `wake`, `portfolio`, `thinking-mode`)
- migration machinery (DEC-011)
- destroy command
- GitHub "Use this template"
- hub/garden repo
- converting the master repo into a scaffolded instance
- deleting Justfile/scripts on master before MS-101 is accepted
- committing work product (CLI never git-commits instance files)
- content grading / thinking-quality linters (DEC-005)
- state-transition commands (PHASE-01 sets `state=spark` at birth only)
- `just init` on master
- renaming master's `research-program.toml` to `mycelium.toml`

### Master vs instance (do not confuse)

Master remains an ADRP v1 instance for its own evolution. Do not convert master's `research-program.toml` to `mycelium.toml`. Do not run `just init` on master. `framework/` stays master-only and is NEVER emitted. PHASE-01 replaces the operational surface for *new idea instances* with the CLI. Justfile/scripts retire from the *instance* story once the CLI covers them. They stay on master until MS-101 is accepted.

## 2. Authority and do-not-reopen

### Governing documents

| Document | Status | Rule |
| --- | --- | --- |
| `framework/blueprint.md` | Accepted 2026-08-14 | Do not rewrite vision. Do not reopen DEC-001–011. |
| DEC-001 | Accepted | Evolve ADRP in place into Mycelium. |
| DEC-002 | Accepted | Durable tiered record is the product. Tiers: focused / standard / high-assurance. |
| DEC-003 | Accepted | One repository per idea. No hub. `idea` topic on publish. |
| DEC-004 | Accepted | Self-contained multi-runtime. Instance files are enough to think. |
| DEC-005 | Accepted | Convention over configuration. Checks validate containers, never contents. |
| DEC-006 | Accepted | Lifecycle: spark → exploring ⇄ simmering → clarified → handed-off; any → archived. |
| DEC-007 | Accepted | Sparring / agreement states / cruxes — PHASE-03, not this brief. |
| DEC-008 | Accepted | Perspective ladder — PHASE-04, not this brief. |
| DEC-009 | Accepted | Name is Mycelium. |
| DEC-010 | Accepted | Mycelium is a CLI. Happy path: emit + git init + `gh repo create` + `idea` topic. CLI never commits work product. |
| DEC-011 | Accepted | No migration machinery. Two version fields. Runtime reads instance files, never embed, for check/generate. |
| DEC-012 | Binding this brief | New instances use `mycelium.toml`. Engineering lands the DEC file. Text: Appendix A. |
| DEC-013 | Binding this brief | REFUSE (not warn) when allocating outside a declared stage-scoped range (OQ-007). Text: Appendix B. |

### Arvo-settled open questions (binding)

- **OQ-001:** keep `spark` / `exploring` / `simmering` / `clarified` / `handed-off` / `archived` exactly as written in DEC-006 / blueprint. Do not bikeshed synonyms.
- **OQ-002:** new instances use `mycelium.toml` (2.0 name). Recorded as DEC-012.
- **OQ-007:** default to REFUSE (not warn) when allocating outside a declared stage-scoped range. Recorded as DEC-013.

### Process override for THIS project

Blueprint "humans-own-git" is overridden for the *master* repo's engineering workflow: Arvo merges Quality-green PRs and accepts the phase. Engineering opens PRs. Engineering does NOT push to main. The CLI still never git-commits *instance* work product. Those two rules are not in conflict.

### Environment

Cursor cloud env: Go 1.26 at `/usr/local/bin/go` (PR #3 merged). Do not block on toolchain. Do not add a version-manager dance.

### Do not reopen

Do not reopen the product shape, the language, the dependency floor, the state vocabulary, the manifest filename, the refuse-vs-warn range rule, the no-commit rule, the instance-files-are-truth rule, or the hermetic/credentialed split of MS-101. If a later discovery contradicts a locked item, stop and write a DEC; do not silently patch this brief in a code PR.

## 3. Current main (what exists / what must not be broken)

Repo on `main` is still ADRP v1. SHA at brief time is not a pin; treat `main` as the live v1 tree.

### What exists

| Path | Role | PHASE-01 fate |
| --- | --- | --- |
| `Justfile` | v1 operator (`just check`, `just init`, …) | Keep on master until MS-101 accepted. Do not delete. |
| `scripts/*.py` | v1 check/init/status | Keep on master. |
| `research-program.toml` | master's v1 manifest | Keep. Do not rename. |
| `program/contracts/*` | v1 contracts | Keep. ADD 2.0 files beside them. |
| `program/operator/*` | v1 bootstrap | Keep on master. |
| `program/reference/*` | v1 reference, including `rigor-tiers.md` | Keep. |
| `program/templates/*.md` | v1 templates, no front matter | Keep unregistered v1 files. Rewrite *registered* types to 2.0 names + `+++` front matter (Slice 2). |
| `.agents/skills/research-*` | v1 skills | Stay on master. Do NOT emit into new instances. |
| `framework/` | blueprint + DEC-001–011 | Master-only. NEVER emitted. |
| `docs/` | v1 placeholders | Master-only. Not emitted. |
| `decisions/` | v1 empty placeholder | Master-only. Not emitted. |
| `AGENTS.md`, `README.md` | v1 operator docs | Keep on master. Instance gets its own emitted copies. |
| `cmd/`, `internal/`, `go.mod` | absent on main today | Added by Slice 1. |

v1 registered-looking templates that are NOT PHASE-01 generator types (keep as browsable `program/templates/` + old contracts; do not register schemas this phase): `attachment-manifest`, `bootstrap-task`, `launch-message`, `stage-package`, `validation-report`, `validation-task`, `focused-research-prompt`. Phase and milestone ARE registered (see §6).

v1 `program/contracts/identifiers.md` lists DEC/REC/REQ/FND/RSK/OQ/SPK/PHASE/MS and marks EVD/ASM optional. Slice 0/2 updates that contract: ASM and EVD become first-class; keep v1 ranges as examples.

### What must not be broken

- `just check` on master must keep passing. It does not look at `go.mod`. Adding Go files is fine.
- Master's `research-program.toml` stays the v1 manifest.
- Do not run `just init` on master.
- Do not emit `framework/` into an instance.
- Do not delete Justfile/scripts before MS-101 is accepted.
- Do not convert master into a scaffolded 2.0 instance.

### Rollback posture

If a CLI PR is bad: revert that PR. Master remains ADRP v1. That is the rollback.

## 4. Target architecture

### Language and module

| Item | Value |
| --- | --- |
| Language | Go 1.26 |
| Module path | `github.com/robertguss/mycelium` |
| Binary name | `mycelium` |
| Entry | `cmd/mycelium/main.go` |
| Link mode | `CGO_ENABLED=0` static binary |
| CLI version (source default) | `0.1.0-dev` until PHASE-05 |
| Version stamp | `ldflags -X github.com/robertguss/mycelium/internal/version.Version=...` |
| Emitted methodology version | `2.0.0` |

### Dependencies

stdlib-first. Allowed third-party: ONLY `github.com/pelletier/go-toml/v2` (manifest + schema.toml + front matter). Pin the latest tagged v2 at add time via `go get`. No cobra. No viper. No yaml library. No testify. No jinja. No github.com/google/go-github. Tests use stdlib `testing` + table-driven cases.

Subcommands: stdlib `flag` + a small `internal/cli` dispatcher.

### Packages (binding layout)

```text
cmd/mycelium/main.go          # os.Exit(cli.Main(os.Args))
internal/cli/                 # dispatch, global flags, usage
internal/version/             # Version string; default 0.1.0-dev
internal/embed/               # go:embed all:program
internal/embed/program/       # generate-copied embed root (see below)
internal/clock/               # Clock interface; real UTC; test inject; MYCELIUM_NOW
internal/execrun/             # Runner interface for git/gh (name avoids exec keyword clash)
internal/metadata/            # +++ TOML front-matter reader
internal/idpath/              # pure ID ↔ path
internal/manifest/            # mycelium.toml parse/validate
internal/schema/              # sidecar schema.toml load
internal/slug/                # slugify
internal/logfmt/              # parse/append log lines
internal/teach/               # teaching-error writer (stderr)
internal/lock/                # .mycelium/lock flock
internal/journal/             # .mycelium/journal.json
internal/op/                  # protocol wrapper (preflight/lock/stage/commit/rollback)
internal/scaffold/            # new idea
internal/generate/            # new <type>
internal/check/               # conformance
internal/tiercmd/             # mycelium tier
internal/publish/             # mycelium publish / new idea publish half
internal/clitest/             # hermetic binary helper (test-only)
```

Do not invent extra public packages. Unexported helpers may live next to their tests. No `pkg/`.

### Embed (go:embed cannot escape the package dir)

**Architect default:** package `internal/embed` contains `embed.go` with `//go:embed all:program` and an `internal/embed/program/` directory. Repo-root `program/` is the authoritative browsable tree. `go generate` in `internal/embed` copies repo-root `program/` → `internal/embed/program/` using stdlib only (`io`, `os`, `path/filepath`). Filter `*.go` out of the copy so a generate helper never ships into instances. Slice 1 may embed a stub (`internal/embed/program/.keep`) so `go test` / `go build` work before Slice 2 authors 2.0 content. After Slice 2, generate is run before commit.

Runtime scaffold emits from the embed FS. Runtime `check` and `new <type>` read the *instance* `program/` files, never the embed FS (DEC-011).

### Clock

```text
type Clock interface { Now() time.Time }
```

Real implementation returns `time.Now().UTC()`. Tests inject a fixed clock. Env `MYCELIUM_NOW` (RFC3339) overrides the real clock when set. Dates in artifacts and the log are `YYYY-MM-DD` UTC.

### Exec adapter

```text
type Runner interface {
    LookPath(name string) (string, error)
    Run(ctx context.Context, name string, args ...string, opts RunOpts) (Result, error)
}
```

`RunOpts` carries `Dir` (cwd) and env extras. Real runner wraps `os/exec`. Tests install a fake that records `(name, args, dir)` and returns programmed results. Hermetic tests assert `gh` was never invoked under `--offline` or `MYCELIUM_OFFLINE=1`. `git` and `gh` are invoked as binaries, not as API clients.

### Version command output

**Architect default:** `mycelium version` prints exactly one line to stdout: the stamped version string (default `0.1.0-dev`) and a trailing newline. No prefix. Tests compare equality.

### Dispatcher rules

- `mycelium` / `mycelium -h` / `mycelium --help` → usage on stdout, exit 0.
- Unknown command → teaching error, exit 1.
- Missing required args → teaching error, exit 1.
- Exit 0 on success, exit 1 on every failure (usage, conformance, refuse). No other exit codes this phase.

`new` is one command: if the first argument is `idea`, run scaffold; otherwise treat it as a type key and run the generator.

### Instance root resolution

**Architect default:** `--dir PATH` is the instance root (not a parent). If `--dir` is absent:

- `new idea` uses `./<slug>` relative to cwd.
- `new <type>`, `check`, `tier`, `publish` walk from cwd upward looking for `mycelium.toml`, stopping at a `.git` directory or filesystem root. If none, refuse with a teaching error.

`--dir` parent must exist. The command creates only the final directory (`new idea`) or refuses (`new idea` if the target exists).

## 5. Manifest (`mycelium.toml`)

New instances use this filename (DEC-012). Master keeps `research-program.toml`.

### Field list (lock)

```toml
schema_version = 1
idea_name = "..."
slug = "..."
state = "spark"
tier = "focused"
methodology_version = "2.0.0"
generated_by_cli_version = "0.1.0-dev"
created_date = "YYYY-MM-DD"
updated_date = "YYYY-MM-DD"
revisit = ""          # required non-empty when state=simmering
github_repo = ""      # "owner/name" after publish; empty if local-only

[identifiers]
# optional stage-scoped ranges, e.g.
# findings = "FND-001..FND-099"
# recommendations = "REC-001..REC-099"
# requirements = "REQ-001..REQ-299"

[[deviations]]
# optional; omit the array when empty
# convention = "extra-top-level:notes.md"
# reason = "scratch pad; not an artifact"
```

### Field rules

| Field | Type | Required | Notes |
| --- | --- | --- | --- |
| `schema_version` | int | yes | Must be `1` this phase. Other values: teaching error. |
| `idea_name` | string | yes | Non-empty. Human name as passed to `new idea`. |
| `slug` | string | yes | kebab-case. Must equal `slugify(idea_name)`. |
| `state` | string | yes | One of the six DEC-006 names. Birth value: `spark`. |
| `tier` | string | yes | `focused` \| `standard` \| `high-assurance`. Birth default: `focused`. |
| `methodology_version` | string | yes | Emitted `2.0.0`. Orthogonal to CLI version (DEC-011). |
| `generated_by_cli_version` | string | yes | Stamped from `internal/version`. |
| `created_date` | date | yes | `YYYY-MM-DD` UTC at scaffold. Never rewritten. |
| `updated_date` | date | yes | `YYYY-MM-DD` UTC. Bumped on every successful mutating op. |
| `revisit` | string | yes (key) | Empty string allowed except when `state=simmering`. |
| `github_repo` | string | yes (key) | Empty until publish. Then `owner/name` (no `https://`). |
| `[identifiers]` | table | optional | Keys: `findings`, `recommendations`, `requirements`. Values: `NS-XXX..NS-YYY`. |
| `[[deviations]]` | array of tables | optional | Each row: `convention` (non-empty) + `reason` (non-empty). Silent deviation is a check failure. |

State and tier are orthogonal. Two version fields are orthogonal. Do not infer tier from state. Do not infer methodology version from CLI version.

Range value grammar: `^([A-Z]+)-([0-9]+)\.\.([A-Z]+)-([0-9]+)$`. Both namespaces must match the key's NS. Start integer <= end integer. One range per key this phase.

**Architect default:** there is no `mycelium range` command this phase. Fixture CI / tests write `[identifiers]` into the fixture manifest with the stdlib before generating FND/REC/REQ.

Unknown top-level keys in the manifest: **Architect default:** refuse on parse (forward-incompatible is a DEC, not a silent ignore). Unknown keys under `[identifiers]`: refuse. Extra keys on a deviation row: refuse.

## 6. Artifact catalog + ID-to-path + schemas

### Registered types (PHASE-01 only)

These are the ONLY generator types this phase. Adding a type = template + schema, zero Go changes to generator/checker.

| type key | NS | home dir | filename | digits | stage-scoped? |
| --- | --- | --- | --- | --- | --- |
| `decision` | DEC | `decisions` | `DEC-###-<slug>.md` | 3 | no |
| `assumption` | ASM | `assumptions` | `ASM-###-<slug>.md` | 3 | no |
| `evidence` | EVD | `evidence` | `EVD-###-<slug>.md` | 3 | no |
| `spike` | SPK | `spikes` | `SPK-###-<slug>.md` | 3 | no |
| `finding` | FND | `findings` | `FND-###-<slug>.md` | 3 | YES |
| `recommendation` | REC | `recommendations` | `REC-###-<slug>.md` | 3 | YES |
| `requirement` | REQ | `requirements` | `REQ-###-<slug>.md` | 3 | YES |
| `question` | OQ | `questions` | `OQ-###-<slug>.md` | 3 | no |
| `risk` | RSK | `risks` | `RSK-###-<slug>.md` | 3 | no |
| `phase` | PHASE | `phases` | `PHASE-##-<slug>.md` | 2 | no |
| `milestone` | MS | `milestones` | `MS-###-<slug>.md` | 3 | no |

PHASE uses two digits (`PHASE-01`). All others three digits, zero-padded. Fixture CI generates one of each.

### ID-to-path (pure function, both directions)

Package `internal/idpath`. No filesystem. No clock.

```text
ID        = NS "-" DIGITS          e.g. DEC-001, PHASE-01
Path      = home "/" NS "-" DIGITS "-" slug ".md"
ParseID   = ^([A-Z]+)-([0-9]+)$
ParsePath = ^home/NS-DIGITS-slug\.md$
```

Rules:

- `Format(ns, n, slug)` zero-pads `n` to the namespace's digit count.
- `Parse("DEC-001")` → `{NS:DEC, N:1}`. `Parse("DEC-1")` is a valid *id token* for link resolution (numeric), but a file named `DEC-1-foo.md` fails the filename pattern.
- `PathFor(id, slug)` and `ParsePath(path)` are inverses for every registered type and every slug that `slugify` accepts.
- Conformance enforces both ways: every artifact file matches its pattern and lives in its home; every in-repo ID reference resolves to a matching file.

Link resolution (**Architect default**): scan artifact bodies (and `log.md`, `README.md`, `CONTEXT.md`, `AGENTS.md`) for `\b(DEC|ASM|EVD|SPK|FND|REC|REQ|OQ|RSK|PHASE|MS)-[0-9]+\b`. Require a matching file in that NS home (`NS-<zero-padded>-*.md`). Do not crawl the web. Do not resolve URLs. A reference to `DEC-001` matches `decisions/DEC-001-*.md`. Missing file → teaching error naming the convention and `program/contracts/naming.md`.

### Allocation

Filesystem is the registry. Scan the type's home dir for files matching the filename pattern. Next ID = `max(N)+1`, or `1` if the dir is missing/empty. Never reuse. Refuse overwrite.

**Architect default: ALLOW gaps.** Next = max(N)+1 over files present. Sequence check does not require dense `001..N`. IDs must be unique. Files must match the pattern. No tombstones this phase. A deleted max ID may be reissued (file is gone). A still-present file is never overwritten. Humans still must not silently remove history; the generator never deletes.

Stability rule from v1 identifiers contract still applies to *humans*: do not silently remove history. The generator will not delete. Lowering a tier will not delete.

### Stage-scoped allocation (FND / REC / REQ) — DEC-013 / OQ-007

`[identifiers]` must declare a range for that type BEFORE allocate. If no range, or the next ID would fall outside all declared ranges → **REFUSE** (not warn). Teaching error names DEC-013, `program/contracts/identifiers.md`, and the missing/overflowed range.

Fixture CI must write ranges into the fixture manifest before generating those three types (or generate them after a helper that declares ranges).

### Front matter

Blueprint: metadata is front matter; one metadata reader; body cannot masquerade as metadata.

**Architect default:** TOML front matter fenced by `+++` (not YAML `---`). The one TOML parser (`pelletier/go-toml/v2`) serves manifest + schema + front matter. No yaml dependency.

```text
+++
id = "DEC-001"
title = "..."
status = "Proposed"
date = "2026-08-14"
owner = "Robert Guss"
+++

# DEC-001 — Title

## Context
...
```

Reader rules:

- File must start with `+++` at byte 0 (optional single leading UTF-8 BOM is stripped; nothing else).
- Closing `+++` is the first subsequent line that is exactly `+++`.
- Bytes between the fences are TOML. Required keys come from the sidecar schema.
- A `+++` line in the body after the closing fence is body text. It does not start a second front-matter block.
- Missing opening fence, missing close, or unparseable TOML → teaching error.

Existing v1 `decision-record.md` has no front matter. PHASE-01 rewrites registered templates to add `+++` front matter. Do not retrofit master's `framework/` DECs.

### Tokens

Templates use `{{ID}}` `{{TITLE}}` `{{SLUG}}` `{{DATE}}` only. Plain `strings.Replace`. No text/template. No jinja. `DATE` = `YYYY-MM-DD` UTC.

### Sidecar schema shape (every registered type)

File: `program/templates/<type>.schema.toml` sitting next to `program/templates/<type>.md`.

```toml
namespace = "DEC"
home = "decisions"
filename_pattern = "DEC-{NNN}-{slug}.md"
stage_scoped = false
digits = 3
required_front_matter = ["id", "title", "status", "date", "owner"]
required_sections = [
  "Context",
  "Decision",
  "Rationale",
  "Consequences",
  "Alternatives Considered",
  "Risks",
  "Revisit Triggers",
  "Approval",
]

[enums.status]
values = ["Proposed", "Accepted", "Superseded", "Rejected"]
```

`filename_pattern` tokens: `{NNN}` = 3-digit, `{NN}` = 2-digit, `{slug}` = kebab. `{NS}` is not required because NS is a field.

Required sections are H2 headings (`## Name`). Extra H2s are allowed. Missing required H2 → teaching error. Heading match is exact, case-sensitive.

If a required front-matter field has an `[enums.<field>]` table, the value must be one of `values`. If no enum table, the field must be present and a non-empty string (dates match `^\d{4}-\d{2}-\d{2}$`).

`id` in front matter must equal the filename's ID (`NS-DIGITS`). `title` should equal the generator title; check does not rewrite it. **Architect default:** check requires `id` match only; title drift is not a failure (containers, not contents).

### Required front matter and H2s per type

| type | required_front_matter | enums | required_sections (H2) |
| --- | --- | --- | --- |
| decision | id, title, status, date, owner | status: Proposed, Accepted, Superseded, Rejected | Context; Decision; Rationale; Consequences; Alternatives Considered; Risks; Revisit Triggers; Approval |
| assumption | id, title, status, date, attached_to | status: Open, Held, Falsified, Retired | Statement; Falsifier; Implications; Revisit Triggers |
| evidence | id, title, status, date, source | status: Draft, Recorded, Revalidated, Superseded | Claim; Source; Observation; Limitations; Revalidation Trigger |
| spike | id, title, status, date, decision_at_stake | status: Planned, Completed, Inconclusive, Superseded | Method; Commands and Artifacts; Results; Limitations; Architectural Consequence; Cleanup; Reproduction Instructions |
| finding | id, title, severity, confidence, date | severity: Critical, High, Medium, Low; confidence: High, Medium, Low | Problem; Evidence; Failure Scenario; Impact; Root Cause; Required Correction; Residual Risk |
| recommendation | id, title, classification, confidence, date | classification: Default, Required, Optional, Exception, Experimental, Watchlist, Rejected; confidence: High, Medium, Low | Recommendation; Requirements and Constraints; Rationale; Evidence; Tradeoffs; Alternatives Considered; Revisit Triggers |
| requirement | id, title, priority, date, phase | priority: Must, Should, May | Requirement; Rationale; Acceptance Evidence; Exceptions |
| question | id, title, agreement, date | agreement: open, aligned, agree-to-disagree | Question; Context; Positions; Crux; Disposition |
| risk | id, title, severity, likelihood, date | severity: Critical, High, Medium, Low; likelihood: High, Medium, Low | Description; Impact; Mitigation; Residual Risk; Revisit Triggers |
| phase | id, title, status, date | status: Planned, In Progress, Accepted, Abandoned | Entry Criteria; Scope; Explicit Non-Goals; Exit Criteria |
| milestone | id, title, phase, date | (none) | Outcome; Prerequisites; Acceptance Evidence |

Templates emit those fences, fields (token-filled), H1 title line, and the H2s with a one-line placeholder comment under each (`<!-- fill -->`). Placeholder comments are body, not metadata.

### v1 template rename (Slice 2)

| v1 file | 2.0 registered file |
| --- | --- |
| `program/templates/decision-record.md` | `program/templates/decision.md` + `decision.schema.toml` |
| `program/templates/evidence-spike.md` | `program/templates/spike.md` + `spike.schema.toml` |
| `program/templates/finding.md` | `program/templates/finding.md` + `finding.schema.toml` |
| `program/templates/recommendation.md` | `program/templates/recommendation.md` + `recommendation.schema.toml` |
| `program/templates/requirement.md` | `program/templates/requirement.md` + `requirement.schema.toml` |
| `program/templates/phase.md` | `program/templates/phase.md` + `phase.schema.toml` |
| `program/templates/milestone.md` | `program/templates/milestone.md` + `milestone.schema.toml` |

Keep the v1 filenames as stubs that point at the new names? **Architect default:** replace in place for the seven above (rewrite content). Add new files for `assumption`, `evidence`, `question`, `risk`. Do not delete unregistered v1 templates.

## 7. Tiers

Names from `program/reference/rigor-tiers.md` and DEC-002: `focused` | `standard` | `high-assurance`.

Default scaffold tier: `focused` (thinnest; spark-capable). Any state may hold any tier. Structure is emitted per tier, never per state.

`mycelium tier` is idempotent: update manifest; emit only structure the new tier newly requires; never overwrite existing work. Lowering a tier deletes nothing (no-deletion) and only relaxes which artifacts bind.

### Machine-readable files

`program/tiers/focused.toml`, `program/tiers/standard.toml`, `program/tiers/high-assurance.toml`. The checker and `tier` command read the *instance* copies.

```toml
# focused.toml
name = "focused"
emits = []
binds = ["manifest", "log.md", "CONTEXT.md", "AGENTS.md", "program/"]
```

```toml
# standard.toml
name = "standard"
emits = ["decisions/", "assumptions/", "evidence/", "questions/", "risks/"]
binds = [
  "manifest", "log.md", "CONTEXT.md", "AGENTS.md", "program/",
  "decisions/", "assumptions/", "evidence/", "questions/", "risks/",
]
```

```toml
# high-assurance.toml
name = "high-assurance"
emits = [
  "decisions/", "assumptions/", "evidence/", "questions/", "risks/",
  "spikes/", "findings/", "recommendations/", "requirements/",
  "phases/", "milestones/",
]
binds = [
  "manifest", "log.md", "CONTEXT.md", "AGENTS.md", "program/",
  "decisions/", "assumptions/", "evidence/", "questions/", "risks/",
  "spikes/", "findings/", "recommendations/", "requirements/",
  "phases/", "milestones/",
]
```

`emits` is relative to instance root. `focused.emits` is empty because the spark skeleton is not a tier emit; it is always emitted by `new idea`.

### Bind vs emit

| tier | emits (beyond spark skeleton) | binds |
| --- | --- | --- |
| focused | none | manifest, `log.md`, `CONTEXT.md`, `AGENTS.md`, `program/` (schemas present) |
| standard | `decisions/`, `assumptions/`, `evidence/`, `questions/`, `risks/` | focused binds + those dirs exist (may be empty) |
| high-assurance | standard + `spikes/`, `findings/`, `recommendations/`, `requirements/`, `phases/`, `milestones/` | standard binds + those dirs exist (may be empty) |

If an artifact file exists at ANY tier, it must pass its schema. Extra valid artifacts at focused are allowed. Missing type dirs at focused are OK.

### Directory creation on promote

**Architect default:** when a dir is newly required, create it and write a one-line `README.md` if and only if that path does not exist:

```text
# <Heading>

Home for <NS>-### artifacts.
```

Do not overwrite an existing README or any other file. Empty dirs without README are acceptable to check (dir exists). The README is so git will keep the dir.

### Lowering a tier

Deletes nothing. `binds` relax. A standard instance that already has `decisions/DEC-001-*.md` and is lowered to focused still checks that DEC against its schema. The `decisions/` dir is not a focused *bind*, but the file exists, so the schema applies.

### `mycelium tier` behavior

```text
mycelium tier <tier> [--dir PATH]
```

- Unknown tier → teaching error listing the three names.
- Same tier already set and all emit dirs present → no-op, exit 0, stdout: `already <tier>`.
- Same tier but missing emit dirs → create the missing ones, bump `updated_date`, append a log line, exit 0.
- Raise or lower → set `tier`, emit newly required dirs only, bump `updated_date`, append one log line, exit 0.
- Runs under the operation protocol.

## 8. Commands

Exact CLI. Flags. Exit codes. Teaching errors.

Global: exit 0 success, exit 1 failure. Teaching errors on stderr. Success text on stdout.

Env:

| Env | Effect |
| --- | --- |
| `MYCELIUM_OFFLINE=1` | Treat as `--offline` on every command that can publish. Never exec `gh`. Never open network. |
| `MYCELIUM_NOW` | RFC3339 clock override. |

### 8.1 `mycelium version`

```text
mycelium version
```

No flags. Prints the stamped version line. Ignores `--dir`. Ignores offline.

### 8.2 `mycelium new idea <name>`

```text
mycelium new idea <name> [--dir PATH] [--offline] [--publish] [--tier focused|standard|high-assurance]
```

| Flag | Default | Rule |
| --- | --- | --- |
| `--dir PATH` | `./<slug>` relative to cwd | Exact destination directory. Does not append slug. |
| `--offline` | off (unless `MYCELIUM_OFFLINE=1`) | Never exec `gh`. Never open network. Always legal. |
| `--publish` | off | Require authenticated `gh`. Fail with a teaching error if not. Implies publish after local success. |
| `--tier` | `focused` | Must be one of the three names. |

`<name>` may be multiple words (use quotes). Slug = kebab-case of name. Refuse empty slug. Refuse if target exists (file or directory, even empty). Do not mkdir `~/ideas`. `~/ideas/<slug>` is the documented human convention, not a hardcoded default.

Always: emit + `git init`. NEVER git-commit. NEVER `git add`.

Publish half — **Architect default** (DEC-010 happy path vs MS-101 hermetic gate):

1. `new idea` ALWAYS emits + `git init`.
2. Default (no `--offline`, no `--publish`): if `gh` is on PATH and `gh auth status` exits 0, also publish. If `gh` is missing or unauthenticated, SUCCEED locally and print the next step: `mycelium publish`. Do not fail the scaffold.
3. `--offline` or `MYCELIUM_OFFLINE=1`: never exec `gh`, never open network. Hermetic tests MUST pass `--offline` or the env.
4. `--publish`: require authenticated `gh`; fail with a teaching error if not. Do not leave a half-published remote without cleanup (see §10).
5. `--offline` and `--publish` together → teaching error (contradictory).

`git` is the `git` binary. `gh` is the `gh` binary.

### 8.3 `mycelium new <type> "<Title>"`

```text
mycelium new <type> "<Title>" [--dir PATH]
```

Run from an instance root (or `--dir`). `<type>` is a registered type key from the instance's `program/templates/*.schema.toml`. Unknown type → teaching error listing registered keys.

Allocates next ID, refuses overwrite, refuses out-of-range (DEC-013), fills tokens, writes the file, appends one log line, bumps `updated_date`. Never runs git. Runs under the operation protocol.

**Architect default:** `--dir` is supported so tests need not chdir. Also walks parents for `mycelium.toml` when `--dir` is absent.

### 8.4 `mycelium check`

```text
mycelium check [--dir PATH] [--abort-journal]
```

Read-only unless `--abort-journal`.

`--abort-journal` (**Architect default**, documented rollback): delete staged temps listed in `.mycelium/journal.json`, delete the journal, delete a stale lock file. Do **not** delete already-renamed artifacts (no-deletion). Print the surviving paths. Exit 0 if the journal/lock are gone; exit 1 if there was nothing to abort (teaching error: no journal).

Success (no abort): print a short OK summary and exit 0:

```text
mycelium check: ok
instance: <slug>
state: spark
tier: focused
artifacts: <n>
```

Failure: one or more teaching errors on stderr, exit 1. Print all independent failures found (do not stop at the first), then exit 1. **Architect default:** cap at 20 errors, then `mycelium: further errors omitted`.

### 8.5 `mycelium tier <tier>`

```text
mycelium tier <tier> [--dir PATH]
```

See §7. Operation protocol.

### 8.6 `mycelium publish`

```text
mycelium publish [--dir PATH]
```

Publishes an already-local instance. Idempotent if the remote already exists and the `idea` topic is present. Requires authenticated `gh`. Never commits. Records `github_repo`. Appends a log line. Operation protocol.

### Commands that do not exist this phase

`status`, `status --all`, `supersede`, `destroy`, `upgrade`, `init` (as a mycelium command), `range`, `wake`, `council`. Quality refuses PRs that add them.

## 9. Operation protocol

Applies to: `new idea` (local emit half), `new <type>`, `tier`, and the generate/publish half of `publish` / `new idea` publish. `version` and `check` (without abort) do not use it.

Supported FS floor: local filesystem with atomic rename. Network filesystems are outside the floor. Lock + journal still bound the damage there; do not claim atomicity.

### Steps

1. **Preflight** — manifest+log parse (except first-time scaffold, where they do not exist yet), schema resolves, target path free, flags legal, instance root found.
2. **Lock** — exclusive repo lock at `<root>/.mycelium/lock` via `fcntl`/`flock` on a file descriptor. Write `pid=<n>\nstarted=<RFC3339>\n` into the file. Create `.mycelium/` if needed.
3. **Stage** — write outputs as temporary files under `.mycelium/stage/<op-id>/` and record intent in `.mycelium/journal.json`.
4. **Commit** — atomic rename in a fixed order: artifact file(s), then log, then manifest. For scaffold, rename skeleton files, then `program/`, then log, then manifest last.
5. **Rollback** — failure before the first rename: delete staged files, delete journal, change nothing committed. After a partial commit: journal survives; re-run of the *same* command resumes under the original ID rather than allocating a new one.
6. **Detection** — `check` detects leftover journal or stale lock and names recovery.

Recovery text (teaching error):

```text
mycelium: interrupted operation
convention: operation-protocol
contract: program/contracts/operation-protocol.md
fix: re-run the same command to complete, or mycelium check --abort-journal to roll back
```

### Journal schema

`.mycelium/journal.json` (stdlib `encoding/json`):

```text
{
  "schema_version": 1,
  "op": "scaffold" | "new" | "tier" | "publish",
  "type": "decision" | null,
  "title": "..." | "",
  "original_id": "DEC-001" | "",
  "started_at": "<RFC3339>",
  "staged_dir": ".mycelium/stage/<op-id>",
  "renames": [
    {"from": "...", "to": "...", "done": false}
  ],
  "log_line": "....",
  "argv": ["new", "decision", "Title"]
}
```

Resume: leftover journal + incoming argv matches `argv` (or same `op`+`type`+`title`) → reuse `original_id` and finish undone renames. Leftover journal + different command → refuse, name `--abort-journal`.

After all renames succeed: delete journal, delete staged dir, release lock, remove lock file.

### Lock / stale lock

flock is released when the process exits. A leftover lock *file* with no live PID is stale.

**Architect default:** stale = lock file exists AND (`pid` missing OR `os.FindProcess` + signal 0 returns ESRCH) AND (optional) `started` older than 30 minutes *or* PID is dead regardless of age. Dead PID → stale immediately. Live PID → lock held.

`check` fails on leftover journal OR stale lock. `check` does **not** fail solely because a live lock is held by another process; it prints a one-line notice on stdout and continues the rest of the suite. **Architect default.**

`--abort-journal` also removes a stale lock file.

### Scaffold vs git init

`git init` happens AFTER successful file commit, BEFORE publish. If `git init` fails, instance files exist; print a teaching error; do not roll back files. User can `git init` by hand. **Architect default:** `git init -b main`.

### Publish and the protocol

Publish updates manifest (`github_repo`, `updated_date`) and log. Those two writes use the protocol. `gh` exec happens after preflight and lock, before manifest rename: if `gh` fails, no manifest/log change. If `gh` succeeded and manifest rename fails, journal survives; resume records the same `github_repo` and does not create a second repo (idempotent `gh`).

## 10. Scaffold (`new idea`)

### Emit list (spark skeleton)

Always emitted. `framework/` is NEVER emitted.

| Path | Notes |
| --- | --- |
| `README.md` | Idea name, one-line purpose stub, `mycelium check` pointer. |
| `mycelium.toml` | §5. `state=spark`, `tier` from flag (default focused), versions stamped, dates from clock, `revisit=""`, `github_repo=""`. |
| `log.md` | Heading + one scaffold line. Nearly empty is OK. |
| `CONTEXT.md` | `# Glossary` and a blank line. Empty glossary is OK at spark. |
| `AGENTS.md` | Instance-facing. Teaches CLI + manual floor. No master-only notes. No Justfile. No `research-program.toml`. |
| `.agents/skills/mycelium-cli/SKILL.md` | PHASE-01 thin skill: commands, flags, manual floor, teaching-error shape, "do not git commit unless the human asks". |
| `.gitignore` | Contains `.mycelium/lock` only. Journal is NOT ignored (interrupted ops stay visible). |
| `program/` | Embedded program tree (contracts, templates+schemas, tiers, skills source, skeleton, README, reference). |

If `--tier standard` or `high-assurance`, also emit that tier's dirs (+ README stubs) at birth.

### Do not emit

`framework/`, `cmd/`, `internal/`, `go.mod`, `go.sum`, `Justfile`, `scripts/`, `research-program.toml`, master's `.git`, master's `decisions/`, master's `docs/` v1 placeholders, v1 `.agents/skills/research-*`, `*.go` from embed generate.

### Git

`git init -b main` in the instance dir. No `git add`. No `git commit`. No `git config` writes.

### Publish split (critical for MS-101)

DEC-010 user-facing happy path is emit + git init + `gh repo create` + `idea` topic. MS-101 hermetic gate requires no network. The flag/env split in §8.2 is the resolution.

Publish steps (when allowed):

1. `gh auth status` (fail only if `--publish`; otherwise skip publish).
2. `gh repo create <slug> --private` under the authenticated user. Do **not** pass `--push` (there is no commit). Do **not** pass `--source` if that implies a push.
3. `git remote add origin <url>` if origin is missing. If origin already points at that URL, continue.
4. `gh repo edit <owner/name> --add-topic idea`.
5. Set `github_repo = "owner/name"` via the protocol.
6. Append log line: `publish` / `-` / `github.com/owner/name`.

**Architect default:** repos are private. Description: `idea: <idea_name>` truncated to 80 chars.

Idempotent `mycelium publish`:

- Remote exists + topic present + `github_repo` set → exit 0, stdout `already published: owner/name`.
- Remote exists, topic missing → add topic, update manifest if needed.
- No remote → create + remote + topic + manifest.

`--publish` failure after `gh repo create` succeeded: attempt `gh repo delete <name> --yes` only when the name matches `^mycelium-ms101-[0-9]+$` (fixture tests). For user slugs, do **not** auto-delete. Print the repo URL and the teaching error. Journal records the created name so resume does not create a second repo.

### Slugify

**Architect default (DEC-014):** PHASE-01 slugify is the documented latin/compatibility fold in `internal/slug` (`latinFold` + ASCII `[a-zA-Z0-9]` + map space/`_` to `-`, drop other runes including unlisted letters, collapse `--`, trim `-`, lowercase). Max 80 characters. Empty → refuse. Refuse `.` and `..`. Full Unicode NFKD is deferred; do not grow `latinFold` this phase. See `framework/decisions/DEC-014-phase-01-slugify-latin-fold.md`.

### Default directory

`./<slug>` relative to cwd. `--dir` overrides (exact path). Do not mkdir `~/ideas`.

### Stdout on success (local only)

```text
created ./<slug>
state: spark
tier: focused
next: cd <slug> && mycelium new decision "First thought"
publish: mycelium publish
```

If published in the same invocation, replace the last line with `published: owner/name (topic: idea)`.

## 11. Generator (`new <type>`)

Data-driven. Reads instance `program/templates/<type>.md` and `<type>.schema.toml`. Does not switch on type in Go beyond "load schema by key". The type key is the filename stem.

### Algorithm

1. Resolve instance root. Preflight: manifest+log parse, schema exists, title non-empty.
2. Slugify title. Refuse empty slug.
3. If `stage_scoped` and no `[identifiers]` range for that home key → REFUSE (DEC-013).
4. Lock.
5. Scan home dir. `next = max(N)+1` or `1`. Format with `digits`.
6. If stage-scoped and next ID outside the declared range → REFUSE (DEC-013). Unlock. No writes.
7. Destination `home/NS-DIGITS-slug.md`. If it exists → REFUSE overwrite (teaching error: rename or pick another title).
8. Stage: replace `{{ID}}` `{{TITLE}}` `{{SLUG}}` `{{DATE}}` in the template. Write temp artifact. Write temp log (original + one new line). Write temp manifest (`updated_date` only, unless a future field needs it; identifiers are not auto-written).
9. Journal with `original_id`.
10. Commit renames: artifact, log, manifest.
11. Print path + next-step one-liner. Exit 0.

Home dir is created if missing (even at focused, when the user generates an extra valid artifact). That is allowed.

### Log line

```text
YYYY-MM-DD\tnew\tDEC-001\tTitle here
```

One line per artifact.

### Refuse cases (all teaching errors, exit 1)

- unknown type
- not in an instance
- empty title / empty slug
- overwrite
- no range (stage-scoped)
- out of range (stage-scoped)
- leftover journal for a *different* op
- lock held by a live other process (wait? **Architect default:** refuse immediately, do not block forever; print pid)

**Architect default:** lock acquire timeout is 0 (non-blocking flock). If busy, refuse.

## 12. Check

Structure only (DEC-005). Never grades prose. Never requires N decisions. Never requires glossary terms. Never requires cruxes (PHASE-03). Never requires wake briefs (PHASE-02).

Runtime reads instance files, never embed.

### Must implement

1. **ID uniqueness + ID-to-path both directions.** Every file in a registered home matches the pattern. Every parsed ID maps back to that path. Duplicate IDs fail.
2. **Link resolution** for in-repo ID references (DEC-001 style). Regex in §6. Require a matching file. Do not crawl the web.
3. **Required front matter + required H2 sections** per that type's sidecar schema. Body `+++` does not count as a fence.
4. **Legal state values** + legal stored-state invariants:
   - state ∈ {spark, exploring, simmering, clarified, handed-off, archived}
   - `simmering` ⇒ `revisit` non-empty
   - `clarified` or `handed-off` ⇒ FAIL this phase with a teaching error naming the missing PHASE-02 state-transition command and the PHASE-06 packet (`handed-off` also requires a packet; packet is PHASE-06; neither state is reachable via CLI this phase)
   - unknown state ⇒ FAIL
   - `archived` is a legal stored state this phase (no extra invariant)
5. **Legal tier values** ∈ {focused, standard, high-assurance}.
6. **Tier-appropriate presence (binds).** Missing bound files/dirs fail. Extra valid artifacts at a low tier pass.
7. **Parseable log prefixes.** Every non-empty, non-`#` heading line matches `^\d{4}-\d{2}-\d{2}\t(scaffold|new|tier|publish|check)\t(\S+)\t`. Blank lines and ATX headings are allowed. **Architect default:** `check` itself does not append a log line.
8. **Leftover journal / stale lock.** Fail and name recovery (`re-run` or `mycelium check --abort-journal`).
9. **Declared-deviation rule.** Undeclared extra top-level convention-breaking paths fail. Declared rows need `convention` + `reason`.
10. **Stage-scoped IDs fall inside a declared range.** A FND/REC/REQ file whose number is outside every declared range for that key fails even if a human wrote it.
11. **Teaching errors:** name convention, link contract path, suggest fix command.

### Transition table (encoded now, commanded later)

```text
spark       → exploring, archived
exploring   → simmering, clarified, archived
simmering   → exploring, archived
clarified   → handed-off, archived
handed-off  → archived
archived    → (none)
```

PHASE-01 has no state-transition command. Check still encodes the table (unit-tested). Stored `clarified` / `handed-off` fail as above. A future command will consult this table.

### Allowed top-level paths (undeclared extras fail)

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

### What check does not do

- Grade thinking quality
- Require a minimum artifact count
- Require glossary entries
- Require cruxes on questions (PHASE-03)
- Require wake briefs (PHASE-02)
- Read embedded schemas
- Mutate files (except `--abort-journal`)

## 13. Vertical slices (binding build order)

Each slice ends checkable. PR-per-slice, sequential, rebase on main. Arvo merges Quality-green PRs. Do not stack unpublished slices on one branch unless Quality is backed up — prefer one live PR at a time. Slice 0 can land with Slice 1 if faster.

### Slice 0 — Commissioning artifacts

This brief + DEC-012 + DEC-013 + acceptance matrix stubs + new contracts. Docs/contracts only. No Go yet.

Land (see §16 for the full path list):

- `framework/phases/PHASE-01-implementation-brief.md` (this file)
- `framework/phases/PHASE-01-acceptance.md`
- `framework/decisions/DEC-012-manifest-filename-mycelium-toml.md`
- `framework/decisions/DEC-013-stage-range-refuse.md`
- `program/contracts/naming.md`
- `program/contracts/manifest.md`
- `program/contracts/lifecycle.md`
- `program/contracts/operation-protocol.md`
- `program/contracts/conformance.md`
- update `program/contracts/identifiers.md` (add ASM/EVD; keep v1 ranges as examples)

Done: files exist on a PR. Quality reads them against this brief.

### Slice 1 — CLI skeleton + embed + version

`go.mod` (1.26), `cmd/mycelium`, `internal/cli` dispatcher, `internal/version`, `internal/embed` of `program/` (even if program/ 2.0 content is still partial), `mycelium version` prints version.

Done (hermetic): `go test ./...` and `go build -o /tmp/mycelium ./cmd/mycelium` and `/tmp/mycelium version` equals the stamped version. No network.

### Slice 2 — program/ 2.0 content

Schemas, rewritten registered templates (with `+++` front matter), `tiers/*.toml`, naming/manifest/lifecycle/operation-protocol/conformance contracts (if not all in Slice 0), instance skeleton files, mycelium-cli skill. Master v1 `program/contracts` stay. ADD 2.0 files. Do not delete v1 operator/reference yet (master still uses them).

Done: every registered type has `template.md` + `schema.toml`; tiers parse; skeleton files exist under `program/skeleton/`.

Skeleton sources (copied to instance root on scaffold):

```text
program/skeleton/README.md
program/skeleton/log.md
program/skeleton/CONTEXT.md
program/skeleton/AGENTS.md
program/skeleton/gitignore
program/skills/mycelium-cli/SKILL.md
```

`mycelium.toml` is generated, not a static skeleton file (fields are instance-specific).

### Slice 3 — metadata reader + ID-to-path + manifest parse (pure)

`internal/metadata`, `internal/idpath`, `internal/manifest`, `internal/schema`. Unit tests only. No FS mutation yet except testdata.

Done: table-driven tests for ID↔path, front-matter bounds (body cannot masquerade), manifest required fields, range membership.

### Slice 4 — lock / journal / operation protocol

`internal/lock`, `internal/journal`, `internal/op`.

Done: unit tests on a temp dir: crash-before-rename leaves nothing; crash-after-partial-rename leaves a journal; resume uses original ID; stale lock detected.

### Slice 5 — `mycelium new idea --offline` scaffold

Emit skeleton + `program/` + `mycelium.toml` (`state=spark`, `tier=focused`, versions stamped) + `git init`, no commit, no gh.

Done (hermetic): binary scaffolds into a temp dir; check (Slice 6 may land same PR if needed, else a smoke that files exist + git repo); `framework/` absent from instance; `research-program.toml` absent; Justfile/scripts absent.

### Slice 6 — `mycelium check`

Schema-driven, tier-aware, teaching errors.

Done (hermetic): spark instance passes check; illegal state fails with teaching error; ID-to-path mismatch fails; leftover journal fails; undeclared extra convention-breaking file fails.

### Slice 7 — `mycelium new <type>`

Data-driven generator.

Done (hermetic): generate one of each registered type (declare ranges first for FND/REC/REQ); check green; refuse overwrite; refuse out-of-range; next ID = max+1; log line appended; tokens replaced.

### Slice 8 — `mycelium tier`

Done (hermetic): focused→standard emits new dirs, does not rewrite existing DEC; standard→focused deletes nothing, check still green (binds relax); idempotent second call is a no-op.

### Slice 9 — Fixture-instance CI (hermetic job)

CI job: build binary, scaffold `--offline`, declare ranges, generate one of every registered type, `mycelium check`, assert no network (offline flag + a test that `gh` / `git-remote` are not invoked).

Done: GitHub Actions workflow on push/PR, no secrets required, green on the PR. Path: `.github/workflows/phase-01-hermetic.yml`.

### Slice 10 — Authenticated GitHub integration (SEPARATE job)

`new idea --publish` OR `new idea` + `publish` creates a repo under the authenticated user, name `mycelium-ms101-<unix>` (**Architect default**), adds `idea` topic, records `github_repo` in the manifest, cleanup on failure AND on success (delete the fixture repo). Skip if GH token absent (job is allowed-to-skip, not a silent pass-as-green of the publish path — report `skipped: no credentials`).

Done: `workflow_dispatch` or secret-gated job (`.github/workflows/phase-01-github.yml`). Hermetic job does NOT depend on it. MS-101(b) is this job passing once with credentials.

## 14. Done / verified mapped onto MS-101

MS-101 has two parts. Both required to accept PHASE-01. Five-minute spark-to-first-thought is a USER SLO, not the phase gate.

### (a) Hermetic local (phase gate, no network)

- `go test ./...` green with network disabled (`MYCELIUM_OFFLINE=1` and/or `--offline` on all scaffold commands).
- Built binary scaffolds a conformant spark instance (`new idea --offline`); `mycelium check` exits 0.
- Fixture CI generates one of every registered artifact type and `mycelium check` exits 0.
- Instance does not contain `framework/`, `Justfile`, `scripts/`, or `research-program.toml`.
- Teaching errors covered by tests for: overwrite, out-of-range, illegal state, leftover journal, ID-to-path mismatch.

### (b) Authenticated GitHub integration (separately credentialed)

- A test with `gh` credentials publishes the repo with the `idea` topic and cleans up on failure (and success).
- This test MUST NOT run in the hermetic job.
- Absence of credentials skips this job; it does not fail MS-101(a).
- PHASE-01 is not accepted until (b) has passed at least once.

### Slice → MS-101 map

| Slice | MS-101 clause |
| --- | --- |
| 0 | commissioning; not a runtime clause |
| 1 | (a) `go test` / `go build` / `version` |
| 2 | (a) schemas + skeleton exist for fixture |
| 3 | (a) unit tests inside `go test ./...` |
| 4 | (a) leftover journal tests |
| 5 | (a) scaffold `--offline` + absence asserts |
| 6 | (a) `check` exits 0; teaching errors |
| 7 | (a) fixture generates every type |
| 8 | (a) tier tests inside `go test ./...` |
| 9 | (a) fixture CI job green |
| 10 | (b) credentialed job passed once |

PHASE-01 is accepted when (a) is green on main and (b) has a recorded pass. Arvo accepts the phase. Engineering does not self-accept.

## 15. Automated test plan

Engineering MUST write these tests. Quality thermos against this list. Do NOT require Playwright, Docker, or live GitHub in default `go test ./...`.

### Unit (no network, no gh)

- idpath both directions, all 11 types, bad IDs (`DEC-1` as filename, unknown NS, empty, `DEC-0001` four digits)
- metadata reader: `+++` fences; body `+++` does not count; required fields; BOM stripped; missing close
- manifest parse/validate; missing field; bad state; bad tier; bad `schema_version`; unknown key
- range membership (inside / outside / missing declaration / start>end / NS mismatch)
- slugify (spaces, punctuation, unicode accents, empty, overflow)
- token replace (all four tokens; leftover `{{FOO}}` stays; no engine)
- lock contention (two procs, one wins)
- journal resume / abort
- transition table (legal/illegal)

### Integration / hermetic CLI (exec the built binary, temp dirs)

- version equals stamped
- `new idea --offline` → spark `check` 0
- `new idea --offline` into existing dir → refuse
- generate all 11 types + `check` 0 (write ranges first)
- generate FND without range → refuse
- generate FND with range, then one past the end → refuse
- overwrite refuse
- tier up / tier down / idempotent
- leftover journal → check fails with recovery text
- `--offline` never execs `gh` (internal/execrun adapter; tests assert `gh` was not called)
- instance lacks `framework/`, `Justfile`, `scripts/`, `research-program.toml`
- teaching-error shape (four lines: mycelium / convention / contract / fix) for overwrite, out-of-range, illegal state, leftover journal, ID-to-path mismatch

### Credentialed (build tag `github_integration` OR env `MYCELIUM_GH_TEST=1`)

- publish + topic + cleanup
- second publish is idempotent
- cleanup on induced failure (e.g. topic-add fails after create)

Default `go test ./...` MUST NOT run the credentialed tests. Use `//go:build github_integration` on those files, or a skip when `MYCELIUM_GH_TEST` is unset. **Architect default:** build tag `github_integration` plus a skip if `gh auth status` fails.

Fixture repo name: `mycelium-ms101-<unix>`. `t.Cleanup` deletes it. Workflow `always()` also deletes.

## 16. PHASE-01 acceptance matrix / in-repo contract paths

Engineering lands these paths. Slice 0 may land the docs/contracts; later slices land code and workflows.

| Path | What it is |
| --- | --- |
| `framework/phases/PHASE-01-implementation-brief.md` | this brief |
| `framework/phases/PHASE-01-acceptance.md` | matrix; rows = MS-101 checks + per-slice done bars |
| `framework/decisions/DEC-012-manifest-filename-mycelium-toml.md` | Appendix A, landed as a DEC |
| `framework/decisions/DEC-013-stage-range-refuse.md` | Appendix B, landed as a DEC |
| `program/contracts/naming.md` | ID-to-path |
| `program/contracts/manifest.md` | `mycelium.toml` fields |
| `program/contracts/lifecycle.md` | state machine + PHASE-01 stored-state rules |
| `program/contracts/operation-protocol.md` | lock/journal/stage/commit/rollback |
| `program/contracts/conformance.md` | check rules + teaching errors |
| `program/contracts/identifiers.md` | update: add ASM/EVD; keep v1 ranges as examples |
| `program/tiers/focused.toml` | §7 |
| `program/tiers/standard.toml` | §7 |
| `program/tiers/high-assurance.toml` | §7 |
| `program/templates/<type>.md` + `<type>.schema.toml` | 11 types |
| `program/skeleton/*` | emitted files |
| `program/skills/mycelium-cli/SKILL.md` | thin skill |
| `.github/workflows/phase-01-hermetic.yml` | Slice 9 |
| `.github/workflows/phase-01-github.yml` | Slice 10, secret-gated |

### Acceptance matrix rows (copy into `PHASE-01-acceptance.md`)

Each row: id, check, evidence, owner (Engineering | CI | Arvo).

| id | check | evidence |
| --- | --- | --- |
| A-S0 | Slice 0 files exist and match this brief | PR diff; Quality read |
| A-S1 | `go test ./...` and `go build` and `version` | CI + local hermetic |
| A-S2 | 11 templates + 11 schemas + 3 tiers + skeleton + skill | tree exists |
| A-S3 | table-driven pure tests green | `go test` |
| A-S4 | crash/resume/stale-lock tests green | `go test` |
| A-S5 | `new idea --offline` spark; forbidden paths absent | hermetic CLI |
| A-S6 | check 0 on spark; teaching errors for illegal state, ID mismatch, journal, undeclared extra | hermetic CLI |
| A-S7 | 11 types generated; overwrite/out-of-range refuse; log line; tokens | hermetic CLI |
| A-S8 | tier up/down/idempotent; no deletion | hermetic CLI |
| A-S9 | hermetic workflow green on the PR | GitHub Actions |
| A-S10 | credentialed job passed once; hermetic does not depend on it | Actions log |
| MS-101a | all (a) bullets in §14 | A-S1..A-S9 |
| MS-101b | all (b) bullets in §14 | A-S10 |
| SLO-5m | spark-to-first-thought under five minutes | user SLO; not the phase gate |

## 17. Decided / Architect defaults

No TBD. Open items are decided inline and labeled **Architect default**. Do not bikeshed them in a code PR. Write a DEC to change one.

Index of defaults that are easy to miss:

- `+++` TOML front matter. ALLOW ID gaps; next = max+1; no tombstones.
- `mycelium check --abort-journal` is the rollback. No-deletion of renamed files.
- Default `new idea`: publish only if `gh` is authenticated and `--offline` is absent; otherwise succeed local and print `mycelium publish`.
- Hermetic tests must pass `--offline` or `MYCELIUM_OFFLINE=1`. `--offline` + `--publish` is a refuse.
- Default dir `./<slug>`. `--dir` is exact. Do not mkdir `~/ideas`.
- `git init -b main`. Never add/commit. `git`/`gh` via `internal/execrun` only.
- Embed: `internal/embed` + `go generate` copy. Filter `*.go`. Check reads instance files.
- Only `github.com/pelletier/go-toml/v2`. stdlib tests. Exit 0/1. `version` prints the string only.
- Clock UTC + `MYCELIUM_NOW`. Log `YYYY-MM-DD\t<op>\t<ID-or-->\t<title-or-note>`.
- Non-blocking flock. Live lock: check notices. Journal/stale lock: check fails.
- No `mycelium range`. Fixture writes `[identifiers]`. One range per key.
- Credentialed tests: build tag `github_integration`. Fixture repo `mycelium-ms101-<unix>`.
- clarified/handed-off stored → check fails (PHASE-02/06). archived stored → legal.
- Extra valid artifacts at focused: allowed. Lowering a tier deletes nothing.
- Manifest unknown keys: refuse. Title drift: not a failure. `id` must match filename.
- `.gitignore` contains `.mycelium/lock` only. Help/`-h`: usage, exit 0.
- Slice 0+1 may combine. One live PR at a time. `just check` on master must pass.

## 18. Risks, rollback, what Quality should refuse

### Risks and mitigations

| Risk | Mitigation |
| --- | --- |
| Master `just check` requires the v1 tree | Do not remove v1 files this phase. Rollback: revert the CLI PR; master remains ADRP v1. |
| `go.mod` on master is new | `just check` does not look at `go.mod`. Fine. |
| Accidental `just init` on master | Called out here. Engineering must not do this. Quality refuses a PR that rewrites master's v1 manifest/tree via init. |
| Publish test leaking fixture repos | Unique names `mycelium-ms101-<unix>` + `t.Cleanup` + workflow `always()` delete. |
| Embed path mistakes emitting `framework/` | Hermetic test asserts absence of `framework/`, Justfile, scripts, `research-program.toml`. |
| Partial operation on crash | Journal tests (Slice 4) are the mitigation. |
| Hermetic job silently depends on secrets | Two workflow files. Hermetic has no secrets. Quality refuses a single combined job. |
| Generator/checker type switch in Go | Adding a type must be template + schema only. Review the Go diff for new type names. |
| Content-quality lint creep | DEC-005. Refuse. |

### Rollback

Revert the offending PR on master. Do not `git push --force` to main. Do not delete Justfile/scripts as a "cleanup" rollback. Master v1 is the safe state until MS-101 is accepted.

### Quality should refuse

Refuse to approve if:

- any PHASE-02+ command or skill shipped as if in scope
- hermetic tests require network or `gh`
- CLI git-commits instance work product
- `framework/` is emitted into an instance
- generator or checker hardcodes types (must be schema-driven)
- a YAML / Jinja / cobra / viper dependency appears
- content-quality lint is added
- Justfile/scripts deleted from master before MS-101 accepted
- PR pushed straight to main
- MS-101(a) and (b) are conflated into one job
- out-of-range allocation warns instead of refuses
- instance manifest is still named `research-program.toml`
- runtime check reads embedded schemas instead of instance files
- overwrite of existing artifacts
- lowering a tier deletes files
- `just init` was run on master
- master's `research-program.toml` was renamed
- a destroy command or GitHub-template flow appears
- teaching errors omit convention / contract / fix
- Slice 10 is required to go green for the hermetic workflow

## 19. Overnight execution order

PR-per-slice, sequential, rebase on main. Arvo merges Quality-green PRs. Engineering opens PRs. Engineering does NOT push to main. Prefer one live PR at a time so review stays small. Slice 0 can land with Slice 1 if faster. Do not stack unpublished slices on one branch unless Quality is backed up.

Suggested order (same calendar night is fine; do not skip done bars):

1. Slice 0 PR (docs/contracts). Optional: combine with Slice 1.
2. Slice 1 PR — wait for merge (or combine with 0).
3. Slice 2 PR — program/ 2.0 content. Can overlap review of 1 only if 1 is already Quality-green and queued.
4. Slice 3 PR — pure parsers.
5. Slice 4 PR — protocol.
6. Slice 5 PR — scaffold `--offline`. Smoke if Slice 6 is not ready.
7. Slice 6 PR — check. May combine with 5 if both are small and 4 is merged.
8. Slice 7 PR — generator.
9. Slice 8 PR — tier.
10. Slice 9 PR — hermetic workflow. Must be green on that PR.
11. Slice 10 PR — secret-gated workflow + build-tagged tests. Must not break Slice 9. MS-101(b) needs one credentialed pass; Arvo or Robert triggers `workflow_dispatch`.

Each PR title: `PHASE-01 Slice N: <done-bar noun>`. Each PR body links this brief and the slice done bar. No drive-by refactors. No v1 deletions.

## 20. Handoff

### What Engineering starts with

This file. Only this file. Clone `https://github.com/robertguss/mycelium` at current `main`. Read `framework/blueprint.md` and DEC-001–011 for authority, not for a second plan. Execute Slice 0 first unless combining 0+1.

Cursor cloud: Go 1.26 at `/usr/local/bin/go`. Use it. `CGO_ENABLED=0`.

### What Engineering must not do

- Do not implement PHASE-02+ commands or skills.
- Do not convert master into a 2.0 instance.
- Do not rename master's `research-program.toml`.
- Do not run `just init` on master.
- Do not delete Justfile/scripts.
- Do not emit `framework/`.
- Do not git-commit instance work product from the CLI.
- Do not push to main.
- Do not add cobra / viper / yaml / testify / go-github.
- Do not grade content.
- Do not warn on out-of-range allocation.
- Do not combine MS-101(a) and (b) into one job.
- Do not open a design debate in the PR. Open items were decided in §17.
- Do not write a second brief.

### What Quality reads

This brief, the acceptance matrix, the DEC-012/013 files, and the PR diff. Thermos: §15 tests exist and match; §18 refuse list is clean; MS-101(a) hermetic; Slice 10 isolated.

### What Arvo does

Merges Quality-green PRs. Accepts PHASE-01 when MS-101(a) is green on main and MS-101(b) has passed once.

## Appendix A — DEC-012 text

Engineering lands this file at `framework/decisions/DEC-012-manifest-filename-mycelium-toml.md`.

```text
# DEC-012 — New instances use mycelium.toml as the manifest filename

- **Status:** Accepted
- **Date:** 2026-08-14
- **Owner:** Robert Guss
- **Supersedes:** None (settles blueprint OQ-002)
- **Related recommendations:** None
- **Related evidence:** DEC-010 (CLI scaffold); DEC-011 (two version fields)

## Context

Blueprint OQ-002 asked whether new instances keep `research-program.toml`
or rename as a 2.0 migration. The master repo remains an ADRP v1 instance
and already has `research-program.toml`. New idea repos are Mycelium 2.0
from birth.

## Decision

1. Scaffolded instances use `mycelium.toml` as the sole manifest filename.
2. The master repository keeps `research-program.toml` and is not converted.
3. Runtime commands detect an instance by the presence of `mycelium.toml`.
   They do not read `research-program.toml` as a 2.0 manifest.
4. No migration machinery is added to rename existing files (DEC-011).

## Rationale

The 2.0 name matches the product. Keeping master's v1 filename avoids
pretending the master is a scaffolded idea and keeps `just check` working.

## Consequences

PHASE-01 contracts, checker, and scaffold all say `mycelium.toml`.
Documentation that still says `research-program.toml` refers to master v1.

## Alternatives Considered

Keep `research-program.toml` everywhere (rejects the 2.0 name).
Rename master too (converts master into an instance; rejected).

## Risks

Agents open master and look for `mycelium.toml`. Mitigation: AGENTS.md on
master stays v1; instance AGENTS.md is emitted and names `mycelium.toml`.

## Revisit Triggers

If master itself is ever re-scaffolded as a Mycelium instance, reopen.

## Approval

Settled by Arvo 2026-08-14 (OQ-002). Recorded by Architect in the
PHASE-01 implementation brief.
```

## Appendix B — DEC-013 text

Engineering lands this file at `framework/decisions/DEC-013-stage-range-refuse.md`.

```text
# DEC-013 — Refuse allocation outside a declared stage-scoped range

- **Status:** Accepted
- **Date:** 2026-08-14
- **Owner:** Robert Guss
- **Supersedes:** None (settles blueprint OQ-007)
- **Related recommendations:** None
- **Related evidence:** program/contracts/identifiers.md (v1 ranges)

## Context

FND, REC, and REQ are stage-scoped. The blueprint allocated ranges per
stage so identifiers do not collide across reviews. OQ-007 asked whether
the generator should warn or refuse when the next ID falls outside a
declared range, or when no range is declared.

## Decision

1. Stage-scoped types (finding, recommendation, requirement) require a
   declared range in the instance manifest `[identifiers]` table before
   allocation.
2. If no range is declared, the generator REFUSES.
3. If the next ID would fall outside all declared ranges, the generator
   REFUSES.
4. A warning is not sufficient. Check also fails existing files whose
   IDs sit outside every declared range for that key.
5. Non-stage-scoped types (DEC, ASM, EVD, SPK, OQ, RSK, PHASE, MS) do
   not require a range.

## Rationale

A warning is a log line agents ignore. A refuse is a teaching error with
a contract link. Ranges exist to keep stage traces disjoint; silent
overflow defeats them.

## Consequences

Fixture CI must write ranges before generating FND/REC/REQ. There is no
`mycelium range` command in PHASE-01; tests edit the manifest.

## Alternatives Considered

Warn and allocate (rejected: unenforceable). Auto-extend the range
(rejected: hides the stage boundary).

## Risks

Sparks that want a finding before declaring a range are blocked.
Mitigation: focused sparks do not bind findings; declaring a range is
one manifest edit.

## Revisit Triggers

If a real idea needs multiple disjoint ranges per type in PHASE-01,
reopen to allow an array of ranges. This brief allows one range per key.

## Approval

Settled by Arvo 2026-08-14 (OQ-007 → REFUSE). Recorded by Architect in
the PHASE-01 implementation brief.
```

## Appendix C — `mycelium.toml` example (spark / focused)

```toml
schema_version = 1
idea_name = "Garden lighting"
slug = "garden-lighting"
state = "spark"
tier = "focused"
methodology_version = "2.0.0"
generated_by_cli_version = "0.1.0-dev"
created_date = "2026-08-14"
updated_date = "2026-08-14"
revisit = ""
github_repo = ""

[identifiers]
# fixture CI writes these before generating FND/REC/REQ:
# findings = "FND-001..FND-099"
# recommendations = "REC-001..REC-099"
# requirements = "REQ-001..REQ-099"
```

## Appendix D — Sidecar schema example (decision)

`program/templates/decision.schema.toml`:

```toml
namespace = "DEC"
home = "decisions"
filename_pattern = "DEC-{NNN}-{slug}.md"
stage_scoped = false
digits = 3
required_front_matter = ["id", "title", "status", "date", "owner"]
required_sections = [
  "Context",
  "Decision",
  "Rationale",
  "Consequences",
  "Alternatives Considered",
  "Risks",
  "Revisit Triggers",
  "Approval",
]

[enums.status]
values = ["Proposed", "Accepted", "Superseded", "Rejected"]
```

`program/templates/decision.md` (abridged):

```text
+++
id = "{{ID}}"
title = "{{TITLE}}"
status = "Proposed"
date = "{{DATE}}"
owner = ""
+++

# {{ID}} — {{TITLE}}

## Context

<!-- fill -->

## Decision

<!-- fill -->

## Rationale

<!-- fill -->

## Consequences

<!-- fill -->

## Alternatives Considered

<!-- fill -->

## Risks

<!-- fill -->

## Revisit Triggers

<!-- fill -->

## Approval

<!-- fill -->
```

## Appendix E — Teaching-error format

Stderr. Exit 1. Four lines, in this order:

```text
mycelium: <one-line failure>
convention: <name>
contract: program/contracts/<file>.md
fix: <command or rename>
```

Examples:

```text
mycelium: refuse overwrite decisions/DEC-001-garden-lighting.md
convention: no-overwrite
contract: program/contracts/naming.md
fix: choose a new title or rename the existing file

mycelium: FND-100 is outside declared range FND-001..FND-099
convention: stage-range
contract: program/contracts/identifiers.md
fix: widen [identifiers].findings in mycelium.toml, then retry

mycelium: state=handed-off is not reachable in PHASE-01
convention: lifecycle
contract: program/contracts/lifecycle.md
fix: restore state to spark|exploring|simmering|archived (PHASE-02/06 commands are not shipped)

mycelium: interrupted operation
convention: operation-protocol
contract: program/contracts/operation-protocol.md
fix: re-run the same command to complete, or mycelium check --abort-journal to roll back

mycelium: reference DEC-014 has no file
convention: id-to-path
contract: program/contracts/naming.md
fix: add the artifact (mycelium new decision "...") or remove the reference
```

`check` success is a short OK summary on stdout (see §8.4), exit 0.

## Appendix F — Target file tree

### Master (after PHASE-01; v1 files retained)

```text
cmd/mycelium/main.go
internal/cli/ ...
internal/embed/embed.go
internal/embed/program/          # generate copy
internal/version/ ...
go.mod
go.sum
.github/workflows/phase-01-hermetic.yml
.github/workflows/phase-01-github.yml
framework/blueprint.md
framework/decisions/DEC-001 … DEC-013
framework/phases/PHASE-01-implementation-brief.md
framework/phases/PHASE-01-acceptance.md
program/contracts/               # v1 + naming/manifest/lifecycle/operation-protocol/conformance
program/templates/*.md           # 11 registered + remaining v1 unregistered
program/templates/*.schema.toml  # 11
program/tiers/{focused,standard,high-assurance}.toml
program/skeleton/
program/skills/mycelium-cli/SKILL.md
program/operator/                # v1, keep
program/reference/               # v1, keep
Justfile                         # v1, keep
scripts/                         # v1, keep
research-program.toml            # v1 master manifest, keep
.agents/skills/research-*        # v1, keep, do not emit
```

### Emitted instance (spark / focused, local-only)

```text
README.md
mycelium.toml
log.md
CONTEXT.md
AGENTS.md
.gitignore
.agents/skills/mycelium-cli/SKILL.md
program/contracts/
program/templates/
program/tiers/
program/skeleton/
program/skills/mycelium-cli/SKILL.md
program/reference/               # present if embedded
.git/                            # init only; no commit
```

Absent from the instance: `framework/`, `cmd/`, `internal/`, `go.mod`, `Justfile`, `scripts/`, `research-program.toml`, v1 `research-*` skills, master's `docs/`, master's `decisions/`.

After `mycelium tier standard`: add `decisions/`, `assumptions/`, `evidence/`, `questions/`, `risks/` (README stubs if new).

After generating one decision: `decisions/DEC-001-<slug>.md` plus one log line.

After `mycelium publish`: `github_repo` set; `origin` remote; GitHub `idea` topic.

End of PHASE-01 implementation brief. Engineering executes from this file only.
