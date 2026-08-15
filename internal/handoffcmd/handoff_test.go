package handoffcmd_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/execrun"
	"github.com/robertguss/mycelium/internal/handoff"
	"github.com/robertguss/mycelium/internal/handoffcmd"
	"github.com/robertguss/mycelium/internal/journal"
	"github.com/robertguss/mycelium/internal/lock"
	"github.com/robertguss/mycelium/internal/manifest"
	"github.com/robertguss/mycelium/internal/metadata"
	"github.com/robertguss/mycelium/internal/scaffold"
)

func TestRunHappyAndRefuse(t *testing.T) {
	cwd := t.TempDir()
	deps := fixedDeps(t)
	inst := scaffoldOffline(t, cwd, "Cmd Coverage")
	toClarified(t, deps, inst)

	var stdout, stderr bytes.Buffer
	deps.Stdout = &stdout
	deps.Stderr = &stderr
	code := handoffcmd.Run(handoffcmd.Options{Dir: inst, Cwd: cwd, Argv: []string{"handoff", "--dir", inst}}, deps)
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "mycelium handoff: ok") {
		t.Fatalf("stdout=%q", stdout.String())
	}
	findings := handoff.Check(os.DirFS(filepath.Join(inst, "handoff")))
	if len(findings) > 0 {
		t.Fatalf("%v", findings)
	}

	stdout.Reset()
	stderr.Reset()
	code = handoffcmd.Run(handoffcmd.Options{Dir: inst, Cwd: cwd}, deps)
	if code != 1 {
		t.Fatalf("re-handoff exit %d", code)
	}
	if !strings.Contains(stderr.String(), "already exists") {
		t.Fatalf("stderr=%q", stderr.String())
	}
}

