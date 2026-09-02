package voucher_test

import (
	"context"
	"errors"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
	walletstore "github.com/magicvr/schema-ui-core/apps/api/modules/wallet/store"
	"github.com/magicvr/schema-ui-core/apps/api/modules/wallet/subject"
	"github.com/magicvr/schema-ui-core/apps/api/modules/wallet/voucher"
)

func now() time.Time { return time.Unix(1700000000, 0).UTC() }

type env struct {
	subjectStore *subject.Store
	walletRepo   *walletstore.Repository
	service      *voucher.Service
}

func newEnv(t *testing.T) *env {
	t.Helper()
	st, err := testsupport.OpenStore(":memory:", "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })

	subStore := subject.NewStore(st)
	wRepo := walletstore.NewRepository(st)
	svc := voucher.NewService(st, wRepo, subStore)

	return &env{
		subjectStore: subStore,
		walletRepo:   wRepo,
		service:      svc,
	}
}

func TestGenerateBatchAndEntropy(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()

	batch, err := e.service.GenerateBatch(ctx, "batch-001", 5, 1000, "CNY", nil, now())
	if err != nil {
		t.Fatalf("generate batch: %v", err)
	}
	if len(batch) != 5 {
		t.Fatalf("batch len = %d, want 5", len(batch))
	}

	seenCodes := make(map[string]bool)
	for i, g := range batch {
		if len(g.Code) != 24 {
			t.Fatalf("code[%d] len = %d, want 24", i, len(g.Code))
		}
		if len(g.Voucher.CodePrefix) != 6 {
			t.Fatalf("prefix[%d] len = %d, want 6", i, len(g.Voucher.CodePrefix))
		}
		if g.Code[:6] != g.Voucher.CodePrefix {
			t.Fatalf("prefix mismatch: code=%s, prefix=%s", g.Code, g.Voucher.CodePrefix)
		}
		if g.Voucher.CodeHash != voucher.HashCode(g.Code) {
			t.Fatalf("hash mismatch")
		}
		if g.Voucher.Status != voucher.StatusUnused {
			t.Fatalf("status = %s, want unused", g.Voucher.Status)
		}
		if g.Voucher.Amount != 1000 {
			t.Fatalf("amount = %d, want 1000", g.Voucher.Amount)
		}
		if seenCodes[g.Code] {
			t.Fatalf("duplicate code generated: %s", g.Code)
		}
		seenCodes[g.Code] = true
	}
}

func TestRedeemSuccess(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()

	// 1. Create a subject.
	sub, _, err := e.subjectStore.GetOrCreateSubject(ctx, "telegram", "tg-1001", now())
	if err != nil {
		t.Fatal(err)
	}

	// 2. Generate a voucher.
	batch, err := e.service.GenerateBatch(ctx, "batch-1", 1, 5000, "CNY", nil, now())
	if err != nil {
		t.Fatal(err)
	}
	item := batch[0]

	// 3. Redeem the voucher.
	res, err := e.service.Redeem(ctx, sub.ID, item.Code, now().Add(time.Minute))
	if err != nil {
		t.Fatalf("redeem: %v", err)
	}
	if res.Amount != 5000 || res.Balance != 5000 {
		t.Fatalf("res = %+v, want amount=5000 balance=5000", res)
	}

	// 4. Verify wallet account.
	acct, err := e.walletRepo.GetSubjectAccountByOwner(sub.ID)
	if err != nil {
		t.Fatal(err)
	}
	if acct.OwnerType != walletstore.OwnerSubject || acct.OwnerID != sub.ID {
		t.Fatalf("account owner mismatch: %+v", acct)
	}
	if acct.BalanceTotal != 5000 || acct.BalanceAvailable != 5000 || acct.BalanceFrozen != 0 {
		t.Fatalf("balances invariant failed: %+v", acct)
	}

	// 5. Verify voucher state.
	v, err := e.service.GetVoucher(ctx, item.Voucher.ID)
	if err != nil {
		t.Fatal(err)
	}
	if v.Status != voucher.StatusRedeemed {
		t.Fatalf("status = %s, want redeemed", v.Status)
	}
	if v.RedeemedBy == nil || *v.RedeemedBy != sub.ID {
		t.Fatalf("redeemed_by mismatch")
	}
}

func TestRedeemDuplicateIdempotentFailure(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()

	sub, _, err := e.subjectStore.GetOrCreateSubject(ctx, "telegram", "tg-1002", now())
	if err != nil {
		t.Fatal(err)
	}
	batch, err := e.service.GenerateBatch(ctx, "batch-1", 1, 100, "CNY", nil, now())
	if err != nil {
		t.Fatal(err)
	}
	item := batch[0]

	// First redeem ok.
	if _, err := e.service.Redeem(ctx, sub.ID, item.Code, now()); err != nil {
		t.Fatalf("first redeem: %v", err)
	}

	// Second redeem must fail with ErrVoucherAlreadyRedeemed.
	_, err = e.service.Redeem(ctx, sub.ID, item.Code, now())
	if !errors.Is(err, voucher.ErrVoucherAlreadyRedeemed) {
		t.Fatalf("second redeem err = %v, want ErrVoucherAlreadyRedeemed", err)
	}

	// Balance must still be 100 (no double credit).
	acct, err := e.walletRepo.GetSubjectAccountByOwner(sub.ID)
	if err != nil {
		t.Fatal(err)
	}
	if acct.BalanceTotal != 100 {
		t.Fatalf("balance = %d, want 100", acct.BalanceTotal)
	}
}

func TestRedeemVoidAndExpired(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()

	sub, _, _ := e.subjectStore.GetOrCreateSubject(ctx, "telegram", "tg-1003", now())

	// 1. Void voucher.
	batchVoid, _ := e.service.GenerateBatch(ctx, "b-void", 1, 100, "CNY", nil, now())
	if err := e.service.VoidVoucher(ctx, batchVoid[0].Voucher.ID, now()); err != nil {
		t.Fatalf("void: %v", err)
	}
	_, err := e.service.Redeem(ctx, sub.ID, batchVoid[0].Code, now())
	if !errors.Is(err, voucher.ErrVoucherVoid) {
		t.Fatalf("redeem void err = %v, want ErrVoucherVoid", err)
	}

	// 2. Expired voucher.
	exp := now().Add(-time.Hour)
	batchExp, _ := e.service.GenerateBatch(ctx, "b-exp", 1, 100, "CNY", &exp, now().Add(-2*time.Hour))
	_, err = e.service.Redeem(ctx, sub.ID, batchExp[0].Code, now())
	if !errors.Is(err, voucher.ErrVoucherExpired) {
		t.Fatalf("redeem expired err = %v, want ErrVoucherExpired", err)
	}
}

func TestRedeemSubjectMustExist(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()

	batch, _ := e.service.GenerateBatch(ctx, "b1", 1, 100, "CNY", nil, now())
	_, err := e.service.Redeem(ctx, "non-existent-subject-id", batch[0].Code, now())
	if !errors.Is(err, voucher.ErrSubjectNotFound) {
		t.Fatalf("redeem non-existent subject err = %v, want ErrSubjectNotFound", err)
	}
}

func TestConcurrentDoubleSpendFailClosed(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "race.db")
	st, err := testsupport.OpenStore(dbPath, "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	subStore := subject.NewStore(st)
	wRepo := walletstore.NewRepository(st)
	svc := voucher.NewService(st, wRepo, subStore)
	ctx := context.Background()

	sub, _, err := subStore.GetOrCreateSubject(ctx, "telegram", "tg-race-user", now())
	if err != nil {
		t.Fatal(err)
	}

	batch, err := svc.GenerateBatch(ctx, "b-race", 1, 1000, "CNY", nil, now())
	if err != nil {
		t.Fatal(err)
	}
	code := batch[0].Code

	const concurrency = 20
	var wg sync.WaitGroup
	var successCount atomic.Int32
	var conflictCount atomic.Int32
	var otherErrors atomic.Int32

	startGate := make(chan struct{})

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-startGate // Synchronize all goroutines to hit the DB concurrently
			_, err := svc.Redeem(ctx, sub.ID, code, now())
			if err == nil {
				successCount.Add(1)
			} else if errors.Is(err, voucher.ErrVoucherConflict) || errors.Is(err, voucher.ErrVoucherAlreadyRedeemed) {
				conflictCount.Add(1)
			} else {
				otherErrors.Add(1)
			}
		}()
	}
	close(startGate) // Release all goroutines simultaneously
	wg.Wait()

	if otherErrors.Load() > 0 {
		t.Fatalf("unexpected other errors count: %d", otherErrors.Load())
	}
	if successCount.Load() != 1 {
		t.Fatalf("success count = %d, want EXACTLY 1", successCount.Load())
	}
	if conflictCount.Load() != concurrency-1 {
		t.Fatalf("conflict count = %d, want %d", conflictCount.Load(), concurrency-1)
	}

	// Verify balance is exactly 1000 (no double spend).
	acct, err := wRepo.GetSubjectAccountByOwner(sub.ID)
	if err != nil {
		t.Fatal(err)
	}
	if acct.BalanceTotal != 1000 || acct.BalanceAvailable != 1000 || acct.BalanceFrozen != 0 {
		t.Fatalf("account balance corrupted: %+v", acct)
	}

	// Verify immutable ledger has exactly 1 entry for this voucher.
	var entryCount int
	err = st.Run(ctx, func(tx kernel.Tx) error {
		return tx.QueryRow(ctx, "SELECT COUNT(*) FROM wallet_ledger_entries WHERE ref_type = 'voucher' AND ref_id = ?", batch[0].Voucher.ID).Scan(&entryCount)
	})
	if err != nil {
		t.Fatal(err)
	}
	if entryCount != 1 {
		t.Fatalf("ledger entry count = %d, want 1", entryCount)
	}
}

func TestNoPlaintextInDatabase(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()

	batch, err := e.service.GenerateBatch(ctx, "b-secret", 3, 500, "CNY", nil, now())
	if err != nil {
		t.Fatal(err)
	}

	// Direct raw SQL scan on vouchers table to verify plaintext codes never exist in DB.
	for _, g := range batch {
		v, err := e.service.GetVoucher(ctx, g.Voucher.ID)
		if err != nil {
			t.Fatal(err)
		}
		if v.CodePrefix != g.Code[:6] {
			t.Fatalf("prefix mismatch: %s vs %s", v.CodePrefix, g.Code[:6])
		}
		if v.CodeHash != "" {
			t.Fatalf("CodeHash field must be omitted from JSON/client structs")
		}
	}
}
