package scaffold_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/execrun"
	"github.com/robertguss/mycelium/internal/manifest"
	"github.com/robertguss/mycelium/internal/op"
	"github.com/robertguss/mycelium/internal/scaffold"
	"github.com/robertguss/mycelium/internal/version"
)

func TestNewIdeaOfflineScaffold(t *testing.T) {
	root := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	var stdout, stderr bytes.Buffer
	code := scaffold.Run(scaffold.Options{
		Name:    "Garden lighting",
		Offline: true,
		Cwd:     root,
		Argv:    []string{"new", "idea", "Garden lighting", "--offline"},
	}, scaffold.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 14, 12, 0, 0, 0, time.UTC)},
		Runner: rec,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	inst := filepath.Join(root, "garden-lighting")
	mustExist(t, inst, "README.md")
	mustExist(t, inst, "mycelium.toml")
	mustExist(t, inst, "log.md")
	mustExist(t, inst, "CONTEXT.md")
	mustExist(t, inst, "AGENTS.md")
	mustExist(t, inst, ".gitignore")
	mustExist(t, inst, ".agents/skills/mycelium-cli/SKILL.md")
	mustExist(t, inst, "program/README.md")
	mustExist(t, inst, "program/tiers/focused.toml")

	for _, forbidden := range []string{
		"framework", "Justfile", "scripts", "research-program.toml",
		"cmd", "internal", "go.mod", "go.sum",
	} {
		if _, err := os.Stat(filepath.Join(inst, forbidden)); !os.IsNotExist(err) {
			t.Fatalf("forbidden path present: %s (err=%v)", forbidden, err)
		}
	}

	b, err := os.ReadFile(filepath.Join(inst, "mycelium.toml"))
	if err != nil {
		t.Fatal(err)
	}
	m, err := manifest.Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	if m.State != "spark" || m.Tier != "focused" {
		t.Fatalf("state/tier = %s/%s", m.State, m.Tier)
	}
	if m.MethodologyVersion != "2.0.0" {
		t.Fatalf("methodology_version=%q", m.MethodologyVersion)
	}
	if m.GeneratedByCLIVersion != version.Version {
		t.Fatalf("cli version=%q want %q", m.GeneratedByCLIVersion, version.Version)
	}
	if m.CreatedDate != "2026-08-14" || m.UpdatedDate != "2026-08-14" {
		t.Fatalf("dates created=%q updated=%q", m.CreatedDate, m.UpdatedDate)
	}
	if m.GithubRepo != "" || m.Revisit != "" {
		t.Fatalf("github_repo/revisit should be empty")
	}

	gi, err := os.ReadFile(filepath.Join(inst, ".gitignore"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(gi)) != ".mycelium/lock" {
		t.Fatalf(".gitignore=%q", gi)
	}

	assertGitRepoNoCommits(t, inst)

	if rec.Called("gh") {
		t.Fatal("gh must never be invoked under --offline")
	}
	for _, c := range rec.Calls {
		if c.Name == "git" {
			for _, a := range c.Args {
				if a == "commit" || a == "add" {
					t.Fatalf("git must not %s: %+v", a, c)
				}
			}
		}
	}

	out := stdout.String()
	for _, want := range []string{
		"created ./garden-lighting",
		"state: spark",
		"tier: focused",
		`next: cd garden-lighting && mycelium new decision "First thought"`,
		"publish: mycelium publish",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("stdout missing %q\n%s", want, out)
		}
	}
}

