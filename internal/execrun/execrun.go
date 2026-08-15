// Package execrun runs external binaries (git, gh) behind an injectable Runner.
package execrun

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// RunOpts configures one Run invocation.
type RunOpts struct {
	Dir string
	Env []string // appended to the process environment when non-empty
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

// Real wraps os/exec.
type Real struct{}

// LookPath implements Runner.
func (Real) LookPath(name string) (string, error) {
	return exec.LookPath(name)
}

// Run implements Runner.
func (Real) Run(ctx context.Context, name string, args []string, opts RunOpts) (Result, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = opts.Dir
	if len(opts.Env) > 0 {
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
