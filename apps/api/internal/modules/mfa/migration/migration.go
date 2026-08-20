// Package migration owns the admin.mfa schema (S-10 · GOAL-017 D-002 §2):
// per-user TOTP state (pending/active, encrypted secret, recovery hashes,
// last used time step) and one-time login proofs with a failure counter.
package migration

import (
	"context"
	"fmt"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// ModuleID is the S-10 MFA module owner.
const ModuleID = "admin.mfa"

// mfaDDL (0029): user_mfa rows exist only for enrolled users; status
// transitions pending → active on confirm. mfa_proofs are the one-time login
// second-factor proofs (5-minute TTL; fail_count reaching 5 destroys the
// proof — A-004 F-001 response).
var mfaDDL = []string{
	`CREATE TABLE user_mfa (
  user_id                TEXT PRIMARY KEY REFERENCES users(id),
  status                 TEXT NOT NULL CHECK (status IN ('pending','active')),
  totp_secret_ciphertext TEXT NOT NULL,
  recovery_codes_hash    TEXT NOT NULL,
  last_used_step         INTEGER NOT NULL DEFAULT 0,
  created_at             INTEGER NOT NULL,
  updated_at             INTEGER NOT NULL
)`,
	`CREATE TABLE mfa_proofs (
  id         TEXT PRIMARY KEY,
  user_id    TEXT NOT NULL REFERENCES users(id),
  fail_count INTEGER NOT NULL DEFAULT 0,
  expires_at INTEGER NOT NULL,
  created_at INTEGER NOT NULL
)`,
}

// Descriptors returns the immutable 0029 MFA history.
func Descriptors() []kernel.MigrationContribution {
	return []kernel.MigrationContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "user_mfa"},
			Version:              29,
			Name:                 "user_mfa",
			Checksum:             kernel.MigrationChecksum(mfaDDL, "0029:user-mfa:v1"),
			Apply:                migrateMFA,
		},
	}
}

func migrateMFA(tx kernel.Tx) error {
	for _, stmt := range mfaDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create mfa tables: %w", err)
		}
	}
	return nil
}
