// GOAL-003 C3 acceptance tests (workspace-031 · GOAL-002 D-002 v1.0.0 §4):
// the single-transaction purchase path — freeze → deduct_frozen → purchase
// voucher → entitlement in ONE kernel.Tx, idempotent replay, bounded retry on
// races, and the §4.4 concurrency invariants.
package service_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/pgtest"
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

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "digitaloffer.db"), "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return newTestEnvOn(t, st)
}

func newTestEnvOn(t *testing.T, st kernel.Store) *testEnv {
	t.Helper()
	rec := &stubRecorder{}
	return &testEnv{
		t:          t,
		wallet:     walletstore.NewRepository(st),
		subjects:   subject.NewStore(st),
		svc:        service.NewService(store.NewRepository(st), walletstore.NewRepository(st), subject.NewStore(st), rec),
		operations: rec,
	}
}

// seedSubject registers one external subject and funds its wallet account.
func (e *testEnv) seedSubject(name string, balance int64, currency string) string {
	e.t.Helper()
	sub, _, err := e.subjects.GetOrCreateSubject(context.Background(), "test-issuer", name, now)
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
		Name: "offer-" + form + "-" + status, PriceAmount: price, Currency: "CNY",
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

// TestPurchaseSingleTransaction proves §4.2: a successful purchase commits
// freeze + deduct_frozen ledger entries, the purchase voucher and the
// entitlement in one visible state, with wallet refs back-linked.
func TestPurchaseSingleTransaction(t *testing.T) {
	for _, form := range []string{store.FormDuration, store.FormCount} {
		t.Run(form, func(t *testing.T) {
			env := newTestEnv(t)
			subjectID := env.seedSubject("buyer-"+form, 1000, "")
			offer := env.seedOffer(form, 400, store.StatusOnSale)

			res, err := env.svc.Purchase(context.Background(), subjectID, offer.ID, "req-1", time.Now().UTC())
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
			switch form {
			case store.FormDuration:
				if res.Entitlement.ExpiresAt == nil || !res.Entitlement.ExpiresAt.After(now) {
					t.Fatalf("duration entitlement expires = %v", res.Entitlement.ExpiresAt)
				}
			case store.FormCount:
				if res.Entitlement.RemainingCount == nil || *res.Entitlement.RemainingCount != 5 {
					t.Fatalf("count entitlement remaining = %v", res.Entitlement.RemainingCount)
				}
			}

			// Wallet ledger: exactly freeze + deduct_frozen, total dropped by
			// 400, no frozen residue, ref back-link present.
			acct, err := env.wallet.GetSubjectAccountByOwner(subjectID)
			if err != nil {
				t.Fatalf("load account: %v", err)
			}
			if acct.BalanceTotal != 600 || acct.BalanceAvailable != 600 || acct.BalanceFrozen != 0 {
				t.Fatalf("balances = %+v", acct)
			}
			for _, entryID := range []string{res.Purchase.FreezeEntryID, res.Purchase.DeductEntryID} {
				if entryID == "" {
					t.Fatal("purchase must reference both ledger entries")
				}
			}
			entries, _, err := env.wallet.ListEntries(acct.ID, "", "", 1, 50)
			if err != nil {
				t.Fatalf("list entries: %v", err)
			}
			kinds := map[string]int{}
			for _, e := range entries {
				kinds[e.EntryType]++
				switch e.EntryType {
				case walletstore.EntryFreeze, walletstore.EntryDeductFrozen:
					if e.RefType != "biz_offer_purchase" || e.RefID != res.Purchase.ID {
						t.Fatalf("ledger ref = %s/%s, want biz_offer_purchase/%s", e.RefType, e.RefID, res.Purchase.ID)
					}
				}
			}
			if kinds[walletstore.EntryFreeze] != 1 || kinds[walletstore.EntryDeductFrozen] != 1 {
				t.Fatalf("ledger kinds = %v, want exactly one freeze and one deduct_frozen", kinds)
			}
		})
	}
}

// TestPurchaseFailurePaths proves the §4.2 terminal rejections leave NO
// purchase voucher, NO entitlement and NO frozen residue.
func TestPurchaseFailurePaths(t *testing.T) {
	cases := []struct {
		name    string
		balance int64
		setup   func(e *testEnv, subjectID string) (string, string)
		want    error
	}{
		{"insufficient funds", 100, func(e *testEnv, subjectID string) (string, string) {
			offer := e.seedOffer(store.FormCount, 400, store.StatusOnSale)
			return subjectID, offer.ID
		}, walletstore.ErrInsufficient},
		{"offer not on sale", 1000, func(e *testEnv, subjectID string) (string, string) {
			offer := e.seedOffer(store.FormCount, 400, store.StatusDraft)
			return subjectID, offer.ID
		}, store.ErrOfferNotOnSale},
		{"unknown subject", 1000, func(e *testEnv, subjectID string) (string, string) {
			offer := e.seedOffer(store.FormCount, 400, store.StatusOnSale)
			return "subject-does-not-exist", offer.ID
		}, store.ErrSubjectNotFound},
		{"unknown offer", 1000, func(e *testEnv, subjectID string) (string, string) {
			return subjectID, "offer-does-not-exist"
		}, store.ErrNotFound},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			env := newTestEnv(t)
			seededID := env.seedSubject("buyer", tc.balance, "")
			wantSubjectID, offerID := tc.setup(env, seededID)
			_, err := env.svc.Purchase(context.Background(), wantSubjectID, offerID, "req-x", time.Now().UTC())
			if !errors.Is(err, tc.want) {
				t.Fatalf("err = %v, want %v", err, tc.want)
			}
			purchases, total, err := env.svc.ListPurchases(store.PurchaseFilter{Page: 1, PageSize: 50})
			if err != nil {
				t.Fatal(err)
			}
			if total != 0 || len(purchases) != 0 {
				t.Fatalf("failed purchase must not persist a voucher: %d rows", total)
			}
			acct, err := env.wallet.GetSubjectAccountByOwner(wantSubjectID)
			if err == nil && (acct.BalanceFrozen != 0 || acct.BalanceTotal != tc.balance) {
				t.Fatalf("balances after failure = %+v, want untouched %d", acct, tc.balance)
			}
		})
	}
}

