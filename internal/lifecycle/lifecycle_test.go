package lifecycle_test

import (
	"testing"

	"github.com/robertguss/mycelium/internal/lifecycle"
)

func TestLegalNext(t *testing.T) {
	cases := map[string][]string{
		"spark":      {"exploring", "archived"},
		"exploring":  {"simmering", "clarified", "archived"},
		"simmering":  {"exploring", "archived"},
		"clarified":  {"handed-off", "archived"},
		"handed-off": {"archived"},
		"archived":   nil,
		"nope":       nil,
	}
	for from, want := range cases {
		got := lifecycle.LegalNext(from)
		if len(got) != len(want) {
			t.Fatalf("%s: got %v want %v", from, got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("%s[%d]=%q want %q", from, i, got[i], want[i])
			}
		}
	}
}

func TestLegalEdges(t *testing.T) {
	cases := []struct {
		from, to string
		want     bool
	}{
		{"spark", "exploring", true},
		{"spark", "archived", true},
		{"spark", "clarified", false},
		{"spark", "simmering", false},
		{"spark", "handed-off", false},
		{"exploring", "simmering", true},
		{"exploring", "clarified", true},
		{"exploring", "archived", true},
		{"exploring", "exploring", false},
		{"exploring", "handed-off", false},
		{"simmering", "exploring", true},
		{"simmering", "archived", true},
		{"simmering", "clarified", false},
		{"clarified", "archived", true},
		{"clarified", "handed-off", true},
		{"handed-off", "archived", true},
		{"handed-off", "exploring", false},
		{"archived", "exploring", false},
		{"archived", "archived", false},
	}
	for _, tc := range cases {
		if got := lifecycle.Legal(tc.from, tc.to); got != tc.want {
			t.Fatalf("Legal(%q,%q)=%v want %v", tc.from, tc.to, got, tc.want)
		}
	}
}

func TestAllowedTargets(t *testing.T) {
	want := []string{"exploring", "simmering", "clarified", "handed-off", "archived"}
	got := lifecycle.AllowedTargets()
	if len(got) != len(want) {
		t.Fatalf("got %v want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("[%d]=%q want %q", i, got[i], want[i])
		}
	}
	for _, bad := range []string{"spark"} {
		for _, tgt := range got {
			if tgt == bad {
				t.Fatalf("AllowedTargets must not include %q", bad)
			}
		}
	}
}

func TestIsWake(t *testing.T) {
	cases := []struct {
		from, to string
		want     bool
	}{
		{"simmering", "exploring", true},
		{"simmering", "archived", false},
		{"exploring", "exploring", false},
		{"spark", "exploring", false},
	}
	for _, tc := range cases {
		if got := lifecycle.IsWake(tc.from, tc.to); got != tc.want {
			t.Fatalf("IsWake(%q,%q)=%v want %v", tc.from, tc.to, got, tc.want)
		}
	}
}

func TestRevisitFlags(t *testing.T) {
	cases := []struct {
		to        string
		required  bool
		forbidden bool
	}{
		{"simmering", true, false},
		{"exploring", false, true},
		{"clarified", false, true},
		{"archived", false, true},
		{"handed-off", false, true},
	}
	for _, tc := range cases {
		if got := lifecycle.RevisitRequired(tc.to); got != tc.required {
			t.Fatalf("RevisitRequired(%q)=%v want %v", tc.to, got, tc.required)
		}
		if got := lifecycle.RevisitForbidden(tc.to); got != tc.forbidden {
			t.Fatalf("RevisitForbidden(%q)=%v want %v", tc.to, got, tc.forbidden)
		}
	}
}
