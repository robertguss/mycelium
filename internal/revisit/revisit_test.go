package revisit_test

import (
	"testing"
	"time"

	"github.com/robertguss/mycelium/internal/revisit"
)

func TestParse(t *testing.T) {
	cases := []struct {
		in      string
		wantK   revisit.Kind
		wantD   string // YYYY-MM-DD when Date
		wantE   string // event token when Event
		wantErr bool
	}{
		{in: "2026-08-08", wantK: revisit.Date, wantD: "2026-08-08"},
		{in: "event:after-iphone-launch", wantK: revisit.Event, wantE: "event:after-iphone-launch"},
		{in: "event:budget-review", wantK: revisit.Event, wantE: "event:budget-review"},
		{in: "event:q4", wantK: revisit.Event, wantE: "event:q4"},
		{in: "", wantErr: true},
		{in: "in two weeks", wantErr: true},
		{in: "2026-08-08T00:00:00Z", wantErr: true},
		{in: "2026-02-30", wantErr: true},
		{in: "2026-8-8", wantErr: true},
		{in: "2026/08/08", wantErr: true},
		{in: "EVENT:foo", wantErr: true},
		{in: "event:After-Launch", wantErr: true},
		{in: "event:", wantErr: true},
		{in: "event:after_iphone", wantErr: true},
		{in: "after-iphone-launch", wantErr: true},
	}
	for _, tc := range cases {
		k, d, e, err := revisit.Parse(tc.in)
		if tc.wantErr {
			if err == nil {
				t.Fatalf("Parse(%q): want error, got kind=%v date=%v event=%q", tc.in, k, d, e)
			}
			continue
		}
		if err != nil {
			t.Fatalf("Parse(%q): %v", tc.in, err)
		}
		if k != tc.wantK {
			t.Fatalf("Parse(%q): kind=%v want %v", tc.in, k, tc.wantK)
		}
		if tc.wantK == revisit.Date {
			if got := d.UTC().Format("2006-01-02"); got != tc.wantD {
				t.Fatalf("Parse(%q): date=%q want %q", tc.in, got, tc.wantD)
			}
			if e != "" {
				t.Fatalf("Parse(%q): event=%q want empty", tc.in, e)
			}
		}
		if tc.wantK == revisit.Event {
			if e != tc.wantE {
				t.Fatalf("Parse(%q): event=%q want %q", tc.in, e, tc.wantE)
			}
			if !d.IsZero() {
				t.Fatalf("Parse(%q): date should be zero", tc.in)
			}
		}
	}
}

func TestDueOverdue(t *testing.T) {
	date := mustDate(t, "2026-08-08")
	cases := []struct {
		name    string
		kind    revisit.Kind
		now     string
		due     bool
		overdue bool
	}{
		{name: "before", kind: revisit.Date, now: "2026-08-07", due: false, overdue: false},
		{name: "on", kind: revisit.Date, now: "2026-08-08", due: true, overdue: false},
		{name: "after", kind: revisit.Date, now: "2026-08-09", due: true, overdue: true},
		{name: "event-any", kind: revisit.Event, now: "2026-08-09", due: false, overdue: false},
		{name: "event-on-date", kind: revisit.Event, now: "2026-08-08", due: false, overdue: false},
	}
	for _, tc := range cases {
		now := mustDate(t, tc.now)
		if got := revisit.Due(tc.kind, date, now); got != tc.due {
			t.Fatalf("%s Due=%v want %v", tc.name, got, tc.due)
		}
		if got := revisit.Overdue(tc.kind, date, now); got != tc.overdue {
			t.Fatalf("%s Overdue=%v want %v", tc.name, got, tc.overdue)
		}
	}
}

func TestExtractTriggerDate(t *testing.T) {
	cases := []struct {
		name string
		body string
		want string // empty → ok=false
	}{
		{name: "bare", body: "2026-08-06\n", want: "2026-08-06"},
		{name: "leading-ws", body: "  2026-08-06\n", want: "2026-08-06"},
		{name: "suffix-text", body: "2026-08-06 something\n", want: "2026-08-06"},
		{name: "iso-datetime", body: "2026-08-06T00:00:00Z\n", want: ""},
		{name: "first-wins", body: "2026-08-05\n2026-08-06\n", want: "2026-08-05"},
		{name: "mid-line-ignored", body: "trigger 2026-08-06\n", want: ""},
		{name: "later-line", body: "none yet\n2026-08-06\n", want: "2026-08-06"},
		{name: "missing", body: "no date here\n", want: ""},
		{name: "empty", body: "", want: ""},
	}
	for _, tc := range cases {
		got, ok := revisit.ExtractTriggerDate(tc.body)
		if tc.want == "" {
			if ok {
				t.Fatalf("%s: want ok=false, got %v", tc.name, got)
			}
			continue
		}
		if !ok {
			t.Fatalf("%s: want %q, got ok=false", tc.name, tc.want)
		}
		if s := got.UTC().Format("2006-01-02"); s != tc.want {
			t.Fatalf("%s: got %q want %q", tc.name, s, tc.want)
		}
	}
}

func mustDate(t *testing.T, s string) time.Time {
	t.Helper()
	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		t.Fatal(err)
	}
	return d.UTC()
}
