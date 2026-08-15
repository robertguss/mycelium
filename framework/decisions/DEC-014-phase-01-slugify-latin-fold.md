# DEC-014 — PHASE-01 slugify is a latin/compatibility fold, not NFKD

- **Status:** Accepted
- **Date:** 2026-08-14
- **Owner:** Robert Guss
- **Supersedes:** None (narrows PHASE-01 implementation brief §10 slugify)
- **Related recommendations:** None
- **Related evidence:** Quality probe on PR #8 (ﬁle, fullwidth Test, Đà Nẵng)

## Context

PHASE-01 implementation brief §10 originally named Unicode NFKD for
`slugify`. The PHASE-01 dependency floor allows only
`github.com/pelletier/go-toml/v2`. `internal/slug` therefore ships a fixed
`latinFold` map plus ASCII `[a-zA-Z0-9]` keep and space/`_` → `-`, without
`golang.org/x/text` and without an in-package NFKD.

Quality probed three outcomes that differ from full NFKD and filed a
NEEDS-WORK asking either for NFKD or an accepted DEC documenting the fold.

## Decision

1. PHASE-01 slugify is the existing documented latin/compatibility fold in
   `internal/slug` (`latinFold` + ASCII `[a-zA-Z0-9]` + space/`_` → `-`).
2. Full Unicode NFKD (brief §10 original sentence) is deferred until a real
   instance needs it.
3. Do not grow the `latinFold` map in this phase; unlisted runes are dropped.
4. Known PHASE-01 outcomes (specified behavior, not bugs):
   - ﬁ (U+FB01) is dropped (`ﬁle` → `le`)
   - fullwidth Latin is dropped (`Ｔｅｓｔ` → empty / `ErrEmpty`)
   - `Đà Nẵng` → `a-nng` (`Đ` and `ẵ` unlisted)
5. Revisit when a real idea name needs NFKD.

## Rationale

Growing the fold map or adding `x/text` for hypothetical names is scope
creep against the dependency floor. Locking the probed outcomes in tests
makes the contract honest without claiming NFKD.

## Consequences

Brief §10 and `internal/slug` comments cite this DEC. Callers that need
true NFKD must wait for a later DEC or bring a name that the fold already
handles.

## Alternatives Considered

Implement in-package NFKD (rejected: large, easy to get wrong).
Add `golang.org/x/text` (rejected: breaks dependency floor).
Grow `latinFold` for the probed runes (rejected: unbounded map growth).

## Risks

A real idea name with unlisted letters yields a surprising or empty slug.
Mitigation: refuse empty; revisit when that name appears.

## Revisit Triggers

A real idea name needs NFKD (or a broader fold) to produce a usable slug.

## Approval

Settled by Arvo 2026-08-14 (Quality DEC option on PR #8).
