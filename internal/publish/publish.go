// Package publish implements mycelium publish and new-idea publish half.
package publish

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/robertguss/mycelium/internal/check"
	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/execrun"
	"github.com/robertguss/mycelium/internal/journal"
	"github.com/robertguss/mycelium/internal/logfmt"
	"github.com/robertguss/mycelium/internal/manifest"
	"github.com/robertguss/mycelium/internal/op"
	"github.com/robertguss/mycelium/internal/teach"
)

// Mode selects auth + network behavior.
type Mode int

const (
	// ModeOffline never execs gh (MYCELIUM_OFFLINE / --offline).
	ModeOffline Mode = iota
	// ModeOptional publishes when gh is authenticated; otherwise skips.
	ModeOptional
	// ModeRequired fails when gh is missing or unauthenticated (--publish / publish cmd).
	ModeRequired
)

// Kind is the publish outcome for callers (scaffold banner).
type Kind int

const (
	KindFailed Kind = iota
	KindSkipped
	KindAlready
	KindPublished
)

// Options for one publish run.
type Options struct {
	Dir        string // start path for FindRoot; empty → Cwd
	Cwd        string
	Mode       Mode
	RepoName   string // empty → manifest slug; tests set mycelium-ms101-<unix>
	Argv       []string
	Quiet      bool // suppress stdout; scaffold owns the success banner
	RequireGit bool // publish cmd requires .git; scaffold already inited
}

// Deps are injectable collaborators.
type Deps struct {
	Clock  clock.Clock
	Runner execrun.Runner
	Stdout io.Writer
	Stderr io.Writer
}

// Outcome is exit code plus publish kind for callers.
type Outcome struct {
	Code      int
	Kind      Kind
	OwnerRepo string // "owner/name" when known
}

var fixtureNameRE = regexp.MustCompile(`^mycelium-ms101-[0-9]+$`)

