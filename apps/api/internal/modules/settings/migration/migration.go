// Package migration owns the admin.settings migration 0007 (site_settings).
// R6 C6.2 slice 3: the DDL + Apply move out of the store into the owning
// module; the store ledger references this package for the descriptor, and the
// module's CompiledPersistence returns it. Checksums are computed with the
// shared kernel.MigrationChecksum so the ledger stays byte-compatible.
package migration

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
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

// migrate0007 creates site_settings and seeds the default singleton row.
func migrate0007(tx *sql.Tx) error {
	for _, stmt := range siteSettingsDDL {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("create site_settings: %w", err)
		}
	}
	now := time.Now().UTC().Unix()
	if _, err := tx.Exec(
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
// default theme. All new columns are TEXT NOT NULL DEFAULT '' so existing rows
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
func migrate0010(tx *sql.Tx) error {
	for _, stmt := range siteSettingsV2DDL {
		if _, err := tx.Exec(stmt); err != nil {
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
func migrate0040(tx *sql.Tx) error {
	for _, stmt := range siteFooterDDL {
		if _, err := tx.Exec(stmt); err != nil {
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
	}
}

// siteOperationLogRetentionDDL (0046): admin-editable audit log retention.
var siteOperationLogRetentionDDL = []string{
	`ALTER TABLE site_settings ADD COLUMN operation_log_retention_days INTEGER NOT NULL DEFAULT 90`,
	`ALTER TABLE site_settings ADD COLUMN operation_log_expiration_action TEXT NOT NULL DEFAULT 'archive'`,
}

func migrate0046(tx *sql.Tx) error {
	for _, stmt := range siteOperationLogRetentionDDL {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("extend site_settings retention: %w", err)
		}
	}
	return nil
}
