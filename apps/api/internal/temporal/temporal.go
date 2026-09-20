// Package temporal is the shared time codec for the VP-040 time-column
// contract (workspace-040, R2 / GOAL-003).
//
// Authority for every rule below:
//
//   - GOAL-002/attachments/r1-c2-column-contract-draft-v0.1.md §1 "Canonical
//     value contract" and §3 "Codec boundary";
//   - GOAL-002/attachments/r1-c2-c3-guardrails-v0.1.md §1;
//   - Root D-008 (precision / zero policy) and Root D-015 (new writes use
//     UTC().Truncate(time.Microsecond); milliseconds family corrected after
//     measurement on PostgreSQL 15/16/17);
//   - GOAL-002/attachments/r1-public-wire-inventory-v0.1.md L17 (input accepts
//     legal RFC3339 with 0/3/6/9 fractional digits and equivalent offsets;
//     output is fixed 6 digits);
//   - GOAL-002 D-018 (the row-copy SQL is cross-checked against this codec, so
//     FromUnix / FromUnixMilli must not deviate from time.Unix / time.UnixMilli);
//   - GOAL-003 D-001 (this package's API shape).
//
// Storage contract this package implements:
//
//	domain value   = time.Time, semantic UTC instant
//	PostgreSQL     = timestamptz(6)
//	SQLite         = TEXT, exactly "YYYY-MM-DDTHH:MM:SS.ffffffZ" (27 chars)
//
// The package depends on the standard library only. It must never import a
// driver package, database/sql, or any driver-specific time type: the domain
// surface exposes time.Time and nothing else.
package temporal

import (
	"errors"
	"fmt"
	"time"
)

// Layout is the canonical stored form: fixed width, UTC, exactly six
// fractional digits, "Z" suffix. 27 characters.
const Layout = "2006-01-02T15:04:05.000000Z"

// CanonicalLen is the character length of a canonical value. Anything else is
// not a canonical stored value.
const CanonicalLen = len("2006-01-02T15:04:05.000000Z")

// Sentinel errors. Callers wrap these with table/column/row context
// (fmt.Errorf("...: %w", err)); migrations fail closed on them rather than
// converting silently.
var (
	// ErrRange reports a value outside the range the canonical form can
	// represent (the four-digit year field limits us to years 0..9999).
	ErrRange = errors.New("temporal: value out of supported range")
	// ErrFormat reports text that is not canonical fixed-6 UTC RFC3339.
	ErrFormat = errors.New("temporal: value is not canonical fixed-6 UTC RFC3339")
	// ErrNoZone reports a timestamp with no timezone offset.
	ErrNoZone = errors.New("temporal: timestamp has no timezone offset")
	// ErrFraction reports an unsupported fractional-second precision.
	ErrFraction = errors.New("temporal: unsupported fractional-second precision")
)

// minYear and maxYear bound the four-digit year field of Layout.
const (
	minYear = 0
	maxYear = 9999
)

// FromUnix converts legacy integer seconds to the domain value.
//
// The semantics are exactly time.Unix(sec, 0).UTC(). It deliberately does not
// round or clamp: GOAL-002 D-018 cross-checks the migration SQL against this
// function, so any deviation would silently corrupt that equivalence.
func FromUnix(sec int64) time.Time {
	return time.Unix(sec, 0).UTC()
}

// FromUnixMilli converts legacy integer milliseconds to the domain value.
//
// The semantics are exactly time.UnixMilli(ms).UTC(), which floors seconds and
// keeps a non-negative millisecond remainder; negative inputs therefore convert
// to the preceding second (for example -1 becomes 1969-12-31T23:59:59.999Z).
// A negative epoch is a legal instant, not a sentinel (Root D-015).
func FromUnixMilli(ms int64) time.Time {
	return time.UnixMilli(ms).UTC()
}

// Truncate normalizes a value for writing: UTC, truncated toward zero to
// microsecond precision. Root D-008 / D-015 require every new domain write to
// pass through this before binding.
func Truncate(t time.Time) time.Time {
	return t.UTC().Truncate(time.Microsecond)
}

