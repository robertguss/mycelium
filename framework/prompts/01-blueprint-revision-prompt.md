# Blueprint revision prompt — Mycelium framework blueprint

- **Stage:** blueprint-revision (framework evolution loop; "Session B")
- **Commissioned:** 2026-08-14 by Robert Guss
- **Commissioning basis:**
  [`01-blueprint-adversarial-review.md`](../reviews/01-blueprint-adversarial-review.md),
  all thirteen findings accepted per its Dispositions section
- **Output:** revised [`framework/blueprint.md`](../blueprint.md), edited in
  place (git history is the change record)
- **Allowed file scope:** `framework/blueprint.md` only
- **Session rule:** fresh session working from this packet alone; do not
  execute any downstream stage

## Role

You are the revision editor for the Mycelium framework blueprint. You apply
accepted review findings faithfully. You do not invent new decisions, do not
relitigate dispositioned choices, and do not expand scope beyond the
corrections commissioned here.

## Authority and precedence

1. Accepted decision records `DEC-001` through `DEC-011` (highest).
2. The review's Dispositions section: Robert's accepted findings and chosen
   resolutions.
3. Each finding's Required Correction and Proposed Specification Diff.
4. The current proposed blueprint text, which yields wherever it conflicts
   with the above.

Chat history and model memory are not authoritative. Everything needed is in
the attached packet.

## Task

Revise `framework/blueprint.md` so every accepted finding's Required
Correction is applied. Where findings embedded choices, these resolutions are
binding:

1. **FND-001 — front matter.** Rewrite the metadata clause in Templates and
   generation: artifact metadata in emitted templates is front matter defined
   and validated by each type's sidecar schema. Name the single parser
   contract that both the generator and `mycelium check` consume. Remove the
   bullet-style clause. Do not retrofit the master's own `framework/`
   artifacts; that is out of scope.
2. **FND-002.** Give the instance manifest separate lifecycle-state and
   rigor-tier fields. Define legal tier transitions and name an idempotent
   operation that adds newly required structure without rewriting existing
   work; assign that operation to a phase.
3. **FND-003.** Specify `mycelium status --all` as live GitHub enumeration
   (`idea` topic plus manifest reads) with a documented local fallback and
   explicit handling of remote-only, local-only, archived, unauthenticated,
   and temporarily unavailable repositories. Assign its first usable version
   to a phase.
4. **FND-004.** Replace the absolute atomicity claim with a concrete
   operation protocol: preflight validation, locking, temporary-file writes,
   atomic renames where supported, commit ordering, rollback and retry
   semantics, and `mycelium check` detecting an interrupted operation with a
   documented recovery path. State the supported filesystem floor.
5. **FND-005 through FND-010 — milestones.** Rewrite MS-101, MS-201, MS-301,
   MS-401, MS-501, and MS-601 per each finding's Required Correction, at
   blueprint-wording granularity. Record explicitly that the full acceptance
   matrices land in each phase's contract when that phase is commissioned
   (Robert's scoping rule), so the blueprint does not become a test plan.
   In particular: MS-101 separates the hermetic local scaffold from the
   authenticated GitHub integration and keeps five minutes as a user SLO;
   MS-201 gains an injectable clock and a deterministic fixture; MS-301
   tests one aligned and one disputed fixture instead of requiring a crux
   per session; MS-401 becomes one row of a perspective-ladder acceptance
   matrix covering the DEC-008 contract; MS-501 splits functional
   acceptance (supersession links, old-manifest tolerance, release
   integrity) from the install SLO; MS-601 names a canonical packet fixture
   with a bounded implementation target and executable acceptance tests.
6. **FND-011 — both versions.** The instance manifest records
   `methodology_version` and `generated_by_cli_version` as separate fields.
7. **FND-012.** Add `.agents/skills/` and any other emitted runtime-adapter
   paths to the instance anatomy. State which emitted files are portable
   authority and which are generated conveniences.
8. **FND-013 — adversarial review, not council.** Rewrite OQ-006 as
   dispositioned: the 2026-08-14 review is recorded as an independent,
   manually commissioned adversarial review; council-contract dogfood moves
   to PHASE-04. Add commissioning-prompt provenance (durable in-repo
   prompts, model provenance, named ladder rung) to PHASE-04's scope.

## Rules

- Keep `Status: Proposed — awaiting acceptance by Robert Guss`. Acceptance is
  Robert's act, after independent validation.
- Update the blueprint's Date line and add a one-line revision note citing
  the review and its Dispositions section.
- Never reuse or renumber IDs. Disposition open questions explicitly (OQ-006
  at minimum; touch others only where an accepted finding requires it).
- Keep citations portable (relative Markdown links).
- No corrections beyond the thirteen accepted findings. The laziness guard
  applies: prefer the smallest text that makes each correction binding.

## Deliverables

1. The revised `framework/blueprint.md`.
2. A brief execution summary outside the artifact: findings applied, sections
   touched, anything unmet and why, any remaining blocker.

---

## Appendix A — Fresh-session launch message

> You are executing the blueprint revision stage of the Mycelium framework
> evolution.
>
> The attached files correspond to these authoritative repository artifacts:
>
> 1. `framework/prompts/01-blueprint-revision-prompt.md` — the commissioning
>    prompt for this stage
> 2. `framework/blueprint.md` — the proposed blueprint to revise
> 3. `framework/reviews/01-blueprint-adversarial-review.md` — the accepted
>    findings and Robert's dispositions
> 4. `framework/decisions/DEC-001` through `DEC-011` — accepted decision
>    records, the highest authority
> 5. `AGENTS.md` — repository operating rules
>
> Read every attached artifact completely before beginning. Apply their
> authority and precedence rules exactly. Execute the complete task
> commissioned by `framework/prompts/01-blueprint-revision-prompt.md`.
>
> If this session has write access to the repository, edit
> `framework/blueprint.md` in place and touch nothing else. If it does not,
> produce the complete standalone contents intended for
> `framework/blueprint.md`.
>
> Do not ask clarifying questions unless a true blocker exists under the
> commissioning prompt. Do not begin a downstream stage.
>
> At the end provide: the revised artifact, a brief execution summary, any
> unmet requirement and why, and any remaining blocker.

## Appendix B — Post-stage validation task

Run in a separate fresh session (`research-validate` skill) after the
revision exists.

Read: `AGENTS.md`, this prompt, the review including its Dispositions,
`DEC-001` through `DEC-011`, and the revised blueprint.

Validate:

- Every finding `FND-001` through `FND-013` has its Required Correction
  applied, faithful to the disposition resolutions: front matter, both
  manifest version fields, and this run recorded as an independent
  adversarial review.
- Milestone rewrites stay at blueprint-wording granularity and record the
  acceptance-matrix deferral to phase contracts.
- No revised text contradicts `DEC-001` through `DEC-011`.
- Status remains Proposed; Date and revision note updated.
- IDs unique and none reused; OQ-006 dispositioned; all links resolve.
- The git diff touches `framework/blueprint.md` only.

Fix only mechanical defects. Report substantive defects and mark the revision
requires-revision rather than approving it.

## Appendix C — Recommended commit boundaries (humans run git)

```text
docs: add blueprint adversarial review with dispositions   # review + OQ-005 disposition in blueprint
docs: add blueprint revision prompt                        # framework/prompts/ + framework/handoffs/
docs: publish revised framework blueprint                  # after Session B passes validation
docs: accept framework blueprint                           # status flip on Robert's acceptance
```
