// Package handoffcmd implements mycelium handoff (packet generator + state flip).
package handoffcmd

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/robertguss/mycelium/internal/check"
	"github.com/robertguss/mycelium/internal/clock"
	myceliumembed "github.com/robertguss/mycelium/internal/embed"
	"github.com/robertguss/mycelium/internal/handoff"
	"github.com/robertguss/mycelium/internal/indexmd"
	"github.com/robertguss/mycelium/internal/logfmt"
	"github.com/robertguss/mycelium/internal/manifest"
	"github.com/robertguss/mycelium/internal/metadata"
	"github.com/robertguss/mycelium/internal/op"
	"github.com/robertguss/mycelium/internal/teach"
)

const (
	transitionTitle = "clarified -> handed-off"
	defaultSystem   = "pstack/poteto"
	defaultBudget   = "30m"
)

var (
	mdLinkRE  = regexp.MustCompile(`\[([^\]]*)\]\(([^)]+)\)`)
	idTokenRE = regexp.MustCompile(`\b((?:DEC|OQ|EVD)-[0-9]+)\b`)
	h2LineRE  = regexp.MustCompile(`(?m)^## (.+)$`)
)

// Options for one handoff operation.
type Options struct {
	Dir  string
	Cwd  string
	Argv []string
}

// Deps are injectable collaborators.
type Deps struct {
	Clock  clock.Clock
	Stdout io.Writer
	Stderr io.Writer
}

