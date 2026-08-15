package manifest_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/manifest"
)

const validTOML = `schema_version = 1
idea_name = "Garden lighting"
slug = "garden-lighting"
state = "spark"
tier = "focused"
methodology_version = "2.0.0"
generated_by_cli_version = "0.1.0-dev"
created_date = "2026-08-14"
updated_date = "2026-08-14"
revisit = ""
github_repo = ""
`

func TestParseValid(t *testing.T) {
	m, err := manifest.Parse([]byte(validTOML))
	if err != nil {
		t.Fatal(err)
	}
	if m.Slug != "garden-lighting" || m.State != "spark" || m.SchemaVersion != 1 {
		t.Fatalf("unexpected: %+v", m)
	}
}

func TestParseWithRangesAndDeviations(t *testing.T) {
	in := validTOML + `
[identifiers]
findings = "FND-001..FND-099"
recommendations = "REC-001..REC-099"
requirements = "REQ-001..REQ-299"

[[deviations]]
convention = "extra-top-level:notes.md"
reason = "scratch pad"
`
	m, err := manifest.Parse([]byte(in))
	if err != nil {
		t.Fatal(err)
	}
	if err := m.InRange("findings", 50); err != nil {
		t.Fatal(err)
	}
	if len(m.Deviations) != 1 || m.Deviations[0].Convention == "" {
		t.Fatalf("deviations: %+v", m.Deviations)
	}
}

func TestMissingField(t *testing.T) {
	in := strings.Replace(validTOML, "tier = \"focused\"\n", "", 1)
	_, err := manifest.Parse([]byte(in))
	if !errors.Is(err, manifest.ErrRequired) {
		t.Fatalf("err = %v", err)
	}
}

func TestBadState(t *testing.T) {
	in := strings.Replace(validTOML, `state = "spark"`, `state = "boiling"`, 1)
	_, err := manifest.Parse([]byte(in))
	if !errors.Is(err, manifest.ErrInvalid) {
		t.Fatalf("err = %v", err)
	}
}

func TestBadTier(t *testing.T) {
	in := strings.Replace(validTOML, `tier = "focused"`, `tier = "enterprise"`, 1)
	_, err := manifest.Parse([]byte(in))
	if !errors.Is(err, manifest.ErrInvalid) {
		t.Fatalf("err = %v", err)
	}
}

func TestBadSchemaVersion(t *testing.T) {
	in := strings.Replace(validTOML, `schema_version = 1`, `schema_version = 2`, 1)
	_, err := manifest.Parse([]byte(in))
	if !errors.Is(err, manifest.ErrInvalid) {
		t.Fatalf("err = %v", err)
	}
}

func TestUnknownTopKey(t *testing.T) {
	in := validTOML + "extra = true\n"
	_, err := manifest.Parse([]byte(in))
	if !errors.Is(err, manifest.ErrUnknownKey) {
		t.Fatalf("err = %v", err)
	}
}

func TestUnknownIdentifierKey(t *testing.T) {
	in := validTOML + "\n[identifiers]\nphases = \"PHASE-01..PHASE-09\"\n"
	_, err := manifest.Parse([]byte(in))
	if !errors.Is(err, manifest.ErrUnknownKey) {
		t.Fatalf("err = %v", err)
	}
}

func TestExtraDeviationKey(t *testing.T) {
	in := validTOML + `
[[deviations]]
convention = "x"
reason = "y"
note = "nope"
`
	_, err := manifest.Parse([]byte(in))
	if !errors.Is(err, manifest.ErrUnknownKey) {
		t.Fatalf("err = %v", err)
	}
}

func TestSimmeringNeedsRevisit(t *testing.T) {
	in := strings.Replace(validTOML, `state = "spark"`, `state = "simmering"`, 1)
	_, err := manifest.Parse([]byte(in))
	if !errors.Is(err, manifest.ErrInvalid) {
		t.Fatalf("err = %v", err)
	}
	in = strings.Replace(in, `revisit = ""`, `revisit = "2026-09-01"`, 1)
	if _, err := manifest.Parse([]byte(in)); err != nil {
		t.Fatal(err)
	}
}

func TestSlugMustMatchIdeaName(t *testing.T) {
	in := strings.Replace(validTOML, `slug = "garden-lighting"`, `slug = "wrong"`, 1)
	_, err := manifest.Parse([]byte(in))
	if !errors.Is(err, manifest.ErrInvalid) {
		t.Fatalf("err = %v", err)
	}
}

func TestBadDate(t *testing.T) {
	tests := []struct {
		name string
		old  string
		new  string
	}{
		{name: "created not iso", old: `created_date = "2026-08-14"`, new: `created_date = "08/14/2026"`},
		{name: "updated garbage", old: `updated_date = "2026-08-14"`, new: `updated_date = "yesterday"`},
		{name: "created missing dashes", old: `created_date = "2026-08-14"`, new: `created_date = "20260814"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.Replace(validTOML, tt.old, tt.new, 1)
			_, err := manifest.Parse([]byte(in))
			if !errors.Is(err, manifest.ErrInvalid) {
				t.Fatalf("err = %v, want ErrInvalid", err)
			}
		})
	}
}

func TestRangeMembership(t *testing.T) {
	in := validTOML + `
[identifiers]
findings = "FND-001..FND-099"
`
	m, err := manifest.Parse([]byte(in))
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name    string
		key     string
		n       int
		wantErr error
	}{
		{name: "inside", key: "findings", n: 1},
		{name: "inside end", key: "findings", n: 99},
		{name: "outside", key: "findings", n: 100, wantErr: manifest.ErrRange},
		{name: "missing declaration", key: "recommendations", n: 1, wantErr: manifest.ErrRange},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := m.InRange(tt.key, tt.n)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("err = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestParseRangeErrors(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		ns   string
	}{
		{name: "start>end", raw: "FND-099..FND-001", ns: "FND"},
		{name: "NS mismatch", raw: "REC-001..REC-099", ns: "FND"},
		{name: "mixed NS", raw: "FND-001..REC-099", ns: "FND"},
		{name: "bad grammar", raw: "FND-001-FND-099", ns: "FND"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := manifest.ParseRange(tt.raw, tt.ns)
			if !errors.Is(err, manifest.ErrRange) {
				t.Fatalf("err = %v", err)
			}
		})
	}
}

func TestEncodeRoundTrip(t *testing.T) {
	m := manifest.Manifest{
		SchemaVersion:         1,
		IdeaName:              "Garden lighting",
		Slug:                  "garden-lighting",
		State:                 "spark",
		Tier:                  "focused",
		MethodologyVersion:    "2.0.0",
		GeneratedByCLIVersion: "0.1.0-dev",
		CreatedDate:           "2026-08-14",
		UpdatedDate:           "2026-08-14",
		Revisit:               "",
		GithubRepo:            "",
	}
	b, err := manifest.Encode(m)
	if err != nil {
		t.Fatal(err)
	}
	got, err := manifest.Parse(b)
	if err != nil {
		t.Fatalf("parse encoded: %v\n%s", err, b)
	}
	if got.IdeaName != m.IdeaName || got.Slug != m.Slug || got.State != m.State {
		t.Fatalf("round-trip mismatch: %+v", got)
	}
	if strings.Contains(string(b), "[identifiers]") {
		t.Fatalf("empty identifiers should be omitted:\n%s", b)
	}
}
