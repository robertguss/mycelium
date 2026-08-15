// Package scaffold implements mycelium new idea (local emit + git init).
package scaffold

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/robertguss/mycelium/internal/clock"
	myceliumembed "github.com/robertguss/mycelium/internal/embed"
	"github.com/robertguss/mycelium/internal/execrun"
	"github.com/robertguss/mycelium/internal/journal"
	"github.com/robertguss/mycelium/internal/logfmt"
	"github.com/robertguss/mycelium/internal/manifest"
	"github.com/robertguss/mycelium/internal/op"
	"github.com/robertguss/mycelium/internal/slug"
	"github.com/robertguss/mycelium/internal/teach"
	"github.com/robertguss/mycelium/internal/version"
)

const methodologyVersion = "2.0.0"

var allowedTiers = map[string]struct{}{
	"focused": {}, "standard": {}, "high-assurance": {},
}

// tierEmits are dirs (+ README stubs) beyond the spark skeleton.
var tierEmits = map[string][]string{
	"focused": {},
	"standard": {
		"decisions", "assumptions", "evidence", "questions", "risks",
	},
	"high-assurance": {
		"decisions", "assumptions", "evidence", "questions", "risks",
		"spikes", "findings", "recommendations", "requirements",
		"phases", "milestones",
	},
}

// Options for one new-idea scaffold.
type Options struct {
	Name    string
	Dir     string // exact destination; empty → ./<slug> under Cwd
	Offline bool
	Publish bool
	Tier    string // default focused
	Cwd     string
	Argv    []string // journal argv (without program name), e.g. new idea Name --offline
}

// Deps are injectable collaborators.
type Deps struct {
	Clock  clock.Clock
	Runner execrun.Runner
	Stdout io.Writer
	Stderr io.Writer
}

