// GOAL-004 C1/C2 acceptance tests (workspace-031 · GOAL-002 D-002 v1.2.0 §5/
// §6/§8): the unified entitlement Check aggregate, atomic Consume with void
// linearization, the §6 Telegram commands over the Disabled/Capture seams and
// the §8 query bucket.
package service_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/channel/telegram"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/digitaloffer/service"
	"github.com/magicvr/schema-ui-core/apps/api/modules/digitaloffer/store"
)

func TestCheckAggregation(t *testing.T) {
	env := newTestEnv(t)
	subjectID := env.seedSubject("buyer", 1000, "")
	offer := env.seedOffer(store.FormCount, 400, store.StatusOnSale)

	// No rows at all.
	res, err := env.svc.Check(context.Background(), subjectID, offer.ID, time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	if res.Valid || res.Reason != service.ReasonNoEntitlement {
		t.Fatalf("empty check = %+v, want no_entitlement", res)
	}

	// Valid after purchase.
	if _, err := env.svc.Purchase(context.Background(), subjectID, offer.ID, "req-check", time.Now().UTC()); err != nil {
		t.Fatal(err)
	}
	res, err = env.svc.Check(context.Background(), subjectID, offer.ID, time.Now().UTC())
	if err != nil || !res.Valid || res.Reason != service.ReasonNone {
		t.Fatalf("valid check = %+v err %v", res, err)
	}

	// Exhausted after consuming every use.
	ent := seedCountEntitlement(t, env, subjectID, 100) // separate offer, 5 uses
	for i := 0; i < 5; i++ {
		if err := env.svc.Consume(context.Background(), subjectID, ent.OfferID, 1, time.Now().UTC()); err != nil {
			t.Fatalf("consume %d: %v", i, err)
		}
	}
	res, err = env.svc.Check(context.Background(), subjectID, ent.OfferID, time.Now().UTC())
	if err != nil || res.Valid || res.Reason != service.ReasonExhausted {
		t.Fatalf("exhausted check = %+v err %v", res, err)
	}

	// Voided after voiding the only active row of a third offer.
	ent2 := seedCountEntitlement(t, env, subjectID, 100)
	if _, err := env.svc.VoidEntitlement(context.Background(), service.Actor{ID: "a", Name: "A"}, ent2.ID, time.Now().UTC()); err != nil {
		t.Fatal(err)
	}
	res, err = env.svc.Check(context.Background(), subjectID, ent2.OfferID, time.Now().UTC())
	if err != nil || res.Valid || res.Reason != service.ReasonVoided {
		t.Fatalf("voided check = %+v err %v", res, err)
	}
}

func TestCheckDurationExpiry(t *testing.T) {
	env := newTestEnv(t)
	subjectID := env.seedSubject("buyer", 1000, "")
	offer := env.seedOffer(store.FormDuration, 400, store.StatusOnSale)

	purchaseAt := now
	if _, err := env.svc.Purchase(context.Background(), subjectID, offer.ID, "req-dur", purchaseAt); err != nil {
		t.Fatal(err)
	}
	// Inside the window: valid.
	validAt := purchaseAt.Add(time.Duration(offer.DurationSeconds-1) * time.Second)
	if res, err := env.svc.Check(context.Background(), subjectID, offer.ID, validAt); err != nil || !res.Valid {
		t.Fatalf("in-window check = %+v err %v, want valid", res, err)
	}
	// After expiry: expired.
	expiredAt := purchaseAt.Add(time.Duration(offer.DurationSeconds+1) * time.Second)
	res, err := env.svc.Check(context.Background(), subjectID, offer.ID, expiredAt)
	if err != nil || res.Valid || res.Reason != service.ReasonExpired {
		t.Fatalf("expired check = %+v err %v, want expired", res, err)
	}
}

func TestConsumeScenarios(t *testing.T) {
	t.Run("insufficient total is terminal with zero deduction", func(t *testing.T) {
		env := newTestEnv(t)
		subjectID := env.seedSubject("buyer", 1000, "")
		ent := seedCountEntitlement(t, env, subjectID, 100) // 5 uses
		if err := env.svc.Consume(context.Background(), subjectID, ent.OfferID, 6, time.Now().UTC()); !errors.Is(err, store.ErrEntitlementInsufficient) {
			t.Fatalf("err = %v, want ErrEntitlementInsufficient", err)
		}
		got, err := env.svc.GetEntitlement(ent.ID)
		if err != nil || *got.RemainingCount != 5 {
			t.Fatalf("remaining = %+v err %v, want untouched 5", got, err)
		}
	})

	t.Run("multi-row consumption uses oldest first and both rows", func(t *testing.T) {
		env := newTestEnv(t)
		subjectID := env.seedSubject("buyer", 2000, "")
		offer := env.seedOffer(store.FormCount, 100, store.StatusOnSale)
		first, err := env.svc.Purchase(context.Background(), subjectID, offer.ID, env.requestID("req-multi-1"), time.Now().UTC())
		if err != nil {
			t.Fatal(err)
		}
		second, err := env.svc.Purchase(context.Background(), subjectID, offer.ID, env.requestID("req-multi-2"), time.Now().UTC())
		if err != nil {
			t.Fatal(err)
		}
		if err := env.svc.Consume(context.Background(), subjectID, offer.ID, 6, time.Now().UTC()); err != nil {
			t.Fatalf("consume 6: %v", err)
		}
		row1, _ := env.svc.GetEntitlement(first.Entitlement.ID)
		row2, _ := env.svc.GetEntitlement(second.Entitlement.ID)
		if *row1.RemainingCount != 0 || *row2.RemainingCount != 4 {
			t.Fatalf("remaining = %d/%d, want 0/4 (oldest first)", *row1.RemainingCount, *row2.RemainingCount)
		}
	})

	t.Run("concurrent void and consume linearize without overspend", func(t *testing.T) {
		env := newTestEnv(t)
		subjectID := env.seedSubject("buyer", 1000, "")
		ent := seedCountEntitlement(t, env, subjectID, 100) // 5 uses

		var wg sync.WaitGroup
		wg.Add(2)
		var consumeErr, voidErr error
		go func() {
			defer wg.Done()
			consumeErr = env.svc.Consume(context.Background(), subjectID, ent.OfferID, 5, time.Now().UTC())
		}()
		go func() {
			defer wg.Done()
			_, voidErr = env.svc.VoidEntitlement(context.Background(), service.Actor{ID: "a", Name: "A"}, ent.ID, time.Now().UTC())
		}()
		wg.Wait()
		if voidErr != nil {
			t.Fatalf("void: %v", voidErr)
		}
		got, err := env.svc.GetEntitlement(ent.ID)
		if err != nil {
			t.Fatal(err)
		}
		switch {
		case consumeErr == nil:
			// Consume won: 0 remaining, then voided.
			if *got.RemainingCount != 0 || got.Status != store.EntitlementVoided {
				t.Fatalf("consume-then-void state = %+v", got)
			}
		case errors.Is(consumeErr, store.ErrEntitlementInsufficient):
			// Void won: row voided with 5 remaining.
			if got.Status != store.EntitlementVoided || *got.RemainingCount != 5 {
				t.Fatalf("void-then-consume state = %+v", got)
			}
		default:
			t.Fatalf("consume err = %v", consumeErr)
		}
	})

	t.Run("duration entitlements never consume", func(t *testing.T) {
		env := newTestEnv(t)
		subjectID := env.seedSubject("buyer", 1000, "")
		offer := env.seedOffer(store.FormDuration, 400, store.StatusOnSale)
		if _, err := env.svc.Purchase(context.Background(), subjectID, offer.ID, "req-dur2", time.Now().UTC()); err != nil {
			t.Fatal(err)
		}
		if err := env.svc.Consume(context.Background(), subjectID, offer.ID, 1, time.Now().UTC()); !errors.Is(err, store.ErrEntitlementInsufficient) {
			t.Fatalf("duration consume err = %v, want ErrEntitlementInsufficient", err)
		}
	})
}

// dispatchOnce drives one constructed update through the production dispatcher
// path (handlers receive the update; replies land in the capture sender that
// was registered via RegisterTelegram).
func dispatchOnce(t *testing.T, dispatcher *telegram.Dispatcher, sender *telegram.CaptureSender, text, subjectID string, updateID int64) {
	t.Helper()
	upd := kernel.TelegramUpdate{
		ChatID:    "12345",
		UserID:    "777",
		SubjectID: subjectID,
		Command:   strings.TrimPrefix(strings.Fields(text)[0], "/"),
		Text:      text,
		UpdateID:  updateID,
	}
	if err := dispatcher.Dispatch(context.Background(), upd, sender); err != nil {
		t.Fatalf("dispatch %q: %v", text, err)
	}
}

// TestTelegramCommandsDisabled proves §6: with the DisabledDispatcher the
// registration is a no-op success and the module stays Bot-API independent.
func TestTelegramCommandsDisabled(t *testing.T) {
	env := newTestEnv(t)
	if err := env.svc.RegisterTelegram(telegram.NewDisabledDispatcher(), telegram.NewDisabledSender()); err != nil {
		t.Fatalf("disabled registration: %v", err)
	}
}

// TestTelegramCommandsReply proves §6: price lists on-sale offers, buy
// purchases with the tg: idempotent request id and fails closed without a
// mapped subject, entitlements lists only valid rows.
func TestTelegramCommandsReply(t *testing.T) {
	env := newTestEnv(t)
	capture := telegram.NewCaptureSender()
	dispatcher := telegram.NewDispatcher()
	if err := env.svc.RegisterTelegram(dispatcher, capture); err != nil {
		t.Fatal(err)
	}
	subjectID := env.seedSubject("tg-buyer", 1000, "")
	offerA := env.seedOffer(store.FormCount, 400, store.StatusOnSale)
	env.seedOffer(store.FormDuration, 900, store.StatusOnSale)
	env.seedOffer(store.FormCount, 50, store.StatusOffSale)

	// price: on-sale items only, off-sale excluded.
	dispatchOnce(t, dispatcher, capture, "price", "", 1)
	sent := capture.Last()
	if sent == nil || !strings.Contains(sent.Text, fmt.Sprintf("id: %s", offerA.ID)) || strings.Contains(sent.Text, "offer-count-off_sale") {
		t.Fatalf("price reply = %+v", sent)
	}

	// buy without a mapped subject fails closed.
	dispatchOnce(t, dispatcher, capture, "/buy "+offerA.ID, "", 2)
	if last := capture.Last(); !strings.Contains(last.Text, "身份") {
		t.Fatalf("anonymous buy reply = %s", last.Text)
	}

	// buy with a subject purchases exactly once; a duplicate update replays.
	dispatchOnce(t, dispatcher, capture, "/buy "+offerA.ID, subjectID, 4242)
	if last := capture.Last(); !strings.Contains(last.Text, "购买成功") {
		t.Fatalf("buy reply = %s", last.Text)
	}
	dispatchOnce(t, dispatcher, capture, "/buy "+offerA.ID, subjectID, 4242)
	if last := capture.Last(); !strings.Contains(last.Text, "幂等重放") {
		t.Fatalf("duplicate buy reply = %s", last.Text)
	}
	if acct := env.accountOf(subjectID); acct.BalanceTotal != 600 {
		t.Fatalf("balance = %d, want 600 (one deduction only)", acct.BalanceTotal)
	}

	// unknown offer replies the reason text.
	dispatchOnce(t, dispatcher, capture, "/buy no-such-offer", subjectID, 3)
	if last := capture.Last(); !strings.Contains(last.Text, "未找到") {
		t.Fatalf("unknown offer reply = %s", last.Text)
	}

	// entitlements: the count entitlement shows remaining uses.
	dispatchOnce(t, dispatcher, capture, "entitlements", subjectID, 4)
	if last := capture.Last(); !strings.Contains(last.Text, "剩余 5 次") {
		t.Fatalf("entitlements reply = %s", last.Text)
	}
}

// TestTelegramPriceRateLimited proves §8: the price/entitlements query bucket
// (bizoffer|price|<subject_id>, 1 min / 30) throttles queries per subject.
func TestTelegramPriceRateLimited(t *testing.T) {
	env := newTestEnv(t) // limiter wired
	capture := telegram.NewCaptureSender()
	dispatcher := telegram.NewDispatcher()
	if err := env.svc.RegisterTelegram(dispatcher, capture); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 30; i++ {
		dispatchOnce(t, dispatcher, capture, "price", "", int64(i))
	}
	dispatchOnce(t, dispatcher, capture, "price", "", 1000)
	if last := capture.Last(); !strings.Contains(last.Text, "过于频繁") {
		t.Fatalf("exhausted query reply = %s", last.Text)
	}
}

var seedSeq atomic.Int64

// seedCountEntitlement drives a purchase so the subject holds count uses.
func seedCountEntitlement(t *testing.T, env *testEnv, subjectID string, price int64) *store.Entitlement {
	t.Helper()
	offer := env.seedOffer(store.FormCount, price, store.StatusOnSale)
	reqID := fmt.Sprintf("req-%s-%d-%d", subjectID, price, seedSeq.Add(1))
	res, err := env.svc.Purchase(context.Background(), subjectID, offer.ID, reqID, time.Now().UTC())
	if err != nil {
		t.Fatalf("purchase: %v", err)
	}
	return res.Entitlement
}
