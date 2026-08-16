// GOAL-021 F-001（A-002 recommended）：0033 重建必须保留既有流水行。
// 手工构造 0031 版 ledger 表（旧 CHECK）+ 两行流水，执行 0033 重建 SQL 后
// 断言行数/数据保留/新 CHECK 接受 deduct_frozen。
package store

import (
	"database/sql"
	"path/filepath"
	"testing"
)

func TestMigrate0033PreservesWalletLedgerRows(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ledger-old.db")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// 0031-era schema: entry_type CHECK without deduct_frozen.
	createLegacy := `CREATE TABLE wallet_ledger_entries (` +
		`  id                      TEXT PRIMARY KEY,` +
		`  account_id              TEXT NOT NULL,` +
		`  entry_type              TEXT NOT NULL CHECK (entry_type IN ('adjust','freeze','unfreeze')),` +
		`  amount_delta            INTEGER NOT NULL CHECK (amount_delta != 0),` +
		`  balance_after_total     INTEGER NOT NULL CHECK (balance_after_total >= 0),` +
		`  balance_after_available INTEGER NOT NULL CHECK (balance_after_available >= 0),` +
		`  balance_after_frozen    INTEGER NOT NULL CHECK (balance_after_frozen >= 0),` +
		`  ref_type                TEXT,` +
		`  ref_id                  TEXT,` +
		`  idempotency_key         TEXT,` +
		`  memo                    TEXT NOT NULL,` +
		`  actor_id                TEXT NOT NULL,` +
		`  actor_name              TEXT NOT NULL,` +
		`  created_at              INTEGER NOT NULL,` +
		`  UNIQUE (account_id, idempotency_key),` +
		`  CHECK (balance_after_total = balance_after_available + balance_after_frozen)` +
		`)`
	if _, err := db.Exec(createLegacy); err != nil {
		t.Fatal(err)
	}
	insert := `INSERT INTO wallet_ledger_entries (id, account_id, entry_type, amount_delta, balance_after_total, balance_after_available, balance_after_frozen, ref_type, ref_id, idempotency_key, memo, actor_id, actor_name, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	for _, r := range [][]any{
		{"e1", "acct-1", "adjust", 1000, 1000, 1000, 0, nil, nil, nil, "grant", "a1", "A", 1},
		{"e2", "acct-1", "freeze", 300, 1000, 700, 300, nil, nil, nil, "hold", "a1", "A", 2},
	} {
		if _, err := db.Exec(insert, r...); err != nil {
			t.Fatalf("insert legacy row: %v", err)
		}
	}

	// Execute the 0033 rebuild inline (rename → recreate → copy → drop).
	createNew := `CREATE TABLE wallet_ledger_entries (` +
		`  id                      TEXT PRIMARY KEY,` +
		`  account_id              TEXT NOT NULL,` +
		`  entry_type              TEXT NOT NULL CHECK (entry_type IN ('adjust','freeze','unfreeze','deduct_frozen')),` +
		`  amount_delta            INTEGER NOT NULL CHECK (amount_delta != 0),` +
		`  balance_after_total     INTEGER NOT NULL CHECK (balance_after_total >= 0),` +
		`  balance_after_available INTEGER NOT NULL CHECK (balance_after_available >= 0),` +
		`  balance_after_frozen    INTEGER NOT NULL CHECK (balance_after_frozen >= 0),` +
		`  ref_type                TEXT,` +
		`  ref_id                  TEXT,` +
		`  idempotency_key         TEXT,` +
		`  memo                    TEXT NOT NULL,` +
		`  actor_id                TEXT NOT NULL,` +
		`  actor_name              TEXT NOT NULL,` +
		`  created_at              INTEGER NOT NULL,` +
		`  UNIQUE (account_id, idempotency_key),` +
		`  CHECK (balance_after_total = balance_after_available + balance_after_frozen)` +
		`)`
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(`ALTER TABLE wallet_ledger_entries RENAME TO wallet_ledger_entries_old`); err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(createNew); err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(`INSERT INTO wallet_ledger_entries (id, account_id, entry_type, amount_delta, balance_after_total, balance_after_available, balance_after_frozen, ref_type, ref_id, idempotency_key, memo, actor_id, actor_name, created_at) SELECT id, account_id, entry_type, amount_delta, balance_after_total, balance_after_available, balance_after_frozen, ref_type, ref_id, idempotency_key, memo, actor_id, actor_name, created_at FROM wallet_ledger_entries_old`); err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(`DROP TABLE wallet_ledger_entries_old`); err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(`CREATE INDEX idx_wallet_ledger_account ON wallet_ledger_entries(account_id, created_at DESC)`); err != nil {
		t.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	// Rows preserved.
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM wallet_ledger_entries`).Scan(&count); err != nil || count != 2 {
		t.Fatalf("rows after 0033 = %d (err %v), want 2", count, err)
	}
	var entryType, memo string
	if err := db.QueryRow(`SELECT entry_type, memo FROM wallet_ledger_entries WHERE id = 'e2'`).Scan(&entryType, &memo); err != nil {
		t.Fatal(err)
	}
	if entryType != "freeze" || memo != "hold" {
		t.Fatalf("legacy row mutated: %s/%s", entryType, memo)
	}

	// New CHECK accepts deduct_frozen.
	if _, err := db.Exec(insert, "e3", "acct-1", "deduct_frozen", 100, 900, 700, 200, nil, nil, nil, "settle", "a1", "A", 3); err != nil {
		t.Fatalf("deduct_frozen rejected after 0033: %v", err)
	}
}
