# Research Charter Contract

The Research Charter defines the evidence and decision methodology inherited by
all later stages.

## Required sections

- Research philosophy
- Scope discipline
- Source hierarchy
- Citation rules
- Current-information rules
- Evidence-spike protocol
- Evidence Ledger format
- Recommendation format
- Evaluation rubric
- Confidence model
- Risk and open-question format
- Replication and reconciliation protocol
- Synthesis rules
- Adversarial-review rules
- Validation rules
- Handoff rules
- Anti-patterns
- Completion standards

## Default source hierarchy

1. Official specifications, standards, primary documentation, source
   repositories, and first-party release information.
2. Peer-reviewed research, authoritative institutional publications,
   maintainer-authored design records, and official security advisories.
3. High-quality independent technical analysis, production case studies, and
   reproducible benchmarks.
4. Community reports, issue discussions, forum posts, and practitioner
   anecdotes.
5. Vendor marketing and unsourced summaries.

Lower-tier evidence may reveal failure modes but should not carry a load-bearing
recommendation alone when stronger evidence is available.

## Current verification

Any claim that may have changed must be verified as of the actual research date
(tools, versions, pricing, APIs, laws, platform behavior, compatibility,
advisories, deployment features, licensing, etc.).

## Portable citations

Prefer Markdown links, numbered footnotes, and source ledgers with URLs and
access dates. Do not rely solely on ephemeral UI citation tokens.
