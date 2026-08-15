package generate_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/robertguss/mycelium/internal/check"
	"github.com/robertguss/mycelium/internal/cli"
	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/execrun"
	"github.com/robertguss/mycelium/internal/generate"
	"github.com/robertguss/mycelium/internal/idpath"
	"github.com/robertguss/mycelium/internal/journal"
	"github.com/robertguss/mycelium/internal/manifest"
	"github.com/robertguss/mycelium/internal/op"
)

func fixedDeps(t *testing.T, cwd string) cli.Deps {
	t.Helper()
	return cli.Deps{
		Clock:     clock.Fixed{T: time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC)},
		Runner:    &execrun.Recording{Inner: execrun.Real{}},
		Getwd:     func() (string, error) { return cwd, nil },
		LookupEnv: func(string) string { return "" },
	}
}

func scaffoldOffline(t *testing.T, cwd, name string) string {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := cli.Run(
		[]string{"mycelium", "new", "idea", name, "--offline"},
		&stdout, &stderr, fixedDeps(t, cwd),
	)
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr.String())
	}
	slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	inst := filepath.Join(cwd, slug)
	if _, err := os.Stat(filepath.Join(inst, "mycelium.toml")); err != nil {
		t.Fatalf("instance missing at %s: %v", inst, err)
	}
	return inst
}

func writeRanges(t *testing.T, inst string) {
	t.Helper()
	path := filepath.Join(inst, "mycelium.toml")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	m, err := manifest.Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	m.Identifiers = map[string]manifest.Range{
		"findings":        mustRange(t, "FND-001..FND-010", "FND"),
		"recommendations": mustRange(t, "REC-001..REC-010", "REC"),
		"requirements":    mustRange(t, "REQ-001..REQ-010", "REQ"),
	}
	out, err := manifest.Encode(m)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, out, 0o644); err != nil {
		t.Fatal(err)
	}
}

func mustRange(t *testing.T, raw, ns string) manifest.Range {
	t.Helper()
	rg, err := manifest.ParseRange(raw, ns)
	if err != nil {
		t.Fatal(err)
	}
	return rg
}

func runNew(t *testing.T, deps cli.Deps, args ...string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := cli.Run(append([]string{"mycelium"}, args...), &stdout, &stderr, deps)
	return code, stdout.String(), stderr.String()
}

