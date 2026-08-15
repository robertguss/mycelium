// Package statuscmd implements mycelium status and status --all (local scan).
package statuscmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/robertguss/mycelium/internal/check"
	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/execrun"
	"github.com/robertguss/mycelium/internal/manifest"
	"github.com/robertguss/mycelium/internal/revisit"
	"github.com/robertguss/mycelium/internal/scaffold"
	"github.com/robertguss/mycelium/internal/teach"
)

// Options for one status run.
type Options struct {
	Dir      string
	Cwd      string
	All      bool
	Root     string
	Offline  bool
	Archived bool
}

// Deps are injectable collaborators.
type Deps struct {
	Clock       clock.Clock
	Stdout      io.Writer
	Stderr      io.Writer
	Runner      execrun.Runner
	LookupEnv   func(string) string
	UserHomeDir func() (string, error)
}

// Idea is one portfolio row (local-only in Slice 5).
type Idea struct {
	Slug    string
	State   string
	Tier    string
	Revisit string
	Flag    string
	Github  string
}

// Run prints status. Read-only. Exit 0 on success, 1 on teaching error.
func Run(opts Options, deps Deps) int {
	deps = normalizeDeps(deps)
	if opts.All {
		return runAll(opts, deps)
	}
	return runSingle(opts, deps)
}

func normalizeDeps(deps Deps) Deps {
	if deps.Clock == nil {
		deps.Clock = clock.System()
	}
	if deps.Stdout == nil {
		deps.Stdout = io.Discard
	}
	if deps.Stderr == nil {
		deps.Stderr = io.Discard
	}
	if deps.Runner == nil {
		deps.Runner = execrun.Real{}
	}
	if deps.LookupEnv == nil {
		deps.LookupEnv = os.Getenv
	}
	if deps.UserHomeDir == nil {
		deps.UserHomeDir = os.UserHomeDir
	}
	return deps
}

func runSingle(opts Options, deps Deps) int {
	cwd, err := resolveCwd(opts.Cwd)
	if err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("cannot resolve cwd: %v", err),
			"command-flags",
			"framework/phases/PHASE-01-implementation-brief.md",
			"retry from a readable working directory",
		)
	}
	start := cwd
	if opts.Dir != "" {
		start = absAgainst(opts.Dir, cwd)
	}
	root, err := check.FindRoot(start)
	if err != nil {
		return teach.Write(deps.Stderr,
			"not a mycelium instance (no mycelium.toml found)",
			"instance-root",
			"program/contracts/manifest.md",
			"run from an instance directory or pass --dir PATH",
		)
	}

	m, err := loadManifest(root)
	if err != nil {
		return teach.Write(deps.Stderr,
			err.Error(),
			"manifest",
			"program/contracts/manifest.md",
			"fix mycelium.toml required fields and retry",
		)
	}

	github := m.GithubRepo
	if github == "" {
		github = "unpublished"
	}

	fmt.Fprintf(deps.Stdout, "slug: %s\n", m.Slug)
	fmt.Fprintf(deps.Stdout, "state: %s\n", m.State)
	fmt.Fprintf(deps.Stdout, "tier: %s\n", m.Tier)
	fmt.Fprintf(deps.Stdout, "revisit: %s\n", m.Revisit)
	fmt.Fprintf(deps.Stdout, "due: %s\n", dueToken(m.Revisit, deps.Clock.Now()))
	fmt.Fprintf(deps.Stdout, "github: %s\n", github)
	return 0
}

func runAll(opts Options, deps Deps) int {
	cwd, err := resolveCwd(opts.Cwd)
	if err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("cannot resolve cwd: %v", err),
			"command-flags",
			"framework/phases/PHASE-01-implementation-brief.md",
			"retry from a readable working directory",
		)
	}

	offline := opts.Offline || scaffold.OfflineFromEnv(deps.LookupEnv)
	reason := partialReason(offline, deps.Runner)
	ideas := scanLocal(opts, deps, cwd)
	SortIdeas(ideas, deps.Clock.Now())

	fmt.Fprintf(deps.Stdout, "partial: local-only (%s)\n", reason)
	for _, idea := range ideas {
		fmt.Fprintf(deps.Stdout, "%s\t%s\t%s\t%s\t%s\t%s\n",
			idea.Slug, idea.State, idea.Tier, idea.Revisit, idea.Flag, idea.Github)
	}
	overdue := countOverdue(ideas, deps.Clock.Now())
	fmt.Fprintf(deps.Stdout, "%d ideas (%d overdue, partial)\n", len(ideas), overdue)
	return 0
}

func partialReason(offline bool, runner execrun.Runner) string {
	if offline {
		return "offline"
	}
	if _, err := runner.LookPath("gh"); err != nil {
		return "gh missing"
	}
	return "offline"
}

