// Package statecmd implements mycelium state and mycelium wake.
package statecmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/robertguss/mycelium/internal/check"
	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/indexmd"
	"github.com/robertguss/mycelium/internal/lifecycle"
	"github.com/robertguss/mycelium/internal/logfmt"
	"github.com/robertguss/mycelium/internal/manifest"
	"github.com/robertguss/mycelium/internal/op"
	"github.com/robertguss/mycelium/internal/revisit"
	"github.com/robertguss/mycelium/internal/teach"
	"github.com/robertguss/mycelium/internal/wakebrief"
)

// Options for state or wake.
type Options struct {
	Target     string // exploring|simmering|clarified|archived; ignored when Wake
	Wake       bool
	Dir        string
	Cwd        string
	Revisit    string
	HasRevisit bool
	Argv       []string
}

// Deps are injectable collaborators.
type Deps struct {
	Clock  clock.Clock
	Stdout io.Writer
	Stderr io.Writer
}

// RunState executes mycelium state <target>.
func RunState(opts Options, deps Deps) int {
	opts.Wake = false
	return run(opts, deps)
}

// RunWake executes mycelium wake (simmering → exploring + brief).
func RunWake(opts Options, deps Deps) int {
	opts.Wake = true
	opts.Target = "exploring"
	opts.HasRevisit = false
	opts.Revisit = ""
	return run(opts, deps)
}

func run(opts Options, deps Deps) int {
	if deps.Clock == nil {
		deps.Clock = clock.System()
	}
	if deps.Stdout == nil {
		deps.Stdout = io.Discard
	}
	if deps.Stderr == nil {
		deps.Stderr = io.Discard
	}

	root, code := resolveRoot(opts, deps)
	if code != 0 {
		return code
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

	logBytes, err := os.ReadFile(filepath.Join(root, "log.md"))
	if err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("cannot read log.md: %v", err),
			"log",
			"program/contracts/conformance.md",
			"restore log.md at the instance root",
		)
	}

	if opts.Wake {
		return runWakeTransition(root, m, logBytes, opts, deps)
	}
	return runStateTransition(root, m, logBytes, opts, deps)
}

func runWakeTransition(root string, m manifest.Manifest, logBytes []byte, opts Options, deps Deps) int {
	if m.State != "simmering" {
		return teach.Write(deps.Stderr,
			"wake is only legal from simmering",
			"lifecycle",
			"program/contracts/lifecycle.md",
			"mycelium state simmering --revisit YYYY-MM-DD   # if you meant to park first",
		)
	}
	if _, _, _, err := revisit.Parse(m.Revisit); err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("revisit %q is not a date or event:<kebab>", m.Revisit),
			"revisit",
			"program/contracts/revisit.md",
			"use YYYY-MM-DD (UTC) or event:after-iphone-launch",
		)
	}
	return commitWake(root, m, logBytes, opts, deps)
}

