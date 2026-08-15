package clitest_test

import (
	"os/exec"
	"strings"
	"testing"
)

// TestROB557ManualGuideRgCommands locks Quality guide items 2 and 7.
// Execs the exact argv from the ROB-557 PR manual test guide (fixed form).
func TestROB557ManualGuideRgCommands(t *testing.T) {
	rgPath, err := exec.LookPath("rg")
	if err != nil {
		t.Skip("rg not on PATH")
	}
	root := repoRoot(t)

	t.Run("item2_forbidden_strings", func(t *testing.T) {
		// Exact guide argv: rg -n -F -e '{{PROJECT_NAME}}' -e 'just init' -e 'Use this template' README.md AGENTS.md
		cmd := exec.Command(rgPath, "-n", "-F",
			"-e", "{{PROJECT_NAME}}",
			"-e", "just init",
			"-e", "Use this template",
			"README.md", "AGENTS.md",
		)
		cmd.Dir = root
		out, err := cmd.CombinedOutput()
		code := 0
		if err != nil {
			ee, ok := err.(*exec.ExitError)
			if !ok {
				t.Fatalf("rg item2: unexpected error %v\n%s", err, out)
			}
			code = ee.ExitCode()
		}
		if code == 2 {
			t.Fatalf("rg item2: parse/error exit 2 (must not)\n%s", out)
		}
		if code != 0 && code != 1 {
			t.Fatalf("rg item2: unexpected exit %d\n%s", code, out)
		}
		if code == 0 || len(strings.TrimSpace(string(out))) != 0 {
			t.Fatalf("rg item2: expected no matches (exit 1 or empty stdout), got exit %d\n%s", code, out)
		}
	})

	t.Run("item7_manifest_rewrite", func(t *testing.T) {
		// Exact guide argv (multiline -U):
		// rg -n -U 'does not carry[[:space:]]+`research-program.toml`; idea instances use `mycelium.toml` \(DEC-012\)' \
		//   program/contracts/manifest.md internal/embed/program/contracts/manifest.md
		pattern := "does not carry[[:space:]]+`research-program.toml`; idea instances use `mycelium.toml` \\(DEC-012\\)"
		cmd := exec.Command(rgPath, "-n", "-U", pattern,
			"program/contracts/manifest.md",
			"internal/embed/program/contracts/manifest.md",
		)
		cmd.Dir = root
		out, err := cmd.CombinedOutput()
		code := 0
		if err != nil {
			ee, ok := err.(*exec.ExitError)
			if !ok {
				t.Fatalf("rg item7: unexpected error %v\n%s", err, out)
			}
			code = ee.ExitCode()
		}
		if code != 0 {
			t.Fatalf("rg item7: expected exit 0 with matches in both files, got %d\n%s", code, out)
		}
		text := string(out)
		if !strings.Contains(text, "program/contracts/manifest.md") {
			t.Fatalf("rg item7: missing match for program/contracts/manifest.md\n%s", text)
		}
		if !strings.Contains(text, "internal/embed/program/contracts/manifest.md") {
			t.Fatalf("rg item7: missing match for internal/embed/program/contracts/manifest.md\n%s", text)
		}
	})
}
