package statuscmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/robertguss/mycelium/internal/clock"
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
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := writeFixture(t, "garden-lighting", tc.state, "focused", tc.revisit, "")
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
	root := writeFixture(t, "garden-lighting", "exploring", "focused", "", "robertguss/garden-lighting")
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

func writeFixture(t *testing.T, slug, state, tier, revisit, github string) string {
	t.Helper()
	root := t.TempDir()
	toml := `schema_version = 1
idea_name = "Garden Lighting"
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
