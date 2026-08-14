// Package migration owns the admin.data-dictionary schema (S-01 · GOAL-008
// D-002 `2): the two-level dict_types / dict_entries tables. Deleting a type
// cascades to its entries (documented decision, D-002 `7).
package migration

import (
	"database/sql"
	"fmt"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// ModuleID is the S-01 dictionary module owner.
const ModuleID = "admin.data-dictionary"

// dictDDL (0019): dictionary types and entries. key is the stable type
// identifier referenced by entries; UNIQUE(dict_key, entry_key) pins the
// two-level enumeration identity.
var dictDDL = []string{
	`CREATE TABLE dict_types (
  id         TEXT PRIMARY KEY,
  key        TEXT NOT NULL UNIQUE,
  name       TEXT NOT NULL,
  enabled    INTEGER NOT NULL DEFAULT 1,
  description TEXT,
  sort       INTEGER NOT NULL DEFAULT 0,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
)`,
	`CREATE TABLE dict_entries (
  id         TEXT PRIMARY KEY,
  dict_key   TEXT NOT NULL REFERENCES dict_types(key) ON DELETE CASCADE,
  entry_key  TEXT NOT NULL,
  label      TEXT NOT NULL,
  enabled    INTEGER NOT NULL DEFAULT 1,
  sort       INTEGER NOT NULL DEFAULT 0,
  remark     TEXT,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL,
  UNIQUE (dict_key, entry_key)
)`,
	`CREATE INDEX idx_dict_entries_dict_key ON dict_entries(dict_key, sort)`,
}

// Descriptors returns the immutable 0019 dictionary history.
func Descriptors() []kernel.MigrationContribution {
	return []kernel.MigrationContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "dictionary"},
			Version:              19,
			Name:                 "dictionary",
			Checksum:             kernel.MigrationChecksum(dictDDL, "0019:dictionary:v1"),
			Apply:                migrateDict,
		},
	}
}

func migrateDict(tx *sql.Tx) error {
	for _, stmt := range dictDDL {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("create dictionary tables: %w", err)
		}
	}
	return nil
}
