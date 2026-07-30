# {{PROJECT_NAME}}

Artifact-driven research program repository (GitHub template).

Evidence-grounded research → coherent architecture → adversarial review →
implementation plan as **phases and milestones**. Git is the system of record.
Stops before a granular coding backlog.

## Quick start

1. **Create** a new repo from this template (GitHub → Use this template), or clone.
2. **Init** (working title; no git):

   ```bash
   just init "my-project"
   # or: just init name="my-project"
   ```

3. **Discovery** — interview one question at a time until framing is solid
   (use the `research-program` skill or follow `program/reference/discovery-protocol.md`).
4. Fill and **accept** `docs/00-program-blueprint.md`, then `docs/01-research-charter.md`.
5. Run the adaptive research graph with just-in-time prompts, fresh sessions,
   validation, and human approval gates.

```bash
just status   # manifest overview + eligible stages
just check    # tree and acceptance sanity
```

## What to read first

| If you are…              | Read                                                  |
| ------------------------ | ----------------------------------------------------- |
| Starting a new program   | `program/operator/getting-started.md`                 |
| Resuming work            | `program/operator/resume-protocol.md` + `just status` |
| An agent in this repo    | `AGENTS.md` + `.agents/skills/`                       |
| Writing an artifact type | `program/contracts/` + `program/templates/`           |
| Deep rules               | `program/reference/`                                  |

## Repository layout

```text
.
├── README.md                 # this file
├── AGENTS.md                 # agent/human operating rules
├── Justfile                  # init / status / check (no git)
├── research-program.toml     # operational manifest
├── .agents/skills/           # portable agent skills
├── program/                  # methodology library (not project conclusions)
├── decisions/                # DEC-### records
└── docs/
    ├── 00-program-blueprint.md
    ├── 01-research-charter.md
    ├── prompts/              # JIT prompts (fixed spine skeletons; tracks added later)
    ├── reports/              # focused research (created JIT)
    ├── reconciliations/
    ├── evidence/             # SPK-###
    ├── specifications/
    ├── plans/
    ├── reviews/
    ├── handoffs/
    └── validations/
```

Focused research tracks are **not** pre-created. After Blueprint, add them from
templates under `program/templates/`.

## Current implementation authority

Until revised artifacts are accepted:

| Role                     | Path                                                         | Status      |
| ------------------------ | ------------------------------------------------------------ | ----------- |
| Implementation authority | `docs/specifications/02-definitive-specification-revised.md` | Placeholder |
| Delivery authority       | `docs/plans/02-implementation-plan-revised.md`               | Placeholder |

Update this section when authorities are accepted.

## Principles (summary)

- Artifacts carry context; chat does not.
- Fresh session per substantive stage.
- Evidence before confidence; optional spikes and replication.
- Independent validation; human approval gates.
- Synthesis decides; review attacks; revision integrates.
- Final plan stops at phases and milestones.

Full methodology: [`program/README.md`](program/README.md).

## Skills

| Skill             | Path                                |
| ----------------- | ----------------------------------- |
| research-program  | `.agents/skills/research-program/`  |
| research-stage    | `.agents/skills/research-stage/`    |
| research-validate | `.agents/skills/research-validate/` |

Compatible with agents that load project skills from `.agents/skills`.

## Methodology version

1.0 — redistributed from the Artifact-Driven Research Program master
specification into the hybrid `program/` layout.
