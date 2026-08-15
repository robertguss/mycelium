// Package manifest parses and validates mycelium.toml.
package manifest

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/pelletier/go-toml/v2"

	"github.com/robertguss/mycelium/internal/slug"
)

var (
	ErrUnknownKey = errors.New("manifest: unknown key")
	ErrRequired   = errors.New("manifest: missing required field")
	ErrInvalid    = errors.New("manifest: invalid value")
	ErrRange      = errors.New("manifest: bad range")
)

var allowedStates = map[string]struct{}{
	"spark": {}, "exploring": {}, "simmering": {},
	"clarified": {}, "handed-off": {}, "archived": {},
}

var allowedTiers = map[string]struct{}{
	"focused": {}, "standard": {}, "high-assurance": {},
}

var identifierKeys = map[string]string{
	"findings":        "FND",
	"recommendations": "REC",
	"requirements":    "REQ",
}

var rangeRE = regexp.MustCompile(`^([A-Z]+)-([0-9]+)\.\.([A-Z]+)-([0-9]+)$`)
var dateRE = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

var requiredTop = []string{
	"schema_version", "idea_name", "slug", "state", "tier",
	"methodology_version", "generated_by_cli_version",
	"created_date", "updated_date", "revisit", "github_repo",
}

var allowedTop = map[string]struct{}{
	"schema_version": {}, "idea_name": {}, "slug": {}, "state": {}, "tier": {},
	"methodology_version": {}, "generated_by_cli_version": {},
	"created_date": {}, "updated_date": {}, "revisit": {}, "github_repo": {},
	"identifiers": {}, "deviations": {},
}

// Range is a stage-scoped ID span for one namespace.
type Range struct {
	NS    string
	Start int
	End   int
	Raw   string
}

// Deviation is one [[deviations]] row.
type Deviation struct {
	Convention string
	Reason     string
}

// Manifest is a validated mycelium.toml.
type Manifest struct {
	SchemaVersion         int
	IdeaName              string
	Slug                  string
	State                 string
	Tier                  string
	MethodologyVersion    string
	GeneratedByCLIVersion string
	CreatedDate           string
	UpdatedDate           string
	Revisit               string
	GithubRepo            string
	Identifiers           map[string]Range
	Deviations            []Deviation
}

