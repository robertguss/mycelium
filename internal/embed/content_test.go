package embed_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	toml "github.com/pelletier/go-toml/v2"

	"github.com/robertguss/mycelium/internal/embed"
)

type typeSchema struct {
	Namespace            string              `toml:"namespace"`
	Home                 string              `toml:"home"`
	FilenamePattern      string              `toml:"filename_pattern"`
	StageScoped          bool                `toml:"stage_scoped"`
	Digits               int                 `toml:"digits"`
	RequiredFrontMatter  []string            `toml:"required_front_matter"`
	RequiredSections     []string            `toml:"required_sections"`
	Enums                map[string]enumList `toml:"enums"`
}

type enumList struct {
	Values []string `toml:"values"`
}

type tierFile struct {
	Name  string   `toml:"name"`
	Emits []string `toml:"emits"`
	Binds []string `toml:"binds"`
}

type typeExpect struct {
	namespace   string
	home        string
	pattern     string
	digits      int
	stageScoped bool
	fm          []string
	sections    []string
}

var registered = map[string]typeExpect{
	"decision": {
		namespace: "DEC", home: "decisions", pattern: "DEC-{NNN}-{slug}.md",
		digits: 3, stageScoped: false,
		fm: []string{"id", "title", "status", "date", "owner"},
		sections: []string{"Context", "Decision", "Rationale", "Consequences", "Alternatives Considered", "Risks", "Revisit Triggers", "Approval"},
	},
	"assumption": {
		namespace: "ASM", home: "assumptions", pattern: "ASM-{NNN}-{slug}.md",
		digits: 3, stageScoped: false,
		fm: []string{"id", "title", "status", "date", "attached_to"},
		sections: []string{"Statement", "Falsifier", "Implications", "Revisit Triggers"},
	},
	"evidence": {
		namespace: "EVD", home: "evidence", pattern: "EVD-{NNN}-{slug}.md",
		digits: 3, stageScoped: false,
		fm: []string{"id", "title", "status", "date", "source"},
		sections: []string{"Claim", "Source", "Observation", "Limitations", "Revalidation Trigger"},
	},
	"spike": {
		namespace: "SPK", home: "spikes", pattern: "SPK-{NNN}-{slug}.md",
		digits: 3, stageScoped: false,
		fm: []string{"id", "title", "status", "date", "decision_at_stake"},
		sections: []string{"Method", "Commands and Artifacts", "Results", "Limitations", "Architectural Consequence", "Cleanup", "Reproduction Instructions"},
	},
	"finding": {
		namespace: "FND", home: "findings", pattern: "FND-{NNN}-{slug}.md",
		digits: 3, stageScoped: true,
		fm: []string{"id", "title", "severity", "confidence", "date"},
		sections: []string{"Problem", "Evidence", "Failure Scenario", "Impact", "Root Cause", "Required Correction", "Residual Risk"},
	},
	"recommendation": {
		namespace: "REC", home: "recommendations", pattern: "REC-{NNN}-{slug}.md",
		digits: 3, stageScoped: true,
		fm: []string{"id", "title", "classification", "confidence", "date"},
		sections: []string{"Recommendation", "Requirements and Constraints", "Rationale", "Evidence", "Tradeoffs", "Alternatives Considered", "Revisit Triggers"},
	},
	"requirement": {
		namespace: "REQ", home: "requirements", pattern: "REQ-{NNN}-{slug}.md",
		digits: 3, stageScoped: true,
		fm: []string{"id", "title", "priority", "date", "phase"},
		sections: []string{"Requirement", "Rationale", "Acceptance Evidence", "Exceptions"},
	},
	"question": {
		namespace: "OQ", home: "questions", pattern: "OQ-{NNN}-{slug}.md",
		digits: 3, stageScoped: false,
		fm: []string{"id", "title", "agreement", "date"},
		sections: []string{"Question", "Context", "Positions", "Crux", "Disposition"},
	},
	"risk": {
		namespace: "RSK", home: "risks", pattern: "RSK-{NNN}-{slug}.md",
		digits: 3, stageScoped: false,
		fm: []string{"id", "title", "severity", "likelihood", "date"},
		sections: []string{"Description", "Impact", "Mitigation", "Residual Risk", "Revisit Triggers"},
	},
	"phase": {
		namespace: "PHASE", home: "phases", pattern: "PHASE-{NN}-{slug}.md",
		digits: 2, stageScoped: false,
		fm: []string{"id", "title", "status", "date"},
		sections: []string{"Entry Criteria", "Scope", "Explicit Non-Goals", "Exit Criteria"},
	},
	"milestone": {
		namespace: "MS", home: "milestones", pattern: "MS-{NNN}-{slug}.md",
		digits: 3, stageScoped: false,
		fm: []string{"id", "title", "phase", "date"},
		sections: []string{"Outcome", "Prerequisites", "Acceptance Evidence"},
	},
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "../.."))
}

