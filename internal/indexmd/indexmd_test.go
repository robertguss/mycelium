package indexmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/indexmd"
	"github.com/robertguss/mycelium/internal/logfmt"
)

func TestRenderStructure(t *testing.T) {
	body := string(indexmd.Render(indexmd.Instance{
		IdeaName:   "Wake Fixture",
		State:      "spark",
		Tier:       "focused",
		Revisit:    "",
		GithubRepo: "",
		Counts:     indexmd.ZeroCounts(),
		LogLines:   []string{logfmt.Line("2026-08-01", "scaffold", "-", "Wake Fixture")},
		Wake:       "none",
	}))
	for _, h2 := range []string{"## State", "## Artifacts", "## Log tail", "## Wake"} {
		if !strings.Contains(body, h2+"\n") {
			t.Fatalf("missing %s in:\n%s", h2, body)
		}
	}
	if !strings.HasPrefix(body, "# Wake Fixture\n") {
		t.Fatalf("H1: %q", body[:40])
	}
	if !strings.Contains(body, "DEC: 0\n") || !strings.Contains(body, "MS: 0\n") {
		t.Fatalf("want zero counts for all NS:\n%s", body)
	}
	if !strings.Contains(body, "## Wake\n\nnone\n") {
		t.Fatalf("wake none:\n%s", body)
	}
}

func TestRenderLogTailCapsAt20(t *testing.T) {
	var lines []string
	for i := 0; i < 25; i++ {
		lines = append(lines, logfmt.Line("2026-08-01", "new", "DEC-001", "x"))
	}
	body := string(indexmd.Render(indexmd.Instance{
		IdeaName: "T",
		State:    "spark",
		Tier:     "focused",
		Counts:   indexmd.ZeroCounts(),
		LogLines: lines,
		Wake:     "none",
	}))
	section := body[strings.Index(body, "## Log tail\n"):strings.Index(body, "## Wake\n")]
	n := 0
	for _, line := range strings.Split(section, "\n") {
		if logfmt.Parseable(line) {
			n++
		}
	}
	if n != 20 {
		t.Fatalf("tail lines=%d want 20", n)
	}
}

func TestRenderWakeLatest(t *testing.T) {
	body := string(indexmd.Render(indexmd.Instance{
		IdeaName: "T",
		State:    "exploring",
		Tier:     "focused",
		Counts:   indexmd.ZeroCounts(),
		Wake:     "briefs/LATEST.md",
	}))
	if !strings.Contains(body, "## Wake\n\nbriefs/LATEST.md\n") {
		t.Fatalf("wake path:\n%s", body)
	}
}

func TestLoadCountsAndWake(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "mycelium.toml"), `schema_version = 1
idea_name = "Garden Lighting"
slug = "garden-lighting"
state = "spark"
tier = "focused"
methodology_version = "2.0.0"
generated_by_cli_version = "0.1.0-dev"
created_date = "2026-08-15"
updated_date = "2026-08-15"
revisit = ""
github_repo = ""
`)
	write(t, filepath.Join(root, "log.md"), "# Log\n\n2026-08-15\tscaffold\t-\tGarden Lighting\n")
	dec := filepath.Join(root, "decisions")
	if err := os.MkdirAll(dec, 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(dec, "DEC-001-first.md"), "x")
	write(t, filepath.Join(dec, "README.md"), "ignore")
	if err := os.MkdirAll(filepath.Join(root, "briefs"), 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(root, "briefs", "LATEST.md"), "wake")

	inst, err := indexmd.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if inst.IdeaName != "Garden Lighting" || inst.Wake != "briefs/LATEST.md" {
		t.Fatalf("inst=%+v", inst)
	}
	found := false
	for _, c := range inst.Counts {
		if c.NS == "DEC" {
			found = true
			if c.Count != 1 {
				t.Fatalf("DEC count=%d", c.Count)
			}
		}
		if c.NS == "ASM" && c.Count != 0 {
			t.Fatalf("ASM count=%d", c.Count)
		}
	}
	if !found {
		t.Fatal("DEC missing from counts")
	}
	if len(inst.LogLines) != 1 {
		t.Fatalf("log lines=%v", inst.LogLines)
	}
}

func write(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}
