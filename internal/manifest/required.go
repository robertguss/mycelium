package manifest

import (
	"os"
	"path/filepath"
	"strings"
)

// DefaultRequiredKeys is the hardcoded required top-level set when an
// instance has no program/contracts/manifest.md (includes github_repo).
func DefaultRequiredKeys() []string {
	out := make([]string, len(requiredTop))
	copy(out, requiredTop)
	return out
}

// RequiredKeysFromContractMarkdown reads the ## Required fields pipe-table
// first column. Keys may be bare or backtick-wrapped.
func RequiredKeysFromContractMarkdown(md []byte) ([]string, error) {
	lines := strings.Split(string(md), "\n")
	inSection := false
	var keys []string
	seen := map[string]struct{}{}
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "## ") {
			heading := strings.TrimSpace(strings.TrimPrefix(trim, "## "))
			if inSection {
				break
			}
			inSection = strings.EqualFold(heading, "Required fields")
			continue
		}
		if !inSection {
			continue
		}
		if !strings.HasPrefix(trim, "|") {
			continue
		}
		cells := splitPipeRow(trim)
		if len(cells) == 0 {
			continue
		}
		field := strings.TrimSpace(cells[0])
		lower := strings.ToLower(field)
		if lower == "field" || isMDSeparator(field) {
			continue
		}
		key := strings.Trim(field, "`")
		key = strings.TrimSpace(key)
		if key == "" || !isIdentKey(key) {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		keys = append(keys, key)
	}
	if len(keys) == 0 {
		return nil, ErrInvalid
	}
	return keys, nil
}

// RequiredKeysForInstance loads root/program/contracts/manifest.md.
// Missing or unparseable → DefaultRequiredKeys (includes github_repo).
// Does not read the binary embed when the instance file exists.
func RequiredKeysForInstance(root string) []string {
	path := filepath.Join(root, "program", "contracts", "manifest.md")
	b, err := os.ReadFile(path)
	if err != nil {
		return DefaultRequiredKeys()
	}
	keys, err := RequiredKeysFromContractMarkdown(b)
	if err != nil || len(keys) == 0 {
		return DefaultRequiredKeys()
	}
	return keys
}

func splitPipeRow(line string) []string {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "|")
	line = strings.TrimSuffix(line, "|")
	parts := strings.Split(line, "|")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		out = append(out, strings.TrimSpace(p))
	}
	return out
}

func isMDSeparator(field string) bool {
	s := strings.ReplaceAll(field, " ", "")
	if s == "" {
		return false
	}
	for _, r := range s {
		if r != '-' && r != ':' {
			return false
		}
	}
	return true
}

func isIdentKey(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		switch {
		case r >= 'a' && r <= 'z':
		case r == '_':
		case i > 0 && r >= '0' && r <= '9':
		default:
			return false
		}
	}
	return true
}
