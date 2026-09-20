package temporal

import (
	"errors"
	"testing"
	"time"
)

// TestFromUnixMatchesGo is the GOAL-002 D-018 cross-check baseline: the
// migration SQL is validated against these functions, so they must equal
// time.Unix exactly.
func TestFromUnixMatchesGo(t *testing.T) {
	cases := []int64{0, 1, -1, 1758320000, -86400, 253402300799, 2147483647}
	for _, sec := range cases {
		got := FromUnix(sec)
		want := time.Unix(sec, 0).UTC()
		if !got.Equal(want) {
			t.Errorf("FromUnix(%d) = %v, want %v", sec, got, want)
		}
		if got.Location() != time.UTC {
			t.Errorf("FromUnix(%d) location = %v, want UTC", sec, got.Location())
		}
	}
}

// TestFromUnixMilliFloorSemantics pins the negative-millisecond behaviour that
// the SQL integer-split expression is required to reproduce.
func TestFromUnixMilliFloorSemantics(t *testing.T) {
	cases := []struct {
		ms   int64
		want string
	}{
		{1758320000123, "2025-09-19T22:13:20.123000Z"},
		{1758320000999, "2025-09-19T22:13:20.999000Z"},
		{0, "1970-01-01T00:00:00.000000Z"},
		{-1, "1969-12-31T23:59:59.999000Z"},
		{-999, "1969-12-31T23:59:59.001000Z"},
		{-1000, "1969-12-31T23:59:59.000000Z"},
		{-1001, "1969-12-31T23:59:58.999000Z"},
		{-86400000, "1969-12-31T00:00:00.000000Z"},
		{-1758320000123, "1914-04-14T01:46:39.877000Z"},
		{253402300799999, "9999-12-31T23:59:59.999000Z"},
	}
	for _, c := range cases {
		got := MustFormat(FromUnixMilli(c.ms))
		if got != c.want {
			t.Errorf("FromUnixMilli(%d) = %s, want %s", c.ms, got, c.want)
		}
		// Must agree with the standard library, not merely with the literals.
		if want := time.UnixMilli(c.ms).UTC().Format(Layout); got != want {
			t.Errorf("FromUnixMilli(%d) diverges from time.UnixMilli: %s vs %s", c.ms, got, want)
		}
	}
}

// TestFormatCanonicalShape checks the fixed-width canonical output.
func TestFormatCanonicalShape(t *testing.T) {
	cases := []struct {
		sec  int64
		want string
	}{
		{0, "1970-01-01T00:00:00.000000Z"},
		{1758320000, "2025-09-19T22:13:20.000000Z"},
		{-1, "1969-12-31T23:59:59.000000Z"},
		{253402300799, "9999-12-31T23:59:59.000000Z"},
	}
	for _, c := range cases {
		got := MustFormat(FromUnix(c.sec))
		if got != c.want {
			t.Errorf("FromUnix(%d) = %s, want %s", c.sec, got, c.want)
		}
		if len(got) != CanonicalLen {
			t.Errorf("FromUnix(%d) length = %d, want %d (%q)", c.sec, len(got), CanonicalLen, got)
		}
	}
}

// TestFormatRejectsUnrepresentableYear covers the four-digit year limit.
func TestFormatRejectsUnrepresentableYear(t *testing.T) {
	if _, err := Format(time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC)); !errors.Is(err, ErrRange) {
		t.Errorf("year 10000: err = %v, want ErrRange", err)
	}
	if _, err := Format(time.Date(-1, 1, 1, 0, 0, 0, 0, time.UTC)); !errors.Is(err, ErrRange) {
		t.Errorf("year -1: err = %v, want ErrRange", err)
	}
	if _, err := Format(time.Date(9999, 12, 31, 23, 59, 59, 999000000, time.UTC)); err != nil {
		t.Errorf("year 9999: unexpected err %v", err)
	}
}

