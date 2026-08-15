package check_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/robertguss/mycelium/internal/check"
	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/execrun"
	"github.com/robertguss/mycelium/internal/journal"
	"github.com/robertguss/mycelium/internal/manifest"
	"github.com/robertguss/mycelium/internal/op"
	"github.com/robertguss/mycelium/internal/scaffold"
)

func scaffoldOffline(t *testing.T, parent, name string) string {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := scaffold.Run(scaffold.Options{
		Name:    name,
		Dir:     filepath.Join(parent, "inst"),
		Offline: true,
		Cwd:     parent,
		Argv:    []string{"new", "idea", name, "--offline"},
	}, scaffold.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC)},
		Runner: execrun.Real{},
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr.String())
	}
	return filepath.Join(parent, "inst")
}

func TestSparkInstanceCheckOK(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Garden Lighting")
	r := check.Run(root)
	if !r.OK {
		t.Fatalf("findings=%v", r.Findings)
	}
	if r.Slug != "garden-lighting" || r.State != "spark" || r.Tier != "focused" {
		t.Fatalf("summary = %+v", r)
	}
	if r.Artifacts != 0 {
		t.Fatalf("artifacts=%d want 0", r.Artifacts)
	}
}

func TestClarifiedStateLegal(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "State Clarified")
	path := filepath.Join(root, "mycelium.toml")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	patched := strings.Replace(string(b), `state = 'spark'`, `state = 'clarified'`, 1)
	if patched == string(b) {
		patched = strings.Replace(string(b), `state = "spark"`, `state = "clarified"`, 1)
	}
	if patched == string(b) {
		t.Fatalf("could not patch state in %q", b)
	}
	if err := os.WriteFile(path, []byte(patched), 0o644); err != nil {
		t.Fatal(err)
	}
	r := check.Run(root)
	if !r.OK {
		t.Fatalf("want clarified legal, got %v", r.Findings)
	}
}

func TestHandedOffStateFail(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "State Handed Off")
	path := filepath.Join(root, "mycelium.toml")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	patched := strings.Replace(string(b), `state = 'spark'`, `state = 'handed-off'`, 1)
	if patched == string(b) {
		patched = strings.Replace(string(b), `state = "spark"`, `state = "handed-off"`, 1)
	}
	if patched == string(b) {
		t.Fatalf("could not patch state in %q", b)
	}
	if err := os.WriteFile(path, []byte(patched), 0o644); err != nil {
		t.Fatal(err)
	}
	r := check.Run(root)
	if r.OK {
		t.Fatal("want failure")
	}
	found := false
	for _, f := range r.Findings {
		if f.Convention == "lifecycle" && strings.Contains(f.What, "handed-off") && strings.Contains(f.What, "PHASE-06") {
			found = true
			if f.Contract != "program/contracts/lifecycle.md" {
				t.Fatalf("contract=%q", f.Contract)
			}
		}
	}
	if !found {
		t.Fatalf("want PHASE-06 lifecycle finding, got %v", r.Findings)
	}
}

func TestSimmeringBadRevisitFail(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "State Fail")
	path := filepath.Join(root, "mycelium.toml")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	m, err := manifest.Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	m.State = "simmering"
	m.Revisit = "in two weeks"
	out, err := manifest.Encode(m)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, out, 0o644); err != nil {
		t.Fatal(err)
	}
	r := check.Run(root)
	if r.OK {
		t.Fatal("want failure")
	}
	found := false
	for _, f := range r.Findings {
		if f.Convention == "revisit" && strings.Contains(f.What, "simmering") {
			found = true
			if f.Contract != "program/contracts/revisit.md" {
				t.Fatalf("contract=%q", f.Contract)
			}
		}
	}
	if !found {
		t.Fatalf("want revisit finding, got %v", r.Findings)
	}
}

func TestIDToPathMismatch(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Path Fail")
	dec := filepath.Join(root, "decisions")
	if err := os.MkdirAll(dec, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dec, "NOT-AN-ID.md"), []byte("bad"), 0o644); err != nil {
		t.Fatal(err)
	}
	r := check.Run(root)
	if r.OK {
		t.Fatal("want failure")
	}
	found := false
	for _, f := range r.Findings {
		if f.Convention == "id-to-path" && strings.Contains(f.What, "NOT-AN-ID.md") {
			found = true
		}
	}
	if !found {
		t.Fatalf("want id-to-path finding, got %v", r.Findings)
	}
}

func TestNestedArtifactReported(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Nested Fail")
	nested := filepath.Join(root, "decisions", "nested")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	rel := "decisions/nested/DEC-001-x.md"
	if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(rel)), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	r := check.Run(root)
	if r.OK {
		t.Fatal("want failure")
	}
	found := false
	for _, f := range r.Findings {
		if f.Convention == "id-to-path" && strings.Contains(f.What, rel) {
			found = true
		}
	}
	if !found {
		t.Fatalf("want nested id-to-path finding, got %v", r.Findings)
	}
}