// Run publishes the instance at Dir/Cwd. Exit 0 on success or optional skip.
func Run(opts Options, deps Deps) Outcome {
	if deps.Clock == nil {
		deps.Clock = clock.System()
	}
	if deps.Runner == nil {
		deps.Runner = execrun.Real{}
	}
	if deps.Stdout == nil {
		deps.Stdout = io.Discard
	}
	if deps.Stderr == nil {
		deps.Stderr = io.Discard
	}

	if opts.Mode == ModeOffline {
		return Outcome{
			Code: teach.Write(deps.Stderr,
				"publish refused while offline",
				"command-flags",
				"framework/phases/PHASE-01-implementation-brief.md",
				"omit --offline / MYCELIUM_OFFLINE=1, or run mycelium publish when network is allowed",
			),
			Kind: KindFailed,
		}
	}

	start := opts.Dir
	if start == "" {
		start = opts.Cwd
	}
	if start == "" {
		wd, err := os.Getwd()
		if err != nil {
			return fail(deps, fmt.Sprintf("cannot resolve cwd: %v", err),
				"command-flags",
				"framework/phases/PHASE-01-implementation-brief.md",
				"retry from a readable working directory")
		}
		start = wd
	}

	root, err := check.FindRoot(start)
	if err != nil {
		return fail(deps, "not a mycelium instance (no mycelium.toml found)",
			"instance-root",
			"framework/phases/PHASE-01-implementation-brief.md",
			"run from an instance directory, or pass --dir PATH")
	}

	manBytes, err := os.ReadFile(filepath.Join(root, "mycelium.toml"))
	if err != nil {
		return fail(deps, fmt.Sprintf("cannot read mycelium.toml: %v", err),
			"manifest",
			"program/contracts/manifest.md",
			"restore mycelium.toml at the instance root")
	}
	m, err := manifest.Parse(manBytes)
	if err != nil {
		return fail(deps, fmt.Sprintf("invalid mycelium.toml: %v", err),
			"manifest",
			"program/contracts/manifest.md",
			"fix mycelium.toml required fields and retry")
	}

	logBytes, err := os.ReadFile(filepath.Join(root, "log.md"))
	if err != nil {
		return fail(deps, fmt.Sprintf("cannot read log.md: %v", err),
			"log",
			"framework/phases/PHASE-01-implementation-brief.md",
			"restore log.md at the instance root")
	}

	if opts.RequireGit {
		if _, err := os.Stat(filepath.Join(root, ".git")); err != nil {
			return fail(deps, "instance has no .git directory",
				"git-init",
				"program/contracts/operation-protocol.md",
				"run: git --git-dir=.git --work-tree=. init -b main")
		}
	}

	ctx := context.Background()
	if _, err := deps.Runner.LookPath("gh"); err != nil {
		if opts.Mode == ModeOptional {
			return Outcome{Code: 0, Kind: KindSkipped}
		}
		return fail(deps, "gh is not on PATH",
			"publish",
			"framework/phases/PHASE-01-implementation-brief.md",
			"install GitHub CLI (gh) and authenticate with: gh auth login")
	}

	auth, err := deps.Runner.Run(ctx, "gh", []string{"auth", "status"}, execrun.RunOpts{Dir: root})
	if err != nil || auth.ExitCode != 0 {
		if opts.Mode == ModeOptional {
			return Outcome{Code: 0, Kind: KindSkipped}
		}
		detail := strings.TrimSpace(string(auth.Stderr))
		if detail == "" {
			detail = strings.TrimSpace(string(auth.Stdout))
		}
		if detail == "" && err != nil {
			detail = err.Error()
		}
		return fail(deps, fmt.Sprintf("gh is not authenticated: %s", detail),
			"publish",
			"framework/phases/PHASE-01-implementation-brief.md",
			"run: gh auth login")
	}

	snap, err := observe(ctx, deps.Runner, root, m)
	if err != nil {
		return fail(deps, err.Error(),
			"publish",
			"framework/phases/PHASE-01-implementation-brief.md",
			"fix git remotes and gh access, then retry")
	}

	if snap.complete() {
		if !opts.Quiet {
			fmt.Fprintf(deps.Stdout, "already published: %s\n", snap.ownerRepo)
		}
		return Outcome{Code: 0, Kind: KindAlready, OwnerRepo: snap.ownerRepo}
	}

	repoName := strings.TrimSpace(opts.RepoName)
	if repoName == "" {
		repoName = m.Slug
	}
	if repoName == "" {
		return fail(deps, "cannot determine repository name",
			"publish",
			"framework/phases/PHASE-01-implementation-brief.md",
			"set slug in mycelium.toml")
	}

	now := deps.Clock.Now().UTC()
	date := clock.Date(now)
	argv := opts.Argv
	if len(argv) == 0 {
		argv = []string{"publish"}
		if opts.Dir != "" {
			argv = append(argv, "--dir", opts.Dir)
		}
	}

	logLine := logfmt.Line(date, "publish", "-", "pending")
	sess, err := op.Begin(root, op.Intent{
		Op:         "publish",
		Title:      m.IdeaName,
		OriginalID: "",
		LogLine:    logLine,
		Argv:       argv,
	}, now)
	if err != nil {
		if errors.Is(err, op.ErrJournalMismatch) {
			return fail(deps, "leftover journal for a different operation",
				"operation-protocol",
				"program/contracts/operation-protocol.md",
				"re-run the original command to complete, or mycelium check --abort-journal to roll back")
		}
		if errors.Is(err, op.ErrLocked) {
			return fail(deps, fmt.Sprintf("lock held by another process: %v", err),
				"operation-protocol",
				"program/contracts/operation-protocol.md",
				"wait for the other process to finish; do not force the lock")
		}
		return fail(deps, fmt.Sprintf("operation begin failed: %v", err),
			"operation-protocol",
			"program/contracts/operation-protocol.md",
			"re-run the same command, or mycelium check --abort-journal")
	}

	createdName := ""
	ownerRepo := snap.ownerRepo
	if ownerRepo == "" && sess.OriginalID() != "" {
		ownerRepo = sess.OriginalID()
		createdName = repoLeaf(ownerRepo)
	}

	if ownerRepo == "" {
		login, loginErr := ghLogin(ctx, deps.Runner, root)
		if loginErr != nil {
			rollbackOrClose(sess)
			return fail(deps, loginErr.Error(),
				"publish",
				"framework/phases/PHASE-01-implementation-brief.md",
				"run: gh auth login")
		}
		desc := "idea: " + m.IdeaName
		if len(desc) > 80 {
			desc = desc[:80]
		}
		createArgs := []string{
			"repo", "create", repoName,
			"--private",
			"--description", desc,
		}
		res, runErr := deps.Runner.Run(ctx, "gh", createArgs, execrun.RunOpts{Dir: root})
		if runErr != nil || res.ExitCode != 0 {
			rollbackOrClose(sess)
			detail := strings.TrimSpace(string(res.Stderr))
			if detail == "" {
				detail = strings.TrimSpace(string(res.Stdout))
			}
			if detail == "" && runErr != nil {
				detail = runErr.Error()
			}
			return fail(deps, fmt.Sprintf("gh repo create failed: %s", detail),
				"publish",
				"framework/phases/PHASE-01-implementation-brief.md",
				"fix the repository name is free and gh has create permission")
		}
		createdName = repoName
		ownerRepo = login + "/" + repoName
		if url := firstHTTPURL(string(res.Stdout)); url != "" {
			if or, ok := ownerRepoFromURL(url); ok {
				ownerRepo = or
			}
		}
		sess.Journal().OriginalID = ownerRepo
		sess.Journal().LogLine = logfmt.Line(date, "publish", "-", "github.com/"+ownerRepo)
		if saveErr := journal.Save(root, sess.Journal()); saveErr != nil {
			cleanupCreated(ctx, deps, root, createdName, ownerRepo)
			rollbackOrClose(sess)
			return fail(deps, fmt.Sprintf("cannot record created repo in journal: %v", saveErr),
				"operation-protocol",
				"program/contracts/operation-protocol.md",
				fmt.Sprintf("repo may exist at https://github.com/%s; re-run mycelium publish or delete it manually", ownerRepo))
		}
	} else {
		createdName = repoLeaf(ownerRepo)
	}

	if !snap.hasOrigin {
		url := "https://github.com/" + ownerRepo + ".git"
		res, runErr := deps.Runner.Run(ctx, "git", []string{"remote", "add", "origin", url}, execrun.RunOpts{Dir: root})
		if runErr != nil || res.ExitCode != 0 {
			cleanupCreated(ctx, deps, root, createdName, ownerRepo)
			rollbackOrClose(sess)
			detail := strings.TrimSpace(string(res.Stderr))
			if detail == "" && runErr != nil {
				detail = runErr.Error()
			}
			return fail(deps, fmt.Sprintf("git remote add failed: %s", detail),
				"publish",
				"framework/phases/PHASE-01-implementation-brief.md",
				fmt.Sprintf("add origin manually: git remote add origin https://github.com/%s.git", ownerRepo))
		}
	}

	if !snap.hasTopic {
		res, runErr := deps.Runner.Run(ctx, "gh", []string{"repo", "edit", ownerRepo, "--add-topic", "idea"}, execrun.RunOpts{Dir: root})
		if runErr != nil || res.ExitCode != 0 {
			cleanupCreated(ctx, deps, root, createdName, ownerRepo)
			rollbackOrClose(sess)
			detail := strings.TrimSpace(string(res.Stderr))
			if detail == "" && runErr != nil {
				detail = runErr.Error()
			}
			msg := fmt.Sprintf("gh repo edit --add-topic failed: %s", detail)
			fix := "retry: gh repo edit " + ownerRepo + " --add-topic idea"
			if !fixtureNameRE.MatchString(createdName) {
				fix = fmt.Sprintf("repo left at https://github.com/%s; %s", ownerRepo, fix)
			}
			return fail(deps, msg, "publish", "framework/phases/PHASE-01-implementation-brief.md", fix)
		}
	}

	if len(sess.Journal().Renames) == 0 {
		m.GithubRepo = ownerRepo
		m.UpdatedDate = date
		manOut, encErr := manifest.Encode(m)
		if encErr != nil {
			cleanupCreated(ctx, deps, root, createdName, ownerRepo)
			rollbackOrClose(sess)
			return fail(deps, fmt.Sprintf("cannot encode manifest: %v", encErr),
				"manifest",
				"program/contracts/manifest.md",
				"report this as a CLI bug")
		}
		finalLine := logfmt.Line(date, "publish", "-", "github.com/"+ownerRepo)
		sess.Journal().LogLine = finalLine
		files := []op.Staged{
			{RelTo: "log.md", Content: appendLogLine(logBytes, finalLine)},
			{RelTo: "mycelium.toml", Content: manOut},
		}
		if stageErr := sess.Stage(files); stageErr != nil {
			cleanupCreated(ctx, deps, root, createdName, ownerRepo)
			rollbackOrClose(sess)
			return fail(deps, fmt.Sprintf("stage failed: %v", stageErr),
				"operation-protocol",
				"program/contracts/operation-protocol.md",
				"re-run the same command after freeing disk space")
		}
	}

	if err := sess.Commit(); err != nil {
		_ = sess.Close()
		return fail(deps, fmt.Sprintf("commit failed: %v", err),
			"operation-protocol",
			"program/contracts/operation-protocol.md",
			"re-run the same command to complete, or mycelium check --abort-journal to roll back")
	}
	_ = os.Remove(filepath.Join(root, ".mycelium"))

	if !opts.Quiet {
		fmt.Fprintf(deps.Stdout, "published: %s (topic: idea)\n", ownerRepo)
	}
	return Outcome{Code: 0, Kind: KindPublished, OwnerRepo: ownerRepo}
}

