// Package statuscmd implements mycelium status (single instance in Slice 4).
package statuscmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/robertguss/mycelium/internal/check"
	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/manifest"
	"github.com/robertguss/mycelium/internal/revisit"
	"github.com/robertguss/mycelium/internal/teach"
)

// Options for one status run.
type Options struct {
	Dir string
	Cwd string
}

// Deps are injectable collaborators.
type Deps struct {
	Clock  clock.Clock
	Stdout io.Writer
	Stderr io.Writer
}

// Run prints single-instance status. Read-only. Exit 0 on success, 1 on teaching error.
func Run(opts Options, deps Deps) int {
	if deps.Clock == nil {
		deps.Clock = clock.System()
	}
	if deps.Stdout == nil {
		deps.Stdout = io.Discard
	}
	if deps.Stderr == nil {
		deps.Stderr = io.Discard
	}

	cwd := opts.Cwd
	var err error
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
	start := cwd
	if opts.Dir != "" {
		if filepath.IsAbs(opts.Dir) {
			start = opts.Dir
		} else {
			start = filepath.Join(cwd, opts.Dir)
		}
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

	manBytes, err := os.ReadFile(filepath.Join(root, "mycelium.toml"))
	if err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("cannot read mycelium.toml: %v", err),
			"manifest",
			"program/contracts/manifest.md",
			"restore mycelium.toml at the instance root",
		)
	}
	m, err := manifest.Parse(manBytes)
	if err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("invalid mycelium.toml: %v", err),
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
