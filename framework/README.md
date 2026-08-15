# framework/ — Mycelium's own evolution

This directory holds the artifacts of the framework's self-evolution: its
blueprint and its decision records. It is the one place in this repository
that is **about the framework itself** rather than methodology that ships
inside the CLI.

Since DEC-010, Mycelium is a scaffolding CLI, not a GitHub template.
Instances receive only what `mycelium new idea` emits, and this directory is
simply never emitted. (The earlier init-time strip rule existed to
compensate for GitHub's indiscriminate template copying and is obsolete.)

## Contents

| Path           | Purpose                                         |
| -------------- | ----------------------------------------------- |
| `blueprint.md` | Governing design for the Mycelium evolution     |
| `decisions/`   | `DEC-###` records for framework-level decisions |

Framework decisions use their own `DEC-###` sequence, independent of any
instance's decision numbering. They follow the standard template at
`program/templates/decision-record.md`.
