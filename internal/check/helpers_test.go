package check

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/manifest"
	"github.com/robertguss/mycelium/internal/schema"
	"github.com/robertguss/mycelium/internal/wakebrief"
)

func TestFindRootBoundaries(t *testing.T) {
	root := t.TempDir()
	nested := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "mycelium.toml"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := FindRoot(nested)
	if err != nil || got != root {
		t.Fatalf("root=%q err=%v", got, err)
	}

	blocked := t.TempDir()
	if err := os.Mkdir(filepath.Join(blocked, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := FindRoot(blocked); !errors.Is(err, ErrNotInstance) {
		t.Fatalf("git boundary err=%v", err)
	}
	if _, err := FindRoot(t.TempDir()); !errors.Is(err, ErrNotInstance) {
		t.Fatalf("filesystem boundary err=%v", err)
	}
}

func TestTierBindChecks(t *testing.T) {
	root := t.TempDir()
	var findings []Finding
	add := collectFindings(&findings)
	checkTierBinds(root, "unknown", add)
	if !hasConvention(findings, "tier") {
		t.Fatalf("unknown tier findings=%v", findings)
	}

	findings = nil
	checkTierBinds(root, "focused", add)
	if !hasText(findings, "tier file missing") {
		t.Fatalf("missing tier findings=%v", findings)
	}
	tierDir := filepath.Join(root, "program", "tiers")
	if err := os.MkdirAll(tierDir, 0o755); err != nil {
		t.Fatal(err)
	}
	tierPath := filepath.Join(tierDir, "focused.toml")
	if err := os.WriteFile(tierPath, []byte("binds = ["), 0o644); err != nil {
		t.Fatal(err)
	}
	findings = nil
	checkTierBinds(root, "focused", add)
	if !hasText(findings, "tier file invalid") {
		t.Fatalf("invalid tier findings=%v", findings)
	}

	body := `binds = ["manifest", "required/", "index.md", "other.md"]`
	if err := os.WriteFile(tierPath, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "mycelium.toml"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	findings = nil
	checkTierBinds(root, "focused", add)
	for _, want := range []string{"required/", "index.md", "other.md"} {
		if !hasText(findings, want) {
			t.Fatalf("missing %q in %v", want, findings)
		}
	}
	if err := os.Mkdir(filepath.Join(root, "required"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"index.md", "other.md"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	findings = nil
	checkTierBinds(root, "focused", add)
	if len(findings) != 0 {
		t.Fatalf("satisfied binds findings=%v", findings)
	}
}

func TestTopLevelAndIndexChecks(t *testing.T) {
	var findings []Finding
	add := collectFindings(&findings)
	checkTopLevel(filepath.Join(t.TempDir(), "missing"), manifest.Manifest{}, nil, false, add)
	if !hasText(findings, "cannot read instance root") {
		t.Fatalf("read findings=%v", findings)
	}

	root := t.TempDir()
	for _, name := range []string{"README.md", "decisions", "reviews", "notes", "waived"} {
		if err := os.MkdirAll(filepath.Join(root, name), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	m := manifest.Manifest{Deviations: []manifest.Deviation{
		{Convention: "extra-top-level:waived/", Reason: "fixture"},
		{},
	}}
	findings = nil
	checkTopLevel(root, m, map[string]struct{}{"decisions": {}}, false, add)
	for _, want := range []string{"deviation row missing", "council pack absent", "undeclared extra top-level path notes"} {
		if !hasText(findings, want) {
			t.Fatalf("missing %q in %v", want, findings)
		}
	}
	findings = nil
	checkTopLevel(root, m, map[string]struct{}{"decisions": {}, "notes": {}}, true, add)
	if hasText(findings, "reviews/") || hasText(findings, "notes") {
		t.Fatalf("allowed paths findings=%v", findings)
	}

	findings = nil
	checkIndex(root, add)
	if !hasText(findings, "index.md missing") {
		t.Fatalf("missing index findings=%v", findings)
	}
	if err := os.WriteFile(filepath.Join(root, "index.md"), []byte("no headings"), 0o644); err != nil {
		t.Fatal(err)
	}
	findings = nil
	checkIndex(root, add)
	if !hasText(findings, "missing H1") || !hasText(findings, `"State"`) {
		t.Fatalf("incomplete index findings=%v", findings)
	}
	complete := "# Index\n\n## State\n\n## Artifacts\n\n## Log tail\n\n## Wake\n"
	if err := os.WriteFile(filepath.Join(root, "index.md"), []byte(complete), 0o644); err != nil {
		t.Fatal(err)
	}
	findings = nil
	checkIndex(root, add)
	if len(findings) != 0 {
		t.Fatalf("complete index findings=%v", findings)
	}
}

func TestArtifactIndexAndValidationBranches(t *testing.T) {
	root := t.TempDir()
	sch := schema.Schema{
		Namespace:           "DEC",
		Home:                "decisions",
		Digits:              3,
		RequiredFrontMatter: []string{"id", "date", "status"},
		RequiredSections:    []string{"Context"},
		Enums:               map[string][]string{"status": {"Proposed"}},
	}
	var findings []Finding
	add := collectFindings(&findings)
	var result Result

	if err := os.WriteFile(filepath.Join(root, "decisions"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	arts := buildArtifactIndex(root, []schema.Schema{sch}, &result, add)
	if len(arts) != 0 || !hasText(findings, "not a directory") {
		t.Fatalf("file home arts=%v findings=%v", arts, findings)
	}
	if err := os.Remove(filepath.Join(root, "decisions")); err != nil {
		t.Fatal(err)
	}
	home := filepath.Join(root, "decisions")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range map[string]string{
		"README.md":               "ignored",
		"bad.md":                  "bad",
		"DEC-001-one.md":          validDecisionDoc("DEC-001"),
		"DEC-001-duplicate.md":    validDecisionDoc("DEC-001"),
		"DEC-002-invalid-meta.md": "not front matter",
	} {
		if err := os.WriteFile(filepath.Join(home, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	findings = nil
	arts = buildArtifactIndex(root, []schema.Schema{sch}, &result, add)
	if len(arts) != 3 || result.Artifacts != 3 {
		t.Fatalf("arts=%v count=%d", arts, result.Artifacts)
	}
	if !hasText(findings, "bad.md") || !hasText(findings, "duplicate id DEC-001") {
		t.Fatalf("index findings=%v", findings)
	}

	mismatch := `+++
id = "DEC-999"
date = "bad-date"
status = "Wrong"
+++

## Other
`
	if err := os.WriteFile(filepath.Join(home, "DEC-001-one.md"), []byte(mismatch), 0o644); err != nil {
		t.Fatal(err)
	}
	missingKeys := "+++\nid = \"DEC-001\"\n+++\n\n## Context\n\nok\n"
	if err := os.WriteFile(filepath.Join(home, "DEC-001-duplicate.md"), []byte(missingKeys), 0o644); err != nil {
		t.Fatal(err)
	}
	findings = nil
	checkFrontMatterAndSections(root, append(arts, artifactFile{Rel: "decisions/missing.md", Home: "decisions"}), map[string]schema.Schema{"decisions": sch}, add)
	for _, want := range []string{"front matter invalid", "does not match filename", "not in enum", "date must be", `missing required H2 "Context"`, "cannot read decisions/missing.md"} {
		if !hasText(findings, want) {
			t.Fatalf("missing %q in %v", want, findings)
		}
	}

	findings = nil
	stageSchema := sch
	stageSchema.StageScoped = true
	checkStageScoped(manifest.Manifest{Identifiers: map[string]manifest.Range{}}, arts, map[string]schema.Schema{"decisions": stageSchema}, add)
	if !hasText(findings, "outside declared range") {
		t.Fatalf("stage findings=%v", findings)
	}
}

func TestLogWakeAndOutputHelpers(t *testing.T) {
	root := t.TempDir()
	var findings []Finding
	add := collectFindings(&findings)
	if got := checkLog(root, add); got != nil || !hasText(findings, "log.md missing") {
		t.Fatalf("missing log got=%q findings=%v", got, findings)
	}
	if err := os.WriteFile(filepath.Join(root, "log.md"), []byte("bad\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	findings = nil
	checkLog(root, add)
	if !hasText(findings, "illegal prefix") {
		t.Fatalf("bad log findings=%v", findings)
	}
	wakeLog := []byte("# Log\n2026-08-15\twake\t-\tResume\n")
	if err := os.WriteFile(filepath.Join(root, "log.md"), wakeLog, 0o644); err != nil {
		t.Fatal(err)
	}
	findings = nil
	got := checkLog(root, add)
	if len(findings) != 0 || !bytes.Equal(got, wakeLog) {
		t.Fatalf("valid log got=%q findings=%v", got, findings)
	}
	if logHasWake(nil) || logHasWake([]byte("not parseable")) || !logHasWake(wakeLog) {
		t.Fatal("wake detection mismatch")
	}

	checkWakeBrief(root, []byte("2026-08-15\tcheck\t-\tCheck\n"), add)
	findings = nil
	checkWakeBrief(root, wakeLog, add)
	if !hasText(findings, "LATEST.md is missing") {
		t.Fatalf("missing wake findings=%v", findings)
	}
	briefDir := filepath.Join(root, "briefs")
	if err := os.MkdirAll(briefDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(briefDir, "LATEST.md"), []byte("## First\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	findings = nil
	checkWakeBrief(root, wakeLog, add)
	if len(findings) == 0 {
		t.Fatal("incomplete brief must fail")
	}
	var complete strings.Builder
	for _, heading := range wakebrief.RequiredH2s() {
		complete.WriteString("## " + heading + "\n\n")
	}
	if err := os.WriteFile(filepath.Join(briefDir, "LATEST.md"), []byte(complete.String()), 0o644); err != nil {
		t.Fatal(err)
	}
	findings = nil
	checkWakeBrief(root, wakeLog, add)
	if len(findings) != 0 {
		t.Fatalf("complete wake findings=%v", findings)
	}

	var out bytes.Buffer
	WriteOK(&out, Result{Slug: "fixture", State: "spark", Tier: "focused", Artifacts: 3})
	for _, want := range []string{"mycelium check: ok", "instance: fixture", "artifacts: 3"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("output=%q missing %q", out.String(), want)
		}
	}
}

func collectFindings(target *[]Finding) func(string, string, string, string) {
	return func(what, convention, contract, fix string) {
		*target = append(*target, Finding{
			What:       what,
			Convention: convention,
			Contract:   contract,
			Fix:        fix,
		})
	}
}

func validDecisionDoc(id string) string {
	return "+++\nid = \"" + id + "\"\ndate = \"2026-08-15\"\nstatus = \"Proposed\"\n+++\n\n## Context\n\nok\n"
}

func hasText(findings []Finding, want string) bool {
	for _, finding := range findings {
		text := finding.What + "\n" + finding.Convention + "\n" + finding.Contract + "\n" + finding.Fix
		if strings.Contains(text, want) {
			return true
		}
	}
	return false
}

func hasConvention(findings []Finding, convention string) bool {
	for _, finding := range findings {
		if finding.Convention == convention {
			return true
		}
	}
	return false
}