func TestRunRichPacket(t *testing.T) {
	cwd := t.TempDir()
	deps := fixedDeps(t)
	inst := scaffoldOffline(t, cwd, "Rich Packet")
	toClarified(t, deps, inst)

	writeAcceptedDEC(t, inst, "DEC-001", "Use Go", "See EVD-001 and [../evidence/EVD-001-bench.md](../evidence/EVD-001-bench.md).")
	writeOQ(t, inst, "OQ-001", "open")
	writeEVD(t, inst, "EVD-001", "bench")
	writeFile(t, filepath.Join(inst, "CONTEXT.md"), "# Glossary\n\n## Term\n\n### Definition\n\nA thing.\n")
	writeFile(t, filepath.Join(inst, "playbooks", "HOW.md"), "# How\n\n## Target\n\nDo it.\n\n## Steps\n\n1. See [../decisions/DEC-001-use-go.md](../decisions/DEC-001-use-go.md)\n2. See [../CONTEXT.md](../CONTEXT.md)\n\n## Done\n\nGreen tests.\n")
	writeFile(t, filepath.Join(inst, "acceptance", "add_test.go"), "package acceptance\n")
	// log without trailing newline exercises appendLogLine branch
	logPath := filepath.Join(inst, "log.md")
	b := mustRead(t, logPath)
	if err := os.WriteFile(logPath, bytes.TrimRight(b, "\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	deps.Stdout = &stdout
	deps.Stderr = &stderr
	code := handoffcmd.Run(handoffcmd.Options{Dir: inst, Cwd: cwd}, deps)
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	packet := string(mustRead(t, filepath.Join(inst, "handoff", "PACKET.md")))
	for _, want := range []string{"DEC-001", "OQ-001", "evidence/SUMMARY.md", "playbooks/HOW.md", "acceptance/add_test.go"} {
		if !strings.Contains(packet, want) {
			t.Fatalf("packet missing %q:\n%s", want, packet)
		}
	}
	play := string(mustRead(t, filepath.Join(inst, "handoff", "playbooks", "HOW.md")))
	if strings.Contains(play, "../decisions/") && !strings.Contains(play, "](../decisions/") {
		t.Fatalf("playbook links not rewritten:\n%s", play)
	}
	if !strings.Contains(play, "../decisions/DEC-001-use-go.md") && !strings.Contains(play, "decisions/DEC-001") {
		// rewrite should produce in-packet relative from playbooks/
		t.Fatalf("playbook missing in-packet decision link:\n%s", play)
	}
	glossary := string(mustRead(t, filepath.Join(inst, "handoff", "glossary.md")))
	if !strings.Contains(glossary, "Glossary") {
		t.Fatalf("glossary=%q", glossary)
	}
	summary := string(mustRead(t, filepath.Join(inst, "handoff", "evidence", "SUMMARY.md")))
	if !strings.Contains(summary, "EVD-001") {
		t.Fatalf("summary=%q", summary)
	}
	findings := handoff.Check(os.DirFS(filepath.Join(inst, "handoff")))
	if len(findings) > 0 {
		t.Fatalf("structure: %v", findings)
	}
}

func TestRunNotClarified(t *testing.T) {
	cwd := t.TempDir()
	deps := fixedDeps(t)
	inst := scaffoldOffline(t, cwd, "Not Clarified")
	var stderr bytes.Buffer
	deps.Stderr = &stderr
	code := handoffcmd.Run(handoffcmd.Options{Dir: inst, Cwd: cwd}, deps)
	if code != 1 {
		t.Fatalf("exit %d", code)
	}
	if !strings.Contains(stderr.String(), "legal only from clarified (got spark)") {
		t.Fatalf("stderr=%q", stderr.String())
	}
}

func TestRunNotInstance(t *testing.T) {
	cwd := t.TempDir()
	deps := fixedDeps(t)
	var stderr bytes.Buffer
	deps.Stderr = &stderr
	code := handoffcmd.Run(handoffcmd.Options{Dir: cwd, Cwd: cwd}, deps)
	if code != 1 {
		t.Fatalf("exit %d", code)
	}
	if !strings.Contains(stderr.String(), "not a mycelium instance") {
		t.Fatalf("stderr=%q", stderr.String())
	}
}

func TestRunMissingManifest(t *testing.T) {
	cwd := t.TempDir()
	deps := fixedDeps(t)
	inst := scaffoldOffline(t, cwd, "No Manifest")
	toClarified(t, deps, inst)
	if err := os.Remove(filepath.Join(inst, "mycelium.toml")); err != nil {
		t.Fatal(err)
	}
	// FindRoot walks for mycelium.toml — without it, not an instance.
	var stderr bytes.Buffer
	deps.Stderr = &stderr
	code := handoffcmd.Run(handoffcmd.Options{Dir: inst, Cwd: cwd}, deps)
	if code != 1 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
}

func TestRunInvalidManifest(t *testing.T) {
	cwd := t.TempDir()
	deps := fixedDeps(t)
	inst := scaffoldOffline(t, cwd, "Bad Manifest")
	toClarified(t, deps, inst)
	writeFile(t, filepath.Join(inst, "mycelium.toml"), "nope = true\n")
	var stderr bytes.Buffer
	deps.Stderr = &stderr
	code := handoffcmd.Run(handoffcmd.Options{Dir: inst, Cwd: cwd}, deps)
	if code != 1 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "invalid mycelium.toml") {
		t.Fatalf("stderr=%q", stderr.String())
	}
}

func TestRunMissingLog(t *testing.T) {
	cwd := t.TempDir()
	deps := fixedDeps(t)
	inst := scaffoldOffline(t, cwd, "No Log")
	toClarified(t, deps, inst)
	if err := os.Remove(filepath.Join(inst, "log.md")); err != nil {
		t.Fatal(err)
	}
	var stderr bytes.Buffer
	deps.Stderr = &stderr
	code := handoffcmd.Run(handoffcmd.Options{Dir: inst, Cwd: cwd}, deps)
	if code != 1 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "cannot read log.md") {
		t.Fatalf("stderr=%q", stderr.String())
	}
}

func TestRunJournalMismatch(t *testing.T) {
	cwd := t.TempDir()
	deps := fixedDeps(t)
	inst := scaffoldOffline(t, cwd, "Journal Mismatch")
	toClarified(t, deps, inst)

	j := &journal.Journal{
		SchemaVersion: 1,
		Op:            "state",
		Title:         "x",
		StartedAt:     "2026-08-15T00:00:00Z",
		StagedDir:     ".mycelium/stage/x",
		Argv:          []string{"state", "archived"},
	}
	if err := journal.Save(inst, j); err != nil {
		t.Fatal(err)
	}

	var stderr bytes.Buffer
	deps.Stderr = &stderr
	code := handoffcmd.Run(handoffcmd.Options{Dir: inst, Cwd: cwd}, deps)
	if code != 1 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "leftover journal") {
		t.Fatalf("stderr=%q", stderr.String())
	}
}

