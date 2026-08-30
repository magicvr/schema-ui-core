// Package migration owns the admin.login-captcha schema (S-11 · GOAL-011
// D-002 `4): one-time arithmetic challenges plus the single-row enable flag.
package migration

import (
	"context"
	"fmt"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
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

// captchaPGDDL is the postgres variant of captchaDDL: Unix time columns
// (expires_at / created_at / updated_at) are BIGINT (R1 v1.4 §3).
var captchaPGDDL = []string{
	`CREATE TABLE captcha_challenges (
  id          TEXT PRIMARY KEY,
  answer_hash TEXT NOT NULL,
  expires_at  BIGINT NOT NULL,
  created_at  BIGINT NOT NULL
)`,
	`CREATE TABLE captcha_config (
  id         INTEGER PRIMARY KEY CHECK (id = 1),
  enabled    INTEGER NOT NULL DEFAULT 0,
  created_at BIGINT NOT NULL,
  updated_at BIGINT NOT NULL
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
			ApplyPostgres:        migrateCaptchaPG,
		},
	}
}

func migrateCaptcha(tx kernel.Tx) error {
	for _, stmt := range captchaDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create captcha tables: %w", err)
		}
	}
	return nil
}

func migrateCaptchaPG(tx kernel.Tx) error {
	for _, stmt := range captchaPGDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create captcha tables (postgres): %w", err)
		}
	}
	return nil
}
