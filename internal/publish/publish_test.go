package publish_test

import (
	"bytes"
	"context"
	"io"
	"os"
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

func TestOfflineNeverExecsGh(t *testing.T) {
	root := t.TempDir()
	scaffoldFake := &execrun.Fake{
		Paths: map[string]string{"git": "/usr/bin/git"},
		RunFunc: func(ctx context.Context, name string, args []string, opts execrun.RunOpts) (execrun.Result, error) {
			if name == "git" && containsArg(args, "init") {
				if err := os.MkdirAll(filepath.Join(opts.Dir, ".git"), 0o755); err != nil {
					return execrun.Result{}, err
				}
				return execrun.Result{}, nil
			}
			t.Fatalf("unexpected scaffold exec %s %v", name, args)
			return execrun.Result{}, nil
		},
	}
	inst := scaffoldOffline(t, root, scaffoldFake)
	pubFake := &execrun.Fake{
		Paths: map[string]string{"gh": "/usr/bin/gh", "git": "/usr/bin/git"},
		RunFunc: func(ctx context.Context, name string, args []string, opts execrun.RunOpts) (execrun.Result, error) {
			t.Fatalf("unexpected exec %s %v", name, args)
			return execrun.Result{}, nil
		},
	}
	out := publish.Run(publish.Options{
		Dir:        inst,
		Mode:       publish.ModeOffline,
		RequireGit: true,
	}, publish.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)},
		Runner: pubFake,
		Stdout: io.Discard,
		Stderr: io.Discard,
	})
	if out.Code != 1 {
		t.Fatalf("exit %d want 1", out.Code)
	}
	if pubFake.Called("gh") {
		t.Fatal("gh must not run offline")
	}
}

func TestOfflinePublishContradictoryViaScaffold(t *testing.T) {
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

func TestPublishHappyPathFake(t *testing.T) {
	root := t.TempDir()
	state := &ghState{login: "alice", topics: map[string][]string{}}
	fake := newPublishFake(t, state)
	inst := scaffoldOffline(t, root, fake)

	out := publish.Run(publish.Options{
		Dir:        inst,
		Mode:       publish.ModeRequired,
		RequireGit: true,
		RepoName:   "mycelium-ms101-1001",
	}, publish.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)},
		Runner: fake,
		Stdout: io.Discard,
		Stderr: io.Discard,
	})
	if out.Code != 0 {
		t.Fatalf("exit %d kind=%v", out.Code, out.Kind)
	}
	if out.OwnerRepo != "alice/mycelium-ms101-1001" {
		t.Fatalf("ownerRepo=%q", out.OwnerRepo)
	}
	m := readManifest(t, inst)
	if m.GithubRepo != "alice/mycelium-ms101-1001" {
		t.Fatalf("github_repo=%q", m.GithubRepo)
	}
	logBytes, err := os.ReadFile(filepath.Join(inst, "log.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(logBytes), "publish\t-\tgithub.com/alice/mycelium-ms101-1001") {
		t.Fatalf("log=%q", logBytes)
	}
	if state.originURL == "" {
		t.Fatal("expected origin remote")
	}
	if !contains(state.topics["alice/mycelium-ms101-1001"], "idea") {
		t.Fatalf("topics=%v", state.topics)
	}
}

func TestPublishIdempotentSecond(t *testing.T) {
	root := t.TempDir()
	state := &ghState{login: "alice", topics: map[string][]string{}}
	fake := newPublishFake(t, state)
	inst := scaffoldOffline(t, root, fake)

	first := publish.Run(publish.Options{
		Dir: inst, Mode: publish.ModeRequired, RequireGit: true, RepoName: "mycelium-ms101-2002",
	}, publish.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)},
		Runner: fake, Stdout: io.Discard, Stderr: io.Discard,
	})
	if first.Code != 0 {
		t.Fatalf("first exit %d", first.Code)
	}
	creates := countRuns(fake, "gh", "repo", "create")
	var stdout bytes.Buffer
	second := publish.Run(publish.Options{
		Dir: inst, Mode: publish.ModeRequired, RequireGit: true, RepoName: "mycelium-ms101-2002",
	}, publish.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)},
		Runner: fake, Stdout: &stdout, Stderr: io.Discard,
	})
	if second.Code != 0 || second.Kind != publish.KindAlready {
		t.Fatalf("second exit=%d kind=%v stdout=%q", second.Code, second.Kind, stdout.String())
	}
	if !strings.Contains(stdout.String(), "already published: alice/mycelium-ms101-2002") {
		t.Fatalf("stdout=%q", stdout.String())
	}
	if countRuns(fake, "gh", "repo", "create") != creates {
		t.Fatal("second publish must not create again")
	}
}

