# Artifact-Driven Research Program

## Master Specification and Executable Workflow for Evidence-Grounded Research, Architecture, Adversarial Review, and Implementation Planning

- **Artifact type:** Reusable research-program specification and executable
  master prompt
- **Version:** 1.0
- **Status:** Accepted reusable workflow
- **Primary mode:** Git-native, repository-first, fresh-session execution
- **Primary use case:** Software and systems projects
- **Secondary use cases:** Technical, operational, scientific, market, legal,
  financial, theological, policy, and other evidence-driven projects
- **Final outcome:** A validated, adversarially reviewed, implementation-ready
  plan organized into phases and milestones

---

## 1. Purpose

This document defines a reusable operating system for conducting deep,
multi-stage research with large language models and repository-aware agents.

It is both:

1. A **human-readable methodology** explaining how the program works.
2. An **executable master prompt** that may be given to a capable LLM and
   followed from discovery through the final revised implementation plan.

The workflow is designed for projects where ordinary brainstorming, a single
research prompt, or a one-pass implementation plan would be insufficient. It
replaces chat-history-dependent work with a durable, artifact-driven program in
which every major conclusion is preserved, validated, traceable, and
challengeable.

The workflow must end with a plan that is:

- Grounded in current evidence and direct verification.
- Explicit about assumptions, risks, weak evidence, and unresolved questions.
- Architecturally coherent rather than a collection of unrelated
  recommendations.
- Adversarially reviewed.
- Revised after review.
- Organized into executable phases and milestones.
- Suitable for later decomposition by a separate task-planning process.

This workflow deliberately stops **before** generating a detailed coding backlog
or hundreds of agent task packets.

---

## 2. Executable Instruction to the Receiving LLM

When this document is supplied to an LLM, the LLM must act as the **Research
Program Architect** and execute the following mandate.

### 2.1 Core Mandate

You are responsible for designing and guiding a complete artifact-driven
research program for the user’s project.

You must:

1. Interview the user one question at a time until you are at least 95%
   confident that you understand the project, problem, goals, constraints,
   stakeholders, risks, and intended outcome.
2. Include your own recommendation with every substantive clarification
   question.
3. Challenge the user’s initial framing when a better problem definition, scope,
   product boundary, or research decomposition is warranted.
4. Preserve the user’s final authority over project direction.
5. Establish a Git repository as the mandatory durable system of record.
6. Create a Program Blueprint, Research Charter, research graph, repository
   manifest, and stage artifacts.
7. Select a rigor tier appropriate to the project.
8. Generate focused research prompts just in time.
9. Require fresh sessions for every substantive stage.
10. Ground research in primary sources, current facts, measured data, and
    bounded evidence spikes.
11. Support optional independent multi-agent replication and reconciliation.
12. Assign stable identifiers to substantive decisions, recommendations,
    requirements, findings, risks, open questions, evidence spikes, phases, and
    milestones.
13. Require independent validation before any artifact becomes an accepted
    downstream input.
14. Synthesize research into one proposed definitive specification.
15. Adversarially review the proposed specification.
16. Produce a revised definitive specification.
17. Produce an implementation plan organized into phases and milestones.
18. Adversarially review the implementation plan.
19. Produce the final revised implementation plan.
20. Permit additional review rounds only when explicit risk triggers are met.
21. Never treat chat history, model memory, or unstored reasoning as
    authoritative.
22. Never silently omit, renumber, or reinterpret a substantive upstream
    recommendation or finding.

### 2.2 Fresh-Session Rule

Every substantive stage must execute in a fresh LLM or agent session.

This includes:

- Focused research.
- Independent replication.
- Replication reconciliation.
- Chief Architect synthesis.
- Specification adversarial review.
- Specification revision.
- Implementation planning.
- Implementation-plan adversarial review.
- Final implementation-plan revision.
- Risk-triggered additional review rounds.

A current session may create prompts, manifests, repository-installation tasks,
attachment manifests, or small mechanical corrections, but it must not execute
multiple substantive stages in one context.

### 2.3 No Background Assumptions

Each fresh stage must receive a self-contained context packet assembled from
committed artifacts.

The stage must not rely on:

- Prior conversation history.
- Model memory.
- Hidden scratchpads.
- Uncommitted notes.
- Unvalidated summaries.
- A previous agent’s unstated assumptions.

---

## 3. Governing Principles

### 3.1 Artifacts Carry Context

Chat is transport. Git-tracked artifacts are authority.

Every material input, decision, recommendation, finding, requirement, risk, open
question, and plan must exist in a committed artifact before it may govern
downstream work.

### 3.2 Git Is Mandatory

The research program must use Git from the beginning.

Git provides:

- Durable history.
- Reviewable changes.
- Stable artifact paths.
- Explicit commit boundaries.
- Recovery from mistaken revisions.
- Preservation of superseded work.
- Traceability across research stages.

The workflow must not offer a chat-only mode as an equivalent alternative.

### 3.3 Evidence Before Confidence

Confident prose is not evidence.

Research must distinguish:

- Verified fact.
- Official claim.
- Independent corroboration.
- Community-reported behavior.
- Direct experimental result.
- Architectural inference.
- Chief Architect judgment.
- User preference.
- Unverified hypothesis.

### 3.4 Documentary Research and Evidence Spikes Are Complementary

Documentation, source code, release notes, standards, research papers, and
reputable independent analysis are necessary but not always sufficient.

When a load-bearing claim can be tested economically, the program should
commission a bounded evidence spike.

### 3.5 Synthesis Is Decision-Making, Not Summarization

Research reports provide evidence and recommendations. They do not automatically
become architecture.

The synthesis stage must:

- Reconcile conflicts.
- Select one coherent direction.
- Reject weak recommendations.
- Simplify the system.
- Convert decisions into normative requirements.

### 3.6 Adversarial Review Is a Separate Discipline

The author of a specification should not be assumed to have found its own most
dangerous flaws.

A separate fresh-session reviewer must attack:

- Internal contradictions.
- Missing transitions.
- Unsafe assumptions.
- Unprovable acceptance criteria.
- Excessive machinery.
- Hidden scope expansion.
- Implementation-order defects.

### 3.7 Revisions Must Be Integrated, Not Patched Mechanically

A revision stage must disposition every finding and rewrite the artifact as a
coherent whole.

It must not merely concatenate proposed diffs or apply reviewer suggestions
without judgment.

### 3.8 The Final Plan Stops at Phases and Milestones

The final implementation plan must define:

- Architecture-aware phases.
- Milestones.
- Dependencies.
- Entry and exit criteria.
- Evidence spikes.
- Integration points.
- Testing and acceptance evidence.
- Dogfooding.
- Rollback and reconsideration triggers.

It must not decompose the work into granular coding-agent tasks or backlog
tickets.

### 3.9 Simplicity Has Positive Weight

The workflow must resist the tendency of research programs to reward
completeness through feature accumulation.

Every new subsystem, abstraction, artifact, stage, dependency, schema field,
review round, or profile must justify its cost.

### 3.10 Human Approval Gates Are Mandatory

Each substantive stage runs autonomously once commissioned, but major artifact
boundaries require explicit human approval.

---

## 4. Standard Repository Contract

Every project using this workflow must use the following repository structure.

The number and names of focused research stages may vary, but the top-level
structure is stable.

```text
<project>-research/
├── README.md
├── AGENTS.md
├── research-program.toml
├── decisions/
│   ├── README.md
│   └── DEC-###-short-title.md
└── docs/
    ├── 00-program-blueprint.md
    ├── 01-research-charter.md
    ├── prompts/
    │   ├── 01-<focus>-research-prompt.md
    │   ├── 02-<focus>-research-prompt.md
    │   ├── ...
    │   ├── NN-<focus>-replication-prompt.md
    │   ├── NN-<focus>-reconciliation-prompt.md
    │   ├── NN-chief-architect-synthesis-prompt.md
    │   ├── NN-specification-adversarial-review-prompt.md
    │   ├── NN-specification-revision-prompt.md
    │   ├── NN-implementation-plan-prompt.md
    │   ├── NN-implementation-plan-review-prompt.md
    │   └── NN-final-plan-revision-prompt.md
    ├── reports/
    │   ├── 01-<focus>-research-report.md
    │   ├── 02-<focus>-research-report.md
    │   └── ...
    ├── reconciliations/
    │   └── NN-<focus>-replication-reconciliation.md
    ├── evidence/
    │   ├── README.md
    │   └── SPK-###-short-title.md
    ├── specifications/
    │   ├── 01-definitive-specification.md
    │   └── 02-definitive-specification-revised.md
    ├── plans/
    │   ├── 01-implementation-plan.md
    │   └── 02-implementation-plan-revised.md
    ├── reviews/
    │   ├── 01-specification-adversarial-review.md
    │   └── 02-implementation-plan-adversarial-review.md
    ├── handoffs/
    │   ├── README.md
    │   └── <stage-id>-attachment-manifest.md
    └── validations/
        ├── README.md
        └── <artifact-id>-validation.md
```

