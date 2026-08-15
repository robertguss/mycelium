package statuscmd_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/execrun"
	"github.com/robertguss/mycelium/internal/statuscmd"
)

func TestMergeIdeasFlags(t *testing.T) {
	local := []statuscmd.Idea{
		{Slug: "garden-lighting", State: "simmering", Tier: "focused", Revisit: "2026-08-08", Flag: "unpublished", Github: "unpublished"},
		{Slug: "local-only", State: "spark", Tier: "focused", Flag: "unpublished", Github: "unpublished"},
		{Slug: "stale-repo", State: "exploring", Tier: "focused", Flag: "unpublished", Github: "unpublished"},
	}
	remotes := []statuscmd.Remote{
		{Name: "garden-lighting", Owner: "robertguss", State: "clarified", Tier: "standard", Revisit: "2026-01-01"},
		{Name: "other-idea", Owner: "robertguss", State: "unread", Tier: "-", Revisit: "-"},
	}
	got := statuscmd.MergeIdeas(local, remotes, false)
	statuscmd.SortIdeas(got, time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC))

	by := map[string]statuscmd.Idea{}
	for _, idea := range got {
		by[idea.Slug] = idea
	}
	if len(got) != 4 {
		t.Fatalf("len=%d want 4: %+v", len(got), got)
	}
	gl := by["garden-lighting"]
	if gl.Flag != "ok" || gl.Github != "robertguss/garden-lighting" {
		t.Fatalf("both: %+v", gl)
	}
	if gl.State != "simmering" || gl.Tier != "focused" || gl.Revisit != "2026-08-08" {
		t.Fatalf("prefer local manifest: %+v", gl)
	}
	oi := by["other-idea"]
	if oi.Flag != "remote" || oi.State != "unread" || oi.Tier != "-" || oi.Revisit != "-" {
		t.Fatalf("remote unread: %+v", oi)
	}
	if oi.Github != "robertguss/other-idea" {
		t.Fatalf("remote github: %+v", oi)
	}
	lo := by["local-only"]
	if lo.Flag != "unpublished" || lo.Github != "unpublished" {
		t.Fatalf("local-only: %+v", lo)
	}
	sr := by["stale-repo"]
	if sr.Flag != "unpublished" || sr.Github != "unpublished" {
		t.Fatalf("local with no remote match stays unpublished: %+v", sr)
	}
}

func TestMergeIdeasArchivedFilter(t *testing.T) {
	local := []statuscmd.Idea{
		{Slug: "old-local", State: "archived", Tier: "focused", Flag: "unpublished", Github: "unpublished"},
		{Slug: "both-arch", State: "exploring", Tier: "focused", Flag: "unpublished", Github: "unpublished"},
		{Slug: "live", State: "spark", Tier: "focused", Flag: "unpublished", Github: "unpublished"},
	}
	remotes := []statuscmd.Remote{
		{Name: "remote-arch", Owner: "robertguss", IsArchived: true, State: "spark", Tier: "focused"},
		{Name: "both-arch", Owner: "robertguss", IsArchived: true, State: "spark", Tier: "focused"},
		{Name: "live", Owner: "robertguss", State: "spark", Tier: "focused"},
	}

	hidden := statuscmd.MergeIdeas(local, remotes, false)
	slugs := map[string]bool{}
	for _, idea := range hidden {
		slugs[idea.Slug] = true
	}
	if slugs["old-local"] || slugs["remote-arch"] || slugs["both-arch"] {
		t.Fatalf("archived leaked: %v", slugs)
	}
	if !slugs["live"] {
		t.Fatal("live missing")
	}
	if len(hidden) != 1 {
		t.Fatalf("hidden len=%d want 1", len(hidden))
	}

	shown := statuscmd.MergeIdeas(local, remotes, true)
	if len(shown) != 4 {
		t.Fatalf("shown len=%d want 4: %+v", len(shown), shown)
	}
}

