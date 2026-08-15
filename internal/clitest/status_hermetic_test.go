package clitest_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/robertguss/mycelium/internal/clitest"
	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/execrun"
)

func TestStatusSingleInstanceDue(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}

	inst := scaffoldIdea(t, work, "Garden Lighting", clk, rec, env)
	mustState(t, clk, rec, env, work, inst, "exploring")
	mustState(t, clk, rec, env, work, inst, "simmering", "--revisit", "2026-08-08")

	before := clock.Fixed{T: time.Date(2026, 8, 7, 0, 0, 0, 0, time.UTC)}
	code, stdout, stderr := runCLI(t, before, rec, env, work, "status", "--dir", inst)
	if code != 0 {
		t.Fatalf("status before due: exit %d stderr=%q", code, stderr)
	}
	assertStatusLines(t, stdout, map[string]string{
		"slug":    "garden-lighting",
		"state":   "simmering",
		"tier":    "focused",
		"revisit": "2026-08-08",
		"due":     "no",
		"github":  "unpublished",
	})

	onDue := clock.Fixed{T: time.Date(2026, 8, 8, 0, 0, 0, 0, time.UTC)}
	code, stdout, stderr = runCLI(t, onDue, rec, env, work, "status", "--dir", inst)
	if code != 0 {
		t.Fatalf("status on due: exit %d stderr=%q", code, stderr)
	}
	assertStatusLines(t, stdout, map[string]string{
		"slug":    "garden-lighting",
		"state":   "simmering",
		"tier":    "focused",
		"revisit": "2026-08-08",
		"due":     "yes",
		"github":  "unpublished",
	})

	clitest.AssertNoNetwork(t, rec)
	if rec.Called("gh") {
		t.Fatal("gh was called")
	}
}

func TestStatusEventRevisit(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}

	inst := scaffoldIdea(t, work, "Event Idea", clk, rec, env)
	mustState(t, clk, rec, env, work, inst, "exploring")
	mustState(t, clk, rec, env, work, inst, "simmering", "--revisit", "event:after-iphone-launch")

	code, stdout, stderr := runCLI(t, clk, rec, env, work, "status", "--dir", inst)
	if code != 0 {
		t.Fatalf("status event: exit %d stderr=%q", code, stderr)
	}
	assertStatusLines(t, stdout, map[string]string{
		"slug":    "event-idea",
		"state":   "simmering",
		"tier":    "focused",
		"revisit": "event:after-iphone-launch",
		"due":     "event",
		"github":  "unpublished",
	})
	clitest.AssertNoNetwork(t, rec)
}

func TestStatusNoInstanceTeachingError(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}

	code, stdout, stderr := runCLI(t, clk, rec, env, work, "status")
	if code != 1 {
		t.Fatalf("exit %d want 1 stderr=%q", code, stderr)
	}
	if stdout != "" {
		t.Fatalf("stdout should be empty, got %q", stdout)
	}
	if !strings.Contains(stderr, "not a mycelium instance") {
		t.Fatalf("stderr=%q", stderr)
	}
	clitest.AssertNoNetwork(t, rec)
}

func TestStatusAllOfflineTwoInstances(t *testing.T) {
	work := t.TempDir()
	ideas := filepath.Join(work, "ideas")
	if err := os.MkdirAll(ideas, 0o755); err != nil {
		t.Fatal(err)
	}
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}

	a := scaffoldIdea(t, ideas, "Garden Lighting", clk, rec, env)
	b := scaffoldIdea(t, ideas, "Wake Fixture", clk, rec, env)
	mustState(t, clk, rec, env, ideas, a, "exploring")
	mustState(t, clk, rec, env, ideas, a, "simmering", "--revisit", "2026-08-08")
	mustState(t, clk, rec, env, ideas, b, "exploring")

	callsBefore := len(rec.Calls)
	due := clock.Fixed{T: time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC)}
	code, stdout, stderr := runCLI(t, due, rec, env, work, "status", "--all", "--offline", "--root", ideas)
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr)
	}
	for _, c := range rec.Calls[callsBefore:] {
		if c.Kind == "run" && (c.Name == "gh" || filepath.Base(c.Name) == "gh") {
			t.Fatalf("gh run: %+v", c)
		}
	}
	lines := strings.Split(strings.TrimSuffix(stdout, "\n"), "\n")
	if len(lines) < 2 {
		t.Fatalf("stdout=%q", stdout)
	}
	if !strings.HasPrefix(lines[0], "partial: local-only (") {
		t.Fatalf("first line=%q", lines[0])
	}
	if lines[0] != "partial: local-only (offline)" {
		t.Fatalf("want offline reason, got %q", lines[0])
	}
	if lines[len(lines)-1] != "2 ideas (1 overdue, partial)" {
		t.Fatalf("summary=%q stdout=%q", lines[len(lines)-1], stdout)
	}
	if !strings.Contains(stdout, "garden-lighting\tsimmering") {
		t.Fatalf("missing garden-lighting: %q", stdout)
	}
	if !strings.Contains(stdout, "wake-fixture\texploring") {
		t.Fatalf("missing wake-fixture: %q", stdout)
	}
	clitest.AssertNoNetwork(t, rec)
}

