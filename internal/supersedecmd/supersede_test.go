package supersedecmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/robertguss/mycelium/internal/cli"
	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/execrun"
	"github.com/robertguss/mycelium/internal/journal"
	"github.com/robertguss/mycelium/internal/lock"
	"github.com/robertguss/mycelium/internal/manifest"
	"github.com/robertguss/mycelium/internal/metadata"
	"github.com/robertguss/mycelium/internal/op"
	"github.com/robertguss/mycelium/internal/supersede"
	"github.com/robertguss/mycelium/internal/supersedecmd"
)

func fixedDeps(t *testing.T, cwd string) cli.Deps {
	t.Helper()
	return cli.Deps{
		Clock:     clock.Fixed{T: time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC)},
		Runner:    &execrun.Recording{Inner: execrun.Real{}},
		Getwd:     func() (string, error) { return cwd, nil },
		LookupEnv: func(string) string { return "" },
	}
}

func scaffoldOffline(t *testing.T, cwd, name string) string {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := cli.Run(
		[]string{"mycelium", "new", "idea", name, "--offline"},
		&stdout, &stderr, fixedDeps(t, cwd),
	)
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr.String())
	}
	slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	return filepath.Join(cwd, slug)
}

func runCLI(t *testing.T, deps cli.Deps, args ...string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := cli.Run(append([]string{"mycelium"}, args...), &stdout, &stderr, deps)
	return code, stdout.String(), stderr.String()
}

func mustDecisions(t *testing.T, deps cli.Deps, inst string, titles ...string) {
	t.Helper()
	for _, title := range titles {
		code, _, stderr := runCLI(t, deps, "new", "decision", title, "--dir", inst)
		if code != 0 {
			t.Fatalf("new decision %q exit %d stderr=%q", title, code, stderr)
		}
	}
}

func runSupersede(t *testing.T, cwd, inst, oldID, newID string, argv []string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	if argv == nil {
		argv = []string{"supersede", oldID, "--by", newID, "--dir", inst}
	}
	code := supersedecmd.Run(supersedecmd.Options{
		OldID: oldID,
		NewID: newID,
		Dir:   inst,
		Cwd:   cwd,
		Argv:  argv,
	}, supersedecmd.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC)},
		Stdout: &stdout,
		Stderr: &stderr,
	})
	return code, stdout.String(), stderr.String()
}

func TestRunHappyPath(t *testing.T) {
	cwd := t.TempDir()
	deps := fixedDeps(t, cwd)
	inst := scaffoldOffline(t, cwd, "Happy Supersede")
	mustDecisions(t, deps, inst, "Use SQLite", "Use SQLite with WAL")
	stateBefore := readManifest(t, inst).State

	code, stdout, stderr := runSupersede(t, cwd, inst, "DEC-001", "DEC-002", nil)
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stdout, "mycelium supersede: ok") ||
		!strings.Contains(stdout, "old: DEC-001") ||
		!strings.Contains(stdout, "new: DEC-002") {
		t.Fatalf("stdout=%q", stdout)
	}

	oldMeta := readMeta(t, inst, "decisions", "DEC-001")
	if oldMeta["status"] != "Superseded" || oldMeta["superseded_by"] != "DEC-002" {
		t.Fatalf("OLD=%v", oldMeta)
	}
	newMeta := readMeta(t, inst, "decisions", "DEC-002")
	if newMeta["supersedes"] != "DEC-001" {
		t.Fatalf("NEW=%v", newMeta)
	}
	logBody := string(mustRead(t, filepath.Join(inst, "log.md")))
	if !strings.Contains(logBody, "\tsupersede\tDEC-001\tDEC-001 -> DEC-002") {
		t.Fatalf("log=%s", logBody)
	}
	m := readManifest(t, inst)
	if m.State != stateBefore {
		t.Fatalf("state changed")
	}
	if m.UpdatedDate != "2026-08-15" {
		t.Fatalf("updated_date=%q", m.UpdatedDate)
	}

	code, _, stderr = runCLI(t, deps, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("check exit %d stderr=%q", code, stderr)
	}
}

