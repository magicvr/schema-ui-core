// Migration runner for the SQLite store.
//
// Module packages own migration SQL and Apply functions. This package owns only
// the platform runner: transaction boundaries, the immutable global ledger,
// upgrade snapshots, integrity checks, and fail-closed applied-history checks.
package store

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	migrationcontract "github.com/magicvr/schema-ui-core/apps/api/internal/migration"
)

type appliedMigration struct {
	version  int
	name     string
	checksum string
}

// migrate applies the caller-supplied compiled-global catalog. No module SQL or
// migration registry is retained in store; catalog.Apply is the only execution
// source.
func (s *Store) migrate(catalog []kernel.MigrationContribution) error {
	if err := s.assertForeignKeysOn(); err != nil {
		return err
	}

	applied, err := s.appliedMigrations()
	if err != nil {
		return err
	}
	if applied == nil {
		// Migration 0001 bootstraps both the baseline schema and the ledger in a
		// single transaction. Existing ledger-less R2 databases are fingerprinted
		// by that module-owned Apply function before the ledger is created.
		if err := s.applyMigration(catalog[0]); err != nil {
			return err
		}
		applied, err = s.appliedMigrations()
		if err != nil {
			return err
		}
	}
	if err := validateApplied(applied, catalog); err != nil {
		return err
	}

	// One recoverable snapshot per pending data-mutating migration (version >= 2,
	// I-011-002 A-002 F-002): each upgrade step keeps an independent rollback
	// point. Snapshot filenames carry millisecond precision so an immediate
	// retry after a failed upgrade cannot collide (D5).
	for _, migration := range pendingMigrations(applied, catalog) {
		if migration.Version >= 2 {
			if err := s.snapshotBeforePending(migration.Version); err != nil {
				return err
			}
		}
		if err := s.applyMigration(migration); err != nil {
			return err
		}
	}
	return s.verifyIntegrity()
}

// applyMigration runs one module-owned Apply function and its ledger insert in
// one transaction. Tombstones have no Apply function and record a no-op history
// entry; all active descriptors must supply Apply.
func (s *Store) applyMigration(migration kernel.MigrationContribution) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin migration %d (%s): %w", migration.Version, migration.Name, err)
	}
	if migration.Apply != nil {
		if err := migration.Apply(tx); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("migration %d (%s): %w", migration.Version, migration.Name, err)
		}
	}
	if _, err := tx.Exec(
		`INSERT INTO schema_migrations (version, name, checksum, applied_at) VALUES (?, ?, ?, ?)`,
		migration.Version, migration.Name, migration.Checksum, time.Now().UTC().Unix(),
	); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("record migration %d (%s): %w", migration.Version, migration.Name, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration %d (%s): %w", migration.Version, migration.Name, err)
	}
	return nil
}

// normalizeCatalog validates the platform-facing contract and returns an
// isolated, version-sorted copy. Production normally receives an already
// validated kernel.CollectPersistence result; the store repeats the ledger
// invariants so direct callers cannot bypass fail-closed startup behavior.
func normalizeCatalog(catalog []kernel.MigrationContribution) ([]kernel.MigrationContribution, error) {
	if len(catalog) == 0 {
		return nil, errors.New("store: compiled migration catalog is empty")
	}
	entries := make([]migrationcontract.Entry, 0, len(catalog))
	clone := append([]kernel.MigrationContribution(nil), catalog...)
	for _, migration := range clone {
		if strings.TrimSpace(migration.ModuleID) == "" || migration.ModuleID != strings.TrimSpace(migration.ModuleID) {
			return nil, fmt.Errorf("store: migration %d (%s) has invalid module id %q", migration.Version, migration.Name, migration.ModuleID)
		}
		if strings.TrimSpace(migration.Name) == "" || migration.Key != migration.Name {
			return nil, fmt.Errorf("store: migration %d has invalid identity key=%q name=%q", migration.Version, migration.Key, migration.Name)
		}
		if migration.Tombstone && migration.Apply != nil {
			return nil, fmt.Errorf("store: tombstone migration %d (%s) must not carry Apply", migration.Version, migration.Name)
		}
		if !migration.Tombstone && migration.Apply == nil {
			return nil, fmt.Errorf("store: migration %d (%s) requires Apply", migration.Version, migration.Name)
		}
		entries = append(entries, migrationcontract.Entry{
			Version: migration.Version, Name: migration.Name,
			ModuleID: migration.ModuleID, Checksum: migration.Checksum,
		})
	}
	if _, err := migrationcontract.Collect(entries); err != nil {
		return nil, fmt.Errorf("store: compiled migration contract: %w", err)
	}
	sort.Slice(clone, func(i, j int) bool { return clone[i].Version < clone[j].Version })
	return clone, nil
}

// validateApplied verifies the ledger is a contiguous prefix of the supplied
// compiled catalog with matching names and checksums.
func validateApplied(applied []appliedMigration, catalog []kernel.MigrationContribution) error {
	known := make(map[int]kernel.MigrationContribution, len(catalog))
	for _, migration := range catalog {
		known[migration.Version] = migration
	}
	if len(applied) == 0 {
		return errors.New("store: schema_migrations exists but is empty (partial bootstrap)")
	}
	if applied[0].version != 1 {
		return fmt.Errorf("store: migration ledger starts at version %d, want 1", applied[0].version)
	}
	for i, row := range applied {
		if i > 0 && row.version != applied[i-1].version+1 {
			return fmt.Errorf("store: migration ledger missing intermediate version before %d", row.version)
		}
		migration, ok := known[row.version]
		if !ok {
			return fmt.Errorf("store: unknown applied migration version %d (%s)", row.version, row.name)
		}
		if migration.Name != row.name {
			return fmt.Errorf("store: applied migration %d name %q, code has %q", row.version, row.name, migration.Name)
		}
		if migration.Checksum != row.checksum {
			return fmt.Errorf("store: migration %d checksum drift (ledger %s, code %s)", row.version, row.checksum, migration.Checksum)
		}
	}
	return nil
}

