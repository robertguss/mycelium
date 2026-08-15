// Package wakebrief collects wake citations and writes the re-entry brief.
package wakebrief

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/idpath"
	"github.com/robertguss/mycelium/internal/logfmt"
	"github.com/robertguss/mycelium/internal/manifest"
	"github.com/robertguss/mycelium/internal/metadata"
	"github.com/robertguss/mycelium/internal/revisit"
)

// EvidenceCite is one due/overdue EVD citation.
type EvidenceCite struct {
	ID   string
	Date string
}

// AssumptionCite is one qualifying ASM citation.
type AssumptionCite struct {
	ID     string
	Status string
	Date   string // empty if included by status only
}

// Brief is the deterministic citation set for one wake.
type Brief struct {
	WakeDate     string
	SimmerDate   string
	SimmerNote   string
	Revisit      string
	SimmerLine   string
	LogSince     []string
	Evidence     []EvidenceCite
	Assumptions  []AssumptionCite
	CreatedDate  string
}

// Collect gathers citations from the instance. Caller must preflight simmering.
func Collect(root string, m manifest.Manifest, logBytes []byte, now time.Time) (Brief, error) {
	wakeDate := clock.Date(now)
	b := Brief{
		WakeDate:    wakeDate,
		Revisit:     m.Revisit,
		CreatedDate: m.CreatedDate,
	}

	lines := logfmt.ParseableLines(logBytes)
	simmerIdx := -1
	for i := len(lines) - 1; i >= 0; i-- {
		date, op, _, note := splitLog(lines[i])
		if op == "state" && strings.HasPrefix(note, "simmering") {
			simmerIdx = i
			b.SimmerLine = lines[i]
			b.SimmerDate = date
			b.SimmerNote = note
			break
		}
	}

	if simmerIdx >= 0 {
		b.LogSince = append([]string(nil), lines[simmerIdx:]...)
	} else {
		for _, line := range lines {
			date, _, _, _ := splitLog(line)
			if date >= m.CreatedDate {
				b.LogSince = append(b.LogSince, line)
			}
		}
		if b.SimmerDate == "" {
			b.SimmerDate = m.CreatedDate
		}
		if b.SimmerNote == "" {
			b.SimmerNote = "simmering revisit=" + m.Revisit
		}
		if b.SimmerLine == "" {
			b.SimmerLine = logfmt.Line(b.SimmerDate, "state", "-", b.SimmerNote)
		}
	}

	evd, err := scanEvidence(root, now)
	if err != nil {
		return Brief{}, err
	}
	b.Evidence = evd

	asm, err := scanAssumptions(root, now)
	if err != nil {
		return Brief{}, err
	}
	b.Assumptions = asm
	return b, nil
}

// Render returns markdown bytes for the wake brief (five required H2s).
func Render(b Brief) []byte {
	var s strings.Builder
	fmt.Fprintf(&s, "# Wake — %s\n\n", b.WakeDate)

	fmt.Fprintf(&s, "## Parked\n\n")
	fmt.Fprintf(&s, "Simmered on %s with revisit %s.\n\n", b.SimmerDate, b.Revisit)

	fmt.Fprintf(&s, "## Log since simmer\n\n")
	if len(b.LogSince) == 0 {
		s.WriteString(b.SimmerLine)
		s.WriteByte('\n')
	} else {
		for _, line := range b.LogSince {
			s.WriteString(line)
			s.WriteByte('\n')
		}
	}
	s.WriteByte('\n')

	fmt.Fprintf(&s, "## Evidence triggers\n\n")
	if len(b.Evidence) == 0 {
		s.WriteString("none\n")
	} else {
		for _, e := range b.Evidence {
			fmt.Fprintf(&s, "%s (revalidation %s)\n", e.ID, e.Date)
		}
	}
	s.WriteByte('\n')

	fmt.Fprintf(&s, "## Assumptions\n\n")
	if len(b.Assumptions) == 0 {
		s.WriteString("none\n")
	} else {
		for _, a := range b.Assumptions {
			if a.Date != "" {
				fmt.Fprintf(&s, "%s (%s; revisit %s)\n", a.ID, a.Status, a.Date)
			} else {
				fmt.Fprintf(&s, "%s (%s)\n", a.ID, a.Status)
			}
		}
	}
	s.WriteByte('\n')

	fmt.Fprintf(&s, "## Suggested next\n\n")
	s.WriteString("<!-- fill -->\n")
	return []byte(s.String())
}