// Parse decodes and validates mycelium.toml bytes.
func Parse(data []byte) (Manifest, error) {
	var raw map[string]any
	if err := toml.Unmarshal(data, &raw); err != nil {
		return Manifest{}, fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	for k := range raw {
		if _, ok := allowedTop[k]; !ok {
			return Manifest{}, fmt.Errorf("%w: %s", ErrUnknownKey, k)
		}
	}
	for _, k := range requiredTop {
		if _, ok := raw[k]; !ok {
			return Manifest{}, fmt.Errorf("%w: %s", ErrRequired, k)
		}
	}

	m := Manifest{Identifiers: map[string]Range{}}

	sv, err := asInt(raw["schema_version"])
	if err != nil {
		return Manifest{}, fmt.Errorf("%w: schema_version", ErrInvalid)
	}
	if sv != 1 {
		return Manifest{}, fmt.Errorf("%w: schema_version must be 1", ErrInvalid)
	}
	m.SchemaVersion = sv

	m.IdeaName, err = asString(raw["idea_name"])
	if err != nil || m.IdeaName == "" {
		return Manifest{}, fmt.Errorf("%w: idea_name", ErrRequired)
	}
	m.Slug, err = asString(raw["slug"])
	if err != nil || m.Slug == "" {
		return Manifest{}, fmt.Errorf("%w: slug", ErrRequired)
	}
	wantSlug, err := slug.Slugify(m.IdeaName)
	if err != nil {
		return Manifest{}, fmt.Errorf("%w: idea_name slugify: %v", ErrInvalid, err)
	}
	if m.Slug != wantSlug {
		return Manifest{}, fmt.Errorf("%w: slug must equal slugify(idea_name) (%q)", ErrInvalid, wantSlug)
	}

	m.State, err = asString(raw["state"])
	if err != nil {
		return Manifest{}, fmt.Errorf("%w: state", ErrInvalid)
	}
	if _, ok := allowedStates[m.State]; !ok {
		return Manifest{}, fmt.Errorf("%w: state %q", ErrInvalid, m.State)
	}
	m.Tier, err = asString(raw["tier"])
	if err != nil {
		return Manifest{}, fmt.Errorf("%w: tier", ErrInvalid)
	}
	if _, ok := allowedTiers[m.Tier]; !ok {
		return Manifest{}, fmt.Errorf("%w: tier %q", ErrInvalid, m.Tier)
	}

	m.MethodologyVersion, err = asString(raw["methodology_version"])
	if err != nil || m.MethodologyVersion == "" {
		return Manifest{}, fmt.Errorf("%w: methodology_version", ErrRequired)
	}
	m.GeneratedByCLIVersion, err = asString(raw["generated_by_cli_version"])
	if err != nil || m.GeneratedByCLIVersion == "" {
		return Manifest{}, fmt.Errorf("%w: generated_by_cli_version", ErrRequired)
	}
	m.CreatedDate, err = asString(raw["created_date"])
	if err != nil || m.CreatedDate == "" {
		return Manifest{}, fmt.Errorf("%w: created_date", ErrRequired)
	}
	if !dateRE.MatchString(m.CreatedDate) {
		return Manifest{}, fmt.Errorf("%w: created_date must be YYYY-MM-DD", ErrInvalid)
	}
	m.UpdatedDate, err = asString(raw["updated_date"])
	if err != nil || m.UpdatedDate == "" {
		return Manifest{}, fmt.Errorf("%w: updated_date", ErrRequired)
	}
	if !dateRE.MatchString(m.UpdatedDate) {
		return Manifest{}, fmt.Errorf("%w: updated_date must be YYYY-MM-DD", ErrInvalid)
	}
	m.Revisit, err = asString(raw["revisit"])
	if err != nil {
		return Manifest{}, fmt.Errorf("%w: revisit", ErrInvalid)
	}
	m.GithubRepo, err = asString(raw["github_repo"])
	if err != nil {
		return Manifest{}, fmt.Errorf("%w: github_repo", ErrInvalid)
	}
	if m.State == "simmering" && m.Revisit == "" {
		return Manifest{}, fmt.Errorf("%w: revisit required when state=simmering", ErrInvalid)
	}

	if idRaw, ok := raw["identifiers"]; ok {
		idMap, ok := idRaw.(map[string]any)
		if !ok {
			return Manifest{}, fmt.Errorf("%w: identifiers", ErrInvalid)
		}
		for k, v := range idMap {
			ns, ok := identifierKeys[k]
			if !ok {
				return Manifest{}, fmt.Errorf("%w: identifiers.%s", ErrUnknownKey, k)
			}
			s, err := asString(v)
			if err != nil {
				return Manifest{}, fmt.Errorf("%w: identifiers.%s", ErrInvalid, k)
			}
			rg, err := ParseRange(s, ns)
			if err != nil {
				return Manifest{}, err
			}
			m.Identifiers[k] = rg
		}
	}

	if devRaw, ok := raw["deviations"]; ok {
		arr, ok := devRaw.([]any)
		if !ok {
			return Manifest{}, fmt.Errorf("%w: deviations", ErrInvalid)
		}
		for i, item := range arr {
			row, ok := item.(map[string]any)
			if !ok {
				return Manifest{}, fmt.Errorf("%w: deviations[%d]", ErrInvalid, i)
			}
			for k := range row {
				if k != "convention" && k != "reason" {
					return Manifest{}, fmt.Errorf("%w: deviations.%s", ErrUnknownKey, k)
				}
			}
			conv, err := asString(row["convention"])
			if err != nil || conv == "" {
				return Manifest{}, fmt.Errorf("%w: deviations[%d].convention", ErrRequired, i)
			}
			reason, err := asString(row["reason"])
			if err != nil || reason == "" {
				return Manifest{}, fmt.Errorf("%w: deviations[%d].reason", ErrRequired, i)
			}
			m.Deviations = append(m.Deviations, Deviation{Convention: conv, Reason: reason})
		}
	}

	return m, nil
}

// ParseRange parses NS-XXX..NS-YYY and checks NS + start<=end.
func ParseRange(raw, expectNS string) (Range, error) {
	m := rangeRE.FindStringSubmatch(raw)
	if m == nil {
		return Range{}, fmt.Errorf("%w: %q", ErrRange, raw)
	}
	ns1, ns2 := m[1], m[3]
	if ns1 != expectNS || ns2 != expectNS {
		return Range{}, fmt.Errorf("%w: NS mismatch want %s got %s..%s", ErrRange, expectNS, ns1, ns2)
	}
	start, err1 := strconv.Atoi(m[2])
	end, err2 := strconv.Atoi(m[4])
	if err1 != nil || err2 != nil {
		return Range{}, fmt.Errorf("%w: %q", ErrRange, raw)
	}
	if start > end {
		return Range{}, fmt.Errorf("%w: start>end in %q", ErrRange, raw)
	}
	return Range{NS: expectNS, Start: start, End: end, Raw: raw}, nil
}

// Contains reports whether n is inside the range (inclusive).
func (r Range) Contains(n int) bool {
	return n >= r.Start && n <= r.End
}

// InRange checks membership for a stage-scoped identifier key.
// Missing declaration → ErrRange.
func (m Manifest) InRange(key string, n int) error {
	rg, ok := m.Identifiers[key]
	if !ok {
		return fmt.Errorf("%w: missing declaration for %s", ErrRange, key)
	}
	if !rg.Contains(n) {
		return fmt.Errorf("%w: %d outside %s", ErrRange, n, rg.Raw)
	}
	return nil
}

func asString(v any) (string, error) {
	s, ok := v.(string)
	if !ok {
		return "", ErrInvalid
	}
	return s, nil
}

func asInt(v any) (int, error) {
	switch n := v.(type) {
	case int64:
		return int(n), nil
	case int:
		return n, nil
	case float64:
		if n != float64(int(n)) {
			return 0, ErrInvalid
		}
		return int(n), nil
	default:
		return 0, ErrInvalid
	}
}

// ValidIdentifierKey reports whether key is findings|recommendations|requirements.
func ValidIdentifierKey(key string) bool {
	_, ok := identifierKeys[key]
	return ok
}

// NSForIdentifierKey returns the NS for an identifiers table key.
func NSForIdentifierKey(key string) (string, bool) {
	ns, ok := identifierKeys[key]
	return ns, ok
}