// Run scaffolds a spark instance. Exit 0 on success, 1 on teaching error.
func Run(opts Options, deps Deps) int {
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

	name := strings.TrimSpace(opts.Name)
	if name == "" {
		return teach.Write(deps.Stderr,
			"idea name is required",
			"command-flags",
			"framework/phases/PHASE-01-implementation-brief.md",
			`run mycelium new idea "<name>" [--offline]`,
		)
	}

	tier := opts.Tier
	if tier == "" {
		tier = "focused"
	}
	if _, ok := allowedTiers[tier]; !ok {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("unknown tier %q", tier),
			"tier",
			"program/contracts/conformance.md",
			"use --tier focused|standard|high-assurance",
		)
	}

	if opts.Offline && opts.Publish {
		return teach.Write(deps.Stderr,
			"--offline and --publish are contradictory",
			"command-flags",
			"framework/phases/PHASE-01-implementation-brief.md",
			"use --offline for hermetic local scaffold, or --publish without --offline",
		)
	}
	if opts.Publish {
		return teach.Write(deps.Stderr,
			`publish is not implemented in this slice`,
			"phase-01-slice-order",
			"framework/phases/PHASE-01-implementation-brief.md",
			"use --offline for local scaffold; publish lands in a later slice",
		)
	}

	ideaSlug, err := slug.Slugify(name)
	if err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("cannot slugify idea name: %v", err),
			"slugify",
			"framework/decisions/DEC-014-phase-01-slugify-latin-fold.md",
			"pass a name with at least one letter or digit",
		)
	}

	cwd := opts.Cwd
	if cwd == "" {
		cwd, err = os.Getwd()
		if err != nil {
			return teach.Write(deps.Stderr,
				fmt.Sprintf("cannot resolve cwd: %v", err),
				"command-flags",
				"framework/phases/PHASE-01-implementation-brief.md",
				"retry from a readable working directory",
			)
		}
	}

	displayPath := opts.Dir
	target := opts.Dir
	if target == "" {
		displayPath = "./" + ideaSlug
		target = filepath.Join(cwd, ideaSlug)
	} else if !filepath.IsAbs(target) {
		target = filepath.Join(cwd, target)
	}
	target = filepath.Clean(target)

	parent := filepath.Dir(target)
	if st, err := os.Stat(parent); err != nil || !st.IsDir() {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("parent directory does not exist: %s", parent),
			"scaffold-dir",
			"framework/phases/PHASE-01-implementation-brief.md",
			"create the parent directory, then retry",
		)
	}

	now := deps.Clock.Now()
	date := clock.Date(now)

	argv := opts.Argv
	if len(argv) == 0 {
		argv = []string{"new", "idea", name}
		if opts.Offline {
			argv = append(argv, "--offline")
		}
		if opts.Tier != "" && opts.Tier != "focused" {
			argv = append(argv, "--tier", opts.Tier)
		}
		if opts.Dir != "" {
			argv = append(argv, "--dir", opts.Dir)
		}
	}

	createdRoot := false
	st, err := os.Stat(target)
	switch {
	case err == nil:
		if !st.IsDir() {
			return teach.Write(deps.Stderr,
				fmt.Sprintf("refuse: target already exists: %s", displayPath),
				"scaffold-dir",
				"framework/phases/PHASE-01-implementation-brief.md",
				"choose a new name, pass --dir to a free path, or remove the existing path",
			)
		}
		existing, jerr := journal.Load(target)
		if jerr != nil && !errors.Is(jerr, journal.ErrNotExist) {
			return teach.Write(deps.Stderr,
				fmt.Sprintf("cannot read journal: %v", jerr),
				"operation-protocol",
				"program/contracts/operation-protocol.md",
				"mycelium check --abort-journal, then retry",
			)
		}
		if existing == nil {
			ok, perr := onlyMyceliumLeftovers(target)
			if perr != nil {
				return teach.Write(deps.Stderr,
					fmt.Sprintf("cannot inspect target: %v", perr),
					"scaffold-dir",
					"framework/phases/PHASE-01-implementation-brief.md",
					"fix filesystem permissions and retry",
				)
			}
			if !ok {
				return teach.Write(deps.Stderr,
					fmt.Sprintf("refuse: target already exists: %s", displayPath),
					"scaffold-dir",
					"framework/phases/PHASE-01-implementation-brief.md",
					"choose a new name, pass --dir to a free path, or remove the existing path",
				)
			}
			if err := os.RemoveAll(filepath.Join(target, ".mycelium")); err != nil {
				return teach.Write(deps.Stderr,
					fmt.Sprintf("cannot prune leftover .mycelium: %v", err),
					"scaffold-dir",
					"framework/phases/PHASE-01-implementation-brief.md",
					"remove the .mycelium directory manually, then retry",
				)
			}
			createdRoot = true
		}
	case os.IsNotExist(err):
		if err := os.Mkdir(target, 0o755); err != nil {
			return teach.Write(deps.Stderr,
				fmt.Sprintf("cannot create target directory: %v", err),
				"scaffold-dir",
				"framework/phases/PHASE-01-implementation-brief.md",
				"fix filesystem permissions and retry",
			)
		}
		createdRoot = true
	default:
		return teach.Write(deps.Stderr,
			fmt.Sprintf("cannot stat target: %v", err),
			"scaffold-dir",
			"framework/phases/PHASE-01-implementation-brief.md",
			"fix filesystem permissions and retry",
		)
	}
	sess, err := op.Begin(target, op.Intent{
		Op:      "scaffold",
		Title:   name,
		LogLine: logfmt.Line(date, "scaffold", "-", name),
		Argv:    argv,
	}, now)
	if err != nil {
		if createdRoot {
			_ = os.RemoveAll(target)
		}
		return teach.Write(deps.Stderr,
			fmt.Sprintf("operation begin failed: %v", err),
			"operation-protocol",
			"program/contracts/operation-protocol.md",
			"re-run the same command, or mycelium check --abort-journal",
		)
	}

	files, err := buildFiles(name, ideaSlug, tier, date, version.Version)
	if err != nil {
		rollbackOrClose(sess)
		if createdRoot {
			_ = os.RemoveAll(target)
		}
		return teach.Write(deps.Stderr,
			fmt.Sprintf("scaffold plan failed: %v", err),
			"scaffold",
			"framework/phases/PHASE-01-implementation-brief.md",
			"report this as a CLI bug; embed program/ may be incomplete",
		)
	}

	if err := sess.Stage(files); err != nil {
		rollbackOrClose(sess)
		if createdRoot {
			_ = os.RemoveAll(target)
		}
		return teach.Write(deps.Stderr,
			fmt.Sprintf("stage failed: %v", err),
			"operation-protocol",
			"program/contracts/operation-protocol.md",
			"re-run the same command after freeing disk space",
		)
	}

	if err := sess.Commit(); err != nil {
		_ = sess.Close()
		return teach.Write(deps.Stderr,
			fmt.Sprintf("commit failed: %v", err),
			"operation-protocol",
			"program/contracts/operation-protocol.md",
			"re-run the same command to complete, or mycelium check --abort-journal to roll back",
		)
	}
	_ = os.Remove(filepath.Join(target, ".mycelium"))

	res, err := deps.Runner.Run(context.Background(), "git",
		[]string{"--git-dir=.git", "--work-tree=.", "init", "-b", "main"},
		execrun.RunOpts{Dir: target},
	)
	if err != nil || res.ExitCode != 0 {
		msg := fmt.Sprintf("git init failed: %v", err)
		if err == nil {
			msg = fmt.Sprintf("git init failed: exit %d: %s", res.ExitCode, strings.TrimSpace(string(res.Stderr)))
		}
		return teach.Write(deps.Stderr,
			msg,
			"git-init",
			"program/contracts/operation-protocol.md",
			fmt.Sprintf("instance files exist at %s; run: git --git-dir=.git --work-tree=. init -b main", displayPath),
		)
	}
	if _, err := os.Stat(filepath.Join(target, ".git")); err != nil {
		return teach.Write(deps.Stderr,
			"git init reported success but instance .git is missing",
			"git-init",
			"program/contracts/operation-protocol.md",
			fmt.Sprintf("instance files exist at %s; unset GIT_DIR and run: git --git-dir=.git --work-tree=. init -b main", displayPath),
		)
	}

	fmt.Fprintf(deps.Stdout, "created %s\n", displayPath)
	fmt.Fprintf(deps.Stdout, "state: spark\n")
	fmt.Fprintf(deps.Stdout, "tier: %s\n", tier)
	fmt.Fprintf(deps.Stdout, "next: cd %s && mycelium new decision \"First thought\"\n", ideaSlug)
	fmt.Fprintf(deps.Stdout, "publish: mycelium publish\n")
	return 0
}