func pendingMigrations(applied []appliedMigration, catalog []kernel.MigrationContribution) []kernel.MigrationContribution {
	appliedSet := make(map[int]bool, len(applied))
	for _, row := range applied {
		appliedSet[row.version] = true
	}
	pending := make([]kernel.MigrationContribution, 0, len(catalog)-len(appliedSet))
	for _, migration := range catalog {
		if !appliedSet[migration.Version] {
			pending = append(pending, migration)
		}
	}
	return pending
}

// appliedMigrations returns nil when no ledger table exists, otherwise the
// applied rows ordered by version. An existing empty ledger is a partial
// bootstrap and is rejected.
func (s *Store) appliedMigrations() ([]appliedMigration, error) {
	var exists bool
	if err := s.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM sqlite_master WHERE type = 'table' AND name = 'schema_migrations')`,
	).Scan(&exists); err != nil {
		return nil, fmt.Errorf("store: check migration ledger: %w", err)
	}
	if !exists {
		return nil, nil
	}
	rows, err := s.db.Query(`SELECT version, name, checksum FROM schema_migrations ORDER BY version`)
	if err != nil {
		return nil, fmt.Errorf("store: read migration ledger: %w", err)
	}
	defer rows.Close()
	var applied []appliedMigration
	for rows.Next() {
		var row appliedMigration
		if err := rows.Scan(&row.version, &row.name, &row.checksum); err != nil {
			return nil, fmt.Errorf("store: scan migration ledger: %w", err)
		}
		applied = append(applied, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: read migration ledger: %w", err)
	}
	if len(applied) == 0 {
		return nil, errors.New("store: schema_migrations exists but is empty (partial bootstrap)")
	}
	return applied, nil
}

func (s *Store) assertForeignKeysOn() error {
	if _, err := s.db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		return fmt.Errorf("store: enable foreign_keys: %w", err)
	}
	var on int
	if err := s.db.QueryRow(`PRAGMA foreign_keys`).Scan(&on); err != nil {
		return fmt.Errorf("store: read foreign_keys: %w", err)
	}
	if on != 1 {
		return errors.New("store: PRAGMA foreign_keys could not be enabled")
	}
	return nil
}

// snapshotBeforePending produces a recoverable copy of a non-empty file
// database before the first pending data-mutating migration of an upgrade
// batch. The filename carries millisecond precision so an immediate retry
// after a failed upgrade cannot collide with an existing snapshot (D5).
func (s *Store) snapshotBeforePending(firstPendingVersion int) error {
	if s.path == "" || s.path == ":memory:" {
		return nil
	}
	hasData, err := s.dbHasRows()
	if err != nil {
		return err
	}
	if !hasData {
		return nil
	}
	target := fmt.Sprintf("%s.pre-v%04d-%s.sqlite", s.path, firstPendingVersion, time.Now().UTC().Format("20060102T150405.000Z"))
	if _, err := s.db.Exec("VACUUM INTO '" + strings.ReplaceAll(target, "'", "''") + "'"); err != nil {
		return fmt.Errorf("pre-v%04d snapshot to %s: %w", firstPendingVersion, target, err)
	}
	if err := checkIntegrityFile(target); err != nil {
		return fmt.Errorf("pre-v%04d snapshot %s invalid: %w", firstPendingVersion, target, err)
	}
	return nil
}

func (s *Store) dbHasRows() (bool, error) {
	rows, err := s.db.Query(`SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%' AND name != 'schema_migrations'`)
	if err != nil {
		return false, fmt.Errorf("snapshot: list tables: %w", err)
	}
	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			rows.Close()
			return false, fmt.Errorf("snapshot: scan table: %w", err)
		}
		tables = append(tables, name)
	}
	if err := rows.Close(); err != nil {
		return false, fmt.Errorf("snapshot: list tables: %w", err)
	}
	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("snapshot: list tables: %w", err)
	}
	for _, name := range tables {
		var count int
		if err := s.db.QueryRow(`SELECT COUNT(*) FROM "` + strings.ReplaceAll(name, `"`, `""`) + `"`).Scan(&count); err != nil {
			return false, fmt.Errorf("snapshot: count %s: %w", name, err)
		}
		if count > 0 {
			return true, nil
		}
	}
	return false, nil
}

func checkIntegrityFile(path string) error {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("open snapshot: %w", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	return checkIntegrity(db)
}

func (s *Store) verifyIntegrity() error {
	if err := checkIntegrity(s.db); err != nil {
		return err
	}
	rows, err := s.db.Query(`PRAGMA foreign_key_check`)
	if err != nil {
		return fmt.Errorf("store: foreign_key_check: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		var table, parent string
		var rowID int64
		if err := rows.Scan(&table, &rowID, &parent); err != nil {
			return fmt.Errorf("store: read foreign_key_check: %w", err)
		}
		return fmt.Errorf("store: foreign_key_check violation: %s rowid %d references %s", table, rowID, parent)
	}
	return rows.Err()
}

func checkIntegrity(db *sql.DB) error {
	var result string
	if err := db.QueryRow(`PRAGMA integrity_check`).Scan(&result); err != nil {
		return fmt.Errorf("integrity_check: %w", err)
	}
	if result != "ok" {
		return fmt.Errorf("integrity_check = %q, want ok", result)
	}
	return nil
}
