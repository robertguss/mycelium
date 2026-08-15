# Default Role Mapping

Tool-neutral at the role level. Practical default:

| Role                       | Responsibility                                           | Typical tool                          |
| -------------------------- | -------------------------------------------------------- | ------------------------------------- |
| Research Program Architect | Discovery, Blueprint, Charter, graph, JIT prompts        | Strong reasoning LLM                  |
| Deep Research Agent        | Current source-backed research and spikes                | Deep research-capable LLM             |
| Repository Agent           | Install artifacts, inspect Git, validate paths and diffs | Repository-aware agent                |
| Replication Agent          | Independent run of identical prompt                      | Separate research-capable session     |
| Reconciliation Agent       | Compare replicated reports                               | Fresh high-reasoning session          |
| Chief Architect            | Resolve research into specification                      | Fresh high-reasoning session          |
| Adversarial Reviewer       | Attack specification or plan                             | Separate fresh high-reasoning session |
| Revision Architect         | Disposition findings and rewrite artifact                | Fresh high-reasoning session          |
| Validation Agent           | Independent artifact and repository validation           | Repository-aware agent                |

No role should depend on proprietary product behavior not represented in the
artifact contract.
