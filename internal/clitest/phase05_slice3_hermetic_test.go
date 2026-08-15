package clitest_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/robertguss/mycelium/internal/clitest"
	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/execrun"
)

// PHASE-05 Slice 3 — DEC-011 old-manifest tolerance (G0–G3).
// G3 is a status-only golden: status exit 0, check exit 1. Do not assert check 0 on G3.

func legacyTestdata(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Join(filepath.Dir(file), "testdata", "legacy")
}

func overlayMycelium(t *testing.T, inst, srcToml string) {
	t.Helper()
	b, err := os.ReadFile(srcToml)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inst, "mycelium.toml"), b, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestPhase05Slice3G0CheckAndStatus(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}

	inst := filepath.Join(work, "g0")
	code, _, stderr := runCLI(t, clk, rec, env, work,
		"new", "idea", "Legacy Control", "--offline", "--dir", inst)
	if code != 0 {
		t.Fatalf("scaffold g0 exit %d stderr=%q", code, stderr)
	}
	overlayMycelium(t, inst, filepath.Join(legacyTestdata(t), "g0", "mycelium.toml"))

	code, _, stderr = runCLI(t, clk, rec, env, work, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("g0 check exit %d stderr=%q", code, stderr)
	}
	code, stdout, stderr := runCLI(t, clk, rec, env, work, "status", "--offline", "--dir", inst)
	if code != 0 {
		t.Fatalf("g0 status exit %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stdout, "slug: legacy-control") {
		t.Fatalf("g0 status stdout=%q", stdout)
	}
	clitest.AssertNoNetwork(t, rec)
}

func TestPhase05Slice3G1CheckAndStatus(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}

	inst := filepath.Join(work, "g1")
	code, _, stderr := runCLI(t, clk, rec, env, work,
		"new", "idea", "Legacy One", "--offline", "--dir", inst)
	if code != 0 {
		t.Fatalf("scaffold g1 exit %d stderr=%q", code, stderr)
	}
	td := legacyTestdata(t)
	overlayMycelium(t, inst, filepath.Join(td, "g1", "mycelium.toml"))
	contractSrc := filepath.Join(td, "g1", "program", "contracts", "manifest.md")
	contractDst := filepath.Join(inst, "program", "contracts", "manifest.md")
	b, err := os.ReadFile(contractSrc)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(contractDst, b, 0o644); err != nil {
		t.Fatal(err)
	}

	code, _, stderr = runCLI(t, clk, rec, env, work, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("g1 check exit %d stderr=%q (missing github_repo must pass with frozen contract)", code, stderr)
	}
	code, stdout, stderr := runCLI(t, clk, rec, env, work, "status", "--offline", "--dir", inst)
	if code != 0 {
		t.Fatalf("g1 status exit %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stdout, "slug: legacy-one") {
		t.Fatalf("g1 status stdout=%q", stdout)
	}
	clitest.AssertNoNetwork(t, rec)
}

func TestPhase05Slice3G2StatusAllPartial(t *testing.T) {
	work := t.TempDir()
	root := filepath.Join(work, "ideas")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	g2 := filepath.Join(root, "g2-legacy-master")
	if err := os.MkdirAll(g2, 0o755); err != nil {
		t.Fatal(err)
	}
	src := filepath.Join(legacyTestdata(t), "g2", "research-program.toml")
	b, err := os.ReadFile(src)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(g2, "research-program.toml"), b, 0o644); err != nil {
		t.Fatal(err)
	}

	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}

	code, stdout, stderr := runCLI(t, clk, rec, env, work,
		"status", "--all", "--offline", "--root", root)
	if code != 0 {
		t.Fatalf("g2 status --all exit %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stdout, "partial: legacy-manifest (") {
		t.Fatalf("missing partial: %q", stdout)
	}
	if !strings.Contains(stdout, "research-program.toml without mycelium.toml") {
		t.Fatalf("missing g2 reason: %q", stdout)
	}
	if strings.Contains(stdout, "g2-legacy-master\t") || strings.Contains(stdout, "\tlegacy") {
		// no idea row for the research-program-only directory
		for _, line := range strings.Split(stdout, "\n") {
			if strings.Contains(line, "\t") && strings.Contains(line, "g2") {
				t.Fatalf("g2 listed as idea row: %q", line)
			}
		}
	}
	if !strings.Contains(stdout, "0 ideas (") {
		t.Fatalf("want zero ideas, got %q", stdout)
	}
	clitest.AssertNoNetwork(t, rec)
}

func TestPhase05Slice3G3StatusOnlyNotCheckPass(t *testing.T) {
	// G3 is a status-only golden, not a check-pass.
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}

	inst := filepath.Join(work, "g3")
	code, _, stderr := runCLI(t, clk, rec, env, work,
		"new", "idea", "Legacy Extra", "--offline", "--dir", inst)
	if code != 0 {
		t.Fatalf("scaffold g3 exit %d stderr=%q", code, stderr)
	}
	overlayMycelium(t, inst, filepath.Join(legacyTestdata(t), "g3", "mycelium.toml"))

	code, stdout, stderr := runCLI(t, clk, rec, env, work, "status", "--offline", "--dir", inst)
	if code != 0 {
		t.Fatalf("g3 status exit %d stderr=%q (status-only golden must pass)", code, stderr)
	}
	if !strings.Contains(stdout, "slug: legacy-extra") {
		t.Fatalf("g3 status stdout=%q", stdout)
	}

	code, _, stderr = runCLI(t, clk, rec, env, work, "check", "--dir", inst)
	if code != 1 {
		t.Fatalf("g3 check exit %d want 1 (unknown key refuse; not a check-pass golden) stderr=%q", code, stderr)
	}
	clitest.AssertNoNetwork(t, rec)
}

