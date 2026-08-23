package kernel

import (
	"errors"
	"strings"
)

// IsDuplicateObject reports whether err is a postgres "already exists"
// failure (duplicate table / column / object). SQLite migrations do not
// produce these SQLSTATEs; the helper is for postgres Apply bodies that
// must adopt a schema that already has the object.
//
// Detection is textual so call sites stay driver-agnostic, matching
// IsUniqueViolation: pgx/stdlib formats `ERROR: relation "users" already
// exists (SQLSTATE 42P07)` (and 42701 / 42710 / 42P06).
func IsDuplicateObject(err error) bool {
	for e := err; e != nil; e = errors.Unwrap(e) {
		msg := e.Error()
		if msg == "" {
			continue
		}
		if strings.Contains(msg, "SQLSTATE 42P07") ||
			strings.Contains(msg, "SQLSTATE 42701") ||
			strings.Contains(msg, "SQLSTATE 42710") ||
			strings.Contains(msg, "SQLSTATE 42P06") {
			return true
		}
		if strings.Contains(msg, "already exists") &&
			(strings.Contains(msg, "relation ") ||
				strings.Contains(msg, "column ") ||
				strings.Contains(msg, "constraint ") ||
				strings.Contains(msg, "type ")) {
			return true
		}
	}
	return false
}
