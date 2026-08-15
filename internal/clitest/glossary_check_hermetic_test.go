package clitest_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/clitest"
)

func TestSlice3GlossaryH1OnlyExits0(t *testing.T) {
	bin := clitest.Bin(t)
	workDir := t.TempDir()
	env := []string{"MYCELIUM_OFFLINE=1"}

	code, _, stderr := clitest.Run(t, bin, workDir, env, "new", "idea", "Slice3 H1", "--offline")
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr)
	}
	inst := filepath.Join(workDir, "slice3-h1")
	if err := os.WriteFile(filepath.Join(inst, "CONTEXT.md"), []byte("# Glossary\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	code, _, stderr = clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("check exit %d want 0 stderr=%q", code, stderr)
	}
}

func TestSlice3GlossaryMissingDefinitionExits1(t *testing.T) {
	bin := clitest.Bin(t)
	workDir := t.TempDir()
	env := []string{"MYCELIUM_OFFLINE=1"}

	code, _, stderr := clitest.Run(t, bin, workDir, env, "new", "idea", "Slice3 Miss Def", "--offline")
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr)
	}
	inst := filepath.Join(workDir, "slice3-miss-def")
	if err := os.WriteFile(filepath.Join(inst, "CONTEXT.md"), []byte("# Glossary\n\n## SQLite\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	code, _, stderr = clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 1 {
		t.Fatalf("check exit %d want 1 stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, `CONTEXT.md term "SQLite" missing ### Definition`) {
		t.Fatalf("stderr missing term message: %q", stderr)
	}
	if !strings.Contains(stderr, "program/contracts/glossary.md") {
		t.Fatalf("stderr missing glossary.md: %q", stderr)
	}
	if !strings.Contains(stderr, "convention: glossary") {
		t.Fatalf("stderr missing convention: %q", stderr)
	}
	if !strings.Contains(stderr, "add ### Definition under ## SQLite") {
		t.Fatalf("stderr missing fix: %q", stderr)
	}
}

func TestSlice3GlossaryFillDefinitionExits0(t *testing.T) {
	bin := clitest.Bin(t)
	workDir := t.TempDir()
	env := []string{"MYCELIUM_OFFLINE=1"}

	code, _, stderr := clitest.Run(t, bin, workDir, env, "new", "idea", "Slice3 Fill Def", "--offline")
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr)
	}
	inst := filepath.Join(workDir, "slice3-fill-def")
	body := "# Glossary\n\n## SQLite\n\n### Definition\n\n<!-- fill -->\n"
	if err := os.WriteFile(filepath.Join(inst, "CONTEXT.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	code, _, stderr = clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("check exit %d want 0 stderr=%q", code, stderr)
	}
}

func TestSlice3GlossaryMissingFileExits0(t *testing.T) {
	bin := clitest.Bin(t)
	workDir := t.TempDir()
	env := []string{"MYCELIUM_OFFLINE=1"}

	code, _, stderr := clitest.Run(t, bin, workDir, env, "new", "idea", "Slice3 No Ctx", "--offline")
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr)
	}
	inst := filepath.Join(workDir, "slice3-no-ctx")
	tierPath := filepath.Join(inst, "program", "tiers", "focused.toml")
	tb, err := os.ReadFile(tierPath)
	if err != nil {
		t.Fatal(err)
	}
	patched := strings.ReplaceAll(string(tb), `"CONTEXT.md", `, "")
	patched = strings.ReplaceAll(patched, `, "CONTEXT.md"`, "")
	patched = strings.ReplaceAll(patched, `"CONTEXT.md"`, "")
	if err := os.WriteFile(tierPath, []byte(patched), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(inst, "CONTEXT.md")); err != nil {
		t.Fatal(err)
	}

	code, _, stderr = clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("check exit %d want 0 stderr=%q", code, stderr)
	}
}
