package pack_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/robertguss/mycelium/internal/pack"
)

func TestDiscoverOnePack(t *testing.T) {
	prog := t.TempDir()
	writeCoreDecision(t, prog)
	mustMkdir(t, filepath.Join(prog, "packs", "council"))

	r, err := pack.Discover(prog)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Packs) != 1 || r.Packs[0].Name != "council" {
		t.Fatalf("packs=%v", r.Packs)
	}
	if len(r.Collisions) != 0 {
		t.Fatalf("collisions=%v", r.Collisions)
	}
	p := r.Packs[0]
	assertHas(t, p.Namespaces, "CMP", "RPT", "RCL")
	assertHas(t, p.TypeKeys, "commissioning", "model-report", "reconciliation")
	assertHas(t, p.Homes, "reviews/commissioning", "reviews/reports", "reviews/reconciliations")
}

func TestDiscoverZeroPacks(t *testing.T) {
	prog := t.TempDir()
	writeCoreDecision(t, prog)

	r, err := pack.Discover(prog)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Packs) != 0 {
		t.Fatalf("packs=%v want none", r.Packs)
	}
	if len(r.Collisions) != 0 {
		t.Fatalf("collisions=%v", r.Collisions)
	}
}

func TestDiscoverZeroPacksMissingPacksDir(t *testing.T) {
	prog := t.TempDir()
	writeCoreDecision(t, prog)

	r, err := pack.Discover(prog)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Packs) != 0 {
		t.Fatalf("packs=%v", r.Packs)
	}
}

func TestCollisionNamespace(t *testing.T) {
	prog := t.TempDir()
	writeCoreDecision(t, prog)
	mustMkdir(t, filepath.Join(prog, "packs", "council"))
	writePackSchema(t, prog, "fixture-pack", "collide.schema.toml", `namespace = "CMP"
home = "fixture-home"
filename_pattern = "CMP-{NNN}-{slug}.md"
stage_scoped = false
digits = 3
required_front_matter = ["id"]
required_sections = ["Body"]
`)

	r, err := pack.Discover(prog)
	if err != nil {
		t.Fatal(err)
	}
	if !hasCollision(r, "namespace", "CMP") {
		t.Fatalf("want NS CMP collision, got %v", r.Collisions)
	}
}

func TestCollisionTypeKey(t *testing.T) {
	prog := t.TempDir()
	writeCoreDecision(t, prog)
	mustMkdir(t, filepath.Join(prog, "packs", "council"))
	writePackSchema(t, prog, "fixture-pack", "commissioning.schema.toml", `namespace = "ZZZ"
home = "fixture-home"
filename_pattern = "ZZZ-{NNN}-{slug}.md"
stage_scoped = false
digits = 3
required_front_matter = ["id"]
required_sections = ["Body"]
`)

	r, err := pack.Discover(prog)
	if err != nil {
		t.Fatal(err)
	}
	if !hasCollision(r, "type-key", "commissioning") {
		t.Fatalf("want type-key commissioning collision, got %v", r.Collisions)
	}
}

func TestCollisionHome(t *testing.T) {
	prog := t.TempDir()
	writeCoreDecision(t, prog)
	mustMkdir(t, filepath.Join(prog, "packs", "council"))
	writePackSchema(t, prog, "fixture-pack", "other.schema.toml", `namespace = "ZZZ"
home = "reviews/commissioning"
filename_pattern = "ZZZ-{NNN}-{slug}.md"
stage_scoped = false
digits = 3
required_front_matter = ["id"]
required_sections = ["Body"]
`)

	r, err := pack.Discover(prog)
	if err != nil {
		t.Fatal(err)
	}
	if !hasCollision(r, "home", "reviews/commissioning") {
		t.Fatalf("want home collision, got %v", r.Collisions)
	}
}

func TestCollisionVsCore(t *testing.T) {
	prog := t.TempDir()
	writeCoreDecision(t, prog)
	writePackSchema(t, prog, "fixture-pack", "steal.schema.toml", `namespace = "DEC"
home = "stolen"
filename_pattern = "DEC-{NNN}-{slug}.md"
stage_scoped = false
digits = 3
required_front_matter = ["id"]
required_sections = ["Body"]
`)

	r, err := pack.Discover(prog)
	if err != nil {
		t.Fatal(err)
	}
	if !hasCollision(r, "namespace", "DEC") {
		t.Fatalf("want NS DEC vs core collision, got %v", r.Collisions)
	}
	found := false
	for _, c := range r.Collisions {
		if c.Key == "DEC" && (c.Source == pack.CoreSource || c.Other == pack.CoreSource) {
			found = true
		}
	}
	if !found {
		t.Fatalf("collision must name %s: %v", pack.CoreSource, r.Collisions)
	}
}

func TestReviewsAllowed(t *testing.T) {
	if !pack.ReviewsAllowed([]pack.Pack{{Name: "council"}}) {
		t.Fatal("council present → true")
	}
	if pack.ReviewsAllowed(nil) {
		t.Fatal("absent → false")
	}
	if pack.ReviewsAllowed([]pack.Pack{{Name: "other"}}) {
		t.Fatal("non-council → false")
	}
}

func writeCoreDecision(t *testing.T, prog string) {
	t.Helper()
	dir := filepath.Join(prog, "templates")
	mustMkdir(t, dir)
	body := `namespace = "DEC"
home = "decisions"
filename_pattern = "DEC-{NNN}-{slug}.md"
stage_scoped = false
digits = 3
required_front_matter = ["id", "title", "date"]
required_sections = ["Context", "Decision", "Consequences"]
`
	if err := os.WriteFile(filepath.Join(dir, "decision.schema.toml"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writePackSchema(t *testing.T, prog, packName, file, body string) {
	t.Helper()
	dir := filepath.Join(prog, "packs", packName, "templates")
	mustMkdir(t, dir)
	if err := os.WriteFile(filepath.Join(dir, file), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func mustMkdir(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
}

func hasCollision(r pack.Result, kind, key string) bool {
	for _, c := range r.Collisions {
		if c.Kind == kind && c.Key == key {
			return true
		}
	}
	return false
}

func assertHas(t *testing.T, got []string, want ...string) {
	t.Helper()
	set := map[string]struct{}{}
	for _, g := range got {
		set[g] = struct{}{}
	}
	for _, w := range want {
		if _, ok := set[w]; !ok {
			t.Fatalf("missing %q in %v", w, got)
		}
	}
}
