// Package cli dispatches mycelium commands.
package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/robertguss/mycelium/internal/check"
	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/execrun"
	"github.com/robertguss/mycelium/internal/generate"
	"github.com/robertguss/mycelium/internal/op"
	"github.com/robertguss/mycelium/internal/scaffold"
	"github.com/robertguss/mycelium/internal/teach"
	"github.com/robertguss/mycelium/internal/version"
)

const usage = `mycelium — convention-over-configuration thinking CLI

Usage:
  mycelium version
  mycelium new idea <name> [--dir PATH] [--offline] [--publish] [--tier focused|standard|high-assurance]
  mycelium new <type> "<Title>" [--dir PATH]
  mycelium check [--dir PATH] [--abort-journal]
  mycelium -h | --help

PHASE-01 commands (later slices): tier, publish
`

// Deps are injectable collaborators for hermetic tests.
type Deps struct {
	Clock     clock.Clock
	Runner    execrun.Runner
	LookupEnv func(string) string
	Getwd     func() (string, error)
}

// Main runs the CLI with argv (os.Args style: argv[0] is the program name).
// Exit 0 on success, exit 1 on failure.
func Main(argv []string) int {
	return Run(argv, os.Stdout, os.Stderr, Deps{})
}

// Run is Main with injectable stdout/stderr and deps.
func Run(argv []string, stdout, stderr io.Writer, deps Deps) int {
	if deps.LookupEnv == nil {
		deps.LookupEnv = os.Getenv
	}
	if deps.Getwd == nil {
		deps.Getwd = os.Getwd
	}
	if deps.Clock == nil {
		deps.Clock = clock.System()
	}
	if deps.Runner == nil {
		deps.Runner = execrun.Real{}
	}

	args := argv[1:]
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Fprint(stdout, usage)
		return 0
	}

	cmd := args[0]
	switch cmd {
	case "version":
		return cmdVersion(args[1:], stdout, stderr)
	case "help":
		fmt.Fprint(stdout, usage)
		return 0
	case "new":
		return cmdNew(args[1:], stdout, stderr, deps)
	case "check":
		return cmdCheck(args[1:], stdout, stderr, deps)
	case "tier", "publish":
		return teach.Write(stderr,
			fmt.Sprintf("%q is not implemented in this slice", cmd),
			"phase-01-slice-order",
			"framework/phases/PHASE-01-implementation-brief.md",
			"use mycelium version, mycelium new idea --offline, or mycelium check; wait for the slice that ships this command",
		)
	default:
		return teach.Write(stderr,
			fmt.Sprintf("unknown command %q", cmd),
			"command-surface",
			"framework/phases/PHASE-01-implementation-brief.md",
			"run mycelium --help to list available commands",
		)
	}
}

func cmdVersion(args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 {
		a := args[0]
		if a == "-h" || a == "--help" {
			fmt.Fprintln(stdout, "Usage: mycelium version")
			return 0
		}
		return teach.Write(stderr,
			fmt.Sprintf("version accepts no arguments (got %q)", strings.Join(args, " ")),
			"command-flags",
			"framework/phases/PHASE-01-implementation-brief.md",
			"run mycelium version with no extra arguments",
		)
	}
	fmt.Fprintln(stdout, version.Version)
	return 0
}

func cmdNew(args []string, stdout, stderr io.Writer, deps Deps) int {
	if len(args) == 0 {
		return teach.Write(stderr,
			"new requires a subcommand: idea or a registered type",
			"command-flags",
			"framework/phases/PHASE-01-implementation-brief.md",
			`try: mycelium new idea "My idea" --offline`,
		)
	}
	if args[0] == "idea" {
		return cmdNewIdea(args[1:], stdout, stderr, deps)
	}
	return cmdNewType(args, stdout, stderr, deps)
}

