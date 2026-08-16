// Package supersedecmd implements mycelium supersede (artifact cross-links).
package supersedecmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/robertguss/mycelium/internal/check"
	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/idpath"
	"github.com/robertguss/mycelium/internal/indexmd"
	"github.com/robertguss/mycelium/internal/logfmt"
	"github.com/robertguss/mycelium/internal/manifest"
	"github.com/robertguss/mycelium/internal/metadata"
	"github.com/robertguss/mycelium/internal/op"
	"github.com/robertguss/mycelium/internal/supersede"
	"github.com/robertguss/mycelium/internal/teach"
)

// Options for one supersede operation.
type Options struct {
	OldID string
	NewID string
	Dir   string // start path for FindRoot; empty → Cwd
	Cwd   string
	Argv  []string // journal argv (without program name)
}

// Deps are injectable collaborators.
type Deps struct {
	Clock  clock.Clock
	Stdout io.Writer
	Stderr io.Writer
}

// Run supersedes OLD by NEW. Exit 0 / 1.
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

	oldTok := strings.TrimSpace(opts.OldID)
	newTok := strings.TrimSpace(opts.NewID)
	if oldTok == "" || newTok == "" {
		return teach.Write(deps.Stderr,
			"supersede requires <OLD-ID> --by <NEW-ID>",
			"command-flags",
			"program/contracts/conformance.md",
			"mycelium supersede DEC-001 --by DEC-002",
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
				"program/contracts/conformance.md",
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

	if r := supersede.CheckPair(oldTok, newTok, nil, nil); r != nil {
		return teachRefusal(deps.Stderr, r)
	}
	oldID, oldNS, _ := supersede.ParseID(oldTok)
	newID, _, _ := supersede.ParseID(newTok)

	oldRel, oldData, err := findArtifact(root, oldID, oldNS)
	if err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("no artifact %s", oldID),
			"id-to-path",
			"program/contracts/naming.md",
			`mycelium new <type> "…" then retry`,
		)
	}
	newRel, newData, err := findArtifact(root, newID, oldNS)
	if err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("no artifact %s", newID),
			"id-to-path",
			"program/contracts/naming.md",
			`mycelium new <type> "…" then retry`,
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
	logBytes, err := os.ReadFile(filepath.Join(root, "log.md"))
	if err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("cannot read log.md: %v", err),
			"log",
			"program/contracts/conformance.md",
			"restore log.md at the instance root",
		)
	}

	typ, err := idpath.LookupNS(oldNS)
	if err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("unknown namespace %s", oldNS),
			"id-to-path",
			"program/contracts/naming.md",
			"choose DEC, ASM, EVD, or SPK IDs",
		)
	}

	now := deps.Clock.Now().UTC()
	date := clock.Date(now)
	title := oldID + " -> " + newID
	logLine := logfmt.Line(date, "supersede", oldID, title)

	argv := opts.Argv
	if len(argv) == 0 {
		argv = []string{"supersede", oldID, "--by", newID}
		if opts.Dir != "" {
			argv = append(argv, "--dir", opts.Dir)
		}
	}

	sess, err := op.Begin(root, op.Intent{
		Op:         "supersede",
		Type:       typ.Key,
		Title:      title,
		OriginalID: oldID,
		LogLine:    logLine,
		Argv:       argv,
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
		oldDoc, err := metadata.Parse(oldData)
		if err != nil {
			rollbackOrClose(sess)
			return teach.Write(deps.Stderr,
				fmt.Sprintf("cannot parse %s: %v", oldID, err),
				"front-matter",
				"program/contracts/conformance.md",
				"fix +++ TOML front matter and retry",
			)
		}
		newDoc, err := metadata.Parse(newData)
		if err != nil {
			rollbackOrClose(sess)
			return teach.Write(deps.Stderr,
				fmt.Sprintf("cannot parse %s: %v", newID, err),
				"front-matter",
				"program/contracts/conformance.md",
				"fix +++ TOML front matter and retry",
			)
		}
		if r := supersede.CheckPair(oldTok, newTok, oldDoc.Meta, newDoc.Meta); r != nil {
			rollbackOrClose(sess)
			return teachRefusal(deps.Stderr, r)
		}

		oldOut, err := supersede.ApplyOld(oldData, newID)
		if err != nil {
			rollbackOrClose(sess)
			return teach.Write(deps.Stderr,
				fmt.Sprintf("cannot update %s: %v", oldID, err),
				"supersede",
				"program/contracts/conformance.md",
				"fix front matter and retry",
			)
		}
		newOut, err := supersede.ApplyNew(newData, oldID)
		if err != nil {
			rollbackOrClose(sess)
			return teach.Write(deps.Stderr,
				fmt.Sprintf("cannot update %s: %v", newID, err),
				"supersede",
				"program/contracts/conformance.md",
				"fix front matter and retry",
			)
		}

		m.UpdatedDate = date
		manOut, err := manifest.Encode(m)
		if err != nil {
			rollbackOrClose(sess)
			return teach.Write(deps.Stderr,
				fmt.Sprintf("cannot encode manifest: %v", err),
				"manifest",
				"program/contracts/manifest.md",
				"report this as a CLI bug",
			)
		}
		newLog := appendLogLine(logBytes, logLine)
		idx, err := indexmd.Load(root)
		if err != nil {
			rollbackOrClose(sess)
			return teach.Write(deps.Stderr,
				fmt.Sprintf("cannot build index.md: %v", err),
				"index",
				"program/contracts/index.md",
				"restore mycelium.toml and log.md, then retry",
			)
		}
		idx.LogLines = append(idx.LogLines, logLine)

		// Commit order: OLD, NEW, index.md, log.md, mycelium.toml last.
		files := []op.Staged{
			{RelTo: oldRel, Content: oldOut},
			{RelTo: newRel, Content: newOut},
			{RelTo: "index.md", Content: indexmd.Render(idx)},
			{RelTo: "log.md", Content: newLog},
			{RelTo: "mycelium.toml", Content: manOut},
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

	fmt.Fprintln(deps.Stdout, "mycelium supersede: ok")
	fmt.Fprintf(deps.Stdout, "old: %s\n", oldID)
	fmt.Fprintf(deps.Stdout, "new: %s\n", newID)
	return 0
}

func teachRefusal(stderr io.Writer, r *supersede.Refusal) int {
	return teach.Write(stderr,
		r.What,
		"supersede",
		"program/contracts/conformance.md",
		r.Fix,
	)
}

func findArtifact(root, idStr, ns string) (rel string, data []byte, err error) {
	t, err := idpath.LookupNS(ns)
	if err != nil {
		return "", nil, err
	}
	dir := filepath.Join(root, t.Home)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", nil, err
	}
	prefix := idStr + "-"
	for _, e := range entries {
		if e.IsDir() || !strings.HasPrefix(e.Name(), prefix) || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		rel = filepath.ToSlash(filepath.Join(t.Home, e.Name()))
		data, err = os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			return "", nil, err
		}
		return rel, data, nil
	}
	return "", nil, fmt.Errorf("not found")
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
