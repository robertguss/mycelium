# Cursor council adapter

Markdown procedure only. The CLI does not invoke this adapter. Credits:
DEC-008; v1 `program/contracts/replication-reconciliation.md`;
karpathy/llm-council (three-stage shape, subordinated to v1 reconciliation).

## Preconditions

- A commissioning file exists with identical `## Prompt` for every model and a
  filled `## Attachments` manifest.
- `adapter = "cursor"` on the CMP (and matching on each RPT).
- Runtime can fan out parallel subagents. If it cannot, skip rungs 2–3;
  sparring still applies.

## Procedure

Fan out parallel subagents. Give every subagent the identical `## Prompt` body
and the attachment manifest. Each subagent saves one RPT via
`mycelium new model-report` (or hand-edit to the same shape).

Recommended karpathy/llm-council three-stage shape (procedure, not a fourth
artifact type):

1. Independent first opinions — each model writes its own RPT (`## Position`,
   `## Findings`, `## Dissent`) without seeing the others.
2. Optional anonymized cross-review — notes land in RPT `## Findings` (or a
   second pass on the same RPT). Do **not** add a fourth artifact type.
3. Synthesis = the RCL via `mycelium new reconciliation`, subordinated to v1
   reconciliation H2s including `## Retained dissent`.

No chairman-smoothing. Dissent retained. Do not choose by majority vote,
length, confidence of prose, or model reputation.

## Manual floor

If Cursor fan-out is unavailable, use `adapters/manual.md` or skip this rung.
