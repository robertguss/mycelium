package clitest_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/clitest"
)

// TestMS301HermeticFixtures is the PHASE-03 MS-301 gate (Appendix E / §13).
func TestMS301HermeticFixtures(t *testing.T) {
	bin := clitest.Bin(t)
	stubDir, ghMarker := installGhNeverStub(t)
	work := t.TempDir()
	env := ms301Env(stubDir)

	run := func(args ...string) (int, string, string) {
		t.Helper()
		return clitest.Run(t, bin, work, env, args...)
	}

	t.Run("disputed", func(t *testing.T) {
		inst := filepath.Join(work, "disputed-fixture")
		code, _, stderr := run("new", "idea", "Disputed Fixture", "--offline", "--dir", inst)
		if code != 0 {
			t.Fatalf("new idea: exit %d stderr=%q", code, stderr)
		}
		code, _, stderr = run("new", "question", "Use SQLite", "--dir", inst)
		if code != 0 {
			t.Fatalf("new question: exit %d stderr=%q", code, stderr)
		}
		oq := mustFindOQ(t, inst)
		patchOQAgreement(t, oq, "agree-to-disagree")
		patchOQDisputedRecord(t, oq)

		code, _, stderr = run("check", "--dir", inst)
		if code != 0 {
			t.Fatalf("disputed check: exit %d stderr=%q", code, stderr)
		}

		// Negative: same shape without ## Crux.
		missing := filepath.Join(work, "disputed-missing-crux")
		copyTree(t, inst, missing)
		deleteOQCrux(t, mustFindOQ(t, missing))
		code, _, stderr = run("check", "--dir", missing)
		if code != 1 {
			t.Fatalf("missing-crux check: exit %d want 1 stderr=%q", code, stderr)
		}
		if !strings.Contains(stderr, "## Crux") {
			t.Fatalf("stderr missing ## Crux: %q", stderr)
		}
		if !strings.Contains(stderr, "program/contracts/sparring.md") {
			t.Fatalf("stderr missing sparring.md: %q", stderr)
		}
	})

	t.Run("aligned", func(t *testing.T) {
		inst := filepath.Join(work, "aligned-fixture")
		code, _, stderr := run("new", "idea", "Aligned Fixture", "--offline", "--dir", inst)
		if code != 0 {
			t.Fatalf("new idea: exit %d stderr=%q", code, stderr)
		}
		code, _, stderr = run("new", "question", "Keep the name", "--dir", inst)
		if code != 0 {
			t.Fatalf("new question: exit %d stderr=%q", code, stderr)
		}
		oq := mustFindOQ(t, inst)
		patchOQAgreement(t, oq, "aligned")
		body := string(readFile(t, oq))
		if !strings.Contains(body, "## Positions") {
			t.Fatal("aligned OQ missing ## Positions")
		}
		if strings.Contains(body, "## Crux") || strings.Contains(body, "## Reasons") {
			t.Fatalf("aligned must not add Crux/Reasons:\n%s", body)
		}

		code, _, stderr = run("check", "--dir", inst)
		if code != 0 {
			t.Fatalf("aligned check: exit %d stderr=%q", code, stderr)
		}
	})

	t.Run("invalid-agreement", func(t *testing.T) {
		inst := filepath.Join(work, "invalid-agreement")
		code, _, stderr := run("new", "idea", "Invalid Agreement", "--offline", "--dir", inst)
		if code != 0 {
			t.Fatalf("new idea: exit %d stderr=%q", code, stderr)
		}
		code, _, stderr = run("new", "question", "Maybe later", "--dir", inst)
		if code != 0 {
			t.Fatalf("new question: exit %d stderr=%q", code, stderr)
		}
		patchOQAgreement(t, mustFindOQ(t, inst), "maybe")

		code, _, stderr = run("check", "--dir", inst)
		if code != 1 {
			t.Fatalf("invalid agreement check: exit %d want 1 stderr=%q", code, stderr)
		}
	})

	t.Run("spark-zero-questions", func(t *testing.T) {
		inst := filepath.Join(work, "spark-zero")
		code, _, stderr := run("new", "idea", "Spark Zero", "--offline", "--dir", inst)
		if code != 0 {
			t.Fatalf("new idea: exit %d stderr=%q", code, stderr)
		}
		if entries, err := os.ReadDir(filepath.Join(inst, "questions")); err == nil && len(entries) > 0 {
			t.Fatalf("spark must have zero OQ files, got %d", len(entries))
		}

		code, _, stderr = run("check", "--dir", inst)
		if code != 0 {
			t.Fatalf("spark-zero check: exit %d stderr=%q", code, stderr)
		}
	})

	t.Run("bare-positions-fill", func(t *testing.T) {
		inst := filepath.Join(work, "bare-positions")
		code, _, stderr := run("new", "idea", "Bare Positions", "--offline", "--dir", inst)
		if code != 0 {
			t.Fatalf("new idea: exit %d stderr=%q", code, stderr)
		}
		code, _, stderr = run("new", "question", "Fill later", "--dir", inst)
		if code != 0 {
			t.Fatalf("new question: exit %d stderr=%q", code, stderr)
		}
		oq := mustFindOQ(t, inst)
		body := string(readFile(t, oq))
		if !strings.Contains(body, `agreement = "open"`) {
			t.Fatalf("want agreement open:\n%s", body)
		}
		if !strings.Contains(body, "## Positions") || !strings.Contains(body, "<!-- fill -->") {
			t.Fatalf("want Positions <!-- fill -->:\n%s", body)
		}
		if strings.Contains(body, "## Crux") || strings.Contains(body, "## Reasons") {
			t.Fatalf("fill Positions must not emit Crux/Reasons:\n%s", body)
		}

		code, _, stderr = run("check", "--dir", inst)
		if code != 0 {
			t.Fatalf("bare-positions check: exit %d stderr=%q", code, stderr)
		}
	})

	if _, err := os.Stat(ghMarker); err == nil {
		t.Fatal("gh was invoked (marker present)")
	} else if !os.IsNotExist(err) {
		t.Fatalf("gh marker stat: %v", err)
	}
}

