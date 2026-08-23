// GOAL-037 / F-008 根治：0050（wallet_ledger_order_repair）迁移测试——既库
// 乱序流水重排、健康库 no-op、坏数据 fail-closed。挂在 store 包（可用
// MigrationCatalog 构造从 v1 连续的全局目录）。
package store

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	walletstore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/wallet/store"
)

func seedTx050(t *testing.T, st *Store, q string, args ...any) {
	t.Helper()
	if err := st.WithTx(context.Background(), func(tx *sql.Tx) error {
		_, err := tx.Exec(q, args...)
		return err
	}); err != nil {
		t.Fatalf("seed: %v", err)
	}
}

// seedDisorderedLedger050 writes an account plus a same-millisecond disordered
// group: the freeze entry sorts BEFORE its funding adjust under legacy
// random-suffix ids, while balance-after snapshots still form the CORRECT
// chain (adjust +1000 → freeze −300 → unfreeze +100).
func seedDisorderedLedger050(t *testing.T, st *Store) {
	t.Helper()
	now := time.Date(2026, 8, 23, 10, 0, 0, 0, time.UTC)
	idA := "000001a02e7a2e1faaaaaaaaaaaaaaaaaaaaaaaa" // sorts first (freeze)
	idB := "000001a02e7a2e1fzzzzzzzzzzzzzzzzzzzzzzzz" // sorts second (adjust)
	idC := "000001a02e7a2e20bbbbbbbbbbbbbbbbbbbbbbbb" // next millisecond (unfreeze)
	seedTx050(t, st, `INSERT INTO wallet_accounts (id, owner_type, owner_id, currency, balance_total, balance_available, balance_frozen, status, version, created_at, updated_at)
		 VALUES ('acct-1', 'user', 'u1', 'CNY', 1000, 800, 200, 'active', 6, 1, 1)`)
	for _, rec := range []struct {
		id     string
		etype  string
		delta  int64
		afterT int64
		afterA int64
		afterF int64
	}{
		{idB, "adjust", 1000, 1000, 1000, 0},
		{idA, "freeze", 300, 1000, 700, 300},
		{idC, "unfreeze", 100, 1000, 800, 200},
	} {
		seedTx050(t, st, `INSERT INTO wallet_ledger_entries (id, account_id, entry_type, amount_delta, balance_after_total, balance_after_available, balance_after_frozen, memo, actor_id, actor_name, created_at)
		 VALUES (?, 'acct-1', ?, ?, ?, ?, ?, 'm', 'a', 'a', ?)`,
			rec.id, rec.etype, rec.delta, rec.afterT, rec.afterA, rec.afterF, now.Unix())
	}
}

// catalogBefore excludes version 50 (and later): explicit version filter is
// safer than slicing by index (catalog[i] = version i+1).
func catalogBefore050(catalog []kernel.MigrationContribution) []kernel.MigrationContribution {
	out := []kernel.MigrationContribution{}
	for _, m := range catalog {
		if m.Version < 50 {
			out = append(out, m)
		}
	}
	return out
}