func cmdNewType(args []string, stdout, stderr io.Writer, deps Deps) int {
	opts, err := parseNewTypeFlags(args)
	if err == errHelp {
		fmt.Fprintln(stdout, `Usage: mycelium new <type> "<Title>" [--dir PATH]`)
		return 0
	}
	if err != nil {
		return teach.Write(stderr,
			err.Error(),
			"command-flags",
			"framework/phases/PHASE-01-implementation-brief.md",
			`usage: mycelium new <type> "<Title>" [--dir PATH]`,
		)
	}
	cwd, err := deps.Getwd()
	if err != nil {
		return teach.Write(stderr,
			fmt.Sprintf("cannot resolve cwd: %v", err),
			"command-flags",
			"framework/phases/PHASE-01-implementation-brief.md",
			"retry from a readable working directory",
		)
	}
	argv := append([]string{"new"}, args...)
	return generate.Run(generate.Options{
		TypeKey: opts.typeKey,
		Title:   opts.title,
		Dir:     opts.dir,
		Cwd:     cwd,
		Argv:    argv,
	}, generate.Deps{
		Clock:  deps.Clock,
		Stdout: stdout,
		Stderr: stderr,
	})
}

func cmdNewIdea(args []string, stdout, stderr io.Writer, deps Deps) int {
	opts, err := parseNewIdeaFlags(args)
	if err == errHelp {
		fmt.Fprintln(stdout, "Usage: mycelium new idea <name> [--dir PATH] [--offline] [--publish] [--tier focused|standard|high-assurance]")
		return 0
	}
	if err != nil {
		return teach.Write(stderr,
			err.Error(),
			"command-flags",
			"framework/phases/PHASE-01-implementation-brief.md",
			`usage: mycelium new idea <name> [--dir PATH] [--offline] [--publish] [--tier focused|standard|high-assurance]`,
		)
	}
	off := opts.offline || scaffold.OfflineFromEnv(deps.LookupEnv)
	cwd, err := deps.Getwd()
	if err != nil {
		return teach.Write(stderr,
			fmt.Sprintf("cannot resolve cwd: %v", err),
			"command-flags",
			"framework/phases/PHASE-01-implementation-brief.md",
			"retry from a readable working directory",
		)
	}

	argv := append([]string{"new", "idea"}, args...)
	return scaffold.Run(scaffold.Options{
		Name:    opts.name,
		Dir:     opts.dir,
		Offline: off,
		Publish: opts.publish,
		Tier:    opts.tier,
		Cwd:     cwd,
		Argv:    argv,
	}, scaffold.Deps{
		Clock:  deps.Clock,
		Runner: deps.Runner,
		Stdout: stdout,
		Stderr: stderr,
	})
}

func cmdCheck(args []string, stdout, stderr io.Writer, deps Deps) int {
	opts, err := parseCheckFlags(args)
	if err == errHelp {
		fmt.Fprintln(stdout, "Usage: mycelium check [--dir PATH] [--abort-journal]")
		return 0
	}
	if err != nil {
		return teach.Write(stderr,
			err.Error(),
			"command-flags",
			"framework/phases/PHASE-01-implementation-brief.md",
			"usage: mycelium check [--dir PATH] [--abort-journal]",
		)
	}
	cwd, err := deps.Getwd()
	if err != nil {
		return teach.Write(stderr,
			fmt.Sprintf("cannot resolve cwd: %v", err),
			"command-flags",
			"framework/phases/PHASE-01-implementation-brief.md",
			"retry from a readable working directory",
		)
	}
	start := cwd
	if opts.dir != "" {
		if filepath.IsAbs(opts.dir) {
			start = opts.dir
		} else {
			start = filepath.Join(cwd, opts.dir)
		}
	}
	root, err := check.FindRoot(start)
	if err != nil {
		return teach.Write(stderr,
			"not a mycelium instance (no mycelium.toml found)",
			"instance-root",
			"program/contracts/manifest.md",
			"run from an instance directory or pass --dir PATH",
		)
	}

	if opts.abort {
		if err := check.AbortJournal(root, stdout); err != nil {
			if errors.Is(err, op.ErrNothingToAbort) {
				return teach.Write(stderr,
					"nothing to abort (no journal, stale lock, or orphan stage)",
					"operation-protocol",
					"program/contracts/operation-protocol.md",
					"no action needed",
				)
			}
			if errors.Is(err, op.ErrLocked) {
				return teach.Write(stderr,
					fmt.Sprintf("lock held by another process: %v", err),
					"operation-protocol",
					"program/contracts/operation-protocol.md",
					"wait for the other process to finish; do not abort under a live lock",
				)
			}
			return teach.Write(stderr,
				fmt.Sprintf("abort failed: %v", err),
				"operation-protocol",
				"program/contracts/operation-protocol.md",
				"inspect .mycelium/ and retry",
			)
		}
		return 0
	}

	r := check.Run(root)
	if r.LiveLockNotice != "" {
		fmt.Fprintln(stdout, r.LiveLockNotice)
	}
	if !r.OK {
		return teach.WriteFindings(stderr, r.Findings)
	}
	check.WriteOK(stdout, r)
	return 0
}