func TestNewIdeaOfflineRefuseExisting(t *testing.T) {
	root := t.TempDir()
	inst := filepath.Join(root, "garden-lighting")
	if err := os.Mkdir(inst, 0o755); err != nil {
		t.Fatal(err)
	}
	rec := &execrun.Recording{Inner: execrun.Real{}}
	var stdout, stderr bytes.Buffer
	code := scaffold.Run(scaffold.Options{
		Name:    "Garden lighting",
		Offline: true,
		Cwd:     root,
	}, scaffold.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)},
		Runner: rec,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if code != 1 {
		t.Fatalf("exit %d want 1", code)
	}
	if !strings.Contains(stderr.String(), "already exists") {
		t.Fatalf("stderr=%q", stderr.String())
	}
	if rec.Called("gh") {
		t.Fatal("gh must not run")
	}
	if rec.Called("git") {
		t.Fatal("git must not run when refuse")
	}
}

func TestOfflineNeverExecsGh(t *testing.T) {
	root := t.TempDir()
	fake := &execrun.Fake{
		Paths: map[string]string{"git": "/usr/bin/git", "gh": "/usr/bin/gh"},
		RunFunc: func(ctx context.Context, name string, args []string, opts execrun.RunOpts) (execrun.Result, error) {
			if name == "gh" {
				t.Fatal("gh invoked")
			}
			if name == "git" && containsArg(args, "init") {
				if err := os.MkdirAll(filepath.Join(opts.Dir, ".git"), 0o755); err != nil {
					return execrun.Result{}, err
				}
				return execrun.Result{}, nil
			}
			return execrun.Result{ExitCode: 1, Stderr: []byte("unexpected")}, nil
		},
	}
	var stdout, stderr bytes.Buffer
	code := scaffold.Run(scaffold.Options{
		Name:    "No Gh",
		Offline: true,
		Cwd:     root,
	}, scaffold.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)},
		Runner: fake,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	if fake.Called("gh") {
		t.Fatal("gh recorded")
	}
}

func TestInstanceLacksMasterOnlyPaths(t *testing.T) {
	root := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	var stdout, stderr bytes.Buffer
	code := scaffold.Run(scaffold.Options{
		Name:    "Sparse Check",
		Offline: true,
		Cwd:     root,
	}, scaffold.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)},
		Runner: rec,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	inst := filepath.Join(root, "sparse-check")
	entries, err := os.ReadDir(inst)
	if err != nil {
		t.Fatal(err)
	}
	names := map[string]bool{}
	for _, e := range entries {
		names[e.Name()] = true
	}
	for _, bad := range []string{"framework", "Justfile", "scripts", "research-program.toml"} {
		if names[bad] {
			t.Fatalf("instance contains %s", bad)
		}
	}
}

func TestOfflinePublishContradictory(t *testing.T) {
	var stderr bytes.Buffer
	code := scaffold.Run(scaffold.Options{
		Name:    "X",
		Offline: true,
		Publish: true,
		Cwd:     t.TempDir(),
	}, scaffold.Deps{Stderr: &stderr, Stdout: io.Discard})
	if code != 1 {
		t.Fatalf("exit %d", code)
	}
	if !strings.Contains(stderr.String(), "contradictory") {
		t.Fatalf("stderr=%q", stderr.String())
	}
}

func TestTierStandardEmitsDirs(t *testing.T) {
	root := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	var stdout, stderr bytes.Buffer
	code := scaffold.Run(scaffold.Options{
		Name:    "Tiered",
		Offline: true,
		Tier:    "standard",
		Cwd:     root,
	}, scaffold.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)},
		Runner: rec,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	inst := filepath.Join(root, "tiered")
	for _, d := range []string{"decisions", "assumptions", "evidence", "questions", "risks"} {
		mustExist(t, inst, d+"/README.md")
	}
	b, _ := os.ReadFile(filepath.Join(inst, "mycelium.toml"))
	m, err := manifest.Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	if m.Tier != "standard" {
		t.Fatalf("tier=%q", m.Tier)
	}
}

