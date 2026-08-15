# DEC-010 — Mycelium ships as a single-binary CLI (Go), not a GitHub template

- **Status:** Accepted
- **Date:** 2026-08-14
- **Owner:** Robert Guss
- **Supersedes:** DEC-004 (distribution mechanism only; its self-containment,
  version-pinning, and multi-runtime requirements are reaffirmed)
- **Related recommendations:** None
- **Related evidence:** Rails (`rails new` / `rails generate` is a CLI);
  cookiecutter (scaffold-only counterexample); Go `embed` package

## Context

The template-repo distribution model forced ceremony that existed only to
compensate for GitHub's indiscriminate template copying: the `framework/`
strip rule, init-time deletion, full placeholder trees at spark, and a
python3 + just dependency floor. Robert proposed replacing the template with
a scaffolding CLI and noted he builds his CLI tools in Go and wants a single
binary.

## Decision

1. Mycelium is a CLI (`mycelium`), written in Go, distributed as a single
   static binary. The master repository is its source repo, evolved in place
   per DEC-001.
2. Methodology content (contracts, templates + schemas, tiers, packs,
   skills, instance skeleton) lives as normal browsable files under
   `program/` in the master, is compiled into the binary via `go:embed`, and
   is emitted into every scaffolded instance.
3. Core commands: `mycelium new idea` (scaffold + git init + `gh repo
   create` + `idea` topic), `mycelium new <type>` (artifact generation),
   `mycelium check` (conformance), `mycelium status --all` (portfolio),
   `mycelium supersede` (state transitions). The CLI never commits work
   product; humans own git for thinking artifacts.
4. Two safety rules bound the design. The CLI creates, never thinks:
   instances remain self-contained thinking environments whose contracts,
   templates, and skills are present as files, so an agent with no binary
   can still operate manually. And at runtime, instance files are truth:
   `check` and `status` validate against the instance's own emitted schemas,
   never against the binary's embedded copies.
5. The `framework/` strip rule is obsolete: `framework/` is simply never
   emitted. Scaffolding is tier-aware: a spark receives only spark
   structure; deeper trees are emitted when a tier demands them.

## Rationale

The CLI chooses what to emit, so every compensation for template copying is
deleted rather than maintained. Dependencies drop from python3 + just +
template flow to one binary, strengthening multi-runtime portability
(DEC-004's own goal). Generator and checker become one tool reading one
schema, closing the drift gap by construction. DEC-003's portfolio and topic
hygiene become automated instead of conventional. Go specifically: single
static cross-compiled binary, `go:embed` is purpose-built for shipping the
methodology inside the tool while keeping it browsable in the repo, and it
is Robert's native CLI stack. Rust offers no advantage here without
fluency; Python breaks single-binary; cookiecutter is scaffold-only and
drags in Jinja, which DEC-005's no-templating-engine stance already
rejected.

## Consequences

The repo ceases to be a "Use this template" repo. PHASE-01 becomes CLI
foundation work; release engineering (tagged releases, install docs) enters
scope. Instances no longer carry a Justfile or scripts; the emitted
contracts plus the CLI are the operational surface, and a shipped skill
teaches agents both the commands and the manual floor.

## Alternatives Considered

Remaining a GitHub template (keeps strip ceremony and dependency floor).
Wrapping cookiecutter (scaffold-only; generation and checking would be
homeless). Python CLI via pipx (not a single binary; dependency management
on bare VMs).

## Risks

Binary availability on air-gapped or fresh machines. Mitigation: the manual
floor plus one-line install; CI installs are a single curl. Version skew
between binary and old instances. Mitigation: the instance-files-are-truth
rule.

## Revisit Triggers

If the manual floor proves too weak in practice (agents without the binary
producing nonconforming instances at a rate conformance cannot catch),
revisit emitting a minimal script shim.

## Approval

Approved by Robert Guss, 2026-08-14 ("this is a cli tool similar to rails...
I am totally ok with that... I want a single binary ideally"). Agent
concurred with recorded receipts rather than deference.