func scanLocal(opts Options, deps Deps, cwd string) []Idea {
	bySlug := map[string]Idea{}

	rootPath := resolveIdeasRoot(opts, deps)
	entries, err := os.ReadDir(rootPath)
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			child := filepath.Join(rootPath, e.Name())
			idea, ok, teachErr := tryLoadIdea(child)
			if teachErr != "" {
				writeScanTeach(deps.Stderr, teachErr)
				continue
			}
			if !ok {
				continue
			}
			if idea.State == "archived" && !opts.Archived {
				continue
			}
			bySlug[idea.Slug] = idea
		}
	}

	anchor := cwd
	if opts.Dir != "" {
		anchor = absAgainst(opts.Dir, cwd)
	}
	if inst, err := check.FindRoot(anchor); err == nil {
		idea, ok, teachErr := tryLoadIdea(inst)
		if teachErr != "" {
			writeScanTeach(deps.Stderr, teachErr)
		} else if ok {
			if idea.State != "archived" || opts.Archived {
				bySlug[idea.Slug] = idea
			}
		}
	}

	out := make([]Idea, 0, len(bySlug))
	for _, idea := range bySlug {
		out = append(out, idea)
	}
	return out
}

func resolveIdeasRoot(opts Options, deps Deps) string {
	home, _ := deps.UserHomeDir()
	if opts.Root != "" {
		return expandTilde(opts.Root, home)
	}
	if v := strings.TrimSpace(deps.LookupEnv("MYCELIUM_IDEAS_ROOT")); v != "" {
		return expandTilde(v, home)
	}
	if home == "" {
		return "ideas"
	}
	return filepath.Join(home, "ideas")
}

func expandTilde(path, home string) string {
	if path == "~" {
		if home == "" {
			return path
		}
		return home
	}
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, `~\`) {
		if home == "" {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

func tryLoadIdea(dir string) (Idea, bool, string) {
	tom := filepath.Join(dir, "mycelium.toml")
	if _, err := os.Stat(tom); err != nil {
		return Idea{}, false, ""
	}
	m, err := loadManifest(dir)
	if err != nil {
		return Idea{}, false, fmt.Sprintf("cannot read instance %s: %v", dir, err)
	}
	return Idea{
		Slug:    m.Slug,
		State:   m.State,
		Tier:    m.Tier,
		Revisit: m.Revisit,
		Flag:    "unpublished",
		Github:  "unpublished",
	}, true, ""
}

func writeScanTeach(stderr io.Writer, what string) {
	fmt.Fprintf(stderr, "mycelium: %s\n", what)
	fmt.Fprintf(stderr, "convention: status-local-scan\n")
	fmt.Fprintf(stderr, "contract: program/contracts/status.md\n")
	fmt.Fprintf(stderr, "fix: fix or remove the unreadable instance and retry\n")
}

func loadManifest(root string) (manifest.Manifest, error) {
	manBytes, err := os.ReadFile(filepath.Join(root, "mycelium.toml"))
	if err != nil {
		return manifest.Manifest{}, fmt.Errorf("cannot read mycelium.toml: %w", err)
	}
	m, err := manifest.Parse(manBytes)
	if err != nil {
		return manifest.Manifest{}, fmt.Errorf("invalid mycelium.toml: %w", err)
	}
	return m, nil
}

func resolveCwd(cwd string) (string, error) {
	if cwd != "" {
		return cwd, nil
	}
	return os.Getwd()
}

func absAgainst(path, cwd string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(cwd, path)
}

func dueToken(revisitRaw string, now time.Time) string {
	if revisitRaw == "" {
		return "no"
	}
	kind, date, _, err := revisit.Parse(revisitRaw)
	if err != nil {
		return "no"
	}
	if kind == revisit.Event {
		return "event"
	}
	if revisit.Due(kind, date, now) {
		return "yes"
	}
	return "no"
}

// SortIdeas orders portfolio rows per §9.6 (buckets, then slug ASC).
func SortIdeas(ideas []Idea, now time.Time) {
	sort.SliceStable(ideas, func(i, j int) bool {
		bi := sortBucket(ideas[i], now)
		bj := sortBucket(ideas[j], now)
		if bi != bj {
			return bi < bj
		}
		return ideas[i].Slug < ideas[j].Slug
	})
}

func sortBucket(idea Idea, now time.Time) int {
	if idea.State == "archived" {
		return 8
	}
	if idea.State == "simmering" {
		kind, date, _, err := revisit.Parse(idea.Revisit)
		if err == nil && kind == revisit.Date {
			if revisit.Overdue(kind, date, now) {
				return 1
			}
			if utcDate(now).Equal(utcDate(date)) {
				return 2
			}
			return 7
		}
		if err == nil && kind == revisit.Event {
			return 6
		}
		return 7
	}
	switch idea.State {
	case "exploring":
		return 3
	case "spark":
		return 4
	case "clarified":
		return 5
	default:
		return 7
	}
}

func countOverdue(ideas []Idea, now time.Time) int {
	n := 0
	for _, idea := range ideas {
		if idea.State != "simmering" {
			continue
		}
		kind, date, _, err := revisit.Parse(idea.Revisit)
		if err != nil || kind != revisit.Date {
			continue
		}
		if revisit.Overdue(kind, date, now) {
			n++
		}
	}
	return n
}

func utcDate(t time.Time) time.Time {
	y, m, d := t.UTC().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