func TestGenerateAllElevenTypesThenCheck(t *testing.T) {
	cwd := t.TempDir()
	inst := scaffoldOffline(t, cwd, "All Types")
	writeRanges(t, inst)
	deps := fixedDeps(t, cwd)

	titles := map[string]string{}
	for _, typ := range idpath.Types() {
		titles[typ.Key] = "Sample " + typ.Key
	}
	for _, typ := range idpath.Types() {
		code, out, errText := runNew(t, deps, "new", typ.Key, titles[typ.Key], "--dir", inst)
		if code != 0 {
			t.Fatalf("new %s exit %d stderr=%q", typ.Key, code, errText)
		}
		if !strings.Contains(out, "created ") {
			t.Fatalf("stdout missing created: %q", out)
		}
	}

	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"mycelium", "check", "--dir", inst}, &stdout, &stderr, deps)
	if code != 0 {
		t.Fatalf("check exit %d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "mycelium check: ok") {
		t.Fatalf("check stdout=%q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "artifacts: 11") {
		t.Fatalf("want artifacts: 11 in %q", stdout.String())
	}
}

func TestRefuseOverwrite(t *testing.T) {
	cwd := t.TempDir()
	inst := scaffoldOffline(t, cwd, "Overwrite")
	deps := fixedDeps(t, cwd)
	// Directory occupies the path for Stat but nextID skips dirs.
	dest := filepath.Join(inst, "decisions", "DEC-001-same-title.md")
	if err := os.MkdirAll(dest, 0o755); err != nil {
		t.Fatal(err)
	}
	code, _, errText := runNew(t, deps, "new", "decision", "Same Title", "--dir", inst)
	if code != 1 {
		t.Fatalf("exit %d want 1 stderr=%q", code, errText)
	}
	if !strings.Contains(errText, "refuse overwrite") {
		t.Fatalf("stderr=%q", errText)
	}
	assertTeachingShape(t, errText)
}

// H1: leftover journal To=README.md + matching new <type> resume must not clobber README.
func TestH1ResumeRefusesClobberREADME(t *testing.T) {
	cwd := t.TempDir()
	inst := scaffoldOffline(t, cwd, "H1 Readme")
	deps := fixedDeps(t, cwd)
	now := time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC)

	readme := filepath.Join(inst, "README.md")
	keep := []byte("KEEP-README-BYTES\n")
	if err := os.WriteFile(readme, keep, 0o644); err != nil {
		t.Fatal(err)
	}

	sess, err := op.Begin(inst, op.Intent{
		Op:         "new",
		Type:       "decision",
		Title:      "Any Title",
		OriginalID: "DEC-001",
		LogLine:    "2026-08-15\tnew\tDEC-001\tAny Title",
		Argv:       []string{"new", "decision", "Any Title", "--dir", inst},
		OpID:       "h1-readme",
	}, now)
	if err != nil {
		t.Fatal(err)
	}
	if err := sess.Stage([]op.Staged{
		{RelTo: "README.md", Content: []byte("PWNED-README\n")},
		{RelTo: "log.md", Content: []byte("log\n")},
		{RelTo: "mycelium.toml", Content: []byte("m\n")},
	}); err != nil {
		t.Fatal(err)
	}
	if err := sess.Close(); err != nil {
		t.Fatal(err)
	}

	code, _, errText := runNew(t, deps, "new", "decision", "Any Title", "--dir", inst)
	if code != 1 {
		t.Fatalf("exit %d want 1 stderr=%q", code, errText)
	}
	got, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(keep) {
		t.Fatalf("README clobbered: got %q", got)
	}
	if strings.Contains(string(got), "PWNED-README") {
		t.Fatal("README contains PWNED-README")
	}
	_ = errText
}

// H1: leftover journal To=decisions/DEC-001-keep-this.md + resume must leave file unchanged.
func TestH1ResumeRefusesClobberExistingDEC(t *testing.T) {
	cwd := t.TempDir()
	inst := scaffoldOffline(t, cwd, "H1 Dec")
	deps := fixedDeps(t, cwd)
	now := time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC)

	rel := "decisions/DEC-001-keep-this.md"
	dest := filepath.Join(inst, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		t.Fatal(err)
	}
	keep := []byte("KEEP-THIS-DEC\n")
	if err := os.WriteFile(dest, keep, 0o644); err != nil {
		t.Fatal(err)
	}

	sess, err := op.Begin(inst, op.Intent{
		Op:         "new",
		Type:       "decision",
		Title:      "Keep This",
		OriginalID: "DEC-001",
		LogLine:    "2026-08-15\tnew\tDEC-001\tKeep This",
		Argv:       []string{"new", "decision", "Keep This", "--dir", inst},
		OpID:       "h1-dec",
	}, now)
	if err != nil {
		t.Fatal(err)
	}
	if err := sess.Stage([]op.Staged{
		{RelTo: rel, Content: []byte("PWNED-DEC\n")},
		{RelTo: "log.md", Content: []byte("log\n")},
		{RelTo: "mycelium.toml", Content: []byte("m\n")},
	}); err != nil {
		t.Fatal(err)
	}
	if err := sess.Close(); err != nil {
		t.Fatal(err)
	}

	code, _, errText := runNew(t, deps, "new", "decision", "Keep This", "--dir", inst)
	if code != 1 {
		t.Fatalf("exit %d want 1 stderr=%q", code, errText)
	}
	if !strings.Contains(errText, "refuse overwrite") && !strings.Contains(errText, "commit failed") {
		t.Fatalf("stderr=%q", errText)
	}
	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(keep) {
		t.Fatalf("DEC clobbered: got %q", got)
	}
}

