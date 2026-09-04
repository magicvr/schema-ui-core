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

var telegramIngressDDL = []string{
	`CREATE TABLE telegram_sessions (
  bot_id          INTEGER NOT NULL,
  chat_id         INTEGER NOT NULL,
  chat_type       TEXT    NOT NULL DEFAULT '',
  title           TEXT    NOT NULL DEFAULT '',
  username        TEXT    NOT NULL DEFAULT '',
  last_message_at INTEGER NOT NULL,
  created_at      INTEGER NOT NULL,
  updated_at      INTEGER NOT NULL,
  PRIMARY KEY (bot_id, chat_id)
)`,
	`CREATE INDEX idx_telegram_sessions_activity
  ON telegram_sessions (bot_id, last_message_at DESC, chat_id DESC)`,
	`CREATE TABLE telegram_inbound_messages (
  bot_id            INTEGER NOT NULL,
  update_id         INTEGER NOT NULL,
  chat_id           INTEGER NOT NULL,
  user_id           INTEGER,
  message_id        INTEGER,
  callback_query_id TEXT,
  direction         TEXT    NOT NULL DEFAULT 'inbound' CHECK (direction = 'inbound'),
  message_kind      TEXT    NOT NULL,
  text              TEXT,
  callback_data     TEXT,
  sender_username   TEXT,
  received_at       INTEGER NOT NULL,
  PRIMARY KEY (bot_id, update_id)
)`,
	`CREATE INDEX idx_telegram_inbound_messages_chat_received
  ON telegram_inbound_messages (bot_id, chat_id, received_at DESC, update_id DESC)`,
}

var telegramIngressPGDDL = []string{
	`CREATE TABLE telegram_sessions (
  bot_id          BIGINT NOT NULL,
  chat_id         BIGINT NOT NULL,
  chat_type       TEXT   NOT NULL DEFAULT '',
  title           TEXT   NOT NULL DEFAULT '',
  username        TEXT   NOT NULL DEFAULT '',
  last_message_at BIGINT NOT NULL,
  created_at      BIGINT NOT NULL,
  updated_at      BIGINT NOT NULL,
  PRIMARY KEY (bot_id, chat_id)
)`,
	`CREATE INDEX idx_telegram_sessions_activity
  ON telegram_sessions (bot_id, last_message_at DESC, chat_id DESC)`,
	`CREATE TABLE telegram_inbound_messages (
  bot_id            BIGINT NOT NULL,
  update_id         BIGINT NOT NULL,
  chat_id           BIGINT NOT NULL,
  user_id           BIGINT,
  message_id        BIGINT,
  callback_query_id TEXT,
  direction         TEXT   NOT NULL DEFAULT 'inbound' CHECK (direction = 'inbound'),
  message_kind      TEXT   NOT NULL,
  text              TEXT,
  callback_data     TEXT,
  sender_username   TEXT,
  received_at       BIGINT NOT NULL,
  PRIMARY KEY (bot_id, update_id)
)`,
	`CREATE INDEX idx_telegram_inbound_messages_chat_received
  ON telegram_inbound_messages (bot_id, chat_id, received_at DESC, update_id DESC)`,
}

var telegramOutboundDDL = []string{
	`CREATE TABLE telegram_outbound_messages (
  bot_id        INTEGER NOT NULL,
  request_id    TEXT    NOT NULL,
  retry_root    TEXT    NOT NULL,
  retry_of      TEXT,
  chat_id       INTEGER NOT NULL,
  text          TEXT    NOT NULL,
  status        TEXT    NOT NULL CHECK (status IN ('pending', 'sent', 'failed')),
  error_message TEXT,
  created_at    INTEGER NOT NULL,
  updated_at    INTEGER NOT NULL,
  PRIMARY KEY (bot_id, request_id)
 )`,
	`CREATE INDEX idx_telegram_outbound_messages_chat_created
  ON telegram_outbound_messages (bot_id, chat_id, created_at DESC, request_id DESC)`,
	`CREATE UNIQUE INDEX idx_telegram_outbound_messages_pending_root
  ON telegram_outbound_messages (bot_id, retry_root) WHERE status = 'pending'`,
}

var telegramOutboundPGDDL = []string{
	`CREATE TABLE telegram_outbound_messages (
  bot_id        BIGINT NOT NULL,
  request_id    TEXT   NOT NULL,
  retry_root    TEXT   NOT NULL,
  retry_of      TEXT,
  chat_id       BIGINT NOT NULL,
  text          TEXT   NOT NULL,
  status        TEXT   NOT NULL CHECK (status IN ('pending', 'sent', 'failed')),
  error_message TEXT,
  created_at    BIGINT NOT NULL,
  updated_at    BIGINT NOT NULL,
  PRIMARY KEY (bot_id, request_id)
 )`,
	`CREATE INDEX idx_telegram_outbound_messages_chat_created
  ON telegram_outbound_messages (bot_id, chat_id, created_at DESC, request_id DESC)`,
	`CREATE UNIQUE INDEX idx_telegram_outbound_messages_pending_root
  ON telegram_outbound_messages (bot_id, retry_root) WHERE status = 'pending'`,
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

func migrateTelegramIngress(tx kernel.Tx) error {
	for _, stmt := range telegramIngressDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create telegram ingress tables: %w", err)
		}
	}
	return nil
}

func migrateTelegramIngressPG(tx kernel.Tx) error {
	for _, stmt := range telegramIngressPGDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create telegram ingress tables (postgres): %w", err)
		}
	}
	return nil
}

func migrateTelegramOutbound(tx kernel.Tx) error {
	for _, stmt := range telegramOutboundDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create telegram outbound table: %w", err)
		}
	}
	return nil
}

func migrateTelegramOutboundPG(tx kernel.Tx) error {
	for _, stmt := range telegramOutboundPGDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create telegram outbound table (postgres): %w", err)
		}
	}
	return nil
}

// Descriptors returns the compiled-global migrations for telegram_config
// (v66), its connection settings extension (v67), the R3 ingress tables
// (v68), and the R3 operator outbound ledger (v69).
func Descriptors() ([]kernel.MigrationContribution, error) {
	configChecksum := kernel.MigrationChecksum(telegramConfigDDL, "0066:telegram-config:v1")
	connectionChecksum := kernel.MigrationChecksum(telegramConfigConnectionDDL, "0067:telegram-config-connection:v1")
	ingressChecksum := kernel.MigrationChecksum(telegramIngressDDL, "0068:telegram-ingress:v1")
	outboundChecksum := kernel.MigrationChecksum(telegramOutboundDDL, "0069:telegram-outbound:v1")
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
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "telegram_ingress"},
			Version:              68,
			Name:                 "telegram_ingress",
			Checksum:             ingressChecksum,
			Apply:                migrateTelegramIngress,
			ApplyPostgres:        migrateTelegramIngressPG,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "telegram_outbound"},
			Version:              69,
			Name:                 "telegram_outbound",
			Checksum:             outboundChecksum,
			Apply:                migrateTelegramOutbound,
			ApplyPostgres:        migrateTelegramOutboundPG,
		},
	}, nil
}
