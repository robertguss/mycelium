package clitest_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/clitest"
	"github.com/robertguss/mycelium/internal/execrun"
)

func TestScaffoldEmitsLifecycleSkills(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}

	inst := scaffoldIdea(t, work, "Skill Emit", clk, rec, env)
	clitest.AssertNoNetwork(t, rec)

	for _, skill := range []string{"mycelium-cli", "spark", "wake", "portfolio"} {
		path := filepath.Join(inst, ".agents", "skills", skill, "SKILL.md")
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("missing %s: %v", path, err)
		}
	}
}

func TestMutatingCommandsDoNotRetrofitSkills(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}

	inst := scaffoldIdea(t, work, "No Retrofit", clk, rec, env)
	for _, skill := range []string{"spark", "wake", "portfolio"} {
		if err := os.RemoveAll(filepath.Join(inst, ".agents", "skills", skill)); err != nil {
			t.Fatal(err)
		}
	}
	assertSkillsAbsent(t, inst, "spark", "wake", "portfolio")

	code, _, stderr := runCLI(t, clk, rec, env, work, "index", "--dir", inst)
	if code != 0 {
		t.Fatalf("index exit %d stderr=%q", code, stderr)
	}
	assertSkillsAbsent(t, inst, "spark", "wake", "portfolio")

	code, _, stderr = runCLI(t, clk, rec, env, work, "tier", "standard", "--dir", inst)
	if code != 0 {
		t.Fatalf("tier exit %d stderr=%q", code, stderr)
	}
	assertSkillsAbsent(t, inst, "spark", "wake", "portfolio")

	code, _, stderr = runCLI(t, clk, rec, env, work, "state", "exploring", "--dir", inst)
	if code != 0 {
		t.Fatalf("state exploring exit %d stderr=%q", code, stderr)
	}
	assertSkillsAbsent(t, inst, "spark", "wake", "portfolio")

	clitest.AssertNoNetwork(t, rec)
}

func assertSkillsAbsent(t *testing.T, inst string, skills ...string) {
	t.Helper()
	for _, skill := range skills {
		path := filepath.Join(inst, ".agents", "skills", skill)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("skill %s should stay absent after mutate (err=%v)", skill, err)
		}
	}
}