func TestRegisteredSchemasAndTemplates(t *testing.T) {
	root := moduleRoot(t)
	tplDir := filepath.Join(root, "program", "templates")
	for key, want := range registered {
		t.Run(key, func(t *testing.T) {
			schemaPath := filepath.Join(tplDir, key+".schema.toml")
			raw, err := os.ReadFile(schemaPath)
			if err != nil {
				t.Fatalf("read schema: %v", err)
			}
			var sch typeSchema
			if err := toml.Unmarshal(raw, &sch); err != nil {
				t.Fatalf("parse schema: %v", err)
			}
			if sch.Namespace != want.namespace {
				t.Errorf("namespace: got %q want %q", sch.Namespace, want.namespace)
			}
			if sch.Home != want.home {
				t.Errorf("home: got %q want %q", sch.Home, want.home)
			}
			if sch.FilenamePattern != want.pattern {
				t.Errorf("filename_pattern: got %q want %q", sch.FilenamePattern, want.pattern)
			}
			if sch.Digits != want.digits {
				t.Errorf("digits: got %d want %d", sch.Digits, want.digits)
			}
			if sch.StageScoped != want.stageScoped {
				t.Errorf("stage_scoped: got %v want %v", sch.StageScoped, want.stageScoped)
			}
			if !equalStrings(sch.RequiredFrontMatter, want.fm) {
				t.Errorf("required_front_matter: got %#v want %#v", sch.RequiredFrontMatter, want.fm)
			}
			if !equalStrings(sch.RequiredSections, want.sections) {
				t.Errorf("required_sections: got %#v want %#v", sch.RequiredSections, want.sections)
			}

			tplPath := filepath.Join(tplDir, key+".md")
			body, err := os.ReadFile(tplPath)
			if err != nil {
				t.Fatalf("read template: %v", err)
			}
			text := string(body)
			if !strings.HasPrefix(text, "+++") {
				t.Fatal("template must start with +++")
			}
			for _, tok := range []string{"{{ID}}", "{{TITLE}}", "{{SLUG}}", "{{DATE}}"} {
				if !strings.Contains(text, tok) {
					t.Errorf("template missing token %s", tok)
				}
			}
			for _, sec := range want.sections {
				heading := "## " + sec
				if !strings.Contains(text, heading) {
					t.Errorf("template missing H2 %q", sec)
				}
			}
		})
	}
}

func TestTierFiles(t *testing.T) {
	root := moduleRoot(t)
	cases := []struct {
		file  string
		name  string
		emits []string
		binds []string
	}{
		{
			file: "focused.toml", name: "focused",
			emits: []string{},
			binds: []string{"manifest", "log.md", "CONTEXT.md", "AGENTS.md", "program/", "index.md"},
		},
		{
			file: "standard.toml", name: "standard",
			emits: []string{"decisions/", "assumptions/", "evidence/", "questions/", "risks/"},
			binds: []string{
				"manifest", "log.md", "CONTEXT.md", "AGENTS.md", "program/", "index.md",
				"decisions/", "assumptions/", "evidence/", "questions/", "risks/",
			},
		},
		{
			file: "high-assurance.toml", name: "high-assurance",
			emits: []string{
				"decisions/", "assumptions/", "evidence/", "questions/", "risks/",
				"spikes/", "findings/", "recommendations/", "requirements/",
				"phases/", "milestones/",
			},
			binds: []string{
				"manifest", "log.md", "CONTEXT.md", "AGENTS.md", "program/", "index.md",
				"decisions/", "assumptions/", "evidence/", "questions/", "risks/",
				"spikes/", "findings/", "recommendations/", "requirements/",
				"phases/", "milestones/",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			raw, err := os.ReadFile(filepath.Join(root, "program", "tiers", tc.file))
			if err != nil {
				t.Fatalf("read: %v", err)
			}
			var tier tierFile
			if err := toml.Unmarshal(raw, &tier); err != nil {
				t.Fatalf("parse: %v", err)
			}
			if tier.Name != tc.name {
				t.Errorf("name: got %q want %q", tier.Name, tc.name)
			}
			if !equalStrings(tier.Emits, tc.emits) {
				t.Errorf("emits: got %#v want %#v", tier.Emits, tc.emits)
			}
			if !equalStrings(tier.Binds, tc.binds) {
				t.Errorf("binds: got %#v want %#v", tier.Binds, tc.binds)
			}
		})
	}
}

func TestSkeletonAndSkillExist(t *testing.T) {
	root := moduleRoot(t)
	paths := []string{
		"program/skeleton/README.md",
		"program/skeleton/log.md",
		"program/skeleton/CONTEXT.md",
		"program/skeleton/AGENTS.md",
		"program/skeleton/gitignore",
		"program/skeleton/index.md",
		"program/skills/mycelium-cli/SKILL.md",
		"program/skills/spark/SKILL.md",
		"program/skills/wake/SKILL.md",
		"program/skills/portfolio/SKILL.md",
	}
	for _, rel := range paths {
		p := filepath.Join(root, rel)
		info, err := os.Stat(p)
		if err != nil {
			t.Errorf("missing %s: %v", rel, err)
			continue
		}
		if info.Size() == 0 {
			t.Errorf("%s is empty", rel)
		}
	}
	gi, err := os.ReadFile(filepath.Join(root, "program/skeleton/gitignore"))
	if err != nil {
		t.Fatalf("gitignore: %v", err)
	}
	if strings.TrimSpace(string(gi)) != ".mycelium/lock" {
		t.Errorf("gitignore content: got %q want .mycelium/lock", strings.TrimSpace(string(gi)))
	}
}

func TestEmbedProgramHasRegisteredContent(t *testing.T) {
	for key := range registered {
		if _, err := embed.Program.ReadFile("program/templates/" + key + ".schema.toml"); err != nil {
			t.Errorf("embed missing schema %s: %v", key, err)
		}
		if _, err := embed.Program.ReadFile("program/templates/" + key + ".md"); err != nil {
			t.Errorf("embed missing template %s: %v", key, err)
		}
	}
	for _, name := range []string{"focused.toml", "standard.toml", "high-assurance.toml"} {
		if _, err := embed.Program.ReadFile("program/tiers/" + name); err != nil {
			t.Errorf("embed missing tier %s: %v", name, err)
		}
	}
	for _, skill := range []string{
		"mycelium-cli", "spark", "wake", "portfolio",
	} {
		path := "program/skills/" + skill + "/SKILL.md"
		if _, err := embed.Program.ReadFile(path); err != nil {
			t.Errorf("embed missing skill %s: %v", skill, err)
		}
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
