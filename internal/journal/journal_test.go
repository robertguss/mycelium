package journal_test

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/robertguss/mycelium/internal/journal"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	root := t.TempDir()
	j := &journal.Journal{
		SchemaVersion: 1,
		Op:            "new",
		Title:         "T",
		OriginalID:    "DEC-001",
		StartedAt:     "2026-08-15T12:00:00Z",
		StagedDir:     ".mycelium/stage/op1",
		Renames: []journal.Rename{
			{From: ".mycelium/stage/op1/a.md", To: "decisions/DEC-001.md", Done: false},
		},
		LogLine: "2026-08-15\tnew\tDEC-001\tT",
		Argv:    []string{"new", "decision", "T"},
	}
	j.SetType("decision")
	if err := journal.Save(root, j); err != nil {
		t.Fatal(err)
	}
	got, err := journal.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if got.OriginalID != "DEC-001" || got.TypeString() != "decision" {
		t.Fatalf("got %+v", got)
	}
	if got.Renames[0].To != "decisions/DEC-001.md" {
		t.Fatalf("rename = %+v", got.Renames[0])
	}
}

func TestTypeNull(t *testing.T) {
	root := t.TempDir()
	j := &journal.Journal{SchemaVersion: 1, Op: "scaffold", Argv: []string{"new", "idea", "x"}}
	if err := journal.Save(root, j); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(journal.Path(root))
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		t.Fatal(err)
	}
	if raw["type"] != nil {
		t.Fatalf("type should be null, got %#v", raw["type"])
	}
}

func TestLoadMissing(t *testing.T) {
	_, err := journal.Load(t.TempDir())
	if !errors.Is(err, journal.ErrNotExist) {
		t.Fatalf("got %v", err)
	}
}

func TestMatches(t *testing.T) {
	j := &journal.Journal{
		Op:    "new",
		Title: "T",
		Argv:  []string{"new", "decision", "T"},
	}
	j.SetType("decision")
	if !journal.Matches(j, []string{"new", "decision", "T"}, "new", "decision", "T") {
		t.Fatal("argv match")
	}
	if !journal.Matches(j, []string{"other"}, "new", "decision", "T") {
		t.Fatal("op+type+title match")
	}
	if journal.Matches(j, []string{"other"}, "tier", "", "") {
		t.Fatal("should not match different op")
	}
}
