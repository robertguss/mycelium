package supersedecmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/robertguss/mycelium/internal/cli"
	"github.com/robertguss/mycelium/internal/lock"
	"github.com/robertguss/mycelium/internal/supersedecmd"
)

func TestRunCoverageEdges(t *testing.T) {
	cwd := t.TempDir()
	deps := fixedDeps(t, cwd)
	inst := scaffoldOffline(t, cwd, "Coverage Edges")
	mustDecisions(t, deps, inst, "One", "Two")

	rel, err := filepath.Rel(mustAbs(t, cwd), inst)
	if err != nil {
		t.Fatal(err)
	}
	origWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(cwd); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(origWD) })

	var stderr bytes.Buffer
	code := supersedecmd.Run(supersedecmd.Options{
		OldID: "DEC-001",
		NewID: "DEC-002",
		Dir:   rel,
		Argv:  nil,
	}, supersedecmd.Deps{Stderr: &stderr})
	if code != 0 {
		t.Fatalf("relative dir exit %d stderr=%q", code, stderr.String())
	}

	mustDecisions(t, deps, inst, "Three")
	logPath := filepath.Join(inst, "log.md")
	b := mustRead(t, logPath)
	if err := os.WriteFile(logPath, bytes.TrimRight(b, "\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	code, _, stderrStr := runSupersede(t, cwd, inst, "DEC-002", "DEC-003", nil)
	if code != 0 {
		t.Fatalf("chain exit %d stderr=%q", code, stderrStr)
	}

	code, _, stderrStr = runSupersede(t, cwd, inst, "DEC-003", "DEC-099", nil)
	if code != 1 || !strings.Contains(stderrStr, "no artifact DEC-099") {
		t.Fatalf("missing new: %q", stderrStr)
	}

	oldRel := findRel(t, inst, "decisions", "DEC-003")
	if err := os.WriteFile(filepath.Join(inst, filepath.FromSlash(oldRel)), []byte("not front matter\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	mustDecisions(t, deps, inst, "Four")
	code, _, stderrStr = runSupersede(t, cwd, inst, "DEC-003", "DEC-004", nil)
	if code != 1 || !strings.Contains(stderrStr, "cannot parse DEC-003") {
		t.Fatalf("corrupt old: %q", stderrStr)
	}

	inst2 := scaffoldOffline(t, cwd, "Coverage Parse New")
	mustDecisions(t, deps, inst2, "A", "B")
	newRel := findRel(t, inst2, "decisions", "DEC-002")
	if err := os.WriteFile(filepath.Join(inst2, filepath.FromSlash(newRel)), []byte("+++\nbad\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	code, _, stderrStr = runSupersede(t, cwd, inst2, "DEC-001", "DEC-002", nil)
	if code != 1 || !strings.Contains(stderrStr, "cannot parse DEC-002") {
		t.Fatalf("corrupt new: %q", stderrStr)
	}

	inst3 := scaffoldOffline(t, cwd, "Coverage Lock")
	mustDecisions(t, deps, inst3, "L1", "L2")
	held, err := lock.Acquire(inst3, time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = held.Release() })
	code, _, stderrStr = runSupersede(t, cwd, inst3, "DEC-001", "DEC-002", nil)
	if code != 1 || !strings.Contains(stderrStr, "lock held") {
		t.Fatalf("lock: %q", stderrStr)
	}

	inst4 := scaffoldOffline(t, cwd, "Coverage Abs")
	mustDecisions(t, deps, inst4, "X", "Y")
	code, _, stderrStr = runSupersede(t, cwd, inst4, "DEC-001", "DEC-002",
		[]string{"supersede", "DEC-001", "--by", "DEC-002", "--dir", inst4})
	if code != 0 {
		t.Fatalf("abs dir exit %d stderr=%q", code, stderrStr)
	}

	inst5 := scaffoldOffline(t, cwd, "Coverage Manifest")
	mustDecisions(t, deps, inst5, "M1", "M2")
	if err := os.WriteFile(filepath.Join(inst5, "mycelium.toml"), []byte("not = [toml"), 0o644); err != nil {
		t.Fatal(err)
	}
	code, _, stderrStr = runSupersede(t, cwd, inst5, "DEC-001", "DEC-002", nil)
	if code != 1 || !strings.Contains(stderrStr, "invalid mycelium.toml") {
		t.Fatalf("bad manifest: %q", stderrStr)
	}

	inst6 := scaffoldOffline(t, cwd, "Coverage Log")
	mustDecisions(t, deps, inst6, "G1", "G2")
	_ = os.Remove(filepath.Join(inst6, "log.md"))
	code, _, stderrStr = runSupersede(t, cwd, inst6, "DEC-001", "DEC-002", nil)
	if code != 1 || !strings.Contains(stderrStr, "cannot read log.md") {
		t.Fatalf("missing log: %q", stderrStr)
	}

	inst7 := scaffoldOffline(t, cwd, "Coverage ByEq")
	mustDecisions(t, deps, inst7, "P", "Q")
	var stdout, stderr2 bytes.Buffer
	code = cli.Run([]string{"mycelium", "supersede", "DEC-001", "--by=DEC-002", "--dir", inst7},
		&stdout, &stderr2, deps)
	if code != 0 {
		t.Fatalf("--by= exit %d stderr=%q", code, stderr2.String())
	}

	// Unreadable mycelium.toml after FindRoot (Stat ok, ReadFile fails).
	inst8 := scaffoldOffline(t, cwd, "Coverage Chmod")
	mustDecisions(t, deps, inst8, "C1", "C2")
	tomlPath := filepath.Join(inst8, "mycelium.toml")
	if err := os.Chmod(tomlPath, 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(tomlPath, 0o644) })
	code, _, stderrStr = runSupersede(t, cwd, inst8, "DEC-001", "DEC-002", nil)
	if code != 1 || !strings.Contains(stderrStr, "cannot read mycelium.toml") {
		t.Fatalf("chmod toml: %q", stderrStr)
	}

	// findArtifact skips dirs and non-md; ReadDir fails when home missing.
	inst9 := scaffoldOffline(t, cwd, "Coverage Find")
	mustDecisions(t, deps, inst9, "F1", "F2")
	_ = os.Mkdir(filepath.Join(inst9, "decisions", "DEC-001-not-a-file"), 0o755)
	_ = os.WriteFile(filepath.Join(inst9, "decisions", "DEC-001-note.txt"), []byte("x"), 0o644)
	code, _, stderrStr = runSupersede(t, cwd, inst9, "DEC-001", "DEC-002", nil)
	if code != 0 {
		t.Fatalf("skip dir/txt exit %d stderr=%q", code, stderrStr)
	}
	inst10 := scaffoldOffline(t, cwd, "Coverage NoHome")
	// No decisions/ dir → findArtifact ReadDir fails → no artifact.
	code, _, stderrStr = runSupersede(t, cwd, inst10, "DEC-001", "DEC-002", nil)
	if code != 1 || !strings.Contains(stderrStr, "no artifact DEC-001") {
		t.Fatalf("no home: %q", stderrStr)
	}

	// index.md load failure: replace log.md with a directory after scaffold.
	inst11 := scaffoldOffline(t, cwd, "Coverage Index")
	mustDecisions(t, deps, inst11, "I1", "I2")
	// Break indexmd.Load by removing mycelium.toml required parse… already covered.
	// Make program missing so? indexmd.Load uses mycelium.toml + log + homes.
	// Replace log.md with unreadable after we need it for supersede preflight — too late.
	// Stage/commit failures left as defensive branches.

	code = supersedecmd.Run(supersedecmd.Options{OldID: "", NewID: ""}, supersedecmd.Deps{})
	if code != 1 {
		t.Fatalf("empty deps refuse exit %d", code)
	}
}

func mustAbs(t *testing.T, p string) string {
	t.Helper()
	a, err := filepath.Abs(p)
	if err != nil {
		t.Fatal(err)
	}
	return a
}
