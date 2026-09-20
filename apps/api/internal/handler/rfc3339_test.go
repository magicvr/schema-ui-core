package handler

import (
	"testing"
	"time"
)

// TestFormatWireTime pins the public output contract (Root D-003/D-008): fixed
// six fractional digits, UTC, "Z" only, and truncation toward zero — never
// rounding — so a nanosecond-bearing instant cannot leak a rounded microsecond.
func TestFormatWireTime(t *testing.T) {
	cases := []struct {
		name  string
		input time.Time
		want  string
	}{
		{"whole second", time.Date(2026, 8, 17, 12, 0, 0, 0, time.UTC), "2026-08-17T12:00:00.000000Z"},
		{"millisecond", time.Date(2026, 8, 17, 12, 0, 0, 123_000_000, time.UTC), "2026-08-17T12:00:00.123000Z"},
		{
			"truncates sub-microsecond digits instead of rounding",
			time.Date(2025, 9, 19, 22, 13, 20, 123_456_789, time.UTC),
			"2025-09-19T22:13:20.123456Z",
		},
		{
			"non-UTC input is normalised",
			time.Date(2026, 8, 17, 20, 0, 0, 0, time.FixedZone("CST", 8*3600)),
			"2026-08-17T12:00:00.000000Z",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := FormatWireTime(tc.input)
			if got != tc.want {
				t.Fatalf("FormatWireTime(%s) = %q, want %q", tc.input, got, tc.want)
			}
			if len(got) != 27 {
				t.Fatalf("wire output %q is not 27 characters", got)
			}
		})
	}
}

// TestFormatWireTimePtr keeps absent instants absent instead of fabricating one.
func TestFormatWireTimePtr(t *testing.T) {
	if got, ok := FormatWireTimePtr(nil); ok || got != "" {
		t.Fatalf("nil pointer = (%q, %v), want (\"\", false)", got, ok)
	}
	instant := time.Date(2026, 8, 17, 12, 0, 0, 0, time.UTC)
	got, ok := FormatWireTimePtr(&instant)
	if !ok || got != "2026-08-17T12:00:00.000000Z" {
		t.Fatalf("pointer = (%q, %v)", got, ok)
	}
}

// TestParseWireTimeCompatibility is the D-005 input matrix: 0/3/6/9 fractional
// digits and equivalent offsets are accepted and normalised, while a zoneless
// value is rejected rather than guessed.
func TestParseWireTimeCompatibility(t *testing.T) {
	accepted := []struct {
		in   string
		want string
	}{
		{"2026-08-17T12:00:00Z", "2026-08-17T12:00:00.000000Z"},
		{"2026-08-17T12:00:00.123Z", "2026-08-17T12:00:00.123000Z"},
		{"2026-08-17T12:00:00.123456Z", "2026-08-17T12:00:00.123456Z"},
		{"2026-08-17T12:00:00.123456789Z", "2026-08-17T12:00:00.123456Z"},
		{"2026-08-17T12:00:00+00:00", "2026-08-17T12:00:00.000000Z"},
		{"2026-08-17T20:00:00+08:00", "2026-08-17T12:00:00.000000Z"},
	}
	for _, tc := range accepted {
		parsed, err := ParseWireTime(tc.in)
		if err != nil {
			t.Fatalf("ParseWireTime(%q) errored: %v", tc.in, err)
		}
		if got := FormatWireTime(parsed); got != tc.want {
			t.Fatalf("ParseWireTime(%q) -> %q, want %q", tc.in, got, tc.want)
		}
	}
	rejected := []string{
		"2026-08-17T12:00:00",     // no zone
		"2026-08-17 12:00:00Z",    // space separator
		"2026-08-17T12:00:00.12Z", // unsupported 2-digit fraction
		"not-a-time",
		"",
	}
	for _, in := range rejected {
		if _, err := ParseWireTime(in); err == nil {
			t.Fatalf("ParseWireTime(%q) must be rejected", in)
		}
	}
}
