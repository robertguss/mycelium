package sparring_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/sparring"
)

func TestParseAgreement(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in      string
		want    sparring.Agreement
		wantErr bool
	}{
		{"open", sparring.Open, false},
		{"aligned", sparring.Aligned, false},
		{"agree-to-disagree", sparring.AgreeToDisagree, false},
		{"aligned ", sparring.Aligned, false},
		{"", "", true},
		{"maybe", "", true},
		{"Agree-To-Disagree", "", true},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			t.Parallel()
			got, err := sparring.ParseAgreement(tc.in)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("ParseAgreement(%q) err=nil want error", tc.in)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseAgreement(%q) err=%v", tc.in, err)
			}
			if got != tc.want {
				t.Fatalf("ParseAgreement(%q)=%q want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestRequiredH2(t *testing.T) {
	t.Parallel()
	base := []string{"Question", "Context", "Positions", "Disposition"}
	disputed := []string{"Question", "Context", "Positions", "Disposition", "Reasons", "Crux"}
	cases := []struct {
		a    sparring.Agreement
		want []string
	}{
		{sparring.Open, base},
		{sparring.Aligned, base},
		{sparring.AgreeToDisagree, disputed},
	}
	for _, tc := range cases {
		t.Run(string(tc.a), func(t *testing.T) {
			t.Parallel()
			got := sparring.RequiredH2(tc.a)
			if !slices.Equal(got, tc.want) {
				t.Fatalf("RequiredH2(%q)=%v want %v", tc.a, got, tc.want)
			}
		})
	}
}

func TestRequiredH3(t *testing.T) {
	t.Parallel()
	pair := []string{"Human", "Agent"}
	cases := []struct {
		name string
		a    sparring.Agreement
		h2   string
		want []string
	}{
		{"open-positions", sparring.Open, "Positions", nil},
		{"aligned-crux", sparring.Aligned, "Crux", nil},
		{"disputed-positions", sparring.AgreeToDisagree, "Positions", pair},
		{"disputed-reasons", sparring.AgreeToDisagree, "Reasons", pair},
		{"disputed-crux", sparring.AgreeToDisagree, "Crux", pair},
		{"disputed-question", sparring.AgreeToDisagree, "Question", nil},
		{"disputed-disposition", sparring.AgreeToDisagree, "Disposition", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := sparring.RequiredH3(tc.a, tc.h2)
			if !slices.Equal(got, tc.want) {
				t.Fatalf("RequiredH3(%q,%q)=%v want %v", tc.a, tc.h2, got, tc.want)
			}
		})
	}
}

func TestMissingHeadings(t *testing.T) {
	t.Parallel()
	openBody := strings.Join([]string{
		"## Question",
		"",
		"<!-- fill -->",
		"",
		"## Context",
		"",
		"<!-- fill -->",
		"",
		"## Positions",
		"",
		"<!-- fill -->",
		"",
		"## Disposition",
		"",
		"<!-- fill -->",
	}, "\n")
	alignedExtraCrux := openBody + "\n\n## Crux\n\nextra\n"
	disputedComplete := strings.Join([]string{
		"## Question",
		"",
		"q",
		"",
		"## Context",
		"",
		"c",
		"",
		"## Positions",
		"",
		"### Human",
		"",
		"h",
		"",
		"### Agent",
		"",
		"a",
		"",
		"## Reasons",
		"",
		"### Human",
		"",
		"h",
		"",
		"### Agent",
		"",
		"a",
		"",
		"## Crux",
		"",
		"### Human",
		"",
		"h",
		"",
		"### Agent",
		"",
		"a",
		"",
		"## Disposition",
		"",
		"d",
	}, "\n")
	disputedMissingCrux := strings.Join([]string{
		"## Question",
		"",
		"q",
		"",
		"## Context",
		"",
		"c",
		"",
		"## Positions",
		"",
		"### Human",
		"",
		"h",
		"",
		"### Agent",
		"",
		"a",
		"",
		"## Reasons",
		"",
		"### Human",
		"",
		"h",
		"",
		"### Agent",
		"",
		"a",
		"",
		"## Disposition",
		"",
		"d",
	}, "\n")
	disputedMissingHuman := strings.Join([]string{
		"## Question",
		"",
		"q",
		"",
		"## Context",
		"",
		"c",
		"",
		"## Positions",
		"",
		"### Agent",
		"",
		"a",
		"",
		"## Reasons",
		"",
		"### Human",
		"",
		"h",
		"",
		"### Agent",
		"",
		"a",
		"",
		"## Crux",
		"",
		"### Human",
		"",
		"h",
		"",
		"### Agent",
		"",
		"a",
		"",
		"## Disposition",
		"",
		"d",
	}, "\n")

	cases := []struct {
		name string
		a    sparring.Agreement
		body string
		want []string
	}{
		{"open-fill-positions", sparring.Open, openBody, nil},
		{"aligned-extra-crux", sparring.Aligned, alignedExtraCrux, nil},
		{"disputed-complete", sparring.AgreeToDisagree, disputedComplete, nil},
		{"disputed-missing-crux", sparring.AgreeToDisagree, disputedMissingCrux, []string{"## Crux"}},
		{"disputed-missing-human", sparring.AgreeToDisagree, disputedMissingHuman, []string{"### Human under ## Positions"}},
		{"open-missing-context", sparring.Open, "## Question\n\nx\n\n## Positions\n\ny\n\n## Disposition\n\nz\n", []string{"## Context"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := sparring.MissingHeadings(tc.a, tc.body)
			if !slices.Equal(got, tc.want) {
				t.Fatalf("MissingHeadings=%v want %v", got, tc.want)
			}
		})
	}
}

