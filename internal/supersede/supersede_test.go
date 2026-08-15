package supersede_test

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/cli"
	"github.com/robertguss/mycelium/internal/metadata"
	"github.com/robertguss/mycelium/internal/schema"
	"github.com/robertguss/mycelium/internal/supersede"
)

func TestParseIDIdeaStatesRefuse(t *testing.T) {
	t.Parallel()
	for _, tok := range []string{"spark", "exploring", "simmering", "clarified", "handed-off", "archived"} {
		_, _, r := supersede.ParseID(tok)
		if r == nil {
			t.Fatalf("%s: want refusal", tok)
		}
		want := tok + " is an idea state, not an artifact"
		if r.What != want {
			t.Fatalf("%s what=%q want %q", tok, r.What, want)
		}
		if !strings.Contains(r.Fix, "mycelium state") {
			t.Fatalf("%s fix=%q", tok, r.Fix)
		}
	}
}

func TestParseIDArtifacts(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in, wantID, wantNS string
	}{
		{"DEC-001", "DEC-001", "DEC"},
		{"DEC-1", "DEC-001", "DEC"},
		{"ASM-014", "ASM-014", "ASM"},
		{"EVD-002", "EVD-002", "EVD"},
		{"SPK-003", "SPK-003", "SPK"},
		{"OQ-001", "OQ-001", "OQ"},
		{"PHASE-01", "PHASE-01", "PHASE"},
		{"REC-001", "REC-001", "REC"},
	}
	for _, tc := range cases {
		id, ns, r := supersede.ParseID(tc.in)
		if r != nil {
			t.Fatalf("%s: unexpected refuse %v", tc.in, r)
		}
		if id != tc.wantID || ns != tc.wantNS {
			t.Fatalf("%s: got %s/%s want %s/%s", tc.in, id, ns, tc.wantID, tc.wantNS)
		}
	}
}

func TestEligibility(t *testing.T) {
	t.Parallel()
	yes := []string{"DEC", "ASM", "EVD", "SPK"}
	no := []string{"OQ", "PHASE", "REC", "REQ", "RSK", "FND", "MS", "CMP", "RPT", "RCL"}
	for _, ns := range yes {
		if !supersede.Eligible(ns) {
			t.Fatalf("%s should be eligible", ns)
		}
	}
	for _, ns := range no {
		if supersede.Eligible(ns) {
			t.Fatalf("%s should not be eligible", ns)
		}
	}
}

