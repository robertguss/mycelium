# PHASE-01 Acceptance Matrix

- **Status:** Commissioning matrix (Slice 0)
- **Date:** 2026-08-14
- **Brief:** [PHASE-01-implementation-brief.md](PHASE-01-implementation-brief.md)

This PR aims at **A-S0** and **A-S1**. Later rows stay unchecked until their
slices land.

| id | check | evidence | owner | this PR |
| --- | --- | --- | --- | --- |
| A-S0 | Slice 0 files exist and match the brief | PR diff; Quality read | Engineering / Quality | **aim** |
| A-S1 | go test ./... and go build and version | CI + local hermetic | Engineering / CI | **aim** |
| A-S2 | 11 templates + 11 schemas + 3 tiers + skeleton + skill | tree exists | Engineering | unchecked |
| A-S3 | table-driven pure tests green | go test | Engineering / CI | unchecked |
| A-S4 | crash/resume/stale-lock tests green | go test | Engineering / CI | unchecked |
| A-S5 | new idea --offline spark; forbidden paths absent | hermetic CLI | Engineering / CI | unchecked |
| A-S6 | check 0 on spark; teaching errors for illegal state, ID mismatch, journal, undeclared extra | hermetic CLI | Engineering / CI | unchecked |
| A-S7 | 11 types generated; overwrite/out-of-range refuse; log line; tokens | hermetic CLI | Engineering / CI | unchecked |
| A-S8 | tier up/down/idempotent; no deletion | hermetic CLI | Engineering / CI | unchecked |
| A-S9 | hermetic workflow green on the PR | GitHub Actions | CI | unchecked |
| A-S10 | credentialed job passed once; hermetic does not depend on it | Actions log | CI / Arvo | unchecked |
| MS-101a | all (a) bullets in brief §14 | A-S1..A-S9 | Arvo accepts | unchecked |
| MS-101b | all (b) bullets in brief §14 | A-S10 | Arvo accepts | unchecked |
| SLO-5m | spark-to-first-thought under five minutes | user SLO; not the phase gate | Robert | unchecked |
