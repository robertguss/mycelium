# PHASE-04 Acceptance Matrix

In-repo checklist copied from
[`PHASE-04-implementation-brief.md`](PHASE-04-implementation-brief.md) §15.
The brief is binding.

PHASE-04 is accepted when MS-401 is green in `go test ./...` on `main`.
GitHub Actions is not a gate. Arvo accepts the phase. A live Cursor council
run is human evidence, not the gate.

| id | check | evidence | owner |
| --- | --- | --- | --- |
| A-S0 | Slice 0 files exist and match this brief | PR diff; Quality read | Engineering |
| A-S1 | pack presence; collision; enable/disable; `reviews/` extra-top-level | hermetic `go test` | Engineering |
| A-S2 | pack schemas; `mycelium new` discovers pack types; unknown when pack absent | hermetic `go test` + file read | Engineering |
| A-S3 | provenance / hash / cardinality / `opt_in` / `cost_class` / `SEED-DISSENT` | hermetic `go test` | Engineering |
| A-S4 | pack skills + adapters + capability note; mycelium-cli names the surface; no `mycelium council` verb to run | hermetic `go test` + file read | Engineering |
| A-S5 | MS-401 matrix green | `go test ./...` | Engineering |
| MS-401-SO | second-opinion row | A-S5 | Engineering |
| MS-401-CUR | Cursor-council row (seeded files, `adapter=cursor`; no live Cursor) | A-S5 | Engineering |
| MS-401-MAN | manual-floor row | A-S5 | Engineering |
| MS-401-OPT | explicit opt-in and stated cost class (happy + negatives) | A-S5 | Engineering |
| MS-401-HASH | prompt identity and model provenance evidenced | A-S5 | Engineering |
| MS-401-IND | independent per-model reports (two RPT files) | A-S5 | Engineering |
| MS-401-DIS | seeded dissent surviving reconciliation | A-S5 | Engineering |
| MS-401-PKG | council-pack enable/disable without touching core checks | A-S5 / A-S1 | Engineering |

No live-Cursor-gate row. No Actions-job row. Quality should refuse a PR that
adds an Actions job as the MS-401 gate. Do not cite OQ-006.
