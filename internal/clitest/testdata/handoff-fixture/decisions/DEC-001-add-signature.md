+++
id = "DEC-001"
title = "Add signature"
status = "Accepted"
date = "2026-08-15"
owner = "TBD"
+++

# DEC-001 — Add signature

## Context

Bounded target needs a stable function signature for the handoff fixture.

## Decision

Implement `Add(a, b int) int` in `add.go`.

## Rationale

Smallest useful golden target for hermetic acceptance tests.

## Consequences

Acceptance tests call `Add` with the locked signature.

## Alternatives Considered

none

## Risks

none

## Revisit Triggers

none

## Approval

Accepted for the canonical handoff fixture.