func TestCleanupOnTopicFailureFixture(t *testing.T) {
	root := t.TempDir()
	state := &ghState{login: "alice", topics: map[string][]string{}, failTopic: true}
	fake := newPublishFake(t, state)
	inst := scaffoldOffline(t, root, fake)

	var stderr bytes.Buffer
	out := publish.Run(publish.Options{
		Dir: inst, Mode: publish.ModeRequired, RequireGit: true, RepoName: "mycelium-ms101-3003",
	}, publish.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)},
		Runner: fake, Stdout: io.Discard, Stderr: &stderr,
	})
	if out.Code != 1 {
		t.Fatalf("exit %d want 1", out.Code)
	}
	if !state.deleted {
		t.Fatal("fixture repo must be deleted after topic failure")
	}
	if !strings.Contains(stderr.String(), "add-topic failed") {
		t.Fatalf("stderr=%q", stderr.String())
	}
}

func TestNoAutoDeleteUserSlug(t *testing.T) {
	root := t.TempDir()
	state := &ghState{login: "alice", topics: map[string][]string{}, failTopic: true}
	fake := newPublishFake(t, state)
	inst := scaffoldOffline(t, root, fake)

	var stderr bytes.Buffer
	out := publish.Run(publish.Options{
		Dir: inst, Mode: publish.ModeRequired, RequireGit: true, // uses slug from manifest
	}, publish.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)},
		Runner: fake, Stdout: io.Discard, Stderr: &stderr,
	})
	if out.Code != 1 {
		t.Fatalf("exit %d", out.Code)
	}
	if state.deleted {
		t.Fatal("user slug must not auto-delete")
	}
	if !strings.Contains(stderr.String(), "https://github.com/alice/") {
		t.Fatalf("stderr should teach URL: %q", stderr.String())
	}
}

func TestOptionalSkipsWhenUnauthenticated(t *testing.T) {
	root := t.TempDir()
	fake := &execrun.Fake{
		Paths: map[string]string{"gh": "/usr/bin/gh", "git": "/bin/git"},
		RunFunc: func(ctx context.Context, name string, args []string, opts execrun.RunOpts) (execrun.Result, error) {
			if name == "gh" && len(args) >= 2 && args[0] == "auth" {
				return execrun.Result{ExitCode: 1, Stderr: []byte("not logged in")}, nil
			}
			if name == "git" && containsArg(args, "init") {
				_ = os.MkdirAll(filepath.Join(opts.Dir, ".git"), 0o755)
				return execrun.Result{}, nil
			}
			return execrun.Result{}, nil
		},
	}
	var stdout bytes.Buffer
	code := scaffold.Run(scaffold.Options{
		Name: "Optional Skip", Cwd: root,
	}, scaffold.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)},
		Runner: fake, Stdout: &stdout, Stderr: io.Discard,
	})
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	if !strings.Contains(stdout.String(), "publish: mycelium publish") {
		t.Fatalf("stdout=%q", stdout.String())
	}
	if countRuns(fake, "gh", "repo", "create") != 0 {
		t.Fatal("must not create when unauthenticated")
	}
}

type ghState struct {
	login     string
	topics    map[string][]string
	repos     map[string]bool
	originURL string
	failTopic bool
	deleted   bool
}

func newPublishFake(t *testing.T, state *ghState) *execrun.Fake {
	t.Helper()
	if state.repos == nil {
		state.repos = map[string]bool{}
	}
	return &execrun.Fake{
		Paths: map[string]string{"gh": "/usr/bin/gh", "git": "/bin/git"},
		RunFunc: func(ctx context.Context, name string, args []string, opts execrun.RunOpts) (execrun.Result, error) {
			switch name {
			case "git":
				return handleGit(t, state, args, opts)
			case "gh":
				return handleGh(t, state, args)
			default:
				t.Fatalf("unexpected binary %s", name)
				return execrun.Result{}, nil
			}
		},
	}
}

