// Package release provides CHANGELOG heading and SHA256SUMS helpers for
// scripts/release.sh hermetic tests. Release itself is the shell script, not
// a mycelium CLI verb.
package release

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// HasVersionHeading reports whether changelog text contains a Keep-a-Changelog
// heading for version. Accepts "## [0.1.0]" and "## [0.1.0] - YYYY-MM-DD".
func HasVersionHeading(changelog, version string) bool {
	if version == "" {
		return false
	}
	plain := "## [" + version + "]"
	datedPrefix := plain + " - "
	sc := bufio.NewScanner(strings.NewReader(changelog))
	for sc.Scan() {
		line := strings.TrimRight(sc.Text(), "\r")
		if line == plain || strings.HasPrefix(line, datedPrefix) {
			return true
		}
	}
	return false
}

// MissingHeadingMessage is the teaching line scripts/release.sh should print
// (or a Go-side equivalent) when the CHANGELOG lacks the version heading.
func MissingHeadingMessage(version string) string {
	return fmt.Sprintf(
		"release refused: CHANGELOG.md missing heading ## [%s] (optional ' - YYYY-MM-DD' suffix allowed)",
		version,
	)
}

// SumsEntry is one SHA256SUMS line: hex digest and basename.
type SumsEntry struct {
	Hex  string
	Name string
}

// ParseSHA256SUMS parses classic `sha256sum` output (`<hex>  <name>`).
func ParseSHA256SUMS(text string) ([]SumsEntry, error) {
	var out []SumsEntry
	sc := bufio.NewScanner(strings.NewReader(text))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		hexPart, name, ok := strings.Cut(line, "  ")
		if !ok {
			return nil, fmt.Errorf("SHA256SUMS: bad line %q (want \"<hex>  <name>\")", line)
		}
		hexPart = strings.TrimSpace(hexPart)
		name = strings.TrimSpace(name)
		if len(hexPart) != 64 {
			return nil, fmt.Errorf("SHA256SUMS: hex length %d, want 64", len(hexPart))
		}
		if _, err := hex.DecodeString(hexPart); err != nil {
			return nil, fmt.Errorf("SHA256SUMS: invalid hex: %w", err)
		}
		if name == "" || strings.Contains(name, "/") {
			return nil, fmt.Errorf("SHA256SUMS: name must be a basename, got %q", name)
		}
		out = append(out, SumsEntry{Hex: strings.ToLower(hexPart), Name: name})
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("SHA256SUMS: empty")
	}
	return out, nil
}

// FileSHA256 returns the lowercase hex SHA-256 of path.
func FileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// MatchSHA256SUMS verifies each entry against files in dir (basename only).
// Returns nil when every listed file's digest matches.
func MatchSHA256SUMS(dir string, entries []SumsEntry) error {
	if len(entries) == 0 {
		return fmt.Errorf("SHA256SUMS: no entries")
	}
	for _, e := range entries {
		path := filepath.Join(dir, e.Name)
		got, err := FileSHA256(path)
		if err != nil {
			return fmt.Errorf("%s: %w", e.Name, err)
		}
		if got != e.Hex {
			return fmt.Errorf("%s: digest mismatch want %s got %s", e.Name, e.Hex, got)
		}
	}
	return nil
}
