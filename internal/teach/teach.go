// Package teach writes four-line teaching errors to stderr.
package teach

import (
	"fmt"
	"io"
)

const maxFindings = 20

// Finding is one teaching-error block.
type Finding struct {
	What       string
	Convention string
	Contract   string
	Fix        string
}

// Write prints the PHASE-01 teaching-error shape and returns exit code 1.
func Write(stderr io.Writer, what, convention, contract, fix string) int {
	return WriteFindings(stderr, []Finding{{
		What: what, Convention: convention, Contract: contract, Fix: fix,
	}})
}

// WriteFindings prints up to 20 teaching errors, then an omit line if more remain.
// Always returns exit code 1.
func WriteFindings(stderr io.Writer, findings []Finding) int {
	n := len(findings)
	if n > maxFindings {
		n = maxFindings
	}
	for i := 0; i < n; i++ {
		f := findings[i]
		fmt.Fprintf(stderr, "mycelium: %s\n", f.What)
		fmt.Fprintf(stderr, "convention: %s\n", f.Convention)
		fmt.Fprintf(stderr, "contract: %s\n", f.Contract)
		fmt.Fprintf(stderr, "fix: %s\n", f.Fix)
	}
	if len(findings) > maxFindings {
		fmt.Fprintln(stderr, "mycelium: further errors omitted")
	}
	return 1
}
