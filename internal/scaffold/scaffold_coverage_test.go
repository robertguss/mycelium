package scaffold_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/robertguss/mycelium/internal/clock"
	"github.com/robertguss/mycelium/internal/execrun"
	"github.com/robertguss/mycelium/internal/scaffold"
)

func TestOfflineFromEnv(t *testing.T) {
	cases := map[string]bool{
		"1": true, "true": true, "YES": true, "on": true, " True ": true,
		"": false, "0": false, "false": false, "no": false, "maybe": false,
	}
	for in, want := range cases {
		got := scaffold.OfflineFromEnv(func(string) string { return in })
		if got != want {
			t.Fatalf("OfflineFromEnv(%q)=%v want %v", in, got, want)
		}
	}
	if scaffold.OfflineFromEnv(nil) {
		t.Fatal("nil getenv should read real env; expected false without MYCELIUM_OFFLINE")
	}
}

func TestScaffoldValidationErrors(t *testing.T) {
	root := t.TempDir()
	clk := clock.Fixed{T: time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)}

	t.Run("empty name", func(t *testing.T) {
		var stderr bytes.Buffer
		code := scaffold.Run(scaffold.Options{Name: "  ", Offline: true, Cwd: root}, scaffold.Deps{
			Clock: clk, Stderr: &stderr, Stdout: io.Discard,
		})
		if code != 1 || !strings.Contains(stderr.String(), "idea name is required") {
			t.Fatalf("code=%d stderr=%q", code, stderr.String())
		}
	})

	t.Run("unknown tier", func(t *testing.T) {
		var stderr bytes.Buffer
		code := scaffold.Run(scaffold.Options{Name: "X", Tier: "nope", Offline: true, Cwd: root}, scaffold.Deps{
			Clock: clk, Stderr: &stderr, Stdout: io.Discard,
		})
		if code != 1 || !strings.Contains(stderr.String(), "unknown tier") {
			t.Fatalf("code=%d stderr=%q", code, stderr.String())
		}
	})

	t.Run("bad slug", func(t *testing.T) {
		var stderr bytes.Buffer
		code := scaffold.Run(scaffold.Options{Name: "!!!", Offline: true, Cwd: root}, scaffold.Deps{
			Clock: clk, Stderr: &stderr, Stdout: io.Discard,
		})
		if code != 1 || !strings.Contains(stderr.String(), "cannot slugify") {
			t.Fatalf("code=%d stderr=%q", code, stderr.String())
		}
	})

	t.Run("missing parent", func(t *testing.T) {
		var stderr bytes.Buffer
		code := scaffold.Run(scaffold.Options{
			Name: "Orphan", Dir: filepath.Join(root, "missing-parent", "child"), Offline: true, Cwd: root,
		}, scaffold.Deps{Clock: clk, Stderr: &stderr, Stdout: io.Discard})
		if code != 1 || !strings.Contains(stderr.String(), "parent directory does not exist") {
			t.Fatalf("code=%d stderr=%q", code, stderr.String())
		}
	})

	t.Run("target is file", func(t *testing.T) {
		path := filepath.Join(root, "as-file")
		if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		var stderr bytes.Buffer
		code := scaffold.Run(scaffold.Options{
			Name: "File Target", Dir: path, Offline: true, Cwd: root,
		}, scaffold.Deps{Clock: clk, Stderr: &stderr, Stdout: io.Discard})
		if code != 1 || !strings.Contains(stderr.String(), "already exists") {
			t.Fatalf("code=%d stderr=%q", code, stderr.String())
		}
	})

	t.Run("nonempty dir without journal", func(t *testing.T) {
		path := filepath.Join(root, "busy")
		if err := os.Mkdir(path, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(path, "README.md"), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		var stderr bytes.Buffer
		code := scaffold.Run(scaffold.Options{
			Name: "Busy", Dir: path, Offline: true, Cwd: root,
		}, scaffold.Deps{Clock: clk, Stderr: &stderr, Stdout: io.Discard})
		if code != 1 || !strings.Contains(stderr.String(), "already exists") {
			t.Fatalf("code=%d stderr=%q", code, stderr.String())
		}
	})
}

