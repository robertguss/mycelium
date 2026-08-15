// Package revisit parses revisit triggers and due/overdue math (PHASE-02).
package revisit

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

// Kind is the revisit shape.
type Kind int

const (
	Date Kind = iota + 1
	Event
)

const dateLayout = "2006-01-02"

var (
	errEmpty  = errors.New("revisit: empty")
	eventRE   = regexp.MustCompile(`^event:[a-z0-9]+(?:-[a-z0-9]+)*$`)
	triggerRE = regexp.MustCompile(`(?m)^\s*(\d{4}-\d{2}-\d{2})\b`)
)

// Parse validates revisit grammar.
// Date returns UTC midnight; Event returns a zero time and the event token.
func Parse(s string) (Kind, time.Time, string, error) {
	if s == "" {
		return 0, time.Time{}, "", errEmpty
	}
	if t, err := time.Parse(dateLayout, s); err == nil {
		return Date, t.UTC(), "", nil
	}
	if eventRE.MatchString(s) {
		return Event, time.Time{}, s, nil
	}
	return 0, time.Time{}, "", fmt.Errorf("revisit: %q is not a date or event:<kebab>", s)
}

// Due is true for date shape when now's UTC date is on or after date.
// Event is never due via this helper.
func Due(kind Kind, date, now time.Time) bool {
	if kind != Date {
		return false
	}
	return !utcDate(now).Before(utcDate(date))
}

// Overdue is true for date shape when now's UTC date is strictly after date.
// Event is never overdue.
func Overdue(kind Kind, date, now time.Time) bool {
	if kind != Date {
		return false
	}
	return utcDate(now).After(utcDate(date))
}

// ExtractTriggerDate returns the first line-leading YYYY-MM-DD in sectionBody.
func ExtractTriggerDate(sectionBody string) (time.Time, bool) {
	m := triggerRE.FindStringSubmatch(sectionBody)
	if m == nil {
		return time.Time{}, false
	}
	t, err := time.Parse(dateLayout, m[1])
	if err != nil {
		return time.Time{}, false
	}
	return t.UTC(), true
}

func utcDate(t time.Time) time.Time {
	y, m, d := t.UTC().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