func ms301Env(stubDir string) []string {
	base := stripEnvKeys(os.Environ(), "MYCELIUM_NOW", "MYCELIUM_OFFLINE", "PATH", "MYCELIUM_BIN")
	path := stubDir + string(os.PathListSeparator) + os.Getenv("PATH")
	return append(base,
		"PATH="+path,
		"MYCELIUM_OFFLINE=1",
	)
}

func mustFindOQ(t *testing.T, inst string) string {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(inst, "questions", "OQ-*.md"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 {
		t.Fatalf("want 1 OQ file, got %d in %s", len(matches), inst)
	}
	return matches[0]
}

func patchOQAgreement(t *testing.T, path, agreement string) {
	t.Helper()
	body := string(readFile(t, path))
	updated, n := replaceAgreementLine(body, agreement)
	if n != 1 {
		t.Fatalf("agreement replace count=%d in %s", n, path)
	}
	if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
		t.Fatal(err)
	}
}

func replaceAgreementLine(body, agreement string) (string, int) {
	const prefix = "agreement = "
	n := 0
	lines := strings.Split(body, "\n")
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), prefix) {
			lines[i] = prefix + `"` + agreement + `"`
			n++
		}
	}
	return strings.Join(lines, "\n"), n
}

// patchOQDisputedRecord sets Positions/Reasons/Crux H3 pairs. Bodies are placeholders;
// tests never assert their words (DEC-005).
func patchOQDisputedRecord(t *testing.T, path string) {
	t.Helper()
	body := string(readFile(t, path))
	const positions = `## Positions

### Human

h

### Agent

a
`
	const disagreement = `## Reasons

### Human

rh

### Agent

ra

## Crux

### Human

ch

### Agent

ca

`
	start := strings.Index(body, "## Positions")
	if start < 0 {
		t.Fatal("## Positions missing")
	}
	disp := strings.Index(body[start:], "## Disposition")
	if disp < 0 {
		t.Fatal("## Disposition missing")
	}
	disp += start
	updated := body[:start] + positions + disagreement + body[disp:]
	if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
		t.Fatal(err)
	}
}

func deleteOQCrux(t *testing.T, path string) {
	t.Helper()
	body := string(readFile(t, path))
	start := strings.Index(body, "## Crux")
	if start < 0 {
		t.Fatal("## Crux missing before delete")
	}
	rest := body[start+len("## Crux"):]
	next := strings.Index(rest, "\n## ")
	if next < 0 {
		t.Fatal("no H2 after ## Crux")
	}
	updated := body[:start] + rest[next+1:]
	if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
		t.Fatal(err)
	}
}
