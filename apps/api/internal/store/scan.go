package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/temporal"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// Read-side temporal adapter (workspace-040 R2 / GOAL-004 D-001 §1).
//
// database/sql converts a time.Time *source* into a *time.Time destination, but
// it does not parse a string: scanning the canonical SQLite text
// "YYYY-MM-DDTHH:MM:SS.ffffffZ" into *time.Time fails with
// "unsupported Scan, storing driver.Value type string into type *time.Time".
// Modules must therefore be free to scan domain time types while the dialect
// difference (canonical TEXT vs native timestamptz) stays inside the store, so
// every sqlTx/pgTx query result is wrapped and its Scan destinations are
// normalized here:
//
//   - *time.Time, *sql.NullTime and **time.Time destinations are read through a
//     temporary and converted from whatever the driver returned (string / []byte
//     / time.Time / nil);
//   - a stored sqlite text value must be the canonical fixed-6 UTC form —
//     anything else fails closed rather than being silently accepted (the
//     non-canonical-text irreversible point of the C2 contract);
//   - every other destination is passed through untouched.

type scanRows struct{ rows *sql.Rows }

func (r scanRows) Next() bool { return r.rows.Next() }

func (r scanRows) Scan(dest ...any) error { return scanInto(r.rows, dest) }

func (r scanRows) Close() error { return r.rows.Close() }

func (r scanRows) Err() error { return r.rows.Err() }

type scanRow struct{ row *sql.Row }

func (r scanRow) Scan(dest ...any) error { return scanInto(r.row, dest) }

type rowScanner interface{ Scan(dest ...any) error }

// instantTarget remembers one domain-time destination and the temporary the
// driver actually fills.
type instantTarget struct {
	index int
	dest  any
	raw   *any
}

func scanInto(scanner rowScanner, dest []any) error {
	targets := instantTargets(dest)
	if len(targets) == 0 {
		return scanner.Scan(dest...)
	}
	args := make([]any, len(dest))
	copy(args, dest)
	for i := range targets {
		raw := new(any)
		targets[i].raw = raw
		args[targets[i].index] = raw
	}
	if err := scanner.Scan(args...); err != nil {
		return err
	}
	for i := range targets {
		instant, valid, err := storedInstant(*targets[i].raw)
		if err != nil {
			return fmt.Errorf("scan column %d: %w", targets[i].index, err)
		}
		if err := assignInstant(targets[i].dest, instant, valid); err != nil {
			return fmt.Errorf("scan column %d: %w", targets[i].index, err)
		}
	}
	return nil
}

// instantTargets returns the domain-time destinations of one Scan call.
func instantTargets(dest []any) []instantTarget {
	var targets []instantTarget
	for index, d := range dest {
		switch d.(type) {
		case *time.Time, *sql.NullTime, **time.Time:
			targets = append(targets, instantTarget{index: index, dest: d})
		}
	}
	return targets
}

// storedInstant converts one driver value into a domain instant. valid=false
// reports SQL NULL.
func storedInstant(value any) (time.Time, bool, error) {
	switch v := value.(type) {
	case nil:
		return time.Time{}, false, nil
	case time.Time:
		return v.UTC(), true, nil
	case string:
		return parseStoredCanonical(v)
	case []byte:
		return parseStoredCanonical(string(v))
	default:
		return time.Time{}, false, fmt.Errorf("unsupported stored time value type %T", value)
	}
}

// parseStoredCanonical accepts only the canonical stored form: it round-trips
// through the codec byte-identically. Non-canonical text (locally offset,
// 3/9-digit fractions, or any other spelling) is a contract violation and fails
// closed.
func parseStoredCanonical(s string) (time.Time, bool, error) {
	parsed, err := temporal.Parse(s)
	if err != nil {
		return time.Time{}, false, fmt.Errorf("stored time %q: %w", s, err)
	}
	canonical, err := temporal.Format(parsed)
	if err != nil {
		return time.Time{}, false, fmt.Errorf("stored time %q: %w", s, err)
	}
	if canonical != s {
		return time.Time{}, false, fmt.Errorf("stored time %q is not canonical fixed-6 UTC", s)
	}
	return parsed, true, nil
}

func assignInstant(dest any, instant time.Time, valid bool) error {
	switch d := dest.(type) {
	case *time.Time:
		if !valid {
			return fmt.Errorf("converting NULL to time.Time is unsupported")
		}
		*d = instant
	case *sql.NullTime:
		*d = sql.NullTime{Time: instant, Valid: valid}
	case **time.Time:
		if !valid {
			*d = nil
			return nil
		}
		value := instant
		*d = &value
	default:
		return fmt.Errorf("unsupported temporal destination %T", dest)
	}
	return nil
}

// kernel.Rows / kernel.Row conformance is asserted at compile time so a future
// signature change in the kernel port cannot silently drop the adapter.
var (
	_ kernel.Rows = scanRows{}
	_ kernel.Row  = scanRow{}
)
