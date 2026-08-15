// Package check implements mycelium check (structure-only conformance).
package check

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pelletier/go-toml/v2"

	"github.com/robertguss/mycelium/internal/idpath"
	"github.com/robertguss/mycelium/internal/journal"
	"github.com/robertguss/mycelium/internal/lifecycle"
	"github.com/robertguss/mycelium/internal/lock"
	"github.com/robertguss/mycelium/internal/manifest"
	"github.com/robertguss/mycelium/internal/metadata"
	"github.com/robertguss/mycelium/internal/op"
	"github.com/robertguss/mycelium/internal/schema"
	"github.com/robertguss/mycelium/internal/teach"
)

// Finding is an alias for teach.Finding.
type Finding = teach.Finding

// Result is the outcome of a read-only check run.
type Result struct {
	OK             bool
	Slug           string
	State          string
	Tier           string
	Artifacts      int
	Findings       []Finding
	LiveLockNotice string
}

var (
	ErrNotInstance = errors.New("check: not a mycelium instance")
	linkRE         = regexp.MustCompile(`\b(DEC|ASM|EVD|SPK|FND|REC|REQ|OQ|RSK|PHASE|MS)-[0-9]+\b`)
	logLineRE      = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}\t(scaffold|new|tier|publish|check)\t(\S+)\t`)
	dateRE         = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	h2RE           = regexp.MustCompile(`(?m)^## (.+)$`)
)

var alwaysAllowedTop = map[string]struct{}{
	"README.md": {}, "mycelium.toml": {}, "log.md": {}, "CONTEXT.md": {},
	"AGENTS.md": {}, ".gitignore": {}, "LICENSE": {}, "CHANGELOG.md": {},
	".agents": {}, ".mycelium": {}, ".git": {}, ".github": {}, "program": {},
}

// LegalNext is the PHASE-02 commanded edge table (delegates to lifecycle).
func LegalNext(from string) []string {
	return lifecycle.LegalNext(from)
}

// FindRoot walks start upward looking for mycelium.toml.
// Checks toml before the .git stop so instance roots that have both still match.
func FindRoot(start string) (string, error) {
	cur, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(cur, "mycelium.toml")); err == nil {
			return cur, nil
		}
		if _, err := os.Stat(filepath.Join(cur, ".git")); err == nil {
			return "", ErrNotInstance
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			return "", ErrNotInstance
		}
		cur = parent
	}
}

// Run executes the conformance suite against an instance root.
// Runtime reads instance files only (never embed).
func Run(root string) Result {
	var r Result
	add := func(what, conv, contract, fix string) {
		r.Findings = append(r.Findings, Finding{
			What: what, Convention: conv, Contract: contract, Fix: fix,
		})
	}

	mb, err := os.ReadFile(filepath.Join(root, "mycelium.toml"))
	if err != nil {
		add("mycelium.toml missing or unreadable", "manifest", "program/contracts/manifest.md", "run mycelium new idea --offline or restore mycelium.toml")
		return r
	}
	m, err := manifest.Parse(mb)
	if err != nil {
		add(fmt.Sprintf("mycelium.toml invalid: %v", err), "manifest", "program/contracts/manifest.md", "fix mycelium.toml required fields and values")
		return r
	}
	r.Slug = m.Slug
	r.State = m.State
	r.Tier = m.Tier

	checkLifecycle(m, add)
	checkJournalLock(root, &r, add)

	schemas, err := loadSchemas(root)
	if err != nil {
		add(fmt.Sprintf("cannot load schemas: %v", err), "schema", "program/contracts/conformance.md", "restore program/templates/*.schema.toml")
		return r
	}
	homeByNS := map[string]string{}
	schemaByHome := map[string]schema.Schema{}
	homes := map[string]struct{}{}
	for _, s := range schemas {
		homeByNS[s.Namespace] = s.Home
		schemaByHome[s.Home] = s
		homes[s.Home] = struct{}{}
	}

	checkTierBinds(root, m.Tier, add)
	checkTopLevel(root, m, homes, add)

	index := buildArtifactIndex(root, schemas, &r, add)
	checkFrontMatterAndSections(root, index, schemaByHome, add)
	checkStageScoped(m, index, schemaByHome, add)
	checkLog(root, add)
	checkLinks(root, index, homeByNS, add)

	if len(r.Findings) == 0 {
		r.OK = true
	}
	return r
}