func TestStatusAllMissingRootPartial(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}
	missing := filepath.Join(work, "missing-ideas")

	callsBefore := len(rec.Calls)
	code, stdout, stderr := runCLI(t, clk, rec, env, work, "status", "--all", "--offline", "--root", missing)
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr)
	}
	for _, c := range rec.Calls[callsBefore:] {
		if c.Kind == "run" && (c.Name == "gh" || filepath.Base(c.Name) == "gh") {
			t.Fatalf("gh run: %+v", c)
		}
	}
	want := "partial: local-only (offline)\n0 ideas (0 overdue, partial)\n"
	if stdout != want {
		t.Fatalf("stdout=%q want %q", stdout, want)
	}
}

func TestStatusRootWithoutAllTeachingError(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}

	code, _, stderr := runCLI(t, clk, rec, env, work, "status", "--root", work)
	if code != 1 {
		t.Fatalf("exit %d want 1 stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "--root requires --all") {
		t.Fatalf("stderr=%q", stderr)
	}
}

func TestStatusOfflineNoOp(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}

	inst := scaffoldIdea(t, work, "Offline Status", clk, rec, env)
	mustState(t, clk, rec, env, work, inst, "exploring")
	mustState(t, clk, rec, env, work, inst, "simmering", "--revisit", "2026-08-08")

	callsBefore := len(rec.Calls)
	dueClk := clock.Fixed{T: time.Date(2026, 8, 8, 12, 0, 0, 0, time.UTC)}
	code, stdout, stderr := runCLI(t, dueClk, rec, env, work, "status", "--offline", "--dir", inst)
	if code != 0 {
		t.Fatalf("status --offline: exit %d stderr=%q", code, stderr)
	}
	assertStatusLines(t, stdout, map[string]string{
		"slug":    "offline-status",
		"state":   "simmering",
		"tier":    "focused",
		"revisit": "2026-08-08",
		"due":     "yes",
		"github":  "unpublished",
	})
	for _, c := range rec.Calls[callsBefore:] {
		if c.Kind == "run" && filepath.Base(c.Name) == "gh" {
			t.Fatalf("gh run after status: %+v", c)
		}
	}
	clitest.AssertNoNetwork(t, rec)
}

func TestStatusMYCELIUM_NOWBinary(t *testing.T) {
	bin := clitest.Bin(t)
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)}
	envMap := map[string]string{"MYCELIUM_OFFLINE": "1"}

	inst := scaffoldIdea(t, work, "Clock Env", clk, rec, envMap)
	mustState(t, clk, rec, envMap, work, inst, "exploring")
	mustState(t, clk, rec, envMap, work, inst, "simmering", "--revisit", "2026-08-08")

	code, stdout, stderr := clitest.Run(t, bin, work, []string{
		"MYCELIUM_OFFLINE=1",
		"MYCELIUM_NOW=2026-08-07T00:00:00Z",
	}, "status", "--dir", inst)
	if code != 0 {
		t.Fatalf("binary status: exit %d stderr=%q", code, stderr)
	}
	assertStatusLines(t, stdout, map[string]string{
		"slug":    "clock-env",
		"state":   "simmering",
		"tier":    "focused",
		"revisit": "2026-08-08",
		"due":     "no",
		"github":  "unpublished",
	})

	code, stdout, stderr = clitest.Run(t, bin, work, []string{
		"MYCELIUM_OFFLINE=1",
		"MYCELIUM_NOW=2026-08-08T00:00:00Z",
	}, "status", "--dir", inst)
	if code != 0 {
		t.Fatalf("binary status due: exit %d stderr=%q", code, stderr)
	}
	assertStatusLines(t, stdout, map[string]string{
		"slug":    "clock-env",
		"state":   "simmering",
		"tier":    "focused",
		"revisit": "2026-08-08",
		"due":     "yes",
		"github":  "unpublished",
	})
}

func TestStatusWalkUpFromCwd(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}

	inst := scaffoldIdea(t, work, "Walk Up", clk, rec, env)
	mustState(t, clk, rec, env, work, inst, "exploring")
	nested := filepath.Join(inst, "nested", "deeper")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}

	code, stdout, stderr := runCLI(t, clk, rec, env, nested, "status")
	if code != 0 {
		t.Fatalf("walk-up status: exit %d stderr=%q", code, stderr)
	}
	assertStatusLines(t, stdout, map[string]string{
		"slug":    "walk-up",
		"state":   "exploring",
		"tier":    "focused",
		"revisit": "",
		"due":     "no",
		"github":  "unpublished",
	})
}

func assertStatusLines(t *testing.T, stdout string, want map[string]string) {
	t.Helper()
	lines := strings.Split(strings.TrimSuffix(stdout, "\n"), "\n")
	if len(lines) != 6 {
		t.Fatalf("want 6 lines, got %d: %q", len(lines), stdout)
	}
	keys := []string{"slug", "state", "tier", "revisit", "due", "github"}
	for i, key := range keys {
		prefix := key + ": "
		if !strings.HasPrefix(lines[i], prefix) {
			t.Fatalf("line %d=%q want prefix %q", i, lines[i], prefix)
		}
		got := strings.TrimPrefix(lines[i], prefix)
		if got != want[key] {
			t.Fatalf("%s=%q want %q (stdout=%q)", key, got, want[key], stdout)
		}
	}
}
