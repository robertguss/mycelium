# Mycelium — Framework Blueprint

- **Status:** Proposed — awaiting acceptance by Robert Guss
- **Date:** 2026-08-14 (revised 2026-08-14)
- **Revision:** 1 — applies all thirteen accepted findings of the
  [blueprint adversarial review](reviews/01-blueprint-adversarial-review.md) per
  its [Dispositions](reviews/01-blueprint-adversarial-review.md#dispositions)
- **Authors:** Robert Guss with his agent (discovery interview, one question at
  a time per `program/reference/discovery-protocol.md`)
- **Methodology version targeted:** 2.0 (evolving v1.0 in place)
- **Decision records:** `framework/decisions/DEC-001` through `DEC-011`

This blueprint governs the evolution of this repository from the Artifact-Driven
Research Program (v1) into **Mycelium**, a convention-over-configuration
framework for thinking. It defines what will be built and in what order. It does
not conduct the work itself.

## Vision

Mycelium is to thinking what Rails is to web applications: not a template you
copy but a tool you run. One master repository builds a single-binary CLI
(`mycelium`) that carries the conventions, contracts, templates, checks, and
skills, and scaffolds each idea its own instance repository, thin at birth,
growing rigor in place as the idea earns it. A human and their agents spar
inside it, and everything that matters — terms, assumptions, positions,
decisions, dissent, evidence, and time itself — lands in git as auditable
artifacts.

The name: mycelium is the underground network that does the real work long
before anything fruits above ground. Ideas simmer the same way.

Mycelium is the precursor to implementation. Its terminal artifact is a handoff
packet that an implementation system (pstack / poteto-mode, or any other) can
pick up and build from.

## Problem statement

Robert's ideation today happens across chat sessions that are not authoritative,
with no durable capture of disagreements, assumptions, or the reasoning behind
decisions. Sessions span hours, days, and weeks; ideas need to simmer; both
human and agent need to get back up to speed on return. Single-model
conversations also share a single set of blind spots, and the existing research
methodology (v1 of this repo) is model-blind and program-scale only — there is
no home for a thirty-minute spark.

## Intended users

1. Robert, thinking across many ideas in parallel over long horizons.
2. His agents, in any runtime: Cursor today, Claude Code, cloud VMs, and
   whatever ships next month. The framework must not assume any one of them.
3. Later, anyone who clones the template.

## Goals

1. Every idea gets a self-contained, git-auditable home in minutes.
2. The agent is a sparring partner with mandatory positions, not an interviewer
   or a yes-man.
3. Disagreement is captured as first-class gold: both positions, both sets of
   reasons, and the crux that would change each mind.
4. Multi-model perspective is available on demand and off by default.
5. Ideas can simmer for weeks and wake warm.
6. Conventions make the structure boring so the thinking can be interesting.
7. Deterministic checks keep agents honest about structure.
8. Everything works in any agent runtime that can read files and run git.

## Non-goals

- A hub or garden repository aggregating multiple ideas (rejected, DEC-003).
- A standalone council web app (llm-council is a reference, not a product).
- A Cursor plugin as the system of record (rejected, DEC-004).
- Deterministic grading of _thinking quality_ — checks validate containers,
  never contents (DEC-005).
- A granular coding backlog for implementation — handoff stops at the packet,
  per v1's "stops before a granular coding backlog."
- Migration machinery for existing instances (deferred until a real need
  arrives, DEC-011).
- A GitHub "Use this template" flow (superseded by the CLI, DEC-010).

## Locked constraints

1. Multi-runtime portability. Thinking requires only files + git; the CLI
   accelerates but never gates it (the CLI creates, never thinks).
2. Self-contained instances. Clone = everything an agent needs, scaffolded in at
   creation, not fetched at runtime.
3. One repository per idea. Rigor tiers vary inside a constant spine.
4. Councils are opt-in, never required, cost-visible before running.
5. Humans own git. Agents write files; humans (or explicitly delegated agents)
   commit. The CLI never commits work product.
6. Methodology is version-pinned per instance at scaffold time, recorded in the
   manifest as `methodology_version`, distinct from `generated_by_cli_version`.
   Runtime commands treat instance files as truth, never the binary's embedded
   copies, so old instances keep working under new binaries. Migration machinery
   is deferred (DEC-011).
7. v1's authority rules stand: git-tracked artifacts are authoritative; chat and
   model memory are not.

## Success criteria

1. Spark-to-first-captured-thought under five minutes, including repo creation.
2. Re-entry after a multi-week gap: one wake brief restores working context for
   human and agent without rereading raw logs.
3. Conformance suite green means: no orphaned IDs, no illegal state transitions,
   no dissent record missing a crux, no simmering idea missing a revisit
   trigger.
4. A handoff packet is sufficient for a fresh implementation agent with no
   access to this repo's chat history.
5. The framework's own maintenance cost stays below the thinking it saves; if
   meta-work dominates, subtract (Laziness guard).

## Lineage and sources

| Source                 | What Mycelium takes                                                                                                                                                                                                                        |
| ---------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| ADRP v1 (this repo)    | Governance spine, contracts, evidence model, rigor tiers, replication/reconciliation rules, ID system, fresh-session rule, human gates                                                                                                     |
| pstack (poteto)        | Model-role configuration, sparring candor ("no is an acceptable answer"), arena/interrogate mechanics as council prompt content, second-opinion doctrine                                                                                   |
| karpathy/llm-council   | Three-stage council shape: independent first opinions, anonymized cross-review, synthesis — subordinated to v1's reconciliation rules (no majority vote, dissent retained)                                                                 |
| karpathy llm-wiki gist | index.md + log.md navigation, LLM-as-bookkeeper, decisions as living objects with falsifiable assumptions (Vigil)                                                                                                                          |
| mattpocock/skills      | Grilling's frontier-of-decisions interview with a recommendation per question; wayfinder's fog-of-war and destination-first charting; domain-modeling's glossary + ADR discipline; handoff's compaction ethos                              |
| Ruby on Rails          | Convention over configuration as identity; derivable naming; generators; timestamped migrations; environments (as rigor-tier config); engines (as presence-registered packs); dummy-app testing (as the fixture instance); teaching errors |

## Domain model

### Idea lifecycle (state machine, stored in the instance manifest)

```text
spark → exploring ⇄ simmering → clarified → handed-off
any state → archived
```

- **spark** — repo exists, framing may not. Passes conformance nearly empty.
- **exploring** — active sessions; frontier of open questions being worked.
- **simmering** — deliberately parked. Requires a revisit trigger (date or
  event). Not fog (question unstateable) and not blocked (waiting on a
  decision); simmering means "could decide now, choosing not to."
- **wake** is the simmering → exploring transition. Its ritual: reread index and
  log tail, check evidence revalidation triggers and assumption records against
  what changed, then brief the human.
- **clarified** — destination reached; handoff packet buildable.
- **handed-off** — packet delivered to an implementation system.
- **archived** — dead or absorbed; record preserved, never deleted.

Conformance validates transitions (e.g., no handed-off without a packet; no
simmering without a revisit trigger).

### Objects

- **Term** — glossary entry (`CONTEXT.md` per domain-modeling conventions).
  Agents challenge drift against it.
- **Assumption** (`ASM-###`, promoted from v1-optional to required where
  decisions exist) — falsifiable, attached to decisions and ideas, checked at
  wake and on new evidence.
- **Position** — a stance held by a named party (human, resident agent, or a
  council member) with reasons. Every substantive question gets the agent's
  position; bare questions are a conformance smell in sparring transcripts.
- **Question** — carries an agreement state: `open` | `aligned` |
  `agree-to-disagree`. The third state is terminal and honorable.
- **Disagreement record** — both positions, both reasons, and the **crux**: what
  evidence would change each mind. Cruxes are eligible to become research stages
  or spikes.
- **Decision** (`DEC-###`) — v1 template, plus assumptions listed and dissent
  retained inline.
- **Evidence** (`EVD-###`, `SPK-###`) — v1 evidence model unchanged,
  revalidation triggers now consumed by the wake ritual.
- **Council** — a commissioned replication run whose replicas are _different
  models_. Reuses v1's replication and reconciliation contracts and directories
  verbatim; reports land per-model, reconciliation retains dissent, and
  selection by majority vote or model reputation stays banned.
- **Handoff packet** — the terminal artifact: framing, locked decisions,
  glossary, open questions with agreement states, evidence summary, suggested
  implementation playbooks.

## Repository anatomy

### Master (this repo)

```text
cmd/, internal/ Go CLI source (single static binary: mycelium)
program/        methodology: contracts, templates + sidecar schemas, tiers,
                packs, skills, instance skeleton — normal browsable files,
                go:embed'd into the binary, emitted on scaffold
framework/      Mycelium's own blueprint + decisions (never emitted)
```

The v1 `Justfile` and `scripts/` retire once their duties move into the CLI;
instances carry neither. The emitted contracts plus the binary are the whole
operational surface, and a shipped skill teaches agents both the commands and
the manual floor.

### Instance (per idea)

`mycelium new idea "name"` scaffolds it: README, manifest, `log.md`, empty
`CONTEXT.md`, `AGENTS.md`, `.agents/skills/` (the emitted skills, teaching the
commands and the manual floor), and `program/` contracts, templates, and
schemas, then git init, `gh repo create`, and the `idea` topic.

The manifest separates two orthogonal fields: lifecycle **state** (`spark` at
birth, transitions per the DEC-006 machine) and rigor **tier** (per DEC-002).
Any state may hold any tier. It also records two distinct version fields,
`methodology_version` (DEC-004's pin) and `generated_by_cli_version` (DEC-011's
seam), with no one-to-one release invariant declared between them.

Structure is emitted per tier, never per state: a low tier carries little, and
deeper docs/ trees arrive when a tier that requires them is set. (v1's full tree
of placeholders existed to compensate for indiscriminate template copying; the
CLI ends that.) Tier transitions are legal in both directions and independent of
state transitions. `mycelium tier <tier>` (PHASE-01) is the idempotent promotion
operation: it updates the manifest and emits only the structure the new tier
newly requires, never overwriting or rewriting existing work; lowering a tier
deletes nothing (no-deletion rule) and only relaxes which artifacts bind.
Conformance requirements are **tier-aware**: a minimal tier passes with almost
nothing; ledger discipline binds only at the tier that demands it.

Of the emitted files, the portable authority is `program/` (contracts,
templates, sidecar schemas, tier config) plus the manifest and `AGENTS.md`.
`.agents/skills/` is the canonical emitted runtime-adapter path — skill source
lives under the master's `program/` tree — and it, like any further per-runtime
adapter files, is a generated convenience, always subordinate to the portable
contracts.

## Conventions and conformance

Adopted 2026-08-14: the full Rails borrowing set (items 1–8 below plus the
declared-deviation rule), approved by Robert as a block.

1. Naming, IDs, and layout are derivable, not searchable. The ID-to-path rule is
   a pure function: every namespace has exactly one home directory and one
   filename pattern (`DEC-014` → `decisions/DEC-014-<kebab-slug>.md`).
   Directories are plural lowercase nouns, one artifact type each; IDs are
   `UPPER-###`; slugs are kebab-case. Conformance enforces the mapping both
   directions: files must match their pattern, references must resolve.
2. Rigor tiers are machine-readable config (`program/tiers/*.toml`), not prose.
   The manifest pins a tier; the checker reads tier config to learn which
   artifacts bind. (Rails environments, translated.)
3. Generators (`mycelium new <type> "Title"`) create the full bundle: file from
   template, next ID allocated, manifest entry when stage-tracked, log line
   appended — all under the operation protocol (see Templates and generation),
   which makes an interrupted bundle detectable and recoverable instead of
   claiming half-created artifacts are impossible. There is no destroy;
   state-transition commands (`mycelium supersede`) honor the no-deletion rule.
4. No migrations (DEC-011). Instead: every manifest records
   `methodology_version` and `generated_by_cli_version` as separate fields,
   runtime commands validate against the instance's own emitted schemas
   (instance files are truth), and the master keeps a `CHANGELOG.md` as each
   release's human-readable face.
5. Packs are presence-registered: a directory under `program/packs/<name>/`
   containing its own templates, contracts, and checks. Drop in to enable,
   remove to disable, no registry file. Conformance fails namespace collisions
   between packs.
6. The framework tests itself against a fixture instance: CI scaffolds a
   throwaway instance with the freshly built binary, generates one of every
   artifact type, and runs the checks. Convention breakage surfaces in the
   master's CI, never first in a real idea repo. (Rails engines' dummy app,
   translated.)
