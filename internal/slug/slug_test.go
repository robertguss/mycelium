package slug_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/slug"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    string
		wantErr error
	}{
		{name: "spaces", in: "Garden lighting", want: "garden-lighting"},
		{name: "underscores", in: "garden_lighting", want: "garden-lighting"},
		{name: "punctuation", in: "Hello, World!", want: "hello-world"},
		{name: "accents", in: "Café naïve", want: "cafe-naive"},
		{name: "umlaut", in: "Zürich", want: "zurich"},
		{name: "collapse dashes", in: "a  --  b", want: "a-b"},
		{name: "trim dashes", in: " -foo- ", want: "foo"},
		{name: "empty", in: "", wantErr: slug.ErrEmpty},
		{name: "punctuation only", in: "...", wantErr: slug.ErrEmpty},
		{name: "dot", in: ".", wantErr: slug.ErrEmpty},
		{name: "dotdot", in: "..", wantErr: slug.ErrEmpty},
		{name: "overflow truncated", in: strings.Repeat("a", 100), want: strings.Repeat("a", 80)},
		{name: "overflow mid-word", in: strings.Repeat("ab", 50), want: strings.Repeat("ab", 40)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := slug.Slugify(tt.in)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Slugify(%q) err = %v, want %v", tt.in, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Slugify(%q) unexpected err: %v", tt.in, err)
			}
			if got != tt.want {
				t.Fatalf("Slugify(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
