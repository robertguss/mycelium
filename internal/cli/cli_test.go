package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/robertguss/mycelium/internal/cli"
	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/clitest"
	"github.com/robertguss/mycelium/internal/execrun"
	"github.com/robertguss/mycelium/internal/manifest"
	"github.com/robertguss/mycelium/internal/version"
)

func TestVersionStdoutEquality(t *testing.T) {
	cases := []struct {
		name string
		argv []string
		want string
	}{
		{"bare version", []string{"mycelium", "version"}, version.Version + "\n"},
		{"default stamp", []string{"mycelium", "version"}, "0.1.0-dev\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := cli.Run(tc.argv, &stdout, &stderr, cli.Deps{})
			if code != 0 {
				t.Fatalf("exit %d; stderr=%q", code, stderr.String())
			}
			if stdout.String() != tc.want {
				t.Fatalf("stdout=%q want %q", stdout.String(), tc.want)
			}
			if stderr.Len() != 0 {
				t.Fatalf("stderr not empty: %q", stderr.String())
			}
		})
	}
}

func TestHelpExitsZero(t *testing.T) {
	cases := [][]string{
		{"mycelium"},
		{"mycelium", "-h"},
		{"mycelium", "--help"},
		{"mycelium", "help"},
	}
	for _, argv := range cases {
		t.Run(strings.Join(argv, " "), func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := cli.Run(argv, &stdout, &stderr, cli.Deps{})
			if code != 0 {
				t.Fatalf("exit %d; stderr=%q", code, stderr.String())
			}
			if !strings.Contains(stdout.String(), "mycelium") {
				t.Fatalf("usage missing from stdout: %q", stdout.String())
			}
		})
	}
}

func TestUnknownCommandTeachingError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"mycelium", "not-a-command"}, &stdout, &stderr, cli.Deps{})
	if code != 1 {
		t.Fatalf("exit %d want 1", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout should be empty, got %q", stdout.String())
	}
	errText := stderr.String()
	for _, prefix := range []string{"mycelium:", "convention:", "contract:", "fix:"} {
		if !strings.Contains(errText, prefix) {
			t.Fatalf("teaching error missing %q in %q", prefix, errText)
		}
	}
	lines := strings.Split(strings.TrimSuffix(errText, "\n"), "\n")
	if len(lines) != 4 {
		t.Fatalf("want 4 teaching lines, got %d: %q", len(lines), errText)
	}
}

func TestPhase01PublishOfflineRefuses(t *testing.T) {
	root := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	deps := cli.Deps{
		Clock:     clock.Fixed{T: time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC)},
		Runner:    rec,
		Getwd:     func() (string, error) { return root, nil },
		LookupEnv: func(string) string { return "" },
	}
	var stdout, stderr bytes.Buffer
	code := cli.Run(
		[]string{"mycelium", "new", "idea", "Pub Later", "--offline"},
		&stdout, &stderr, deps,
	)
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr.String())
	}
	inst := filepath.Join(root, "pub-later")
	stdout.Reset()
	stderr.Reset()
	deps.LookupEnv = func(k string) string {
		if k == "MYCELIUM_OFFLINE" {
			return "1"
		}
		return ""
	}
	code = cli.Run(
		[]string{"mycelium", "publish", "--dir", inst},
		&stdout, &stderr, deps,
	)
	if code != 1 {
		t.Fatalf("publish exit %d want 1", code)
	}
	if !strings.Contains(stderr.String(), "offline") {
		t.Fatalf("stderr=%q", stderr.String())
	}
	clitest.AssertNoNetwork(t, rec)
}

func TestCLINewIdeaOfflineThenCheck(t *testing.T) {
	root := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	deps := cli.Deps{
		Clock:     clock.Fixed{T: time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC)},
		Runner:    rec,
		Getwd:     func() (string, error) { return root, nil },
		LookupEnv: func(string) string { return "" },
	}
	var stdout, stderr bytes.Buffer
	code := cli.Run(
		[]string{"mycelium", "new", "idea", "Check Me", "--offline"},
		&stdout, &stderr, deps,
	)
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr.String())
	}
	inst := filepath.Join(root, "check-me")
	stdout.Reset()
	stderr.Reset()
	code = cli.Run(
		[]string{"mycelium", "check", "--dir", inst},
		&stdout, &stderr, deps,
	)
	if code != 0 {
		t.Fatalf("check exit %d stderr=%q", code, stderr.String())
	}
	got := stdout.String()
	wantLines := []string{
		"mycelium check: ok",
		"instance: check-me",
		"state: spark",
		"tier: focused",
		"artifacts: 0",
	}
	for _, line := range wantLines {
		if !strings.Contains(got, line) {
			t.Fatalf("stdout missing %q in %q", line, got)
		}
	}
}