func TestRunLocked(t *testing.T) {
	cwd := t.TempDir()
	deps := fixedDeps(t)
	inst := scaffoldOffline(t, cwd, "Locked")
	toClarified(t, deps, inst)

	held, err := lock.Acquire(inst, deps.Clock.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = held.Release() })

	var stderr bytes.Buffer
	deps.Stderr = &stderr
	code := handoffcmd.Run(handoffcmd.Options{Dir: inst, Cwd: cwd}, deps)
	if code != 1 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "lock held") {
		t.Fatalf("stderr=%q", stderr.String())
	}
}

func TestRunRelativeAndAbsoluteDir(t *testing.T) {
	cwd := t.TempDir()
	deps := fixedDeps(t)
	inst := scaffoldOffline(t, cwd, "Rel Abs")
	toClarified(t, deps, inst)
	rel, err := filepath.Rel(cwd, inst)
	if err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	deps.Stdout = &stdout
	deps.Stderr = &stderr
	code := handoffcmd.Run(handoffcmd.Options{Dir: rel, Cwd: cwd}, deps)
	if code != 0 {
		t.Fatalf("rel exit %d stderr=%q", code, stderr.String())
	}

	inst2 := scaffoldOffline(t, cwd, "Abs Dir")
	toClarified(t, deps, inst2)
	stdout.Reset()
	stderr.Reset()
	code = handoffcmd.Run(handoffcmd.Options{Dir: inst2, Cwd: cwd}, deps)
	if code != 0 {
		t.Fatalf("abs exit %d stderr=%q", code, stderr.String())
	}
}

func TestRunEmptyCwdUsesGetwd(t *testing.T) {
	cwd := t.TempDir()
	deps := fixedDeps(t)
	inst := scaffoldOffline(t, cwd, "Getwd")
	toClarified(t, deps, inst)
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(inst); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	var stdout, stderr bytes.Buffer
	deps.Stdout = &stdout
	deps.Stderr = &stderr
	code := handoffcmd.Run(handoffcmd.Options{Dir: "", Cwd: ""}, deps)
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
}

func TestRunNilDeps(t *testing.T) {
	code := handoffcmd.Run(handoffcmd.Options{Dir: t.TempDir(), Cwd: t.TempDir()}, handoffcmd.Deps{})
	if code != 1 {
		t.Fatalf("exit %d", code)
	}
}

func TestRunBuildPacketFail(t *testing.T) {
	cwd := t.TempDir()
	deps := fixedDeps(t)
	inst := scaffoldOffline(t, cwd, "Build Fail")
	toClarified(t, deps, inst)
	// decisions as a file → listAcceptedDecisions ReadDir fails.
	_ = os.RemoveAll(filepath.Join(inst, "decisions"))
	writeFile(t, filepath.Join(inst, "decisions"), "not-a-dir\n")
	var stderr bytes.Buffer
	deps.Stderr = &stderr
	code := handoffcmd.Run(handoffcmd.Options{Dir: inst, Cwd: cwd}, deps)
	if code != 1 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "cannot build handoff packet") {
		t.Fatalf("stderr=%q", stderr.String())
	}
}

func TestRunEvidenceBareIDFilename(t *testing.T) {
	cwd := t.TempDir()
	deps := fixedDeps(t)
	inst := scaffoldOffline(t, cwd, "Bare EVD Name")
	toClarified(t, deps, inst)
	dir := filepath.Join(inst, "evidence")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(dir, "EVD-007.md"), "# bare\n")
	writeAcceptedDEC(t, inst, "DEC-001", "Bare", "EVD-007 cited.")
	var stderr bytes.Buffer
	deps.Stderr = &stderr
	code := handoffcmd.Run(handoffcmd.Options{Dir: inst, Cwd: cwd}, deps)
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
}

func TestRunNestedAcceptance(t *testing.T) {
	cwd := t.TempDir()
	deps := fixedDeps(t)
	inst := scaffoldOffline(t, cwd, "Nested Acc")
	toClarified(t, deps, inst)
	writeFile(t, filepath.Join(inst, "acceptance", "sub", "t_test.go"), "package sub\n")
	var stderr bytes.Buffer
	deps.Stderr = &stderr
	code := handoffcmd.Run(handoffcmd.Options{Dir: inst, Cwd: cwd}, deps)
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	if _, err := os.Stat(filepath.Join(inst, "handoff", "acceptance", "sub", "t_test.go")); err != nil {
		t.Fatal(err)
	}
}

