package handoff_test

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/robertguss/mycelium/internal/check"
	"github.com/robertguss/mycelium/internal/cli"
	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/execrun"
	"github.com/robertguss/mycelium/internal/handoff"
	"github.com/robertguss/mycelium/internal/scaffold"
)

func appendixBPacket() string {
	return `+++
id = "HO-001"
date = "2026-08-15"
implementation_system = "pstack/poteto"
time_budget = "30m"
+++

# Handoff packet

## Framing

Bounded target: implement Add(a, b int) int in add.go.

## Locked decisions

- DEC-001 (Accepted) — see decisions/

## Glossary

See glossary.md (copied from CONTEXT.md).

## Open questions

none

## Evidence summary

See evidence/SUMMARY.md.

## Implementation playbooks

See playbooks/PLAYBOOK.md.

## Implementation system

pstack/poteto (manual is the floor).

## Time budget

30m

## Acceptance

See acceptance/. Executable tests for Add.
`
}

func appendixBTree() fstest.MapFS {
	return fstest.MapFS{
		"PACKET.md":   &fstest.MapFile{Data: []byte(appendixBPacket())},
		"glossary.md": &fstest.MapFile{Data: []byte("term: Add\n")},
		"decisions":   &fstest.MapFile{Mode: 0o755 | os.ModeDir},
		"decisions/DEC-001-add-signature.md": &fstest.MapFile{Data: []byte(`+++
id = "DEC-001"
title = "Add signature"
status = "Accepted"
date = "2026-08-15"
owner = "TBD"
+++

# DEC-001

## Context

Add(a, b int) int
`)},
		"questions":           &fstest.MapFile{Mode: 0o755 | os.ModeDir},
		"evidence":            &fstest.MapFile{Mode: 0o755 | os.ModeDir},
		"evidence/SUMMARY.md": &fstest.MapFile{Data: []byte("none\n")},
		"playbooks":           &fstest.MapFile{Mode: 0o755 | os.ModeDir},
		"playbooks/PLAYBOOK.md": &fstest.MapFile{Data: []byte(`# Playbook

## Target

implement Add in add.go

## Steps

1. write Add

## Done

tests green
`)},
		"acceptance":           &fstest.MapFile{Mode: 0o755 | os.ModeDir},
		"acceptance/README.md": &fstest.MapFile{Data: []byte("tests here\n")},
	}
}

func TestAppendixBStructurePass(t *testing.T) {
	t.Parallel()
	findings := handoff.Check(appendixBTree())
	if len(findings) != 0 {
		t.Fatalf("want pass, got %v", findings)
	}
}

func TestMissingH2Fails(t *testing.T) {
	t.Parallel()
	body := strings.Replace(appendixBPacket(), "## Acceptance\n\nSee acceptance/. Executable tests for Add.\n", "", 1)
	fsys := appendixBTree()
	fsys["PACKET.md"] = &fstest.MapFile{Data: []byte(body)}
	findings := handoff.Check(fsys)
	if !hasWhat(findings, `missing required H2 "Acceptance"`) {
		t.Fatalf("want missing Acceptance, got %v", findings)
	}
}

func TestMissingTimeBudgetFails(t *testing.T) {
	t.Parallel()
	body := strings.Replace(appendixBPacket(), "time_budget = \"30m\"\n", "", 1)
	findings := handoff.ValidatePacketBytes([]byte(body))
	if !hasWhat(findings, "missing required front matter key time_budget") {
		t.Fatalf("want missing time_budget, got %v", findings)
	}
}

func TestOutsideDecisionLinkFails(t *testing.T) {
	t.Parallel()
	body := strings.Replace(
		appendixBPacket(),
		"- DEC-001 (Accepted) — see decisions/",
		"- [DEC-001](../decisions/DEC-001-add-signature.md)",
		1,
	)
	fsys := appendixBTree()
	fsys["PACKET.md"] = &fstest.MapFile{Data: []byte(body)}
	findings := handoff.Check(fsys)
	if !hasWhat(findings, "links outside handoff/") {
		t.Fatalf("want outside-link fail, got %v", findings)
	}
}

func TestInPacketDECResolves(t *testing.T) {
	t.Parallel()
	body := strings.Replace(
		appendixBPacket(),
		"- DEC-001 (Accepted) — see decisions/",
		"- [DEC-001](decisions/DEC-001-add-signature.md)",
		1,
	)
	fsys := appendixBTree()
	fsys["PACKET.md"] = &fstest.MapFile{Data: []byte(body)}
	findings := handoff.Check(fsys)
	if len(findings) != 0 {
		t.Fatalf("want in-packet link pass, got %v", findings)
	}
}

