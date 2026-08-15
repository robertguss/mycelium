// Package metadata reads +++ TOML front matter from artifact files.
package metadata

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

var utf8BOM = []byte{0xEF, 0xBB, 0xBF}

var (
	ErrMissingOpen  = errors.New("metadata: missing opening +++")
	ErrMissingClose = errors.New("metadata: missing closing +++")
	ErrTOML         = errors.New("metadata: invalid TOML")
	ErrRequired     = errors.New("metadata: missing required key")
)

// Document is front matter plus the body after the closing fence.
type Document struct {
	Meta map[string]any
	Body string
}

// Parse splits data into TOML front matter and body.
// Opening +++ must be at byte 0 after an optional UTF-8 BOM.
// Closing +++ is the first subsequent line that is exactly "+++".
// A +++ line in the body is body text, not a second front-matter block.
func Parse(data []byte) (Document, error) {
	data = bytes.TrimPrefix(data, utf8BOM)
	if !bytes.HasPrefix(data, []byte("+++")) {
		return Document{}, ErrMissingOpen
	}
	// Opening fence must be a full first line of exactly +++.
	rest := data[3:]
	if len(rest) == 0 {
		return Document{}, ErrMissingClose
	}
	if rest[0] == '\n' {
		rest = rest[1:]
	} else if len(rest) >= 2 && rest[0] == '\r' && rest[1] == '\n' {
		rest = rest[2:]
	} else {
		return Document{}, ErrMissingOpen
	}

	lines := splitKeepEnds(rest)
	closeIdx := -1
	for i, line := range lines {
		if lineContent(line) == "+++" {
			closeIdx = i
			break
		}
	}
	if closeIdx < 0 {
		return Document{}, ErrMissingClose
	}

	var tomlBuf strings.Builder
	for i := 0; i < closeIdx; i++ {
		tomlBuf.Write(stripEOL(lines[i]))
		tomlBuf.WriteByte('\n')
	}

	meta := map[string]any{}
	if err := toml.Unmarshal([]byte(tomlBuf.String()), &meta); err != nil {
		return Document{}, fmt.Errorf("%w: %v", ErrTOML, err)
	}

	var body strings.Builder
	for i := closeIdx + 1; i < len(lines); i++ {
		body.Write(lines[i])
	}
	return Document{Meta: meta, Body: body.String()}, nil
}

// RequireKeys reports the first missing or empty-string required key.
func RequireKeys(meta map[string]any, keys []string) error {
	for _, k := range keys {
		v, ok := meta[k]
		if !ok {
			return fmt.Errorf("%w: %s", ErrRequired, k)
		}
		if s, ok := v.(string); ok && s == "" {
			return fmt.Errorf("%w: %s", ErrRequired, k)
		}
	}
	return nil
}

func splitKeepEnds(data []byte) [][]byte {
	if len(data) == 0 {
		return nil
	}
	var lines [][]byte
	start := 0
	for i := 0; i < len(data); i++ {
		if data[i] == '\n' {
			lines = append(lines, data[start:i+1])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}

func stripEOL(line []byte) []byte {
	line = bytes.TrimSuffix(line, []byte("\n"))
	line = bytes.TrimSuffix(line, []byte("\r"))
	return line
}

func lineContent(line []byte) string {
	return string(stripEOL(line))
}
