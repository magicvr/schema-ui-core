package migration

import (
	"context"
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

// recordsPersistPGDDL is the postgres variant of recordsPersistDDL:
// updated_at (Unix time) is BIGINT (R1 v1.4 §3).
var recordsPersistPGDDL = []string{
	`CREATE TABLE records (
  id         TEXT PRIMARY KEY,
  name       TEXT NOT NULL CHECK (length(trim(name)) > 0),
  status     TEXT NOT NULL CHECK (length(trim(status)) > 0),
  owner      TEXT NOT NULL CHECK (length(trim(owner)) > 0),
  updated_at BIGINT NOT NULL
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

// mailOutboxDDL creates the mock-channel outbound record table
// (VP-017 R6 / workspace-017 GOAL-007; contract frozen by GOAL-006 D-002 §3):
// one row per accepted mock Send, listed newest-first by created_at and
// evicted oldest-first beyond the bounded retention cap.
var mailOutboxDDL = []string{
	`CREATE TABLE mail_outbox (
  id         TEXT PRIMARY KEY,
  to_addr    TEXT NOT NULL,
  subject    TEXT NOT NULL,
  body       TEXT NOT NULL,
  created_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_mail_outbox_created_at ON mail_outbox(created_at)`,
}

// mailOutboxPGDDL is the postgres variant of mailOutboxDDL:
// created_at (Unix time) is BIGINT (R1 v1.4 §3 convention).
var mailOutboxPGDDL = []string{
	`CREATE TABLE mail_outbox (
  id         TEXT PRIMARY KEY,
  to_addr    TEXT NOT NULL,
  subject    TEXT NOT NULL,
  body       TEXT NOT NULL,
  created_at BIGINT NOT NULL
)`,
	`CREATE INDEX idx_mail_outbox_created_at ON mail_outbox(created_at)`,
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
			ApplyPostgres:        migrateRecordsPersistPG,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "records_retire"},
			Version:              6,
			Name:                 "records_retire",
			Checksum:             kernel.MigrationChecksum(recordsRetireDDL, "0006:records-retire:v1"),
			Apply:                migrateRecordsRetire,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "mail_outbox"},
			Version:              51,
			Name:                 "mail_outbox",
			Checksum:             kernel.MigrationChecksum(mailOutboxDDL, "0051:mail-outbox:v1"),
			Apply:                migrateMailOutbox,
			ApplyPostgres:        migrateMailOutboxPG,
		},
	}
}

func migrateRecordsPersist(tx kernel.Tx) error {
	for _, stmt := range recordsPersistDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create records: %w", err)
		}
	}
	return nil
}

func migrateRecordsPersistPG(tx kernel.Tx) error {
	for _, stmt := range recordsPersistPGDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create records (postgres): %w", err)
		}
	}
	return nil
}

func migrateRecordsRetire(tx kernel.Tx) error {
	for _, stmt := range recordsRetireDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("records retire: %w", err)
		}
	}
	return nil
}

func migrateMailOutbox(tx kernel.Tx) error {
	for _, stmt := range mailOutboxDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create mail_outbox: %w", err)
		}
	}
	return nil
}

func migrateMailOutboxPG(tx kernel.Tx) error {
	for _, stmt := range mailOutboxPGDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create mail_outbox (postgres): %w", err)
		}
	}
	return nil
}