func TestCLICheckTeachingErrorFourLineShape(t *testing.T) {
	root := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	deps := cli.Deps{
		Clock:     clock.Fixed{T: time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC)},
		Runner:    rec,
		Getwd:     func() (string, error) { return root, nil },
		LookupEnv: func(string) string { return "" },
	}
	var stdout, stderr bytes.Buffer
	if code := cli.Run([]string{"mycelium", "new", "idea", "Shape", "--offline"}, &stdout, &stderr, deps); code != 0 {
		t.Fatalf("scaffold: %s", stderr.String())
	}
	inst := filepath.Join(root, "shape")
	if err := os.WriteFile(filepath.Join(inst, "scratch.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	stderr.Reset()
	code := cli.Run([]string{"mycelium", "check", "--dir", inst}, &stdout, &stderr, deps)
	if code != 1 {
		t.Fatalf("exit %d want 1", code)
	}
	errText := stderr.String()
	lines := strings.Split(strings.TrimSuffix(errText, "\n"), "\n")
	if len(lines) < 4 {
		t.Fatalf("want >=4 teaching lines, got %d: %q", len(lines), errText)
	}
	if !strings.HasPrefix(lines[0], "mycelium: ") {
		t.Fatalf("line0=%q", lines[0])
	}
	if !strings.HasPrefix(lines[1], "convention: ") {
		t.Fatalf("line1=%q", lines[1])
	}
	if !strings.HasPrefix(lines[2], "contract: ") {
		t.Fatalf("line2=%q", lines[2])
	}
	if !strings.HasPrefix(lines[3], "fix: ") {
		t.Fatalf("line3=%q", lines[3])
	}
}

func TestCLINewIdeaOffline(t *testing.T) {
	root := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	var stdout, stderr bytes.Buffer
	code := cli.Run(
		[]string{"mycelium", "new", "idea", "CLI Spark", "--offline"},
		&stdout, &stderr,
		cli.Deps{
			Clock:     clock.Fixed{T: time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)},
			Runner:    rec,
			Getwd:     func() (string, error) { return root, nil },
			LookupEnv: func(string) string { return "" },
		},
	)
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	inst := filepath.Join(root, "cli-spark")
	b, err := os.ReadFile(filepath.Join(inst, "mycelium.toml"))
	if err != nil {
		t.Fatal(err)
	}
	m, err := manifest.Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	if m.State != "spark" || m.Tier != "focused" {
		t.Fatalf("%+v", m)
	}
	if rec.Called("gh") {
		t.Fatal("gh invoked")
	}
	if version.Version != "0.1.0-dev" {
		t.Fatalf("version=%q", version.Version)
	}
}

func TestCLINewIdeaMYCELIUM_OFFLINE(t *testing.T) {
	root := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	var stdout, stderr bytes.Buffer
	code := cli.Run(
		[]string{"mycelium", "new", "idea", "Env Offline"},
		&stdout, &stderr,
		cli.Deps{
			Clock:  clock.Fixed{T: time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)},
			Runner: rec,
			Getwd:  func() (string, error) { return root, nil },
			LookupEnv: func(k string) string {
				if k == "MYCELIUM_OFFLINE" {
					return "1"
				}
				return ""
			},
		},
	)
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	if rec.Called("gh") {
		t.Fatal("gh invoked under MYCELIUM_OFFLINE=1")
	}
}

func TestCLIGIT_DIRWithDirFlag(t *testing.T) {
	root := t.TempDir()
	sibling := filepath.Join(root, "sibling.git")
	target := filepath.Join(root, "cli-target")
	if err := os.Mkdir(sibling, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GIT_DIR", sibling)
	t.Setenv("GIT_WORK_TREE", root)

	rec := &execrun.Recording{Inner: execrun.Real{}}
	var stdout, stderr bytes.Buffer
	code := cli.Run(
		[]string{"mycelium", "new", "idea", "CLI Dir", "--offline", "--dir", target},
		&stdout, &stderr,
		cli.Deps{
			Clock:     clock.Fixed{T: time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)},
			Runner:    rec,
			Getwd:     func() (string, error) { return root, nil },
			LookupEnv: func(string) string { return "" },
		},
	)
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	if _, err := os.Stat(filepath.Join(target, ".git")); err != nil {
		t.Fatalf("target/.git missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(sibling, "HEAD")); err == nil {
		t.Fatal("sibling GIT_DIR should be untouched")
	}
	if rec.Called("gh") {
		t.Fatal("gh must not run")
	}
}

func TestCLIEmptySlugRefuses(t *testing.T) {
	root := t.TempDir()
	var stdout, stderr bytes.Buffer
	code := cli.Run(
		[]string{"mycelium", "new", "idea", "!!!", "--offline"},
		&stdout, &stderr,
		cli.Deps{
			Clock:     clock.Fixed{T: time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)},
			Runner:    &execrun.Recording{Inner: execrun.Real{}},
			Getwd:     func() (string, error) { return root, nil },
			LookupEnv: func(string) string { return "" },
		},
	)
	if code != 1 {
		t.Fatalf("exit %d want 1", code)
	}
	if !strings.Contains(stderr.String(), "slugify") && !strings.Contains(stderr.String(), "cannot slugify") {
		t.Fatalf("stderr=%q", stderr.String())
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("empty-slug must not create target, got %v", entries)
	}
}