type snapshot struct {
	hasOrigin  bool
	hasTopic   bool
	ownerRepo  string
	githubRepo string
}

func (s snapshot) complete() bool {
	return s.hasOrigin && s.hasTopic && s.githubRepo != "" && s.ownerRepo != "" && s.githubRepo == s.ownerRepo
}

func observe(ctx context.Context, runner execrun.Runner, root string, m manifest.Manifest) (snapshot, error) {
	s := snapshot{githubRepo: strings.TrimSpace(m.GithubRepo)}
	res, err := runner.Run(ctx, "git", []string{"remote", "get-url", "origin"}, execrun.RunOpts{Dir: root})
	if err == nil && res.ExitCode == 0 {
		url := strings.TrimSpace(string(res.Stdout))
		if url != "" {
			s.hasOrigin = true
			if or, ok := ownerRepoFromURL(url); ok {
				s.ownerRepo = or
			}
		}
	}
	if s.ownerRepo == "" && s.githubRepo != "" {
		s.ownerRepo = s.githubRepo
	}
	if s.ownerRepo != "" {
		topics, terr := listTopics(ctx, runner, root, s.ownerRepo)
		if terr != nil {
			return s, terr
		}
		for _, t := range topics {
			if t == "idea" {
				s.hasTopic = true
				break
			}
		}
	}
	return s, nil
}

