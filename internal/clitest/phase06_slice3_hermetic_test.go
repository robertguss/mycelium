package clitest_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/robertguss/mycelium/internal/clitest"
	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/execrun"
)

// PHASE-06 Slice 3 — pstack/poteto bridge docs emit on --offline scaffold.
// Asserts AGENTS.md Implementation systems, reference mapping file, and
// mycelium-cli skill naming mycelium handoff + handoff/PACKET.md.

func TestPhase06Slice3OfflineScaffoldBridge(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1", "MYCELIUM_NOW": "2026-08-15T00:00:00Z"}

	inst := filepath.Join(work, "slice3-bridge")
	code, _, stderr := runCLI(t, clk, rec, env, work,
		"new", "idea", "Slice 3 Bridge", "--offline", "--dir", inst)
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr)
	}

	agents := string(readFile(t, filepath.Join(inst, "AGENTS.md")))
	if !strings.Contains(agents, "## Implementation systems") {
		t.Fatalf("AGENTS.md missing Implementation systems heading:\n%s", agents)
	}
	for _, claim := range []string{
		"ONLY",
		"handoff/",
		"No chat history",
		"pstack/poteto-mode",
		"manual",
		"floor",
	} {
		if !strings.Contains(agents, claim) {
			t.Fatalf("AGENTS.md missing isolation claim %q", claim)
		}
	}
	if !strings.Contains(agents, "mycelium handoff") {
		t.Fatalf("AGENTS.md commands missing mycelium handoff")
	}
	if !strings.Contains(agents, "state handed-off") || !strings.Contains(agents, "packet already exists") {
		t.Fatalf("AGENTS.md missing state handed-off / packet already exists guidance")
	}
	if strings.Contains(agents, "mycelium council") {
		t.Fatalf("AGENTS.md must not name mycelium council as a command to run")
	}

	refPath := filepath.Join(inst, "program", "reference", "implementation-systems.md")
	ref := string(readFile(t, refPath))
	for _, want := range []string{
		"Framing",
		"Locked decisions",
		"Glossary",
		"Open questions",
		"Evidence summary",
		"Implementation playbooks",
		"Implementation system",
		"Time budget",
		"Acceptance",
		"pstack/poteto",
		"manual",
		"30m",
		"ONLY",
		"handoff/",
		"No chat history",
	} {
		if !strings.Contains(ref, want) {
			t.Fatalf("implementation-systems.md missing %q", want)
		}
	}

	skill := string(readFile(t, filepath.Join(inst, ".agents", "skills", "mycelium-cli", "SKILL.md")))
	if !strings.Contains(skill, "mycelium handoff") {
		t.Fatalf("mycelium-cli skill missing mycelium handoff")
	}
	if !strings.Contains(skill, "handoff/PACKET.md") {
		t.Fatalf("mycelium-cli skill missing handoff/PACKET.md")
	}
	if !strings.Contains(skill, "handed-off") || !strings.Contains(skill, "archived") {
		t.Fatalf("mycelium-cli skill missing terminal handed-off except archived")
	}
	if strings.Contains(strings.ToLower(skill), "mycelium council") {
		t.Fatalf("mycelium-cli skill must not name mycelium council as a command to run")
	}

	// No pstack binary / wrapper in scaffold.
	if _, err := os.Stat(filepath.Join(inst, "pstack")); err == nil {
		t.Fatal("scaffold must not emit a pstack binary")
	}
	if _, err := os.Stat(filepath.Join(inst, "bin", "pstack")); err == nil {
		t.Fatal("scaffold must not emit bin/pstack")
	}

	code, _, stderr = runCLI(t, clk, rec, env, work, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("check after scaffold exit %d stderr=%q", code, stderr)
	}
	clitest.AssertNoNetwork(t, rec)
}
