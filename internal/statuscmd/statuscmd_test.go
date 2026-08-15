package statuscmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/execrun"
	"github.com/robertguss/mycelium/internal/statuscmd"
)

func TestRunDueMatrix(t *testing.T) {
	cases := []struct {
		name    string
		state   string
		revisit string
		now     time.Time
		wantDue string
	}{
		{"before", "simmering", "2026-08-08", time.Date(2026, 8, 7, 0, 0, 0, 0, time.UTC), "no"},
		{"on", "simmering", "2026-08-08", time.Date(2026, 8, 8, 0, 0, 0, 0, time.UTC), "yes"},
		{"after", "simmering", "2026-08-08", time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC), "yes"},
		{"event", "simmering", "event:after-iphone-launch", time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC), "event"},
		{"empty", "exploring", "", time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC), "no"},
		{"leftover-date-exploring", "exploring", "2026-08-08", time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC), "no"},
		{"leftover-event-clarified", "clarified", "event:after-iphone-launch", time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC), "no"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := writeFixtureDir(t, t.TempDir(), "garden-lighting", tc.state, "focused", tc.revisit, "")
			var stdout, stderr bytes.Buffer
			code := statuscmd.Run(statuscmd.Options{Dir: root, Cwd: t.TempDir()}, statuscmd.Deps{
				Clock:  clock.Fixed{T: tc.now},
				Stdout: &stdout,
				Stderr: &stderr,
			})
			if code != 0 {
				t.Fatalf("exit %d stderr=%q", code, stderr.String())
			}
			got := field(t, stdout.String(), "due")
			if got != tc.wantDue {
				t.Fatalf("due=%q want %q\n%s", got, tc.wantDue, stdout.String())
			}
			if field(t, stdout.String(), "github") != "unpublished" {
				t.Fatalf("github line: %s", stdout.String())
			}
		})
	}
}

func TestRunGithubRepo(t *testing.T) {
	root := writeFixtureDir(t, t.TempDir(), "garden-lighting", "exploring", "focused", "", "robertguss/garden-lighting")
	var stdout, stderr bytes.Buffer
	code := statuscmd.Run(statuscmd.Options{Dir: root, Cwd: t.TempDir()}, statuscmd.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)},
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	if field(t, stdout.String(), "github") != "robertguss/garden-lighting" {
		t.Fatalf("stdout=%q", stdout.String())
	}
}

func TestRunNoInstance(t *testing.T) {
	empty := t.TempDir()
	var stdout, stderr bytes.Buffer
	code := statuscmd.Run(statuscmd.Options{Dir: empty, Cwd: empty}, statuscmd.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)},
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if code != 1 {
		t.Fatalf("exit %d want 1", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout=%q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "not a mycelium instance") {
		t.Fatalf("stderr=%q", stderr.String())
	}
}

func TestRunAllOfflineTwoInstances(t *testing.T) {
	ideas := t.TempDir()
	writeFixtureDir(t, filepath.Join(ideas, "garden-lighting"), "garden-lighting", "simmering", "focused", "2026-08-08", "")
	writeFixtureDir(t, filepath.Join(ideas, "wake-fixture"), "wake-fixture", "exploring", "focused", "", "")
	// Master-only child must be ignored.
	master := filepath.Join(ideas, "mycelium-master")
	if err := os.MkdirAll(master, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(master, "research-program.toml"), []byte("name = \"x\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	fake := &execrun.Fake{Paths: map[string]string{"gh": "/usr/bin/gh"}}
	var stdout, stderr bytes.Buffer
	now := time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC)
	code := statuscmd.Run(statuscmd.Options{
		All:     true,
		Root:    ideas,
		Offline: true,
		Cwd:     t.TempDir(),
	}, statuscmd.Deps{
		Clock:     clock.Fixed{T: now},
		Stdout:    &stdout,
		Stderr:    &stderr,
		Runner:    fake,
		LookupEnv: func(string) string { return "" },
	})
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	if fake.Called("gh") {
		t.Fatalf("gh was run: %+v", fake.Calls)
	}
	lines := strings.Split(strings.TrimSuffix(stdout.String(), "\n"), "\n")
	if len(lines) != 4 {
		t.Fatalf("want 4 lines, got %d: %q", len(lines), stdout.String())
	}
	if lines[0] != "partial: local-only (offline)" {
		t.Fatalf("first line=%q", lines[0])
	}
	if lines[1] != "garden-lighting\tsimmering\tfocused\t2026-08-08\tunpublished\tunpublished" {
		t.Fatalf("row1=%q", lines[1])
	}
	if lines[2] != "wake-fixture\texploring\tfocused\t\tunpublished\tunpublished" {
		t.Fatalf("row2=%q", lines[2])
	}
	if lines[3] != "2 ideas (1 overdue, partial)" {
		t.Fatalf("summary=%q", lines[3])
	}
}

func TestRunAllIncludeDirOutsideRoot(t *testing.T) {
	ideas := t.TempDir()
	writeFixtureDir(t, filepath.Join(ideas, "inside"), "inside", "spark", "focused", "", "")
	outside := writeFixtureDir(t, t.TempDir(), "outside", "exploring", "focused", "", "")

	fake := &execrun.Fake{}
	var stdout, stderr bytes.Buffer
	code := statuscmd.Run(statuscmd.Options{
		All:     true,
		Root:    ideas,
		Dir:     outside,
		Offline: true,
		Cwd:     t.TempDir(),
	}, statuscmd.Deps{
		Clock:     clock.Fixed{T: time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC)},
		Stdout:    &stdout,
		Stderr:    &stderr,
		Runner:    fake,
		LookupEnv: func(string) string { return "" },
	})
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	if fake.Called("gh") {
		t.Fatal("gh was run")
	}
	out := stdout.String()
	if !strings.Contains(out, "outside\texploring") {
		t.Fatalf("missing outside instance: %q", out)
	}
	if !strings.Contains(out, "inside\tspark") {
		t.Fatalf("missing inside instance: %q", out)
	}
	if !strings.Contains(out, "2 ideas (0 overdue, partial)") {
		t.Fatalf("summary: %q", out)
	}
}