// Format renders t in the canonical stored form.
//
// Format is strict about the output shape (UTC, six digits, Z) and returns
// ErrRange when the instant cannot be represented by the four-digit year field.
//
// Callers on the write path must pass a Truncate'd (or Value-wrapped) instant.
// Format does not itself truncate: formatting a raw time.Time with sub-
// microsecond precision would let time.Format round to six digits, which is not
// the contract's truncate-toward-zero rule.
func Format(t time.Time) (string, error) {
	u := t.UTC()
	if y := u.Year(); y < minYear || y > maxYear {
		return "", fmt.Errorf("%w: year %d outside %d..%d", ErrRange, y, minYear, maxYear)
	}
	return u.Format(Layout), nil
}

// MustFormat is Format for call sites that have already validated the range
// (for example after Truncate on a value that came from the database). It
// panics on an out-of-range instant, so it must not be used on unvalidated
// external input.
func MustFormat(t time.Time) string {
	s, err := Format(t)
	if err != nil {
		panic(err)
	}
	return s
}

// Parse reads a legal RFC3339 timestamp and normalizes it to a UTC domain value.
//
// The input surface is deliberately wider than the output surface (inbound
// tolerance, canonical outbound):
//
//   - fractional seconds: 0, 3, 6 or 9 digits;
//   - offsets: "Z" and any "±HH:MM" form, normalized to UTC;
//   - a missing offset, or a bare local timestamp, is rejected with ErrNoZone
//     rather than guessed.
//
// Rejecting non-zero offsets is a transport-layer policy that belongs to the
// API input parser (R3, GOAL-003 D-001 §3), not to this package.
func Parse(s string) (time.Time, error) {
	if !hasZone(s) {
		return time.Time{}, fmt.Errorf("%w: %q", ErrNoZone, s)
	}
	if err := checkFraction(s); err != nil {
		return time.Time{}, err
	}
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %q: %v", ErrFormat, s, err)
	}
	return t.UTC(), nil
}

// hasZone reports whether s ends with a timezone designator: "Z", "z", or a
// "±HH:MM" offset. Anything else is treated as zoneless.
func hasZone(s string) bool {
	if len(s) == 0 {
		return false
	}
	switch s[len(s)-1] {
	case 'Z', 'z':
		return true
	}
	// Look for a '+' or '-' after the time component (position 10 onwards), so
	// that date separators are not mistaken for an offset sign.
	for i := 10; i < len(s); i++ {
		if s[i] == '+' || s[i] == '-' {
			return true
		}
	}
	return false
}

// checkFraction validates that the fractional-second part, if present, has 0, 3,
// 6 or 9 digits.
func checkFraction(s string) error {
	dot := -1
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			dot = i
			break
		}
	}
	if dot < 0 {
		return nil
	}
	digits := 0
	for i := dot + 1; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			break
		}
		digits++
	}
	switch digits {
	case 0, 3, 6, 9:
		return nil
	default:
		return fmt.Errorf("%w: %d fractional digits in %q (want 0, 3, 6 or 9)", ErrFraction, digits, s)
	}
}

// Value is a non-nullable instant that has already been normalized for storage.
type Value struct {
	t time.Time
}

// NewValue normalizes t for writing (UTC, microsecond-truncated toward zero).
func NewValue(t time.Time) Value {
	return Value{t: Truncate(t)}
}

// Time returns the normalized instant.
func (v Value) Time() time.Time { return v.t }

// String returns the canonical stored form. It never fails for values built by
// NewValue; a zero Value (the Go zero time) formats as its year 1 instant.
func (v Value) String() string { return MustFormat(v.t) }

// NullValue is a nullable instant, symmetric with sql.NullTime but carrying the
// contract's canonical formatting.
type NullValue struct {
	t     time.Time
	valid bool
}

// NewNullValue normalizes t for writing and marks the value present.
func NewNullValue(t time.Time) NullValue {
	return NullValue{t: Truncate(t), valid: true}
}

// NullValueFromValid mirrors sql.NullTime: absent values stay absent.
func NullValueFromValid(t time.Time, valid bool) NullValue {
	if !valid {
		return NullValue{}
	}
	return NewNullValue(t)
}

// Valid reports whether the value is present.
func (n NullValue) Valid() bool { return n.valid }

// Time returns the normalized instant, or the zero time when absent.
func (n NullValue) Time() time.Time { return n.t }

// String returns the canonical form, or "NULL" when absent.
func (n NullValue) String() string {
	if !n.valid {
		return "NULL"
	}
	return MustFormat(n.t)
}
