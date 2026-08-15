// Package embed holds the go:embed copy of program/ for scaffold emit.
package embed

import stdembed "embed"

// Program is the embedded methodology tree (contracts, templates, tiers).
// Slice 1 ships a stub; Slice 2 authors real 2.0 content.
//
//go:embed all:program
var Program stdembed.FS