func rollbackOrClose(sess *op.Session) {
	if err := sess.Rollback(); errors.Is(err, op.ErrPartialCommit) {
		_ = sess.Close()
	}
}

// onlyMyceliumLeftovers reports whether root contains solely a .mycelium/ dir
// (failed prior attempt with no journal). Empty dirs and any other files refuse.
func onlyMyceliumLeftovers(root string) (bool, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return false, err
	}
	if len(entries) != 1 {
		return false, nil
	}
	e := entries[0]
	return e.Name() == ".mycelium" && e.IsDir(), nil
}

func buildFiles(ideaName, ideaSlug, tier, date, cliVersion string) ([]op.Staged, error) {
	var files []op.Staged

	readme, err := myceliumembed.Program.ReadFile("program/skeleton/README.md")
	if err != nil {
		return nil, err
	}
	readme = []byte(strings.Replace(string(readme), "# Idea name", "# "+ideaName, 1))
	files = append(files, op.Staged{RelTo: "README.md", Content: readme})

	for _, name := range []struct {
		embed string
		dest  string
	}{
		{"program/skeleton/CONTEXT.md", "CONTEXT.md"},
		{"program/skeleton/AGENTS.md", "AGENTS.md"},
		{"program/skeleton/gitignore", ".gitignore"},
		{"program/skills/mycelium-cli/SKILL.md", ".agents/skills/mycelium-cli/SKILL.md"},
	} {
		b, err := myceliumembed.Program.ReadFile(name.embed)
		if err != nil {
			return nil, err
		}
		files = append(files, op.Staged{RelTo: name.dest, Content: b})
	}

	prog, err := collectProgram()
	if err != nil {
		return nil, err
	}
	files = append(files, prog...)

	for _, dir := range tierEmits[tier] {
		title := strings.ToUpper(dir[:1]) + dir[1:]
		stub := fmt.Sprintf("# %s\n", title)
		files = append(files, op.Staged{
			RelTo:   filepath.ToSlash(filepath.Join(dir, "README.md")),
			Content: []byte(stub),
		})
	}

	logBody := "# Log\n\n" + logfmt.Line(date, "scaffold", "-", ideaName) + "\n"
	files = append(files, op.Staged{RelTo: "log.md", Content: []byte(logBody)})

	m := manifest.Manifest{
		SchemaVersion:         1,
		IdeaName:              ideaName,
		Slug:                  ideaSlug,
		State:                 "spark",
		Tier:                  tier,
		MethodologyVersion:    methodologyVersion,
		GeneratedByCLIVersion: cliVersion,
		CreatedDate:           date,
		UpdatedDate:           date,
		Revisit:               "",
		GithubRepo:            "",
	}
	mb, err := manifest.Encode(m)
	if err != nil {
		return nil, err
	}
	if _, err := manifest.Parse(mb); err != nil {
		return nil, fmt.Errorf("encoded manifest invalid: %w", err)
	}
	files = append(files, op.Staged{RelTo: "mycelium.toml", Content: mb})

	return files, nil
}

func collectProgram() ([]op.Staged, error) {
	var out []op.Staged
	err := fs.WalkDir(myceliumembed.Program, "program", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".go") {
			return nil
		}
		b, err := myceliumembed.Program.ReadFile(path)
		if err != nil {
			return err
		}
		out = append(out, op.Staged{RelTo: filepath.ToSlash(path), Content: b})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OfflineFromEnv is true when MYCELIUM_OFFLINE is a truthy value.
func OfflineFromEnv(getenv func(string) string) bool {
	if getenv == nil {
		getenv = os.Getenv
	}
	v := strings.TrimSpace(strings.ToLower(getenv("MYCELIUM_OFFLINE")))
	switch v {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}
