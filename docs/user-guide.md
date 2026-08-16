# Mycelium user guide

This is the operator guide for people who will run ideas. The
[README](../README.md) is the landing page. After you scaffold an idea, the
in-repo start page is `program/operator/getting-started.md`.

Mycelium is meant to be used **with an LLM**, not as a pile of terminal
commands. The CLI is the ledger. The agent is the sparring partner.

## Contents

1. [What mycelium is](#what-mycelium-is)
2. [Install](#install)
3. [Work with an agent](#work-with-an-agent)
4. [Anatomy of an idea](#anatomy-of-an-idea)
5. [Lifecycle](#lifecycle)
6. [Rigor tiers](#rigor-tiers)
7. [Pick a scenario](#pick-a-scenario)
8. [Commands](#commands)
9. [Artifacts](#artifacts)
10. [Skills](#skills)
11. [Checks and teaching errors](#checks-and-teaching-errors)
12. [Git](#git)
13. [This repo vs an idea](#this-repo-vs-an-idea)

## What mycelium is

One idea gets one git repository. Inside it you and an agent capture:

- what the words mean (`CONTEXT.md`)
- what is still open (`questions/`)
- what you assume (`assumptions/`)
- what you decided (`decisions/`)
- what would change each mind (cruxes on disputed questions)
- what you measured (`evidence/`, `spikes/`)
- where the idea sits (`mycelium.toml` `state` + `tier`)

Chat is transport. Git-tracked files are authority. The agent must not treat
a previous conversation as the record.

Mycelium stops before implementation. The last artifact is
`handoff/PACKET.md`, which a fresh implementation agent can build from
without this chat history.

## Install

See [`install.md`](install.md). Short path:

```bash
# linux-amd64 or darwin-arm64 release binary
curl -fsSL https://github.com/robertguss/mycelium/releases/latest/download/mycelium-$(uname -s | tr A-Z a-z)-$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/') -o ~/.local/bin/mycelium && chmod +x ~/.local/bin/mycelium

mycelium version
```

Other platforms: build from source with Go 1.26 and `CGO_ENABLED=0`.

Optional:

| Variable | Effect |
| --- | --- |
| `MYCELIUM_OFFLINE=1` | Never talk to GitHub. Same as `--offline` on `new idea`, `publish`, and `status --all`. |
| `MYCELIUM_IDEAS_ROOT` | Directory `status --all` scans (default `~/ideas`). |
| `MYCELIUM_NOW` | RFC3339 clock override (tests and hermetic wakes). |
| `MYCELIUM_CONFIG` | Directory override for `~/.config/mycelium/` (council panel presets). |

By convention, keep idea repos under `~/ideas/<slug>`. The CLI does not
require that layout.

## Work with an agent

This is the default loop.

1. **Install** the CLI once.
2. **Scaffold** an idea (`mycelium new idea "…" --offline` until you want GitHub).
3. **Open the idea folder** in Cursor, Claude Code, Grok, or any runtime that
   can read files and run shell commands.
4. **State your goal** in the first message. Do not dump the whole methodology.
   Examples:
   - “Capture the spark.”
   - “We need to decide whether X. Spar with me.”
   - “Park this until 2026-09-01.”
   - “Wake this idea.”
   - “This needs a real research program.”
   - “Hand this off for implementation.”
5. The agent should follow `.agents/skills/session/SKILL.md` (or the matching
   ritual skill) and write artifacts with `mycelium new …`.
6. **You** approve framing, accept decisions, and run `git commit`.

If the runtime does not load skills, point it at `AGENTS.md` and
`program/operator/getting-started.md`.

You can type every `mycelium` command yourself. That is the **manual floor**.
It is always legal. It is not how the tool is designed to be used day to day.

## Anatomy of an idea

`mycelium new idea "Garden lighting" --offline` creates a directory named
`garden-lighting` (slug of the title) unless you pass `--dir`.

Typical tree at birth (`focused` tier):

```text
garden-lighting/
├── mycelium.toml          # operational index: state, tier, dates
├── README.md
├── AGENTS.md              # agent rules for this idea
├── CONTEXT.md             # glossary
├── index.md               # generated map
├── log.md                 # append-only operations
├── .agents/skills/        # session, spark, thinking, …
├── program/               # methodology copy frozen into this idea
├── briefs/                # written on wake
└── …type homes as you create artifacts
```

`mycelium.toml` is how the CLI recognizes an instance. Commands walk up from
the current directory (or `--dir`) until they find it.

The CLI **git inits**. It never `git add`s or `git commit`s. `--publish`
(or a later `mycelium publish`) creates the GitHub repo and applies the
`idea` topic when `gh` is authenticated.

## Lifecycle

Exact machine:

```text
spark → exploring ⇄ simmering → clarified → handed-off
any (except archived) → archived
archived → (none)
```

| From | To | Command | Extra |
| --- | --- | --- | --- |
| `spark` | `exploring` | `mycelium state exploring` | |
| `exploring` | `simmering` | `mycelium state simmering --revisit VALUE` | Revisit required |
| `exploring` | `clarified` | `mycelium state clarified` | No packet yet |
| `simmering` | `exploring` | `mycelium wake` | Writes `briefs/LATEST.md` |
| `clarified` | `handed-off` | `mycelium handoff` | Writes `handoff/PACKET.md` |
| any except `archived` | `archived` | `mycelium state archived` | Deletes nothing |

`--revisit` is either a date (`2026-09-01`) or `event:<kebab>`
(`event:after-iphone-launch`). It is legal only when the target is
`simmering`.

`mycelium state exploring` from `simmering` also writes the wake brief.
Silent wake is refused.

`mycelium state handed-off` is legal only if `handoff/PACKET.md` already
exists and passes structure. Prefer `mycelium handoff`.

`handed-off` and `archived` are terminal except `handed-off → archived`.

## Rigor tiers

State and tier are independent. Birth is `state = spark`, `tier = focused`.

| Tier | When | What binds |
| --- | --- | --- |
| `focused` | A spark, a reversible choice, a familiar problem | Thin spine. Almost empty still passes `check`. |
| `standard` | A real software project | Decisions, assumptions, evidence, questions, risks. |
| `high-assurance` | Failure is expensive or hard to reverse | Plus spikes, findings, recommendations, requirements, phases, milestones. |

```bash
mycelium tier standard
```

Raising a tier emits directories the new tier needs. It never overwrites
existing work. Lowering a tier deletes nothing; it only relaxes what
`check` requires.

The **governance spine** (discovery → blueprint → charter → research graph
→ spec → plan → handoff) is the long program. It is not required to leave
`spark`. Use it when the idea earns it. Details live under `program/`
inside the idea.

## Pick a scenario

### 1. Capture a spark (15–30 minutes)

Goal: the idea exists in git and the first thought is written down.

```bash
mycelium new idea "garden lighting" --offline
```

Open the folder in your agent:

> This is a new mycelium idea. Capture the spark. Interview me one question
> at a time. Do not invent scope.

The agent should:

1. Read `index.md` and `mycelium.toml`.
2. Interview you. Recommend on each question. Leave decisions with you.
3. `mycelium state exploring`
4. `mycelium new question "…"` or `mycelium new decision "…"`
5. `mycelium check`

Done when `state` is `exploring`, one artifact exists, and `check` is green.
Commit if you want the thought to survive.

Do **not** start a blueprint or a research graph in this session.

### 2. Work a live idea (one sitting)

Goal: move one question, decision, or assumption forward.

Open the already-exploring idea. Say the question. The `thinking` skill is
always on:

- The agent takes a position on every substantive question.
- Agreement is `open`, `aligned`, or `agree-to-disagree`.
- Disagreement keeps both positions, both reasons, and a **crux** (what
  evidence would change each mind).
- Vague words get sharpened in `CONTEXT.md`.
- Presuppositions become `mycelium new assumption`.
- A crux may become a spike. The agent does not auto-promote it.

End the sitting with `mycelium check`. Park it (scenario 3) if you are not
done.

### 3. Park it (simmer)

Goal: stop on purpose, not because you forgot.

Simmering means you *could* decide now and are choosing not to. It is not
blocked and it is not fog.

```bash
mycelium state simmering --revisit 2026-09-01
# or
mycelium state simmering --revisit event:after-iphone-launch
```

Or ask the agent: “Park this until September 1.”

### 4. Come back after a gap (wake)

```bash
mycelium wake
```

Read `briefs/LATEST.md` first. Do not start from `log.md`. Then work
(scenario 2).

The brief cites what changed: the simmer line, due evidence, open
assumptions, and a suggested next step.

### 5. Survey the portfolio

```bash
mycelium status --all
# hermetic / no gh:
mycelium status --all --offline --root ~/ideas
```

Merges local `~/ideas` (or `MYCELIUM_IDEAS_ROOT` / `--root`) with GitHub
repos tagged `idea` when `gh` works. If GitHub is unavailable the listing
is marked `partial: local-only`. Treat that as incomplete.

`--archived` includes archived ideas. Default view hides them.

### 6. A longer research program (days to weeks)

Use this when the cost of being wrong is high, the problem is unfamiliar,
or you need a spec and a plan — not just a captured thought.

1. Raise the tier (`mycelium tier standard` or `high-assurance`).
2. Open the idea in a **fresh** agent session.
3. Follow `program/operator/getting-started.md` and
   `program/reference/governance-spine.md`:
   discovery interview → program blueprint → research charter → focused
   research tracks → optional replication → synthesis → spec review →
   implementation plan → handoff.
4. One substantive stage per session. Chat is not authority.
5. Human approval gates are listed in `program/operator/approval-gates.md`.

The CLI still only writes containers. The agent and you do the thinking.
`mycelium check` still only grades structure.

### 7. Second opinion or council

Off by default. Suggest only when a decision is hard to reverse, evidence
conflicts, or confidence is still low after a spike.

- **Second opinion:** one other model, same prompt, one report, no
  reconciliation. Cheap.
- **Council:** N models, independent reports, then a reconciliation that
  **retains dissent**. Never majority vote. State the cost class first.

There is no `mycelium council` command. The surface is
`mycelium new commissioning`, `mycelium new model-report`,
`mycelium new reconciliation`, plus the `second-opinion` and `council`
skills. A runtime that cannot fan out pastes the same prompt into N chat
UIs (the manual floor).

### 8. Hand off to implementation

When the destination is actually reached:

```bash
mycelium state clarified
mycelium handoff
```

`handoff/` is self-contained. The implementer receives **only** that tree.
Default implementation system is pstack/poteto-mode; `manual` is always
legal. Mapping: `program/reference/implementation-systems.md`.

If `handoff/PACKET.md` already exists, use `mycelium state handed-off`
instead of running `handoff` again.

### 9. Replace an artifact or kill the idea

Supersede is **artifact-level**. It does not change idea `state`.

```bash
mycelium supersede DEC-001 --by DEC-004
```

Eligible: `DEC`, `ASM`, `EVD`, `SPK`. Open questions are not superseded;
open a new question.

```bash
mycelium state archived
```

Deletes nothing. `archived` is terminal.

## Commands

Global: exit `0` on success, `1` on failure. Most commands accept
`--dir PATH`. `--help` works on the binary and on each verb.

| Command | Purpose |
| --- | --- |
| `mycelium version` | Print the stamped version. |
| `mycelium new idea <name>` | Scaffold an idea. Always `git init`. Never commits. Flags: `--dir`, `--offline`, `--publish`, `--tier focused\|standard\|high-assurance`. |
| `mycelium new <type> "<Title>"` | Allocate the next ID and write the artifact from the instance templates. |
| `mycelium check` | Structure conformance. `--abort-journal` rolls back an interrupted operation. |
| `mycelium tier <tier>` | Set the rigor tier. Emit new dirs only. Never delete. |
| `mycelium state <target>` | Lifecycle transition. `--revisit` only when target is `simmering`. |
| `mycelium wake` | Simmering → exploring and write the re-entry brief. |
| `mycelium status` | Six-line single-idea status. |
| `mycelium status --all` | Portfolio. `--root`, `--offline`, `--archived`. |
| `mycelium index` | Rebuild `index.md`. Does not append a log line. |
| `mycelium publish` | Create or update the GitHub repo and `idea` topic. |
| `mycelium supersede <OLD-ID> --by <NEW-ID>` | Bidirectional supersede + log line. |
| `mycelium handoff` | Write `handoff/PACKET.md`, then `state = handed-off`. Legal only from `clarified`. |

There is no `mycelium think`, `spar`, `session`, `council`, `ladder`, or
`replicate` verb. Those are skills and artifact types.

### `new idea` flags

| Flag | Meaning |
| --- | --- |
| `--offline` | No GitHub. Still `git init`. |
| `--publish` | Create the GitHub repo now (needs `gh`). Contradicts `--offline`. |
| `--tier` | Birth tier. Default `focused`. |
| `--dir PATH` | Exact destination instead of `./<slug>`. |

### `new <type>` registered types

| Type key | ID | Home |
| --- | --- | --- |
| `decision` | `DEC-###` | `decisions/` |
| `assumption` | `ASM-###` | `assumptions/` |
| `evidence` | `EVD-###` | `evidence/` |
| `spike` | `SPK-###` | `spikes/` |
| `question` | `OQ-###` | `questions/` |
| `risk` | `RSK-###` | `risks/` |
| `finding` | `FND-###` | `findings/` |
| `recommendation` | `REC-###` | `recommendations/` |
| `requirement` | `REQ-###` | `requirements/` |
| `phase` | `PHASE-##` | `phases/` |
| `milestone` | `MS-###` | `milestones/` |
| `commissioning` | `CMP-###` | `reviews/commissioning/` |
| `model-report` | `RPT-###` | `reviews/reports/` |
| `reconciliation` | `RCL-###` | `reviews/reconciliations/` |

`finding`, `recommendation`, and `requirement` are stage-scoped: the
manifest must declare an ID range before `new` will allocate one.

The filesystem is the registry. Next ID is `max(N)+1` in that home.
Interrupted `new` / `tier` / `state` / `handoff` operations leave a
journal; `check` tells you how to finish or `--abort-journal`.

## Artifacts

Every generated file has TOML front matter (`+++`) and required headings
from its sidecar schema in `program/templates/`. You may edit Markdown by
hand. Keep front matter valid, keep the headings, keep the filename on the
ID-to-path rule:

```text
<home>/<NS>-<zero-padded-digits>-<kebab-slug>.md
```

Example: `decisions/DEC-001-use-dusk-sensors.md`.

Then run `mycelium check`.

## Skills

New scaffolds emit these under `.agents/skills/`:

| Skill | When |
| --- | --- |
| `session` | Start of a sitting. Picks the ritual from your goal. |
| `spark` | Fresh idea. First thought. |
| `thinking` | Always-on sparring in any non-archived state. |
| `simmer` | Park with a revisit trigger. |
| `wake` | Re-enter a simmering idea. |
| `clarify` | Destination reached. |
| `handoff` | Write the implementation packet. |
| `portfolio` | `status --all`. Read-only. |
| `mycelium-cli` | Command reference. |
| `second-opinion` | One other model. Pack present. |
| `council` | Multi-model replication. Pack present. |

`tier`, `index`, `state`, and `wake` do not retrofit skills into older
instances. Copy them from a newly scaffolded idea, or re-scaffold.

## Checks and teaching errors

`mycelium check` grades **structure only**: IDs, paths, required headings,
legal state, revisit grammar, pack rules, packet shape. It never grades
whether the thinking is any good. That is your job, plus optional
adversarial review and councils.

Failures look like this:

```text
mycelium: illegal transition spark → clarified
convention: lifecycle
contract: program/contracts/lifecycle.md
fix: allowed next states: exploring, archived
```

Fix the named convention, then re-run the command it suggests.

## Git

- The CLI never commits idea work product.
- Agents commit only when you ask.
- Humans own `git add` / `commit` / `push`.
- Prefer one coherent artifact per commit. See
  `program/reference/commit-boundaries.md` inside an idea.

## This repo vs an idea

You are reading the **product** repository. It builds the binary.

An **idea** repository is what `mycelium new idea` creates. It has
`mycelium.toml`. Work there. Do not turn this product repo into an idea.

Live rules for an idea are the files under that idea’s `program/contracts/`.
A newer CLI will still generate and check an old idea by **that idea’s**
emitted copy, not by the binary’s embedded copy.

## Read next (inside an idea)

1. `program/operator/getting-started.md`
2. `.agents/skills/session/SKILL.md`
3. `program/contracts/lifecycle.md` when a transition is refused
4. `program/reference/governance-spine.md` only for the long program
