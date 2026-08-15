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
	"github.com/robertguss/mycelium/internal/handoff"
	"github.com/robertguss/mycelium/internal/journal"
	"github.com/robertguss/mycelium/internal/op"
)

func TestPhase06Slice2HandoffHappyPath(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1", "MYCELIUM_NOW": "2026-08-15T00:00:00Z"}

	inst := scaffoldIdea(t, work, "Handoff Happy", clk, rec, env)
	mustState(t, clk, rec, env, work, inst, "exploring")
	mustState(t, clk, rec, env, work, inst, "clarified")

	code, stdout, stderr := runCLI(t, clk, rec, env, work, "handoff", "--dir", inst)
	if code != 0 {
		t.Fatalf("handoff exit %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stdout, "mycelium handoff: ok") ||
		!strings.Contains(stdout, "state: handed-off") ||
		!strings.Contains(stdout, "packet: handoff/PACKET.md") {
		t.Fatalf("stdout=%q", stdout)
	}

	m := loadManifest(t, inst)
	if m.State != "handed-off" {
		t.Fatalf("state=%q", m.State)
	}
	packet := filepath.Join(inst, "handoff", "PACKET.md")
	if _, err := os.Stat(packet); err != nil {
		t.Fatal(err)
	}
	findings := handoff.Check(os.DirFS(filepath.Join(inst, "handoff")))
	if len(findings) > 0 {
		t.Fatalf("packet structure: %v", findings)
	}

	logBody := string(readFile(t, filepath.Join(inst, "log.md")))
	wantLine := "2026-08-15\thandoff\tHO-001\tclarified -> handed-off"
	if !strings.Contains(logBody, wantLine) {
		t.Fatalf("log missing %q\n%s", wantLine, logBody)
	}

	code, _, stderr = runCLI(t, clk, rec, env, work, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("check after handoff exit %d stderr=%q", code, stderr)
	}
	clitest.AssertNoNetwork(t, rec)
}

func TestPhase06Slice2RefuseTable(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}

	// clarified, no packet: state handed-off refuses; check still passes.
	nopkt := scaffoldIdea(t, work, "No Packet", clk, rec, env)
	mustState(t, clk, rec, env, work, nopkt, "exploring")
	mustState(t, clk, rec, env, work, nopkt, "clarified")
	snap := snapshotTree(t, nopkt)

	code, _, stderr := runCLI(t, clk, rec, env, work, "state", "handed-off", "--dir", nopkt)
	if code != 1 {
		t.Fatalf("state handed-off without packet exit %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "state=handed-off requires a handoff packet") {
		t.Fatalf("stderr=%q", stderr)
	}
	if !strings.Contains(stderr, "mycelium handoff") {
		t.Fatalf("stderr=%q", stderr)
	}
	assertTreeUnchanged(t, nopkt, snap)

	code, _, stderr = runCLI(t, clk, rec, env, work, "check", "--dir", nopkt)
	if code != 0 {
		t.Fatalf("clarified without packet check exit %d stderr=%q", code, stderr)
	}

	// handoff from exploring refuses.
	expl := scaffoldIdea(t, work, "Exploring Handoff", clk, rec, env)
	mustState(t, clk, rec, env, work, expl, "exploring")
	snap2 := snapshotTree(t, expl)
	code, _, stderr = runCLI(t, clk, rec, env, work, "handoff", "--dir", expl)
	if code != 1 {
		t.Fatalf("handoff from exploring exit %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "handoff is legal only from clarified (got exploring)") {
		t.Fatalf("stderr=%q", stderr)
	}
	assertTreeUnchanged(t, expl, snap2)

	// Happy handoff, then refuse re-handoff / exploring / handed-off again.
	ok := scaffoldIdea(t, work, "Refuse After", clk, rec, env)
	mustState(t, clk, rec, env, work, ok, "exploring")
	mustState(t, clk, rec, env, work, ok, "clarified")
	code, _, stderr = runCLI(t, clk, rec, env, work, "handoff", "--dir", ok)
	if code != 0 {
		t.Fatalf("setup handoff exit %d stderr=%q", code, stderr)
	}
	snap3 := snapshotTree(t, ok)

	code, _, stderr = runCLI(t, clk, rec, env, work, "handoff", "--dir", ok)
	if code != 1 {
		t.Fatalf("re-handoff exit %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "handoff/PACKET.md already exists") {
		t.Fatalf("stderr=%q", stderr)
	}
	if !strings.Contains(stderr, "mycelium state handed-off") {
		t.Fatalf("stderr=%q", stderr)
	}
	assertTreeUnchanged(t, ok, snap3)

	code, _, stderr = runCLI(t, clk, rec, env, work, "state", "exploring", "--dir", ok)
	if code != 1 {
		t.Fatalf("handed-off → exploring exit %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "illegal transition handed-off → exploring") {
		t.Fatalf("stderr=%q", stderr)
	}
	assertTreeUnchanged(t, ok, snap3)

	code, _, stderr = runCLI(t, clk, rec, env, work, "state", "handed-off", "--dir", ok)
	if code != 1 {
		t.Fatalf("handed-off → handed-off exit %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "illegal transition handed-off → handed-off") {
		t.Fatalf("stderr=%q", stderr)
	}

	// stored handed-off without packet → check FAIL
	bare := scaffoldIdea(t, work, "Bare Handed", clk, rec, env)
	patchManifestState(t, bare, "handed-off")
	code, _, stderr = runCLI(t, clk, rec, env, work, "check", "--dir", bare)
	if code != 1 {
		t.Fatalf("bare handed-off check exit %d want 1 stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "handoff packet") {
		t.Fatalf("stderr=%q", stderr)
	}

	// optional: archived from handed-off still legal
	code, _, stderr = runCLI(t, clk, rec, env, work, "state", "archived", "--dir", ok)
	if code != 0 {
		t.Fatalf("handed-off → archived exit %d stderr=%q", code, stderr)
	}
	if loadManifest(t, ok).State != "archived" {
		t.Fatal("want archived")
	}
	clitest.AssertNoNetwork(t, rec)
}

func TestPhase06Slice2StateHandedOffWithPacket(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}

	inst := scaffoldIdea(t, work, "State Flip", clk, rec, env)
	mustState(t, clk, rec, env, work, inst, "exploring")
	mustState(t, clk, rec, env, work, inst, "clarified")

	// Generate packet via handoff, then roll state back to clarified without deleting packet.
	code, _, stderr := runCLI(t, clk, rec, env, work, "handoff", "--dir", inst)
	if code != 0 {
		t.Fatalf("handoff exit %d stderr=%q", code, stderr)
	}
	// Manually set state back to clarified (simulates packet already present).
	patchManifestState(t, inst, "clarified")
	// Remove the handoff log line so we can assert the state path writes one.
	// Keep packet intact.
	packetBefore := readFile(t, filepath.Join(inst, "handoff", "PACKET.md"))

	code, stdout, stderr := runCLI(t, clk, rec, env, work, "state", "handed-off", "--dir", inst)
	if code != 0 {
		t.Fatalf("state handed-off exit %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stdout, "mycelium state: ok") ||
		!strings.Contains(stdout, "state: handed-off") ||
		!strings.Contains(stdout, "packet: handoff/PACKET.md") {
		t.Fatalf("stdout=%q", stdout)
	}
	packetAfter := readFile(t, filepath.Join(inst, "handoff", "PACKET.md"))
	if string(packetBefore) != string(packetAfter) {
		t.Fatal("state handed-off must not regenerate packet")
	}
	if loadManifest(t, inst).State != "handed-off" {
		t.Fatal("want handed-off")
	}
	code, _, stderr = runCLI(t, clk, rec, env, work, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("check exit %d stderr=%q", code, stderr)
	}
	clitest.AssertNoNetwork(t, rec)
}

func TestPhase06Slice2HandoffResumeJournal(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}

	inst := scaffoldIdea(t, work, "Resume Handoff", clk, rec, env)
	mustState(t, clk, rec, env, work, inst, "exploring")
	mustState(t, clk, rec, env, work, inst, "clarified")

	argv := []string{"handoff", "--dir", inst}
	title := "clarified -> handed-off"
	now := clk.Now().UTC()
	sess, err := op.Begin(inst, op.Intent{
		Op:      "handoff",
		Title:   title,
		LogLine: "2026-08-15\thandoff\tHO-001\t" + title,
		Argv:    argv,
	}, now)
	if err != nil {
		t.Fatal(err)
	}
	// Stage a minimal partial set then leave journal (simulate crash before finish).
	packet := []byte(`+++
id = "HO-001"
date = "2026-08-15"
implementation_system = "pstack/poteto"
time_budget = "30m"
+++

# Handoff packet

## Framing

resume

## Locked decisions

none

## Glossary

See [glossary.md](glossary.md).

## Open questions

none

## Evidence summary

none

## Implementation playbooks

- see [playbooks/PLAYBOOK.md](playbooks/PLAYBOOK.md)

## Implementation system

pstack/poteto

## Time budget

30m

## Acceptance

- see [acceptance/README.md](acceptance/README.md)
`)
	files := []op.Staged{
		{RelTo: "handoff/PACKET.md", Content: packet},
		{RelTo: "handoff/glossary.md", Content: []byte("none\n")},
		{RelTo: "handoff/decisions/.keep", Content: []byte{}},
		{RelTo: "handoff/questions/.keep", Content: []byte{}},
		{RelTo: "handoff/evidence/SUMMARY.md", Content: []byte("none\n")},
		{RelTo: "handoff/playbooks/PLAYBOOK.md", Content: []byte("# Playbook\n\n## Target\n\nnone\n\n## Steps\n\nnone\n\n## Done\n\nnone\n")},
		{RelTo: "handoff/acceptance/README.md", Content: []byte("none\n")},
		{RelTo: "index.md", Content: readFile(t, filepath.Join(inst, "index.md"))},
		{RelTo: "log.md", Content: append(readFile(t, filepath.Join(inst, "log.md")), []byte("2026-08-15\thandoff\tHO-001\tclarified -> handed-off\n")...)},
		{RelTo: "mycelium.toml", Content: []byte(strings.Replace(string(readFile(t, filepath.Join(inst, "mycelium.toml"))), `state = "clarified"`, `state = "handed-off"`, 1))},
	}
	// Also try single quotes
	if !strings.Contains(string(files[len(files)-1].Content), "handed-off") {
		files[len(files)-1].Content = []byte(strings.Replace(string(readFile(t, filepath.Join(inst, "mycelium.toml"))), `state = 'clarified'`, `state = 'handed-off'`, 1))
	}
	if err := sess.Stage(files); err != nil {
		t.Fatal(err)
	}
	if err := sess.CommitPartial(2); err != nil {
		t.Fatal(err)
	}
	_ = sess.Close()
	if _, err := journal.Load(inst); err != nil {
		t.Fatalf("want leftover journal: %v", err)
	}

	code, stdout, stderr := runCLI(t, clk, rec, env, work, "handoff", "--dir", inst)
	if code != 0 {
		t.Fatalf("resume handoff exit %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stdout, "mycelium handoff: ok") {
		t.Fatalf("stdout=%q", stdout)
	}
	if _, err := journal.Load(inst); err == nil {
		t.Fatal("journal should be gone after resume")
	}
	if loadManifest(t, inst).State != "handed-off" {
		t.Fatal("want handed-off after resume")
	}
	code, _, stderr = runCLI(t, clk, rec, env, work, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("check after resume exit %d stderr=%q", code, stderr)
	}
	clitest.AssertNoNetwork(t, rec)
}

func TestPhase06Slice2ExtraArgsRefuse(t *testing.T) {
	work := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	clk := clock.Fixed{T: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)}
	env := map[string]string{"MYCELIUM_OFFLINE": "1"}
	inst := scaffoldIdea(t, work, "Extra Args", clk, rec, env)
	mustState(t, clk, rec, env, work, inst, "exploring")
	mustState(t, clk, rec, env, work, inst, "clarified")

	code, _, stderr := runCLI(t, clk, rec, env, work, "handoff", "extra", "--dir", inst)
	if code != 1 {
		t.Fatalf("exit %d want 1 stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "handoff accepts only --dir PATH") {
		t.Fatalf("stderr=%q", stderr)
	}
}

func snapshotTree(t *testing.T, inst string) map[string][]byte {
	t.Helper()
	return snapshotInstance(t, inst)
}

func assertTreeUnchanged(t *testing.T, inst string, snap map[string][]byte) {
	t.Helper()
	assertUnchanged(t, inst, snap)
}
