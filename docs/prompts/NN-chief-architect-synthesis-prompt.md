# Chief Architect Synthesis Prompt — {{PROJECT_NAME}}

- **Artifact ID:** PROMPT-synthesis
- **Program:** {{PROJECT_NAME}}
- **Stage:** synthesis
- **Status:** Placeholder — generate just-in-time after research is accepted
- **Required output:** docs/specifications/01-definitive-specification.md
- **Requirement range:** REQ-001..REQ-299

> Do not freeze this prompt until prerequisite research stages are accepted.
> Replace `NN` with the actual sequence number when installed.
> Contract: `program/contracts/synthesis.md`.

## Role

Act as Chief Architect. Synthesis is decision-making, not summarization.

## Mission

Produce one coherent proposed definitive specification at the required output
path. Disposition every substantive `REC-###`.

## Authority and Precedence

See `program/contracts/authority-and-precedence.md` and the Program Blueprint.

## Required Inputs

- Accepted Blueprint and Charter
- All materially relevant accepted research reports (full text)
- Accepted Decision Records
- Attachment manifest for this stage

## Locked Context

_Fill just-in-time from Blueprint._

## Stage Boundary

### Included

Synthesis into one specification; requirement allocation; high-level phase
boundaries.

### Excluded

Focused research; adversarial review; granular implementation tasks.

## Methodology

- Read all material inputs completely.
- Disposition every REC.
- Resolve conflicts; reject weak machinery; simplify.
- Leave no foundational decision to implementers without an explicit bounded
  spike.

## Required Output Structure

See `docs/specifications/01-definitive-specification.md` skeleton and
`program/contracts/definitive-specification.md`.

## Completion Checklist

- [ ] Every REC dispositioned
- [ ] Standalone specification produced
- [ ] Status: Proposed — pending adversarial review
- [ ] Allowed file scope only
