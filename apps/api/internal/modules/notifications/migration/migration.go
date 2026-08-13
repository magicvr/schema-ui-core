package migration

import (
	"database/sql"
	"fmt"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// ModuleID is the F-04 notification module owner (GOAL-006 · workspace-011).
const ModuleID = "admin.notifications"

// notificationsDDL (0016): per-user in-app notification rows. read_at is NULL
// while unread; created_at drives ordering and the per-user pruning cap.
var notificationsDDL = []string{
	`CREATE TABLE notifications (
  id         TEXT PRIMARY KEY,
  user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  event      TEXT NOT NULL CHECK (event IN ('account.locked','account.disabled','account.unlocked','account.password-changed')),
  title      TEXT NOT NULL,
  body       TEXT NOT NULL,
  read_at    INTEGER,
  created_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_notifications_user_created ON notifications(user_id, created_at DESC)`,
}

// notificationsEnabledDDL (0017): per-user in-app notification master switch.
// 0 = new notifications are not produced for this user (settings API).
var notificationsEnabledDDL = []string{
	`ALTER TABLE users ADD COLUMN notifications_enabled INTEGER NOT NULL DEFAULT 1`,
}

// Descriptors returns the immutable 0016-0017 migration history.
func Descriptors() []kernel.MigrationContribution {
	return []kernel.MigrationContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "notifications"},
			Version:              16,
			Name:                 "notifications",
			Checksum:             kernel.MigrationChecksum(notificationsDDL, "0016:notifications:v1"),
			Apply:                migrateNotifications,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "notifications_enabled"},
			Version:              17,
			Name:                 "notifications_enabled",
			Checksum:             kernel.MigrationChecksum(notificationsEnabledDDL, "0017:notifications-enabled:v1"),
			Apply:                migrateNotificationsEnabled,
		},
	}
}

func migrateNotifications(tx *sql.Tx) error {
	for _, stmt := range notificationsDDL {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("create notifications: %w", err)
		}
	}
	return nil
}

func migrateNotificationsEnabled(tx *sql.Tx) error {
	for _, stmt := range notificationsEnabledDDL {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("add users.notifications_enabled: %w", err)
		}
	}
	return nil
}
