package idpath_test

import (
	"errors"
	"testing"

	"github.com/robertguss/mycelium/internal/idpath"
	"github.com/robertguss/mycelium/internal/slug"
)

func TestRoundTripAllTypes(t *testing.T) {
	slugs := []string{"garden-lighting", "a", "x-y-z"}
	for _, typ := range idpath.Types() {
		for _, s := range slugs {
			t.Run(typ.NS+"/"+s, func(t *testing.T) {
				path, err := idpath.Format(typ.NS, 1, s)
				if err != nil {
					t.Fatalf("Format: %v", err)
				}
				id, gotSlug, err := idpath.ParsePath(path)
				if err != nil {
					t.Fatalf("ParsePath(%q): %v", path, err)
				}
				if id.NS != typ.NS || id.N != 1 || gotSlug != s {
					t.Fatalf("ParsePath = {%s %d %q}, want {%s 1 %q}", id.NS, id.N, gotSlug, typ.NS, s)
				}
				path2, err := idpath.PathFor(id, gotSlug)
				if err != nil {
					t.Fatalf("PathFor: %v", err)
				}
				if path2 != path {
					t.Fatalf("PathFor inverse: got %q want %q", path2, path)
				}
			})
		}
	}
}

func TestRoundTripSlugifyAccepts(t *testing.T) {
	names := []string{"Garden lighting", "Café", "Hello_World"}
	for _, typ := range idpath.Types() {
		for _, name := range names {
			s, err := slug.Slugify(name)
			if err != nil {
				t.Fatalf("Slugify(%q): %v", name, err)
			}
			path, err := idpath.Format(typ.NS, 42, s)
			if err != nil {
				t.Fatalf("Format: %v", err)
			}
			id, gotSlug, err := idpath.ParsePath(path)
			if err != nil {
				t.Fatalf("ParsePath(%q): %v", path, err)
			}
			if gotSlug != s || id.NS != typ.NS || id.N != 42 {
				t.Fatalf("round-trip mismatch for %s / %q", typ.NS, s)
			}
		}
	}
}

func TestParseIDToken(t *testing.T) {
	tests := []struct {
		in      string
		wantNS  string
		wantN   int
		wantErr error
	}{
		{in: "DEC-001", wantNS: "DEC", wantN: 1},
		{in: "DEC-1", wantNS: "DEC", wantN: 1},
		{in: "PHASE-01", wantNS: "PHASE", wantN: 1},
		{in: "PHASE-1", wantNS: "PHASE", wantN: 1},
		{in: "MS-099", wantNS: "MS", wantN: 99},
		{in: "", wantErr: idpath.ErrEmpty},
		{in: "FOO-001", wantErr: idpath.ErrUnknown},
		{in: "dec-001", wantErr: idpath.ErrFormat},
		{in: "DEC-", wantErr: idpath.ErrFormat},
		{in: "DEC-001-extra", wantErr: idpath.ErrFormat},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := idpath.Parse(tt.in)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Parse(%q) err = %v, want %v", tt.in, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse(%q): %v", tt.in, err)
			}
			if got.NS != tt.wantNS || got.N != tt.wantN {
				t.Fatalf("Parse(%q) = %+v, want %s/%d", tt.in, got, tt.wantNS, tt.wantN)
			}
		})
	}
}

func TestBadFilenames(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{name: "unpadded DEC-1", path: "decisions/DEC-1-foo.md"},
		{name: "four digits DEC", path: "decisions/DEC-0001-foo.md"},
		{name: "three digits PHASE", path: "phases/PHASE-001-foo.md"},
		{name: "wrong home", path: "findings/DEC-001-foo.md"},
		{name: "unknown NS home", path: "nope/DEC-001-foo.md"},
		{name: "empty", path: ""},
		{name: "no slug", path: "decisions/DEC-001-.md"},
		{name: "no ext", path: "decisions/DEC-001-foo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := idpath.ParsePath(tt.path)
			if err == nil {
				t.Fatalf("ParsePath(%q) succeeded, want error", tt.path)
			}
		})
	}
}

func TestFormatZeroPad(t *testing.T) {
	got, err := idpath.Format("DEC", 1, "x")
	if err != nil {
		t.Fatal(err)
	}
	if got != "decisions/DEC-001-x.md" {
		t.Fatalf("got %q", got)
	}
	got, err = idpath.Format("PHASE", 1, "x")
	if err != nil {
		t.Fatal(err)
	}
	if got != "phases/PHASE-01-x.md" {
		t.Fatalf("got %q", got)
	}
}

func TestFormatDigitWidth(t *testing.T) {
	_, err := idpath.Format("DEC", 1000, "x")
	if !errors.Is(err, idpath.ErrFormat) {
		t.Fatalf("Format(DEC,1000) err = %v, want ErrFormat", err)
	}
	_, err = idpath.FormatID("DEC", 1000)
	if !errors.Is(err, idpath.ErrFormat) {
		t.Fatalf("FormatID(DEC,1000) err = %v, want ErrFormat", err)
	}
	_, err = idpath.Format("PHASE", 100, "x")
	if !errors.Is(err, idpath.ErrFormat) {
		t.Fatalf("Format(PHASE,100) err = %v, want ErrFormat", err)
	}

	path, err := idpath.Format("DEC", 999, "max")
	if err != nil {
		t.Fatal(err)
	}
	if path != "decisions/DEC-999-max.md" {
		t.Fatalf("got %q", path)
	}
	id, slug, err := idpath.ParsePath(path)
	if err != nil {
		t.Fatal(err)
	}
	if id.NS != "DEC" || id.N != 999 || slug != "max" {
		t.Fatalf("ParsePath = %+v %q", id, slug)
	}
	path2, err := idpath.PathFor(id, slug)
	if err != nil || path2 != path {
		t.Fatalf("PathFor inverse: %q %v", path2, err)
	}
}

func TestQuestionNotStageScoped(t *testing.T) {
	oq, err := idpath.LookupNS("OQ")
	if err != nil {
		t.Fatal(err)
	}
	if oq.StageScoped {
		t.Fatal("OQ must not be stage_scoped")
	}
	fnd, err := idpath.LookupNS("FND")
	if err != nil {
		t.Fatal(err)
	}
	if !fnd.StageScoped {
		t.Fatal("FND must be stage_scoped")
	}
}