// TestTruncateTowardZero pins the write-path rule, including pre-epoch instants.
func TestTruncateTowardZero(t *testing.T) {
	cases := []struct {
		in   time.Time
		want string
	}{
		{time.Date(2025, 9, 19, 22, 13, 20, 123456789, time.UTC), "2025-09-19T22:13:20.123456Z"},
		// Negative instants must truncate toward zero, not toward -inf.
		{time.Date(1969, 12, 31, 23, 59, 59, 999999999, time.UTC), "1969-12-31T23:59:59.999999Z"},
		{time.Date(1969, 12, 31, 23, 59, 59, 123456789, time.UTC), "1969-12-31T23:59:59.123456Z"},
	}
	for _, c := range cases {
		if got := MustFormat(Truncate(c.in)); got != c.want {
			t.Errorf("Truncate(%v) = %s, want %s", c.in, got, c.want)
		}
	}
}

// TestTruncateNormalizesZone checks that a non-UTC location still stores UTC.
func TestTruncateNormalizesZone(t *testing.T) {
	loc := time.FixedZone("UTC+8", 8*3600)
	in := time.Date(2025, 9, 20, 6, 13, 20, 0, loc) // == 2025-09-19T22:13:20Z
	if got := MustFormat(Truncate(in)); got != "2025-09-19T22:13:20.000000Z" {
		t.Errorf("Truncate(UTC+8 instant) = %s, want 2025-09-19T22:13:20.000000Z", got)
	}
}

// TestParseAcceptsCompatibleInput covers the inbound compatibility surface:
// 0/3/6/9 fractional digits and equivalent offsets. Each case states the
// instant it actually denotes. The four fractional widths denote *different*
// instants, so the assertion is per case rather than against one shared string.
func TestParseAcceptsCompatibleInput(t *testing.T) {
	cases := []struct {
		in   string
		want time.Time
	}{
		{in: "2025-09-19T22:13:20Z", want: time.Date(2025, 9, 19, 22, 13, 20, 0, time.UTC)},
		{in: "2025-09-19T22:13:20.123Z", want: time.Date(2025, 9, 19, 22, 13, 20, 123000000, time.UTC)},
		{in: "2025-09-19T22:13:20.123000Z", want: time.Date(2025, 9, 19, 22, 13, 20, 123000000, time.UTC)},
		{in: "2025-09-19T22:13:20.123456789Z", want: time.Date(2025, 9, 19, 22, 13, 20, 123456789, time.UTC)},
		{in: "2025-09-19T22:13:20.123+00:00", want: time.Date(2025, 9, 19, 22, 13, 20, 123000000, time.UTC)},
		// +08:00 spelling of 2025-09-19T22:13:20.123Z: same instant, other zone.
		{in: "2025-09-20T06:13:20.123+08:00", want: time.Date(2025, 9, 19, 22, 13, 20, 123000000, time.UTC)},
	}
	for _, c := range cases {
		got, err := Parse(c.in)
		if err != nil {
			t.Errorf("Parse(%q) unexpected err: %v", c.in, err)
			continue
		}
		if got.Location() != time.UTC {
			t.Errorf("Parse(%q) location = %v, want UTC", c.in, got.Location())
		}
		if !got.Equal(c.want) {
			t.Errorf("Parse(%q) = %v, want %v", c.in, got, c.want)
		}
	}

	// The spellings that denote one instant must agree after normalization. This
	// equivalence, not cross-case equality above, is what migrations rely on.
	same := []string{
		"2025-09-19T22:13:20.123Z",
		"2025-09-19T22:13:20.123000Z",
		"2025-09-19T22:13:20.123+00:00",
		"2025-09-20T06:13:20.123+08:00",
	}
	var first string
	for i, in := range same {
		got, err := Parse(in)
		if err != nil {
			t.Fatalf("Parse(%q): %v", in, err)
		}
		s := MustFormat(Truncate(got))
		if i == 0 {
			first = s
			continue
		}
		if s != first {
			t.Errorf("equivalent spellings disagree: %q -> %s vs %s", in, s, first)
		}
	}
}