func TestRunQuestionsNotDir(t *testing.T) {
	cwd := t.TempDir()
	deps := fixedDeps(t)
	inst := scaffoldOffline(t, cwd, "OQ File")
	toClarified(t, deps, inst)
	_ = os.RemoveAll(filepath.Join(inst, "questions"))
	writeFile(t, filepath.Join(inst, "questions"), "file\n")
	var stderr bytes.Buffer
	deps.Stderr = &stderr
	code := handoffcmd.Run(handoffcmd.Options{Dir: inst, Cwd: cwd}, deps)
	if code != 1 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
}

func TestRunPlaybooksNotDir(t *testing.T) {
	cwd := t.TempDir()
	deps := fixedDeps(t)
	inst := scaffoldOffline(t, cwd, "PB File")
	toClarified(t, deps, inst)
	writeFile(t, filepath.Join(inst, "playbooks"), "file\n")
	var stderr bytes.Buffer
	deps.Stderr = &stderr
	// File named playbooks is ignored; stub playbook is emitted.
	code := handoffcmd.Run(handoffcmd.Options{Dir: inst, Cwd: cwd}, deps)
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	if _, err := os.Stat(filepath.Join(inst, "handoff", "playbooks", "PLAYBOOK.md")); err != nil {
		t.Fatal(err)
	}
}

func TestRunMissingCONTEXT(t *testing.T) {
	cwd := t.TempDir()
	deps := fixedDeps(t)
	inst := scaffoldOffline(t, cwd, "No Context File")
	toClarified(t, deps, inst)
	_ = os.Remove(filepath.Join(inst, "CONTEXT.md"))
	var stderr bytes.Buffer
	deps.Stderr = &stderr
	code := handoffcmd.Run(handoffcmd.Options{Dir: inst, Cwd: cwd}, deps)
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	g := string(mustRead(t, filepath.Join(inst, "handoff", "glossary.md")))
	if strings.TrimSpace(g) != "none" {
		t.Fatalf("glossary=%q", g)
	}
}

func TestRunEvidenceNoFrontMatter(t *testing.T) {
	cwd := t.TempDir()
	deps := fixedDeps(t)
	inst := scaffoldOffline(t, cwd, "EVD No Meta")
	toClarified(t, deps, inst)
	dir := filepath.Join(inst, "evidence")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	// No front matter → idFromFilename path.
	writeFile(t, filepath.Join(dir, "EVD-002-raw.md"), "# Raw evidence\n")
	writeAcceptedDEC(t, inst, "DEC-001", "Cite", "Mentions EVD-002.")
	var stderr bytes.Buffer
	deps.Stderr = &stderr
	code := handoffcmd.Run(handoffcmd.Options{Dir: inst, Cwd: cwd}, deps)
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	if _, err := os.Stat(filepath.Join(inst, "handoff", "evidence", "EVD-002-raw.md")); err != nil {
		t.Fatal(err)
	}
}

func TestRunRewriteLinkVariants(t *testing.T) {
	cwd := t.TempDir()
	deps := fixedDeps(t)
	inst := scaffoldOffline(t, cwd, "Link Variants")
	toClarified(t, deps, inst)
	writeAcceptedDEC(t, inst, "DEC-001", "Links", "none")
	writeFile(t, filepath.Join(inst, "playbooks", "P.md"), `# P

## Target

See [abs](/etc/passwd) and [web](https://example.com) and [mail](mailto:a@b.c) and [hash](#sec) and [titled](../decisions/DEC-001-links.md "title").

## Steps

none

## Done

none
`)
	var stderr bytes.Buffer
	deps.Stderr = &stderr
	code := handoffcmd.Run(handoffcmd.Options{Dir: inst, Cwd: cwd}, deps)
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	play := string(mustRead(t, filepath.Join(inst, "handoff", "playbooks", "P.md")))
	if !strings.Contains(play, "https://example.com") || !strings.Contains(play, "mailto:") {
		t.Fatalf("external links lost:\n%s", play)
	}
}

