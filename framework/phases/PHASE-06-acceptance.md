# PHASE-06 Acceptance Matrix

In-repo checklist copied from
[`PHASE-06-implementation-brief.md`](PHASE-06-implementation-brief.md) §14.
The brief is binding.

PHASE-06 is accepted when MS-601 is green in `go test ./...` on `main`.
GitHub Actions is not a gate. An isolated-agent run and a real-project
handoff are not the gate. Arvo accepts the phase.

| id | check | evidence | owner |
| --- | --- | --- | --- |
| A-S0 | Slice 0 files exist and match this brief | PR diff; Quality read | Engineering |
| A-S1 | parser / structure checker; no CLI verb yet; no state lift | hermetic `go test` | Engineering |
| A-S2 | `mycelium handoff` + `state handed-off` IFF + item 8 + item 24 + storage lift | hermetic `go test` | Engineering |
| A-S3 | AGENTS.md Implementation systems + mycelium-cli + reference | file read + scaffold emit | Engineering |
| A-S4 | fixture copies + executable acceptance + golden impl | hermetic `go test` | Engineering |
| A-S5 | MS-601 matrix green + 85% coverage | `go test ./...` on Quality's computer | Engineering |
| MS-601-GEN | structurally complete `handoff/` | A-S5 / A-S4 | Engineering |
| MS-601-CMD | `handoff` from `clarified` → `handed-off` + log op | A-S5 / A-S2 | Engineering |
| MS-601-REFUSE | `state handed-off` without packet refuses | A-S5 / A-S2 | Engineering |
| MS-601-CHECK | stored `handed-off` IFF packet | A-S5 / A-S2 | Engineering |
| MS-601-GOLD | acceptance tests + golden impl | A-S5 / A-S4 | Engineering |
| MS-601-COV | product + new packages ≥85% | A-S5 | Engineering |

No isolated-agent-gate row. No dogfood-gate row. No Actions-job row.
Quality should refuse a PR that adds an Actions job as the MS-601 gate.
Quality should not refuse a missing PHASE-06 workflow.
