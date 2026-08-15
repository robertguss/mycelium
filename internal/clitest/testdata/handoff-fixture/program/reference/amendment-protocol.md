# Program Amendment Protocol

The approved Blueprint and graph may be amended, but never silently.

## Material amendment triggers

- Missing research track discovered
- Stage proves redundant
- Dependency incorrectly ordered
- New evidence invalidates a premise
- Replication or spike materially changes scope
- Legal, security, or operational constraint emerges

## Required amendment steps

1. Document the new evidence.
2. Propose the exact graph or scope change.
3. Analyze impact on existing reports, identifiers, prompts, dependencies, and
   downstream artifacts.
4. Obtain explicit human approval.
5. Create a `DEC-###` when foundational or invalidating prior authority.
6. Update the Blueprint.
7. Update `research-program.toml`.
8. Mark affected stages `requires-revalidation`, `superseded`, or `blocked`.
9. Re-run only stages whose assumptions or inputs materially changed.
10. Preserve all original artifacts in Git history.

## Non-material prompt refinement

A just-in-time prompt may be refined without amending the program when it does
not change: stage objective, scope, authority, dependencies, expected output,
identifier range, or downstream role.
