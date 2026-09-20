package handler

import (
	"errors"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/temporal"
)

// Public wire time contract (VP-040 R3; Root D-003 / D-005 / D-009).
//
// Output: exactly "YYYY-MM-DDTHH:MM:SS.ffffffZ" — UTC, six fractional digits,
// "Z" only. This is the single formatter for every structured JSON/HTTP DTO
// time field, config-package metadata and filesystem ModTime that the D-009
// scope names; per-handler inline layouts are forbidden so the shape cannot
// drift.
//
// Input: legal RFC3339 with 0, 3, 6 or 9 fractional digits and an equivalent
// offset, normalized to UTC; a zoneless or ambiguous local timestamp is
// rejected (ambiguous local time must never be guessed).
//
// Human prose (email bodies, audit detail JSON strings, arbitrary payload
// fields) stays outside this contract per D-009 and keeps its own rules.

// FormatWireTime renders a domain instant in the fixed-6 wire form.
//
// It truncates toward zero to microsecond precision through the shared codec,
// so a nanosecond-bearing time.Time can never round into the next microsecond
// (Root D-008: truncate, never round).
func FormatWireTime(t time.Time) string {
	return temporal.FormatWire(t)
}

// FormatWireTimePtr renders an optional instant. The bool reports whether a
// value was present, so callers keep control of the absent-field spelling
// (null, omitted, or an explicit empty string) instead of receiving a
// fabricated instant.
func FormatWireTimePtr(t *time.Time) (string, bool) {
	if t == nil {
		return "", false
	}
	return FormatWireTime(*t), true
}

// ErrWireTimeInvalid reports input that is not an accepted wire timestamp.
var ErrWireTimeInvalid = errors.New("handler: invalid wire timestamp")

// ParseWireTime parses an input wire timestamp per D-005: 0/3/6/9 fractional
// digits and any explicit offset are accepted and normalized to UTC, while a
// zoneless value is rejected.
func ParseWireTime(value string) (time.Time, error) {
	parsed, err := temporal.Parse(value)
	if err != nil {
		return time.Time{}, ErrWireTimeInvalid
	}
	return parsed, nil
}
