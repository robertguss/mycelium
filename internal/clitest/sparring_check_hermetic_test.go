package clitest_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/clitest"
)

func TestSlice2DisputedMissingCruxExits1(t *testing.T) {
	bin := clitest.Bin(t)
	workDir := t.TempDir()
	env := []string{"MYCELIUM_OFFLINE=1"}

	code, _, stderr := clitest.Run(t, bin, workDir, env, "new", "idea", "Slice2 Neg", "--offline")
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr)
	}
	inst := filepath.Join(workDir, "slice2-neg")
	qdir := filepath.Join(inst, "questions")
	if err := os.MkdirAll(qdir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := `+++
id = "OQ-001"
title = "Use SQLite"
agreement = "agree-to-disagree"
date = "2026-08-15"
+++

# OQ-001 — Use SQLite

## Question

q

## Context

c

## Positions

### Human

h

### Agent

a

## Reasons

### Human

rh

### Agent

ra

## Disposition

d
`
	if err := os.WriteFile(filepath.Join(qdir, "OQ-001-use-sqlite.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	code, _, stderr = clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 1 {
		t.Fatalf("check exit %d want 1 stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "## Crux") {
		t.Fatalf("stderr missing ## Crux: %q", stderr)
	}
	if !strings.Contains(stderr, "program/contracts/sparring.md") {
		t.Fatalf("stderr missing sparring.md: %q", stderr)
	}
	if !strings.Contains(stderr, "convention: sparring") {
		t.Fatalf("stderr missing convention: %q", stderr)
	}
}

func TestSlice2DisputedCompleteExits0(t *testing.T) {
	bin := clitest.Bin(t)
	workDir := t.TempDir()
	env := []string{"MYCELIUM_OFFLINE=1"}

	code, _, stderr := clitest.Run(t, bin, workDir, env, "new", "idea", "Slice2 Pos", "--offline")
	if code != 0 {
		t.Fatalf("scaffold exit %d stderr=%q", code, stderr)
	}
	inst := filepath.Join(workDir, "slice2-pos")
	qdir := filepath.Join(inst, "questions")
	if err := os.MkdirAll(qdir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := `+++
id = "OQ-001"
title = "Use SQLite"
agreement = "agree-to-disagree"
date = "2026-08-15"
+++

# OQ-001 — Use SQLite

## Question

q

## Context

c

## Positions

### Human

h

### Agent

a

## Reasons

### Human

rh

### Agent

ra

## Crux

### Human

ch

### Agent

ca

## Disposition

d
`
	if err := os.WriteFile(filepath.Join(qdir, "OQ-001-use-sqlite.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	code, _, stderr = clitest.Run(t, bin, workDir, env, "check", "--dir", inst)
	if code != 0 {
		t.Fatalf("check exit %d want 0 stderr=%q", code, stderr)
	}
}