func TestRunAllMissingRoot(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "no-such-ideas")
	fake := &execrun.Fake{}
	var stdout, stderr bytes.Buffer
	code := statuscmd.Run(statuscmd.Options{
		All:     true,
		Root:    missing,
		Offline: true,
		Cwd:     t.TempDir(),
	}, statuscmd.Deps{
		Clock:     clock.Fixed{T: time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC)},
		Stdout:    &stdout,
		Stderr:    &stderr,
		Runner:    fake,
		LookupEnv: func(string) string { return "" },
	})
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	lines := strings.Split(strings.TrimSuffix(stdout.String(), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("want 2 lines, got %d: %q", len(lines), stdout.String())
	}
	if lines[0] != "partial: local-only (offline)" {
		t.Fatalf("first=%q", lines[0])
	}
	if lines[1] != "0 ideas (0 overdue, partial)" {
		t.Fatalf("summary=%q", lines[1])
	}
}

func TestRunAllNoMutation(t *testing.T) {
	ideas := t.TempDir()
	inst := writeFixtureDir(t, filepath.Join(ideas, "stable"), "stable", "spark", "focused", "", "")
	before, err := os.ReadFile(filepath.Join(inst, "mycelium.toml"))
	if err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := statuscmd.Run(statuscmd.Options{
		All: true, Root: ideas, Offline: true, Cwd: t.TempDir(),
	}, statuscmd.Deps{
		Clock:     clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)},
		Stdout:    &stdout,
		Stderr:    &stderr,
		Runner:    &execrun.Fake{},
		LookupEnv: func(string) string { return "" },
	})
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	after, err := os.ReadFile(filepath.Join(inst, "mycelium.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("mycelium.toml mutated")
	}
	entries, err := os.ReadDir(inst)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("unexpected files in instance: %v", entries)
	}
}

func TestRunAllOmitsArchivedUnlessFlag(t *testing.T) {
	ideas := t.TempDir()
	writeFixtureDir(t, filepath.Join(ideas, "live"), "live", "spark", "focused", "", "")
	writeFixtureDir(t, filepath.Join(ideas, "old"), "old", "archived", "focused", "", "")

	run := func(archived bool) string {
		var stdout, stderr bytes.Buffer
		code := statuscmd.Run(statuscmd.Options{
			All: true, Root: ideas, Offline: true, Archived: archived, Cwd: t.TempDir(),
		}, statuscmd.Deps{
			Clock:     clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)},
			Stdout:    &stdout,
			Stderr:    &stderr,
			Runner:    &execrun.Fake{},
			LookupEnv: func(string) string { return "" },
		})
		if code != 0 {
			t.Fatalf("archived=%v exit %d stderr=%q", archived, code, stderr.String())
		}
		return stdout.String()
	}

	hidden := run(false)
	if strings.Contains(hidden, "old\t") {
		t.Fatalf("archived shown without flag: %q", hidden)
	}
	if !strings.Contains(hidden, "1 ideas (0 overdue, partial)") {
		t.Fatalf("hidden summary: %q", hidden)
	}
	shown := run(true)
	if !strings.Contains(shown, "old\tarchived") {
		t.Fatalf("archived missing with flag: %q", shown)
	}
	if !strings.Contains(shown, "2 ideas (0 overdue, partial)") {
		t.Fatalf("shown summary: %q", shown)
	}
}

