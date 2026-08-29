package migration

import (
	"context"
	"fmt"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
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

// accountAvatarDDL adds the self-service avatar column (W13 T-05): the avatar
// URL points into the account avatar asset store (server-produced rasters);
// "" = no avatar. Same additive ALTER pattern as accountEnableDDL.
var accountAvatarDDL = []string{
	`ALTER TABLE users ADD COLUMN avatar_url TEXT NOT NULL DEFAULT ''`,
}

// Descriptors returns the immutable admin.account migration history:
// 0013 (enable state) + 0035 (avatar url, appended to the global ledger).
func Descriptors() []kernel.MigrationContribution {
	return []kernel.MigrationContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "account_enable_state"},
			Version:              13,
			Name:                 "account_enable_state",
			Checksum:             kernel.MigrationChecksum(accountEnableDDL, "0013:account-enable-state:v1"),
			Apply:                migrateAccountEnable,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "account_avatar_url"},
			Version:              35,
			Name:                 "account_avatar_url",
			Checksum:             kernel.MigrationChecksum(accountAvatarDDL, "0035:account-avatar-url:v1"),
			Apply:                migrateAccountAvatar,
		},
	}
}

func migrateAccountEnable(tx kernel.Tx) error {
	for _, stmt := range accountEnableDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("add users.enabled: %w", err)
		}
	}
	return nil
}

func migrateAccountAvatar(tx kernel.Tx) error {
	for _, stmt := range accountAvatarDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("add users.avatar_url: %w", err)
		}
	}
	return nil
}
