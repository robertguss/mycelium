package schema_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/schema"
)

const decisionSchema = `namespace = "DEC"
home = "decisions"
filename_pattern = "DEC-{NNN}-{slug}.md"
stage_scoped = false
digits = 3
required_front_matter = ["id", "title", "status", "date", "owner"]
required_sections = [
  "Context",
  "Decision",
]

[enums.status]
values = ["Proposed", "Accepted", "Superseded", "Rejected"]
`

func TestParseDecisionSchema(t *testing.T) {
	s, err := schema.Parse([]byte(decisionSchema))
	if err != nil {
		t.Fatal(err)
	}
	if s.Namespace != "DEC" || s.Home != "decisions" || s.Digits != 3 || s.StageScoped {
		t.Fatalf("unexpected: %+v", s)
	}
	if len(s.RequiredFrontMatter) != 5 {
		t.Fatalf("front matter: %v", s.RequiredFrontMatter)
	}
	vals := s.Enums["status"]
	if len(vals) != 4 || vals[0] != "Proposed" {
		t.Fatalf("enums: %v", vals)
	}
}

func TestParseStageScoped(t *testing.T) {
	in := `namespace = "FND"
home = "findings"
filename_pattern = "FND-{NNN}-{slug}.md"
stage_scoped = true
digits = 3
required_front_matter = ["id"]
required_sections = ["Problem"]

[enums.severity]
values = ["Critical", "High"]
`
	s, err := schema.Parse([]byte(in))
	if err != nil {
		t.Fatal(err)
	}
	if !s.StageScoped {
		t.Fatal("want stage_scoped")
	}
}

func TestMissingRequired(t *testing.T) {
	_, err := schema.Parse([]byte(`home = "decisions"
filename_pattern = "x"
stage_scoped = false
digits = 3
required_front_matter = ["id"]
required_sections = ["Context"]
`))
	if !errors.Is(err, schema.ErrRequired) {
		t.Fatalf("err = %v", err)
	}
}

func TestParseRejectsInvalidSchemas(t *testing.T) {
	cases := []struct {
		name string
		body string
		want error
	}{
		{name: "invalid toml", body: `namespace = [`, want: schema.ErrInvalid},
		{name: "missing home", body: `namespace = "DEC"`, want: schema.ErrRequired},
		{name: "missing filename pattern", body: `namespace = "DEC"
home = "decisions"`, want: schema.ErrRequired},
		{name: "missing stage scoped", body: `namespace = "DEC"
home = "decisions"
filename_pattern = "x"`, want: schema.ErrRequired},
		{name: "invalid stage scoped", body: `namespace = "DEC"
home = "decisions"
filename_pattern = "x"
stage_scoped = "false"`, want: schema.ErrInvalid},
		{name: "missing digits", body: `namespace = "DEC"
home = "decisions"
filename_pattern = "x"
stage_scoped = false`, want: schema.ErrRequired},
		{name: "zero digits", body: `namespace = "DEC"
home = "decisions"
filename_pattern = "x"
stage_scoped = false
digits = 0`, want: schema.ErrInvalid},
		{name: "missing front matter", body: `namespace = "DEC"
home = "decisions"
filename_pattern = "x"
stage_scoped = false
digits = 3`, want: schema.ErrRequired},
		{name: "invalid front matter item", body: `namespace = "DEC"
home = "decisions"
filename_pattern = "x"
stage_scoped = false
digits = 3
required_front_matter = ["id", 1]
required_sections = ["Body"]`, want: schema.ErrRequired},
		{name: "missing sections", body: `namespace = "DEC"
home = "decisions"
filename_pattern = "x"
stage_scoped = false
digits = 3
required_front_matter = ["id"]`, want: schema.ErrRequired},
		{name: "invalid enums table", body: `namespace = "DEC"
home = "decisions"
filename_pattern = "x"
stage_scoped = false
digits = 3
required_front_matter = ["id"]
required_sections = ["Body"]
enums = "bad"`, want: schema.ErrInvalid},
		{name: "invalid enum entry", body: decisionSchema + `
[enums.owner]
value = "bad"`, want: schema.ErrRequired},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := schema.Parse([]byte(tt.body))
			if !errors.Is(err, tt.want) {
				t.Fatalf("err=%v want %v", err, tt.want)
			}
		})
	}
}

func TestLoadFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "decision.schema.toml")
	if err := os.WriteFile(path, []byte(decisionSchema), 0o644); err != nil {
		t.Fatal(err)
	}
	s, err := schema.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if s.Namespace != "DEC" {
		t.Fatalf("got %+v", s)
	}
	if _, err := schema.Load(filepath.Join(dir, "missing.schema.toml")); !os.IsNotExist(err) {
		t.Fatalf("missing load err=%v", err)
	}
}