func TestRunAllUnreadableChildTeachContinue(t *testing.T) {
	ideas := t.TempDir()
	writeFixtureDir(t, filepath.Join(ideas, "good"), "good", "spark", "focused", "", "")
	bad := filepath.Join(ideas, "bad")
	if err := os.MkdirAll(bad, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bad, "mycelium.toml"), []byte("not = valid\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := statuscmd.Run(statuscmd.Options{
		All: true, Root: ideas, Offline: true, Cwd: t.TempDir(),
	}, statuscmd.Deps{
		Clock:     clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)},
		Stdout:    &stdout,
		Stderr:    &stderr,
		Runner:    &execrun.Fake{},
		LookupEnv: func(string) string { return "" },
	})
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	if !strings.Contains(stderr.String(), "cannot read instance") {
		t.Fatalf("stderr=%q", stderr.String())
	}
	if strings.Contains(stdout.String(), "bad\t") {
		t.Fatalf("bad row on stdout: %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "good\tspark") {
		t.Fatalf("good missing: %q", stdout.String())
	}
}

func TestRunAllGhMissingReason(t *testing.T) {
	ideas := t.TempDir()
	fake := &execrun.Fake{} // no gh path
	var stdout, stderr bytes.Buffer
	code := statuscmd.Run(statuscmd.Options{
		All: true, Root: ideas, Offline: false, Cwd: t.TempDir(),
	}, statuscmd.Deps{
		Clock:     clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)},
		Stdout:    &stdout,
		Stderr:    &stderr,
		Runner:    fake,
		LookupEnv: func(string) string { return "" },
	})
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	if fake.Called("gh") {
		t.Fatal("gh Run must not happen")
	}
	if !strings.HasPrefix(stdout.String(), "partial: local-only (gh missing)\n") {
		t.Fatalf("stdout=%q", stdout.String())
	}
}

func TestSortIdeasBuckets(t *testing.T) {
	now := time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC)
	ideas := []statuscmd.Idea{
		{Slug: "z-arch", State: "archived", Flag: "unpublished", Github: "unpublished"},
		{Slug: "b-future", State: "simmering", Revisit: "2026-08-20", Flag: "unpublished", Github: "unpublished"},
		{Slug: "a-overdue", State: "simmering", Revisit: "2026-08-08", Flag: "unpublished", Github: "unpublished"},
		{Slug: "c-today", State: "simmering", Revisit: "2026-08-09", Flag: "unpublished", Github: "unpublished"},
		{Slug: "e-spark", State: "spark", Flag: "unpublished", Github: "unpublished"},
		{Slug: "d-exploring", State: "exploring", Flag: "unpublished", Github: "unpublished"},
		{Slug: "f-clarified", State: "clarified", Flag: "unpublished", Github: "unpublished"},
		{Slug: "g-event", State: "simmering", Revisit: "event:launch", Flag: "unpublished", Github: "unpublished"},
		{Slug: "a2-overdue", State: "simmering", Revisit: "2026-08-01", Flag: "unpublished", Github: "unpublished"},
	}
	statuscmd.SortIdeas(ideas, now)
	want := []string{
		"a-overdue", "a2-overdue", "c-today", "d-exploring", "e-spark",
		"f-clarified", "g-event", "b-future", "z-arch",
	}
	for i, slug := range want {
		if ideas[i].Slug != slug {
			t.Fatalf("pos %d = %q want %q (full=%v)", i, ideas[i].Slug, slug, slugs(ideas))
		}
	}
}

func slugs(ideas []statuscmd.Idea) []string {
	out := make([]string, len(ideas))
	for i, idea := range ideas {
		out[i] = idea.Slug
	}
	return out
}

func writeFixtureDir(t *testing.T, root, slug, state, tier, revisit, github string) string {
	t.Helper()
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	ideaName := strings.ReplaceAll(slug, "-", " ")
	toml := `schema_version = 1
idea_name = "` + ideaName + `"
slug = "` + slug + `"
state = "` + state + `"
tier = "` + tier + `"
methodology_version = "2.0.0"
generated_by_cli_version = "0.1.0-dev"
created_date = "2026-08-01"
updated_date = "2026-08-01"
revisit = "` + revisit + `"
github_repo = "` + github + `"
`
	if err := os.WriteFile(filepath.Join(root, "mycelium.toml"), []byte(toml), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func field(t *testing.T, stdout, key string) string {
	t.Helper()
	prefix := key + ": "
	for _, line := range strings.Split(stdout, "\n") {
		if strings.HasPrefix(line, prefix) {
			return strings.TrimPrefix(line, prefix)
		}
	}
	t.Fatalf("missing %s in %q", key, stdout)
	return ""
}
