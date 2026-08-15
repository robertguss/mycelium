package execrun_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/robertguss/mycelium/internal/execrun"
)

func TestFakeRecordsAndBlocksMissing(t *testing.T) {
	f := &execrun.Fake{Paths: map[string]string{"git": "/bin/git"}}
	if _, err := f.LookPath("git"); err != nil {
		t.Fatal(err)
	}
	if _, err := f.LookPath("gh"); err == nil {
		t.Fatal("expected missing gh")
	}
	res, err := f.Run(context.Background(), "git", []string{"init", "-b", "main"}, execrun.RunOpts{Dir: "/tmp"})
	if err != nil {
		t.Fatal(err)
	}
	if res.ExitCode != 0 {
		t.Fatalf("exit %d", res.ExitCode)
	}
	if !f.Called("git") {
		t.Fatal("git not recorded")
	}
	if f.Called("gh") {
		t.Fatal("gh should not be recorded as run")
	}
}

func TestRecordingWrapsInner(t *testing.T) {
	inner := &execrun.Fake{}
	rec := &execrun.Recording{Inner: inner}
	_, _ = rec.Run(context.Background(), "gh", []string{"auth", "status"}, execrun.RunOpts{})
	if !rec.Called("gh") {
		t.Fatal("expected gh call")
	}
	if !inner.Called("gh") {
		t.Fatal("inner should see gh")
	}
}

func TestRealRunScrubsGIT_DIRWhenDirSet(t *testing.T) {
	root := t.TempDir()
	trap := filepath.Join(root, "trap.git")
	work := filepath.Join(root, "work")
	if err := os.MkdirAll(trap, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(work, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GIT_DIR", trap)
	t.Setenv("GIT_WORK_TREE", root)

	res, err := (execrun.Real{}).Run(
		context.Background(),
		"git",
		[]string{"init", "-b", "main"},
		execrun.RunOpts{Dir: work},
	)
	if err != nil {
		t.Fatal(err)
	}
	if res.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%s", res.ExitCode, res.Stderr)
	}
	if _, err := os.Stat(filepath.Join(work, ".git")); err != nil {
		t.Fatalf("expected .git under Dir: %v", err)
	}
	if _, err := os.Stat(filepath.Join(trap, "HEAD")); err == nil {
		t.Fatal("GIT_DIR trap should not have been initialized")
	}
}
