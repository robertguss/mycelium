// Package journal reads and writes .mycelium/journal.json.
package journal

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
)

const RelPath = ".mycelium/journal.json"

var ErrNotExist = errors.New("journal: not exist")

// Rename is one staged→final atomic rename intent.
type Rename struct {
	From string `json:"from"`
	To   string `json:"to"`
	Done bool   `json:"done"`
}

// Journal is the on-disk interrupted-operation record.
type Journal struct {
	SchemaVersion int      `json:"schema_version"`
	Op            string   `json:"op"`
	Type          *string  `json:"type"`
	Title         string   `json:"title"`
	OriginalID    string   `json:"original_id"`
	StartedAt     string   `json:"started_at"`
	StagedDir     string   `json:"staged_dir"`
	Renames       []Rename `json:"renames"`
	LogLine       string   `json:"log_line"`
	Argv          []string `json:"argv"`
}

// Path returns the absolute journal path under root.
func Path(root string) string {
	return filepath.Join(root, RelPath)
}

// Load reads the journal. Returns ErrNotExist when absent.
func Load(root string) (*Journal, error) {
	b, err := os.ReadFile(Path(root))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotExist
		}
		return nil, err
	}
	var j Journal
	if err := json.Unmarshal(b, &j); err != nil {
		return nil, err
	}
	return &j, nil
}

// Save writes the journal atomically (temp + rename in .mycelium/).
func Save(root string, j *Journal) error {
	dir := filepath.Join(root, ".mycelium")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	tmp, err := os.CreateTemp(dir, "journal-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(b); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	return os.Rename(tmpName, Path(root))
}

// Remove deletes the journal if present.
func Remove(root string) error {
	err := os.Remove(Path(root))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// TypeString returns the type field or "" when null.
func (j *Journal) TypeString() string {
	if j == nil || j.Type == nil {
		return ""
	}
	return *j.Type
}

// SetType sets type to s, or null when s is empty.
func (j *Journal) SetType(s string) {
	if s == "" {
		j.Type = nil
		return
	}
	j.Type = &s
}

// Matches reports whether incoming argv matches the journal argv,
// or the same op+type+title triplet.
func Matches(j *Journal, argv []string, op, typ, title string) bool {
	if j == nil {
		return false
	}
	if slices.Equal(j.Argv, argv) {
		return true
	}
	return j.Op == op && j.TypeString() == typ && j.Title == title
}
