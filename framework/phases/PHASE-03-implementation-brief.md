# PHASE-03 Implementation Brief — Sparring

- **Status:** Binding
- **Date:** 2026-08-15
- **Audience:** Engineering (pstack / poteto-mode)
- **Authority:** `framework/blueprint.md` (Accepted 2026-08-14). DEC-001 through DEC-014 stand. **DEC-007 is this phase** (already Accepted). Implement it. Do **not** record DEC-015 (see Appendix A).
- **Repo:** https://github.com/robertguss/mycelium
- **Pin:** Engineering starts from `main` @ `c9ad97874a94e314ceaf1eb2245e433175990f64` (PHASE-01 + PHASE-02 accepted on main). Do not implement from a later SHA unless Arvo re-pins in writing.
- **Product:** single-binary Go CLI `mycelium`. Master builds the CLI. `program/` is `go:embed`'d and emitted on scaffold. PHASE-03 adds the thinking-mode skill and agreement-conditional disagreement-record grammar on existing OQ files. **No new top-level CLI verb.**
- **Phase gate:** MS-301 via hermetic `go test ./...` only (two fixture sessions). GitHub Actions is **not** a gate (Robert waived CI). Do **not** add `.github/workflows/phase-03-*.yml`. Do **not** extend `phase-01-hermetic.yml` as a phase gate. Quality should refuse a PR that adds an Actions job as the MS-301 gate. Quality should **not** refuse a missing PHASE-03 workflow (absence is correct).
- **How to use this file:** Engineering executes from THIS FILE ONLY. No "see chat". No TBD. Open items are decided here and labeled **Architect default**.

Headings: §§1–19 then Appendices A–F (no DEC-015, disputed OQ example, aligned OQ example, `CONTEXT.md` glossary example, MS-301 fixture recipe, target file tree).

Cloud env name is exactly `robertguss/mycelium`. Go 1.26 at `/usr/local/go` or `/usr/local/bin/go`. `CGO_ENABLED=0`. stdlib + `github.com/pelletier/go-toml/v2` only. No cobra / viper / yaml / testify / go-github.

## 1. Scope / out of scope

Tonight is PHASE-03 only. Later phases get their own briefs after MS-301 is accepted. Do not implement, stub-ship, or "leave a hook command" for later phases. Do not implement PHASE-01 or PHASE-02 leftovers.

### In scope

- Thinking-mode skill `program/skills/thinking/SKILL.md` (front matter `name: thinking`): mandatory positions, agreement states, disagreement records with cruxes, glossary challenge, assumption audit. Grilling and domain-modeling conventions absorbed and credited. Emit on **new scaffolds only**.
- Agreement-conditional structure on existing OQ files (`questions/OQ-###-*.md`). Disagreement record lives **ON the OQ file**. No new namespace.
- New contract `program/contracts/sparring.md` plus glossary rules for `CONTEXT.md`.
- Check updates listed in §9 (conditional headings; glossary H2 ⇒ H3 `Definition`; optional DEC `## Dissent`).
- Template / schema deltas in §11 (`question.md` + `question.schema.toml`; optional Dissent hint on `decision.md`; `CONTEXT.md` H1 stays).
- MS-301 hermetic fixture tests in `go test ./...`. No Actions job.
- Commissioning artifacts: this brief, sparring + glossary contracts, conformance items 15–17, acceptance stub.
- `program/skills/mycelium-cli/SKILL.md` and `program/skeleton/AGENTS.md` updates that name the `thinking` skill.

### Out of scope (Quality refuses PRs that add them)