// Run executes mycelium handoff. Exit 0 / 1.
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

	if packetPasses(root) {
		return teach.Write(deps.Stderr,
			"handoff/PACKET.md already exists",
			"lifecycle",
			"program/contracts/lifecycle.md",
			"mycelium state handed-off [--dir PATH]",
		)
	}

	if m.State != "clarified" {
		return teach.Write(deps.Stderr,
			fmt.Sprintf("handoff is legal only from clarified (got %s)", m.State),
			"lifecycle",
			"program/contracts/lifecycle.md",
			"mycelium state clarified, then mycelium handoff",
		)
	}

	now := deps.Clock.Now().UTC()
	date := clock.Date(now)
	logLine := logfmt.Line(date, "handoff", handoff.PacketID, transitionTitle)

	argv := opts.Argv
	if len(argv) == 0 {
		argv = []string{"handoff"}
		if opts.Dir != "" {
			argv = append(argv, "--dir", opts.Dir)
		}
	}

	sess, err := op.Begin(root, op.Intent{
		Op:      "handoff",
		Title:   transitionTitle,
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
		files, err := buildPacket(root, m, date)
		if err != nil {
			rollbackOrClose(sess)
			return teach.Write(deps.Stderr,
				fmt.Sprintf("cannot build handoff packet: %v", err),
				"handoff-packet",
				"program/contracts/handoff-packet.md",
				"fix instance artifacts and retry mycelium handoff",
			)
		}

		m.State = "handed-off"
		m.Revisit = ""
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
		idx.State = "handed-off"
		idx.Revisit = ""
		idx.LogLines = append(idx.LogLines, logLine)

		staged := make([]op.Staged, 0, len(files)+3)
		staged = append(staged, files...)
		staged = append(staged,
			op.Staged{RelTo: "index.md", Content: indexmd.Render(idx)},
			op.Staged{RelTo: "log.md", Content: newLog},
			op.Staged{RelTo: "mycelium.toml", Content: manOut},
		)
		if err := sess.Stage(staged); err != nil {
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

	fmt.Fprintln(deps.Stdout, "mycelium handoff: ok")
	fmt.Fprintln(deps.Stdout, "state: handed-off")
	fmt.Fprintln(deps.Stdout, "packet: handoff/PACKET.md")
	return 0
}

type namedFile struct {
	Name string
	Data []byte
	ID   string
	Meta map[string]any
}

func buildPacket(root string, m manifest.Manifest, date string) ([]op.Staged, error) {
	decs, err := listAcceptedDecisions(root)
	if err != nil {
		return nil, err
	}
	oqs, err := listQuestions(root)
	if err != nil {
		return nil, err
	}
	evds, err := collectEvidence(root, decs)
	if err != nil {
		return nil, err
	}
	playFiles, err := buildPlaybooks(root)
	if err != nil {
		return nil, err
	}
	accFiles, err := buildAcceptance(root)
	if err != nil {
		return nil, err
	}

	packetBody, err := renderPacket(m, date, decs, oqs, evds, playFiles, accFiles)
	if err != nil {
		return nil, err
	}

	var out []op.Staged
	out = append(out, op.Staged{RelTo: "handoff/PACKET.md", Content: []byte(packetBody)})

	for _, d := range decs {
		out = append(out, op.Staged{
			RelTo:   path.Join("handoff", "decisions", d.Name),
			Content: rewriteLinks(d.Data, "decisions"),
		})
	}
	glossary, err := buildGlossary(root)
	if err != nil {
		return nil, err
	}
	out = append(out, op.Staged{RelTo: "handoff/glossary.md", Content: glossary})

	for _, q := range oqs {
		out = append(out, op.Staged{
			RelTo:   path.Join("handoff", "questions", q.Name),
			Content: append([]byte(nil), q.Data...),
		})
	}

	out = append(out, op.Staged{
		RelTo:   "handoff/evidence/SUMMARY.md",
		Content: []byte(buildEvidenceSummary(evds)),
	})
	for _, e := range evds {
		out = append(out, op.Staged{
			RelTo:   path.Join("handoff", "evidence", e.Name),
			Content: append([]byte(nil), e.Data...),
		})
	}

	out = append(out, playFiles...)
	out = append(out, accFiles...)

	if len(decs) == 0 {
		out = append(out, op.Staged{RelTo: "handoff/decisions/.keep", Content: []byte{}})
	}
	if len(oqs) == 0 {
		out = append(out, op.Staged{RelTo: "handoff/questions/.keep", Content: []byte{}})
	}
	return out, nil
}

func renderPacket(m manifest.Manifest, date string, decs, oqs, evds []namedFile, plays, accs []op.Staged) (string, error) {
	tpl, err := myceliumembed.Program.ReadFile("program/templates/handoff-packet.md")
	if err != nil {
		return "", err
	}
	s := string(tpl)
	s = strings.ReplaceAll(s, "{{ID}}", handoff.PacketID)
	s = strings.ReplaceAll(s, "{{DATE}}", date)

	bodies := map[string]string{
		"Framing":                 "Idea: " + m.IdeaName + ".",
		"Locked decisions":        formatDECList(decs),
		"Glossary":                "See [glossary.md](glossary.md).",
		"Open questions":          formatOQList(oqs),
		"Evidence summary":        formatEvidence(evds),
		"Implementation playbooks": formatPathList(plays, "playbooks/"),
		"Implementation system":   defaultSystem,
		"Time budget":             defaultBudget,
		"Acceptance":              formatPathList(accs, "acceptance/"),
	}
	return setH2Bodies(s, bodies), nil
}

func formatDECList(decs []namedFile) string {
	if len(decs) == 0 {
		return "none"
	}
	var b strings.Builder
	for i, d := range decs {
		if i > 0 {
			b.WriteByte('\n')
		}
		line := "- " + d.ID
		if t := metaString(d.Meta["title"]); t != "" {
			line += " — " + t
		}
		line += " — see [decisions/" + d.Name + "](decisions/" + d.Name + ")"
		b.WriteString(line)
	}
	return b.String()
}

func formatOQList(oqs []namedFile) string {
	if len(oqs) == 0 {
		return "none"
	}
	var b strings.Builder
	for i, q := range oqs {
		if i > 0 {
			b.WriteByte('\n')
		}
		line := "- " + q.ID
		if a := metaString(q.Meta["agreement"]); a != "" {
			line += " (agreement=" + a + ")"
		}
		line += " — see [questions/" + q.Name + "](questions/" + q.Name + ")"
		b.WriteString(line)
	}
	return b.String()
}

func formatEvidence(evds []namedFile) string {
	if len(evds) == 0 {
		return "none"
	}
	return "See [evidence/SUMMARY.md](evidence/SUMMARY.md)."
}

func formatPathList(files []op.Staged, prefix string) string {
	var lines []string
	for _, f := range files {
		rel := strings.TrimPrefix(f.RelTo, "handoff/")
		if strings.HasPrefix(rel, prefix) {
			lines = append(lines, "- see ["+rel+"]("+rel+")")
		}
	}
	if len(lines) == 0 {
		return "none"
	}
	return strings.Join(lines, "\n")
}

func setH2Bodies(doc string, bodies map[string]string) string {
	matches := h2LineRE.FindAllStringSubmatchIndex(doc, -1)
	if len(matches) == 0 {
		return doc
	}
	var b strings.Builder
	prev := 0
	for i, m := range matches {
		title := doc[m[2]:m[3]]
		bodyEnd := len(doc)
		if i+1 < len(matches) {
			bodyEnd = matches[i+1][0]
		}
		b.WriteString(doc[prev:m[0]])
		if body, ok := bodies[title]; ok {
			b.WriteString("## ")
			b.WriteString(title)
			b.WriteString("\n\n")
			b.WriteString(strings.TrimRight(body, "\n"))
			b.WriteString("\n\n")
		} else {
			b.WriteString(doc[m[0]:bodyEnd])
		}
		prev = bodyEnd
	}
	b.WriteString(doc[prev:])
	return strings.TrimRight(b.String(), "\n") + "\n"
}

func buildGlossary(root string) ([]byte, error) {
	b, err := os.ReadFile(filepath.Join(root, "CONTEXT.md"))
	if err != nil {
		if os.IsNotExist(err) {
			return []byte("none\n"), nil
		}
		return nil, err
	}
	if strings.TrimSpace(string(b)) == "" {
		return []byte("none\n"), nil
	}
	return append([]byte(nil), b...), nil
}

func buildEvidenceSummary(evds []namedFile) string {
	if len(evds) == 0 {
		return "none\n"
	}
	var b strings.Builder
	b.WriteString("# Evidence summary\n\n")
	for _, e := range evds {
		b.WriteString("- ")
		b.WriteString(e.ID)
		b.WriteString(" — see ")
		b.WriteString(e.Name)
		b.WriteByte('\n')
	}
	return b.String()
}

func buildPlaybooks(root string) ([]op.Staged, error) {
	src := filepath.Join(root, "playbooks")
	st, err := os.Stat(src)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	} else if st.IsDir() {
		files, err := readDirFiles(src)
		if err != nil {
			return nil, err
		}
		if len(files) > 0 {
			var out []op.Staged
			for _, f := range files {
				out = append(out, op.Staged{
					RelTo:   path.Join("handoff", "playbooks", f.Name),
					Content: rewriteLinks(f.Data, "playbooks"),
				})
			}
			return out, nil
		}
	}
	tpl, err := myceliumembed.Program.ReadFile("program/templates/handoff-playbook.md")
	if err != nil {
		return nil, err
	}
	return []op.Staged{{
		RelTo:   "handoff/playbooks/PLAYBOOK.md",
		Content: append([]byte(nil), tpl...),
	}}, nil
}

func buildAcceptance(root string) ([]op.Staged, error) {
	src := filepath.Join(root, "acceptance")
	st, err := os.Stat(src)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	} else if st.IsDir() {
		files, err := readDirFiles(src)
		if err != nil {
			return nil, err
		}
		if len(files) > 0 {
			var out []op.Staged
			for _, f := range files {
				out = append(out, op.Staged{
					RelTo:   path.Join("handoff", "acceptance", f.Name),
					Content: append([]byte(nil), f.Data...),
				})
			}
			return out, nil
		}
	}
	return []op.Staged{{
		RelTo:   "handoff/acceptance/README.md",
		Content: []byte("none\n"),
	}}, nil
}

