package embed_test

import (
	"io/fs"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/embed"
)

func TestProgramEmbedHasNoGoFiles(t *testing.T) {
	err := fs.WalkDir(embed.Program, "program", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(d.Name()), ".go") {
			t.Errorf("embed contains filtered .go file: %s", path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk embed: %v", err)
	}
}

func TestProgramEmbedNotEmpty(t *testing.T) {
	entries, err := fs.ReadDir(embed.Program, "program")
	if err != nil {
		t.Fatalf("ReadDir program: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("embedded program/ is empty")
	}
}
