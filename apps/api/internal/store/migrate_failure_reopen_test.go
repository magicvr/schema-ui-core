// A-012 F-007 closure: runner-level evidence that one migration's Apply is
// atomic across BOTH dialects — when a multi-statement Apply fails after its
// first DDL, the whole migration (schema objects AND ledger row) is rolled
// back with zero residue, and a corrected same-version catalog reopens the
// same database cleanly with exactly one ledger record per version.
//
// The runner already wraps Apply + ledger insert in one transaction
// (migrate.go applyMigration / postgres.go applyMigrationPG); these tests pin
// the failure path as an executable slice instead of an implementation
// inference (SQLite mandatory per A-012; PostgreSQL mirrors the same
// semantics through the live PG integration harness).
package store

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/pgtest"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// f007BootstrapDDL is the v1 bootstrap: the ledger plus one empty base table.
// Identical DDL executes on sqlite and postgres (placeholders rebound).
var f007BootstrapDDL = []string{
	`CREATE TABLE schema_migrations (version INTEGER PRIMARY KEY, name TEXT NOT NULL, checksum TEXT NOT NULL, applied_at BIGINT NOT NULL)`,
	`CREATE TABLE f007_users (id TEXT PRIMARY KEY, username TEXT NOT NULL UNIQUE)`,
}

// f007MultiStmtDDL is the multi-statement v2 migration body. The first
// statement creates an object; a subsequent injected error must roll BOTH the
// object and the ledger row back.
var f007MultiStmtDDL = []string{
	`CREATE TABLE f007_fail_items (id TEXT PRIMARY KEY, offer_id TEXT NOT NULL, name TEXT NOT NULL)`,
	`CREATE INDEX idx_f007_fail_items_offer ON f007_fail_items(offer_id)`,
}

// f007Catalog builds a two-version catalog. When fail is true, v2 executes
// both DDL statements and then returns an injected error mid-Apply; when
// false, v2 is the corrected same-version body (identical statements and
// checksum — the fix is only the absent injected error).
func f007Catalog(fail bool) []kernel.MigrationContribution {
	return []kernel.MigrationContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: "r3.test", Key: "bootstrap"},
			Version:              1,
			Name:                 "bootstrap",
			Checksum:             kernel.MigrationChecksum(f007BootstrapDDL, "f007:bootstrap:v1"),
			Apply: func(tx kernel.Tx) error {
				for _, stmt := range f007BootstrapDDL {
					if _, err := tx.Exec(context.Background(), stmt); err != nil {
						return err
					}
				}
				return nil
			},
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: "r3.test", Key: "multi_stmt_failure"},
			Version:              2,
			Name:                 "multi_stmt_failure",
			Checksum:             kernel.MigrationChecksum(f007MultiStmtDDL, "f007:multi-stmt:v1"),
			Apply: func(tx kernel.Tx) error {
				for _, stmt := range f007MultiStmtDDL {
					if _, err := tx.Exec(context.Background(), stmt); err != nil {
						return err
					}
				}
				if fail {
					return errors.New("injected mid-apply failure")
				}
				return nil
			},
		},
	}
}