func TestCheckPairRefuseTable(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name           string
		old, new       string
		oldMeta        map[string]any
		newMeta        map[string]any
		wantWhatSubstr string
		wantFixSubstr  string
	}{
		{
			name: "missing-old", old: "", new: "DEC-002",
			wantWhatSubstr: "supersede requires <OLD-ID> --by <NEW-ID>",
			wantFixSubstr:  "mycelium supersede DEC-001 --by DEC-002",
		},
		{
			name: "self", old: "DEC-001", new: "DEC-001",
			wantWhatSubstr: "cannot supersede an ID with itself",
			wantFixSubstr:  "pass two different IDs",
		},
		{
			name: "self-padded", old: "DEC-1", new: "DEC-001",
			wantWhatSubstr: "cannot supersede an ID with itself",
		},
		{
			name: "namespace", old: "DEC-001", new: "ASM-001",
			wantWhatSubstr: "supersede requires the same namespace (got DEC vs ASM)",
			wantFixSubstr:  "pick two IDs in one namespace",
		},
		{
			name: "oq", old: "OQ-001", new: "OQ-002",
			wantWhatSubstr: "type OQ is not supersedable",
			wantFixSubstr:  "open a new question; do not supersede an OQ",
		},
		{
			name: "phase", old: "PHASE-01", new: "PHASE-02",
			wantWhatSubstr: "type PHASE is not supersedable",
		},
		{
			name: "rec", old: "REC-001", new: "REC-002",
			wantWhatSubstr: "type REC is not supersedable",
		},
		{
			name: "idea-state-old", old: "spark", new: "DEC-001",
			wantWhatSubstr: "spark is an idea state, not an artifact",
		},
		{
			name: "idea-state-new", old: "DEC-001", new: "archived",
			wantWhatSubstr: "archived is an idea state, not an artifact",
		},
		{
			name: "old-already", old: "DEC-001", new: "DEC-003",
			oldMeta:        map[string]any{"status": "Superseded", "superseded_by": "DEC-002"},
			wantWhatSubstr: "DEC-001 is already Superseded by DEC-002",
			wantFixSubstr:  "supersede the current record (DEC-002) --by <newer>",
		},
		{
			name: "new-already", old: "DEC-001", new: "DEC-002",
			newMeta:        map[string]any{"supersedes": "DEC-000"},
			wantWhatSubstr: "DEC-002 already supersedes DEC-000",
			wantFixSubstr:  "one-to-one this phase; pick a different NEW",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := supersede.CheckPair(tc.old, tc.new, tc.oldMeta, tc.newMeta)
			if r == nil {
				t.Fatal("want refusal")
			}
			if !strings.Contains(r.What, tc.wantWhatSubstr) {
				t.Fatalf("what=%q want substr %q", r.What, tc.wantWhatSubstr)
			}
			if tc.wantFixSubstr != "" && !strings.Contains(r.Fix, tc.wantFixSubstr) {
				t.Fatalf("fix=%q want substr %q", r.Fix, tc.wantFixSubstr)
			}
		})
	}
}

func TestCheckPairHappy(t *testing.T) {
	t.Parallel()
	r := supersede.CheckPair("DEC-001", "DEC-002",
		map[string]any{"status": "Accepted"},
		map[string]any{"status": "Accepted"},
	)
	if r != nil {
		t.Fatalf("unexpected refuse: %+v", r)
	}
}

func TestApplyPairDECBytes(t *testing.T) {
	t.Parallel()
	oldIn := []byte(`+++
id = "DEC-001"
title = "Use SQLite"
status = "Accepted"
date = "2026-08-15"
owner = "TBD"
+++

# DEC-001 — Use SQLite

## Context

none
`)
	newIn := []byte(`+++
id = "DEC-002"
title = "Use SQLite with WAL"
status = "Accepted"
date = "2026-08-15"
owner = "TBD"
+++

# DEC-002 — Use SQLite with WAL

## Context

none
`)
	oldOut, newOut, r, err := supersede.ApplyPair(oldIn, newIn, "DEC-001", "DEC-002")
	if err != nil {
		t.Fatal(err)
	}
	if r != nil {
		t.Fatalf("refuse: %+v", r)
	}
	oldDoc, err := metadata.Parse(oldOut)
	if err != nil {
		t.Fatal(err)
	}
	newDoc, err := metadata.Parse(newOut)
	if err != nil {
		t.Fatal(err)
	}
	if oldDoc.Meta["status"] != "Superseded" {
		t.Fatalf("old status=%v", oldDoc.Meta["status"])
	}
	if oldDoc.Meta["superseded_by"] != "DEC-002" {
		t.Fatalf("old superseded_by=%v", oldDoc.Meta["superseded_by"])
	}
	if newDoc.Meta["supersedes"] != "DEC-001" {
		t.Fatalf("new supersedes=%v", newDoc.Meta["supersedes"])
	}
	if newDoc.Meta["status"] != "Accepted" {
		t.Fatalf("new status changed: %v", newDoc.Meta["status"])
	}
	if !bytes.Contains(oldOut, []byte(`status = "Superseded"`)) {
		t.Fatalf("old bytes missing status line:\n%s", oldOut)
	}
	if !bytes.Contains(oldOut, []byte(`superseded_by = "DEC-002"`)) {
		t.Fatalf("old bytes missing superseded_by:\n%s", oldOut)
	}
	if !bytes.Contains(newOut, []byte(`supersedes = "DEC-001"`)) {
		t.Fatalf("new bytes missing supersedes:\n%s", newOut)
	}
	if !strings.Contains(oldDoc.Body, "## Context") {
		t.Fatal("old body lost")
	}
}

