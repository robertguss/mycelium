package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/idpath"
)

func TestCheckSupersedeIFFHappy(t *testing.T) {
	root := t.TempDir()
	writeMinimalPair(t, root, true)

	var findings []Finding
	arts := []artifactFile{
		{Rel: "decisions/DEC-001-old.md", Home: "decisions", IDStr: "DEC-001", ID: mustParseID(t, "DEC-001")},
		{Rel: "decisions/DEC-002-new.md", Home: "decisions", IDStr: "DEC-002", ID: mustParseID(t, "DEC-002")},
	}
	checkSupersedeIFF(root, arts, collectFindings(&findings))
	if len(findings) != 0 {
		t.Fatalf("findings=%v", findings)
	}
}

func TestCheckSupersedeIFFBrokenBidirectional(t *testing.T) {
	root := t.TempDir()
	decDir := filepath.Join(root, "decisions")
	if err := os.MkdirAll(decDir, 0o755); err != nil {
		t.Fatal(err)
	}
	old := `+++
id = "DEC-001"
status = "Superseded"
superseded_by = "DEC-002"
+++

## Context

none
`
	new := `+++
id = "DEC-002"
status = "Accepted"
+++

## Context

none
`
	if err := os.WriteFile(filepath.Join(decDir, "DEC-001-old.md"), []byte(old), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(decDir, "DEC-002-new.md"), []byte(new), 0o644); err != nil {
		t.Fatal(err)
	}
	var findings []Finding
	arts := []artifactFile{
		{Rel: "decisions/DEC-001-old.md", Home: "decisions", IDStr: "DEC-001", ID: mustParseID(t, "DEC-001")},
		{Rel: "decisions/DEC-002-new.md", Home: "decisions", IDStr: "DEC-002", ID: mustParseID(t, "DEC-002")},
	}
	checkSupersedeIFF(root, arts, collectFindings(&findings))
	if !hasText(findings, "peer supersedes") {
		t.Fatalf("findings=%v", findings)
	}
}

func TestCheckSupersedeIFFOneToOneInbound(t *testing.T) {
	root := t.TempDir()
	decDir := filepath.Join(root, "decisions")
	if err := os.MkdirAll(decDir, 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"DEC-001-a.md": `+++
id = "DEC-001"
status = "Superseded"
superseded_by = "DEC-003"
+++

## Context

none
`,
		"DEC-002-b.md": `+++
id = "DEC-002"
status = "Superseded"
superseded_by = "DEC-003"
+++

## Context

none
`,
		"DEC-003-c.md": `+++
id = "DEC-003"
status = "Accepted"
supersedes = "DEC-001"
+++

## Context

none
`,
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(decDir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	var findings []Finding
	arts := []artifactFile{
		{Rel: "decisions/DEC-001-a.md", Home: "decisions", IDStr: "DEC-001", ID: mustParseID(t, "DEC-001")},
		{Rel: "decisions/DEC-002-b.md", Home: "decisions", IDStr: "DEC-002", ID: mustParseID(t, "DEC-002")},
		{Rel: "decisions/DEC-003-c.md", Home: "decisions", IDStr: "DEC-003", ID: mustParseID(t, "DEC-003")},
	}
	checkSupersedeIFF(root, arts, collectFindings(&findings))
	if !hasText(findings, "multiple inbound superseded_by") {
		t.Fatalf("findings=%v", findings)
	}
}

func TestLogLineREAcceptsSupersede(t *testing.T) {
	line := "2026-08-15\tsupersede\tDEC-001\tDEC-001 -> DEC-002"
	if !logLineRE.MatchString(line) {
		t.Fatalf("logLineRE rejected %q", line)
	}
	handoff := "2026-08-15\thandoff\tHO-001\tclarified -> handed-off"
	if !logLineRE.MatchString(handoff) {
		t.Fatalf("logLineRE rejected handoff %q", handoff)
	}
	bad := "2026-08-15\tcouncil\tDEC-001\tnope"
	if logLineRE.MatchString(bad) {
		t.Fatal("council must not match")
	}
}

func writeMinimalPair(t *testing.T, root string, linked bool) {
	t.Helper()
	decDir := filepath.Join(root, "decisions")
	if err := os.MkdirAll(decDir, 0o755); err != nil {
		t.Fatal(err)
	}
	old := `+++
id = "DEC-001"
status = "Superseded"
superseded_by = "DEC-002"
+++

## Context

none
`
	new := `+++
id = "DEC-002"
status = "Accepted"
supersedes = "DEC-001"
+++

## Context

none
`
	if !linked {
		new = strings.Replace(new, "supersedes = \"DEC-001\"\n", "", 1)
	}
	if err := os.WriteFile(filepath.Join(decDir, "DEC-001-old.md"), []byte(old), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(decDir, "DEC-002-new.md"), []byte(new), 0o644); err != nil {
		t.Fatal(err)
	}
}

func mustParseID(t *testing.T, s string) idpath.ID {
	t.Helper()
	id, err := idpath.Parse(s)
	if err != nil {
		t.Fatal(err)
	}
	return id
}
