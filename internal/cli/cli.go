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
	"github.com/robertguss/mycelium/internal/indexcmd"
	"github.com/robertguss/mycelium/internal/op"
	"github.com/robertguss/mycelium/internal/publish"
	"github.com/robertguss/mycelium/internal/scaffold"
	"github.com/robertguss/mycelium/internal/statecmd"
	"github.com/robertguss/mycelium/internal/statuscmd"
	"github.com/robertguss/mycelium/internal/teach"
	"github.com/robertguss/mycelium/internal/tiercmd"
	"github.com/robertguss/mycelium/internal/version"
)

const usage = `mycelium — convention-over-configuration thinking CLI

Usage:
  mycelium version
  mycelium new idea <name> [--dir PATH] [--offline] [--publish] [--tier focused|standard|high-assurance]
  mycelium new <type> "<Title>" [--dir PATH]
  mycelium check [--dir PATH] [--abort-journal]
  mycelium tier <tier> [--dir PATH]
  mycelium publish [--dir PATH]
  mycelium index [--dir PATH]
  mycelium state <exploring|simmering|clarified|archived> [--dir PATH] [--revisit VALUE]
  mycelium wake [--dir PATH]
  mycelium status [--dir PATH] [--all] [--root PATH] [--archived] [--offline]
  mycelium -h | --help
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
	case "tier":
		return cmdTier(args[1:], stdout, stderr, deps)
	case "publish":
		return cmdPublish(args[1:], stdout, stderr, deps)
	case "index":
		return cmdIndex(args[1:], stdout, stderr, deps)
	case "state":
		return cmdState(args[1:], stdout, stderr, deps)
	case "wake":
		return cmdWake(args[1:], stdout, stderr, deps)
	case "status":
		return cmdStatus(args[1:], stdout, stderr, deps)
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

func cmdTier(args []string, stdout, stderr io.Writer, deps Deps) int {
	opts, err := parseTierFlags(args)
	if err == errHelp {
		fmt.Fprintln(stdout, "Usage: mycelium tier <tier> [--dir PATH]")
		return 0
	}
	if err != nil {
		return teach.Write(stderr,
			err.Error(),
			"command-flags",
			"framework/phases/PHASE-01-implementation-brief.md",
			"usage: mycelium tier <tier> [--dir PATH] where <tier> is focused, standard, or high-assurance",
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
	argv := append([]string{"tier"}, args...)
	return tiercmd.Run(tiercmd.Options{
		Tier: opts.tier,
		Dir:  opts.dir,
		Cwd:  cwd,
		Argv: argv,
	}, tiercmd.Deps{
		Clock:  deps.Clock,
		Stdout: stdout,
		Stderr: stderr,
	})
}

func cmdPublish(args []string, stdout, stderr io.Writer, deps Deps) int {
	opts, err := parsePublishFlags(args)
	if err == errHelp {
		fmt.Fprintln(stdout, "Usage: mycelium publish [--dir PATH]")
		return 0
	}
	if err != nil {
		return teach.Write(stderr,
			err.Error(),
			"command-flags",
			"framework/phases/PHASE-01-implementation-brief.md",
			"usage: mycelium publish [--dir PATH]",
		)
	}
	off := scaffold.OfflineFromEnv(deps.LookupEnv)
	cwd, err := deps.Getwd()
	if err != nil {
		return teach.Write(stderr,
			fmt.Sprintf("cannot resolve cwd: %v", err),
			"command-flags",
			"framework/phases/PHASE-01-implementation-brief.md",
			"retry from a readable working directory",
		)
	}
	mode := publish.ModeRequired
	if off {
		mode = publish.ModeOffline
	}
	argv := append([]string{"publish"}, args...)
	out := publish.Run(publish.Options{
		Dir:        opts.dir,
		Cwd:        cwd,
		Mode:       mode,
		Argv:       argv,
		RequireGit: true,
	}, publish.Deps{
		Clock:  deps.Clock,
		Runner: deps.Runner,
		Stdout: stdout,
		Stderr: stderr,
	})
	return out.Code
}

func cmdIndex(args []string, stdout, stderr io.Writer, deps Deps) int {
	opts, err := parseIndexFlags(args)
	if err == errHelp {
		fmt.Fprintln(stdout, `Usage: mycelium index [--dir PATH]

Examples:
  mycelium index
  mycelium index --dir ./my-idea`)
		return 0
	}
	if err != nil {
		return teach.Write(stderr,
			err.Error(),
			"command-flags",
			"framework/phases/PHASE-01-implementation-brief.md",
			"usage: mycelium index [--dir PATH]",
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
	argv := append([]string{"index"}, args...)
	return indexcmd.Run(indexcmd.Options{
		Dir:  opts.dir,
		Cwd:  cwd,
		Argv: argv,
	}, indexcmd.Deps{
		Clock:  deps.Clock,
		Stdout: stdout,
		Stderr: stderr,
	})
}

func cmdState(args []string, stdout, stderr io.Writer, deps Deps) int {
	opts, err := parseStateFlags(args)
	if err == errHelp {
		fmt.Fprintln(stdout, `Usage: mycelium state <exploring|simmering|clarified|archived> [--dir PATH] [--revisit VALUE]

