package kernel

import (
	"errors"
	"strings"
)

// IsUniqueViolation reports whether err is a unique-constraint violation on
// either supported store dialect (W9 F-001/F-011).
//
// Detection is textual so module repositories stay driver-agnostic (no
// driver-specific error types in the kernel contract):
//   - modernc/sqlite: "UNIQUE constraint failed: table.col" and the newer
//     "constraint failed: UNIQUE ..." shape (result code 2067/1555).
//   - pgx/stdlib:     `duplicate key value violates unique constraint "..."
//     (SQLSTATE 23505)`.
//
// The unwrap chain is walked so wrapped errors (%w) are still recognized.
// A caller that must distinguish WHICH constraint fired keeps its existing
// constraint-name substring check on top of this predicate (see
// authsession service-credential name detection).
func IsUniqueViolation(err error) bool {
	for e := err; e != nil; e = errors.Unwrap(e) {
		msg := e.Error()
		if msg == "" {
			continue
		}
		if strings.Contains(msg, "UNIQUE constraint failed") ||
			strings.Contains(msg, "constraint failed: UNIQUE") ||
			strings.Contains(msg, "duplicate key value violates unique constraint") ||
			strings.Contains(msg, "SQLSTATE 23505") {
			return true
		}
	}
	return false
}