// indexExistsDB reports whether a named index exists in the sqlite catalog.
func indexExistsDB(t *testing.T, db *sql.DB, name string) bool {
	t.Helper()
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name=?`, name).Scan(&n); err != nil {
		t.Fatalf("check index %s: %v", name, err)
	}
	return n == 1
}

// TestMigrateMidApplyFailureNoResidueThenReopen is the SQLite slice of
// A-012 F-007: failed multi-statement Apply leaves no schema object and no
// ledger row; a corrected same-version catalog reopens and applies cleanly.
func TestMigrateMidApplyFailureNoResidueThenReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "f007_fail_reopen.db")

	// Attempt 1: v2 fails after its first DDL statement.
	if _, err := OpenWithCatalog(path, f007Catalog(true)); err == nil || !strings.Contains(err.Error(), "injected mid-apply failure") {
		t.Fatalf("open with failing v2 = %v, want injected mid-apply failure", err)
	}

	// Zero residue: v1 committed; v2's table, index and ledger row absent.
	check := rawOpen(t, path)
	if !tableExistsDB(t, check, "schema_migrations") {
		check.Close()
		t.Fatal("schema_migrations must exist after v1 committed")
	}
	if tableExistsDB(t, check, "f007_fail_items") {
		check.Close()
		t.Fatal("failed v2 must leave no table residue")
	}
	if indexExistsDB(t, check, "idx_f007_fail_items_offer") {
		check.Close()
		t.Fatal("failed v2 must leave no index residue")
	}
	var ledgerCount int
	if err := check.QueryRow(`SELECT COUNT(*) FROM schema_migrations`).Scan(&ledgerCount); err != nil {
		check.Close()
		t.Fatalf("read ledger: %v", err)
	}
	if ledgerCount != 1 {
		check.Close()
		t.Fatalf("ledger rows after failed v2 = %d, want exactly 1 (v1 only)", ledgerCount)
	}
	var versions string
	if err := check.QueryRow(`SELECT group_concat(version) FROM schema_migrations`).Scan(&versions); err != nil {
		check.Close()
		t.Fatalf("read ledger versions: %v", err)
	}
	if versions != "1" {
		check.Close()
		t.Fatalf("ledger versions = %q, want 1", versions)
	}
	if err := check.Close(); err != nil {
		t.Fatal(err)
	}

	// Attempt 2: the corrected same-version catalog reopens the same file and
	// applies v2 cleanly — exactly one ledger record, schema objects present.
	st, err := OpenWithCatalog(path, f007Catalog(false))
	if err != nil {
		t.Fatalf("reopen with corrected catalog: %v", err)
	}
	defer st.Close()
	if err := st.verifyIntegrity(); err != nil {
		t.Fatalf("post-reopen integrity: %v", err)
	}
	applied, err := st.appliedMigrations()
	if err != nil {
		t.Fatalf("applied after reopen: %v", err)
	}
	if len(applied) != 2 || applied[0].version != 1 || applied[1].version != 2 {
		t.Fatalf("applied after reopen = %+v, want {1,2}", applied)
	}
	var v2Rows int
	if err := st.db.QueryRow(`SELECT COUNT(*) FROM schema_migrations WHERE version = 2`).Scan(&v2Rows); err != nil {
		t.Fatal(err)
	}
	if v2Rows != 1 {
		t.Fatalf("v2 ledger rows after reopen = %d, want exactly 1", v2Rows)
	}
	if !tableExistsDB(t, st.db, "f007_fail_items") {
		t.Fatal("corrected v2 must create f007_fail_items")
	}
	if !indexExistsDB(t, st.db, "idx_f007_fail_items_offer") {
		t.Fatal("corrected v2 must create idx_f007_fail_items_offer")
	}
}

// TestPostgresMigrateMidApplyFailureNoResidueThenReopen mirrors the SQLite
// slice on the live postgres integration harness (A-012 F-007 requirement 3:
// PostgreSQL same semantics enter the integration suite when PG is available).
func TestPostgresMigrateMidApplyFailureNoResidueThenReopen(t *testing.T) {
	dsn := pgtest.DSN()
	if dsn == "" {
		t.Skip("postgres test env not set (PG_TEST_*); skipping postgres mid-apply failure/reopen")
	}
	ctx := context.Background()
	scratch := scratchDSN(t, dsn, "r3f007fail")

	// Attempt 1: failed multi-statement Apply rolls back schema + ledger.
	if _, err := Open(ctx, OpenOptions{
		Dialect:        kernel.DialectPostgres,
		DSN:            scratch,
		ConnectTimeout: 10 * time.Second,
	}, f007Catalog(true)); err == nil || !strings.Contains(err.Error(), "injected mid-apply failure") {
		t.Fatalf("open with failing v2 = %v, want injected mid-apply failure", err)
	}

	check, err := sql.Open("pgx", scratch)
	if err != nil {
		t.Fatal(err)
	}
	var ledgerCount int
	if err := check.QueryRowContext(ctx, `SELECT count(*) FROM schema_migrations`).Scan(&ledgerCount); err != nil {
		check.Close()
		t.Fatal(err)
	}
	if ledgerCount != 1 {
		check.Close()
		t.Fatalf("ledger rows after failed v2 = %d, want exactly 1 (v1 only)", ledgerCount)
	}
	var tableExists bool
	if err := check.QueryRowContext(ctx, `SELECT to_regclass('f007_fail_items') IS NOT NULL`).Scan(&tableExists); err != nil {
		check.Close()
		t.Fatal(err)
	}
	if tableExists {
		check.Close()
		t.Fatal("failed v2 must leave no table residue on postgres")
	}
	if err := check.Close(); err != nil {
		t.Fatal(err)
	}

	// Attempt 2: corrected same-version catalog reopens and applies cleanly.
	st, err := Open(ctx, OpenOptions{
		Dialect:        kernel.DialectPostgres,
		DSN:            scratch,
		ConnectTimeout: 10 * time.Second,
	}, f007Catalog(false))
	if err != nil {
		t.Fatalf("reopen with corrected catalog: %v", err)
	}
	defer st.Close()
	pg := st.(*postgres)

	var v2Rows int
	if err := pg.db.QueryRowContext(ctx, `SELECT count(*) FROM schema_migrations WHERE version = 2`).Scan(&v2Rows); err != nil {
		t.Fatal(err)
	}
	if v2Rows != 1 {
		t.Fatalf("v2 ledger rows after reopen = %d, want exactly 1", v2Rows)
	}
	if err := pg.db.QueryRowContext(ctx, `SELECT to_regclass('f007_fail_items') IS NOT NULL`).Scan(&tableExists); err != nil {
		t.Fatal(err)
	}
	if !tableExists {
		t.Fatal("corrected v2 must create f007_fail_items on postgres")
	}
	var versions string
	if err := pg.db.QueryRowContext(ctx, `SELECT string_agg(version::text, ',' ORDER BY version) FROM schema_migrations`).Scan(&versions); err != nil {
		t.Fatal(err)
	}
	if versions != "1,2" {
		t.Fatalf("postgres ledger versions = %q, want 1,2", versions)
	}
}
