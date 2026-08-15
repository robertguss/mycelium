# PHASE-06 Implementation Brief — Handoff

- **Status:** Binding
- **Date:** 2026-08-15
- **Audience:** Engineering (pstack / poteto-mode)
- **Authority:** `framework/blueprint.md` (Accepted 2026-08-14). DEC-001 through DEC-014 stand. OQ-004 is settled in this brief (Appendix A). Do **not** record DEC-015. Do **not** reopen DEC-012 / DEC-013 / DEC-014.
- **Repo:** https://github.com/robertguss/mycelium
- **Pin:** Engineering starts from `main` @ `a486cc9aa04a9ef22e7d0e564df3f2ebe692b1bb` (PHASE-01–05 accepted). Do not implement from a later SHA unless Arvo re-pins in writing.
- **Product:** single-binary Go CLI `mycelium`. Master builds the CLI. `program/` is `go:embed`'d and emitted on scaffold. PHASE-06 adds the implementation packet (`handoff/PACKET.md`), **`mycelium handoff`**, and DEC-006 reachability of **`handed-off`**. **No new pack.** **No `mycelium council`.** CLI version stays `0.1.0-dev`.
- **Phase gate:** MS-601 via files + hermetic `go test ./...` + Quality's local guide. GitHub Actions is **not** a gate (Robert waived CI). Do **not** add `.github/workflows/phase-06-*.yml`. Do **not** extend `phase-01-hermetic.yml` as a phase gate. Quality should refuse a PR that adds an Actions job as the MS-601 gate. Quality should **not** refuse a missing PHASE-06 workflow (absence is correct).
- **Not the gate:** A live isolated-agent run of the packet is **human evidence** (Quality's `## Manual test guide`, optional). A real-project handoff is **dogfood**, not the gate. Do not put either in a done bar.
- **How to use this file:** Engineering executes from THIS FILE ONLY. No "see chat". No TBD. Open items are decided here and labeled **Architect default**. Not-decided is empty (§15).
- **Coverage floor:** **85%** statement/line on product packages (exclude generated / vendor / data-only fixtures). **New packages must themselves be ≥85%.** Stamp: **85**.

Headings: §§1–18 then Appendices A–F (no new DEC; PACKET.md example; pstack/poteto mapping; canonical Add fixture; MS-601 recipe; file tree / DO NOT ADD).

Cloud env name is exactly `robertguss/mycelium`. Go 1.26 at `/usr/local/go` or `/usr/local/bin/go`. `CGO_ENABLED=0`. stdlib + `github.com/pelletier/go-toml/v2` only. No cobra / viper / yaml / testify / go-github. Linear parent for slice tickets: **ROB-511** (Arvo cuts tickets after this brief is stamped).

## 1. Scope / out of scope

Tonight is PHASE-06 only. Do not reimplement PHASE-01–05. Do not convert master. Do not commission a later phase.

### In scope

- New packet contract: `program/contracts/handoff-packet.md`. This is the PHASE-06 packet. Do **not** reuse `program/contracts/handoffs.md` as the packet contract.
- Five-line banner at the top of `program/contracts/handoffs.md`: that file is v1 ADRP session-attachment manifests; the implementation packet is `handoff-packet.md`.
- Instance packet tree `handoff/` (never `docs/handoffs/`, never `framework/`): `PACKET.md`, `decisions/`, `glossary.md`, `questions/`, `evidence/`, `playbooks/`, `acceptance/`.
- New verb: `mycelium handoff [--dir PATH]`. Legal **only** from `clarified`. Writes the packet **then** sets `state = handed-off`. Silent handoff forbidden.
- `mycelium state handed-off` becomes a legal argv target **IFF** `handoff/PACKET.md` already exists and passes structure. Missing packet → refuse, teaching error names `mycelium handoff`.
- DEC-006 machine unchanged. `handed-off` becomes reachable **this phase**. `handed-off` is terminal except `→ archived`. Conformance: no `handed-off` without a packet.
- Check lift in the same slice as the command (Slice 2): stored `handed-off` is LEGAL iff packet structure passes; stored `handed-off` without packet still FAIL; `clarified` without packet still LEGAL. Lift the PHASE-02/05 "handed-off always FAIL" rule.
- Log op `handoff`. Journal op adds `handoff`. Protocol: preflight / lock / stage / commit / rollback. Commit order: `handoff/**`, `index.md`, `log.md`, `mycelium.toml` last.
- pstack/poteto bridge (docs, not a runtime dependency) + `AGENTS.md` **Implementation systems** section + `mycelium-cli` skill.
- Canonical fixture (bounded `Add(a, b int) int` target) + executable `handoff/acceptance/` tests + golden implementation in testdata.
- MS-601 hermetic fixture matrix in `go test ./...` (`internal/clitest`). Coverage **≥85%** on product packages; new packages themselves **≥85%**.
- Commissioning artifacts: this brief, `handoff-packet.md`, lifecycle rewrite, `handoffs.md` banner, acceptance stub.

### Out of scope (one refuse list; details in §16)

- Isolated-agent run or real-project handoff as a merge gate (dogfood / optional human evidence only).
- A live agent inside `go test`. Live GitHub, live Cursor, live network, browser e2e.
- Migrations: no `program/migrations/`, no `just upgrade`, no `applied_migrations` (DEC-011).
- New pack. `mycelium council`. Portable council CLI. Reopening OQ-003.
- A PHASE-06 Actions workflow. Extending `phase-01-hermetic.yml` as a phase gate. `.github/workflows/phase-06-*.yml`.
- Reusing `program/contracts/handoffs.md` as the packet contract. Emitting the packet under `docs/handoffs/` or `framework/`.
- Reopening DEC-012 / DEC-013 / DEC-014. Recording DEC-015.
- Emitting `framework/`. CLI `git add` / `git commit` of instance work product.
- Converting master (`research-program.toml`, `just init`, deleting Justfile / v1 scripts).
- Growing `latinFold`, NFKD, `golang.org/x/text` (DEC-014).
- Changing CLI version off `0.1.0-dev`. Changing `methodology_version` off `2.0.0`.
- Implementing MS-101(b). Commissioning a `GH_TOKEN` job.

### Master vs instance (unchanged)

Master remains an ADRP v1 instance. Do not convert `research-program.toml` to `mycelium.toml`. Do not run `just init` on master. `framework/` is master-only and is NEVER emitted. Justfile and v1 `scripts/*.py` stay. Runtime detects idea instances by `mycelium.toml`.

## 2. Authority and do-not-reopen

### Governing documents

| Document | Rule |
| --- | --- |
| `framework/blueprint.md` (Accepted 2026-08-14) | Do not rewrite vision. DEC-001–014 stand. PHASE-06 ceiling is blueprint lines ~440–446 (packet contract + generator; pstack/poteto bridge; AGENTS.md implementation-systems; MS-601). Packet object ~168–170. OQ-004 ~465, **settled here**. DEC-006 machine ~123–143. |
| **DEC-006** | spark → exploring ⇄ simmering → clarified → handed-off; any → archived. Do **not** change the machine. `handed-off` becomes reachable this phase. Conformance: no handed-off without a packet. |
| DEC-005 | Checks validate containers, never contents. Packet check validates H2s, copies, front matter, and in-packet link resolution — not Framing prose. |
| DEC-010 | CLI never git-commits instance work product. |
| DEC-011 | No migrations. Instance files are truth. |
| DEC-012 / DEC-013 / DEC-014 | Do not reopen (`mycelium.toml`; refuse out-of-range; `latinFold` only, no NFKD, no `x/text`, do not grow the map). |
| This brief | Binding 2026-08-15. PHASE-06 only. Architect defaults are binding. OQ-004 settled here. No DEC-015. |

### Process override (unchanged)

Blueprint "humans-own-git" is overridden for the *master* repo's engineering process: Arvo merges Quality-green PRs and accepts the phase. Engineering opens PRs. Engineering does NOT push to main. The CLI still never git-commits *instance* work product.

### Do not reopen

Do not reopen the product shape, language, dependency floor, state vocabulary, manifest filename, refuse-vs-warn range rule, no-commit rule, instance-files-are-truth rule, slugify/DEC-014, publish, MS-101(b), PHASE-03 sparring, PHASE-04 ladder / OQ-003 (only `council` is a pack; no portable council CLI), or PHASE-05 supersede / release / DEC-011 tolerance. If a later discovery contradicts a locked item, stop and write a DEC; do not silently patch this brief in a code PR.

Do not reopen DEC-012, DEC-013, or DEC-014. Do **not** record DEC-015. Do not treat v1 `handoffs.md` as the packet.

## 3. Floor on main (do not reimplement)

Pin: `a486cc9aa04a9ef22e7d0e564df3f2ebe692b1bb`. Treat this SHA as the floor. Reuse packages. Do not rewrite working PHASE-01–05 commands.

### Already shipped (do not rebuild)

Reuse: `cmd/mycelium`, `internal/{cli,version,embed,clock,execrun,metadata,idpath,manifest,schema,slug,logfmt,teach,lock,journal,op,scaffold,generate,check,tiercmd,publish,clitest,revisit,lifecycle,indexmd,wakebrief,statecmd,statuscmd,sparring,pack,ladder,supersede,supersedecmd}`.

| Touch | Fate |
| --- | --- |
| `internal/cli` | **Add one verb:** `handoff`. Existing verbs stay: `version`, `help`, `new`, `check`, `tier`, `publish`, `index`, `state`, `wake`, `status`, `supersede`. Do not add `council` / `upgrade` / `migrate`. |
| `internal/handoff` | **New (Slice 1).** Packet parser / schema / structure checker. No CLI. No state lift. |
| `internal/handoffcmd` | **New (Slice 2).** Command + protocol bind. |
| `internal/lifecycle` / `statecmd` | **Extend (Slice 2).** `handed-off` becomes a legal argv target IFF packet structure passes. `handed-off` → anything except `archived` refuses. |
| `internal/check` | **Extend (Slice 2).** Lift "handed-off always FAIL". Item 8 regex gains `handoff`. New item 24 (packet IFF). Do not rewrite the package. |
| `internal/op` | **Extend** `allowReplace` so a `handoff` journal may replace paths under `handoff/` plus the regenerable set (`log.md`, `mycelium.toml`, `index.md`). Do not open general overwrite for `new`. |
| `internal/journal` | Op string `handoff` is legal. `Title` = `clarified -> handed-off`. No new journal schema version. |
| `internal/logfmt` | Log-op regex gains `handoff`. |
| `internal/indexmd` | Reuse. `handoff` rewrites `index.md` the same way `state` does. **No new required H2.** |
| `internal/slug` | Do not touch (DEC-014). |
| `internal/version` | Stay `"0.1.0-dev"`. `methodology_version` stays `2.0.0`. |
| `internal/embed` | Re-run `go generate` after `program/` edits. |
| `program/contracts/conformance.md` | Item 8 regex + item 24 + lift timing. Items 1–23 stay; do not renumber. |
| `program/contracts/lifecycle.md` | Rewrite stored-`handed-off` rule (legal iff packet). Teaching error now names `mycelium handoff`. |
| `program/contracts/handoffs.md` | **Banner only.** Do not rewrite the v1 body. Do not treat it as the packet. |
| `program/contracts/handoff-packet.md` | **New.** The packet contract. |
| `program/skills/{spark,wake,portfolio,thinking}` + council pack | Do not rewrite. Update `mycelium-cli` + `AGENTS.md` only. |
| `phase-01-hermetic.yml` / `phase-01-github.yml` | Leave alone. Do **not** add a PHASE-06 workflow. Actions is not a gate. |
| Justfile / v1 `scripts/*.py` / `research-program.toml` / `CHANGELOG.md` / `scripts/release.sh` | Keep. Do not delete. Do not `just init`. |
| `framework/` | Master-only. NEVER emitted. |

### Pin facts (do not "discover" otherwise)

- CLI verbs on the pin: `version` / `help` / `new` / `check` / `tier` / `publish` / `index` / `state` / `wake` / `status` / `supersede`. **No `handoff` yet.**
- `internal/version.Version = "0.1.0-dev"`.
- Lifecycle still refuses `clarified → handed-off`. `mycelium state handed-off` is not a legal argv target.
- Stored `handed-off` still **FAIL**s check. Teaching error still says the packet command is not shipped.
- Log-op regex on the pin: `scaffold\|new\|tier\|publish\|check\|state\|wake\|supersede`.
- `op.allowReplace` on the pin allows `log.md`, `mycelium.toml`, `index.md`, `briefs/*.md`, plus the two supersede artifact paths. `handoff/**` is not replaceable until Slice 2.
- `program/contracts/handoffs.md` on the pin is the v1 session-attachment manifest. It is **not** the packet. No `handoff-packet.md` yet. No `handoff/` tree on a scaffold.
- Master uses `research-program.toml`. Runtime detects instances by `mycelium.toml`.

### What must not be broken

`just check` on master; hermetic `go test ./...`; no `framework/` emit; no master conversion. PHASE-01–05 fixtures stay green. Spark with zero questions stays green. Instances without the council pack stay green. `mycelium supersede` stays green. `clarified` without a packet stays LEGAL. Until Slice 2 binds, stored `handed-off` still FAILs and `state handed-off` still refuses.

If a PHASE-06 PR is bad: revert that PR. Floor is the pin SHA.

## 4. Packet format (OQ-004) + v1 `handoffs.md` distinction

OQ-004 (handoff packet format + pstack playbook mapping) is **settled in this brief**. No TBD. No DEC-015.

### v1 file is not the packet

`program/contracts/handoffs.md` is the v1 ADRP **session-attachment manifest** contract (`/tmp/handoffs-v1.md` on this commission). It is about what a fresh *thinking* session attaches. It is **not** `handoff/PACKET.md`.

**Architect default — five-line banner** (exact, top of `handoffs.md`, before the existing H1 or replacing the H1 with this block + the old H1):

```text
# Context Handoff and Attachment Manifests

> **v1 session-attachment manifests — not the PHASE-06 packet.**
> This file is the ADRP v1 attachment-manifest contract.
> The implementation packet is `program/contracts/handoff-packet.md` (`handoff/PACKET.md`).
> Do not treat attachment manifests as the handoff packet.
```

Do not rewrite the v1 body. Do not point the generator at `handoffs.md`. Do not emit the packet under `docs/handoffs/`.

### New contract

`program/contracts/handoff-packet.md` is the only packet contract. Check and generator are data-driven off it (containers only, DEC-005).

### Instance directory (binding)

```text
handoff/
  PACKET.md
  decisions/          # copies of Accepted DECs cited
  glossary.md         # from CONTEXT.md
  questions/          # OQ copies with agreement states
  evidence/           # summary + cited EVD copies
  playbooks/          # implementation playbooks
  acceptance/         # executable tests for the bounded target
```

Never `docs/handoffs/`. Never `framework/`. `handoff/` is an allowed top-level instance path (like `briefs/`).

### `PACKET.md` front matter (binding)

| Key | Rule |
| --- | --- |
| `id` | **Architect default:** `HO-001` (one packet per instance this phase). |
| `date` | `YYYY-MM-DD` from `internal/clock` / `MYCELIUM_NOW`. |
| `implementation_system` | Enum: `pstack/poteto` \| `manual`. Generator writes `pstack/poteto`. `manual` is the floor. |
| `time_budget` | Required. Regex `^[0-9]+[mh]$`. Generator writes `30m`. Fixture uses `30m`. |

Unknown front-matter keys: **Architect default:** allowed (append-only), same spirit as optional DEC keys. Required keys above must be present and valid.

### Required H2s (binding; order fixed)

Framing; Locked decisions; Glossary; Open questions; Evidence summary; Implementation playbooks; Implementation system; Time budget; Acceptance.

Blueprint packet object (~168–170) maps onto those H2s: framing, locked decisions, glossary, open questions with agreement states, evidence summary, suggested implementation playbooks — plus the three PHASE-06 additions (Implementation system, Time budget, Acceptance).

Check validates **containers only** (DEC-005): H2s present, front matter valid, listed copies exist, ID links resolve inside the packet. Do not grade Framing prose. Do not grade playbook quality.

### Self-contained (binding)

The packet is self-contained. Check **FAILS** if `PACKET.md` or any file under `handoff/playbooks/` links to an instance path **outside** `handoff/` that was not copied in. ID links (`DEC-001`, `OQ-004`, `EVD-001`) must resolve to a file **inside** the packet copies (`handoff/decisions/`, `handoff/questions/`, `handoff/evidence/`). A link to `../decisions/DEC-001-*.md` (instance tree) is a FAIL even if that file exists outside the packet.

**Architect default — generator rewrite:** when copying, rewrite ID and relative links in `PACKET.md` and playbooks so they point at the in-packet copies. Do not leave instance-root paths.

### What the generator copies (Architect default)

`mycelium handoff` builds `handoff/` from the instance, then flips state.

| Packet path | Source |
| --- | --- |
| `PACKET.md` | Template `program/templates/handoff-packet.md` filled from instance + defaults (`pstack/poteto`, `30m`, `HO-001`). H2 bodies are structural lists (IDs, titles, paths) — not graded prose. |
| `decisions/` | Copy every instance DEC whose `status = "Accepted"`. If none, directory exists empty and Locked decisions H2 lists `none`. |
| `glossary.md` | Copy instance `CONTEXT.md` (or emit a file whose body is `none` if CONTEXT is empty). |
| `questions/` | Copy every instance OQ, preserving `agreement`. If none, directory exists empty and Open questions H2 lists `none`. |
| `evidence/` | Copy every EVD cited from Accepted DECs or from instance `evidence/`. If none, directory exists empty and Evidence summary H2 lists `none`. Write `handoff/evidence/SUMMARY.md` listing cited IDs. |
| `playbooks/` | Copy instance `playbooks/` if present and non-empty; else emit `handoff/playbooks/PLAYBOOK.md` from `program/templates/handoff-playbook.md` (required H2s: Target, Steps, Done). Canonical fixture **pre-places** playbooks so the generator copies them. |
| `acceptance/` | Copy instance `acceptance/` if present and non-empty; else emit `handoff/acceptance/README.md` stating `none` (structure pass, not a golden-impl pass). Canonical fixture **pre-places** executable tests so the generator copies them. |

A clarified instance with no playbooks and no acceptance still produces a **structurally** complete packet (stubs + `none`). MS-601 row 5 uses the fixture that *does* carry real tests.

**Architect default — no extra flags this phase.** Grammar is only `mycelium handoff [--dir PATH]`. Defaults are `pstack/poteto` and `30m`. A human may edit `PACKET.md` front matter in place after generation; check accepts either enum value and any `time_budget` matching the regex. Re-running `handoff` after the packet exists refuses (see §5).

## 5. `mycelium handoff` + handed-off reachability

This is the first command that may set `state = handed-off`.

### Grammar

```text
mycelium handoff [--dir PATH]
```

| Token | Rule |
| --- | --- |
| `--dir` / `--dir=` | Optional. Existing instance-root walk. |
| `-h` / `--help` | Usage, exit 0. |
| Other flags / extra positionals | Refuse (teaching error). |

### Legal only from `clarified`

| From | Command | Result |
| --- | --- | --- |
| `clarified`, no `handoff/PACKET.md` | `mycelium handoff` | Write packet, then `state = handed-off`. |
| `clarified`, `handoff/PACKET.md` exists and passes | `mycelium handoff` | **Refuse.** Packet already exists. Teach: `mycelium state handed-off`. |
| `clarified`, `handoff/PACKET.md` exists and passes | `mycelium state handed-off` | Flip state only. Do not regenerate. Not silent: packet is present. |
| `clarified`, packet missing or structure fail | `mycelium state handed-off` | **Refuse.** Teach: `mycelium handoff`. No writes. |
| any state except `clarified` | `mycelium handoff` | **Refuse.** Teach: legal only from `clarified`. |
| `handed-off` | `mycelium handoff` | **Refuse.** Already handed-off. |
| `handed-off` | `mycelium state exploring` (etc.) | **Refuse.** Terminal except `archived`. |
| `handed-off` | `mycelium state archived` | Legal (DEC-006 any → archived). No deletion. |
| `handed-off` | `mycelium state handed-off` | **Refuse.** Already in target. (Idempotent no-op is **not** this phase; refuse.) |

Silent handoff is forbidden: no path may set `handed-off` without a packet that passes structure.

### Teaching errors (four-line `teach` format; exit 1)

| Case | `what` (stderr line 1) | `fix` |
| --- | --- | --- |
| `state handed-off` without packet (or structure fail) | `state=handed-off requires a handoff packet` | `mycelium handoff [--dir PATH]` |
| `handoff` from non-`clarified` | `handoff is legal only from clarified (got <state>)` | `mycelium state clarified, then mycelium handoff` |
| `handoff` and `PACKET.md` already exists | `handoff/PACKET.md already exists` | `mycelium state handed-off [--dir PATH]` |
| `handoff` extra args / unknown flag | `handoff accepts only --dir PATH` | `mycelium handoff [--dir PATH]` |
| `handed-off` → non-archived | `illegal transition handed-off → <target>` | `handed-off is terminal except mycelium state archived` |
| Not an instance | `not a mycelium instance (no mycelium.toml found)` | `--dir PATH` |
| Leftover journal for a different op | existing `ErrJournalMismatch` text | `mycelium check --abort-journal` |

The pin's "packet command is not shipped" line is **replaced** (Slice 2). New text names `mycelium handoff`. Contract line may name `program/contracts/handoff-packet.md` or `lifecycle.md`. Do not delete files on refuse. Do not write a log line on refuse.

### Operation protocol

Same protocol as `new` / `state` / `supersede` (blueprint ~315): preflight / lock / stage / commit / rollback. Journal `op = "handoff"`.

| Step | Rule |
| --- | --- |
| Preflight | Manifest + log parse; state is `clarified`; no existing passing `handoff/PACKET.md`; no conflicting leftover journal. Nothing written. |
| Lock | Exclusive repo lock. |
| Stage | Write `handoff/**`, rewritten `index.md`, appended `log.md`, bumped `mycelium.toml` as staged files. Journal argv = the invocation. `Title` = `clarified -> handed-off`. |
| Commit order | **`handoff/**` (PACKET.md first, then copies), `index.md`, `log.md`, `mycelium.toml` last.** `mycelium.toml` sets `state = "handed-off"` and bumps `updated_date`. |
| Rollback | Failure before the first rename: remove staged files, change nothing. After a partial commit: journal survives; re-running the same argv resumes. |
| Detection | `mycelium check` already detects leftover journal / stale lock (item 9). `handoff` journals use the same recovery (`complete` or `--abort-journal`). |

**Architect default — `allowReplace`:** for `op == "handoff"` only, Commit may replace paths under `handoff/` plus `index.md` / `log.md` / `mycelium.toml`. Do not allow `new` to overwrite.

**Architect default — log line:**

```text
YYYY-MM-DD<TAB>handoff<TAB>HO-001<TAB>clarified -> handed-off
```

Date from `internal/clock` / `MYCELIUM_NOW`. ID column is `HO-001`. Note is exactly `clarified -> handed-off`.

**Architect default — success stdout:**

```text
mycelium handoff: ok
state: handed-off
packet: handoff/PACKET.md
```

Exit 0. Then `mycelium check --dir PATH` exits 0 on the happy path.

**Architect default — `state handed-off` success** (packet already present): same protocol **without** regenerating `handoff/**`. Commit order: `index.md`, `log.md`, `mycelium.toml` last. Log op is still `handoff` (the lifecycle event is the handoff, whether invoked as the verb or as `state handed-off`). Journal op is `handoff`. Stdout:

```text
mycelium state: ok
state: handed-off
packet: handoff/PACKET.md
```

**Architect default — index.md:** rewrite via existing `internal/indexmd`. Required H2s stay State / Artifacts / Log tail / Wake. State H2 reflects `handed-off`. Do not add a required H2.

**Architect default — resume match:** `journal.Matches` on argv, or `op+title` with `Title` = `clarified -> handed-off`.

## 6. pstack/poteto bridge + AGENTS.md

No pstack binary. No poteto CLI wrapper. No network. The bridge is documentation the isolated implementer reads **inside the packet** and in the emitted `AGENTS.md`.

### Binding mapping (full table in Appendix C)

| Packet section | pstack / poteto |
| --- | --- |
| Framing | why/how context; poteto-mode constraints |
| Locked decisions | decided list + DEC copies |
| Glossary | shared language |
| Open questions | agreement states; poteto candor |
| Evidence summary | citations |
| Implementation playbooks | how/ vertical slices for the bounded target |
| Implementation system | default `pstack/poteto`; `manual` is the floor |
| Time budget | required; fixture uses `30m` |
| Acceptance | executable tests in `handoff/acceptance/` |

### `AGENTS.md` — Implementation systems (binding)

`program/skeleton/AGENTS.md` gains an **Implementation systems** section (Slice 3). Required claims, verbatim-intent:

- The isolated implementer receives **ONLY** `handoff/`.
- No chat history.
- No instance source beyond the packet.
- Default system is **pstack/poteto-mode**.
- `manual` (read the packet, implement, run `handoff/acceptance/`) is the floor and always legal.
- Do not fetch the rest of the instance. Do not reopen locked decisions. Do not treat v1 attachment manifests as the packet.

**Architect default — emitted reference:** `program/reference/implementation-systems.md` (the mapping table + the four isolation rules) is emitted on scaffold and copied into `handoff/playbooks/` **only if** the instance has no playbooks of its own *and* the stub playbook needs a pointer. Prefer: AGENTS.md section + `handoff-packet.md` contract contain the table; the reference file is the portable copy. New scaffolds pick it up via embed. `tier` / `index` do not retrofit old instances.

## 7. Commands (what exists / what does not)

| Command | This phase |
| --- | --- |
| `version` / `help` / `new` / `check` / `tier` / `publish` / `index` / `state` / `wake` / `status` / `supersede` | Exist on the pin. Do not reopen. `check` + `state` gain §5 / §8 deltas. |
| **`mycelium handoff`** | **Add.** §5. |
| `mycelium state handed-off` | **Becomes legal IFF** `handoff/PACKET.md` exists and passes structure. |
| `mycelium council` / `second-opinion` / `ladder` / `replicate` | Do **not** add. OQ-003 stands. |
| `mycelium upgrade` / `migrate` / `release` / `install` | Do **not** add. |

Usage string in `internal/cli` gains one line:

```text
  mycelium handoff [--dir PATH]
```

Unknown-command teach for `council` stays "unknown command" (do not special-case). `handoff` is now a known verb.

## 8. Check updates

Items 1–23 stay; do not renumber. Add item 24. Extend item 8. Lift the PHASE-02/05 stored-`handed-off` FAIL.

| Item | Delta | Binds |
| --- | --- | --- |
| 8 | Log-op regex becomes `scaffold\|new\|tier\|publish\|check\|state\|wake\|supersede\|handoff`. | Slice 2 |
| 9 | Unchanged. A leftover `handoff` journal is an interrupted operation. | already |
| storage | `handed-off` is **LEGAL** iff packet structure passes. `handed-off` without packet still **FAIL**. `clarified` without packet still **LEGAL**. | Slice 2 |
| path | `handoff/` is an allowed top-level path. | Slice 2 |
| **24** | **New.** Packet structure IFF. If `handoff/PACKET.md` exists **or** `state = handed-off`: required front matter; nine H2s; directories present; copies exist for every ID listed; ID / path links resolve **inside** `handoff/`; `time_budget` matches `^[0-9]+[mh]$`; `implementation_system` is `pstack/poteto` or `manual`. If `state = handed-off` and `PACKET.md` is missing → FAIL. If `handoff/` is absent and state is not `handed-off` → item 24 does not fire. | Slice 2 |

### Lift timing (conformance.md)

| Slice | Check / lifecycle behavior |
| --- | --- |
| 1 | Parser + structure checker exist as functions. **No** CLI. **No** item 8 / 24 bind. **No** storage-rule lift. Stored `handed-off` still FAIL. `state handed-off` still refuses. A hand-written `handoff` log line would still fail check — do not write one in Slice 1 fixtures. |
| 2 | Command + `state handed-off` IFF + item 8 + item 24 + storage lift. Happy path + refuse table. |
| 3 | No new check item. AGENTS.md + skill + reference. |
| 4 | No new check item. Fixture + golden impl. Acceptance tests are executable (not a check-of-prose). |
| 5 | MS-601 matrix harness in `internal/clitest` runs every §14 row. Coverage **≥85%**. |

### What check must not do (additions)

| Temptation | Verdict |
| --- | --- |
| Grade Framing / playbook / glossary prose | **No.** DEC-005. |
| Fail `clarified` because `handoff/` is absent | **No.** |
| Pass stored `handed-off` without `handoff/PACKET.md` | **No.** |
| Resolve ID links to the instance tree outside `handoff/` | **No.** Self-contained. |
| Require a new index.md H2 | **No.** |
| Add `council` / `upgrade` to the log-op regex | **No.** |
| Call network / `gh` / read `GH_TOKEN` | **No.** |
| Treat an isolated-agent run or dogfood handoff as a check | **No.** |
| Treat an Actions job as a check | **No.** |

Teaching errors stay four lines, cap 20. Success stdout for `check` is unchanged.

## 9. Skills / AGENTS / mycelium-cli

No new skill. No new pack. No retrofit into existing instances.

| File | Delta |
| --- | --- |
| `program/skills/mycelium-cli/SKILL.md` | Document `mycelium handoff [--dir PATH]`. Name refuse cases (not clarified; packet already exists; `state handed-off` without packet). State that `handed-off` is terminal except `archived`. Name `handoff/PACKET.md` as the packet. Do not mention a `mycelium council` command as something to run. Do not point at v1 `handoffs.md` as the packet. |
| `program/skeleton/AGENTS.md` | **Implementation systems** section (§6). One short paragraph in the commands list: idea handoff is `mycelium handoff`; `state handed-off` only if the packet already exists. |
| `program/reference/implementation-systems.md` | **New.** Mapping table + isolation rules. Emitted on scaffold. |
| council / thinking / spark / wake / portfolio skills | Do not touch. |

New scaffolds pick up the updated files via embed. `tier` / `index` do not retrofit.

## 10. Templates / schema

`methodology_version` stays `2.0.0`. Re-run embed generate after edits.

| File | Delta |
| --- | --- |
| `program/templates/handoff-packet.md` | **New.** Front matter keys + nine H2s with `{{ID}}` / `{{DATE}}` tokens. `implementation_system` default `pstack/poteto`. `time_budget` default `30m`. H2 bodies start as `none` / ID lists the generator fills. |
| `program/templates/handoff-playbook.md` | **New.** Stub playbook H2s: Target, Steps, Done. Used only when instance `playbooks/` is absent. |
| `program/templates/handoff-packet.schema.toml` | **New** sidecar if the pin's schema DSL can express a non-`new <type>` artifact. **Architect default:** if the DSL is ID-namespace-only, the structure checker in `internal/handoff` hard-codes the four front-matter keys + nine H2s from the contract; do not invent a new schema DSL. Prefer a sidecar if `optional_front_matter` / `required_sections` already work. |
| Existing type schemas / templates | **Do not edit.** Handoff is not a `mycelium new <type>` artifact. |
| `CONTEXT.md` / pack templates | Do not edit. |

Do not add a `HO` namespace to `mycelium new`. `HO-001` is the packet id, allocated by the handoff generator, not by `new`.

## 11. Vertical slices 0–5

PR-per-slice, sequential, rebase on main. Arvo merges Quality-green PRs. Engineering opens PRs. Engineering does NOT push to main. One live PR at a time.

Each PR title: `PHASE-06 Slice N: <done-bar noun>`. Each PR body links this brief, the slice done bar, Linear **ROB-511**, and a `## Manual test guide` (Quality pulls the SHA onto their computer; Actions is not the proof). No drive-by refactors. No v1 deletions. No Actions job as a done bar. Do **not** add `.github/workflows/phase-06-*.yml`.

### Slice 0 — Commissioning (docs only)

This brief + `program/contracts/handoff-packet.md` + lifecycle rewrite (stored-`handed-off` IFF; teaching error names `mycelium handoff`) + `handoffs.md` five-line banner + `framework/phases/PHASE-06-acceptance.md` (rows = §14). No product code. No Go.

Done: files exist on a PR. Quality reads them against this brief. **Architect does not open this PR. Engineering lands it.**

### Slice 1 — Packet parser / schema + structure checker (no state lift)

- `internal/handoff`: parse `PACKET.md` front matter; validate four keys; validate nine H2s; walk copies; resolve ID / path links **inside** `handoff/` only; self-contained FAIL on outside links.
- Template + optional sidecar (§10). Embed generate.
- **No** `internal/cli` case. **No** item 8 / 24 bind. **No** storage-rule lift. **No** log line written by a command.

Done (hermetic `go test`): Appendix B bytes pass structure; missing H2 fails; missing `time_budget` fails; link to `../decisions/DEC-001-*.md` fails; in-packet `DEC-001` copy resolves; `mycelium handoff` is still "unknown command"; stored `handed-off` still FAILs check; `state handed-off` still refuses.

### Slice 2 — `mycelium handoff` + `state handed-off` IFF + check lift

- `internal/handoffcmd` + `internal/cli` verb.
- Protocol: preflight / lock / stage / commit order in §5. `allowReplace` extension.
- Item 8 regex + item 24 + storage lift (PHASE-02/05 "handed-off always FAIL" **removed in this slice**).
- `statecmd` / `lifecycle`: `handed-off` legal IFF packet passes; terminal except `archived`.
- `indexmd` rewrite. Log line. `state = handed-off`. Teaching error names `mycelium handoff`.

Done (hermetic `go test`): `handoff` from `clarified` → `state=handed-off`, `handoff/PACKET.md` present, log op `handoff`, `check` exit 0; `state handed-off` without packet → refuse, no writes; stored `handed-off` without packet → `check` FAIL; stored `handed-off` with packet → `check` PASS; `clarified` without packet → `check` PASS; `handed-off → exploring` refuse; leftover `handoff` journal resumes.

### Slice 3 — pstack/poteto bridge + AGENTS.md + mycelium-cli

- `program/skeleton/AGENTS.md` Implementation systems section.
- `program/reference/implementation-systems.md` (mapping table + isolation).
- `program/skills/mycelium-cli/SKILL.md` documents `mycelium handoff`.
- Embed generate.

Done: new `--offline` scaffold contains the section and the reference file; skill names `mycelium handoff` and `handoff/PACKET.md`; no pstack binary; no new pack.

### Slice 4 — Canonical fixture + acceptance tests + golden impl

- Fixture instance (Appendix D): bounded target `Add(a, b int) int` in a specified file; playbook + executable acceptance tests live in the fixture so the generator copies them into `handoff/`.
- Golden implementation in testdata passes those tests (proves the target is bounded and the tests are real).
- Time budget fixture value: `30m`.

Done (hermetic `go test`): generator from the fixture writes a structurally complete `handoff/`; copied playbooks + acceptance tests present; golden impl passes `handoff/acceptance/`; a broken impl fails them. No live agent.

### Slice 5 — MS-601 matrix harness + coverage ≥85%

**Architect default:** Slice 5 is the matrix harness in `internal/clitest` that runs **all** MS-601 rows (Appendix E), plus a cover-profile assertion that **new packages** (`internal/handoff`, `internal/handoffcmd`) are each **≥85%** statement/line and that product packages meet the **85%** floor (exclude generated / vendor / data-only fixtures).

Done: `go test ./...` runs the MS-601 matrix green and the 85% floor holds. That **is** the gate. An isolated-agent run is **not** the gate. A real-project handoff is **dogfood**, not the gate. An Actions log is **not** the gate.

## 12. Done / verified mapped onto MS-601

MS-601 is the hermetic phase gate. Blueprint MS-601 names an isolated-agent run. House test shape (locked 2026-08-15, restated by Arvo on this commission) says live Cursor is NOT a gate.

**Architect default — reconciliation:** the merge gate is the hermetic matrix below. The isolated-agent run goes in Quality's `## Manual test guide` as **optional human evidence**. A real-project handoff is **dogfood**, not the gate. Do not put a live agent in `go test`.

### MS-601 expected (authoritative; recipe in Appendix E)

| Row | Setup | Expect |
| --- | --- | --- |
| MS-601-GEN | Fixture instance (Appendix D); generator / `mycelium handoff` | Structurally complete `handoff/` (nine H2s, copies, self-contained links) |
| MS-601-CMD | Fixture in `clarified`; `mycelium handoff` | `state=handed-off`; log op `handoff`; `mycelium check` exit 0 |
| MS-601-REFUSE | `clarified`, no packet; `mycelium state handed-off` | Refuse; no writes |
| MS-601-CHECK | Stored `handed-off` without packet; stored `handed-off` with packet | FAIL; PASS. `clarified` without packet PASS |
| MS-601-GOLD | Fixture `handoff/acceptance/` + golden impl in testdata | Tests executable; golden passes; broken impl fails |
| MS-601-COV | Product packages + new packages | **≥85%** statement/line; new packages themselves **≥85%** |

`gh` never invoked. No `GH_TOKEN`. No Actions job. No live Cursor. No live network. Tests never touch the real home directory except as already allowed by PHASE-02 status fixtures (temp `UserHomeDir`).

### Slice → MS-601 map

| Slice | MS-601 clause |
| --- | --- |
| 0 | commissioning; not a runtime clause |
| 1 | parser / structure; not yet the command |
| 2 | MS-601-CMD / REFUSE / CHECK |
| 3 | bridge docs; not a runtime clause |
| 4 | MS-601-GEN (fixture copies) + MS-601-GOLD |
| 5 | the matrix + 85% in `go test ./...` **is** the gate |

PHASE-06 is accepted when MS-601 is green in `go test ./...` on main **and** Quality has thermos'd the table on their computer (local guide, `## Manual test guide`). Arvo accepts the phase. Engineering does not self-accept. Isolated-agent / dogfood may be attached as human evidence; they do not replace the gate.

## 13. Automated test plan

Engineering MUST write these tests. Quality thermos against this list on **their computer**. That local `go test ./...` is the verify bar, not an Actions log. Every ready PR includes `## Manual test guide`. Quality pulls the SHA onto their computer. Actions is not the proof.

### House test shape (locked 2026-08-15; numbered 1–5)

1. Fast unit tests on parsers/state/rules (most of the **85%** coverage floor).
2. Hermetic CLI: real binary, temp instance, fake clock/remotes, including fail cases.
3. Thin MS-601 fixture matrix.
4. Live GitHub, live Cursor, live network, browser e2e are NOT gates. Human evidence only.
5. No UI on this project — no browser tests.

Also stamp coverage **≥85%**. Do not put live e2e, isolated-agent, dogfood, or CI-as-proof in a done bar.

### Unit (no network, no gh, no home directory)

| Area | Cases |
| --- | --- |
| front matter | four keys; `time_budget` regex; enum `pstack/poteto` \| `manual`; missing key fail |
| H2s | all nine present; missing one fail; order required |
| self-contained | in-packet DEC resolves; `../decisions/...` fails; missing copy fails |
| storage rules | `handed-off` + packet pass; `handed-off` − packet fail; `clarified` − packet pass |
| transitions | `clarified → handed-off` legal iff packet; `handed-off → exploring` refuse; `handed-off → archived` legal |
| coverage | new packages ≥85%; product packages ≥85% |

### Hermetic CLI (built binary, temp dirs)

| Case | Expect |
| --- | --- |
| `mycelium handoff` from `clarified` | packet written; `state=handed-off`; log op `handoff`; check 0 |
| `mycelium state handed-off` without packet | exit 1; files unchanged; teach names `mycelium handoff` |
| `mycelium state handed-off` with passing packet | state flip; no regenerate; check 0 |
| `mycelium handoff` when `PACKET.md` exists | exit 1; teach `state handed-off` |
| `mycelium handoff` from `exploring` | exit 1; no writes |
| `handed-off` → `exploring` | exit 1 |
| stored `handed-off` without packet | `check` exit 1 |
| stored `handed-off` with packet | `check` exit 0 |
| leftover `handoff` journal | resumes |
| no `gh` | `gh` never invoked |
| MS-601 | Appendix E; `go test ./...` |
| 85% | cover profile on new + product packages |

Do not require live GitHub, live Cursor, a live agent, a real-project handoff, or an Actions job for `go test ./...`.

## 14. Acceptance / MS-601 fixture matrix

Slice 0 lands `framework/phases/PHASE-06-acceptance.md` with these rows. Later slices land the packages and testdata in Appendix F. No workflow file. No DEC-015 file. No `mycelium council` package. No `program/migrations/`. Do not reuse `handoffs.md` as the packet.

Each row: id, check, evidence, owner (Engineering | Arvo). **CI is not an owner.** Robert waived CI.

| id | check | evidence |
| --- | --- | --- |
| A-S0 | Slice 0 files exist and match this brief | PR diff; Quality read |
| A-S1 | parser / structure checker; no CLI verb yet; no state lift | hermetic `go test` |
| A-S2 | `mycelium handoff` + `state handed-off` IFF + item 8 + item 24 + storage lift | hermetic `go test` |
| A-S3 | AGENTS.md Implementation systems + mycelium-cli + reference | file read + scaffold emit |
| A-S4 | fixture copies + executable acceptance + golden impl | hermetic `go test` |
| A-S5 | MS-601 matrix green + **85%** coverage | `go test ./...` on Quality's computer |
| MS-601-GEN | structurally complete `handoff/` | A-S5 / A-S4 |
| MS-601-CMD | `handoff` from `clarified` → `handed-off` + log op | A-S5 / A-S2 |
| MS-601-REFUSE | `state handed-off` without packet refuses | A-S5 / A-S2 |
| MS-601-CHECK | stored `handed-off` IFF packet | A-S5 / A-S2 |
| MS-601-GOLD | acceptance tests + golden impl | A-S5 / A-S4 |
| MS-601-COV | product + new packages **≥85%** | A-S5 |

No isolated-agent-gate row. No dogfood-gate row. No Actions-job row. Quality should refuse a PR that adds an Actions job as the MS-601 gate. Quality should **not** refuse a missing PHASE-06 workflow.

## 15. Decided / Architect defaults / not decided

No TBD. Open items are decided inline and labeled **Architect default**. Do not bikeshed them in a code PR. Write a DEC to change one. No DEC-015 is required for these. OQ-004 is settled here.

Index of defaults that are easy to miss:

- **Packet contract:** `program/contracts/handoff-packet.md`. v1 `handoffs.md` gets a five-line banner only. Instance tree is `handoff/` (never `docs/handoffs/`, never `framework/`).
- **Front matter:** `id=HO-001`, `date` from clock, `implementation_system=pstack/poteto` (enum includes `manual`), `time_budget=30m` (`^[0-9]+[mh]$`). No extra CLI flags.
- **H2s:** Framing; Locked decisions; Glossary; Open questions; Evidence summary; Implementation playbooks; Implementation system; Time budget; Acceptance.
- **Self-contained:** ID / path links resolve inside `handoff/` copies. Outside instance paths FAIL.
- **New verb:** `mycelium handoff [--dir PATH]`. Only from `clarified`. Writes packet then `state=handed-off`. If `PACKET.md` already exists, refuse and teach `state handed-off`.
- **`state handed-off`:** legal IFF packet exists and passes. Missing → refuse, names `mycelium handoff`. Log op `handoff` either way.
- **Terminal:** `handed-off` → `archived` only. DEC-006 unchanged.
- **Commit order:** `handoff/**`, `index.md`, `log.md`, `mycelium.toml` last. CLI never git-commits.
- **MS-601 gate:** hermetic matrix + **85%** coverage. Isolated-agent = optional `## Manual test guide` human evidence. Real-project handoff = **dogfood**, not the gate. Not Actions.
- **No new pack.** No `mycelium council`. OQ-003 stands. `methodology_version` stays `2.0.0`. Version stays `0.1.0-dev`. Do not touch `internal/slug` (DEC-014).

**Not decided:** none — defaults above. (Empty on purpose.)

## 16. Risks, rollback, Quality should refuse

### Risks

| Risk | Mitigation |
| --- | --- |
| Silent `handed-off` (no packet) | §5. Quality should refuse. |
| v1 `handoffs.md` reused as the packet | §4 banner + Appendix F. Quality should refuse. |
| Packet emitted under `docs/handoffs/` or `framework/` | §4. Quality should refuse. |
| Isolated-agent or dogfood treated as the gate | §12. Quality should refuse. |
| Live agent inside `go test` | House shape item 4. Quality should refuse. |
| Actions job added as the MS-601 gate | Robert waived CI. Quality should refuse. Absence of a PHASE-06 workflow is correct. |
| Coverage below **85%** on new or product packages | Slice 5. Quality should refuse. |
| Migrations sneak back | DEC-011. Quality should refuse. |
| New pack / `mycelium council` | OQ-003. Quality should refuse. |
| Growing `latinFold` or adding `x/text` | DEC-014. Do not touch `internal/slug`. |
| CLI git-commits instance work product | DEC-010. Quality should refuse. |
| Emitting `framework/` or converting master | Absence tests stay. |
| `allowReplace` opened for `new` | §5: handoff targets only. |

### Rollback

Revert the offending PR on master. Do not `git push --force` to main. Floor is `a486cc9aa04a9ef22e7d0e564df3f2ebe692b1bb`. Do not delete Justfile / v1 scripts as a "cleanup" rollback. If Slice 2 must revert, stored `handed-off` FAIL and `state handed-off` refuse return (pin behavior).

### Quality should refuse

Refuse to approve if:

- an Actions job is added as the MS-601 gate, or `.github/workflows/phase-06-*.yml` is added, or `phase-01-hermetic.yml` is extended as a phase gate (Quality should **not** refuse a missing PHASE-06 workflow — absence is correct)
- isolated-agent, dogfood, live GitHub, live Cursor, live network, browser e2e, or CI-as-proof is required as the done bar
- a live agent is placed in `go test`
- `program/migrations/`, `just upgrade`, or `applied_migrations` appears
- `handed-off` succeeds without a packet, or check passes stored `handed-off` without `handoff/PACKET.md`
- `clarified` without a packet is made illegal
- v1 `program/contracts/handoffs.md` is treated as the packet contract (no banner, or generator reads it as `PACKET.md`)
- packet lands under `docs/handoffs/` or `framework/`
- a `mycelium council` / `upgrade` / `migrate` verb or a new pack appears
- `framework/` is emitted; CLI git-commits instance work product
- `latinFold` grows or NFKD is implemented (DEC-014)
- cobra / viper / yaml / testify / go-github / `golang.org/x/text` appears
- Justfile / v1 scripts deleted, `just init` run on master, or `research-program.toml` renamed
- DEC-012 / DEC-013 / DEC-014 reopened; DEC-015 recorded
- coverage on a new package or on product packages is below **85%**
- hermetic tests call network or real `gh`
- PR lacks `## Manual test guide`
- PR pushed straight to main

### Quality local guide (the verify bar)

Quality pulls the PR SHA onto their computer, runs `go test ./...` (and the cover profile for the **85%** floor), and thermos the MS-601 table in §14 / Appendix E. That local run is the verify bar, not an Actions log. Follow the PR's `## Manual test guide`. Optional: run an isolated agent against `handoff/` only (human evidence). Do not ask Engineering for a workflow badge. Quality should refuse Actions-as-MS-601-gate; Quality should not refuse a missing PHASE-06 workflow.

## 17. Execution order + Linear ROB-511

Same order as §11 (slices 0→5). PR-per-slice, sequential, rebase on main. One live PR at a time. Slice 2 is the command + check lift — do not combine it with 3–5. Slice 5 must be green in `go test ./...` on its PR (not Actions).

Title: `PHASE-06 Slice N: <done-bar noun>`. Body links this brief, the slice done bar, parent **ROB-511**, and `## Manual test guide`. Arvo cuts ROB-511 child tickets after this brief is stamped. No drive-by refactors, v1 deletions, PHASE-01–05 leftover work, or a later-phase command. Engineering opens PRs; Arvo merges Quality-green PRs; Engineering does NOT push to main.

Cursor cloud env name is exactly `robertguss/mycelium`. Go 1.26. `CGO_ENABLED=0`.

## 18. Handoff

### What Engineering starts with

This file. Only this file. Start from `https://github.com/robertguss/mycelium` at `a486cc9aa04a9ef22e7d0e564df3f2ebe692b1bb` (PHASE-01–05 accepted). Read `framework/blueprint.md` and DEC-001–014 for authority, not for a second plan. Execute Slice 0 first: land this brief at `framework/phases/PHASE-06-implementation-brief.md` plus `handoff-packet.md`, the lifecycle rewrite, the `handoffs.md` banner, and `framework/phases/PHASE-06-acceptance.md`. **Architect does not open the docs PR. Engineering lands it.** Do not implement from a later SHA unless Arvo re-pins. Cloud env `robertguss/mycelium`. Go 1.26. `CGO_ENABLED=0`. Linear parent: **ROB-511**.

### What Engineering must not do

See §16. Do not open a design debate in the PR. Do not write a second brief. Do not write DEC-015. Do not add an Actions job as the MS-601 gate. Do not start from a later SHA and "fast-forward" this brief. Do not invent `mycelium council`. Do not put an isolated-agent run or a real-project handoff in a done bar. Do not reuse v1 `handoffs.md` as the packet. Do not convert master. Do not emit `framework/`.

### What Quality reads

This brief, the acceptance matrix, the conformance delta, and the PR diff. Thermos: §13 house test shape + tests exist; §16 refuse list is clean; MS-601 hermetic `go test ./...` on Quality's computer; **85%** coverage; no Actions job as the MS-601 gate; no migrations; no new pack; packet is `handoff-packet.md` not v1 `handoffs.md`.

Quality should refuse a PR that adds an Actions job as the MS-601 gate; Quality should **not** refuse a missing PHASE-06 workflow.

### What Arvo does

Cuts ROB-511 child tickets after this brief is stamped. Merges Quality-green PRs. Accepts PHASE-06 when MS-601 is green on main and Quality has thermos'd locally. May attach an isolated-agent run or a real-project handoff as dogfood / human evidence; must not treat them as the gate. Does not re-pin without writing the new SHA. Does not record DEC-015.

## Appendix A — No new DEC (OQ-004 settled here)

No DEC-015. PHASE-06 implements DEC-006 reachability of `handed-off` (machine unchanged; packet IFF) and settles **OQ-004** (packet format + pstack/poteto mapping) as Architect defaults in §§4–6 and Appendix C. PHASE-06 does not reopen DEC-012, DEC-013, or DEC-014. Remaining choices are Architect defaults in §15. Engineering lands **zero** new files under `framework/decisions/`.

If a later discovery contradicts a locked item, stop and write a DEC; do not silently patch this brief in a code PR.

## Appendix B — PACKET.md example

Fixture instance (not master `framework/`). Tests must **not** grade the words — only front matter, H2s, copies, self-contained links, log line, and `mycelium check` exit 0.

```text
---
id = "HO-001"
date = "2026-08-15"
implementation_system = "pstack/poteto"
time_budget = "30m"
---

# Handoff packet

## Framing

Bounded target: implement Add(a, b int) int in add.go.

## Locked decisions

- DEC-001 (Accepted) — see decisions/

## Glossary

See glossary.md (copied from CONTEXT.md).

## Open questions

none

## Evidence summary

See evidence/SUMMARY.md.

## Implementation playbooks

See playbooks/PLAYBOOK.md.

## Implementation system

pstack/poteto (manual is the floor).

## Time budget

30m

## Acceptance

See acceptance/. Executable tests for Add.
```

`handoff/decisions/` contains the Accepted DEC-001 copy. `handoff/glossary.md` exists. `handoff/questions/` exists (empty ok). `handoff/evidence/SUMMARY.md` exists. `handoff/playbooks/PLAYBOOK.md` exists. `handoff/acceptance/` contains executable tests. Links in `PACKET.md` / playbooks point at those copies, not at `../decisions/`.

Happy path:

```text
mycelium new idea "Handoff Fixture" --offline --dir PATH
# bring instance to clarified; place fixture playbooks + acceptance (Appendix D)
mycelium handoff --dir PATH
mycelium check --dir PATH
# exit 0
# mycelium.toml state = handed-off
# log op = handoff
```

Negatives in the same test file: `state handed-off` on a clarified instance with no `handoff/PACKET.md` → exit 1, no writes; stored `handed-off` with `handoff/` deleted → `check` exit 1; `handoff` from `exploring` → exit 1; `state exploring` from `handed-off` → exit 1.

## Appendix C — pstack/poteto mapping table (full)

| Packet section | pstack / poteto | Isolated implementer does |
| --- | --- | --- |
| Framing | why/how context; poteto-mode constraints | Read first. Do not reopen. |
| Locked decisions | decided list + DEC copies in `handoff/decisions/` | Treat as Accepted. Do not restage. |
| Glossary | shared language (`handoff/glossary.md`) | Use these terms. Challenge drift only inside the packet. |
| Open questions | agreement states; poteto candor ("no is an acceptable answer") | Honor `open` / `aligned` / `agree-to-disagree`. Do not silently flip. |
| Evidence summary | citations (`handoff/evidence/`) | Cite; do not invent EVDs. |
| Implementation playbooks | how/ vertical slices for the bounded target | Execute the slices. Default system is pstack/poteto-mode. |
| Implementation system | default `pstack/poteto`; `manual` is the floor | If pstack is unavailable, `manual` still satisfies the contract. |
| Time budget | required; fixture uses `30m` | Stop at the budget. Do not expand scope. |
| Acceptance | executable tests in `handoff/acceptance/` | Run them. Green = the bounded target is done. |

Isolation (binding): implementer receives ONLY `handoff/`; no chat history; no instance source beyond the packet.

## Appendix D — Canonical fixture sketch (Add)

Tiny bounded target so MS-601-GOLD is real and hermetic.

**Architect default** layout under `internal/clitest/testdata/handoff-fixture/` (a frozen clarified instance):

```text
mycelium.toml                 # state = clarified
CONTEXT.md                    # term: Add
decisions/DEC-001-add-signature.md   # status = Accepted; Add(a, b int) int
playbooks/PLAYBOOK.md         # Target: implement Add in add.go; Steps; Done
acceptance/add_test.go        # executable tests: Add(2,2)=4; Add(-1,1)=0; Add(0,0)=0
```

Generator copies those into `handoff/` (plus `PACKET.md`, `glossary.md`, `evidence/SUMMARY.md`, `questions/`).

**Golden implementation** under `internal/clitest/testdata/golden-add/`:

```text
add.go    # func Add(a, b int) int { return a + b }
```

Slice 4 test: copy golden `add.go` next to the copied `handoff/acceptance/` tests and run them (exit 0). Copy a broken `add.go` (`return 0`) and run them (exit 1). That proves the tests are executable and the target is bounded. This is **not** product CLI code. Do not import the golden package from `cmd/mycelium`.

Playbook names the file `add.go` and the signature `Add(a, b int) int`. Check does not grade the playbook words (DEC-005). The golden-impl row grades the *tests*, not the playbook prose.

Time budget in the fixture `PACKET.md`: `30m`. Implementation system: `pstack/poteto`.

## Appendix E — MS-601 fixture recipe

Hermetic. No network. No `gh`. No `GH_TOKEN`. `go test ./...` only. Do **not** add an Actions job or `.github/workflows/phase-06-*.yml`. Do **not** extend `phase-01-hermetic.yml` as a phase gate. Temp dirs; freshly built binary; fake clock (`MYCELIUM_NOW=2026-08-15T00:00:00Z` or `internal/clock`). Do not invoke a live agent. Do not invoke a real-project handoff.

**Architect default:** matrix test lives at `internal/clitest/ms601_hermetic_test.go` (execs the binary). Parser tests live in `internal/handoff`. Coverage assertion may live in `internal/handoff` / `internal/handoffcmd` or the matrix file.

### MS-601-GEN / MS-601-CMD

```text
# start from Appendix D fixture (or scaffold + place playbooks/acceptance + state clarified)
mycelium handoff --dir PATH/ho
mycelium check --dir PATH/ho
# exit 0
# assert handoff/PACKET.md H2s + front matter
# assert copies + self-contained links
# assert mycelium.toml state = handed-off
# assert log op = handoff
```

### MS-601-REFUSE / MS-601-CHECK

```text
# clarified, no handoff/
mycelium state handed-off --dir PATH/nopkt
# exit 1; tree unchanged

# hand-set state = handed-off with no handoff/
mycelium check --dir PATH/stored-bare
# exit 1

# stored handed-off + valid handoff/
mycelium check --dir PATH/stored-ok
# exit 0

# clarified, no handoff/
mycelium check --dir PATH/clarified
# exit 0
```

### MS-601-GOLD

Run fixture `handoff/acceptance/` against `testdata/golden-add` (pass) and against a broken stub (fail).

### MS-601-COV

`go test -coverprofile` on product packages (exclude generated / vendor / data-only fixtures) and on `internal/handoff` + `internal/handoffcmd` separately. Assert **≥85%** statement/line. Stamp: **85**.

`gh` invoked → FAIL the test. No isolated-agent process. No dogfood repo. No Actions job.

## Appendix F — File tree / DO NOT ADD

### Master (additions on top of the PHASE-01–05 tree; v1 files retained)

```text
internal/handoff/                      # Slice 1 parser / structure
internal/handoffcmd/                   # Slice 2 command
internal/clitest/ms601_hermetic_test.go
internal/clitest/testdata/handoff-fixture/
internal/clitest/testdata/golden-add/
program/contracts/handoff-packet.md
program/contracts/handoffs.md          # five-line banner only
program/contracts/lifecycle.md         # handed-off IFF + teach mycelium handoff
program/contracts/conformance.md       # item 8 + item 24 + lift timing
program/reference/implementation-systems.md
program/templates/handoff-packet.md
program/templates/handoff-playbook.md
program/skills/mycelium-cli/SKILL.md
program/skeleton/AGENTS.md             # Implementation systems
framework/phases/PHASE-06-implementation-brief.md
framework/phases/PHASE-06-acceptance.md
internal/embed/program/                # regenerate after program/ edits
```

### DO NOT ADD

```text
.github/workflows/phase-06-*.yml
.github/workflows/phase-06-hermetic.yml
.github/workflows/phase-06-ms601.yml
framework/decisions/DEC-015-*.md
program/migrations/
program/packs/{handoff,registry.toml}
a mycelium council / upgrade / migrate command package
docs/handoffs/                         # not the instance packet path
reuse of program/contracts/handoffs.md as the packet
framework/ emit into instances
internal/slug changes
```

Do **not** add a PHASE-06 workflow. Do **not** extend `phase-01-hermetic.yml` as a phase gate. Quality should refuse a PR that adds an Actions job as the MS-601 gate. Quality should **not** refuse a missing PHASE-06 workflow (absence is correct). Do not delete Justfile / v1 scripts / `research-program.toml` / PHASE-01 workflows. Do not add a `mycelium council` command. Do not add migrations. Do not convert master. Do not emit `framework/`. Do not treat v1 `handoffs.md` as `handoff/PACKET.md`.

### Emitted instance (spark / focused, local-only, PHASE-06 scaffold)

Unchanged from PHASE-05 plus updated `mycelium-cli` / `AGENTS.md` / `implementation-systems.md`. After `mycelium handoff` (Slice 2+): `handoff/` written; `log.md` appended; `index.md` rewritten; `mycelium.toml` `state = handed-off` + `updated_date`.

Absent: `framework/`, `cmd/`, `internal/`, `go.mod`, `Justfile`, `scripts/`, `research-program.toml`, a `council` CLI verb, `program/migrations/`, `docs/handoffs/`.

Unexported helpers may live next to their tests. No `pkg/`. Do not touch `internal/slug`.

End of PHASE-06 implementation brief. Engineering executes from this file only.
