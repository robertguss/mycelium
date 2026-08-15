package clock_test

import (
	"testing"
	"time"

	"github.com/robertguss/mycelium/internal/clock"
)

func TestFixed(t *testing.T) {
	want := time.Date(2026, 8, 14, 15, 4, 5, 0, time.UTC)
	c := clock.Fixed{T: want}
	if got := c.Now(); !got.Equal(want) {
		t.Fatalf("Now=%v want %v", got, want)
	}
	if got := clock.Date(c.Now()); got != "2026-08-14" {
		t.Fatalf("Date=%q", got)
	}
}

func TestSystemRespectsMYCELIUM_NOW(t *testing.T) {
	t.Setenv("MYCELIUM_NOW", "2026-01-02T03:04:05Z")
	got := clock.System().Now()
	want := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	if !got.Equal(want) {
		t.Fatalf("Now=%v want %v", got, want)
	}
}

func TestSystemInvalidMYCELIUM_NOWFallsBack(t *testing.T) {
	t.Setenv("MYCELIUM_NOW", "not-rfc3339")
	before := time.Now().UTC().Add(-time.Second)
	got := clock.System().Now()
	after := time.Now().UTC().Add(time.Second)
	if got.Before(before) || got.After(after) {
		t.Fatalf("fallback Now=%v out of range", got)
	}
}