func runStateTransition(root string, m manifest.Manifest, logBytes []byte, opts Options, deps Deps) int {
	target := strings.TrimSpace(opts.Target)
	if target == "" {
		return teach.Write(deps.Stderr,
			"state requires a target",
			"command-flags",
			"program/contracts/lifecycle.md",
			"usage: mycelium state <exploring|simmering|clarified|archived> [--dir PATH] [--revisit VALUE]",
		)
	}
	if target == "handed-off" {
		return teach.Write(deps.Stderr,
			"state=handed-off requires a PHASE-06 handoff packet",
			"lifecycle",
			"program/contracts/lifecycle.md",
			"stay in clarified, or mycelium state archived; packet command is not shipped",
		)
	}
	if !isAllowedTarget(target) {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("unknown state target %q", target),
			"lifecycle",
			"program/contracts/lifecycle.md",
			"allowed targets: exploring, simmering, clarified, archived (handed-off is PHASE-06)",
		)
	}
	if opts.HasRevisit && lifecycle.RevisitForbidden(target) {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("--revisit is only legal when target is simmering (got %s)", target),
			"revisit",
			"program/contracts/revisit.md",
			"omit --revisit, or use mycelium state simmering --revisit YYYY-MM-DD",
		)
	}
	if lifecycle.RevisitRequired(target) {
		if !opts.HasRevisit || strings.TrimSpace(opts.Revisit) == "" {
			return teach.Write(deps.Stderr,
				"state simmering requires --revisit",
				"revisit",
				"program/contracts/revisit.md",
				"mycelium state simmering --revisit YYYY-MM-DD",
			)
		}
		if _, _, _, err := revisit.Parse(opts.Revisit); err != nil {
			return teach.Write(deps.Stderr,
				fmt.Sprintf("revisit %q is not a date or event:<kebab>", opts.Revisit),
				"revisit",
				"program/contracts/revisit.md",
				"use YYYY-MM-DD (UTC) or event:after-iphone-launch",
			)
		}
	}

	if m.State == "archived" {
		return teach.Write(deps.Stderr,
			"archived is terminal",
			"lifecycle",
			"program/contracts/lifecycle.md",
			"no further state transitions are legal from archived",
		)
	}
	if m.State == "handed-off" {
		return teach.Write(deps.Stderr,
			"state=handed-off requires a PHASE-06 handoff packet",
			"lifecycle",
			"program/contracts/lifecycle.md",
			"stay in clarified, or mycelium state archived; packet command is not shipped",
		)
	}

	same := m.State == target
	if same {
		if target == "simmering" && opts.HasRevisit && opts.Revisit != m.Revisit {
			return commitGeneric(root, m, logBytes, target, opts.Revisit, opts, deps)
		}
		fmt.Fprintf(deps.Stdout, "already %s\n", target)
		return 0
	}

	if !lifecycle.Legal(m.State, target) {
		next := lifecycle.LegalNext(m.State)
		fix := "allowed next states: " + strings.Join(next, ", ")
		if len(next) == 0 {
			fix = "no commanded next states from " + m.State
		}
		return teach.Write(deps.Stderr,
			fmt.Sprintf("illegal transition %s → %s", m.State, target),
			"lifecycle",
			"program/contracts/lifecycle.md",
			fix,
		)
	}

	if lifecycle.IsWake(m.State, target) {
		if _, _, _, err := revisit.Parse(m.Revisit); err != nil {
			return teach.Write(deps.Stderr,
				fmt.Sprintf("revisit %q is not a date or event:<kebab>", m.Revisit),
				"revisit",
				"program/contracts/revisit.md",
				"use YYYY-MM-DD (UTC) or event:after-iphone-launch",
			)
		}
		return commitWake(root, m, logBytes, opts, deps)
	}

	revisitVal := ""
	if target == "simmering" {
		revisitVal = opts.Revisit
	}
	return commitGeneric(root, m, logBytes, target, revisitVal, opts, deps)
}

func commitWake(root string, m manifest.Manifest, logBytes []byte, opts Options, deps Deps) int {
	now := deps.Clock.Now().UTC()
	date := clock.Date(now)
	brief, err := wakebrief.Collect(root, m, logBytes, now)
	if err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("cannot collect wake brief: %v", err),
			"wake",
			"program/contracts/wake.md",
			"fix evidence/ and assumptions/ readability and retry",
		)
	}
	briefBytes := wakebrief.Render(brief)
	datedRel := wakebrief.DatedPath(now)
	logLine := logfmt.Line(date, "wake", "-", "exploring")

	m.State = "exploring"
	m.Revisit = ""
	m.UpdatedDate = date
	manOut, err := manifest.Encode(m)
	if err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("cannot encode manifest: %v", err),
			"manifest",
			"program/contracts/manifest.md",
			"fix mycelium.toml and retry",
		)
	}
	newLog := appendLogLine(logBytes, logLine)

	idx, err := indexmd.Load(root)
	if err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("cannot build index.md: %v", err),
			"index",
			"program/contracts/index.md",
			"restore mycelium.toml and log.md, then retry",
		)
	}
	idx.State = "exploring"
	idx.Revisit = ""
	idx.Wake = "briefs/LATEST.md"
	idx.LogLines = append(idx.LogLines, logLine)

	argv := opts.Argv
	if len(argv) == 0 {
		if opts.Wake {
			argv = []string{"wake"}
		} else {
			argv = []string{"state", "exploring"}
		}
		if opts.Dir != "" {
			argv = append(argv, "--dir", opts.Dir)
		}
	}

	files := []op.Staged{
		{RelTo: datedRel, Content: briefBytes},
		{RelTo: "briefs/LATEST.md", Content: append([]byte(nil), briefBytes...)},
		{RelTo: "index.md", Content: indexmd.Render(idx)},
		{RelTo: "log.md", Content: newLog},
		{RelTo: "mycelium.toml", Content: manOut},
	}
	if code := protocolCommit(root, "wake", logLine, argv, files, now, deps); code != 0 {
		return code
	}
	fmt.Fprintf(deps.Stdout, "woke %s\n", datedRel)
	fmt.Fprintln(deps.Stdout, "state: exploring")
	return 0
}