- Any new top-level CLI verb: `mycelium think`, `spar`, `session`, `council`, `handoff`, `supersede`. Surface is existing `mycelium new question`, `mycelium new assumption`, `mycelium new decision`, `mycelium new spike`, `mycelium check`, plus the new `thinking` skill.
- A `DSG-###` (or any new) namespace. Disagreement is not a new type.
- A `session.md` type. A "session" is the set of OQ files in a fixture instance.
- A required `index.md` H2. PHASE-02 H2s stay: State, Artifacts, Log tail, Wake.
- Requiring any OQ to exist. Spark instances without questions still pass.
- Requiring `## Crux` or `## Reasons` when `agreement` is `aligned` or `open`. The aligned MS-301 session passes with **no disagreement record required**.
- Content-scoring of positions, reasons, cruxes, or "substantive vs bare" (DEC-005). That judgment belongs to the human or an adversarial reviewer, **never an automated content score**.
- Assumption-audit file type. Audit is skill-only via existing `mycelium new assumption`.
- Auto-promotion of a crux to a spike or research stage. Skill may say a crux is eligible; no command this phase.
- PHASE-04 perspective ladder / second-opinion / council (DEC-008).
- PHASE-05 `mycelium supersede` and release/install work.
- PHASE-06 handoff packet. `handed-off` stays unreachable.
- Implementing MS-101(b). Do not commission a `GH_TOKEN` job. Do not reopen publish.
- partial: local-only is NOT required this phase (that's PHASE-02).
- Growing `latinFold`, adding NFKD, or adding `golang.org/x/text` (DEC-014).
- Converting master (`research-program.toml`, `just init`, deleting Justfile/scripts).
- Emitting `framework/`.
- CLI `git add` / `git commit` of instance work product.
- Retrofitting the `thinking` skill into existing instances. No retrofit command.
- A seven-real-day dogfood as the phase gate.
- Adding `.github/workflows/phase-03-*.yml` or extending `phase-01-hermetic.yml` as a phase gate. Quality should refuse a PR that adds an Actions job as the MS-301 gate.

### Master vs instance (unchanged)

Master remains an ADRP v1 instance for its own evolution. Do not convert master's `research-program.toml` to `mycelium.toml`. Do not run `just init` on master. `framework/` stays master-only and is NEVER emitted. Justfile/scripts stay on master. PHASE-03 changes the operational surface for *idea instances* only.

## 2. Authority and do-not-reopen

### Governing documents

| Document | Rule |
| --- | --- |
| `framework/blueprint.md` (Accepted 2026-08-14) | Do not rewrite vision. DEC-001–014 stand. PHASE-03 text at blueprint lines ~411–419 is the scope ceiling. |
| DEC-005 | Checks validate containers, never contents. No automated content score. |
| DEC-006 | spark → exploring ⇄ simmering → clarified → handed-off; any → archived. Do not change the machine. |
| **DEC-007** | **This phase.** Mandatory positions, agreement states, crux-bearing disagreement records, glossary challenge, assumption audit. Already Accepted. Implement it. |
| DEC-008 | Ladder / second-opinion / council = PHASE-04. Not this brief. |
| DEC-010 / DEC-011 | CLI never commits. No migration. Runtime reads instance files. |
| DEC-012 / DEC-013 / DEC-014 | Do not reopen (`mycelium.toml`; refuse out-of-range; `latinFold` only, no NFKD, no `x/text`, do not grow the map). |
| This brief | Binding 2026-08-15. PHASE-03 only. Architect defaults are binding. No DEC-015. |

### Process override (unchanged)

Blueprint "humans-own-git" is overridden for the *master* repo's engineering process: Arvo merges Quality-green PRs and accepts the phase. Engineering opens PRs. Engineering does NOT push to main. The CLI still never git-commits *instance* work product.

### Do not reopen

Do not reopen the product shape, the language, the dependency floor, the state vocabulary, the manifest filename, the refuse-vs-warn range rule, the no-commit rule, the instance-files-are-truth rule, slugify/DEC-014, publish, or MS-101(b). If a later discovery contradicts a locked item, stop and write a DEC; do not silently patch this brief in a code PR.

Do not reopen:

- DEC-012 (`mycelium.toml`)
- DEC-013 (refuse out-of-range)
- DEC-014 (slugify = existing `latinFold`; no NFKD; no `x/text`; do not grow the map)

Do **not** record DEC-015. DEC-007 already exists and is this phase.

### Leftovers stay leftovers

- MS-101(b) has **not** passed (`GH_TOKEN` missing; skip ≠ pass). Do **not** implement it. Do **not** commission a `GH_TOKEN` job.
- PHASE-02 seven-real-day dogfood is human evidence, not a gate, and is **not** this phase.
- `partial: local-only` is NOT required this phase (that's PHASE-02).
- Do **not** reimplement shipped PHASE-01 or PHASE-02 work.

### Phases you must not commission

| Phase | Why it is not this brief |
| --- | --- |
| PHASE-04 | DEC-008 ladder / second-opinion / council |
| PHASE-05 | `mycelium supersede` |
| PHASE-06 | Handoff packet; `handed-off` becomes legal only then |

## 3. What PHASE-01+02 left on main (floor; do not reimplement)

Pin: `c9ad97874a94e314ceaf1eb2245e433175990f64`. Treat this SHA as the floor. Reuse packages. Do not rewrite working PHASE-01 or PHASE-02 commands.

### Already shipped (do not rebuild)

Reuse: `cmd/mycelium`, `internal/{cli,version,embed,clock,execrun,metadata,idpath,manifest,schema,slug,logfmt,teach,lock,journal,op,scaffold,generate,check,tiercmd,publish,clitest,revisit,lifecycle,indexmd,wakebrief,statecmd,statuscmd}`.

| Touch | Fate |
| --- | --- |
| `internal/cli` | **No new verb.** Commands already: `version`, `help`, `new`, `check`, `tier`, `publish`, `index`, `state`, `wake`, `status`. Do not add `think` / `spar` / `session` / `council` / `handoff`. |
| `internal/check` | **Extend** with agreement-conditional OQ rules (Slice 2), glossary `CONTEXT.md` rules (Slice 3), optional DEC Dissent (Slice 4). Do not rewrite the package. |
| `internal/generate` / `internal/schema` | **Extend** as needed for the question schema `required_sections` drop of always-required `Crux`. Do not invent a schema DSL for IFF rules. |
| `internal/scaffold` | Emit the new `thinking` skill on **new scaffolds only** (Slice 5). No retrofit. |
| `internal/slug` | Do not touch (DEC-014). |
| `internal/clock` + `MYCELIUM_NOW` | Reuse if a fixture wants stable dates. **Not required** for MS-301 unless you choose to pin dates on generated files. |
| `internal/execrun` | Unused by MS-301. Hermetic tests still must not call `gh`. |
| `internal/lock` / `journal` / `op` | Do **not** add a `think` / `spar` op. `new` already covers `mycelium new question`. |
| `internal/lifecycle` / `statecmd` / `wakebrief` / `indexmd` / `statuscmd` | Do not rewrite. Lifecycle storage rules stay PHASE-02. |
| `program/templates/question.md` + `question.schema.toml` | **Delta** (§11). Type `OQ`, home `questions/`, enum `open\|aligned\|agree-to-disagree` already exist. |
| `program/skeleton/CONTEXT.md` | Already `# Glossary` plus a blank line. Empty glossary stays legal. |
| `program/contracts/{conformance,lifecycle,identifiers}.md` | Extend conformance with items 15–17. Do not rewrite lifecycle. Do not add a `DSG` namespace to identifiers. |
| `program/skills/{spark,wake,portfolio,mycelium-cli}` | Do not rewrite spark/wake/portfolio. Update `mycelium-cli` to name `thinking`. |
| `phase-01-hermetic.yml` / `phase-01-github.yml` | Leave alone. Do **not** add a PHASE-03 workflow. Do **not** extend `phase-01-hermetic.yml` as a phase gate. Actions is not a gate. |
| Justfile / scripts / `research-program.toml` | Keep. Do not delete. Do not `just init`. |
| `framework/` | Master-only. NEVER emitted. |
| `internal/version` | Stay `0.1.0-dev` unless already different on the pin — do not bump as a phase ritual. `methodology_version` stays `2.0.0`. Re-run embed generate after `program/` edits. |

### PHASE-01+02 behaviors that stay true

- Birth state is `spark`. Stored `handed-off` still FAILS check. `clarified` is legal. `simmering` requires revisit grammar.
- `index.md` required H2s stay: State, Artifacts, Log tail, Wake. Questions show up under Artifacts if the renderer already lists type homes. **Do not add a required H2.**
- Clock is injectable (`internal/clock` + `MYCELIUM_NOW` RFC3339).
- `git` / `gh` go through `internal/execrun`.
- Log lines are tab-separated (`internal/logfmt`). Log-ops stay `scaffold|new|tier|publish|check|state|wake`.
- Exit 0 success, exit 1 every failure. Teaching errors: four lines on stderr (`mycelium` / `convention` / `contract` / `fix`). Cap 20.
- `--dir PATH` is the instance root. Absent: walk upward for `mycelium.toml`, stop at `.git` or filesystem root.
- CLI never git-commits. `git init` only on scaffold, already shipped.
- Runtime check/generate read instance `program/`, never embed (DEC-011). Dependency floor unchanged. Module `github.com/robertguss/mycelium`. `CGO_ENABLED=0`. Go 1.26.
- Question type is already registered: namespace `OQ`, home `questions/`, filename `OQ-{NNN}-{slug}.md`, required front matter `id`, `title`, `agreement`, `date`.
- On the pin, `question.schema.toml` lists `Crux` in `required_sections`. That always-required Crux is the PHASE-03 schema delta — drop it from the always-required list (Slice 1). Conditional Crux binds in Slice 2.
- `mycelium new question` already exists. Do not reimplement it.

### What must not be broken

`just check` on master; hermetic `go test ./...`; no `framework/` emit; no master conversion. PHASE-01/02 fixtures that already have `## Crux` on an `open` OQ stay green (extra Crux/Reasons on `open`/`aligned` do not fail). PHASE-01/02 fixtures that have DECs without `## Dissent` stay green. Spark instances with zero questions stay green.

If a PHASE-03 PR is bad: revert that PR. Floor is the pin SHA.

## 4. Sparring stance + agreement states (the DEC-007 table)

DEC-007 is Accepted. This section encodes it as check vs skill. Check stays structure-only (DEC-005). The skill carries the stance so it survives model and runtime changes.

### DEC-007 clauses

| # | Stance (DEC-007) | Check this phase | Skill this phase |
| --- | --- | --- | --- |
| 1 | Every substantive question carries the agent's position or recommendation. Bare questions are a smell. | `## Positions` is always required as a **container**. Body is not graded. A one-word or `<!-- fill -->` body does **not** fail. | Agent writes a real position. Substantive-vs-bare judgment belongs to the human or an adversarial reviewer, never an automated content score (DEC-005). |
| 2 | Every substantive question carries an agreement state: `open`, `aligned`, or `agree-to-disagree`. The third is terminal and honorable. | Front matter `agreement` must be exactly one of those three. Invalid value → FAIL. | Set the field. `aligned` and `agree-to-disagree` are **terminal-by-convention** (do not edit back to `open`; open a new OQ). Check does not keep agreement history. |
| 3 | Disagreement records capture both positions, both sets of reasons, and the cruxes: what evidence would change each mind. | IFF `agreement == "agree-to-disagree"`: required H2/H3 grammar in §5. IFF `aligned` or `open`: `## Crux` and `## Reasons` are **not** required. | Fill Human + Agent under Positions, Reasons, and Crux when disputed. Aligned session: **no disagreement record required**. |
| 4 | Dissent is retained forever, including the agent's after being overruled. | Optional DEC `## Dissent`. If the heading exists, the section must contain at least one resolvable `OQ-###` or `ASM-###`. | Keep dissent on the DEC when overruled. Do not require Dissent on existing DECs. |
| 5 | Assumption audit is a standing move. | **None.** No AUDIT file type. Check does not require a periodic audit file. | Agent uses existing `mycelium new assumption`. |
| 6 | Glossary discipline: terms challenged on drift, sharpened on vagueness, recorded on resolution. | `CONTEXT.md`: H1 `# Glossary` stays. Empty (H1 only) is legal. Any `## <Term>` must have H3 `Definition`. Do not require N terms. Do not grade definitions. | Skill does the challenge. Credit domain-modeling. |

### Agreement enum (do not bikeshed)

```text
open | aligned | agree-to-disagree
```

| Value | Meaning | Terminal-by-convention? | Disagreement record (Crux + Reasons + H3s)? |
| --- | --- | --- | --- |
| `open` | Still in play. Positions container present. | No. Edit in place. | **Not** required. Extra Crux/Reasons do not fail. |
| `aligned` | Parties agree. Honorable close. | **Yes.** Do not edit back to `open`; open a new OQ. | **Not** required. This is how the aligned MS-301 session passes with **no disagreement record required**. |
| `agree-to-disagree` | Parties disagree and stop. Terminal and honorable (DEC-007). | **Yes.** Do not edit back to `open`; open a new OQ. | **Required.** Record retains **both positions**, both sets of reasons, and **cruxes**. |

Invalid `agreement` value → FAIL. Teaching error names `program/templates/question.schema.toml` (enum) and the illegal token.

### Terminal-by-convention (not time-travel)

**Architect default:** `aligned` and `agree-to-disagree` are terminal **by skill + contract**, not by a check that keeps agreement history. Check does not store previous `agreement` values and does not fail a file that was flipped back to `open` (undetectable without history). The `thinking` skill and `program/contracts/sparring.md` say: do not edit those states back to `open`; open a new OQ instead.

### Always-on, no session type

Sparring is **always-on** (skill) in any non-archived state. Check does **not** require any OQ to exist. Spark instances without questions still pass.

There is **no** `session.md` type. A "session" in MS-301 is the set of OQ files in a fixture instance.

### Crux → next experiment

Skill may say a crux is eligible to become `mycelium new spike` or a research stage. **No auto-promotion command this phase.**

## 5. Disagreement record grammar (conditional headings)

The disagreement record lives **ON the OQ file**. Do not add `DSG-###`. Do not add a home directory.

### Always (every OQ)

| Kind | Binding |
| --- | --- |
| Front matter | `id`, `title`, `agreement`, `date` (already on the pin). |
| `agreement` enum | `open` / `aligned` / `agree-to-disagree`. Invalid → FAIL. |
| H2 `## Question` | Required. Body not graded. |
| H2 `## Context` | Required. Body not graded. |
| H2 `## Positions` | Required as a **container**. Heading must exist. Check does **not** grade whether the body is a real position (DEC-005). A one-word or `<!-- fill -->` body does **not** fail. |
| H2 `## Disposition` | Required. Body not graded. |

Heading match is exact and case-sensitive (PHASE-01 H2 rule). `## positions` fails. Extra H2s are allowed.

### IFF `agreement == "agree-to-disagree"`

All of the following are required. Missing any → FAIL. Teaching error names the **missing heading** and `program/contracts/sparring.md`.

| Heading | Where | Nested H3s |
| --- | --- | --- |
| `## Positions` | already always-required | **Must** contain `### Human` and `### Agent` |
| `## Reasons` | required only in this state | **Must** contain `### Human` and `### Agent` |
| `## Crux` | required only in this state | **Must** contain `### Human` and `### Agent` |

H3 names are exact: `### Human` and `### Agent`. Not `### The Human`, `### human`, `### HUMAN`, `### Resident agent`.

**Architect default — section body:** H3s must appear in the parent H2's section body (bytes after that exact H2 until the next H2 or EOF). Same container rule PHASE-02 used for wake trigger extract. Extra H3s in the section do not fail. H3 body may be `<!-- fill -->` or one word — tests must **not** grade the words.

### IFF `agreement == "aligned"` or `open`

| Heading | Check |
| --- | --- |
| `## Positions` | Required as a container. H3 `Human` / `Agent` are **not** required. |
| `## Crux` | **Not** required. |
| `## Reasons` | **Not** required. |
| Extra `## Crux` / `## Reasons` | Do **not** fail. |
| Extra H3s | Do **not** fail. |

This is how the aligned MS-301 session passes with **no disagreement record required**.

### Schema cannot express IFF

**Architect default:** `question.schema.toml` `required_sections` becomes `["Question", "Context", "Positions", "Disposition"]`. Drop `Crux` from the always-required list. Do **not** add `Reasons` to `required_sections`. Do **not** invent a schema DSL for conditionals. IFF rules live in `program/contracts/sparring.md` and `internal/sparring`. Check binds them in Slice 2.

Slice 1 side-effect (accepted): once the schema drops always-required `Crux`, the existing schema-driven check stops requiring `Crux` on every OQ. That is intended. The IFF "require Crux when `agree-to-disagree`" does **not** bind until Slice 2.

### Teaching errors (binding shape)

Four lines, stderr, exit 1.

```text
mycelium: OQ-001 missing ## Crux (required when agreement=agree-to-disagree)
convention: sparring
contract: program/contracts/sparring.md
fix: add ## Crux with ### Human and ### Agent

mycelium: OQ-001 missing ### Human under ## Positions
convention: sparring
contract: program/contracts/sparring.md
fix: add ### Human and ### Agent under ## Positions

mycelium: OQ-001 missing ## Reasons (required when agreement=agree-to-disagree)
convention: sparring
contract: program/contracts/sparring.md
fix: add ## Reasons with ### Human and ### Agent

mycelium: agreement "maybe" is not open|aligned|agree-to-disagree
convention: question-front-matter
contract: program/templates/question.schema.toml
fix: set agreement to open, aligned, or agree-to-disagree
```

### Pure parser (Slice 1)

New package `internal/sparring`. No filesystem. No CLI. Table-driven tests only. Check does not call it until Slice 2.

```text
ParseAgreement(s string) (state, error)          # open | aligned | agree-to-disagree
RequiredH2(agreement) []string                   # always Question, Context, Positions, Disposition; plus Reasons, Crux iff agree-to-disagree
RequiredH3(agreement, h2) []string               # Human, Agent iff agree-to-disagree and h2 in {Positions, Reasons, Crux}
MissingHeadings(agreement, body) []string        # exact missing heading names for teaching errors
SectionBody(body, h2) string                     # bytes after H2 until next H2 or EOF
```

Do not rewrite `internal/check` in Slice 1. Do not bind the IFF rules in check in Slice 1.

## 6. Glossary + assumption audit conventions

### Glossary (`CONTEXT.md`)

Pin skeleton is already:

```text
# Glossary

```

| Rule | Binding |
| --- | --- |
| H1 | `# Glossary` stays. **Architect default:** if `CONTEXT.md` exists, it must contain a line that is exactly `# Glossary`. |
| Empty glossary | H1 only is **legal**. Check does not require N terms. |
| Any `## <Term>` | That term's section **must** contain H3 `Definition`. Missing → FAIL. Teaching error names the term heading and `program/contracts/glossary.md`. |
| Definition body | Not graded. `<!-- fill -->` or one word does **not** fail. |
| File existence | Do **not** add a new required-file bind. New scaffolds already emit `CONTEXT.md`. If the file is missing, do not fail solely for that. If it exists, apply the rules. |
| Content | Check does not grade definitions, challenge quality, or drift. Skill does. |

Credit **domain-modeling** (glossary + ADR) in the `thinking` skill. Do not copy upstream domain-modeling files into `program/`.

Teaching error:

```text
mycelium: CONTEXT.md term "SQLite" missing ### Definition
convention: glossary
contract: program/contracts/glossary.md
fix: add ### Definition under ## SQLite
```

### Assumption audit (skill-only)

| Rule | Binding |
| --- | --- |
| File type | **None.** Do not add an `AUDIT` type or namespace. |
| Command | Agent uses existing `mycelium new assumption`. |
| Check | Does **not** require a periodic audit file. Does **not** require N assumptions. |
| Skill | Standing move: both parties periodically dump what they are taking for granted; the agent infers and reads back the human's unstated presuppositions. |

Do not reimplement the assumption template or `mycelium new assumption`.

## 7. Thinking-mode skill (grilling + domain-modeling credited)

### Identity

| Field | Binding |
| --- | --- |
| Path | `program/skills/thinking/SKILL.md` |
| Front matter | `name: thinking` |
| Emit | `.agents/skills/thinking/SKILL.md` on **new scaffolds only** |
| Retrofit | **None.** No retrofit command. `tier` / `index` / `state` / `wake` / `check` do not emit it. |
| Credits (must appear in the skill body) | mattpocock grilling (recommendation per question; decisions stay the user's); domain-modeling (glossary + ADR); pstack poteto candor ("no is an acceptable answer"); DEC-007 |
| Upstream files | Do **not** copy grilling / domain-modeling / poteto files into `program/`. Credit, do not vendor. |

### When it applies

Always-on in any **non-archived** state. On `archived`, the skill says: do not open new OQs; existing records stay. Check does not enforce skill behavior and does not require OQs.

### Binding procedure (encode this; do not invent a second flow)

1. Read `index.md` and `CONTEXT.md` (not the whole tree).
2. On every substantive question: take a position; capture it with `mycelium new question` (or edit the OQ). Bare questions are a smell — the human or an adversarial reviewer judges substance, never an automated content score (DEC-005).
3. Set `agreement` to `open`, `aligned`, or `agree-to-disagree`.
4. If `agree-to-disagree`: fill `## Positions`, `## Reasons`, and `## Crux`, each with `### Human` and `### Agent`. The record retains **both positions**, both sets of reasons, and **cruxes**.
5. If `aligned`: `## Positions` H2 present; **no disagreement record required** (no Crux/Reasons required).
6. If `open`: keep working; Positions container present; Crux/Reasons not required.
7. Do not edit `aligned` or `agree-to-disagree` back to `open`. Open a new OQ instead.
8. Glossary challenge: on drift or vagueness, sharpen the term; record under `CONTEXT.md` as `## <Term>` + `### Definition`.
9. Assumption audit: periodically dump presuppositions via `mycelium new assumption`. No AUDIT file.
10. A crux is eligible to become `mycelium new spike` or a research stage. Do not auto-promote. No new command.
11. Recommendation per question; **decisions stay the user's** (grilling). "no is an acceptable answer" (poteto candor).
12. Run `mycelium check` before handing back.
13. Do not `git commit` unless the human asks. The CLI never commits.

### What the skill must not say

- Do not tell the agent to run `mycelium think` / `spar` / `session` / `council` / `handoff` (those verbs do not exist).
- Do not tell the agent to create a `DSG-###` file or a `session.md`.
- Do not tell the agent to flip `handed-off`.
- Do not tell the agent that CI / an Actions job is the done bar.

## 8. Commands (existing only; what does not exist)

Exact CLI. No new top-level verb. Flags and exit codes stay PHASE-01/02.

Global: exit 0 success, exit 1 failure. Teaching errors on stderr (four lines). Success text on stdout.

Env (unchanged; MS-301 does not require them):

| Env | Effect this phase |
| --- | --- |
| `MYCELIUM_OFFLINE=1` | Unchanged. Unused by MS-301. |
| `MYCELIUM_NOW` | RFC3339 clock override. **Optional / not required** for MS-301 unless a test pins dates on generated files. |
| `MYCELIUM_IDEAS_ROOT` | Unchanged. Portfolio `partial: local-only` is NOT required this phase (that's PHASE-02). |

### Commands that exist after PHASE-03 (same argv as the pin)

| Command | New? | PHASE-03 delta |
| --- | --- | --- |
| `mycelium version` | no | none |
| `mycelium new idea <name> [--dir] [--offline] [--publish] [--tier]` | no | Slice 5: emit `thinking` skill on new scaffolds |
| `mycelium new question <title> [--dir]` | no | Slice 1: template/schema omit always-required Crux |
| `mycelium new assumption <title> [--dir]` | no | skill uses it for the audit; no command change |
| `mycelium new decision <title> [--dir]` | no | Slice 4: optional Dissent **hint** in the template, not a live required heading |
| `mycelium new spike <title> [--dir]` | no | skill may point a crux here; no auto-promotion |
| `mycelium check [--dir] [--abort-journal]` | no | Slices 2–4 bind new structure rules |
| `mycelium tier <tier> [--dir]` | no | does **not** emit `thinking` |
| `mycelium publish [--dir]` | no | do not reopen |
| `mycelium state <target> [--dir] [--revisit]` | no | do not rewrite |
| `mycelium wake [--dir]` | no | do not rewrite |
| `mycelium status [--dir] [--all] [--root] [--archived] [--offline]` | no | do not rewrite |
| `mycelium index [--dir]` | no | do not add an H2; do not emit `thinking` |

### Commands that do not exist this phase

`think`, `spar`, `session`, `council`, `handoff`, `supersede`, `explore` / `simmer` as separate verbs, `destroy`, `range`, any `new disagreement` / `new session` type.

Quality should refuse PRs that add them.

### 8.1 `mycelium new question` (already shipped)

```text
mycelium new question <title> [--dir PATH]
```

Allocates the next `OQ-###`, writes `questions/OQ-###-<slug>.md` from the (updated) template, appends a `new` log line, regenerates `index.md`. Journal `op=new`. Never git-commits.

After Slice 1 the emitted file has `agreement = "open"` and H2s Question, Context, Positions, Disposition. It does **not** emit `## Crux` or `## Reasons`.

The human or agent then edits front matter and headings. That edit is not a CLI verb.

### 8.2 `mycelium check` (delta is §9)

```text
mycelium check [--dir PATH] [--abort-journal]
```

No new flags. Same teaching-error shape. Same `--abort-journal`. Success stdout unchanged.

### 8.3 Instance root resolution

Unchanged from PHASE-01/02.

## 9. Check updates (what changes from PHASE-02 conformance)

Structure only (DEC-005). Runtime still reads instance files, never embed.

### PHASE-02 items 1–14 stay

Do not rewrite lifecycle storage rules. Do not lift `handed-off` FAIL. Do not change `index.md` required H2s. Do not change log-ops. Do not require a wake brief on instances that never simmered.

### New must-implement items (15–17)

| # | Rule | Slice |
| --- | --- | --- |
| 15 | Agreement-conditional OQ headings per §5 / `program/contracts/sparring.md`. Invalid `agreement` → FAIL (enum already on the pin; keep it). | Bind in Slice 2 (schema drop of always-required Crux lands in Slice 1). |
| 16 | If `CONTEXT.md` exists: H1 `# Glossary`; any `## <Term>` ⇒ H3 `Definition`. Empty glossary legal. Do not require N terms. Do not grade definitions. | Slice 3 |
| 17 | If a DEC contains `## Dissent`, the section must contain at least one resolvable `OQ-###` or `ASM-###`. Heading absent → pass. | Slice 4 |

### What check must not do

| Temptation | Verdict |
| --- | --- |
| Require any OQ on spark / exploring / any state | **No.** Spark with zero questions still passes. |
| Require `## Crux` or `## Reasons` on `aligned` or `open` | **No.** |
| Require H3 Human/Agent on `aligned` or `open` | **No.** |
| Grade Positions / Reasons / Crux / Definition prose | **No.** DEC-005. |
| Score substantive vs bare | **No.** Human or adversarial reviewer. |
| Require a periodic assumption-audit file | **No.** |
| Require `## Dissent` on existing or new DECs | **No.** |
| Add a required `index.md` H2 | **No.** |
| Fail extra Crux/Reasons on `open`/`aligned` | **No.** |
| Keep agreement history / fail a flip back to `open` | **No.** |
| Add `think` / `spar` to the log-op regex | **No.** |
| Call network / `gh` / read `GH_TOKEN` | **No.** |

### Lift timing

| Slice | Check behavior |
| --- | --- |
| 1 | Schema `required_sections` drops `Crux`. Existing schema-driven check therefore stops requiring Crux on every OQ. **No** IFF bind yet. An `agree-to-disagree` OQ missing `## Crux` still **passes** check in this slice. Parser + table tests only. |
| 2 | Check calls `internal/sparring` for each `questions/OQ-*.md`. IFF rules bind. Missing Crux/Reasons/H3s on `agree-to-disagree` → FAIL. `aligned`/`open` without Crux/Reasons → pass. Invalid agreement → FAIL. Positions `<!-- fill -->` on `open`/`aligned` → pass. |
| 3 | `CONTEXT.md` glossary rules bind. |
| 4 | Optional DEC Dissent rule binds. PHASE-01/02 fixtures without Dissent stay green. |

### `program/contracts/conformance.md`

Add items 15–17 so Quality can thermos a numbered list. State the lift timing in the contract the same way PHASE-02 did. Do not renumber items 1–14.

### Link scan

Unchanged. Type homes already include `questions/` and `assumptions/`. Dissent IDs resolve through the existing scan. **Architect default:** do not newly require `CONTEXT.md` in the link scan.

## 10. Skills + mycelium-cli + AGENTS.md

Emit path: `.agents/skills/<name>/SKILL.md` like `mycelium-cli`.

Source of truth: `program/skills/<name>/SKILL.md`. Scaffold copies into the instance. Re-run `go generate` in `internal/embed`.

### When skills are emitted

| Event | Emit `thinking`? |
| --- | --- |
| New `mycelium new idea` scaffold (PHASE-03+) | **Yes** |
| `mycelium tier` | **No** |
| `mycelium index` | **No** |
| `mycelium state` / `mycelium wake` / `mycelium check` | **No** (do not retrofit) |
| One-shot `program/skills/` copy command | **Not a PHASE-03 command** |

**Architect default:** do not rewrite `spark` / `wake` / `portfolio`. Existing PHASE-01/02 instances: re-scaffold or copy `thinking` **manually**. Document that sentence in `program/skeleton/AGENTS.md` and in `mycelium-cli` SKILL.md.

### `program/skills/mycelium-cli/SKILL.md` update

- Add a row / paragraph: sparring surface is `mycelium new question`, `mycelium new assumption`, `mycelium new decision`, `mycelium new spike`, `mycelium check`, plus the `thinking` skill. **No** `think` / `spar` verb.
- Extend the emit list: new scaffolds emit `.agents/skills/{mycelium-cli,spark,wake,portfolio,thinking}/SKILL.md`.
- Keep the manual floor, teaching-error shape, and "do not git commit unless the human asks".

### `program/skeleton/AGENTS.md` update

- Name the `thinking` skill in the Skills section emit list.
- Same no-retrofit sentence.
- Do not list `think` / `spar` / `session` / `council` / `handoff` as commands.

### Slice

Slice 5 lands `thinking` + mycelium-cli update + AGENTS.md + embed generate.

## 11. Templates / schema deltas (question, decision, CONTEXT)

### `program/templates/question.schema.toml`

On the pin:

```text
required_sections = ["Question", "Context", "Positions", "Crux", "Disposition"]
[enums.agreement] values = ["open", "aligned", "agree-to-disagree"]
```

After Slice 1:

```text
namespace = "OQ"
home = "questions"
filename_pattern = "OQ-{NNN}-{slug}.md"
stage_scoped = false
digits = 3
required_front_matter = ["id", "title", "agreement", "date"]
required_sections = [
  "Question",
  "Context",
  "Positions",
  "Disposition",
]

[enums.agreement]
values = ["open", "aligned", "agree-to-disagree"]
```

Do not add `Reasons`. Do not add a schema-level IFF table. Enum stays.

Existing tests that assert `Crux` is always in `required_sections` update in Slice 1.

### `program/templates/question.md`

On the pin the template emits `## Crux`. After Slice 1 it does **not**.

**Architect default — emitted shape** (tokens filled by the existing generator):

```text
+++
id = "{{ID}}"
title = "{{TITLE}}"
agreement = "open"
date = "{{DATE}}"
+++

# {{ID}} — {{TITLE}}

<!-- slug: {{SLUG}} -->

## Question

<!-- fill -->

## Context

<!-- fill -->

## Positions

<!-- fill -->

## Disposition

<!-- fill -->
```

No `## Crux`. No `## Reasons`. Agent/human adds those headings when flipping to `agree-to-disagree`.

Existing OQ files that still have `## Crux` on `open`/`aligned` stay green.

### `program/templates/decision.md` + `decision.schema.toml`

`required_sections` does **not** gain `Dissent`.

**Architect default:** do **not** emit a live `## Dissent` heading (a live empty heading would fail Slice 4 on every new DEC). Append an HTML comment:

```text
<!-- Optional section (not required): ## Dissent — if you add this heading, cite at least one resolvable OQ-### or ASM-### -->
```

Existing DECs without the heading stay green (PHASE-01/02 fixtures must stay green).

If a human/agent adds `## Dissent`, Slice 4 requires at least one resolvable `OQ-###` or `ASM-###` in that section body. Other prose is allowed and not graded. Citing only a `DEC-###` does **not** satisfy the dissent rule. Unresolved IDs still fail the existing link scan.

Teaching error:

```text
mycelium: DEC-001 ## Dissent has no resolvable OQ-### or ASM-###
convention: dissent
contract: program/contracts/sparring.md
fix: cite an existing OQ-### or ASM-### in ## Dissent, or remove the heading
```

### `program/skeleton/CONTEXT.md`

No delta required. H1 `# Glossary` stays. Empty glossary stays legal.

### `program/templates/assumption.md`

No delta. Audit uses the existing type.

## 12. Vertical slices 0–6 with build order, each checkable, go test / files only

PR-per-slice, sequential, rebase on main. Arvo merges Quality-green PRs. Engineering opens PRs. Engineering does NOT push to main. Prefer one live PR at a time. Do not stack unpublished slices on one branch unless Quality is backed up.

Each PR title: `PHASE-03 Slice N: <done-bar noun>`. Each PR body links this brief and the slice done bar. No drive-by refactors. No v1 deletions. No PHASE-04+ commands. No Actions job as a done bar.

### Slice 0 — Commissioning (docs only)

This brief + new contracts + acceptance stub. No product code. No Go.

Land: `framework/phases/PHASE-03-implementation-brief.md`, `framework/phases/PHASE-03-acceptance.md` (rows = §15), new `program/contracts/sparring.md`, new `program/contracts/glossary.md`, updates to `conformance.md` (items 15–17 + lift timing; may finish in Slices 2–4).

Done: files exist on a PR. Quality reads them against this brief. No product code. Do **not** add `.github/workflows/phase-03-*.yml`.

### Slice 1 — Template / schema / contract + pure parser (no IFF check-rule yet)

- `question.schema.toml` drops always-required `Crux`.
- `question.md` omits `## Crux` / `## Reasons`.
- `internal/sparring` parser + table tests (§5).
- Existing schema tests that treated Crux as always-required are updated.

Done: `go test` for `ParseAgreement` / `RequiredH2` / `RequiredH3` / `MissingHeadings` / `SectionBody`. New `mycelium new question` (if exercised) emits no Crux heading. Check does **not** yet fail an `agree-to-disagree` OQ missing Crux. No new CLI verb.

### Slice 2 — Check binds the conditional rules

- `internal/check` calls `internal/sparring` for each OQ.
- IFF rules bind. Invalid agreement FAIL. `aligned`/`open` without Crux/Reasons PASS. Positions `<!-- fill -->` on `open`/`aligned` PASS. Spark with zero questions PASS.
- Teaching errors name the missing heading and `program/contracts/sparring.md`.

Done (hermetic `go test`): disputed-shape missing `## Crux` → check exits 1; same shape complete → check exits 0; aligned without Crux/Reasons → check exits 0; invalid agreement → check exits 1; fill Positions on `open` → check exits 0. Files only + `go test`. No Actions job.

### Slice 3 — Glossary `CONTEXT.md` check

- If `CONTEXT.md` exists: require `# Glossary`; any H2 ⇒ H3 `Definition`.
- Empty (H1 only) PASS. Missing Definition FAIL. Do not require N terms. Do not grade definitions.

Done (hermetic `go test`): H1-only `CONTEXT.md` → check 0; one `## Term` without `### Definition` → check 1; one `## Term` with `### Definition` (body `<!-- fill -->`) → check 0.

### Slice 4 — Optional DEC Dissent

- Template comment hint (§11). Schema `required_sections` unchanged.
- If `## Dissent` exists: at least one resolvable `OQ-###` or `ASM-###`.
- PHASE-01/02 fixtures without Dissent stay green.

Done (hermetic `go test`): DEC without Dissent → check 0; DEC with Dissent citing a real OQ → check 0; DEC with Dissent and no OQ/ASM token → check 1.

### Slice 5 — `thinking` skill + mycelium-cli + AGENTS.md + embed generate

- `program/skills/thinking/SKILL.md` with the credits and procedure in §7.
- Update `mycelium-cli` and `AGENTS.md`.
- `go generate` in `internal/embed`.
- New scaffolds emit `.agents/skills/thinking/SKILL.md`.
- `tier` / `index` / `state` / `wake` do not retrofit.

Done (hermetic `go test`): new scaffold `--offline` has the thinking skill file; a fixture that lacked it still lacks it after `tier` / `index`. Credits present in the skill body (string asserts are enough; do not grade prose quality).

### Slice 6 — MS-301 two fixtures in `go test ./...`

Hermetic fixture recipe = Appendix E. Two temp instances. `go test ./...` only. No network. No `gh`. No `GH_TOKEN`. Do **not** add `.github/workflows/phase-03-*.yml`. Do **not** extend `phase-01-hermetic.yml` as a phase gate. Quality should refuse a PR that adds an Actions job as the MS-301 gate.

Done: `go test ./...` runs the MS-301 fixture green. That **is** the gate.

## 13. Done / verified mapped onto MS-301

MS-301 is the hermetic phase gate (two fixture sessions). **No seven-real-day dogfood.** **No Actions job.** Robert waived CI. Done bar is hermetic `go test ./...` only.

Blueprint MS-301 (do not expand): two fixture sessions pass — one deliberately disputed, whose record retains **both positions**, both sets of reasons, and **cruxes**; one honestly aligned, which passes with **no disagreement record required**. Substantive-vs-bare judgment belongs to the human or an adversarial reviewer, never an automated content score (DEC-005).

### MS-301 expected (authoritative; recipe in Appendix E)

| Session | Setup | `mycelium check` |
| --- | --- | --- |
| Disputed | focused scaffold; `mycelium new question` titled like "Use SQLite"; `agreement = "agree-to-disagree"`; Positions / Reasons / Crux each have H3 Human + H3 Agent with any non-empty text (tests must not grade the words) | exit 0 |
| Disputed negative (same test file) | same shape missing `## Crux` | exit 1 |
| Aligned | scaffold; `mycelium new question`; `agreement = "aligned"`; Positions H2 present; **no** Crux/Reasons headings | exit 0 |
| Invalid agreement | any OQ with `agreement` not in the enum | exit 1 |
| Spark, zero questions | `mycelium new idea --offline` then check; no OQ files | exit 0 |
| Bare Positions container | `agreement` is `open` or `aligned`; Positions body is `<!-- fill -->` | exit 0 |

`gh` never invoked. No `GH_TOKEN`. No Actions job. `MYCELIUM_NOW` is optional / not required.

### Slice → MS-301 map

| Slice | MS-301 clause |
| --- | --- |
| 0 | commissioning; not a runtime clause |
| 1 | parser + schema drop; IFF not bound yet |
| 2 | IFF rules so disputed / aligned / negative / invalid / fill-Positions behave |
| 3 | glossary; not required for the two sessions, must not break them |
| 4 | Dissent; not required for the two sessions; must not fail PHASE-01/02 DECs |
| 5 | thinking skill on new scaffolds; not a check clause |
| 6 | the two-fixture test in `go test ./...` **is** the gate |

PHASE-03 is accepted when MS-301 is green in `go test ./...` on main. Arvo accepts the phase. Engineering does not self-accept.

## 14. Automated test plan

Engineering MUST write these tests. Quality thermos against this list. Do NOT require Playwright, Docker, live GitHub, `GH_TOKEN`, or an Actions job in default `go test ./...`.

`partial: local-only` is NOT required this phase (that's PHASE-02).

### Unit (no network, no gh)

| Area | Cases |
| --- | --- |
| agreement parse | `open`, `aligned`, `agree-to-disagree`, empty, `maybe`, `Agree-To-Disagree`, `aligned ` |
| RequiredH2 | open/aligned → Question, Context, Positions, Disposition only; agree-to-disagree adds Reasons + Crux |
| RequiredH3 | empty unless agree-to-disagree; then Human+Agent on Positions, Reasons, Crux |
| MissingHeadings | each missing H2/H3 named; extra Crux on aligned → no missing; fill Positions on open → no missing |
| SectionBody | H3 inside parent; H3 after next H2 does not count; exact case |
| glossary | H1-only → no missing; H2 without Definition → missing; H2 with Definition → ok; body not graded |
| dissent IDs | extract `OQ-###` / `ASM-###` from a Dissent section; DEC-only token is not enough |

### Hermetic CLI (built binary, temp dirs)

| Case | Expect |
| --- | --- |
| `mycelium new question "Use SQLite"` after Slice 1 | file exists; no `## Crux` / `## Reasons` in the template output; `agreement = "open"` |
| disputed complete | Appendix E; `mycelium check` exits 0 |
| disputed missing `## Crux` | `mycelium check` exits 1; stderr names `## Crux` and `program/contracts/sparring.md` |
| disputed missing `### Human` under Positions | check exits 1; names the H3 |
| disputed missing `## Reasons` | check exits 1 |
| aligned, no Crux/Reasons | check exits 0 |
| aligned with extra Crux | check exits 0 (extra does not fail) |
| invalid agreement | check exits 1 |
| spark, zero questions | check exits 0 |
| Positions body `<!-- fill -->` on `open` or `aligned` | check exits 0 |
| CONTEXT.md H1 only | check exits 0 |
| CONTEXT.md H2 without Definition | check exits 1 |
| DEC without Dissent | check exits 0 |
| DEC Dissent citing real OQ | check exits 0 |
| DEC Dissent with no OQ/ASM | check exits 1 |
| new scaffold emits thinking skill | file present under `.agents/skills/thinking/SKILL.md` |
| `tier` / `index` do not retrofit thinking | absence preserved |
| MS-301 | Appendix E; both sessions; `go test ./...` |
| no `gh` | fake runner or no execrun use; `gh` never invoked |

Do not require live GitHub for `go test ./...`. Do not commission a `GH_TOKEN` job. Do not add an Actions job as the MS-301 gate.

## 15. Acceptance matrix / in-repo contract paths

Slice 0 lands the paths listed in §12 Slice 0. Later slices also land: updated question template/schema; `program/skills/thinking/SKILL.md`; updated `program/skills/mycelium-cli/SKILL.md`; updated `program/skeleton/AGENTS.md`; `internal/sparring`; check extensions; `internal/clitest/ms301_hermetic_test.go`. No workflow file. No DEC-015 file.

### Acceptance matrix rows (copy into `PHASE-03-acceptance.md`)

Each row: id, check, evidence, owner (Engineering | Arvo). **CI is not an owner.** Robert waived CI.

| id | check | evidence |
| --- | --- | --- |
| A-S0 | Slice 0 files exist and match this brief | PR diff; Quality read |
| A-S1 | schema drops always-required Crux; parser table tests green; new question template omits Crux/Reasons | `go test` + file read |
| A-S2 | IFF check rules; disputed fail/pass; aligned pass; invalid fail; fill-Positions pass; spark-zero-OQ pass | hermetic `go test` |
| A-S3 | CONTEXT.md glossary rules | hermetic `go test` |
| A-S4 | optional Dissent; existing DECs green | hermetic `go test` |
| A-S5 | thinking skill emitted on new scaffold; no retrofit; credits present; mycelium-cli + AGENTS.md name it | hermetic `go test` + file read |
| A-S6 | MS-301 two fixtures green | `go test ./...` |
| MS-301 | all §13 expected bullets | A-S6 (uses 1–2) |

No DOGFOOD row. No Actions-job row. Quality should refuse a PR that adds an Actions job as the MS-301 gate.

## 16. Decided / Architect defaults

No TBD. Open items are decided inline and labeled **Architect default**. Do not bikeshed them in a code PR. Write a DEC to change one. No DEC-015 is required for these. DEC-007 already exists.

Index of defaults that are easy to miss:

- **No new top-level CLI verb.** Surface = existing `new question` / `new assumption` / `new decision` / `new spike` / `check` + the `thinking` skill.
- Disagreement record lives **ON the OQ file**. No `DSG-###`. No new namespace.
- `agreement` enum stays `open` / `aligned` / `agree-to-disagree`.
- Always-required H2s: Question, Context, Disposition. Positions always required as a container (body not graded).
- IFF `agree-to-disagree`: Positions + Reasons + Crux, each with H3 `Human` and H3 `Agent`. Missing any → FAIL. Teaching error names the missing heading and `program/contracts/sparring.md`.
- IFF `aligned` or `open`: Crux and Reasons **not** required. Extra headings do not fail. Aligned MS-301 session: **no disagreement record required**.
- `aligned` and `agree-to-disagree` are terminal-by-convention (skill + contract). Check does not keep agreement history.
- No `session.md`. A session is the set of OQ files in a fixture instance.
- Sparring is always-on (skill) in any non-archived state. Check does not require any OQ. Spark with zero questions still passes.
- Do not add a required `index.md` H2. PHASE-02 H2s stay.
- `CONTEXT.md`: H1 `# Glossary` stays. Empty (H1 only) legal. Any H2 ⇒ H3 `Definition`. Do not require N terms. Do not grade definitions. No new required-file bind.
- Assumption audit is skill-only via `mycelium new assumption`. No AUDIT type.
- Decision template: HTML comment hint, **not** a live `## Dissent` heading. If the heading exists, at least one resolvable `OQ-###` or `ASM-###`. Do not require Dissent on existing DECs.
- Skill path `program/skills/thinking/SKILL.md`, front matter `name: thinking`. New scaffolds only. No retrofit command. Credits: mattpocock grilling (recommendation per question; decisions stay the user's); domain-modeling (glossary + ADR); pstack poteto candor ("no is an acceptable answer"); DEC-007. Do not copy upstream files into `program/`.
- Crux eligible for `mycelium new spike` or a research stage. No auto-promotion command.
- Slice 1 = template/schema/contract + parser (no IFF check bind). Slice 2 = check IFF. Slice 3 = glossary. Slice 4 = Dissent. Slice 5 = skill + embed. Slice 6 = MS-301 in `go test ./...`.
- Schema cannot express IFF; `required_sections` drops `Crux`; no Reasons in the schema list; no schema DSL.
- New-question template omits Crux/Reasons. Default `agreement = "open"`.
- H3 names exact `### Human` / `### Agent`. Section body = bytes after parent H2 until next H2 or EOF.
- House floor: do not rewrite `internal/{cli,check,scaffold,generate,schema,journal,lock,op,lifecycle,statecmd,wakebrief,indexmd,statuscmd,slug,clock,execrun}`. Extend check and generate/schema. New package `internal/sparring` only. Do not touch `internal/slug` (DEC-014).
- `methodology_version` stays `2.0.0`. CLI version stays `0.1.0-dev` unless already different on the pin.
- Log ops: do not add `think` / `spar`. `new` already covers `mycelium new question`.
- `MYCELIUM_NOW` is optional / not required for MS-301.
- `partial: local-only` is NOT required this phase (that's PHASE-02).
- Done bar is hermetic `go test ./...` only. Do not add `.github/workflows/phase-03-*.yml`. Do not extend `phase-01-hermetic.yml` as a phase gate. Quality should refuse a PR that adds an Actions job as the MS-301 gate. Quality should **not** refuse a missing PHASE-03 workflow (absence is correct).
- Do not commission PHASE-04–06. `handed-off` stays unreachable.
- Do not record DEC-015. Do not reopen DEC-012 / DEC-013 / DEC-014.

## 17. Risks, rollback, what Quality should refuse

### Risks and mitigations

| Risk | Mitigation |
| --- | --- |
| Check grades Positions prose or scores substantive vs bare | DEC-005. Container-only. Quality refuse. |
| `## Crux` required on `aligned` (aligned fixture cannot pass) | IFF table. Hermetic aligned case. Quality refuse. |
| New CLI verb (`think` / `spar` / `session`) | §8. Quality refuse. |
| New `DSG` namespace | §5. Quality refuse. |
| Required `index.md` H2 | §9. Quality refuse. |
| OQs required on spark | Spark-zero-OQ test. Quality refuse. |
| Schema IFF DSL or Crux left always-required | Slice 1 schema drop. Quality refuse a leftover always-required Crux after Slice 1. |
| Live `## Dissent` on the DEC template fails every new DEC | Comment hint, not a live heading. |
| PHASE-04 ladder / PHASE-06 packet shipped as a hook | §1 / §2. Quality refuse. |
| Actions job added as the MS-301 gate | Robert waived CI. Quality should refuse. Absence of a PHASE-03 workflow is correct. |
| `GH_TOKEN` job / MS-101(b) reopened | Leftover stays leftover. Quality refuse. |
| Growing `latinFold` or adding `x/text` | DEC-014. Do not touch `internal/slug`. |
| Emitting `framework/` or deleting Justfile | Absence tests stay. Master v1 stays. |
| CLI git-commits instance work product | DEC-010. Quality refuse. |
| Content-scoring "both positions" wording | Tests assert headings + exit codes, not words. |
| Flipping `handed-off` to legal | Command refuse + check FAIL stay. |

### Rollback

Revert the offending PR on master. Do not `git push --force` to main. Floor is `c9ad97874a94e314ceaf1eb2245e433175990f64`. Do not delete Justfile/scripts as a "cleanup" rollback.

### Quality should refuse

Refuse to approve if:

- an Actions job is added as the MS-301 gate, or `.github/workflows/phase-03-*.yml` is added, or `phase-01-hermetic.yml` is extended as a phase gate (Quality should **not** refuse a missing PHASE-03 workflow — absence is correct)
- a new CLI verb is shipped (`think`, `spar`, `session`, `council`, `handoff`, `supersede`)
- a `DSG` (or any new) namespace is added
- a required `index.md` H2 is added
- positions / reasons / cruxes / substantive-vs-bare are content-scored
- check requires any OQ on spark (or any state)
- check requires `## Crux` or `## Reasons` when `agreement` is `aligned` or `open`
- PHASE-04 ladder / second-opinion / council is shipped
- PHASE-06 packet is shipped, or `handed-off` succeeds without a packet, or check stops failing stored `handed-off`
- a `GH_TOKEN` job is commissioned as a gate, or MS-101(b) is implemented
- `framework/` is emitted into an instance
- CLI git-commits instance work product
- `latinFold` grows or NFKD is implemented (DEC-014)
- cobra / viper / yaml / testify / go-github / `golang.org/x/text` appears
- Justfile/scripts deleted from master, `just init` was run on master, or `research-program.toml` was renamed
- DEC-012 / DEC-013 / DEC-014 reopened in a code PR
- DEC-015 is recorded
- a `session.md` type or AUDIT type appears
- a retrofit-skills command appears
- hermetic tests call network or real `gh`
- PR pushed straight to main
- seven-real-day dogfood is required as the gate

## 18. Execution order (PR-per-slice)

Same order as §12 (slices 0→6). PR-per-slice, sequential, rebase on main. One live PR at a time. Slice 2 is the check bind — do not combine it with 5–6. Slice 6 must be green in `go test ./...` on its PR (not Actions).

Title: `PHASE-03 Slice N: <done-bar noun>`. Body links this brief and the slice done bar. No drive-by refactors, v1 deletions, PHASE-01/02 leftover work, or PHASE-04+ commands. Engineering opens PRs; Arvo merges Quality-green PRs; Engineering does NOT push to main.

Cursor cloud env name is exactly `robertguss/mycelium`. Go 1.26. `CGO_ENABLED=0`.

## 19. Handoff

### What Engineering starts with

This file. Only this file. Start from `https://github.com/robertguss/mycelium` at `c9ad97874a94e314ceaf1eb2245e433175990f64` (current `main` at pin time). Read `framework/blueprint.md` and DEC-001–014 for authority, not for a second plan. Execute Slice 0 first. Do not implement from a later SHA unless Arvo re-pins.

Cursor cloud: env `robertguss/mycelium`. Go 1.26. `CGO_ENABLED=0`.

### What Engineering must not do

See §17 (Quality should refuse) and §16 (Architect defaults). Do not open a design debate in the PR. Do not write a second brief. Do not write DEC-015. Do not commission PHASE-04–06. Do not add an Actions job as the MS-301 gate. Do not start from a later SHA and "fast-forward" this brief.

### What Quality reads

This brief, the acceptance matrix, `program/contracts/sparring.md`, `program/contracts/glossary.md`, the conformance delta, and the PR diff. Thermos: §14 tests exist and match; §17 refuse list is clean; MS-301 hermetic `go test ./...`; no `GH_TOKEN` gate; no Actions job as the MS-301 gate; no new CLI verb; no `DSG` namespace.

Quality should **not** refuse a missing PHASE-03 workflow (absence is correct). Quality should refuse a PR that adds an Actions job as the MS-301 gate.

### What Arvo does

Merges Quality-green PRs. Accepts PHASE-03 when MS-301 is green on main. Does not re-pin without writing the new SHA. Does not record DEC-015.

## Appendix A — No new DEC (DEC-007 already exists)

No DEC-015. PHASE-03 implements **DEC-007** (Accepted 2026-08-14). PHASE-03 does not reopen DEC-012, DEC-013, or DEC-014. Remaining choices are Architect defaults in §16. Engineering lands **zero** new files under `framework/decisions/`.

If a later discovery contradicts a locked item, stop and write a DEC; do not silently patch this brief in a code PR.

## Appendix B — Disputed OQ example

Fixture shape after `mycelium new question "Use SQLite"` plus edits. Tests must **not** grade the words under the H3s — only headings, `agreement`, and `mycelium check` exit 0. The record retains **both positions**, both sets of reasons, and **cruxes**.

`questions/OQ-001-use-sqlite.md`:

```text
+++
id = "OQ-001"
title = "Use SQLite"
agreement = "agree-to-disagree"
date = "2026-08-15"
+++

# OQ-001 — Use SQLite

<!-- slug: use-sqlite -->

## Question

Should this idea use SQLite as the store?

## Context

Focused fixture. Words are not graded.

## Positions

### Human

Use SQLite.

### Agent

Do not use SQLite.

## Reasons

### Human

One file. Enough.

### Agent

Need a server later.

## Crux

### Human

A second writer appears.

### Agent

The idea stays single-process.

## Disposition

agree-to-disagree. Stop. Open a new OQ if the crux fires.
```

Negative case (same test file): delete the `## Crux` heading and its H3s → `mycelium check` exits 1; stderr names `## Crux` and `program/contracts/sparring.md`.

## Appendix C — Aligned OQ example

Honestly aligned session. Positions H2 present. **No** Crux/Reasons headings. **No disagreement record required.** `mycelium check` exits 0.

`questions/OQ-001-keep-the-name.md`:

```text
+++
id = "OQ-001"
title = "Keep the name"
agreement = "aligned"
date = "2026-08-15"
+++

# OQ-001 — Keep the name

<!-- slug: keep-the-name -->

## Question

Keep the fixture idea name?

## Context

Aligned fixture. Words are not graded.

## Positions

<!-- fill -->

## Disposition

aligned.
```

A one-word or `<!-- fill -->` Positions body does **not** fail. Extra `## Crux` / `## Reasons` (if someone adds them) do **not** fail.

## Appendix D — CONTEXT.md glossary example

### Empty (legal)

```text
# Glossary
```

`mycelium check` exits 0. Check does not require N terms.

### One term (legal)

```text
# Glossary

## SQLite

### Definition

<!-- fill -->
```

`mycelium check` exits 0. Definition body is not graded.

### One term missing Definition (FAIL)

```text
# Glossary

## SQLite
```

`mycelium check` exits 1. Teaching error names `## SQLite` and `program/contracts/glossary.md`.

Skill (not check) challenges drift and vagueness and records the term on resolution. Credit domain-modeling in the skill.

## Appendix E — MS-301 fixture recipe

Hermetic. No network. No `gh`. No `GH_TOKEN`. `go test ./...` only. Do **not** add an Actions job. Do **not** add `.github/workflows/phase-03-*.yml`. Do **not** extend `phase-01-hermetic.yml` as a phase gate.

Two hermetic temp instances. Injectable nothing required unless you reuse `MYCELIUM_NOW` for dates on generated files (`MYCELIUM_NOW` is optional / not required).

Work in temp dirs. Binary = freshly built `mycelium`. Edits may use the test helper / stdlib. Do not add a `mycelium edit` command.

**Architect default:** test lives at `internal/clitest/ms301_hermetic_test.go` (execs the binary). Parser tests live in `internal/sparring`.

### Disputed instance

```text
mycelium new idea "Disputed Fixture" --offline --dir PATH/disputed-fixture
# focused default is fine; state = spark

mycelium new question "Use SQLite" --dir PATH/disputed-fixture
# questions/OQ-001-use-sqlite.md

# EDIT front matter: agreement = "agree-to-disagree"
# EDIT ## Positions to contain ### Human and ### Agent (any non-empty text)
# ADD ## Reasons with ### Human and ### Agent (any non-empty text)
# ADD ## Crux with ### Human and ### Agent (any non-empty text)
# Keep ## Question, ## Context, ## Disposition (template already has them)

mycelium check --dir PATH/disputed-fixture
# exit 0
```

Negative case in the **same test file** (copy of the disputed instance, or a second temp dir):

```text
# same shape, then DELETE ## Crux (and its H3s)
mycelium check --dir PATH/disputed-missing-crux
# exit 1
# stderr contains "## Crux" and "program/contracts/sparring.md"
```

### Aligned instance

```text
mycelium new idea "Aligned Fixture" --offline --dir PATH/aligned-fixture

mycelium new question "Keep the name" --dir PATH/aligned-fixture

# EDIT front matter: agreement = "aligned"
# ## Positions H2 present (template already has it)
# Do NOT add ## Crux
# Do NOT add ## Reasons

mycelium check --dir PATH/aligned-fixture
# exit 0
```

### Also in the same `go test ./...` (may share the file)

| Case | Expect |
| --- | --- |
| Edit `agreement = "maybe"` (or any non-enum) | `mycelium check` exits 1 |
| Spark scaffold, zero questions, no OQ files | `mycelium check` exits 0 |
| `agreement = "open"` or `"aligned"`; Positions body is `<!-- fill -->`; no Crux/Reasons | `mycelium check` exits 0 |

Tests must not grade the words under Human/Agent. `gh` never invoked. No seven-real-day dogfood.

## Appendix F — Target file tree additions (and explicit DO NOT ADD workflow files)

### Master (additions on top of the PHASE-01+02 tree; v1 files retained)

```text
internal/sparring/                         # agreement-conditional headings + glossary helpers + table tests
internal/clitest/ms301_hermetic_test.go
program/contracts/sparring.md              # new
program/contracts/glossary.md              # new
program/contracts/conformance.md           # items 15–17 + lift timing
program/templates/question.md              # updated (no always-emitted Crux/Reasons)
program/templates/question.schema.toml     # Crux dropped from required_sections
program/templates/decision.md              # optional Dissent HTML comment
program/skills/thinking/SKILL.md           # new
program/skills/mycelium-cli/SKILL.md       # updated
program/skeleton/AGENTS.md                 # updated
framework/phases/PHASE-03-implementation-brief.md
framework/phases/PHASE-03-acceptance.md
internal/embed/program/                    # regenerate after program/ edits
```

### DO NOT ADD

```text
.github/workflows/phase-03-*.yml
.github/workflows/phase-03-hermetic.yml
.github/workflows/phase-03-ms301.yml
framework/decisions/DEC-015-*.md
program/templates/session.md
program/templates/disagreement.md
program/templates/audit.md
a DSG namespace or questions/ sibling home for disagreements
internal/slug changes
think / spar / session / council / handover / handoff command packages
```

Do **not** add a PHASE-03 workflow. Do **not** extend `phase-01-hermetic.yml` as a phase gate. Quality should refuse a PR that adds an Actions job as the MS-301 gate. Quality should **not** refuse a missing PHASE-03 workflow (absence is correct). Do not delete Justfile/scripts/`research-program.toml`/PHASE-01 workflows.

### Emitted instance (spark / focused, local-only, PHASE-03 scaffold)

```text
README.md  mycelium.toml  log.md  index.md  CONTEXT.md  AGENTS.md  .gitignore
.agents/skills/{mycelium-cli,spark,wake,portfolio,thinking}/SKILL.md
program/ …
.git/          # init only; no commit
```

Absent: `framework/`, `cmd/`, `internal/`, `go.mod`, `Justfile`, `scripts/`, `research-program.toml`, v1 `research-*` skills, `session.md`, `DSG-*` files.

After `mycelium new question` (Slice 1+ template): `questions/OQ-001-*.md` with Question / Context / Positions / Disposition; `agreement = "open"`; no Crux/Reasons unless the agent adds them.

Unexported helpers may live next to their tests. No `pkg/`. No extra public command packages. Do not touch `internal/slug`.

End of PHASE-03 implementation brief. Engineering executes from this file only.