// TestParseRejects covers the rejection rules.
func TestParseRejects(t *testing.T) {
	t.Run("zoneless", func(t *testing.T) {
		for _, in := range []string{"2025-09-19T22:13:20", "2025-09-19T22:13:20.123", "2025-09-19 22:13:20"} {
			if _, err := Parse(in); !errors.Is(err, ErrNoZone) {
				t.Errorf("Parse(%q) err = %v, want ErrNoZone", in, err)
			}
		}
	})
	t.Run("bad fraction width", func(t *testing.T) {
		for _, in := range []string{"2025-09-19T22:13:20.1Z", "2025-09-19T22:13:20.12Z", "2025-09-19T22:13:20.1234Z", "2025-09-19T22:13:20.12345Z"} {
			if _, err := Parse(in); !errors.Is(err, ErrFraction) {
				t.Errorf("Parse(%q) err = %v, want ErrFraction", in, err)
			}
		}
	})
	t.Run("invalid calendar value", func(t *testing.T) {
		if _, err := Parse("2025-02-30T00:00:00Z"); !errors.Is(err, ErrFormat) {
			t.Errorf("Parse(Feb 30) err = %v, want ErrFormat", err)
		}
	})
	t.Run("empty", func(t *testing.T) {
		if _, err := Parse(""); !errors.Is(err, ErrNoZone) {
			t.Errorf("Parse(\"\") err = %v, want ErrNoZone", err)
		}
	})
}

// TestRoundTripCanonical proves Format/Parse are inverse on canonical values.
func TestRoundTripCanonical(t *testing.T) {
	for _, sec := range []int64{0, 1, -1, 1758320000, -86400, 253402300799} {
		v := NewValue(FromUnix(sec))
		s := v.String()
		back, err := Parse(s)
		if err != nil {
			t.Fatalf("Parse(%s): %v", s, err)
		}
		if !back.Equal(v.Time()) {
			t.Errorf("round trip %d: %v != %v", sec, back, v.Time())
		}
	}
}

// TestValueAndNullValue covers the write-path wrappers.
func TestValueAndNullValue(t *testing.T) {
	in := time.Date(2025, 9, 19, 22, 13, 20, 123456789, time.UTC)
	if got := NewValue(in).String(); got != "2025-09-19T22:13:20.123456Z" {
		t.Errorf("NewValue.String() = %s", got)
	}
	if got := NewNullValue(in).String(); got != "2025-09-19T22:13:20.123456Z" {
		t.Errorf("NewNullValue.String() = %s", got)
	}
	absent := NullValueFromValid(time.Time{}, false)
	if absent.Valid() {
		t.Error("absent NullValue reports Valid()")
	}
	if got := absent.String(); got != "NULL" {
		t.Errorf("absent NullValue.String() = %s, want NULL", got)
	}
	present := NullValueFromValid(in, true)
	if !present.Valid() {
		t.Error("present NullValue reports !Valid()")
	}
}

// TestCanonicalConstants pins the contract constants. The "no driver imports"
// property is not testable from inside the package by inspection; it is verified
// at review time by the dependency query recorded in this subgoal's execution
// ledger:
//
//	go list -f "{{range .Imports}}{{.}} {{end}}" ./internal/temporal/
//	=> errors fmt time
func TestCanonicalConstants(t *testing.T) {
	if CanonicalLen != 27 {
		t.Errorf("CanonicalLen = %d, want 27", CanonicalLen)
	}
	if Layout != "2006-01-02T15:04:05.000000Z" {
		t.Errorf("Layout = %q", Layout)
	}
	if len("2006-01-02T15:04:05.000000Z") != CanonicalLen {
		t.Errorf("Layout literal length %d != CanonicalLen %d", len("2006-01-02T15:04:05.000000Z"), CanonicalLen)
	}
}
