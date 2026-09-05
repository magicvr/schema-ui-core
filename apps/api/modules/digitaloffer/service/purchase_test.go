// GOAL-003 C3 acceptance tests (workspace-031 · GOAL-002 D-002 v1.1.0 §4):
// the single-transaction purchase path — freeze → deduct_frozen → purchase
// voucher → entitlement in ONE kernel.Tx, idempotent replay, bounded retry on
// races, retry-exhaustion zero-residue and the §8 purchase bucket. The
// scenario matrix runs against SQLite always and against real PostgreSQL when
// PG_TEST_* is configured (A-002 F-002 dual-database demand).
package service_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/pgtest"
	"github.com/magicvr/schema-ui-core/apps/api/internal/ratelimit"
	storepkg "github.com/magicvr/schema-ui-core/apps/api/internal/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/compiled"
	"github.com/magicvr/schema-ui-core/apps/api/modules/digitaloffer/service"
	"github.com/magicvr/schema-ui-core/apps/api/modules/digitaloffer/store"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
	walletstore "github.com/magicvr/schema-ui-core/apps/api/modules/wallet/store"
	"github.com/magicvr/schema-ui-core/apps/api/modules/wallet/subject"
)

var now = time.Unix(1700000000, 0).UTC()

type testEnv struct {
	t          *testing.T
	st         kernel.Store
	salt       string
	wallet     *walletstore.Repository
	subjects   *subject.Store
	svc        *service.Service
	operations *stubRecorder
}

// stubRecorder captures fail-closed audit rows (D-002 §7).
type stubRecorder struct {
	mu   sync.Mutex
	ops  []string
	fail error
}

func (s *stubRecorder) RecordOperationTx(_ kernel.Tx, op operationlog.Operation) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.fail != nil {
		return s.fail
	}
	s.ops = append(s.ops, op.Event)
	return nil
}

func (s *stubRecorder) events() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.ops...)
}

func (s *stubRecorder) reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ops = nil
}

func (s *stubRecorder) setFail(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.fail = err
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "digitaloffer.db"), "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return newTestEnvOn(t, st, true)
}

// newTestEnvOn builds the service over an explicit store. withLimiter wires
// the §8 purchase bucket (production shape and dedicated limiter scenarios);
// money-path matrix scenarios run without it so request counting does not
// intersect transaction semantics.
func newTestEnvOn(t *testing.T, st kernel.Store, withLimiter bool) *testEnv {
	t.Helper()
	rec := &stubRecorder{}
	var limiters kernel.RateLimiterProvider
	if withLimiter {
		limiters = ratelimit.NewProvider()
	}
	return &testEnv{
		t:          t,
		st:         st,
		wallet:     walletstore.NewRepository(st),
		subjects:   subject.NewStore(st),
		svc:        service.NewService(store.NewRepository(st), walletstore.NewRepository(st), subject.NewStore(st), rec, limiters),
		operations: rec,
	}
}

// qualify namespaces seeds and request ids so scenario runs sharing one
// PostgreSQL store never collide on (issuer, external_id) or request guards.
func (e *testEnv) qualify(name string) string {
	if e.salt == "" {
		return name
	}
	return name + "-" + e.salt
}

func (e *testEnv) requestID(base string) string { return e.qualify(base) }

// seedSubject registers one external subject and funds its wallet account.
func (e *testEnv) seedSubject(name string, balance int64, currency string) string {
	e.t.Helper()
	sub, _, err := e.subjects.GetOrCreateSubject(context.Background(), "test-issuer", e.qualify(name), now)
	if err != nil {
		e.t.Fatalf("seed subject: %v", err)
	}
	acct, _, err := e.wallet.GetOrCreateSubjectAccount(sub.ID, now)
	if err != nil {
		e.t.Fatalf("seed account: %v", err)
	}
	if currency != "" && acct.Currency != currency {
		e.t.Fatalf("account currency = %s, want %s", acct.Currency, currency)
	}
	if balance > 0 {
		entryID := fmt.Sprintf("%s-fund-%s", sub.ID, name)
		if _, _, err := e.wallet.Mutate(acct.ID, walletstore.LedgerEntryInput{
			EntryType: walletstore.EntryAdjust, AmountDelta: balance,
			Memo: "seed funding", ActorID: "test", ActorName: "Test",
		}, entryID, now); err != nil {
			e.t.Fatalf("seed funding: %v", err)
		}
	}
	return sub.ID
}

