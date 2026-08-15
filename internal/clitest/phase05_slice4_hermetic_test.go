package clitest_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/clitest"
	"github.com/robertguss/mycelium/internal/release"
	"github.com/robertguss/mycelium/internal/version"
)

// PHASE-05 Slice 4 — CHANGELOG + release.sh + install doc (MS-501-REL pieces).
// No live GitHub Release, no tag, no Actions, no gh.

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Dir(file)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found")
		}
		dir = parent
	}
}

func TestPhase05Slice4ChangelogFixture(t *testing.T) {
	root := repoRoot(t)
	body, err := os.ReadFile(filepath.Join(root, "CHANGELOG.md"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	if !strings.Contains(text, "## [Unreleased]") {
		t.Fatal("CHANGELOG.md missing ## [Unreleased]")
	}
	if !strings.Contains(text, "## [0.1.0] - 2026-08-15") {
		t.Fatal("CHANGELOG.md missing ## [0.1.0] - 2026-08-15 fixture heading")
	}
	if !release.HasVersionHeading(text, "0.1.0") {
		t.Fatal("HasVersionHeading(0.1.0) false on fixture CHANGELOG")
	}
	if release.HasVersionHeading(text, "9.9.9") {
		t.Fatal("HasVersionHeading(9.9.9) unexpectedly true")
	}
}

func TestPhase05Slice4VersionStillDev(t *testing.T) {
	if version.Version != "0.1.0-dev" {
		t.Fatalf("version.Version = %q, want 0.1.0-dev", version.Version)
	}
	root := repoRoot(t)
	src, err := os.ReadFile(filepath.Join(root, "internal", "version", "version.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(src), `Version = "0.1.0-dev"`) {
		t.Fatalf("committed version.go must keep 0.1.0-dev:\n%s", src)
	}
}

func TestPhase05Slice4InstallDoc(t *testing.T) {
	root := repoRoot(t)
	path := filepath.Join(root, "docs", "install.md")
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	if !strings.Contains(text, "releases/latest/download/mycelium-") {
		t.Fatal("docs/install.md missing releases/latest/download/mycelium-")
	}
	if !strings.Contains(text, "~/.local/bin/mycelium") {
		t.Fatal("docs/install.md missing ~/.local/bin/mycelium")
	}
	oneLiner := `curl -fsSL https://github.com/robertguss/mycelium/releases/latest/download/mycelium-$(uname -s | tr A-Z a-z)-$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/') -o ~/.local/bin/mycelium && chmod +x ~/.local/bin/mycelium`
	if !strings.Contains(text, oneLiner) {
		t.Fatal("docs/install.md missing verbatim install one-liner")
	}
}

func TestPhase05Slice4ReleaseHappy(t *testing.T) {
	root := repoRoot(t)
	dist := filepath.Join(root, "dist")
	t.Cleanup(func() {
		_ = os.RemoveAll(dist)
	})
	_ = os.RemoveAll(dist)

	cmd := exec.Command("bash", "scripts/release.sh", "0.1.0")
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0", "PATH="+os.Getenv("PATH"))
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("release.sh 0.1.0 failed: %v\n%s", err, out)
	}

	linuxBin := filepath.Join(dist, "mycelium-linux-amd64")
	darwinBin := filepath.Join(dist, "mycelium-darwin-arm64")
	sumsPath := filepath.Join(dist, "SHA256SUMS")
	for _, p := range []string{linuxBin, darwinBin, sumsPath} {
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("missing %s: %v\n%s", p, err, out)
		}
	}

	sumsBody, err := os.ReadFile(sumsPath)
	if err != nil {
		t.Fatal(err)
	}
	entries, err := release.ParseSHA256SUMS(string(sumsBody))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("SHA256SUMS entries=%d want 2\n%s", len(entries), sumsBody)
	}
	names := map[string]bool{}
	for _, e := range entries {
		names[e.Name] = true
	}
	if !names["mycelium-linux-amd64"] || !names["mycelium-darwin-arm64"] {
		t.Fatalf("SHA256SUMS names=%v", names)
	}
	if err := release.MatchSHA256SUMS(dist, entries); err != nil {
		t.Fatalf("checksum mismatch: %v", err)
	}

	// Committed source version must remain 0.1.0-dev even after ldflags stamp.
	if version.Version != "0.1.0-dev" {
		t.Fatalf("in-tree version mutated to %q", version.Version)
	}
}

func TestPhase05Slice4ReleaseRefuse(t *testing.T) {
	root := repoRoot(t)
	tmp := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmp, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	script, err := os.ReadFile(filepath.Join(root, "scripts", "release.sh"))
	if err != nil {
		t.Fatal(err)
	}
	dstScript := filepath.Join(tmp, "scripts", "release.sh")
	if err := os.WriteFile(dstScript, script, 0o755); err != nil {
		t.Fatal(err)
	}
	changelog := "# Changelog\n\n## [Unreleased]\n"
	if err := os.WriteFile(filepath.Join(tmp, "CHANGELOG.md"), []byte(changelog), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("bash", "scripts/release.sh", "9.9.9")
	cmd.Dir = tmp
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected refuse for 9.9.9, got success:\n%s", out)
	}
	combined := string(out)
	if !strings.Contains(combined, "## [9.9.9]") {
		t.Fatalf("refuse must name missing heading ## [9.9.9]:\n%s", combined)
	}
	if !strings.Contains(combined, "release refused") {
		t.Fatalf("refuse must say release refused:\n%s", combined)
	}

	dist := filepath.Join(tmp, "dist")
	if entries, _ := os.ReadDir(dist); len(entries) > 0 {
		t.Fatalf("refuse must write nothing under dist/; got %v", entries)
	}
	if _, err := os.Stat(filepath.Join(dist, "SHA256SUMS")); !os.IsNotExist(err) {
		t.Fatalf("SHA256SUMS must not exist after refuse: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dist, "mycelium-linux-amd64")); !os.IsNotExist(err) {
		t.Fatalf("linux binary must not exist after refuse: %v", err)
	}
}

func TestPhase05Slice4NoReleaseOrInstallVerb(t *testing.T) {
	bin := clitest.Bin(t)
	work := t.TempDir()
	for _, verb := range []string{"release", "install"} {
		code, _, stderr := clitest.Run(t, bin, work, nil, verb)
		if code == 0 {
			t.Fatalf("%s must be unknown command, got exit 0", verb)
		}
		if !strings.Contains(stderr, "unknown command") {
			t.Fatalf("%s stderr=%q, want unknown command", verb, stderr)
		}
	}
}

func TestPhase05Slice4NoPhase05Workflow(t *testing.T) {
	root := repoRoot(t)
	wf := filepath.Join(root, ".github", "workflows")
	entries, err := os.ReadDir(wf)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, "phase-05-") && strings.HasSuffix(name, ".yml") {
			t.Fatalf("must not add PHASE-05 workflow: %s", name)
		}
	}
}
