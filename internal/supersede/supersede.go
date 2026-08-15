// Package supersede holds pure parsers and cross-link rules for mycelium supersede.
// Slice 1: no filesystem CLI, no log line, no journal.
package supersede

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/robertguss/mycelium/internal/idpath"
	"github.com/robertguss/mycelium/internal/metadata"
)

// Idea-state tokens refuse as artifact IDs (§4).
var ideaStates = map[string]struct{}{
	"spark": {}, "exploring": {}, "simmering": {},
	"clarified": {}, "handed-off": {}, "archived": {},
}

// Eligible namespaces after §10 schema deltas (status enum includes Superseded).
var eligibleNS = map[string]struct{}{
	"DEC": {}, "ASM": {}, "EVD": {}, "SPK": {},
}

// Refusal is a teaching-error what/fix pair from §4. Convention/contract bind later.
type Refusal struct {
	What string
	Fix  string
}

// IsIdeaState reports whether tok is an idea-lifecycle state token.
func IsIdeaState(tok string) bool {
	_, ok := ideaStates[tok]
	return ok
}

// Eligible reports whether namespace may be superseded this phase.
func Eligible(ns string) bool {
	_, ok := eligibleNS[ns]
	return ok
}

// ParseID parses an artifact ID token. Idea-state tokens refuse.
// Returns a zero-padded canonical ID via idpath.FormatID.
func ParseID(tok string) (canonical string, ns string, refuse *Refusal) {
	if IsIdeaState(tok) {
		return "", "", &Refusal{
			What: tok + " is an idea state, not an artifact",
			Fix:  "mycelium state <target> for lifecycle; mycelium supersede for artifacts",
		}
	}
	id, err := idpath.Parse(tok)
	if err != nil {
		return "", "", &Refusal{
			What: fmt.Sprintf("no artifact %s", tok),
			Fix:  `mycelium new <type> "…" then retry`,
		}
	}
	canon, err := idpath.FormatID(id.NS, id.N)
	if err != nil {
		return "", "", &Refusal{
			What: fmt.Sprintf("no artifact %s", tok),
			Fix:  `mycelium new <type> "…" then retry`,
		}
	}
	return canon, id.NS, nil
}

// CheckPair validates OLD/NEW tokens and optional front-matter maps for a one-to-one supersede.
// oldMeta/newMeta may be nil when only ID-level rules are tested.
func CheckPair(oldTok, newTok string, oldMeta, newMeta map[string]any) *Refusal {
	if oldTok == "" || newTok == "" {
		return &Refusal{
			What: "supersede requires <OLD-ID> --by <NEW-ID>",
			Fix:  "mycelium supersede DEC-001 --by DEC-002",
		}
	}
	oldID, oldNS, r := ParseID(oldTok)
	if r != nil {
		return r
	}
	newID, newNS, r := ParseID(newTok)
	if r != nil {
		return r
	}
	if oldID == newID {
		return &Refusal{
			What: "cannot supersede an ID with itself",
			Fix:  "pass two different IDs",
		}
	}
	if oldNS != newNS {
		return &Refusal{
			What: fmt.Sprintf("supersede requires the same namespace (got %s vs %s)", oldNS, newNS),
			Fix:  "pick two IDs in one namespace",
		}
	}
	if !Eligible(oldNS) {
		fix := "choose DEC, ASM, EVD, or SPK IDs"
		if oldNS == "OQ" {
			fix = "open a new question; do not supersede an OQ"
		}
		return &Refusal{
			What: fmt.Sprintf("type %s is not supersedable", oldNS),
			Fix:  fix,
		}
	}
	if oldMeta != nil {
		status := metaString(oldMeta, "status")
		if status == "Superseded" {
			existing := metaString(oldMeta, "superseded_by")
			if existing == "" {
				existing = "(unknown)"
			}
			return &Refusal{
				What: fmt.Sprintf("%s is already Superseded by %s", oldID, existing),
				Fix:  fmt.Sprintf("supersede the current record (%s) --by <newer>", existing),
			}
		}
	}
	if newMeta != nil {
		existing := metaString(newMeta, "supersedes")
		if existing != "" {
			return &Refusal{
				What: fmt.Sprintf("%s already supersedes %s", newID, existing),
				Fix:  "one-to-one this phase; pick a different NEW",
			}
		}
	}
	return nil
}

