# PHASE-05 Acceptance Matrix

In-repo checklist copied from
[`PHASE-05-implementation-brief.md`](PHASE-05-implementation-brief.md) §14.
The brief is binding.

PHASE-05 is accepted when MS-501 is green in `go test ./...` on `main`.
GitHub Actions is not a gate. Install SLO and a live GitHub Release are not
the gate. Arvo accepts the phase.

| id | check | evidence | owner |
| --- | --- | --- | --- |
| A-S0 | Slice 0 files exist and match this brief | PR diff; Quality read | Engineering |
| A-S1 | parsers / eligibility / schema deltas; no CLI verb yet | hermetic `go test` | Engineering |
| A-S2 | `mycelium supersede` + item 8 + item 23 + refuse table | hermetic `go test` | Engineering |
| A-S3 | G0–G3 + `partial: legacy-manifest` + no scan abort | hermetic `go test` | Engineering |
| A-S4 | CHANGELOG heading + release refuse/checksum + install doc | hermetic `go test` + file read | Engineering |
| A-S5 | MS-501 matrix green | `go test ./...` on Quality's computer | Engineering |
| MS-501-SUP | bidirectional cross-links + log entry | A-S5 / A-S2 | Engineering |
| MS-501-G0 | current manifest check+status | A-S5 / A-S3 | Engineering |
| MS-501-G1 | missing `github_repo` + frozen contract | A-S5 / A-S3 | Engineering |
| MS-501-G2 | `research-program.toml` only; no crash | A-S5 / A-S3 | Engineering |
| MS-501-G3 | unknown key: status-only golden | A-S5 / A-S3 | Engineering |
| MS-501-REL | checksummed binaries + CHANGELOG heading | A-S5 / A-S4 | Engineering |

**G3 note:** MS-501-G3 / A-S3 treat G3 as a status-only golden — `mycelium status` exit 0 and `mycelium check` exit 1; do not assert check exit 0 on G3.

No Install-SLO-gate row. No GitHub-Release-gate row. No Actions-job row.
Quality should refuse a PR that adds an Actions job as the MS-501 gate.
Quality should not refuse a missing PHASE-05 workflow.