func checkLifecycle(m manifest.Manifest, add func(string, string, string, string)) {
	switch m.State {
	case "spark", "exploring", "simmering", "archived":
		// ok (simmering+revisit already enforced by manifest.Parse)
	case "clarified", "handed-off":
		add(
			fmt.Sprintf("state=%s is not reachable in PHASE-01", m.State),
			"lifecycle",
			"program/contracts/lifecycle.md",
			"restore state to spark|exploring|simmering|archived (PHASE-02/06 commands are not shipped)",
		)
	default:
		add(
			fmt.Sprintf("unknown state %q", m.State),
			"lifecycle",
			"program/contracts/lifecycle.md",
			"set state to spark|exploring|simmering|archived",
		)
	}
}

func checkJournalLock(root string, r *Result, add func(string, string, string, string)) {
	hasJournal, stale, err := op.Detect(root)
	if err != nil {
		add(fmt.Sprintf("cannot detect journal/lock: %v", err), "operation-protocol", "program/contracts/operation-protocol.md", "inspect .mycelium/")
		return
	}
	info, err := lock.Inspect(root)
	if err != nil {
		add(fmt.Sprintf("cannot inspect lock: %v", err), "operation-protocol", "program/contracts/operation-protocol.md", "inspect .mycelium/lock")
		return
	}
	if info.State == lock.Live {
		r.LiveLockNotice = fmt.Sprintf("mycelium check: live lock held by pid=%d (continuing)", info.PID)
	}
	if hasJournal || stale {
		add(
			"interrupted operation",
			"operation-protocol",
			"program/contracts/operation-protocol.md",
			"re-run the same command to complete, or mycelium check --abort-journal to roll back",
		)
	}
}

type artifactFile struct {
	Rel   string
	Home  string
	ID    idpath.ID
	IDStr string
}

func loadSchemas(root string) ([]schema.Schema, error) {
	dir := filepath.Join(root, "program", "templates")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []schema.Schema
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".schema.toml") {
			continue
		}
		s, err := schema.Load(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", e.Name(), err)
		}
		out = append(out, s)
	}
	if len(out) == 0 {
		return nil, errors.New("no *.schema.toml found")
	}
	return out, nil
}

type tierFile struct {
	Binds []string `toml:"binds"`
}

func checkTierBinds(root, tier string, add func(string, string, string, string)) {
	switch tier {
	case "focused", "standard", "high-assurance":
	default:
		add(fmt.Sprintf("unknown tier %q", tier), "tier", "program/contracts/conformance.md", "set tier to focused|standard|high-assurance")
		return
	}
	path := filepath.Join(root, "program", "tiers", tier+".toml")
	b, err := os.ReadFile(path)
	if err != nil {
		add(fmt.Sprintf("tier file missing: %s", path), "tier", "program/contracts/conformance.md", "restore program/tiers/"+tier+".toml")
		return
	}
	var tf tierFile
	if err := toml.Unmarshal(b, &tf); err != nil {
		add(fmt.Sprintf("tier file invalid: %v", err), "tier", "program/contracts/conformance.md", "fix program/tiers/"+tier+".toml")
		return
	}
	for _, bind := range tf.Binds {
		switch {
		case bind == "manifest":
			if _, err := os.Stat(filepath.Join(root, "mycelium.toml")); err != nil {
				add("missing bound path mycelium.toml", "tier-binds", "program/contracts/conformance.md", "restore mycelium.toml")
			}
		case strings.HasSuffix(bind, "/"):
			dir := strings.TrimSuffix(bind, "/")
			if _, err := os.Stat(filepath.Join(root, dir)); err != nil {
				add(fmt.Sprintf("missing bound directory %s/", dir), "tier-binds", "program/contracts/conformance.md", "restore "+dir+"/ or raise tier after emit")
			}
		default:
			if _, err := os.Stat(filepath.Join(root, bind)); err != nil {
				add(fmt.Sprintf("missing bound path %s", bind), "tier-binds", "program/contracts/conformance.md", "restore "+bind)
			}
		}
	}
}

func checkTopLevel(root string, m manifest.Manifest, homes map[string]struct{}, add func(string, string, string, string)) {
	entries, err := os.ReadDir(root)
	if err != nil {
		add(fmt.Sprintf("cannot read instance root: %v", err), "conformance", "program/contracts/conformance.md", "ensure the instance directory is readable")
		return
	}
	declared := map[string]struct{}{}
	for _, d := range m.Deviations {
		if strings.HasPrefix(d.Convention, "extra-top-level:") {
			declared[strings.TrimPrefix(d.Convention, "extra-top-level:")] = struct{}{}
		}
		if d.Convention == "" || d.Reason == "" {
			add("deviation row missing convention or reason", "deviation", "program/contracts/manifest.md", "fill convention and reason on every [[deviations]] row")
		}
	}
	for _, e := range entries {
		name := e.Name()
		if _, ok := alwaysAllowedTop[name]; ok {
			continue
		}
		if _, ok := homes[name]; ok {
			continue
		}
		if _, ok := declared[name]; ok {
			continue
		}
		add(
			fmt.Sprintf("undeclared extra top-level path %s", name),
			"extra-top-level",
			"program/contracts/conformance.md",
			fmt.Sprintf("remove %s or declare [[deviations]] convention = \"extra-top-level:%s\"", name, name),
		)
	}
}