func TestAbortJournalViaCheckSurvivesInstanceFrom(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Abort Survive")
	readme := filepath.Join(root, "README.md")
	manifestPath := filepath.Join(root, "mycelium.toml")
	readmeBefore, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	manifestBefore, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, ".mycelium", "stage", "op1"), 0o755); err != nil {
		t.Fatal(err)
	}
	j := &journal.Journal{
		SchemaVersion: 1,
		Op:            "scaffold",
		StartedAt:     "2026-08-15T12:00:00Z",
		StagedDir:     ".mycelium/stage/op1",
		Argv:          []string{"new", "idea", "X"},
		Renames: []journal.Rename{
			{From: "README.md", To: "README.md", Done: false},
			{From: "mycelium.toml", To: "mycelium.toml", Done: false},
		},
	}
	if err := journal.Save(root, j); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	if err := check.AbortJournal(root, &stdout); err != nil {
		t.Fatal(err)
	}
	if b, err := os.ReadFile(readme); err != nil || !bytes.Equal(b, readmeBefore) {
		t.Fatalf("README.md changed: %q err=%v", b, err)
	}
	if b, err := os.ReadFile(manifestPath); err != nil || !bytes.Equal(b, manifestBefore) {
		t.Fatalf("mycelium.toml changed: %q err=%v", b, err)
	}
	if _, err := journal.Load(root); !errors.Is(err, journal.ErrNotExist) {
		t.Fatal("journal should be gone")
	}
}

func TestLeftoverJournalFails(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Journal Fail")
	j := &journal.Journal{
		SchemaVersion: 1,
		Op:            "new",
		StartedAt:     "2026-08-15T12:00:00Z",
		StagedDir:     ".mycelium/stage/leftover",
		Argv:          []string{"new", "decision", "X"},
	}
	if err := os.MkdirAll(filepath.Join(root, ".mycelium", "stage", "leftover"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := journal.Save(root, j); err != nil {
		t.Fatal(err)
	}
	r := check.Run(root)
	if r.OK {
		t.Fatal("want failure")
	}
	found := false
	for _, f := range r.Findings {
		if f.Convention == "operation-protocol" && strings.Contains(f.Fix, "--abort-journal") {
			found = true
		}
	}
	if !found {
		t.Fatalf("want journal recovery finding, got %v", r.Findings)
	}
}

func TestUndeclaredExtraTopLevel(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Extra Fail")
	if err := os.WriteFile(filepath.Join(root, "notes.md"), []byte("scratch"), 0o644); err != nil {
		t.Fatal(err)
	}
	r := check.Run(root)
	if r.OK {
		t.Fatal("want failure")
	}
	found := false
	for _, f := range r.Findings {
		if f.Convention == "extra-top-level" && strings.Contains(f.What, "notes.md") {
			found = true
		}
	}
	if !found {
		t.Fatalf("want extra-top-level finding, got %v", r.Findings)
	}
}

func TestLegalNextTable(t *testing.T) {
	cases := map[string][]string{
		"spark":      {"exploring", "archived"},
		"exploring":  {"simmering", "clarified", "archived"},
		"simmering":  {"exploring", "archived"},
		"clarified":  {"archived"},
		"handed-off": nil,
		"archived":   nil,
		"nope":       nil,
	}
	for from, want := range cases {
		got := check.LegalNext(from)
		if len(got) != len(want) {
			t.Fatalf("%s: got %v want %v", from, got, want)
			continue
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("%s[%d]=%q want %q", from, i, got[i], want[i])
			}
		}
	}
}

func TestAbortStagedDirDotDoesNotWipeViaCheck(t *testing.T) {
	root := t.TempDir()
	marker := filepath.Join(root, "KEEP.md")
	if err := os.WriteFile(marker, []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, ".mycelium"), 0o755); err != nil {
		t.Fatal(err)
	}
	j := &journal.Journal{
		SchemaVersion: 1,
		Op:            "scaffold",
		StartedAt:     "2026-08-15T12:00:00Z",
		StagedDir:     ".",
		Argv:          []string{"new", "idea", "X"},
	}
	if err := journal.Save(root, j); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	if err := check.AbortJournal(root, &stdout); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Fatal("instance wiped")
	}
	if _, err := journal.Load(root); !errors.Is(err, journal.ErrNotExist) {
		t.Fatal("journal should be gone")
	}
}

func TestAbortNothing(t *testing.T) {
	err := check.AbortJournal(t.TempDir(), io.Discard)
	if !errors.Is(err, op.ErrNothingToAbort) {
		t.Fatalf("got %v", err)
	}
}

