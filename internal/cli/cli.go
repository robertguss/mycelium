// Package cli dispatches mycelium commands.
package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/robertguss/mycelium/internal/version"
)

// phase01Names are commissioned PHASE-01 commands not implemented in Slice 1.
var phase01Names = map[string]struct{}{
	"new":     {},
	"check":   {},
	"tier":    {},
	"publish": {},
}

const usage = `mycelium — convention-over-configuration thinking CLI

Usage:
  mycelium version
  mycelium -h | --help

PHASE-01 commands (later slices): new, check, tier, publish
`

// Main runs the CLI with argv (os.Args style: argv[0] is the program name).
// Exit 0 on success, 1 on failure.
func Main(argv []string) int {
	return RunForTest(argv, os.Stdout, os.Stderr)
}

// RunForTest is Main with injectable stdout/stderr for hermetic tests.
func RunForTest(argv []string, stdout, stderr io.Writer) int {
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
	}

	if _, ok := phase01Names[cmd]; ok {
		return teach(stderr,
			fmt.Sprintf("%q is not implemented in this slice", cmd),
			"phase-01-slice-order",
			"framework/phases/PHASE-01-implementation-brief.md",
			"use mycelium version; wait for the slice that ships this command",
		)
	}

	return teach(stderr,
		fmt.Sprintf("unknown command %q", cmd),
		"command-surface",
		"framework/phases/PHASE-01-implementation-brief.md",
		"run mycelium --help to list available commands",
	)
}

func cmdVersion(args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 {
		a := args[0]
		if a == "-h" || a == "--help" {
			fmt.Fprintln(stdout, "Usage: mycelium version")
			return 0
		}
		return teach(stderr,
			fmt.Sprintf("version accepts no arguments (got %q)", strings.Join(args, " ")),
			"command-flags",
			"framework/phases/PHASE-01-implementation-brief.md",
			"run mycelium version with no extra arguments",
		)
	}
	fmt.Fprintln(stdout, version.Version)
	return 0
}

func teach(stderr io.Writer, what, convention, contract, fix string) int {
	fmt.Fprintf(stderr, "mycelium: %s\n", what)
	fmt.Fprintf(stderr, "convention: %s\n", convention)
	fmt.Fprintf(stderr, "contract: %s\n", contract)
	fmt.Fprintf(stderr, "fix: %s\n", fix)
	return 1
}