func TestFrontMatterTable(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		mutate  func(string) string
		wantSub string
	}{
		{
			name: "bad-impl-system",
			mutate: func(s string) string {
				return strings.Replace(s, `implementation_system = "pstack/poteto"`, `implementation_system = "cursor"`, 1)
			},
			wantSub: "implementation_system",
		},
		{
			name: "bad-time-budget",
			mutate: func(s string) string {
				return strings.Replace(s, `time_budget = "30m"`, `time_budget = "30minutes"`, 1)
			},
			wantSub: "time_budget",
		},
		{
			name: "bad-date",
			mutate: func(s string) string {
				return strings.Replace(s, `date = "2026-08-15"`, `date = "08-15-2026"`, 1)
			},
			wantSub: "date must be YYYY-MM-DD",
		},
		{
			name: "wrong-id",
			mutate: func(s string) string {
				return strings.Replace(s, `id = "HO-001"`, `id = "HO-002"`, 1)
			},
			wantSub: "must be HO-001",
		},
		{
			name: "missing-id",
			mutate: func(s string) string {
				return strings.Replace(s, `id = "HO-001"`+"\n", "", 1)
			},
			wantSub: "missing required front matter key id",
		},
		{
			name: "manual-ok",
			mutate: func(s string) string {
				return strings.Replace(s, `implementation_system = "pstack/poteto"`, `implementation_system = "manual"`, 1)
			},
			wantSub: "",
		},
		{
			name: "unknown-key-ok",
			mutate: func(s string) string {
				return strings.Replace(s, `time_budget = "30m"`, "time_budget = \"30m\"\nextra = \"ok\"", 1)
			},
			wantSub: "",
		},
		{
			name: "hour-budget-ok",
			mutate: func(s string) string {
				return strings.Replace(s, `time_budget = "30m"`, `time_budget = "2h"`, 1)
			},
			wantSub: "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			findings := handoff.ValidatePacketBytes([]byte(tc.mutate(appendixBPacket())))
			if tc.wantSub == "" {
				if len(findings) != 0 {
					t.Fatalf("want pass, got %v", findings)
				}
				return
			}
			if !hasWhat(findings, tc.wantSub) {
				t.Fatalf("want substr %q in %v", tc.wantSub, findings)
			}
		})
	}
}

func TestH2OrderRequired(t *testing.T) {
	t.Parallel()
	body := appendixBPacket()
	body = strings.Replace(body, "## Framing\n\nBounded target: implement Add(a, b int) int in add.go.\n\n", "", 1)
	body = strings.Replace(body, "## Acceptance\n\nSee acceptance/. Executable tests for Add.\n",
		"## Acceptance\n\nSee acceptance/.\n\n## Framing\n\nlate framing\n", 1)
	findings := handoff.ValidatePacketBytes([]byte(body))
	if !hasWhat(findings, "out of order") {
		t.Fatalf("want order fail, got %v", findings)
	}
}

func TestMissingDECCopyFails(t *testing.T) {
	t.Parallel()
	fsys := appendixBTree()
	delete(fsys, "decisions/DEC-001-add-signature.md")
	findings := handoff.Check(fsys)
	if !hasWhat(findings, "no copy under handoff/decisions/") {
		t.Fatalf("want missing DEC copy, got %v", findings)
	}
}

func TestMissingGlossaryFails(t *testing.T) {
	t.Parallel()
	fsys := appendixBTree()
	delete(fsys, "glossary.md")
	findings := handoff.Check(fsys)
	if !hasWhat(findings, "glossary.md missing") {
		t.Fatalf("want glossary missing, got %v", findings)
	}
}

func TestMissingPacketFails(t *testing.T) {
	t.Parallel()
	fsys := fstest.MapFS{
		"decisions": &fstest.MapFile{Mode: 0o755 | os.ModeDir},
	}
	findings := handoff.Check(fsys)
	if !hasWhat(findings, "PACKET.md missing") {
		t.Fatalf("want PACKET missing, got %v", findings)
	}
}

func TestPlaybookOutsideLinkFails(t *testing.T) {
	t.Parallel()
	fsys := appendixBTree()
	fsys["playbooks/PLAYBOOK.md"] = &fstest.MapFile{Data: []byte(`# Playbook

## Target

see [DEC](../../decisions/DEC-001-add-signature.md)

## Steps

none

## Done

none
`)}
	findings := handoff.Check(fsys)
	if !hasWhat(findings, "links outside handoff/") {
		t.Fatalf("want playbook outside-link fail, got %v", findings)
	}
}