func TestSectionBody(t *testing.T) {
	t.Parallel()
	body := strings.Join([]string{
		"## Positions",
		"",
		"### Human",
		"",
		"inside",
		"",
		"### Agent",
		"",
		"also",
		"",
		"## Reasons",
		"",
		"### Human",
		"",
		"outside",
	}, "\n")
	got := sparring.SectionBody(body, "Positions")
	if !strings.Contains(got, "### Human") || !strings.Contains(got, "inside") {
		t.Fatalf("Positions body missing inner H3: %q", got)
	}
	if strings.Contains(got, "outside") || strings.Contains(got, "## Reasons") {
		t.Fatalf("Positions body leaked next H2: %q", got)
	}
	if sparring.SectionBody(body, "positions") != "" {
		t.Fatal("case mismatch should miss")
	}
	if sparring.SectionBody(body, "Missing") != "" {
		t.Fatal("absent H2 should be empty")
	}
	reasons := sparring.SectionBody(body, "Reasons")
	if !strings.Contains(reasons, "outside") {
		t.Fatalf("Reasons body=%q", reasons)
	}
	if strings.Contains(reasons, "inside") {
		t.Fatalf("Reasons body pulled prior section: %q", reasons)
	}
}

func TestHasGlossaryH1(t *testing.T) {
	t.Parallel()
	if !sparring.HasGlossaryH1("# Glossary\n") {
		t.Fatal("H1-only should have glossary H1")
	}
	if !sparring.HasGlossaryH1("preamble\n# Glossary\n## SQLite\n") {
		t.Fatal("H1 after preamble should match")
	}
	if sparring.HasGlossaryH1("# glossary\n") {
		t.Fatal("case mismatch must miss")
	}
	if sparring.HasGlossaryH1("## Glossary\n") {
		t.Fatal("H2 must not count as H1")
	}
	if sparring.HasGlossaryH1("# Glossary Extra\n") {
		t.Fatal("exact line only")
	}
}

func TestMissingGlossaryDefinitions(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		body string
		want []string
	}{
		{"h1-only", "# Glossary\n", nil},
		{"h2-without-definition", "# Glossary\n\n## SQLite\n", []string{"SQLite"}},
		{"h2-with-definition-fill", "# Glossary\n\n## SQLite\n\n### Definition\n\n<!-- fill -->\n", nil},
		{"h2-with-one-word", "# Glossary\n\n## Term\n\n### Definition\n\nx\n", nil},
		{"definition-after-next-h2", "# Glossary\n\n## SQLite\n\n## Postgres\n\n### Definition\n\nok\n", []string{"SQLite"}},
		{"h3-wrong-name", "# Glossary\n\n## SQLite\n\n### Def\n\nok\n", []string{"SQLite"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := sparring.MissingGlossaryDefinitions(tc.body)
			if !slices.Equal(got, tc.want) {
				t.Fatalf("MissingGlossaryDefinitions=%v want %v", got, tc.want)
			}
		})
	}
}

func TestDissentIDs(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		sec  string
		want []string
	}{
		{"empty", "", nil},
		{"dec-only", "Overruled; see DEC-001.\n", nil},
		{"oq", "See OQ-001 for the crux.\n", []string{"OQ-001"}},
		{"asm", "Assumption ASM-002 stands.\n", []string{"ASM-002"}},
		{"both", "OQ-001 and ASM-003 disagree.\n", []string{"OQ-001", "ASM-003"}},
		{"dec-and-oq", "DEC-001 lost; OQ-007 remains.\n", []string{"OQ-007"}},
		{"prose-ok", "I still think this is wrong. Cite OQ-012.\n", []string{"OQ-012"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := sparring.DissentIDs(tc.sec)
			if !slices.Equal(got, tc.want) {
				t.Fatalf("DissentIDs=%v want %v", got, tc.want)
			}
		})
	}
}

func TestHasH2DissentSection(t *testing.T) {
	t.Parallel()
	body := "## Approval\n\nx\n\n## Dissent\n\nOQ-001\n"
	if !sparring.HasH2(body, "Dissent") {
		t.Fatal("want HasH2 Dissent")
	}
	if sparring.HasH2(body, "Crux") {
		t.Fatal("want no Crux H2")
	}
	sec := sparring.SectionBody(body, "Dissent")
	if !slices.Equal(sparring.DissentIDs(sec), []string{"OQ-001"}) {
		t.Fatalf("section extract=%q ids=%v", sec, sparring.DissentIDs(sec))
	}
}
