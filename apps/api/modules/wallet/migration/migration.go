// Package migration owns the admin.wallet schema (S-14 · GOAL-019 D-002
// §4): wallet accounts with the three-balance invariant, the immutable ledger
// entries (apply-table snapshots) and the reconciliation-run log.
package migration

import (
	"context"
	"fmt"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
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

// walletPGDDL is the postgres variant of walletDDL: money columns
// (balance_* / amount_delta / balance_after_*) and Unix time columns
// (created_at / updated_at) are BIGINT (R1 v1.4 §3; wallet money exceeds int4).
var walletPGDDL = []string{
	`CREATE TABLE wallet_accounts (
  id                TEXT PRIMARY KEY,
  owner_type        TEXT NOT NULL CHECK (owner_type IN ('user','business','system')),
  owner_id          TEXT NOT NULL,
  currency          TEXT NOT NULL DEFAULT 'CNY',
  balance_total     BIGINT NOT NULL DEFAULT 0 CHECK (balance_total >= 0),
  balance_available BIGINT NOT NULL DEFAULT 0 CHECK (balance_available >= 0),
  balance_frozen    BIGINT NOT NULL DEFAULT 0 CHECK (balance_frozen >= 0),
  status            TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','disabled')),
  version           INTEGER NOT NULL DEFAULT 0,
  created_at        BIGINT NOT NULL,
  updated_at        BIGINT NOT NULL,
  UNIQUE (owner_type, owner_id, currency),
  CHECK (balance_total = balance_available + balance_frozen)
)`,
	`CREATE TABLE wallet_ledger_entries (
  id                      TEXT PRIMARY KEY,
  account_id              TEXT NOT NULL,
  entry_type              TEXT NOT NULL CHECK (entry_type IN ('adjust','freeze','unfreeze')),
  amount_delta            BIGINT NOT NULL CHECK (amount_delta != 0),
  balance_after_total     BIGINT NOT NULL CHECK (balance_after_total >= 0),
  balance_after_available BIGINT NOT NULL CHECK (balance_after_available >= 0),
  balance_after_frozen    BIGINT NOT NULL CHECK (balance_after_frozen >= 0),
  ref_type                TEXT,
  ref_id                  TEXT,
  idempotency_key         TEXT,
  memo                    TEXT NOT NULL,
  actor_id                TEXT NOT NULL,
  actor_name              TEXT NOT NULL,
  created_at              BIGINT NOT NULL,
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
  created_at     BIGINT NOT NULL
)`,
}

// walletLedgerDeductPGDDL is the postgres variant of walletLedgerDeductDDL:
// BIGINT money/time columns.
var walletLedgerDeductPGDDL = []string{
	`CREATE TABLE wallet_ledger_entries (
  id                      TEXT PRIMARY KEY,
  account_id              TEXT NOT NULL,
  entry_type              TEXT NOT NULL CHECK (entry_type IN ('adjust','freeze','unfreeze','deduct_frozen')),
  amount_delta            BIGINT NOT NULL CHECK (amount_delta != 0),
  balance_after_total     BIGINT NOT NULL CHECK (balance_after_total >= 0),
  balance_after_available BIGINT NOT NULL CHECK (balance_after_available >= 0),
  balance_after_frozen    BIGINT NOT NULL CHECK (balance_after_frozen >= 0),
  ref_type                TEXT,
  ref_id                  TEXT,
  idempotency_key         TEXT,
  memo                    TEXT NOT NULL,
  actor_id                TEXT NOT NULL,
  actor_name              TEXT NOT NULL,
  created_at              BIGINT NOT NULL,
  UNIQUE (account_id, idempotency_key),
  CHECK (balance_after_total = balance_after_available + balance_after_frozen)
)`,
	`CREATE INDEX idx_wallet_ledger_account ON wallet_ledger_entries(account_id, created_at DESC)`,
}

// walletLedgerDeductDDL (0033 · GOAL-021 D-001 §3): the ledger entry_type CHECK
// gains 'deduct_frozen' (consume from the frozen bucket). SQLite cannot alter a
// CHECK, so the table is rebuilt (rename → create → copy → drop) like the
// operationlog rebuild pattern; data and constraints are preserved verbatim.
var walletLedgerDeductDDL = []string{
	`CREATE TABLE wallet_ledger_entries (
  id                      TEXT PRIMARY KEY,
  account_id              TEXT NOT NULL,
  entry_type              TEXT NOT NULL CHECK (entry_type IN ('adjust','freeze','unfreeze','deduct_frozen')),
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
}

func migrateWalletLedgerDeduct(tx kernel.Tx) error {
	return walletLedgerDeduct(rebuildCommon{tx: tx, ddl: walletLedgerDeductDDL})
}

// migrateWalletLedgerDeductPG is the postgres variant of migrateWalletLedgerDeduct.
func migrateWalletLedgerDeductPG(tx kernel.Tx) error {
	return walletLedgerDeduct(rebuildCommon{tx: tx, ddl: walletLedgerDeductPGDDL})
}

type rebuildCommon struct {
	tx  kernel.Tx
	ddl []string
}

func walletLedgerDeduct(r rebuildCommon) error {
	tx, ddl := r.tx, r.ddl
	if _, err := tx.Exec(context.Background(), `ALTER TABLE wallet_ledger_entries RENAME TO wallet_ledger_entries_old`); err != nil {
		return fmt.Errorf("rename wallet_ledger_entries: %w", err)
	}
	if _, err := tx.Exec(context.Background(), ddl[0]); err != nil {
		return fmt.Errorf("recreate wallet_ledger_entries: %w", err)
	}
	if _, err := tx.Exec(context.Background(),
		`INSERT INTO wallet_ledger_entries (id, account_id, entry_type, amount_delta, balance_after_total, balance_after_available, balance_after_frozen, ref_type, ref_id, idempotency_key, memo, actor_id, actor_name, created_at)
		 SELECT id, account_id, entry_type, amount_delta, balance_after_total, balance_after_available, balance_after_frozen, ref_type, ref_id, idempotency_key, memo, actor_id, actor_name, created_at FROM wallet_ledger_entries_old`,
	); err != nil {
		return fmt.Errorf("migrate wallet_ledger_entries rows: %w", err)
	}
	if _, err := tx.Exec(context.Background(), `DROP TABLE wallet_ledger_entries_old`); err != nil {
		return fmt.Errorf("drop wallet_ledger_entries_old: %w", err)
	}
	if _, err := tx.Exec(context.Background(), ddl[1]); err != nil {
		return fmt.Errorf("recreate wallet ledger index: %w", err)
	}
	return nil
}

// Descriptors returns the immutable 0031 + 0033 wallet history.
func Descriptors() []kernel.MigrationContribution {
	return []kernel.MigrationContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "wallet"},
			Version:              31,
			Name:                 "wallet",
			Checksum:             kernel.MigrationChecksum(walletDDL, "0031:wallet:v1"),
			Apply:                migrateWallet,
			ApplyPostgres:        migrateWalletPG,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "wallet_ledger_deduct"},
			Version:              33,
			Name:                 "wallet_ledger_deduct",
			Checksum:             kernel.MigrationChecksum(walletLedgerDeductDDL, "0033:wallet-ledger-deduct:v1"),
			Apply:                migrateWalletLedgerDeduct,
			ApplyPostgres:        migrateWalletLedgerDeductPG,
		},
		{
			// 0050 · GOAL-037 / F-008 根治：一次性重排既有库中"同一毫秒
			// 乱序"的流水 id，恢复 D-002 §1 的 (created_at, id) 回放序契约。
			// 逻辑为方言中立 Go（order_repair.go），无 DDL。
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: orderRepairKey},
			Version:              orderRepairVersion,
			Name:                 orderRepairKey,
			Checksum:             kernel.MigrationChecksum(nil, "0050:wallet-ledger-order-repair:v1"),
			Apply:                migrateOrderRepair,
			ApplyPostgres:        migrateOrderRepairPG,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "wallet_voucher_and_subject"},
			Version:              64,
			Name:                 "wallet_voucher_and_subject",
			Checksum:             kernel.MigrationChecksum(walletVoucherAndSubjectDDL, "0064:wallet-voucher-and-subject:v1"),
			Apply:                migrateWalletVoucherAndSubject,
			ApplyPostgres:        migrateWalletVoucherAndSubjectPG,
		},
	}
}

func migrateWallet(tx kernel.Tx) error {
	for _, stmt := range walletDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create wallet tables: %w", err)
		}
	}
	return nil
}

func migrateWalletPG(tx kernel.Tx) error {
	for _, stmt := range walletPGDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create wallet tables (postgres): %w", err)
		}
	}
	return nil
}

// walletVoucherAndSubjectDDL (0064): establishes external subject mapping (subjects)
// and prepaid instruments (vouchers), plus extends wallet_accounts owner_type CHECK
// with 'subject'.
var walletVoucherAndSubjectDDL = []string{
	`CREATE TABLE subjects (
  id          TEXT PRIMARY KEY,
  issuer      TEXT NOT NULL,
  external_id TEXT NOT NULL,
  created_at  INTEGER NOT NULL,
  UNIQUE (issuer, external_id)
)`,
	`CREATE INDEX idx_subjects_issuer_external ON subjects(issuer, external_id)`,
	`CREATE TABLE vouchers (
  id           TEXT PRIMARY KEY,
  batch_id     TEXT NOT NULL,
  code_hash    TEXT NOT NULL,
  code_prefix  TEXT NOT NULL,
  amount       INTEGER NOT NULL CHECK (amount > 0),
  currency     TEXT NOT NULL DEFAULT 'CNY',
  status       TEXT NOT NULL DEFAULT 'unused' CHECK (status IN ('unused','redeemed','void')),
  expires_at   INTEGER,
  redeemed_by  TEXT,
  redeemed_at  INTEGER,
  created_at   INTEGER NOT NULL,
  updated_at   INTEGER NOT NULL,
  UNIQUE (code_hash)
)`,
	`CREATE INDEX idx_vouchers_batch ON vouchers(batch_id, created_at DESC)`,
	`CREATE INDEX idx_vouchers_status ON vouchers(status)`,
	`CREATE TABLE wallet_accounts (
  id                TEXT PRIMARY KEY,
  owner_type        TEXT NOT NULL CHECK (owner_type IN ('user','business','system','subject')),
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
}

var walletVoucherAndSubjectPGDDL = []string{
	`CREATE TABLE subjects (
  id          TEXT PRIMARY KEY,
  issuer      TEXT NOT NULL,
  external_id TEXT NOT NULL,
  created_at  BIGINT NOT NULL,
  UNIQUE (issuer, external_id)
)`,
	`CREATE INDEX idx_subjects_issuer_external ON subjects(issuer, external_id)`,
	`CREATE TABLE vouchers (
  id           TEXT PRIMARY KEY,
  batch_id     TEXT NOT NULL,
  code_hash    TEXT NOT NULL,
  code_prefix  TEXT NOT NULL,
  amount       BIGINT NOT NULL CHECK (amount > 0),
  currency     TEXT NOT NULL DEFAULT 'CNY',
  status       TEXT NOT NULL DEFAULT 'unused' CHECK (status IN ('unused','redeemed','void')),
  expires_at   BIGINT,
  redeemed_by  TEXT,
  redeemed_at  BIGINT,
  created_at   BIGINT NOT NULL,
  updated_at   BIGINT NOT NULL,
  UNIQUE (code_hash)
)`,
	`CREATE INDEX idx_vouchers_batch ON vouchers(batch_id, created_at DESC)`,
	`CREATE INDEX idx_vouchers_status ON vouchers(status)`,
}

func migrateWalletVoucherAndSubject(tx kernel.Tx) error {
	for i := 0; i < 5; i++ {
		if _, err := tx.Exec(context.Background(), walletVoucherAndSubjectDDL[i]); err != nil {
			return fmt.Errorf("create subject and voucher tables: %w", err)
		}
	}
	// Rebuild wallet_accounts to include 'subject' in owner_type CHECK
	if _, err := tx.Exec(context.Background(), `ALTER TABLE wallet_accounts RENAME TO wallet_accounts_old`); err != nil {
		return fmt.Errorf("rename wallet_accounts: %w", err)
	}
	if _, err := tx.Exec(context.Background(), walletVoucherAndSubjectDDL[5]); err != nil {
		return fmt.Errorf("recreate wallet_accounts: %w", err)
	}
	copySQL := `INSERT INTO wallet_accounts (id, owner_type, owner_id, currency, balance_total, balance_available, balance_frozen, status, version, created_at, updated_at)
SELECT id, owner_type, owner_id, currency, balance_total, balance_available, balance_frozen, status, version, created_at, updated_at FROM wallet_accounts_old`
	if _, err := tx.Exec(context.Background(), copySQL); err != nil {
		return fmt.Errorf("copy wallet_accounts rows: %w", err)
	}
	if _, err := tx.Exec(context.Background(), `DROP TABLE wallet_accounts_old`); err != nil {
		return fmt.Errorf("drop wallet_accounts_old: %w", err)
	}
	return nil
}

func migrateWalletVoucherAndSubjectPG(tx kernel.Tx) error {
	for _, stmt := range walletVoucherAndSubjectPGDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create subject and voucher tables (postgres): %w", err)
		}
	}
	if _, err := tx.Exec(context.Background(), `ALTER TABLE wallet_accounts DROP CONSTRAINT IF EXISTS wallet_accounts_owner_type_check`); err != nil {
		return fmt.Errorf("drop wallet_accounts owner_type check: %w", err)
	}
	if _, err := tx.Exec(context.Background(), `ALTER TABLE wallet_accounts ADD CONSTRAINT wallet_accounts_owner_type_check CHECK (owner_type IN ('user','business','system','subject'))`); err != nil {
		return fmt.Errorf("add wallet_accounts owner_type check: %w", err)
	}
	return nil
}
