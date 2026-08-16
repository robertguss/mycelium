# Resume Protocol

Open this idea in the agent. Say “resume” or “wake this.” Do not paste a
chat transcript and call it context.

## Any idea (default)

1. Run `mycelium status` (and `mycelium check` if the tree may be dirty).
2. Read `index.md`. If `state` is `simmering`, run `mycelium wake` and
   read `briefs/LATEST.md` **before** `log.md`.
3. Follow `.agents/skills/session/SKILL.md` from the human’s goal.
4. Recommend the next legal move. Do not infer completion from chat.

Helpers: `mycelium status`, `mycelium check`, `mycelium wake`.

## Long program (blueprint / charter / stages in play)

When this idea is running the governance spine, also:

1. Verify the working tree.
2. Read `mycelium.toml`, `README.md`, `AGENTS.md`, and any accepted
   Blueprint / Charter.
3. Confirm every stage marked `accepted` has a valid artifact and accepting
   commit (when commits are used).
4. Detect placeholders incorrectly marked complete.
5. Detect missing outputs and invalid status transitions.
6. Identify currently eligible stages. Respect parallel dependencies.
7. Recommend the next legal stage.
8. Generate the just-in-time package for that stage: prompt, installation
   task, attachment manifest, launch message, validation task, recommended
   commit message.

**Do not infer completion from chat history.**

## Stage statuses

`planned` | `prompt-ready` | `in-progress` | `awaiting-validation` |
`requires-revision` | `accepted` | `blocked` | `requires-revalidation` |
`superseded` | `cancelled`

```text
planned → prompt-ready → in-progress → awaiting-validation
  → accepted | requires-revision | blocked
requires-revision → in-progress
accepted → requires-revalidation | superseded
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
`accepted`. Statuses `requires-revalidation`, `blocked`, or `superseded`
do not satisfy a dependency.
