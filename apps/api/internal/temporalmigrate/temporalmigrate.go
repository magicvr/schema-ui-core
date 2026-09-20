// Package temporalmigrate carries the shared, dialect-neutral mechanics of the
// workspace-040 v73–v87 timestamp conversions (Root GOAL-001 R2 / GOAL-003
// M2). It owns three things and nothing else:
//
//   - the ordered-statement rule: the canonical checksum input of a conversion
//     descriptor is m0 (sentinel preflight) → m1–m3 (constraint release, table
//     rebuild, constraint/index rebuild) → m4 (post-rebuild verification), in
//     that literal order (r1-c2-descriptor-ledger-v1.0-fc.md §2, D-017);
//   - the executable m0 and m4 steps over kernel.Tx;
//   - a plain ordered executor for the m1–m3 rebuild slice.
//
// The frozen conversion expressions themselves are NOT built here: every
// CREATE TABLE / INSERT … SELECT / ALTER … USING statement is a literal in the
// owning module package, so the checksum covers exactly the SQL that runs
// (D-018 option C).
//
// Dialect note: m0 uses only portable SQL. m4 verifies the SQLite physical
// shape (sqlite_master / pragma_table_info) and is therefore only used by the
// SQLite Apply bodies; the postgres Apply bodies verify through
// information_schema in their own literal statements.
package temporalmigrate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// Guard is a fail-closed precondition that a table retired by an earlier
// version is still absent. The frozen v73 scope names one: the retired
// `records` table (r1-c2-descriptor-ledger-v1.0-fc.md §1). The normal catalog
// path drops it in v6, so the guard only fires on an abnormal path (a skipped
// v6, or a restore that brought it back) — exactly the case where the
// conversion would otherwise leave `records.updated_at` unconverted.
type Guard struct {
	Table string
}

// GuardStatements returns the exact m0 guard SQL (one statement per table).
func GuardStatements(guards []Guard) []string {
	out := make([]string, 0, len(guards))
	for _, guard := range guards {
		out = append(out, fmt.Sprintf(
			`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = '%s'`, guard.Table))
	}
	return out
}

// PostgresGuardStatement is the postgres counterpart of GuardStatements.
func PostgresGuardStatement(guard Guard) string {
	return fmt.Sprintf(
		`SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = current_schema() AND table_name = '%s'`,
		guard.Table)
}

// RunGuards executes the m0 guards on SQLite.
func RunGuards(tx kernel.Tx, guards []Guard, label string) error {
	ctx := context.Background()
	for _, guard := range guards {
		count, err := count1(tx, ctx, fmt.Sprintf(
			`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = '%s'`, guard.Table))
		if err != nil {
			return fmt.Errorf("%s: probe retired table %s: %w", label, guard.Table, err)
		}
		if count != 0 {
			return fmt.Errorf("%s: retired table %s is present; migration refused", label, guard.Table)
		}
	}
	return nil
}

// RunPostgresGuards executes the m0 guards on postgres.
func RunPostgresGuards(tx kernel.Tx, guards []Guard, label string) error {
	ctx := context.Background()
	for _, guard := range guards {
		count, err := count1(tx, ctx, PostgresGuardStatement(guard))
		if err != nil {
			return fmt.Errorf("%s: probe retired table %s: %w", label, guard.Table, err)
		}
		if count != 0 {
			return fmt.Errorf("%s: retired table %s is present; migration refused", label, guard.Table)
		}
	}
	return nil
}

// Preflight is the m0 sentinel census of one converted column: the migration
// refuses to convert a column whose legacy values cannot be classified.
//
// Policy (Root D-012 / D-015, corrected by A-042 and re-reviewed by A-046):
//
//   - Voucher columns (vouchers.expires_at / vouchers.redeemed_at) treat a
//     negative epoch as data corruption: bucket_negative > 0 fails closed.
//   - Every other column treats a negative epoch as a legal instant before
//     1970: the count is reported but never blocks.
//   - A zero in a sentinel column maps to NULL; a zero elsewhere converts as
//     the epoch instant.
type Preflight struct {
	Table  string
	Column string
	// Voucher marks the two legacy voucher columns (Root D-012).
	Voucher bool
}

// Buckets is one preflight census result.
type Buckets struct {
	Zero     int64
	Negative int64
	Positive int64
}

