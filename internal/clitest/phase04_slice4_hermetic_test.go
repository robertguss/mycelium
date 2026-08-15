package clitest_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/clitest"
)

func TestPhase04Slice4ScaffoldEmitsPackSkills(t *testing.T) {
	bin := clitest.Bin(t)
	stubDir := installNetworkStubs(t)
	home := t.TempDir()
	workDir := t.TempDir()
	env := hermeticEnv(stubDir, home)

	code, _, stderr := clitest.Run(t, bin, workDir, env,
		"new", "idea", "Slice4 Skills", "--offline", "--dir", filepath.Join(workDir, "slice4-skills"))
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr)
	}
	assertNoGh(t, stderr)
	inst := filepath.Join(workDir, "slice4-skills")

	if _, err := os.Stat(filepath.Join(inst, "program", "packs", "council")); err != nil {
		t.Fatalf("council pack missing: %v", err)
	}
	for _, skill := range []string{"council", "second-opinion"} {
		path := filepath.Join(inst, ".agents", "skills", skill, "SKILL.md")
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("missing pack skill %s: %v", skill, err)
		}
	}

	agents, err := os.ReadFile(filepath.Join(inst, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	agentsBody := string(agents)
	if !strings.Contains(agentsBody, "cannot fan out skip rungs 2–3") &&
		!strings.Contains(agentsBody, "cannot fan out skip rungs 2-3") {
		t.Fatalf("AGENTS.md missing capability note: %s", agentsBody)
	}
	if !strings.Contains(agentsBody, "sparring still applies") {
		t.Fatalf("AGENTS.md missing sparring clause: %s", agentsBody)
	}
	if !strings.Contains(agentsBody, "council") || !strings.Contains(agentsBody, "second-opinion") {
		t.Fatalf("AGENTS.md must name council and second-opinion skills")
	}
	if strings.Contains(agentsBody, "mycelium council") ||
		strings.Contains(agentsBody, "mycelium ladder") ||
		strings.Contains(agentsBody, "mycelium replicate") {
		t.Fatalf("AGENTS.md must not list missing verbs as commands")
	}

	councilBody := readSkill(t, inst, "council")
	secondBody := readSkill(t, inst, "second-opinion")
	for _, credit := range []string{
		"DEC-008",
		"program/contracts/replication-reconciliation.md",
		"karpathy/llm-council",
		"pstack",
	} {
		if !strings.Contains(councilBody, credit) {
			t.Fatalf("council skill missing credit %q", credit)
		}
	}
	for _, credit := range []string{"DEC-008", "pstack"} {
		if !strings.Contains(secondBody, credit) {
			t.Fatalf("second-opinion skill missing credit %q", credit)
		}
	}
	assertNoCouncilVerbToRun(t, councilBody)
	assertNoCouncilVerbToRun(t, secondBody)

	code, _, stderr = clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("fresh scaffold check exit %d stderr=%q", code, stderr)
	}
	assertNoGh(t, stderr)
	assertNoHomeTouch(t, home)
}

func TestPhase04Slice4NoRetrofitPackSkills(t *testing.T) {
	bin := clitest.Bin(t)
	stubDir := installNetworkStubs(t)
	home := t.TempDir()
	workDir := t.TempDir()
	env := hermeticEnv(stubDir, home)

	code, _, stderr := clitest.Run(t, bin, workDir, env,
		"new", "idea", "Slice4 Retrofit", "--offline", "--dir", filepath.Join(workDir, "slice4-retrofit"))
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr)
	}
	inst := filepath.Join(workDir, "slice4-retrofit")
	for _, skill := range []string{"council", "second-opinion"} {
		if err := os.RemoveAll(filepath.Join(inst, ".agents", "skills", skill)); err != nil {
			t.Fatal(err)
		}
	}
	assertSkillsAbsent(t, inst, "council", "second-opinion")

	code, _, stderr = clitest.Run(t, bin, workDir, env, "index", "--dir", inst)
	if code != 0 {
		t.Fatalf("index exit %d stderr=%q", code, stderr)
	}
	assertSkillsAbsent(t, inst, "council", "second-opinion")

	code, _, stderr = clitest.Run(t, bin, workDir, env, "tier", "standard", "--dir", inst)
	if code != 0 {
		t.Fatalf("tier exit %d stderr=%q", code, stderr)
	}
	assertSkillsAbsent(t, inst, "council", "second-opinion")

	assertNoGh(t, stderr)
	assertNoHomeTouch(t, home)
}

func TestPhase04Slice4CouncilVerbUnknown(t *testing.T) {
	bin := clitest.Bin(t)
	stubDir := installNetworkStubs(t)
	home := t.TempDir()
	workDir := t.TempDir()
	env := hermeticEnv(stubDir, home)

	code, _, stderr := clitest.Run(t, bin, workDir, env, "council")
	if code == 0 {
		t.Fatalf("mycelium council must be unknown; got exit 0 stderr=%q", stderr)
	}
	assertNoGh(t, stderr)
	assertNoHomeTouch(t, home)
}

func readSkill(t *testing.T, inst, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(inst, ".agents", "skills", name, "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func assertNoCouncilVerbToRun(t *testing.T, body string) {
	t.Helper()
	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.Contains(trimmed, "mycelium council") {
			continue
		}
		lower := strings.ToLower(trimmed)
		if strings.Contains(lower, "do not") ||
			strings.Contains(lower, "does not exist") ||
			strings.Contains(lower, "those verbs do not exist") ||
			strings.Contains(lower, "no portable") {
			continue
		}
		if strings.HasPrefix(trimmed, "- Run `mycelium council`") {
			continue
		}
		t.Fatalf("skill tells agent to run mycelium council: %q", trimmed)
	}
}