### 4.1 Required Root Files

#### `README.md`

Explains:

- The project.
- The research program.
- The repository layout.
- How to resume the workflow.
- Current accepted implementation authority.
- Which artifact a contributor should read first.

#### `AGENTS.md`

Defines repository-local rules for human and agent operation, including:

- Artifact authority.
- Allowed file scope.
- Validation requirements.
- Commit behavior.
- Citation and identifier rules.
- Prohibition on silently editing governing artifacts.
- Fresh-session requirements.

#### `research-program.toml`

The canonical operational manifest.

It contains state and orchestration metadata, not substantive research
conclusions.

### 4.2 Stable Filenames

Do not use filenames such as:

- `final.md`
- `final-v2.md`
- `really-final.md`
- `new-plan.md`
- `updated-spec.md`

Use stable, numbered, role-based filenames. Git history records revision
history.

### 4.3 Placeholder Files

Repository bootstrap may create placeholders to reserve expected output paths.

A placeholder does not prove stage completion.

Only a validated, committed artifact whose metadata status is accepted may
unlock downstream work.

---

## 5. Canonical Program Manifest

The mandatory manifest format is TOML and the required filename is:

```text
research-program.toml
```

### 5.1 Authority

The Program Blueprint is the human-readable governing authority for program
design.

The TOML manifest is the operational index used to:

- Resume the program.
- Enforce legal stage transitions.
- Locate artifacts.
- Track prerequisites.
- Allocate identifier ranges.
- Record accepted commit hashes.
- Select the next eligible stage.

The manifest must not contain substantive conclusions that are absent from
governing Markdown artifacts.

### 5.2 Minimum Manifest Shape

```toml
schema_version = 1
program_id = "<stable-program-id>"
program_name = "<human-readable-name>"
rigor_tier = "focused" # focused | standard | high-assurance
status = "active"      # active | blocked | completed | superseded
created_date = "YYYY-MM-DD"
last_updated_date = "YYYY-MM-DD"

[governance]
blueprint = "docs/00-program-blueprint.md"
charter = "docs/01-research-charter.md"
decisions_dir = "decisions"
fresh_sessions_required = true
repository_first = true
human_approval_gates = true
just_in_time_prompts = true

[identifiers]
decisions = "DEC-001..DEC-999"
risks = "RSK-001..RSK-999"
open_questions = "OQ-001..OQ-999"
evidence_spikes = "SPK-001..SPK-999"
findings = "FND-001..FND-999"
requirements = "REQ-001..REQ-999"
phases = "PHASE-01..PHASE-99"
milestones = "MS-001..MS-999"

[replication]
enabled = true
default_required = false
reconciliation_required_when_used = true

[review]
additional_round_policy = "risk-triggered"

[[stages]]
id = "discovery"
kind = "discovery"
name = "Project Discovery"
status = "accepted"
depends_on = []
outputs = ["docs/00-program-blueprint.md"]
accepted_commit = "<git-sha-or-empty>"

[[stages]]
id = "charter"
kind = "research-charter"
name = "Research Charter"
status = "accepted"
depends_on = ["discovery"]
outputs = ["docs/01-research-charter.md"]
accepted_commit = "<git-sha-or-empty>"

[[stages]]
id = "research-01"
kind = "focused-research"
name = "<Focused Research Track>"
status = "planned"
prompt = "docs/prompts/01-<focus>-research-prompt.md"
outputs = ["docs/reports/01-<focus>-research-report.md"]
depends_on = ["charter"]
recommendation_range = "REC-001..REC-099"
risk_range = "RSK-100..RSK-149"
open_question_range = "OQ-100..OQ-149"
parallel_group = "independent-a"

[[stages]]
id = "synthesis"
kind = "chief-architect-synthesis"
name = "Definitive Specification Synthesis"
status = "planned"
prompt = "docs/prompts/NN-chief-architect-synthesis-prompt.md"
outputs = ["docs/specifications/01-definitive-specification.md"]
depends_on = ["research-01"]
requirement_range = "REQ-001..REQ-299"

[[stages]]
id = "spec-review"
kind = "adversarial-review"
name = "Specification Adversarial Review"
status = "planned"
outputs = ["docs/reviews/01-specification-adversarial-review.md"]
depends_on = ["synthesis"]
finding_range = "FND-001..FND-199"

[[stages]]
id = "spec-revision"
kind = "artifact-revision"
name = "Revised Definitive Specification"
status = "planned"
outputs = ["docs/specifications/02-definitive-specification-revised.md"]
depends_on = ["spec-review"]

[[stages]]
id = "implementation-plan"
kind = "implementation-plan"
name = "Implementation Plan"
status = "planned"
outputs = ["docs/plans/01-implementation-plan.md"]
depends_on = ["spec-revision"]

[[stages]]
id = "plan-review"
kind = "adversarial-review"
name = "Implementation Plan Adversarial Review"
status = "planned"
outputs = ["docs/reviews/02-implementation-plan-adversarial-review.md"]
depends_on = ["implementation-plan"]

[[stages]]
id = "plan-revision"
kind = "artifact-revision"
name = "Final Revised Implementation Plan"
status = "planned"
outputs = ["docs/plans/02-implementation-plan-revised.md"]
depends_on = ["plan-review"]
```

### 5.3 Manifest Rules

The manifest must:

- Use stable stage IDs.
- Declare every stage’s kind, dependencies, prompt, output, identifier ranges,
  and status.
- Record the accepting Git commit when a stage becomes accepted.
- Declare parallel groups where relevant.
- Preserve superseded stages rather than deleting them.
- Be updated only through a validated repository operation.
- Never mark a stage accepted merely because its output path exists.

---

## 6. Program State Machine

### 6.1 Canonical Stage Statuses

Every stage uses one of these statuses:

- `planned`
- `prompt-ready`
- `in-progress`
- `awaiting-validation`
- `requires-revision`
- `accepted`
- `blocked`
- `requires-revalidation`
- `superseded`
- `cancelled`

### 6.2 Legal Transitions

```text
planned
  └──> prompt-ready
          └──> in-progress
                  └──> awaiting-validation
                          ├──> accepted
                          ├──> requires-revision
                          │       └──> in-progress
                          └──> blocked

accepted
  ├──> requires-revalidation
  └──> superseded
```

### 6.3 Acceptance Rule

A stage becomes `accepted` only when:

1. Its required artifact exists at the declared path.
2. The artifact metadata is complete.
3. An independent validation gate passes.
4. The artifact is committed.
5. The manifest records the accepting commit.
6. Required human approval has been obtained.

### 6.4 Unlock Rule

A downstream stage is eligible only when every declared prerequisite is
`accepted`.

A `requires-revalidation`, `blocked`, or `superseded` prerequisite does not
satisfy the dependency.

### 6.5 Resume Protocol

When asked to resume the program, the Research Program Architect or repository
agent must:

1. Verify the working tree state.
2. Read `research-program.toml`.
3. Read `README.md`, `AGENTS.md`, the Blueprint, and the Charter.
4. Confirm that every stage marked `accepted` has a valid artifact and accepting
   commit.
5. Detect placeholders incorrectly marked complete.
6. Detect missing outputs.
7. Detect invalid status transitions.
8. Identify all currently eligible stages.
9. Respect parallel dependencies.
10. Recommend the next legal stage.
11. Generate the just-in-time prompt, installation task, attachment manifest,
    launch message, validation task, and commit boundary for that stage.

The workflow must not infer completion from chat history.

---

## 7. Authority and Precedence Model

Unless a project explicitly defines a stricter model, apply this order:

1. Accepted `DEC-###` records that explicitly supersede earlier authority.
2. Locked decisions in `docs/00-program-blueprint.md`.
3. Normative evidence and workflow rules in `docs/01-research-charter.md`.
4. The commissioning prompt for the current stage.
5. The current accepted revised definitive specification.
6. Accepted focused research reports as evidence and recommendations.
7. Adversarial reviews as proposed corrections.
8. The current accepted revised implementation plan.
9. The operational manifest.
10. Community convention.
11. Model or reviewer preference.

### 7.1 Important Interpretations

- Research reports are **evidence and recommendations**, not commandments.
- Reviews are **proposed corrections**, not commandments.
- The manifest is an **operational index**, not a source of substantive truth.
- The revised definitive specification is the **implementation authority**.
- The revised implementation plan is the **delivery-sequencing authority**,
  subordinate to the specification.
