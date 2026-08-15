package clitest_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/clitest"
	"github.com/robertguss/mycelium/internal/manifest"
)

func TestPackPresenceHermetic(t *testing.T) {
	bin := clitest.Bin(t)
	stubDir := installNetworkStubs(t)
	home := t.TempDir()
	workDir := t.TempDir()
	env := hermeticEnv(stubDir, home)

	code, _, stderr := clitest.Run(t, bin, workDir, env, "new", "idea", "Pack On", "--offline")
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr)
	}
	inst := filepath.Join(workDir, "pack-on")
	if _, err := os.Stat(filepath.Join(inst, "program", "packs", "council")); err != nil {
		t.Fatalf("council pack missing on scaffold: %v", err)
	}

	code, _, stderr = clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("pack present no reviews: exit %d stderr=%q", code, stderr)
	}
	assertNoGh(t, stderr)
	assertNoHomeTouch(t, home)
}

func TestPackAbsentNoReviewsHermetic(t *testing.T) {
	bin := clitest.Bin(t)
	stubDir := installNetworkStubs(t)
	home := t.TempDir()
	workDir := t.TempDir()
	env := hermeticEnv(stubDir, home)

	code, _, stderr := clitest.Run(t, bin, workDir, env, "new", "idea", "Pack Off", "--offline")
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr)
	}
	inst := filepath.Join(workDir, "pack-off")
	if err := os.RemoveAll(filepath.Join(inst, "program", "packs", "council")); err != nil {
		t.Fatal(err)
	}

	code, _, stderr = clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("pack absent no reviews: exit %d stderr=%q", code, stderr)
	}
	assertNoGh(t, stderr)
	assertNoHomeTouch(t, home)
}

func TestPackAbsentReviewsLeftoverHermetic(t *testing.T) {
	bin := clitest.Bin(t)
	stubDir := installNetworkStubs(t)
	home := t.TempDir()
	workDir := t.TempDir()
	env := hermeticEnv(stubDir, home)

	code, _, stderr := clitest.Run(t, bin, workDir, env, "new", "idea", "Pack Leftover", "--offline")
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr)
	}
	inst := filepath.Join(workDir, "pack-leftover")
	if err := os.RemoveAll(filepath.Join(inst, "program", "packs", "council")); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(inst, "reviews"), 0o755); err != nil {
		t.Fatal(err)
	}

	code, _, stderr = clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 1 {
		t.Fatalf("want exit 1, got %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "reviews/") {
		t.Fatalf("stderr must name reviews/: %q", stderr)
	}
	if !strings.Contains(stderr, "extra-top-level") {
		t.Fatalf("stderr must name extra-top-level: %q", stderr)
	}
	if !strings.Contains(stderr, "council pack absent") {
		t.Fatalf("stderr must say council pack absent: %q", stderr)
	}
	assertNoGh(t, stderr)
	assertNoHomeTouch(t, home)
}

func TestPackAbsentReviewsDeviationHermetic(t *testing.T) {
	bin := clitest.Bin(t)
	stubDir := installNetworkStubs(t)
	home := t.TempDir()
	workDir := t.TempDir()
	env := hermeticEnv(stubDir, home)

	code, _, stderr := clitest.Run(t, bin, workDir, env, "new", "idea", "Pack Dev", "--offline")
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr)
	}
	inst := filepath.Join(workDir, "pack-dev")
	if err := os.RemoveAll(filepath.Join(inst, "program", "packs", "council")); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(inst, "reviews"), 0o755); err != nil {
		t.Fatal(err)
	}
	addDeviation(t, inst, "extra-top-level:reviews/", "leftover reviews after disabling council pack")

	code, _, stderr = clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("deviation should pass: exit %d stderr=%q", code, stderr)
	}
	assertNoGh(t, stderr)
	assertNoHomeTouch(t, home)
}

func TestPackCollisionHermetic(t *testing.T) {
	bin := clitest.Bin(t)
	stubDir := installNetworkStubs(t)
	home := t.TempDir()
	workDir := t.TempDir()
	env := hermeticEnv(stubDir, home)

	code, _, stderr := clitest.Run(t, bin, workDir, env, "new", "idea", "Pack Collide", "--offline")
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr)
	}
	inst := filepath.Join(workDir, "pack-collide")
	dir := filepath.Join(inst, "program", "packs", "fixture-pack", "templates")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	schema := `namespace = "CMP"
home = "fixture-home"
filename_pattern = "CMP-{NNN}-{slug}.md"
stage_scoped = false
digits = 3
required_front_matter = ["id"]
required_sections = ["Body"]
`
	if err := os.WriteFile(filepath.Join(dir, "collide.schema.toml"), []byte(schema), 0o644); err != nil {
		t.Fatal(err)
	}

	code, _, stderr = clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 1 {
		t.Fatalf("want collision exit 1, got %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "pack-collision") && !strings.Contains(stderr, "collision") {
		t.Fatalf("stderr must name collision: %q", stderr)
	}
	if !strings.Contains(stderr, "CMP") {
		t.Fatalf("stderr must name CMP: %q", stderr)
	}
	if !strings.Contains(stderr, "council") || !strings.Contains(stderr, "fixture-pack") {
		t.Fatalf("stderr must name both packs: %q", stderr)
	}
	assertNoGh(t, stderr)
	assertNoHomeTouch(t, home)
}

func hermeticEnv(stubDir, home string) []string {
	return []string{
		"MYCELIUM_OFFLINE=1",
		"HOME=" + home,
		"PATH=" + stubDir + string(os.PathListSeparator) + os.Getenv("PATH"),
	}
}

func addDeviation(t *testing.T, inst, convention, reason string) {
	t.Helper()
	path := filepath.Join(inst, "mycelium.toml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	m, err := manifest.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	m.Deviations = append(m.Deviations, manifest.Deviation{
		Convention: convention,
		Reason:     reason,
	})
	out, err := manifest.Encode(m)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, out, 0o644); err != nil {
		t.Fatal(err)
	}
}

func assertNoGh(t *testing.T, stderr string) {
	t.Helper()
	if strings.Contains(stderr, "gh ") || strings.Contains(strings.ToLower(stderr), "github.com") {
		t.Fatalf("gh must not be invoked: stderr=%q", stderr)
	}
}

func assertNoHomeTouch(t *testing.T, home string) {
	t.Helper()
	entries, err := os.ReadDir(home)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("tests must not write home dir; found %v", entries)
	}
}
