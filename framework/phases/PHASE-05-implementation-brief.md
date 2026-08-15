# PHASE-05 Implementation Brief — Distribution and lifecycle commands

- **Status:** Binding
- **Date:** 2026-08-15
- **Audience:** Engineering (pstack / poteto-mode)
- **Authority:** `framework/blueprint.md` (Accepted 2026-08-14). DEC-001 through DEC-014 stand. DEC-011 is this phase's risk guard. Do **not** record DEC-015 (Appendix A). Do **not** reopen DEC-012 / DEC-013 / DEC-014.
- **Repo:** https://github.com/robertguss/mycelium
- **Pin:** Engineering starts from `main` @ `657a14653da1c41fd4c0590a5b0aa625eaa9adde` (PHASE-01–04 accepted). Do not implement from a later SHA unless Arvo re-pins in writing.
- **Product:** single-binary Go CLI `mycelium`. Master builds the CLI. `program/` is `go:embed`'d and emitted on scaffold. PHASE-05 adds **`mycelium supersede`**, `CHANGELOG.md` discipline, a hermetic release script with `SHA256SUMS`, install docs, and DEC-011 old-manifest tolerance on `status` / `check`. **No new pack.** **No `mycelium council`.** **`handed-off` stays unreachable.**
- **Phase gate:** MS-501 via files + hermetic `go test ./...` + Quality's local guide. GitHub Actions is **not** a gate (Robert waived CI). Do **not** add `.github/workflows/phase-05-*.yml`. Do **not** extend `phase-01-hermetic.yml` as a phase gate. Quality should refuse a PR that adds an Actions job as the MS-501 gate. Quality should **not** refuse a missing PHASE-05 workflow (absence is correct).
- **Not the gate:** Install SLO (clean VM, one-line install to a scaffolded instance in under a minute) is **human evidence**. Cutting the real GitHub tag and uploading a GitHub Release are **human evidence**. Do not put them in a done bar.
- **How to use this file:** Engineering executes from THIS FILE ONLY. No "see chat". No TBD. Open items are decided here and labeled **Architect default**. Not-decided is empty (§15).

Headings: §§1–18 then Appendices A–F (no new DEC; supersede before/after; golden old-manifest snippets; CHANGELOG + SHA256SUMS; MS-501 fixture recipe; file tree / DO NOT ADD).

Cloud env name is exactly `robertguss/mycelium`. Go 1.26 at `/usr/local/go` or `/usr/local/bin/go`. `CGO_ENABLED=0`. stdlib + `github.com/pelletier/go-toml/v2` only. No cobra / viper / yaml / testify / go-github. Linear parent for slice tickets: **ROB-512** (Arvo cuts tickets after this brief is stamped).

## 1. Scope / out of scope

Tonight is PHASE-05 only. Do not implement, stub-ship, or leave a hook for PHASE-06. Do not reimplement PHASE-01–04. Do not convert master.

### In scope

- First new top-level verb since PHASE-02: `mycelium supersede <OLD-ID> --by <NEW-ID> [--dir PATH]`. Artifact-level, not idea-lifecycle. Teaching errors on illegal use.
- Bidirectional front-matter cross-links (`status = "Superseded"`, `superseded_by`, `supersedes`) plus a `supersede` log line. No required new H2.
- Eligible-type schema deltas (§10). OQ is not eligible.
- Operation protocol for `supersede` (preflight / lock / stage / commit / rollback). Journal op adds `supersede`.
- DEC-011 risk guard: `status` and `status --all` never abort a whole scan because one instance is unreadable. Emit `partial: legacy-manifest (<path>: <reason>)` and continue. Golden fixtures G0–G3 (§5, Appendix C).
- Master `CHANGELOG.md` (Keep a Changelog, semver). CLI version stays `0.1.0-dev` until a human cuts `v0.1.0`.
- `scripts/release.sh` (Justfile recipe may call it): `CGO_ENABLED=0` builds `linux-amd64` and `darwin-arm64` only; writes `dist/SHA256SUMS`; refuses if CHANGELOG lacks `## [<version>]`.
- `docs/install.md` names the one-liner. Install SLO is human evidence. Hermetic: file exists and names the one-liner.
- Check updates in §8 (log-op `supersede`; bidirectional IFF; instance-contract required keys for G1).
- `program/skills/mycelium-cli/SKILL.md` and `program/skeleton/AGENTS.md` name `mycelium supersede`. No new skill. No new pack.
- MS-501 hermetic fixture matrix in `go test ./...` (`internal/clitest`). No Actions job.
- Commissioning artifacts: this brief, acceptance stub, conformance lift timing.

### Out of scope (one refuse list; details in §16)

- PHASE-06 handoff packet. `handed-off` stays unreachable.
- Migrations: no `program/migrations/`, no `just upgrade`, no `applied_migrations` (DEC-011).
- New pack. `mycelium council`. Portable council CLI. Reopening OQ-003.
- Windows binary. goreleaser. A PHASE-05 Actions workflow. Extending `phase-01-hermetic.yml` as a phase gate.
- Install SLO or a live GitHub Release / real tag as a merge gate.
- Changing `mycelium.toml` `state`. Implementing `handed-off`.
- Reopening DEC-012 / DEC-013 / DEC-014. Recording DEC-015.
- Emitting `framework/`. CLI `git add` / `git commit` of instance work product.
- Converting master (`research-program.toml`, `just init`, deleting Justfile / v1 scripts).
- Growing `latinFold`, NFKD, `golang.org/x/text` (DEC-014).
- Implementing MS-101(b). Commissioning a `GH_TOKEN` job.

### Master vs instance (unchanged)

Master remains an ADRP v1 instance. Do not convert `research-program.toml` to `mycelium.toml`. Do not run `just init` on master. `framework/` is master-only and is NEVER emitted. Justfile and v1 `scripts/*.py` stay. Runtime detects idea instances by `mycelium.toml`.

## 2. Authority and do-not-reopen

### Governing documents

| Document | Rule |
| --- | --- |
| `framework/blueprint.md` (Accepted 2026-08-14) | Do not rewrite vision. DEC-001–014 stand. PHASE-05 ceiling is blueprint lines ~432–439 (`mycelium supersede`; tagged releases; `CHANGELOG.md`; portfolio scanner tolerant of older shapes). Supersede grammar ~302–303. Operation protocol includes supersede ~315. DEC-011 seams ~199–200, ~240–244. Portfolio hardening ~375–376. |
| **DEC-011** | **This phase's risk guard.** No migrations. Instance files are truth. `CHANGELOG.md` is the release face. Portfolio scanner must tolerate older manifest shapes. |
| DEC-005 | Checks validate containers, never contents. |
| DEC-006 | spark → exploring ⇄ simmering → clarified → handed-off; any → archived. Do not change the machine. `handed-off` stays unreachable. |
| DEC-010 | CLI never git-commits instance work product. |
| DEC-012 / DEC-013 / DEC-014 | Do not reopen (`mycelium.toml`; refuse out-of-range; `latinFold` only, no NFKD, no `x/text`, do not grow the map). |
| This brief | Binding 2026-08-15. PHASE-05 only. Architect defaults are binding. No DEC-015. |

