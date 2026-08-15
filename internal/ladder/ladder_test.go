package ladder_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/ladder"
)

func TestParseRung(t *testing.T) {
	cases := []struct {
		in      string
		want    ladder.Rung
		wantErr bool
	}{
		{"second-opinion", ladder.RungSecondOpinion, false},
		{"council", ladder.RungCouncil, false},
		{"", "", true},
		{"Second-Opinion", "", true},
		{"second-opinion ", "", true},
		{" second-opinion", "", true},
		{"COUNCIL", "", true},
	}
	for _, tc := range cases {
		got, err := ladder.ParseRung(tc.in)
		if tc.wantErr {
			if err == nil {
				t.Fatalf("ParseRung(%q) err=nil want error", tc.in)
			}
			continue
		}
		if err != nil {
			t.Fatalf("ParseRung(%q): %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("ParseRung(%q)=%q want %q", tc.in, got, tc.want)
		}
	}
}

func TestParseAdapter(t *testing.T) {
	cases := []struct {
		in      string
		want    ladder.Adapter
		wantErr bool
	}{
		{"cursor", ladder.AdapterCursor, false},
		{"manual", ladder.AdapterManual, false},
		{"", "", true},
		{"Cursor", "", true},
		{"manual ", "", true},
	}
	for _, tc := range cases {
		got, err := ladder.ParseAdapter(tc.in)
		if tc.wantErr {
			if err == nil {
				t.Fatalf("ParseAdapter(%q) err=nil want error", tc.in)
			}
			continue
		}
		if err != nil {
			t.Fatalf("ParseAdapter(%q): %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("ParseAdapter(%q)=%q want %q", tc.in, got, tc.want)
		}
	}
}

func TestParseCostClass(t *testing.T) {
	cases := []struct {
		in      string
		want    ladder.CostClass
		wantErr bool
	}{
		{"cheap", ladder.CostCheap, false},
		{"quick", ladder.CostQuick, false},
		{"standard", ladder.CostStandard, false},
		{"high-stakes", ladder.CostHighStakes, false},
		{"", "", true},
		{"Cheap", "", true},
		{"standard ", "", true},
		{"high_stakes", "", true},
	}
	for _, tc := range cases {
		got, err := ladder.ParseCostClass(tc.in)
		if tc.wantErr {
			if err == nil {
				t.Fatalf("ParseCostClass(%q) err=nil want error", tc.in)
			}
			continue
		}
		if err != nil {
			t.Fatalf("ParseCostClass(%q): %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("ParseCostClass(%q)=%q want %q", tc.in, got, tc.want)
		}
	}
}

func TestCostClassOK(t *testing.T) {
	cases := []struct {
		rung  ladder.Rung
		class ladder.CostClass
		ok    bool
	}{
		{ladder.RungSecondOpinion, ladder.CostCheap, true},
		{ladder.RungSecondOpinion, ladder.CostQuick, false},
		{ladder.RungSecondOpinion, ladder.CostStandard, false},
		{ladder.RungSecondOpinion, ladder.CostHighStakes, false},
		{ladder.RungCouncil, ladder.CostCheap, false},
		{ladder.RungCouncil, ladder.CostQuick, true},
		{ladder.RungCouncil, ladder.CostStandard, true},
		{ladder.RungCouncil, ladder.CostHighStakes, true},
		{ladder.Rung(""), ladder.CostCheap, false},
	}
	for _, tc := range cases {
		if got := ladder.CostClassOK(tc.rung, tc.class); got != tc.ok {
			t.Fatalf("CostClassOK(%q,%q)=%v want %v", tc.rung, tc.class, got, tc.ok)
		}
	}
}

func TestOptInOK(t *testing.T) {
	cases := []struct {
		name string
		v    any
		ok   bool
	}{
		{"true", true, true},
		{"false", false, false},
		{"string true", "true", false},
		{"missing nil", nil, false},
		{"int 1", 1, false},
	}
	for _, tc := range cases {
		if got := ladder.OptInOK(tc.v); got != tc.ok {
			t.Fatalf("OptInOK(%v)=%v want %v", tc.v, got, tc.ok)
		}
	}
}

func TestPromptSHA256(t *testing.T) {
	const appendixB = "Should this idea use SQLite as the store? Answer independently. Do not see other reports."
	const wantB = "ec87bfc2afd545807ca87b5c29cae8e77262cb3c746fc63e4539d8daeb2a77de"
	const emptyHash = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	body := "## Prompt\n\n" + appendixB + "\n\n## Attachments\n\nnone\n"
	if got := ladder.PromptSHA256(body); got != wantB {
		t.Fatalf("appendix B hash=%q want %q", got, wantB)
	}

	padded := "## Prompt\n\n  \n" + appendixB + "\n\n  \n## Attachments\n\nnone\n"
	if got := ladder.PromptSHA256(padded); got != wantB {
		t.Fatalf("whitespace pad hash=%q want %q", got, wantB)
	}

	empty := "## Prompt\n\n## Attachments\n\nnone\n"
	if got := ladder.PromptSHA256(empty); got != emptyHash {
		t.Fatalf("empty prompt hash=%q want %q", got, emptyHash)
	}

	firstWins := "## Prompt\n\n" + appendixB + "\n\n## Prompt\n\nother\n\n## Attachments\n\nnone\n"
	if got := ladder.PromptSHA256(firstWins); got != wantB {
		t.Fatalf("first H2 wins hash=%q want %q", got, wantB)
	}

	got := ladder.PromptSHA256(body)
	if strings.ToUpper(got) == got && got != "" {
		t.Fatal("hash must be lowercase")
	}
	if len(got) != 64 {
		t.Fatalf("hash len=%d want 64", len(got))
	}
	if strings.HasPrefix(got, "sha256:") {
		t.Fatal("must not prefix sha256:")
	}
}

func TestCardinality(t *testing.T) {
	cases := []struct {
		rung    ladder.Rung
		nRPT    int
		nRCL    int
		wantErr bool
	}{
		{ladder.RungSecondOpinion, 0, 0, false},
		{ladder.RungSecondOpinion, 1, 0, false},
		{ladder.RungSecondOpinion, 2, 0, true},
		{ladder.RungSecondOpinion, 1, 1, true},
		{ladder.RungSecondOpinion, 0, 1, true},
		{ladder.RungCouncil, 0, 0, false},
		{ladder.RungCouncil, 1, 0, true},
		{ladder.RungCouncil, 2, 0, true},
		{ladder.RungCouncil, 1, 1, true},
		{ladder.RungCouncil, 2, 1, false},
		{ladder.RungCouncil, 3, 1, false},
		{ladder.RungCouncil, 2, 2, true},
		{ladder.RungCouncil, 0, 1, true},
		{ladder.Rung("nope"), 0, 0, true},
	}
	for _, tc := range cases {
		err := ladder.Cardinality(tc.rung, tc.nRPT, tc.nRCL)
		if tc.wantErr && err == nil {
			t.Fatalf("Cardinality(%q,%d,%d) err=nil want error", tc.rung, tc.nRPT, tc.nRCL)
		}
		if !tc.wantErr && err != nil {
			t.Fatalf("Cardinality(%q,%d,%d): %v", tc.rung, tc.nRPT, tc.nRCL, err)
		}
		if tc.wantErr && err != nil && !errors.Is(err, ladder.ErrCard) {
			t.Fatalf("Cardinality error should wrap ErrCard: %v", err)
		}
	}
}

func TestSeedDissentOK(t *testing.T) {
	if err := ladder.SeedDissentOK([]string{"none", "fine"}, "none"); err != nil {
		t.Fatalf("no token: %v", err)
	}
	if err := ladder.SeedDissentOK([]string{"SEED-DISSENT"}, "kept SEED-DISSENT here"); err != nil {
		t.Fatalf("token retained: %v", err)
	}
	err := ladder.SeedDissentOK([]string{"SEED-DISSENT"}, "none")
	if !errors.Is(err, ladder.ErrSeedDissent) {
		t.Fatalf("want ErrSeedDissent, got %v", err)
	}
	if err := ladder.SeedDissentOK(nil, ""); err != nil {
		t.Fatalf("empty: %v", err)
	}
}

func TestSectionBody(t *testing.T) {
	body := "## Prompt\n\nfirst\n\n## Prompt\n\nsecond\n\n## Attachments\n\nnone\n"
	got := ladder.SectionBody(body, "Prompt")
	if !strings.Contains(got, "first") || strings.Contains(got, "second") {
		t.Fatalf("first H2 wins: %q", got)
	}
	if ladder.SectionBody(body, "prompt") != "" {
		t.Fatal("H2 match is case-sensitive")
	}
	if ladder.SectionBody(body, "Missing") != "" {
		t.Fatal("missing H2 should be empty")
	}
	att := ladder.SectionBody(body, "Attachments")
	if !strings.Contains(att, "none") {
		t.Fatalf("Attachments=%q", att)
	}
}

func TestRequiredH2s(t *testing.T) {
	if got := ladder.RequiredCMPH2(); len(got) != 3 || got[0] != "Prompt" {
		t.Fatalf("CMP H2s=%v", got)
	}
	if got := ladder.RequiredRPTH2(); len(got) != 3 || got[2] != "Dissent" {
		t.Fatalf("RPT H2s=%v", got)
	}
	rcl := ladder.RequiredRCLH2()
	if len(rcl) != 10 {
		t.Fatalf("RCL H2 count=%d want 10", len(rcl))
	}
	if rcl[len(rcl)-1] != "Retained dissent" {
		t.Fatalf("last RCL H2=%q", rcl[len(rcl)-1])
	}
}
