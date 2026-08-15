package wakebrief_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/robertguss/mycelium/internal/manifest"
	"github.com/robertguss/mycelium/internal/wakebrief"
)

func TestRenderRequiredH2s(t *testing.T) {
	b := wakebrief.Brief{
		WakeDate:   "2026-08-09",
		SimmerDate: "2026-08-01",
		Revisit:    "2026-08-08",
		SimmerLine: "2026-08-01\tstate\t-\tsimmering revisit=2026-08-08",
		LogSince:   []string{"2026-08-01\tstate\t-\tsimmering revisit=2026-08-08"},
		Evidence:   []wakebrief.EvidenceCite{{ID: "EVD-001", Date: "2026-08-06"}},
		Assumptions: []wakebrief.AssumptionCite{
			{ID: "ASM-001", Status: "Held", Date: "2026-08-05"},
		},
	}
	out := string(wakebrief.Render(b))
	for _, h2 := range wakebrief.RequiredH2s() {
		if !strings.Contains(out, "## "+h2) {
			t.Fatalf("missing H2 %q in:\n%s", h2, out)
		}
	}
	for _, want := range []string{"EVD-001", "ASM-001", "2026-08-01", "2026-08-08"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in:\n%s", want, out)
		}
	}
	if strings.Contains(out, "ASM-002") {
		t.Fatal("must not cite ASM-002")
	}
}

func TestCollectCitations(t *testing.T) {
	root := citationFixture(t)
	m := mustManifest(t, root)
	logBytes := mustRead(t, filepath.Join(root, "log.md"))
	now := time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC)
	brief, err := wakebrief.Collect(root, m, logBytes, now)
	if err != nil {
		t.Fatal(err)
	}
	if brief.SimmerDate != "2026-08-01" {
		t.Fatalf("simmer date=%q", brief.SimmerDate)
	}
	if len(brief.Evidence) != 1 || brief.Evidence[0].ID != "EVD-001" {
		t.Fatalf("evidence=%v want only EVD-001 (EVD-002 future excluded)", brief.Evidence)
	}
	ids := map[string]bool{}
	for _, a := range brief.Assumptions {
		ids[a.ID] = true
	}
	if !ids["ASM-001"] || !ids["ASM-003"] {
		t.Fatalf("assumptions=%v want ASM-001 (Held+date) and ASM-003 (Open)", brief.Assumptions)
	}
	if ids["ASM-002"] {
		t.Fatalf("ASM-002 Retired must be absent: %v", brief.Assumptions)
	}
	out := string(wakebrief.Render(brief))
	if strings.Contains(out, "ASM-002") || strings.Contains(out, "EVD-002") {
		t.Fatalf("render must omit ASM-002 and EVD-002:\n%s", out)
	}
}

func citationFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	write(t, filepath.Join(root, "mycelium.toml"), `schema_version = 1
idea_name = "Wake Fixture"
slug = "wake-fixture"
created_date = "2026-08-01"
updated_date = "2026-08-01"
state = "simmering"
tier = "focused"
methodology_version = "2.0.0"
generated_by_cli_version = "0.1.0-dev"
revisit = "2026-08-08"
github_repo = ""
`)
	write(t, filepath.Join(root, "log.md"), `# Log

2026-08-01	scaffold	-	wake-fixture
2026-08-01	state	-	exploring
2026-08-01	state	-	simmering revisit=2026-08-08
`)
	mustMkdir(t, filepath.Join(root, "evidence"))
	mustMkdir(t, filepath.Join(root, "assumptions"))
	write(t, filepath.Join(root, "evidence", "EVD-001-vendor.md"), artifactEVD("EVD-001", "Vendor", "2026-08-06"))
	write(t, filepath.Join(root, "evidence", "EVD-002-future.md"), artifactEVD("EVD-002", "Future", "2026-08-20"))
	write(t, filepath.Join(root, "assumptions", "ASM-001-api.md"), artifactASM("ASM-001", "API", "Held", "2026-08-05"))
	write(t, filepath.Join(root, "assumptions", "ASM-002-budget.md"), artifactASM("ASM-002", "Budget", "Retired", ""))
	write(t, filepath.Join(root, "assumptions", "ASM-003-open.md"), artifactASM("ASM-003", "Open item", "Open", ""))
	return root
}

func artifactEVD(id, title, trigger string) string {
	return `+++
id = "` + id + `"
title = "` + title + `"
status = "Recorded"
date = "2026-08-01"
source = "changelog"
+++

## Claim

x

## Source

y

## Observation

z

## Limitations

w

## Revalidation Trigger

` + trigger + `
`
}

func artifactASM(id, title, status, trigger string) string {
	body := "none yet"
	if trigger != "" {
		body = trigger
	}
	return `+++
id = "` + id + `"
title = "` + title + `"
status = "` + status + `"
date = "2026-08-01"
attached_to = "DEC-001"
+++

## Statement

x

## Falsifier

y

## Implications

z

## Revisit Triggers

` + body + `
`
}

func mustManifest(t *testing.T, root string) manifest.Manifest {
	t.Helper()
	mb := mustRead(t, filepath.Join(root, "mycelium.toml"))
	m, err := manifest.Parse(mb)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func write(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func mustMkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
}