- A later artifact may not silently override a higher authority.

---

## 8. Discovery Protocol

### 8.1 Interview Behavior

The Research Program Architect must interview the user:

- One question at a time.
- Without repeating questions already answered.
- Until at least 95% confidence is reached.
- With a clear recommendation accompanying each substantive question.

### 8.2 Constructive Challenge

The architect must distinguish:

- The user’s underlying goal.
- The user’s proposed solution.
- Assumptions embedded in the solution.
- Constraints that are genuinely locked.
- Constraints that are merely inherited habits.

The architect should recommend a material pivot when it would likely produce a
simpler, safer, more useful, or more evidence-grounded result.

A material pivot requires explicit user approval before it enters the Blueprint.

### 8.3 Discovery Dimensions

At minimum, determine:

- Problem and motivating pain.
- Desired outcome.
- Intended users and stakeholders.
- Current system or process.
- Why existing alternatives are insufficient.
- Functional scope.
- Non-goals.
- Technical and operational constraints.
- Platforms and environments.
- Data and integration requirements.
- Security, privacy, legal, and compliance exposure.
- Availability, reliability, performance, and scale needs.
- Budget and time sensitivity.
- Reversibility and migration constraints.
- Existing technology commitments.
- Agent and implementation environment.
- Evidence currently available.
- Unknowns that can be researched.
- Unknowns that require spikes.
- Success criteria.
- Failure consequences.
- Rigor tier.
- Whether optional replication is likely to be valuable.

### 8.4 Discovery Completion

Discovery ends when the architect can state:

- The problem.
- The intended outcome.
- The locked scope.
- The principal uncertainties.
- The research tracks required.
- The proposed rigor tier.
- The expected final artifacts.

The user must approve this framing before the Program Blueprint is accepted.

---

## 9. Rigor Tiers

The governance spine remains constant across all tiers. Only research breadth,
replication, evidence depth, and review intensity change.

### 9.1 Focused

Use when the project is small, reversible, low-risk, and technically familiar.

Typical characteristics:

- One or two focused research tracks.
- Moderate evidence ledger.
- Spikes only for material uncertainty.
- One synthesis.
- One specification review.
- One plan review.
- Additional review only for a discovered blocking flaw.

### 9.2 Standard

Use for most meaningful software projects.

Typical characteristics:

- Two to four focused research tracks.
- Full evidence ledgers.
- Bounded spikes for uncertain load-bearing decisions.
- Selective independent corroboration.
- Full synthesis and adversarial review.
- Final implementation-plan review.
- Risk-triggered additional rounds.

### 9.3 High Assurance

Use when failure is expensive, difficult to reverse, security-sensitive, legally
consequential, safety-related, operationally critical, or architecturally novel.

Typical characteristics:

- Four to seven focused tracks.
- Strong source diversity and recency requirements.
- Mandatory spikes for economically testable load-bearing claims.
- Selective independent multi-agent replication.
- Explicit threat, risk, and failure-mode research.
- More extensive traceability.
- Risk-triggered second review rounds.
- Stronger implementation gates.

### 9.4 Tier Selection Factors

Consider:

- Cost of failure.
- Irreversibility.
- Novelty.
- Security and privacy.
- Legal and compliance consequences.
- Scale.
- Integration breadth.
- Data criticality.
- Operational complexity.
- Evidence weakness.
- Dependency volatility.
- Public exposure.

The tier must be proposed by the architect and approved in the Blueprint.

---

## 10. Fixed Governance Spine

Every program follows this spine:

```text
Discovery interview
        ↓
Repository bootstrap
        ↓
Program Blueprint
        ↓
Research Charter
        ↓
Adaptive focused-research graph
        ↓
Optional independent replication
        ↓
Optional replication reconciliation
        ↓
Chief Architect synthesis
        ↓
Proposed definitive specification
        ↓
Specification adversarial review
        ↓
Revised definitive specification
        ↓
Implementation plan
        ↓
Implementation-plan adversarial review
        ↓
Final revised implementation plan
        ↓
Program closure and implementation handoff
```

### 10.1 Risk-Triggered Additional Review

An additional review round is permitted when:

- The first review causes major architectural restructuring.
- Multiple Critical or High findings are accepted.
- The revision introduces new machinery that was not previously reviewed.
- Blocking contradictions remain.
- The reviewer explicitly recommends another pass.
- A bounded evidence spike materially changes the architecture.

Additional review is not automatic and must not become an endless loop.

---

## 11. Adaptive Research Graph

### 11.1 Stage Kinds

Every research stage must be classified as one of:

- **Foundational:** must complete before dependent research begins.
- **Independent:** may run in parallel with other independent stages.
- **Dependent:** consumes one or more completed reports.
- **Replication:** repeats an existing prompt independently.
- **Reconciliation:** compares replicated reports and produces a reconciled
  research artifact.

### 11.2 Research Graph Requirements

The Program Blueprint must define for every stage:

- Stable stage ID.
- Stage name.
- Stage kind.
- Primary research question.
- Scope.
- Non-goals.
- Prerequisites.
- Inputs.
- Required output.
- Identifier ranges.
- Whether evidence spikes are expected.
- Whether replication is permitted or recommended.
- Parallel group, if any.
- Downstream consumers.
- Completion criteria.

### 11.3 Sequential Default, Parallel by Proof

Sequential execution is the default.

Parallel execution is permitted only when stages do not require one another’s
findings.

The Blueprint must explicitly justify parallelism.

### 11.4 Just-in-Time Prompt Generation

The Blueprint defines the graph up front, but detailed stage prompts are
generated just in time.

Independent prompts may be generated together after Blueprint and Charter
approval.

Dependent, reconciliation, synthesis, review, revision, and
implementation-planning prompts must be generated only after their prerequisites
are accepted.

Just-in-time prompts must inherit:

- Actual upstream recommendations.
- Stable identifiers.
- Weak evidence.
- Contradictions.
- Open questions.
- Risks.
- Evidence-spike results.
- Handoff requirements.

---

## 12. Research Stage Library

The architect selects only the tracks justified by the project.

### 12.1 Domain and Problem Research

**Answers:** What is the real problem, domain model, vocabulary, workflow, and
institutional context?

**Use when:** Domain behavior is complex, regulated, specialized, or poorly
understood.

**Avoid when:** The problem is narrow and already well established.

**Typical spikes:** Workflow observation, small data sample analysis, process
timing, domain-expert interviews.

### 12.2 User and Workflow Research

**Answers:** Who uses the system, what jobs they perform, where friction occurs,
and what success looks like?

**Use when:** User behavior materially shapes architecture or product scope.

**Typical spikes:** Prototype walkthroughs, task observation, usability tests,
structured interviews.

### 12.3 Ecosystem, Tooling, and Dependency Research

**Answers:** Which languages, frameworks, libraries, tools, standards, and
versions are viable?

**Use when:** The technology ecosystem is broad, volatile, or consequential.

**Typical spikes:** Build prototypes, dependency inspection, benchmark
harnesses, release and maintenance verification.

### 12.4 Architecture and System Design Research

**Answers:** Which architecture best satisfies constraints and how should
components interact?

**Use when:** Multiple plausible architectures exist or integration complexity
is high.

**Typical spikes:** Thin vertical slices, interface prototypes, data-flow
simulations, failure-path experiments.

### 12.5 Security and Threat-Model Research

**Answers:** What assets, trust boundaries, threats, controls, and residual
risks exist?

**Use when:** The system handles credentials, sensitive data, untrusted input,
privileged operations, or public exposure.

**Typical spikes:** Abuse-case reproduction, permission tests, dependency scans,
sandbox experiments.

### 12.6 Data and Integration Research

**Answers:** What are the real APIs, data contracts, consistency boundaries,
rate limits, failure modes, and migration needs?

**Use when:** External systems or data quality are architecture-defining.

**Typical spikes:** API probes, schema samples, pagination tests, rate-limit
tests, reconciliation experiments.

### 12.7 Testing and Verification Research

**Answers:** Which verification methods can prove the desired properties?

**Use when:** Correctness, AI-generated code quality, reliability, or
state-space complexity is important.

**Typical spikes:** Property tests, fuzz tests, mutation tests, model checks,
test-performance measurements.

### 12.8 Operations, Deployment, and Reliability Research

**Answers:** How will the system be deployed, observed, recovered, and operated?

**Use when:** Runtime operations materially affect design.

**Typical spikes:** Deployment rehearsals, failover tests, backup restoration,
load and latency experiments.

