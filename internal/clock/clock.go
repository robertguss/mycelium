// Package clock provides injectable UTC time for artifacts and logs.
package clock

import (
	"os"
	"time"
)

// Clock returns the current time (always use UTC for dates).
type Clock interface {
	Now() time.Time
}

// Fixed returns now every call. For tests.
type Fixed struct {
	T time.Time
}

// Now implements Clock.
func (f Fixed) Now() time.Time {
	return f.T.UTC()
}

type system struct{}

// System returns a clock that uses time.Now().UTC(),
// overridden by MYCELIUM_NOW (RFC3339) when set and valid.
func System() Clock {
	return system{}
}

// Now implements Clock.
func (system) Now() time.Time {
	if s := os.Getenv("MYCELIUM_NOW"); s != "" {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			return t.UTC()
		}
	}
	return time.Now().UTC()
}

// Date returns YYYY-MM-DD for t in UTC.
func Date(t time.Time) string {
	return t.UTC().Format("2006-01-02")
}
