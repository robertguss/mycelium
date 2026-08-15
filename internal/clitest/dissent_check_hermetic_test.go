package clitest_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/clitest"
)

func TestSlice4DECWithoutDissentExits0(t *testing.T) {
	bin := clitest.Bin(t)
	workDir := t.TempDir()
	env := []string{"MYCELIUM_OFFLINE=1"}

	code, _, stderr := clitest.Run(t, bin, workDir, env, "new", "idea", "Slice4 No Dissent", "--offline")
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr)
	}
	inst := filepath.Join(workDir, "slice4-no-dissent")
	code, _, stderr = clitest.Run(t, bin, workDir, env, "new", "decision", "Ship it", "--dir", inst)
	if code != 0 {
		t.Fatalf("new decision exit %d stderr=%q", code, stderr)
	}
	decPath := filepath.Join(inst, "decisions", "DEC-001-ship-it.md")
	b, err := os.ReadFile(decPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), "\n## Dissent\n") {
		t.Fatalf("template must not emit live ## Dissent heading:\n%s", b)
	}
	if !strings.Contains(string(b), "<!-- Optional section (not required): ## Dissent") {
		t.Fatalf("template missing Dissent HTML comment hint:\n%s", b)
	}

	code, _, stderr = clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("check exit %d want 0 stderr=%q", code, stderr)
	}
}

func TestSlice4DECDissentCitingRealOQExits0(t *testing.T) {
	bin := clitest.Bin(t)
	workDir := t.TempDir()
	env := []string{"MYCELIUM_OFFLINE=1"}

	code, _, stderr := clitest.Run(t, bin, workDir, env, "new", "idea", "Slice4 Dissent OK", "--offline")
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr)
	}
	inst := filepath.Join(workDir, "slice4-dissent-ok")
	code, _, stderr = clitest.Run(t, bin, workDir, env, "new", "question", "Use SQLite", "--dir", inst)
	if code != 0 {
		t.Fatalf("new question exit %d stderr=%q", code, stderr)
	}
	code, _, stderr = clitest.Run(t, bin, workDir, env, "new", "decision", "Ship it", "--dir", inst)
	if code != 0 {
		t.Fatalf("new decision exit %d stderr=%q", code, stderr)
	}
	decPath := filepath.Join(inst, "decisions", "DEC-001-ship-it.md")
	b, err := os.ReadFile(decPath)
	if err != nil {
		t.Fatal(err)
	}
	patched := string(b) + "\n## Dissent\n\nStill disagree; see OQ-001.\n"
	if err := os.WriteFile(decPath, []byte(patched), 0o644); err != nil {
		t.Fatal(err)
	}

	code, _, stderr = clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("check exit %d want 0 stderr=%q", code, stderr)
	}
}

func TestSlice4DECDissentNoTokenExits1(t *testing.T) {
	bin := clitest.Bin(t)
	workDir := t.TempDir()
	env := []string{"MYCELIUM_OFFLINE=1"}

	code, _, stderr := clitest.Run(t, bin, workDir, env, "new", "idea", "Slice4 Dissent Bare", "--offline")
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr)
	}
	inst := filepath.Join(workDir, "slice4-dissent-bare")
	code, _, stderr = clitest.Run(t, bin, workDir, env, "new", "decision", "Ship it", "--dir", inst)
	if code != 0 {
		t.Fatalf("new decision exit %d stderr=%q", code, stderr)
	}
	decPath := filepath.Join(inst, "decisions", "DEC-001-ship-it.md")
	b, err := os.ReadFile(decPath)
	if err != nil {
		t.Fatal(err)
	}
	patched := string(b) + "\n## Dissent\n\nI object.\n"
	if err := os.WriteFile(decPath, []byte(patched), 0o644); err != nil {
		t.Fatal(err)
	}

	code, _, stderr = clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 1 {
		t.Fatalf("check exit %d want 1 stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "DEC-001 ## Dissent has no resolvable OQ-### or ASM-###") {
		t.Fatalf("stderr missing what: %q", stderr)
	}
	if !strings.Contains(stderr, "convention: dissent") {
		t.Fatalf("stderr missing convention: %q", stderr)
	}
	if !strings.Contains(stderr, "program/contracts/sparring.md") {
		t.Fatalf("stderr missing contract: %q", stderr)
	}
	if !strings.Contains(stderr, "cite an existing OQ-### or ASM-### in ## Dissent, or remove the heading") {
		t.Fatalf("stderr missing fix: %q", stderr)
	}
}

func TestSlice4DECDissentOnlyDECExits1(t *testing.T) {
	bin := clitest.Bin(t)
	workDir := t.TempDir()
	env := []string{"MYCELIUM_OFFLINE=1"}

	code, _, stderr := clitest.Run(t, bin, workDir, env, "new", "idea", "Slice4 Dissent DEC", "--offline")
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr)
	}
	inst := filepath.Join(workDir, "slice4-dissent-dec")
	code, _, stderr = clitest.Run(t, bin, workDir, env, "new", "decision", "Ship it", "--dir", inst)
	if code != 0 {
		t.Fatalf("new decision exit %d stderr=%q", code, stderr)
	}
	decPath := filepath.Join(inst, "decisions", "DEC-001-ship-it.md")
	b, err := os.ReadFile(decPath)
	if err != nil {
		t.Fatal(err)
	}
	patched := string(b) + "\n## Dissent\n\nSee DEC-001 for context.\n"
	if err := os.WriteFile(decPath, []byte(patched), 0o644); err != nil {
		t.Fatal(err)
	}

	code, _, stderr = clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 1 {
		t.Fatalf("check exit %d want 1 stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "DEC-001 ## Dissent has no resolvable OQ-### or ASM-###") {
		t.Fatalf("stderr missing dissent what: %q", stderr)
	}
}
