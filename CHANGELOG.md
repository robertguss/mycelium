# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- `docs/user-guide.md` — operator walkthrough (LLM-first, scenario-based).
- Session skills: `session`, `simmer`, `clarify`, `handoff`.

### Changed

- Product README and `docs/install.md` rewritten for new users (install,
  first idea, agent workflow). From-source install documented.
- Shipped `program/operator/` and `program/README.md` rewritten as the
  in-instance start path. Skills teach spark / work / park / wake /
  research / handoff.
- `thinking`, `council`, and `second-opinion` no longer claim
  `mycelium handoff` does not exist.
- Product `AGENTS.md` authority is shipped contracts + the user guide.

### Removed

- `framework/` (blueprint, DECs, phase briefs, reviews, prompts). Git is
  the archive. Teaching-error `contract:` paths now point at
  `program/contracts/`. `ForbiddenPaths` still lists `framework`.

## [0.1.0] - 2026-08-15

### Added

- PHASE-01 through PHASE-04 CLI surface.
- `mycelium supersede` (PHASE-05).