// H2: leftover original_id + existing dest + empty renames → refuse before Stage.
func TestH2ResumeEmptyRenamesRefusesExistingDest(t *testing.T) {
	cwd := t.TempDir()
	inst := scaffoldOffline(t, cwd, "H2 Empty")
	deps := fixedDeps(t, cwd)
	now := time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC)

	rel := "decisions/DEC-001-title.md"
	dest := filepath.Join(inst, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		t.Fatal(err)
	}
	userBytes := []byte("USER-WROTE-THIS-DEC\n")
	if err := os.WriteFile(dest, userBytes, 0o644); err != nil {
		t.Fatal(err)
	}

	// Journal with original_id but empty renames (pre-Stage crash shape).
	j := &journal.Journal{
		SchemaVersion: 1,
		Op:            "new",
		Title:         "Title",
		OriginalID:    "DEC-001",
		StartedAt:     now.UTC().Format(time.RFC3339),
		StagedDir:     ".mycelium/stage/h2-empty",
		LogLine:       "2026-08-15\tnew\tDEC-001\tTitle",
		Argv:          []string{"new", "decision", "Title", "--dir", inst},
	}
	j.SetType("decision")
	if err := os.MkdirAll(filepath.Join(inst, ".mycelium", "stage", "h2-empty"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := journal.Save(inst, j); err != nil {
		t.Fatal(err)
	}

	code, _, errText := runNew(t, deps, "new", "decision", "Title", "--dir", inst)
	if code != 1 {
		t.Fatalf("exit %d want 1 stderr=%q", code, errText)
	}
	if !strings.Contains(errText, "refuse overwrite") {
		t.Fatalf("stderr=%q", errText)
	}
	assertTeachingShape(t, errText)
	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(userBytes) {
		t.Fatalf("dest changed: %q", got)
	}
}

func TestRefuseTitleNewlineAndTab(t *testing.T) {
	cwd := t.TempDir()
	inst := scaffoldOffline(t, cwd, "Bad Title")
	deps := fixedDeps(t, cwd)
	for _, title := range []string{"line\nbreak", "has\ttab"} {
		code, _, errText := runNew(t, deps, "new", "decision", title, "--dir", inst)
		if code != 1 {
			t.Fatalf("title %q exit %d want 1", title, code)
		}
		if !strings.Contains(errText, "newline or tab") {
			t.Fatalf("title %q stderr=%q", title, errText)
		}
		assertTeachingShape(t, errText)
	}
	// Ensure no log injection line was appended.
	b, err := os.ReadFile(filepath.Join(inst, "log.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), "\tnew\t") && (strings.Contains(string(b), "break") || strings.Contains(string(b), "tab")) {
		// Only scaffold line should exist; refuse before logfmt.Line / Stage.
		lines := strings.Split(strings.TrimSpace(string(b)), "\n")
		for _, line := range lines {
			if strings.Contains(line, "\tnew\t") {
				t.Fatalf("log got new line after title refuse: %q", line)
			}
		}
	}
}

func TestRefuseFindingWithoutRange(t *testing.T) {
	cwd := t.TempDir()
	inst := scaffoldOffline(t, cwd, "No Range")
	deps := fixedDeps(t, cwd)
	code, _, errText := runNew(t, deps, "new", "finding", "Missing Range", "--dir", inst)
	if code != 1 {
		t.Fatalf("exit %d want 1", code)
	}
	if !strings.Contains(errText, "no [identifiers] range") {
		t.Fatalf("stderr=%q", errText)
	}
	assertTeachingShape(t, errText)
	if _, err := os.Stat(filepath.Join(inst, "findings")); !os.IsNotExist(err) {
		t.Fatalf("must not create findings without range (err=%v)", err)
	}
}

func TestRefuseFindingPastRangeEnd(t *testing.T) {
	cwd := t.TempDir()
	inst := scaffoldOffline(t, cwd, "Past End")
	deps := fixedDeps(t, cwd)

	path := filepath.Join(inst, "mycelium.toml")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	m, err := manifest.Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	m.Identifiers = map[string]manifest.Range{
		"findings": mustRange(t, "FND-001..FND-002", "FND"),
	}
	out, err := manifest.Encode(m)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, out, 0o644); err != nil {
		t.Fatal(err)
	}

	for i, title := range []string{"One", "Two"} {
		code, _, errText := runNew(t, deps, "new", "finding", title, "--dir", inst)
		if code != 0 {
			t.Fatalf("alloc %d: %s", i+1, errText)
		}
	}
	code, _, errText := runNew(t, deps, "new", "finding", "Three", "--dir", inst)
	if code != 1 {
		t.Fatalf("exit %d want 1", code)
	}
	if !strings.Contains(errText, "outside declared range") {
		t.Fatalf("stderr=%q", errText)
	}
	assertTeachingShape(t, errText)
}

func TestNextIDIsMaxPlusOneAllowsGaps(t *testing.T) {
	cwd := t.TempDir()
	inst := scaffoldOffline(t, cwd, "Gaps")
	deps := fixedDeps(t, cwd)

	// Seed DEC-001 and DEC-003 (gap at 002).
	if err := os.MkdirAll(filepath.Join(inst, "decisions"), 0o755); err != nil {
		t.Fatal(err)
	}
	seed001 := `+++
id = "DEC-001"
title = "Seed One"
status = "Proposed"
date = "2026-08-15"
owner = ""
+++

# DEC-001 — Seed One

## Context

## Decision

## Rationale

## Consequences

## Alternatives Considered

## Risks

## Revisit Triggers

## Approval
`
	seed003 := strings.ReplaceAll(seed001, "DEC-001", "DEC-003")
	seed003 = strings.ReplaceAll(seed003, "Seed One", "Seed Three")
	if err := os.WriteFile(filepath.Join(inst, "decisions", "DEC-001-seed-one.md"), []byte(seed001), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inst, "decisions", "DEC-003-seed-three.md"), []byte(seed003), 0o644); err != nil {
		t.Fatal(err)
	}

	code, out, errText := runNew(t, deps, "new", "decision", "After Gap", "--dir", inst)
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, errText)
	}
	if !strings.Contains(out, "decisions/DEC-004-after-gap.md") {
		t.Fatalf("want DEC-004, got %q", out)
	}
	if _, err := os.Stat(filepath.Join(inst, "decisions", "DEC-004-after-gap.md")); err != nil {
		t.Fatal(err)
	}
}

