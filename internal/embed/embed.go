// Package embed holds the go:embed copy of program/ for scaffold emit.
package embed

import stdembed "embed"

//go:generate go run gen_copy.go

// Program is the embedded methodology tree (contracts, templates, tiers).
// Authoritative browsable tree is repo-root program/; this package holds a
// generate-copied snapshot (*.go filtered out).
//
//go:embed all:program
var Program stdembed.FS
