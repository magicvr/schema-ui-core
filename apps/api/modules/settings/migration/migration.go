// Package migration owns the admin.settings migration 0007 (site_settings).
// R6 C6.2 slice 3: the DDL + Apply move out of the store into the owning
// module; the store ledger references this package for the descriptor, and the
// module's CompiledPersistence returns it. Checksums are computed with the
// shared kernel.MigrationChecksum so the ledger stays byte-compatible.
package migration

import (
	"context"
	"fmt"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

const ModuleID = "admin.settings"

// DefaultSiteTitle is the frozen default branding title (matches the checked-in
// manifest; moved here with 0007 so the migration is self-contained).
const DefaultSiteTitle = "Schema UI Core"

// siteSettingsDDL is the singleton branding table (GOAL-013): site title +
// optional logo URL text (upload plugin not in scope). It is the canonical
// checksum input for 0007.
var siteSettingsDDL = []string{
	`CREATE TABLE site_settings (
  id         TEXT PRIMARY KEY CHECK (id = 'default'),
  site_title TEXT NOT NULL,
  logo_url   TEXT NOT NULL DEFAULT '',
  updated_at INTEGER NOT NULL
)`,
}

// siteSettingsPGDDL is the postgres variant of siteSettingsDDL: updated_at
// (Unix time) is BIGINT (R1 v1.4 §3).
var siteSettingsPGDDL = []string{
	`CREATE TABLE site_settings (
  id         TEXT PRIMARY KEY CHECK (id = 'default'),
  site_title TEXT NOT NULL,
  logo_url   TEXT NOT NULL DEFAULT '',
  updated_at BIGINT NOT NULL
)`,
}

// migrate0007PG is the postgres variant of migrate0007 (BIGINT updated_at).
func migrate0007PG(tx kernel.Tx) error {
	for _, stmt := range siteSettingsPGDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create site_settings (postgres): %w", err)
		}
	}
	now := time.Now().UTC().Unix()
	if _, err := tx.Exec(context.Background(),
		`INSERT INTO site_settings (id, site_title, logo_url, updated_at) VALUES ('default', ?, '', ?)`,
		DefaultSiteTitle, now,
	); err != nil {
		return fmt.Errorf("seed site_settings (postgres): %w", err)
	}
	return nil
}

// migrate0007 creates site_settings and seeds the default singleton row.
func migrate0007(tx kernel.Tx) error {
	for _, stmt := range siteSettingsDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create site_settings: %w", err)
		}
	}
	now := time.Now().UTC().Unix()
	if _, err := tx.Exec(context.Background(),
		`INSERT INTO site_settings (id, site_title, logo_url, updated_at) VALUES ('default', ?, '', ?)`,
		DefaultSiteTitle, now,
	); err != nil {
		return fmt.Errorf("seed site_settings: %w", err)
	}
	return nil
}

const transformID = "0007:site-settings:v1"

// siteSettingsV2DDL extends the singleton with the VP-007 system-settings
// fields: light/dark logo + favicon URLs, default locale, default timezone and
// default theme. All new columns are TEXT NOT NULL DEFAULT ” so existing rows
// migrate in place; empty strings carry "unset" semantics.
var siteSettingsV2DDL = []string{
	`ALTER TABLE site_settings ADD COLUMN logo_url_light TEXT NOT NULL DEFAULT ''`,
	`ALTER TABLE site_settings ADD COLUMN logo_url_dark TEXT NOT NULL DEFAULT ''`,
	`ALTER TABLE site_settings ADD COLUMN favicon_url TEXT NOT NULL DEFAULT ''`,
	`ALTER TABLE site_settings ADD COLUMN default_locale TEXT NOT NULL DEFAULT ''`,
	`ALTER TABLE site_settings ADD COLUMN site_timezone TEXT NOT NULL DEFAULT ''`,
	`ALTER TABLE site_settings ADD COLUMN default_theme TEXT NOT NULL DEFAULT ''`,
}

const transformV2ID = "0010:site-settings:v2"

// DefaultOperationLogRetentionDays is the admin-editable default (90 days).
const DefaultOperationLogRetentionDays = 90

// DefaultOperationLogExpirationAction archives expired rows instead of deleting.
const DefaultOperationLogExpirationAction = "archive"

const (
	MinOperationLogRetentionDays = 1
	MaxOperationLogRetentionDays = 3650
	ExpirationActionArchive      = "archive"
	ExpirationActionDelete       = "delete"
)