func commitGeneric(root string, m manifest.Manifest, logBytes []byte, target, revisitVal string, opts Options, deps Deps) int {
	now := deps.Clock.Now().UTC()
	date := clock.Date(now)
	note := target
	if target == "simmering" {
		note = "simmering revisit=" + revisitVal
	}
	logLine := logfmt.Line(date, "state", "-", note)

	m.State = target
	m.Revisit = revisitVal
	m.UpdatedDate = date
	manOut, err := manifest.Encode(m)
	if err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("cannot encode manifest: %v", err),
			"manifest",
			"program/contracts/manifest.md",
			"fix mycelium.toml and retry",
		)
	}
	newLog := appendLogLine(logBytes, logLine)

	idx, err := indexmd.Load(root)
	if err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("cannot build index.md: %v", err),
			"index",
			"program/contracts/index.md",
			"restore mycelium.toml and log.md, then retry",
		)
	}
	idx.State = target
	idx.Revisit = revisitVal
	idx.LogLines = append(idx.LogLines, logLine)

	argv := opts.Argv
	if len(argv) == 0 {
		argv = []string{"state", target}
		if opts.Dir != "" {
			argv = append(argv, "--dir", opts.Dir)
		}
		if target == "simmering" {
			argv = append(argv, "--revisit", revisitVal)
		}
	}

	files := []op.Staged{
		{RelTo: "index.md", Content: indexmd.Render(idx)},
		{RelTo: "log.md", Content: newLog},
		{RelTo: "mycelium.toml", Content: manOut},
	}
	if code := protocolCommit(root, "state", logLine, argv, files, now, deps); code != 0 {
		return code
	}
	fmt.Fprintf(deps.Stdout, "state: %s\n", target)
	fmt.Fprintf(deps.Stdout, "revisit: %s\n", revisitVal)
	return 0
}

func protocolCommit(root, opName, logLine string, argv []string, files []op.Staged, now time.Time, deps Deps) int {
	sess, err := op.Begin(root, op.Intent{
		Op:      opName,
		LogLine: logLine,
		Argv:    argv,
	}, now)
	if err != nil {
		if errors.Is(err, op.ErrJournalMismatch) {
			return teach.Write(deps.Stderr,
				"leftover journal for a different operation",
				"operation-protocol",
				"program/contracts/operation-protocol.md",
				"re-run the original command to complete, or mycelium check --abort-journal to roll back",
			)
		}
		if errors.Is(err, op.ErrLocked) {
			return teach.Write(deps.Stderr,
				fmt.Sprintf("lock held by another process: %v", err),
				"operation-protocol",
				"program/contracts/operation-protocol.md",
				"wait for the other process to finish; do not force the lock",
			)
		}
		return teach.Write(deps.Stderr,
			fmt.Sprintf("operation begin failed: %v", err),
			"operation-protocol",
			"program/contracts/operation-protocol.md",
			"re-run the same command, or mycelium check --abort-journal",
		)
	}

	if len(sess.Journal().Renames) == 0 {
		if err := sess.Stage(files); err != nil {
			rollbackOrClose(sess)
			return teach.Write(deps.Stderr,
				fmt.Sprintf("stage failed: %v", err),
				"operation-protocol",
				"program/contracts/operation-protocol.md",
				"re-run the same command after freeing disk space",
			)
		}
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
	_ = os.Remove(filepath.Join(root, ".mycelium"))
	return 0
}

func resolveRoot(opts Options, deps Deps) (string, int) {
	cwd := opts.Cwd
	var err error
	if cwd == "" {
		cwd, err = os.Getwd()
		if err != nil {
			return "", teach.Write(deps.Stderr,
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
		return "", teach.Write(deps.Stderr,
			"not a mycelium instance (no mycelium.toml found)",
			"instance-root",
			"program/contracts/manifest.md",
			"run from an instance directory or pass --dir PATH",
		)
	}
	return root, 0
}

func isAllowedTarget(t string) bool {
	for _, a := range lifecycle.AllowedTargets() {
		if a == t {
			return true
		}
	}
	return false
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
