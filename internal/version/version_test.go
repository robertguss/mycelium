package version_test

import (
	"testing"

	"github.com/robertguss/mycelium/internal/version"
)

func TestDefaultVersion(t *testing.T) {
	if version.Version != "0.1.0-dev" {
		t.Fatalf("Version = %q, want %q", version.Version, "0.1.0-dev")
	}
}
