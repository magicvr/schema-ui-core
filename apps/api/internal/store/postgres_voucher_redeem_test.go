package store

// VP-029 A-005 F-003（A-008 处理）：PostgreSQL 上的真实 Redeem 与并发首次开户
// 端到端——此前并发双花只在文件 SQLite 上证明，PG 侧仅有开户 ON CONFLICT 片段。
// 本测试在 PG 上跑：两张不同卡、同一新主体并发核销 → 单主体户、余额=两笔面额
// 之和、账本恰 2 条；重复核销同码 fail-closed 且不双记。
// 由 PG_TEST_* 环境门控（无 PG 时 skip，CI postgres job 点亮）。

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	walletstore "github.com/magicvr/schema-ui-core/apps/api/modules/wallet/store"
	"github.com/magicvr/schema-ui-core/apps/api/modules/wallet/subject"
	"github.com/magicvr/schema-ui-core/apps/api/modules/wallet/voucher"
)

func TestPostgresVoucherRedeemAndConcurrentSubject(t *testing.T) {
	st := postgresScratchDB(t, "vp029redeem") // skip when PG_TEST_* unset
	ctx := context.Background()
	now := time.Now().UTC()

	subStore := subject.NewStore(st)
	wRepo := walletstore.NewRepository(st)
	svc := voucher.NewService(st, wRepo, subStore)

	// 1. Register one brand-new subject (concurrent first-time account open).
	sub, _, err := subStore.GetOrCreateSubject(ctx, "pg-issuer", "ext-redeem-1", now)
	if err != nil {
		t.Fatal(err)
	}

	// 2. Two distinct vouchers (different cards) across two batches.
	batchA, err := svc.GenerateBatch(ctx, "pg-batch-a", 1, 500, "CNY", nil, now)
	if err != nil {
		t.Fatal(err)
	}
	batchB, err := svc.GenerateBatch(ctx, "pg-batch-b", 1, 1500, "CNY", nil, now)
	if err != nil {
		t.Fatal(err)
	}
	codes := []string{batchA[0].Code, batchB[0].Code}

	// 3. Redeem BOTH concurrently against the SAME new subject on PostgreSQL.
	errCh := make(chan error, len(codes))
	var wg sync.WaitGroup
	for _, code := range codes {
		wg.Add(1)
		go func(c string) {
			defer wg.Done()
			_, err := svc.Redeem(ctx, sub.ID, c, now.Add(2*time.Second))
			errCh <- err
		}(code)
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			t.Fatalf("concurrent PG redeem: %v", err)
		}
	}

	// 4. Exactly ONE subject account row, balance = 500 + 1500.
	acct, err := wRepo.GetSubjectAccountByOwner(sub.ID)
	if err != nil {
		t.Fatal(err)
	}
	if acct.BalanceTotal != 2000 || acct.BalanceAvailable != 2000 || acct.BalanceFrozen != 0 {
		t.Fatalf("balance after concurrent redeems = %+v, want total/available 2000", acct)
	}
	var accountRows int
	if err := st.Run(ctx, func(tx kernel.Tx) error {
		return tx.QueryRow(ctx,
			`SELECT COUNT(*) FROM wallet_accounts WHERE owner_type = ? AND owner_id = ?`,
			"subject", sub.ID,
		).Scan(&accountRows)
	}); err != nil {
		t.Fatal(err)
	}
	if accountRows != 1 {
		t.Fatalf("subject account rows = %d, want 1", accountRows)
	}

	// 5. Ledger carries exactly two voucher credits (no double-credit).
	var entryRows int
	if err := st.Run(ctx, func(tx kernel.Tx) error {
		return tx.QueryRow(ctx,
			`SELECT COUNT(*) FROM wallet_ledger_entries WHERE account_id = ? AND ref_type = 'voucher'`,
			acct.ID,
		).Scan(&entryRows)
	}); err != nil {
		t.Fatal(err)
	}
	if entryRows != 2 {
		t.Fatalf("voucher ledger entries = %d, want 2", entryRows)
	}

	// 6. Repeat redeem of the same code fails closed and does not double-credit.
	if _, err := svc.Redeem(ctx, sub.ID, codes[0], now.Add(3*time.Second)); !errors.Is(err, voucher.ErrVoucherAlreadyRedeemed) {
		t.Fatalf("duplicate redeem err = %v, want ErrVoucherAlreadyRedeemed", err)
	}
	acct2, err := wRepo.GetSubjectAccountByOwner(sub.ID)
	if err != nil {
		t.Fatal(err)
	}
	if acct2.BalanceTotal != 2000 {
		t.Fatalf("balance after duplicate redeem = %d, want 2000 (no double credit)", acct2.BalanceTotal)
	}
}
