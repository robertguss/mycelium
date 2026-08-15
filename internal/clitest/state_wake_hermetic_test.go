package clitest_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/robertguss/mycelium/internal/cli"
	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/clitest"
	"github.com/robertguss/mycelium/internal/execrun"
	"github.com/robertguss/mycelium/internal/manifest"
	"github.com/robertguss/mycelium/internal/metadata"
)

func TestStateWakeHermetic(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}

	inst := scaffoldIdea(t, work, "Wake Slice", clk, rec, env)
	clitest.AssertNoNetwork(t, rec)

	code, stdout, stderr := runCLI(t, clk, rec, env, work, "state", "exploring", "--dir", inst)
	if code != 0 {
		t.Fatalf("state exploring exit %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stdout, "state: exploring") {
		t.Fatalf("stdout=%q", stdout)
	}

	code, _, stderr = runCLI(t, clk, rec, env, work, "state", "simmering", "--dir", inst)
	if code != 1 {
		t.Fatalf("simmer without revisit: exit %d want 1 stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "program/contracts/revisit.md") {
		t.Fatalf("want revisit contract, got %q", stderr)
	}

	code, _, stderr = runCLI(t, clk, rec, env, work, "state", "clarified", "--revisit", "2026-08-08", "--dir", inst)
	if code != 1 {
		t.Fatalf("--revisit on clarified: exit %d want 1 stderr=%q", code, stderr)
	}

	code, _, stderr = runCLI(t, clk, rec, env, work, "state", "handed-off", "--dir", inst)
	if code != 1 {
		t.Fatalf("handed-off: exit %d want 1 stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "handoff packet") {
		t.Fatalf("want handoff packet teaching error, got %q", stderr)
	}
	if !strings.Contains(stderr, "mycelium handoff") {
		t.Fatalf("want mycelium handoff in fix, got %q", stderr)
	}

	code, stdout, stderr = runCLI(t, clk, rec, env, work, "state", "clarified", "--dir", inst)
	if code != 0 {
		t.Fatalf("state clarified exit %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stdout, "state: clarified") {
		t.Fatalf("stdout=%q", stdout)
	}
	code, _, stderr = runCLI(t, clk, rec, env, work, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("check after clarified exit %d stderr=%q", code, stderr)
	}

	// Reset to exploring for wake path via archive? clarified → exploring illegal.
	// Use a fresh instance for simmer/wake coverage.
	inst2 := scaffoldIdea(t, work, "Wake Ritual", clk, rec, env)
	mustState(t, clk, rec, env, work, inst2, "exploring")

	code, _, stderr = runCLI(t, clk, rec, env, work, "wake", "--dir", inst2)
	if code != 1 {
		t.Fatalf("wake from exploring: exit %d want 1 stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "wake is only legal from simmering") {
		t.Fatalf("stderr=%q", stderr)
	}
	if !strings.Contains(stderr, "program/contracts/lifecycle.md") {
		t.Fatalf("stderr=%q", stderr)
	}

	code, _, stderr = runCLI(t, clk, rec, env, work,
		"state", "simmering", "--revisit", "2026-08-08", "--dir", inst2)
	if code != 0 {
		t.Fatalf("simmer exit %d stderr=%q", code, stderr)
	}

	code, stdout, stderr = runCLI(t, clk, rec, env, work, "state", "exploring", "--dir", inst2)
	if code != 0 {
		t.Fatalf("already exploring? after simmer then exploring should wake; exit %d stderr=%q stdout=%q", code, stderr, stdout)
	}
	// We just simmered then immediately woke via state exploring — verify wake stdout.
	if !strings.Contains(stdout, "woke briefs/WAKE-2026-08-01.md") {
		t.Fatalf("want wake stdout, got %q", stdout)
	}
	clitest.AssertNoNetwork(t, rec)
}

func TestWakeBriefSharedAndOverwrite(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}

	inst := scaffoldIdea(t, work, "Brief Fixture", clk, rec, env)
	mustState(t, clk, rec, env, work, inst, "exploring")

	code, _, stderr := runCLI(t, clk, rec, env, work, "new", "decision", "Park the idea", "--dir", inst)
	if code != 0 {
		t.Fatalf("new decision exit %d stderr=%q", code, stderr)
	}
	code, _, stderr = runCLI(t, clk, rec, env, work, "new", "assumption", "API stays stable", "--dir", inst)
	if code != 0 {
		t.Fatalf("new assumption exit %d stderr=%q", code, stderr)
	}
	code, _, stderr = runCLI(t, clk, rec, env, work, "new", "evidence", "Vendor changelog", "--dir", inst)
	if code != 0 {
		t.Fatalf("new evidence exit %d stderr=%q", code, stderr)
	}
	code, _, stderr = runCLI(t, clk, rec, env, work, "new", "assumption", "Budget is unlimited", "--dir", inst)
	if code != 0 {
		t.Fatalf("new assumption 2 exit %d stderr=%q", code, stderr)
	}

	patchASM(t, inst, "ASM-001", "Held", "2026-08-05")
	patchEVD(t, inst, "EVD-001", "2026-08-06")
	patchASM(t, inst, "ASM-002", "Retired", "")

	code, _, stderr = runCLI(t, clk, rec, env, work,
		"state", "simmering", "--revisit", "2026-08-08", "--dir", inst)
	if code != 0 {
		t.Fatalf("simmer exit %d stderr=%q", code, stderr)
	}

	// Copy for state exploring wake vs wake command.
	copyA := filepath.Join(work, "wake-via-wake")
	copyB := filepath.Join(work, "wake-via-state")
	copyTree(t, inst, copyA)
	copyTree(t, inst, copyB)

	wakeClk := clock.Fixed{T: time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC)}
	code, stdout, stderr := runCLI(t, wakeClk, rec, env, work, "wake", "--dir", copyA)
	if code != 0 {
		t.Fatalf("wake exit %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stdout, "woke briefs/WAKE-2026-08-09.md") || !strings.Contains(stdout, "state: exploring") {
		t.Fatalf("stdout=%q", stdout)
	}

	code, stdout, stderr = runCLI(t, wakeClk, rec, env, work, "state", "exploring", "--dir", copyB)
	if code != 0 {
		t.Fatalf("state exploring wake exit %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stdout, "woke briefs/WAKE-2026-08-09.md") {
		t.Fatalf("silent wake forbidden; stdout=%q", stdout)
	}

	briefA := readFile(t, filepath.Join(copyA, "briefs", "WAKE-2026-08-09.md"))
	briefB := readFile(t, filepath.Join(copyB, "briefs", "WAKE-2026-08-09.md"))
	latestA := readFile(t, filepath.Join(copyA, "briefs", "LATEST.md"))
	if !bytes.Equal(briefA, latestA) {
		t.Fatal("LATEST must equal dated brief")
	}
	// Same citation set + H2s (bytes may match; require H2s + IDs).
	for _, body := range [][]byte{briefA, briefB} {
		s := string(body)
		for _, h2 := range []string{"## Parked", "## Log since simmer", "## Evidence triggers", "## Assumptions", "## Suggested next"} {
			if !strings.Contains(s, h2) {
				t.Fatalf("missing %s in brief:\n%s", h2, s)
			}
		}
		if !strings.Contains(s, "EVD-001") || !strings.Contains(s, "ASM-001") {
			t.Fatalf("missing citations:\n%s", s)
		}
		if strings.Contains(s, "ASM-002") {
			t.Fatalf("must not cite ASM-002:\n%s", s)
		}
	}

	m := loadManifest(t, copyA)
	if m.State != "exploring" || m.Revisit != "" {
		t.Fatalf("manifest state=%q revisit=%q", m.State, m.Revisit)
	}

	code, _, stderr = runCLI(t, wakeClk, rec, env, work, "check", "--dir", copyA)
	if code != 0 {
		t.Fatalf("check after wake exit %d stderr=%q", code, stderr)
	}

	// Same-day overwrite: re-simmer and wake again.
	code, _, stderr = runCLI(t, wakeClk, rec, env, work,
		"state", "simmering", "--revisit", "2026-08-10", "--dir", copyA)
	if code != 0 {
		t.Fatalf("re-simmer exit %d stderr=%q", code, stderr)
	}
	code, _, stderr = runCLI(t, wakeClk, rec, env, work, "wake", "--dir", copyA)
	if code != 0 {
		t.Fatalf("second wake exit %d stderr=%q", code, stderr)
	}
	brief2 := readFile(t, filepath.Join(copyA, "briefs", "WAKE-2026-08-09.md"))
	latest2 := readFile(t, filepath.Join(copyA, "briefs", "LATEST.md"))
	if !bytes.Equal(brief2, latest2) {
		t.Fatal("overwrite: LATEST must equal dated brief")
	}
	if !strings.Contains(string(brief2), "simmering revisit=2026-08-10") {
		t.Fatalf("second wake should cite new simmer:\n%s", brief2)
	}

	// Stored handed-off fails check (binary path).
	handed := scaffoldIdea(t, work, "Handed Fail", clk, rec, env)
	patchManifestState(t, handed, "handed-off")
	code, _, stderr = runCLI(t, clk, rec, env, work, "check", "--dir", handed)
	if code != 1 {
		t.Fatalf("handed-off check exit %d want 1 stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "handoff packet") {
		t.Fatalf("stderr=%q", stderr)
	}

	clitest.AssertNoNetwork(t, rec)
	for _, c := range rec.Calls {
		if c.Name == "gh" {
			t.Fatalf("gh invoked: %#v", c)
		}
	}
}

func TestWakeFromSparkRefuse(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}
	inst := scaffoldIdea(t, work, "Spark Wake", clk, rec, env)
	code, _, stderr := runCLI(t, clk, rec, env, work, "wake", "--dir", inst)
	if code != 1 {
		t.Fatalf("exit %d want 1", code)
	}
	if !strings.Contains(stderr, "wake is only legal from simmering") {
		t.Fatalf("stderr=%q", stderr)
	}
}

func TestSameStateNoopAndSimmerRevisitUpdate(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}
	inst := scaffoldIdea(t, work, "Same State", clk, rec, env)
	mustState(t, clk, rec, env, work, inst, "exploring")
	code, stdout, stderr := runCLI(t, clk, rec, env, work, "state", "exploring", "--dir", inst)
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr)
	}
	if strings.TrimSpace(stdout) != "already exploring" {
		t.Fatalf("stdout=%q", stdout)
	}

	mustState(t, clk, rec, env, work, inst, "simmering", "--revisit", "2026-08-08")
	code, stdout, stderr = runCLI(t, clk, rec, env, work,
		"state", "simmering", "--revisit", "2026-08-15", "--dir", inst)
	if code != 0 {
		t.Fatalf("revisit update exit %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stdout, "state: simmering") || !strings.Contains(stdout, "revisit: 2026-08-15") {
		t.Fatalf("stdout=%q", stdout)
	}
	m := loadManifest(t, inst)
	if m.Revisit != "2026-08-15" {
		t.Fatalf("revisit=%q", m.Revisit)
	}
}

func scaffoldIdea(t *testing.T, work, name string, clk clock.Clock, rec *execrun.Recording, env map[string]string) string {
	t.Helper()
	code, _, stderr := runCLI(t, clk, rec, env, work, "new", "idea", name, "--offline")
	if code != 0 {
		t.Fatalf("scaffold %q exit %d stderr=%q", name, code, stderr)
	}
	slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	return filepath.Join(work, slug)
}

func mustState(t *testing.T, clk clock.Clock, rec *execrun.Recording, env map[string]string, work, inst, target string, extra ...string) {
	t.Helper()
	args := append([]string{"state", target, "--dir", inst}, extra...)
	code, _, stderr := runCLI(t, clk, rec, env, work, args...)
	if code != 0 {
		t.Fatalf("state %s exit %d stderr=%q", target, code, stderr)
	}
}

func runCLI(t *testing.T, clk clock.Clock, rec *execrun.Recording, env map[string]string, cwd string, args ...string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := cli.Run(append([]string{"mycelium"}, args...), &stdout, &stderr, cli.Deps{
		Clock:  clk,
		Runner: rec,
		Getwd:  func() (string, error) { return cwd, nil },
		LookupEnv: func(k string) string {
			if env != nil {
				return env[k]
			}
			return ""
		},
	})
	return code, stdout.String(), stderr.String()
}

func loadManifest(t *testing.T, inst string) manifest.Manifest {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(inst, "mycelium.toml"))
	if err != nil {
		t.Fatal(err)
	}
	m, err := manifest.Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func patchManifestState(t *testing.T, inst, state string) {
	t.Helper()
	path := filepath.Join(inst, "mycelium.toml")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	m, err := manifest.Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	m.State = state
	m.Revisit = ""
	out, err := manifest.Encode(m)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, out, 0o644); err != nil {
		t.Fatal(err)
	}
}

func patchASM(t *testing.T, inst, id, status, triggerDate string) {
	t.Helper()
	dir := filepath.Join(inst, "assumptions")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), id+"-") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		doc, err := metadata.Parse(raw)
		if err != nil {
			t.Fatal(err)
		}
		doc.Meta["status"] = status
		body := doc.Body
		if triggerDate != "" {
			body = replaceSection(body, "Revisit Triggers", triggerDate+"\n")
		}
		writeMeta(t, path, doc.Meta, body)
		return
	}
	t.Fatalf("no assumption file for %s", id)
}