func TestLogLineAppended(t *testing.T) {
	cwd := t.TempDir()
	inst := scaffoldOffline(t, cwd, "Log Line")
	deps := fixedDeps(t, cwd)
	code, _, errText := runNew(t, deps, "new", "decision", "Logged Thought", "--dir", inst)
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, errText)
	}
	b, err := os.ReadFile(filepath.Join(inst, "log.md"))
	if err != nil {
		t.Fatal(err)
	}
	want := "2026-08-15\tnew\tDEC-001\tLogged Thought"
	if !strings.Contains(string(b), want) {
		t.Fatalf("log missing %q\n%s", want, b)
	}
}

func TestTokensReplacedLeftoverStays(t *testing.T) {
	cwd := t.TempDir()
	inst := scaffoldOffline(t, cwd, "Tokens")
	deps := fixedDeps(t, cwd)
	code, _, errText := runNew(t, deps, "new", "decision", "Token Check", "--dir", inst)
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, errText)
	}
	b, err := os.ReadFile(filepath.Join(inst, "decisions", "DEC-001-token-check.md"))
	if err != nil {
		t.Fatal(err)
	}
	body := string(b)
	for _, tok := range []string{"{{ID}}", "{{TITLE}}", "{{SLUG}}", "{{DATE}}"} {
		if strings.Contains(body, tok) {
			t.Fatalf("token %s still present:\n%s", tok, body)
		}
	}
	if !strings.Contains(body, `id = "DEC-001"`) {
		t.Fatalf("ID not replaced:\n%s", body)
	}
	if !strings.Contains(body, `title = "Token Check"`) {
		t.Fatalf("TITLE not replaced:\n%s", body)
	}
	if !strings.Contains(body, "token-check") {
		t.Fatalf("SLUG not replaced:\n%s", body)
	}
	if !strings.Contains(body, `date = "2026-08-15"`) {
		t.Fatalf("DATE not replaced:\n%s", body)
	}
}

