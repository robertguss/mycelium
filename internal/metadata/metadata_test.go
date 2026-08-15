package metadata_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/metadata"
)

func TestParseFrontMatter(t *testing.T) {
	in := "+++\nid = \"DEC-001\"\ntitle = \"T\"\nstatus = \"Proposed\"\n+++\n\n# Body\n"
	doc, err := metadata.Parse([]byte(in))
	if err != nil {
		t.Fatal(err)
	}
	if doc.Meta["id"] != "DEC-001" || doc.Meta["title"] != "T" {
		t.Fatalf("meta = %#v", doc.Meta)
	}
	if !strings.Contains(doc.Body, "# Body") {
		t.Fatalf("body = %q", doc.Body)
	}
}

func TestBodyPlusPlusNotFrontMatter(t *testing.T) {
	in := "+++\nid = \"DEC-001\"\n+++\n\nbody\n+++\nid = \"FAKE\"\n+++\n"
	doc, err := metadata.Parse([]byte(in))
	if err != nil {
		t.Fatal(err)
	}
	if doc.Meta["id"] != "DEC-001" {
		t.Fatalf("id = %v", doc.Meta["id"])
	}
	if _, ok := doc.Meta["FAKE"]; ok {
		t.Fatal("body +++ must not start second front-matter block")
	}
	if !strings.Contains(doc.Body, "+++") || !strings.Contains(doc.Body, "FAKE") {
		t.Fatalf("body should keep +++ lines, got %q", doc.Body)
	}
}

func TestBOMStripped(t *testing.T) {
	in := append([]byte{0xEF, 0xBB, 0xBF}, []byte("+++\nid = \"x\"\n+++\n")...)
	doc, err := metadata.Parse(in)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Meta["id"] != "x" {
		t.Fatalf("meta = %#v", doc.Meta)
	}
}

func TestMissingClose(t *testing.T) {
	_, err := metadata.Parse([]byte("+++\nid = \"x\"\n"))
	if !errors.Is(err, metadata.ErrMissingClose) {
		t.Fatalf("err = %v", err)
	}
}

func TestMissingOpen(t *testing.T) {
	tests := []string{
		"id = \"x\"\n+++\n",
		" ---\nid = \"x\"\n---\n",
		" +++\nid = \"x\"\n+++\n",
	}
	for _, in := range tests {
		_, err := metadata.Parse([]byte(in))
		if !errors.Is(err, metadata.ErrMissingOpen) {
			t.Fatalf("Parse(%q) err = %v, want missing open", in, err)
		}
	}
}

func TestBadTOML(t *testing.T) {
	_, err := metadata.Parse([]byte("+++\nid = \n+++\n"))
	if !errors.Is(err, metadata.ErrTOML) {
		t.Fatalf("err = %v", err)
	}
}

func TestRequireKeys(t *testing.T) {
	in := "+++\nid = \"DEC-001\"\ntitle = \"T\"\n+++\n"
	doc, err := metadata.Parse([]byte(in))
	if err != nil {
		t.Fatal(err)
	}
	if err := metadata.RequireKeys(doc.Meta, []string{"id", "title"}); err != nil {
		t.Fatal(err)
	}
	err = metadata.RequireKeys(doc.Meta, []string{"id", "status"})
	if !errors.Is(err, metadata.ErrRequired) {
		t.Fatalf("err = %v", err)
	}
}

func TestBodyCannotMasquerade(t *testing.T) {
	// Without opening fence, body-looking content is not metadata.
	_, err := metadata.Parse([]byte("# Title\n\nid = \"DEC-001\"\n"))
	if !errors.Is(err, metadata.ErrMissingOpen) {
		t.Fatalf("err = %v", err)
	}
}
