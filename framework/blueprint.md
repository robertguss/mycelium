# Mycelium — Framework Blueprint

- **Status:** Proposed — awaiting acceptance by Robert Guss
- **Date:** 2026-08-14
- **Authors:** Robert Guss with his agent (discovery interview, one question at
  a time per `program/reference/discovery-protocol.md`)
- **Methodology version targeted:** 2.0 (evolving v1.0 in place)
- **Decision records:** `framework/decisions/DEC-001` through `DEC-011`

This blueprint governs the evolution of this repository from the
Artifact-Driven Research Program (v1) into **Mycelium**, a
convention-over-configuration framework for thinking. It defines what will be
built and in what order. It does not conduct the work itself.

## Vision

Mycelium is to thinking what Rails is to web applications: not a template you
copy but a tool you run. One master repository builds a single-binary CLI
(`mycelium`) that carries the conventions, contracts, templates, checks, and
skills, and scaffolds each idea its own instance repository, thin at birth,
growing rigor in place as the idea earns it. A human and their agents spar inside it, and
everything that matters — terms, assumptions, positions, decisions, dissent,
evidence, and time itself — lands in git as auditable artifacts.

The name: mycelium is the underground network that does the real work long
before anything fruits above ground. Ideas simmer the same way.

Mycelium is the precursor to implementation. Its terminal artifact is a
handoff packet that an implementation system (pstack / poteto-mode, or any
other) can pick up and build from.

## Problem statement

Robert's ideation today happens across chat sessions that are not
authoritative, with no durable capture of disagreements, assumptions, or the
reasoning behind decisions. Sessions span hours, days, and weeks; ideas need
to simmer; both human and agent need to get back up to speed on return.
Single-model conversations also share a single set of blind spots, and the
existing research methodology (v1 of this repo) is model-blind and
program-scale only — there is no home for a thirty-minute spark.

## Intended users

1. Robert, thinking across many ideas in parallel over long horizons.
2. His agents, in any runtime: Cursor today, Claude Code, cloud VMs, and
   whatever ships next month. The framework must not assume any one of them.
3. Later, anyone who clones the template.

## Goals

1. Every idea gets a self-contained, git-auditable home in minutes.
2. The agent is a sparring partner with mandatory positions, not an
   interviewer or a yes-man.
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
- Deterministic grading of *thinking quality* — checks validate containers,
  never contents (DEC-005).
- A granular coding backlog for implementation — handoff stops at the packet,
  per v1's "stops before a granular coding backlog."
- Migration machinery for existing instances (deferred until a real need
  arrives, DEC-011).
- A GitHub "Use this template" flow (superseded by the CLI, DEC-010).

## Locked constraints

1. Multi-runtime portability. Thinking requires only files + git; the CLI
   accelerates but never gates it (the CLI creates, never thinks).
2. Self-contained instances. Clone = everything an agent needs, scaffolded
   in at creation, not fetched at runtime.
3. One repository per idea. Rigor tiers vary inside a constant spine.
4. Councils are opt-in, never required, cost-visible before running.
5. Humans own git. Agents write files; humans (or explicitly delegated
   agents) commit. The CLI never commits work product.
6. Methodology is version-pinned per instance at scaffold time. Runtime
   commands treat instance files as truth, never the binary's embedded
   copies, so old instances keep working under new binaries. Migration
   machinery is deferred (DEC-011).
7. v1's authority rules stand: git-tracked artifacts are authoritative; chat
   and model memory are not.

## Success criteria

1. Spark-to-first-captured-thought under five minutes, including repo
   creation.
2. Re-entry after a multi-week gap: one wake brief restores working context
   for human and agent without rereading raw logs.
3. Conformance suite green means: no orphaned IDs, no illegal state
   transitions, no dissent record missing a crux, no simmering idea missing a
   revisit trigger.
4. A handoff packet is sufficient for a fresh implementation agent with no
   access to this repo's chat history.
5. The framework's own maintenance cost stays below the thinking it saves;
   if meta-work dominates, subtract (Laziness guard).

## Lineage and sources

