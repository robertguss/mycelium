package clitest_test

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/clitest"
	"github.com/robertguss/mycelium/internal/handoff"
)

// Coverage stamp (PHASE-06 brief). Do not lower.
const ms601CoverageFloor = 85

// TestMS601HermeticMatrix is the PHASE-06 MS-601 gate (Appendix E).
// Execs the real binary. No network, no gh, no Actions, no live agent.
func TestMS601HermeticMatrix(t *testing.T) {
	bin := clitest.Bin(t)
	stubDir, ghMarker := installGhNeverStub(t)
	home := t.TempDir()
	env := append(hermeticEnv(stubDir, home),
		"MYCELIUM_NOW=2026-08-15T00:00:00Z",
		"GH_TOKEN=",
		"MYCELIUM_CONFIG="+filepath.Join(home, "config"),
	)
	td := phase06Slice4Testdata(t)

	run := func(t *testing.T, work string, args ...string) (int, string, string) {
		t.Helper()
		code, stdout, stderr := clitest.Run(t, bin, work, env, args...)
		assertNoGh(t, stderr)
		if _, err := os.Stat(ghMarker); err == nil {
			t.Fatal("gh was invoked")
		} else if !os.IsNotExist(err) {
			t.Fatalf("stat gh marker: %v", err)
		}
		assertNoHomeTouch(t, home)
		return code, stdout, stderr
	}

	t.Run("MS-601-GEN", func(t *testing.T) {
		work := t.TempDir()
		inst := filepath.Join(work, "ho")
		copyTree(t, filepath.Join(td, "handoff-fixture"), inst)

		code, stdout, stderr := run(t, work, "handoff", "--dir", inst)
		if code != 0 {
			t.Fatalf("handoff exit %d stderr=%q", code, stderr)
		}
		if !strings.Contains(stdout, "mycelium handoff: ok") {
			t.Fatalf("stdout=%q", stdout)
		}

		packetPath := filepath.Join(inst, "handoff", "PACKET.md")
		packet := string(readFile(t, packetPath))
		for _, want := range []string{
			`id = "HO-001"`,
			`date = "2026-08-15"`,
			`implementation_system = "pstack/poteto"`,
			`time_budget = "30m"`,
			"## Framing",
			"## Locked decisions",
			"## Glossary",
			"## Open questions",
			"## Evidence summary",
			"## Implementation playbooks",
			"## Implementation system",
			"## Time budget",
			"## Acceptance",
		} {
			if !strings.Contains(packet, want) {
				t.Fatalf("PACKET.md missing %q:\n%s", want, packet)
			}
		}
		findings := handoff.Check(os.DirFS(filepath.Join(inst, "handoff")))
		if len(findings) > 0 {
			t.Fatalf("packet structure: %v", findings)
		}
		for _, rel := range []string{
			"handoff/playbooks/PLAYBOOK.md",
			"handoff/acceptance/add_test.go",
			"handoff/decisions/DEC-001-add-signature.md",
			"handoff/glossary.md",
			"handoff/evidence/SUMMARY.md",
		} {
			if _, err := os.Stat(filepath.Join(inst, rel)); err != nil {
				t.Fatalf("missing %s: %v", rel, err)
			}
		}
	})

	t.Run("MS-601-CMD", func(t *testing.T) {
		work := t.TempDir()
		inst := filepath.Join(work, "ho")
		copyTree(t, filepath.Join(td, "handoff-fixture"), inst)

		code, stdout, stderr := run(t, work, "handoff", "--dir", inst)
		if code != 0 {
			t.Fatalf("handoff exit %d stderr=%q", code, stderr)
		}
		if !strings.Contains(stdout, "mycelium handoff: ok") ||
			!strings.Contains(stdout, "state: handed-off") ||
			!strings.Contains(stdout, "packet: handoff/PACKET.md") {
			t.Fatalf("stdout=%q", stdout)
		}
		if loadManifest(t, inst).State != "handed-off" {
			t.Fatalf("state=%q", loadManifest(t, inst).State)
		}
		logBody := string(readFile(t, filepath.Join(inst, "log.md")))
		wantLine := "2026-08-15\thandoff\tHO-001\tclarified -> handed-off"
		if !strings.Contains(logBody, wantLine) {
			t.Fatalf("log missing %q\n%s", wantLine, logBody)
		}
		code, _, stderr = run(t, work, "check", "--dir", inst)
		if code != 0 {
			t.Fatalf("check exit %d stderr=%q", code, stderr)
		}
	})

	t.Run("MS-601-REFUSE", func(t *testing.T) {
		work := t.TempDir()
		inst := filepath.Join(work, "nopkt")
		code, _, stderr := run(t, work, "new", "idea", "No Packet Refuse", "--offline", "--dir", inst)
		if code != 0 {
			t.Fatalf("new idea exit %d stderr=%q", code, stderr)
		}
		code, _, stderr = run(t, work, "state", "exploring", "--dir", inst)
		if code != 0 {
			t.Fatalf("state exploring exit %d stderr=%q", code, stderr)
		}
		code, _, stderr = run(t, work, "state", "clarified", "--dir", inst)
		if code != 0 {
			t.Fatalf("state clarified exit %d stderr=%q", code, stderr)
		}
		snap := snapshotTree(t, inst)

		code, _, stderr = run(t, work, "state", "handed-off", "--dir", inst)
		if code != 1 {
			t.Fatalf("state handed-off without packet exit %d want 1 stderr=%q", code, stderr)
		}
		if !strings.Contains(stderr, "state=handed-off requires a handoff packet") {
			t.Fatalf("stderr=%q", stderr)
		}
		if !strings.Contains(stderr, "mycelium handoff") {
			t.Fatalf("stderr=%q", stderr)
		}
		assertTreeUnchanged(t, inst, snap)
	})

	t.Run("MS-601-CHECK", func(t *testing.T) {
		work := t.TempDir()

		// stored handed-off without packet → check FAIL
		bare := filepath.Join(work, "stored-bare")
		code, _, stderr := run(t, work, "new", "idea", "Stored Bare", "--offline", "--dir", bare)
		if code != 0 {
			t.Fatalf("scaffold bare exit %d stderr=%q", code, stderr)
		}
		patchManifestState(t, bare, "handed-off")
		code, _, stderr = run(t, work, "check", "--dir", bare)
		if code != 1 {
			t.Fatalf("bare handed-off check exit %d want 1 stderr=%q", code, stderr)
		}
		if !strings.Contains(stderr, "handoff packet") {
			t.Fatalf("stderr=%q", stderr)
		}

		// stored handed-off + valid handoff/ → check PASS
		ok := filepath.Join(work, "stored-ok")
		copyTree(t, filepath.Join(td, "handoff-fixture"), ok)
		code, _, stderr = run(t, work, "handoff", "--dir", ok)
		if code != 0 {
			t.Fatalf("handoff setup exit %d stderr=%q", code, stderr)
		}
		code, _, stderr = run(t, work, "check", "--dir", ok)
		if code != 0 {
			t.Fatalf("stored-ok check exit %d stderr=%q", code, stderr)
		}

		// clarified, no handoff/ → check PASS
		clar := filepath.Join(work, "clarified")
		code, _, stderr = run(t, work, "new", "idea", "Clarified Check", "--offline", "--dir", clar)
		if code != 0 {
			t.Fatalf("scaffold clar exit %d stderr=%q", code, stderr)
		}
		code, _, stderr = run(t, work, "state", "exploring", "--dir", clar)
		if code != 0 {
			t.Fatalf("state exploring exit %d stderr=%q", code, stderr)
		}
		code, _, stderr = run(t, work, "state", "clarified", "--dir", clar)
		if code != 0 {
			t.Fatalf("state clarified exit %d stderr=%q", code, stderr)
		}
		if _, err := os.Stat(filepath.Join(clar, "handoff")); !os.IsNotExist(err) {
			t.Fatalf("clarified fixture must not have handoff/: %v", err)
		}
		code, _, stderr = run(t, work, "check", "--dir", clar)
		if code != 0 {
			t.Fatalf("clarified check exit %d stderr=%q", code, stderr)
		}
	})

	t.Run("MS-601-GOLD", func(t *testing.T) {
		work := t.TempDir()
		inst := filepath.Join(work, "gold-fixture")
		copyTree(t, filepath.Join(td, "handoff-fixture"), inst)

		code, _, stderr := run(t, work, "handoff", "--dir", inst)
		if code != 0 {
			t.Fatalf("handoff exit %d stderr=%q", code, stderr)
		}
		accSrc := filepath.Join(inst, "handoff", "acceptance", "add_test.go")
		if _, err := os.Stat(accSrc); err != nil {
			t.Fatal(err)
		}

		goldDir := filepath.Join(work, "gold-run")
		if err := os.MkdirAll(goldDir, 0o755); err != nil {
			t.Fatal(err)
		}
		mustCopyFile(t, accSrc, filepath.Join(goldDir, "add_test.go"))
		mustCopyFile(t, filepath.Join(td, "golden-add", "add.go"), filepath.Join(goldDir, "add.go"))
		writeFile(t, filepath.Join(goldDir, "go.mod"), "module acceptance\n\ngo 1.26\n")

		cmd := exec.Command("go", "test", ".")
		cmd.Dir = goldDir
		cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("golden go test failed: %v\n%s", err, out)
		}

		writeFile(t, filepath.Join(goldDir, "add.go"), "package acceptance\n\nfunc Add(a, b int) int { return 0 }\n")
		cmd = exec.Command("go", "test", ".")
		cmd.Dir = goldDir
		cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
		out, err = cmd.CombinedOutput()
		if err == nil {
			t.Fatalf("broken impl must fail go test; output:\n%s", out)
		}
		if ee, ok := err.(*exec.ExitError); !ok || ee.ExitCode() == 0 {
			t.Fatalf("broken impl exit=%v output:\n%s", err, out)
		}
	})

	t.Run("MS-601-COV", func(t *testing.T) {
		root := repoRoot(t)
		// Stamp: 85 (Appendix E).
		// Product packages = new packages (handoff, handoffcmd). Exclude generated
		// embed, vendor, and clitest harness / data-only fixtures from this run.
		handoffPct := packageCoverPercent(t, root, "./internal/handoff/")
		if handoffPct < float64(ms601CoverageFloor) {
			t.Fatalf("internal/handoff coverage %.1f%% < %d", handoffPct, ms601CoverageFloor)
		}
		handoffcmdPct := packageCoverPercent(t, root, "./internal/handoffcmd/")
		if handoffcmdPct < float64(ms601CoverageFloor) {
			t.Fatalf("internal/handoffcmd coverage %.1f%% < %d", handoffcmdPct, ms601CoverageFloor)
		}
		productPct := packagesCoverPercent(t, root, []string{"./internal/handoff/", "./internal/handoffcmd/"})
		if productPct < float64(ms601CoverageFloor) {
			t.Fatalf("product packages coverage %.1f%% < %d", productPct, ms601CoverageFloor)
		}
		t.Logf("coverage stamp=%d handoff=%.1f%% handoffcmd=%.1f%% product=%.1f%%",
			ms601CoverageFloor, handoffPct, handoffcmdPct, productPct)
	})

	t.Run("MS-601-guards", func(t *testing.T) {
		root := repoRoot(t)
		wf := filepath.Join(root, ".github", "workflows")
		entries, err := os.ReadDir(wf)
		if err != nil {
			t.Fatal(err)
		}
		for _, e := range entries {
			name := e.Name()
			if strings.HasPrefix(name, "phase-06-") && strings.HasSuffix(name, ".yml") {
				t.Fatalf("must not add PHASE-06 workflow: %s", name)
			}
		}
	})
}

