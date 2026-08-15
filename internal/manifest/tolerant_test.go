package manifest_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/manifest"
)

const tolerantBase = `schema_version = 1
idea_name = "Legacy One"
slug = "legacy-one"
state = "spark"
tier = "focused"
methodology_version = "2.0.0"
generated_by_cli_version = "0.1.0-dev"
created_date = "2026-08-01"
updated_date = "2026-08-01"
revisit = ""
`

func TestParseTolerantMissingGithubRepo(t *testing.T) {
	m, err := manifest.ParseTolerant([]byte(tolerantBase))
	if err != nil {
		t.Fatal(err)
	}
	if m.Slug != "legacy-one" || m.GithubRepo != "" {
		t.Fatalf("got %+v", m)
	}
	_, err = manifest.Parse([]byte(tolerantBase))
	if !errors.Is(err, manifest.ErrRequired) {
		t.Fatalf("strict Parse should still require github_repo, got %v", err)
	}
}

func TestParseTolerantUnknownKey(t *testing.T) {
	in := tolerantBase + "github_repo = \"\"\nlegacy_note = \"append-only\"\n"
	m, err := manifest.ParseTolerant([]byte(in))
	if err != nil {
		t.Fatal(err)
	}
	if m.Slug != "legacy-one" {
		t.Fatalf("slug=%q", m.Slug)
	}
	_, err = manifest.Parse([]byte(in))
	if !errors.Is(err, manifest.ErrUnknownKey) {
		t.Fatalf("strict Parse should refuse unknown key, got %v", err)
	}
}

func TestParseTolerantNotTOML(t *testing.T) {
	_, err := manifest.ParseTolerant([]byte("{{{not toml"))
	if !errors.Is(err, manifest.ErrInvalid) {
		t.Fatalf("err=%v", err)
	}
}

func TestParseTolerantNoSlug(t *testing.T) {
	_, err := manifest.ParseTolerant([]byte("not = \"valid\"\n"))
	if !errors.Is(err, manifest.ErrRequired) {
		t.Fatalf("err=%v want ErrRequired", err)
	}
}

func TestRequiredKeysFromContractOmitsGithub(t *testing.T) {
	md := []byte(`# Manifest

## Required fields

| Field | Rule |
| --- | --- |
| ` + "`schema_version`" + ` | Must be 1 |
| ` + "`idea_name`" + ` | title |
| ` + "`slug`" + ` | slug |
| ` + "`state`" + ` | state |
| ` + "`tier`" + ` | tier |
| ` + "`methodology_version`" + ` | pin |
| ` + "`generated_by_cli_version`" + ` | cli |
| ` + "`created_date`" + ` | date |
| ` + "`updated_date`" + ` | date |
| ` + "`revisit`" + ` | revisit |

## Refuse on unknown keys

- Unknown top-level keys → refuse.
`)
	keys, err := manifest.RequiredKeysFromContractMarkdown(md)
	if err != nil {
		t.Fatal(err)
	}
	for _, k := range keys {
		if k == "github_repo" {
			t.Fatalf("github_repo should be omitted: %v", keys)
		}
	}
	want := []string{
		"schema_version", "idea_name", "slug", "state", "tier",
		"methodology_version", "generated_by_cli_version",
		"created_date", "updated_date", "revisit",
	}
	if len(keys) != len(want) {
		t.Fatalf("keys=%v want %v", keys, want)
	}
	for i := range want {
		if keys[i] != want[i] {
			t.Fatalf("keys[%d]=%q want %q", i, keys[i], want[i])
		}
	}

	m, err := manifest.ParseWithRequired([]byte(tolerantBase), keys)
	if err != nil {
		t.Fatal(err)
	}
	if m.GithubRepo != "" {
		t.Fatalf("github_repo=%q", m.GithubRepo)
	}
}

func TestRequiredKeysCurrentContractIncludesGithub(t *testing.T) {
	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	// Prefer reading the repo contract via a temp instance copy path.
	b, err := os.ReadFile(filepath.Join(root, "program", "contracts", "manifest.md"))
	if err != nil {
		t.Fatal(err)
	}
	keys, err := manifest.RequiredKeysFromContractMarkdown(b)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, k := range keys {
		if k == "github_repo" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("current contract should list github_repo: %v", keys)
	}
}

func TestRequiredKeysForInstanceFallback(t *testing.T) {
	dir := t.TempDir()
	keys := manifest.RequiredKeysForInstance(dir)
	def := manifest.DefaultRequiredKeys()
	if strings.Join(keys, ",") != strings.Join(def, ",") {
		t.Fatalf("fallback keys=%v want %v", keys, def)
	}
	found := false
	for _, k := range keys {
		if k == "github_repo" {
			found = true
		}
	}
	if !found {
		t.Fatal("fallback must include github_repo")
	}
}

func TestRequiredKeysForInstanceReadsFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "program", "contracts")
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
	md := "# x\n\n## Required fields\n\n| Field | Rule |\n| --- | --- |\n| `slug` | s |\n| `state` | s |\n"
	if err := os.WriteFile(filepath.Join(path, "manifest.md"), []byte(md), 0o644); err != nil {
		t.Fatal(err)
	}
	keys := manifest.RequiredKeysForInstance(dir)
	if len(keys) != 2 || keys[0] != "slug" || keys[1] != "state" {
		t.Fatalf("keys=%v", keys)
	}
}

func TestValidIdentifierKeyAndNS(t *testing.T) {
	if !manifest.ValidIdentifierKey("findings") {
		t.Fatal("findings should be valid")
	}
	if manifest.ValidIdentifierKey("phases") {
		t.Fatal("phases should be invalid")
	}
	ns, ok := manifest.NSForIdentifierKey("recommendations")
	if !ok || ns != "REC" {
		t.Fatalf("ns=%q ok=%v", ns, ok)
	}
	if _, ok := manifest.NSForIdentifierKey("nope"); ok {
		t.Fatal("expected missing")
	}
}

func TestEncodeWithIdentifiers(t *testing.T) {
	m := manifest.Manifest{
		SchemaVersion:         1,
		IdeaName:              "Legacy One",
		Slug:                  "legacy-one",
		State:                 "spark",
		Tier:                  "focused",
		MethodologyVersion:    "2.0.0",
		GeneratedByCLIVersion: "0.1.0-dev",
		CreatedDate:           "2026-08-01",
		UpdatedDate:           "2026-08-01",
		Revisit:               "",
		GithubRepo:            "",
		Identifiers: map[string]manifest.Range{
			"findings": {NS: "FND", Start: 1, End: 9, Raw: "FND-001..FND-009"},
		},
		Deviations: []manifest.Deviation{{Convention: "x", Reason: "y"}},
	}
	b, err := manifest.Encode(m)
	if err != nil {
		t.Fatal(err)
	}
	got, err := manifest.Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	if got.Identifiers["findings"].Raw != "FND-001..FND-009" {
		t.Fatalf("identifiers=%v", got.Identifiers)
	}
	if len(got.Deviations) != 1 {
		t.Fatalf("deviations=%v", got.Deviations)
	}
}

func TestParseWithRequiredEmptyGithubWhenOmitted(t *testing.T) {
	keys := []string{
		"schema_version", "idea_name", "slug", "state", "tier",
		"methodology_version", "generated_by_cli_version",
		"created_date", "updated_date", "revisit",
	}
	m, err := manifest.ParseWithRequired([]byte(tolerantBase), keys)
	if err != nil {
		t.Fatal(err)
	}
	if m.GithubRepo != "" {
		t.Fatalf("github=%q", m.GithubRepo)
	}
}

func TestRequiredKeysEmptyTableInvalid(t *testing.T) {
	_, err := manifest.RequiredKeysFromContractMarkdown([]byte("## Required fields\n\nno table\n"))
	if !errors.Is(err, manifest.ErrInvalid) {
		t.Fatalf("err=%v", err)
	}
}

func TestRequiredKeysForInstanceBadMarkdownFallsBack(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "program", "contracts")
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(path, "manifest.md"), []byte("## Other\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	keys := manifest.RequiredKeysForInstance(dir)
	if len(keys) != len(manifest.DefaultRequiredKeys()) {
		t.Fatalf("keys=%v", keys)
	}
}

func TestParseFloatSchemaVersion(t *testing.T) {
	in := strings.Replace(tolerantBase, "schema_version = 1", "schema_version = 1.0", 1)
	in += "github_repo = \"\"\n"
	m, err := manifest.Parse([]byte(in))
	if err != nil {
		t.Fatal(err)
	}
	if m.SchemaVersion != 1 {
		t.Fatalf("sv=%d", m.SchemaVersion)
	}
}

func TestParseBadGithubRepoType(t *testing.T) {
	in := tolerantBase + "github_repo = 1\n"
	_, err := manifest.Parse([]byte(in))
	if !errors.Is(err, manifest.ErrInvalid) {
		t.Fatalf("err=%v", err)
	}
}

func TestRequiredKeysSkipsJunkAndDupes(t *testing.T) {
	md := []byte(`## Required fields

| Field | Rule |
| --- | --- |
| | empty |
| NotAKey | bad |
| ` + "`slug`" + ` | one |
| ` + "`slug`" + ` | dupe |
| ` + "`state2`" + ` | with digit |

## Next
`)
	keys, err := manifest.RequiredKeysFromContractMarkdown(md)
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 2 || keys[0] != "slug" || keys[1] != "state2" {
		t.Fatalf("keys=%v", keys)
	}
}