func TestDirFlagExactDestination(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "custom-dest")
	rec := &execrun.Recording{Inner: execrun.Real{}}
	var stdout, stderr bytes.Buffer
	code := scaffold.Run(scaffold.Options{
		Name:    "Custom Dest",
		Dir:     dest,
		Offline: true,
		Cwd:     root,
	}, scaffold.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)},
		Runner: rec,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	mustExist(t, dest, "mycelium.toml")
	if _, err := os.Stat(filepath.Join(root, "custom-dest-slug")); !os.IsNotExist(err) {
		t.Fatal("must not append slug to --dir")
	}
}

func mustExist(t *testing.T, root, rel string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel))); err != nil {
		t.Fatalf("missing %s: %v", rel, err)
	}
}

func assertGitRepoNoCommits(t *testing.T, dir string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(dir, ".git")); err != nil {
		t.Fatalf("not a git repo: %v", err)
	}
	gitEnv := execrun.WithoutGitOverrides(os.Environ())
	cmd := exec.Command("git", "-C", dir, "rev-list", "--all")
	cmd.Env = gitEnv
	out, err := cmd.CombinedOutput()
	if err != nil {
		if len(bytes.TrimSpace(out)) == 0 {
			return
		}
		t.Fatalf("git rev-list: %v (%s)", err, out)
	}
	if len(bytes.TrimSpace(out)) != 0 {
		t.Fatalf("expected no commits, got %q", out)
	}
	stCmd := exec.Command("git", "-C", dir, "status", "--porcelain")
	stCmd.Env = gitEnv
	st, err := stCmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(st, []byte("mycelium.toml")) {
		t.Fatalf("expected untracked mycelium.toml in status %q", st)
	}
}

