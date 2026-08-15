// Package logfmt formats and recognizes mycelium log.md lines.
package logfmt

import "strings"

// Line builds one tab-separated log entry:
//
//	YYYY-MM-DD\t<op>\t<ID-or-->\t<title-or-note>
func Line(date, op, idOrDash, note string) string {
	return strings.Join([]string{date, op, idOrDash, note}, "\t")
}