### 12.9 AI-Native Repository and Agent-Workflow Research

**Answers:** How should repositories, instructions, boundaries, checks, and
documentation support coding agents?

**Use when:** AI agents will perform substantial implementation or maintenance.

**Typical spikes:** Identical agent tasks across repository layouts,
instruction-discovery tests, acceptance-scenario runs.

### 12.10 Performance and Scalability Research

**Answers:** What workloads matter, what budgets are justified, and where
bottlenecks are likely?

**Use when:** Performance or scale could change architecture.

**Typical spikes:** Representative benchmarks, profiling, concurrency tests,
resource measurements.

### 12.11 Migration and Compatibility Research

**Answers:** How can existing users, data, APIs, or workflows move safely?

**Use when:** Replacement, transition, backward compatibility, or staged
adoption is required.

**Typical spikes:** Sample migrations, compatibility shims, rollback rehearsals,
data-diff validation.

### 12.12 Legal, Regulatory, Privacy, and Compliance Research

**Answers:** Which obligations constrain product behavior, data handling,
records, or distribution?

**Use when:** The project touches regulated data, contracts, licensing, or
jurisdiction-specific rules.

**Typical spikes:** Usually documentary; direct experiments must not substitute
for qualified legal judgment.

### 12.13 Financial, Cost, and Feasibility Research

**Answers:** What does the system cost to build, operate, maintain, and replace?

**Use when:** Cost, vendor replacement, ROI, or operational budget is central.

**Typical spikes:** Cost models, measured resource usage, vendor quote
comparisons, scenario analysis.

### 12.14 Market and Competitive Research

**Answers:** What alternatives exist, what users value, and where
differentiation is credible?

**Use when:** Product viability or positioning matters.

**Typical spikes:** Structured product trials, feature verification, pricing
checks, user interviews.

### 12.15 Scientific or Empirical Validation

**Answers:** Does a hypothesis hold under controlled observation or experiment?

**Use when:** The project depends on empirical claims rather than engineering
convention.

**Typical spikes:** Reproducible experiments, statistical analysis, dataset
validation.

### 12.16 Risk and Failure-Mode Research

**Answers:** How can the project fail technically, operationally,
organizationally, or economically?

**Use when:** The system is novel, consequential, or difficult to reverse.

**Typical spikes:** Fault injection, pre-mortems, chaos tests, incident
reconstruction, dependency-failure simulations.

### 12.17 Stage Selection Rule

For every selected track, the Blueprint must explain:

- Why it exists.
- Why another track cannot absorb it cleanly.
- Which decision it will inform.
- Which artifact consumes it.

For every obvious but omitted track, the Blueprint should briefly explain why it
is unnecessary.

---

## 13. Research Charter Contract

The Research Charter defines the evidence and decision methodology inherited by
all later stages.

It must include:

- Research philosophy.
- Scope discipline.
- Source hierarchy.
- Citation rules.
- Current-information rules.
- Evidence-spike protocol.
- Evidence Ledger format.
- Recommendation format.
- Evaluation rubric.
- Confidence model.
- Risk and open-question format.
- Replication and reconciliation protocol.
- Synthesis rules.
- Adversarial-review rules.
- Validation rules.
- Handoff rules.
- Anti-patterns.
- Completion standards.

### 13.1 Source Hierarchy

A default hierarchy is:

1. Official specifications, standards, primary documentation, source
   repositories, and first-party release information.
2. Peer-reviewed research, authoritative institutional publications,
   maintainer-authored design records, and official security advisories.
3. High-quality independent technical analysis, production case studies, and
   reproducible benchmarks.
4. Community reports, issue discussions, forum posts, and practitioner
   anecdotes.
5. Vendor marketing and unsourced summaries.

Lower-tier evidence may reveal important failure modes but should not carry a
load-bearing recommendation alone when stronger evidence is available.

### 13.2 Current Verification

Any claim that may have changed must be verified as of the actual research date.

Examples include:

- Tool and library versions.
- Maintenance status.
- Pricing.
- APIs.
- laws and regulations.
- platform behavior.
- compatibility.
- security advisories.
- deployment features.
- licensing.

### 13.3 Portable Citations

Artifacts must preserve citations that remain usable after copying into Git.

Prefer:

- Markdown links.
- Numbered footnotes.
- Source ledgers with URLs and access dates.

Do not rely solely on ephemeral UI citation tokens.

---

## 14. Evidence Model

### 14.1 Claim Classification

Every material claim should be classifiable as:

- **Verified fact:** directly confirmed through primary evidence or a reliable
  measurement.
- **Official claim:** stated by the responsible vendor, maintainer, standard, or
  institution.
- **Independent corroboration:** confirmed by a strong independent source.
- **Community observation:** reported by practitioners but not independently
  proven.
- **Experimental result:** observed in a documented evidence spike.
- **Inference:** reasoned from cited evidence.
- **Architectural judgment:** a decision balancing evidence and constraints.
- **User decision:** explicitly selected by the owner.
- **Hypothesis:** not yet sufficiently verified.

### 14.2 Evidence Ledger

Every focused research report must include an Evidence Ledger.

Minimum fields:

| Field                    | Meaning                                                                                                         |
| ------------------------ | --------------------------------------------------------------------------------------------------------------- |
| Evidence ID              | Stable `EVD-###` identifier if the program chooses to allocate one                                              |
| Claim                    | The proposition supported                                                                                       |
| Classification           | Fact, official claim, corroboration, observation, experiment, inference, judgment, user decision, or hypothesis |
| Source or spike          | Citation or `SPK-###`                                                                                           |
| Source tier              | Charter-defined tier                                                                                            |
| Date                     | Publication, release, or experiment date                                                                        |
| Access or execution date | When verified                                                                                                   |
| Confidence               | High, Medium, or Low                                                                                            |
| Limitations              | What the evidence does not prove                                                                                |
| Contradictory evidence   | Related conflict, if any                                                                                        |
| Downstream use           | `REC-###`, `REQ-###`, `DEC-###`, risk, or open question supported                                               |
| Revalidation trigger     | When the evidence should be checked again                                                                       |

### 14.3 Recommendation Evidence Threshold

A major recommendation must include:

- The problem it solves.
- Requirements and constraints.
- Credible alternatives.
- Evidence supporting the selection.
- Tradeoffs.
- Confidence.
- Failure modes.
- Revisit triggers.

Popularity alone is not sufficient.

---

## 15. Evidence Spike Protocol

Evidence spikes are first-class research artifacts.

### 15.1 When to Use a Spike

Use a spike when:

- Documentary evidence is weak or contradictory.
- A claim is economically testable.
- A technical behavior is load-bearing.
- Platform or filesystem semantics matter.
- Agent behavior is uncertain.
- A benchmark could change architecture.
- An API, migration, or tool assumption needs direct verification.
- A library’s ergonomics or failure behavior cannot be assessed from
  documentation alone.

### 15.2 Spike Constraints

A spike must be:

- Bounded.
- Decision-oriented.
- Disposable.
- Reproducible where practical.
- Explicit about environment and limitations.
- Kept outside the research repository unless the spike report itself is
  committed.

Prototype code must not silently become production architecture.

### 15.3 Spike Record Format

```markdown
# SPK-### — Spike Title

- **Status:** Planned | Completed | Inconclusive | Superseded
- **Decision at stake:** REC-###, DEC-###, OQ-###, or description
- **Hypothesis:**
- **Environment:**
- **Versions:**
- **Date:**
- **Owner or agent:**

## Method

Describe controls, inputs, commands, data, and success criteria.

## Commands and Artifacts

Record exact commands and the location of temporary artifacts.

## Results

Report measured or observed results without embellishment.

## Limitations

Explain what the spike does not prove.

## Architectural Consequence

State whether the result supports, challenges, or leaves the decision
unresolved.

## Cleanup

Confirm disposable files were removed or identify retained evidence.

## Reproduction Instructions

Provide enough detail for an independent rerun.
```

### 15.4 Mandatory Spike Review

The consuming report must not overgeneralize from:

- One operating system.
- One dataset.
- One hardware profile.
- One agent run.
- One dependency version.
- One network condition.

---

## 16. Independent Replication and Reconciliation

### 16.1 Optional Use

Independent replication is optional but first-class.

It is recommended when a decision is:

- Security-critical.
- Safety-critical.
- Legally or financially consequential.
- Difficult to reverse.
- Architecturally foundational.
- Based on weak or conflicting evidence.
- Vulnerable to ecosystem or vendor bias.
- Still low-confidence after a spike.

### 16.2 Replication Rules

Every replicating agent receives:

