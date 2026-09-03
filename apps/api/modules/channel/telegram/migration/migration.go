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

// Descriptors returns the compiled-global migration for telegram_config (v66).
func Descriptors() ([]kernel.MigrationContribution, error) {
	checksum := kernel.MigrationChecksum(telegramConfigDDL, "0066:telegram-config:v1")
	return []kernel.MigrationContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "telegram_config"},
			Version:              66,
			Name:                 "telegram_config",
			Checksum:             checksum,
			Apply:                migrateTelegramConfig,
			ApplyPostgres:        migrateTelegramConfigPG,
		},
	}, nil
}
