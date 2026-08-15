package op_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/robertguss/mycelium/internal/journal"
	"github.com/robertguss/mycelium/internal/lock"
	"github.com/robertguss/mycelium/internal/op"
)

func fixedNow() time.Time {
	return time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC)
}

func intent(id string) op.Intent {
	return op.Intent{
		Op:         "new",
		Type:       "decision",
		Title:      "Ship It",
		OriginalID: id,
		LogLine:    "2026-08-15\tnew\tDEC-001\tShip It",
		Argv:       []string{"new", "decision", "Ship It"},
		OpID:       "test-op",
	}
}

func TestCrashBeforeRenameLeavesNothing(t *testing.T) {
	root := t.TempDir()
	s, err := op.Begin(root, intent("DEC-001"), fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Stage([]op.Staged{
		{RelTo: "decisions/DEC-001.md", Content: []byte("artifact")},
		{RelTo: "log.md", Content: []byte("log")},
		{RelTo: "mycelium.toml", Content: []byte("manifest")},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := journal.Load(root); err != nil {
		t.Fatal("journal should exist after stage")
	}
	if err := s.Rollback(); err != nil {
		t.Fatal(err)
	}
	if _, err := journal.Load(root); !errors.Is(err, journal.ErrNotExist) {
		t.Fatalf("journal should be gone, got %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "decisions", "DEC-001.md")); !os.IsNotExist(err) {
		t.Fatal("no committed artifact")
	}
	if _, err := os.Stat(filepath.Join(root, ".mycelium", "stage")); !os.IsNotExist(err) {
		t.Fatal("stage dir should be gone")
	}
	info, err := lock.Inspect(root)
	if err != nil {
		t.Fatal(err)
	}
	if info.State != lock.Absent {
		t.Fatalf("lock state = %v", info.State)
	}
}

func TestCrashAfterPartialRenameLeavesJournal(t *testing.T) {
	root := t.TempDir()
	s, err := op.Begin(root, intent("DEC-001"), fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Stage([]op.Staged{
		{RelTo: "decisions/DEC-001.md", Content: []byte("artifact")},
		{RelTo: "log.md", Content: []byte("log")},
		{RelTo: "mycelium.toml", Content: []byte("manifest")},
	}); err != nil {
		t.Fatal(err)
	}
	if err := s.CommitPartial(1); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	j, err := journal.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if !j.Renames[0].Done || j.Renames[1].Done {
		t.Fatalf("renames = %+v", j.Renames)
	}
	if _, err := os.Stat(filepath.Join(root, "decisions", "DEC-001.md")); err != nil {
		t.Fatal("first rename should have landed")
	}
	if _, err := os.Stat(filepath.Join(root, "log.md")); !os.IsNotExist(err) {
		t.Fatal("second rename must not have landed")
	}
}

func TestResumeUsesOriginalID(t *testing.T) {
	root := t.TempDir()
	s1, err := op.Begin(root, intent("DEC-001"), fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	if err := s1.Stage([]op.Staged{
		{RelTo: "decisions/DEC-001.md", Content: []byte("a")},
		{RelTo: "log.md", Content: []byte("l")},
	}); err != nil {
		t.Fatal(err)
	}
	if err := s1.CommitPartial(1); err != nil {
		t.Fatal(err)
	}
	if err := s1.Close(); err != nil {
		t.Fatal(err)
	}

	// Resume with empty OriginalID in intent — must reuse DEC-001.
	resume := intent("")
	s2, err := op.Begin(root, resume, fixedNow().Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if s2.OriginalID() != "DEC-001" {
		t.Fatalf("original_id = %q, want DEC-001", s2.OriginalID())
	}
	if err := s2.Commit(); err != nil {
		t.Fatal(err)
	}
	if _, err := journal.Load(root); !errors.Is(err, journal.ErrNotExist) {
		t.Fatal("journal should be cleared after resume commit")
	}
	b, err := os.ReadFile(filepath.Join(root, "log.md"))
	if err != nil || string(b) != "l" {
		t.Fatalf("log = %q err=%v", b, err)
	}
}

func TestStaleLockDetected(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".mycelium"), 0o755); err != nil {
		t.Fatal(err)
	}
	dead := []byte("pid=2147483646\nstarted=2020-01-01T00:00:00Z\n")
	if err := os.WriteFile(lock.Path(root), dead, 0o644); err != nil {
		t.Fatal(err)
	}
	_, stale, err := op.Detect(root)
	if err != nil {
		t.Fatal(err)
	}
	if !stale {
		t.Fatal("expected stale lock")
	}
	// Begin should clear stale and acquire.
	s, err := op.Begin(root, intent("DEC-001"), fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	info, err := lock.Inspect(root)
	if err != nil {
		t.Fatal(err)
	}
	if info.State != lock.Live {
		t.Fatalf("state = %v", info.State)
	}
}

func TestLiveLockAcquireRefuses(t *testing.T) {
	root := t.TempDir()
	s, err := op.Begin(root, intent("DEC-001"), fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	_, err = op.Begin(root, intent("DEC-002"), fixedNow())
	if !errors.Is(err, op.ErrLocked) {
		t.Fatalf("want ErrLocked, got %v", err)
	}
}

func TestAbortJournal(t *testing.T) {
	root := t.TempDir()
	s, err := op.Begin(root, intent("DEC-001"), fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Stage([]op.Staged{
		{RelTo: "decisions/DEC-001.md", Content: []byte("a")},
		{RelTo: "log.md", Content: []byte("l")},
		{RelTo: "mycelium.toml", Content: []byte("m")},
	}); err != nil {
		t.Fatal(err)
	}
	if err := s.CommitPartial(1); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	// Leave a stale lock file too.
	if err := os.WriteFile(lock.Path(root), []byte("pid=2147483646\nstarted=2020-01-01T00:00:00Z\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := op.Abort(root); err != nil {
		t.Fatal(err)
	}
	if _, err := journal.Load(root); !errors.Is(err, journal.ErrNotExist) {
		t.Fatal("journal should be gone")
	}
	if _, err := os.Stat(filepath.Join(root, "decisions", "DEC-001.md")); err != nil {
		t.Fatal("already-renamed artifact must survive abort")
	}
	if _, err := os.Stat(filepath.Join(root, ".mycelium", "stage")); !os.IsNotExist(err) {
		t.Fatal("staged temps should be gone")
	}
	info, err := lock.Inspect(root)
	if err != nil {
		t.Fatal(err)
	}
	if info.State != lock.Absent {
		t.Fatalf("stale lock should be removed, got %v", info.State)
	}
	if _, err := os.Stat(filepath.Join(root, "log.md")); !os.IsNotExist(err) {
		t.Fatal("unrenamed dest must not appear")
	}
}

func TestLeftoverJournalDifferentArgvRefuses(t *testing.T) {
	root := t.TempDir()
	s, err := op.Begin(root, intent("DEC-001"), fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Stage([]op.Staged{{RelTo: "decisions/DEC-001.md", Content: []byte("a")}}); err != nil {
		t.Fatal(err)
	}
	if err := s.CommitPartial(0); err != nil {
		t.Fatal(err)
	}
	// CommitPartial(0) still needs journal on disk — Stage already saved it.
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	other := op.Intent{
		Op:    "tier",
		Title: "",
		Argv:  []string{"tier", "standard"},
	}
	_, err = op.Begin(root, other, fixedNow())
	if !errors.Is(err, op.ErrJournalMismatch) {
		t.Fatalf("want ErrJournalMismatch, got %v", err)
	}
}

func TestAbortNothing(t *testing.T) {
	err := op.Abort(t.TempDir())
	if !errors.Is(err, op.ErrNothingToAbort) {
		t.Fatalf("got %v", err)
	}
}