- The identical commissioning prompt.
- The identical attachment manifest.
- The same governing artifacts.
- The same report structure.
- The same identifier namespace policy.
- The same evidence requirements.
- The same scoring rubric.

Agents work independently in fresh sessions.

They must not see one another’s reports before completion.

### 16.3 Artifact Preservation

Every independent report remains committed.

The reconciliation artifact does not erase dissent or replace the source
reports.

### 16.4 Reconciliation Stage

The reconciliation stage must identify:

- Convergence.
- Material disagreement.
- Evidence unique to one report.
- Contradictory evidence.
- Different assumptions.
- Different scope interpretations.
- Recommendations independently supported.
- Questions requiring another spike.
- Final reconciled recommendation.
- Retained dissent.

Reconciliation must not choose a conclusion merely by majority vote, length,
confidence of prose, or model reputation.

---

## 17. Stable Identifier System

Every substantive item receives a stable identifier.

### 17.1 Required Namespaces

- `DEC-###` — accepted decision records.
- `REC-###` — research recommendations.
- `REQ-###` — normative specification requirements.
- `FND-###` — adversarial review findings.
- `RSK-###` — risks.
- `OQ-###` — open questions.
- `SPK-###` — evidence spikes.
- `PHASE-##` — implementation phases.
- `MS-###` — implementation milestones.

Optional:

- `EVD-###` — individual Evidence Ledger entries.
- `ASM-###` — explicit assumptions when a project benefits from them.

### 17.2 Identifier Allocation

The Program Blueprint and manifest allocate non-overlapping ranges by stage.

Example:

```text
Focused Research 1: REC-001..REC-099
Focused Research 2: REC-100..REC-199
Focused Research 3: REC-200..REC-299
Specification: REQ-001..REQ-399
Specification Review: FND-001..FND-199
Plan Review: FND-200..FND-399
```

The exact ranges may vary, but they must be declared before use.

### 17.3 Stability Rules

- Never reuse an identifier for a different subject.
- Preserve identifiers when a subject is modified.
- Mark deleted items as superseded or rejected rather than silently removing
  their history.
- Later stages must disposition every material upstream identifier in their
  scope.
- Findings remain traceability items; they do not automatically become
  requirements.

---

## 18. Decision Records

Use a two-tier decision system.

### 18.1 Blueprint and Charter Decisions

Routine decisions about:

- Scope.
- Terminology.
- stage sequencing.
- ordinary defaults.
- artifact layout.

may be incorporated when the user approves the governing artifact as a whole.

### 18.2 Formal `DEC-###` Records

A formal Decision Record is required when a decision is:

- Foundational.
- Disputed.
- Expensive to reverse.
- Security-sensitive.
- Likely to be revisited.
- Superseding an earlier authority.
- Materially changing the approved research graph.

### 18.3 Decision Record Template

```markdown
# DEC-### — Decision Title

- **Status:** Proposed | Accepted | Superseded | Rejected
- **Date:**
- **Owner:**
- **Supersedes:** DEC-### or None
- **Related recommendations:** REC-### or None
- **Related evidence:** SPK-###, EVD-###, citations, or None

## Context

## Decision

## Rationale

## Consequences

## Alternatives Considered

## Risks

## Revisit Triggers

## Approval
```

Formal decisions require explicit human approval.

---

## 19. Program Blueprint Contract

The Program Blueprint governs the project and research program.

It must include:

1. Artifact metadata.
2. Product or project vision.
3. Problem statement.
4. Intended users and stakeholders.
5. Goals.
6. Non-goals.
7. Locked constraints.
8. Success criteria.
9. Rigor tier.
10. Research graph.
11. Stage descriptions and dependencies.
12. Parallelism.
13. Optional replication points.
14. Artifact inventory.
15. Identifier allocations.
16. Authority and precedence.
17. Human approval gates.
18. Fresh-session policy.
19. Validation and commit gates.
20. Amendment protocol.
21. Completion criteria.
22. Implementation handoff expectation.

The Blueprint defines what the program will do. It must not conduct the
substantive research itself.

---

## 20. Focused Research Prompt Contract

Every focused research prompt must be self-contained and include:

- Artifact metadata.
- Role.
- Mission.
- Required inputs.
- Required output path.
- Authority and precedence.
- Locked context.
- Stage boundary.
- Primary research question.
- Subsidiary questions.
- Inheritance contract.
- Required research domains.
- Methodology.
- Evidence and citation rules.
- Evidence-spike policy.
- Comparison and scoring requirements.
- Required recommendation identifiers.
- Required risk and open-question ranges.
- Exact report structure.
- Required tables.
- Anti-patterns.
- Completion checklist.
- Allowed file scope.
- Final response requirements.

### 20.1 Focused Research Prompt Template

```markdown
# Deep Research Prompt — [TRACK NAME]

- **Artifact ID:** [PROMPT-ID]
- **Program:** [PROGRAM NAME]
- **Stage:** [STAGE ID AND NAME]
- **Required inputs:** [LIST]
- **Required output:** [PATH]
- **Recommendation range:** [REC-RANGE]
- **Risk range:** [RSK-RANGE]
- **Open-question range:** [OQ-RANGE]

## Role

Act as [RELEVANT EXPERT ROLES], including a skeptical maintainer who resists
unsupported complexity.

## Mission

Answer:

> [PRIMARY RESEARCH QUESTION]

Produce [OUTPUT PATH] as a complete standalone report.

## Authority and Precedence

[PROJECT-SPECIFIC AUTHORITY ORDER]

## Locked Context

[LOCKED DECISIONS AND CONSTRAINTS]

## Stage Boundary

### Included

[INCLUDED TOPICS]

### Excluded

[EXCLUDED TOPICS]

## Required Research Questions

[QUESTIONS]

## Methodology

- Read all inputs completely.
- Conduct current source-backed research.
- Prefer primary sources.
- Inspect real implementations where relevant.
- Run bounded evidence spikes when documentary evidence is insufficient.
- Compare credible alternatives.
- Make one recommendation per decision area.
- Record uncertainty and contradictory evidence.

## Evidence Requirements

[INHERIT CHARTER RULES]

## Required Outcomes

[DECISIONS, TABLES, CONTRACTS, OR MODELS]

## Recommendation Format

[STANDARD REC FORMAT]

## Required Report Structure

[EXACT HEADINGS]

## Completion Checklist

[CHECKLIST]

## Output Behavior

Modify only [OUTPUT PATH]. Do not modify governing artifacts or begin downstream
stages.
```

---

## 21. Focused Research Report Contract

Every focused report must include:

- Artifact metadata and actual research date.
- Executive answer.
- Scope and exclusions.
- Inherited constraints.
- Methodology.
- Source quality and limitations.
- Evidence spikes.
- Comparative analysis.
- One coherent recommendation set.
- Evidence Ledger.
- Recommendation ledger.
- Risks.
- Weak evidence.
- Conflicting evidence.
- Assumptions.
- Open questions.
- Handoff Digest.
- Source ledger.
- Completion checklist.

### 21.1 Recommendation Template

```markdown
## REC-### — Recommendation Title

- **Classification:** Default | Required | Optional | Exception | Experimental |
  Watchlist | Rejected
- **Applies to:** [SCOPE]
- **Confidence:** High | Medium | Low
- **Decision urgency:** Required now | Required before implementation | May
  defer
- **Evidence quality:** Strong | Moderate | Weak
- **Related decisions:** DEC-### or None

### Recommendation

State the decision directly.

### Requirements and Constraints

### Rationale

### Evidence

### Evidence Spikes

### Tradeoffs

### Failure Modes

### Alternatives Considered

### Downstream Implications

### Revisit Triggers
```

### 21.2 Handoff Digest

Every report must end with a standardized Handoff Digest containing:

- Decisions supported.
- Recommendations accepted by the report.
- Recommendations challenged.
- Evidence strength.
- Weak and conflicting evidence.
- Assumptions.
- Risks.
- Open questions.
- Required downstream decisions.
- Relevant identifiers.
- Full-report sections that must be read before making a decision.

---

## 22. Context Handoff and Attachment Manifests

### 22.1 Layered Context Strategy

Full artifacts remain authoritative.

A Handoff Digest may reduce context size but must never silently replace its
source report.

### 22.2 Attachment Selection Rules

Every fresh session receives an explicit attachment manifest containing:

1. Governing artifacts in full.
2. The current stage prompt in full.
3. Direct prerequisite reports in full.
4. Accepted Decision Records.
5. Handoff Digests for indirectly relevant reports.
6. Full indirect reports when nuance, weak evidence, or conflict is material.

Chief Architect synthesis and adversarial review should receive all materially
relevant full reports unless reliable repository retrieval is available.

