// Package version holds the stamped CLI version string.
package version

// Version is the CLI version. Override at link time:
//
//	-ldflags "-X github.com/robertguss/mycelium/internal/version.Version=…"
var Version = "0.1.0-dev"