### Process override (unchanged)

Blueprint "humans-own-git" is overridden for the *master* repo's engineering process: Arvo merges Quality-green PRs and accepts the phase. Engineering opens PRs. Engineering does NOT push to main. The CLI still never git-commits *instance* work product.

### Do not reopen

Do not reopen the product shape, language, dependency floor, state vocabulary, manifest filename, refuse-vs-warn range rule, no-commit rule, instance-files-are-truth rule, slugify/DEC-014, publish, MS-101(b), PHASE-03 sparring, PHASE-04 ladder / OQ-003 (only `council` is a pack; no portable council CLI). If a later discovery contradicts a locked item, stop and write a DEC; do not silently patch this brief in a code PR.

Do not reopen DEC-012, DEC-013, or DEC-014. Do **not** record DEC-015. Do not commission PHASE-06.

## 3. Floor on main (do not reimplement)

Pin: `657a14653da1c41fd4c0590a5b0aa625eaa9adde`. Treat this SHA as the floor. Reuse packages. Do not rewrite working PHASE-01–04 commands.

### Already shipped (do not rebuild)

Reuse: `cmd/mycelium`, `internal/{cli,version,embed,clock,execrun,metadata,idpath,manifest,schema,slug,logfmt,teach,lock,journal,op,scaffold,generate,check,tiercmd,publish,clitest,revisit,lifecycle,indexmd,wakebrief,statecmd,statuscmd,sparring,pack,ladder}`.

| Touch | Fate |
| --- | --- |
| `internal/cli` | **Add one verb:** `supersede`. Existing verbs stay: `version`, `help`, `new`, `check`, `tier`, `publish`, `index`, `state`, `wake`, `status`. Do not add `council` / `handoff` / `upgrade` / `release`. |
| `internal/supersede` | **New (Slice 1).** Parsers / eligibility / refuse / cross-link rules. No CLI in Slice 1. |
| `internal/supersedecmd` | **New (Slice 2).** Command + protocol bind. |
| `internal/op` | **Extend** `allowReplace` so a `supersede` journal may replace the two existing artifact files (plus the existing regenerable set: `log.md`, `mycelium.toml`, `index.md`). Do not open general overwrite for `new`. |
| `internal/journal` | Op string `supersede` is legal. `OriginalID` = OLD-ID. `Title` = `<OLD-ID> -> <NEW-ID>`. No new journal schema version. |
| `internal/logfmt` / `internal/check` | Log-op regex gains `supersede` (item 8). Item 23 IFF binds in Slice 2. |
| `internal/manifest` | **Extend** with `ParseTolerant` for `status` (Slice 3). Strict parse stays for current-instance check except G1 required-key rule in §5. |
| `internal/statuscmd` | **Extend** (Slice 3): never abort a scan; emit `partial: legacy-manifest (<path>: <reason>)`. Existing `partial: local-only` stays. |
| `internal/check` | **Extend** items 8 + 23 + G1 required-key rule. Do not rewrite the package. |
| `internal/indexmd` | Reuse. `supersede` rewrites `index.md` the same way `new` / `state` do. **No new required H2.** |
| `internal/lifecycle` / `statecmd` | Do not touch. Supersede does not change `state`. |
| `internal/slug` | Do not touch (DEC-014). |
| `internal/version` | Stay `"0.1.0-dev"`. Do not bump as a phase ritual. `methodology_version` stays `2.0.0`. |
| `internal/embed` | Re-run `go generate` after `program/` edits. |
| `program/contracts/conformance.md` | Item 8 regex + item 23 + lift timing. Items 1–22 stay; do not renumber. |
| `program/skills/{spark,wake,portfolio,thinking}` + council pack | Do not rewrite. Update `mycelium-cli` + `AGENTS.md` only. |
| `phase-01-hermetic.yml` / `phase-01-github.yml` | Leave alone. Do **not** add a PHASE-05 workflow. Actions is not a gate. |
| Justfile / v1 `scripts/*.py` / `research-program.toml` | Keep. Add `scripts/release.sh` and a Justfile recipe that calls it. Do not delete v1 scripts. Do not `just init`. |
| `framework/` | Master-only. NEVER emitted. |
| `CHANGELOG.md` | **Absent on the pin.** Add at repo root (Slice 4). |

### Pin facts (do not "discover" otherwise)

- CLI verbs on the pin: `version` / `help` / `new` / `check` / `tier` / `publish` / `index` / `state` / `wake` / `status`. **No supersede yet.**
- `internal/version.Version = "0.1.0-dev"`. **No `CHANGELOG.md` on the pin.**
- Decision schema already has `status` enum `Proposed \| Accepted \| Superseded \| Rejected`. EVD and SPK already include `Superseded`. ASM does not (add it). OQ uses `agreement`, not `status`.
- Manifest required fields (current contract): `schema_version`, `idea_name`, `slug`, `state`, `tier`, `methodology_version`, `generated_by_cli_version`, `created_date`, `updated_date`, `revisit`, `github_repo`. Unknown top-level keys refuse on **check**.
- Master uses `research-program.toml`. Runtime detects instances by `mycelium.toml`.
- Log-op regex on the pin: `scaffold\|new\|tier\|publish\|check\|state\|wake`.
- `status --all` already continues past a bad child (writes a four-line teach to stderr). Replace that path with `partial: legacy-manifest` (do not abort). Single-instance `status` still uses strict `loadManifest` on the pin — switch it to tolerant parse (§5).
- `op.allowReplace` on the pin allows `log.md`, `mycelium.toml`, `index.md`, `briefs/*.md` only. Artifact files are not replaceable until Slice 2 extends this for `supersede` only.

### What must not be broken

`just check` on master; hermetic `go test ./...`; no `framework/` emit; no master conversion. PHASE-01–04 fixtures stay green. Spark with zero questions stays green. Instances without the council pack stay green. Stored `handed-off` still FAILS check. `mycelium state handed-off` still refuses.

If a PHASE-05 PR is bad: revert that PR. Floor is the pin SHA.

## 4. `mycelium supersede` (grammar, refuse, protocol)

This is the first new top-level verb since PHASE-02.

### Grammar

```text
mycelium supersede <OLD-ID> --by <NEW-ID> [--dir PATH]
```

| Token | Rule |
| --- | --- |
| `OLD-ID` / `NEW-ID` | Exact artifact IDs (`DEC-001`, `ASM-014`). Same namespace. Both must already exist. No deletion. No create. |
| `--by` / `--by=` | Required. Names NEW. |
| `--dir` / `--dir=` | Optional. Existing instance-root walk. |
| `-h` / `--help` | Usage, exit 0. |
| Other flags / extra positionals | Refuse (teaching error). |

**Architect default — not idea-lifecycle.** The command does not change `mycelium.toml` `state`. It does not implement `handed-off`. It does not accept idea-state tokens (`spark`, `exploring`, `simmering`, `clarified`, `handed-off`, `archived`) as IDs.

**Architect default — one-to-one this phase.** NEW may have at most one `supersedes`. OLD may have at most one `superseded_by`. Chains are legal only by superseding the *current* record (OLD is the live tip, not a past pair). Do not rewrite a past pair.