| Source | What Mycelium takes |
| ------ | ------------------- |
| ADRP v1 (this repo) | Governance spine, contracts, evidence model, rigor tiers, replication/reconciliation rules, ID system, fresh-session rule, human gates |
| pstack (poteto) | Model-role configuration, sparring candor ("no is an acceptable answer"), arena/interrogate mechanics as council prompt content, second-opinion doctrine |
| karpathy/llm-council | Three-stage council shape: independent first opinions, anonymized cross-review, synthesis — subordinated to v1's reconciliation rules (no majority vote, dissent retained) |
| karpathy llm-wiki gist | index.md + log.md navigation, LLM-as-bookkeeper, decisions as living objects with falsifiable assumptions (Vigil) |
| mattpocock/skills | Grilling's frontier-of-decisions interview with a recommendation per question; wayfinder's fog-of-war and destination-first charting; domain-modeling's glossary + ADR discipline; handoff's compaction ethos |
| Ruby on Rails | Convention over configuration as identity; derivable naming; generators; timestamped migrations; environments (as rigor-tier config); engines (as presence-registered packs); dummy-app testing (as the fixture instance); teaching errors |

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
- **wake** is the simmering → exploring transition. Its ritual: reread index
  and log tail, check evidence revalidation triggers and assumption records
  against what changed, then brief the human.
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
- **Disagreement record** — both positions, both reasons, and the **crux**:
  what evidence would change each mind. Cruxes are eligible to become
  research stages or spikes.
- **Decision** (`DEC-###`) — v1 template, plus assumptions listed and dissent
  retained inline.
- **Evidence** (`EVD-###`, `SPK-###`) — v1 evidence model unchanged,
  revalidation triggers now consumed by the wake ritual.
- **Council** — a commissioned replication run whose replicas are
  *different models*. Reuses v1's replication and reconciliation contracts
  and directories verbatim; reports land per-model, reconciliation retains
  dissent, and selection by majority vote or model reputation stays banned.
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
operational surface, and a shipped skill teaches agents both the commands
and the manual floor.

### Instance (per idea)

`mycelium new idea "name"` scaffolds it: README, manifest (with
`state = "spark"` and the generating CLI version), `log.md`, empty
`CONTEXT.md`, `AGENTS.md`, `program/` contracts and templates, then git
init, `gh repo create`, and the `idea` topic. Structure is emitted per tier:
a spark carries only what a spark needs, and deeper docs/ trees arrive when
a tier that requires them is set. (v1's full tree of placeholders existed to
compensate for indiscriminate template copying; the CLI ends that.)
Conformance requirements are **tier-aware**: a spark passes with almost
nothing; ledger discipline binds only at the tier that demands it.

## Conventions and conformance

Adopted 2026-08-14: the full Rails borrowing set (items 1–8 below plus the
declared-deviation rule), approved by Robert as a block.

1. Naming, IDs, and layout are derivable, not searchable. The ID-to-path
   rule is a pure function: every namespace has exactly one home directory
   and one filename pattern (`DEC-014` → `decisions/DEC-014-<kebab-slug>.md`).
   Directories are plural lowercase nouns, one artifact type each; IDs are
   `UPPER-###`; slugs are kebab-case. Conformance enforces the mapping both
   directions: files must match their pattern, references must resolve.
2. Rigor tiers are machine-readable config (`program/tiers/*.toml`), not
   prose. The manifest pins a tier; the checker reads tier config to learn
   which artifacts bind. (Rails environments, translated.)
3. Generators (`mycelium new <type> "Title"`) create the full bundle: file
   from template, next ID allocated, manifest entry when stage-tracked, log
   line appended. Half-created artifacts stop being possible. There is no
   destroy; state-transition commands (`mycelium supersede`) honor the
   no-deletion rule.
4. No migrations (DEC-011). Instead: every manifest records the generating
   CLI version, runtime commands validate against the instance's own emitted
   schemas (instance files are truth), and the master keeps a `CHANGELOG.md`
   as each release's human-readable face.
5. Packs are presence-registered: a directory under `program/packs/<name>/`
   containing its own templates, contracts, and checks. Drop in to enable,
   remove to disable, no registry file. Conformance fails namespace
   collisions between packs.