func (e *testEnv) seedOffer(form string, price int64, status string) *store.Offer {
	e.t.Helper()
	in := service.CreateOfferInput{
		Name: e.qualify("offer-" + form + "-" + status), PriceAmount: price, Currency: "CNY",
		EntitlementForm: form,
	}
	if form == store.FormDuration {
		in.DurationSeconds = 30 * 24 * 3600
	} else {
		in.CountPerPurchase = 5
	}
	if status == store.StatusOnSale {
		in.OnSale = true
	}
	offer, err := e.svc.CreateOffer(context.Background(), service.Actor{ID: "admin-1", Name: "Admin"}, in, now)
	if err != nil {
		e.t.Fatalf("seed offer: %v", err)
	}
	return offer
}

func (e *testEnv) accountOf(subjectID string) *walletstore.Account {
	e.t.Helper()
	acct, err := e.wallet.GetSubjectAccountByOwner(subjectID)
	if err != nil {
		e.t.Fatalf("load account: %v", err)
	}
	return acct
}

func (e *testEnv) ledgerKinds(accountID string) map[string]int {
	e.t.Helper()
	entries, _, err := e.wallet.ListEntries(accountID, "", "", 1, 200)
	if err != nil {
		e.t.Fatalf("list entries: %v", err)
	}
	kinds := map[string]int{}
	for _, entry := range entries {
		kinds[entry.EntryType]++
		if entry.EntryType == walletstore.EntryFreeze || entry.EntryType == walletstore.EntryDeductFrozen {
			if entry.RefType != "biz_offer_purchase" {
				e.t.Fatalf("ledger ref = %s/%s, want biz_offer_purchase back-link", entry.RefType, entry.RefID)
			}
		}
	}
	return kinds
}

