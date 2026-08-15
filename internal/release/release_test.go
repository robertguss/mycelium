package release_test

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/release"
)

func TestHasVersionHeading(t *testing.T) {
	cases := []struct {
		name      string
		changelog string
		version   string
		want      bool
	}{
		{"plain", "## [0.1.0]\n", "0.1.0", true},
		{"dated", "## [0.1.0] - 2026-08-15\n", "0.1.0", true},
		{"unreleased only", "## [Unreleased]\n", "0.1.0", false},
		{"wrong version", "## [0.2.0]\n", "0.1.0", false},
		{"prefix trap", "## [0.1.0-rc1]\n", "0.1.0", false},
		{"empty version", "## [0.1.0]\n", "", false},
		{"fixture body", "# Changelog\n\n## [Unreleased]\n\n## [0.1.0] - 2026-08-15\n\n### Added\n\n- x\n", "0.1.0", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := release.HasVersionHeading(tc.changelog, tc.version)
			if got != tc.want {
				t.Fatalf("HasVersionHeading(%q) = %v, want %v", tc.version, got, tc.want)
			}
		})
	}
}

func TestMissingHeadingMessage(t *testing.T) {
	msg := release.MissingHeadingMessage("9.9.9")
	if !strings.Contains(msg, "## [9.9.9]") {
		t.Fatalf("message must name missing heading: %q", msg)
	}
	if !strings.Contains(msg, "release refused") {
		t.Fatalf("message must refuse: %q", msg)
	}
}

func TestParseAndMatchSHA256SUMS(t *testing.T) {
	dir := t.TempDir()
	linuxBody := []byte("linux-bytes")
	darwinBody := []byte("darwin-bytes")
	linuxHex := sha256Hex(linuxBody)
	darwinHex := sha256Hex(darwinBody)
	if err := os.WriteFile(filepath.Join(dir, "mycelium-linux-amd64"), linuxBody, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "mycelium-darwin-arm64"), darwinBody, 0o644); err != nil {
		t.Fatal(err)
	}

	sums := linuxHex + "  mycelium-linux-amd64\n" + darwinHex + "  mycelium-darwin-arm64\n"
	entries, err := release.ParseSHA256SUMS(sums)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("entries=%d", len(entries))
	}
	if err := release.MatchSHA256SUMS(dir, entries); err != nil {
		t.Fatalf("match: %v", err)
	}

	// Tamper one file → mismatch.
	if err := os.WriteFile(filepath.Join(dir, "mycelium-linux-amd64"), []byte("tampered"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := release.MatchSHA256SUMS(dir, entries); err == nil {
		t.Fatal("expected digest mismatch after tamper")
	}
}

func TestParseSHA256SUMSErrors(t *testing.T) {
	cases := []string{
		"",
		"not-a-sum",
		"abcd  name", // short hex + single space
		strings.Repeat("a", 64) + "  path/with/slash",
		strings.Repeat("z", 64) + "  name", // invalid hex
	}
	for _, raw := range cases {
		if _, err := release.ParseSHA256SUMS(raw); err == nil {
			t.Fatalf("expected error for %q", raw)
		}
	}
}

func TestMatchSHA256SUMSEmpty(t *testing.T) {
	if err := release.MatchSHA256SUMS(t.TempDir(), nil); err == nil {
		t.Fatal("expected error for empty entries")
	}
}

func TestFileSHA256Missing(t *testing.T) {
	if _, err := release.FileSHA256(filepath.Join(t.TempDir(), "nope")); err == nil {
		t.Fatal("expected missing file error")
	}
}

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
