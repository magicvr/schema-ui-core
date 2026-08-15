// Package migration owns the admin.wallet schema (S-14 · GOAL-019 D-002
// §4): wallet accounts with the three-balance invariant, the immutable ledger
// entries (apply-table snapshots) and the reconciliation-run log.
package migration

import (
	"database/sql"
	"fmt"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// ModuleID is the S-14 wallet module owner.
const ModuleID = "admin.wallet"

// walletDDL (0031): account rows carry the invariant
// balance_total = balance_available + balance_frozen (CHECK + reconciliation
// chain), an optimistic-lock version and a (owner_type, owner_id, currency)
// uniqueness. Ledger entries are append-only with balance-after snapshots and
// a composite (account_id, idempotency_key) uniqueness (D-002 v1.1.0 §1).
var walletDDL = []string{
	`CREATE TABLE wallet_accounts (
  id                TEXT PRIMARY KEY,
  owner_type        TEXT NOT NULL CHECK (owner_type IN ('user','business','system')),
  owner_id          TEXT NOT NULL,
  currency          TEXT NOT NULL DEFAULT 'CNY',
  balance_total     INTEGER NOT NULL DEFAULT 0 CHECK (balance_total >= 0),
  balance_available INTEGER NOT NULL DEFAULT 0 CHECK (balance_available >= 0),
  balance_frozen    INTEGER NOT NULL DEFAULT 0 CHECK (balance_frozen >= 0),
  status            TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','disabled')),
  version           INTEGER NOT NULL DEFAULT 0,
  created_at        INTEGER NOT NULL,
  updated_at        INTEGER NOT NULL,
  UNIQUE (owner_type, owner_id, currency),
  CHECK (balance_total = balance_available + balance_frozen)
)`,
	`CREATE TABLE wallet_ledger_entries (
  id                      TEXT PRIMARY KEY,
  account_id              TEXT NOT NULL,
  entry_type              TEXT NOT NULL CHECK (entry_type IN ('adjust','freeze','unfreeze')),
  amount_delta            INTEGER NOT NULL CHECK (amount_delta != 0),
  balance_after_total     INTEGER NOT NULL CHECK (balance_after_total >= 0),
  balance_after_available INTEGER NOT NULL CHECK (balance_after_available >= 0),
  balance_after_frozen    INTEGER NOT NULL CHECK (balance_after_frozen >= 0),
  ref_type                TEXT,
  ref_id                  TEXT,
  idempotency_key         TEXT,
  memo                    TEXT NOT NULL,
  actor_id                TEXT NOT NULL,
  actor_name              TEXT NOT NULL,
  created_at              INTEGER NOT NULL,
  UNIQUE (account_id, idempotency_key),
  CHECK (balance_after_total = balance_after_available + balance_after_frozen)
)`,
	`CREATE INDEX idx_wallet_ledger_account ON wallet_ledger_entries(account_id, created_at DESC)`,
	`CREATE TABLE wallet_reconciliation_runs (
  id             TEXT PRIMARY KEY,
  account_id     TEXT,
  result         TEXT NOT NULL CHECK (result IN ('consistent','inconsistent')),
  mismatch_count INTEGER NOT NULL DEFAULT 0,
  details        TEXT NOT NULL DEFAULT '{}',
  actor_id       TEXT NOT NULL,
  created_at     INTEGER NOT NULL
)`,
}

// Descriptors returns the immutable 0031 wallet history.
func Descriptors() []kernel.MigrationContribution {
	return []kernel.MigrationContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "wallet"},
			Version:              31,
			Name:                 "wallet",
			Checksum:             kernel.MigrationChecksum(walletDDL, "0031:wallet:v1"),
			Apply:                migrateWallet,
		},
	}
}

func migrateWallet(tx *sql.Tx) error {
	for _, stmt := range walletDDL {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("create wallet tables: %w", err)
		}
	}
	return nil
}
