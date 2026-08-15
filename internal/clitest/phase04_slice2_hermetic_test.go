package clitest_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/clitest"
)

func TestPackTemplateGenerationHermetic(t *testing.T) {
	bin := clitest.Bin(t)
	stubDir := installNetworkStubs(t)
	home := t.TempDir()
	workDir := t.TempDir()
	env := hermeticEnv(stubDir, home)

	code, _, stderr := clitest.Run(t, bin, workDir, env, "new", "idea", "Pack Templates", "--offline")
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr)
	}
	assertNoGh(t, stderr)
	inst := filepath.Join(workDir, "pack-templates")

	code, _, stderr = clitest.Run(t, bin, workDir, env, "new", "commissioning", "SQLite", "--dir", inst)
	if code != 0 {
		t.Fatalf("new commissioning exit %d stderr=%q", code, stderr)
	}
	assertNoGh(t, stderr)
	assertGeneratedFile(t, filepath.Join(inst, "reviews", "commissioning", "CMP-001-*.md"), []string{
		`rung = "second-opinion"`,
		"opt_in = true",
		`cost_class = "cheap"`,
		`adapter = "manual"`,
		"## Prompt",
		"## Attachments",
		"## Cost",
	})

	code, _, stderr = clitest.Run(t, bin, workDir, env, "new", "model-report", "SQLite report", "--dir", inst)
	if code != 0 {
		t.Fatalf("new model-report exit %d stderr=%q", code, stderr)
	}
	assertNoGh(t, stderr)
	assertGeneratedFile(t, filepath.Join(inst, "reviews", "reports", "RPT-001-*.md"), []string{
		`model = "fill-me"`,
		`commissioning = "CMP-000"`,
		`rung = "second-opinion"`,
		`adapter = "manual"`,
		`prompt_sha256 = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"`,
		"## Position",
		"## Findings",
		"## Dissent",
	})

	code, _, stderr = clitest.Run(t, bin, workDir, env, "new", "reconciliation", "SQLite council", "--dir", inst)
	if code != 0 {
		t.Fatalf("new reconciliation exit %d stderr=%q", code, stderr)
	}
	assertNoGh(t, stderr)
	assertGeneratedFile(t, filepath.Join(inst, "reviews", "reconciliations", "RCL-001-*.md"), []string{
		`commissioning = "CMP-000"`,
		`rung = "council"`,
		"## Convergence",
		"## Material disagreement",
		"## Evidence unique to one report",
		"## Contradictory evidence",
		"## Different assumptions",
		"## Different scope interpretations",
		"## Recommendations independently supported",
		"## Questions requiring another spike",
		"## Final reconciled recommendation",
		"## Retained dissent",
	})
	assertNoHomeTouch(t, home)
}

func TestPackTypeUnknownWhenPackAbsentHermetic(t *testing.T) {
	bin := clitest.Bin(t)
	stubDir := installNetworkStubs(t)
	home := t.TempDir()
	workDir := t.TempDir()
	env := hermeticEnv(stubDir, home)

	code, _, stderr := clitest.Run(t, bin, workDir, env, "new", "idea", "Pack Disabled Types", "--offline")
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr)
	}
	assertNoGh(t, stderr)
	inst := filepath.Join(workDir, "pack-disabled-types")
	if err := os.RemoveAll(filepath.Join(inst, "program", "packs", "council")); err != nil {
		t.Fatal(err)
	}

	code, _, stderr = clitest.Run(t, bin, workDir, env, "new", "commissioning", "SQLite", "--dir", inst)
	if code != 1 {
		t.Fatalf("new commissioning exit %d want 1 stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, `unknown type "commissioning"`) {
		t.Fatalf("stderr=%q", stderr)
	}
	assertNoGh(t, stderr)
	assertNoHomeTouch(t, home)
}

func assertGeneratedFile(t *testing.T, pattern string, wants []string) {
	t.Helper()
	matches, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 {
		t.Fatalf("matches for %q: %v", pattern, matches)
	}
	data, err := os.ReadFile(matches[0])
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	for _, want := range wants {
		if !strings.Contains(text, want) {
			t.Fatalf("%s missing %q:\n%s", matches[0], want, text)
		}
	}
}
