package handler

import (
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/temporal"
)

// TestR3BWireRoundTripIsZoneIndependent is the Go half of the VP-020 round-trip
// matrix (I-040-004, carrier ruled by I-041-007).
//
// It pins the two properties the session/user timezone feature depends on:
//
//  1. The wire value is an instant, not a wall clock: converting the same
//     instant into any display zone — including zones the headless server has
//     never seen — yields the byte-identical canonical string, so a deployment
//     or session timezone can never shift what is stored or transmitted.
//  2. format → parse → format is the identity, and parsing recovers exactly the
//     stored instant at microsecond precision (the display layer's rounding
//     never reaches the persisted value: Root D-008 truncates toward zero).
func TestR3BWireRoundTripIsZoneIndependent(t *testing.T) {
	shanghai := time.FixedZone("CST", 8*3600)
	newYork, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("load America/New_York: %v", err)
	}
	kathmandu, err := time.LoadLocation("Asia/Kathmandu") // +05:45, a non-hour offset
	if err != nil {
		t.Fatalf("load Asia/Kathmandu: %v", err)
	}
	zones := []*time.Location{time.UTC, shanghai, newYork, kathmandu}

	instants := []struct {
		name string
		at   time.Time
	}{
		{"trailing-zero microseconds", time.Date(2026, 9, 20, 12, 57, 15, 900_000_000, time.UTC)},
		{"whole second", time.Date(2026, 9, 20, 12, 57, 15, 0, time.UTC)},
		{"sub-microsecond digits are truncated, not rounded", time.Date(2026, 9, 20, 12, 57, 15, 123_456_789, time.UTC)},
		{"negative epoch is a legal instant", time.Date(1965, 3, 4, 5, 6, 7, 654_321_000, time.UTC)},
		{"dst spring-forward edge (US)", time.Date(2026, 3, 8, 7, 30, 0, 1_000, time.UTC)},
		{"dst fall-back edge (US)", time.Date(2026, 11, 1, 5, 30, 0, 999_999_000, time.UTC)},
	}

	for _, tc := range instants {
		t.Run(tc.name, func(t *testing.T) {
			wire := FormatWireTime(tc.at)
			if !wireInstantPattern.MatchString(wire) {
				t.Fatalf("FormatWireTime(%s) = %q, not the canonical fixed-6 shape", tc.at, wire)
			}

			// (1) Display zone independence.
			for _, zone := range zones {
				if got := FormatWireTime(tc.at.In(zone)); got != wire {
					t.Fatalf("rendering the instant in %s changed the wire value: %q vs %q", zone, got, wire)
				}
			}

			// (2) format → parse → format identity, at microsecond precision.
			parsed, err := ParseWireTime(wire)
			if err != nil {
				t.Fatalf("ParseWireTime(%q) errored: %v", wire, err)
			}
			if !parsed.Equal(temporal.Truncate(tc.at)) {
				t.Fatalf("round trip lost the instant: parsed %s, want %s (truncated toward zero)",
					parsed.UTC(), temporal.Truncate(tc.at).UTC())
			}
			if got := FormatWireTime(parsed); got != wire {
				t.Fatalf("format → parse → format = %q, want the identical %q", got, wire)
			}

			// A display-zone conversion of the parsed instant round-trips too.
			for _, zone := range zones {
				if got := FormatWireTime(parsed.In(zone)); got != wire {
					t.Fatalf("parsed instant displayed in %s = %q, want %q", zone, got, wire)
				}
			}
		})
	}
}

// TestR3BSessionTimezoneDoesNotMoveStoredOrSentInstants pins the negative half of
// the VP-020 contract at the HTTP boundary: the same request answered while the
// process timezone differs still carries the same canonical instant. It uses
// time.Local because that is the only ambient zone a handler could accidentally
// pick up.
func TestR3BSessionTimezoneDoesNotMoveStoredOrSentInstants(t *testing.T) {
	original := time.Local
	t.Cleanup(func() { time.Local = original })

	at := time.Date(2026, 9, 20, 12, 57, 15, 900_000_000, time.UTC)
	for _, name := range []string{"UTC", "America/New_York", "Asia/Kathmandu"} {
		loc, err := time.LoadLocation(name)
		if err != nil {
			t.Fatalf("load %s: %v", name, err)
		}
		time.Local = loc
		if got := FormatWireTime(at); got != "2026-09-20T12:57:15.900000Z" {
			t.Fatalf("with time.Local=%s the wire value = %q, want the canonical UTC instant", name, got)
		}
	}
}