func buildArtifactIndex(root string, schemas []schema.Schema, r *Result, add func(string, string, string, string)) []artifactFile {
	var out []artifactFile
	seen := map[string]string{} // IDStr → rel path
	for _, s := range schemas {
		dir := filepath.Join(root, s.Home)
		st, err := os.Stat(dir)
		if err != nil {
			continue // home absent OK at any tier
		}
		if !st.IsDir() {
			add(fmt.Sprintf("%s exists but is not a directory", s.Home), "id-to-path", "program/contracts/naming.md", "make "+s.Home+" a directory")
			continue
		}
		err = filepath.WalkDir(dir, func(path string, d os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() {
				return nil
			}
			name := d.Name()
			if name == "README.md" {
				return nil
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			rel = filepath.ToSlash(rel)
			id, _, err := idpath.ParsePath(rel)
			if err != nil {
				add(
					fmt.Sprintf("file %s does not match ID-to-path pattern", rel),
					"id-to-path",
					"program/contracts/naming.md",
					"rename to "+s.Namespace+"-<digits>-<slug>.md or remove",
				)
				return nil
			}
			idStr, err := idpath.FormatID(id.NS, id.N)
			if err != nil {
				add(fmt.Sprintf("cannot format id for %s: %v", rel, err), "id-to-path", "program/contracts/naming.md", "fix the filename digits")
				return nil
			}
			if prev, ok := seen[idStr]; ok {
				add(
					fmt.Sprintf("duplicate id %s in %s and %s", idStr, prev, rel),
					"id-uniqueness",
					"program/contracts/naming.md",
					"rename or remove one of the duplicate files",
				)
			} else {
				seen[idStr] = rel
			}
			out = append(out, artifactFile{Rel: rel, Home: s.Home, ID: id, IDStr: idStr})
			return nil
		})
		if err != nil {
			add(fmt.Sprintf("cannot walk %s/: %v", s.Home, err), "id-to-path", "program/contracts/naming.md", "fix permissions on "+s.Home+"/")
		}
	}
	r.Artifacts = len(out)
	return out
}

func checkFrontMatterAndSections(root string, arts []artifactFile, byHome map[string]schema.Schema, add func(string, string, string, string)) {
	for _, a := range arts {
		s, ok := byHome[a.Home]
		if !ok {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(a.Rel)))
		if err != nil {
			add(fmt.Sprintf("cannot read %s: %v", a.Rel, err), "front-matter", "program/contracts/conformance.md", "restore "+a.Rel)
			continue
		}
		doc, err := metadata.Parse(b)
		if err != nil {
			add(fmt.Sprintf("%s front matter invalid: %v", a.Rel, err), "front-matter", "program/contracts/conformance.md", "fix +++ TOML front matter in "+a.Rel)
			continue
		}
		if err := metadata.RequireKeys(doc.Meta, s.RequiredFrontMatter); err != nil {
			add(fmt.Sprintf("%s missing required front matter: %v", a.Rel, err), "front-matter", "program/contracts/conformance.md", "add required keys in "+a.Rel)
		}
		if idVal, ok := doc.Meta["id"].(string); ok && idVal != a.IDStr {
			add(
				fmt.Sprintf("%s front-matter id %q does not match filename id %s", a.Rel, idVal, a.IDStr),
				"id-to-path",
				"program/contracts/naming.md",
				"set id = \""+a.IDStr+"\" or rename the file",
			)
		}
		for field, vals := range s.Enums {
			v, ok := doc.Meta[field].(string)
			if !ok {
				continue
			}
			if !contains(vals, v) {
				add(
					fmt.Sprintf("%s field %s=%q not in enum", a.Rel, field, v),
					"front-matter-enum",
					"program/contracts/conformance.md",
					fmt.Sprintf("set %s to one of %v", field, vals),
				)
			}
		}
		for _, k := range s.RequiredFrontMatter {
			if k == "date" {
				if v, ok := doc.Meta["date"].(string); ok && !dateRE.MatchString(v) {
					add(fmt.Sprintf("%s date must be YYYY-MM-DD", a.Rel), "front-matter", "program/contracts/conformance.md", "fix date in "+a.Rel)
				}
			}
		}
		present := map[string]struct{}{}
		for _, m := range h2RE.FindAllStringSubmatch(doc.Body, -1) {
			present[m[1]] = struct{}{}
		}
		for _, sec := range s.RequiredSections {
			if _, ok := present[sec]; !ok {
				add(
					fmt.Sprintf("%s missing required H2 %q", a.Rel, sec),
					"required-sections",
					"program/contracts/conformance.md",
					"add ## "+sec+" to "+a.Rel,
				)
			}
		}
	}
}