func TestPlaybookInPacketRelativeLinkOK(t *testing.T) {
	t.Parallel()
	fsys := appendixBTree()
	fsys["playbooks/PLAYBOOK.md"] = &fstest.MapFile{Data: []byte(`# Playbook

## Target

see [DEC](../decisions/DEC-001-add-signature.md)

## Steps

none

## Done

none
`)}
	findings := handoff.Check(fsys)
	if len(findings) != 0 {
		t.Fatalf("want relative in-packet pass, got %v", findings)
	}
}

func TestRequiredH2s(t *testing.T) {
	t.Parallel()
	got := handoff.RequiredH2s()
	if len(got) != 9 {
		t.Fatalf("len=%d", len(got))
	}
	if got[0] != "Framing" || got[8] != "Acceptance" {
		t.Fatalf("%v", got)
	}
	if handoff.PacketID != "HO-001" {
		t.Fatalf("PacketID=%s", handoff.PacketID)
	}
}

func TestTemplatesExist(t *testing.T) {
	t.Parallel()
	root := repoRoot(t)
	packet, err := os.ReadFile(filepath.Join(root, "program", "templates", "handoff-packet.md"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(packet)
	for _, key := range []string{"id =", "date =", "implementation_system =", "time_budget ="} {
		if !strings.Contains(s, key) {
			t.Fatalf("packet template missing %s", key)
		}
	}
	if !strings.Contains(s, `implementation_system = "pstack/poteto"`) {
		t.Fatal("default implementation_system")
	}
	if !strings.Contains(s, `time_budget = "30m"`) {
		t.Fatal("default time_budget")
	}
	for _, h2 := range handoff.RequiredH2s() {
		if !strings.Contains(s, "## "+h2) {
			t.Fatalf("packet template missing ## %s", h2)
		}
	}
	play, err := os.ReadFile(filepath.Join(root, "program", "templates", "handoff-playbook.md"))
	if err != nil {
		t.Fatal(err)
	}
	ps := string(play)
	for _, h2 := range []string{"Target", "Steps", "Done"} {
		if !strings.Contains(ps, "## "+h2) {
			t.Fatalf("playbook template missing ## %s", h2)
		}
	}
	// No HO schema sidecar — handoff is not mycelium new <type>.
	if _, err := os.Stat(filepath.Join(root, "program", "templates", "handoff-packet.schema.toml")); err == nil {
		t.Fatal("must not add handoff-packet.schema.toml (ID-namespace DSL)")
	}
}

func TestMyceliumHandoffKnownVerb(t *testing.T) {
	t.Parallel()
	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"mycelium", "handoff", "--help"}, &stdout, &stderr, cli.Deps{})
	if code != 0 {
		t.Fatalf("exit %d want 0 stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "mycelium handoff") {
		t.Fatalf("stdout=%q", stdout.String())
	}
}

func TestStateHandedOffStillRefusesWithoutPacket(t *testing.T) {
	t.Parallel()
	parent := t.TempDir()
	inst := scaffoldOffline(t, parent, "Slice2 Refuse")
	var stdout, stderr bytes.Buffer
	code := cli.Run(
		[]string{"mycelium", "state", "handed-off", "--dir", inst},
		&stdout,
		&stderr,
		cli.Deps{Clock: clock.Fixed{T: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)}},
	)
	if code != 1 {
		t.Fatalf("exit %d want 1 stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "handoff packet") {
		t.Fatalf("stderr=%q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "mycelium handoff") {
		t.Fatalf("stderr=%q", stderr.String())
	}
}

func TestStoredHandedOffStillFailsCheck(t *testing.T) {
	t.Parallel()
	root := scaffoldOffline(t, t.TempDir(), "Slice2 Handed Off")
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
		t.Fatal("want stored handed-off FAIL")
	}
	found := false
	for _, f := range r.Findings {
		if strings.Contains(f.What, "handed-off") && strings.Contains(f.What, "handoff packet") {
			found = true
		}
	}
	if !found {
		t.Fatalf("want handoff packet finding, got %v", r.Findings)
	}
}

func hasWhat(fs []handoff.Finding, substr string) bool {
	for _, f := range fs {
		if strings.Contains(f.What, substr) {
			return true
		}
	}
	return false
}

func scaffoldOffline(t *testing.T, parent, name string) string {
	t.Helper()
	dir := filepath.Join(parent, "inst")
	var stdout, stderr bytes.Buffer
	code := scaffold.Run(scaffold.Options{
		Name:    name,
		Dir:     dir,
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
	return dir
}

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("no caller")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}
