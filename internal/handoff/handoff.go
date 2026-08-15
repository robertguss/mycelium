// Package handoff holds pure parsers and structure checks for handoff/PACKET.md.
// Slice 1: no CLI, no state lift, no journal, no log line.
package handoff

import (
	"fmt"
	"io/fs"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/robertguss/mycelium/internal/metadata"
)

// PacketID is the Architect-default packet id this phase (one packet per instance).
const PacketID = "HO-001"

// RequiredH2s are the nine PACKET.md headings in fixed order.
func RequiredH2s() []string {
	return []string{
		"Framing",
		"Locked decisions",
		"Glossary",
		"Open questions",
		"Evidence summary",
		"Implementation playbooks",
		"Implementation system",
		"Time budget",
		"Acceptance",
	}
}

var (
	requiredKeys = []string{"id", "date", "implementation_system", "time_budget"}

	dateRE       = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	timeBudgetRE = regexp.MustCompile(`^[0-9]+[mh]$`)
	h2RE         = regexp.MustCompile(`(?m)^## (.+)$`)
	mdLinkRE     = regexp.MustCompile(`\[([^\]]*)\]\(([^)]+)\)`)
	idTokenRE    = regexp.MustCompile(`\b([A-Z]+)-([0-9]+)\b`)
	implSystems  = map[string]struct{}{"pstack/poteto": {}, "manual": {}}
	copyHomeByNS = map[string]string{"DEC": "decisions", "OQ": "questions", "EVD": "evidence"}
	requiredDirs = []string{"decisions", "questions", "evidence", "playbooks", "acceptance"}
)

// Finding is one teaching reason (what / fix).
type Finding struct {
	What string
	Fix  string
}

// Check validates a handoff/ tree. fsys is rooted at handoff/ (PACKET.md at root).
// Returns nil when the packet structure passes.
func Check(fsys fs.FS) []Finding {
	var out []Finding
	add := func(what, fix string) {
		out = append(out, Finding{What: what, Fix: fix})
	}

	raw, err := fs.ReadFile(fsys, "PACKET.md")
	if err != nil {
		add("handoff/PACKET.md missing", "create handoff/PACKET.md from program/templates/handoff-packet.md")
		return out
	}

	doc, ferrs := validateFrontMatter(raw)
	out = append(out, ferrs...)
	if doc.Body == "" && len(ferrs) > 0 && hasMissingOpen(ferrs) {
		return out
	}
	out = append(out, validateH2s(doc.Body)...)

	for _, d := range requiredDirs {
		if !isDir(fsys, d) {
			add("handoff/"+d+"/ missing", "create handoff/"+d+"/")
		}
	}
	if !isFile(fsys, "glossary.md") {
		add("handoff/glossary.md missing", "create handoff/glossary.md")
	}

	// Self-contained: ID copies + path links inside PACKET.md and playbooks/.
	out = append(out, checkSelfContained(fsys, "PACKET.md", doc.Body)...)
	if isDir(fsys, "playbooks") {
		_ = fs.WalkDir(fsys, "playbooks", func(p string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil || d.IsDir() {
				return nil
			}
			if !strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
				return nil
			}
			b, err := fs.ReadFile(fsys, p)
			if err != nil {
				add("cannot read handoff/"+p, "fix handoff/"+p)
				return nil
			}
			out = append(out, checkSelfContained(fsys, p, string(b))...)
			return nil
		})
	}

	return out
}

// ValidatePacketBytes checks PACKET.md front matter and H2s only (no tree walk).
func ValidatePacketBytes(data []byte) []Finding {
	doc, ferrs := validateFrontMatter(data)
	out := append([]Finding{}, ferrs...)
	if hasMissingOpen(ferrs) {
		return out
	}
	return append(out, validateH2s(doc.Body)...)
}

func validateFrontMatter(data []byte) (metadata.Document, []Finding) {
	var out []Finding
	doc, err := metadata.Parse(data)
	if err != nil {
		out = append(out, Finding{
			What: "PACKET.md front matter invalid: " + err.Error(),
			Fix:  "use +++ TOML front matter with id, date, implementation_system, time_budget",
		})
		return metadata.Document{}, out
	}
	for _, k := range requiredKeys {
		v, ok := doc.Meta[k]
		if !ok {
			out = append(out, Finding{
				What: "PACKET.md missing required front matter key " + k,
				Fix:  "set " + k + " in PACKET.md front matter",
			})
			continue
		}
		s := metaString(v)
		if s == "" {
			out = append(out, Finding{
				What: "PACKET.md missing required front matter key " + k,
				Fix:  "set " + k + " in PACKET.md front matter",
			})
		}
	}
	if id := metaString(doc.Meta["id"]); id != "" && id != PacketID {
		out = append(out, Finding{
			What: fmt.Sprintf("PACKET.md id=%q must be %s", id, PacketID),
			Fix:  "set id = \"" + PacketID + "\"",
		})
	}
	if date := metaString(doc.Meta["date"]); date != "" && !dateRE.MatchString(date) {
		out = append(out, Finding{
			What: "PACKET.md date must be YYYY-MM-DD",
			Fix:  "set date to YYYY-MM-DD",
		})
	}
	if sys := metaString(doc.Meta["implementation_system"]); sys != "" {
		if _, ok := implSystems[sys]; !ok {
			out = append(out, Finding{
				What: fmt.Sprintf("PACKET.md implementation_system=%q is not pstack/poteto|manual", sys),
				Fix:  "set implementation_system to pstack/poteto or manual",
			})
		}
	}
	if tb := metaString(doc.Meta["time_budget"]); tb != "" && !timeBudgetRE.MatchString(tb) {
		out = append(out, Finding{
			What: fmt.Sprintf("PACKET.md time_budget=%q must match ^[0-9]+[mh]$", tb),
			Fix:  "set time_budget like 30m or 2h",
		})
	}
	return doc, out
}

