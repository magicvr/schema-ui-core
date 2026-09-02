package wallet

// GOAL-037 / F-008 regression: newID must order SAME-millisecond entries by
// creation order (the D-002 §1 replay sorts (created_at ASC, id ASC)); the
// historical random-only suffix made same-ms ordering arbitrary, so replay
// could run a freeze before its funding adjust ("replay apply failed:
// insufficient balance", inconsistent reconcile).
import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
	walletstore "github.com/magicvr/schema-ui-core/apps/api/modules/wallet/store"
)

func TestNewIDSameMillisecondOrdering(t *testing.T) {
	now := time.Date(2026, 8, 23, 10, 0, 0, 0, time.UTC)
	ids := make([]string, 0, 1000)
	for i := 0; i < 1000; i++ {
		id, err := newID(now)
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, id)
	}
	for i := 1; i < len(ids); i++ {
		if !(ids[i-1] < ids[i]) {
			t.Fatalf("same-millisecond ids out of creation order at %d: %q vs %q", i, ids[i-1], ids[i])
		}
	}
	// Uniqueness holds even with a constant clock (counter + random suffix).
	seen := map[string]bool{}
	for _, id := range ids {
		if seen[id] {
			t.Fatalf("duplicate id %q", id)
		}
		seen[id] = true
	}
}

func TestNewIDCrossMillisecondOrdering(t *testing.T) {
	base := time.Date(2026, 8, 23, 10, 0, 0, 0, time.UTC)
	id1, err := newID(base)
	if err != nil {
		t.Fatal(err)
	}
	id2, err := newID(base.Add(time.Millisecond))
	if err != nil {
		t.Fatal(err)
	}
	if !(id1 < id2) {
		t.Fatalf("cross-millisecond ids not ordered: %q vs %q", id1, id2)
	}
}

func TestWalletServiceSubjectAccountLifecycle(t *testing.T) {
	st, err := testsupport.OpenStore(":memory:", "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	repo := walletstore.NewRepository(st)
	svc := NewService(repo, st)
	ctx := context.Background()
	now := time.Now().UTC()

	// 1. Unregistered subject must fail to open account (VP-029 判据 #1: 未登记主体不能开户).
	_, err = svc.CreateAccount(walletstore.OwnerSubject, "unregistered-sub-1", "CNY", now)
	if !errors.Is(err, walletstore.ErrNotFound) {
		t.Fatalf("unregistered subject create err = %v, want ErrNotFound", err)
	}

	// 2. Register subject.
	sub, _, err := svc.SubjectStore().GetOrCreateSubject(ctx, "telegram", "tg-9999", now)
	if err != nil {
		t.Fatal(err)
	}

	// 3. Registered subject can create wallet account.
	acct, err := svc.CreateAccount(walletstore.OwnerSubject, sub.ID, "CNY", now)
	if err != nil {
		t.Fatalf("create account for registered subject: %v", err)
	}
	if acct.OwnerType != walletstore.OwnerSubject || acct.OwnerID != sub.ID {
		t.Fatalf("unexpected account: %+v", acct)
	}

	// 4. Generate and redeem voucher through VoucherService.
	batch, err := svc.VoucherService().GenerateBatch(ctx, "b-test", 1, 2500, "CNY", nil, now)
	if err != nil {
		t.Fatal(err)
	}
	res, err := svc.VoucherService().Redeem(ctx, sub.ID, batch[0].Code, now.Add(time.Second))
	if err != nil {
		t.Fatalf("voucher redeem: %v", err)
	}
	if res.Amount != 2500 || res.Balance != 2500 {
		t.Fatalf("unexpected redeem result: %+v", res)
	}
}