// TestPurchaseIdempotentReplay proves §4.4: the same (subject, request) pair
// replays the SAME voucher with no second deduction; a different offer under
// the same request id is a hard conflict.
func TestPurchaseIdempotentReplay(t *testing.T) {
	env := newTestEnv(t)
	subjectID := env.seedSubject("buyer", 1000, "")
	offer := env.seedOffer(store.FormCount, 400, store.StatusOnSale)

	first, err := env.svc.Purchase(context.Background(), subjectID, offer.ID, "req-1", time.Now().UTC())
	if err != nil {
		t.Fatalf("first purchase: %v", err)
	}
	second, err := env.svc.Purchase(context.Background(), subjectID, offer.ID, "req-1", time.Now().UTC())
	if err != nil {
		t.Fatalf("replay: %v", err)
	}
	if !second.Replayed || second.Purchase.ID != first.Purchase.ID {
		t.Fatalf("replay mismatch: %+v vs %+v", second, first)
	}
	acct, _ := env.wallet.GetSubjectAccountByOwner(subjectID)
	if acct.BalanceTotal != 600 {
		t.Fatalf("balance after replay = %d, want 600", acct.BalanceTotal)
	}

	other := env.seedOffer(store.FormCount, 100, store.StatusOnSale)
	if _, err := env.svc.Purchase(context.Background(), subjectID, other.ID, "req-1", time.Now().UTC()); !errors.Is(err, store.ErrRequestIdConflict) {
		t.Fatalf("cross-offer replay err = %v, want ErrRequestIdConflict", err)
	}
}

// TestPurchaseConcurrentSameRequest proves the §4.4 convergence invariant:
// concurrent double-fire of one request yields exactly one voucher, one
// deduction and one entitlement (SQLite writers serialize; the loser replays).
func TestPurchaseConcurrentSameRequest(t *testing.T) {
	env := newTestEnv(t)
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
			results[i], errs[i] = env.svc.Purchase(context.Background(), subjectID, offer.ID, "req-dup", time.Now().UTC())
		}(i)
	}
	wg.Wait()
	for i, err := range errs {
		if err != nil {
			t.Fatalf("attempt %d: %v", i, err)
		}
	}
	ids := map[string]int{}
	for _, r := range results {
		ids[r.Purchase.ID]++
	}
	if len(ids) != 1 {
		t.Fatalf("distinct vouchers = %d, want exactly 1", len(ids))
	}
	acct, _ := env.wallet.GetSubjectAccountByOwner(subjectID)
	if acct.BalanceTotal != 600 || acct.BalanceFrozen != 0 {
		t.Fatalf("balances = %+v, want one deduction only", acct)
	}
	_, total, _ := env.svc.ListEntitlements(store.EntitlementFilter{SubjectID: subjectID, Page: 1, PageSize: 50})
	if total != 1 {
		t.Fatalf("entitlements = %d, want 1", total)
	}
}

// TestPurchaseConcurrentBalanceRace proves 判据 2: with balance for exactly
// one purchase, concurrent distinct requests yield exactly one success and no
// frozen residue.
func TestPurchaseConcurrentBalanceRace(t *testing.T) {
	env := newTestEnv(t)
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
	acct, _ := env.wallet.GetSubjectAccountByOwner(subjectID)
	if acct.BalanceTotal != 0 || acct.BalanceAvailable != 0 || acct.BalanceFrozen != 0 {
		t.Fatalf("balances = %+v, want fully consumed with no frozen residue", acct)
	}
}

