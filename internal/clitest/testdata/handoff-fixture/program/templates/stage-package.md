# Standard Just-in-Time Stage Package

For every substantive stage, the Research Program Architect must produce five
items:

1. **Canonical stage prompt** (installed at reserved path under `docs/prompts/`).
2. **Repository installation task** for placing the prompt at its reserved path.
3. **Attachment manifest** (`docs/handoffs/<stage-id>-attachment-manifest.md`).
4. **Fresh-session launch message**.
5. **Post-stage validation task** and recommended commit message.

This package makes each transition reproducible and prevents missing context.

## Installation task shape

```markdown
# Install prompt for [STAGE ID]

- Source: (paste or path)
- Destination: docs/prompts/[filename].md
- Update research-program.toml: status → prompt-ready when installed
- Do not execute the research stage in the install session
```

## Validation task

See `validation-task.md`.
