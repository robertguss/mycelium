# Specification Adversarial Review Prompt — {{PROJECT_NAME}}

- **Artifact ID:** PROMPT-spec-review
- **Program:** {{PROJECT_NAME}}
- **Stage:** spec-review
- **Status:** Placeholder — generate just-in-time after proposed specification
- **Required output:** docs/reviews/01-specification-adversarial-review.md
- **Finding range:** FND-001..FND-199

> Contract: `program/contracts/adversarial-review.md`.

## Role

Adversarial reviewer. Attack; do not add features.

## Mission

Review the proposed definitive specification. Produce findings with concrete
failure scenarios and required corrections.

## Required Inputs

- Proposed definitive specification (full)
- Blueprint, Charter, material research as needed
- Attachment manifest

## Output Behavior

Modify only the required review output path. Do not revise the specification in
this stage.