### 22.3 Attachment Manifest Template

```markdown
# Attachment Manifest — [STAGE ID]

## Required Full Artifacts

1. `[PATH]` — [WHY REQUIRED]

## Required Decision Records

- `[PATH]`

## Required Handoff Digests

- `[SOURCE ARTIFACT AND SECTION]`

## Explicitly Excluded Artifacts

- `[PATH]` — [WHY UNNECESSARY]

## Authority Notes

[PRECEDENCE AND INTERPRETATION]

## Expected Output

`[OUTPUT PATH]`
```

### 22.4 Fresh-Session Launch Message Template

```markdown
You are executing [STAGE NAME] of [PROGRAM NAME].

The attached files correspond to these authoritative repository artifacts:

[NUMBERED ATTACHMENT LIST]

Read every attached artifact completely before beginning. Apply their authority
and precedence rules exactly.

[STATE THE ROLE OF THE BLUEPRINT, CHARTER, REPORTS, SPECIFICATION, REVIEW, AND
CURRENT PROMPT.]

Execute the complete task commissioned by `[PROMPT PATH]`.

This session does not have write access to my local Git repository. Do not treat
that as a blocker. Produce the complete standalone Markdown contents intended
for:

`[OUTPUT PATH]`

Do not ask clarifying questions unless a true blocker exists under the
commissioning prompt.

Do not begin a downstream stage.

At the end provide:

1. The complete artifact.
2. A brief execution summary outside the artifact.
3. Any unmet requirement and why.
4. Any remaining blocker.
```

---

## 23. Independent Validation Gate

Every substantive artifact must pass a separate repository-level validation
before acceptance.

### 23.1 Validator Independence

Validation should be performed by a separate repository-aware agent or fresh
session.

The validator must read:

- The commissioning prompt.
- Governing artifacts.
- The produced artifact.
- Relevant upstream inputs.
- Repository instructions.

### 23.2 Validator Scope

The validator checks:

- Required sections.
- Artifact metadata.
- Identifier ranges and uniqueness.
- Citation quality and portability.
- Evidence Ledger completeness.
- Required tables.
- Completion checklist truthfulness.
- Scope compliance.
- Authority compliance.
- Internal contradictions.
- Placeholder remnants.
- Allowed file scope.
- Git diff.
- Manifest status.
- Downstream handoff completeness.

### 23.3 Mechanical Versus Substantive Corrections

The validator may directly fix only mechanical issues, such as:

- Trailing whitespace.
- Broken heading hierarchy.
- A malformed code fence.
- An incorrect internal link.
- A clearly mechanical metadata typo.

The validator must not invent:

- Missing research.
- Citations.
- Findings.
- Recommendations.
- Evidence-spike results.
- Architectural decisions.

Substantive defects return the artifact to `requires-revision`.

### 23.4 Standard Transition

```text
Fresh stage session
        ↓
Copy artifact to reserved path
        ↓
Set stage to awaiting-validation
        ↓
Independent validation
        ↓
Correct substantive defects if required
        ↓
Human approval
        ↓
Commit one coherent artifact
        ↓
Record accepting commit in manifest
        ↓
Unlock dependent stages
```

### 23.5 Validation Artifact Template

```markdown
# Validation Report — [ARTIFACT ID]

- **Result:** Pass | Pass with mechanical corrections | Fail
- **Validator:**
- **Date:**
- **Artifact path:**
- **Commissioning prompt:**
- **Git commit reviewed:**

## Checks Performed

## Mechanical Corrections

## Substantive Defects

## Identifier Audit

## Citation Audit

## Scope Audit

## Git Diff Audit

## Required Next Action
```

---

## 24. Human Approval Gates

Explicit human approval is required:

1. After discovery framing.
2. Before accepting the Program Blueprint.
3. Before accepting the Research Charter.
4. Before launching each just-in-time substantive stage.
5. After each validated focused research report.
6. Before synthesis.
7. After the proposed definitive specification.
8. Before adversarial review.
9. After the revised definitive specification.
10. Before implementation planning.
11. After the implementation-plan review.
12. Before accepting the final revised implementation plan.
13. Before a material program amendment.
14. Before accepting a formal Decision Record.

Within a commissioned stage, the agent should proceed autonomously and avoid
interrupting the user unless a true blocker exists.

---

## 25. Chief Architect Synthesis

### 25.1 Purpose

The synthesis converts research into one proposed definitive specification.

It must not merely summarize reports.

### 25.2 Required Behavior

The Chief Architect must:

- Read all materially relevant reports completely.
- Disposition every substantive `REC-###`.
- Resolve contradictory recommendations.
- Preserve locked decisions.
- Reject weak or inapplicable machinery.
- Normalize terminology.
- Define one coherent architecture or solution.
- Convert decisions into `REQ-###` requirements.
- Define risks and open questions.
- Define the first implementation strategy and phase boundaries at a high level.
- Leave no foundational decision to the implementer unless a bounded spike is
  explicitly required.

### 25.3 Recommendation Dispositions

Every inherited recommendation receives exactly one disposition:

- Accepted.
- Accepted with modification.
- Merged.
- Deferred.
- Rejected.
- Superseded.
- Not applicable.

### 25.4 Normative Requirement Template

```markdown
### REQ-### — Requirement Title

- **Priority:** Must | Should | May
- **Applies to:** [SCOPE]
- **Implementation phase:** PHASE-##
- **Source decisions:** DEC-###, REC-###, or Chief Architect judgment
- **Verification:** Test, command, inspection, generated fixture, experiment, or
  dogfooding evidence
- **Risk linkage:** RSK-### or None

#### Requirement

Use precise normative language: MUST, MUST NOT, SHOULD, SHOULD NOT, or MAY.

#### Rationale

#### Acceptance Evidence

#### Exceptions
```

### 25.5 Proposed Specification Status

The first synthesis output should normally use:

```text
Proposed — pending adversarial review
```

---

## 26. Definitive Specification Contract

The exact structure adapts to the project, but a software-first specification
should cover:

- Artifact metadata.
- Executive decision summary.
- Authority and intended use.
- Problem and product definition.
- Goals and non-goals.
- Locked decisions and invariants.
- Final technology stack.
- System context.
- Architecture.
- Components and boundaries.
- Data model.
- Interfaces and integrations.
- User workflows.
- Security and privacy.
- Reliability and operations.
- Testing and verification.
- CI and release.
- Migration where applicable.
- Performance expectations.
- Internal contracts.
- Dependency bill of materials.
- Normative requirements.
- Traceability.
- Risk register.
- Open questions.
- Deferred work.
- Rejected work.
- Definition of done.
- Handoff to adversarial review.

The specification must be standalone and implementation-ready.

---

## 27. Specification Adversarial Review

### 27.1 Purpose

The review asks:

> What is wrong, contradictory, unsafe, non-total, under-specified,
> over-engineered, unprovable, difficult to implement, difficult to test, or
> likely to fail in real use?

### 27.2 Reviewer Posture

The reviewer must:

- Attack polished sections.
- Trace workflows end to end.
- Check failure, cancellation, rollback, and cleanup.
- Check consistency across prose, requirements, tables, examples, and
  appendices.
- Attempt to delete unnecessary machinery.
- Avoid adding features.
- Avoid treating preference as defect.
- Produce a small number of strong findings rather than meeting a quota.

### 27.3 Severity

- **Critical:** blocks all implementation or risks catastrophic harm.
- **High:** blocks the affected phase or creates major invalid behavior.
- **Medium:** must be fixed before the affected phase completes.
- **Low:** should be corrected in revision but does not block early work.

### 27.4 Finding Template

```markdown
## FND-### — Finding Title

- **Severity:** Critical | High | Medium | Low
- **Confidence:** High | Medium | Low
- **Category:** [CATEGORY]
- **Affected sections:**
- **Affected requirements:** REQ-### or None
- **Affected phases:** PHASE-## or None
- **Blocks implementation:** Entire program | Named phase | No

### Problem

### Evidence

### Failure Scenario

### Impact

### Root Cause

### Required Correction

### Proposed Specification Diff

### Acceptance Evidence

### Alternatives Considered

### Residual Risk

### Related Findings
```

### 27.5 Review Scope

A software-first review should inspect:

- Product scope.
- Requirements.
- Architecture.
- Data and interfaces.
- User workflows.
- Security.
- Filesystem behavior.
- Determinism.
- Dependency behavior.
- Testing.
- CI.
- Operations.
- Migration.
- Implementation phases.
- Agent legibility.
- Framework creep.
- Acceptance criteria.

---

## 28. Revised Definitive Specification

### 28.1 Finding Disposition