// migrate0010 applies the VP-007 settings column extension.
func migrate0010(tx kernel.Tx) error {
	for _, stmt := range siteSettingsV2DDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("extend site_settings: %w", err)
		}
	}
	return nil
}

// siteFooterDDL (0040 · W16-F10): optional footer copyright and ICP number.
var siteFooterDDL = []string{
	`ALTER TABLE site_settings ADD COLUMN copyright_text TEXT NOT NULL DEFAULT ''`,
	`ALTER TABLE site_settings ADD COLUMN icp_number TEXT NOT NULL DEFAULT ''`,
}

// migrate0040 applies the footer settings column extension.
func migrate0040(tx kernel.Tx) error {
	for _, stmt := range siteFooterDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("extend site_settings footer: %w", err)
		}
	}
	return nil
}

// Descriptors returns the admin.settings migration descriptors (R6 C6.2).
func Descriptors() []kernel.MigrationContribution {
	return []kernel.MigrationContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "site_settings"},
			Version:              7,
			Name:                 "site_settings",
			Checksum:             kernel.MigrationChecksum(siteSettingsDDL, transformID),
			Apply:                migrate0007,
			ApplyPostgres:        migrate0007PG,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "site_settings_v2"},
			Version:              10,
			Name:                 "site_settings_v2",
			Checksum:             kernel.MigrationChecksum(siteSettingsV2DDL, transformV2ID),
			Apply:                migrate0010,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "site_footer"},
			Version:              40,
			Name:                 "site_footer",
			Checksum:             kernel.MigrationChecksum(siteFooterDDL, "0040:site-footer:v1"),
			Apply:                migrate0040,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "site_operation_log_retention"},
			Version:              46,
			Name:                 "site_operation_log_retention",
			Checksum:             kernel.MigrationChecksum(siteOperationLogRetentionDDL, "0046:site-operation-log-retention:v1"),
			Apply:                migrate0046,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "site_default_currency"},
			Version:              62,
			Name:                 "site_default_currency",
			Checksum:             kernel.MigrationChecksum(siteCurrencyDDL, "0062:site-default-currency:v1"),
			Apply:                migrate0062,
		},
		{
			// 0063 · R4 零冲突升级演练样本（GOAL-005 S2）：site_settings.updated_at 索引。
			// 双方言同一语法（CREATE INDEX IF NOT EXISTS），单一 Apply。
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "site_settings_updated_at_index"},
			Version:              63,
			Name:                 "site_settings_updated_at_index",
			Checksum:             kernel.MigrationChecksum(siteUpdatedAtIndexDDL, transformV0050ID),
			Apply:                migrate0063,
		},
	}
}

// siteUpdatedAtIndexDDL (0063 · R4 演练样本)。SQLite 与 PostgreSQL 均支持
// CREATE INDEX IF NOT EXISTS——单一 Apply 覆盖双方言。
var siteUpdatedAtIndexDDL = []string{
	`CREATE INDEX IF NOT EXISTS idx_site_settings_updated_at ON site_settings (updated_at)`,
}

const transformV0050ID = "0063:site-settings:updated-at-index:v1"

func migrate0063(tx kernel.Tx) error {
	for _, stmt := range siteUpdatedAtIndexDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("index site_settings updated_at: %w", err)
		}
	}
	return nil
}

// siteCurrencyDDL (0062 · workspace-020 R3): site-wide default currency
// (ISO 4217 code; part of the format-semantics contract GOAL-002 D-001 §4.1).
// Empty string keeps the "unset" semantics consistent with locale/timezone.
var siteCurrencyDDL = []string{
	`ALTER TABLE site_settings ADD COLUMN default_currency TEXT NOT NULL DEFAULT ''`,
}

// migrate0062 applies the default-currency column extension.
func migrate0062(tx kernel.Tx) error {
	for _, stmt := range siteCurrencyDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("extend site_settings currency: %w", err)
		}
	}
	return nil
}

// siteOperationLogRetentionDDL (0046): admin-editable audit log retention.
var siteOperationLogRetentionDDL = []string{
	`ALTER TABLE site_settings ADD COLUMN operation_log_retention_days INTEGER NOT NULL DEFAULT 90`,
	`ALTER TABLE site_settings ADD COLUMN operation_log_expiration_action TEXT NOT NULL DEFAULT 'archive'`,
}

func migrate0046(tx kernel.Tx) error {
	for _, stmt := range siteOperationLogRetentionDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("extend site_settings retention: %w", err)
		}
	}
	return nil
}