func patchEVD(t *testing.T, inst, id, triggerDate string) {
	t.Helper()
	dir := filepath.Join(inst, "evidence")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), id+"-") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		doc, err := metadata.Parse(raw)
		if err != nil {
			t.Fatal(err)
		}
		body := replaceSection(doc.Body, "Revalidation Trigger", triggerDate+"\n")
		writeMeta(t, path, doc.Meta, body)
		return
	}
	t.Fatalf("no evidence file for %s", id)
}

func replaceSection(body, heading, newBody string) string {
	want := "## " + heading
	lines := strings.Split(body, "\n")
	start := -1
	for i, line := range lines {
		if strings.TrimRight(line, "\r") == want {
			start = i
			break
		}
	}
	if start < 0 {
		return body + "\n" + want + "\n\n" + newBody
	}
	end := len(lines)
	for i := start + 1; i < len(lines); i++ {
		if strings.HasPrefix(strings.TrimRight(lines[i], "\r"), "## ") {
			end = i
			break
		}
	}
	var out []string
	out = append(out, lines[:start+1]...)
	out = append(out, "", strings.TrimRight(newBody, "\n"))
	out = append(out, lines[end:]...)
	return strings.Join(out, "\n")
}

func writeMeta(t *testing.T, path string, meta map[string]any, body string) {
	t.Helper()
	var b strings.Builder
	b.WriteString("+++\n")
	order := []string{"id", "title", "status", "date", "attached_to", "source"}
	seen := map[string]bool{}
	for _, k := range order {
		v, ok := meta[k]
		if !ok {
			continue
		}
		seen[k] = true
		writeTOMLKey(&b, k, v)
	}
	for k, v := range meta {
		if seen[k] {
			continue
		}
		writeTOMLKey(&b, k, v)
	}
	b.WriteString("+++\n")
	if !strings.HasPrefix(body, "\n") {
		b.WriteByte('\n')
	}
	b.WriteString(body)
	if !strings.HasSuffix(body, "\n") {
		b.WriteByte('\n')
	}
	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeTOMLKey(b *strings.Builder, k string, v any) {
	b.WriteString(k)
	b.WriteString(" = ")
	switch x := v.(type) {
	case string:
		b.WriteByte('"')
		b.WriteString(x)
		b.WriteByte('"')
	case bool:
		if x {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
	case int64:
		fmt.Fprintf(b, "%d", x)
	case float64:
		fmt.Fprintf(b, "%v", x)
	default:
		fmt.Fprintf(b, "%q", fmt.Sprint(x))
	}
	b.WriteByte('\n')
}

func copyTree(t *testing.T, src, dst string) {
	t.Helper()
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return os.WriteFile(target, b, info.Mode())
	})
	if err != nil {
		t.Fatal(err)
	}
}
