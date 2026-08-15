# PHASE-04 Implementation Brief — Perspective ladder

- **Status:** Binding
- **Date:** 2026-08-15
- **Audience:** Engineering (pstack / poteto-mode)
- **Authority:** `framework/blueprint.md` (Accepted 2026-08-14). DEC-001 through DEC-014 stand. **DEC-008 is this phase** (already Accepted). Implement it. Do **not** record DEC-015 (see Appendix A).
- **Repo:** https://github.com/robertguss/mycelium
- **Pin:** Engineering starts from `main` @ `d520560d138937a8904c43623524d2d2b9ef82f3` (PHASE-01–03 accepted; MS-301 on main via PR #40). Do not implement from a later SHA unless Arvo re-pins in writing.
- **Product:** single-binary Go CLI `mycelium`. Master builds the CLI. `program/` is `go:embed`'d and emitted on scaffold. PHASE-04 adds the council pack (`program/packs/council`) — presence-is-registration — plus commissioning / report / reconciliation grammar for the perspective ladder. **No portable council CLI.** **No new top-level CLI verb.**
- **Phase gate:** MS-401 via hermetic `go test ./...` only (perspective-ladder acceptance matrix). GitHub Actions is **not** a gate (Robert waived CI). Do **not** add `.github/workflows/phase-04-*.yml`. Do **not** extend `phase-01-hermetic.yml` as a phase gate. Quality should refuse a PR that adds an Actions job as the MS-401 gate. Quality should **not** refuse a missing PHASE-04 workflow (absence is correct).
- **How to use this file:** Engineering executes from THIS FILE ONLY. No "see chat". No TBD. Open items are decided here and labeled **Architect default**.

Headings: §§1–19 then Appendices A–F (no DEC-015, second-opinion CMP+RPT example, council CMP+two RPTs+RCL with SEED-DISSENT, `panels.toml` example, MS-401 fixture recipe, target file tree).

Cloud env name is exactly `robertguss/mycelium`. Go 1.26 at `/usr/local/go` or `/usr/local/bin/go`. `CGO_ENABLED=0`. stdlib + `github.com/pelletier/go-toml/v2` only. No cobra / viper / yaml / testify / go-github.

## 1. Scope / out of scope

Tonight is PHASE-04 only. Later phases get their own briefs after MS-401 is accepted. Do not implement, stub-ship, or "leave a hook command" for later phases. Do not implement PHASE-01–03 leftovers. Do not reimplement shipped sparring.

### In scope

- Council pack at `program/packs/council/` (first pack; presence-is-registration; no registry file). Drop in to enable, delete the directory to disable.
- OQ-003 remaining sliver, locked here: **only `council` is a pack this phase** (§5).
- Commissioning + report + reconciliation contracts, templates, and sidecar schemas for model-diverse replication.
- Second-opinion move (rung 2) and council move (rung 3). Sparring (rung 1) stays as shipped.
- Cursor council adapter and manual-floor adapter as **markdown procedures**, not Go.
- Panel presets documented at `~/.config/mycelium/panels.toml` (override `$MYCELIUM_CONFIG` as a directory). Check does **not** read user config.
- Namespaces `CMP-###`, `RPT-###`, `RCL-###` with homes under `reviews/`.
- Check updates listed in §9 (pack presence, collision, `reviews/` extra-top-level, provenance / `prompt_sha256` / cardinality / `opt_in` / `cost_class` / `SEED-DISSENT`).
- Pack skills `council` and `second-opinion`, emitted on **new scaffolds only** when the pack is present in emitted `program/`.
- `program/skills/mycelium-cli/SKILL.md` and `program/skeleton/AGENTS.md` capability note: runtimes that cannot fan out skip rungs 2–3; sparring still applies.
- MS-401 hermetic fixture tests in `go test ./...`. No Actions job.
- Commissioning artifacts: this brief, pack contracts, conformance items 18–22, acceptance stub.
- Master `go:embed` includes `program/packs/council/`. New instances receive the pack. Core checks must not require the pack.

### Out of scope (Quality refuses PRs that add them)

- Any new top-level CLI verb: `mycelium council`, `second-opinion`, `ladder`, `replicate`, `handoff`, `supersede`. **No portable council CLI** in v1 (DEC-008 clause 5). Surface is existing `mycelium new <type>` (discovers pack templates when the pack is present), `mycelium check` (loads pack checks when `program/packs/council/` exists), plus pack skills.
- A CLI that calls a model, Cursor, or the network for ladder rungs.
- Extra packs beyond `council`. Explicitly **not** packs this phase: thinking, spark, wake, portfolio, mycelium-cli, check, generate, lifecycle, status, index, publish, sparring, handoff (PHASE-06), supersede (PHASE-05).
- Forking `program/contracts/replication-reconciliation.md`. Pack reconciliation contract points at it verbatim.
- Majority-vote, prose-confidence, or model-reputation selectors.
- Requiring a council (or any CMP/RPT/RCL) to leave `spark` or to pass core check.
- Content-scoring of Position / Findings / Dissent / Retained dissent / prompt prose (DEC-005).
- Requiring `~/.config/mycelium/panels.toml` (or `$MYCELIUM_CONFIG/panels.toml`) to exist. Check does not read it.
- Enforcing panel size beyond council ≥2 RPTs.
- A live Cursor council run as the MS-401 gate. Live run is human evidence for Arvo, **not** the gate.
- Citing OQ-006's blueprint adversarial review as MS-401 evidence (different prompt; neither a council nor a second-opinion).
- PHASE-05 `mycelium supersede` and release/install work.
- PHASE-06 handoff packet. `handed-off` stays unreachable.
- Implementing MS-101(b). Do not commission a `GH_TOKEN` job. Do not reopen publish.
- Growing `latinFold`, adding NFKD, or adding `golang.org/x/text` (DEC-014).
- Converting master (`research-program.toml`, `just init`, deleting Justfile/scripts).
- Emitting `framework/`.
- CLI `git add` / `git commit` of instance work product.
- Retrofitting pack skills into existing instances. No retrofit command.
- Adding `.github/workflows/phase-04-*.yml` or extending `phase-01-hermetic.yml` as a phase gate. Quality should refuse a PR that adds an Actions job as the MS-401 gate.
- Reopening PHASE-03 sparring or DEC-007 agreement grammar.

### Master vs instance (unchanged)

Master remains an ADRP v1 instance for its own evolution. Do not convert master's `research-program.toml` to `mycelium.toml`. Do not run `just init` on master. `framework/` stays master-only and is NEVER emitted. Justfile/scripts stay on master. PHASE-04 changes the operational surface for *idea instances* only.

## 2. Authority and do-not-reopen

### Governing documents

| Document | Rule |
| --- | --- |
| `framework/blueprint.md` (Accepted 2026-08-14) | Do not rewrite vision. DEC-001–014 stand. PHASE-04 text at blueprint lines ~420–431 is the scope ceiling. OQ-003 sliver (~462–464) is locked in §5. OQ-006 (~469–474) is not MS-401 evidence. |
| DEC-005 | Checks validate containers, never contents. No automated content score. |
| DEC-006 | spark → exploring ⇄ simmering → clarified → handed-off; any → archived. Do not change the machine. `handed-off` stays unreachable. |
| DEC-007 | Sparring. Already shipped (PHASE-03). Do not reopen the grammar. |
| **DEC-008** | **This phase.** Three-rung ladder; councils opt-in; engine-agnostic contracts; Cursor adapter first; **No portable council CLI** in v1; panel presets in user config. Already Accepted. Implement it. |
| DEC-010 / DEC-011 | CLI never commits. No migration. Runtime reads instance files. |
| DEC-012 / DEC-013 / DEC-014 | Do not reopen (`mycelium.toml`; refuse out-of-range; `latinFold` only, no NFKD, no `x/text`, do not grow the map). |
| `program/contracts/replication-reconciliation.md` | v1 contract. Reuse verbatim. Do not fork. |
| This brief | Binding 2026-08-15. PHASE-04 only. Architect defaults are binding. No DEC-015. |

### Process override (unchanged)

Blueprint "humans-own-git" is overridden for the *master* repo's engineering process: Arvo merges Quality-green PRs and accepts the phase. Engineering opens PRs. Engineering does NOT push to main. The CLI still never git-commits *instance* work product.

### Do not reopen

Do not reopen the product shape, the language, the dependency floor, the state vocabulary, the manifest filename, the refuse-vs-warn range rule, the no-commit rule, the instance-files-are-truth rule, slugify/DEC-014, publish, MS-101(b), or PHASE-03 sparring grammar. If a later discovery contradicts a locked item, stop and write a DEC; do not silently patch this brief in a code PR.

Do not reopen DEC-012, DEC-013, DEC-014, or DEC-007. Do **not** record DEC-015. DEC-008 already exists and is this phase.

### Leftovers stay leftovers; phases you must not commission

MS-101(b) has **not** passed (`GH_TOKEN` missing; skip ≠ pass) — do not implement it or commission a `GH_TOKEN` job. PHASE-02 seven-real-day dogfood and a live Cursor council run are human evidence, not gates. Do not reimplement shipped PHASE-01–03 work. Do not commission PHASE-05 (`mycelium supersede`) or PHASE-06 (handoff packet; `handed-off` stays unreachable).

## 3. What PHASE-01–03 left on main (floor; do not reimplement)

Pin: `d520560d138937a8904c43623524d2d2b9ef82f3`. Treat this SHA as the floor. Reuse packages. Do not rewrite working PHASE-01, PHASE-02, or PHASE-03 commands. There is **no** `program/packs/` on the pin. Council is the first pack.

### Already shipped (do not rebuild)

Reuse: `cmd/mycelium`, `internal/{cli,version,embed,clock,execrun,metadata,idpath,manifest,schema,slug,logfmt,teach,lock,journal,op,scaffold,generate,check,tiercmd,publish,clitest,revisit,lifecycle,indexmd,wakebrief,statecmd,statuscmd,sparring}`.

| Touch | Fate |
| --- | --- |
| `internal/cli` | **No new verb.** Commands already: `version`, `help`, `new`, `check`, `tier`, `publish`, `index`, `state`, `wake`, `status`. Do not add `council` / `second-opinion` / `ladder` / `replicate` / `handoff`. |
| `internal/check` | **Extend** with pack presence, collision, `reviews/` extra-top-level (Slice 1), then ladder IFF rules (Slice 3). Do not rewrite the package. |
| `internal/generate` / `internal/schema` | **Extend** to discover `program/packs/*/templates/*.schema.toml` when the pack directory exists (Slice 2). Do not hardcode council type keys in core. |
| `internal/idpath` | **Extend** only if nested homes (`reviews/commissioning`) need a multi-segment `home`. Do not flatten pack homes to top-level. |
| `internal/scaffold` | Emit pack skills on **new scaffolds only** when the pack is present in emitted `program/` (Slice 4). No retrofit. |
| `internal/embed` | Re-run `go generate` after `program/` edits. Embed includes `program/packs/council/`. |
| `internal/slug` | Do not touch (DEC-014). |
| `internal/clock` + `MYCELIUM_NOW` | Reuse if a fixture wants stable dates. **Not required** for MS-401 unless you choose to pin dates on generated files. |
| `internal/execrun` | Unused by MS-401. Hermetic tests still must not call `gh`. |
| `internal/lock` / `journal` / `op` | Do **not** add a `council` / `replicate` / `ladder` op. `new` already covers `mycelium new commissioning`. |
| `internal/lifecycle` / `statecmd` / `wakebrief` / `indexmd` / `statuscmd` | Do not rewrite. Lifecycle storage rules stay PHASE-02. Do not add a required `index.md` H2. |
| `internal/sparring` | Do not reopen. PHASE-03 grammar stays. |
| `program/contracts/{conformance,lifecycle,identifiers,naming,replication-reconciliation,sparring}.md` | Extend conformance with items 18–22. Add CMP/RPT/RCL to identifiers + naming. **Do not fork** the v1 replication contract. Do not edit sparring grammar. |
| `program/skills/{spark,wake,portfolio,thinking,mycelium-cli}` | Do not rewrite spark/wake/portfolio/thinking. Update `mycelium-cli` with the capability note and pack surface. |
| `phase-01-hermetic.yml` / `phase-01-github.yml` | Leave alone. Do **not** add a PHASE-04 workflow. Do **not** extend `phase-01-hermetic.yml` as a phase gate. Actions is not a gate. |
| Justfile / scripts / `research-program.toml` | Keep. Do not delete. Do not `just init`. |
| `framework/` | Master-only. NEVER emitted. |
| `internal/version` | Stay `0.1.0-dev` unless already different on the pin — do not bump as a phase ritual. `methodology_version` stays `2.0.0`. Re-run embed generate after `program/` edits. |

### PHASE-01–03 behaviors that stay true

Birth state is `spark`. Stored `handed-off` still FAILS check. `index.md` H2s stay State / Artifacts / Log tail / Wake — **do not add a required H2.** Log-ops stay `scaffold|new|tier|publish|check|state|wake`. Exit 0/1; four-line teaching errors; cap 20. `--dir` and instance-root walk unchanged. CLI never git-commits. Runtime reads instance `program/` (DEC-011). Module `github.com/robertguss/mycelium`. `CGO_ENABLED=0`. Go 1.26. Items 1–17 stay (do not renumber). Existing `new question|assumption|decision|spike` stay. Spark with zero questions still passes; sparring does not require a council.

### What must not be broken

`just check` on master; hermetic `go test ./...`; no `framework/` emit; no master conversion. PHASE-01/02/03 fixtures stay green. Spark instances with zero questions stay green. Instances without `program/packs/council/` and without `reviews/` stay green. PHASE-03 aligned/disputed OQ fixtures stay green.

If a PHASE-04 PR is bad: revert that PR. Floor is the pin SHA.

## 4. Perspective ladder + DEC-008 table

DEC-008 is Accepted. This section encodes it as check vs skill vs adapter. Check stays structure-only (DEC-005). Skills and markdown adapters carry the stance so it survives model and runtime changes. The CLI never thinks and never calls a model.

### Three rungs (do not bikeshed)

| Rung | Cost | Who | Check this phase | Skill / adapter this phase |
| --- | --- | --- | --- | --- |
| 1. Sparring | free, always on | resident agent | Already shipped (PHASE-03). Do not change. | `thinking` skill. Unchanged. |
| 2. Second-opinion | cheap, one word to invoke | exactly one different model, identical commissioning prompt | CMP `rung = "second-opinion"`; `opt_in = true`; `cost_class = "cheap"`; exactly one matching RPT when any RPT exists; no RCL; `prompt_sha256` matches | `second-opinion` pack skill + chosen adapter |
| 3. Council | expensive, opt-in, rare | N different models; full replication + reconciliation | CMP `rung = "council"`; `opt_in = true`; `cost_class` in `quick\|standard\|high-stakes`; ≥2 matching RPTs + exactly one matching RCL once any RPT/RCL exists; `prompt_sha256` matches; `SEED-DISSENT` substring rule | `council` pack skill + chosen adapter |

### DEC-008 clauses

| # | Stance (DEC-008) | Check this phase | Skill / adapter this phase |
| --- | --- | --- | --- |
| 1 | Three-rung ladder: sparring, second-opinion, council. | Rung enum on CMP/RPT/RCL. Sparring files unchanged. | Skills name the rungs. Sparring still applies when rungs 2–3 are skipped. |
| 2 | Councils are opt-in, never a required stage, never auto-run. Agent suggests only when v1 replication triggers fire, stating panel size and cost class first. | `opt_in` must be `true` (TOML boolean). Missing or `false` → FAIL. Check does **not** require any CMP to exist. Spark without a council still passes. Check does **not** read panel presets. | Skill states cost class before running. Skill cites v1 triggers. |
| 3 | Reuse v1 replication + reconciliation verbatim. Replicas on different models. Reports land per-model. Reconciliation retains dissent. Majority vote, prose confidence, and model reputation stay banned. | Independent files: N RPT paths. RCL required H2s include `Retained dissent`. Check does not content-score the method and does not require distinct `model` strings. | Skill + RCL contract state the bans. Adapter does not chairman-smooth. |
| 4 | Contracts are engine-agnostic: commissioning-prompt and report file shapes only. Producers are swappable adapters. Cursor parallel subagents first; manual floor satisfies the contract with zero tooling. | `adapter` enum `cursor\|manual`. CLI does not invoke either. | `adapters/cursor.md` and `adapters/manual.md` are procedures. |
| 5 | **No portable council CLI** in v1. `AGENTS.md` carries a capability note so non-fan-out runtimes skip rungs 2–3 gracefully. | No `council` verb. Unknown-command path stays. | Capability note in AGENTS.md + mycelium-cli. Skills must not tell the agent to run `mycelium council`. |
| 6 | Panel presets (quick / standard / high-stakes) live in user-level config. | **None.** Check does not read user config. Check does not enforce panel size beyond council ≥2. | Skill documents defaults (quick=2, standard=3, high-stakes=4) and the file path. Missing file is legal. |

### v1 replication triggers (skill cites; check does not fire them)

Recommend a council when a decision is security-critical, safety-critical, legally or financially consequential, difficult to reverse, architecturally foundational, based on weak or conflicting evidence, vulnerable to ecosystem/vendor bias, or still low-confidence after a spike. Source: `program/contracts/replication-reconciliation.md`. Do not fork that file.

### Engine-agnosticism

Contracts define file shapes. _How_ reports get produced is a swappable adapter. The CLI never calls a model, Cursor, or the network for ladder rungs. Runtimes that cannot fan out skip rungs 2–3; sparring still applies.

## 5. Pack presence-registration (council only; OQ-003 sliver)

Blueprint item 5 and OQ-003: packs are presence-registered directories under `program/packs/<name>/`. No registry file. The remaining sliver — which capabilities beyond council become packs — **is locked here**.

### OQ-003 sliver (binding; do not leave open)

**Only `council` is a pack this phase.**

Explicitly **not** packs (this phase or as a drive-by in a PHASE-04 PR):

```text
thinking
spark
wake
portfolio
mycelium-cli
check
generate
lifecycle
status
index
publish
sparring
handoff
supersede
```

`handoff` stays PHASE-06. `supersede` stays PHASE-05. Do not create `program/packs/<any-of-the-above>/`. Quality should refuse extra packs.

### Presence-is-registration

| Rule | Binding |
| --- | --- |
| Enable | Directory `program/packs/council/` exists in the **instance** `program/` tree. |
| Disable | Delete `program/packs/council` on that instance. |
| Registry file | **None.** Do not add `program/packs/registry.toml` or a manifest pack list. |
| Master embed | `go:embed` includes `program/packs/council/`. New scaffolds receive the pack. |
| Core checks | Must **not** require the pack. An instance with the pack deleted and no `reviews/` still passes items 1–17. |
| Discovery | Check, generate, and schema loaders scan `program/packs/*/`. |
| Skills emit | `.agents/skills/{council,second-opinion}/` on **new scaffolds only** when the pack is present in emitted `program/`. |

### Collision rule (write it even though only one pack exists)

Conformance fails namespace collisions between packs, and between a pack and core.

**Architect default — collision keys.** Two registrations collide when they share any of:

- ID namespace (`CMP`, `RPT`, `RCL`, or any core NS: `DEC`, `ASM`, `EVD`, `SPK`, `FND`, `REC`, `REQ`, `OQ`, `RSK`, `PHASE`, `MS`)
- Type key (filename stem of `*.schema.toml`: `commissioning`, `model-report`, `reconciliation`, or any core type)
- Home directory (`reviews/commissioning`, `questions`, …)

Collision → FAIL. Teaching error names both pack directories (or pack vs `program/templates/`) and `program/contracts/conformance.md`.

Only `council` ships. Still write the rule. The Slice 1 fixture fabricates a second pack **inside a temp instance only**. Do not add a second pack under master's `program/packs/`.

### `reviews/` extra-top-level

`reviews/` is an allowed top-level path **only when the council pack is present**.

| Pack present? | `reviews/` present? | Deviation `extra-top-level:reviews/`? | `mycelium check` |
| --- | --- | --- | --- |
| yes | no | n/a | PASS (core + pack; no CMP required) |
| yes | yes | n/a | PASS on the path; then apply pack type rules inside |
| no | no | n/a | PASS (core check still passes) |
| no | yes | no | **FAIL** extra-top-level |
| no | yes | yes | PASS on the path (leftover waived) |

**Architect default:** extra-top-level looks at first-level names only. Pack type homes are nested (`reviews/commissioning/`, `reviews/reports/`, `reviews/reconciliations/`). Do not flatten them to top-level `commissioning/` / `reports/` / `reconciliations/`. `reviews/README.md` is allowed when the pack is present (not an ID-to-path home). Inside a type home, files must match the filename pattern or be `README.md`.

### Pack file tree (binding; do not add extra files as a drive-by)

```text
program/packs/council/
  README.md
  contracts/commissioning.md
  contracts/report.md
  contracts/reconciliation.md
  templates/commissioning.md
  templates/commissioning.schema.toml
  templates/model-report.md
  templates/model-report.schema.toml
  templates/reconciliation.md
  templates/reconciliation.schema.toml
  skills/council/SKILL.md
  skills/second-opinion/SKILL.md
  adapters/cursor.md
  adapters/manual.md
```

`contracts/reconciliation.md` is a **thin pointer**: obey `program/contracts/replication-reconciliation.md` verbatim. List the required H2s (v1 list, including `Retained dissent`). Do not copy the v1 body. Do not fork the v1 file.

## 6. Commissioning / report / reconciliation grammar + provenance

New namespaces are pack-owned. Digits 3. Not stage-scoped. Add them to `program/contracts/identifiers.md` and `program/contracts/naming.md`. Update the link-scan regex to include `CMP|RPT|RCL`.

| Type key | NS | Home | Filename pattern | Digits | Stage-scoped |
| --- | --- | --- | --- | --- | --- |
| commissioning | `CMP` | `reviews/commissioning/` | `CMP-###-slug.md` | 3 | no |
| model-report | `RPT` | `reviews/reports/` | `RPT-###-slug.md` | 3 | no |
| reconciliation | `RCL` | `reviews/reconciliations/` | `RCL-###-slug.md` | 3 | no |

Type keys match template filename stems. `mycelium new commissioning`, `mycelium new model-report`, `mycelium new reconciliation`. There is no `mycelium new council` type.

Heading match is exact and case-sensitive (PHASE-01 H2 rule). Extra H2s are allowed. Check does not grade prose (DEC-005).

### 6.1 Commissioning (CMP)

Required front matter:

| Field | Binding |
| --- | --- |
| `id` | `CMP-###` |
| `title` | string |
| `date` | `YYYY-MM-DD` |
| `rung` | `second-opinion` \| `council` |
| `opt_in` | TOML boolean. **Must be `true`.** Missing or `false` → FAIL. String `"true"` → FAIL. |
| `cost_class` | IFF `rung == "second-opinion"`: must be `cheap`. IFF `rung == "council"`: must be `quick` \| `standard` \| `high-stakes`. Cross-rung values FAIL. |
| `adapter` | `cursor` \| `manual` |

`prompt_sha256` is **not** required on the CMP. **Architect default:** check computes the hash at check time; every matching RPT must repeat the same hex.

Required H2s: `Prompt`, `Attachments`, `Cost`. Bodies are not graded. `Attachments` may be `none`. `Cost` restates the class in prose; check does not parse it.

### 6.2 Prompt identity

```text
prompt_sha256 = hex(sha256(TrimSpace(SectionBody(body, "Prompt"))))
```

- `SectionBody` = bytes after the exact H2 `## Prompt` until the next H2 or EOF (same container rule PHASE-02/03 used).
- `TrimSpace` = surrounding whitespace only.
- Hex = 64 **lowercase** `[0-9a-f]` characters. No `sha256:` prefix. Uppercase FAIL.
- First `## Prompt` wins. Extra `## Prompt` headings do not change the hash.
- Empty-after-trim is structurally legal (hash of empty string). Skill says write a real prompt. Check does not grade it.
- Stdlib only: `crypto/sha256` + `encoding/hex`.

RPT front matter `prompt_sha256` must equal the CMP-computed hex. Mismatch → FAIL. Teaching error names the expected hex, the RPT id, and `program/packs/council/contracts/report.md`.

### 6.3 Model report (RPT)

Required front matter:

| Field | Binding |
| --- | --- |
| `id` | `RPT-###` |
| `title` | string |
| `date` | `YYYY-MM-DD` |
| `model` | string; **not graded**; not an enum; empty FAIL |
| `commissioning` | `CMP-###` that resolves to an existing file |
| `rung` | must equal the CMP's `rung` |
| `adapter` | must equal the CMP's `adapter` |
| `prompt_sha256` | 64 lowercase hex; must equal the CMP-computed hash |

Required H2s: `Position`, `Findings`, `Dissent`. **`## Dissent` must exist.** Body may be `none`. Check does not grade prose (DEC-005).

**Architect default:** check does **not** require distinct `model` strings across RPTs. Same-model twice is a smell for the human, not a FAIL.

### 6.4 Reconciliation (RCL)

Required front matter:

| Field | Binding |
| --- | --- |
| `id` | `RCL-###` |
| `title` | string |
| `date` | `YYYY-MM-DD` |
| `commissioning` | `CMP-###` that resolves |
| `rung` | `council` **only**. `second-opinion` on an RCL → FAIL. |

Second-opinion does **not** require an RCL. An RCL whose CMP is `rung = "second-opinion"` → FAIL.

Required H2s — v1 list, exact strings, including **`Retained dissent`**:

```text
Convergence
Material disagreement
Evidence unique to one report
Contradictory evidence
Different assumptions
Different scope interpretations
Recommendations independently supported
Questions requiring another spike
Final reconciled recommendation
Retained dissent
```

No majority vote / reputation / prose-confidence selectors — stated in the pack reconciliation contract and the council skill. Check does **not** content-score the method.

### 6.5 Rung cardinality

**Architect default — WIP CMP is legal.** Cardinality binds IFF any matching RPT or RCL exists for that CMP. A CMP alone (commissioned, not yet run) passes, provided `opt_in`, `cost_class`, `rung`, and `adapter` are legal.

| CMP `rung` | Matching files | `mycelium check` |
| --- | --- | --- |
| `second-opinion` | no RPT, no RCL | PASS |
| `second-opinion` | exactly one RPT, no RCL | PASS |
| `second-opinion` | two or more RPTs | FAIL |
| `second-opinion` | any RCL | FAIL |
| `council` | no RPT, no RCL | PASS |
| `council` | ≥1 RPT or any RCL, but not (≥2 RPTs **and** exactly one RCL) | FAIL |
| `council` | ≥2 RPTs and exactly one RCL | PASS |
| `council` | ≥2 RPTs and two or more RCLs | FAIL |

"Matching" means the file's `commissioning` field equals the CMP id. Orphan RPT/RCL (missing CMP, or `commissioning` does not resolve) → FAIL.

Prompt identity: every matching RPT's `prompt_sha256` equals the CMP-computed hash.

### 6.6 Seeded dissent (MS-401)

If **any** matching RPT's `## Dissent` section body contains the exact token `SEED-DISSENT`, the matching RCL's `## Retained dissent` section body **must** contain `SEED-DISSENT`. This is a substring presence check, not a quality score.

| Situation | Check |
| --- | --- |
| No RPT contains `SEED-DISSENT` | Rule does not fire |
| RPT contains `SEED-DISSENT`, matching RCL retains it | PASS |
| RPT contains `SEED-DISSENT`, matching RCL lacks it | FAIL |
| `second-opinion` RPT contains `SEED-DISSENT` (no RCL required) | Rule does not require an RCL |

Section body = bytes after the parent H2 until the next H2 or EOF.

### 6.7 Teaching errors (binding shape)

Four lines, stderr, exit 1.

```text
mycelium: CMP-001 opt_in must be true (got false)
convention: council-opt-in
contract: program/packs/council/contracts/commissioning.md
fix: set opt_in = true, or delete the commissioning file

mycelium: CMP-001 cost_class "cheap" is not quick|standard|high-stakes when rung=council
convention: council-cost-class
contract: program/packs/council/contracts/commissioning.md
fix: set cost_class to quick, standard, or high-stakes

mycelium: RPT-001 prompt_sha256 mismatch (want ec87bfc2afd545807ca87b5c29cae8e77262cb3c746fc63e4539d8daeb2a77de)
convention: prompt-identity
contract: program/packs/council/contracts/report.md
fix: set prompt_sha256 to the sha256 hex of CMP-001 ## Prompt (trim surrounding whitespace)

mycelium: CMP-001 council requires >=2 model reports and exactly one reconciliation
convention: council-cardinality
contract: program/packs/council/contracts/commissioning.md
fix: add matching RPT-### files (>=2) and exactly one RCL-###, or remove started reports

mycelium: RCL-001 ## Retained dissent missing SEED-DISSENT
convention: seeded-dissent
contract: program/packs/council/contracts/reconciliation.md
fix: retain the SEED-DISSENT token in ## Retained dissent

mycelium: extra top-level path reviews/ (council pack absent)
convention: extra-top-level
contract: program/contracts/conformance.md
fix: delete reviews/, restore program/packs/council/, or declare extra-top-level:reviews/

mycelium: pack namespace collision: CMP claimed by council and fixture-pack
convention: pack-collision
contract: program/contracts/conformance.md
fix: remove the colliding pack directory
```

### 6.8 Pure parser (Slice 1 / Slice 3)

New package `internal/ladder`. No filesystem. No CLI. Table-driven tests. Check does not call the IFF helpers until Slice 3.

```text
ParseRung(s string) (rung, error)                 # second-opinion | council
ParseAdapter(s string) (adapter, error)           # cursor | manual
ParseCostClass(s string) (class, error)           # cheap | quick | standard | high-stakes
CostClassOK(rung, class) bool                     # IFF table in §6.1
OptInOK(v any) bool                               # TOML boolean true only
PromptSHA256(body string) string                  # hex(sha256(TrimSpace(SectionBody(body, "Prompt"))))
RequiredCMPH2() []string                          # Prompt, Attachments, Cost
RequiredRPTH2() []string                          # Position, Findings, Dissent
RequiredRCLH2() []string                          # v1 list including Retained dissent
Cardinality(rung, nRPT, nRCL) error               # §6.5; n=0,0 is ok (WIP)
SeedDissentOK(rptDissentBodies []string, rclRetained string) error
SectionBody(body, h2) string                      # reuse the PHASE-03 rule
```

New package `internal/pack`. Discovers `program/packs/*/` directories, loads pack schemas, reports collisions, answers `ReviewsAllowed(packs) bool` (true IFF `council` is present).

Do not rewrite `internal/check` in Slice 1 beyond calling `internal/pack` for presence / collision / `reviews/` extra-top-level. Do not bind CMP/RPT/RCL IFF rules until Slice 3.

## 7. Second-opinion move + adapters (Cursor, manual) + panel presets

### 7.1 Second-opinion move

Exactly one different model. Identical commissioning prompt. `rung = "second-opinion"`. `cost_class = "cheap"`. `opt_in = true`. Exactly one matching RPT once any RPT exists. No RCL. Agreement is high-signal; disagreement surfaces a fork (skill prose; check does not score agreement).

Credit **pstack second-opinion doctrine** and **DEC-008** in the skill. Do not copy upstream files into `program/`.

### 7.2 Adapters are markdown procedures, not Go

Do not add `internal/cursor`, an OpenRouter client, or any network caller for ladder rungs. The CLI does not invoke either adapter.

| Adapter | Path | Procedure (encode this; do not invent a second flow) |
| --- | --- | --- |
| Cursor | `program/packs/council/adapters/cursor.md` | Fan out parallel subagents. Each subagent receives the identical `## Prompt` body and the attachment manifest. Each saves one RPT file. Absorb **karpathy/llm-council** three-stage shape as the *recommended procedure*: (1) independent first opinions, (2) anonymized cross-review notes (optional; land in RPT `## Findings` or a second pass — **do not add a fourth artifact type**), (3) synthesis = the RCL, **subordinated to v1 reconciliation**. No chairman-smoothing. Dissent retained. |
| Manual | `program/packs/council/adapters/manual.md` | Paste the identical prompt into N chat UIs. Save N RPT files. Write the RCL by hand. Zero tooling. |

Both adapters credit DEC-008, the v1 replication contract, and karpathy/llm-council. Cursor adapter is first; manual floor always works.

### 7.3 Panel presets

| Rule | Binding |
| --- | --- |
| Path | `~/.config/mycelium/panels.toml` |
| Override | `$MYCELIUM_CONFIG` is a **directory**. File is `$MYCELIUM_CONFIG/panels.toml`. |
| Keys | `quick` / `standard` / `high-stakes`, each with `models = ["..."]` |
| Missing file | **Legal.** |
| Skill defaults | quick=2 models, standard=3, high-stakes=4 |
| Check | Does **not** read user config. Does **not** enforce panel size beyond council ≥2. |
| Tests | **Never** touch the real home directory. Do not write `~/.config/mycelium`. |
| Go reader | **None this phase.** No `internal/panels` package. No CLI flag that loads the file. Example = Appendix D. |

## 8. Commands (existing only; what does not exist)

Exact CLI. No new top-level verb. Flags and exit codes stay PHASE-01/02/03. **No portable council CLI.**

Global: exit 0 success, exit 1 failure. Teaching errors on stderr (four lines). Success text on stdout.

Env (unchanged; MS-401 does not require them):

| Env | Effect this phase |
| --- | --- |
| `MYCELIUM_OFFLINE=1` | Unchanged. Unused by MS-401. |
| `MYCELIUM_NOW` | RFC3339 clock override. **Optional / not required** for MS-401 unless a test pins dates. |
| `MYCELIUM_IDEAS_ROOT` | Unchanged. |
| `MYCELIUM_CONFIG` | Documented directory override for `panels.toml`. Check does **not** read it. Tests must not depend on it. |

### Commands that exist after PHASE-04 (same argv as the pin, plus pack type discovery)

| Command | New? | PHASE-04 delta |
| --- | --- | --- |
| `mycelium new idea …` | no | Slice 4: emit pack + pack skills on new scaffolds |
| `mycelium new <type> <title> [--dir]` | no | Slice 2: discovers pack types when `program/packs/council/` exists |
| `mycelium check [--dir] [--abort-journal]` | no | Slices 1 and 3 bind pack + ladder rules |
| `tier` / `index` / `state` / `wake` / `status` / `publish` / `version` | no | no pack-skill emit; do not rewrite; `handed-off` still refused; no new `index.md` H2 |

### Commands that do not exist this phase

`mycelium council`, `second-opinion` (as a verb), `ladder`, `replicate`, `handoff`, `supersede`, `explore` / `simmer` as separate verbs, `destroy`, `range`, any `new council` type.

Quality should refuse PRs that add them. A PR that adds a `council` top-level verb is a DEC-008 clause-5 violation.

### 8.1 `mycelium new commissioning` (discovered, not a new verb)

```text
mycelium new commissioning <title> [--dir PATH]
```

Available only when the pack is present in the instance `program/`. Allocates the next `CMP-###`, writes `reviews/commissioning/CMP-###-<slug>.md` from the pack template, appends a `new` log line, regenerates `index.md`. Journal `op=new`. Never git-commits. Never calls a model.

**Architect default — emitted defaults** so a lone CMP is a legal second-opinion stub: `rung = "second-opinion"`, `opt_in = true`, `cost_class = "cheap"`, `adapter = "manual"`. Agent edits for council.

Same shape for `mycelium new model-report` and `mycelium new reconciliation`. Templates emit fill-me `commissioning` / `prompt_sha256` / `model` values the agent must edit. Generator does **not** auto-hash the CMP (no new flags; no magic "latest CMP").

When the pack is absent, these type keys are unknown (existing unknown-type teaching error).

### 8.2 `mycelium check` (delta is §9)

```text
mycelium check [--dir PATH] [--abort-journal]
```

No new flags. Same teaching-error shape, `--abort-journal`, and success stdout. Does not read user config. Does not call a model, Cursor, or the network. Instance-root resolution is unchanged from PHASE-01/02/03.

## 9. Check updates (what changes from PHASE-03 conformance)

Structure only (DEC-005). Runtime still reads instance files, never embed. Pin conformance already has **17** must-implement items. Do not renumber 1–17. Add 18–22.

### PHASE-03 items 1–17 stay

Do not rewrite lifecycle storage rules. Do not lift `handed-off` FAIL. Do not change `index.md` required H2s. Do not change log-ops. Do not reopen sparring IFF rules. Do not require a wake brief on instances that never simmered.

### New must-implement items (18–22)

| # | Rule | Slice |
| --- | --- | --- |
| 18 | Pack presence-is-registration: load `program/packs/<name>/` when the directory exists. Namespace, type-key, or home collision between packs (or pack vs core) → FAIL. | Slice 1 |
| 19 | `reviews/` is an allowed top-level path only when `program/packs/council/` is present. Without the pack and without `reviews/`, core check still passes. Without the pack but with `reviews/` leftover: FAIL extra-top-level unless deviation `extra-top-level:reviews/`. | Slice 1 |
| 20 | When the pack is present and CMP files exist: required front matter + H2s per §6.1; `opt_in` must be `true`; `cost_class` IFF `rung`; `adapter` enum. | Schema keys/H2s in Slice 2; IFF in Slice 3 |
| 21 | When the pack is present and RPT files exist: required front matter + H2s including `Dissent`; `commissioning` resolves; `rung`/`adapter` match the CMP; `prompt_sha256` equals the check-computed CMP hash. | Schema keys/H2s in Slice 2; hash/match in Slice 3 |
| 22 | Rung cardinality per §6.5. RCL required H2s including `Retained dissent`. `SEED-DISSENT` substring rule per §6.6. RCL `rung` is `council` only. | Slice 3 |

### What check must not do

| Temptation | Verdict |
| --- | --- |
| Require the council pack | **No.** Core items 1–17 pass without it. |
| Require any CMP / RPT / RCL | **No.** Spark with zero reviews still passes. |
| Require a council to leave `spark` | **No.** |
| Require `## Crux` changes or reopen DEC-007 | **No.** |
| Grade Position / Findings / Dissent / Retained dissent / Prompt prose | **No.** DEC-005. |
| Content-score reports or reconciliation method | **No.** |
| Require distinct `model` strings | **No.** |
| Read `~/.config/mycelium` or `$MYCELIUM_CONFIG` | **No.** |
| Enforce panel size beyond council ≥2 | **No.** |
| Require `panels.toml` to exist | **No.** |
| Call a model / Cursor / network / `gh` / read `GH_TOKEN` | **No.** |
| Add `council` / `replicate` / `ladder` to the log-op regex | **No.** |
| Fail a lone CMP (no RPT yet) | **No.** WIP is legal. |
| Treat OQ-006 as a council | **No.** |
| Add a required `index.md` H2 | **No.** |

### Lift timing

| Slice | Check behavior |
| --- | --- |
| 1 | Pack presence + collision + `reviews/` extra-top-level. **No** CMP/RPT/RCL content checks yet. A CMP file, if someone writes one early, is only subject to extra-top-level / nested-home existence, not hash or cardinality. Slice 1 fixtures do **not** create CMP/RPT/RCL files. |
| 2 | Pack schemas registered. Item 7 (schema-driven front matter + sections) applies to pack types. `opt_in` must be *present* but "must be `true`" and cost-class IFF and hash and cardinality do **not** bind yet. `mycelium new` discovers pack types. |
| 3 | Check calls `internal/ladder`. IFF / hash / cardinality / `SEED-DISSENT` bind. |
| 4 | Skills + adapters; no new check rule. |
| 5 | MS-401 matrix fixtures in `go test ./...` **are** the gate. |

### `program/contracts/conformance.md`

Add items 18–22 so Quality can thermos a numbered list. State the lift timing in the contract the same way PHASE-03 did. Do not renumber items 1–17. Bump the heading "Must-implement checks (17)" to "(22)".

### Link scan

Add `CMP|RPT|RCL` to the naming-contract regex. Homes are the three nested reviews directories. Unresolved `commissioning = "CMP-###"` fails the existing scan once the NS is registered (Slice 2+).

## 10. Skills + adapters + AGENTS.md capability note

Emit path: `.agents/skills/<name>/SKILL.md` like `mycelium-cli`.

Source of truth for pack skills: `program/packs/council/skills/<name>/SKILL.md`. Scaffold copies into the instance when the pack is present in emitted `program/`. Re-run `go generate` in `internal/embed`.

### When pack skills are emitted

| Event | Emit `council` + `second-opinion`? |
| --- | --- |
| New `mycelium new idea` scaffold (PHASE-04+) with pack in emitted `program/` | **Yes** |
| New scaffold after the human deleted the pack from master's embed (do not) | n/a — master embed keeps the pack |
| Instance whose `program/packs/council/` was deleted | Skills already on disk under `.agents/` may remain (`.agents/` extra files are allowed). Types and pack checks are gone. |
| `mycelium tier` / `index` / `state` / `wake` / `check` | **No** (do not retrofit) |
| One-shot pack-skills copy command | **Not a PHASE-04 command** |

**Architect default:** do not rewrite `spark` / `wake` / `portfolio` / `thinking`. Existing PHASE-01–03 instances: re-scaffold or copy the pack **manually**. Document that sentence in `program/skeleton/AGENTS.md` and in `mycelium-cli` SKILL.md.

### `program/packs/council/skills/council/SKILL.md`

Front matter `name: council`. Binding procedure:

1. Read `index.md` and `CONTEXT.md` (not the whole tree).
2. Suggest a council only when a v1 replication trigger fires. State panel size and `cost_class` first.
3. `mycelium new commissioning "<title>"`. Set `rung = "council"`, `opt_in = true`, `cost_class` to `quick` / `standard` / `high-stakes`, `adapter` to `cursor` or `manual`.
4. Fill `## Prompt` (identical for every model), `## Attachments`, `## Cost`.
5. Follow `adapters/cursor.md` or `adapters/manual.md`. Do not ask the CLI to call a model.
6. `mycelium new model-report` once per model. Set `commissioning`, `rung`, `adapter`, `model`, and `prompt_sha256` to the sha256 hex of the trimmed `## Prompt` body. Fill `## Position`, `## Findings`, `## Dissent` (body may be `none`).
7. `mycelium new reconciliation`. Set `commissioning` and `rung = "council"`. Fill every v1 H2, including `## Retained dissent`. Do **not** choose by majority vote, length, confidence of prose, or model reputation.
8. Run `mycelium check` before handing back. On `prompt_sha256` mismatch the teaching error names the expected hex.
9. Do not `git commit` unless the human asks. The CLI never commits.
10. If the runtime cannot fan out, skip this rung; sparring still applies.

Credits (must appear in the skill body): DEC-008; v1 `program/contracts/replication-reconciliation.md`; karpathy/llm-council (three-stage shape, subordinated); pstack second-opinion doctrine (for the neighboring move). Do not copy upstream files into `program/`.

What the skill must not say: do not tell the agent to run `mycelium council` / `ladder` / `replicate` / `handoff` (those verbs do not exist). Do not tell the agent that CI / an Actions job is the done bar. Do not tell the agent to flip `handed-off`.

### `program/packs/council/skills/second-opinion/SKILL.md`

Front matter `name: second-opinion`. Same surface (`mycelium new commissioning` + one `model-report`). `rung = "second-opinion"`. `cost_class = "cheap"`. Exactly one RPT. No RCL. Credit pstack second-opinion doctrine and DEC-008. Same "do not run `mycelium council`" ban.

### `program/skills/mycelium-cli/SKILL.md` update

- Add a paragraph: ladder surface is `mycelium new commissioning`, `mycelium new model-report`, `mycelium new reconciliation`, `mycelium check`, plus pack skills `council` and `second-opinion`, when `program/packs/council/` is present. **No portable council CLI.** **No** `council` / `ladder` / `replicate` verb.
- Capability note: runtimes that cannot fan out skip rungs 2–3; sparring still applies.
- Extend the emit list: new scaffolds emit `.agents/skills/{mycelium-cli,spark,wake,portfolio,thinking,council,second-opinion}/SKILL.md` when the pack is present.
- Keep the manual floor, teaching-error shape, and "do not git commit unless the human asks".

### `program/skeleton/AGENTS.md` update

- Name the `council` and `second-opinion` skills in the Skills section emit list.
- **Capability note (verbatim intent):** runtimes that cannot fan out skip rungs 2–3; sparring still applies.
- Same no-retrofit sentence.
- Do not list `mycelium council` / `ladder` / `replicate` / `handoff` as commands.

## 11. Templates / schema deltas + pack file tree

Pack templates live under `program/packs/council/templates/`. Do **not** add CMP/RPT/RCL schemas to `program/templates/` (core). Discovery is the point of presence-is-registration.

**Architect default — discovery.** `internal/schema` and `internal/generate` scan:

```text
program/templates/*.schema.toml
program/packs/*/templates/*.schema.toml
```

Type key = filename stem (`commissioning.schema.toml` → `commissioning`). When the pack directory is absent, pack type keys are unknown. Do not hardcode `commissioning` / `model-report` / `reconciliation` in core.

### `commissioning.schema.toml`

```text
namespace = "CMP"
home = "reviews/commissioning"
filename_pattern = "CMP-{NNN}-{slug}.md"
stage_scoped = false
digits = 3
required_front_matter = ["id", "title", "date", "rung", "opt_in", "cost_class", "adapter"]
required_sections = ["Prompt", "Attachments", "Cost"]

[enums.rung]
values = ["second-opinion", "council"]

[enums.adapter]
values = ["cursor", "manual"]

[enums.cost_class]
values = ["cheap", "quick", "standard", "high-stakes"]
```

`opt_in` is a required key, not an enum. "Must be TOML `true`" binds in Slice 3. Cost-class IFF binds in Slice 3. Schema enum lists all four `cost_class` values so Slice 2 can accept either rung's legal value.

### `model-report.schema.toml`

```text
namespace = "RPT"
home = "reviews/reports"
filename_pattern = "RPT-{NNN}-{slug}.md"
stage_scoped = false
digits = 3
required_front_matter = ["id", "title", "date", "model", "commissioning", "rung", "adapter", "prompt_sha256"]
required_sections = ["Position", "Findings", "Dissent"]

[enums.rung]
values = ["second-opinion", "council"]

[enums.adapter]
values = ["cursor", "manual"]
```

### `reconciliation.schema.toml`

```text
namespace = "RCL"
home = "reviews/reconciliations"
filename_pattern = "RCL-{NNN}-{slug}.md"
stage_scoped = false
digits = 3
required_front_matter = ["id", "title", "date", "commissioning", "rung"]
required_sections = [
  "Convergence",
  "Material disagreement",
  "Evidence unique to one report",
  "Contradictory evidence",
  "Different assumptions",
  "Different scope interpretations",
  "Recommendations independently supported",
  "Questions requiring another spike",
  "Final reconciled recommendation",
  "Retained dissent",
]

[enums.rung]
values = ["council"]
```

Do not invent a schema DSL for cardinality or hash IFF. Those live in pack contracts and `internal/ladder`.

### Emitted template shapes

**Architect default — CMP** (tokens filled by the existing generator):

```text
+++
id = "{{ID}}"
title = "{{TITLE}}"
date = "{{DATE}}"
rung = "second-opinion"
opt_in = true
cost_class = "cheap"
adapter = "manual"
+++

# {{ID}} — {{TITLE}}

<!-- slug: {{SLUG}} -->

## Prompt

<!-- fill -->

## Attachments

none

## Cost

cheap
```

**Architect default — RPT:** emit `rung = "second-opinion"`, `adapter = "manual"`, `commissioning = "CMP-000"`, `prompt_sha256 = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"` (sha256 of empty) as visible fill-me values, plus H2s Position / Findings / Dissent (`none` allowed). `mycelium new model-report` tests assert the file and headings exist, not that an unedited RPT passes instance check. MS-401 fixtures **edit** the fields.

**Architect default — RCL:** emit `rung = "council"`, `commissioning = "CMP-000"`, all ten v1 H2s with `<!-- fill -->`. Same edit-before-pass rule.

### Core templates

No delta to `question.md`, `decision.md`, or `CONTEXT.md`. Do not reopen PHASE-03.

### Pack README + contracts

`program/packs/council/README.md`: this is the council pack; presence-is-registration; drop in to enable; delete the directory to disable; no portable council CLI.

Contracts: commissioning.md + report.md encode §6. reconciliation.md points at `program/contracts/replication-reconciliation.md` verbatim, lists the ten H2s, states RCL is council-only, states the `SEED-DISSENT` substring rule, restates the banned selectors.

## 12. Vertical slices 0–5 with build order, each checkable, go test / files only

PR-per-slice, sequential, rebase on main. Arvo merges Quality-green PRs. Engineering opens PRs. Engineering does NOT push to main. Prefer one live PR at a time. Do not stack unpublished slices on one branch unless Quality is backed up.

Each PR title: `PHASE-04 Slice N: <done-bar noun>`. Each PR body links this brief and the slice done bar. No drive-by refactors. No v1 deletions. No PHASE-05+ commands. No Actions job as a done bar. Do **not** add `.github/workflows/phase-04-*.yml`.

### Slice 0 — Commissioning (docs only)

This brief + acceptance stub + pack contract stubs. No product code. No Go.

Land: `framework/phases/PHASE-04-implementation-brief.md`, `framework/phases/PHASE-04-acceptance.md` (rows = §15), pack contract files under `program/packs/council/contracts/` (may be docs-only stubs that already state the pointer to the v1 replication contract), updates to `conformance.md` (items 18–22 + lift timing; may finish in Slices 1–3), identifiers + naming rows (may finish in Slice 2).

Done: files exist on a PR. Quality reads them against this brief. No product code.

### Slice 1 — Pack presence + collision + enable/disable (no CMP/RPT checks yet)

- `internal/pack` discovers `program/packs/*/`. Collision rule. `ReviewsAllowed`.
- `internal/check` calls it. Item 18–19 bind.
- `program/packs/council/` directory exists on master (README + empty-enough tree so presence is real). Embed generate so new scaffolds receive it.
- Fixtures do **not** create CMP/RPT/RCL files. Hash / cardinality / `opt_in` IFF do **not** bind.

Done (hermetic `go test`): pack + no `reviews/` → 0; no pack + no `reviews/` → 0; no pack + `reviews/` → 1; same + `extra-top-level:reviews/` → 0; fabricated second pack claiming `CMP` → 1.

### Slice 2 — Templates / schemas + `mycelium new` discovers pack types

- Pack templates + sidecar schemas as in §11.
- Generate/schema scan `program/packs/*/templates/*.schema.toml`.
- `mycelium new commissioning` / `model-report` / `reconciliation` work when the pack is present; unknown when absent.
- Identifiers + naming + link regex gain `CMP|RPT|RCL`.
- Item 7 applies to pack types (required keys + H2s). Slice 3 IFF not bound: `opt_in = false` may still pass in this slice.

Done (hermetic `go test`): pack present → `mycelium new commissioning "SQLite"` writes `reviews/commissioning/CMP-001-*.md` with the §11 defaults; pack absent → type unknown (exit 1); three pack schemas parse.

### Slice 3 — Check binds provenance / hash / cardinality / opt-in / cost_class / SEED-DISSENT

- `internal/ladder` parser + table tests (§6.8).
- `internal/check` calls it for each CMP/RPT/RCL when the pack is present.
- Items 20–22 bind. Teaching errors match §6.7.

Done (hermetic `go test`): `opt_in = false` → 1; council `cost_class = "cheap"` → 1; second-opinion + two RPTs → 1; council + one RPT → 1; council + ≥2 RPTs + one RCL → 0; hash mismatch → 1 (stderr names expected hex); `SEED-DISSENT` in RPT not RCL → 1; lone CMP → 0; no pack, no `reviews/` → 0.

### Slice 4 — Skills + adapters + AGENTS capability note + mycelium-cli

- Pack skills + adapters as in §7 and §10.
- Update `mycelium-cli` and `AGENTS.md`.
- `go generate` in `internal/embed`.
- New scaffolds emit `.agents/skills/{council,second-opinion}/SKILL.md` when the pack is present.
- `tier` / `index` / `state` / `wake` do not retrofit.

Done (hermetic `go test`): new scaffold `--offline` has the pack directory and both pack skills; a fixture that lacked them still lacks them after `tier` / `index`. Credits present (string asserts; do not grade prose). Capability note present in AGENTS.md. Skills do not mention a `mycelium council` command as something to run, except to say it does not exist.

### Slice 5 — MS-401 matrix fixtures in `go test ./...`

Hermetic fixture recipe = Appendix E. Temp instances. `go test ./...` only. No network. No `gh`. No `GH_TOKEN`. Do **not** add `.github/workflows/phase-04-*.yml`. Do **not** extend `phase-01-hermetic.yml` as a phase gate. Quality should refuse a PR that adds an Actions job as the MS-401 gate.

Done: `go test ./...` runs the MS-401 matrix green. That **is** the gate. A live Cursor council run is **not** the gate.

## 13. Done / verified mapped onto MS-401

MS-401 is the hermetic phase gate (perspective-ladder acceptance matrix). **No live Cursor run as the gate.** **No Actions job.** Robert waived CI. Done bar is hermetic `go test ./...` as far as it can go.

Blueprint MS-401 (do not expand): the end-to-end council run is **one row** of a perspective-ladder acceptance matrix covering the DEC-008 contract — second-opinion, Cursor-council, and manual-floor rows; explicit opt-in and stated cost class; prompt identity and model provenance evidenced; independent per-model reports; a seeded dissent surviving reconciliation; council-pack enable/disable without touching core checks.

### MS-401 expected (authoritative; recipe in Appendix E)

| Row | Setup | `mycelium check` |
| --- | --- | --- |
| second-opinion | pack present; CMP `rung=second-opinion` `opt_in=true` `cost_class=cheap` `adapter=manual`; exactly one matching RPT with matching `prompt_sha256`; no RCL | exit 0 |
| Cursor-council | pack present; CMP `rung=council` `opt_in=true` `cost_class=standard` `adapter=cursor`; ≥2 RPTs (independent files, `model` set); exactly one RCL; one RPT `## Dissent` contains `SEED-DISSENT`; RCL `## Retained dissent` contains `SEED-DISSENT`; hashes match | exit 0 |
| manual-floor | same council shape with `adapter=manual` | exit 0 |
| opt-in / cost_class negatives | `opt_in=false`; council `cost_class=cheap`; second-opinion `cost_class=standard` | exit 1 each |
| prompt identity / provenance | RPT `prompt_sha256` ≠ CMP-computed hash; RPT `rung` ≠ CMP `rung` | exit 1 each |
| independent per-model reports | council with two RPT files (not one file with two headings) | exit 0 on the happy council row |
| seeded dissent | RPT has `SEED-DISSENT`; RCL `## Retained dissent` lacks it | exit 1 |
| enable/disable | pack present, no `reviews/` → 0; pack deleted, no `reviews/` → 0 (core still passes); pack deleted, `reviews/` leftover → 1 | as table |

`gh` never invoked. No `GH_TOKEN`. No Actions job. `MYCELIUM_NOW` is optional / not required. Tests never touch the real home directory. OQ-006 is not a row.

### Slice → MS-401 map

| Slice | MS-401 clause |
| --- | --- |
| 0 | commissioning; not a runtime clause |
| 1 | enable/disable + collision; no CMP/RPT yet |
| 2 | pack types exist so fixtures can be generated |
| 3 | provenance / hash / cardinality / opt-in / cost_class / SEED-DISSENT |
| 4 | skills + adapters + capability note; not a check clause |
| 5 | the matrix in `go test ./...` **is** the gate |

PHASE-04 is accepted when MS-401 is green in `go test ./...` on main. Arvo accepts the phase. Engineering does not self-accept. A live Cursor council run may be attached as human evidence; it does not replace the gate.

## 14. Automated test plan

Engineering MUST write these tests. Quality thermos against this list. Do NOT require Playwright, Docker, live GitHub, `GH_TOKEN`, a live Cursor fan-out, or an Actions job in default `go test ./...`.

### Unit (no network, no gh, no home directory)

| Area | Cases |
| --- | --- |
| pack discover | one pack; zero packs; collision on NS / type key / home |
| ReviewsAllowed | council present → true; absent → false |
| ParseRung / ParseAdapter / ParseCostClass | legal tokens; empty; `Second-Opinion`; trailing space |
| CostClassOK | second-opinion+cheap ok; second-opinion+quick fail; council+cheap fail; council+standard ok |
| OptInOK | `true` ok; `false` fail; `"true"` fail; missing fail |
| PromptSHA256 | Appendix B body → `ec87bfc2afd545807ca87b5c29cae8e77262cb3c746fc63e4539d8daeb2a77de`; surrounding whitespace ignored; empty → `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| Cardinality | §6.5 table |
| SeedDissentOK | token in RPT+RCL ok; token in RPT only fail; no token ok |
| SectionBody | first H2 wins; H2 after next H2 does not count |
| schema parse | three pack schemas |

### Hermetic CLI (built binary, temp dirs)

| Case | Expect |
| --- | --- |
| pack present, no `reviews/` | `mycelium check` exits 0 |
| pack absent, no `reviews/` | check exits 0 |
| pack absent, `reviews/` leftover | check exits 1; names `reviews/` and extra-top-level |
| pack absent, `reviews/` + `extra-top-level:reviews/` | check exits 0 |
| fabricated second pack claiming `CMP` | check exits 1; names collision |
| `mycelium new commissioning "SQLite"` with pack | file exists; defaults in §11 |
| same type with pack deleted | exit 1; unknown type |
| `mycelium new model-report` / `reconciliation` with pack | files exist under the nested homes |
| lone CMP (legal defaults) after Slice 3 | check exits 0 |
| second-opinion complete (Appendix B) | check exits 0 |
| second-opinion + two RPTs | check exits 1 |
| council complete with `SEED-DISSENT` (Appendix C) | check exits 0 |
| council missing RCL | check exits 1 |
| council one RPT | check exits 1 |
| `opt_in = false` | check exits 1 |
| council `cost_class = "cheap"` | check exits 1 |
| hash mismatch | check exits 1; stderr names expected hex |
| `SEED-DISSENT` in RPT, absent in RCL | check exits 1 |
| new scaffold emits pack + pack skills | files present |
| `tier` / `index` do not retrofit pack skills | absence preserved |
| MS-401 | Appendix E; `go test ./...` |
| no `gh` | fake runner or no execrun use; `gh` never invoked |
| no home dir | tests do not read or write `~/.config/mycelium` |

Do not require live GitHub or live Cursor for `go test ./...`. Do not commission a `GH_TOKEN` job. Do not add an Actions job as the MS-401 gate.

## 15. Acceptance matrix / in-repo contract paths (the perspective-ladder matrix)

Slice 0 lands the paths listed in §12 Slice 0. Later slices also land: `program/packs/council/` tree; identifiers + naming + conformance deltas; `internal/pack`; `internal/ladder`; generate/schema/check extensions; updated `program/skills/mycelium-cli/SKILL.md`; updated `program/skeleton/AGENTS.md`; `internal/clitest/ms401_hermetic_test.go`. No workflow file. No DEC-015 file. No `mycelium council` command package.

In-repo contract paths:

```text
program/packs/council/contracts/commissioning.md
program/packs/council/contracts/report.md
program/packs/council/contracts/reconciliation.md
program/contracts/replication-reconciliation.md
program/contracts/conformance.md
program/contracts/identifiers.md
program/contracts/naming.md
```

### Acceptance matrix rows (copy into `PHASE-04-acceptance.md`)

Each row: id, check, evidence, owner (Engineering | Arvo). **CI is not an owner.** Robert waived CI.

| id | check | evidence |
| --- | --- | --- |
| A-S0 | Slice 0 files exist and match this brief | PR diff; Quality read |
| A-S1 | pack presence; collision; enable/disable; `reviews/` extra-top-level | hermetic `go test` |
| A-S2 | pack schemas; `mycelium new` discovers pack types; unknown when pack absent | hermetic `go test` + file read |
| A-S3 | provenance / hash / cardinality / `opt_in` / `cost_class` / `SEED-DISSENT` | hermetic `go test` |
| A-S4 | pack skills + adapters + capability note; mycelium-cli names the surface; no `mycelium council` verb to run | hermetic `go test` + file read |
| A-S5 | MS-401 matrix green | `go test ./...` |
| MS-401-SO | second-opinion row | A-S5 |
| MS-401-CUR | Cursor-council row (seeded files, `adapter=cursor`; no live Cursor) | A-S5 |
| MS-401-MAN | manual-floor row | A-S5 |
| MS-401-OPT | explicit opt-in and stated cost class (happy + negatives) | A-S5 |
| MS-401-HASH | prompt identity and model provenance evidenced | A-S5 |
| MS-401-IND | independent per-model reports (two RPT files) | A-S5 |
| MS-401-DIS | seeded dissent surviving reconciliation | A-S5 |
| MS-401-PKG | council-pack enable/disable without touching core checks | A-S5 / A-S1 |

No live-Cursor-gate row. No Actions-job row. Quality should refuse a PR that adds an Actions job as the MS-401 gate. Do not cite OQ-006 as a matrix row.

## 16. Decided / Architect defaults

No TBD. Open items are decided inline and labeled **Architect default**. Do not bikeshed them in a code PR. Write a DEC to change one. No DEC-015 is required for these. DEC-008 already exists.

Index of defaults that are easy to miss (details live in §§4–11; this list is the bikeshed lock):

- **No new top-level CLI verb.** **No portable council CLI.** Surface = existing `new <type>` / `check` + pack skills. CLI never calls a model, Cursor, or the network for ladder rungs.
- OQ-003: **only `council` is a pack.** §5 not-packs list is closed. presence-is-registration; no registry; collision rule still written. Embed includes `program/packs/council/`. Disable = delete that directory. Core checks must not require the pack.
- Namespaces `CMP` / `RPT` / `RCL`; homes nested under `reviews/`; digits 3; not stage-scoped; type keys `commissioning` / `model-report` / `reconciliation`. `reviews/` allowed top-level only when the pack is present (else FAIL unless `extra-top-level:reviews/`).
- CMP / RPT / RCL grammar, `opt_in`, `cost_class` IFF, `prompt_sha256` (check computes; RPT repeats), cardinality, and `SEED-DISSENT` are §6. Adapters are markdown (§7). Panel file optional; check does not read user config; no Go reader; tests never touch the real home directory.
- Skills in the pack; new scaffolds only. Capability note: runtimes that cannot fan out skip rungs 2–3; sparring still applies. Credits: DEC-008, v1 replication contract, karpathy/llm-council, pstack second-opinion doctrine. Do not copy upstream files into `program/`.
- Discovery scans `program/templates/*.schema.toml` and `program/packs/*/templates/*.schema.toml`. Nested `home` legal. New packages: `internal/pack`, `internal/ladder` only. Do not touch `internal/slug` (DEC-014). Do not fork the v1 replication contract. Sparring (rung 1) stays. `methodology_version` stays `2.0.0`. CLI version stays `0.1.0-dev` unless already different. No new log op.
- Lift: Slice 1 presence/collision; Slice 2 templates + `new`; Slice 3 IFF bind; Slice 4 skills/adapters/AGENTS; Slice 5 MS-401 in `go test ./...`. Gate is hermetic `go test`. Quality should refuse an Actions job as the MS-401 gate; Quality should **not** refuse a missing PHASE-04 workflow. Do not add `.github/workflows/phase-04-*.yml`. Do not extend `phase-01-hermetic.yml` as a phase gate.
- Do not commission PHASE-05–06. `handed-off` stays unreachable. No DEC-015. Do not reopen DEC-012 / DEC-013 / DEC-014.

## 17. Risks, rollback, what Quality should refuse

### Risks and mitigations

| Risk | Mitigation |
| --- | --- |
| Portable council CLI / model-invoking CLI | DEC-008 clause 5. §8. Quality should refuse. |
| Extra packs beyond council | §5 OQ-003 lock. Quality should refuse. |
| Forking the v1 replication contract | Thin pointer. Quality should refuse. |
| Majority-vote selector in RCL skill or check | Banned in contract + skill. Check does not content-score. Quality should refuse. |
| Requiring a council to leave spark / to pass core check | Enable/disable fixtures. Quality should refuse. |
| Content-scoring reports | DEC-005. Quality should refuse. |
| Actions job added as the MS-401 gate | Robert waived CI. Quality should refuse. Absence of a PHASE-04 workflow is correct. |
| `GH_TOKEN` job / MS-101(b) reopened | Leftover stays leftover. Quality should refuse. |
| Check reads `~/.config/mycelium` or tests touch home | §7.3. Quality should refuse. |
| Live Cursor run treated as the gate | §13. Quality should refuse. |
| Growing `latinFold` or adding `x/text` | DEC-014. Do not touch `internal/slug`. |
| Emitting `framework/` or deleting Justfile | Absence tests stay. Master v1 stays. |
| CLI git-commits instance work product | DEC-010. Quality should refuse. |
| Flipping `handed-off` to legal | Command refuse + check FAIL stay. |
| Reopening DEC-007 grammar | PHASE-03 stays. Quality should refuse. |
| Chairman-smoothing in the Cursor adapter | Adapter subordinated to v1. Dissent retained. |

### Rollback

Revert the offending PR on master. Do not `git push --force` to main. Floor is `d520560d138937a8904c43623524d2d2b9ef82f3`. Do not delete Justfile/scripts as a "cleanup" rollback.

### Quality should refuse

Refuse to approve if:

- an Actions job is added as the MS-401 gate, or `.github/workflows/phase-04-*.yml` is added, or `phase-01-hermetic.yml` is extended as a phase gate (Quality should **not** refuse a missing PHASE-04 workflow — absence is correct)
- a new CLI verb is shipped (`mycelium council`, `second-opinion`, `ladder`, `replicate`, `handoff`, `supersede`) or the CLI calls a model / Cursor / the network for ladder rungs
- extra packs beyond `council` appear
- `program/contracts/replication-reconciliation.md` is forked or rewritten
- a majority-vote / reputation / prose-confidence selector is implemented
- check requires a council (or any CMP) to leave `spark` or to pass core check
- Position / Findings / Dissent / Retained dissent / prompts are content-scored
- a `GH_TOKEN` job is commissioned as a gate, or MS-101(b) is implemented
- `framework/` is emitted into an instance
- CLI git-commits instance work product
- `latinFold` grows or NFKD is implemented (DEC-014)
- cobra / viper / yaml / testify / go-github / `golang.org/x/text` appears
- Justfile/scripts deleted from master, `just init` was run on master, or `research-program.toml` was renamed
- DEC-012 / DEC-013 / DEC-014 reopened in a code PR
- DEC-015 is recorded
- PHASE-05 / PHASE-06 commands appear, or `handed-off` succeeds without a packet, or check stops failing stored `handed-off`
- user-config is required to exist, or check reads `~/.config/mycelium`, or tests touch the real home directory
- a live Cursor run is required as the MS-401 gate
- hermetic tests call network or real `gh`
- PR pushed straight to main
- DEC-007 sparring grammar is reopened

## 18. Execution order (PR-per-slice)

Same order as §12 (slices 0→5). PR-per-slice, sequential, rebase on main. One live PR at a time. Slice 3 is the check bind — do not combine it with 4–5. Slice 5 must be green in `go test ./...` on its PR (not Actions).

Title: `PHASE-04 Slice N: <done-bar noun>`. Body links this brief and the slice done bar. No drive-by refactors, v1 deletions, PHASE-01–03 leftover work, or PHASE-05+ commands. Engineering opens PRs; Arvo merges Quality-green PRs; Engineering does NOT push to main.

Cursor cloud env name is exactly `robertguss/mycelium`. Go 1.26. `CGO_ENABLED=0`.

## 19. Handoff

### What Engineering starts with

This file. Only this file. Start from `https://github.com/robertguss/mycelium` at `d520560d138937a8904c43623524d2d2b9ef82f3` (PHASE-01–03 accepted; MS-301 on main via PR #40). Read `framework/blueprint.md` and DEC-001–014 for authority, not for a second plan. Execute Slice 0 first. Do not implement from a later SHA unless Arvo re-pins. Cloud env `robertguss/mycelium`. Go 1.26. `CGO_ENABLED=0`.

### What Engineering must not do

See §17 (Quality should refuse) and §16 (Architect defaults). Do not open a design debate in the PR. Do not write a second brief. Do not write DEC-015. Do not commission PHASE-05–06. Do not add an Actions job as the MS-401 gate. Do not start from a later SHA and "fast-forward" this brief. Do not clone a later SHA "to save time." Do not invent a `mycelium council` command.

### What Quality reads

This brief, the acceptance matrix, the pack contracts, the conformance delta, and the PR diff. Thermos: §14 tests exist and match; §17 refuse list is clean; MS-401 hermetic `go test ./...`; no `GH_TOKEN` gate; no Actions job as the MS-401 gate; no new CLI verb; no extra packs; v1 replication contract unforked.

Quality should refuse a PR that adds an Actions job as the MS-401 gate; Quality should **not** refuse a missing PHASE-04 workflow.

### What Arvo does

Merges Quality-green PRs. Accepts PHASE-04 when MS-401 is green on main. May attach a live Cursor council run as human evidence; must not treat it as the gate. Does not re-pin without writing the new SHA. Does not record DEC-015.

## Appendix A — No new DEC (DEC-008 already exists)

No DEC-015. PHASE-04 implements **DEC-008** (Accepted 2026-08-14). PHASE-04 does not reopen DEC-012, DEC-013, or DEC-014. Remaining choices are Architect defaults in §16. Engineering lands **zero** new files under `framework/decisions/`.

If a later discovery contradicts a locked item, stop and write a DEC; do not silently patch this brief in a code PR.

## Appendix B — Second-opinion CMP + RPT example

Fixture shape after `mycelium new commissioning "SQLite second opinion"` plus edits, then `mycelium new model-report "Model A"`. Tests must **not** grade the words — only headings, front matter, hash, cardinality, and `mycelium check` exit 0.

`## Prompt` body (trimmed) hashes to
`ec87bfc2afd545807ca87b5c29cae8e77262cb3c746fc63e4539d8daeb2a77de`.

`reviews/commissioning/CMP-001-sqlite-second-opinion.md`:

```text
+++
id = "CMP-001"
title = "SQLite second opinion"
date = "2026-08-15"
rung = "second-opinion"
opt_in = true
cost_class = "cheap"
adapter = "manual"
+++

# CMP-001 — SQLite second opinion

<!-- slug: sqlite-second-opinion -->

## Prompt

Should this idea use SQLite as the store? Answer independently. Do not see other reports.

## Attachments

none

## Cost

cheap
```

`reviews/reports/RPT-001-model-a.md`:

```text
+++
id = "RPT-001"
title = "Model A"
date = "2026-08-15"
model = "model-a"
commissioning = "CMP-001"
rung = "second-opinion"
adapter = "manual"
prompt_sha256 = "ec87bfc2afd545807ca87b5c29cae8e77262cb3c746fc63e4539d8daeb2a77de"
+++

# RPT-001 — Model A

<!-- slug: model-a -->

## Position

Use SQLite.

## Findings

One file. Enough.

## Dissent

none
```

No RCL. `mycelium check` exits 0. Negative in the same test file: add a second RPT matching CMP-001 → exit 1. Hash-mismatch negative: flip one hex nibble → exit 1; stderr names the expected hex.

## Appendix C — Council CMP + two RPTs + RCL with SEED-DISSENT

Cursor-council row (`adapter = "cursor"`). Manual-floor row is the same files with `adapter = "manual"` on the CMP and both RPTs. `## Prompt` body (trimmed) hashes to
`8997334f7f2f0bf821bce8ccc4a8d6cf027317c6c66d821200a032a6a11ce098`.

`reviews/commissioning/CMP-001-sqlite-council.md`:

```text
+++
id = "CMP-001"
title = "SQLite council"
date = "2026-08-15"
rung = "council"
opt_in = true
cost_class = "standard"
adapter = "cursor"
+++

# CMP-001 — SQLite council

<!-- slug: sqlite-council -->

## Prompt

Review the SQLite store decision. Work independently. Do not see other reports. Retain dissent.

## Attachments

none

## Cost

standard
```

`reviews/reports/RPT-001-model-a.md`:

```text
+++
id = "RPT-001"
title = "Model A"
date = "2026-08-15"
model = "model-a"
commissioning = "CMP-001"
rung = "council"
adapter = "cursor"
prompt_sha256 = "8997334f7f2f0bf821bce8ccc4a8d6cf027317c6c66d821200a032a6a11ce098"
+++

# RPT-001 — Model A

<!-- slug: model-a -->

## Position

Use SQLite.

## Findings

One file. Enough.

## Dissent

SEED-DISSENT
```

`reviews/reports/RPT-002-model-b.md`: same shape with `id = "RPT-002"`, `title = "Model B"`, `model = "model-b"`, `## Position` "Do not use SQLite.", `## Dissent` `none`. Same `prompt_sha256`, `commissioning`, `rung`, `adapter`.

`reviews/reconciliations/RCL-001-sqlite-council.md`:

```text
+++
id = "RCL-001"
title = "SQLite council reconciliation"
date = "2026-08-15"
commissioning = "CMP-001"
rung = "council"
+++

# RCL-001 — SQLite council reconciliation

<!-- slug: sqlite-council-reconciliation -->

## Convergence

Both reports address the store.

## Material disagreement

Whether SQLite is enough.

## Evidence unique to one report

none

## Contradictory evidence

none

## Different assumptions

Single-process vs later writers.

## Different scope interpretations

none

## Recommendations independently supported

Spike if a second writer appears.

## Questions requiring another spike

none

## Final reconciled recommendation

Use SQLite now. Revisit on a second writer.

## Retained dissent

SEED-DISSENT
```

`mycelium check` exits 0. Negative: delete `SEED-DISSENT` from `## Retained dissent` → exit 1. Negative: delete RCL → exit 1. Negative: delete RPT-002 → exit 1. Bodies under the H2s are not graded except the `SEED-DISSENT` substring.

These are **seeded files**. They are not a live Cursor run. Do not cite OQ-006 as this row.

## Appendix D — `panels.toml` example

Missing file is legal. Check does not read this file. Tests never write it under the real home directory. Skill documents the defaults (quick=2, standard=3, high-stakes=4).

`$MYCELIUM_CONFIG` is a directory override; the file is then `$MYCELIUM_CONFIG/panels.toml`. Otherwise `~/.config/mycelium/panels.toml`.

```text
[quick]
models = ["model-a", "model-b"]

[standard]
models = ["model-a", "model-b", "model-c"]

[high-stakes]
models = ["model-a", "model-b", "model-c", "model-d"]
```

No Go reader this phase. Example is documentation.

## Appendix E — MS-401 fixture recipe (all matrix rows)

Hermetic. No network. No `gh`. No `GH_TOKEN`. `go test ./...` only. Do **not** add an Actions job or `.github/workflows/phase-04-*.yml`. Do **not** extend `phase-01-hermetic.yml` as a phase gate. Tests never touch the real home directory. Temp dirs; freshly built binary; stdlib edits. Do not add `mycelium edit`. Do not invoke Cursor.

**Architect default:** matrix test lives at `internal/clitest/ms401_hermetic_test.go` (execs the binary). Parser tests live in `internal/ladder` and `internal/pack`.

### Enable / disable (also Slice 1)

```text
mycelium new idea "Pack On" --offline --dir PATH/pack-on
# emitted program/packs/council/ present; no reviews/
mycelium check --dir PATH/pack-on
# exit 0

# copy or scaffold, then DELETE program/packs/council
mycelium check --dir PATH/pack-off
# exit 0   (core checks still pass)

# pack-off + mkdir reviews/
mycelium check --dir PATH/pack-off-reviews
# exit 1   (extra-top-level)

# same + mycelium.toml deviation extra-top-level:reviews/
mycelium check --dir PATH/pack-off-reviews-dev
# exit 0
```

Collision: in a temp instance that has the council pack, write `program/packs/fixture-pack/templates/collide.schema.toml` claiming `namespace = "CMP"`. `mycelium check` exits 1.

### second-opinion row

```text
mycelium new idea "SO Fixture" --offline --dir PATH/so-fixture
mycelium new commissioning "SQLite second opinion" --dir PATH/so-fixture
mycelium new model-report "Model A" --dir PATH/so-fixture
# EDIT CMP + RPT to Appendix B
mycelium check --dir PATH/so-fixture
# exit 0
```

Negatives in the same test file: second RPT → exit 1; hash nibble flip → exit 1; `opt_in = false` → exit 1; `cost_class = "standard"` on second-opinion → exit 1.

### Cursor-council row (seeded files; not a live Cursor run)

```text
mycelium new idea "Council Cursor" --offline --dir PATH/council-cursor
mycelium new commissioning "SQLite council" --dir PATH/council-cursor
mycelium new model-report "Model A" --dir PATH/council-cursor
mycelium new model-report "Model B" --dir PATH/council-cursor
mycelium new reconciliation "SQLite council reconciliation" --dir PATH/council-cursor
# EDIT to Appendix C (adapter=cursor; SEED-DISSENT on RPT-001 and RCL)
mycelium check --dir PATH/council-cursor
# exit 0
```

Negatives: drop RCL → exit 1; drop RPT-002 → exit 1; strip `SEED-DISSENT` from RCL → exit 1; `cost_class = "cheap"` → exit 1.

### manual-floor row

Same as Cursor-council with `adapter = "manual"` on CMP + both RPTs. Exit 0. This **is** the manual-floor matrix row. It does not paste into a real UI.

### Also in the same `go test ./...`

| Case | Expect |
| --- | --- |
| Lone CMP, legal defaults, no RPT | exit 0 (WIP) |
| Spark scaffold, pack present, zero CMP | exit 0 |
| Pack deleted, no `reviews/` | exit 0 |
| `prompt_sha256` uppercase | exit 1 |
| RPT `rung` ≠ CMP `rung` | exit 1 |
| RCL `rung = "second-opinion"` | exit 1 |
| `gh` invoked | FAIL the test |

`MYCELIUM_NOW` is optional / not required. No seven-real-day dogfood. No live Cursor. No Actions job.

## Appendix F — Target file tree additions (and explicit DO NOT ADD workflow files / council CLI)

### Master (additions on top of the PHASE-01–03 tree; v1 files retained)

```text
internal/pack/                    # presence, collision, ReviewsAllowed
internal/ladder/                  # hash, IFF, cardinality, SEED-DISSENT
internal/clitest/ms401_hermetic_test.go
program/packs/council/            # entire §5 tree; do not add extra files
program/contracts/conformance.md  # items 18–22 + lift timing
program/contracts/identifiers.md  # CMP RPT RCL
program/contracts/naming.md       # rows + link regex
program/skills/mycelium-cli/SKILL.md
program/skeleton/AGENTS.md        # capability note
framework/phases/PHASE-04-implementation-brief.md
framework/phases/PHASE-04-acceptance.md
internal/embed/program/           # regenerate after program/ edits
```

### DO NOT ADD

```text
.github/workflows/phase-04-*.yml
.github/workflows/phase-04-hermetic.yml
.github/workflows/phase-04-ms401.yml
framework/decisions/DEC-015-*.md
program/packs/{thinking,spark,wake,portfolio,handoff,supersede,registry.toml}
a mycelium council / second-opinion / ladder / replicate command package
a forked copy of program/contracts/replication-reconciliation.md
internal/cursor/ or any model-invoking client
internal/panels/ (no Go reader this phase)
internal/slug changes
```

Do **not** add a PHASE-04 workflow. Do **not** extend `phase-01-hermetic.yml` as a phase gate. Quality should refuse a PR that adds an Actions job as the MS-401 gate. Quality should **not** refuse a missing PHASE-04 workflow (absence is correct). Do not delete Justfile/scripts/`research-program.toml`/PHASE-01 workflows. Do not add a `mycelium council` command.

### Emitted instance (spark / focused, local-only, PHASE-04 scaffold)

```text
README.md  mycelium.toml  log.md  index.md  CONTEXT.md  AGENTS.md  .gitignore
.agents/skills/{mycelium-cli,spark,wake,portfolio,thinking,council,second-opinion}/SKILL.md
program/ … including program/packs/council/
.git/          # init only; no commit
```

Absent: `framework/`, `cmd/`, `internal/`, `go.mod`, `Justfile`, `scripts/`, `research-program.toml`, v1 `research-*` skills, a `council` CLI verb.

After `mycelium new commissioning` (Slice 2+): `reviews/commissioning/CMP-001-*.md` with the §11 defaults. After disable (delete `program/packs/council`): pack types unknown; leftover `reviews/` FAILs extra-top-level unless deviation.

Unexported helpers may live next to their tests. No `pkg/`. No extra public command packages. Do not touch `internal/slug`.

End of PHASE-04 implementation brief. Engineering executes from this file only.
