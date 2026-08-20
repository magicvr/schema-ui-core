// Package migration owns the admin.data-permission schema (S-09 · GOAL-016
// D-002 §4): per-resource scope policies (owner column + default scope) and
// user × resource scope assignments.
package migration

import (
	"context"
	"fmt"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// ModuleID is the S-09 data-permission module owner.
const ModuleID = "admin.data-permission"

// dataPermissionDDL (0027): scope policies declare the owner column and the
// default scope for a resource; assignments override per user. scope_type is
// constrained to all/self (org is deferred to B-10, I-011-001 §5 / I-004).
var dataPermissionDDL = []string{
	`CREATE TABLE data_scope_policies (
  resource      TEXT PRIMARY KEY,
  owner_column  TEXT NOT NULL,
  default_scope TEXT NOT NULL CHECK (default_scope IN ('all','self')),
  enabled       INTEGER NOT NULL DEFAULT 1,
  updated_at    INTEGER NOT NULL
)`,
	`CREATE TABLE user_data_scopes (
  user_id    TEXT NOT NULL,
  resource   TEXT NOT NULL,
  scope_type TEXT NOT NULL CHECK (scope_type IN ('all','self')),
  updated_at INTEGER NOT NULL,
  PRIMARY KEY (user_id, resource)
)`,
}

// dataPermissionPGDDL is the postgres variant of dataPermissionDDL:
// updated_at (Unix time) is BIGINT (R1 v1.4 §3).
var dataPermissionPGDDL = []string{
	`CREATE TABLE data_scope_policies (
  resource      TEXT PRIMARY KEY,
  owner_column  TEXT NOT NULL,
  default_scope TEXT NOT NULL CHECK (default_scope IN ('all','self')),
  enabled       INTEGER NOT NULL DEFAULT 1,
  updated_at    BIGINT NOT NULL
)`,
	`CREATE TABLE user_data_scopes (
  user_id    TEXT NOT NULL,
  resource   TEXT NOT NULL,
  scope_type TEXT NOT NULL CHECK (scope_type IN ('all','self')),
  updated_at BIGINT NOT NULL,
  PRIMARY KEY (user_id, resource)
)`,
}

// Descriptors returns the immutable 0027 data-permission history.
func Descriptors() []kernel.MigrationContribution {
	return []kernel.MigrationContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "data_permission"},
			Version:              27,
			Name:                 "data_permission",
			Checksum:             kernel.MigrationChecksum(dataPermissionDDL, "0027:data-permission:v1"),
			Apply:                migrateDataPermission,
			ApplyPostgres:        migrateDataPermissionPG,
		},
	}
}

func migrateDataPermission(tx kernel.Tx) error {
	for _, stmt := range dataPermissionDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create data permission tables: %w", err)
		}
	}
	return nil
}

func migrateDataPermissionPG(tx kernel.Tx) error {
	for _, stmt := range dataPermissionPGDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create data permission tables (postgres): %w", err)
		}
	}
	return nil
}