func TestRunAllGitHubComplete(t *testing.T) {
	ideas := t.TempDir()
	writeFixtureDir(t, filepath.Join(ideas, "garden-lighting"), "garden-lighting", "simmering", "focused", "2026-08-08", "")
	writeFixtureDir(t, filepath.Join(ideas, "local-only"), "local-only", "spark", "focused", "", "")

	manifestB64 := base64.StdEncoding.EncodeToString([]byte(`schema_version = 1
idea_name = "Garden Lighting"
slug = "garden-lighting"
state = "clarified"
tier = "standard"
methodology_version = "2.0.0"
generated_by_cli_version = "0.1.0-dev"
created_date = "2026-08-01"
updated_date = "2026-08-01"
revisit = "2026-01-01"
github_repo = "robertguss/garden-lighting"
`))

	fake := &execrun.Fake{
		Paths: map[string]string{"gh": "/usr/bin/gh"},
		RunFunc: func(ctx context.Context, name string, args []string, opts execrun.RunOpts) (execrun.Result, error) {
			if name != "gh" {
				return execrun.Result{ExitCode: 1}, nil
			}
			joined := strings.Join(args, " ")
			switch {
			case joined == "api user --jq .login":
				return execrun.Result{Stdout: []byte("robertguss\n")}, nil
			case strings.HasPrefix(joined, "search repos topic:idea user:robertguss"):
				return execrun.Result{Stdout: []byte(`[
  {"name":"garden-lighting","url":"https://github.com/robertguss/garden-lighting","isArchived":false,"owner":{"login":"robertguss"}},
  {"name":"other-idea","url":"https://github.com/robertguss/other-idea","isArchived":false,"owner":{"login":"robertguss"}}
]`)}, nil
			case joined == "api repos/robertguss/garden-lighting/contents/mycelium.toml --jq .content":
				return execrun.Result{Stdout: []byte(manifestB64 + "\n")}, nil
			case joined == "api repos/robertguss/other-idea/contents/mycelium.toml --jq .content":
				return execrun.Result{ExitCode: 1, Stderr: []byte("404")}, nil
			default:
				return execrun.Result{ExitCode: 1, Stderr: []byte("unexpected: " + joined)}, nil
			}
		},
	}

	var stdout, stderr bytes.Buffer
	code := statuscmd.Run(statuscmd.Options{
		All: true, Root: ideas, Offline: false, Cwd: t.TempDir(),
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
	if strings.HasPrefix(stdout.String(), "partial:") {
		t.Fatalf("unexpected partial: %q", stdout.String())
	}
	lines := strings.Split(strings.TrimSuffix(stdout.String(), "\n"), "\n")
	if len(lines) != 4 {
		t.Fatalf("want 4 lines, got %d: %q", len(lines), stdout.String())
	}
	// Sort §9.6: overdue simmering → spark → unread (rest).
	if lines[0] != "garden-lighting\tsimmering\tfocused\t2026-08-08\tok\trobertguss/garden-lighting" {
		t.Fatalf("row0=%q", lines[0])
	}
	if lines[1] != "local-only\tspark\tfocused\t\tunpublished\tunpublished" {
		t.Fatalf("row1=%q", lines[1])
	}
	if lines[2] != "other-idea\tunread\t-\t-\tremote\trobertguss/other-idea" {
		t.Fatalf("row2=%q", lines[2])
	}
	if lines[3] != "3 ideas (1 overdue, complete)" {
		t.Fatalf("summary=%q", lines[3])
	}
}

func TestRunAllGhFailedReason(t *testing.T) {
	ideas := t.TempDir()
	writeFixtureDir(t, filepath.Join(ideas, "garden-lighting"), "garden-lighting", "simmering", "focused", "2026-08-08", "")

	fake := &execrun.Fake{
		Paths: map[string]string{"gh": "/usr/bin/gh"},
		RunFunc: func(ctx context.Context, name string, args []string, opts execrun.RunOpts) (execrun.Result, error) {
			return execrun.Result{ExitCode: 1, Stderr: []byte("not logged in")}, nil
		},
	}
	var stdout, stderr bytes.Buffer
	code := statuscmd.Run(statuscmd.Options{
		All: true, Root: ideas, Offline: false, Cwd: t.TempDir(),
	}, statuscmd.Deps{
		Clock:     clock.Fixed{T: time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC)},
		Stdout:    &stdout,
		Stderr:    &stderr,
		Runner:    fake,
		LookupEnv: func(string) string { return "" },
	})
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	if !fake.Called("gh") {
		t.Fatal("expected gh attempt")
	}
	lines := strings.Split(strings.TrimSuffix(stdout.String(), "\n"), "\n")
	if lines[0] != "partial: local-only (gh failed)" {
		t.Fatalf("first=%q", lines[0])
	}
	if lines[len(lines)-1] != "1 ideas (1 overdue, partial)" {
		t.Fatalf("summary=%q", lines[len(lines)-1])
	}
	if !strings.Contains(stdout.String(), "unpublished\tunpublished") {
		t.Fatalf("want unpublished local row: %q", stdout.String())
	}
}

func TestRunAllOfflineNeverCallsGh(t *testing.T) {
	ideas := t.TempDir()
	writeFixtureDir(t, filepath.Join(ideas, "x"), "x", "spark", "focused", "", "")
	fake := &execrun.Fake{Paths: map[string]string{"gh": "/usr/bin/gh"}}
	var stdout, stderr bytes.Buffer
	code := statuscmd.Run(statuscmd.Options{
		All: true, Root: ideas, Offline: true, Cwd: t.TempDir(),
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
		t.Fatalf("gh called under --offline: %+v", fake.Calls)
	}
	for _, c := range fake.Calls {
		if c.Kind == "lookpath" && c.Name == "gh" {
			t.Fatalf("LookPath(gh) under --offline: %+v", fake.Calls)
		}
	}
	if !strings.HasPrefix(stdout.String(), "partial: local-only (offline)\n") {
		t.Fatalf("stdout=%q", stdout.String())
	}
}

func TestRunAllEnvOfflineNeverCallsGh(t *testing.T) {
	ideas := t.TempDir()
	writeFixtureDir(t, filepath.Join(ideas, "x"), "x", "spark", "focused", "", "")
	fake := &execrun.Fake{Paths: map[string]string{"gh": "/usr/bin/gh"}}
	var stdout, stderr bytes.Buffer
	code := statuscmd.Run(statuscmd.Options{
		All: true, Root: ideas, Offline: false, Cwd: t.TempDir(),
	}, statuscmd.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)},
		Stdout: &stdout,
		Stderr: &stderr,
		Runner: fake,
		LookupEnv: func(k string) string {
			if k == "MYCELIUM_OFFLINE" {
				return "1"
			}
			return ""
		},
	})
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	if fake.Called("gh") {
		t.Fatalf("gh called under MYCELIUM_OFFLINE=1: %+v", fake.Calls)
	}
	if !strings.HasPrefix(stdout.String(), "partial: local-only (offline)\n") {
		t.Fatalf("stdout=%q", stdout.String())
	}
}