func TestScaffoldResumesPartialCommit(t *testing.T) {
	root := t.TempDir()
	name := "Resume Me"
	target := filepath.Join(root, "resume-me")
	if err := os.Mkdir(target, 0o755); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)
	argv := []string{"new", "idea", name, "--offline"}
	sess, err := op.Begin(target, op.Intent{
		Op:      "scaffold",
		Title:   name,
		LogLine: "2026-08-14\tscaffold\t-\tResume Me",
		Argv:    argv,
	}, now)
	if err != nil {
		t.Fatal(err)
	}
	files := []op.Staged{
		{RelTo: "README.md", Content: []byte("# Resume Me\n")},
		{RelTo: "CONTEXT.md", Content: []byte("# Glossary\n\n")},
		{RelTo: "mycelium.toml", Content: mustEncodeManifest(t, name, "resume-me")},
	}
	if err := sess.Stage(files); err != nil {
		t.Fatal(err)
	}
	if err := sess.CommitPartial(1); err != nil {
		t.Fatal(err)
	}
	if err := sess.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(target, "README.md")); err != nil {
		t.Fatal("expected partial rename of README.md")
	}
	if _, err := os.Stat(filepath.Join(target, "mycelium.toml")); !os.IsNotExist(err) {
		t.Fatal("mycelium.toml should still be staged")
	}

	rec := &execrun.Recording{Inner: execrun.Real{}}
	var stdout, stderr bytes.Buffer
	code := scaffold.Run(scaffold.Options{
		Name:    name,
		Offline: true,
		Cwd:     root,
		Argv:    argv,
	}, scaffold.Deps{
		Clock:  clock.Fixed{T: now},
		Runner: rec,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	mustExist(t, target, "mycelium.toml")
	mustExist(t, target, "CONTEXT.md")
	assertGitRepoNoCommits(t, target)
	if rec.Called("gh") {
		t.Fatal("gh invoked")
	}
}

func mustEncodeManifest(t *testing.T, ideaName, ideaSlug string) []byte {
	t.Helper()
	b, err := manifest.Encode(manifest.Manifest{
		SchemaVersion:         1,
		IdeaName:              ideaName,
		Slug:                  ideaSlug,
		State:                 "spark",
		Tier:                  "focused",
		MethodologyVersion:    "2.0.0",
		GeneratedByCLIVersion: version.Version,
		CreatedDate:           "2026-08-14",
		UpdatedDate:           "2026-08-14",
	})
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func containsArg(args []string, want string) bool {
	for _, a := range args {
		if a == want {
			return true
		}
	}
	return false
}

func TestGitInitIgnoresGIT_DIR(t *testing.T) {
	root := t.TempDir()
	sibling := filepath.Join(root, "hijacked.git")
	target := filepath.Join(root, "instance-dest")
	if err := os.Mkdir(sibling, 0o755); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadDir(sibling)
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("GIT_DIR", sibling)
	t.Setenv("GIT_WORK_TREE", root)

	rec := &execrun.Recording{Inner: execrun.Real{}}
	var stdout, stderr bytes.Buffer
	code := scaffold.Run(scaffold.Options{
		Name:    "Git Dir Trap",
		Dir:     target,
		Offline: true,
		Cwd:     root,
		Argv:    []string{"new", "idea", "Git Dir Trap", "--dir", target, "--offline"},
	}, scaffold.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)},
		Runner: rec,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	if _, err := os.Stat(filepath.Join(target, ".git")); err != nil {
		t.Fatalf("target/.git missing: %v", err)
	}
	after, err := os.ReadDir(sibling)
	if err != nil {
		t.Fatal(err)
	}
	if len(after) != len(before) {
		t.Fatalf("sibling GIT_DIR touched: before=%d after=%d", len(before), len(after))
	}
	if _, err := os.Stat(filepath.Join(sibling, "HEAD")); err == nil {
		t.Fatalf("git init wrote into GIT_DIR sibling %s", sibling)
	}
	assertGitRepoNoCommits(t, target)
	if rec.Called("gh") {
		t.Fatal("gh must not run")
	}
	for _, c := range rec.Calls {
		if c.Name != "git" {
			continue
		}
		joined := strings.Join(c.Args, " ")
		if !strings.Contains(joined, "--git-dir=.git") || !strings.Contains(joined, "--work-tree=.") {
			t.Fatalf("git args missing explicit dir flags: %v", c.Args)
		}
	}
}

func TestPruneMyceliumLeftoversAndContinue(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "leftover-dest")
	if err := os.MkdirAll(filepath.Join(target, ".mycelium", "stage", "x"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(target, ".mycelium", "lock"), []byte("pid=1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	rec := &execrun.Recording{Inner: execrun.Real{}}
	var stdout, stderr bytes.Buffer
	code := scaffold.Run(scaffold.Options{
		Name:    "Leftover Dest",
		Dir:     target,
		Offline: true,
		Cwd:     root,
	}, scaffold.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)},
		Runner: rec,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	mustExist(t, target, "mycelium.toml")
	mustExist(t, target, ".git")
	if _, err := os.Stat(filepath.Join(target, ".mycelium")); !os.IsNotExist(err) {
		t.Fatal("expected .mycelium pruned after success")
	}
}

func TestGitInitMissingDotGitRefuses(t *testing.T) {
	root := t.TempDir()
	fake := &execrun.Fake{
		RunFunc: func(ctx context.Context, name string, args []string, opts execrun.RunOpts) (execrun.Result, error) {
			return execrun.Result{}, nil
		},
	}
	var stdout, stderr bytes.Buffer
	code := scaffold.Run(scaffold.Options{
		Name:    "No Dot Git",
		Offline: true,
		Cwd:     root,
	}, scaffold.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)},
		Runner: fake,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if code != 1 {
		t.Fatalf("exit %d want 1", code)
	}
	if !strings.Contains(stderr.String(), "instance .git is missing") {
		t.Fatalf("stderr=%q", stderr.String())
	}
	inst := filepath.Join(root, "no-dot-git")
	if _, err := os.Stat(filepath.Join(inst, "mycelium.toml")); err != nil {
		t.Fatal("files should remain after git-init failure")
	}
	if _, err := os.Stat(filepath.Join(inst, ".git")); !os.IsNotExist(err) {
		t.Fatal(".git should be absent")
	}
}