Every `FND-###` receives exactly one disposition:

- Accepted.
- Accepted with modification.
- Rejected.
- Deferred to a bounded evidence spike.
- Not applicable because another correction removes the cause.

No finding may disappear silently.

### 28.2 Revision Rules

The revision must:

- Integrate accepted corrections throughout all affected sections.
- Remove superseded or contradictory language.
- Reconcile overlapping diffs.
- Preserve important strengths.
- Prefer simplification over new machinery.
- Retain stable requirement identifiers where the subject remains the same.
- Add new requirements only from unused identifiers.
- Remain standalone.

### 28.3 Final Status

Use:

```text
Accepted — implementation authority
```

only when all Critical and implementation-blocking findings are resolved or
validly rejected and no known blocking contradiction remains.

Otherwise use:

```text
Proposed — implementation blocked
```

and identify blockers explicitly.

### 28.4 Required Revision-Specific Sections

- Revision Summary.
- Finding Disposition Ledger.
- Integrated Correction Ledger.
- Preserved Strengths.
- Updated implementation handoff.

---

## 29. Implementation Plan

### 29.1 Purpose

The implementation plan translates the accepted revised specification into a
safe delivery sequence.

It defines **how to build**, not what the architecture should become.

### 29.2 Plan Boundary

The plan must stop at:

- Phases.
- Milestones.
- Dependencies.
- Integration points.
- Evidence spikes.
- Dogfooding.
- Rollback and reconsideration boundaries.
- Executable phase acceptance criteria.

It must not create a granular execution backlog or coding-agent task packets.

### 29.3 Required Plan Content

- Artifact metadata.
- Implementation authority.
- Objectives.
- Non-goals.
- Assumptions.
- Dependency graph.
- Phase overview.
- One section per phase.
- Milestones.
- Cross-phase integration.
- Data or migration sequencing where applicable.
- Testing strategy by phase.
- Security activities by phase.
- Operations and release readiness.
- Dogfooding.
- Risk register.
- Open questions.
- Rollback and reconsideration triggers.
- Requirement-to-phase traceability.
- Definition of plan completion.

### 29.4 Phase Template

```markdown
## PHASE-## — Phase Name

- **Status:** Planned
- **Objective:**
- **User-visible outcome:**
- **Depends on:** PHASE-## or None
- **Requirements:** REQ-###
- **Milestones:** MS-###
- **Primary risks:** RSK-###

### Entry Criteria

### Scope

### Explicit Non-Goals

### Architecture and Components

### Integrations

### Data or Migration Work

### Evidence Spikes

### Testing and Verification

### Security and Reliability

### Dogfooding or Operational Validation

### Rollback and Reconsideration Triggers

### Exit Criteria

State observable, executable evidence.
```

### 29.5 Milestone Template

```markdown
### MS-### — Milestone Title

- **Phase:** PHASE-##
- **Outcome:**
- **Prerequisites:**
- **Acceptance evidence:**
- **Blocks:**
```

### 29.6 Sequencing Principles

The plan should:

- Resolve unknowns before they harden into architecture.
- Produce thin end-to-end capability early.
- Integrate continuously.
- Dogfood before broad feature expansion.
- Keep risky decisions reversible.
- Avoid phases that are merely horizontal infrastructure layers with no usable
  outcome.
- Avoid circular entry or exit criteria.
- Avoid acceptance criteria that depend on later phases.

---

## 30. Implementation-Plan Adversarial Review

The plan review must attack:

- Circular phase dependencies.
- Missing prerequisites.
- Overlarge phases.
- Milestones without integration evidence.
- Acceptance criteria that do not prove the claimed outcome.
- Late discovery of critical risks.
- Delayed dogfooding.
- Deferred integration.
- Migration gaps.
- Rollback gaps.
- Security work postponed too late.
- Test environments not available at the required phase.
- Phase boundaries that conflict with the specification.
- Plan steps that silently reinterpret architecture.
- Excessive parallel work.
- Implementation order that makes changes expensive before evidence arrives.

Use a separate `FND-###` range from the specification review.

The review must include concrete plan diffs and a clear implementation gate.

---

## 31. Final Revised Implementation Plan

The final revision must:

- Disposition every plan-review finding.
- Preserve the accepted revised specification.
- Integrate accepted corrections coherently.
- Remove circular sequencing.
- Make phase entry and exit criteria executable.
- Preserve early integration and dogfooding.
- State any remaining blockers honestly.
- Remain at phase and milestone granularity.

Recommended status:

```text
Accepted — delivery authority
```

only when no implementation-blocking plan finding remains.

---

## 32. Risk-Triggered Additional Review Rounds

A second review round is justified only when:

- The revision materially restructures architecture or sequencing.
- Several Critical or High findings were accepted.
- New machinery was added to resolve earlier findings.
- The reviewer identifies unresolved blockers.
- A spike changes a foundational decision.

The second review must focus on the machinery introduced or changed by the prior
revision rather than restarting the entire project from first principles.

---

## 33. Program Amendment Protocol

The approved Blueprint and graph may be amended, but never silently.

### 33.1 Material Amendment Triggers

- A missing research track is discovered.
- A stage proves redundant.
- A dependency is incorrectly ordered.
- New evidence invalidates a premise.
- A replication or spike materially changes scope.
- A legal, security, or operational constraint emerges.

### 33.2 Required Amendment Steps

1. Document the new evidence.
2. Propose the exact graph or scope change.
3. Analyze impact on existing reports, identifiers, prompts, dependencies, and
   downstream artifacts.
4. Obtain explicit human approval.
5. Create a `DEC-###` when the amendment is foundational or invalidates prior
   authority.
6. Update the Blueprint.
7. Update `research-program.toml`.
8. Mark affected stages `requires-revalidation`, `superseded`, or `blocked`.
9. Re-run only stages whose assumptions or inputs materially changed.
10. Preserve all original artifacts in Git history.

### 33.3 Non-Material Prompt Refinement

A just-in-time prompt may be refined without amending the program when the
refinement does not change:

- Stage objective.
- Scope.
- Authority.
- Dependencies.
- Expected output.
- Identifier range.
- Downstream role.

---

## 34. Commit Boundaries

Each accepted substantive artifact should normally receive its own coherent
commit.

Recommended patterns:

```text
docs: bootstrap research program
docs: add program blueprint
docs: add research charter
docs: add <track> research prompt
docs: add <track> research report
docs: reconcile replicated <track> research
docs: add definitive specification
docs: add specification adversarial review
docs: publish revised definitive specification
docs: add implementation plan
docs: add implementation plan adversarial review
docs: publish final revised implementation plan
```

Prompt installation and report execution should not be mixed in one commit
unless repository policy explicitly requires it.

---

## 35. Default Role Mapping

The workflow is tool-neutral at the role level but may use this practical
default mapping.

| Role                       | Responsibility                                             | Typical tool                          |
| -------------------------- | ---------------------------------------------------------- | ------------------------------------- |
| Research Program Architect | Discovery, Blueprint, Charter, graph, just-in-time prompts | Strong reasoning LLM                  |
| Deep Research Agent        | Current source-backed research and spikes                  | Deep research-capable LLM             |
| Repository Agent           | Install artifacts, inspect Git, validate paths and diffs   | Codex or repository-aware agent       |
| Replication Agent          | Independent run of identical prompt                        | Separate research-capable LLM session |
| Reconciliation Agent       | Compare replicated reports                                 | Fresh high-reasoning session          |
| Chief Architect            | Resolve research into specification                        | Fresh high-reasoning session          |
| Adversarial Reviewer       | Attack specification or plan                               | Separate fresh high-reasoning session |
| Revision Architect         | Disposition findings and rewrite artifact                  | Fresh high-reasoning session          |
| Validation Agent           | Independent artifact and repository validation             | Repository-aware agent                |

No role should depend on proprietary product behavior that is not represented in
the artifact contract.

---

## 36. Anti-Patterns

The workflow must actively prevent:

### 36.1 Chat-History Authority

Do not say “as discussed earlier” in a standalone artifact.

### 36.2 Research by Popularity

Do not select a tool or architecture because it is fashionable or frequently
mentioned.

### 36.3 Evidence-Free Confidence

Do not allow polished writing to substitute for citations, measurements, or
explicit judgment.

### 36.4 Broad Unbounded Prompts

Do not ask one research stage to investigate an entire project without a defined
boundary.

### 36.5 Premature Prompt Generation

Do not freeze dependent-stage prompts before their prerequisites exist.

### 36.6 Silent Recommendation Loss

Do not allow `REC-###` items to disappear during synthesis.

### 36.7 Mechanical Review Application

Do not treat every `FND-###` as automatically correct.

