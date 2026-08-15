// Package logfmt formats and recognizes mycelium log.md lines.
package logfmt

import "strings"

// Line builds one tab-separated log entry:
//
//	YYYY-MM-DD\t<op>\t<ID-or-->\t<title-or-note>
func Line(date, op, idOrDash, note string) string {
	return strings.Join([]string{date, op, idOrDash, note}, "\t")
}

// Parseable reports whether line is a log entry (not a blank or heading).
func Parseable(line string) bool {
	line = strings.TrimRight(line, "\r")
	if line == "" || strings.HasPrefix(line, "#") {
		return false
	}
	parts := strings.SplitN(line, "\t", 4)
	if len(parts) < 4 {
		return false
	}
	if !isDate(parts[0]) || parts[1] == "" || parts[2] == "" {
		return false
	}
	return true
}

// ParseableLines returns parseable log lines from a log.md body, in file order.
func ParseableLines(body []byte) []string {
	raw := strings.Split(string(body), "\n")
	var out []string
	for _, line := range raw {
		if Parseable(line) {
			out = append(out, strings.TrimRight(line, "\r"))
		}
	}
	return out
}

func isDate(s string) bool {
	if len(s) != 10 || s[4] != '-' || s[7] != '-' {
		return false
	}
	for i := 0; i < 10; i++ {
		if i == 4 || i == 7 {
			continue
		}
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
