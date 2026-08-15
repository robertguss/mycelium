package scaffold

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/robertguss/mycelium/internal/op"
)

func TestPackPresent(t *testing.T) {
	if packPresent(nil, "program/packs/council/") {
		t.Fatal("empty files must be absent")
	}
	if packPresent([]op.Staged{{RelTo: "program/README.md"}}, "program/packs/council/") {
		t.Fatal("unrelated file must not count")
	}
	if !packPresent([]op.Staged{{RelTo: "program/packs/council/README.md"}}, "program/packs/council/") {
		t.Fatal("prefix match must count")
	}
}

func TestOnlyMyceliumLeftoversEdges(t *testing.T) {
	root := t.TempDir()
	ok, err := onlyMyceliumLeftovers(filepath.Join(root, "missing"))
	if err == nil || ok {
		t.Fatalf("missing dir: ok=%v err=%v", ok, err)
	}
	empty := filepath.Join(root, "empty")
	if err := os.Mkdir(empty, 0o755); err != nil {
		t.Fatal(err)
	}
	ok, err = onlyMyceliumLeftovers(empty)
	if err != nil || ok {
		t.Fatalf("empty dir: ok=%v err=%v", ok, err)
	}
	two := filepath.Join(root, "two")
	if err := os.MkdirAll(filepath.Join(two, ".mycelium"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(two, "extra"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err = onlyMyceliumLeftovers(two)
	if err != nil || ok {
		t.Fatalf("two entries: ok=%v err=%v", ok, err)
	}
	fileOnly := filepath.Join(root, "file-only")
	if err := os.Mkdir(fileOnly, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(fileOnly, "not-dir"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err = onlyMyceliumLeftovers(fileOnly)
	if err != nil || ok {
		t.Fatalf("single file: ok=%v err=%v", ok, err)
	}
}

func TestRollbackOrClose(t *testing.T) {
	root := t.TempDir()
	now := time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)
	sess, err := op.Begin(root, op.Intent{
		Op: "scaffold", Title: "x", LogLine: "2026-08-14 scaffold - x", Argv: []string{"new", "idea", "x"},
	}, now)
	if err != nil {
		t.Fatal(err)
	}
	rollbackOrClose(sess)

	root2 := t.TempDir()
	sess2, err := op.Begin(root2, op.Intent{
		Op: "scaffold", Title: "y", LogLine: "2026-08-14 scaffold - y", Argv: []string{"new", "idea", "y"},
	}, now)
	if err != nil {
		t.Fatal(err)
	}
	if err := sess2.Stage([]op.Staged{{RelTo: "README.md", Content: []byte("# y\n")}}); err != nil {
		t.Fatal(err)
	}
	if err := sess2.Commit(); err != nil {
		t.Fatal(err)
	}
	rollbackOrClose(sess2) // ErrPartialCommit → Close
}
