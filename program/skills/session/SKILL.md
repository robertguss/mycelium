---
name: session
description: >
  Choose the Mycelium ritual from the human's goal: new spark, work a live
  idea, park, wake, research program, second opinion, council, clarify, or
  handoff. Use at the start of a sitting in a mycelium.toml repo, or when
  the human states a goal but not a command.
---

# Session

Read `index.md` and `mycelium.toml` (not the whole tree). Pick **one**
ritual. Follow that skill to its done bar.

| Goal | Skill |
| --- | --- |
| New idea / first thought / “capture the spark” | `spark` |
| Decide, spar, work a question | `thinking` |
| Park / come back later | `simmer` |
| Return after a gap / “wake this” | `wake` |
| Destination reached / idea is decided | `clarify` |
| Ready to implement / write the packet | `handoff` |
| Survey many ideas | `portfolio` |
| Multi-day research program | `thinking` + `program/operator/getting-started.md` |
| One other model | `second-opinion` |
| Multi-model replication | `council` |

`thinking` also applies in every non-archived working session.

If `state` is `simmering` and the human is returning, choose `wake` even
if they said “let’s keep going.”

Done: the named skill has been invoked and that skill’s done bar is met.
Do not `git commit` unless the human asks.