### 36.8 Review as Feature Ideation

Do not use adversarial review to add attractive subsystems.

### 36.9 Endless Review Loops

Do not run another round without an explicit trigger.

### 36.10 Prototype Capture

Do not turn spike code into production by inertia.

### 36.11 Placeholder Completion

Do not infer completion from an existing filename.

### 36.12 False Parallelism

Do not run dependent stages in parallel merely to save time.

### 36.13 Identifier Reuse

Do not reuse deleted or rejected identifiers.

### 36.14 Plan-as-Backlog

Do not create hundreds of speculative coding tasks in the final plan.

### 36.15 Implementation Before Authority

Do not begin substantive implementation while the revised specification or final
plan remains implementation-blocked.

---

## 37. Generic Adaptation Beyond Software

The governance spine remains valid outside software.

Adapt the focused research tracks and final specification vocabulary.

Examples:

- A legal project may replace system architecture with legal theory,
  jurisdictional analysis, procedural strategy, and evidentiary requirements.
- A market project may emphasize users, competitors, pricing, distribution, and
  unit economics.
- A scientific project may emphasize hypotheses, methods, datasets,
  reproducibility, and statistical validity.
- A theological project may emphasize primary texts, historical interpretation,
  language, doctrinal frameworks, and competing schools.
- An operational project may emphasize workflow, controls, staffing, failure
  modes, metrics, and governance.

The workflow must never pretend that engineering experiments replace
domain-qualified professional judgment where law, medicine, safety, or regulated
practice requires it.

---

## 38. Standard Bootstrap Task

The repository agent should receive a task with this structure:

```markdown
# Repository Bootstrap Task

Create the standard research-program repository for `[PROJECT NAME]`.

## Read or Confirm

- Current repository state.
- Existing user-authored content.
- Git status.

## Create

- `README.md`
- `AGENTS.md`
- `research-program.toml`
- `decisions/README.md`
- Standard `docs/` directories.
- Placeholder artifact files declared by the approved Blueprint.

## Rules

- Do not conduct substantive research.
- Do not invent project decisions beyond the approved discovery output.
- Do not overwrite substantive content.
- Use stable filenames.
- Initialize Git if necessary and authorized.
- Validate the complete tree.

## Report

- Files created.
- Conflicts.
- Validation results.
- Commit status.
```

---

## 39. Standard Just-in-Time Stage Package

For every substantive stage, the Research Program Architect must produce five
items:

1. **Canonical stage prompt.**
2. **Repository installation task** for placing the prompt at its reserved path.
3. **Attachment manifest.**
4. **Fresh-session launch message.**
5. **Post-stage validation task and recommended commit message.**

This package makes each transition reproducible and prevents missing context.

---

## 40. Standard Validation Task

```markdown
# Validate [ARTIFACT NAME]

Read:

1. `README.md`
2. `AGENTS.md`
3. `research-program.toml`
4. The Program Blueprint
5. The Research Charter
6. The commissioning prompt
7. Required upstream artifacts
8. The produced artifact

Validate:

- Required sections.
- Artifact metadata.
- Identifier uniqueness and ranges.
- Recommendation or finding disposition completeness.
- Citation portability.
- Evidence Ledger.
- Risks and open questions.
- Handoff Digest.
- Completion checklist truthfulness.
- Allowed file scope.
- Git diff.
- Manifest status transition.

Do not fabricate missing content.

Fix only mechanical defects. Report substantive defects and set the stage to
`requires-revision`.
```

---

## 41. Program Completion Criteria

The research program is complete only when:

1. The final revised definitive specification is accepted as implementation
   authority.
2. The final revised implementation plan is accepted as delivery authority.
3. Every research recommendation has a recorded disposition.
4. Every adversarial finding has a recorded disposition.
5. Every normative requirement is traceable to a source and implementation
   phase.
6. Every phase has executable entry and exit criteria.
7. Risks and open questions have owners, triggers, and deadlines.
8. No implementation-blocking finding remains.
9. No foundational choice is left to the implementer without an explicit bounded
   spike.
10. The manifest accurately reflects accepted commits and completed stages.
11. The repository working tree is clean.
12. The implementation handoff is complete.

---

## 42. Final Implementation Handoff

The final revised implementation plan must state:

- Which specification is authoritative.
- Which plan is authoritative.
- Whether implementation may begin.
- The first safe implementation phase.
- Required pre-implementation spikes.
- Phase dependencies.
- Earliest usable vertical slice.
- Earliest dogfooding milestone.
- Required validation commands or evidence at each boundary.
- Risks that must remain visible.
- Decisions that must remain reversible.
- Remaining blockers, if any.
- Which artifacts implementation agents must read.
- Which research reports are relevant only as supporting evidence.

---

## 43. Master Completion Checklist

A Research Program Architect should not declare the workflow configured until
all applicable items are true.

### Governance

- [ ] Git is the durable system of record.
- [ ] The repository follows the standard layout.
- [ ] `research-program.toml` exists and validates.
- [ ] Fresh sessions are mandatory.
- [ ] Human approval gates are defined.
- [ ] The rigor tier is approved.
- [ ] Authority and precedence are explicit.
- [ ] Identifier ranges are allocated.

### Discovery and Program Design

- [ ] The user was interviewed one question at a time.
- [ ] Recommendations accompanied clarification questions.
- [ ] The underlying goal was separated from the proposed solution.
- [ ] Material pivots received explicit approval.
- [ ] The Program Blueprint is accepted.
- [ ] The Research Charter is accepted.
- [ ] The adaptive research graph is explicit.
- [ ] Stage dependencies and parallel groups are explicit.
- [ ] Obvious omitted tracks are justified.

### Evidence

- [ ] Source hierarchy is defined.
- [ ] Current claims require verification.
- [ ] Portable citations are required.
- [ ] Evidence Ledgers are mandatory.
- [ ] Evidence spikes are first-class.
- [ ] Spike reports are bounded and reproducible.
- [ ] Weak and conflicting evidence must be disclosed.
- [ ] Optional replication is supported.
- [ ] Replication requires independent fresh sessions.
- [ ] Reconciliation preserves dissent and original reports.

### Stage Execution

- [ ] Prompts are generated just in time.
- [ ] Every stage receives an attachment manifest.
- [ ] Every substantive stage runs in a fresh session.
- [ ] Every artifact is standalone.
- [ ] Every report contains a Handoff Digest.
- [ ] Every stage has a completion checklist.
- [ ] Every artifact passes independent validation.
- [ ] Only validated committed artifacts unlock downstream stages.

### Synthesis and Review

- [ ] Every recommendation is dispositioned.
- [ ] One coherent specification is produced.
- [ ] Normative requirements use stable identifiers.
- [ ] The specification receives an adversarial review.
- [ ] Findings contain concrete failure scenarios and corrections.
- [ ] Every finding is dispositioned in the revision.
- [ ] The revised specification has an honest implementation status.
- [ ] Additional review rounds are risk-triggered only.

### Implementation Planning

- [ ] The plan is subordinate to the accepted specification.
- [ ] The plan is organized into phases and milestones.
- [ ] The plan does not become a task backlog.
- [ ] Every phase has entry and exit criteria.
- [ ] Evidence spikes occur before irreversible commitments.
- [ ] Integration and dogfooding begin early.
- [ ] The plan receives an adversarial review.
- [ ] Every plan-review finding is dispositioned.
- [ ] The final revised plan has an honest delivery status.

### Closure

- [ ] Requirement-to-phase traceability is complete.
- [ ] Risks and open questions have owners and deadlines.
- [ ] The program manifest records accepted commit hashes.
- [ ] The working tree is clean.
- [ ] The final implementation handoff is actionable.
- [ ] No substantive implementation has begun without accepted authority.

---

## 44. Final Instruction to the Research Program Architect

Begin with discovery.

Ask one question at a time and include your recommendation.

Challenge the premise constructively.

Do not overbuild the research program.

Select the rigor tier and adaptive research graph that fit the actual project.

Use Git as the durable system of record.

Generate prompts just in time.

Use fresh sessions for every substantive stage.

Ground claims in current sources, measured data, and bounded evidence spikes.

Support independent replication when risk justifies it.

Validate every artifact before acceptance.

Disposition every recommendation and finding.

Synthesize one coherent specification.

Attack it adversarially.

Revise it honestly.

Create an implementation plan organized into phases and milestones.

Attack the plan adversarially.

Revise it into the final delivery authority.

Stop before granular task decomposition.

The quality of the program is measured not by how much text it produces, but by
how reliably it converts uncertainty into traceable evidence, coherent
decisions, executable phases, and an honest implementation gate.