// TestPurchaseCurrencyMismatch proves the §4.4 account-currency gate.
func TestPurchaseCurrencyMismatch(t *testing.T) {
	env := newTestEnv(t)
	subjectID := env.seedSubject("buyer", 1000, "CNY")
	in := service.CreateOfferInput{Name: "usd-offer", PriceAmount: 100, Currency: "USD", EntitlementForm: store.FormCount, CountPerPurchase: 1, OnSale: true}
	offer, err := env.svc.CreateOffer(context.Background(), service.Actor{ID: "admin-1", Name: "Admin"}, in, now)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := env.svc.Purchase(context.Background(), subjectID, offer.ID, "req-usd", time.Now().UTC()); !errors.Is(err, store.ErrCurrencyMismatch) {
		t.Fatalf("err = %v, want ErrCurrencyMismatch", err)
	}
}

// TestAdminWritesFailClosedAudit proves §7: offer create/update/status and
// entitlement void append audit rows in the same transaction; an audit failure
// rolls the domain write back.
func TestAdminWritesFailClosedAudit(t *testing.T) {
	env := newTestEnv(t)
	actor := service.Actor{ID: "admin-1", Name: "Admin"}
	offer := env.seedOffer(store.FormCount, 400, store.StatusOnSale)

	env.operations.mu.Lock()
	env.operations.ops = nil
	env.operations.mu.Unlock()
	price := int64(500)
	if _, err := env.svc.UpdateOffer(context.Background(), actor, offer.ID, service.UpdateOfferInput{PriceAmount: &price}, offer.Version, now); err != nil {
		t.Fatalf("update offer: %v", err)
	}
	env.operations.mu.Lock()
	events := append([]string(nil), env.operations.ops...)
	env.operations.mu.Unlock()
	if len(events) != 1 || events[0] != service.EventOfferUpdate {
		t.Fatalf("audit events = %v, want exactly one %s", events, service.EventOfferUpdate)
	}

	// Void with audit, then repeat: idempotent success, no second audit row.
	subjectID := env.seedSubject("buyer", 1000, "")
	res, err := env.svc.Purchase(context.Background(), subjectID, offer.ID, "req-void", time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	env.operations.mu.Lock()
	env.operations.ops = nil
	env.operations.mu.Unlock()
	if _, err := env.svc.VoidEntitlement(context.Background(), actor, res.Entitlement.ID, time.Now().UTC()); err != nil {
		t.Fatalf("void: %v", err)
	}
	env.operations.mu.Lock()
	events = append([]string(nil), env.operations.ops...)
	env.operations.mu.Unlock()
	if len(events) != 1 || events[0] != service.EventEntitlementVoid {
		t.Fatalf("void audit = %v, want exactly one %s", events, service.EventEntitlementVoid)
	}
	env.operations.mu.Lock()
	env.operations.ops = nil
	env.operations.mu.Unlock()
	again, err := env.svc.VoidEntitlement(context.Background(), actor, res.Entitlement.ID, time.Now().UTC())
	if err != nil || !again.AlreadyVoided {
		t.Fatalf("repeat void = %+v err %v, want idempotent success", again, err)
	}
	env.operations.mu.Lock()
	events = append([]string(nil), env.operations.ops...)
	env.operations.mu.Unlock()
	if len(events) != 0 {
		t.Fatalf("repeat void must not append audit rows, got %v", events)
	}

	// Audit failure injection: the domain write must roll back with it.
	env.operations.mu.Lock()
	env.operations.fail = errors.New("forced audit failure")
	env.operations.mu.Unlock()
	in := service.CreateOfferInput{Name: "rolled-back", PriceAmount: 1, Currency: "CNY", EntitlementForm: store.FormCount, CountPerPurchase: 1}
	if _, err := env.svc.CreateOffer(context.Background(), actor, in, now); err == nil {
		t.Fatal("create must fail when the audit write fails")
	}
	env.operations.mu.Lock()
	env.operations.fail = nil
	env.operations.mu.Unlock()
	offers, total, err := env.svc.ListOffers(store.OfferFilter{Q: "rolled-back", Page: 1, PageSize: 10})
	if err == nil && (total != 0 || len(offers) != 0) {
		t.Fatalf("domain write must roll back with the audit failure, found %d rows", total)
	}
}

// TestPurchasePostgresAcceptance runs the §4.2/§4.4 happy path plus the
// concurrent double-fire invariant against PostgreSQL (env-gated, D-002
// dual-dialect demand).
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

	env := newTestEnvOn(t, st)
	subjectID := env.seedSubject("pg-buyer", 1000, "")
	offer := env.seedOffer(store.FormCount, 400, store.StatusOnSale)

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
	acct, err := env.wallet.GetSubjectAccountByOwner(subjectID)
	if err != nil {
		t.Fatal(err)
	}
	if acct.BalanceTotal != 600 || acct.BalanceFrozen != 0 {
		t.Fatalf("pg balances = %+v, want one deduction only", acct)
	}
}