// DatedPath is the relative path for today's wake brief.
func DatedPath(now time.Time) string {
	return filepath.ToSlash(filepath.Join("briefs", "WAKE-"+clock.Date(now)+".md"))
}

// RequiredH2s are the five structure-only wake brief headings.
func RequiredH2s() []string {
	return []string{
		"Parked",
		"Log since simmer",
		"Evidence triggers",
		"Assumptions",
		"Suggested next",
	}
}

func scanEvidence(root string, now time.Time) ([]EvidenceCite, error) {
	dir := filepath.Join(root, "evidence")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []EvidenceCite
	for _, e := range entries {
		if e.IsDir() || e.Name() == "README.md" || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		rel := filepath.ToSlash(filepath.Join("evidence", e.Name()))
		id, _, err := idpath.ParsePath(rel)
		if err != nil {
			continue
		}
		idStr, err := idpath.FormatID(id.NS, id.N)
		if err != nil {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		doc, err := metadata.Parse(raw)
		if err != nil {
			continue
		}
		body := sectionBody(doc.Body, "Revalidation Trigger")
		trig, ok := revisit.ExtractTriggerDate(body)
		if !ok {
			continue
		}
		if !revisit.Due(revisit.Date, trig, now) {
			continue
		}
		out = append(out, EvidenceCite{ID: idStr, Date: clock.Date(trig)})
	}
	return out, nil
}

func scanAssumptions(root string, now time.Time) ([]AssumptionCite, error) {
	dir := filepath.Join(root, "assumptions")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []AssumptionCite
	for _, e := range entries {
		if e.IsDir() || e.Name() == "README.md" || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		rel := filepath.ToSlash(filepath.Join("assumptions", e.Name()))
		id, _, err := idpath.ParsePath(rel)
		if err != nil {
			continue
		}
		idStr, err := idpath.FormatID(id.NS, id.N)
		if err != nil {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		doc, err := metadata.Parse(raw)
		if err != nil {
			continue
		}
		status, _ := doc.Meta["status"].(string)
		body := sectionBody(doc.Body, "Revisit Triggers")
		trig, hasDate := revisit.ExtractTriggerDate(body)
		dateDue := hasDate && revisit.Due(revisit.Date, trig, now)
		statusQual := status == "Open" || status == "Held"
		if !statusQual && !dateDue {
			continue
		}
		c := AssumptionCite{ID: idStr, Status: status}
		if hasDate && dateDue {
			c.Date = clock.Date(trig)
		} else if hasDate && statusQual {
			c.Date = clock.Date(trig)
		}
		out = append(out, c)
	}
	return out, nil
}

func sectionBody(body, heading string) string {
	want := "## " + heading
	lines := strings.Split(body, "\n")
	start := -1
	for i, line := range lines {
		trim := strings.TrimRight(line, "\r")
		if trim == want {
			start = i + 1
			break
		}
	}
	if start < 0 {
		return ""
	}
	var b strings.Builder
	for i := start; i < len(lines); i++ {
		trim := strings.TrimRight(lines[i], "\r")
		if strings.HasPrefix(trim, "## ") {
			break
		}
		b.WriteString(lines[i])
		if i+1 < len(lines) {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func splitLog(line string) (date, op, id, note string) {
	parts := strings.SplitN(line, "\t", 4)
	if len(parts) < 4 {
		return "", "", "", ""
	}
	return parts[0], parts[1], parts[2], parts[3]
}
