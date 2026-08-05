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

// WithTx exposes the platform transaction boundary to module repositories.
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