// runPurchaseMatrix executes the D-002 §4.2/§4.4 acceptance scenarios that
// must hold on every dialect. The env factory decides the store: the SQLite
// runner opens a fresh store per scenario, the PostgreSQL runner reuses one
// real PG database (A-004 F-002 — the matrix itself must run on both).
func runPurchaseMatrix(t *testing.T, newEnv func(t *testing.T) *testEnv) {
	t.Run("single transaction", func(t *testing.T) {
		env := newEnv(t)
		subjectID := env.seedSubject("buyer", 1000, "")
		offer := env.seedOffer(store.FormCount, 400, store.StatusOnSale)

		res, err := env.svc.Purchase(context.Background(), subjectID, offer.ID, env.requestID("req-1"), time.Now().UTC())
		if err != nil {
			t.Fatalf("purchase: %v", err)
		}
		if res.Replayed {
			t.Fatal("first purchase must not be a replay")
		}
		if res.Purchase.Status != "fulfilled" || res.Purchase.Amount != 400 {
			t.Fatalf("purchase = %+v", res.Purchase)
		}
		if res.Entitlement.PurchaseID != res.Purchase.ID || res.Entitlement.SubjectID != subjectID {
			t.Fatalf("entitlement = %+v", res.Entitlement)
		}
		if res.Entitlement.RemainingCount == nil || *res.Entitlement.RemainingCount != 5 {
			t.Fatalf("count entitlement remaining = %v", res.Entitlement.RemainingCount)
		}
		acct := env.accountOf(subjectID)
		if acct.BalanceTotal != 600 || acct.BalanceAvailable != 600 || acct.BalanceFrozen != 0 {
			t.Fatalf("balances = %+v", acct)
		}
		kinds := env.ledgerKinds(acct.ID)
		if kinds[walletstore.EntryFreeze] != 1 || kinds[walletstore.EntryDeductFrozen] != 1 {
			t.Fatalf("ledger kinds = %v, want exactly one freeze and one deduct_frozen", kinds)
		}
	})

	t.Run("single transaction duration form", func(t *testing.T) {
		env := newEnv(t)
		subjectID := env.seedSubject("buyer", 1000, "")
		offer := env.seedOffer(store.FormDuration, 400, store.StatusOnSale)
		res, err := env.svc.Purchase(context.Background(), subjectID, offer.ID, env.requestID("req-1"), time.Now().UTC())
		if err != nil {
			t.Fatalf("purchase: %v", err)
		}
		if res.Entitlement.ExpiresAt == nil || !res.Entitlement.ExpiresAt.After(now) {
			t.Fatalf("duration entitlement expires = %v", res.Entitlement.ExpiresAt)
		}
	})

	t.Run("terminal failure paths leave zero residue", func(t *testing.T) {
		cases := []struct {
			name    string
			balance int64
			setup   func(e *testEnv, seededID string) (string, string)
			want    error
		}{
			{"insufficient funds", 100, func(e *testEnv, _ string) (string, string) {
				offer := e.seedOffer(store.FormCount, 400, store.StatusOnSale)
				return e.seedSubject("payer", 100, ""), offer.ID
			}, walletstore.ErrInsufficient},
			{"offer not on sale", 1000, func(e *testEnv, subjectID string) (string, string) {
				offer := e.seedOffer(store.FormCount, 400, store.StatusDraft)
				return subjectID, offer.ID
			}, store.ErrOfferNotOnSale},
			{"unknown subject", 1000, func(e *testEnv, _ string) (string, string) {
				offer := e.seedOffer(store.FormCount, 400, store.StatusOnSale)
				return "subject-does-not-exist", offer.ID
			}, store.ErrSubjectNotFound},
			{"unknown offer", 1000, func(e *testEnv, subjectID string) (string, string) {
				return subjectID, "offer-does-not-exist"
			}, store.ErrNotFound},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				env := newEnv(t)
				seededID := env.seedSubject("buyer", tc.balance, "")
				wantSubjectID, offerID := tc.setup(env, seededID)
				_, err := env.svc.Purchase(context.Background(), wantSubjectID, offerID, "req-x", time.Now().UTC())
				if !errors.Is(err, tc.want) {
					t.Fatalf("err = %v, want %v", err, tc.want)
				}
				purchases, total, err := env.svc.ListPurchases(store.PurchaseFilter{SubjectID: wantSubjectID, Page: 1, PageSize: 50})
				if err != nil {
					t.Fatal(err)
				}
				if total != 0 || len(purchases) != 0 {
					t.Fatalf("failed purchase must not persist a voucher: %d rows", total)
				}
				if acct, err := env.wallet.GetSubjectAccountByOwner(wantSubjectID); err == nil {
					if acct.BalanceFrozen != 0 || acct.BalanceTotal != tc.balance {
						t.Fatalf("balances after failure = %+v, want untouched %d", acct, tc.balance)
					}
				}
			})
		}
	})

	t.Run("idempotent replay and cross-offer conflict", func(t *testing.T) {
		env := newEnv(t)
		subjectID := env.seedSubject("buyer", 1000, "")
		offer := env.seedOffer(store.FormCount, 400, store.StatusOnSale)

		first, err := env.svc.Purchase(context.Background(), subjectID, offer.ID, env.requestID("req-1"), time.Now().UTC())
		if err != nil {
			t.Fatalf("first purchase: %v", err)
		}
		second, err := env.svc.Purchase(context.Background(), subjectID, offer.ID, env.requestID("req-1"), time.Now().UTC())
		if err != nil {
			t.Fatalf("replay: %v", err)
		}
		if !second.Replayed || second.Purchase.ID != first.Purchase.ID {
			t.Fatalf("replay mismatch: %+v vs %+v", second, first)
		}
		if acct := env.accountOf(subjectID); acct.BalanceTotal != 600 {
			t.Fatalf("balance after replay = %d, want 600", acct.BalanceTotal)
		}

		other := env.seedOffer(store.FormCount, 100, store.StatusOnSale)
		if _, err := env.svc.Purchase(context.Background(), subjectID, other.ID, env.requestID("req-1"), time.Now().UTC()); !errors.Is(err, store.ErrRequestIdConflict) {
			t.Fatalf("cross-offer replay err = %v, want ErrRequestIdConflict", err)
		}
	})

	t.Run("concurrent same request converges", func(t *testing.T) {
		env := newEnv(t)
		subjectID := env.seedSubject("buyer", 1000, "")
		offer := env.seedOffer(store.FormCount, 400, store.StatusOnSale)

		const n = 6
		results := make([]*service.PurchaseResult, n)
		errs := make([]error, n)
		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				results[i], errs[i] = env.svc.Purchase(context.Background(), subjectID, offer.ID, env.requestID("req-dup"), time.Now().UTC())
			}(i)
		}
		wg.Wait()
		ids := map[string]int{}
		for i, err := range errs {
			if err != nil {
				t.Fatalf("attempt %d: %v", i, err)
			}
			ids[results[i].Purchase.ID]++
		}
		if len(ids) != 1 {
			t.Fatalf("distinct vouchers = %d, want exactly 1", len(ids))
		}
		if acct := env.accountOf(subjectID); acct.BalanceTotal != 600 || acct.BalanceFrozen != 0 {
			t.Fatalf("balances = %+v, want one deduction only", acct)
		}
		_, total, _ := env.svc.ListEntitlements(store.EntitlementFilter{SubjectID: subjectID, Page: 1, PageSize: 50})
		if total != 1 {
			t.Fatalf("entitlements = %d, want 1", total)
		}
	})

	t.Run("concurrent balance race leaves no frozen residue", func(t *testing.T) {
		env := newEnv(t)
		subjectID := env.seedSubject("buyer", 400, "")
		offer := env.seedOffer(store.FormCount, 400, store.StatusOnSale)

		const n = 6
		success := 0
		errs := make([]error, n)
		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				_, err := env.svc.Purchase(context.Background(), subjectID, offer.ID, fmt.Sprintf("req-%d", i), time.Now().UTC())
				errs[i] = err
			}(i)
		}
		wg.Wait()
		for i, err := range errs {
			if err == nil {
				success++
				continue
			}
			if !errors.Is(err, walletstore.ErrInsufficient) {
				t.Fatalf("attempt %d: non-terminal error %v", i, err)
			}
		}
		if success != 1 {
			t.Fatalf("successful purchases = %d, want exactly 1", success)
		}
		if acct := env.accountOf(subjectID); acct.BalanceTotal != 0 || acct.BalanceAvailable != 0 || acct.BalanceFrozen != 0 {
			t.Fatalf("balances = %+v, want fully consumed with no frozen residue", acct)
		}
	})

	t.Run("currency mismatch rejected", func(t *testing.T) {
		env := newEnv(t)
		subjectID := env.seedSubject("buyer", 1000, "CNY")
		offer, err := env.svc.CreateOffer(context.Background(), service.Actor{ID: "admin-1", Name: "Admin"}, service.CreateOfferInput{
			Name: "usd-offer", PriceAmount: 100, Currency: "USD", EntitlementForm: store.FormCount, CountPerPurchase: 1, OnSale: true,
		}, now)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := env.svc.Purchase(context.Background(), subjectID, offer.ID, env.requestID("req-usd"), time.Now().UTC()); !errors.Is(err, store.ErrCurrencyMismatch) {
			t.Fatalf("err = %v, want ErrCurrencyMismatch", err)
		}
	})

	t.Run("retry exhaustion leaves zero residue", func(t *testing.T) {
		env := newEnv(t)
		subjectID := env.seedSubject("buyer", 1000, "")
		offer := env.seedOffer(store.FormCount, 400, store.StatusOnSale)

		attempts := 0
		env.svc.SetPurchaseFaultHookForTest(func(step string) error {
			if step == "deduct" {
				attempts++
				return errors.New("injected transient failure")
			}
			return nil
		})
		_, err := env.svc.Purchase(context.Background(), subjectID, offer.ID, env.requestID("req-exhaust"), time.Now().UTC())
		if err == nil || !strings.Contains(err.Error(), "attempts exhausted") {
			t.Fatalf("err = %v, want bounded exhaustion", err)
		}
		if attempts != 3 {
			t.Fatalf("deduct attempts = %d, want exactly 3", attempts)
		}
		acct := env.accountOf(subjectID)
		if acct.BalanceTotal != 1000 || acct.BalanceFrozen != 0 || acct.BalanceAvailable != 1000 {
			t.Fatalf("balances = %+v, want untouched after exhaustion", acct)
		}
		kinds := env.ledgerKinds(acct.ID)
		if kinds[walletstore.EntryFreeze] != 0 || kinds[walletstore.EntryDeductFrozen] != 0 {
			t.Fatalf("ledger kinds = %v, want no purchase entries after exhaustion", kinds)
		}
		_, total, _ := env.svc.ListPurchases(store.PurchaseFilter{SubjectID: subjectID, Page: 1, PageSize: 50})
		if total != 0 {
			t.Fatalf("exhausted purchase must not persist a voucher")
		}
		_, total, _ = env.svc.ListEntitlements(store.EntitlementFilter{SubjectID: subjectID, Page: 1, PageSize: 50})
		if total != 0 {
			t.Fatalf("exhausted purchase must not grant entitlements")
		}
	})

	t.Run("terminal error is not retried", func(t *testing.T) {
		env := newEnv(t)
		subjectID := env.seedSubject("buyer", 1000, "")
		offer := env.seedOffer(store.FormCount, 400, store.StatusOnSale)

		calls := 0
		env.svc.SetPurchaseFaultHookForTest(func(step string) error {
			calls++
			return walletstore.ErrInsufficient
		})
		if _, err := env.svc.Purchase(context.Background(), subjectID, offer.ID, env.requestID("req-term"), time.Now().UTC()); !errors.Is(err, walletstore.ErrInsufficient) {
			t.Fatalf("err = %v, want terminal ErrInsufficient", err)
		}
		if calls != 1 {
			t.Fatalf("hook calls = %d, want exactly 1 (terminal errors are not retried)", calls)
		}
		if acct := env.accountOf(subjectID); acct.BalanceTotal != 1000 || acct.BalanceFrozen != 0 {
			t.Fatalf("balances = %+v, want untouched", acct)
		}
	})

	t.Run("transient failure then success has no duplicates", func(t *testing.T) {
		env := newEnv(t)
		subjectID := env.seedSubject("buyer", 1000, "")
		offer := env.seedOffer(store.FormCount, 400, store.StatusOnSale)

		calls := 0
		env.svc.SetPurchaseFaultHookForTest(func(step string) error {
			if step == "purchase" {
				calls++
				if calls == 1 {
					return errors.New("injected transient failure")
				}
			}
			return nil
		})
		res, err := env.svc.Purchase(context.Background(), subjectID, offer.ID, env.requestID("req-recover"), time.Now().UTC())
		if err != nil {
			t.Fatalf("purchase: %v", err)
		}
		if res.Replayed {
			t.Fatal("retry success must be a fresh commit, not a replay")
		}
		if acct := env.accountOf(subjectID); acct.BalanceTotal != 600 || acct.BalanceFrozen != 0 {
			t.Fatalf("balances = %+v, want exactly one deduction", acct)
		}
		kinds := env.ledgerKinds(env.accountOf(subjectID).ID)
		if kinds[walletstore.EntryFreeze] != 1 || kinds[walletstore.EntryDeductFrozen] != 1 {
			t.Fatalf("ledger kinds = %v, want exactly one freeze and one deduct_frozen", kinds)
		}
		_, total, _ := env.svc.ListEntitlements(store.EntitlementFilter{SubjectID: subjectID, Page: 1, PageSize: 50})
		if total != 1 {
			t.Fatalf("entitlements = %d, want 1", total)
		}
	})

	t.Run("admin writes fail closed with audit", func(t *testing.T) {
		env := newEnv(t)
		actor := service.Actor{ID: "admin-1", Name: "Admin"}
		offer := env.seedOffer(store.FormCount, 400, store.StatusOnSale)

		env.operations.reset()
		price := int64(500)
		if _, err := env.svc.UpdateOffer(context.Background(), actor, offer.ID, service.UpdateOfferInput{PriceAmount: &price}, offer.Version, now); err != nil {
			t.Fatalf("update offer: %v", err)
		}
		if events := env.operations.events(); len(events) != 1 || events[0] != service.EventOfferUpdate {
			t.Fatalf("audit events = %v, want exactly one %s", events, service.EventOfferUpdate)
		}

		subjectID := env.seedSubject("buyer", 1000, "")
		res, err := env.svc.Purchase(context.Background(), subjectID, offer.ID, env.requestID("req-void"), time.Now().UTC())
		if err != nil {
			t.Fatal(err)
		}
		env.operations.reset()
		if _, err := env.svc.VoidEntitlement(context.Background(), actor, res.Entitlement.ID, time.Now().UTC()); err != nil {
			t.Fatalf("void: %v", err)
		}
		if events := env.operations.events(); len(events) != 1 || events[0] != service.EventEntitlementVoid {
			t.Fatalf("void audit = %v, want exactly one %s", events, service.EventEntitlementVoid)
		}
		env.operations.reset()
		again, err := env.svc.VoidEntitlement(context.Background(), actor, res.Entitlement.ID, time.Now().UTC())
		if err != nil || !again.AlreadyVoided {
			t.Fatalf("repeat void = %+v err %v, want idempotent success", again, err)
		}
		if events := env.operations.events(); len(events) != 0 {
			t.Fatalf("repeat void must not append audit rows, got %v", events)
		}

		env.operations.setFail(errors.New("forced audit failure"))
		in := service.CreateOfferInput{Name: "rolled-back", PriceAmount: 1, Currency: "CNY", EntitlementForm: store.FormCount, CountPerPurchase: 1}
		if _, err := env.svc.CreateOffer(context.Background(), actor, in, now); err == nil {
			t.Fatal("create must fail when the audit write fails")
		}
		env.operations.setFail(nil)
		if _, total, err := env.svc.ListOffers(store.OfferFilter{Q: "rolled-back", Page: 1, PageSize: 10}); (err == nil && total != 0) || err != nil {
			t.Fatalf("domain write must roll back with the audit failure (total=%d err=%v)", total, err)
		}
	})

	t.Run("search handles multibyte and oversized Q", func(t *testing.T) {
		env := newEnv(t)
		actor := service.Actor{ID: "admin-1", Name: "Admin"}
		if _, err := env.svc.CreateOffer(context.Background(), actor, service.CreateOfferInput{
			Name: "会员套餐 🎫", PriceAmount: 100, Currency: "CNY", EntitlementForm: store.FormCount, CountPerPurchase: 1, OnSale: true,
		}, now); err != nil {
			t.Fatal(err)
		}
		// 60 CJK runes = 180 bytes: must not produce invalid UTF-8 or an error.
		if _, _, err := env.svc.ListOffers(store.OfferFilter{Q: strings.Repeat("汉", 60), Page: 1, PageSize: 10}); err != nil {
			t.Fatalf("multibyte oversized Q: %v", err)
		}
		// Emoji are multi-rune: 60 emoji Q must survive rune truncation and
		// still match the offer created above.
		if _, _, err := env.svc.ListOffers(store.OfferFilter{Q: strings.Repeat("🎫", 60), Page: 1, PageSize: 10}); err != nil {
			t.Fatalf("emoji oversized Q: %v", err)
		}
		offers, total, err := env.svc.ListOffers(store.OfferFilter{Q: "会员套餐", Page: 1, PageSize: 10})
		if err != nil || total != 1 || len(offers) != 1 {
			t.Fatalf("multibyte match = %d rows err %v, want 1", total, err)
		}
		if _, total, err := env.svc.ListOffers(store.OfferFilter{Q: "🎫", Page: 1, PageSize: 10}); err != nil || total != 1 {
			t.Fatalf("emoji Q match = %d rows err %v, want 1", total, err)
		}
		// The rune boundary: a 100-rune Q that prefixes the oversized term
		// stays valid; truncation must never split a rune.
		long := strings.Repeat("汉", 99) + "a"
		if got := len([]rune(long)); got != 100 {
			t.Fatalf("boundary setup = %d runes, want 100", got)
		}
		if _, _, err := env.svc.ListOffers(store.OfferFilter{Q: long + "a", Page: 1, PageSize: 10}); err != nil {
			t.Fatalf("101-rune Q must truncate cleanly to 100 runes: %v", err)
		}
		if _, _, err := env.svc.ListOffers(store.OfferFilter{Q: strings.Repeat("a", 101), Page: 1, PageSize: 10}); err != nil {
			t.Fatalf("oversized ascii Q: %v", err)
		}
	})
}

