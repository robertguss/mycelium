// Package execrun runs external binaries (git, gh) behind an injectable Runner.
package execrun

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunOpts configures one Run invocation.
type RunOpts struct {
	Dir string
	Env []string // appended after the base environment when non-empty
}

// Result is the completed process output.
type Result struct {
	ExitCode int
	Stdout   []byte
	Stderr   []byte
}

// Runner looks up and runs binaries. Tests install fakes that record calls.
type Runner interface {
	LookPath(name string) (string, error)
	Run(ctx context.Context, name string, args []string, opts RunOpts) (Result, error)
}

// gitEnvKeys redirect git away from cmd.Dir when inherited.
var gitEnvKeys = map[string]struct{}{
	"GIT_DIR":                          {},
	"GIT_WORK_TREE":                    {},
	"GIT_INDEX_FILE":                   {},
	"GIT_OBJECT_DIRECTORY":             {},
	"GIT_COMMON_DIR":                   {},
	"GIT_ALTERNATE_OBJECT_DIRECTORIES": {},
}

// WithoutGitOverrides drops GIT_DIR / GIT_WORK_TREE and related location keys.
func WithoutGitOverrides(env []string) []string {
	out := make([]string, 0, len(env))
	for _, e := range env {
		key, _, ok := strings.Cut(e, "=")
		if !ok {
			out = append(out, e)
			continue
		}
		if _, drop := gitEnvKeys[key]; drop {
			continue
		}
		out = append(out, e)
	}
	return out
}

// Real wraps os/exec.
type Real struct{}

// LookPath implements Runner.
func (Real) LookPath(name string) (string, error) {
	return exec.LookPath(name)
}

// Run implements Runner.
// When opts.Dir is set, git location env overrides are stripped so the child
// uses Dir instead of an inherited GIT_DIR / GIT_WORK_TREE.
func (Real) Run(ctx context.Context, name string, args []string, opts RunOpts) (Result, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = opts.Dir
	switch {
	case opts.Dir != "":
		cmd.Env = WithoutGitOverrides(os.Environ())
		if len(opts.Env) > 0 {
			cmd.Env = append(cmd.Env, opts.Env...)
		}
	case len(opts.Env) > 0:
		cmd.Env = append(os.Environ(), opts.Env...)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	res := Result{
		Stdout: stdout.Bytes(),
		Stderr: stderr.Bytes(),
	}
	if err == nil {
		return res, nil
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		res.ExitCode = ee.ExitCode()
		return res, nil
	}
	return res, err
}

// Call is one recorded LookPath or Run.
type Call struct {
	Kind string // "lookpath" | "run"
	Name string
	Args []string
	Dir  string
}

// Recording wraps an inner Runner and records every call.
type Recording struct {
	Inner Runner
	Calls []Call
}

// LookPath implements Runner.
func (r *Recording) LookPath(name string) (string, error) {
	r.Calls = append(r.Calls, Call{Kind: "lookpath", Name: name})
	if r.Inner == nil {
		return "", errors.New("execrun: no inner runner")
	}
	return r.Inner.LookPath(name)
}

// Run implements Runner.
func (r *Recording) Run(ctx context.Context, name string, args []string, opts RunOpts) (Result, error) {
	cp := append([]string(nil), args...)
	r.Calls = append(r.Calls, Call{Kind: "run", Name: name, Args: cp, Dir: opts.Dir})
	if r.Inner == nil {
		return Result{}, errors.New("execrun: no inner runner")
	}
	return r.Inner.Run(ctx, name, args, opts)
}

// Called reports whether any recorded Run used binary name.
func (r *Recording) Called(name string) bool {
	for _, c := range r.Calls {
		if c.Kind == "run" && c.Name == name {
			return true
		}
	}
	return false
}

// Fake is a programmable Runner for unit tests (no real exec).
type Fake struct {
	Paths map[string]string
	// RunFunc, if set, handles Run. Otherwise returns ExitCode 0.
	RunFunc func(ctx context.Context, name string, args []string, opts RunOpts) (Result, error)
	Calls   []Call
}

// LookPath implements Runner.
func (f *Fake) LookPath(name string) (string, error) {
	f.Calls = append(f.Calls, Call{Kind: "lookpath", Name: name})
	if f.Paths != nil {
		if p, ok := f.Paths[name]; ok {
			return p, nil
		}
	}
	return "", fmt.Errorf("execrun: %q not found", name)
}

// Run implements Runner.
func (f *Fake) Run(ctx context.Context, name string, args []string, opts RunOpts) (Result, error) {
	cp := append([]string(nil), args...)
	f.Calls = append(f.Calls, Call{Kind: "run", Name: name, Args: cp, Dir: opts.Dir})
	if f.RunFunc != nil {
		return f.RunFunc(ctx, name, args, opts)
	}
	return Result{}, nil
}

// Called reports whether any recorded Run used binary name.
func (f *Fake) Called(name string) bool {
	for _, c := range f.Calls {
		if c.Kind == "run" && c.Name == name {
			return true
		}
	}
	return false
}
