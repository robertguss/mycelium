# DEC-011 — Migration machinery is deferred; instance files are runtime truth

- **Status:** Accepted
- **Date:** 2026-08-14
- **Owner:** Robert Guss
- **Supersedes:** DEC-005 (migration clause only) and DEC-004's
  "explicit migrations" upgrade path
- **Related recommendations:** None
- **Related evidence:** None

## Context

DEC-005 adopted Rails-style migrations (`program/migrations/`,
`applied_migrations`, `just upgrade`) for evolving live instances. Robert
challenged the premise: his research projects are finite; once the research
is done he does not return to upgrade them, so retroactive methodology
migration solves a problem he does not have.

## Decision

1. No migration machinery is built. `program/migrations/`, `just upgrade`,
   and `applied_migrations` are removed from scope.
2. Two cheap seams are kept so migrations can be added later without
   archaeology: every scaffolded manifest records the generating CLI
   version, and the CLI's runtime commands (`check`, `status`) validate
   against the instance's own emitted schemas rather than the binary's
   embedded copies. Old instances therefore keep working under new binaries
   by construction, which is precisely why migrations stay unnecessary.
3. `CHANGELOG.md` discipline in the master remains, as each release's
   human-readable face.

## Rationale

YAGNI, applied honestly. Ideas that simmer across a CLI version bump wake
into a binary that reads their own files as truth; finished ideas are
archives, and archives are not upgraded. Version pinning (DEC-004,
reaffirmed by DEC-010) already made "instance frozen at its methodology
version" the epistemically correct default.

## Consequences

Methodology improvements reach only new instances. A long-simmering idea
that wants a new capability adopts it manually (copy the new template in)
or is re-scaffolded, both acceptable at current scale. PHASE-05 loses its
migration content and becomes distribution and lifecycle commands.

## Alternatives Considered

Full Rails-style migrations (built for a future that may not arrive;
rejected as speculative weight). No seams at all (saves nothing measurable
and forfeits the cheap insurance of a version stamp).

## Risks

A future manifest-format change could break `mycelium status --all` when it
scans old instances. Mitigation: treat the manifest format as append-only
and keep the portfolio scanner tolerant of older shapes.

## Revisit Triggers

The recorded crux: the day a manifest or schema format change would break
portfolio scanning or conformance of existing instances, migrations return
to scope. Also revisit if Robert finds himself manually upgrading more than
two live instances in a quarter.

## Approval

Approved by Robert Guss, 2026-08-14 ("once the research is done, that's it...
lets defer that for now"). Agent concurred; the instance-files-are-truth
rule makes the deferral architecturally safe rather than hopeful.
