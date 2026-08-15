# PHASE-02 Acceptance Matrix

In-repo checklist copied from
[`PHASE-02-implementation-brief.md`](PHASE-02-implementation-brief.md) §15.
The brief is binding.

PHASE-02 is accepted when MS-201 is green in `go test ./...` on `main`.
GitHub Actions is not a gate. Arvo accepts the phase. The 7-real-day dogfood
wake is human evidence, not the gate.

| id | check | evidence | owner |
| --- | --- | --- | --- |
| A-S0 | Slice 0 files exist and match this brief | PR diff; Quality read | Engineering |
| A-S1 | revisit / due / overdue / trigger / table tests green | `go test` | Engineering |
| A-S2 | new scaffold emits `index.md`; `mycelium index` repairs; check binds | hermetic CLI | Engineering |
| A-S3 | `state` / `wake` edges; brief written; clarified legal; handed-off fails | hermetic CLI | Engineering |
| A-S4 | single `status` + due/overdue | hermetic CLI | Engineering |
| A-S5 | `status --all --offline` two-instance root; partial line; no gh | hermetic CLI | Engineering |
| A-S6 | merge flags + unread + archived filter via fake runner | `go test` | Engineering |
| A-S7 | three skills emitted on new scaffold; no retrofit on state/index | hermetic CLI | Engineering |
| A-S8 | MS-201 fixture green | `go test ./...` | Engineering |
| MS-201 | all §13 expected bullets | A-S8 (uses 1–4) | Engineering |
| DOGFOOD-7d | one real 7-day wake | Arvo / Robert; **not the gate** | Arvo |