func validateH2s(body string) []Finding {
	var out []Finding
	matches := h2RE.FindAllStringSubmatch(body, -1)
	seen := map[string]int{}
	for i, m := range matches {
		title := m[1]
		if _, ok := seen[title]; !ok {
			seen[title] = i
		}
	}
	req := RequiredH2s()
	lastIdx := -1
	for _, want := range req {
		idx, ok := seen[want]
		if !ok {
			out = append(out, Finding{
				What: fmt.Sprintf("PACKET.md missing required H2 %q", want),
				Fix:  "add ## " + want + " to PACKET.md",
			})
			continue
		}
		if idx <= lastIdx {
			out = append(out, Finding{
				What: fmt.Sprintf("PACKET.md H2 %q out of order", want),
				Fix:  "keep H2s in order: " + strings.Join(req, "; "),
			})
		}
		lastIdx = idx
	}
	return out
}

func checkSelfContained(fsys fs.FS, rel, body string) []Finding {
	var out []Finding
	baseDir := path.Dir(rel)
	if baseDir == "." {
		baseDir = ""
	}

	for _, m := range mdLinkRE.FindAllStringSubmatch(body, -1) {
		href := strings.TrimSpace(m[2])
		if href == "" || strings.HasPrefix(href, "#") {
			continue
		}
		if isExternalHref(href) {
			continue
		}
		// Strip optional title: path "title"
		if i := strings.IndexAny(href, " \t"); i >= 0 {
			href = href[:i]
		}
		resolved, ok := resolveInside(baseDir, href)
		if !ok {
			out = append(out, Finding{
				What: fmt.Sprintf("%s links outside handoff/: %s", rel, href),
				Fix:  "copy the target into handoff/ and link the in-packet path",
			})
			continue
		}
		if !exists(fsys, resolved) {
			out = append(out, Finding{
				What: fmt.Sprintf("%s link target missing: %s", rel, resolved),
				Fix:  "add handoff/" + resolved + " or fix the link",
			})
		}
	}

	seenID := map[string]struct{}{}
	for _, m := range idTokenRE.FindAllStringSubmatch(body, -1) {
		ns, num := m[1], m[2]
		if ns == "HO" {
			continue
		}
		home, ok := copyHomeByNS[ns]
		if !ok {
			continue
		}
		canon := formatID(ns, num)
		if _, dup := seenID[canon]; dup {
			continue
		}
		seenID[canon] = struct{}{}
		if !hasIDCopy(fsys, home, canon) {
			out = append(out, Finding{
				What: fmt.Sprintf("%s cites %s but no copy under handoff/%s/", rel, canon, home),
				Fix:  "copy " + canon + " into handoff/" + home + "/",
			})
		}
	}
	return out
}

func resolveInside(baseDir, href string) (string, bool) {
	href = path.Clean(href)
	if path.IsAbs(href) || strings.HasPrefix(href, "/") {
		return "", false
	}
	var joined string
	if baseDir == "" {
		joined = href
	} else {
		joined = path.Clean(path.Join(baseDir, href))
	}
	if joined == ".." || strings.HasPrefix(joined, "../") {
		return "", false
	}
	return joined, true
}

func isExternalHref(href string) bool {
	lower := strings.ToLower(href)
	return strings.HasPrefix(lower, "http://") ||
		strings.HasPrefix(lower, "https://") ||
		strings.HasPrefix(lower, "mailto:")
}

func hasIDCopy(fsys fs.FS, home, canon string) bool {
	entries, err := fs.ReadDir(fsys, home)
	if err != nil {
		return false
	}
	prefix := canon + "-"
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if name == canon+".md" || strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}

func formatID(ns, num string) string {
	n, err := strconv.Atoi(num)
	if err != nil || n < 0 {
		return ns + "-" + num
	}
	return fmt.Sprintf("%s-%03d", ns, n)
}

func isDir(fsys fs.FS, name string) bool {
	info, err := fs.Stat(fsys, name)
	return err == nil && info.IsDir()
}

func isFile(fsys fs.FS, name string) bool {
	info, err := fs.Stat(fsys, name)
	return err == nil && !info.IsDir()
}

func exists(fsys fs.FS, name string) bool {
	_, err := fs.Stat(fsys, name)
	return err == nil
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

func hasMissingOpen(fs []Finding) bool {
	for _, f := range fs {
		if strings.Contains(f.What, "missing opening") || strings.Contains(f.What, "front matter invalid") {
			return true
		}
	}
	return false
}