// Verify describes the m4 post-rebuild assertions of one converted table.
type Verify struct {
	Table string
	// Columns are the columns converted by this descriptor; each must end up
	// TEXT on SQLite.
	Columns []string
	// Children are the FK children recreated back onto Table; each must
	// reference the new table, never the retired <table>_old name.
	Children []string
}

// PreflightStatement returns the exact m0 SQL for one sentinel column.
func PreflightStatement(check Preflight) string {
	return fmt.Sprintf(
		`SELECT SUM(CASE WHEN "%s" = 0 THEN 1 ELSE 0 END) AS bucket_zero, `+
			`SUM(CASE WHEN "%s" < 0 THEN 1 ELSE 0 END) AS bucket_negative, `+
			`SUM(CASE WHEN "%s" > 0 THEN 1 ELSE 0 END) AS bucket_positive FROM "%s"`,
		check.Column, check.Column, check.Column, check.Table)
}

// PreflightStatements returns the ordered m0 statements (one per column).
func PreflightStatements(checks []Preflight) []string {
	out := make([]string, 0, len(checks))
	for _, check := range checks {
		out = append(out, PreflightStatement(check))
	}
	return out
}

// VerifyStatements returns the ordered m4 statements of one descriptor.
//
// The two PRAGMA checks are emitted once per descriptor, mirroring RunVerify
// exactly, so the hashed statement list stays identical to what executes
// (A-002 F-I-002 of GOAL-003).
func VerifyStatements(checks []Verify) []string {
	var out []string
	for _, check := range checks {
		out = append(out, fmt.Sprintf(
			`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = '%s_old'`,
			check.Table))
		for _, column := range check.Columns {
			out = append(out, fmt.Sprintf(
				`SELECT COUNT(*) FROM pragma_table_info('%s') WHERE name = '%s' AND type = 'TEXT'`,
				check.Table, column))
		}
		for _, child := range check.Children {
			out = append(out, fmt.Sprintf(
				`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = '%s' `+
					`AND (instr(lower(sql), 'references %s(') > 0 OR instr(lower(sql), 'references "%s"(') > 0)`,
				child, check.Table, check.Table))
		}
	}
	if len(checks) > 0 {
		out = append(out,
			`PRAGMA foreign_key_check`,
			`PRAGMA integrity_check`)
	}
	return out
}

// OrderedWithGuards returns the canonical checksum input of a descriptor: m0
// (retired-table guards, then the sentinel census), the m1–m3 rebuild slice,
// then m4, in that literal order (D-017).
func OrderedWithGuards(guards []Guard, preflight []Preflight, rebuild []string, verify []Verify) []string {
	out := GuardStatements(guards)
	out = append(out, PreflightStatements(preflight)...)
	out = append(out, rebuild...)
	return append(out, VerifyStatements(verify)...)
}

// Ordered is OrderedWithGuards for descriptors that declare no retired-table
// guard.
func Ordered(preflight []Preflight, rebuild []string, verify []Verify) []string {
	return OrderedWithGuards(nil, preflight, rebuild, verify)
}

// RunPreflight executes the m0 census and fails closed when a voucher column
// carries a negative epoch.
func RunPreflight(tx kernel.Tx, checks []Preflight) error {
	for _, check := range checks {
		buckets, err := census(tx, check)
		if err != nil {
			return err
		}
		if check.Voucher && buckets.Negative > 0 {
			return fmt.Errorf(
				"preflight %s.%s: %d negative epoch value(s) are data corruption for a voucher column (Root D-012); migration refused",
				check.Table, check.Column, buckets.Negative)
		}
	}
	return nil
}

func census(tx kernel.Tx, check Preflight) (Buckets, error) {
	var zero, negative, positive sql.NullInt64
	if err := tx.QueryRow(context.Background(), PreflightStatement(check)).
		Scan(&zero, &negative, &positive); err != nil {
		return Buckets{}, fmt.Errorf("preflight %s.%s: %w", check.Table, check.Column, err)
	}
	return Buckets{Zero: zero.Int64, Negative: negative.Int64, Positive: positive.Int64}, nil
}

// Exec runs an ordered statement slice of one descriptor. Every statement is a
// literal of the owning descriptor, so the executed SQL is the hashed SQL.
func Exec(tx kernel.Tx, stmts []string, label string) error {
	for _, stmt := range stmts {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("%s: %w", label, err)
		}
	}
	return nil
}

