package clitest

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/robertguss/mycelium/internal/execrun"
	"github.com/robertguss/mycelium/internal/manifest"
)

var (
	binOnce sync.Once
	binPath string
	binErr  error
)

// ForbiddenPaths are master-only paths that instances must exclude.
var ForbiddenPaths = []string{
	"framework",
	"Justfile",
	"scripts",
	"research-program.toml",
}

// Bin returns a mycelium binary built once for the test process.
func Bin(t *testing.T) string {
	t.Helper()
	if path := os.Getenv("MYCELIUM_BIN"); path != "" {
		return path
	}
	binOnce.Do(func() {
		root, err := moduleRoot()
		if err != nil {
			binErr = err
			return
		}
		dir, err := os.MkdirTemp("", "mycelium-clitest-")
		if err != nil {
			binErr = err
			return
		}
		binPath = filepath.Join(dir, "mycelium")
		cmd := exec.Command("go", "build", "-o", binPath, "./cmd/mycelium")
		cmd.Dir = root
		cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
		out, err := cmd.CombinedOutput()
		if err != nil {
			binErr = fmt.Errorf("build mycelium: %w: %s", err, bytes.TrimSpace(out))
		}
	})
	if binErr != nil {
		t.Fatal(binErr)
	}
	return binPath
}

// Run executes a mycelium binary and captures its result.
func Run(t *testing.T, bin string, dir string, env []string, args ...string) (code int, stdout, stderr string) {
	t.Helper()
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), env...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err == nil {
		return 0, outBuf.String(), errBuf.String()
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode(), outBuf.String(), errBuf.String()
	}
	t.Fatalf("run mycelium: %v", err)
	return 0, "", ""
}

// NetworkViolations describes recorded gh and git remote runs.
func NetworkViolations(rec *execrun.Recording) []string {
	if rec == nil {
		return nil
	}
	var violations []string
	for _, call := range rec.Calls {
		if call.Kind != "run" {
			continue
		}
		name := filepath.Base(call.Name)
		if name == "gh" || (name == "git" && len(call.Args) > 0 && call.Args[0] == "remote") {
			violations = append(violations, strings.Join(append([]string{name}, call.Args...), " "))
		}
	}
	return violations
}

// AssertNoNetwork fails when a recording contains a network-capable run.
func AssertNoNetwork(t *testing.T, rec *execrun.Recording) {
	t.Helper()
	if violations := NetworkViolations(rec); len(violations) > 0 {
		t.Fatalf("network commands executed: %s", strings.Join(violations, "; "))
	}
}

// WriteRanges adds stage-scoped identifier ranges to an instance.
func WriteRanges(t *testing.T, inst string) {
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
	m.Identifiers = map[string]manifest.Range{
		"findings":        mustRange(t, "FND-001..FND-099", "FND"),
		"recommendations": mustRange(t, "REC-001..REC-099", "REC"),
		"requirements":    mustRange(t, "REQ-001..REQ-299", "REQ"),
	}
	out, err := manifest.Encode(m)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, out, 0o644); err != nil {
		t.Fatal(err)
	}
}

// AssertNoForbiddenPaths checks that an instance excludes master-only paths.
func AssertNoForbiddenPaths(t *testing.T, inst string) {
	t.Helper()
	for _, rel := range ForbiddenPaths {
		if _, err := os.Stat(filepath.Join(inst, rel)); !os.IsNotExist(err) {
			t.Fatalf("forbidden path present: %s (err=%v)", rel, err)
		}
	}
}

func moduleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("go.mod not found")
		}
		dir = parent
	}
}

func mustRange(t *testing.T, raw, ns string) manifest.Range {
	t.Helper()
	rg, err := manifest.ParseRange(raw, ns)
	if err != nil {
		t.Fatal(err)
	}
	return rg
}