func TestPhase05Slice3MixedRootNoAbort(t *testing.T) {
	work := t.TempDir()
	root := filepath.Join(work, "ideas")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}
	td := legacyTestdata(t)

	// G0
	g0 := filepath.Join(root, "g0")
	code, _, stderr := runCLI(t, clk, rec, env, work,
		"new", "idea", "Legacy Control", "--offline", "--dir", g0)
	if code != 0 {
		t.Fatalf("g0 scaffold exit %d stderr=%q", code, stderr)
	}
	overlayMycelium(t, g0, filepath.Join(td, "g0", "mycelium.toml"))

	// G1
	g1 := filepath.Join(root, "g1")
	code, _, stderr = runCLI(t, clk, rec, env, work,
		"new", "idea", "Legacy One", "--offline", "--dir", g1)
	if code != 0 {
		t.Fatalf("g1 scaffold exit %d stderr=%q", code, stderr)
	}
	overlayMycelium(t, g1, filepath.Join(td, "g1", "mycelium.toml"))
	b, err := os.ReadFile(filepath.Join(td, "g1", "program", "contracts", "manifest.md"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(g1, "program", "contracts", "manifest.md"), b, 0o644); err != nil {
		t.Fatal(err)
	}

	// G2
	g2 := filepath.Join(root, "g2")
	if err := os.MkdirAll(g2, 0o755); err != nil {
		t.Fatal(err)
	}
	b, err = os.ReadFile(filepath.Join(td, "g2", "research-program.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(g2, "research-program.toml"), b, 0o644); err != nil {
		t.Fatal(err)
	}

	// G3
	g3 := filepath.Join(root, "g3")
	code, _, stderr = runCLI(t, clk, rec, env, work,
		"new", "idea", "Legacy Extra", "--offline", "--dir", g3)
	if code != 0 {
		t.Fatalf("g3 scaffold exit %d stderr=%q", code, stderr)
	}
	overlayMycelium(t, g3, filepath.Join(td, "g3", "mycelium.toml"))

	// unreadable
	bad := filepath.Join(root, "unreadable")
	if err := os.MkdirAll(bad, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bad, "mycelium.toml"), []byte("{{{not-toml"), 0o644); err != nil {
		t.Fatal(err)
	}

	code, stdout, stderr := runCLI(t, clk, rec, env, work,
		"status", "--all", "--offline", "--root", root)
	if code != 0 {
		t.Fatalf("mixed status --all exit %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stdout, "legacy-control\t") {
		t.Fatalf("g0 not listed: %q", stdout)
	}
	if !strings.Contains(stdout, "legacy-one\t") {
		t.Fatalf("g1 not listed: %q", stdout)
	}
	if !strings.Contains(stdout, "legacy-extra\t") {
		t.Fatalf("g3 not listed: %q", stdout)
	}
	if !strings.Contains(stdout, "partial: legacy-manifest (") {
		t.Fatalf("missing legacy-manifest: %q", stdout)
	}
	if !strings.Contains(stdout, "research-program.toml without mycelium.toml") {
		t.Fatalf("missing g2 reason: %q", stdout)
	}
	// count idea rows (tab-separated, not partial/summary)
	ideas := 0
	for _, line := range strings.Split(stdout, "\n") {
		if strings.Count(line, "\t") >= 4 && !strings.HasPrefix(line, "partial:") {
			ideas++
		}
	}
	if ideas != 3 {
		t.Fatalf("want 3 idea rows (G0/G1/G3), got %d\n%s", ideas, stdout)
	}
	if !strings.Contains(stdout, "3 ideas (") {
		t.Fatalf("summary: %q", stdout)
	}
	clitest.AssertNoNetwork(t, rec)
	if rec.Called("gh") {
		t.Fatal("gh was called")
	}
}

func TestPhase05Slice3HandedOffStillUnreachable(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}
	inst := scaffoldIdea(t, work, "Handed Off Guard", clk, rec, env)
	code, _, stderr := runCLI(t, clk, rec, env, work, "state", "handed-off", "--dir", inst)
	if code != 1 {
		t.Fatalf("handed-off exit %d want 1 stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "handed-off") {
		t.Fatalf("stderr=%q", stderr)
	}
}
