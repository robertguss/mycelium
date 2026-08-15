// Package schema loads sidecar *.schema.toml files for registered types.
package schema

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

var (
	ErrInvalid  = errors.New("schema: invalid")
	ErrRequired = errors.New("schema: missing required field")
)

// Schema is one type's sidecar schema.toml.
type Schema struct {
	Namespace           string
	Home                string
	FilenamePattern     string
	StageScoped         bool
	Digits              int
	RequiredFrontMatter []string
	RequiredSections    []string
	Enums               map[string][]string
}

// Entry is one filesystem-registered type and its source paths.
type Entry struct {
	Key          string
	SchemaPath   string
	TemplatePath string
	Schema       Schema
}

// Parse decodes a *.schema.toml document from bytes.
func Parse(data []byte) (Schema, error) {
	var raw map[string]any
	if err := toml.Unmarshal(data, &raw); err != nil {
		return Schema{}, fmt.Errorf("%w: %v", ErrInvalid, err)
	}

	s := Schema{Enums: map[string][]string{}}

	var err error
	s.Namespace, err = reqString(raw, "namespace")
	if err != nil {
		return Schema{}, err
	}
	s.Home, err = reqString(raw, "home")
	if err != nil {
		return Schema{}, err
	}
	s.FilenamePattern, err = reqString(raw, "filename_pattern")
	if err != nil {
		return Schema{}, err
	}
	ss, ok := raw["stage_scoped"]
	if !ok {
		return Schema{}, fmt.Errorf("%w: stage_scoped", ErrRequired)
	}
	s.StageScoped, ok = ss.(bool)
	if !ok {
		return Schema{}, fmt.Errorf("%w: stage_scoped", ErrInvalid)
	}
	digits, err := asInt(raw["digits"])
	if err != nil {
		return Schema{}, fmt.Errorf("%w: digits", ErrRequired)
	}
	if digits < 1 {
		return Schema{}, fmt.Errorf("%w: digits", ErrInvalid)
	}
	s.Digits = digits

	s.RequiredFrontMatter, err = asStringSlice(raw["required_front_matter"])
	if err != nil {
		return Schema{}, fmt.Errorf("%w: required_front_matter", ErrRequired)
	}
	s.RequiredSections, err = asStringSlice(raw["required_sections"])
	if err != nil {
		return Schema{}, fmt.Errorf("%w: required_sections", ErrRequired)
	}

	if enumsRaw, ok := raw["enums"]; ok {
		enumsMap, ok := enumsRaw.(map[string]any)
		if !ok {
			return Schema{}, fmt.Errorf("%w: enums", ErrInvalid)
		}
		for field, v := range enumsMap {
			entry, ok := v.(map[string]any)
			if !ok {
				return Schema{}, fmt.Errorf("%w: enums.%s", ErrInvalid, field)
			}
			vals, err := asStringSlice(entry["values"])
			if err != nil {
				return Schema{}, fmt.Errorf("%w: enums.%s.values", ErrRequired, field)
			}
			s.Enums[field] = vals
		}
	}

	return s, nil
}

// Load reads and parses a schema.toml file from disk.
func Load(path string) (Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Schema{}, err
	}
	return Parse(data)
}

// Discover scans core templates followed by pack templates. Duplicate stems
// retain the first registration.
func Discover(root string) ([]Entry, error) {
	seen := map[string]struct{}{}
	var out []Entry
	coreDir := filepath.Join(root, "program", "templates")
	if err := discoverDir(coreDir, false, seen, &out); err != nil {
		return nil, err
	}

	packsDir := filepath.Join(root, "program", "packs")
	packs, err := os.ReadDir(packsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return out, nil
		}
		return nil, err
	}
	for _, pack := range packs {
		if !pack.IsDir() {
			continue
		}
		dir := filepath.Join(packsDir, pack.Name(), "templates")
		if err := discoverDir(dir, true, seen, &out); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func discoverDir(dir string, optional bool, seen map[string]struct{}, out *[]Entry) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if optional && os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".schema.toml") {
			continue
		}
		key := strings.TrimSuffix(name, ".schema.toml")
		if _, ok := seen[key]; ok {
			continue
		}
		schemaPath := filepath.Join(dir, name)
		s, err := Load(schemaPath)
		if err != nil {
			return fmt.Errorf("%s: %w", schemaPath, err)
		}
		seen[key] = struct{}{}
		*out = append(*out, Entry{
			Key:          key,
			SchemaPath:   schemaPath,
			TemplatePath: filepath.Join(dir, key+".md"),
			Schema:       s,
		})
	}
	return nil
}

func reqString(raw map[string]any, key string) (string, error) {
	v, ok := raw[key]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrRequired, key)
	}
	s, ok := v.(string)
	if !ok || s == "" {
		return "", fmt.Errorf("%w: %s", ErrInvalid, key)
	}
	return s, nil
}

func asStringSlice(v any) ([]string, error) {
	arr, ok := v.([]any)
	if !ok {
		return nil, ErrInvalid
	}
	out := make([]string, 0, len(arr))
	for _, item := range arr {
		s, ok := item.(string)
		if !ok {
			return nil, ErrInvalid
		}
		out = append(out, s)
	}
	return out, nil
}

func asInt(v any) (int, error) {
	switch n := v.(type) {
	case int64:
		return int(n), nil
	case int:
		return n, nil
	case float64:
		if n != float64(int(n)) {
			return 0, ErrInvalid
		}
		return int(n), nil
	default:
		return 0, ErrInvalid
	}
}
