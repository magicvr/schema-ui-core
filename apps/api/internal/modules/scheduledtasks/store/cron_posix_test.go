package store

import (
	"testing"
	"time"
)

// W9 A-005 R-F-003: POSIX day semantics lock for F-025 — when BOTH dom and dow
// are restricted, the day matches on EITHER (Vixie OR); a full-range field is
// behaviorally "*" and keeps AND.
func TestCronMatchesDomDowPOSIX(t *testing.T) {
	at := func(y int, mo time.Month, d int) time.Time {
		return time.Date(y, mo, d, 0, 0, 0, 0, time.UTC)
	}
	// Calendar sanity for the fixed dates used below.
	if at(2026, 1, 5).Weekday() != time.Monday || at(2026, 1, 12).Weekday() != time.Monday || at(2026, 1, 1).Weekday() != time.Thursday {
		t.Fatal("calendar assumption broken: 2026-01-05/12 Monday, 2026-01-01 Thursday")
	}
	cases := []struct {
		expr  string
		when  time.Time
		match bool
		note  string
	}{
		{"0 0 1 * 1", at(2026, 1, 1), true, "restricted dom hit (OR leg 1)"},
		{"0 0 1 * 1", at(2026, 1, 5), true, "restricted dow hit (OR leg 2)"},
		{"0 0 1 * 1", at(2026, 1, 15), false, "neither hits"},
		{"0 0 1 * *", at(2026, 1, 5), false, "dow unrestricted: AND keeps dom-only"},
		{"0 0 1 * *", at(2026, 1, 1), true, "dom only"},
		{"0 0 * * 1", at(2026, 1, 12), true, "dow only"},
		{"0 0 * * 1", at(2026, 1, 1), false, "dom unrestricted: AND keeps dow-only"},
		{"0 0 1-31 * 1", at(2026, 1, 12), true, "full-range dom behaves as *: dow decides"},
		{"0 0 1-31 * 1", at(2026, 1, 1), false, "full-range dom, not Monday (Jan 1 is a Thursday)"},
		{"0 0 1-31 * 1", at(2026, 1, 15), false, "full-range dom, not Monday"},
	}
	for _, tc := range cases {
		f, err := ParseCron(tc.expr)
		if err != nil {
			t.Fatalf("ParseCron(%q): %v", tc.expr, err)
		}
		if got := f.Matches(tc.when); got != tc.match {
			t.Fatalf("%s at %s (%s): Matches = %v, want %v (%s)", tc.expr, tc.when.Format("2006-01-02"), tc.when.Weekday(), got, tc.match, tc.note)
		}
	}
}
