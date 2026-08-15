//go:build github_integration

package publish_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/execrun"
	"github.com/robertguss/mycelium/internal/manifest"
	"github.com/robertguss/mycelium/internal/publish"
	"github.com/robertguss/mycelium/internal/scaffold"
)

func TestGitHubPublishTopicCleanup(t *testing.T) {
	skipUnlessGhAuth(t)
	root := t.TempDir()
	runner := execrun.Real{}
	inst := scaffoldOfflineReal(t, root, runner)
	name := fmt.Sprintf("mycelium-ms101-%d", time.Now().Unix())
	t.Cleanup(func() { deleteFixtureRepo(t, runner, name) })

	out := publish.Run(publish.Options{
		Dir: inst, Mode: publish.ModeRequired, RequireGit: true, RepoName: name,
	}, publish.Deps{
		Clock:  clock.System(),
		Runner: runner,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	if out.Code != 0 {
		t.Fatalf("publish exit %d", out.Code)
	}
	m := mustManifest(t, inst)
	if m.GithubRepo == "" || !strings.HasSuffix(m.GithubRepo, "/"+name) {
		t.Fatalf("github_repo=%q", m.GithubRepo)
	}
	assertTopic(t, m.GithubRepo, "idea")

	var stdout bytes.Buffer
	second := publish.Run(publish.Options{
		Dir: inst, Mode: publish.ModeRequired, RequireGit: true, RepoName: name,
	}, publish.Deps{
		Clock:  clock.System(),
		Runner: runner,
		Stdout: &stdout,
		Stderr: os.Stderr,
	})
	if second.Code != 0 || second.Kind != publish.KindAlready {
		t.Fatalf("idempotent exit=%d kind=%v stdout=%q", second.Code, second.Kind, stdout.String())
	}
	if !strings.Contains(stdout.String(), "already published:") {
		t.Fatalf("stdout=%q", stdout.String())
	}
}

func TestGitHubCleanupOnInducedTopicFailure(t *testing.T) {
	skipUnlessGhAuth(t)
	root := t.TempDir()
	base := execrun.Real{}
	name := fmt.Sprintf("mycelium-ms101-%d", time.Now().Unix())
	t.Cleanup(func() { deleteFixtureRepo(t, base, name) })

	wrapped := &topicFailRunner{real: base}
	inst := scaffoldOfflineReal(t, root, wrapped)
	out := publish.Run(publish.Options{
		Dir: inst, Mode: publish.ModeRequired, RequireGit: true, RepoName: name,
	}, publish.Deps{
		Clock:  clock.System(),
		Runner: wrapped,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	if out.Code != 1 {
		t.Fatalf("exit %d want 1", out.Code)
	}
	if repoExists(t, base, name) {
		t.Fatalf("fixture %s should be deleted after induced failure", name)
	}
}

type topicFailRunner struct {
	real execrun.Real
}

func (f *topicFailRunner) LookPath(name string) (string, error) {
	return f.real.LookPath(name)
}

func (f *topicFailRunner) Run(ctx context.Context, name string, args []string, opts execrun.RunOpts) (execrun.Result, error) {
	if name == "gh" && len(args) >= 2 && args[0] == "repo" && args[1] == "edit" {
		return execrun.Result{ExitCode: 1, Stderr: []byte("induced topic failure")}, nil
	}
	return f.real.Run(ctx, name, args, opts)
}

func skipUnlessGhAuth(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("gh"); err != nil {
		t.Skip("skipped: no credentials")
	}
	cmd := exec.Command("gh", "auth", "status")
	if err := cmd.Run(); err != nil {
		t.Skip("skipped: no credentials")
	}
}

func scaffoldOfflineReal(t *testing.T, root string, runner execrun.Runner) string {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := scaffold.Run(scaffold.Options{
		Name: "Gh Integration", Offline: true, Cwd: root,
	}, scaffold.Deps{
		Clock:  clock.System(),
		Runner: runner,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr.String())
	}
	return filepath.Join(root, "gh-integration")
}

func mustManifest(t *testing.T, inst string) manifest.Manifest {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(inst, "mycelium.toml"))
	if err != nil {
		t.Fatal(err)
	}
	m, err := manifest.Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func assertTopic(t *testing.T, ownerRepo, want string) {
	t.Helper()
	cmd := exec.Command("gh", "api", "repos/"+ownerRepo, "--jq", ".topics")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("topics: %v: %s", err, out)
	}
	if !strings.Contains(string(out), `"`+want+`"`) {
		t.Fatalf("topics=%s want %q", out, want)
	}
}

func deleteFixtureRepo(t *testing.T, runner execrun.Runner, name string) {
	t.Helper()
	loginOut, err := runner.Run(context.Background(), "gh", []string{"api", "user", "--jq", ".login"}, execrun.RunOpts{})
	if err != nil || loginOut.ExitCode != 0 {
		return
	}
	login := strings.TrimSpace(string(loginOut.Stdout))
	target := login + "/" + name
	_, _ = runner.Run(context.Background(), "gh", []string{"repo", "delete", target, "--yes"}, execrun.RunOpts{})
}

func repoExists(t *testing.T, runner execrun.Runner, name string) bool {
	t.Helper()
	loginOut, err := runner.Run(context.Background(), "gh", []string{"api", "user", "--jq", ".login"}, execrun.RunOpts{})
	if err != nil || loginOut.ExitCode != 0 {
		return false
	}
	login := strings.TrimSpace(string(loginOut.Stdout))
	res, err := runner.Run(context.Background(), "gh", []string{"api", "repos/" + login + "/" + name}, execrun.RunOpts{})
	if err != nil {
		return false
	}
	return res.ExitCode == 0
}
