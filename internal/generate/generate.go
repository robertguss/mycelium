// Package generate implements mycelium new <type> (data-driven artifact emit).
package generate

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/robertguss/mycelium/internal/check"
	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/idpath"
	"github.com/robertguss/mycelium/internal/indexmd"
	"github.com/robertguss/mycelium/internal/journal"
	"github.com/robertguss/mycelium/internal/logfmt"
	"github.com/robertguss/mycelium/internal/manifest"
	"github.com/robertguss/mycelium/internal/op"
	"github.com/robertguss/mycelium/internal/schema"
	"github.com/robertguss/mycelium/internal/slug"
	"github.com/robertguss/mycelium/internal/teach"
)

var logLineRE = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}\t(scaffold|new|tier|publish|check)\t(\S+)\t`)

// Options for one new-<type> generation.
type Options struct {
	TypeKey string
	Title   string
	Dir     string // start path for FindRoot; empty → Cwd
	Cwd     string
	Argv    []string // journal argv (without program name)
}

// Deps are injectable collaborators.
type Deps struct {
	Clock  clock.Clock
	Stdout io.Writer
	Stderr io.Writer
}

// Run generates one artifact. Exit 0 on success, 1 on teaching error.
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

	typeKey := strings.TrimSpace(opts.TypeKey)
	title := strings.TrimSpace(opts.Title)
	if typeKey == "" {
		return teach.Write(deps.Stderr,
			"new requires a type key",
			"command-flags",
			"framework/phases/PHASE-01-implementation-brief.md",
			`usage: mycelium new <type> "<Title>" [--dir PATH]`,
		)
	}
	if title == "" {
		return teach.Write(deps.Stderr,
			"title is required",
			"command-flags",
			"framework/phases/PHASE-01-implementation-brief.md",
			`usage: mycelium new <type> "<Title>" [--dir PATH]`,
		)
	}
	if strings.ContainsAny(title, "\n\t") {
		return teach.Write(deps.Stderr,
			"title must not contain newline or tab characters",
			"log-injection",
			"program/contracts/conformance.md",
			`pass a single-line title without tab characters`,
		)
	}

	ideaSlug, err := slug.Slugify(title)
	if err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("cannot slugify title: %v", err),
			"slugify",
			"framework/decisions/DEC-014-phase-01-slugify-latin-fold.md",
			"pass a title with at least one letter or digit",
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

	keys, err := listTypeKeys(root)
	if err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("cannot list registered types: %v", err),
			"schema",
			"program/contracts/conformance.md",
			"restore program/templates/*.schema.toml",
		)
	}
	if !contains(keys, typeKey) {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("unknown type %q (registered: %s)", typeKey, strings.Join(keys, ", ")),
			"registered-types",
			"program/contracts/naming.md",
			"pass a registered type key from program/templates/*.schema.toml",
		)
	}

	sch, err := schema.Load(filepath.Join(root, "program", "templates", typeKey+".schema.toml"))
	if err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("cannot load schema for %q: %v", typeKey, err),
			"schema",
			"program/contracts/conformance.md",
			"restore program/templates/"+typeKey+".schema.toml",
		)
	}
	tplPath := filepath.Join(root, "program", "templates", typeKey+".md")
	tplBytes, err := os.ReadFile(tplPath)
	if err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("cannot read template for %q: %v", typeKey, err),
			"schema",
			"program/contracts/conformance.md",
			"restore program/templates/"+typeKey+".md",
		)
	}

	mb, err := os.ReadFile(filepath.Join(root, "mycelium.toml"))
	if err != nil {
		return teach.Write(deps.Stderr,
			"mycelium.toml missing or unreadable",
			"manifest",
			"program/contracts/manifest.md",
			"restore mycelium.toml",
		)
	}
	m, err := manifest.Parse(mb)
	if err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("mycelium.toml invalid: %v", err),
			"manifest",
			"program/contracts/manifest.md",
			"fix mycelium.toml required fields and values",
		)
	}
	logBytes, err := os.ReadFile(filepath.Join(root, "log.md"))
	if err != nil {
		return teach.Write(deps.Stderr,
			"log.md missing or unreadable",
			"log",
			"program/contracts/conformance.md",
			"restore log.md",
		)
	}
	if err := validateLog(logBytes); err != nil {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("log.md invalid: %v", err),
			"log-prefix",
			"program/contracts/conformance.md",
			"fix log.md prefixes",
		)
	}

	if sch.StageScoped {
		if _, ok := m.Identifiers[sch.Home]; !ok {
			return teach.Write(deps.Stderr,
				fmt.Sprintf("no [identifiers] range declared for %s", sch.Home),
				"stage-range",
				"program/contracts/identifiers.md",
				fmt.Sprintf("add [identifiers].%s = \"%s-001..%s-00N\" to mycelium.toml (DEC-013)", sch.Home, sch.Namespace, sch.Namespace),
			)
		}
	}

	now := deps.Clock.Now()
	date := clock.Date(now)
	argv := opts.Argv
	if len(argv) == 0 {
		argv = []string{"new", typeKey, title}
		if opts.Dir != "" {
			argv = append(argv, "--dir", opts.Dir)
		}
	}

	sess, err := op.Begin(root, op.Intent{
		Op:    "new",
		Type:  typeKey,
		Title: title,
		Argv:  argv,
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

	idStr := sess.OriginalID()
	var relPath string
	if idStr == "" {
		next, scanErr := nextID(root, sch)
		if scanErr != nil {
			rollbackOrClose(sess)
			return teach.Write(deps.Stderr,
				fmt.Sprintf("cannot scan %s: %v", sch.Home, scanErr),
				"id-allocation",
				"program/contracts/naming.md",
				"fix files under "+sch.Home+"/",
			)
		}
		if sch.StageScoped {
			if err := m.InRange(sch.Home, next); err != nil {
				rollbackOrClose(sess)
				return teach.Write(deps.Stderr,
					fmt.Sprintf("next ID %s-%0*d is outside declared range (%v)", sch.Namespace, sch.Digits, next, err),
					"stage-range",
					"program/contracts/identifiers.md",
					fmt.Sprintf("widen [identifiers].%s in mycelium.toml (DEC-013)", sch.Home),
				)
			}
		}
		idStr, err = idpath.FormatID(sch.Namespace, next)
		if err != nil {
			rollbackOrClose(sess)
			return teach.Write(deps.Stderr,
				fmt.Sprintf("cannot format ID: %v", err),
				"id-allocation",
				"program/contracts/naming.md",
				"report this as a CLI bug",
			)
		}
		relPath, err = idpath.Format(sch.Namespace, next, ideaSlug)
		if err != nil {
			rollbackOrClose(sess)
			return teach.Write(deps.Stderr,
				fmt.Sprintf("cannot format path: %v", err),
				"id-allocation",
				"program/contracts/naming.md",
				"report this as a CLI bug",
			)
		}
		sess.Journal().LogLine = logfmt.Line(date, "new", idStr, title)
		if err := sess.SetOriginalID(idStr); err != nil {
			rollbackOrClose(sess)
			return teach.Write(deps.Stderr,
				fmt.Sprintf("cannot write journal: %v", err),
				"operation-protocol",
				"program/contracts/operation-protocol.md",
				"re-run the same command, or mycelium check --abort-journal",
			)
		}
	} else {
		parsed, perr := idpath.Parse(idStr)
		if perr != nil {
			rollbackOrClose(sess)
			return teach.Write(deps.Stderr,
				fmt.Sprintf("journal original_id invalid: %v", perr),
				"operation-protocol",
				"program/contracts/operation-protocol.md",
				"mycelium check --abort-journal, then retry",
			)
		}
		relPath, err = idpath.Format(parsed.NS, parsed.N, ideaSlug)
		if err != nil {
			rollbackOrClose(sess)
			return teach.Write(deps.Stderr,
				fmt.Sprintf("cannot format path: %v", err),
				"id-allocation",
				"program/contracts/naming.md",
				"mycelium check --abort-journal, then retry",
			)
		}
		if sess.Journal().LogLine == "" {
			sess.Journal().LogLine = logfmt.Line(date, "new", idStr, title)
		}
	}

	if err := refuseOverwriteUnlessDone(root, relPath, sess.Journal().Renames); err != nil {
		if len(sess.Journal().Renames) > 0 {
			_ = sess.Close()
		} else {
			rollbackOrClose(sess)
		}
		return teach.Write(deps.Stderr,
			fmt.Sprintf("refuse overwrite: %s already exists", relPath),
			"overwrite",
			"program/contracts/naming.md",
			"rename the existing file, pick another title, or mycelium check --abort-journal",
		)
	}

	// Resume with existing renames: Stage is a no-op; Commit finishes.
	if len(sess.Journal().Renames) == 0 {
		body := ReplaceTokens(string(tplBytes), idStr, title, ideaSlug, date)
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
		newLog := appendLogLine(logBytes, sess.Journal().LogLine)
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
		idx.Inc(sch.Namespace)
		idx.LogLines = append(idx.LogLines, sess.Journal().LogLine)
		files := []op.Staged{
			{RelTo: relPath, Content: []byte(body)},
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

	fmt.Fprintf(deps.Stdout, "created %s\n", relPath)
	fmt.Fprintf(deps.Stdout, "next: fill required sections, then mycelium check\n")
	return 0
}

// ReplaceTokens substitutes {{ID}} {{TITLE}} {{SLUG}} {{DATE}}; other {{…}} stay.
func ReplaceTokens(tpl, id, title, slugStr, date string) string {
	out := tpl
	out = strings.ReplaceAll(out, "{{ID}}", id)
	out = strings.ReplaceAll(out, "{{TITLE}}", title)
	out = strings.ReplaceAll(out, "{{SLUG}}", slugStr)
	out = strings.ReplaceAll(out, "{{DATE}}", date)
	return out
}

// refuseOverwriteUnlessDone refuses when dest exists and the matching rename
// is not already Done (H2: resume must not clobber a user file).
// Filesystem already-Done (to exists, staged from gone, Done=false) is treated
// as Done so crash-after-artifact-rename can resume via Commit.
func refuseOverwriteUnlessDone(root, relPath string, renames []journal.Rename) error {
	dest := filepath.Join(root, filepath.FromSlash(relPath))
	_, err := os.Stat(dest)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	want := filepath.ToSlash(filepath.Clean(filepath.FromSlash(relPath)))
	for _, r := range renames {
		got := filepath.ToSlash(filepath.Clean(filepath.FromSlash(r.To)))
		if got != want {
			continue
		}
		if r.Done {
			return nil
		}
		from := filepath.Join(root, filepath.FromSlash(r.From))
		if _, ferr := os.Stat(from); os.IsNotExist(ferr) {
			return nil
		}
	}
	return errOverwrite
}

var errOverwrite = errors.New("generate: destination exists")

func rollbackOrClose(sess *op.Session) {
	if err := sess.Rollback(); errors.Is(err, op.ErrPartialCommit) {
		_ = sess.Close()
	}
}

func listTypeKeys(root string) ([]string, error) {
	dir := filepath.Join(root, "program", "templates")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var keys []string
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".schema.toml") {
			continue
		}
		key := strings.TrimSuffix(name, ".schema.toml")
		if _, err := os.Stat(filepath.Join(dir, key+".md")); err != nil {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	if len(keys) == 0 {
		return nil, errors.New("no type schemas found")
	}
	return keys, nil
}

func contains(keys []string, key string) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}

func validateLog(b []byte) error {
	for i, line := range strings.Split(string(b), "\n") {
		trim := strings.TrimRight(line, "\r")
		if trim == "" || strings.HasPrefix(trim, "#") {
			continue
		}
		if !logLineRE.MatchString(trim) {
			return fmt.Errorf("line %d has illegal prefix", i+1)
		}
	}
	return nil
}

func nextID(root string, sch schema.Schema) (int, error) {
	home := filepath.Join(root, sch.Home)
	entries, err := os.ReadDir(home)
	if err != nil {
		if os.IsNotExist(err) {
			return 1, nil
		}
		return 0, err
	}
	maxN := 0
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if name == "README.md" {
			continue
		}
		rel := sch.Home + "/" + name
		id, _, err := idpath.ParsePath(rel)
		if err != nil {
			continue
		}
		if id.NS != sch.Namespace {
			continue
		}
		if id.N > maxN {
			maxN = id.N
		}
	}
	return maxN + 1, nil
}

func appendLogLine(existing []byte, line string) []byte {
	s := string(existing)
	if !strings.HasSuffix(s, "\n") && s != "" {
		s += "\n"
	}
	s += line + "\n"
	return []byte(s)
}
