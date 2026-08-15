package clitest_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/robertguss/mycelium/internal/clitest"
	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/execrun"
	"github.com/robertguss/mycelium/internal/handoff"
)

// PHASE-06 Slice 4 — canonical fixture + acceptance + golden impl (MS-601-GEN / MS-601-GOLD).
// Does not land the full MS-601 matrix (Slice 5).

func phase06Slice4Testdata(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Join(filepath.Dir(file), "testdata")
}

func TestPhase06Slice4MS601GEN(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1", "MYCELIUM_NOW": "2026-08-15T00:00:00Z"}

	inst := filepath.Join(work, "handoff-fixture")
	copyTree(t, filepath.Join(phase06Slice4Testdata(t), "handoff-fixture"), inst)

	code, _, stderr := runCLI(t, clk, rec, env, work, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("pre-handoff check exit %d stderr=%q", code, stderr)
	}
	if loadManifest(t, inst).State != "clarified" {
		t.Fatalf("fixture state=%q want clarified", loadManifest(t, inst).State)
	}

	code, stdout, stderr := runCLI(t, clk, rec, env, work, "handoff", "--dir", inst)
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
		"DEC-001",
		"playbooks/PLAYBOOK.md",
		"acceptance/add_test.go",
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
	glossary := string(readFile(t, filepath.Join(inst, "handoff", "glossary.md")))
	if !strings.Contains(glossary, "term: Add") {
		t.Fatalf("glossary missing term Add: %q", glossary)
	}
	play := string(readFile(t, filepath.Join(inst, "handoff", "playbooks", "PLAYBOOK.md")))
	if !strings.Contains(play, "add.go") || !strings.Contains(play, "Add(a, b int) int") {
		t.Fatalf("playbook missing target: %q", play)
	}

	logBody := string(readFile(t, filepath.Join(inst, "log.md")))
	wantLine := "2026-08-15\thandoff\tHO-001\tclarified -> handed-off"
	if !strings.Contains(logBody, wantLine) {
		t.Fatalf("log missing %q\n%s", wantLine, logBody)
	}

	code, _, stderr = runCLI(t, clk, rec, env, work, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("post-handoff check exit %d stderr=%q", code, stderr)
	}
	clitest.AssertNoNetwork(t, rec)
}

func TestPhase06Slice4MS601GOLD(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1", "MYCELIUM_NOW": "2026-08-15T00:00:00Z"}
	td := phase06Slice4Testdata(t)

	inst := filepath.Join(work, "gold-fixture")
	copyTree(t, filepath.Join(td, "handoff-fixture"), inst)

	code, _, stderr := runCLI(t, clk, rec, env, work, "handoff", "--dir", inst)
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
	clitest.AssertNoNetwork(t, rec)
}

func mustCopyFile(t *testing.T, src, dst string) {
	t.Helper()
	b, err := os.ReadFile(src)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dst, b, 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}
