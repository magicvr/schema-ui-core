// Package migration owns the admin.login-captcha schema (S-11 · GOAL-011
// D-002 `4): one-time arithmetic challenges plus the single-row enable flag.
package migration

import (
	"database/sql"
	"fmt"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// ModuleID is the S-11 captcha module owner.
const ModuleID = "admin.login-captcha"

// captchaDDL (0023): challenge rows (answer stored hashed; one-time use via
// delete on verify) and the single-row config switch (default disabled —
// D-001 `5).
var captchaDDL = []string{
	`CREATE TABLE captcha_challenges (
  id          TEXT PRIMARY KEY,
  answer_hash TEXT NOT NULL,
  expires_at  INTEGER NOT NULL,
  created_at  INTEGER NOT NULL
)`,
	`CREATE TABLE captcha_config (
  id         INTEGER PRIMARY KEY CHECK (id = 1),
  enabled    INTEGER NOT NULL DEFAULT 0,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
)`,
}

// Descriptors returns the immutable 0023 captcha history.
func Descriptors() []kernel.MigrationContribution {
	return []kernel.MigrationContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "login_captcha"},
			Version:              23,
			Name:                 "login_captcha",
			Checksum:             kernel.MigrationChecksum(captchaDDL, "0023:login-captcha:v1"),
			Apply:                migrateCaptcha,
		},
	}
}

func migrateCaptcha(tx *sql.Tx) error {
	for _, stmt := range captchaDDL {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("create captcha tables: %w", err)
		}
	}
	return nil
}
