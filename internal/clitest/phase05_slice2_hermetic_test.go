package clitest_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/clitest"
	"github.com/robertguss/mycelium/internal/execrun"
	"github.com/robertguss/mycelium/internal/journal"
	"github.com/robertguss/mycelium/internal/lock"
	"github.com/robertguss/mycelium/internal/metadata"
	"github.com/robertguss/mycelium/internal/op"
	"github.com/robertguss/mycelium/internal/supersede"
)

func TestPhase05Slice2SupersedeHermetic(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}

	inst := scaffoldIdea(t, work, "Supersede Fixture", clk, rec, env)
	stateBefore := loadManifest(t, inst).State

	mustNewDecision(t, clk, rec, env, work, inst, "Use SQLite")
	mustNewDecision(t, clk, rec, env, work, inst, "Use SQLite with WAL")
	patchDECStatus(t, inst, "DEC-001", "Accepted")
	patchDECStatus(t, inst, "DEC-002", "Accepted")

	code, stdout, stderr := runCLI(t, clk, rec, env, work,
		"supersede", "DEC-001", "--by", "DEC-002", "--dir", inst)
	if code != 0 {
		t.Fatalf("supersede exit %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stdout, "mycelium supersede: ok") ||
		!strings.Contains(stdout, "old: DEC-001") ||
		!strings.Contains(stdout, "new: DEC-002") {
		t.Fatalf("stdout=%q", stdout)
	}

	oldDoc := readArtifactMeta(t, inst, "decisions", "DEC-001")
	if oldDoc["status"] != "Superseded" || oldDoc["superseded_by"] != "DEC-002" {
		t.Fatalf("OLD meta=%v", oldDoc)
	}
	newDoc := readArtifactMeta(t, inst, "decisions", "DEC-002")
	if newDoc["supersedes"] != "DEC-001" {
		t.Fatalf("NEW supersedes=%v", newDoc["supersedes"])
	}
	if newDoc["status"] != "Accepted" {
		t.Fatalf("NEW status changed: %v", newDoc["status"])
	}

	logBody := string(readFile(t, filepath.Join(inst, "log.md")))
	wantLine := "2026-08-15\tsupersede\tDEC-001\tDEC-001 -> DEC-002"
	if !strings.Contains(logBody, wantLine) {
		t.Fatalf("log missing %q\n%s", wantLine, logBody)
	}

	m := loadManifest(t, inst)
	if m.State != stateBefore {
		t.Fatalf("state changed %q -> %q", stateBefore, m.State)
	}
	if m.UpdatedDate != "2026-08-15" {
		t.Fatalf("updated_date=%q", m.UpdatedDate)
	}

	code, _, stderr = runCLI(t, clk, rec, env, work, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("check after supersede exit %d stderr=%q", code, stderr)
	}
	clitest.AssertNoNetwork(t, rec)
}

func TestPhase05Slice2SupersedeRefuseTable(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}

	inst := scaffoldIdea(t, work, "Refuse Table", clk, rec, env)
	mustNewDecision(t, clk, rec, env, work, inst, "First")
	mustNewDecision(t, clk, rec, env, work, inst, "Second")
	mustNewDecision(t, clk, rec, env, work, inst, "Third")
	code, _, stderr := runCLI(t, clk, rec, env, work, "new", "question", "Open Q", "--dir", inst)
	if code != 0 {
		t.Fatalf("new question exit %d stderr=%q", code, stderr)
	}
	code, _, stderr = runCLI(t, clk, rec, env, work, "new", "question", "Open Q2", "--dir", inst)
	if code != 0 {
		t.Fatalf("new question 2 exit %d stderr=%q", code, stderr)
	}
	code, _, stderr = runCLI(t, clk, rec, env, work, "new", "assumption", "A1", "--dir", inst)
	if code != 0 {
		t.Fatalf("new assumption exit %d stderr=%q", code, stderr)
	}

	snapshot := snapshotInstance(t, inst)

	cases := []struct {
		name string
		args []string
		what string
	}{
		{"missing-by", []string{"supersede", "DEC-001", "--dir", inst}, "supersede requires <OLD-ID> --by <NEW-ID>"},
		{"self", []string{"supersede", "DEC-001", "--by", "DEC-001", "--dir", inst}, "cannot supersede an ID with itself"},
		{"idea-state", []string{"supersede", "spark", "--by", "DEC-001", "--dir", inst}, "spark is an idea state, not an artifact"},
		{"oq", []string{"supersede", "OQ-001", "--by", "OQ-002", "--dir", inst}, "type OQ is not supersedable"},
		{"missing-old", []string{"supersede", "DEC-999", "--by", "DEC-001", "--dir", inst}, "no artifact DEC-999"},
		{"diff-ns", []string{"supersede", "DEC-001", "--by", "ASM-001", "--dir", inst}, "same namespace"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			code, _, stderr := runCLI(t, clk, rec, env, work, tc.args...)
			if code != 1 {
				t.Fatalf("exit %d want 1 stderr=%q", code, stderr)
			}
			if !strings.Contains(stderr, tc.what) {
				t.Fatalf("stderr=%q want %q", stderr, tc.what)
			}
			assertUnchanged(t, inst, snapshot)
		})
	}

	// Happy path once, then refuse already-superseded / NEW already has supersedes / rewrite past pair.
	code, _, stderr = runCLI(t, clk, rec, env, work,
		"supersede", "DEC-001", "--by", "DEC-002", "--dir", inst)
	if code != 0 {
		t.Fatalf("setup supersede exit %d stderr=%q", code, stderr)
	}
	snap2 := snapshotInstance(t, inst)

	code, _, stderr = runCLI(t, clk, rec, env, work,
		"supersede", "DEC-001", "--by", "DEC-003", "--dir", inst)
	if code != 1 {
		t.Fatalf("already superseded exit %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "DEC-001 is already Superseded by DEC-002") {
		t.Fatalf("stderr=%q", stderr)
	}
	assertUnchanged(t, inst, snap2)

	code, _, stderr = runCLI(t, clk, rec, env, work,
		"supersede", "DEC-003", "--by", "DEC-002", "--dir", inst)
	if code != 1 {
		t.Fatalf("NEW already supersedes exit %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "DEC-002 already supersedes DEC-001") {
		t.Fatalf("stderr=%q", stderr)
	}
	assertUnchanged(t, inst, snap2)

	// Chain tip: DEC-002 → DEC-003 ok.
	code, _, stderr = runCLI(t, clk, rec, env, work,
		"supersede", "DEC-002", "--by", "DEC-003", "--dir", inst)
	if code != 0 {
		t.Fatalf("chain tip exit %d stderr=%q", code, stderr)
	}
	mid := readArtifactMeta(t, inst, "decisions", "DEC-002")
	if mid["status"] != "Superseded" || mid["superseded_by"] != "DEC-003" {
		t.Fatalf("DEC-002 meta=%v", mid)
	}
	tip := readArtifactMeta(t, inst, "decisions", "DEC-003")
	if tip["supersedes"] != "DEC-002" {
		t.Fatalf("DEC-003 supersedes=%v", tip["supersedes"])
	}

	code, _, stderr = runCLI(t, clk, rec, env, work, "state", "handed-off", "--dir", inst)
	if code != 1 {
		t.Fatalf("handed-off exit %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "handoff packet") {
		t.Fatalf("stderr=%q", stderr)
	}
	if !strings.Contains(stderr, "mycelium handoff") {
		t.Fatalf("stderr=%q", stderr)
	}
	clitest.AssertNoNetwork(t, rec)
}

func TestPhase05Slice2SupersedeResumeJournal(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}

	inst := scaffoldIdea(t, work, "Resume Supersede", clk, rec, env)
	mustNewDecision(t, clk, rec, env, work, inst, "Old Rec")
	mustNewDecision(t, clk, rec, env, work, inst, "New Rec")

	argv := []string{"supersede", "DEC-001", "--by", "DEC-002", "--dir", inst}
	title := "DEC-001 -> DEC-002"
	now := clk.Now().UTC()
	sess, err := op.Begin(inst, op.Intent{
		Op:         "supersede",
		Type:       "decision",
		Title:      title,
		OriginalID: "DEC-001",
		LogLine:    "2026-08-15\tsupersede\tDEC-001\t" + title,
		Argv:       argv,
		OpID:       "resume-supersede",
	}, now)
	if err != nil {
		t.Fatal(err)
	}
	oldRel := findRel(t, inst, "decisions", "DEC-001")
	newRel := findRel(t, inst, "decisions", "DEC-002")
	oldData := readFile(t, filepath.Join(inst, filepath.FromSlash(oldRel)))
	newData := readFile(t, filepath.Join(inst, filepath.FromSlash(newRel)))
	oldOut, err := supersede.ApplyOld(oldData, "DEC-002")
	if err != nil {
		t.Fatal(err)
	}
	newOut, err := supersede.ApplyNew(newData, "DEC-001")
	if err != nil {
		t.Fatal(err)
	}
	if err := sess.Stage([]op.Staged{
		{RelTo: oldRel, Content: oldOut},
		{RelTo: newRel, Content: newOut},
		{RelTo: "index.md", Content: readFile(t, filepath.Join(inst, "index.md"))},
		{RelTo: "log.md", Content: append(readFile(t, filepath.Join(inst, "log.md")), []byte("2026-08-15\tsupersede\tDEC-001\t"+title+"\n")...)},
		{RelTo: "mycelium.toml", Content: readFile(t, filepath.Join(inst, "mycelium.toml"))},
	}); err != nil {
		t.Fatal(err)
	}
	if err := sess.CommitPartial(1); err != nil {
		t.Fatal(err)
	}
	_ = sess.Close()

	j, err := journal.Load(inst)
	if err != nil {
		t.Fatal(err)
	}
	if j.Op != "supersede" || j.OriginalID != "DEC-001" {
		t.Fatalf("journal=%+v", j)
	}
	info, err := lock.Inspect(inst)
	if err != nil {
		t.Fatal(err)
	}
	if info.State == lock.Live {
		t.Fatal("live lock after Close")
	}

	code, stdout, stderr := runCLI(t, clk, rec, env, work, argv...)
	if code != 0 {
		t.Fatalf("resume exit %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stdout, "mycelium supersede: ok") {
		t.Fatalf("stdout=%q", stdout)
	}
	if _, err := journal.Load(inst); err == nil {
		t.Fatal("journal should be gone after resume")
	}
	clitest.AssertNoNetwork(t, rec)
}

func mustNewDecision(t *testing.T, clk clock.Clock, rec *execrun.Recording, env map[string]string, work, inst, title string) {
	t.Helper()
	code, _, stderr := runCLI(t, clk, rec, env, work, "new", "decision", title, "--dir", inst)
	if code != 0 {
		t.Fatalf("new decision %q exit %d stderr=%q", title, code, stderr)
	}
}

func patchDECStatus(t *testing.T, inst, id, status string) {
	t.Helper()
	dir := filepath.Join(inst, "decisions")
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
		writeMeta(t, path, doc.Meta, doc.Body)
		return
	}
	t.Fatalf("no decision file for %s", id)
}

func readArtifactMeta(t *testing.T, inst, home, id string) map[string]any {
	t.Helper()
	rel := findRel(t, inst, home, id)
	doc, err := metadata.Parse(readFile(t, filepath.Join(inst, filepath.FromSlash(rel))))
	if err != nil {
		t.Fatal(err)
	}
	return doc.Meta
}

func findRel(t *testing.T, inst, home, id string) string {
	t.Helper()
	dir := filepath.Join(inst, home)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	prefix := id + "-"
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), prefix) && strings.HasSuffix(e.Name(), ".md") {
			return filepath.ToSlash(filepath.Join(home, e.Name()))
		}
	}
	t.Fatalf("no file for %s in %s", id, home)
	return ""
}

func snapshotInstance(t *testing.T, inst string) map[string][]byte {
	t.Helper()
	out := map[string][]byte{}
	err := filepath.WalkDir(inst, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			if d.Name() == ".git" || d.Name() == ".mycelium" {
				return filepath.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(inst, path)
		if err != nil {
			return err
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		out[filepath.ToSlash(rel)] = append([]byte(nil), b...)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func assertUnchanged(t *testing.T, inst string, snap map[string][]byte) {
	t.Helper()
	got := snapshotInstance(t, inst)
	if len(got) != len(snap) {
		t.Fatalf("file count %d -> %d", len(snap), len(got))
	}
	for rel, want := range snap {
		have, ok := got[rel]
		if !ok {
			t.Fatalf("missing %s after refuse", rel)
		}
		if string(have) != string(want) {
			t.Fatalf("%s mutated after refuse", rel)
		}
	}
}