func handleGit(t *testing.T, state *ghState, args []string, opts execrun.RunOpts) (execrun.Result, error) {
	t.Helper()
	if containsArg(args, "init") {
		if err := os.MkdirAll(filepath.Join(opts.Dir, ".git"), 0o755); err != nil {
			return execrun.Result{}, err
		}
		return execrun.Result{}, nil
	}
	if len(args) >= 2 && args[0] == "remote" && args[1] == "get-url" {
		if state.originURL == "" {
			return execrun.Result{ExitCode: 2, Stderr: []byte("no such remote")}, nil
		}
		return execrun.Result{Stdout: []byte(state.originURL + "\n")}, nil
	}
	if len(args) >= 4 && args[0] == "remote" && args[1] == "add" && args[2] == "origin" {
		state.originURL = args[3]
		return execrun.Result{}, nil
	}
	return execrun.Result{ExitCode: 1, Stderr: []byte("unexpected git")}, nil
}

func handleGh(t *testing.T, state *ghState, args []string) (execrun.Result, error) {
	t.Helper()
	if len(args) >= 2 && args[0] == "auth" && args[1] == "status" {
		return execrun.Result{}, nil
	}
	if len(args) >= 4 && args[0] == "api" && args[1] == "user" {
		return execrun.Result{Stdout: []byte(state.login + "\n")}, nil
	}
	if len(args) >= 2 && args[0] == "api" && strings.HasPrefix(args[1], "repos/") {
		or := strings.TrimPrefix(args[1], "repos/")
		if !state.repos[or] {
			return execrun.Result{ExitCode: 1, Stderr: []byte("not found")}, nil
		}
		topics := state.topics[or]
		var b strings.Builder
		b.WriteByte('[')
		for i, topic := range topics {
			if i > 0 {
				b.WriteByte(',')
			}
			b.WriteByte('"')
			b.WriteString(topic)
			b.WriteByte('"')
		}
		b.WriteByte(']')
		return execrun.Result{Stdout: []byte(b.String())}, nil
	}
	if len(args) >= 3 && args[0] == "repo" && args[1] == "create" {
		name := args[2]
		or := state.login + "/" + name
		state.repos[or] = true
		state.topics[or] = nil
		return execrun.Result{Stdout: []byte("https://github.com/" + or + "\n")}, nil
	}
	if len(args) >= 5 && args[0] == "repo" && args[1] == "edit" {
		or := args[2]
		if state.failTopic {
			return execrun.Result{ExitCode: 1, Stderr: []byte("topic boom")}, nil
		}
		state.topics[or] = append(state.topics[or], "idea")
		return execrun.Result{}, nil
	}
	if len(args) >= 3 && args[0] == "repo" && args[1] == "delete" {
		state.deleted = true
		delete(state.repos, args[2])
		return execrun.Result{}, nil
	}
	t.Fatalf("unexpected gh %v", args)
	return execrun.Result{}, nil
}

func scaffoldOffline(t *testing.T, root string, runner execrun.Runner) string {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := scaffold.Run(scaffold.Options{
		Name: "Pub Fixture", Offline: true, Cwd: root,
	}, scaffold.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)},
		Runner: runner, Stdout: &stdout, Stderr: &stderr,
	})
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr.String())
	}
	inst := filepath.Join(root, "pub-fixture")
	return inst
}

func readManifest(t *testing.T, inst string) manifest.Manifest {
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

func contains(ss []string, want string) bool {
	for _, s := range ss {
		if s == want {
			return true
		}
	}
	return false
}

func containsArg(args []string, want string) bool {
	for _, a := range args {
		if a == want {
			return true
		}
	}
	return false
}

func countRuns(f *execrun.Fake, name string, prefix ...string) int {
	n := 0
	for _, c := range f.Calls {
		if c.Kind != "run" || c.Name != name {
			continue
		}
		if len(c.Args) < len(prefix) {
			continue
		}
		ok := true
		for i, p := range prefix {
			if c.Args[i] != p {
				ok = false
				break
			}
		}
		if ok {
			n++
		}
	}
	return n
}