// TestPurchaseMatrixSQLite runs the full matrix on the default SQLite dialect
// (fresh store per scenario).
func TestPurchaseMatrixSQLite(t *testing.T) {
	runPurchaseMatrix(t, func(t *testing.T) *testEnv {
		return newTestEnvOn(t, mustSQLite(t), false)
	})
}

// TestPurchaseRateLimited proves §8: the Service API purchase bucket
// (bizoffer|purchase|<subject_id>, 1 min / 10, AllowRecord, never Clear).
func TestPurchaseRateLimited(t *testing.T) {
	env := newTestEnv(t) // limiter wired (production shape)
	subjectID := env.seedSubject("buyer", 100, "")
	offer := env.seedOffer(store.FormCount, 1, store.StatusOnSale)

	base := time.Unix(1700000000, 0).UTC()
	for i := 0; i < 10; i++ {
		if _, err := env.svc.Purchase(context.Background(), subjectID, offer.ID, fmt.Sprintf("req-%d", i), base.Add(time.Duration(i)*time.Second)); err != nil {
			t.Fatalf("purchase %d within budget: %v", i, err)
		}
	}
	// Budget exhausted: denied without touching the domain.
	if _, err := env.svc.Purchase(context.Background(), subjectID, offer.ID, "req-denied", base.Add(11*time.Second)); !errors.Is(err, service.ErrRateLimited) {
		t.Fatalf("err = %v, want ErrRateLimited", err)
	}
	purchases, total, _ := env.svc.ListPurchases(store.PurchaseFilter{Page: 1, PageSize: 50})
	if total != 10 || len(purchases) != 10 {
		t.Fatalf("purchases = %d, want 10", total)
	}
	// Subject isolation: a different subject has its own bucket.
	other := env.seedSubject("other", 100, "")
	if _, err := env.svc.Purchase(context.Background(), other, offer.ID, "req-other", base.Add(12*time.Second)); err != nil {
		t.Fatalf("other subject must have an independent bucket: %v", err)
	}
	// Window recovery: after the sliding window passes, the first subject is
	// allowed again (no key-wide Clear ever invoked).
	if _, err := env.svc.Purchase(context.Background(), subjectID, offer.ID, "req-later", base.Add(61*time.Second)); err != nil {
		t.Fatalf("window recovery: %v", err)
	}
}