func TestRunAllSearchFallbackThenComplete(t *testing.T) {
	ideas := t.TempDir()
	writeFixtureDir(t, filepath.Join(ideas, "only-remote-local"), "only-remote-local", "spark", "focused", "", "")

	manifestB64 := base64.StdEncoding.EncodeToString([]byte(`schema_version = 1
idea_name = "Remote Idea"
slug = "remote-idea"
state = "exploring"
tier = "focused"
methodology_version = "2.0.0"
generated_by_cli_version = "0.1.0-dev"
created_date = "2026-08-01"
updated_date = "2026-08-01"
revisit = ""
github_repo = "robertguss/remote-idea"
`))

	var sawPrimary, sawOwnerSearch bool
	fake := &execrun.Fake{
		Paths: map[string]string{"gh": "/usr/bin/gh"},
		RunFunc: func(ctx context.Context, name string, args []string, opts execrun.RunOpts) (execrun.Result, error) {
			joined := strings.Join(args, " ")
			switch {
			case joined == "api user --jq .login":
				return execrun.Result{Stdout: []byte("robertguss\n")}, nil
			case strings.HasPrefix(joined, "search repos topic:idea user:robertguss"):
				sawPrimary = true
				return execrun.Result{ExitCode: 1, Stderr: []byte("API rate limit")}, nil
			case strings.Contains(joined, "search repos --owner robertguss --topic idea"):
				sawOwnerSearch = true
				return execrun.Result{Stdout: []byte(`[
  {"name":"remote-idea","url":"https://github.com/robertguss/remote-idea","isArchived":false,"owner":{"login":"robertguss"}}
]`)}, nil
			case joined == "api repos/robertguss/remote-idea/contents/mycelium.toml --jq .content":
				return execrun.Result{Stdout: []byte(manifestB64)}, nil
			default:
				return execrun.Result{ExitCode: 1, Stderr: []byte(fmt.Sprintf("unexpected %q", joined))}, nil
			}
		},
	}

	var stdout, stderr bytes.Buffer
	code := statuscmd.Run(statuscmd.Options{
		All: true, Root: ideas, Cwd: t.TempDir(),
	}, statuscmd.Deps{
		Clock:     clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)},
		Stdout:    &stdout,
		Stderr:    &stderr,
		Runner:    fake,
		LookupEnv: func(string) string { return "" },
	})
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	if !sawPrimary || !sawOwnerSearch {
		t.Fatalf("fallback path not exercised primary=%v owner=%v", sawPrimary, sawOwnerSearch)
	}
	if strings.Contains(stdout.String(), "partial:") {
		t.Fatalf("want complete: %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "remote-idea\texploring\tfocused\t\tremote\trobertguss/remote-idea") {
		t.Fatalf("missing remote row: %q", stdout.String())
	}
	if !strings.HasSuffix(strings.TrimSpace(stdout.String()), "2 ideas (0 overdue, complete)") {
		t.Fatalf("summary: %q", stdout.String())
	}
}