Examples:
  mycelium state exploring
  mycelium state simmering --revisit 2026-08-08
  mycelium state simmering --revisit event:after-iphone-launch
  mycelium state clarified --dir ./my-idea
  mycelium state archived`)
		return 0
	}
	if err != nil {
		return teach.Write(stderr,
			err.Error(),
			"command-flags",
			"program/contracts/lifecycle.md",
			"usage: mycelium state <exploring|simmering|clarified|archived> [--dir PATH] [--revisit VALUE]",
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
	argv := append([]string{"state"}, args...)
	return statecmd.RunState(statecmd.Options{
		Target:     opts.target,
		Dir:        opts.dir,
		Cwd:        cwd,
		Revisit:    opts.revisit,
		HasRevisit: opts.hasRevisit,
		Argv:       argv,
	}, statecmd.Deps{
		Clock:  deps.Clock,
		Stdout: stdout,
		Stderr: stderr,
	})
}

func cmdWake(args []string, stdout, stderr io.Writer, deps Deps) int {
	opts, err := parseWakeFlags(args)
	if err == errHelp {
		fmt.Fprintln(stdout, `Usage: mycelium wake [--dir PATH]

Examples:
  mycelium wake
  mycelium wake --dir ./my-idea`)
		return 0
	}
	if err != nil {
		return teach.Write(stderr,
			err.Error(),
			"command-flags",
			"program/contracts/wake.md",
			"usage: mycelium wake [--dir PATH]",
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
	argv := append([]string{"wake"}, args...)
	return statecmd.RunWake(statecmd.Options{
		Dir:  opts.dir,
		Cwd:  cwd,
		Argv: argv,
	}, statecmd.Deps{
		Clock:  deps.Clock,
		Stdout: stdout,
		Stderr: stderr,
	})
}

func cmdStatus(args []string, stdout, stderr io.Writer, deps Deps) int {
	opts, err := parseStatusFlags(args)
	if err == errHelp {
		fmt.Fprintln(stdout, `Usage: mycelium status [--dir PATH] [--all] [--root PATH] [--archived] [--offline]

Examples:
  mycelium status
  mycelium status --dir ./my-idea
  mycelium status --all --offline --root ~/ideas`)
		return 0
	}
	if err != nil {
		return teach.Write(stderr,
			err.Error(),
			"command-flags",
			"program/contracts/status.md",
			"usage: mycelium status [--dir PATH] [--all] [--root PATH] [--archived] [--offline]",
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
	return statuscmd.Run(statuscmd.Options{
		Dir:      opts.dir,
		Cwd:      cwd,
		All:      opts.all,
		Root:     opts.root,
		Offline:  opts.offline,
		Archived: opts.archived,
	}, statuscmd.Deps{
		Clock:     deps.Clock,
		Stdout:    stdout,
		Stderr:    stderr,
		Runner:    deps.Runner,
		LookupEnv: deps.LookupEnv,
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

type publishFlags struct {
	dir string
}

type indexFlags struct {
	dir string
}

type stateFlags struct {
	target     string
	dir        string
	revisit    string
	hasRevisit bool
}

type wakeFlags struct {
	dir string
}

type statusFlags struct {
	dir      string
	all      bool
	root     string
	hasRoot  bool
	offline  bool
	archived bool
}

type tierFlags struct {
	tier string
	dir  string
}

func parseTierFlags(args []string) (tierFlags, error) {
	var out tierFlags
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
		return out, fmt.Errorf("tier requires a tier name")
	}
	if len(positionals) > 1 {
		return out, fmt.Errorf("unexpected argument %q", positionals[1])
	}
	out.tier = positionals[0]
	return out, nil
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

func parsePublishFlags(args []string) (publishFlags, error) {
	var out publishFlags
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
			return out, fmt.Errorf("unexpected argument %q", a)
		}
	}
	return out, nil
}

func parseIndexFlags(args []string) (indexFlags, error) {
	var out indexFlags
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
			return out, fmt.Errorf("unexpected argument %q", a)
		}
	}
	return out, nil
}

func parseStateFlags(args []string) (stateFlags, error) {
	var out stateFlags
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
		case a == "--revisit":
			if i+1 >= len(args) {
				return out, fmt.Errorf("--revisit requires a value")
			}
			i++
			out.revisit = args[i]
			out.hasRevisit = true
		case strings.HasPrefix(a, "--revisit="):
			out.revisit = strings.TrimPrefix(a, "--revisit=")
			out.hasRevisit = true
		case strings.HasPrefix(a, "-"):
			return out, fmt.Errorf("unknown flag %q", a)
		default:
			positionals = append(positionals, a)
		}
	}
	if len(positionals) == 0 {
		return out, fmt.Errorf("state requires a target")
	}
	if len(positionals) > 1 {
		return out, fmt.Errorf("unexpected argument %q", positionals[1])
	}
	out.target = positionals[0]
	return out, nil
}

func parseWakeFlags(args []string) (wakeFlags, error) {
	var out wakeFlags
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
			return out, fmt.Errorf("unexpected argument %q", a)
		}
	}
	return out, nil
}

func parseStatusFlags(args []string) (statusFlags, error) {
	var out statusFlags
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-h" || a == "--help":
			return out, errHelp
		case a == "--all":
			out.all = true
		case a == "--offline":
			out.offline = true
		case a == "--archived":
			out.archived = true
		case a == "--dir":
			if i+1 >= len(args) {
				return out, fmt.Errorf("--dir requires a path")
			}
			i++
			out.dir = args[i]
		case strings.HasPrefix(a, "--dir="):
			out.dir = strings.TrimPrefix(a, "--dir=")
		case a == "--root":
			if i+1 >= len(args) {
				return out, fmt.Errorf("--root requires a path")
			}
			i++
			out.root = args[i]
			out.hasRoot = true
		case strings.HasPrefix(a, "--root="):
			out.root = strings.TrimPrefix(a, "--root=")
			out.hasRoot = true
		case strings.HasPrefix(a, "-"):
			return out, fmt.Errorf("unknown flag %q", a)
		default:
			return out, fmt.Errorf("unexpected argument %q", a)
		}
	}
	if out.hasRoot && !out.all {
		return out, fmt.Errorf("--root requires --all")
	}
	if out.archived && !out.all {
		return out, fmt.Errorf("status --archived requires --all")
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
