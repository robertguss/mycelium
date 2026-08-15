package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/cli"
	"github.com/robertguss/mycelium/internal/version"
)

func TestVersionStdoutEquality(t *testing.T) {
	cases := []struct {
		name string
		argv []string
		want string
	}{
		{"bare version", []string{"mycelium", "version"}, version.Version + "\n"},
		{"default stamp", []string{"mycelium", "version"}, "0.1.0-dev\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := cli.Run(tc.argv, &stdout, &stderr)
			if code != 0 {
				t.Fatalf("exit %d; stderr=%q", code, stderr.String())
			}
			if stdout.String() != tc.want {
				t.Fatalf("stdout=%q want %q", stdout.String(), tc.want)
			}
			if stderr.Len() != 0 {
				t.Fatalf("stderr not empty: %q", stderr.String())
			}
		})
	}
}

func TestHelpExitsZero(t *testing.T) {
	cases := [][]string{
		{"mycelium"},
		{"mycelium", "-h"},
		{"mycelium", "--help"},
		{"mycelium", "help"},
	}
	for _, argv := range cases {
		t.Run(strings.Join(argv, " "), func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := cli.Run(argv, &stdout, &stderr)
			if code != 0 {
				t.Fatalf("exit %d; stderr=%q", code, stderr.String())
			}
			if !strings.Contains(stdout.String(), "mycelium") {
				t.Fatalf("usage missing from stdout: %q", stdout.String())
			}
		})
	}
}

func TestUnknownCommandTeachingError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"mycelium", "not-a-command"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("exit %d want 1", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout should be empty, got %q", stdout.String())
	}
	errText := stderr.String()
	for _, prefix := range []string{"mycelium:", "convention:", "contract:", "fix:"} {
		if !strings.Contains(errText, prefix) {
			t.Fatalf("teaching error missing %q in %q", prefix, errText)
		}
	}
	lines := strings.Split(strings.TrimSuffix(errText, "\n"), "\n")
	if len(lines) != 4 {
		t.Fatalf("want 4 teaching lines, got %d: %q", len(lines), errText)
	}
}

func TestPhase01NotImplementedTeachingError(t *testing.T) {
	for _, cmd := range []string{"new", "check", "tier", "publish"} {
		t.Run(cmd, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := cli.Run([]string{"mycelium", cmd}, &stdout, &stderr)
			if code != 1 {
				t.Fatalf("exit %d want 1", code)
			}
			errText := stderr.String()
			if !strings.Contains(errText, "not implemented in this slice") {
				t.Fatalf("stderr=%q", errText)
			}
			for _, prefix := range []string{"mycelium:", "convention:", "contract:", "fix:"} {
				if !strings.Contains(errText, prefix) {
					t.Fatalf("missing %q in %q", prefix, errText)
				}
			}
		})
	}
}