// RunVerify executes the m4 assertions: the retired table is gone, every
// converted column is TEXT, every recreated child references the live parent,
// and the database is internally consistent.
func RunVerify(tx kernel.Tx, checks []Verify, label string) error {
	ctx := context.Background()
	for _, check := range checks {
		if count, err := count1(tx, ctx, fmt.Sprintf(
			`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = '%s_old'`,
			check.Table)); err != nil {
			return fmt.Errorf("%s: verify retired table %s_old: %w", label, check.Table, err)
		} else if count != 0 {
			return fmt.Errorf("%s: retired table %s_old still exists", label, check.Table)
		}
		for _, column := range check.Columns {
			count, err := count1(tx, ctx, fmt.Sprintf(
				`SELECT COUNT(*) FROM pragma_table_info('%s') WHERE name = '%s' AND type = 'TEXT'`,
				check.Table, column))
			if err != nil {
				return fmt.Errorf("%s: verify %s.%s: %w", label, check.Table, column, err)
			}
			if count != 1 {
				return fmt.Errorf("%s: %s.%s is not TEXT after the rebuild", label, check.Table, column)
			}
		}
		for _, child := range check.Children {
			count, err := count1(tx, ctx, fmt.Sprintf(
				`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = '%s' `+
					`AND (instr(lower(sql), 'references %s(') > 0 OR instr(lower(sql), 'references "%s"(') > 0)`,
				child, check.Table, check.Table))
			if err != nil {
				return fmt.Errorf("%s: verify child %s: %w", label, child, err)
			}
			if count != 1 {
				return fmt.Errorf("%s: child %s does not reference %s", label, child, check.Table)
			}
		}
	}
	if err := assertNoForeignKeyViolations(tx, ctx, label); err != nil {
		return err
	}
	return assertIntegrity(tx, ctx, label)
}

func assertNoForeignKeyViolations(tx kernel.Tx, ctx context.Context, label string) error {
	rows, err := tx.Query(ctx, `PRAGMA foreign_key_check`)
	if err != nil {
		return fmt.Errorf("%s: foreign_key_check: %w", label, err)
	}
	defer rows.Close()
	if rows.Next() {
		return fmt.Errorf("%s: foreign_key_check reported a violation", label)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("%s: foreign_key_check: %w", label, err)
	}
	return nil
}

func assertIntegrity(tx kernel.Tx, ctx context.Context, label string) error {
	var result string
	if err := tx.QueryRow(ctx, `PRAGMA integrity_check`).Scan(&result); err != nil {
		return fmt.Errorf("%s: integrity_check: %w", label, err)
	}
	if strings.TrimSpace(result) != "ok" {
		return fmt.Errorf("%s: integrity_check = %q, want ok", label, result)
	}
	return nil
}

// PgVerify describes one postgres m4 column assertion (the PG variant is not
// part of the checksum input, D-017; the assertion carries the dialect
// difference instead).
type PgVerify struct {
	Table   string
	Column  string
	NonNull bool
}

// PostgresVerifyStatement returns the exact information_schema assertion for
// one converted postgres column.
func PostgresVerifyStatement(check PgVerify) string {
	nullable := "YES"
	if check.NonNull {
		nullable = "NO"
	}
	return fmt.Sprintf(
		`SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = current_schema() `+
			`AND table_name = '%s' AND column_name = '%s' `+
			`AND data_type = 'timestamp with time zone' AND datetime_precision = 6 AND is_nullable = '%s'`,
		check.Table, check.Column, nullable)
}

// RunPostgresVerify executes the postgres m4 column assertions.
func RunPostgresVerify(tx kernel.Tx, checks []PgVerify, label string) error {
	ctx := context.Background()
	for _, check := range checks {
		count, err := count1(tx, ctx, PostgresVerifyStatement(check))
		if err != nil {
			return fmt.Errorf("%s: verify %s.%s: %w", label, check.Table, check.Column, err)
		}
		if count != 1 {
			return fmt.Errorf("%s: %s.%s is not timestamptz(6) with the expected nullability", label, check.Table, check.Column)
		}
	}
	return nil
}

func count1(tx kernel.Tx, ctx context.Context, query string) (int64, error) {
	var count int64
	err := tx.QueryRow(ctx, query).Scan(&count)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	return count, nil
}