func TestRunRefuseTable(t *testing.T) {
	cwd := t.TempDir()
	deps := fixedDeps(t, cwd)
	inst := scaffoldOffline(t, cwd, "Refuse Supersede")
	mustDecisions(t, deps, inst, "A", "B", "C")
	code, _, stderr := runCLI(t, deps, "new", "question", "Q1", "--dir", inst)
	if code != 0 {
		t.Fatalf("oq exit %d %s", code, stderr)
	}
	code, _, stderr = runCLI(t, deps, "new", "question", "Q2", "--dir", inst)
	if code != 0 {
		t.Fatalf("oq2 exit %d %s", code, stderr)
	}
	code, _, stderr = runCLI(t, deps, "new", "assumption", "Asm", "--dir", inst)
	if code != 0 {
		t.Fatalf("asm exit %d %s", code, stderr)
	}

	cases := []struct {
		old, neu, want string
	}{
		{"", "DEC-002", "supersede requires <OLD-ID> --by <NEW-ID>"},
		{"DEC-001", "DEC-001", "cannot supersede an ID with itself"},
		{"spark", "DEC-001", "spark is an idea state, not an artifact"},
		{"OQ-001", "OQ-002", "type OQ is not supersedable"},
		{"DEC-999", "DEC-001", "no artifact DEC-999"},
		{"DEC-001", "ASM-001", "same namespace"},
	}
	for _, tc := range cases {
		code, _, stderr := runSupersede(t, cwd, inst, tc.old, tc.neu, nil)
		if code != 1 {
			t.Fatalf("%s/%s exit %d", tc.old, tc.neu, code)
		}
		if !strings.Contains(stderr, tc.want) {
			t.Fatalf("%s/%s stderr=%q want %q", tc.old, tc.neu, stderr, tc.want)
		}
	}

	code, _, stderr = runSupersede(t, cwd, inst, "DEC-001", "DEC-002", nil)
	if code != 0 {
		t.Fatalf("setup supersede: %s", stderr)
	}
	code, _, stderr = runSupersede(t, cwd, inst, "DEC-001", "DEC-003", nil)
	if code != 1 || !strings.Contains(stderr, "already Superseded") {
		t.Fatalf("already: exit %d stderr=%q", code, stderr)
	}
	code, _, stderr = runSupersede(t, cwd, inst, "DEC-003", "DEC-002", nil)
	if code != 1 || !strings.Contains(stderr, "already supersedes") {
		t.Fatalf("one-to-one: exit %d stderr=%q", code, stderr)
	}
}

func TestRunNotInstance(t *testing.T) {
	dir := t.TempDir()
	code, _, stderr := runSupersede(t, dir, dir, "DEC-001", "DEC-002", nil)
	if code != 1 || !strings.Contains(stderr, "not a mycelium instance") {
		t.Fatalf("exit %d stderr=%q", code, stderr)
	}
}

func TestRunResumeAfterPartial(t *testing.T) {
	cwd := t.TempDir()
	deps := fixedDeps(t, cwd)
	inst := scaffoldOffline(t, cwd, "Resume Cmd")
	mustDecisions(t, deps, inst, "Old", "New")

	argv := []string{"supersede", "DEC-001", "--by", "DEC-002", "--dir", inst}
	title := "DEC-001 -> DEC-002"
	now := time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC)
	sess, err := op.Begin(inst, op.Intent{
		Op: "supersede", Type: "decision", Title: title, OriginalID: "DEC-001",
		LogLine: "2026-08-15\tsupersede\tDEC-001\t" + title, Argv: argv, OpID: "resume-cmd",
	}, now)
	if err != nil {
		t.Fatal(err)
	}
	oldRel := findRel(t, inst, "decisions", "DEC-001")
	newRel := findRel(t, inst, "decisions", "DEC-002")
	oldOut, err := supersede.ApplyOld(mustRead(t, filepath.Join(inst, filepath.FromSlash(oldRel))), "DEC-002")
	if err != nil {
		t.Fatal(err)
	}
	newOut, err := supersede.ApplyNew(mustRead(t, filepath.Join(inst, filepath.FromSlash(newRel))), "DEC-001")
	if err != nil {
		t.Fatal(err)
	}
	if err := sess.Stage([]op.Staged{
		{RelTo: oldRel, Content: oldOut},
		{RelTo: newRel, Content: newOut},
		{RelTo: "index.md", Content: mustRead(t, filepath.Join(inst, "index.md"))},
		{RelTo: "log.md", Content: append(mustRead(t, filepath.Join(inst, "log.md")), []byte("2026-08-15\tsupersede\tDEC-001\t"+title+"\n")...)},
		{RelTo: "mycelium.toml", Content: mustRead(t, filepath.Join(inst, "mycelium.toml"))},
	}); err != nil {
		t.Fatal(err)
	}
	if err := sess.CommitPartial(1); err != nil {
		t.Fatal(err)
	}
	_ = sess.Close()
	if _, err := journal.Load(inst); err != nil {
		t.Fatal(err)
	}
	info, _ := lock.Inspect(inst)
	if info.State == lock.Live {
		t.Fatal("live lock")
	}

	code, stdout, stderr := runSupersede(t, cwd, inst, "DEC-001", "DEC-002", argv)
	if code != 0 {
		t.Fatalf("resume exit %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stdout, "mycelium supersede: ok") {
		t.Fatalf("stdout=%q", stdout)
	}
	if _, err := journal.Load(inst); err == nil {
		t.Fatal("journal should be gone")
	}
}