func TestLoadEmbeddedDecisionSchema(t *testing.T) {
	path := filepath.Join("..", "embed", "program", "templates", "decision.schema.toml")
	s, err := schema.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if s.Namespace != "DEC" || s.Digits != 3 {
		t.Fatalf("got %+v", s)
	}
	if len(s.Enums["status"]) == 0 {
		t.Fatal("missing status enum")
	}
}

func TestCouncilPackSchemasParse(t *testing.T) {
	cases := map[string]string{
		"commissioning":  "CMP",
		"model-report":   "RPT",
		"reconciliation": "RCL",
	}
	for key, namespace := range cases {
		t.Run(key, func(t *testing.T) {
			path := filepath.Join("..", "..", "program", "packs", "council", "templates", key+".schema.toml")
			s, err := schema.Load(path)
			if err != nil {
				t.Fatal(err)
			}
			if s.Namespace != namespace || s.Digits != 3 || s.StageScoped {
				t.Fatalf("got %+v", s)
			}
		})
	}
}

func TestDiscoverCoreAndPackTemplates(t *testing.T) {
	root := t.TempDir()
	coreDir := filepath.Join(root, "program", "templates")
	packDir := filepath.Join(root, "program", "packs", "council", "templates")
	if err := os.MkdirAll(coreDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(packDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(coreDir, "decision.schema.toml"), []byte(decisionSchema), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(coreDir, "decision.md"), []byte("decision"), 0o644); err != nil {
		t.Fatal(err)
	}
	commissioning := `namespace = "CMP"
home = "reviews/commissioning"
filename_pattern = "CMP-{NNN}-{slug}.md"
stage_scoped = false
digits = 3
required_front_matter = ["id"]
required_sections = ["Prompt"]
`
	if err := os.WriteFile(filepath.Join(packDir, "commissioning.schema.toml"), []byte(commissioning), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(packDir, "commissioning.md"), []byte("commissioning"), 0o644); err != nil {
		t.Fatal(err)
	}

	entries, err := schema.Discover(root)
	if err != nil {
		t.Fatal(err)
	}
	if !hasEntry(entries, "decision") || !hasEntry(entries, "commissioning") {
		t.Fatalf("entries=%v", entries)
	}

	if err := os.RemoveAll(filepath.Join(root, "program", "packs")); err != nil {
		t.Fatal(err)
	}
	entries, err = schema.Discover(root)
	if err != nil {
		t.Fatal(err)
	}
	if !hasEntry(entries, "decision") || hasEntry(entries, "commissioning") {
		t.Fatalf("entries without pack=%v", entries)
	}
}

func TestDiscoverErrorsAndFirstRegistrationWins(t *testing.T) {
	if _, err := schema.Discover(t.TempDir()); err == nil {
		t.Fatal("missing core templates must fail")
	}

	root := t.TempDir()
	coreDir := filepath.Join(root, "program", "templates")
	if err := os.MkdirAll(coreDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(coreDir, "decision.schema.toml"), []byte(decisionSchema), 0o644); err != nil {
		t.Fatal(err)
	}
	packsPath := filepath.Join(root, "program", "packs")
	if err := os.WriteFile(packsPath, []byte("not a directory"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := schema.Discover(root); err == nil {
		t.Fatal("unreadable packs directory must fail")
	}
	if err := os.Remove(packsPath); err != nil {
		t.Fatal(err)
	}

	packDir := filepath.Join(packsPath, "fixture", "templates")
	if err := os.MkdirAll(packDir, 0o755); err != nil {
		t.Fatal(err)
	}
	duplicate := strings.Replace(decisionSchema, `namespace = "DEC"`, `namespace = "ZZZ"`, 1)
	if err := os.WriteFile(filepath.Join(packDir, "decision.schema.toml"), []byte(duplicate), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(packsPath, "README.md"), []byte("ignored"), 0o644); err != nil {
		t.Fatal(err)
	}
	entries, err := schema.Discover(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Schema.Namespace != "DEC" {
		t.Fatalf("first registration not retained: %+v", entries)
	}

	if err := os.WriteFile(filepath.Join(packDir, "broken.schema.toml"), []byte("not toml ="), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := schema.Discover(root); err == nil {
		t.Fatal("invalid pack schema must fail")
	}
}

func hasEntry(entries []schema.Entry, key string) bool {
	for _, entry := range entries {
		if entry.Key == key {
			return true
		}
	}
	return false
}