func TestScaffoldGitInitFailure(t *testing.T) {
	root := t.TempDir()
	fake := &execrun.Fake{
		RunFunc: func(ctx context.Context, name string, args []string, opts execrun.RunOpts) (execrun.Result, error) {
			return execrun.Result{ExitCode: 7, Stderr: []byte("init boom")}, nil
		},
	}
	var stderr bytes.Buffer
	code := scaffold.Run(scaffold.Options{Name: "Git Fail", Offline: true, Cwd: root}, scaffold.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)},
		Runner: fake,
		Stdout: io.Discard,
		Stderr: &stderr,
	})
	if code != 1 {
		t.Fatalf("exit %d", code)
	}
	if !strings.Contains(stderr.String(), "git init failed") {
		t.Fatalf("stderr=%q", stderr.String())
	}
}

func TestScaffoldOptionalPublishPath(t *testing.T) {
	root := t.TempDir()
	fake := &execrun.Fake{
		Paths: map[string]string{"git": "/usr/bin/git", "gh": "/usr/bin/gh"},
		RunFunc: func(ctx context.Context, name string, args []string, opts execrun.RunOpts) (execrun.Result, error) {
			if name == "git" && containsArg(args, "init") {
				if err := os.MkdirAll(filepath.Join(opts.Dir, ".git"), 0o755); err != nil {
					return execrun.Result{}, err
				}
				return execrun.Result{}, nil
			}
			if name == "gh" {
				return execrun.Result{ExitCode: 1, Stderr: []byte("not logged in")}, nil
			}
			return execrun.Result{ExitCode: 1, Stderr: []byte("unexpected " + name)}, nil
		},
	}
	var stdout, stderr bytes.Buffer
	code := scaffold.Run(scaffold.Options{
		Name: "Publish Optional", Offline: false, Publish: false, Cwd: root,
		Dir: filepath.Join(root, "pub-opt"),
	}, scaffold.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)},
		Runner: fake,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "publish: mycelium publish") {
		t.Fatalf("stdout=%q", stdout.String())
	}
}

func TestScaffoldRequiredPublishReports(t *testing.T) {
	root := t.TempDir()
	origin := ""
	fake := &execrun.Fake{
		Paths: map[string]string{"git": "/usr/bin/git", "gh": "/usr/bin/gh"},
		RunFunc: func(ctx context.Context, name string, args []string, opts execrun.RunOpts) (execrun.Result, error) {
			if name == "git" {
				if containsArg(args, "init") {
					if err := os.MkdirAll(filepath.Join(opts.Dir, ".git"), 0o755); err != nil {
						return execrun.Result{}, err
					}
					return execrun.Result{}, nil
				}
				if len(args) >= 2 && args[0] == "remote" && args[1] == "get-url" {
					if origin == "" {
						return execrun.Result{ExitCode: 2, Stderr: []byte("no such remote")}, nil
					}
					return execrun.Result{Stdout: []byte(origin + "\n")}, nil
				}
				if len(args) >= 4 && args[0] == "remote" && args[1] == "add" {
					origin = args[3]
					return execrun.Result{}, nil
				}
				return execrun.Result{}, nil
			}
			if name == "gh" {
				joined := strings.Join(args, " ")
				switch {
				case strings.Contains(joined, "auth status"):
					return execrun.Result{}, nil
				case strings.Contains(joined, "api user"):
					return execrun.Result{Stdout: []byte("acme\n")}, nil
				case strings.Contains(joined, "repo create"):
					return execrun.Result{Stdout: []byte("https://github.com/acme/pub-req\n")}, nil
				default:
					return execrun.Result{}, nil
				}
			}
			return execrun.Result{ExitCode: 1, Stderr: []byte("unexpected")}, nil
		},
	}
	var stdout, stderr bytes.Buffer
	code := scaffold.Run(scaffold.Options{
		Name: "Publish Required", Offline: false, Publish: true, Tier: "standard",
		Cwd: root, Dir: filepath.Join(root, "pub-req"),
	}, scaffold.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)},
		Runner: fake,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "published:") {
		t.Fatalf("stdout missing published: %q", stdout.String())
	}
}

