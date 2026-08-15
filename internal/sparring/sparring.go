package sparring

import (
	"fmt"
	"strings"
)

type Agreement string

const (
	Open            Agreement = "open"
	Aligned         Agreement = "aligned"
	AgreeToDisagree Agreement = "agree-to-disagree"
)

func ParseAgreement(s string) (Agreement, error) {
	switch strings.TrimSpace(s) {
	case string(Open):
		return Open, nil
	case string(Aligned):
		return Aligned, nil
	case string(AgreeToDisagree):
		return AgreeToDisagree, nil
	default:
		return "", fmt.Errorf("agreement %q is not open|aligned|agree-to-disagree", s)
	}
}

func RequiredH2(a Agreement) []string {
	out := []string{"Question", "Context", "Positions", "Disposition"}
	if a == AgreeToDisagree {
		out = append(out, "Reasons", "Crux")
	}
	return out
}

func RequiredH3(a Agreement, h2 string) []string {
	if a != AgreeToDisagree {
		return nil
	}
	switch h2 {
	case "Positions", "Reasons", "Crux":
		return []string{"Human", "Agent"}
	default:
		return nil
	}
}

func MissingHeadings(a Agreement, body string) []string {
	var missing []string
	for _, h2 := range RequiredH2(a) {
		sec := SectionBody(body, h2)
		if !hasExactH2(body, h2) {
			missing = append(missing, "## "+h2)
			continue
		}
		for _, h3 := range RequiredH3(a, h2) {
			if !hasExactH3(sec, h3) {
				missing = append(missing, "### "+h3+" under ## "+h2)
			}
		}
	}
	return missing
}

func SectionBody(body, h2 string) string {
	want := "## " + h2
	lines := strings.Split(body, "\n")
	start := -1
	for i, line := range lines {
		trim := strings.TrimRight(line, "\r")
		if trim == want {
			start = i + 1
			break
		}
	}
	if start < 0 {
		return ""
	}
	var b strings.Builder
	for i := start; i < len(lines); i++ {
		trim := strings.TrimRight(lines[i], "\r")
		if strings.HasPrefix(trim, "## ") {
			break
		}
		b.WriteString(lines[i])
		if i+1 < len(lines) {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func hasExactH2(body, h2 string) bool {
	want := "## " + h2
	for _, line := range strings.Split(body, "\n") {
		if strings.TrimRight(line, "\r") == want {
			return true
		}
	}
	return false
}

func hasExactH3(section, h3 string) bool {
	want := "### " + h3
	for _, line := range strings.Split(section, "\n") {
		if strings.TrimRight(line, "\r") == want {
			return true
		}
	}
	return false
}

func HasGlossaryH1(content string) bool {
	for _, line := range strings.Split(content, "\n") {
		if strings.TrimRight(line, "\r") == "# Glossary" {
			return true
		}
	}
	return false
}

func MissingGlossaryDefinitions(content string) []string {
	var missing []string
	for _, term := range glossaryH2Terms(content) {
		sec := SectionBody(content, term)
		if !hasExactH3(sec, "Definition") {
			missing = append(missing, term)
		}
	}
	return missing
}

func glossaryH2Terms(content string) []string {
	var terms []string
	for _, line := range strings.Split(content, "\n") {
		trim := strings.TrimRight(line, "\r")
		if !strings.HasPrefix(trim, "## ") {
			continue
		}
		terms = append(terms, strings.TrimPrefix(trim, "## "))
	}
	return terms
}
