package schema_test

import (
	"errors"
	"os"
	"path/filepath"
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
