// Package migration owns the admin.recycle-bin schema (S-12 · GOAL-012
// D-002 §1): the recycle_items snapshot table (deleted rows with JSON payload,
// partial-unique per resource+id while unrestored).
package migration

import (
	"context"
	"fmt"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// ModuleID is the S-12 recycle bin module owner.
const ModuleID = "admin.recycle-bin"

// recycleDDL (0025): snapshot rows; restored_at NULL means the row is still
// restorable. Partial unique index prevents duplicate snapshots for the same
// (resource, resource_id) while the item has not been restored.
var recycleDDL = []string{
	`CREATE TABLE recycle_items (
  id          TEXT PRIMARY KEY,
  resource    TEXT NOT NULL,
  resource_id TEXT NOT NULL,
  payload     TEXT NOT NULL,
  actor_id    TEXT NOT NULL,
  actor_name  TEXT NOT NULL,
  deleted_at  INTEGER NOT NULL,
  restored_at INTEGER
)`,
	`CREATE UNIQUE INDEX idx_recycle_items_active ON recycle_items(resource, resource_id) WHERE restored_at IS NULL`,
	`CREATE INDEX idx_recycle_items_deleted_at ON recycle_items(deleted_at DESC)`,
}

// recyclePGDDL is the postgres variant of recycleDDL: deleted_at/restored_at
// (Unix time) are BIGINT (R1 v1.4 §3); the partial unique index is supported
// by postgres too.
var recyclePGDDL = []string{
	`CREATE TABLE recycle_items (
  id          TEXT PRIMARY KEY,
  resource    TEXT NOT NULL,
  resource_id TEXT NOT NULL,
  payload     TEXT NOT NULL,
  actor_id    TEXT NOT NULL,
  actor_name  TEXT NOT NULL,
  deleted_at  BIGINT NOT NULL,
  restored_at BIGINT
)`,
	`CREATE UNIQUE INDEX idx_recycle_items_active ON recycle_items(resource, resource_id) WHERE restored_at IS NULL`,
	`CREATE INDEX idx_recycle_items_deleted_at ON recycle_items(deleted_at DESC)`,
}

// Descriptors returns the immutable 0025 recycle-bin history.
func Descriptors() []kernel.MigrationContribution {
	return []kernel.MigrationContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "recycle_items"},
			Version:              25,
			Name:                 "recycle_items",
			Checksum:             kernel.MigrationChecksum(recycleDDL, "0025:recycle-items:v1"),
			Apply:                migrateRecycle,
			ApplyPostgres:        migrateRecyclePG,
		},
	}
}

func migrateRecycle(tx kernel.Tx) error {
	for _, stmt := range recycleDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create recycle tables: %w", err)
		}
	}
	return nil
}

func migrateRecyclePG(tx kernel.Tx) error {
	for _, stmt := range recyclePGDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create recycle tables (postgres): %w", err)
		}
	}
	return nil
}
