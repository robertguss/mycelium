package clitest_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/clitest"
	"github.com/robertguss/mycelium/internal/release"
	"github.com/robertguss/mycelium/internal/version"
)

// TestMS501HermeticMatrix is the PHASE-05 MS-501 gate list (Appendix E).
// Slices 2–4 already cover pieces; this file is the single named matrix.
// Execs the binary. No network, no gh, no Actions, no live release.
func TestMS501HermeticMatrix(t *testing.T) {
	bin := clitest.Bin(t)
	stubDir, ghMarker := installGhNeverStub(t)
	home := t.TempDir()
	env := append(hermeticEnv(stubDir, home),
		"MYCELIUM_NOW=2026-08-15T00:00:00Z",
		"GH_TOKEN=",
		"MYCELIUM_CONFIG="+filepath.Join(home, "config"),
	)

	run := func(t *testing.T, work string, args ...string) (int, string, string) {
		t.Helper()
		code, stdout, stderr := clitest.Run(t, bin, work, env, args...)
		assertNoGh(t, stderr)
		if _, err := os.Stat(ghMarker); err == nil {
			t.Fatal("gh was invoked")
		} else if !os.IsNotExist(err) {
			t.Fatalf("stat gh marker: %v", err)
		}
		assertNoHomeTouch(t, home)
		return code, stdout, stderr
	}

	t.Run("MS-501-SUP", func(t *testing.T) {
		work := t.TempDir()
		inst := filepath.Join(work, "sup")
		if _, err := os.Stat(inst); !os.IsNotExist(err) {
			t.Fatalf("--dir path must not exist yet: %v", err)
		}

		code, _, stderr := run(t, work, "new", "idea", "Supersede Fixture", "--offline", "--dir", inst)
		if code != 0 {
			t.Fatalf("new idea exit %d stderr=%q", code, stderr)
		}
		stateBefore := loadManifest(t, inst).State

		code, _, stderr = run(t, work, "new", "decision", "Use SQLite", "--dir", inst)
		if code != 0 {
			t.Fatalf("new decision 1 exit %d stderr=%q", code, stderr)
		}
		code, _, stderr = run(t, work, "new", "decision", "Use SQLite with WAL", "--dir", inst)
		if code != 0 {
			t.Fatalf("new decision 2 exit %d stderr=%q", code, stderr)
		}
		patchDECStatus(t, inst, "DEC-001", "Accepted")
		patchDECStatus(t, inst, "DEC-002", "Accepted")

		code, stdout, stderr := run(t, work, "supersede", "DEC-001", "--by", "DEC-002", "--dir", inst)
		if code != 0 {
			t.Fatalf("supersede exit %d stderr=%q", code, stderr)
		}
		if !strings.Contains(stdout, "mycelium supersede: ok") {
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

		logBody := string(readFile(t, filepath.Join(inst, "log.md")))
		wantLine := "2026-08-15\tsupersede\tDEC-001\tDEC-001 -> DEC-002"
		if !strings.Contains(logBody, wantLine) {
			t.Fatalf("log missing %q\n%s", wantLine, logBody)
		}

		if after := loadManifest(t, inst).State; after != stateBefore {
			t.Fatalf("state changed: before=%q after=%q", stateBefore, after)
		}

		code, _, stderr = run(t, work, "check", "--dir", inst)
		if code != 0 {
			t.Fatalf("check after supersede exit %d stderr=%q", code, stderr)
		}
	})

	t.Run("MS-501-G0", func(t *testing.T) {
		work := t.TempDir()
		inst := filepath.Join(work, "g0")
		code, _, stderr := run(t, work, "new", "idea", "Legacy Control", "--offline", "--dir", inst)
		if code != 0 {
			t.Fatalf("scaffold g0 exit %d stderr=%q", code, stderr)
		}
		overlayMycelium(t, inst, filepath.Join(legacyTestdata(t), "g0", "mycelium.toml"))

		code, _, stderr = run(t, work, "check", "--dir", inst)
		if code != 0 {
			t.Fatalf("g0 check exit %d stderr=%q", code, stderr)
		}
		code, stdout, stderr := run(t, work, "status", "--offline", "--dir", inst)
		if code != 0 {
			t.Fatalf("g0 status exit %d stderr=%q", code, stderr)
		}
		if !strings.Contains(stdout, "slug: legacy-control") {
			t.Fatalf("g0 status stdout=%q", stdout)
		}
	})

	t.Run("MS-501-G1", func(t *testing.T) {
		work := t.TempDir()
		inst := filepath.Join(work, "g1")
		code, _, stderr := run(t, work, "new", "idea", "Legacy One", "--offline", "--dir", inst)
		if code != 0 {
			t.Fatalf("scaffold g1 exit %d stderr=%q", code, stderr)
		}
		td := legacyTestdata(t)
		overlayMycelium(t, inst, filepath.Join(td, "g1", "mycelium.toml"))
		contractSrc := filepath.Join(td, "g1", "program", "contracts", "manifest.md")
		contractDst := filepath.Join(inst, "program", "contracts", "manifest.md")
		b, err := os.ReadFile(contractSrc)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(contractDst, b, 0o644); err != nil {
			t.Fatal(err)
		}

		code, _, stderr = run(t, work, "check", "--dir", inst)
		if code != 0 {
			t.Fatalf("g1 check exit %d stderr=%q", code, stderr)
		}
		code, stdout, stderr := run(t, work, "status", "--offline", "--dir", inst)
		if code != 0 {
			t.Fatalf("g1 status exit %d stderr=%q", code, stderr)
		}
		if !strings.Contains(stdout, "slug: legacy-one") {
			t.Fatalf("g1 status stdout=%q", stdout)
		}
	})

	t.Run("MS-501-G2", func(t *testing.T) {
		work := t.TempDir()
		root := filepath.Join(work, "ideas")
		if err := os.MkdirAll(root, 0o755); err != nil {
			t.Fatal(err)
		}
		g2 := filepath.Join(root, "g2-legacy-master")
		if err := os.MkdirAll(g2, 0o755); err != nil {
			t.Fatal(err)
		}
		src := filepath.Join(legacyTestdata(t), "g2", "research-program.toml")
		b, err := os.ReadFile(src)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(g2, "research-program.toml"), b, 0o644); err != nil {
			t.Fatal(err)
		}

		code, stdout, stderr := run(t, work, "status", "--all", "--offline", "--root", root)
		if code != 0 {
			t.Fatalf("g2 status --all exit %d stderr=%q", code, stderr)
		}
		combined := stdout + stderr
		if !strings.Contains(combined, "partial: legacy-manifest") {
			t.Fatalf("missing partial: legacy-manifest: %q", combined)
		}
		if !strings.Contains(combined, "research-program.toml without mycelium.toml") {
			t.Fatalf("missing g2 reason: %q", combined)
		}
		for _, line := range strings.Split(stdout, "\n") {
			if strings.Count(line, "\t") >= 4 && !strings.HasPrefix(line, "partial:") && strings.Contains(line, "g2") {
				t.Fatalf("g2 listed as idea row: %q", line)
			}
		}
		if !strings.Contains(stdout, "0 ideas (") {
			t.Fatalf("want zero ideas, got %q", stdout)
		}
	})

	t.Run("MS-501-G3", func(t *testing.T) {
		// G3 is status-only: status exit 0, check exit 1. Do not assert check 0.
		work := t.TempDir()
		inst := filepath.Join(work, "g3")
		code, _, stderr := run(t, work, "new", "idea", "Legacy Extra", "--offline", "--dir", inst)
		if code != 0 {
			t.Fatalf("scaffold g3 exit %d stderr=%q", code, stderr)
		}
		overlayMycelium(t, inst, filepath.Join(legacyTestdata(t), "g3", "mycelium.toml"))

		code, stdout, stderr := run(t, work, "status", "--offline", "--dir", inst)
		if code != 0 {
			t.Fatalf("g3 status exit %d stderr=%q", code, stderr)
		}
		if !strings.Contains(stdout, "slug: legacy-extra") {
			t.Fatalf("g3 status stdout=%q", stdout)
		}

		code, _, stderr = run(t, work, "check", "--dir", inst)
		if code != 1 {
			t.Fatalf("g3 check exit %d want 1 (not a check-pass golden) stderr=%q", code, stderr)
		}
	})

	t.Run("MS-501-G0-G1-G2-G3-root", func(t *testing.T) {
		// Appendix E recipe: all four under one --root; G0/G1/G3 listed; G2 not.
		work := t.TempDir()
		root := filepath.Join(work, "ideas")
		if err := os.MkdirAll(root, 0o755); err != nil {
			t.Fatal(err)
		}
		td := legacyTestdata(t)

		g0 := filepath.Join(root, "g0")
		code, _, stderr := run(t, work, "new", "idea", "Legacy Control", "--offline", "--dir", g0)
		if code != 0 {
			t.Fatalf("g0 scaffold exit %d stderr=%q", code, stderr)
		}
		overlayMycelium(t, g0, filepath.Join(td, "g0", "mycelium.toml"))

		g1 := filepath.Join(root, "g1")
		code, _, stderr = run(t, work, "new", "idea", "Legacy One", "--offline", "--dir", g1)
		if code != 0 {
			t.Fatalf("g1 scaffold exit %d stderr=%q", code, stderr)
		}
		overlayMycelium(t, g1, filepath.Join(td, "g1", "mycelium.toml"))
		b, err := os.ReadFile(filepath.Join(td, "g1", "program", "contracts", "manifest.md"))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(g1, "program", "contracts", "manifest.md"), b, 0o644); err != nil {
			t.Fatal(err)
		}

		g2 := filepath.Join(root, "g2")
		if err := os.MkdirAll(g2, 0o755); err != nil {
			t.Fatal(err)
		}
		b, err = os.ReadFile(filepath.Join(td, "g2", "research-program.toml"))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(g2, "research-program.toml"), b, 0o644); err != nil {
			t.Fatal(err)
		}

		g3 := filepath.Join(root, "g3")
		code, _, stderr = run(t, work, "new", "idea", "Legacy Extra", "--offline", "--dir", g3)
		if code != 0 {
			t.Fatalf("g3 scaffold exit %d stderr=%q", code, stderr)
		}
		overlayMycelium(t, g3, filepath.Join(td, "g3", "mycelium.toml"))

		code, _, stderr = run(t, work, "check", "--dir", g0)
		if code != 0 {
			t.Fatalf("root g0 check exit %d stderr=%q", code, stderr)
		}
		code, _, stderr = run(t, work, "check", "--dir", g1)
		if code != 0 {
			t.Fatalf("root g1 check exit %d stderr=%q", code, stderr)
		}
		code, _, stderr = run(t, work, "check", "--dir", g3)
		if code != 1 {
			t.Fatalf("root g3 check exit %d want 1 stderr=%q", code, stderr)
		}

		code, stdout, stderr := run(t, work, "status", "--all", "--offline", "--root", root)
		if code != 0 {
			t.Fatalf("mixed status --all exit %d stderr=%q", code, stderr)
		}
		combined := stdout + stderr
		if !strings.Contains(combined, "partial: legacy-manifest") {
			t.Fatalf("missing partial: legacy-manifest: %q", combined)
		}
		if !strings.Contains(stdout, "legacy-control\t") {
			t.Fatalf("g0 not listed: %q", stdout)
		}
		if !strings.Contains(stdout, "legacy-one\t") {
			t.Fatalf("g1 not listed: %q", stdout)
		}
		if !strings.Contains(stdout, "legacy-extra\t") {
			t.Fatalf("g3 not listed: %q", stdout)
		}
		ideas := 0
		for _, line := range strings.Split(stdout, "\n") {
			if strings.Count(line, "\t") >= 4 && !strings.HasPrefix(line, "partial:") {
				ideas++
				if strings.Contains(line, "g2") && !strings.Contains(line, "legacy-") {
					t.Fatalf("g2 listed as idea row: %q", line)
				}
			}
		}
		if ideas != 3 {
			t.Fatalf("want 3 idea rows (G0/G1/G3), got %d\n%s", ideas, stdout)
		}
	})

	t.Run("MS-501-REL", func(t *testing.T) {
		root := repoRoot(t)
		if version.Version != "0.1.0-dev" {
			t.Fatalf("version.Version = %q, want 0.1.0-dev", version.Version)
		}

		installPath := filepath.Join(root, "docs", "install.md")
		installBody, err := os.ReadFile(installPath)
		if err != nil {
			t.Fatal(err)
		}
		installText := string(installBody)
		if !strings.Contains(installText, "releases/latest/download/mycelium-") {
			t.Fatal("docs/install.md missing releases/latest/download/mycelium-")
		}
		if !strings.Contains(installText, "~/.local/bin/mycelium") {
			t.Fatal("docs/install.md missing ~/.local/bin/mycelium")
		}

		dist := filepath.Join(root, "dist")
		t.Cleanup(func() { _ = os.RemoveAll(dist) })
		_ = os.RemoveAll(dist)

		cmd := exec.Command("bash", "scripts/release.sh", "0.1.0")
		cmd.Dir = root
		cmd.Env = append(os.Environ(), "CGO_ENABLED=0", "PATH="+stubDir+string(os.PathListSeparator)+os.Getenv("PATH"))
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("release.sh 0.1.0 failed: %v\n%s", err, out)
		}
		if _, err := os.Stat(ghMarker); err == nil {
			t.Fatal("gh was invoked during release.sh")
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
		if err := release.MatchSHA256SUMS(dist, entries); err != nil {
			t.Fatalf("checksum mismatch: %v", err)
		}
		_ = os.RemoveAll(dist)

		// Refuse: CHANGELOG with only Unreleased → no dist writes.
		tmp := t.TempDir()
		if err := os.MkdirAll(filepath.Join(tmp, "scripts"), 0o755); err != nil {
			t.Fatal(err)
		}
		script, err := os.ReadFile(filepath.Join(root, "scripts", "release.sh"))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(tmp, "scripts", "release.sh"), script, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(tmp, "CHANGELOG.md"), []byte("# Changelog\n\n## [Unreleased]\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		refuse := exec.Command("bash", "scripts/release.sh", "9.9.9")
		refuse.Dir = tmp
		refuse.Env = append(os.Environ(), "CGO_ENABLED=0", "PATH="+stubDir+string(os.PathListSeparator)+os.Getenv("PATH"))
		refuseOut, refuseErr := refuse.CombinedOutput()
		if refuseErr == nil {
			t.Fatalf("expected refuse for 9.9.9, got success:\n%s", refuseOut)
		}
		if entries, _ := os.ReadDir(filepath.Join(tmp, "dist")); len(entries) > 0 {
			t.Fatalf("refuse must write nothing under dist/; got %v", entries)
		}
		if version.Version != "0.1.0-dev" {
			t.Fatalf("in-tree version mutated to %q", version.Version)
		}
	})

	t.Run("MS-501-guards", func(t *testing.T) {
		work := t.TempDir()
		inst := filepath.Join(work, "guard")
		code, _, stderr := run(t, work, "new", "idea", "Guard Fixture", "--offline", "--dir", inst)
		if code != 0 {
			t.Fatalf("scaffold exit %d stderr=%q", code, stderr)
		}
		code, _, stderr = run(t, work, "state", "handed-off", "--dir", inst)
		if code != 1 {
			t.Fatalf("handed-off exit %d want 1 stderr=%q", code, stderr)
		}
		for _, verb := range []string{"council", "release", "install"} {
			code, _, stderr = run(t, work, verb)
			if code == 0 {
				t.Fatalf("%s must be unknown, got exit 0", verb)
			}
			if !strings.Contains(stderr, "unknown command") {
				t.Fatalf("%s stderr=%q", verb, stderr)
			}
		}
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
	})
}
