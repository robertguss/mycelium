# Repository Bootstrap

## Standard layout

```text
<project>/
├── README.md
├── AGENTS.md
├── Justfile
├── research-program.toml
├── .agents/skills/          # portable agent skills
├── program/                 # methodology (this library)
├── decisions/
│   ├── README.md
│   └── DEC-###-short-title.md
└── docs/
    ├── 00-program-blueprint.md
    ├── 01-research-charter.md
    ├── prompts/
    ├── reports/
    ├── reconciliations/
    ├── evidence/
    ├── specifications/
    ├── plans/
    ├── reviews/
    ├── handoffs/
    └── validations/
```

Focused research track files (`01-<focus>-…`) are **not** pre-created. After
Blueprint acceptance, create them just-in-time from
[`../templates/`](../templates/).

## Stable filenames

Do **not** use: `final.md`, `final-v2.md`, `really-final.md`, `new-plan.md`,
`updated-spec.md`. Use stable, numbered, role-based names. Git history records
revisions.

## Placeholder rule

Bootstrap may create placeholders to reserve paths. A placeholder **does not**
prove stage completion. Only a validated, committed artifact whose metadata
status is **accepted** may unlock downstream work.

## Bootstrap task (for repository agents)

See [`../templates/bootstrap-task.md`](../templates/bootstrap-task.md).

Rules:

- Do not conduct substantive research.
- Do not invent project decisions beyond approved discovery output.
- Do not overwrite substantive content.
- Use stable filenames.
- Do not run git unless the human explicitly asks (this template’s `just`
  recipes never run git).
- Validate the complete tree.

## Manifest

`research-program.toml` is the operational index: resume, legal transitions,
artifact paths, identifier ranges, accepting commit hashes. It must **not**
contain substantive conclusions absent from governing Markdown artifacts.

## Required root files

### README.md

Project purpose, program overview, layout, how to resume, current accepted
implementation authority, what to read first.

### AGENTS.md

Artifact authority, allowed file scope, validation, commits, citations,
identifiers, no silent edits to governing artifacts, fresh-session rules.

### research-program.toml

Canonical operational manifest (see root file and
[`../reference/state-machine.md`](../reference/state-machine.md)).
