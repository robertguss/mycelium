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
	"github.com/robertguss/mycelium/internal/manifest"
	"github.com/robertguss/mycelium/internal/version"
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
	// rough: use known slugify for test names we control
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
	// Occupy next destination path without counting as a numbered artifact
	// (directory → skipped by nextID scan; Stat still sees the path).
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
	got := generate.ReplaceTokens(
		"id={{ID}} title={{TITLE}} slug={{SLUG}} date={{DATE}} keep={{FOO}}",
		"DEC-001", "Hello World", "hello-world", "2026-08-15",
	)
	want := "id=DEC-001 title=Hello World slug=hello-world date=2026-08-15 keep={{FOO}}"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}

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

func TestVersionStillDev(t *testing.T) {
	if version.Version != "0.1.0-dev" {
		t.Fatalf("version=%q", version.Version)
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
	// Ensure --dir nested path still finds instance root.
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