func mustSQLite(t *testing.T) kernel.Store {
	t.Helper()
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "matrix.db"), "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return st
}

// TestPurchasePostgresAcceptance runs the §4.2/§4.4 matrix against real
// PostgreSQL (env-gated; D-002 dual-dialect demand, A-002 F-002).
func TestPurchasePostgresAcceptance(t *testing.T) {
	dsn := pgtest.DSN()
	if dsn == "" {
		t.Skip("postgres test env not set (PG_TEST_*); skipping postgres purchase acceptance")
	}
	base, err := url.Parse(dsn)
	if err != nil {
		t.Fatal(err)
	}
	admin, err := sql.Open("pgx", base.String())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = admin.Close() })
	const dbName = "r2purchase"
	if _, err := admin.ExecContext(context.Background(), `DROP DATABASE IF EXISTS `+dbName+` WITH (FORCE)`); err != nil {
		t.Fatal(err)
	}
	if _, err := admin.ExecContext(context.Background(), `CREATE DATABASE `+dbName); err != nil {
		t.Fatal(err)
	}
	u := *base
	u.Path = "/" + dbName
	t.Cleanup(func() {
		_, _ = admin.ExecContext(context.Background(), `DROP DATABASE IF EXISTS `+dbName+` WITH (FORCE)`)
	})
	catalog, err := compiled.PersistenceCatalog()
	if err != nil {
		t.Fatal(err)
	}
	st, err := storepkg.Open(context.Background(), storepkg.OpenOptions{Dialect: kernel.DialectPostgres, DSN: u.String()}, catalog)
	if err != nil {
		t.Fatalf("open postgres store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })

	// Each scenario gets its own namespace on the shared PG database so
	// subjects, offers and request guards never collide across subtests.
	var pgSeq atomic.Uint64
	runPurchaseMatrix(t, func(t *testing.T) *testEnv {
		e := newTestEnvOn(t, st, false)
		e.salt = fmt.Sprintf("pg-%d", pgSeq.Add(1))
		return e
	})
	env := newTestEnvOn(t, st, false)
	env.salt = "pg-final"
	subjectID := env.seedSubject("pg-buyer", 1000, "")
	offer := env.seedOffer(store.FormCount, 400, store.StatusOnSale)

	// PG-specific concurrent double-fire (distinct goroutines, same request).
	const n = 4
	results := make([]*service.PurchaseResult, n)
	errs := make([]error, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results[i], errs[i] = env.svc.Purchase(context.Background(), subjectID, offer.ID, "req-pg-dup", time.Now().UTC())
		}(i)
	}
	wg.Wait()
	ids := map[string]int{}
	for i, err := range errs {
		if err != nil {
			t.Fatalf("pg attempt %d: %v", i, err)
		}
		ids[results[i].Purchase.ID]++
	}
	if len(ids) != 1 {
		t.Fatalf("pg distinct vouchers = %d, want exactly 1", len(ids))
	}
	if acct := env.accountOf(subjectID); acct.BalanceTotal != 600 || acct.BalanceFrozen != 0 {
		t.Fatalf("pg balances = %+v, want one deduction only", acct)
	}
}