func writeOQ(t *testing.T, root, name, body string) {
	t.Helper()
	qdir := filepath.Join(root, "questions")
	if err := os.MkdirAll(qdir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(qdir, name), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func findingHas(fs []check.Finding, substrs ...string) bool {
	for _, f := range fs {
		blob := f.What + "\n" + f.Convention + "\n" + f.Contract + "\n" + f.Fix
		ok := true
		for _, s := range substrs {
			if !strings.Contains(blob, s) {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
	}
	return false
}

const oqFront = `+++
id = "OQ-001"
title = "Use SQLite"
agreement = %q
date = "2026-08-15"
+++

# OQ-001 — Use SQLite

## Question

q

## Context

c

## Positions
%s
## Disposition

d
`

func disputedCompleteBody() string {
	return fmt.Sprintf(oqFront, "agree-to-disagree", `

### Human

h

### Agent

a

## Reasons

### Human

rh

### Agent

ra

## Crux

### Human

ch

### Agent

ca

`)
}

func TestAgreeToDisagreeMissingCruxFails(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Slice2 Missing Crux")
	body := fmt.Sprintf(oqFront, "agree-to-disagree", `

### Human

h

### Agent

a

## Reasons

### Human

rh

### Agent

ra

`)
	writeOQ(t, root, "OQ-001-use-sqlite.md", body)
	r := check.Run(root)
	if r.OK {
		t.Fatal("want fail: disputed missing ## Crux")
	}
	if !findingHas(r.Findings, "## Crux", "program/contracts/sparring.md", "sparring") {
		t.Fatalf("want ## Crux + sparring.md finding, got %v", r.Findings)
	}
}

func TestAgreeToDisagreeCompletePasses(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Slice2 Complete")
	writeOQ(t, root, "OQ-001-use-sqlite.md", disputedCompleteBody())
	r := check.Run(root)
	if !r.OK {
		t.Fatalf("want pass, findings=%v", r.Findings)
	}
}

func TestAgreeToDisagreeMissingHumanUnderPositionsFails(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Slice2 Missing Human")
	body := fmt.Sprintf(oqFront, "agree-to-disagree", `

### Agent

a

## Reasons

### Human

rh

### Agent

ra

## Crux

### Human

ch

### Agent

ca

`)
	writeOQ(t, root, "OQ-001-use-sqlite.md", body)
	r := check.Run(root)
	if r.OK {
		t.Fatal("want fail")
	}
	if !findingHas(r.Findings, "### Human under ## Positions", "program/contracts/sparring.md") {
		t.Fatalf("want ### Human finding, got %v", r.Findings)
	}
}

func TestAgreeToDisagreeMissingReasonsFails(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Slice2 Missing Reasons")
	body := fmt.Sprintf(oqFront, "agree-to-disagree", `

### Human

h

### Agent

a

## Crux

### Human

ch

### Agent

ca

`)
	writeOQ(t, root, "OQ-001-use-sqlite.md", body)
	r := check.Run(root)
	if r.OK {
		t.Fatal("want fail")
	}
	if !findingHas(r.Findings, "## Reasons", "program/contracts/sparring.md", "sparring") {
		t.Fatalf("want ## Reasons finding, got %v", r.Findings)
	}
}

func TestAlignedWithoutCruxReasonsPasses(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Slice2 Aligned Bare")
	body := fmt.Sprintf(oqFront, "aligned", `

<!-- fill -->

`)
	writeOQ(t, root, "OQ-001-use-sqlite.md", body)
	r := check.Run(root)
	if !r.OK {
		t.Fatalf("want pass, findings=%v", r.Findings)
	}
}

func TestAlignedWithExtraCruxPasses(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Slice2 Aligned Extra")
	body := fmt.Sprintf(oqFront, "aligned", `

pos

## Crux

extra

`)
	writeOQ(t, root, "OQ-001-use-sqlite.md", body)
	r := check.Run(root)
	if !r.OK {
		t.Fatalf("want pass, findings=%v", r.Findings)
	}
}

func TestInvalidAgreementFails(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Slice2 Bad Agree")
	body := fmt.Sprintf(oqFront, "maybe", `

<!-- fill -->

`)
	writeOQ(t, root, "OQ-001-use-sqlite.md", body)
	r := check.Run(root)
	if r.OK {
		t.Fatal("want fail")
	}
	if !findingHas(r.Findings, `agreement "maybe" is not open|aligned|agree-to-disagree`, "question-front-matter", "program/templates/question.schema.toml") {
		t.Fatalf("want invalid-agreement teaching shape, got %v", r.Findings)
	}
}

func TestOpenPositionsFillPasses(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Slice2 Open Fill")
	body := fmt.Sprintf(oqFront, "open", `

<!-- fill -->

`)
	writeOQ(t, root, "OQ-001-use-sqlite.md", body)
	r := check.Run(root)
	if !r.OK {
		t.Fatalf("want pass, findings=%v", r.Findings)
	}
}

func TestSparkZeroQuestionsPasses(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Slice2 Spark Zero")
	r := check.Run(root)
	if !r.OK {
		t.Fatalf("want pass, findings=%v", r.Findings)
	}
	if r.Artifacts != 0 {
		t.Fatalf("artifacts=%d want 0", r.Artifacts)
	}
}
