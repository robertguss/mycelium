package tiercmd_test

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
	inst := filepath.Join(cwd, slug)
	if _, err := os.Stat(filepath.Join(inst, "mycelium.toml")); err != nil {
		t.Fatalf("instance missing at %s: %v", inst, err)
	}
	return inst
}

func runCLI(t *testing.T, deps cli.Deps, args ...string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := cli.Run(append([]string{"mycelium"}, args...), &stdout, &stderr, deps)
	return code, stdout.String(), stderr.String()
}

func readManifest(t *testing.T, inst string) manifest.Manifest {
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

func countTierLogLines(t *testing.T, inst string) int {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(inst, "log.md"))
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	for _, line := range strings.Split(string(b), "\n") {
		if strings.Contains(line, "\ttier\t") {
			n++
		}
	}
	return n
}

func mustDir(t *testing.T, inst, rel string) {
	t.Helper()
	fi, err := os.Stat(filepath.Join(inst, rel))
	if err != nil {
		t.Fatalf("missing %s: %v", rel, err)
	}
	if !fi.IsDir() {
		t.Fatalf("%s is not a dir", rel)
	}
}

func mustREADME(t *testing.T, inst, dir, ns string) {
	t.Helper()
	path := filepath.Join(inst, dir, "README.md")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("README missing under %s: %v", dir, err)
	}
	body := string(b)
	if !strings.Contains(body, "Home for "+ns+"-### artifacts.") {
		t.Fatalf("%s README body=%q", dir, body)
	}
}

func assertTeaching(t *testing.T, errText string) {
	t.Helper()
	for _, prefix := range []string{"mycelium:", "convention:", "contract:", "fix:"} {
		if !strings.Contains(errText, prefix) {
			t.Fatalf("teaching error missing %q in %q", prefix, errText)
		}
	}
}

var standardDirs = []string{"decisions", "assumptions", "evidence", "questions", "risks"}
var haExtraDirs = []string{"spikes", "findings", "recommendations", "requirements", "phases", "milestones"}

func TestFocusedToStandardEmitsFiveDirsKeepsExistingDEC(t *testing.T) {
	cwd := t.TempDir()
	inst := scaffoldOffline(t, cwd, "Tier Up")
	deps := fixedDeps(t, cwd)

	// Extra artifact at focused is allowed; generate one DEC before promote.
	code, _, errText := runCLI(t, deps, "new", "decision", "Keep Me", "--dir", inst)
	if code != 0 {
		t.Fatalf("new decision exit %d stderr=%q", code, errText)
	}
	decPath := filepath.Join(inst, "decisions", "DEC-001-keep-me.md")
	before, err := os.ReadFile(decPath)
	if err != nil {
		t.Fatal(err)
	}

	code, out, errText := runCLI(t, deps, "tier", "standard", "--dir", inst)
	if code != 0 {
		t.Fatalf("tier standard exit %d stderr=%q", code, errText)
	}
	if !strings.Contains(out, "focused -> standard") {
		t.Fatalf("stdout=%q", out)
	}

	m := readManifest(t, inst)
	if m.Tier != "standard" {
		t.Fatalf("tier=%q", m.Tier)
	}
	if m.UpdatedDate != "2026-08-15" {
		t.Fatalf("updated_date=%q", m.UpdatedDate)
	}
	for _, d := range standardDirs {
		mustDir(t, inst, d)
	}
	mustREADME(t, inst, "assumptions", "ASM")
	mustREADME(t, inst, "evidence", "EVD")
	mustREADME(t, inst, "questions", "OQ")
	mustREADME(t, inst, "risks", "RSK")

	after, err := os.ReadFile(decPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("existing DEC was rewritten")
	}
	if countTierLogLines(t, inst) != 1 {
		t.Fatalf("want 1 tier log line, got %d", countTierLogLines(t, inst))
	}
}