6. The framework tests itself against a fixture instance: CI scaffolds a
   throwaway instance with the freshly built binary, generates one of every
   artifact type, and runs the checks. Convention breakage surfaces in the
   master's CI, never first in a real idea repo. (Rails engines' dummy app,
   translated.)
7. Every template ships a sidecar schema; generator and checker are both
   data-driven off it (see Templates and generation).
8. Declared deviation: an instance may deviate from a convention only by
   declaring the deviation in its manifest, visibly and auditably. Silent
   deviation is a conformance failure even when the deviation itself would
   be acceptable.
9. The conformance suite validates **structure only**: ID uniqueness and
   sequence, ID-to-path integrity, link resolution, required metadata and
   sections per schema, legal state transitions, tier-appropriate artifact
   presence, parseable log prefixes, crux presence on disagreement records,
   and revisit triggers on simmering ideas. Checks never grade content;
   thinking quality is judged by
   adversarial review, councils, and the human (Goodhart guard).
10. Failures teach: every error names the violated convention, links its
    contract, and suggests the fixing command (did-you-mean style).
11. CI runs the suite on push. Green = the record is structurally
    trustworthy.

## Templates and generation

Every artifact type is defined by exactly two files, side by side in
`program/templates/` (or a pack's `templates/`): the markdown template and a
sidecar schema (`<type>.schema.toml`) declaring the ID namespace, home
directory, filename pattern, required metadata keys, and required sections.

The schema is the shared truth. The generator reads it to know what to
create and where; the checker reads the same file to know what to validate.
The two tools cannot disagree, and adding an artifact type means adding a
template and a schema, with zero code changes to either tool.

Mechanics:

- Templates and schemas live as browsable files under `program/` in the
  master and are `go:embed`'d into the binary. Scaffolding emits them into
  the instance, where they become that instance's frozen truth.
- Templates use `{{ID}}`, `{{TITLE}}`, `{{SLUG}}`, `{{DATE}}` tokens. Plain
  string replacement in Go's stdlib, no templating engine.
- `mycelium new <type> "Title"` allocates the next ID by scanning the type's
  home directory. The filesystem is the registry: because ID-to-path is a
  pure function, a directory scan is unambiguous and there is no central
  index file to rot. The generator refuses to overwrite, fills tokens,
  writes the file, appends the log entry, updates the manifest for
  stage-tracked types, prints the path plus next steps, and never runs git.
- Runtime commands (`new`, `check`, `status`) read the instance's own
  emitted templates and schemas, never the binary's embedded copies. A new
  binary therefore generates and validates old instances by their own
  rules; embedded data is used only at scaffold time.
- `mycelium supersede <ID> --by <ID>` flips status, wires cross-links into
  both records, and logs the transition.
- Stage-scoped namespaces (`REC`/`REQ`/`FND`) respect the ranges the
  blueprint and manifest allocate; enforcement strictness is OQ-007.
- Artifact metadata keeps the human-readable bullet style of the v1
  templates. Machine state stays in the manifest per the manifest-authority
  rule; schemas exist so the checker can parse the bullets, not to move
  state into artifact files.

## Perspective ladder

1. **Sparring** (free, always on) — the resident agent holds positions,
   challenges terms, surfaces presuppositions (the assumption audit), and
   recommends on every substantive question.
2. **Second opinion** (cheap, one word to invoke) — the identical
   commissioning prompt to exactly one different model. Agreement is
   high-signal; disagreement surfaces a fork.
3. **Council** (expensive, opt-in, suggested only when v1's replication
   triggers fire: hard to reverse, weak or conflicting evidence, still
   low-confidence after a spike) — full multi-model replication +
   reconciliation. Cost class stated before running. Panel presets live in
   user-level config (quick / standard / high-stakes).

Engine-agnosticism: contracts define commissioning prompts and report file
shapes; *how* reports get produced is a swappable adapter. Cursor's parallel
multi-model subagents are the first adapter. The manual floor — pasting the
commissioning prompt into N chat UIs and saving N files — satisfies the
contract with zero tooling. Runtimes that cannot fan out simply skip ladder
rungs 2–3; `AGENTS.md` carries the capability note.

## Cross-idea operations

No hub repo. Each instance manifest carries `state` and `revisit`; the
scaffolder applies the `idea` topic on GitHub automatically (DEC-003's topic
hygiene, now enforced by tooling instead of memory). `mycelium status --all`
scans the ideas root (`~/ideas/<slug>` by convention) and answers "what's
simmering, what's due to wake"; a thin skill wraps the command for agents.

## Phases and milestones

Sequential; each phase ends verifiable. Stops at milestones, per house rule.

- **PHASE-01 Foundation.** Go CLI skeleton; `program/` content authored and
  `go:embed`'d; naming and ID-to-path contracts; template sidecar schemas;
  tiers as machine-readable config; `mycelium new idea` scaffolding (emit,
  git init, `gh repo create`, `idea` topic); the data-driven
  `mycelium new <type>` generator; manifest gains idea-lifecycle fields and
  the CLI version stamp; `mycelium check` (schema-driven, tier-aware,
  teaching errors); fixture-instance CI (scaffold, generate, check).
  *MS-101: a machine with the binary goes from nothing to a conformant spark
  instance in under five minutes.*
- **PHASE-02 Lifecycle.** Spark / wake / portfolio skills; log + index
  conventions; re-entry brief; simmer with revisit triggers.
  *MS-201: an idea simmered for 7+ days wakes with a brief citing what
  changed.*
- **PHASE-03 Sparring.** Thinking-mode skill: mandatory positions, agreement
  states, disagreement records with cruxes, glossary challenge, assumption
  audit. Grilling and domain-modeling conventions absorbed and credited.
  *MS-301: a session transcript yields zero bare questions and at least one
  recorded crux.*
- **PHASE-04 Perspective ladder.** Commissioning + report + reconciliation
  contracts for model-diverse replication; second-opinion move; Cursor
  council adapter; panel presets in user config. Council ships as the first
  pack, proving presence-is-registration.
  *MS-401: one full council runs end to end; all artifacts pass conformance;
  dissent retained.*
- **PHASE-05 Distribution and lifecycle commands.** `mycelium supersede`;
  tagged releases with prebuilt binaries and install docs; `CHANGELOG.md`
  discipline; portfolio scanner tolerant of older manifest shapes (DEC-011's
  risk guard).
  *MS-501: a clean VM goes from one-line install to a scaffolded instance in
  under a minute.*
- **PHASE-06 Handoff.** Packet contract + generator; pstack/poteto bridge
  documented; implementation-systems section in AGENTS.md.
  *MS-601: a fresh agent in a clean session implements from a packet alone.*

## Identifier allocations (framework evolution)

- Framework decisions: `DEC-001`–`DEC-099` (in `framework/decisions/`).
- Blueprint adversarial-review findings, when commissioned: `FND-001`–`FND-099`.
- Open questions herein: `OQ-001`–`OQ-019`.

## Open questions

- **OQ-001** Final state-vocabulary bikeshed (`spark/exploring/simmering/...`)
  — settle in PHASE-01 contracts.
- **OQ-002** Manifest filename: keep `research-program.toml` or rename
  (`mycelium.toml`) as a 2.0 migration.
- **OQ-003** Resolved 2026-08-14: packs are presence-registered directories
  under `program/packs/` (see Templates and generation). The remaining
  sliver, which capabilities beyond council become packs, rolls into
  PHASE-04.
- **OQ-004** Handoff packet format details and the pstack playbook mapping.
- **OQ-005** Repo rename timing and instance-topic conventions.
- **OQ-006** Whether the blueprint itself gets a council review before
  PHASE-01 begins (recommended: yes, as the council contract's dogfood).
- **OQ-007** Generator strictness for stage-scoped ID ranges: warn or refuse
  when allocating outside a declared range.

## Authority, amendment, completion

This blueprint governs the evolution once accepted. Accepted `DEC-###`
records in `framework/decisions/` supersede it clause by clause. Amendments
follow `program/reference/amendment-protocol.md`. The evolution is complete
when all six phases are accepted and one real idea has traveled
spark → handed-off entirely inside Mycelium.

Fresh-session note: this blueprint is the synthesized output of the
2026-08-14 discovery interview. Per the fresh-session rule, its adversarial
review (OQ-006) and each build phase should run in fresh sessions with
self-contained packets.
