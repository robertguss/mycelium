package embed_test

import (
	"testing"

	"github.com/robertguss/mycelium/internal/embed"
)

func TestProgramEmbedsKeep(t *testing.T) {
	data, err := embed.Program.ReadFile("program/.keep")
	if err != nil {
		t.Fatalf("ReadFile program/.keep: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("program/.keep is empty")
	}
}