func TestStandardToFocusedDeletesNothingCheckGreen(t *testing.T) {
	cwd := t.TempDir()
	inst := scaffoldOffline(t, cwd, "Tier Down")
	deps := fixedDeps(t, cwd)

	code, _, errText := runCLI(t, deps, "tier", "standard", "--dir", inst)
	if code != 0 {
		t.Fatalf("tier standard exit %d stderr=%q", code, errText)
	}
	code, _, errText = runCLI(t, deps, "new", "decision", "Survive Lower", "--dir", inst)
	if code != 0 {
		t.Fatalf("new decision exit %d stderr=%q", code, errText)
	}
	decPath := filepath.Join(inst, "decisions", "DEC-001-survive-lower.md")
	before, err := os.ReadFile(decPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range standardDirs {
		mustDir(t, inst, d)
	}

	code, out, errText := runCLI(t, deps, "tier", "focused", "--dir", inst)
	if code != 0 {
		t.Fatalf("tier focused exit %d stderr=%q", code, errText)
	}
	if !strings.Contains(out, "standard -> focused") {
		t.Fatalf("stdout=%q", out)
	}
	m := readManifest(t, inst)
	if m.Tier != "focused" {
		t.Fatalf("tier=%q", m.Tier)
	}
	for _, d := range standardDirs {
		mustDir(t, inst, d)
	}
	after, err := os.ReadFile(decPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("DEC deleted or rewritten on lower")
	}

	code, cout, cerr := runCLI(t, deps, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("check exit %d stderr=%q", code, cerr)
	}
	if !strings.Contains(cout, "mycelium check: ok") {
		t.Fatalf("check stdout=%q", cout)
	}
}

func TestFocusedToHighAssuranceEmitsAllDirs(t *testing.T) {
	cwd := t.TempDir()
	inst := scaffoldOffline(t, cwd, "Tier HA")
	deps := fixedDeps(t, cwd)

	code, out, errText := runCLI(t, deps, "tier", "high-assurance", "--dir", inst)
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, errText)
	}
	if !strings.Contains(out, "focused -> high-assurance") {
		t.Fatalf("stdout=%q", out)
	}
	for _, d := range standardDirs {
		mustDir(t, inst, d)
	}
	for _, d := range haExtraDirs {
		mustDir(t, inst, d)
	}
	mustREADME(t, inst, "spikes", "SPK")
	mustREADME(t, inst, "findings", "FND")
	mustREADME(t, inst, "milestones", "MS")
	m := readManifest(t, inst)
	if m.Tier != "high-assurance" {
		t.Fatalf("tier=%q", m.Tier)
	}
}

func TestIdempotentSecondCallAlreadyNoExtraLog(t *testing.T) {
	cwd := t.TempDir()
	inst := scaffoldOffline(t, cwd, "Tier Idem")
	deps := fixedDeps(t, cwd)

	code, _, errText := runCLI(t, deps, "tier", "standard", "--dir", inst)
	if code != 0 {
		t.Fatalf("first tier exit %d stderr=%q", code, errText)
	}
	logBefore, err := os.ReadFile(filepath.Join(inst, "log.md"))
	if err != nil {
		t.Fatal(err)
	}
	manBefore, err := os.ReadFile(filepath.Join(inst, "mycelium.toml"))
	if err != nil {
		t.Fatal(err)
	}
	readmeBefore, err := os.ReadFile(filepath.Join(inst, "decisions", "README.md"))
	if err != nil {
		t.Fatal(err)
	}

	code, out, errText := runCLI(t, deps, "tier", "standard", "--dir", inst)
	if code != 0 {
		t.Fatalf("second tier exit %d stderr=%q", code, errText)
	}
	if strings.TrimSpace(out) != "already standard" {
		t.Fatalf("stdout=%q want already standard", out)
	}
	logAfter, _ := os.ReadFile(filepath.Join(inst, "log.md"))
	manAfter, _ := os.ReadFile(filepath.Join(inst, "mycelium.toml"))
	readmeAfter, _ := os.ReadFile(filepath.Join(inst, "decisions", "README.md"))
	if !bytes.Equal(logBefore, logAfter) {
		t.Fatal("log rewritten on idempotent call")
	}
	if !bytes.Equal(manBefore, manAfter) {
		t.Fatal("manifest rewritten on idempotent call")
	}
	if !bytes.Equal(readmeBefore, readmeAfter) {
		t.Fatal("README rewritten on idempotent call")
	}
	if countTierLogLines(t, inst) != 1 {
		t.Fatalf("want 1 tier log line, got %d", countTierLogLines(t, inst))
	}
}

func TestAlreadyTierRefusesLeftoverJournal(t *testing.T) {
	cwd := t.TempDir()
	inst := scaffoldOffline(t, cwd, "Already Journal")
	deps := fixedDeps(t, cwd)

	code, _, errText := runCLI(t, deps, "tier", "standard", "--dir", inst)
	if code != 0 {
		t.Fatalf("tier standard exit %d stderr=%q", code, errText)
	}
	if err := journal.Save(inst, &journal.Journal{
		SchemaVersion: 1,
		Op:            "new",
		Title:         "Interrupted",
		OriginalID:    "DEC-001",
		StartedAt:     "2026-08-15T12:00:00Z",
		StagedDir:     ".mycelium/stage-test",
		Argv:          []string{"new", "decision", "Interrupted"},
	}); err != nil {
		t.Fatal(err)
	}

	code, out, errText := runCLI(t, deps, "tier", "standard", "--dir", inst)
	if code != 1 {
		t.Fatalf("exit %d want 1 stdout=%q stderr=%q", code, out, errText)
	}
	if strings.Contains(out, "already") {
		t.Fatalf("noop must not succeed with leftover journal: %q", out)
	}
	if !strings.Contains(errText, "leftover journal") {
		t.Fatalf("stderr=%q", errText)
	}
}

