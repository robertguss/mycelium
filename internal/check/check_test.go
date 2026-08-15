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

func TestGlossaryH1OnlyPasses(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Slice3 H1 Only")
	if err := os.WriteFile(filepath.Join(root, "CONTEXT.md"), []byte("# Glossary\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	r := check.Run(root)
	if !r.OK {
		t.Fatalf("want pass, findings=%v", r.Findings)
	}
}

func TestGlossaryTermMissingDefinitionFails(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Slice3 Missing Def")
	if err := os.WriteFile(filepath.Join(root, "CONTEXT.md"), []byte("# Glossary\n\n## SQLite\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	r := check.Run(root)
	if r.OK {
		t.Fatal("want fail: term missing ### Definition")
	}
	if !findingHas(r.Findings, `CONTEXT.md term "SQLite" missing ### Definition`, "glossary", "program/contracts/glossary.md") {
		t.Fatalf("want glossary teaching shape, got %v", r.Findings)
	}
}

func TestGlossaryTermWithFillDefinitionPasses(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Slice3 Fill Def")
	body := "# Glossary\n\n## SQLite\n\n### Definition\n\n<!-- fill -->\n"
	if err := os.WriteFile(filepath.Join(root, "CONTEXT.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	r := check.Run(root)
	if !r.OK {
		t.Fatalf("want pass, findings=%v", r.Findings)
	}
}

func TestGlossaryMissingH1Fails(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Slice3 Missing H1")
	if err := os.WriteFile(filepath.Join(root, "CONTEXT.md"), []byte("## SQLite\n\n### Definition\n\nx\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	r := check.Run(root)
	if r.OK {
		t.Fatal("want fail: missing H1")
	}
	if !findingHas(r.Findings, "CONTEXT.md missing H1 # Glossary", "glossary", "program/contracts/glossary.md") {
		t.Fatalf("want missing-H1 teaching shape, got %v", r.Findings)
	}
}

func TestGlossaryMissingFileNoNewBind(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Slice3 No Context")
	tierPath := filepath.Join(root, "program", "tiers", "focused.toml")
	tb, err := os.ReadFile(tierPath)
	if err != nil {
		t.Fatal(err)
	}
	patched := strings.ReplaceAll(string(tb), `"CONTEXT.md", `, "")
	patched = strings.ReplaceAll(patched, `, "CONTEXT.md"`, "")
	patched = strings.ReplaceAll(patched, `"CONTEXT.md"`, "")
	if err := os.WriteFile(tierPath, []byte(patched), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(root, "CONTEXT.md")); err != nil {
		t.Fatal(err)
	}
	r := check.Run(root)
	if !r.OK {
		t.Fatalf("want pass without CONTEXT.md (no glossary required-file bind), findings=%v", r.Findings)
	}
	for _, f := range r.Findings {
		if f.Convention == "glossary" {
			t.Fatalf("glossary must not fire when CONTEXT.md absent: %v", f)
		}
	}
}

const decBody = `+++
id = "DEC-001"
title = "Ship it"
status = "Proposed"
date = "2026-08-15"
owner = "TBD"
+++

# DEC-001 — Ship it

## Context

c

## Decision

d

## Rationale

r

## Consequences

co

## Alternatives Considered

a

## Risks

ri

## Revisit Triggers

rt

## Approval

ok
`

func writeDEC(t *testing.T, root, name, body string) {
	t.Helper()
	dir := filepath.Join(root, "decisions")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestDECWithoutDissentPasses(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Slice4 No Dissent")
	writeDEC(t, root, "DEC-001-ship-it.md", decBody)
	r := check.Run(root)
	if !r.OK {
		t.Fatalf("want pass without ## Dissent, findings=%v", r.Findings)
	}
}

func TestDECDissentCitingRealOQPasses(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Slice4 Dissent OQ")
	writeOQ(t, root, "OQ-001-use-sqlite.md", fmt.Sprintf(oqFront, "open", "\n<!-- fill -->\n"))
	body := decBody + "\n## Dissent\n\nStill disagree; see OQ-001.\n"
	writeDEC(t, root, "DEC-001-ship-it.md", body)
	r := check.Run(root)
	if !r.OK {
		t.Fatalf("want pass with resolvable OQ in Dissent, findings=%v", r.Findings)
	}
}

func TestDECDissentNoOQOrASMFails(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Slice4 Dissent Empty")
	body := decBody + "\n## Dissent\n\nI object.\n"
	writeDEC(t, root, "DEC-001-ship-it.md", body)
	r := check.Run(root)
	if r.OK {
		t.Fatal("want fail: Dissent without OQ/ASM")
	}
	if !findingHas(r.Findings, "## Dissent has no resolvable OQ-### or ASM-###", "program/contracts/sparring.md", "dissent") {
		t.Fatalf("want dissent teaching shape, got %v", r.Findings)
	}
	if !findingHas(r.Findings, "cite an existing OQ-### or ASM-### in ## Dissent, or remove the heading") {
		t.Fatalf("want fix text, got %v", r.Findings)
	}
}

func TestDECDissentCitingOnlyDECFails(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Slice4 Dissent DEC Only")
	body := decBody + "\n## Dissent\n\nSee DEC-001 for context.\n"
	writeDEC(t, root, "DEC-001-ship-it.md", body)
	r := check.Run(root)
	if r.OK {
		t.Fatal("want fail: DEC-only token does not satisfy dissent")
	}
	if !findingHas(r.Findings, "DEC-001 ## Dissent has no resolvable OQ-### or ASM-###", "dissent", "program/contracts/sparring.md") {
		t.Fatalf("want dissent finding, got %v", r.Findings)
	}
}

func TestUnresolvedCommissioningFrontMatterFails(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Pack Link")
	dir := filepath.Join(root, "reviews", "reports")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := `+++
id = "RPT-001"
title = "SQLite report"
date = "2026-08-15"
model = "fixture"
commissioning = "CMP-999"
rung = "second-opinion"
adapter = "manual"
prompt_sha256 = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
+++

# RPT-001 — SQLite report

## Position

position

## Findings

findings

## Dissent

none
`
	if err := os.WriteFile(filepath.Join(dir, "RPT-001-sqlite-report.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	r := check.Run(root)
	if r.OK {
		t.Fatal("want unresolved commissioning failure")
	}
	if !findingHas(r.Findings, "reference CMP-999 has no file", "reviews/commissioning") {
		t.Fatalf("want unresolved front-matter link finding, got %v", r.Findings)
	}
}

const (
	testPromptB = "Should this idea use SQLite as the store? Answer independently. Do not see other reports."
	testHashB   = "ec87bfc2afd545807ca87b5c29cae8e77262cb3c746fc63e4539d8daeb2a77de"
	testPromptC = "Review the SQLite store decision. Work independently. Do not see other reports. Retain dissent."
	testHashC   = "8997334f7f2f0bf821bce8ccc4a8d6cf027317c6c66d821200a032a6a11ce098"
)

func TestLadderOptInFalse(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Ladder OptIn Unit")
	writeLadderCMP(t, root, "CMP-001-x.md", ladderCMPBody("CMP-001", "second-opinion", false, "cheap", "manual", testPromptB))
	r := check.Run(root)
	if r.OK {
		t.Fatal("want fail")
	}
	if !findingHas(r.Findings, "opt_in must be true", "council-opt-in", "program/packs/council/contracts/commissioning.md") {
		t.Fatalf("findings=%v", r.Findings)
	}
}

func TestLadderCouncilCheap(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Ladder Cheap Unit")
	writeLadderCMP(t, root, "CMP-001-x.md", ladderCMPBody("CMP-001", "council", true, "cheap", "manual", testPromptC))
	r := check.Run(root)
	if r.OK {
		t.Fatal("want fail")
	}
	if !findingHas(r.Findings, `cost_class "cheap"`, "council-cost-class") {
		t.Fatalf("findings=%v", r.Findings)
	}
}

func TestLadderLoneCMP(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Ladder Lone Unit")
	writeLadderCMP(t, root, "CMP-001-x.md", ladderCMPBody("CMP-001", "second-opinion", true, "cheap", "manual", testPromptB))
	r := check.Run(root)
	if !r.OK {
		t.Fatalf("want pass, findings=%v", r.Findings)
	}
}

func TestLadderSecondOpinionTwoRPT(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Ladder Two Unit")
	writeLadderCMP(t, root, "CMP-001-x.md", ladderCMPBody("CMP-001", "second-opinion", true, "cheap", "manual", testPromptB))
	writeLadderRPT(t, root, "RPT-001-a.md", ladderRPTBody("RPT-001", "CMP-001", "second-opinion", "manual", testHashB, "none"))
	writeLadderRPT(t, root, "RPT-002-b.md", ladderRPTBody("RPT-002", "CMP-001", "second-opinion", "manual", testHashB, "none"))
	r := check.Run(root)
	if r.OK {
		t.Fatal("want fail")
	}
	if !findingHas(r.Findings, "second-opinion-cardinality") {
		t.Fatalf("findings=%v", r.Findings)
	}
}

func TestLadderCouncilHappyAndSeed(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Ladder Happy Unit")
	writeLadderCMP(t, root, "CMP-001-x.md", ladderCMPBody("CMP-001", "council", true, "standard", "cursor", testPromptC))
	writeLadderRPT(t, root, "RPT-001-a.md", ladderRPTBody("RPT-001", "CMP-001", "council", "cursor", testHashC, "SEED-DISSENT"))
	writeLadderRPT(t, root, "RPT-002-b.md", ladderRPTBody("RPT-002", "CMP-001", "council", "cursor", testHashC, "none"))
	writeLadderRCL(t, root, "RCL-001-x.md", ladderRCLBody("RCL-001", "CMP-001", "SEED-DISSENT"))
	r := check.Run(root)
	if !r.OK {
		t.Fatalf("want pass, findings=%v", r.Findings)
	}

	writeLadderRCL(t, root, "RCL-001-x.md", ladderRCLBody("RCL-001", "CMP-001", "none"))
	r = check.Run(root)
	if r.OK {
		t.Fatal("want seed fail")
	}
	if !findingHas(r.Findings, "SEED-DISSENT", "seeded-dissent", "program/packs/council/contracts/reconciliation.md") {
		t.Fatalf("findings=%v", r.Findings)
	}
}

func TestLadderHashMismatch(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Ladder Hash Unit")
	writeLadderCMP(t, root, "CMP-001-x.md", ladderCMPBody("CMP-001", "second-opinion", true, "cheap", "manual", testPromptB))
	writeLadderRPT(t, root, "RPT-001-a.md", ladderRPTBody("RPT-001", "CMP-001", "second-opinion", "manual",
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "none"))
	r := check.Run(root)
	if r.OK {
		t.Fatal("want fail")
	}
	if !findingHas(r.Findings, "prompt_sha256 mismatch", testHashB, "prompt-identity") {
		t.Fatalf("findings=%v", r.Findings)
	}
}

func TestLadderCouncilOneRPT(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Ladder One Unit")
	writeLadderCMP(t, root, "CMP-001-x.md", ladderCMPBody("CMP-001", "council", true, "standard", "cursor", testPromptC))
	writeLadderRPT(t, root, "RPT-001-a.md", ladderRPTBody("RPT-001", "CMP-001", "council", "cursor", testHashC, "none"))
	r := check.Run(root)
	if r.OK {
		t.Fatal("want fail")
	}
	if !findingHas(r.Findings, "council-cardinality") {
		t.Fatalf("findings=%v", r.Findings)
	}
}

func TestLadderRungAdapterMismatch(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Ladder Mismatch Unit")
	writeLadderCMP(t, root, "CMP-001-x.md", ladderCMPBody("CMP-001", "second-opinion", true, "cheap", "manual", testPromptB))
	writeLadderRPT(t, root, "RPT-001-a.md", ladderRPTBody("RPT-001", "CMP-001", "council", "cursor", testHashB, "none"))
	r := check.Run(root)
	if r.OK {
		t.Fatal("want fail")
	}
	if !findingHas(r.Findings, "rung") || !findingHas(r.Findings, "adapter") {
		t.Fatalf("findings=%v", r.Findings)
	}
}

func TestLadderOrphanRPT(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Ladder Orphan Unit")
	writeLadderRPT(t, root, "RPT-001-a.md", ladderRPTBody("RPT-001", "CMP-999", "second-opinion", "manual", testHashB, "none"))
	r := check.Run(root)
	if r.OK {
		t.Fatal("want fail")
	}
	if !findingHas(r.Findings, "does not resolve") && !findingHas(r.Findings, "reference CMP-999 has no file") {
		t.Fatalf("findings=%v", r.Findings)
	}
}

func TestLadderOptInStringTrue(t *testing.T) {
	root := scaffoldOffline(t, t.TempDir(), "Ladder OptStr Unit")
	body := `+++
id = "CMP-001"
title = "x"
date = "2026-08-15"
rung = "second-opinion"
opt_in = "true"
cost_class = "cheap"
adapter = "manual"
+++

# CMP-001 — x

## Prompt

` + testPromptB + `

## Attachments

none

## Cost

cheap
`
	writeLadderCMP(t, root, "CMP-001-x.md", body)
	r := check.Run(root)
	if r.OK {
		t.Fatal("want fail on string opt_in")
	}
	if !findingHas(r.Findings, "opt_in must be true", "council-opt-in") {
		t.Fatalf("findings=%v", r.Findings)
	}
}

func writeLadderCMP(t *testing.T, root, name, body string) {
	t.Helper()
	dir := filepath.Join(root, "reviews", "commissioning")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeLadderRPT(t *testing.T, root, name, body string) {
	t.Helper()
	dir := filepath.Join(root, "reviews", "reports")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeLadderRCL(t *testing.T, root, name, body string) {
	t.Helper()
	dir := filepath.Join(root, "reviews", "reconciliations")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func ladderCMPBody(id, rung string, optIn bool, cost, adapter, prompt string) string {
	opt := "true"
	if !optIn {
		opt = "false"
	}
	return fmt.Sprintf(`+++
id = %q
title = "SQLite"
date = "2026-08-15"
rung = %q
opt_in = %s
cost_class = %q
adapter = %q
+++

# %s — SQLite

## Prompt

%s

## Attachments

none

## Cost

%s
`, id, rung, opt, cost, adapter, id, prompt, cost)
}

func ladderRPTBody(id, cmp, rung, adapter, hash, dissent string) string {
	return fmt.Sprintf(`+++
id = %q
title = "Report"
date = "2026-08-15"
model = "model-a"
commissioning = %q
rung = %q
adapter = %q
prompt_sha256 = %q
+++

# %s — Report

## Position

pos

## Findings

find

## Dissent

%s
`, id, cmp, rung, adapter, hash, id, dissent)
}

func ladderRCLBody(id, cmp, retained string) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`+++
id = %q
title = "Rec"
date = "2026-08-15"
commissioning = %q
rung = "council"
+++

# %s — Rec

`, id, cmp, id))
	for _, s := range []string{
		"Convergence",
		"Material disagreement",
		"Evidence unique to one report",
		"Contradictory evidence",
		"Different assumptions",
		"Different scope interpretations",
		"Recommendations independently supported",
		"Questions requiring another spike",
		"Final reconciled recommendation",
	} {
		b.WriteString("## " + s + "\n\nnone\n\n")
	}
	b.WriteString("## Retained dissent\n\n" + retained + "\n")
	return b.String()
}
