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
	"github.com/robertguss/mycelium/internal/logfmt"
	"github.com/robertguss/mycelium/internal/manifest"
	"github.com/robertguss/mycelium/internal/metadata"
	"github.com/robertguss/mycelium/internal/op"
	"github.com/robertguss/mycelium/internal/pack"
	"github.com/robertguss/mycelium/internal/revisit"
	"github.com/robertguss/mycelium/internal/schema"
	"github.com/robertguss/mycelium/internal/sparring"
	"github.com/robertguss/mycelium/internal/teach"
	"github.com/robertguss/mycelium/internal/wakebrief"
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
	linkRE         = regexp.MustCompile(`\b(DEC|ASM|EVD|SPK|FND|REC|REQ|OQ|RSK|PHASE|MS|CMP|RPT|RCL)-[0-9]+\b`)
	logLineRE      = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}\t(scaffold|new|tier|publish|check|state|wake|supersede)\t(\S+)\t`)
	dateRE         = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	h2RE           = regexp.MustCompile(`(?m)^## (.+)$`)
)

var alwaysAllowedTop = map[string]struct{}{
	"README.md": {}, "mycelium.toml": {}, "log.md": {}, "CONTEXT.md": {},
	"AGENTS.md": {}, ".gitignore": {}, "LICENSE": {}, "CHANGELOG.md": {},
	"index.md": {}, "briefs": {},
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
		add(fmt.Sprintf("cannot load schemas: %v", err), "schema", "program/contracts/conformance.md", "restore registered *.schema.toml files under program/")
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

	packResult, err := pack.Discover(filepath.Join(root, "program"))
	if err != nil {
		add(fmt.Sprintf("cannot load packs: %v", err), "pack", "program/contracts/conformance.md", "restore program/packs/ or program/templates/")
		return r
	}
	for _, c := range packResult.Collisions {
		add(
			c.Message(),
			"pack-collision",
			"program/contracts/conformance.md",
			"remove the colliding pack directory",
		)
	}
	reviewsOK := pack.ReviewsAllowed(packResult.Packs)

	checkTierBinds(root, m.Tier, add)
	checkTopLevel(root, m, homes, reviewsOK, add)
	checkIndex(root, add)

	index := buildArtifactIndex(root, schemas, &r, add)
	checkFrontMatterAndSections(root, index, schemaByHome, add)
	checkQuestionSparring(root, index, add)
	checkGlossary(root, add)
	checkDissent(root, index, add)
	checkLadder(root, index, packResult.Packs, add)
	checkStageScoped(m, index, schemaByHome, add)
	logBytes := checkLog(root, add)
	checkWakeBrief(root, logBytes, add)
	checkLinks(root, index, homeByNS, add)
	checkSupersedeIFF(root, index, add)

	if len(r.Findings) == 0 {
		r.OK = true
	}
	return r
}

