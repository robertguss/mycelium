// Package tiercmd implements mycelium tier (raise/lower/repair emit dirs).
package tiercmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"

	"github.com/robertguss/mycelium/internal/check"
	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/lock"
	"github.com/robertguss/mycelium/internal/logfmt"
	"github.com/robertguss/mycelium/internal/manifest"
	"github.com/robertguss/mycelium/internal/op"
	"github.com/robertguss/mycelium/internal/schema"
	"github.com/robertguss/mycelium/internal/teach"
)

var knownTiers = []string{"focused", "standard", "high-assurance"}

// Options for one tier operation.
type Options struct {
	Tier string
	Dir  string // start path for FindRoot; empty → Cwd
	Cwd  string
	Argv []string // journal argv (without program name)
}

// Deps are injectable collaborators.
type Deps struct {
	Clock  clock.Clock
	Stdout io.Writer
	Stderr io.Writer
}

type tierFile struct {
	Name  string   `toml:"name"`
	Emits []string `toml:"emits"`
}

// Run sets the instance tier and emits newly required dirs. Exit 0 / 1.
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

	want := strings.TrimSpace(opts.Tier)
	if want == "" {
		return teach.Write(deps.Stderr,
			"tier requires a tier name",
			"command-flags",
			"framework/phases/PHASE-01-implementation-brief.md",
			"usage: mycelium tier <tier> [--dir PATH] where <tier> is focused, standard, or high-assurance",
		)
	}
	if !isKnownTier(want) {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("unknown tier %q", want),
			"tier",
			"program/reference/rigor-tiers.md",
			"use one of: focused, standard, high-assurance",
		)
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

	tf, err := loadTierFile(root, want)
	if err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("cannot load program/tiers/%s.toml: %v", want, err),
			"tier",
			"program/contracts/conformance.md",
			"restore program/tiers/"+want+".toml from the methodology tree",
		)
	}

	missing, err := missingEmitDirs(root, tf.Emits)
	if err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("cannot inspect emit dirs: %v", err),
			"tier",
			"framework/phases/PHASE-01-implementation-brief.md",
			"fix filesystem permissions under the instance root",
		)
	}

	sameTier := m.Tier == want
	if sameTier && len(missing) == 0 {
		if code := refuseBusyNoop(root, deps); code != 0 {
			return code
		}
		fmt.Fprintf(deps.Stdout, "already %s\n", want)
		return 0
	}

	homeNS, err := loadHomeNamespaces(root)
	if err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("cannot load type schemas: %v", err),
			"schema",
			"program/contracts/conformance.md",
			"restore program/templates/*.schema.toml",
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

	now := deps.Clock.Now().UTC()
	date := clock.Date(now)
	fromTier := m.Tier
	note := fromTier + " -> " + want
	logLine := logfmt.Line(date, "tier", "-", note)

	argv := opts.Argv
	if len(argv) == 0 {
		argv = []string{"tier", want}
		if opts.Dir != "" {
			argv = append(argv, "--dir", opts.Dir)
		}
	}

	sess, err := op.Begin(root, op.Intent{
		Op:      "tier",
		Title:   want,
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
		files, buildErr := buildStaged(root, m, want, date, missing, homeNS, logBytes, logLine)
		if buildErr != nil {
			rollbackOrClose(sess)
			return teach.Write(deps.Stderr,
				buildErr.Error(),
				"tier",
				"framework/phases/PHASE-01-implementation-brief.md",
				"restore missing schemas or fix emit paths and retry",
			)
		}
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

	fmt.Fprintln(deps.Stdout, note)
	return 0
}

func buildStaged(
	root string,
	m manifest.Manifest,
	want, date string,
	missing []string,
	homeNS map[string]string,
	logBytes []byte,
	logLine string,
) ([]op.Staged, error) {
	var files []op.Staged
	for _, dir := range missing {
		ns, ok := homeNS[dir]
		if !ok {
			return nil, fmt.Errorf("no schema namespace for emit dir %q", dir)
		}
		readmeRel := filepath.ToSlash(filepath.Join(dir, "README.md"))
		files = append(files, op.Staged{
			RelTo:   readmeRel,
			Content: []byte(readmeStub(dir, ns)),
		})
	}

	m.Tier = want
	m.UpdatedDate = date
	manOut, err := manifest.Encode(m)
	if err != nil {
		return nil, fmt.Errorf("cannot encode manifest: %w", err)
	}
	newLog := appendLogLine(logBytes, logLine)
	files = append(files,
		op.Staged{RelTo: "log.md", Content: newLog},
		op.Staged{RelTo: "mycelium.toml", Content: manOut},
	)
	return files, nil
}

func readmeStub(dir, ns string) string {
	heading := strings.ToUpper(dir[:1]) + dir[1:]
	return fmt.Sprintf("# %s\n\nHome for %s-### artifacts.\n", heading, ns)
}

func appendLogLine(existing []byte, line string) []byte {
	s := string(existing)
	if !strings.HasSuffix(s, "\n") && s != "" {
		s += "\n"
	}
	s += line + "\n"
	return []byte(s)
}

func isKnownTier(name string) bool {
	for _, t := range knownTiers {
		if t == name {
			return true
		}
	}
	return false
}

func loadTierFile(root, name string) (tierFile, error) {
	path := filepath.Join(root, "program", "tiers", name+".toml")
	b, err := os.ReadFile(path)
	if err != nil {
		return tierFile{}, err
	}
	var tf tierFile
	if err := toml.Unmarshal(b, &tf); err != nil {
		return tierFile{}, err
	}
	if tf.Name != "" && tf.Name != name {
		return tierFile{}, fmt.Errorf("name %q does not match file %s", tf.Name, name)
	}
	return tf, nil
}

func emitDirName(raw string) string {
	s := strings.TrimSpace(raw)
	s = strings.TrimSuffix(s, "/")
	return s
}

func missingEmitDirs(root string, emits []string) ([]string, error) {
	var missing []string
	for _, e := range emits {
		dir := emitDirName(e)
		if dir == "" {
			continue
		}
		abs := filepath.Join(root, dir)
		fi, err := os.Stat(abs)
		if err != nil {
			if os.IsNotExist(err) {
				missing = append(missing, dir)
				continue
			}
			return nil, err
		}
		if !fi.IsDir() {
			return nil, fmt.Errorf("%s exists but is not a directory", dir)
		}
	}
	return missing, nil
}

func loadHomeNamespaces(root string) (map[string]string, error) {
	dir := filepath.Join(root, "program", "templates")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	out := map[string]string{}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".schema.toml") {
			continue
		}
		key := strings.TrimSuffix(name, ".schema.toml")
		sch, err := schema.Load(filepath.Join(dir, key+".schema.toml"))
		if err != nil {
			return nil, err
		}
		if sch.Home != "" && sch.Namespace != "" {
			out[sch.Home] = sch.Namespace
		}
	}
	return out, nil
}

func refuseBusyNoop(root string, deps Deps) int {
	hasJournal, _, err := op.Detect(root)
	if err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("cannot inspect operation state: %v", err),
			"operation-protocol",
			"program/contracts/operation-protocol.md",
			"fix .mycelium/ and retry",
		)
	}
	if hasJournal {
		return teach.Write(deps.Stderr,
			"leftover journal blocks tier",
			"operation-protocol",
			"program/contracts/operation-protocol.md",
			"re-run the original command to complete, or mycelium check --abort-journal to roll back",
		)
	}
	info, err := lock.Inspect(root)
	if err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("cannot inspect lock: %v", err),
			"operation-protocol",
			"program/contracts/operation-protocol.md",
			"fix .mycelium/lock and retry",
		)
	}
	if info.State == lock.Live {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("lock held by another process: pid=%d", info.PID),
			"operation-protocol",
			"program/contracts/operation-protocol.md",
			"wait for the other process to finish; do not force the lock",
		)
	}
	return 0
}

func rollbackOrClose(sess *op.Session) {
	if err := sess.Rollback(); errors.Is(err, op.ErrPartialCommit) {
		_ = sess.Close()
	}
}
