# Repository Bootstrap

An idea instance is created with `mycelium new idea`. Typical layout:

```text
<idea>/
├── README.md
├── AGENTS.md
├── CONTEXT.md
├── mycelium.toml
├── index.md
├── log.md
├── .agents/skills/          # session, spark, thinking, …
├── program/                 # methodology (this library)
├── briefs/                  # written on wake
├── handoff/                 # when handed off
└── <type homes as needed>/  # decisions/, questions/, …
```

Focused research track files are **not** pre-created. After Blueprint
acceptance (when used), create them just-in-time from
[`../templates/`](../templates/).

Do **not** expect a Justfile, `scripts/`, or `research-program.toml` in an
idea repo. Scaffold will not emit `framework/` either.

## Stable filenames

Do **not** use: `final.md`, `final-v2.md`, `really-final.md`, `new-plan.md`,
`updated-spec.md`. Use stable, numbered, role-based names. Git history
records revisions.

## Placeholder rule

Bootstrap may create placeholders to reserve paths. A placeholder **does
not** prove stage completion. Only a validated artifact whose metadata
status is **accepted** may unlock downstream work.

## Bootstrap task (for repository agents)

See [`../templates/bootstrap-task.md`](../templates/bootstrap-task.md).

- Do not conduct substantive research.
- Do not invent project decisions beyond approved discovery output.
- Do not overwrite substantive content.
- Use stable filenames.
- Do not run git unless the human explicitly asks (the CLI never
  git-commits instance work product).
- Validate the complete tree with `mycelium check`.

## Manifest

`mycelium.toml` is the operational index: resume, legal transitions, state,
tier, identifier ranges. It must **not** contain substantive conclusions
absent from governing Markdown artifacts.

## Required root files

### README.md

Idea purpose, how to resume, what to read first. Open this folder in an
agent; do not treat the README as a command cheat-sheet.

### AGENTS.md

Skills, CLI surface, manual floor, teaching errors (see
`program/skeleton/AGENTS.md`).

### mycelium.toml

Canonical operational manifest (see
[`../contracts/manifest.md`](../contracts/manifest.md)).