func packageCoverPercent(t *testing.T, root, pattern string) float64 {
	t.Helper()
	return packagesCoverPercent(t, root, []string{pattern})
}

func packagesCoverPercent(t *testing.T, root string, patterns []string) float64 {
	t.Helper()
	dir := t.TempDir()
	profile := filepath.Join(dir, "cover.out")
	args := append([]string{"test", "-count=1", "-coverprofile=" + profile, "-covermode=set"}, patterns...)
	cmd := exec.Command("go", args...)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go test %v failed: %v\n%s", patterns, err, out)
	}
	return coverProfileTotal(t, profile)
}

func coverProfileTotal(t *testing.T, profile string) float64 {
	t.Helper()
	cmd := exec.Command("go", "tool", "cover", "-func="+profile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go tool cover: %v\n%s", err, out)
	}
	sc := bufio.NewScanner(bytes.NewReader(out))
	var last string
	for sc.Scan() {
		last = sc.Text()
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	// total: (statements) 86.4%
	fields := strings.Fields(last)
	if len(fields) < 3 || !strings.HasPrefix(fields[0], "total:") {
		t.Fatalf("unexpected cover -func footer: %q", last)
	}
	pctStr := strings.TrimSuffix(fields[len(fields)-1], "%")
	pct, err := strconv.ParseFloat(pctStr, 64)
	if err != nil {
		t.Fatalf("parse cover percent %q: %v", pctStr, err)
	}
	return pct
}
