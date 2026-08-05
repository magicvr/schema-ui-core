package migration

import (
	"database/sql"
	"fmt"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

const ModuleID = "core.persistence"

var recordsPersistDDL = []string{
	`CREATE TABLE records (
  id         TEXT PRIMARY KEY,
  name       TEXT NOT NULL CHECK (length(trim(name)) > 0),
  status     TEXT NOT NULL CHECK (length(trim(status)) > 0),
  owner      TEXT NOT NULL CHECK (length(trim(owner)) > 0),
  updated_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_records_name ON records(name)`,
	`CREATE INDEX idx_records_updated_at ON records(updated_at)`,
	`CREATE INDEX idx_records_owner ON records(owner)`,
}

var recordsRetireDDL = []string{
	`DROP TABLE IF EXISTS records`,
	`DELETE FROM role_permissions WHERE permission_id IN ('perm-records-read','perm-records-write')`,
	`DELETE FROM role_menu_items WHERE menu_item_id = 'menu-list-edit-lifecycle'`,
	`DELETE FROM permissions WHERE id IN ('perm-records-read','perm-records-write')`,
	`DELETE FROM menu_items WHERE id = 'menu-list-edit-lifecycle'`,
}

// Descriptors preserves the immutable records migration history after the
// product surface was retired. The frozen Apply behavior remains executable for
// fresh and upgrading databases; this historical owner is not a current Records
// capability.
func Descriptors() []kernel.MigrationContribution {
	return []kernel.MigrationContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "records_persist"},
			Version:              3,
			Name:                 "records_persist",
			Checksum:             kernel.MigrationChecksum(recordsPersistDDL, "0003:records-persist:v1"),
			Apply:                migrateRecordsPersist,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "records_retire"},
			Version:              6,
			Name:                 "records_retire",
			Checksum:             kernel.MigrationChecksum(recordsRetireDDL, "0006:records-retire:v1"),
			Apply:                migrateRecordsRetire,
		},
	}
}

func migrateRecordsPersist(tx *sql.Tx) error {
	for _, stmt := range recordsPersistDDL {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("create records: %w", err)
		}
	}
	return nil
}

func migrateRecordsRetire(tx *sql.Tx) error {
	for _, stmt := range recordsRetireDDL {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("records retire: %w", err)
		}
	}
	return nil
}