var errHelp = fmt.Errorf("help")

type newIdeaFlags struct {
	name    string
	dir     string
	offline bool
	publish bool
	tier    string
}

type newTypeFlags struct {
	typeKey string
	title   string
	dir     string
}

type checkFlags struct {
	dir   string
	abort bool
}

func parseNewTypeFlags(args []string) (newTypeFlags, error) {
	var out newTypeFlags
	var positionals []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-h" || a == "--help":
			return out, errHelp
		case a == "--dir":
			if i+1 >= len(args) {
				return out, fmt.Errorf("--dir requires a path")
			}
			i++
			out.dir = args[i]
		case strings.HasPrefix(a, "--dir="):
			out.dir = strings.TrimPrefix(a, "--dir=")
		case strings.HasPrefix(a, "-"):
			return out, fmt.Errorf("unknown flag %q", a)
		default:
			positionals = append(positionals, a)
		}
	}
	if len(positionals) == 0 {
		return out, fmt.Errorf("new requires a type key and title")
	}
	out.typeKey = positionals[0]
	out.title = strings.TrimSpace(strings.Join(positionals[1:], " "))
	return out, nil
}

func parseCheckFlags(args []string) (checkFlags, error) {
	var out checkFlags
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-h" || a == "--help":
			return out, errHelp
		case a == "--abort-journal":
			out.abort = true
		case a == "--dir":
			if i+1 >= len(args) {
				return out, fmt.Errorf("--dir requires a path")
			}
			i++
			out.dir = args[i]
		case strings.HasPrefix(a, "--dir="):
			out.dir = strings.TrimPrefix(a, "--dir=")
		case strings.HasPrefix(a, "-"):
			return out, fmt.Errorf("unknown flag %q", a)
		default:
			return out, fmt.Errorf("unexpected argument %q", a)
		}
	}
	return out, nil
}

func parseNewIdeaFlags(args []string) (newIdeaFlags, error) {
	out := newIdeaFlags{tier: "focused"}
	var nameParts []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-h" || a == "--help":
			return out, errHelp
		case a == "--offline":
			out.offline = true
		case a == "--publish":
			out.publish = true
		case a == "--dir":
			if i+1 >= len(args) {
				return out, fmt.Errorf("--dir requires a path")
			}
			i++
			out.dir = args[i]
		case strings.HasPrefix(a, "--dir="):
			out.dir = strings.TrimPrefix(a, "--dir=")
		case a == "--tier":
			if i+1 >= len(args) {
				return out, fmt.Errorf("--tier requires a value")
			}
			i++
			out.tier = args[i]
		case strings.HasPrefix(a, "--tier="):
			out.tier = strings.TrimPrefix(a, "--tier=")
		case strings.HasPrefix(a, "-"):
			return out, fmt.Errorf("unknown flag %q", a)
		default:
			nameParts = append(nameParts, a)
		}
	}
	out.name = strings.TrimSpace(strings.Join(nameParts, " "))
	return out, nil
}
