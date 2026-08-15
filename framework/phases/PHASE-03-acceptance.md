# PHASE-03 Acceptance Matrix

In-repo checklist copied from
[`PHASE-03-implementation-brief.md`](PHASE-03-implementation-brief.md) §15.
The brief is binding.

PHASE-03 is accepted when MS-301 is green in `go test ./...` on `main`.
GitHub Actions is not a gate. Arvo accepts the phase. Two hermetic fixture
sessions (one disputed, one aligned) are the gate. Substantive-vs-bare
question judgment is human / adversarial review, never an automated content
score.

| id | check | evidence | owner |
| --- | --- | --- | --- |
| A-S0 | Slice 0 files exist and match this brief | PR diff; Quality read | Engineering |
| A-S1 | schema drops always-required Crux; parser table tests green; new question template omits Crux/Reasons | `go test` + file read | Engineering |
| A-S2 | IFF check rules; disputed fail/pass; aligned pass; invalid fail; fill-Positions pass; spark-zero-OQ pass | hermetic `go test` | Engineering |
| A-S3 | CONTEXT.md glossary rules | hermetic `go test` | Engineering |
| A-S4 | optional Dissent; existing DECs green | hermetic `go test` | Engineering |
| A-S5 | thinking skill emitted on new scaffold; no retrofit; credits present; mycelium-cli + AGENTS.md name it | hermetic `go test` + file read | Engineering |
| A-S6 | MS-301 two fixtures green | `go test ./...` | Engineering |
| MS-301 | all §13 expected bullets | A-S6 (uses 1–2) | Engineering |

No DOGFOOD row. No Actions-job row. Quality should refuse a PR that adds an
Actions job as the MS-301 gate.