func listTopics(ctx context.Context, runner execrun.Runner, root, ownerRepo string) ([]string, error) {
	res, err := runner.Run(ctx, "gh", []string{
		"api", "repos/" + ownerRepo, "--jq", ".topics",
	}, execrun.RunOpts{Dir: root})
	if err != nil {
		return nil, fmt.Errorf("cannot read repo topics: %w", err)
	}
	if res.ExitCode != 0 {
		return nil, nil
	}
	raw := strings.TrimSpace(string(res.Stdout))
	raw = strings.TrimPrefix(raw, "[")
	raw = strings.TrimSuffix(raw, "]")
	if raw == "" {
		return nil, nil
	}
	var out []string
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		part = strings.Trim(part, `"`)
		if part != "" {
			out = append(out, part)
		}
	}
	return out, nil
}

func ghLogin(ctx context.Context, runner execrun.Runner, root string) (string, error) {
	res, err := runner.Run(ctx, "gh", []string{"api", "user", "--jq", ".login"}, execrun.RunOpts{Dir: root})
	if err != nil {
		return "", fmt.Errorf("cannot resolve authenticated user: %v", err)
	}
	if res.ExitCode != 0 {
		detail := strings.TrimSpace(string(res.Stderr))
		if detail == "" {
			detail = strings.TrimSpace(string(res.Stdout))
		}
		return "", fmt.Errorf("cannot resolve authenticated user: %s", detail)
	}
	login := strings.TrimSpace(string(res.Stdout))
	if login == "" {
		return "", errors.New("cannot resolve authenticated user: empty login")
	}
	return login, nil
}

func cleanupCreated(ctx context.Context, deps Deps, root, createdName, ownerRepo string) {
	if !fixtureNameRE.MatchString(createdName) {
		if ownerRepo != "" {
			fmt.Fprintf(deps.Stderr, "mycelium: publish failed after create; repo left at https://github.com/%s\n", ownerRepo)
		}
		return
	}
	target := ownerRepo
	if target == "" {
		target = createdName
	}
	_, _ = deps.Runner.Run(ctx, "gh", []string{"repo", "delete", target, "--yes"}, execrun.RunOpts{Dir: root})
}

func ownerRepoFromURL(url string) (string, bool) {
	u := strings.TrimSpace(url)
	u = strings.TrimSuffix(u, ".git")
	u = strings.TrimSuffix(u, "/")
	switch {
	case strings.HasPrefix(u, "https://github.com/"):
		rest := strings.TrimPrefix(u, "https://github.com/")
		parts := strings.Split(rest, "/")
		if len(parts) >= 2 && parts[0] != "" && parts[1] != "" {
			return parts[0] + "/" + parts[1], true
		}
	case strings.HasPrefix(u, "git@github.com:"):
		rest := strings.TrimPrefix(u, "git@github.com:")
		parts := strings.Split(rest, "/")
		if len(parts) >= 2 && parts[0] != "" && parts[1] != "" {
			return parts[0] + "/" + parts[1], true
		}
	}
	return "", false
}

func firstHTTPURL(s string) string {
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "https://") {
			return line
		}
	}
	return ""
}

func repoLeaf(ownerRepo string) string {
	if i := strings.LastIndex(ownerRepo, "/"); i >= 0 && i+1 < len(ownerRepo) {
		return ownerRepo[i+1:]
	}
	return ownerRepo
}

func appendLogLine(existing []byte, line string) []byte {
	s := string(existing)
	if !strings.HasSuffix(s, "\n") && s != "" {
		s += "\n"
	}
	s += line + "\n"
	return []byte(s)
}

func rollbackOrClose(sess *op.Session) {
	if err := sess.Rollback(); errors.Is(err, op.ErrPartialCommit) {
		_ = sess.Close()
	}
}

func fail(deps Deps, what, convention, contract, fix string) Outcome {
	return Outcome{
		Code: teach.Write(deps.Stderr, what, convention, contract, fix),
		Kind: KindFailed,
	}
}
