//go:build ignore

// Copies repo-root program/ into internal/embed/program/, filtering *.go.
// Run: go generate ./internal/embed
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "embed generate: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	module, err := findModuleRoot(wd)
	if err != nil {
		return err
	}
	src := filepath.Join(module, "program")
	dst := filepath.Join(module, "internal", "embed", "program")
	info, err := os.Stat(src)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("repo-root program/ missing at %s", src)
	}
	if err := os.RemoveAll(dst); err != nil {
		return fmt.Errorf("remove %s: %w", dst, err)
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", dst, err)
	}
	return filepath.WalkDir(src, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(dst, rel), 0o755)
		}
		if strings.HasSuffix(strings.ToLower(d.Name()), ".go") {
			return nil
		}
		return copyFile(path, filepath.Join(dst, rel))
	})
}

func findModuleRoot(start string) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %s", start)
		}
		dir = parent
	}
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}