func TestAlreadyTierRefusesLiveLock(t *testing.T) {
	cwd := t.TempDir()
	inst := scaffoldOffline(t, cwd, "Already Lock")
	deps := fixedDeps(t, cwd)

	code, _, errText := runCLI(t, deps, "tier", "standard", "--dir", inst)
	if code != 0 {
		t.Fatalf("tier standard exit %d stderr=%q", code, errText)
	}
	held, err := lock.Acquire(inst, time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = held.Release() })

	code, out, errText := runCLI(t, deps, "tier", "standard", "--dir", inst)
	if code != 1 {
		t.Fatalf("exit %d want 1 stdout=%q stderr=%q", code, out, errText)
	}
	if strings.Contains(out, "already") {
		t.Fatalf("noop must not succeed with live lock: %q", out)
	}
	if !strings.Contains(errText, "lock held") {
		t.Fatalf("stderr=%q", errText)
	}
}

func TestSameTierMissingEmitDirRepairs(t *testing.T) {
	cwd := t.TempDir()
	inst := scaffoldOffline(t, cwd, "Tier Repair")
	deps := fixedDeps(t, cwd)

	code, _, errText := runCLI(t, deps, "tier", "standard", "--dir", inst)
	if code != 0 {
		t.Fatalf("tier standard exit %d stderr=%q", code, errText)
	}
	if err := os.RemoveAll(filepath.Join(inst, "risks")); err != nil {
		t.Fatal(err)
	}
	// Freeze updated_date to an older day so repair bump is observable.
	m := readManifest(t, inst)
	m.UpdatedDate = "2026-08-01"
	out, err := manifest.Encode(m)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inst, "mycelium.toml"), out, 0o644); err != nil {
		t.Fatal(err)
	}
	beforeLines := countTierLogLines(t, inst)

	code, stdout, errText := runCLI(t, deps, "tier", "standard", "--dir", inst)
	if code != 0 {
		t.Fatalf("repair exit %d stderr=%q", code, errText)
	}
	if strings.Contains(stdout, "already") {
		t.Fatalf("expected repair, got %q", stdout)
	}
	mustDir(t, inst, "risks")
	mustREADME(t, inst, "risks", "RSK")
	m2 := readManifest(t, inst)
	if m2.UpdatedDate != "2026-08-15" {
		t.Fatalf("updated_date=%q want 2026-08-15", m2.UpdatedDate)
	}
	if m2.Tier != "standard" {
		t.Fatalf("tier=%q", m2.Tier)
	}
	if countTierLogLines(t, inst) != beforeLines+1 {
		t.Fatalf("want %d tier log lines, got %d", beforeLines+1, countTierLogLines(t, inst))
	}
}

func TestUnknownTierTeachingError(t *testing.T) {
	cwd := t.TempDir()
	inst := scaffoldOffline(t, cwd, "Tier Unknown")
	deps := fixedDeps(t, cwd)

	code, _, errText := runCLI(t, deps, "tier", "enterprise", "--dir", inst)
	if code != 1 {
		t.Fatalf("exit %d want 1", code)
	}
	assertTeaching(t, errText)
	if !strings.Contains(errText, "unknown tier") {
		t.Fatalf("stderr=%q", errText)
	}
	for _, name := range []string{"focused", "standard", "high-assurance"} {
		if !strings.Contains(errText, name) {
			t.Fatalf("stderr missing %q: %q", name, errText)
		}
	}
}

func TestTierDirWithoutChdir(t *testing.T) {
	cwd := t.TempDir()
	inst := scaffoldOffline(t, cwd, "Tier Dir")
	// Getwd is a unrelated temp dir; --dir points at instance.
	other := t.TempDir()
	deps := fixedDeps(t, other)

	code, out, errText := runCLI(t, deps, "tier", "standard", "--dir", inst)
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, errText)
	}
	if !strings.Contains(out, "focused -> standard") {
		t.Fatalf("stdout=%q", out)
	}
	for _, d := range standardDirs {
		mustDir(t, inst, d)
	}
}

func TestNeverOverwriteExistingREADME(t *testing.T) {
	cwd := t.TempDir()
	inst := scaffoldOffline(t, cwd, "Tier Keepme")
	deps := fixedDeps(t, cwd)

	if err := os.MkdirAll(filepath.Join(inst, "decisions"), 0o755); err != nil {
		t.Fatal(err)
	}
	custom := []byte("# Custom Decisions\n\nDo not touch.\n")
	if err := os.WriteFile(filepath.Join(inst, "decisions", "README.md"), custom, 0o644); err != nil {
		t.Fatal(err)
	}

	code, _, errText := runCLI(t, deps, "tier", "standard", "--dir", inst)
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, errText)
	}
	after, err := os.ReadFile(filepath.Join(inst, "decisions", "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(custom, after) {
		t.Fatalf("README overwritten: %q", after)
	}
	for _, d := range []string{"assumptions", "evidence", "questions", "risks"} {
		mustDir(t, inst, d)
	}
}
