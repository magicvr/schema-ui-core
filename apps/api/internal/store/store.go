// Package store owns the SQLite platform boundary: connection lifecycle,
// transactions, migration execution, snapshots, and readiness state. Domain
// repositories live in their owning modules.
package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"

	_ "modernc.org/sqlite"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// Store is the SQLite platform runner. Domain code consumes its transaction
// boundary and never receives the underlying database handle.
type Store struct {
	db              *sql.DB
	path            string
	fresh           bool
	systemDataReady atomic.Bool
}

// OpenWithCatalog opens the SQLite DB and applies the compiled-global catalog
// supplied by the composition root.
func OpenWithCatalog(path string, catalog []kernel.MigrationContribution) (*Store, error) {
	normalized, err := normalizeCatalog(catalog)
	if err != nil {
		return nil, err
	}
	return open(path, normalized)
}

func open(path string, catalog []kernel.MigrationContribution) (*Store, error) {
	if dir := filepath.Dir(path); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create db dir: %w", err)
		}
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	db.SetMaxOpenConns(1)

	fresh, err := databaseIsEmpty(db)
	if err != nil {
		_ = db.Close()
		return nil, err
	}
	st := &Store{db: db, path: path, fresh: fresh}
	if err := st.migrate(catalog); err != nil {
		_ = db.Close()
		return nil, err
	}
	return st, nil
}

// WasFresh reports whether the database was empty before migration.
func (s *Store) WasFresh() bool { return s.fresh }

// Dialect reports the dialect this store implements (kernel port; R1 v1.4 §2).
func (s *Store) Dialect() kernel.Dialect { return kernel.DialectSQLite }

// Run executes fn inside one transaction and exposes the dialect-neutral
// kernel.Tx (R1 v1.4 §2). This is the kernel port entry point; R4 completed the
// public-surface migration, so modules/jobs/handler speak kernel.Tx — the
// sqlite-only WithTx below is a retained test/adaptation seam, not part of the
// module public contract (A-001 F-001).
//
// Nested Run is forbidden and detected per-callback via a goroutine-local
// marker (R1 v1.4 A-008 F-004 permits a ctx value or the equivalent
// call-stack/goroutine-local): a Store-level flag would misclassify the
// concurrent Runs that the postgres pool is allowed to carry.
func (s *Store) Run(ctx context.Context, fn func(kernel.Tx) error) error {
	if err := enterRun(); err != nil {
		return err
	}
	defer leaveRun()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// R1 v1.4 §2: a panicking fn rolls the transaction back and re-panics, so
	// the tx handle is never abandoned and the caller still observes the panic.
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		}
	}()
	if err := fn(sqlTx{tx: tx}); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}
	return nil
}

// WithTx is the retained sqlite-only adaptation seam (tests / testsupport /
// legacy callers). It is NOT part of the dialect-neutral module contract, which
// goes through Run(func(kernel.Tx)); production wiring never injects it after
// R4 (A-001 F-001 — intentionally kept, documented as hygiene debt).
func (s *Store) WithTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// sqlTx adapts *sql.Tx to kernel.Tx. SQLite keeps the '?' placeholder, so no
// rebinding is applied here.
type sqlTx struct{ tx *sql.Tx }

func (t sqlTx) Exec(ctx context.Context, query string, args ...any) (kernel.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

func (t sqlTx) Query(ctx context.Context, query string, args ...any) (kernel.Rows, error) {
	return t.tx.QueryContext(ctx, query, args...)
}

func (t sqlTx) QueryRow(ctx context.Context, query string, args ...any) kernel.Row {
	return t.tx.QueryRowContext(ctx, query, args...)
}

// MarkSystemDataReady records successful post-finalize reconciliation.
func (s *Store) MarkSystemDataReady() { s.systemDataReady.Store(true) }

// SystemDataReady fails until finalized contributions have been reconciled.
func (s *Store) SystemDataReady() error {
	if !s.systemDataReady.Load() {
		return errors.New("store: system-data reconciliation is not ready")
	}
	return nil
}

func databaseIsEmpty(db *sql.DB) (bool, error) {
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%'`).Scan(&count); err != nil {
		return false, fmt.Errorf("check fresh database: %w", err)
	}
	return count == 0, nil
}

// Close releases the underlying database handle.
func (s *Store) Close() error { return s.db.Close() }

// Ping verifies the SQLite connection with a trivial read.
func (s *Store) Ping(ctx context.Context) error {
	var one int
	if err := s.db.QueryRowContext(ctx, "SELECT 1").Scan(&one); err != nil {
		return fmt.Errorf("sqlite ping: %w", err)
	}
	return nil
}
