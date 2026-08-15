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
	write(t, filepath.Join(root, "evidence", "EVD-001-vendor.md"), `+++
id = "EVD-001"
title = "Vendor"
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

2026-08-06
`)
	write(t, filepath.Join(root, "assumptions", "ASM-001-api.md"), `+++
id = "ASM-001"
title = "API"
status = "Held"
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

2026-08-05
`)
	write(t, filepath.Join(root, "assumptions", "ASM-002-budget.md"), `+++
id = "ASM-002"
title = "Budget"
status = "Retired"
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

none yet
`)

	mb, err := os.ReadFile(filepath.Join(root, "mycelium.toml"))
	if err != nil {
		t.Fatal(err)
	}
	m, err := manifest.Parse(mb)
	if err != nil {
		t.Fatal(err)
	}
	logBytes, err := os.ReadFile(filepath.Join(root, "log.md"))
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC)
	brief, err := wakebrief.Collect(root, m, logBytes, now)
	if err != nil {
		t.Fatal(err)
	}
	if brief.SimmerDate != "2026-08-01" {
		t.Fatalf("simmer date=%q", brief.SimmerDate)
	}
	if len(brief.Evidence) != 1 || brief.Evidence[0].ID != "EVD-001" {
		t.Fatalf("evidence=%v", brief.Evidence)
	}
	if len(brief.Assumptions) != 1 || brief.Assumptions[0].ID != "ASM-001" {
		t.Fatalf("assumptions=%v", brief.Assumptions)
	}
	out := string(wakebrief.Render(brief))
	if strings.Contains(out, "ASM-002") {
		t.Fatal("ASM-002 must be absent")
	}
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
