# Blueprint adversarial review

- **Status:** Complete — all 13 findings dispositioned and accepted 2026-08-14
  (see [Dispositions](#dispositions))
- **Date:** 2026-08-14
- **Reviewer:** GPT-5.6 Sol in Cursor
- **Drafting model family:** Claude/Fable
- **Reviewed artifact:** [`../blueprint.md`](../blueprint.md)
- **Decision authority reviewed:**
  [`DEC-001`](../decisions/DEC-001-evolve-adrp-in-place-into-mycelium.md)
  through [`DEC-011`](../decisions/DEC-011-defer-migrations.md)
- **Finding allocation:** `FND-001` through `FND-099`
- **Recommendation:** Revise before acceptance

## Review posture

This review treats accepted decision records as higher authority than the
proposed blueprint. Findings are proposed corrections. They do not become
requirements unless Robert accepts them.

The review traced the scaffold, generator, checker, portfolio, wake, council,
release, and handoff paths. It also checked failure recovery, external
dependencies, phase boundaries, and whether each milestone proves its phase.

## Overall assessment

The blueprint is coherent at the product level, but it is not ready to accept.
Three defects reach into PHASE-01:

1. The metadata format contradicts an accepted decision.
2. Lifecycle state and rigor tier are mixed together, with no defined way to add
   higher-tier structure to an existing instance.
3. The generator promises atomic behavior but specifies several independent
   writes with no rollback.

The portfolio command also departs from DEC-003's live-query decision. Each
phase has a milestone, but most milestones prove only one happy-path outcome and
leave the phase's load-bearing behavior untested.

## Perspective-ladder observation

This run used a model family different from the drafting model and worked from a
self-contained file packet. That was useful: the accepted metadata clause, tier
transition gap, and portfolio contradiction surfaced without relying on the
drafting conversation.

Under DEC-008's exact definitions, however, this was neither a council nor a
strict second opinion:

- A council requires multiple independent model reports and reconciliation.
- A second opinion receives the same commissioning prompt as the first model.
  This reviewer received an adversarial-review commission, not the discovery and
  drafting prompt.

The run is best recorded as an independent, manually commissioned adversarial
review by a different model. The manual floor itself was workable. Its weak
point was provenance: the packet did not contain a durable copy of the original
commissioning prompt, so prompt identity could not be checked. That gap should
inform PHASE-04.

## Findings

## FND-001 - Blueprint metadata format contradicts DEC-005

- **Severity:** High
- **Confidence:** High
- **Category:** Authority conflict
- **Affected sections:** Templates and generation
- **Affected requirements:** None
- **Affected phases:** PHASE-01
- **Blocks implementation:** PHASE-01

### Problem

The blueprint requires human-readable metadata bullets. DEC-005 requires
front-matter schemas. DEC-011 supersedes only DEC-005's migration clause, so the
front-matter clause still governs.

### Evidence

[`DEC-005`, Decision item 3](../decisions/DEC-005-convention-over-configuration.md)
requires conformance to validate "front-matter schemas." The blueprint's
[Templates and generation](../blueprint.md#templates-and-generation) section
says artifact metadata keeps the v1 bullet style and that the checker parses
those bullets.

### Failure Scenario

PHASE-01 implements bullet parsing. A later validator applies the accepted
front-matter rule and rejects every generated artifact, or two formats become
supported and drift.

### Impact

The generator, schemas, templates, and checker lack one authoritative metadata
grammar.

### Root Cause

The blueprint carries forward the current template appearance without
dispositioning the accepted front-matter clause.

### Required Correction

Choose one metadata grammar. If bullets remain, accept a decision record that
explicitly supersedes DEC-005's front-matter clause. Otherwise change the
blueprint and templates to front matter.

### Proposed Specification Diff

Replace the bullet-style clause with the chosen grammar and name the parser
contract that both generator and checker use.

### Acceptance Evidence

A fixture generated from every sidecar schema parses through the same metadata
reader used by `mycelium check`.

### Alternatives Considered

Supporting both formats was rejected because it recreates the drift the sidecar
design is meant to remove.

### Residual Risk

Markdown bodies can still contain heading-like text. The parser contract must
bound metadata and section detection.

### Related Findings

FND-004

## FND-002 - Lifecycle state and rigor tier lack an orthogonal transition model

- **Severity:** High
- **Confidence:** High
- **Category:** Domain model
- **Affected sections:** Vision; Instance; Conventions and conformance
- **Affected requirements:** None
- **Affected phases:** PHASE-01, PHASE-02
- **Blocks implementation:** PHASE-01

### Problem

The blueprint says a "spark" receives thin structure and deeper trees arrive
when a tier is set. `spark` is a lifecycle state, not a rigor tier. The manifest
fields, legal tier transitions, and operation that emits newly required files
are not defined.

### Evidence

[`DEC-006`, Decision](../decisions/DEC-006-idea-lifecycle-with-simmer.md)
defines `spark` as a lifecycle state.
[`DEC-002`, Decision](../decisions/DEC-002-durable-tiered-record-as-the-product.md)
says rigor artifacts bind at the tier that demands them. The blueprint's
[Instance](../blueprint.md#instance-per-idea) section says "a spark carries only
what a spark needs" and deeper trees arrive when a tier is set.

### Failure Scenario

An idea moves from `spark` to `exploring` while remaining at a low rigor tier,
or raises rigor while still exploring. The CLI cannot tell whether to change
state, tier, or both, and has no specified command for emitting the new
structure.

### Impact

The promise that ideas grow rigor in place is not implementable from the current
contract. Manual copying becomes the hidden upgrade mechanism.

### Root Cause

Lifecycle progression and research rigor were both expressed as "thin to deep"
and collapsed in the scaffold description.

### Required Correction

Define separate manifest fields for lifecycle state and rigor tier. Define legal
tier transitions and an idempotent operation that adds required structure
without rewriting existing work.

### Proposed Specification Diff

Replace "a spark carries only what a spark needs" with a matrix of lifecycle
state versus rigor tier. Assign tier promotion to a named phase and command.

### Acceptance Evidence

Tests independently change state and tier in both orders, preserve existing
files, and leave `mycelium check` green.

### Alternatives Considered

Binding each lifecycle state to one tier was rejected because it conflicts with
DEC-002's purpose: the same kind of idea can warrant different rigor.

### Residual Risk

Because DEC-011 defers migrations, newly introduced future tiers still need a
manual adoption rule.

### Related Findings

FND-005

## FND-003 - Local portfolio scanning conflicts with DEC-003's live query

- **Severity:** High
- **Confidence:** High
- **Category:** Authority conflict
- **Affected sections:** Cross-idea operations; Phases and milestones
- **Affected requirements:** None
- **Affected phases:** PHASE-02, PHASE-05
- **Blocks implementation:** PHASE-02

### Problem

The blueprint limits `mycelium status --all` to scanning `~/ideas/<slug>`.
DEC-003 says the portfolio queries GitHub live using the `idea` topic and
manifest reads. Local-only scanning drops ideas that are not currently cloned.

### Evidence

[`DEC-003`, Decision](../decisions/DEC-003-one-repository-per-idea.md) places
cross-idea state in GitHub metadata and requires a live query.
[`DEC-010`, Decision item 3](../decisions/DEC-010-mycelium-is-a-cli.md) names
`mycelium status --all` as a core command. The blueprint's
[Cross-idea operations](../blueprint.md#cross-idea-operations) section says the
command scans the ideas root.

### Failure Scenario

An idea exists on GitHub with the `idea` topic but has been removed locally. The
portfolio reports no due wake and Robert misses its revisit trigger.

### Impact

The portfolio cannot be trusted as the complete cross-idea view. This reaches
the main mitigation for repository sprawl in DEC-003.

### Root Cause

The predictable local root was treated as the registry even though DEC-003 puts
registry truth in GitHub metadata.

### Required Correction

Specify live GitHub enumeration and manifest retrieval, with a documented local
fallback for offline use. Assign the command's first usable version to a phase.

### Proposed Specification Diff

Replace "scans the ideas root" with the authoritative live workflow and define
how local-only, remote-only, unavailable, and archived repositories appear.

### Acceptance Evidence

Portfolio tests cover remote-only, local-only, archived, unauthenticated, and
temporarily unavailable repositories.

### Alternatives Considered

A local-only registry was rejected by DEC-003. A generated hub remains a revisit
option only if live querying proves untrustworthy.

### Residual Risk

GitHub availability and rate limits can still make the complete view temporarily
unavailable. The command must report that state rather than silently return a
partial list.

### Related Findings

FND-009

## FND-004 - Generator atomicity is asserted but not specified

- **Severity:** High
- **Confidence:** High
- **Category:** Filesystem behavior
- **Affected sections:** Conventions and conformance; Templates and generation
- **Affected requirements:** None
- **Affected phases:** PHASE-01
- **Blocks implementation:** PHASE-01

### Problem

The blueprint says half-created artifacts become impossible, then specifies a
sequence that writes an artifact, appends the log, and updates the manifest.
Those are separate filesystem operations. No transaction, ordering guarantee,
rollback, lock, or repair command is defined.

### Evidence

[Conventions and conformance item 3](../blueprint.md#conventions-and-conformance)
claims generators make the full bundle and that half-created artifacts stop
being possible.
[Templates and generation](../blueprint.md#templates-and-generation) lists the
independent writes.

### Failure Scenario

The artifact write succeeds and the manifest update fails because the manifest
is malformed or the disk fills. The new ID exists, the log and manifest
disagree, and retry allocates another ID or refuses to overwrite. Two concurrent
agents can also scan the same next ID.

### Impact

The generator can create the exact bookkeeping failures it claims to remove.
Recovery behavior becomes implementation guesswork.

### Root Cause

Data-driven generation was specified, but the multi-file commit protocol was
not.

### Required Correction

Define preflight validation, a repository lock, temporary-file writes, atomic
renames where supported, commit ordering, rollback, and retry semantics. Define
how `check` detects an interrupted operation and points to the documented
recovery path.

### Proposed Specification Diff

Replace the absolute atomicity claim with a concrete operation protocol and
document platform limits.

### Acceptance Evidence

Fault-injection tests fail each write step and prove that the documented
recovery path either completes the original ID or returns the repository to its
prior state.

### Alternatives Considered

Writing only the artifact and making log and manifest updates manual was
rejected because it abandons the accepted full-bundle generator.

### Residual Risk

Atomic rename semantics differ across filesystems. The contract should state the
supported floor.

### Related Findings

FND-001

## FND-005 - MS-101 neither matches its prerequisites nor proves PHASE-01

- **Severity:** High
- **Confidence:** High
- **Category:** Milestone testability
- **Affected sections:** PHASE-01; Success criteria
- **Affected requirements:** None
- **Affected phases:** PHASE-01
- **Blocks implementation:** PHASE-01

### Problem

MS-101 says a machine "with the binary" starts from nothing, but scaffolding
also requires `git`, `gh`, GitHub authentication, network access, an account
destination, visibility policy, and permission to create a repository. The
five-minute happy path does not prove most PHASE-01 deliverables.

### Evidence

[`DEC-010`, Decision item 3](../decisions/DEC-010-mycelium-is-a-cli.md) and the
blueprint's [PHASE-01](../blueprint.md#phases-and-milestones) both put
`git init`, `gh repo create`, and topic mutation in the scaffold path. MS-101
mentions only the binary and a conformant spark.

### Failure Scenario

The scaffold passes under five minutes on an already authenticated laptop while
artifact generation, tier-aware checks, teaching errors, lifecycle fields, or
fixture CI are broken. In CI, remote repository creation either cannot run or
leaks test repositories.

### Impact

The phase can be accepted without its core generator and checker working. Timing
results cannot be compared across environments.

### Root Cause

A product speed goal was used as the whole phase acceptance gate.

### Required Correction

Separate a hermetic local scaffold test from an authenticated GitHub integration
test. Define prerequisites, repository visibility, naming, collision handling,
cancellation, partial-failure cleanup, and timing boundaries. Add acceptance
cases for every PHASE-01 deliverable.

### Proposed Specification Diff

Keep the five-minute target as a user SLO. Add an acceptance matrix for local
scaffold, remote publication, generation, checking, tier behavior, and CI.

### Acceptance Evidence

Fixture CI runs without network side effects. A separate disposable-account test
creates and cleans up a remote repository and topic. Both publish timing and
prerequisite data.

### Alternatives Considered

Dropping GitHub creation from the product was rejected because accepted DEC-010
names it as core behavior.

### Residual Risk

Remote latency remains variable, so the remote timing target needs a stated
environment and retry policy.

### Related Findings

FND-002, FND-004

## FND-006 - MS-201 has no deterministic clock or correctness oracle

- **Severity:** Medium
- **Confidence:** High
- **Category:** Milestone testability
- **Affected sections:** PHASE-02
- **Affected requirements:** None
- **Affected phases:** PHASE-02
- **Blocks implementation:** No

### Problem

MS-201 requires an idea to simmer for seven days and produce a brief citing
"what changed." It does not define the clock, the controlled changes, required
citations, or what makes the brief correct.

### Evidence

[`DEC-006`, Decision](../decisions/DEC-006-idea-lifecycle-with-simmer.md)
requires wake to check evidence triggers and assumptions against changes. The
blueprint's [MS-201](../blueprint.md#phases-and-milestones) states only elapsed
time and a brief.

### Failure Scenario

Validation waits a real week and then accepts a fluent brief that misses an
expired evidence source, or tests immediately by changing timestamps in an
unsupported way.

### Impact

PHASE-02 cannot be validated quickly or reproducibly.

### Root Cause

The milestone describes a demonstration, not an executable test contract.

### Required Correction

Define an injectable clock and a fixture with known log entries, evidence
triggers, assumption changes, and expected citations.

### Proposed Specification Diff

Replace the seven-day wait with a simulated elapsed-time scenario plus one human
dogfood run after seven real days.

### Acceptance Evidence

The deterministic fixture produces the expected transition, overdue status,
revalidation results, and source links.

### Alternatives Considered

Using only a real-time dogfood run was rejected because it delays every
revalidation and cannot isolate regressions.

### Residual Risk

The usefulness of the prose brief still needs human judgment.

### Related Findings

None

## FND-007 - MS-301 rewards manufactured disagreement and lacks an evaluation protocol

- **Severity:** High
- **Confidence:** High
- **Category:** Milestone validity
- **Affected sections:** PHASE-03
- **Affected requirements:** None
- **Affected phases:** PHASE-03
- **Blocks implementation:** PHASE-03

### Problem

MS-301 requires at least one crux from a session. An aligned session may have no
honest disagreement. Requiring a crux rewards the performative disagreement that
DEC-007 names as a risk. The milestone also does not say who decides whether a
question is substantive or what evidence makes that judgment repeatable.

### Evidence

[`DEC-007`, Risks](../decisions/DEC-007-sparring-stance-agreement-states-cruxes.md)
warns against manufactured dissent.
[`DEC-005`, Decision item 5](../decisions/DEC-005-convention-over-configuration.md)
limits automated conformance to containers while leaving thinking quality to
adversarial review, councils, and the human. The blueprint's
[MS-301](../blueprint.md#phases-and-milestones) requires one recorded crux and
zero bare questions without assigning the semantic judgment.

### Failure Scenario

An agent invents a weak disagreement so the fixture passes, while an honestly
aligned session fails. Two reviewers can also reach opposite milestone decisions
because "substantive" has no review procedure.

### Impact

The milestone distorts the behavior the phase is meant to cultivate.

### Root Cause

The artifact shape for a genuine disagreement was turned into a quota for every
session.

### Required Correction

Use one deliberately disputed fixture to test crux fields and one aligned
fixture to prove no disagreement is required. Assign substantive-question
judgment to a human or adversarial reviewer and define the review evidence
without turning it into an automated content score.

### Proposed Specification Diff

Make the milestone about correct handling of aligned and disputed sessions, not
a mandatory disagreement count.

### Acceptance Evidence

Structural checks pass both fixtures, and the disputed fixture retains both
positions, reasons, and cruxes.

### Alternatives Considered

An automated natural-language classifier was rejected because DEC-005 assigns
thinking-quality judgment to reviewers and the human.

### Residual Risk

Agents can still fill required fields with boilerplate. DEC-007 correctly leaves
that judgment to humans.

### Related Findings

None

## FND-008 - MS-401 does not prove the accepted perspective-ladder contract

- **Severity:** High
- **Confidence:** High
- **Category:** Milestone coverage
- **Affected sections:** Perspective ladder; PHASE-04
- **Affected requirements:** None
- **Affected phases:** PHASE-04
- **Blocks implementation:** PHASE-04

### Problem

One successful council does not prove opt-in behavior, cost disclosure,
identical commissioning prompts, model diversity, report retention,
reconciliation rules, pack registration, or the manual floor. "Dissent retained"
is vacuous if the selected run contains no dissent.

### Evidence

[`DEC-008`, Decision](../decisions/DEC-008-perspective-ladder-opt-in-councils.md)
requires each behavior. The blueprint's
[MS-401](../blueprint.md#phases-and-milestones) requires only one end-to-end
council whose artifacts pass conformance.

### Failure Scenario

A Cursor adapter auto-runs two instances of the same model, omits cost
preflight, reconciles away disagreement, and still produces structurally valid
files.

### Impact

The phase can pass while violating the authority that created it.

### Root Cause

An integration demo stands in for a contract acceptance suite.

### Required Correction

Add separate second-opinion, Cursor-council, and manual-floor scenarios. Require
model provenance, prompt hashes, explicit opt-in, cost class, independent
reports, seeded dissent, reconciliation, and pack enable/disable tests.

### Proposed Specification Diff

Keep the end-to-end council demonstration, but make it one row in a
perspective-ladder acceptance matrix.

### Acceptance Evidence

Stored manifests prove prompt identity and model diversity. A seeded
disagreement survives reconciliation. Removing the council pack disables it
without affecting core checks. The manual path passes without Cursor fan-out.

### Alternatives Considered

Treating the current review as council dogfood was rejected because it had one
review report and no reconciliation.

### Residual Risk

Model providers may obscure exact model identity. The provenance contract must
define acceptable evidence.

### Related Findings

FND-013

## FND-009 - MS-501 tests installation speed but not PHASE-05 behavior

- **Severity:** High
- **Confidence:** High
- **Category:** Milestone coverage
- **Affected sections:** PHASE-05
- **Affected requirements:** None
- **Affected phases:** PHASE-05
- **Blocks implementation:** PHASE-05

### Problem

MS-501 measures one-line install to scaffold time. It does not test supersession
integrity, release artifacts, changelog discipline, or portfolio compatibility
with older manifest shapes.

### Evidence

[`DEC-011`, Decision and Risks](../decisions/DEC-011-defer-migrations.md)
requires old instances to keep working and names tolerant portfolio scanning as
the risk guard. The blueprint's
[PHASE-05](../blueprint.md#phases-and-milestones) includes those behaviors,
while MS-501 only measures speed.

### Failure Scenario

The installer and scaffold work in under a minute, but `mycelium supersede`
leaves one-way links and `status --all` crashes on an older manifest.

### Impact

The phase can pass with both lifecycle history and backward compatibility
broken.

### Root Cause

The distribution SLO was used as the entire phase gate.

### Required Correction

Add golden old-instance fixtures for `check` and `status`, supersession
cross-link tests, release checksum tests, and changelog assertions. Keep the
speed target separately.

### Proposed Specification Diff

Split MS-501 into functional acceptance and installation SLO clauses.

### Acceptance Evidence

The released binary passes old and current fixture suites, supersession is
bidirectionally linked and logged, and install timing is measured on a named
clean image.

### Alternatives Considered

Deferring compatibility tests until a break occurs was rejected because
DEC-011's no-migration decision depends on compatibility.

### Residual Risk

Fixtures cover known old shapes, not every manually edited instance.

### Related Findings

FND-003

## FND-010 - MS-601 has no bounded implementation or correctness oracle

- **Severity:** High
- **Confidence:** High
- **Category:** Milestone testability
- **Affected sections:** Success criteria; PHASE-06
- **Affected requirements:** None
- **Affected phases:** PHASE-06
- **Blocks implementation:** PHASE-06

### Problem

"Implements from a packet alone" does not name what gets implemented, which
implementation system runs, what source access is prohibited, or how correctness
is judged.

### Evidence

[`DEC-002`, Decision](../decisions/DEC-002-durable-tiered-record-as-the-product.md)
requires a packet sufficient for a fresh implementation agent. The blueprint's
[MS-601](../blueprint.md#phases-and-milestones) repeats the outcome without
defining a test.

### Failure Scenario

A trivial packet that asks for a one-file program passes, while a realistic
packet fails because of an unrelated tool outage. Either result can be claimed
as evidence.

### Impact

The terminal product has no reproducible acceptance gate.

### Root Cause

Packet sufficiency was stated as a property without a canonical task and oracle.

### Required Correction

Define a canonical packet, bounded implementation target, fresh isolated agent,
allowed tools, prohibited chat and source access, time budget, and executable
acceptance tests.

### Proposed Specification Diff

Replace the open-ended milestone with a named fixture implementation and test
suite. Keep a separate real-project handoff as dogfood evidence.

### Acceptance Evidence

At least one fresh agent produces an implementation that passes the fixture
tests using only the packet and documented implementation system.

### Alternatives Considered

Human judgment alone was rejected because the milestone is supposed to make the
phase verifiable.

### Residual Risk

Agent variance remains. The milestone should assess packet completeness, not
rank models.

### Related Findings

None

## FND-011 - Methodology version and generating CLI version are not distinguished

- **Severity:** Medium
- **Confidence:** High
- **Category:** Versioning
- **Affected sections:** Locked constraints; Instance
- **Affected requirements:** None
- **Affected phases:** PHASE-01, PHASE-05
- **Blocks implementation:** No

### Problem

DEC-004 pins the methodology version. DEC-011 requires the generating CLI
version. The blueprint names only the generating CLI version in the scaffold
description and does not define whether the two versions are identical.

### Evidence

[`DEC-004`, Decision](../decisions/DEC-004-template-owned-self-contained-multi-runtime.md)
requires a methodology version pin.
[`DEC-011`, Decision item 2](../decisions/DEC-011-defer-migrations.md) requires
a generating CLI version. The blueprint's
[Instance](../blueprint.md#instance-per-idea) section lists only the latter.

### Failure Scenario

A patch CLI release ships unchanged methodology. An instance records the CLI
version, and later tooling incorrectly infers a methodology change or cannot
identify the emitted contract set.

### Impact

Compatibility and audit logic depend on an unstated one-to-one version
assumption.

### Root Cause

Distribution and methodology release identities were treated as one value.

### Required Correction

Store both versions, or explicitly define and test a permanent one-to-one
version invariant.

### Proposed Specification Diff

Add `methodology_version` and `generated_by_cli_version` to the manifest
contract.

### Acceptance Evidence

Tests cover a CLI-only patch release and a methodology release.

### Alternatives Considered

One shared version remains viable only if every binary release is defined to be
a methodology release.

### Residual Risk

Manually adopted capabilities can still diverge from the pinned methodology
version and need declared deviations.

### Related Findings

FND-009

## FND-012 - The emitted location of runtime skills is ambiguous

- **Severity:** Medium
- **Confidence:** High
- **Category:** Self-containment
- **Affected sections:** Repository anatomy; Instance
- **Affected requirements:** None
- **Affected phases:** PHASE-01
- **Blocks implementation:** No

### Problem

DEC-004 names `.agents/skills/` as part of a self-contained instance. The
blueprint says skills live under the master's `program/` tree, but its instance
scaffold list does not name an emitted skill path.

### Evidence

[`DEC-004`, Decision](../decisions/DEC-004-template-owned-self-contained-multi-runtime.md)
lists `.agents/skills/`.
[`DEC-010`, Decision item 2](../decisions/DEC-010-mycelium-is-a-cli.md)
reaffirms emitted skills. The blueprint's
[Master](../blueprint.md#master-this-repo) anatomy places skills under
`program/`, while [Instance](../blueprint.md#instance-per-idea) omits their
destination.

### Failure Scenario

The scaffold includes skill source files under a directory no runtime discovers.
The repository contains the methodology, but agents never load the wake,
sparring, or manual-floor instructions.

### Impact

Manual operation and multi-runtime self-containment depend on implementer
interpretation.

### Root Cause

The blueprint describes source ownership but not emitted runtime adapter paths.

### Required Correction

Enumerate canonical emitted paths for portable contracts and each supported
runtime adapter. State which files are authority and which are generated
conveniences.

### Proposed Specification Diff

Add `.agents/skills/` to the instance anatomy and define its relationship to
source files under `program/`.

### Acceptance Evidence

The fixture scaffold contains the canonical skill paths, and an agent can
discover the manual floor without the binary.

### Alternatives Considered

Putting all instructions only in `AGENTS.md` was rejected because DEC-004
explicitly retains skills as part of the self-contained instance.

### Residual Risk

Runtime discovery conventions will change. Generated adapters must remain
subordinate to portable contracts.

### Related Findings

FND-008

## FND-013 - OQ-006 and the commissioned run use incompatible ladder terms

- **Severity:** Medium
- **Confidence:** High
- **Category:** Governance
- **Affected sections:** Perspective ladder; Open questions
- **Affected requirements:** None
- **Affected phases:** PHASE-04
- **Blocks implementation:** No

### Problem

OQ-006 asks whether the blueprint gets a council review. The commissioned run
uses one different model. DEC-008 defines that as at most a second opinion, and
only when the commissioning prompt is identical. Calling this run council
dogfood would erase the distinction between the ladder's middle and top rungs.

### Evidence

[`DEC-008`, Decision items 1 and 3](../decisions/DEC-008-perspective-ladder-opt-in-councils.md)
separates one different-model second opinion from multi-model replication and
reconciliation. The blueprint's [OQ-006](../blueprint.md#open-questions) asks
for a council and recommends it as council-contract dogfood.

### Failure Scenario

The blueprint is marked reviewed under OQ-006, and PHASE-04 later cites this
single report as proof that council commissioning and reconciliation already
worked.

### Impact

The first dogfood record teaches the wrong ladder semantics.

### Root Cause

"Different-model review," "second opinion," and "council" were used as
interchangeable labels.

### Required Correction

Disposition OQ-006 explicitly. Either record this as a separate adversarial
review and leave council dogfood for PHASE-04, or commission multiple model
replicas plus reconciliation before blueprint acceptance. If calling it a second
opinion, preserve the original commissioning prompt and prove prompt identity.

### Proposed Specification Diff

Rewrite OQ-006 with the chosen rung and acceptance evidence.

### Acceptance Evidence

The review manifest names the rung, prompt artifact, model provenance, expected
report count, and whether reconciliation is required.

### Alternatives Considered

Treating any different-model read as a council was rejected because it
contradicts DEC-008.

### Residual Risk

This finding does not decide whether a council is worth its cost. Robert owns
that gate.

### Related Findings

FND-008

## Suggested disposition order

1. Resolve FND-001 through FND-005 before PHASE-01 starts.
2. Rewrite each milestone when revising its phase contract, but fix their
   blueprint wording before acceptance so "verifiable" has a stable meaning.
3. Decide FND-013 before describing this run as perspective-ladder dogfood.
4. Address FND-011 and FND-012 in the PHASE-01 manifest and scaffold contract.

No finding in this review is accepted automatically. The blueprint status must
remain Proposed until Robert dispositions the findings and records the
acceptance gate.

## Dispositions

- **Dispositioned by:** Robert Guss
- **Date:** 2026-08-14
- **Method:** Structured disposition session with his agent. Before the
  decisions were put to Robert, the agent independently re-verified the
  findings' load-bearing citations against DEC-002, DEC-003, DEC-004,
  DEC-005, DEC-006, DEC-008, and DEC-011; all quoted authority checked out.
- **Result:** All thirteen findings **accepted**. A blueprint revision stage
  is commissioned. The blueprint remains Proposed until the revision is
  independently validated and Robert accepts it.

### Resolutions chosen where a finding embedded a choice

- **FND-001 — Accepted; front matter wins.** The blueprint and the emitted v2
  templates switch artifact metadata to front matter validated by the sidecar
  schemas, honoring DEC-005 item 3 as written. No superseding decision record
  is needed. The revision must name the single parser contract shared by the
  generator and `mycelium check`. This governs the grammar Mycelium emits and
  validates; retrofitting the master's own `framework/` artifacts is out of
  scope for this finding.
- **FND-011 — Accepted; store both versions.** The instance manifest records
  `methodology_version` and `generated_by_cli_version` as separate fields. No
  one-to-one release invariant is declared.
- **FND-013 — Accepted; this run is an independent adversarial review.** It
  is recorded as neither a council nor a DEC-008 second opinion. OQ-006 is
  dispositioned accordingly in the revision: council dogfood moves to
  PHASE-04. The provenance lesson binds immediately: commissioning prompts
  are stored durably in-repo from this stage onward, starting with the
  revision prompt below.

### Remaining findings

**FND-002 through FND-010 and FND-012 — Accepted as proposed**, with one
scoping rule: milestone corrections (FND-005 through FND-010) are applied at
blueprint-wording granularity, and the full acceptance matrices land in each
phase's contract when that phase is commissioned. The revised blueprint must
record that deferral so "verifiable" keeps a stable meaning without the
blueprint becoming a test plan.

### Commissioned follow-up

The revision stage is commissioned by
[`01-blueprint-revision-prompt.md`](../prompts/01-blueprint-revision-prompt.md)
with attachment manifest
[`blueprint-revision-attachment-manifest.md`](../handoffs/blueprint-revision-attachment-manifest.md).
Allowed file scope for that stage: `framework/blueprint.md` only.