func TestApplyPairRefusesIneligible(t *testing.T) {
	t.Parallel()
	doc := []byte(`+++
id = "OQ-001"
title = "Q"
agreement = "open"
date = "2026-08-15"
+++

## Question

x
`)
	_, _, r, err := supersede.ApplyPair(doc, doc, "OQ-001", "OQ-002")
	if err != nil {
		t.Fatal(err)
	}
	if r == nil || !strings.Contains(r.What, "type OQ is not supersedable") {
		t.Fatalf("got %+v", r)
	}
}

func TestASMSchemaParsesSuperseded(t *testing.T) {
	t.Parallel()
	path := filepath.Join(repoRoot(t), "program", "templates", "assumption.schema.toml")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s, err := schema.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, v := range s.Enums["status"] {
		if v == "Superseded" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("ASM status enum missing Superseded: %v", s.Enums["status"])
	}
	body := string(raw)
	if !strings.Contains(body, "superseded_by") || !strings.Contains(body, "supersedes") {
		t.Fatal("ASM schema must name optional superseded_by/supersedes")
	}
}

func TestOptionalLinkKeysNamedOnEligibleSchemas(t *testing.T) {
	t.Parallel()
	root := filepath.Join(repoRoot(t), "program", "templates")
	for _, name := range []string{"decision.schema.toml", "evidence.schema.toml", "spike.schema.toml"} {
		raw, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatal(err)
		}
		body := string(raw)
		if !strings.Contains(body, "superseded_by") || !strings.Contains(body, "supersedes") {
			t.Fatalf("%s must name optional link keys", name)
		}
		s, err := schema.Parse(raw)
		if err != nil {
			t.Fatal(err)
		}
		found := false
		for _, v := range s.Enums["status"] {
			if v == "Superseded" {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("%s status enum missing Superseded", name)
		}
	}
}

func TestQuestionSchemaNoStatusSuperseded(t *testing.T) {
	t.Parallel()
	raw, err := os.ReadFile(filepath.Join(repoRoot(t), "program", "templates", "question.schema.toml"))
	if err != nil {
		t.Fatal(err)
	}
	s, err := schema.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := s.Enums["status"]; ok {
		t.Fatal("question schema must not have status enum")
	}
	body := string(raw)
	if strings.Contains(body, "Superseded") {
		t.Fatal("question schema must not mention Superseded")
	}
	if strings.Contains(body, "superseded_by") || strings.Contains(body, "supersedes") {
		t.Fatal("question schema must not name supersede link keys")
	}
}

func TestTemplatesDoNotEmitLinkKeys(t *testing.T) {
	t.Parallel()
	root := filepath.Join(repoRoot(t), "program", "templates")
	for _, name := range []string{"decision.md", "assumption.md", "evidence.md", "spike.md"} {
		raw, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatal(err)
		}
		body := string(raw)
		if strings.Contains(body, "superseded_by") || strings.Contains(body, "supersedes =") {
			t.Fatalf("%s must not emit superseded_by/supersedes", name)
		}
	}
}

func TestMyceliumSupersedeHelp(t *testing.T) {
	t.Parallel()
	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"mycelium", "supersede", "-h"}, &stdout, &stderr, cli.Deps{})
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "mycelium supersede") {
		t.Fatalf("stdout=%q", stdout.String())
	}
}

func TestIsIdeaState(t *testing.T) {
	t.Parallel()
	if !supersede.IsIdeaState("simmering") {
		t.Fatal("expected true")
	}
	if supersede.IsIdeaState("DEC-001") {
		t.Fatal("expected false")
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("no caller")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}
