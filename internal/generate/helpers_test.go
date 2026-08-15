package generate

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/manifest"
	"github.com/robertguss/mycelium/internal/schema"
)

const helperSchema = `namespace = "DEC"
home = "decisions"
filename_pattern = "DEC-{NNN}-{slug}.md"
stage_scoped = false
digits = 3
required_front_matter = ["id"]
required_sections = ["Context"]
`

func TestTypeDiscoveryHelpers(t *testing.T) {
	root := t.TempDir()
	core := filepath.Join(root, "program", "templates")
	pack := filepath.Join(root, "program", "packs", "fixture", "templates")
	if err := os.MkdirAll(core, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(pack, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(core, "decision.schema.toml"), []byte(helperSchema), 0o644); err != nil {
		t.Fatal(err)
	}
	keys, err := listTypeKeys(root)
	if err == nil || keys != nil {
		t.Fatalf("schema without template keys=%v err=%v", keys, err)
	}
	if err := os.WriteFile(filepath.Join(core, "decision.md"), []byte("template"), 0o644); err != nil {
		t.Fatal(err)
	}
	packSchema := strings.ReplaceAll(helperSchema, `namespace = "DEC"`, `namespace = "CMP"`)
	packSchema = strings.ReplaceAll(packSchema, `home = "decisions"`, `home = "reviews/commissioning"`)
	if err := os.WriteFile(filepath.Join(pack, "commissioning.schema.toml"), []byte(packSchema), 0o644); err != nil {
		t.Fatal(err)
	}
	keys, err = listTypeKeys(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 1 || keys[0] != "decision" {
		t.Fatalf("keys=%v", keys)
	}
	if err := os.WriteFile(filepath.Join(pack, "commissioning.md"), []byte("template"), 0o644); err != nil {
		t.Fatal(err)
	}
	keys, err = listTypeKeys(root)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(keys, ",") != "commissioning,decision" {
		t.Fatalf("keys=%v", keys)
	}
	entry, err := findType(root, "commissioning")
	if err != nil {
		t.Fatal(err)
	}
	if entry.Schema.Namespace != "CMP" || !strings.HasSuffix(entry.TemplatePath, "commissioning.md") {
		t.Fatalf("entry=%+v", entry)
	}
	if _, err := findType(root, "missing"); err == nil {
		t.Fatal("missing type must fail")
	}
	if err := os.RemoveAll(filepath.Join(root, "program")); err != nil {
		t.Fatal(err)
	}
	if _, err := findType(root, "decision"); err == nil {
		t.Fatal("discovery error must propagate")
	}
}

func TestLogAndIDHelpers(t *testing.T) {
	if err := validateLog([]byte("# Log\n\n2026-08-15\tnew\tDEC-001\tTitle\n")); err != nil {
		t.Fatal(err)
	}
	if err := validateLog([]byte("invalid\n")); err == nil {
		t.Fatal("invalid log must fail")
	}
	if got := string(appendLogLine(nil, "line")); got != "line\n" {
		t.Fatalf("append empty=%q", got)
	}
	if got := string(appendLogLine([]byte("first"), "second")); got != "first\nsecond\n" {
		t.Fatalf("append without newline=%q", got)
	}
	if got := string(appendLogLine([]byte("first\n"), "second")); got != "first\nsecond\n" {
		t.Fatalf("append with newline=%q", got)
	}

	root := t.TempDir()
	sch := schema.Schema{Namespace: "DEC", Home: "decisions", Digits: 3}
	next, err := nextID(root, sch)
	if err != nil || next != 1 {
		t.Fatalf("absent home next=%d err=%v", next, err)
	}
	home := filepath.Join(root, "decisions")
	if err := os.WriteFile(home, []byte("not a directory"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := nextID(root, sch); err == nil {
		t.Fatal("file home must fail")
	}
	if err := os.Remove(home); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(home, "nested"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"README.md",
		"not-an-id.md",
		"FND-099-wrong-home.md",
		"DEC-001-one.md",
		"DEC-003-three.md",
	} {
		if err := os.WriteFile(filepath.Join(home, name), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	next, err = nextID(root, sch)
	if err != nil || next != 4 {
		t.Fatalf("next=%d err=%v", next, err)
	}
}

func TestRunEarlyTeachingErrors(t *testing.T) {
	var stderr bytes.Buffer
	if code := Run(Options{}, Deps{}); code != 1 {
		t.Fatalf("empty type exit=%d", code)
	}
	if code := Run(Options{TypeKey: "decision"}, Deps{Stderr: &stderr}); code != 1 || !strings.Contains(stderr.String(), "title is required") {
		t.Fatalf("empty title exit=%d stderr=%q", code, stderr.String())
	}
	stderr.Reset()
	if code := Run(Options{TypeKey: "decision", Title: "bad\nvalue"}, Deps{Stderr: &stderr}); code != 1 || !strings.Contains(stderr.String(), "newline or tab") {
		t.Fatalf("newline title exit=%d stderr=%q", code, stderr.String())
	}
	stderr.Reset()
	if code := Run(Options{TypeKey: "decision", Title: "!!!"}, Deps{Stderr: &stderr}); code != 1 || !strings.Contains(stderr.String(), "slugify") {
		t.Fatalf("slug exit=%d stderr=%q", code, stderr.String())
	}
	stderr.Reset()
	if code := Run(Options{TypeKey: "decision", Title: "Valid", Cwd: t.TempDir()}, Deps{Stderr: &stderr}); code != 1 || !strings.Contains(stderr.String(), "not a mycelium instance") {
		t.Fatalf("root exit=%d stderr=%q", code, stderr.String())
	}

	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "mycelium.toml"), []byte("invalid = true"), 0o644); err != nil {
		t.Fatal(err)
	}
	stderr.Reset()
	if code := Run(Options{TypeKey: "decision", Title: "Valid", Cwd: root}, Deps{Stderr: &stderr}); code != 1 || !strings.Contains(stderr.String(), "cannot list registered types") {
		t.Fatalf("schemas exit=%d stderr=%q", code, stderr.String())
	}
	writeGeneratorType(t, root)
	stderr.Reset()
	if code := Run(Options{TypeKey: "missing", Title: "Valid", Cwd: root}, Deps{Stderr: &stderr}); code != 1 || !strings.Contains(stderr.String(), "unknown type") {
		t.Fatalf("unknown exit=%d stderr=%q", code, stderr.String())
	}
	stderr.Reset()
	if code := Run(Options{TypeKey: "decision", Title: "Valid", Cwd: root}, Deps{Stderr: &stderr}); code != 1 || !strings.Contains(stderr.String(), "mycelium.toml invalid") {
		t.Fatalf("manifest exit=%d stderr=%q", code, stderr.String())
	}

	writeGeneratorManifest(t, root)
	stderr.Reset()
	if code := Run(Options{TypeKey: "decision", Title: "Valid", Cwd: root}, Deps{Stderr: &stderr}); code != 1 || !strings.Contains(stderr.String(), "log.md missing") {
		t.Fatalf("log missing exit=%d stderr=%q", code, stderr.String())
	}
	if err := os.WriteFile(filepath.Join(root, "log.md"), []byte("invalid\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	stderr.Reset()
	if code := Run(Options{TypeKey: "decision", Title: "Valid", Cwd: root}, Deps{Stderr: &stderr}); code != 1 || !strings.Contains(stderr.String(), "log.md invalid") {
		t.Fatalf("log invalid exit=%d stderr=%q", code, stderr.String())
	}
}

func writeGeneratorType(t *testing.T, root string) {
	t.Helper()
	dir := filepath.Join(root, "program", "templates")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "decision.schema.toml"), []byte(helperSchema), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "decision.md"), []byte("{{ID}}"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeGeneratorManifest(t *testing.T, root string) {
	t.Helper()
	out, err := manifest.Encode(manifest.Manifest{
		SchemaVersion:         1,
		IdeaName:              "Fixture",
		Slug:                  "fixture",
		State:                 "spark",
		Tier:                  "focused",
		MethodologyVersion:    "2.0.0",
		GeneratedByCLIVersion: "0.1.0-dev",
		CreatedDate:           "2026-08-15",
		UpdatedDate:           "2026-08-15",
		Identifiers:           map[string]manifest.Range{},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "mycelium.toml"), out, 0o644); err != nil {
		t.Fatal(err)
	}
}
