package handler

import (
	"errors"
	"strings"
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
// Input: legal RFC3339 with 0, 3, 6 or 9 fractional digits and a ZERO-EQUIVALENT
// offset ("Z", "+00:00", "-00:00"), normalized to UTC; a non-zero offset and a
// zoneless/ambiguous local timestamp are both rejected (Root D-005 §3: "非法时间、
// 非零 offset 解析失败；不接受模糊本地时间字符串"). The codec layer
// (internal/temporal) deliberately keeps accepting non-zero offsets for stored
// payloads; refusing them is this API input parser's policy, which GOAL-003
// D-001 §3 assigned to R3.
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
// digits and a zero-equivalent offset are accepted and normalized to UTC; a
// non-zero offset or a zoneless value is rejected.
//
// The offset is inspected on the raw string, because temporal.Parse normalizes
// to UTC and would otherwise hide the very offset D-005 refuses. "±00:00" and
// "Z" both carry offset zero and stay accepted. A comma fractional separator is
// refused too: Go's RFC3339 parser tolerates it, but RFC3339 itself defines the
// separator as "." and D-005 accepts only legal RFC3339 input.
func ParseWireTime(value string) (time.Time, error) {
	if strings.ContainsRune(value, ',') {
		return time.Time{}, ErrWireTimeInvalid
	}
	parsed, err := temporal.Parse(value)
	if err != nil {
		return time.Time{}, ErrWireTimeInvalid
	}
	if offset := carriedOffsetSeconds(value); offset != 0 {
		return time.Time{}, ErrWireTimeInvalid
	}
	return parsed, nil
}

// carriedOffsetSeconds reports the UTC offset (in seconds) the wire value
// literally carries. It only runs after temporal.Parse has accepted the value,
// so a parse failure here is unreachable; it therefore reports 0 (no offset
// found) rather than inventing a second error path.
func carriedOffsetSeconds(value string) int {
	withZone, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return 0
	}
	_, offset := withZone.Zone()
	return offset
}