func TestRunCanonAndSkipEdges(t *testing.T) {
	cwd := t.TempDir()
	deps := fixedDeps(t)
	inst := scaffoldOffline(t, cwd, "Canon Edges")
	toClarified(t, deps, inst)
	dir := filepath.Join(inst, "decisions")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(dir, "README.md"), "ignore\n")
	writeFile(t, filepath.Join(dir, "DEC-009-bad.md"), "not toml front matter\n")
	writeFile(t, filepath.Join(dir, "DEC-010-num.md"), `+++
id = "DEC-010"
title = 42
status = "Accepted"
date = "2026-08-15"
owner = "TBD"
+++

# DEC-010

Cites EVD-000.
`)
	edir := filepath.Join(inst, "evidence")
	if err := os.MkdirAll(edir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(edir, "EVD-000-z.md"), `+++
id = "EVD-000"
title = "Z"
date = "2026-08-15"
owner = "TBD"
+++

# EVD-000
`)
	qdir := filepath.Join(inst, "questions")
	if err := os.MkdirAll(qdir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(qdir, "README.md"), "x\n")
	writeFile(t, filepath.Join(qdir, "not-an-oq.txt"), "skip\n")

	writeFile(t, filepath.Join(inst, "playbooks", "P.md"), `# P

## Target

[http](http://example.com) [titled](../CONTEXT.md "ctx")

## Steps

none

## Done

none
`)

	var stderr bytes.Buffer
	deps.Stderr = &stderr
	code := handoffcmd.Run(handoffcmd.Options{Dir: inst, Cwd: cwd}, deps)
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	packet := string(mustRead(t, filepath.Join(inst, "handoff", "PACKET.md")))
	if !strings.Contains(packet, "DEC-010") {
		t.Fatalf("packet=%s", packet)
	}
	play := string(mustRead(t, filepath.Join(inst, "handoff", "playbooks", "P.md")))
	if !strings.Contains(play, "glossary.md") {
		t.Fatalf("CONTEXT not mapped:\n%s", play)
	}
}

func fixedDeps(t *testing.T) handoffcmd.Deps {
	t.Helper()
	return handoffcmd.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)},
		Stdout: io.Discard,
		Stderr: io.Discard,
	}
}

func scaffoldOffline(t *testing.T, parent, name string) string {
	t.Helper()
	dir := filepath.Join(parent, "inst-"+strings.ReplaceAll(name, " ", "-"))
	var stdout, stderr bytes.Buffer
	code := scaffold.Run(scaffold.Options{
		Name:    name,
		Dir:     dir,
		Offline: true,
		Cwd:     parent,
		Argv:    []string{"new", "idea", name, "--offline"},
	}, scaffold.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)},
		Runner: execrun.Real{},
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr.String())
	}
	return dir
}

func toClarified(t *testing.T, deps handoffcmd.Deps, inst string) {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(inst, "mycelium.toml"))
	if err != nil {
		t.Fatal(err)
	}
	m, err := manifest.Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	m.State = "clarified"
	m.Revisit = ""
	m.UpdatedDate = "2026-08-15"
	out, err := manifest.Encode(m)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inst, "mycelium.toml"), out, 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeAcceptedDEC(t *testing.T, inst, id, title, extraBody string) {
	t.Helper()
	dir := filepath.Join(inst, "decisions")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	slug := strings.ToLower(strings.ReplaceAll(title, " ", "-"))
	body := `+++
id = "` + id + `"
title = "` + title + `"
status = "Accepted"
date = "2026-08-15"
owner = "TBD"
+++

# ` + id + ` — ` + title + `

## Context

` + extraBody + `

## Decision

none

## Rationale

none

## Consequences

none

## Alternatives Considered

none
`
	path := filepath.Join(dir, id+"-"+slug+".md")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := metadata.Parse([]byte(body)); err != nil {
		t.Fatal(err)
	}
}

func writeOQ(t *testing.T, inst, id, agreement string) {
	t.Helper()
	dir := filepath.Join(inst, "questions")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := `+++
id = "` + id + `"
title = "Open"
agreement = "` + agreement + `"
date = "2026-08-15"
owner = "TBD"
+++

# ` + id + `

## Positions

none
`
	if err := os.WriteFile(filepath.Join(dir, id+"-open.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeEVD(t *testing.T, inst, id, slug string) {
	t.Helper()
	dir := filepath.Join(inst, "evidence")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := `+++
id = "` + id + `"
title = "Evidence"
date = "2026-08-15"
owner = "TBD"
+++

# ` + id + `

## Summary

none
`
	if err := os.WriteFile(filepath.Join(dir, id+"-"+slug+".md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