func TestUnknownTypeTeachingError(t *testing.T) {
	cwd := t.TempDir()
	inst := scaffoldOffline(t, cwd, "Unknown")
	deps := fixedDeps(t, cwd)
	code, _, errText := runNew(t, deps, "new", "widget", "Nope", "--dir", inst)
	if code != 1 {
		t.Fatalf("exit %d want 1", code)
	}
	if !strings.Contains(errText, `unknown type "widget"`) {
		t.Fatalf("stderr=%q", errText)
	}
	if !strings.Contains(errText, "decision") {
		t.Fatalf("registered keys missing from %q", errText)
	}
	assertTeachingShape(t, errText)
}

func TestEmptyTitleAndEmptySlugRefuse(t *testing.T) {
	cwd := t.TempDir()
	inst := scaffoldOffline(t, cwd, "Empty")
	deps := fixedDeps(t, cwd)

	code, _, errText := runNew(t, deps, "new", "decision", "--dir", inst)
	if code != 1 {
		t.Fatalf("empty title exit %d want 1", code)
	}
	if !strings.Contains(errText, "title is required") {
		t.Fatalf("stderr=%q", errText)
	}

	code, _, errText = runNew(t, deps, "new", "decision", "!!!", "--dir", inst)
	if code != 1 {
		t.Fatalf("empty slug exit %d want 1", code)
	}
	if !strings.Contains(errText, "slugify") && !strings.Contains(errText, "cannot slugify") {
		t.Fatalf("stderr=%q", errText)
	}
}

func TestDirFlagWithoutChdir(t *testing.T) {
	cwd := t.TempDir()
	inst := scaffoldOffline(t, cwd, "Dir Flag")
	outside := filepath.Join(cwd, "outside")
	if err := os.Mkdir(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	deps := fixedDeps(t, outside)
	code, out, errText := runNew(t, deps, "new", "assumption", "From Outside", "--dir", inst)
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, errText)
	}
	if !strings.Contains(out, "assumptions/ASM-001-from-outside.md") {
		t.Fatalf("stdout=%q", out)
	}
	if _, err := os.Stat(filepath.Join(inst, "assumptions", "ASM-001-from-outside.md")); err != nil {
		t.Fatal(err)
	}
}

func TestReplaceTokensUnit(t *testing.T) {
	in := "{{ID}}/{{TITLE}}/{{SLUG}}/{{DATE}}/{{FOO}}/{{ID}}"
	got := generate.ReplaceTokens(in, "X", "T", "s", "D")
	if got != "X/T/s/D/{{FOO}}/X" {
		t.Fatalf("got %q", got)
	}
}

func TestGeneratePackageUsesCheckFindRoot(t *testing.T) {
	cwd := t.TempDir()
	inst := scaffoldOffline(t, cwd, "Nested")
	nested := filepath.Join(inst, "program", "templates")
	deps := cli.Deps{
		Clock:     clock.Fixed{T: time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC)},
		Runner:    &execrun.Recording{Inner: execrun.Real{}},
		Getwd:     func() (string, error) { return nested, nil },
		LookupEnv: func(string) string { return "" },
	}
	code, out, errText := runNew(t, deps, "new", "risk", "Nested Start")
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, errText)
	}
	root, err := check.FindRoot(nested)
	if err != nil || root != inst {
		t.Fatalf("FindRoot=%q err=%v want %q", root, err, inst)
	}
	if !strings.Contains(out, "risks/RSK-001-nested-start.md") {
		t.Fatalf("stdout=%q", out)
	}
}

func assertTeachingShape(t *testing.T, errText string) {
	t.Helper()
	for _, prefix := range []string{"mycelium:", "convention:", "contract:", "fix:"} {
		if !strings.Contains(errText, prefix) {
			t.Fatalf("teaching error missing %q in %q", prefix, errText)
		}
	}
}
