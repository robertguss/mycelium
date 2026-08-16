# Program State Machine

This is the **research-stage** status machine (planned → accepted), not
the idea lifecycle. Idea states live in `program/contracts/lifecycle.md`.

## Canonical stage statuses

`planned` | `prompt-ready` | `in-progress` | `awaiting-validation` |
`requires-revision` | `accepted` | `blocked` | `requires-revalidation` |
`superseded` | `cancelled`

## Legal transitions

```text
planned → prompt-ready → in-progress → awaiting-validation
  → accepted | requires-revision | blocked
requires-revision → in-progress
accepted → requires-revalidation | superseded
```

## Manifest rules

- Use stable stage IDs.
- Declare kind, dependencies, prompt, output, identifier ranges, status.
- Record accepting Git commit when accepted.
- Declare parallel groups where relevant.
- Preserve superseded stages; do not delete them.
- Update only through validated repository operation.
- Never mark accepted merely because an output path exists.

## Manifest authority

The Program Blueprint is the human-readable governing authority for program
design. The TOML manifest is the operational index for resume, transitions,
paths, prerequisites, identifier ranges, accepting commits, and next eligible
stage selection. No substantive conclusions that are absent from governing
Markdown.
