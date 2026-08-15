// Package pack discovers program/packs/*/ presence-registration and collisions.
package pack

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/robertguss/mycelium/internal/schema"
)

const (
	// CoreSource labels core template registrations in collision messages.
	CoreSource = "program/templates/"
	// CouncilName is the only pack this phase (OQ-003).
	CouncilName = "council"
)

// Pack is one presence-registered directory under program/packs/<name>/.
type Pack struct {
	Name       string
	Namespaces []string
	TypeKeys   []string
	Homes      []string
}

// Collision is one shared NS, type key, or home between two registrations.
type Collision struct {
	Kind   string // "namespace" | "type-key" | "home"
	Key    string
	Source string // first claimant
	Other  string // second claimant
}

// Result is discover output: packs found and any collisions.
type Result struct {
	Packs      []Pack
	Collisions []Collision
}

// councilClaims apply when program/packs/council/ exists even without templates.
var councilClaims = Pack{
	Name:       CouncilName,
	Namespaces: []string{"CMP", "RPT", "RCL"},
	TypeKeys:   []string{"commissioning", "model-report", "reconciliation"},
	Homes:      []string{"reviews/commissioning", "reviews/reports", "reviews/reconciliations"},
}

// Discover scans programDir for core templates and program/packs/*/ directories.
// Presence-is-registration: no registry file. Collisions FAIL at check time.
func Discover(programDir string) (Result, error) {
	var out Result
	ns := map[string]string{}
	typeKeys := map[string]string{}
	homes := map[string]string{}

	core, err := loadCoreClaims(filepath.Join(programDir, "templates"))
	if err != nil {
		return out, err
	}
	for _, c := range core.Namespaces {
		claim(ns, "namespace", c, CoreSource, &out.Collisions)
	}
	for _, c := range core.TypeKeys {
		claim(typeKeys, "type-key", c, CoreSource, &out.Collisions)
	}
	for _, c := range core.Homes {
		claim(homes, "home", c, CoreSource, &out.Collisions)
	}

	packsDir := filepath.Join(programDir, "packs")
	entries, err := os.ReadDir(packsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return out, nil
		}
		return out, err
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		p, err := loadPack(filepath.Join(packsDir, name), name)
		if err != nil {
			return out, err
		}
		out.Packs = append(out.Packs, p)
		for _, c := range p.Namespaces {
			claim(ns, "namespace", c, name, &out.Collisions)
		}
		for _, c := range p.TypeKeys {
			claim(typeKeys, "type-key", c, name, &out.Collisions)
		}
		for _, c := range p.Homes {
			claim(homes, "home", c, name, &out.Collisions)
		}
	}
	return out, nil
}

// ReviewsAllowed is true IFF a pack directory named council is present.
func ReviewsAllowed(packs []Pack) bool {
	for _, p := range packs {
		if p.Name == CouncilName {
			return true
		}
	}
	return false
}

// Message formats a collision for the four-line teaching error What line.
func (c Collision) Message() string {
	return fmt.Sprintf("pack %s collision: %s claimed by %s and %s", c.Kind, c.Key, c.Source, c.Other)
}

func claim(m map[string]string, kind, key, source string, out *[]Collision) {
	if prev, ok := m[key]; ok {
		if prev == source {
			return
		}
		*out = append(*out, Collision{Kind: kind, Key: key, Source: prev, Other: source})
		return
	}
	m[key] = source
}

func loadCoreClaims(templatesDir string) (Pack, error) {
	p := Pack{Name: CoreSource}
	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return p, nil
		}
		return p, err
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".schema.toml") {
			continue
		}
		stem := strings.TrimSuffix(e.Name(), ".schema.toml")
		s, err := schema.Load(filepath.Join(templatesDir, e.Name()))
		if err != nil {
			return p, fmt.Errorf("%s: %w", e.Name(), err)
		}
		p.TypeKeys = appendUnique(p.TypeKeys, stem)
		p.Namespaces = appendUnique(p.Namespaces, s.Namespace)
		p.Homes = appendUnique(p.Homes, s.Home)
	}
	return p, nil
}

func loadPack(dir, name string) (Pack, error) {
	p := Pack{Name: name}
	if name == CouncilName {
		p.Namespaces = append([]string{}, councilClaims.Namespaces...)
		p.TypeKeys = append([]string{}, councilClaims.TypeKeys...)
		p.Homes = append([]string{}, councilClaims.Homes...)
	}
	templatesDir := filepath.Join(dir, "templates")
	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return p, nil
		}
		return p, err
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".schema.toml") {
			continue
		}
		stem := strings.TrimSuffix(e.Name(), ".schema.toml")
		s, err := schema.Load(filepath.Join(templatesDir, e.Name()))
		if err != nil {
			return p, fmt.Errorf("pack %s %s: %w", name, e.Name(), err)
		}
		p.TypeKeys = appendUnique(p.TypeKeys, stem)
		p.Namespaces = appendUnique(p.Namespaces, s.Namespace)
		p.Homes = appendUnique(p.Homes, s.Home)
	}
	return p, nil
}

func appendUnique(slice []string, v string) []string {
	for _, s := range slice {
		if s == v {
			return slice
		}
	}
	return append(slice, v)
}
