// Package idpath maps artifact IDs to paths and back. No filesystem, no clock.
package idpath

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ID is a parsed artifact identifier (numeric; padding is a path concern).
type ID struct {
	NS string
	N  int
}

// Type describes one registered generator type.
type Type struct {
	Key         string
	NS          string
	Home        string
	Digits      int
	StageScoped bool
}

var (
	ErrEmpty   = errors.New("idpath: empty")
	ErrUnknown = errors.New("idpath: unknown namespace")
	ErrFormat  = errors.New("idpath: bad format")
)

// Registered types (PHASE-01). Order matches the brief catalog.
var types = []Type{
	{Key: "decision", NS: "DEC", Home: "decisions", Digits: 3, StageScoped: false},
	{Key: "assumption", NS: "ASM", Home: "assumptions", Digits: 3, StageScoped: false},
	{Key: "evidence", NS: "EVD", Home: "evidence", Digits: 3, StageScoped: false},
	{Key: "spike", NS: "SPK", Home: "spikes", Digits: 3, StageScoped: false},
	{Key: "finding", NS: "FND", Home: "findings", Digits: 3, StageScoped: true},
	{Key: "recommendation", NS: "REC", Home: "recommendations", Digits: 3, StageScoped: true},
	{Key: "requirement", NS: "REQ", Home: "requirements", Digits: 3, StageScoped: true},
	{Key: "question", NS: "OQ", Home: "questions", Digits: 3, StageScoped: false},
	{Key: "risk", NS: "RSK", Home: "risks", Digits: 3, StageScoped: false},
	{Key: "phase", NS: "PHASE", Home: "phases", Digits: 2, StageScoped: false},
	{Key: "milestone", NS: "MS", Home: "milestones", Digits: 3, StageScoped: false},
}

var byNS map[string]Type
var byHome map[string]Type

func init() {
	byNS = make(map[string]Type, len(types))
	byHome = make(map[string]Type, len(types))
	for _, t := range types {
		byNS[t.NS] = t
		byHome[t.Home] = t
	}
}

// Types returns a copy of the registered type catalog.
func Types() []Type {
	out := make([]Type, len(types))
	copy(out, types)
	return out
}

// LookupNS returns the registered type for a namespace.
func LookupNS(ns string) (Type, error) {
	t, ok := byNS[ns]
	if !ok {
		return Type{}, fmt.Errorf("%w: %s", ErrUnknown, ns)
	}
	return t, nil
}

var idToken = regexp.MustCompile(`^([A-Z]+)-([0-9]+)$`)

// Parse parses an ID token. DEC-001 and DEC-1 are both valid tokens.
func Parse(s string) (ID, error) {
	if s == "" {
		return ID{}, ErrEmpty
	}
	m := idToken.FindStringSubmatch(s)
	if m == nil {
		return ID{}, fmt.Errorf("%w: %q", ErrFormat, s)
	}
	ns := m[1]
	if _, err := LookupNS(ns); err != nil {
		return ID{}, err
	}
	n, err := strconv.Atoi(m[2])
	if err != nil {
		return ID{}, fmt.Errorf("%w: %q", ErrFormat, s)
	}
	return ID{NS: ns, N: n}, nil
}

// FormatID zero-pads n to the namespace digit width (no slug).
// Refuses n < 0 or n that does not fit Digits (n >= 10^Digits).
func FormatID(ns string, n int) (string, error) {
	t, err := LookupNS(ns)
	if err != nil {
		return "", err
	}
	if n < 0 {
		return "", fmt.Errorf("%w: negative", ErrFormat)
	}
	limit := 1
	for i := 0; i < t.Digits; i++ {
		limit *= 10
	}
	if n >= limit {
		return "", fmt.Errorf("%w: %s n=%d exceeds %d digits", ErrFormat, ns, n, t.Digits)
	}
	return fmt.Sprintf("%s-%0*d", ns, t.Digits, n), nil
}

// Format returns home/NS-DIGITS-slug.md with zero-padded digits.
func Format(ns string, n int, slug string) (string, error) {
	id, err := FormatID(ns, n)
	if err != nil {
		return "", err
	}
	t, err := LookupNS(ns)
	if err != nil {
		return "", err
	}
	if slug == "" {
		return "", fmt.Errorf("%w: empty slug", ErrFormat)
	}
	return t.Home + "/" + id + "-" + slug + ".md", nil
}

// PathFor is Format using a parsed ID (zero-pads).
func PathFor(id ID, slug string) (string, error) {
	return Format(id.NS, id.N, slug)
}

// ParsePath parses home/NS-DIGITS-slug.md. Digits must match the type width.
func ParsePath(path string) (ID, string, error) {
	if path == "" {
		return ID{}, "", ErrEmpty
	}
	path = strings.ReplaceAll(path, "\\", "/")
	slash := strings.IndexByte(path, '/')
	if slash <= 0 || slash == len(path)-1 {
		return ID{}, "", fmt.Errorf("%w: %q", ErrFormat, path)
	}
	home := path[:slash]
	rest := path[slash+1:]
	t, ok := byHome[home]
	if !ok {
		return ID{}, "", fmt.Errorf("%w: home %q", ErrUnknown, home)
	}
	if !strings.HasSuffix(rest, ".md") {
		return ID{}, "", fmt.Errorf("%w: %q", ErrFormat, path)
	}
	rest = strings.TrimSuffix(rest, ".md")
	prefix := t.NS + "-"
	if !strings.HasPrefix(rest, prefix) {
		return ID{}, "", fmt.Errorf("%w: %q", ErrFormat, path)
	}
	rest = rest[len(prefix):]
	// Exact digit run of length Digits, then '-', then slug.
	if len(rest) < t.Digits+1 || rest[t.Digits] != '-' {
		return ID{}, "", fmt.Errorf("%w: %q", ErrFormat, path)
	}
	digitPart := rest[:t.Digits]
	slugPart := rest[t.Digits+1:]
	for _, c := range digitPart {
		if c < '0' || c > '9' {
			return ID{}, "", fmt.Errorf("%w: %q", ErrFormat, path)
		}
	}
	// Reject wrong widths that still look numeric (e.g. leading zeros beyond width already excluded by fixed width).
	n, err := strconv.Atoi(digitPart)
	if err != nil {
		return ID{}, "", fmt.Errorf("%w: %q", ErrFormat, path)
	}
	if slugPart == "" {
		return ID{}, "", fmt.Errorf("%w: empty slug", ErrFormat)
	}
	return ID{NS: t.NS, N: n}, slugPart, nil
}
