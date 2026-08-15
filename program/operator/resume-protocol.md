# Resume Protocol

When asked to resume the program, the Research Program Architect or repository
agent must:

1. Verify the working tree state.
2. Read `mycelium.toml` and run `mycelium status` (or `mycelium check`).
3. Read `README.md`, `AGENTS.md`, and any accepted Blueprint / Charter.
4. Confirm every stage marked `accepted` has a valid artifact and accepting
   commit (when commits are used).
5. Detect placeholders incorrectly marked complete.
6. Detect missing outputs.
7. Detect invalid status transitions.
8. Identify all currently eligible stages.
9. Respect parallel dependencies.
10. Recommend the next legal stage.
11. Generate the just-in-time package for that stage: prompt, installation task,
    attachment manifest, launch message, validation task, recommended commit
    message.

**Do not infer completion from chat history.**

Helpers: `mycelium status`, `mycelium check`.

## Stage statuses

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

## Legal transitions

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

## Acceptance rule

A stage becomes `accepted` only when:

1. Required artifact exists at the declared path.
2. Artifact metadata is complete.
3. Independent validation gate passes.
4. Artifact is committed (human/git workflow).
5. Manifest / log records the accepting commit when used.
6. Required human approval has been obtained.

## Unlock rule

A downstream stage is eligible only when every declared prerequisite is
`accepted`. Statuses `requires-revalidation`, `blocked`, or `superseded` do not
satisfy a dependency.
