package migration

import (
	"database/sql"
	"fmt"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

const ModuleID = "core.operationlog"

var operationLogDDL = []string{
	`CREATE TABLE operation_log (
  id         TEXT PRIMARY KEY,
  event      TEXT NOT NULL CHECK (event IN ('records.create','records.update','records.delete','auth.login','auth.logout','auth.refresh')),
  actor_id   TEXT NOT NULL,
  actor_name TEXT NOT NULL,
  record_id  TEXT,
  detail     TEXT,
  created_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC)`,
}

var operationLogExpandDDL = []string{
	`CREATE TABLE operation_log (
  id         TEXT PRIMARY KEY,
  event      TEXT NOT NULL CHECK (event IN ('records.create','records.update','records.delete','auth.login','auth.logout','auth.refresh','users.create','users.update','users.delete','roles.create','roles.update','roles.delete')),
  actor_id   TEXT NOT NULL,
  actor_name TEXT NOT NULL,
  record_id  TEXT,
  detail     TEXT,
  created_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC)`,
}

var operationLogSettingsDDL = []string{
	`CREATE TABLE operation_log (
  id         TEXT PRIMARY KEY,
  event      TEXT NOT NULL CHECK (event IN ('records.create','records.update','records.delete','auth.login','auth.logout','auth.refresh','users.create','users.update','users.delete','roles.create','roles.update','roles.delete','settings.update')),
  actor_id   TEXT NOT NULL,
  actor_name TEXT NOT NULL,
  record_id  TEXT,
  detail     TEXT,
  created_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC)`,
}

// Descriptors returns the immutable 0004, 0005 and 0008 operation-log history.
func Descriptors() []kernel.MigrationContribution {
	return []kernel.MigrationContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log"},
			Version:              4,
			Name:                 "operation_log",
			Checksum:             kernel.MigrationChecksum(operationLogDDL, "0004:operation-log:v1"),
			Apply:                migrateOperationLog,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_expand"},
			Version:              5,
			Name:                 "operation_log_expand",
			Checksum:             kernel.MigrationChecksum(operationLogExpandDDL, "0005:operation-log-expand:v1"),
			Apply:                migrateOperationLogExpand,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_settings"},
			Version:              8,
			Name:                 "operation_log_settings",
			Checksum:             kernel.MigrationChecksum(operationLogSettingsDDL, "0008:operation-log-settings:v1"),
			Apply:                migrateOperationLogSettings,
		},
	}
}

func migrateOperationLog(tx *sql.Tx) error {
	for _, stmt := range operationLogDDL {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("create operation_log: %w", err)
		}
	}
	return nil
}

func migrateOperationLogExpand(tx *sql.Tx) error {
	return rebuildOperationLog(tx, operationLogExpandDDL, "expanded")
}

func migrateOperationLogSettings(tx *sql.Tx) error {
	return rebuildOperationLog(tx, operationLogSettingsDDL, "settings-expanded")
}

func rebuildOperationLog(tx *sql.Tx, ddl []string, label string) error {
	if _, err := tx.Exec(`ALTER TABLE operation_log RENAME TO operation_log_old`); err != nil {
		return fmt.Errorf("rename operation_log: %w", err)
	}
	if _, err := tx.Exec(ddl[0]); err != nil {
		return fmt.Errorf("create operation_log %s: %w", label, err)
	}
	if _, err := tx.Exec(
		`INSERT INTO operation_log (id, event, actor_id, actor_name, record_id, detail, created_at)
		 SELECT id, event, actor_id, actor_name, record_id, detail, created_at FROM operation_log_old`,
	); err != nil {
		return fmt.Errorf("migrate operation_log rows: %w", err)
	}
	if _, err := tx.Exec(`DROP TABLE operation_log_old`); err != nil {
		return fmt.Errorf("drop operation_log_old: %w", err)
	}
	if _, err := tx.Exec(ddl[1]); err != nil {
		return fmt.Errorf("create operation_log index: %w", err)
	}
	return nil
}