func checkStageScoped(m manifest.Manifest, arts []artifactFile, byHome map[string]schema.Schema, add func(string, string, string, string)) {
	for _, a := range arts {
		s, ok := byHome[a.Home]
		if !ok || !s.StageScoped {
			continue
		}
		if err := m.InRange(s.Home, a.ID.N); err != nil {
			add(
				fmt.Sprintf("%s is outside declared range (%v)", a.IDStr, err),
				"stage-range",
				"program/contracts/identifiers.md",
				fmt.Sprintf("widen [identifiers].%s in mycelium.toml or move/remove the file", s.Home),
			)
		}
	}
}

func checkLog(root string, add func(string, string, string, string)) {
	b, err := os.ReadFile(filepath.Join(root, "log.md"))
	if err != nil {
		add("log.md missing or unreadable", "log", "program/contracts/conformance.md", "restore log.md")
		return
	}
	for i, line := range strings.Split(string(b), "\n") {
		trim := strings.TrimRight(line, "\r")
		if trim == "" || strings.HasPrefix(trim, "#") {
			continue
		}
		if !logLineRE.MatchString(trim) {
			add(
				fmt.Sprintf("log.md line %d has illegal prefix", i+1),
				"log-prefix",
				"program/contracts/conformance.md",
				"use YYYY-MM-DD\\t(scaffold|new|tier|publish|check)\\t<id>\\t…",
			)
		}
	}
}

func checkLinks(root string, arts []artifactFile, homeByNS map[string]string, add func(string, string, string, string)) {
	byNum := map[string]struct{}{}
	for _, a := range arts {
		byNum[fmt.Sprintf("%s-%d", a.ID.NS, a.ID.N)] = struct{}{}
	}

	scan := func(rel string) {
		b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			return
		}
		text := string(b)
		if doc, err := metadata.Parse(b); err == nil {
			text = doc.Body
		}
		for _, m := range linkRE.FindAllString(text, -1) {
			id, err := idpath.Parse(m)
			if err != nil {
				continue
			}
			key := fmt.Sprintf("%s-%d", id.NS, id.N)
			if _, ok := byNum[key]; ok {
				continue
			}
			home := homeByNS[id.NS]
			add(
				fmt.Sprintf("reference %s has no file", m),
				"id-to-path",
				"program/contracts/naming.md",
				fmt.Sprintf("add the artifact under %s/ or remove the reference", home),
			)
		}
	}

	for _, rel := range []string{"log.md", "README.md", "CONTEXT.md", "AGENTS.md"} {
		scan(rel)
	}
	for _, a := range arts {
		scan(a.Rel)
	}
}

func contains(ss []string, v string) bool {
	for _, s := range ss {
		if s == v {
			return true
		}
	}
	return false
}

// WriteOK prints the success summary.
func WriteOK(stdout io.Writer, r Result) {
	fmt.Fprintln(stdout, "mycelium check: ok")
	fmt.Fprintf(stdout, "instance: %s\n", r.Slug)
	fmt.Fprintf(stdout, "state: %s\n", r.State)
	fmt.Fprintf(stdout, "tier: %s\n", r.Tier)
	fmt.Fprintf(stdout, "artifacts: %d\n", r.Artifacts)
}

// AbortJournal rolls back staged temps / journal / stale lock.
// Prints surviving already-renamed paths. Returns teach findings on nothing-to-abort.
func AbortJournal(root string, stdout io.Writer) error {
	var survivors []string
	j, err := journal.Load(root)
	if err == nil {
		for _, r := range j.Renames {
			if r.Done {
				survivors = append(survivors, r.To)
			}
		}
	} else if !errors.Is(err, journal.ErrNotExist) {
		return err
	}
	if err := op.Abort(root); err != nil {
		return err
	}
	for _, p := range survivors {
		fmt.Fprintf(stdout, "surviving: %s\n", p)
	}
	return nil
}
