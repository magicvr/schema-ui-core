package migration

import (
	"database/sql"
	"fmt"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// ModuleID is the F-03 account-enable module owner (GOAL-005 · workspace-011).
// The users table is hosted by core.auth-session persistence, but the
// product-state enable semantics (I-011-001 §1 C-01 boundary: 启停属 F-03) are
// owned here; the column lands on users via an additive ALTER.
const ModuleID = "admin.account"

// accountEnableDDL adds the product-state enabled flag (GOAL-005 D-002 §1/§3):
// 1 = enabled (default, backward compatible), 0 = disabled by an admin.
// Disable revokes sessions and bumps token_version at the repository layer;
// this migration only adds the persisted state column.
var accountEnableDDL = []string{
	`ALTER TABLE users ADD COLUMN enabled INTEGER NOT NULL DEFAULT 1`,
}

// Descriptors returns the immutable 0013 migration history for admin.account.
func Descriptors() []kernel.MigrationContribution {
	return []kernel.MigrationContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "account_enable_state"},
			Version:              13,
			Name:                 "account_enable_state",
			Checksum:             kernel.MigrationChecksum(accountEnableDDL, "0013:account-enable-state:v1"),
			Apply:                migrateAccountEnable,
		},
	}
}

func migrateAccountEnable(tx *sql.Tx) error {
	for _, stmt := range accountEnableDDL {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("add users.enabled: %w", err)
		}
	}
	return nil
}