func TestRunAllRemoteArchivedHiddenUnlessFlag(t *testing.T) {
	ideas := t.TempDir()
	fake := &execrun.Fake{
		Paths: map[string]string{"gh": "/usr/bin/gh"},
		RunFunc: func(ctx context.Context, name string, args []string, opts execrun.RunOpts) (execrun.Result, error) {
			joined := strings.Join(args, " ")
			switch {
			case joined == "api user --jq .login":
				return execrun.Result{Stdout: []byte("robertguss\n")}, nil
			case strings.HasPrefix(joined, "search repos topic:idea"):
				return execrun.Result{Stdout: []byte(`[
  {"name":"old-remote","url":"https://github.com/robertguss/old-remote","isArchived":true,"owner":{"login":"robertguss"}},
  {"name":"live-remote","url":"https://github.com/robertguss/live-remote","isArchived":false,"owner":{"login":"robertguss"}}
]`)}, nil
			case strings.Contains(joined, "/contents/mycelium.toml"):
				return execrun.Result{ExitCode: 1}, nil
			default:
				return execrun.Result{ExitCode: 1, Stderr: []byte(joined)}, nil
			}
		},
	}

	run := func(archived bool) string {
		var stdout, stderr bytes.Buffer
		code := statuscmd.Run(statuscmd.Options{
			All: true, Root: ideas, Archived: archived, Cwd: t.TempDir(),
		}, statuscmd.Deps{
			Clock:     clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)},
			Stdout:    &stdout,
			Stderr:    &stderr,
			Runner:    fake,
			LookupEnv: func(string) string { return "" },
		})
		if code != 0 {
			t.Fatalf("archived=%v exit %d stderr=%q", archived, code, stderr.String())
		}
		return stdout.String()
	}

	hidden := run(false)
	if strings.Contains(hidden, "old-remote\t") {
		t.Fatalf("archived remote shown: %q", hidden)
	}
	if !strings.Contains(hidden, "live-remote\tunread") {
		t.Fatalf("live remote missing: %q", hidden)
	}
	if !strings.Contains(hidden, "1 ideas (0 overdue, complete)") {
		t.Fatalf("hidden summary: %q", hidden)
	}

	shown := run(true)
	if !strings.Contains(shown, "old-remote\tunread") {
		t.Fatalf("archived remote missing with flag: %q", shown)
	}
	if !strings.Contains(shown, "2 ideas (0 overdue, complete)") {
		t.Fatalf("shown summary: %q", shown)
	}
}