### Cross-links (bidirectional, structural)

Front matter is the container. **No required new H2.**

| File | Writes |
| --- | --- |
| OLD | `status = "Superseded"` (type has `status`) and `superseded_by = "<NEW-ID>"` |
| NEW | `supersedes = "<OLD-ID>"` (NEW's `status` is unchanged) |

Check link-resolves both IDs (item 6) and binds the IFF in item 23.

### Eligible types

Eligibility rule (closed): a generated type is eligible iff its sidecar schema has a `status` field **and** the enum **includes** `Superseded` after this phase's schema deltas.

| Type | Pin `status` enum | This phase | Eligible? |
| --- | --- | --- | --- |
| DEC | `Proposed \| Accepted \| Superseded \| Rejected` | add optional `superseded_by` / `supersedes` | **Yes** |
| ASM | `Open \| Held \| Falsified \| Retired` | add `Superseded` + optional link keys | **Yes** |
| EVD | already includes `Superseded` | add optional link keys | **Yes** |
| SPK | already includes `Superseded` | add optional link keys | **Yes** |
| OQ | `agreement`, not `status` | none | **No** — refuse; teaching error points at opening a new OQ |
| REC / REQ / RSK / FND | no `status` field | **not applicable** — do not invent a `status` field | **No** |
| PHASE | `Planned \| In Progress \| Accepted \| Abandoned` | none — that enum is phase-lifecycle | **No** |
| MS / CMP / RPT / RCL | no `status` | none | **No** |

### Refuse (teaching errors; four-line `teach` format; exit 1)

| Case | `what` (stderr line 1) | `fix` |
| --- | --- | --- |
| Missing OLD, missing `--by`, or missing NEW | `supersede requires <OLD-ID> --by <NEW-ID>` | `mycelium supersede DEC-001 --by DEC-002` |
| OLD or NEW not found | `no artifact <ID>` | `mycelium new <type> "…" then retry` |
| Different namespace | `supersede requires the same namespace (got DEC vs ASM)` | pick two IDs in one namespace |
| Type ineligible (incl. OQ) | `type <NS> is not supersedable` | OQ: `open a new question; do not supersede an OQ` |
| Idea-state token as an ID | `spark` (etc.) `is an idea state, not an artifact` | `mycelium state <target>` for lifecycle; `mycelium supersede` for artifacts |
| OLD already `status = "Superseded"` | `<OLD-ID> is already Superseded by <existing>` | `supersede the current record (<existing>) --by <newer>` |
| NEW already has `supersedes` set | `<NEW-ID> already supersedes <existing>` | one-to-one this phase; pick a different NEW |
| OLD == NEW | `cannot supersede an ID with itself` | pass two different IDs |
| Not an instance | `not a mycelium instance (no mycelium.toml found)` | `--dir PATH` |
| Leftover journal for a different op | existing `ErrJournalMismatch` text | `mycelium check --abort-journal` |

Do not delete files on refuse. Do not write a log line on refuse.

### Operation protocol

Same protocol as `new` / `tier` (blueprint ~315): preflight / lock / stage / commit / rollback. Journal `op = "supersede"`.

| Step | Rule |
| --- | --- |
| Preflight | Manifest + log parse; both IDs exist; same namespace; eligible; refuse table clean; no conflicting leftover journal. Nothing written. |
| Lock | Exclusive repo lock. |
| Stage | Write updated OLD, updated NEW, rewritten `index.md`, appended `log.md`, bumped `mycelium.toml` as staged files. Journal argv = the invocation. `OriginalID` = OLD-ID. `Title` = `<OLD-ID> -> <NEW-ID>`. |
| Commit order | **OLD file, NEW file, `index.md`, `log.md`, `mycelium.toml` last.** `mycelium.toml` changes `updated_date` only. |
| Rollback | Failure before the first rename: remove staged files, change nothing. After a partial commit: journal survives; re-running the same argv resumes under the original OLD-ID. |
| Detection | `mycelium check` already detects leftover journal / stale lock (item 9). `supersede` journals use the same recovery (`complete` or `--abort-journal`). |

**Architect default — `allowReplace`:** for `op == "supersede"` only, Commit may replace the two artifact `RelTo` paths in the journal. Do not allow `new` to overwrite an existing artifact.

**Architect default — log line:**

```text
YYYY-MM-DD<TAB>supersede<TAB><OLD-ID><TAB><OLD-ID> -> <NEW-ID>
```

Date from `internal/clock` / `MYCELIUM_NOW`. ID column is OLD-ID. Note is exactly `<OLD-ID> -> <NEW-ID>`.

**Architect default — success stdout:**

```text
mycelium supersede: ok
old: DEC-001
new: DEC-002
```

Exit 0. Then `mycelium check --dir PATH` exits 0 on the happy path.

**Architect default — index.md:** rewrite via existing `internal/indexmd`. Required H2s stay State / Artifacts / Log tail / Wake. Do not add a required H2.

**Architect default — resume match:** `journal.Matches` on argv, or `op+type+title` with `Type` = type key (`decision`, `assumption`, `evidence`, `spike`) and `Title` = `<OLD-ID> -> <NEW-ID>`.

## 5. Old-manifest tolerance (DEC-011)

DEC-011 deferred migrations and named the risk: a future manifest-format change could break `mycelium status --all`. Mitigation: treat the manifest format as append-only and keep the portfolio scanner tolerant of older shapes. That mitigation is this phase.

### Binding rules

| Surface | Rule |
| --- | --- |
| `status` and `status --all` | **Never abort the whole scan** because one instance is unreadable. Emit `partial: legacy-manifest (<path>: <reason>)` and continue. Exit 0 if the command itself is well-invoked. |
| `partial: local-only` | Unchanged (offline / `gh` missing / `gh` failed). A scan may emit both partial lines. |
| `status` parser | **Architect default:** `manifest.ParseTolerant`. Ignore unknown top-level keys (append-only). Missing `github_repo` → `""`. If the file is not TOML or yields no usable `slug`, do not crash: emit `partial: legacy-manifest` and skip the idea row. |
| `check` | Still instance-files-are-truth. Required keys come from the instance's `program/contracts/manifest.md` `## Required fields` table (pipe-table first column). Unknown top-level keys still refuse **when the instance contract says so** (current contract does). |
| G2 | A directory with `research-program.toml` and no `mycelium.toml` is **not** a mycelium instance. `status --all` does not crash. Emit `partial: legacy-manifest (<path>: research-program.toml without mycelium.toml)` and skip the idea row. `check` is N/A (FindRoot fails). |

### Golden fixtures (Appendix C has bytes)

| ID | Shape | `mycelium check` | `mycelium status` / `status --all` |
| --- | --- | --- | --- |
| G0 | Current `mycelium.toml` (control) | pass (exit 0) | pass (exit 0); instance listed |
| G1 | `mycelium.toml` **missing** `github_repo`, plus a frozen instance `program/contracts/manifest.md` that does **not** require it | pass (exit 0) | pass (exit 0); instance listed |
| G2 | Directory with `research-program.toml` and **no** `mycelium.toml` | N/A (not an instance) | `status --all` does not crash; skip / legacy row via `partial: legacy-manifest` |
| G3 | `mycelium.toml` with an extra unknown top-level key | **FAIL** (current instance schema refuses unknown keys) | **pass** — status treats as readable (append-only). **G3 is a status-only golden, not a check-pass.** |

State G3 clearly in tests and in `PHASE-05-acceptance.md`: do not assert `mycelium check` exit 0 on G3.

### G1 required-key rule (Architect default)

`internal/check` item 1 reads the instance's `program/contracts/manifest.md`. A key is required iff the `## Required fields` table names it. G1's frozen copy omits `github_repo`. If that file is missing, fall back to the current hardcoded set (includes `github_repo`). Do not parse the binary's embedded contract when an instance file exists.

Do not add `program/migrations/`. Do not add `applied_migrations`. Do not add `just upgrade`.

## 6. CHANGELOG + release + install docs

### `CHANGELOG.md` (master repo root)

Keep a Changelog + semver. **Absent on the pin; add it.**

**Architect default — first headings:**

```text
## [Unreleased]

## [0.1.0] - 2026-08-15
```

The `## [0.1.0] - 2026-08-15` heading is a **fixture** used by tests and by `scripts/release.sh 0.1.0`. CLI version stays `0.1.0-dev` until a human cuts `v0.1.0`. **Cutting the real GitHub tag is human evidence, not the gate.**

### `scripts/release.sh`

**Architect default:** `scripts/release.sh <version>` (Justfile `release version:` may call it). No goreleaser. No Windows this phase.

| Step | Rule |
| --- | --- |
| Refuse | If `CHANGELOG.md` lacks a heading `## [<version>]` (optional ` - YYYY-MM-DD` suffix allowed), exit non-zero. Print a teaching line naming the missing heading. Write nothing under `dist/`. |
| Build | `CGO_ENABLED=0` `go build` → `dist/mycelium-linux-amd64` and `dist/mycelium-darwin-arm64` only (`GOOS`/`GOARCH`). `-ldflags "-X github.com/robertguss/mycelium/internal/version.Version=<version>"` is allowed on the *release* binaries; the committed `internal/version.Version` stays `0.1.0-dev`. |
| Checksums | Write `dist/SHA256SUMS` (two lines, `sha256sum` format: `<hex>  mycelium-<os>-<arch>`). |
| Do not | `git tag`, `git commit`, `gh release`, upload, Windows, linux-arm64, darwin-amd64. |

Hermetic tests (Slice 4): refuse path + checksum match on a temp tree. Do not require a live GitHub Release.

### `docs/install.md`

Documents the one-liner. **Architect default one-liner** (must appear verbatim in the file):

```text
curl -fsSL https://github.com/robertguss/mycelium/releases/latest/download/mycelium-$(uname -s | tr A-Z a-z)-$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/') -o ~/.local/bin/mycelium && chmod +x ~/.local/bin/mycelium
```

Hermetic: file exists and contains `releases/latest/download/mycelium-` plus `~/.local/bin/mycelium`. **Do not require a live VM in `go test`.** Install SLO is human evidence.

No `scripts/install.sh` this phase. No PHASE-05 Actions workflow. GitHub Release upload is human evidence.

## 7. Commands (what exists / what does not)

| Command | This phase |
| --- | --- |
| `version` / `help` / `new` / `check` / `tier` / `publish` / `index` / `state` / `wake` / `status` | Exist on the pin. Do not reopen. `check` + `status` gain §5 / §8 deltas. |
| **`mycelium supersede`** | **Add.** §4. |
| `mycelium council` / `second-opinion` / `ladder` / `replicate` | Do **not** add. OQ-003 stands: only `council` is a pack. No portable council CLI. |
| `mycelium handoff` / `upgrade` / `migrate` / `release` / `install` | Do **not** add. Release is `scripts/release.sh`, not a CLI verb. |
| `mycelium state handed-off` | Still refuses. `handed-off` stays unreachable. |

Usage string in `internal/cli` gains one line:

```text
  mycelium supersede <OLD-ID> --by <NEW-ID> [--dir PATH]
```

Unknown-command teach for `council` / `handoff` stays "unknown command" (do not special-case).

## 8. Check updates

Items 1–22 stay; do not renumber. Add item 23. Extend item 8.

| Item | Delta | Binds |
| --- | --- | --- |
| 1 | Required manifest keys from the **instance** `program/contracts/manifest.md` table (G1). Unknown-key refuse still follows the instance contract (G3 check-fails). | Slice 3 |
| 6 | Resolve `superseded_by` and `supersedes` values as IDs. | Slice 2 |
| 8 | Log-op regex becomes `scaffold\|new\|tier\|publish\|check\|state\|wake\|supersede`. | Slice 2 |
| 9 | Unchanged. A leftover `supersede` journal is an interrupted operation. | already |
| **23** | **New.** Bidirectional IFF + one-to-one. If `status = "Superseded"`: `superseded_by` present, same namespace, resolves, and peer `supersedes` equals this ID. If `supersedes` is set: peer exists, peer `status = "Superseded"`, peer `superseded_by` equals this ID. At most one inbound `superseded_by` per NEW. | Slice 2 |

### Lift timing (conformance.md)

| Slice | Check / status behavior |
| --- | --- |
| 1 | Schema enum + optional keys land. Parsers exist. **No** CLI. **No** item 23 bind. Item 8 regex **not** yet changed (a hand-written `supersede` log line would still fail check — do not write one in Slice 1 fixtures). |
| 2 | Command + item 8 + item 23 + item 6 link bind. Happy DEC pair + refuse table. |
| 3 | Item 1 G1 rule. `status` / `status --all` tolerance. G0–G3. |
| 4 | No new check item. CHANGELOG + release script + checksum tests. |
| 5 | MS-501 matrix harness in `internal/clitest` runs every §14 row. |

### What check must not do (additions)

| Temptation | Verdict |
| --- | --- |
| Change `state` because an artifact was superseded | **No.** |
| Require a new H2 on OLD or NEW | **No.** |
| Grade OLD/NEW prose | **No.** DEC-005. |
| Fail G3 as a *status* fixture | **No.** G3 is status-only. |
| Fail G1 because the *binary* contract requires `github_repo` | **No.** Instance contract wins. |
| Add `handoff` / `upgrade` / `council` to the log-op regex | **No.** |
| Call network / `gh` / read `GH_TOKEN` | **No.** |
| Treat Install SLO or a GitHub Release as a check | **No.** |

Teaching errors stay four lines, cap 20. Success stdout for `check` is unchanged.

## 9. Skills / AGENTS.md / mycelium-cli deltas

No new skill. No new pack. No retrofit into existing instances.

| File | Delta |
| --- | --- |
| `program/skills/mycelium-cli/SKILL.md` | Document `mycelium supersede <OLD-ID> --by <NEW-ID> [--dir PATH]`. Name refuse cases (ineligible / OQ / already superseded / one-to-one / idea-state). State it does not change idea `state` and does not implement `handed-off`. Do not mention a `mycelium council` command as something to run. |
| `program/skeleton/AGENTS.md` | One short paragraph: artifact supersede is `mycelium supersede`; idea lifecycle stays `mycelium state` / `wake`. |
| `program/skills/portfolio/SKILL.md` | One line: `status --all` may print `partial: legacy-manifest` and still list other ideas. Do not rewrite the skill. |
| council / thinking / spark / wake skills | Do not touch. |

New scaffolds pick up the updated `mycelium-cli` + `AGENTS.md` via embed. `tier` / `index` do not retrofit.

## 10. Templates / schema deltas

`methodology_version` stays `2.0.0`. Re-run embed generate after edits.

### Sidecar schemas

| File | Delta |
| --- | --- |
| `program/templates/decision.schema.toml` | Enum already has `Superseded`. Declare optional keys `superseded_by`, `supersedes` (not in `required_front_matter`). |
| `program/templates/assumption.schema.toml` | `[enums.status].values` becomes `["Open", "Held", "Falsified", "Retired", "Superseded"]`. Same optional keys. |
| `program/templates/evidence.schema.toml` | Enum already has `Superseded`. Same optional keys. |
| `program/templates/spike.schema.toml` | Enum already has `Superseded`. Same optional keys. |
| `question.schema.toml` | **No** `status`. **No** `Superseded`. Do not reopen PHASE-03. |
| REC / REQ / RSK / FND / PHASE / MS / pack schemas | **No** `status` invent. No link keys. |

**Architect default — optional-key declaration:** if the schema DSL on the pin has no `optional_front_matter` list, extra front-matter keys are already allowed (required is a minimum). Still *name* `superseded_by` and `supersedes` in a comment or an `optional_front_matter` array if the parser already accepts that key; do not invent a new schema DSL. Check item 23 reads the keys from front matter regardless.

### Templates

`mycelium new decision|assumption|evidence|spike` does **not** emit `superseded_by` or `supersedes` (absent, not empty). Do not add a required H2. A one-line HTML comment in `decision.md` (`<!-- superseded_by / supersedes: set by mycelium supersede -->`) is allowed, not required.

Do not edit `CONTEXT.md`, `question.md`, or pack templates.

## 11. Vertical slices 0–5

PR-per-slice, sequential, rebase on main. Arvo merges Quality-green PRs. Engineering opens PRs. Engineering does NOT push to main. One live PR at a time.

Each PR title: `PHASE-05 Slice N: <done-bar noun>`. Each PR body links this brief, the slice done bar, and Linear **ROB-512**. No drive-by refactors. No v1 deletions. No PHASE-06 commands. No Actions job as a done bar. Do **not** add `.github/workflows/phase-05-*.yml`.

### Slice 0 — Commissioning (docs only)

This brief + `framework/phases/PHASE-05-acceptance.md` (rows = §14) + conformance lift-timing stub (may finish in Slices 1–3). No product code. No Go.

Done: files exist on a PR. Quality reads them against this brief.

### Slice 1 — Parsers / rules only (no CLI)

- `internal/supersede`: parse IDs, eligibility, refuse table as functions, cross-link mutation of front matter (pure: bytes in → bytes out).
- Schema deltas in §10. Embed generate.
- **No** `internal/cli` case. **No** item 8 / 23 bind. **No** log line written by a command.

Done (hermetic `go test`): DEC pair bytes gain `status = "Superseded"`, `superseded_by`, `supersedes`; OQ / PHASE / REC inputs refuse; ASM enum parses `Superseded`; idea-state tokens refuse. `mycelium supersede` is still "unknown command".

### Slice 2 — Command + check bind

- `internal/supersedecmd` + `internal/cli` verb.
- Protocol: preflight / lock / stage / commit order in §4. `allowReplace` extension.
- Item 8 regex + item 6 links + item 23 IFF.
- `indexmd` rewrite. Log line. `updated_date` only.

Done (hermetic `go test`): happy DEC-001 → DEC-002 (Appendix B) then `mycelium check` exit 0; every refuse row exit 1 and no file mutation; leftover `supersede` journal resumes; `state` unchanged; `handed-off` still unreachable.

### Slice 3 — Old-manifest fixtures + status tolerance

- `ParseTolerant`. `status` / `status --all` emit `partial: legacy-manifest` and continue.
- Item 1 G1 required-key rule.
- Golden G0–G3 (Appendix C) on disk under testdata.

Done: G0 check+status pass; G1 check+status pass; G2 `status --all` no crash + `partial: legacy-manifest`; G3 status pass + check fail. A mixed `--root` of G0+G1+G2+G3+one unreadable file still lists G0/G1/G3 and exits 0.

### Slice 4 — CHANGELOG + release script + checksum tests

- `CHANGELOG.md` with `## [Unreleased]` and `## [0.1.0] - 2026-08-15`.
- `scripts/release.sh` + Justfile recipe.
- `docs/install.md` with the one-liner.
- Hermetic: refuse (temp CHANGELOG missing heading) + SHA256SUMS match on built (or temp) files. File-exists assert for `docs/install.md`.

Done: `scripts/release.sh 0.1.0` against the fixture heading writes two binaries + `dist/SHA256SUMS` whose hexes match; `scripts/release.sh 9.9.9` refuses; `internal/version.Version` in the tree is still `0.1.0-dev`. No tag. No upload. No workflow file.

### Slice 5 — MS-501 matrix harness

**Architect default:** Slice 5 is the matrix harness in `internal/clitest` that runs **all** MS-501 rows (Appendix E). Even if Slices 2–4 already cover every row, keep Slice 5 as the single file that is the gate list.

Done: `go test ./...` runs the MS-501 matrix green. That **is** the gate. Install SLO is **not** the gate. A live GitHub Release is **not** the gate. An Actions log is **not** the gate.

## 12. Done / verified mapped onto MS-501

MS-501 is the hermetic phase gate. Blueprint MS-501 (do not expand): functional acceptance — `mycelium supersede` leaves bidirectional cross-links and a log entry; `check` and `status` pass golden old-instance fixtures; releases ship checksummed binaries with a `CHANGELOG.md` entry. Install SLO is a **separate human-evidence clause**, not the merge gate (Arvo locked this).

### MS-501 expected (authoritative; recipe in Appendix E)

| Row | Setup | Expect |
| --- | --- | --- |
| MS-501-SUP | Fixture instance; DEC-001 + DEC-002; `mycelium supersede DEC-001 --by DEC-002` | OLD `status=Superseded` + `superseded_by=DEC-002`; NEW `supersedes=DEC-001`; log line `supersede` / `DEC-001 -> DEC-002`; `mycelium check` exit 0; `state` unchanged |
| MS-501-G0 | Current `mycelium.toml` | `check` + `status` exit 0 |
| MS-501-G1 | Missing `github_repo` + frozen instance manifest contract | `check` + `status` exit 0 |
| MS-501-G2 | `research-program.toml`, no `mycelium.toml` | `status --all` no crash; `partial: legacy-manifest`; not listed as an instance |
| MS-501-G3 | Extra unknown top-level key | `status` exit 0 (readable); `check` exit 1. **Status-only golden.** |
| MS-501-REL | `CHANGELOG.md` has `## [0.1.0]`; `scripts/release.sh 0.1.0` | `dist/mycelium-linux-amd64`, `dist/mycelium-darwin-arm64`, `dist/SHA256SUMS` hexes match; refuse when heading missing |

`gh` never invoked. No `GH_TOKEN`. No Actions job. No live VM. No real tag. Tests never touch the real home directory except as already allowed by PHASE-02 status fixtures (temp `UserHomeDir`).

### Slice → MS-501 map

| Slice | MS-501 clause |
| --- | --- |
| 0 | commissioning; not a runtime clause |
| 1 | parsers / schema; not yet the command |
| 2 | MS-501-SUP (command + check) |
| 3 | MS-501-G0 / G1 / G2 / G3 |
| 4 | MS-501-REL |
| 5 | the matrix in `go test ./...` **is** the gate |

PHASE-05 is accepted when MS-501 is green in `go test ./...` on main **and** Quality has thermos'd the table on their computer (local guide). Arvo accepts the phase. Engineering does not self-accept. Install SLO / real tag / GitHub Release may be attached as human evidence; they do not replace the gate.

## 13. Automated test plan

Engineering MUST write these tests. Quality thermos against this list on **their computer**. That local `go test ./...` is the verify bar, not an Actions log.

### House test shape (locked 2026-08-15)

1. Fast unit tests on parsers/state/rules (most of the 85% coverage floor).
2. Hermetic CLI tests: real binary, temp instance, fake clock/remotes, including fail cases.
3. Thin MS-501 fixture matrix.
4. Live GitHub, live Cursor, live network, browser e2e, and CI are NOT gates. Human evidence only.
5. If UI: at most 1–2 happy-path browser tests (this project has no UI — say so). The rest would be component tests.

This project has **no UI**. Item 5 is: zero browser tests. Do not put live e2e or CI-as-proof in a done bar.

### Unit (no network, no gh, no home directory)

| Area | Cases |
| --- | --- |
| eligibility | DEC/ASM/EVD/SPK yes; OQ/PHASE/REC/REQ/RSK/FND/MS/CMP no |
| refuse table | every §4 row as a function |
| cross-link bytes | Appendix B before → after |
| ParseTolerant | missing `github_repo`; unknown key; not TOML |
| required-keys from instance `manifest.md` | G1 omits `github_repo`; current table includes it |
| CHANGELOG heading | `## [0.1.0]` ok; `## [0.1.0] - 2026-08-15` ok; missing refuse |
| SHA256SUMS parse / match | two-line fixture |

### Hermetic CLI (built binary, temp dirs)

| Case | Expect |
| --- | --- |
| `mycelium supersede DEC-001 --by DEC-002` | Appendix B after; check 0; log line present |
| each refuse row | exit 1; files unchanged |
| `mycelium supersede spark --by DEC-001` | exit 1; idea-state teach |
| `mycelium supersede OQ-001 --by OQ-002` | exit 1; open a new OQ |
| chain: supersede current tip | DEC-002 → DEC-003 ok; rewriting DEC-001 pair refuses |
| `state` after supersede | unchanged |
| G0 / G1 / G2 / G3 | §5 table |
| mixed `--root` | exit 0; `partial: legacy-manifest` present; G0 listed |
| `scripts/release.sh 0.1.0` | two binaries + SHA256SUMS match |
| `scripts/release.sh 9.9.9` | refuse; no `dist/` writes |
| `docs/install.md` | exists; names the one-liner |
| no `gh` | `gh` never invoked |
| MS-501 | Appendix E; `go test ./...` |

Do not require live GitHub, a live VM, a real tag, or an Actions job for `go test ./...`.

## 14. Acceptance / MS-501 fixture matrix

Slice 0 lands `framework/phases/PHASE-05-acceptance.md` with these rows. Later slices land the packages and testdata in Appendix F. No workflow file. No DEC-015 file. No `mycelium council` package. No `program/migrations/`.

Each row: id, check, evidence, owner (Engineering | Arvo). **CI is not an owner.** Robert waived CI.

| id | check | evidence |
| --- | --- | --- |
| A-S0 | Slice 0 files exist and match this brief | PR diff; Quality read |
| A-S1 | parsers / eligibility / schema deltas; no CLI verb yet | hermetic `go test` |
| A-S2 | `mycelium supersede` + item 8 + item 23 + refuse table | hermetic `go test` |
| A-S3 | G0–G3 + `partial: legacy-manifest` + no scan abort | hermetic `go test` |
| A-S4 | CHANGELOG heading + release refuse/checksum + install doc | hermetic `go test` + file read |
| A-S5 | MS-501 matrix green | `go test ./...` on Quality's computer |
| MS-501-SUP | bidirectional cross-links + log entry | A-S5 / A-S2 |
| MS-501-G0 | current manifest check+status | A-S5 / A-S3 |
| MS-501-G1 | missing `github_repo` + frozen contract | A-S5 / A-S3 |
| MS-501-G2 | `research-program.toml` only; no crash | A-S5 / A-S3 |
| MS-501-G3 | unknown key: status-only golden | A-S5 / A-S3 |
| MS-501-REL | checksummed binaries + CHANGELOG heading | A-S5 / A-S4 |

No Install-SLO-gate row. No GitHub-Release-gate row. No Actions-job row. Quality should refuse a PR that adds an Actions job as the MS-501 gate. Quality should **not** refuse a missing PHASE-05 workflow.

## 15. Decided / Architect defaults / not decided

No TBD. Open items are decided inline and labeled **Architect default**. Do not bikeshed them in a code PR. Write a DEC to change one. No DEC-015 is required for these.

Index of defaults that are easy to miss:

- **New verb:** `mycelium supersede <OLD-ID> --by <NEW-ID> [--dir PATH]`. Artifact-level. Does not change `state`. Does not implement `handed-off`. Same namespace. Both IDs exist. No deletion. One-to-one this phase. Chain = supersede the current record.
- **Eligible:** DEC, ASM (add `Superseded`), EVD, SPK. OQ refuse. REC/REQ/RSK/FND not applicable (no `status` field; do not invent one). PHASE/MS/CMP/RPT/RCL no.
- **Cross-links:** OLD `status=Superseded` + `superseded_by`; NEW `supersedes`. No new H2. Commit order: OLD, NEW, `index.md`, `log.md`, `mycelium.toml` last (`updated_date` only).
- **DEC-011:** `ParseTolerant` for status; `partial: legacy-manifest`; G0–G3 as §5 (G3 status-only). Check required keys from instance `manifest.md`. No migrations.
- **Release:** `CHANGELOG.md` + `scripts/release.sh` + `dist/SHA256SUMS` + `docs/install.md`. linux-amd64 + darwin-arm64 only. Version string stays `0.1.0-dev`. No goreleaser. No Windows. No Actions. Install SLO and the real tag / GitHub Release are human evidence.
- **Slice 5** is the `internal/clitest` harness that runs all MS-501 rows. Gate = files + hermetic `go test ./...` + Quality's local guide. Not Actions.
- **No new pack.** No `mycelium council`. OQ-003 stands. `methodology_version` stays `2.0.0`. Do not touch `internal/slug` (DEC-014). Do not commission PHASE-06.

**Not decided:** none — defaults above.

## 16. Risks, rollback, Quality should refuse

### Risks

| Risk | Mitigation |
| --- | --- |
| Supersede mutates idea `state` or unlocks `handed-off` | §4. Quality should refuse. |
| Migrations sneak back (`program/migrations/`, `just upgrade`) | DEC-011. Quality should refuse. |
| `status --all` aborts on one bad instance | G2/G3 + mixed-root fixture. Quality should refuse. |
| G3 treated as a check-pass | §5 table. Quality should refuse. |
| Actions job added as the MS-501 gate | Robert waived CI. Quality should refuse. Absence of a PHASE-05 workflow is correct. |
| Install SLO or live GitHub Release treated as the gate | §12. Quality should refuse. |
| `allowReplace` opened for `new` | §4: supersede targets only. |
| New pack / `mycelium council` | OQ-003. Quality should refuse. |
| Growing `latinFold` or adding `x/text` | DEC-014. Do not touch `internal/slug`. |
| CLI git-commits instance work product | DEC-010. Quality should refuse. |
| Emitting `framework/` or converting master | Absence tests stay. |

### Rollback

Revert the offending PR on master. Do not `git push --force` to main. Floor is `657a14653da1c41fd4c0590a5b0aa625eaa9adde`. Do not delete Justfile / v1 scripts as a "cleanup" rollback.

### Quality should refuse

Refuse to approve if:

- an Actions job is added as the MS-501 gate, or `.github/workflows/phase-05-*.yml` is added, or `phase-01-hermetic.yml` is extended as a phase gate (Quality should **not** refuse a missing PHASE-05 workflow — absence is correct)
- Install SLO, a live GitHub Release, a real `v0.1.0` tag, live GitHub, live Cursor, live network, or CI-as-proof is required as the done bar
- `program/migrations/`, `just upgrade`, or `applied_migrations` appears
- `handed-off` succeeds, or check stops failing stored `handed-off`, or supersede changes `state`
- a `mycelium council` / `handoff` / `upgrade` verb or a new pack appears
- G3 is asserted as a `check` pass, or `status --all` aborts the scan on one bad instance
- `framework/` is emitted; CLI git-commits instance work product
- `latinFold` grows or NFKD is implemented (DEC-014)
- cobra / viper / yaml / testify / go-github / `golang.org/x/text` appears
- Justfile / v1 scripts deleted, `just init` run on master, or `research-program.toml` renamed
- DEC-012 / DEC-013 / DEC-014 reopened; DEC-015 recorded
- PHASE-06 commands appear; Windows binary / goreleaser ships
- hermetic tests call network or real `gh`
- PR pushed straight to main

### Quality local guide (the verify bar)

Quality runs `go test ./...` on their computer and thermos the MS-501 table in §14 / Appendix E. That local run is the verify bar, not an Actions log. Do not ask Engineering for a workflow badge.

## 17. Execution order + Linear ROB-512

Same order as §11 (slices 0→5). PR-per-slice, sequential, rebase on main. One live PR at a time. Slice 2 is the command bind — do not combine it with 3–5. Slice 5 must be green in `go test ./...` on its PR (not Actions).

Title: `PHASE-05 Slice N: <done-bar noun>`. Body links this brief, the slice done bar, and parent **ROB-512**. Arvo cuts ROB-512 child tickets after this brief is stamped. No drive-by refactors, v1 deletions, PHASE-01–04 leftover work, or PHASE-06 commands. Engineering opens PRs; Arvo merges Quality-green PRs; Engineering does NOT push to main.

Cursor cloud env name is exactly `robertguss/mycelium`. Go 1.26. `CGO_ENABLED=0`.

## 18. Handoff

### What Engineering starts with

This file. Only this file. Start from `https://github.com/robertguss/mycelium` at `657a14653da1c41fd4c0590a5b0aa625eaa9adde` (PHASE-01–04 accepted). Read `framework/blueprint.md` and DEC-001–014 for authority, not for a second plan. Execute Slice 0 first: land this brief at `framework/phases/PHASE-05-implementation-brief.md` plus `framework/phases/PHASE-05-acceptance.md`. **Architect does not open the docs PR. Engineering lands it.** Do not implement from a later SHA unless Arvo re-pins. Cloud env `robertguss/mycelium`. Go 1.26. `CGO_ENABLED=0`. Linear parent: **ROB-512**.

### What Engineering must not do

See §16. Do not open a design debate in the PR. Do not write a second brief. Do not write DEC-015. Do not commission PHASE-06. Do not add an Actions job as the MS-501 gate. Do not start from a later SHA and "fast-forward" this brief. Do not invent `mycelium council`. Do not put Install SLO or a live GitHub Release in a done bar.

### What Quality reads

This brief, the acceptance matrix, the conformance delta, and the PR diff. Thermos: §13 house test shape + tests exist; §16 refuse list is clean; MS-501 hermetic `go test ./...` on Quality's computer; no Actions job as the MS-501 gate; no migrations; no new pack.

Quality should refuse a PR that adds an Actions job as the MS-501 gate; Quality should **not** refuse a missing PHASE-05 workflow.

### What Arvo does

Cuts ROB-512 child tickets after this brief is stamped. Merges Quality-green PRs. Accepts PHASE-05 when MS-501 is green on main and Quality has thermos'd locally. May attach Install SLO / a real tag / a GitHub Release as human evidence; must not treat them as the gate. Does not re-pin without writing the new SHA. Does not record DEC-015. Does not commission PHASE-06.

## Appendix A — No new DEC

No DEC-015. PHASE-05 implements **DEC-011**'s risk guard (already Accepted 2026-08-14) plus the blueprint PHASE-05 commands. PHASE-05 does not reopen DEC-012, DEC-013, or DEC-014. Remaining choices are Architect defaults in §15. Engineering lands **zero** new files under `framework/decisions/`.

If a later discovery contradicts a locked item, stop and write a DEC; do not silently patch this brief in a code PR.

## Appendix B — Supersede before/after example (DEC-001 → DEC-002)

Fixture instance (not master `framework/decisions/`). Tests must **not** grade the words — only front matter, log line, and `mycelium check` exit 0.

Before (`decisions/DEC-001-use-sqlite.md` and `decisions/DEC-002-use-sqlite-with-wal.md`), both `status = "Accepted"`, no `superseded_by`, no `supersedes`. Required DEC H2s present with `none` / short fill (DEC-005: do not grade).

```text
mycelium new idea "Supersede Fixture" --offline --dir PATH
mycelium new decision "Use SQLite" --dir PATH
mycelium new decision "Use SQLite with WAL" --dir PATH
# EDIT both to status = "Accepted" if the template emits Proposed
mycelium supersede DEC-001 --by DEC-002 --dir PATH
mycelium check --dir PATH
# exit 0
```

After, OLD front matter contains:

```text
id = "DEC-001"
status = "Superseded"
superseded_by = "DEC-002"
```

NEW front matter contains:

```text
id = "DEC-002"
status = "Accepted"
supersedes = "DEC-001"
```

`log.md` contains a parseable line whose op is `supersede` and whose note is `DEC-001 -> DEC-002`. `mycelium.toml` `state` is still `spark` (or whatever it was). `updated_date` bumped. No file deleted.

Negatives in the same test file: `--by DEC-001` on already-superseded OLD → exit 1; `--by DEC-003` when DEC-002 already has `supersedes` → exit 1; `mycelium supersede OQ-001 --by OQ-002` → exit 1; `mycelium supersede spark --by DEC-002` → exit 1.

## Appendix C — Golden old-manifest snippets

Hermetic testdata. Paths are Architect defaults under `internal/clitest/testdata/legacy/`.

### G0 — current control

A normal `--offline` scaffold is enough (or a frozen copy of current `mycelium.toml` with every required key, including `github_repo = ""`). `check` + `status` exit 0.

### G1 — missing `github_repo` + frozen instance contract

`mycelium.toml` (no `github_repo` key):

```text
schema_version = 1
idea_name = "Legacy One"
slug = "legacy-one"
state = "spark"
tier = "focused"
methodology_version = "2.0.0"
generated_by_cli_version = "0.1.0-dev"
created_date = "2026-08-01"
updated_date = "2026-08-01"
revisit = ""
```

Frozen `program/contracts/manifest.md` `## Required fields` table lists `schema_version`, `idea_name`, `slug`, `state`, `tier`, `methodology_version`, `generated_by_cli_version`, `created_date`, `updated_date`, `revisit` — **not** `github_repo`. Other instance files are a minimal spark (so check has something legal to walk). `check` + `status` exit 0.

### G2 — v1 master-shaped directory

```text
research-program.toml
# any minimal toml; no mycelium.toml
```

`status --all --offline --root PATH` (PATH's child is this directory) exits 0, prints `partial: legacy-manifest (<path>: research-program.toml without mycelium.toml)`, and does not list an idea row for it.

### G3 — unknown top-level key (status-only)

Current-shaped `mycelium.toml` plus:

```text
legacy_note = "append-only"
```

`status` exit 0 (listed). `check` exit 1 (instance contract refuses unknown keys). **Not a check-pass golden.**

## Appendix D — CHANGELOG + SHA256SUMS example

`CHANGELOG.md` (repo root; Keep a Changelog):

```text
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-08-15

### Added

- PHASE-01 through PHASE-04 CLI surface.
- `mycelium supersede` (PHASE-05).
```

`dist/SHA256SUMS` after `scripts/release.sh 0.1.0` (hexes are examples; tests assert match, not these bytes):

```text
<sha256>  mycelium-linux-amd64
<sha256>  mycelium-darwin-arm64
```

Refuse fixture: a temp tree whose CHANGELOG has only `## [Unreleased]` → `scripts/release.sh 0.1.0` exits non-zero and does not write `dist/`.

## Appendix E — MS-501 fixture recipe

Hermetic. No network. No `gh`. No `GH_TOKEN`. `go test ./...` only. Do **not** add an Actions job or `.github/workflows/phase-05-*.yml`. Do **not** extend `phase-01-hermetic.yml` as a phase gate. Temp dirs; freshly built binary; fake clock (`MYCELIUM_NOW=2026-08-15T00:00:00Z` or `internal/clock`). Do not invoke a live VM.

**Architect default:** matrix test lives at `internal/clitest/ms501_hermetic_test.go` (execs the binary). Parser tests live in `internal/supersede`. Status tolerance tests may live in `internal/statuscmd` and/or the matrix file.

### MS-501-SUP

```text
mycelium new idea "Supersede Fixture" --offline --dir PATH/sup
mycelium new decision "Use SQLite" --dir PATH/sup
mycelium new decision "Use SQLite with WAL" --dir PATH/sup
# set both status = "Accepted" if needed
mycelium supersede DEC-001 --by DEC-002 --dir PATH/sup
mycelium check --dir PATH/sup
# exit 0
# assert superseded_by / supersedes / log note / state unchanged
```

### MS-501-G0 / G1 / G2 / G3

Build or copy Appendix C trees under a temp `--root`. Run `mycelium check --dir <g0>` and `--dir <g1>` (exit 0); `mycelium check --dir <g3>` (exit 1); `mycelium status --all --offline --root ROOT` (exit 0; G0/G1/G3 listed; G2 not listed; stdout or stderr contains `partial: legacy-manifest`).

### MS-501-REL

In a temp copy of the repo tree (or a temp dir with `CHANGELOG.md` + `scripts/release.sh` + enough `go.mod` to build), run `scripts/release.sh 0.1.0`. Assert two binaries + `SHA256SUMS` match. Run `scripts/release.sh 9.9.9` and assert refuse. Assert `docs/install.md` exists in the *source* tree and names the one-liner. Do not tag. Do not upload.

`gh` invoked → FAIL the test. No seven-real-day dogfood. No live VM. No Actions job.

## Appendix F — File tree additions / DO NOT ADD

### Master (additions on top of the PHASE-01–04 tree; v1 files retained)

```text
internal/supersede/                 # Slice 1 parsers / rules
internal/supersedecmd/              # Slice 2 command
internal/clitest/ms501_hermetic_test.go
internal/clitest/testdata/legacy/{g0,g1,g2,g3}/
CHANGELOG.md
scripts/release.sh
docs/install.md
program/templates/{decision,assumption,evidence,spike}.schema.toml
program/contracts/conformance.md    # item 8 + item 23 + lift timing
program/skills/mycelium-cli/SKILL.md
program/skeleton/AGENTS.md
program/skills/portfolio/SKILL.md   # one line
framework/phases/PHASE-05-implementation-brief.md
framework/phases/PHASE-05-acceptance.md
internal/embed/program/             # regenerate after program/ edits
Justfile                            # add `release version:` only
```

### DO NOT ADD

```text
.github/workflows/phase-05-*.yml
.github/workflows/phase-05-hermetic.yml
.github/workflows/phase-05-ms501.yml
framework/decisions/DEC-015-*.md
program/migrations/
program/packs/{supersede,handoff,registry.toml}
a mycelium council / handoff / upgrade / release command package
scripts/install.sh
dist/mycelium-windows-amd64
.goreleaser.yml
internal/slug changes
```

Do **not** add a PHASE-05 workflow. Do **not** extend `phase-01-hermetic.yml` as a phase gate. Quality should refuse a PR that adds an Actions job as the MS-501 gate. Quality should **not** refuse a missing PHASE-05 workflow (absence is correct). Do not delete Justfile / v1 scripts / `research-program.toml` / PHASE-01 workflows. Do not add a `mycelium council` command. Do not add Windows. Do not add migrations.

### Emitted instance (spark / focused, local-only, PHASE-05 scaffold)

Unchanged from PHASE-04 plus updated `mycelium-cli` / `AGENTS.md` text. After `mycelium supersede` (Slice 2+): two artifact files rewritten in place; `log.md` appended; `index.md` rewritten; `mycelium.toml` `updated_date` only.

Absent: `framework/`, `cmd/`, `internal/`, `go.mod`, `Justfile`, `scripts/`, `research-program.toml`, a `council` CLI verb, `program/migrations/`.

Unexported helpers may live next to their tests. No `pkg/`. Do not touch `internal/slug`.

End of PHASE-05 implementation brief. Engineering executes from this file only.
