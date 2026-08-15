# Manual council adapter

Markdown procedure only. Zero tooling beyond chat UIs and the mycelium file
surface. The CLI does not invoke this adapter. Credits: DEC-008; v1
`program/contracts/replication-reconciliation.md`; karpathy/llm-council
(three-stage shape, subordinated to v1 reconciliation). Manual floor always
works.

## Preconditions

- A commissioning file exists with identical `## Prompt` for every model and a
  filled `## Attachments` manifest.
- `adapter = "manual"` on the CMP (and matching on each RPT).

## Procedure

1. Paste the identical prompt into N chat UIs (one model / UI per paste).
2. Save N RPT files via `mycelium new model-report` (or hand-edit), one per
   model, with matching `commissioning`, `rung`, `adapter`, `model`, and
   `prompt_sha256`.
3. For a council rung, write the RCL by hand via `mycelium new reconciliation`
   (or hand-edit), filling every v1 H2 including `## Retained dissent`.
4. For second-opinion, stop after exactly one RPT. No RCL.

No chairman-smoothing. Dissent retained. Do not choose by majority vote,
length, confidence of prose, or model reputation.

Optional karpathy/llm-council shape without tooling: (1) independent first
opinions in separate chats, (2) optional anonymized cross-review notes pasted
into RPT `## Findings`, (3) synthesis = the RCL.