func checkLifecycle(m manifest.Manifest, add func(string, string, string, string)) {
	switch m.State {
	case "spark", "exploring", "clarified", "archived":
		// ok; leftover revisit is not a fail
	case "simmering":
		if _, _, _, err := revisit.Parse(m.Revisit); err != nil {
			add(
				fmt.Sprintf("state=simmering requires revisit matching revisit grammar (got %q)", m.Revisit),
				"revisit",
				"program/contracts/revisit.md",
				"set revisit to YYYY-MM-DD (UTC) or event:<kebab>",
			)
		}
	case "handed-off":
		add(
			"state=handed-off requires a PHASE-06 handoff packet",
			"lifecycle",
			"program/contracts/lifecycle.md",
			"stay in clarified, or mycelium state archived; packet command is not shipped",
		)
	default:
		add(
			fmt.Sprintf("unknown state %q", m.State),
			"lifecycle",
			"program/contracts/lifecycle.md",
			"set state to spark|exploring|simmering|clarified|archived",
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
	entries, err := schema.Discover(root)
	if err != nil {
		return nil, err
	}
	out := make([]schema.Schema, 0, len(entries))
	for _, entry := range entries {
		out = append(out, entry.Schema)
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
			fix := "restore " + bind
			if bind == "index.md" {
				fix = "mycelium index"
			}
			if _, err := os.Stat(filepath.Join(root, bind)); err != nil {
				add(fmt.Sprintf("missing bound path %s", bind), "tier-binds", "program/contracts/conformance.md", fix)
			}
		}
	}
}

func checkTopLevel(root string, m manifest.Manifest, homes map[string]struct{}, reviewsAllowed bool, add func(string, string, string, string)) {
	entries, err := os.ReadDir(root)
	if err != nil {
		add(fmt.Sprintf("cannot read instance root: %v", err), "conformance", "program/contracts/conformance.md", "ensure the instance directory is readable")
		return
	}
	declared := map[string]struct{}{}
	for _, d := range m.Deviations {
		if strings.HasPrefix(d.Convention, "extra-top-level:") {
			key := strings.TrimPrefix(d.Convention, "extra-top-level:")
			declared[key] = struct{}{}
			declared[strings.TrimSuffix(key, "/")] = struct{}{}
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
		if name == "reviews" && reviewsAllowed {
			continue
		}
		if _, ok := declared[name]; ok {
			continue
		}
		if name == "reviews" && !reviewsAllowed {
			add(
				"extra top-level path reviews/ (council pack absent)",
				"extra-top-level",
				"program/contracts/conformance.md",
				"delete reviews/, restore program/packs/council/, or declare extra-top-level:reviews/",
			)
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

func checkIndex(root string, add func(string, string, string, string)) {
	b, err := os.ReadFile(filepath.Join(root, "index.md"))
	if err != nil {
		add(
			"index.md missing or unreadable",
			"index",
			"program/contracts/index.md",
			"mycelium index",
		)
		return
	}
	text := string(b)
	hasH1 := false
	for _, line := range strings.Split(text, "\n") {
		trim := strings.TrimRight(line, "\r")
		if strings.HasPrefix(trim, "# ") && !strings.HasPrefix(trim, "## ") {
			hasH1 = true
			break
		}
	}
	if !hasH1 {
		add(
			"index.md missing H1",
			"index",
			"program/contracts/index.md",
			"mycelium index",
		)
	}
	present := map[string]struct{}{}
	for _, m := range h2RE.FindAllStringSubmatch(text, -1) {
		present[m[1]] = struct{}{}
	}
	for _, sec := range []string{"State", "Artifacts", "Log tail", "Wake"} {
		if _, ok := present[sec]; !ok {
			add(
				fmt.Sprintf("index.md missing required H2 %q", sec),
				"index",
				"program/contracts/index.md",
				"mycelium index",
			)
		}
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
				if a.Home == "questions" && field == "agreement" {
					continue
				}
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

func checkQuestionSparring(root string, arts []artifactFile, add func(string, string, string, string)) {
	for _, a := range arts {
		if a.Home != "questions" {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(a.Rel)))
		if err != nil {
			continue
		}
		doc, err := metadata.Parse(b)
		if err != nil {
			continue
		}
		raw, ok := doc.Meta["agreement"].(string)
		if !ok {
			continue
		}
		agr, err := sparring.ParseAgreement(raw)
		if err != nil {
			add(
				fmt.Sprintf("agreement %q is not open|aligned|agree-to-disagree", raw),
				"question-front-matter",
				"program/templates/question.schema.toml",
				"set agreement to open, aligned, or agree-to-disagree",
			)
			continue
		}
		if agr != sparring.AgreeToDisagree {
			continue
		}
		for _, miss := range sparring.MissingHeadings(agr, doc.Body) {
			if !iffMissing(miss) {
				continue
			}
			what, fix := sparringTeach(a.IDStr, miss)
			add(what, "sparring", "program/contracts/sparring.md", fix)
		}
	}
}

func iffMissing(miss string) bool {
	if strings.HasPrefix(miss, "### ") {
		return true
	}
	return miss == "## Reasons" || miss == "## Crux"
}

func sparringTeach(id, miss string) (what, fix string) {
	if strings.HasPrefix(miss, "### ") {
		parent := miss
		if i := strings.Index(miss, " under "); i >= 0 {
			parent = miss[i+len(" under "):]
		}
		return id + " missing " + miss, "add ### Human and ### Agent under " + parent
	}
	return fmt.Sprintf("%s missing %s (required when agreement=agree-to-disagree)", id, miss),
		fmt.Sprintf("add %s with ### Human and ### Agent", miss)
}

func checkGlossary(root string, add func(string, string, string, string)) {
	b, err := os.ReadFile(filepath.Join(root, "CONTEXT.md"))
	if err != nil {
		return
	}
	content := string(b)
	if !sparring.HasGlossaryH1(content) {
		add(
			"CONTEXT.md missing H1 # Glossary",
			"glossary",
			"program/contracts/glossary.md",
			"add a line that is exactly # Glossary",
		)
	}
	for _, term := range sparring.MissingGlossaryDefinitions(content) {
		add(
			fmt.Sprintf("CONTEXT.md term %q missing ### Definition", term),
			"glossary",
			"program/contracts/glossary.md",
			fmt.Sprintf("add ### Definition under ## %s", term),
		)
	}
}

func checkDissent(root string, arts []artifactFile, add func(string, string, string, string)) {
	byNum := map[string]struct{}{}
	for _, a := range arts {
		if a.ID.NS != "OQ" && a.ID.NS != "ASM" {
			continue
		}
		byNum[fmt.Sprintf("%s-%d", a.ID.NS, a.ID.N)] = struct{}{}
	}
	for _, a := range arts {
		if a.Home != "decisions" {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(a.Rel)))
		if err != nil {
			continue
		}
		doc, err := metadata.Parse(b)
		if err != nil {
			continue
		}
		if !sparring.HasH2(doc.Body, "Dissent") {
			continue
		}
		sec := sparring.SectionBody(doc.Body, "Dissent")
		resolvable := false
		for _, tok := range sparring.DissentIDs(sec) {
			id, err := idpath.Parse(tok)
			if err != nil {
				continue
			}
			key := fmt.Sprintf("%s-%d", id.NS, id.N)
			if _, ok := byNum[key]; ok {
				resolvable = true
				break
			}
		}
		if resolvable {
			continue
		}
		add(
			a.IDStr+" ## Dissent has no resolvable OQ-### or ASM-###",
			"dissent",
			"program/contracts/sparring.md",
			"cite an existing OQ-### or ASM-### in ## Dissent, or remove the heading",
		)
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

func checkLog(root string, add func(string, string, string, string)) []byte {
	b, err := os.ReadFile(filepath.Join(root, "log.md"))
	if err != nil {
		add("log.md missing or unreadable", "log", "program/contracts/conformance.md", "restore log.md")
		return nil
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
				"use YYYY-MM-DD\\t(scaffold|new|tier|publish|check|state|wake|supersede)\\t<id>\\t…",
			)
		}
	}
	return b
}

func checkWakeBrief(root string, logBytes []byte, add func(string, string, string, string)) {
	if !logHasWake(logBytes) {
		return
	}
	path := filepath.Join(root, "briefs", "LATEST.md")
	b, err := os.ReadFile(path)
	if err != nil {
		add(
			"log has wake op but briefs/LATEST.md is missing",
			"wake",
			"program/contracts/wake.md",
			"mycelium wake   # rewrite the re-entry brief from simmering",
		)
		return
	}
	present := map[string]struct{}{}
	for _, m := range h2RE.FindAllStringSubmatch(string(b), -1) {
		present[m[1]] = struct{}{}
	}
	for _, sec := range wakebrief.RequiredH2s() {
		if _, ok := present[sec]; !ok {
			add(
				fmt.Sprintf("briefs/LATEST.md missing required H2 %q", sec),
				"wake",
				"program/contracts/wake.md",
				"mycelium wake   # rewrite the re-entry brief",
			)
		}
	}
}

func logHasWake(logBytes []byte) bool {
	if len(logBytes) == 0 {
		return false
	}
	for _, line := range logfmt.ParseableLines(logBytes) {
		parts := strings.SplitN(line, "\t", 4)
		if len(parts) >= 2 && parts[1] == "wake" {
			return true
		}
	}
	return false
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
			for _, value := range doc.Meta {
				if str, ok := value.(string); ok {
					text += "\n" + str
				}
			}
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

	for _, rel := range []string{"log.md", "README.md", "CONTEXT.md", "AGENTS.md", "index.md"} {
		scan(rel)
	}
	if entries, err := os.ReadDir(filepath.Join(root, "briefs")); err == nil {
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			scan(filepath.ToSlash(filepath.Join("briefs", e.Name())))
		}
	}
	for _, a := range arts {
		scan(a.Rel)
	}
}

// checkSupersedeIFF binds conformance item 23: bidirectional IFF + one-to-one.
func checkSupersedeIFF(root string, arts []artifactFile, add func(string, string, string, string)) {
	type artMeta struct {
		Rel          string
		IDStr        string
		NS           string
		Status       string
		SupersededBy string
		Supersedes   string
	}
	byID := map[string]artMeta{}
	inbound := map[string][]string{} // NEW-ID → OLD-IDs that name it in superseded_by

	for _, a := range arts {
		b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(a.Rel)))
		if err != nil {
			continue
		}
		doc, err := metadata.Parse(b)
		if err != nil {
			continue
		}
		am := artMeta{
			Rel:          a.Rel,
			IDStr:        a.IDStr,
			NS:           a.ID.NS,
			Status:       metaString(doc.Meta, "status"),
			SupersededBy: metaString(doc.Meta, "superseded_by"),
			Supersedes:   metaString(doc.Meta, "supersedes"),
		}
		byID[a.IDStr] = am
		if am.SupersededBy != "" {
			peer, err := idpath.Parse(am.SupersededBy)
			if err == nil {
				canon, ferr := idpath.FormatID(peer.NS, peer.N)
				if ferr == nil {
					inbound[canon] = append(inbound[canon], a.IDStr)
				}
			}
		}
	}

	for _, am := range byID {
		if am.Status == "Superseded" {
			if am.SupersededBy == "" {
				add(
					fmt.Sprintf("%s has status Superseded but missing superseded_by", am.IDStr),
					"supersede",
					"program/contracts/conformance.md",
					"mycelium supersede "+am.IDStr+" --by <NEW-ID>",
				)
				continue
			}
			peerID, peerNS, err := parseLinkID(am.SupersededBy)
			if err != nil {
				add(
					fmt.Sprintf("%s superseded_by %q is not a valid ID", am.IDStr, am.SupersededBy),
					"supersede",
					"program/contracts/conformance.md",
					"set superseded_by to a same-namespace artifact ID",
				)
				continue
			}
			if peerNS != am.NS {
				add(
					fmt.Sprintf("%s superseded_by %s crosses namespace (%s vs %s)", am.IDStr, peerID, am.NS, peerNS),
					"supersede",
					"program/contracts/conformance.md",
					"supersede requires the same namespace",
				)
				continue
			}
			peer, ok := byID[peerID]
			if !ok {
				add(
					fmt.Sprintf("%s superseded_by %s has no file", am.IDStr, peerID),
					"supersede",
					"program/contracts/conformance.md",
					"add the peer artifact or fix superseded_by",
				)
				continue
			}
			if peer.Supersedes != am.IDStr {
				add(
					fmt.Sprintf("%s superseded_by %s but peer supersedes is %q (want %s)", am.IDStr, peerID, peer.Supersedes, am.IDStr),
					"supersede",
					"program/contracts/conformance.md",
					"set "+peerID+" supersedes = \""+am.IDStr+"\" or re-run mycelium supersede",
				)
			}
		}

		if am.Supersedes != "" {
			peerID, peerNS, err := parseLinkID(am.Supersedes)
			if err != nil {
				add(
					fmt.Sprintf("%s supersedes %q is not a valid ID", am.IDStr, am.Supersedes),
					"supersede",
					"program/contracts/conformance.md",
					"set supersedes to a same-namespace artifact ID",
				)
				continue
			}
			if peerNS != am.NS {
				add(
					fmt.Sprintf("%s supersedes %s crosses namespace (%s vs %s)", am.IDStr, peerID, am.NS, peerNS),
					"supersede",
					"program/contracts/conformance.md",
					"supersede requires the same namespace",
				)
				continue
			}
			peer, ok := byID[peerID]
			if !ok {
				add(
					fmt.Sprintf("%s supersedes %s has no file", am.IDStr, peerID),
					"supersede",
					"program/contracts/conformance.md",
					"add the peer artifact or fix supersedes",
				)
				continue
			}
			if peer.Status != "Superseded" {
				add(
					fmt.Sprintf("%s supersedes %s but peer status is %q (want Superseded)", am.IDStr, peerID, peer.Status),
					"supersede",
					"program/contracts/conformance.md",
					"set "+peerID+" status = \"Superseded\" or re-run mycelium supersede",
				)
			}
			if peer.SupersededBy != am.IDStr {
				add(
					fmt.Sprintf("%s supersedes %s but peer superseded_by is %q (want %s)", am.IDStr, peerID, peer.SupersededBy, am.IDStr),
					"supersede",
					"program/contracts/conformance.md",
					"set "+peerID+" superseded_by = \""+am.IDStr+"\" or re-run mycelium supersede",
				)
			}
		}
	}

	for newID, olds := range inbound {
		if len(olds) > 1 {
			add(
				fmt.Sprintf("%s has multiple inbound superseded_by links (%s)", newID, strings.Join(olds, ", ")),
				"supersede",
				"program/contracts/conformance.md",
				"one-to-one this phase; keep a single superseded_by pointing at "+newID,
			)
		}
	}
}

func parseLinkID(tok string) (canon, ns string, err error) {
	id, err := idpath.Parse(tok)
	if err != nil {
		return "", "", err
	}
	canon, err = idpath.FormatID(id.NS, id.N)
	if err != nil {
		return "", "", err
	}
	return canon, id.NS, nil
}

func metaString(m map[string]any, key string) string {
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return fmt.Sprint(v)
	}
	return s
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
