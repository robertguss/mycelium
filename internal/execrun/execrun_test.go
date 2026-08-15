package execrun_test

import (
	"context"
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
	// LookPath("gh") must not count as a Run.
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