7. Every template ships a sidecar schema; generator and checker are both
   data-driven off it (see Templates and generation).
8. Declared deviation: an instance may deviate from a convention only by
   declaring the deviation in its manifest, visibly and auditably. Silent
   deviation is a conformance failure even when the deviation itself would be
   acceptable.
9. The conformance suite validates **structure only**: ID uniqueness and
   sequence, ID-to-path integrity, link resolution, required front matter and
   sections per schema, legal state and tier transitions, tier-appropriate
   artifact presence, parseable log prefixes, crux presence on disagreement
   records, and revisit triggers on simmering ideas. Checks never grade content;
   thinking quality is judged by adversarial review, councils, and the human
   (Goodhart guard).
10. Failures teach: every error names the violated convention, links its
    contract, and suggests the fixing command (did-you-mean style).
11. CI runs the suite on push. Green = the record is structurally trustworthy.

## Templates and generation

Every artifact type is defined by exactly two files, side by side in
`program/templates/` (or a pack's `templates/`): the markdown template and a
sidecar schema (`<type>.schema.toml`) declaring the ID namespace, home
directory, filename pattern, required front-matter fields, and required
sections.

The schema is the shared truth. The generator reads it to know what to create
and where; the checker reads the same file to know what to validate. The two
tools cannot disagree, and adding an artifact type means adding a template and a
schema, with zero code changes to either tool.

Mechanics:

- Templates and schemas live as browsable files under `program/` in the master
  and are `go:embed`'d into the binary. Scaffolding emits them into the
  instance, where they become that instance's frozen truth.
- Templates use `{{ID}}`, `{{TITLE}}`, `{{SLUG}}`, `{{DATE}}` tokens. Plain
  string replacement in Go's stdlib, no templating engine.
- `mycelium new <type> "Title"` allocates the next ID by scanning the type's
  home directory. The filesystem is the registry: because ID-to-path is a pure
  function, a directory scan is unambiguous and there is no central index file
  to rot. The generator refuses to overwrite, fills tokens, writes the file,
  appends the log entry, updates the manifest for stage-tracked types, prints
  the path plus next steps, and never runs git — all under the operation
  protocol below.
- Runtime commands (`new`, `check`, `status`) read the instance's own emitted
  templates and schemas, never the binary's embedded copies. A new binary
  therefore generates and validates old instances by their own rules; embedded
  data is used only at scaffold time.
- `mycelium supersede <ID> --by <ID>` flips status, wires cross-links into both
  records, and logs the transition.
- Stage-scoped namespaces (`REC`/`REQ`/`FND`) respect the ranges the blueprint
  and manifest allocate; enforcement strictness is OQ-007.
- Artifact metadata in emitted templates is **front matter**, defined and
  validated by each type's sidecar schema (DEC-005 item 3 as written). One
  parser contract — the **metadata reader**, implemented once in the CLI — is
  consumed by both the generator and `mycelium check`, and it bounds
  front-matter and section detection so body text cannot masquerade as metadata.
  Machine lifecycle state stays in the manifest per the manifest-authority rule.
  This grammar governs what Mycelium emits and validates; the master's own
  `framework/` artifacts are not retrofitted.

Operation protocol (all multi-file commands: `new`, `supersede`, `tier`):

1. **Preflight** — validate that the manifest and log parse, the schema
   resolves, and the target path is free, before anything is written.
2. **Lock** — take an exclusive repository lock for the duration, so concurrent
   generators cannot allocate the same ID.
3. **Stage** — write outputs as temporary files and record the operation's
   intent in a journal.
4. **Commit** — atomic renames where the filesystem supports them, in a fixed
   order: artifact file, then log, then manifest.
5. **Rollback and retry** — failure before the first rename removes the staged
   files and changes nothing; after a partial commit the journal survives, and
   re-running the command resumes the journaled operation under the original ID
   rather than allocating a new one.
6. **Detection and recovery** — `mycelium check` detects a leftover journal or
   stale lock as an interrupted operation and names the documented recovery step
   (complete or roll back), did-you-mean style.

The supported filesystem floor is a local filesystem with atomic rename; network
filesystems are outside the floor, where the lock and journal still bound the
damage.

## Perspective ladder

1. **Sparring** (free, always on) — the resident agent holds positions,
   challenges terms, surfaces presuppositions (the assumption audit), and
   recommends on every substantive question.
2. **Second opinion** (cheap, one word to invoke) — the identical commissioning
   prompt to exactly one different model. Agreement is high-signal; disagreement
   surfaces a fork.
3. **Council** (expensive, opt-in, suggested only when v1's replication triggers
   fire: hard to reverse, weak or conflicting evidence, still low-confidence
   after a spike) — full multi-model replication + reconciliation. Cost class
   stated before running. Panel presets live in user-level config (quick /
   standard / high-stakes).

Engine-agnosticism: contracts define commissioning prompts and report file
shapes; _how_ reports get produced is a swappable adapter. Cursor's parallel
multi-model subagents are the first adapter. The manual floor — pasting the
commissioning prompt into N chat UIs and saving N files — satisfies the contract
with zero tooling. Runtimes that cannot fan out simply skip ladder rungs 2–3;
`AGENTS.md` carries the capability note.

## Cross-idea operations

No hub repo. Each instance manifest carries `state` and `revisit`; the
scaffolder applies the `idea` topic on GitHub automatically (DEC-003's topic
hygiene, now enforced by tooling instead of memory).

`mycelium status --all` answers "what's simmering, what's due to wake" by
querying GitHub live, per DEC-003: enumerate repositories carrying the `idea`
topic (via `gh`), read each manifest (remotely for repos not cloned locally),
and merge a scan of the local ideas root (`~/ideas/<slug>` by convention). Every
case is handled explicitly: remote-only repos appear via the topic query;
local-only repos (created but not yet published, or missing the topic) appear
from the local scan, flagged as unpublished; archived ideas are filtered from
the default view and available behind a flag; when GitHub is unauthenticated or
temporarily unavailable (outage, rate limit), the command degrades to the
documented local fallback — the local scan alone — and marks the output partial,
never silently incomplete. The first usable version ships in PHASE-02 alongside
the thin portfolio skill that wraps the command for agents; PHASE-05 hardens it
with tolerance for older manifest shapes (DEC-011's risk guard).

## Phases and milestones

Sequential; each phase ends verifiable. Stops at milestones, per house rule.
Milestones are stated at blueprint-wording granularity; each phase's full
acceptance matrix lands in that phase's contract when the phase is commissioned,
per the
[review dispositions](reviews/01-blueprint-adversarial-review.md#dispositions),
so "verifiable" keeps a stable meaning without the blueprint becoming a test
plan.

- **PHASE-01 Foundation.** Go CLI skeleton; `program/` content authored and
  `go:embed`'d; naming and ID-to-path contracts; template sidecar schemas and
  the front-matter metadata reader; tiers as machine-readable config;
  `mycelium new idea` scaffolding (emit, git init, `gh repo create`, `idea`
  topic); the data-driven `mycelium new <type>` generator and the operation
  protocol; the idempotent `mycelium tier` operation; manifest gains separate
  lifecycle-state and rigor-tier fields plus `methodology_version` and
  `generated_by_cli_version` stamps; `mycelium check` (schema-driven,
  tier-aware, teaching errors, interrupted-operation detection);
  fixture-instance CI (scaffold, generate, check). _MS-101 — two parts. Hermetic
  local: with no network, the binary scaffolds a conformant spark instance, and
  fixture CI generates and checks one of every artifact type. Authenticated
  GitHub integration: a separately credentialed test publishes the repo with the
  `idea` topic and cleans up on failure. Five minutes from nothing to first
  captured thought stays a user SLO, not the phase gate._
- **PHASE-02 Lifecycle.** Spark / wake / portfolio skills; log + index
  conventions; re-entry brief; simmer with revisit triggers; first usable
  `mycelium status --all` (live GitHub enumeration with documented local
  fallback, per Cross-idea operations). _MS-201: against a deterministic fixture
  — injectable clock, known log entries, evidence triggers, and assumption
  changes — a simulated 7+ day simmer wakes with a brief citing the expected
  changes and sources; one dogfood wake after seven real days follows as human
  evidence._
- **PHASE-03 Sparring.** Thinking-mode skill: mandatory positions, agreement
  states, disagreement records with cruxes, glossary challenge, assumption
  audit. Grilling and domain-modeling conventions absorbed and credited.
  _MS-301: two fixture sessions pass — one deliberately disputed, whose record
  retains both positions, both sets of reasons, and cruxes; one honestly
  aligned, which passes with no disagreement record required.
  Substantive-question and bare-question judgment belongs to the human or an
  adversarial reviewer, never an automated content score (DEC-005's containers
  boundary)._
- **PHASE-04 Perspective ladder.** Commissioning + report + reconciliation
  contracts for model-diverse replication; second-opinion move; Cursor council
  adapter; panel presets in user config. Council ships as the first pack,
  proving presence-is-registration. Council-contract dogfood lands here (moved
  from OQ-006), as does commissioning provenance: commissioning prompts stored
  durably in-repo, model provenance recorded, and the ladder rung named for
  every commissioned review. _MS-401: the end-to-end council run is one row of a
  perspective-ladder acceptance matrix covering the DEC-008 contract —
  second-opinion, Cursor-council, and manual-floor rows; explicit opt-in and
  stated cost class; prompt identity and model provenance evidenced; independent
  per-model reports; a seeded dissent surviving reconciliation; council-pack
  enable/disable without touching core checks._
- **PHASE-05 Distribution and lifecycle commands.** `mycelium supersede`; tagged
  releases with prebuilt binaries and install docs; `CHANGELOG.md` discipline;
  portfolio scanner tolerant of older manifest shapes (DEC-011's risk guard).
  _MS-501 — two clauses. Functional acceptance: `mycelium supersede` leaves
  bidirectional cross-links and a log entry; `check` and `status` pass golden
  old-instance fixtures; releases ship checksummed binaries with a
  `CHANGELOG.md` entry. Install SLO: a clean VM (named image) goes from one-line
  install to a scaffolded instance in under a minute._
- **PHASE-06 Handoff.** Packet contract + generator; pstack/poteto bridge
  documented; implementation-systems section in AGENTS.md. _MS-601: a fresh,
  isolated agent — no chat history, no source access beyond the packet —
  implements a canonical packet fixture's bounded target with the documented
  implementation system inside a stated time budget, and the result passes the
  fixture's executable acceptance tests. A separate real-project handoff serves
  as dogfood evidence, not the gate._

## Identifier allocations (framework evolution)

- Framework decisions: `DEC-001`–`DEC-099` (in `framework/decisions/`).
- Blueprint adversarial-review findings: `FND-001`–`FND-099` (`FND-001` through
  `FND-013` consumed by the
  [2026-08-14 review](reviews/01-blueprint-adversarial-review.md)).
- Open questions herein: `OQ-001`–`OQ-019`.

## Open questions

- **OQ-001** Final state-vocabulary bikeshed (`spark/exploring/simmering/...`) —
  settle in PHASE-01 contracts.
- **OQ-002** Manifest filename: keep `research-program.toml` or rename
  (`mycelium.toml`) as a 2.0 migration.
- **OQ-003** Resolved 2026-08-14: packs are presence-registered directories
  under `program/packs/` (see Templates and generation). The remaining sliver,
  which capabilities beyond council become packs, rolls into PHASE-04.
- **OQ-004** Handoff packet format details and the pstack playbook mapping.
- **OQ-005** Resolved 2026-08-14: the master repository is named `mycelium`;
  scaffolded instance repositories receive the `idea` topic automatically (see
  Cross-idea operations).
- **OQ-006** Resolved 2026-08-14: the blueprint received an independent,
  manually commissioned
  [adversarial review](reviews/01-blueprint-adversarial-review.md) by a
  different model — under DEC-008's definitions neither a council nor a second
  opinion, since the commissioning prompt differed from the drafting prompt.
  Council-contract dogfood moves to PHASE-04.
- **OQ-007** Generator strictness for stage-scoped ID ranges: warn or refuse
  when allocating outside a declared range.

## Authority, amendment, completion

This blueprint governs the evolution once accepted. Accepted `DEC-###` records
in `framework/decisions/` supersede it clause by clause. Amendments follow
`program/reference/amendment-protocol.md`. The evolution is complete when all
six phases are accepted and one real idea has traveled spark → handed-off
entirely inside Mycelium.

Fresh-session note: this blueprint is the synthesized output of the 2026-08-14
discovery interview. Per the fresh-session rule, its adversarial review ran in a
fresh session (OQ-006, resolved), and each build phase should run in fresh
sessions with self-contained packets.
