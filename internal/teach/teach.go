// Package teach writes four-line teaching errors to stderr.
package teach

import (
	"fmt"
	"io"
)

// Write prints the PHASE-01 teaching-error shape and returns exit code 1.
func Write(stderr io.Writer, what, convention, contract, fix string) int {
	fmt.Fprintf(stderr, "mycelium: %s\n", what)
	fmt.Fprintf(stderr, "convention: %s\n", convention)
	fmt.Fprintf(stderr, "contract: %s\n", contract)
	fmt.Fprintf(stderr, "fix: %s\n", fix)
	return 1
}