func TestScaffoldHighAssuranceTierAndDefaults(t *testing.T) {
	root := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	var stdout, stderr bytes.Buffer
	code := scaffold.Run(scaffold.Options{
		Name: "HA", Offline: true, Tier: "high-assurance", Cwd: root,
	}, scaffold.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)},
		Runner: rec,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	inst := filepath.Join(root, "ha")
	for _, d := range []string{"spikes", "findings", "recommendations", "requirements", "phases", "milestones"} {
		mustExist(t, inst, d+"/README.md")
	}
}

func TestScaffoldNilDepsUseDefaults(t *testing.T) {
	root := t.TempDir()
	fake := &execrun.Fake{
		RunFunc: func(ctx context.Context, name string, args []string, opts execrun.RunOpts) (execrun.Result, error) {
			if name == "git" && containsArg(args, "init") {
				if err := os.MkdirAll(filepath.Join(opts.Dir, ".git"), 0o755); err != nil {
					return execrun.Result{}, err
				}
				return execrun.Result{}, nil
			}
			return execrun.Result{}, nil
		},
	}
	code := scaffold.Run(scaffold.Options{
		Name: "Nil Deps", Offline: true, Cwd: root, Dir: filepath.Join(root, "nil-deps"),
	}, scaffold.Deps{Runner: fake})
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
}

func TestScaffoldRelativeDirUnderCwd(t *testing.T) {
	root := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	var stdout, stderr bytes.Buffer
	code := scaffold.Run(scaffold.Options{
		Name: "Rel Dir", Offline: true, Dir: "rel-child", Cwd: root,
	}, scaffold.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)},
		Runner: rec,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	mustExist(t, filepath.Join(root, "rel-child"), "mycelium.toml")
}

func TestScaffoldEmptyCwdUsesGetwd(t *testing.T) {
	root := t.TempDir()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(prev) })

	rec := &execrun.Recording{Inner: execrun.Real{}}
	var stdout, stderr bytes.Buffer
	code := scaffold.Run(scaffold.Options{
		Name: "Cwd Empty", Offline: true, // Cwd intentionally empty
	}, scaffold.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)},
		Runner: rec,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	mustExist(t, filepath.Join(root, "cwd-empty"), "mycelium.toml")
}

func TestScaffoldRequiredPublishFailure(t *testing.T) {
	root := t.TempDir()
	fake := &execrun.Fake{
		Paths: map[string]string{"git": "/usr/bin/git", "gh": "/usr/bin/gh"},
		RunFunc: func(ctx context.Context, name string, args []string, opts execrun.RunOpts) (execrun.Result, error) {
			if name == "git" && containsArg(args, "init") {
				if err := os.MkdirAll(filepath.Join(opts.Dir, ".git"), 0o755); err != nil {
					return execrun.Result{}, err
				}
				return execrun.Result{}, nil
			}
			if name == "gh" {
				return execrun.Result{ExitCode: 1, Stderr: []byte("not logged in")}, nil
			}
			return execrun.Result{ExitCode: 1, Stderr: []byte("unexpected")}, nil
		},
	}
	var stdout, stderr bytes.Buffer
	code := scaffold.Run(scaffold.Options{
		Name: "Pub Fail", Offline: false, Publish: true,
		Cwd: root, Dir: filepath.Join(root, "pub-fail"),
	}, scaffold.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)},
		Runner: fake,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if code != 1 {
		t.Fatalf("exit %d want 1 stderr=%q", code, stderr.String())
	}
}

func TestScaffoldEmitsPackSkillsWhenPackPresent(t *testing.T) {
	root := t.TempDir()
	rec := &execrun.Recording{Inner: execrun.Real{}}
	var stdout, stderr bytes.Buffer
	code := scaffold.Run(scaffold.Options{
		Name: "Pack Skills", Offline: true, Cwd: root,
	}, scaffold.Deps{
		Clock:  clock.Fixed{T: time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)},
		Runner: rec,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if code != 0 {
		t.Fatalf("exit %d stderr=%q", code, stderr.String())
	}
	inst := filepath.Join(root, "pack-skills")
	mustExist(t, inst, "program/packs/council/README.md")
	mustExist(t, inst, ".agents/skills/council/SKILL.md")
	mustExist(t, inst, ".agents/skills/second-opinion/SKILL.md")
	agents, err := os.ReadFile(filepath.Join(inst, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(agents), "sparring still applies") {
		t.Fatalf("capability note missing")
	}
}
