package clitest_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/clitest"
)

func TestIndexRepairAndNewDecisionCounts(t *testing.T) {
	bin := clitest.Bin(t)
	workDir := t.TempDir()
	env := []string{"MYCELIUM_OFFLINE=1", "CGO_ENABLED=0"}

	code, _, stderr := clitest.Run(t, bin, workDir, env, "new", "idea", "Index Slice", "--offline")
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr)
	}
	inst := filepath.Join(workDir, "index-slice")

	indexPath := filepath.Join(inst, "index.md")
	body, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, h2 := range []string{"## State", "## Artifacts", "## Log tail", "## Wake"} {
		if !strings.Contains(string(body), h2) {
			t.Fatalf("scaffold index missing %s", h2)
		}
	}

	code, stdout, stderr := clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("check after scaffold exit %d stderr=%q stdout=%q", code, stderr, stdout)
	}

	if err := os.Remove(indexPath); err != nil {
		t.Fatal(err)
	}
	code, _, stderr = clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 1 {
		t.Fatalf("check after delete index: exit %d want 1 stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "mycelium index") {
		t.Fatalf("want fix mycelium index, got %q", stderr)
	}

	code, stdout, stderr = clitest.Run(t, bin, workDir, env, "index", "--dir", inst)
	if code != 0 {
		t.Fatalf("index exit %d stderr=%q", code, stderr)
	}
	if strings.TrimSpace(stdout) != "wrote index.md" {
		t.Fatalf("stdout=%q", stdout)
	}
	code, _, stderr = clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("check after index repair exit %d stderr=%q", code, stderr)
	}

	code, _, stderr = clitest.Run(t, bin, workDir, env, "new", "decision", "First thought", "--dir", inst)
	if code != 0 {
		t.Fatalf("new decision exit %d stderr=%q", code, stderr)
	}
	body, err = os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "DEC: 1\n") {
		t.Fatalf("want DEC: 1 after new decision:\n%s", body)
	}

	code, stdout, stderr = clitest.Run(t, bin, workDir, env, "index", "--dir", inst)
	if code != 0 {
		t.Fatalf("index idempotent exit %d stderr=%q", code, stderr)
	}
	if strings.TrimSpace(stdout) != "wrote index.md" {
		t.Fatalf("idempotent stdout=%q", stdout)
	}
}
