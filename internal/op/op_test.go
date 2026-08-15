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
	// Begin should flock the existing stale file (same inode), not unlink first.
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

func TestResumeMarksDoneWhenDestExistsSourceGone(t *testing.T) {
	root := t.TempDir()
	s1, err := op.Begin(root, intent("DEC-001"), fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	if err := s1.Stage([]op.Staged{
		{RelTo: "decisions/DEC-001.md", Content: []byte("artifact")},
		{RelTo: "log.md", Content: []byte("log")},
	}); err != nil {
		t.Fatal(err)
	}
	j := s1.Journal()
	from0 := filepath.Join(root, filepath.FromSlash(j.Renames[0].From))
	to0 := filepath.Join(root, filepath.FromSlash(j.Renames[0].To))
	if err := os.MkdirAll(filepath.Dir(to0), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(from0, to0); err != nil {
		t.Fatal(err)
	}
	j.Renames[0].Done = false
	if err := journal.Save(root, j); err != nil {
		t.Fatal(err)
	}
	if err := s1.Close(); err != nil {
		t.Fatal(err)
	}

	s2, err := op.Begin(root, intent(""), fixedNow().Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if err := s2.Commit(); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(to0)
	if err != nil || string(b) != "artifact" {
		t.Fatalf("dest = %q err=%v", b, err)
	}
	if _, err := os.Stat(filepath.Join(root, "log.md")); err != nil {
		t.Fatal("second rename should complete on resume")
	}
	if _, err := journal.Load(root); !errors.Is(err, journal.ErrNotExist) {
		t.Fatal("journal should finish")
	}
}

func TestCommitReplacesExistingDestination(t *testing.T) {
	root := t.TempDir()
	s, err := op.Begin(root, intent("DEC-001"), fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "log.md"), []byte("old\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "mycelium.toml"), []byte("old-m\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := s.Stage([]op.Staged{
		{RelTo: "log.md", Content: []byte("new\n")},
		{RelTo: "mycelium.toml", Content: []byte("new-m\n")},
	}); err != nil {
		t.Fatal(err)
	}
	if err := s.Commit(); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(root, "log.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "new\n" {
		t.Fatalf("log.md got %q", b)
	}
	b, err = os.ReadFile(filepath.Join(root, "mycelium.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "new-m\n" {
		t.Fatalf("mycelium.toml got %q", b)
	}
}

func TestCommitRefusesClobberREADME(t *testing.T) {
	root := t.TempDir()
	s, err := op.Begin(root, intent("DEC-001"), fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	readme := filepath.Join(root, "README.md")
	if err := os.WriteFile(readme, []byte("KEEP-README-BYTES\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := s.Stage([]op.Staged{{RelTo: "README.md", Content: []byte("PWNED-README\n")}}); err != nil {
		t.Fatal(err)
	}
	err = s.Commit()
	if !errors.Is(err, op.ErrRenameConflict) {
		t.Fatalf("want ErrRenameConflict, got %v", err)
	}
	b, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "KEEP-README-BYTES\n" {
		t.Fatalf("README.md clobbered: %q", b)
	}
	_ = s.Close()

	// Resume leftover journal: same refuse, README unchanged.
	s2, err := op.Begin(root, intent("DEC-001"), fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	err = s2.Commit()
	if !errors.Is(err, op.ErrRenameConflict) {
		t.Fatalf("resume want ErrRenameConflict, got %v", err)
	}
	b, err = os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "KEEP-README-BYTES\n" {
		t.Fatalf("README.md clobbered on resume: %q", b)
	}
	_ = s2.Close()
}

func TestCommitRefusesClobberExistingDEC(t *testing.T) {
	root := t.TempDir()
	s, err := op.Begin(root, intent("DEC-001"), fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(root, "decisions", "DEC-001-keep-this.md")
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dest, []byte("KEEP-THIS-DEC\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := s.Stage([]op.Staged{{RelTo: "decisions/DEC-001-keep-this.md", Content: []byte("PWNED-DEC\n")}}); err != nil {
		t.Fatal(err)
	}
	err = s.Commit()
	if !errors.Is(err, op.ErrRenameConflict) {
		t.Fatalf("want ErrRenameConflict, got %v", err)
	}
	b, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "KEEP-THIS-DEC\n" {
		t.Fatalf("DEC clobbered: %q", b)
	}
	_ = s.Close()
}

func TestCommitSupersedeAllowsReplaceArtifact(t *testing.T) {
	root := t.TempDir()
	s, err := op.Begin(root, op.Intent{
		Op:         "supersede",
		Type:       "decision",
		Title:      "DEC-001 -> DEC-002",
		OriginalID: "DEC-001",
		LogLine:    "2026-08-15\tsupersede\tDEC-001\tDEC-001 -> DEC-002",
		Argv:       []string{"supersede", "DEC-001", "--by", "DEC-002"},
		OpID:       "test-supersede",
	}, fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	oldPath := filepath.Join(root, "decisions", "DEC-001-old.md")
	newPath := filepath.Join(root, "decisions", "DEC-002-new.md")
	if err := os.MkdirAll(filepath.Dir(oldPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(oldPath, []byte("OLD-BEFORE\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(newPath, []byte("NEW-BEFORE\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "log.md"), []byte("LOG-BEFORE\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "index.md"), []byte("INDEX-BEFORE\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "mycelium.toml"), []byte("MAN-BEFORE\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := s.Stage([]op.Staged{
		{RelTo: "decisions/DEC-001-old.md", Content: []byte("OLD-AFTER\n")},
		{RelTo: "decisions/DEC-002-new.md", Content: []byte("NEW-AFTER\n")},
		{RelTo: "index.md", Content: []byte("INDEX-AFTER\n")},
		{RelTo: "log.md", Content: []byte("LOG-AFTER\n")},
		{RelTo: "mycelium.toml", Content: []byte("MAN-AFTER\n")},
	}); err != nil {
		t.Fatal(err)
	}
	if err := s.Commit(); err != nil {
		t.Fatalf("supersede commit: %v", err)
	}
	for _, tc := range []struct{ path, want string }{
		{oldPath, "OLD-AFTER\n"},
		{newPath, "NEW-AFTER\n"},
		{filepath.Join(root, "log.md"), "LOG-AFTER\n"},
	} {
		b, err := os.ReadFile(tc.path)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != tc.want {
			t.Fatalf("%s = %q want %q", tc.path, b, tc.want)
		}
	}
}

func TestCommitHandoffAllowsReplaceUnderHandoff(t *testing.T) {
	root := t.TempDir()
	s, err := op.Begin(root, op.Intent{
		Op:      "handoff",
		Title:   "clarified -> handed-off",
		LogLine: "2026-08-15\thandoff\tHO-001\tclarified -> handed-off",
		Argv:    []string{"handoff"},
		OpID:    "test-handoff",
	}, fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	packet := filepath.Join(root, "handoff", "PACKET.md")
	if err := os.MkdirAll(filepath.Dir(packet), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(packet, []byte("BEFORE\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "log.md"), []byte("LOG-BEFORE\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "index.md"), []byte("INDEX-BEFORE\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "mycelium.toml"), []byte("MAN-BEFORE\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := s.Stage([]op.Staged{
		{RelTo: "handoff/PACKET.md", Content: []byte("AFTER\n")},
		{RelTo: "index.md", Content: []byte("INDEX-AFTER\n")},
		{RelTo: "log.md", Content: []byte("LOG-AFTER\n")},
		{RelTo: "mycelium.toml", Content: []byte("MAN-AFTER\n")},
	}); err != nil {
		t.Fatal(err)
	}
	if err := s.Commit(); err != nil {
		t.Fatalf("handoff commit: %v", err)
	}
	b, err := os.ReadFile(packet)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "AFTER\n" {
		t.Fatalf("packet=%q", b)
	}
}

func TestAbortRefusesLiveLock(t *testing.T) {
	root := t.TempDir()
	s, err := op.Begin(root, intent("DEC-001"), fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Stage([]op.Staged{
		{RelTo: "decisions/DEC-001.md", Content: []byte("a")},
		{RelTo: "log.md", Content: []byte("l")},
	}); err != nil {
		t.Fatal(err)
	}
	stageFile := filepath.Join(root, filepath.FromSlash(s.Journal().Renames[1].From))
	if _, err := os.Stat(stageFile); err != nil {
		t.Fatal(err)
	}
	err = op.Abort(root)
	if !errors.Is(err, op.ErrLocked) {
		t.Fatalf("want ErrLocked, got %v", err)
	}
	if _, err := journal.Load(root); err != nil {
		t.Fatal("journal must remain under live lock")
	}
	if _, err := os.Stat(stageFile); err != nil {
		t.Fatal("staged file must remain under live lock")
	}
	info, err := lock.Inspect(root)
	if err != nil {
		t.Fatal(err)
	}
	if info.State != lock.Live {
		t.Fatalf("lock state = %v", info.State)
	}
	_ = s.Close()
}

func TestStageRejectsPathEscape(t *testing.T) {
	root := t.TempDir()
	s, err := op.Begin(root, intent("DEC-001"), fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	err = s.Stage([]op.Staged{{RelTo: "../outside.md", Content: []byte("x")}})
	if !errors.Is(err, op.ErrPathEscape) {
		t.Fatalf("want ErrPathEscape, got %v", err)
	}
	err = s.Stage([]op.Staged{{RelTo: "/abs.md", Content: []byte("x")}})
	if !errors.Is(err, op.ErrPathEscape) {
		t.Fatalf("want ErrPathEscape for abs, got %v", err)
	}
}

func TestBeginRejectsEscapingOpID(t *testing.T) {
	root := t.TempDir()
	bad := intent("DEC-001")
	bad.OpID = "../escape"
	_, err := op.Begin(root, bad, fixedNow())
	if !errors.Is(err, op.ErrPathEscape) {
		t.Fatalf("want ErrPathEscape, got %v", err)
	}
}

func TestAbortAndDetectPruneOrphanStages(t *testing.T) {
	root := t.TempDir()
	orphan := filepath.Join(root, ".mycelium", "stage", "orphan-op")
	if err := os.MkdirAll(orphan, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(orphan, "tmp"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, _, err := op.Detect(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(orphan); !os.IsNotExist(err) {
		t.Fatal("Detect must prune orphan stage dirs")
	}

	if err := os.MkdirAll(orphan, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(orphan, "tmp"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(lock.Path(root), []byte("pid=2147483646\nstarted=2020-01-01T00:00:00Z\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := op.Abort(root); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(orphan); !os.IsNotExist(err) {
		t.Fatal("Abort must prune orphan stage dirs")
	}
}

func TestDetectKeepsJournalStagePrunesOrphans(t *testing.T) {
	root := t.TempDir()
	s, err := op.Begin(root, intent("DEC-001"), fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Stage([]op.Staged{{RelTo: "decisions/DEC-001.md", Content: []byte("a")}}); err != nil {
		t.Fatal(err)
	}
	keep := s.StageDir()
	orphan := filepath.Join(root, ".mycelium", "stage", "other-op")
	if err := os.MkdirAll(orphan, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	hasJ, _, err := op.Detect(root)
	if err != nil || !hasJ {
		t.Fatalf("hasJournal=%v err=%v", hasJ, err)
	}
	if _, err := os.Stat(keep); err != nil {
		t.Fatal("Detect must keep journal staged_dir")
	}
	if _, err := os.Stat(orphan); !os.IsNotExist(err) {
		t.Fatal("Detect must prune other stage dirs")
	}
}

func TestStagedDirDotRefused(t *testing.T) {
	root := t.TempDir()
	marker := filepath.Join(root, "KEEP.md")
	if err := os.WriteFile(marker, []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	j := &journal.Journal{
		SchemaVersion: 1,
		Op:            "scaffold",
		StartedAt:     "2026-08-15T12:00:00Z",
		StagedDir:     ".",
		Argv:          []string{"new", "idea", "X", "--offline"},
	}
	if err := journal.Save(root, j); err != nil {
		t.Fatal(err)
	}
	_, err := op.Begin(root, op.Intent{
		Op:    "scaffold",
		Title: "",
		Argv:  []string{"new", "idea", "X", "--offline"},
	}, fixedNow())
	if !errors.Is(err, op.ErrPathEscape) {
		t.Fatalf("Begin want ErrPathEscape for staged_dir=\".\", got %v", err)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Fatal("instance must survive refused Begin")
	}
}

func TestAbortStagedDirDotDoesNotWipeInstance(t *testing.T) {
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
		Renames: []journal.Rename{
			{From: "staged.tmp", To: "dest.md", Done: false},
		},
	}
	if err := journal.Save(root, j); err != nil {
		t.Fatal(err)
	}
	if err := op.Abort(root); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Fatal("Abort must not RemoveAll instance root when staged_dir is \".\"")
	}
	if _, err := journal.Load(root); !errors.Is(err, journal.ErrNotExist) {
		t.Fatal("journal should still be cleared")
	}
}

func TestAbortDoesNotRemoveInstancePathsInFrom(t *testing.T) {
	root := t.TempDir()
	readme := filepath.Join(root, "README.md")
	manifest := filepath.Join(root, "mycelium.toml")
	if err := os.WriteFile(readme, []byte("readme"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(manifest, []byte("toml"), 0o644); err != nil {
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
	if err := op.Abort(root); err != nil {
		t.Fatal(err)
	}
	if b, err := os.ReadFile(readme); err != nil || string(b) != "readme" {
		t.Fatalf("README.md must survive Abort, got %q err=%v", b, err)
	}
	if b, err := os.ReadFile(manifest); err != nil || string(b) != "toml" {
		t.Fatalf("mycelium.toml must survive Abort, got %q err=%v", b, err)
	}
	if _, err := journal.Load(root); !errors.Is(err, journal.ErrNotExist) {
		t.Fatal("journal should be gone")
	}
}

func TestDetectSkipsPruneWhenLockLive(t *testing.T) {
	root := t.TempDir()
	s, err := op.Begin(root, intent("DEC-001"), fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if err := s.Stage([]op.Staged{{RelTo: "decisions/DEC-001.md", Content: []byte("a")}}); err != nil {
		t.Fatal(err)
	}
	keep := s.StageDir()
	orphan := filepath.Join(root, ".mycelium", "stage", "other-op")
	if err := os.MkdirAll(orphan, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(orphan, "tmp"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	hasJ, stale, err := op.Detect(root)
	if err != nil {
		t.Fatal(err)
	}
	if !hasJ || stale {
		t.Fatalf("hasJournal=%v stale=%v", hasJ, stale)
	}
	if _, err := os.Stat(orphan); err != nil {
		t.Fatal("Detect must not prune orphan stages under a live lock")
	}
	if _, err := os.Stat(keep); err != nil {
		t.Fatal("journal staged_dir must remain")
	}
}
