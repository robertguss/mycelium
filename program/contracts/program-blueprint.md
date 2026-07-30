# Program Blueprint Contract

The Program Blueprint governs the project and research program. It defines what
the program will do. It must **not** conduct the substantive research itself.

## Required sections

1. Artifact metadata
2. Product or project vision
3. Problem statement
4. Intended users and stakeholders
5. Goals
6. Non-goals
7. Locked constraints
8. Success criteria
9. Rigor tier
10. Research graph
11. Stage descriptions and dependencies
12. Parallelism
13. Optional replication points
14. Artifact inventory
15. Identifier allocations
16. Authority and precedence
17. Human approval gates
18. Fresh-session policy
19. Validation and commit gates
20. Amendment protocol
21. Completion criteria
22. Implementation handoff expectation

## Research graph requirements

For every stage, define:

- Stable stage ID
- Stage name
- Stage kind (`foundational` | `independent` | `dependent` | `replication` |
  `reconciliation` — plus fixed spine kinds in the manifest)
- Primary research question
- Scope and non-goals
- Prerequisites and inputs
- Required output path
- Identifier ranges
- Whether evidence spikes are expected
- Whether replication is permitted or recommended
- Parallel group, if any
- Downstream consumers
- Completion criteria

## Stage kinds (focused research)

- **Foundational:** must complete before dependent research begins.
- **Independent:** may run in parallel with other independent stages.
- **Dependent:** consumes one or more completed reports.
- **Replication:** repeats an existing prompt independently.
- **Reconciliation:** compares replicated reports.

## Parallelism

Sequential execution is the default. Parallel execution is permitted only when
stages do not require one another’s findings. The Blueprint must explicitly
justify parallelism.

## Stage selection

For every selected track, explain why it exists, why another track cannot absorb
it, which decision it informs, and which artifact consumes it. For every obvious
but omitted track, briefly explain why it is unnecessary.

## Just-in-time prompts

The Blueprint defines the graph up front; detailed stage prompts are generated
just in time. Independent prompts may be generated together after Blueprint and
Charter approval. Dependent, reconciliation, synthesis, review, revision, and
implementation-planning prompts are generated only after prerequisites are
accepted.
