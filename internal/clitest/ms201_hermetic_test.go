package clitest_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/clitest"
)

// TestMS201HermeticFixture is the PHASE-02 MS-201 gate (Appendix E / §13).
func TestMS201HermeticFixture(t *testing.T) {
	bin := clitest.Bin(t)
	stubDir, ghMarker := installGhNeverStub(t)
	work := t.TempDir()

	run := func(now string, args ...string) (int, string, string) {
		t.Helper()
		return clitest.Run(t, bin, work, ms201Env(stubDir, now), args...)
	}

	const t0 = "2026-08-01T00:00:00Z"

	code, _, stderr := run(t0, "new", "idea", "Wake Fixture", "--offline", "--dir", filepath.Join(work, "wake-fixture"))
	if code != 0 {
		t.Fatalf("new idea: exit %d stderr=%q", code, stderr)
	}
	inst := filepath.Join(work, "wake-fixture")

	m := loadManifest(t, inst)
	if m.Slug != "wake-fixture" || m.State != "spark" || m.CreatedDate != "2026-08-01" {
		t.Fatalf("after scaffold: slug=%q state=%q created_date=%q", m.Slug, m.State, m.CreatedDate)
	}

	code, _, stderr = run(t0, "state", "exploring", "--dir", inst)
	if code != 0 {
		t.Fatalf("state exploring: exit %d stderr=%q", code, stderr)
	}

	code, _, stderr = run(t0, "new", "decision", "Park the idea", "--dir", inst)
	if code != 0 {
		t.Fatalf("new decision: exit %d stderr=%q", code, stderr)
	}
	if _, err := os.Stat(filepath.Join(inst, "decisions", "DEC-001-park-the-idea.md")); err != nil {
		t.Fatalf("DEC-001 missing: %v", err)
	}

	code, _, stderr = run(t0, "new", "assumption", "API stays stable", "--dir", inst)
	if code != 0 {
		t.Fatalf("new assumption ASM-001: exit %d stderr=%q", code, stderr)
	}
	patchASM(t, inst, "ASM-001", "Held", "2026-08-05")

	code, _, stderr = run(t0, "new", "evidence", "Vendor changelog", "--dir", inst)
	if code != 0 {
		t.Fatalf("new evidence: exit %d stderr=%q", code, stderr)
	}
	patchEVD(t, inst, "EVD-001", "2026-08-06")

	code, _, stderr = run(t0, "new", "assumption", "Budget is unlimited", "--dir", inst)
	if code != 0 {
		t.Fatalf("new assumption ASM-002: exit %d stderr=%q", code, stderr)
	}
	patchASM(t, inst, "ASM-002", "Retired", "")

	code, _, stderr = run(t0, "state", "simmering", "--revisit", "2026-08-08", "--dir", inst)
	if code != 0 {
		t.Fatalf("state simmering: exit %d stderr=%q", code, stderr)
	}

	early := filepath.Join(work, "wake-fixture-early")
	copyTree(t, inst, early)
	const earlyNow = "2026-08-07T00:00:00Z"

	code, stdout, stderr := run(earlyNow, "status", "--dir", early)
	if code != 0 {
		t.Fatalf("early status: exit %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stdout, "due: no") {
		t.Fatalf("early status want due: no, got %q", stdout)
	}

	code, _, stderr = run(earlyNow, "wake", "--dir", early)
	if code != 0 {
		t.Fatalf("early wake: exit %d stderr=%q", code, stderr)
	}
	earlyBrief := readFile(t, filepath.Join(early, "briefs", "WAKE-2026-08-07.md"))
	earlyBody := string(earlyBrief)
	if !strings.Contains(earlyBody, "ASM-001") || !strings.Contains(earlyBody, "EVD-001") {
		t.Fatalf("early wake missing ASM-001/EVD-001:\n%s", earlyBody)
	}

	const t1 = "2026-08-09T00:00:00Z"

	code, _, stderr = run(t1, "wake", "--dir", inst)
	if code != 0 {
		t.Fatalf("T1 wake: exit %d stderr=%q", code, stderr)
	}

	m = loadManifest(t, inst)
	if m.State != "exploring" {
		t.Fatalf("after T1: state=%q want exploring", m.State)
	}
	if m.Revisit != "" {
		t.Fatalf("after T1: revisit=%q want empty", m.Revisit)
	}

	datedPath := filepath.Join(inst, "briefs", "WAKE-2026-08-09.md")
	latestPath := filepath.Join(inst, "briefs", "LATEST.md")
	dated := readFile(t, datedPath)
	latest := readFile(t, latestPath)
	if !bytes.Equal(dated, latest) {
		t.Fatal("LATEST.md must equal WAKE-2026-08-09.md bytes")
	}

	body := string(dated)
	for _, h2 := range []string{
		"## Parked",
		"## Log since simmer",
		"## Evidence triggers",
		"## Assumptions",
		"## Suggested next",
	} {
		if !strings.Contains(body, h2) {
			t.Fatalf("T1 brief missing %s:\n%s", h2, body)
		}
	}
	if !strings.Contains(body, "ASM-001") || !strings.Contains(body, "EVD-001") {
		t.Fatalf("T1 brief missing ASM-001/EVD-001:\n%s", body)
	}
	if strings.Contains(body, "ASM-002") {
		t.Fatalf("T1 brief must not cite ASM-002:\n%s", body)
	}
	if !strings.Contains(body, "2026-08-01") && !strings.Contains(body, "2026-08-08") {
		t.Fatalf("T1 brief must cite simmer 2026-08-01 or revisit 2026-08-08:\n%s", body)
	}

	code, _, stderr = run(t1, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("check after T1: exit %d stderr=%q", code, stderr)
	}

	if _, err := os.Stat(ghMarker); err == nil {
		t.Fatal("gh was invoked (marker present)")
	} else if !os.IsNotExist(err) {
		t.Fatalf("gh marker stat: %v", err)
	}
}

func ms201Env(stubDir, now string) []string {
	base := stripEnvKeys(os.Environ(), "MYCELIUM_NOW", "MYCELIUM_OFFLINE", "PATH", "MYCELIUM_BIN")
	path := stubDir + string(os.PathListSeparator) + os.Getenv("PATH")
	return append(base,
		"PATH="+path,
		"MYCELIUM_OFFLINE=1",
		"MYCELIUM_NOW="+now,
	)
}

func stripEnvKeys(env []string, keys ...string) []string {
	drop := make(map[string]bool, len(keys))
	for _, k := range keys {
		drop[k] = true
	}
	out := make([]string, 0, len(env))
	for _, e := range env {
		k, _, ok := strings.Cut(e, "=")
		if ok && drop[k] {
			continue
		}
		out = append(out, e)
	}
	return out
}

func installGhNeverStub(t *testing.T) (stubDir, marker string) {
	t.Helper()
	realGit, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	marker = filepath.Join(dir, "gh-was-invoked")
	writeStub(t, filepath.Join(dir, "gh"), `#!/bin/sh
touch `+shellQuote(marker)+`
echo "gh network access forbidden" >&2
exit 99
`)
	writeStub(t, filepath.Join(dir, "git"), `#!/bin/sh
if [ "$1" = "remote" ]; then
  echo "git remote network access forbidden" >&2
  exit 99
fi
exec `+shellQuote(realGit)+` "$@"
`)
	return dir, marker
}
