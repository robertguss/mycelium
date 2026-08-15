package clitest_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/clitest"
)

const (
	appendixBPrompt = "Should this idea use SQLite as the store? Answer independently. Do not see other reports."
	appendixBHash   = "ec87bfc2afd545807ca87b5c29cae8e77262cb3c746fc63e4539d8daeb2a77de"
	appendixCPrompt = "Review the SQLite store decision. Work independently. Do not see other reports. Retain dissent."
	appendixCHash   = "8997334f7f2f0bf821bce8ccc4a8d6cf027317c6c66d821200a032a6a11ce098"
)

func TestLadderOptInFalseHermetic(t *testing.T) {
	bin, env, workDir, home, inst := ladderScaffold(t, "Ladder OptIn")
	writeCMP(t, inst, "CMP-001-sqlite.md", cmpBody("CMP-001", "SQLite", "second-opinion", false, "cheap", "manual", appendixBPrompt))
	code, _, stderr := clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 1 {
		t.Fatalf("want exit 1, got %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "opt_in must be true") || !strings.Contains(stderr, "council-opt-in") {
		t.Fatalf("stderr=%q", stderr)
	}
	assertNoGh(t, stderr)
	assertNoHomeTouch(t, home)
}

func TestLadderCouncilCheapHermetic(t *testing.T) {
	bin, env, workDir, home, inst := ladderScaffold(t, "Ladder Cheap")
	writeCMP(t, inst, "CMP-001-sqlite.md", cmpBody("CMP-001", "SQLite", "council", true, "cheap", "manual", appendixCPrompt))
	code, _, stderr := clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 1 {
		t.Fatalf("want exit 1, got %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "cost_class") || !strings.Contains(stderr, "council-cost-class") {
		t.Fatalf("stderr=%q", stderr)
	}
	assertNoGh(t, stderr)
	assertNoHomeTouch(t, home)
}

func TestLadderSecondOpinionTwoRPTsHermetic(t *testing.T) {
	bin, env, workDir, home, inst := ladderScaffold(t, "Ladder Two RPT")
	writeCMP(t, inst, "CMP-001-sqlite.md", cmpBody("CMP-001", "SQLite", "second-opinion", true, "cheap", "manual", appendixBPrompt))
	writeRPT(t, inst, "RPT-001-a.md", rptBody("RPT-001", "A", "CMP-001", "second-opinion", "manual", appendixBHash, "none"))
	writeRPT(t, inst, "RPT-002-b.md", rptBody("RPT-002", "B", "CMP-001", "second-opinion", "manual", appendixBHash, "none"))
	code, _, stderr := clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 1 {
		t.Fatalf("want exit 1, got %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "second-opinion") {
		t.Fatalf("stderr=%q", stderr)
	}
	assertNoGh(t, stderr)
	assertNoHomeTouch(t, home)
}

func TestLadderCouncilOneRPTHermetic(t *testing.T) {
	bin, env, workDir, home, inst := ladderScaffold(t, "Ladder One RPT")
	writeCMP(t, inst, "CMP-001-sqlite.md", cmpBody("CMP-001", "SQLite", "council", true, "standard", "cursor", appendixCPrompt))
	writeRPT(t, inst, "RPT-001-a.md", rptBody("RPT-001", "A", "CMP-001", "council", "cursor", appendixCHash, "none"))
	code, _, stderr := clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 1 {
		t.Fatalf("want exit 1, got %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "council-cardinality") {
		t.Fatalf("stderr=%q", stderr)
	}
	assertNoGh(t, stderr)
	assertNoHomeTouch(t, home)
}

func TestLadderCouncilHappyHermetic(t *testing.T) {
	bin, env, workDir, home, inst := ladderScaffold(t, "Ladder Happy")
	writeCMP(t, inst, "CMP-001-sqlite.md", cmpBody("CMP-001", "SQLite", "council", true, "standard", "cursor", appendixCPrompt))
	writeRPT(t, inst, "RPT-001-a.md", rptBody("RPT-001", "A", "CMP-001", "council", "cursor", appendixCHash, "SEED-DISSENT"))
	writeRPT(t, inst, "RPT-002-b.md", rptBody("RPT-002", "B", "CMP-001", "council", "cursor", appendixCHash, "none"))
	writeRCL(t, inst, "RCL-001-sqlite.md", rclBody("RCL-001", "SQLite", "CMP-001", "SEED-DISSENT"))
	code, _, stderr := clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("want exit 0, got %d stderr=%q", code, stderr)
	}
	assertNoGh(t, stderr)
	assertNoHomeTouch(t, home)
}

func TestLadderHashMismatchHermetic(t *testing.T) {
	bin, env, workDir, home, inst := ladderScaffold(t, "Ladder Hash")
	writeCMP(t, inst, "CMP-001-sqlite.md", cmpBody("CMP-001", "SQLite", "second-opinion", true, "cheap", "manual", appendixBPrompt))
	bad := "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
	writeRPT(t, inst, "RPT-001-a.md", rptBody("RPT-001", "A", "CMP-001", "second-opinion", "manual", bad, "none"))
	code, _, stderr := clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 1 {
		t.Fatalf("want exit 1, got %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "prompt_sha256 mismatch") || !strings.Contains(stderr, appendixBHash) {
		t.Fatalf("stderr must name expected hex: %q", stderr)
	}
	if !strings.Contains(stderr, "prompt-identity") {
		t.Fatalf("stderr=%q", stderr)
	}
	assertNoGh(t, stderr)
	assertNoHomeTouch(t, home)
}

func TestLadderSeedDissentMissingHermetic(t *testing.T) {
	bin, env, workDir, home, inst := ladderScaffold(t, "Ladder Seed")
	writeCMP(t, inst, "CMP-001-sqlite.md", cmpBody("CMP-001", "SQLite", "council", true, "standard", "cursor", appendixCPrompt))
	writeRPT(t, inst, "RPT-001-a.md", rptBody("RPT-001", "A", "CMP-001", "council", "cursor", appendixCHash, "SEED-DISSENT"))
	writeRPT(t, inst, "RPT-002-b.md", rptBody("RPT-002", "B", "CMP-001", "council", "cursor", appendixCHash, "none"))
	writeRCL(t, inst, "RCL-001-sqlite.md", rclBody("RCL-001", "SQLite", "CMP-001", "none"))
	code, _, stderr := clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 1 {
		t.Fatalf("want exit 1, got %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "SEED-DISSENT") || !strings.Contains(stderr, "seeded-dissent") {
		t.Fatalf("stderr=%q", stderr)
	}
	assertNoGh(t, stderr)
	assertNoHomeTouch(t, home)
}

func TestLadderLoneCMPHermetic(t *testing.T) {
	bin, env, workDir, home, inst := ladderScaffold(t, "Ladder Lone")
	code, _, stderr := clitest.Run(t, bin, workDir, env, "new", "commissioning", "SQLite", "--dir", inst)
	if code != 0 {
		t.Fatalf("new commissioning exit %d stderr=%q", code, stderr)
	}
	code, _, stderr = clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("lone CMP want 0, got %d stderr=%q", code, stderr)
	}
	assertNoGh(t, stderr)
	assertNoHomeTouch(t, home)
}

func TestLadderNoPackNoReviewsHermetic(t *testing.T) {
	bin, env, workDir, home, inst := ladderScaffold(t, "Ladder No Pack")
	if err := os.RemoveAll(filepath.Join(inst, "program", "packs", "council")); err != nil {
		t.Fatal(err)
	}
	code, _, stderr := clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("want 0, got %d stderr=%q", code, stderr)
	}
	assertNoGh(t, stderr)
	assertNoHomeTouch(t, home)
}

func ladderScaffold(t *testing.T, title string) (bin string, env []string, workDir, home, inst string) {
	t.Helper()
	bin = clitest.Bin(t)
	stubDir := installNetworkStubs(t)
	home = t.TempDir()
	workDir = t.TempDir()
	env = hermeticEnv(stubDir, home)
	code, _, stderr := clitest.Run(t, bin, workDir, env, "new", "idea", title, "--offline", "--dir", filepath.Join(workDir, "inst"))
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr)
	}
	inst = filepath.Join(workDir, "inst")
	return bin, env, workDir, home, inst
}

func writeCMP(t *testing.T, inst, name, body string) {
	t.Helper()
	dir := filepath.Join(inst, "reviews", "commissioning")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeRPT(t *testing.T, inst, name, body string) {
	t.Helper()
	dir := filepath.Join(inst, "reviews", "reports")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeRCL(t *testing.T, inst, name, body string) {
	t.Helper()
	dir := filepath.Join(inst, "reviews", "reconciliations")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func cmpBody(id, title, rung string, optIn bool, cost, adapter, prompt string) string {
	opt := "true"
	if !optIn {
		opt = "false"
	}
	return "+++\n" +
		"id = \"" + id + "\"\n" +
		"title = \"" + title + "\"\n" +
		"date = \"2026-08-15\"\n" +
		"rung = \"" + rung + "\"\n" +
		"opt_in = " + opt + "\n" +
		"cost_class = \"" + cost + "\"\n" +
		"adapter = \"" + adapter + "\"\n" +
		"+++\n\n" +
		"# " + id + " — " + title + "\n\n" +
		"## Prompt\n\n" + prompt + "\n\n" +
		"## Attachments\n\nnone\n\n" +
		"## Cost\n\n" + cost + "\n"
}

func rptBody(id, title, cmp, rung, adapter, hash, dissent string) string {
	return "+++\n" +
		"id = \"" + id + "\"\n" +
		"title = \"" + title + "\"\n" +
		"date = \"2026-08-15\"\n" +
		"model = \"model-" + strings.ToLower(title) + "\"\n" +
		"commissioning = \"" + cmp + "\"\n" +
		"rung = \"" + rung + "\"\n" +
		"adapter = \"" + adapter + "\"\n" +
		"prompt_sha256 = \"" + hash + "\"\n" +
		"+++\n\n" +
		"# " + id + " — " + title + "\n\n" +
		"## Position\n\npos\n\n" +
		"## Findings\n\nfind\n\n" +
		"## Dissent\n\n" + dissent + "\n"
}

func rclBody(id, title, cmp, retained string) string {
	sections := []string{
		"Convergence",
		"Material disagreement",
		"Evidence unique to one report",
		"Contradictory evidence",
		"Different assumptions",
		"Different scope interpretations",
		"Recommendations independently supported",
		"Questions requiring another spike",
		"Final reconciled recommendation",
	}
	var b strings.Builder
	b.WriteString("+++\n")
	b.WriteString("id = \"" + id + "\"\n")
	b.WriteString("title = \"" + title + "\"\n")
	b.WriteString("date = \"2026-08-15\"\n")
	b.WriteString("commissioning = \"" + cmp + "\"\n")
	b.WriteString("rung = \"council\"\n")
	b.WriteString("+++\n\n")
	b.WriteString("# " + id + " — " + title + "\n\n")
	for _, s := range sections {
		b.WriteString("## " + s + "\n\nnone\n\n")
	}
	b.WriteString("## Retained dissent\n\n" + retained + "\n")
	return b.String()
}