// ApplyOld sets status=Superseded and superseded_by=<newID> on OLD artifact bytes.
func ApplyOld(data []byte, newID string) ([]byte, error) {
	return setFrontMatter(data, map[string]string{
		"status":        "Superseded",
		"superseded_by": newID,
	})
}

// ApplyNew sets supersedes=<oldID> on NEW artifact bytes. Status is unchanged.
func ApplyNew(data []byte, oldID string) ([]byte, error) {
	return setFrontMatter(data, map[string]string{
		"supersedes": oldID,
	})
}

// ApplyPair mutates OLD then NEW when CheckPair passes.
func ApplyPair(oldData, newData []byte, oldTok, newTok string) (oldOut, newOut []byte, refuse *Refusal, err error) {
	oldDoc, err := metadata.Parse(oldData)
	if err != nil {
		return nil, nil, nil, err
	}
	newDoc, err := metadata.Parse(newData)
	if err != nil {
		return nil, nil, nil, err
	}
	if r := CheckPair(oldTok, newTok, oldDoc.Meta, newDoc.Meta); r != nil {
		return nil, nil, r, nil
	}
	oldID, _, _ := ParseID(oldTok)
	newID, _, _ := ParseID(newTok)
	oldOut, err = ApplyOld(oldData, newID)
	if err != nil {
		return nil, nil, nil, err
	}
	newOut, err = ApplyNew(newData, oldID)
	if err != nil {
		return nil, nil, nil, err
	}
	return oldOut, newOut, nil, nil
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

// setFrontMatter surgically updates +++ TOML keys, preserving order and body.
func setFrontMatter(data []byte, sets map[string]string) ([]byte, error) {
	if _, err := metadata.Parse(data); err != nil {
		return nil, err
	}
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})
	lines := splitKeepEnds(data)
	if len(lines) == 0 || lineContent(lines[0]) != "+++" {
		return nil, metadata.ErrMissingOpen
	}
	closeIdx := -1
	for i := 1; i < len(lines); i++ {
		if lineContent(lines[i]) == "+++" {
			closeIdx = i
			break
		}
	}
	if closeIdx < 0 {
		return nil, metadata.ErrMissingClose
	}

	seen := map[string]bool{}
	var mid [][]byte
	for i := 1; i < closeIdx; i++ {
		raw := lines[i]
		key, ok := tomlKey(lineContent(raw))
		if ok {
			if val, set := sets[key]; set {
				eol := eolOf(raw)
				mid = append(mid, []byte(fmt.Sprintf("%s = %q%s", key, val, eol)))
				seen[key] = true
				continue
			}
		}
		mid = append(mid, raw)
	}
	for _, key := range orderedKeys(sets, seen) {
		mid = append(mid, []byte(fmt.Sprintf("%s = %q\n", key, sets[key])))
	}

	var out bytes.Buffer
	out.Write(lines[0])
	for _, l := range mid {
		out.Write(l)
	}
	for i := closeIdx; i < len(lines); i++ {
		out.Write(lines[i])
	}
	return out.Bytes(), nil
}

func orderedKeys(sets map[string]string, seen map[string]bool) []string {
	// Prefer status, superseded_by, supersedes when appending new keys.
	prefer := []string{"status", "superseded_by", "supersedes"}
	var out []string
	for _, k := range prefer {
		if _, ok := sets[k]; ok && !seen[k] {
			out = append(out, k)
		}
	}
	for k := range sets {
		if seen[k] {
			continue
		}
		already := false
		for _, p := range out {
			if p == k {
				already = true
				break
			}
		}
		if !already {
			out = append(out, k)
		}
	}
	return out
}

func tomlKey(line string) (string, bool) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return "", false
	}
	eq := strings.IndexByte(line, '=')
	if eq <= 0 {
		return "", false
	}
	key := strings.TrimSpace(line[:eq])
	if key == "" {
		return "", false
	}
	return key, true
}

func splitKeepEnds(data []byte) [][]byte {
	if len(data) == 0 {
		return nil
	}
	var lines [][]byte
	start := 0
	for i := 0; i < len(data); i++ {
		if data[i] == '\n' {
			lines = append(lines, data[start:i+1])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}

func eolOf(line []byte) string {
	if bytes.HasSuffix(line, []byte("\r\n")) {
		return "\r\n"
	}
	if bytes.HasSuffix(line, []byte("\n")) {
		return "\n"
	}
	return "\n"
}

func lineContent(line []byte) string {
	line = bytes.TrimSuffix(line, []byte("\n"))
	line = bytes.TrimSuffix(line, []byte("\r"))
	return string(line)
}