func listAcceptedDecisions(root string) ([]namedFile, error) {
	dir := filepath.Join(root, "decisions")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []namedFile
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		doc, err := metadata.Parse(data)
		if err != nil {
			continue
		}
		if metaString(doc.Meta["status"]) != "Accepted" {
			continue
		}
		out = append(out, namedFile{
			Name: e.Name(),
			Data: data,
			ID:   metaString(doc.Meta["id"]),
			Meta: doc.Meta,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func listQuestions(root string) ([]namedFile, error) {
	dir := filepath.Join(root, "questions")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []namedFile
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		doc, err := metadata.Parse(data)
		if err != nil {
			continue
		}
		out = append(out, namedFile{
			Name: e.Name(),
			Data: data,
			ID:   metaString(doc.Meta["id"]),
			Meta: doc.Meta,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func collectEvidence(root string, decs []namedFile) ([]namedFile, error) {
	want := map[string]struct{}{}
	for _, d := range decs {
		for _, m := range idTokenRE.FindAllStringSubmatch(string(d.Data), -1) {
			if strings.HasPrefix(m[1], "EVD-") {
				want[canonID(m[1])] = struct{}{}
			}
		}
	}
	dir := filepath.Join(root, "evidence")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	byID := map[string]namedFile{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		id := ""
		if doc, err := metadata.Parse(data); err == nil {
			id = metaString(doc.Meta["id"])
		}
		if id == "" {
			id = idFromFilename(e.Name())
		}
		if id == "" {
			continue
		}
		byID[canonID(id)] = namedFile{Name: e.Name(), Data: data, ID: id}
		want[canonID(id)] = struct{}{}
	}
	var out []namedFile
	for id := range want {
		if f, ok := byID[id]; ok {
			out = append(out, f)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func readDirFiles(dir string) ([]namedFile, error) {
	if _, err := os.Stat(dir); err != nil {
		return nil, err
	}
	var out []namedFile
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		out = append(out, namedFile{Name: filepath.ToSlash(rel), Data: data})
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// rewriteLinks rewrites markdown hrefs that escape the packet into in-packet paths.
func rewriteLinks(data []byte, fromDir string) []byte {
	return mdLinkRE.ReplaceAllFunc(data, func(m []byte) []byte {
		sub := mdLinkRE.FindSubmatch(m)
		if len(sub) < 3 {
			return m
		}
		text := string(sub[1])
		href := strings.TrimSpace(string(sub[2]))
		title := ""
		if i := strings.IndexAny(href, " \t"); i >= 0 {
			title = href[i:]
			href = href[:i]
		}
		lower := strings.ToLower(href)
		if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") ||
			strings.HasPrefix(lower, "mailto:") || strings.HasPrefix(href, "#") {
			return m
		}
		// Always map CONTEXT.md → glossary.md (packet copy name).
		base := path.Base(href)
		if base == "CONTEXT.md" {
			if fromDir == "playbooks" || fromDir == "acceptance" {
				href = "../glossary.md"
			} else {
				href = "glossary.md"
			}
		} else {
			joined := path.Clean(path.Join(fromDir, href))
			if joined == ".." || strings.HasPrefix(joined, "../") {
				alt := href
				for strings.HasPrefix(alt, "../") {
					alt = strings.TrimPrefix(alt, "../")
				}
				alt = strings.TrimPrefix(alt, "./")
				if fromDir == "playbooks" || fromDir == "acceptance" {
					href = "../" + alt
				} else {
					href = alt
				}
			}
		}
		out := "[" + text + "](" + href
		if title != "" {
			out += title
		}
		out += ")"
		return []byte(out)
	})
}

func packetPasses(root string) bool {
	if _, err := os.Stat(filepath.Join(root, "handoff", "PACKET.md")); err != nil {
		return false
	}
	return len(handoff.Check(os.DirFS(filepath.Join(root, "handoff")))) == 0
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
				"framework/phases/PHASE-06-implementation-brief.md",
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

func metaString(v any) string {
	if v == nil {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return fmt.Sprint(v)
	}
	return s
}

func canonID(id string) string {
	parts := strings.SplitN(id, "-", 2)
	if len(parts) != 2 {
		return id
	}
	n := strings.TrimLeft(parts[1], "0")
	if n == "" {
		n = "0"
	}
	for len(n) < 3 {
		n = "0" + n
	}
	return parts[0] + "-" + n
}

func idFromFilename(name string) string {
	base := strings.TrimSuffix(name, ".md")
	i := strings.Index(base, "-")
	if i < 0 {
		return ""
	}
	rest := base[i+1:]
	j := strings.Index(rest, "-")
	if j < 0 {
		return base
	}
	return base[:i+1+j]
}
