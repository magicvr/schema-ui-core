package migration

import (
	"context"
	"fmt"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

var telegramConfigDDL = []string{
	`CREATE TABLE telegram_config (
  id                 INTEGER PRIMARY KEY CHECK (id = 1),
  bot_token_enc      TEXT    NOT NULL DEFAULT '',
  webhook_secret_enc TEXT    NOT NULL DEFAULT '',
  updated_at         INTEGER NOT NULL DEFAULT 0
)`,
}

var telegramConfigPGDDL = []string{
	`CREATE TABLE telegram_config (
  id                 INTEGER PRIMARY KEY CHECK (id = 1),
  bot_token_enc      TEXT    NOT NULL DEFAULT '',
  webhook_secret_enc TEXT    NOT NULL DEFAULT '',
  updated_at         BIGINT  NOT NULL DEFAULT 0
)`,
}

var telegramConfigConnectionDDL = []string{
	`ALTER TABLE telegram_config ADD COLUMN mode TEXT NOT NULL DEFAULT 'polling'`,
	`ALTER TABLE telegram_config ADD COLUMN webhook_public_base_url TEXT NOT NULL DEFAULT ''`,
}

func migrateTelegramConfig(tx kernel.Tx) error {
	for _, stmt := range telegramConfigDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create telegram_config: %w", err)
		}
	}
	return nil
}

func migrateTelegramConfigPG(tx kernel.Tx) error {
	for _, stmt := range telegramConfigPGDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create telegram_config (postgres): %w", err)
		}
	}
	return nil
}

func migrateTelegramConfigConnection(tx kernel.Tx) error {
	for _, stmt := range telegramConfigConnectionDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("extend telegram_config: %w", err)
		}
	}
	return nil
}

func migrateTelegramConfigConnectionPG(tx kernel.Tx) error {
	for _, stmt := range telegramConfigConnectionDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("extend telegram_config (postgres): %w", err)
		}
	}
	return nil
}

// Descriptors returns the compiled-global migrations for telegram_config
// (v66) and its connection settings extension (v67).
func Descriptors() ([]kernel.MigrationContribution, error) {
	configChecksum := kernel.MigrationChecksum(telegramConfigDDL, "0066:telegram-config:v1")
	connectionChecksum := kernel.MigrationChecksum(telegramConfigConnectionDDL, "0067:telegram-config-connection:v1")
	return []kernel.MigrationContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "telegram_config"},
			Version:              66,
			Name:                 "telegram_config",
			Checksum:             configChecksum,
			Apply:                migrateTelegramConfig,
			ApplyPostgres:        migrateTelegramConfigPG,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "telegram_config_connection"},
			Version:              67,
			Name:                 "telegram_config_connection",
			Checksum:             connectionChecksum,
			Apply:                migrateTelegramConfigConnection,
			ApplyPostgres:        migrateTelegramConfigConnectionPG,
		},
	}, nil
}