func TestRunJournalMismatch(t *testing.T) {
	cwd := t.TempDir()
	deps := fixedDeps(t, cwd)
	inst := scaffoldOffline(t, cwd, "Mismatch")
	mustDecisions(t, deps, inst, "Old", "New")

	now := time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC)
	sess, err := op.Begin(inst, op.Intent{
		Op: "tier", Title: "standard", LogLine: "2026-08-15\ttier\t-\tfocused -> standard",
		Argv: []string{"tier", "standard"}, OpID: "other-op",
	}, now)
	if err != nil {
		t.Fatal(err)
	}
	if err := sess.Stage([]op.Staged{{RelTo: "log.md", Content: mustRead(t, filepath.Join(inst, "log.md"))}}); err != nil {
		t.Fatal(err)
	}
	_ = sess.Close()

	code, _, stderr := runSupersede(t, cwd, inst, "DEC-001", "DEC-002", nil)
	if code != 1 || !strings.Contains(stderr, "leftover journal for a different operation") {
		t.Fatalf("exit %d stderr=%q", code, stderr)
	}
}

func TestCLIHelpAndFlags(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"mycelium", "supersede", "-h"}, &stdout, &stderr, cli.Deps{})
	if code != 0 || !strings.Contains(stdout.String(), "--by") {
		t.Fatalf("help exit %d out=%q err=%q", code, stdout.String(), stderr.String())
	}
	stdout.Reset()
	code = cli.Run([]string{"mycelium", "--help"}, &stdout, &stderr, cli.Deps{})
	if code != 0 || !strings.Contains(stdout.String(), "mycelium supersede <OLD-ID> --by <NEW-ID>") {
		t.Fatalf("usage missing supersede")
	}
	stderr.Reset()
	code = cli.Run([]string{"mycelium", "supersede", "DEC-001"}, &stdout, &stderr, cli.Deps{})
	if code != 1 || !strings.Contains(stderr.String(), "supersede requires <OLD-ID> --by <NEW-ID>") {
		t.Fatalf("missing --by: %q", stderr.String())
	}
	stderr.Reset()
	code = cli.Run([]string{"mycelium", "supersede", "DEC-001", "--by", "DEC-002", "--nope"}, &stdout, &stderr, cli.Deps{})
	if code != 1 || !strings.Contains(stderr.String(), "unknown flag") {
		t.Fatalf("unknown flag: %q", stderr.String())
	}
}

func readManifest(t *testing.T, inst string) manifest.Manifest {
	t.Helper()
	m, err := manifest.Parse(mustRead(t, filepath.Join(inst, "mycelium.toml")))
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func readMeta(t *testing.T, inst, home, id string) map[string]any {
	t.Helper()
	rel := findRel(t, inst, home, id)
	doc, err := metadata.Parse(mustRead(t, filepath.Join(inst, filepath.FromSlash(rel))))
	if err != nil {
		t.Fatal(err)
	}
	return doc.Meta
}

func findRel(t *testing.T, inst, home, id string) string {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(inst, home))
	if err != nil {
		t.Fatal(err)
	}
	prefix := id + "-"
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), prefix) && strings.HasSuffix(e.Name(), ".md") {
			return filepath.ToSlash(filepath.Join(home, e.Name()))
		}
	}
	t.Fatalf("no %s in %s", id, home)
	return ""
}

func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
