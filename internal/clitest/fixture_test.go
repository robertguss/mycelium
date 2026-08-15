package clitest_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/clitest"
	"github.com/robertguss/mycelium/internal/execrun"
	"github.com/robertguss/mycelium/internal/idpath"
)

func TestHermeticFixtureBinary(t *testing.T) {
	bin := clitest.Bin(t)
	stubDir := installNetworkStubs(t)
	workDir := t.TempDir()
	env := []string{
		"MYCELIUM_OFFLINE=1",
		"PATH=" + stubDir + string(os.PathListSeparator) + os.Getenv("PATH"),
	}

	code, _, stderr := clitest.Run(t, bin, workDir, env, "new", "idea", "CI Fixture", "--offline")
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr)
	}
	inst := filepath.Join(workDir, "ci-fixture")
	clitest.WriteRanges(t, inst)

	for _, typ := range idpath.Types() {
		code, _, stderr = clitest.Run(
			t,
			bin,
			workDir,
			env,
			"new",
			typ.Key,
			"Sample "+typ.Key,
			"--dir",
			inst,
		)
		if code != 0 {
			t.Fatalf("new %s exit %d stderr=%q", typ.Key, code, stderr)
		}
	}

	code, stdout, stderr := clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("check exit %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stdout, "artifacts: 11") {
		t.Fatalf("check stdout missing artifact count: %q", stdout)
	}
	clitest.AssertNoForbiddenPaths(t, inst)
}

func TestNetworkViolations(t *testing.T) {
	rec := &execrun.Recording{
		Calls: []execrun.Call{
			{Kind: "lookpath", Name: "gh"},
			{Kind: "run", Name: "git", Args: []string{"status"}},
			{Kind: "run", Name: "gh", Args: []string{"repo", "create"}},
			{Kind: "run", Name: "/usr/bin/git", Args: []string{"remote", "-v"}},
		},
	}
	got := clitest.NetworkViolations(rec)
	want := []string{"gh repo create", "git remote -v"}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("violations=%q want %q", got, want)
	}
}

func installNetworkStubs(t *testing.T) string {
	t.Helper()
	realGit, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	writeStub(t, filepath.Join(dir, "gh"), `#!/bin/sh
echo "gh network access forbidden" >&2
exit 99
`)
	writeStub(t, filepath.Join(dir, "git"), fmt.Sprintf(`#!/bin/sh
if [ "$1" = "remote" ]; then
  echo "git remote network access forbidden" >&2
  exit 99
fi
exec %s "$@"
`, shellQuote(realGit)))
	return dir
}

func writeStub(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", `'"'"'`) + "'"
}