func TestMigration0050ReordersDisorderedWalletLedger(t *testing.T) {
	path := filepath.Join(t.TempDir(), "repair050.db")
	catalog := MigrationCatalog()
	st, err := OpenWithCatalog(path, catalogBefore050(catalog)) // 1..0049: wallet 0031/0033 applied, 0050 NOT yet
	if err != nil {
		t.Fatalf("open pre-0050: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	seedDisorderedLedger050(t, st)
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	st2, err := OpenWithCatalog(path, catalog) // full: 0050 applies now
	if err != nil {
		t.Fatalf("open with 0050: %v", err)
	}
	defer st2.Close()

	// Guard: the repair migration must actually have been applied.
	applied := 0
	if err := st2.WithTx(context.Background(), func(tx *sql.Tx) error {
		row := tx.QueryRow(`SELECT COUNT(*) FROM schema_migrations WHERE version = 50 AND name = 'wallet_ledger_order_repair'`)
		return row.Scan(&applied)
	}); err != nil {
		t.Fatal(err)
	}
	if applied != 1 {
		t.Fatalf("0050 not applied (schema_migrations rows=%d)", applied)
	}

	repo := walletstore.NewRepository(st2)
	// TEMP diagnosis: dump post-0050 ledger order.
	if err := st2.WithTx(context.Background(), func(tx *sql.Tx) error {
		rows, err := tx.Query(
			`SELECT id, entry_type FROM wallet_ledger_entries WHERE account_id='acct-1' ORDER BY created_at, id`)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var id, etype string
			if err := rows.Scan(&id, &etype); err != nil {
				return err
			}
			t.Logf("POST-0050: %s %s", etype, id)
		}
		return rows.Err()
	}); err != nil {
		t.Fatal(err)
	}
	run, err := repo.ReconcileRun("acct-1", "run-1", "actor", time.Now().UTC())
	if err != nil {
		t.Fatalf("reconcile: %v", err)
	}
	if run.Result != walletstore.ResultConsistent {
		t.Fatalf("reconcile result = %q, want consistent (details=%s)", run.Result, run.Details)
	}

	var adjustID, freezeID string
	if err := st2.WithTx(context.Background(), func(tx *sql.Tx) error {
		rows, err := tx.Query(
			`SELECT id, entry_type FROM wallet_ledger_entries WHERE account_id='acct-1' AND entry_type IN ('adjust','freeze') ORDER BY created_at, id`)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var id, etype string
			if err := rows.Scan(&id, &etype); err != nil {
				return err
			}
			switch etype {
			case "adjust":
				adjustID = id
			case "freeze":
				freezeID = id
			}
		}
		return rows.Err()
	}); err != nil {
		t.Fatal(err)
	}
	if adjustID == "" || freezeID == "" {
		t.Fatalf("missing repaired ids (adjust=%q freeze=%q)", adjustID, freezeID)
	}
	if !(adjustID < freezeID) {
		t.Fatalf("repaired order wrong: adjust id %q should sort before freeze id %q", adjustID, freezeID)
	}
}

func TestMigration0050HealthyLedgerNoOp(t *testing.T) {
	path := filepath.Join(t.TempDir(), "healthy050.db")
	st, err := OpenWithCatalog(path, MigrationCatalog())
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()
	repo := walletstore.NewRepository(st)
	if _, _, err := repo.ListAccounts(walletstore.ListFilter{Page: 1, PageSize: 10}); err != nil {
		t.Fatalf("list accounts: %v", err)
	}
}

func TestMigration0050FailsClosedOnBrokenData(t *testing.T) {
	path := filepath.Join(t.TempDir(), "broken050.db")
	catalog := MigrationCatalog()
	st, err := OpenWithCatalog(path, catalogBefore050(catalog))
	if err != nil {
		t.Fatalf("open pre-0050: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	now := time.Now().UTC()
	seedTx050(t, st, `INSERT INTO wallet_accounts (id, owner_type, owner_id, currency, balance_total, balance_available, balance_frozen, status, version, created_at, updated_at)
		 VALUES ('acct-x', 'user', 'x', 'CNY', 0, 0, 0, 'active', 0, 1, 1)`)
	// Same-millisecond pair whose snapshots cannot match ANY order:
	// adjust +1000 with after=1000 AND freeze −300 with after=(1500,1200,300)
	// (adjust-then-freeze yields (1000,700,300) — snapshot mismatch; freeze
	// first is insufficient) → no valid order → fail-closed.
	seedTx050(t, st, `INSERT INTO wallet_ledger_entries (id, account_id, entry_type, amount_delta, balance_after_total, balance_after_available, balance_after_frozen, memo, actor_id, actor_name, created_at)
		 VALUES ('000001a02e7a2e1faaaaaaaaaaaaaaaaaaaaaaaa', 'acct-x', 'adjust', 1000, 1000, 1000, 0, 'm', 'a', 'a', ?)`, now.Unix())
	seedTx050(t, st, `INSERT INTO wallet_ledger_entries (id, account_id, entry_type, amount_delta, balance_after_total, balance_after_available, balance_after_frozen, memo, actor_id, actor_name, created_at)
		 VALUES ('000001a02e7a2e1fzzzzzzzzzzzzzzzzzzzzzzzz', 'acct-x', 'freeze', 300, 1500, 1200, 300, 'm', 'a', 'a', ?)`, now.Unix())
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := OpenWithCatalog(path, catalog); err == nil {
		t.Fatal("expected fail-closed on broken data, got success")
	}
}