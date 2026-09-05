// GOAL-003 C3 handler smoke tests (workspace-031 · GOAL-002 D-002 v1.0.0 §5.3/
// §7): admin surface permission gating, fail-closed audit pairing surfaced via
// HTTP, the public C-end catalog and the immutable-form rejection.
package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/ratelimit"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/modules/authsession/systemdata"
	"github.com/magicvr/schema-ui-core/apps/api/modules/digitaloffer/service"
	"github.com/magicvr/schema-ui-core/apps/api/modules/digitaloffer/store"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
	walletstore "github.com/magicvr/schema-ui-core/apps/api/modules/wallet/store"
	"github.com/magicvr/schema-ui-core/apps/api/modules/wallet/subject"
)

var errForcedAudit = errors.New("forced audit failure")

// stubTxRecorder lets the handler tests force audit failures (§7 fail-closed).
type stubTxRecorder struct {
	mu   sync.Mutex
	fail error
}

func (s *stubTxRecorder) RecordOperation(_ kernel.Tx, _ operationlog.Operation) error { return nil }

func (s *stubTxRecorder) RecordOperationTx(_ kernel.Tx, _ operationlog.Operation) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.fail
}

func newDigitalOfferEnv(t *testing.T) (*authTestEnv, *service.Service, *stubTxRecorder) {
	t.Helper()
	env := newAuthTestEnv(t)
	// Register the module's permission contributions through the production
	// reconciliation path so the seeded admin role actually holds them.
	perms := []kernel.PermissionContribution{
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "biz.digital-offer", Key: "digitaloffer.read"}, Permission: "digitaloffer.read", Resource: "digitaloffer", Action: "read", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "biz.digital-offer", Key: "digitaloffer.offer.manage"}, Permission: "digitaloffer.offer.manage", Resource: "digitaloffer", Action: "offer.manage", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "biz.digital-offer", Key: "digitaloffer.entitlement.void"}, Permission: "digitaloffer.entitlement.void", Resource: "digitaloffer", Action: "entitlement.void", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
	}
	if err := authsessiondata.Reconcile(context.Background(), env.st, perms, nil); err != nil {
		t.Fatalf("reconcile digitaloffer permissions: %v", err)
	}
	rec := &stubTxRecorder{}
	svc := service.NewService(
		store.NewRepository(env.st),
		walletstore.NewRepository(env.st),
		subject.NewStore(env.st),
		rec,
	)
	for _, r := range DigitalOfferRoutes(env.a, svc, "biz.digital-offer", ratelimit.NewProvider()) {
		env.mux.Handle(r.Method+" "+r.Pattern, r.Handler)
	}
	for _, r := range DigitalOfferPublicRoutes(env.a, svc, "biz.digital-offer", ratelimit.NewProvider()) {
		env.mux.Handle(r.Method+" "+r.Pattern, r.Handler)
	}
	return env, svc, rec
}

func serveDigitalOffer(t *testing.T, env *authTestEnv, token, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	var req *http.Request
	if token == "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
	} else {
		req = bearer(t, token, method, path, body)
	}
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	return rr
}

func TestDigitalOfferAdminSurface(t *testing.T) {
	env, _, _ := newDigitalOfferEnv(t)

	// Unauthenticated: every admin route is 401.
	for _, tc := range []struct{ method, path string }{
		{http.MethodGet, "/api/digitaloffer/offers"},
		{http.MethodPost, "/api/digitaloffer/offers"},
		{http.MethodGet, "/api/digitaloffer/purchases"},
		{http.MethodGet, "/api/digitaloffer/entitlements"},
	} {
		rr := serveDigitalOffer(t, env, "", tc.method, tc.path, "")
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("%s %s = %d, want 401", tc.method, tc.path, rr.Code)
		}
	}

	token := env.login(t, testSeedUsername, testSeedPassword)

	// Create (audited).
	body := `{"name":"Pro membership","description":"30 days","priceAmount":1200,"currency":"CNY","entitlementForm":"duration","durationSeconds":2592000,"onSale":true}`
	rr := serveDigitalOffer(t, env, token, http.MethodPost, "/api/digitaloffer/offers", body)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create offer = %d %s, want 201", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), `"entitlementForm":"duration"`) || !strings.Contains(rr.Body.String(), `"status":"on_sale"`) {
		t.Fatalf("created offer = %s", rr.Body.String())
	}
	offerID := between(rr.Body.String(), `"id":"`, `"`)

	// Immutable form fields are rejected (D-002 §2 / BIZOFFER_FORM_CONFLICT).
	rr = serveDigitalOffer(t, env, token, http.MethodPatch, "/api/digitaloffer/offers/"+offerID,
		`{"entitlementForm":"count","version":0}`)
	if rr.Code != http.StatusConflict || !strings.Contains(rr.Body.String(), "BIZOFFER_FORM_CONFLICT") {
		t.Fatalf("form conflict = %d %s", rr.Code, rr.Body.String())
	}

	// Status change with optimistic lock.
	rr = serveDigitalOffer(t, env, token, http.MethodPatch, "/api/digitaloffer/offers/"+offerID,
		`{"status":"off_sale","version":0}`)
	if rr.Code != http.StatusOK {
		t.Fatalf("status change = %d %s", rr.Code, rr.Body.String())
	}

	// C-end catalog lists on_sale only — the offer was just taken off sale.
	rr = serveDigitalOffer(t, env, "", http.MethodGet, "/api/biz/offers", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("public list = %d %s", rr.Code, rr.Body.String())
	}
	if strings.Contains(rr.Body.String(), "Pro membership") {
		t.Fatal("off_sale offer must not appear in the public catalog")
	}
	rr = serveDigitalOffer(t, env, token, http.MethodPatch, "/api/digitaloffer/offers/"+offerID,
		`{"status":"on_sale","version":1}`)
	if rr.Code != http.StatusOK {
		t.Fatalf("re-on-sale = %d %s", rr.Code, rr.Body.String())
	}
	rr = serveDigitalOffer(t, env, "", http.MethodGet, "/api/biz/offers", "")
	if !strings.Contains(rr.Body.String(), "Pro membership") {
		t.Fatalf("on_sale offer missing from catalog: %s", rr.Body.String())
	}
	if strings.Contains(rr.Body.String(), "version") {
		t.Fatal("public catalog must not leak internal version fields")
	}
}

func TestDigitalOfferAuditFailureRollsBack(t *testing.T) {
	env, _, rec := newDigitalOfferEnv(t)
	token := env.login(t, testSeedUsername, testSeedPassword)

	rec.mu.Lock()
	rec.fail = errForcedAudit
	rec.mu.Unlock()
	rr := serveDigitalOffer(t, env, token, http.MethodPost, "/api/digitaloffer/offers",
		`{"name":"ghost","priceAmount":1,"currency":"CNY","entitlementForm":"count","countPerPurchase":1}`)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("audited create with forced audit failure = %d, want 500", rr.Code)
	}
	rec.mu.Lock()
	rec.fail = nil
	rec.mu.Unlock()

	// The domain row must be gone: a fresh read shows an empty offers page.
	rr = serveDigitalOffer(t, env, token, http.MethodGet, "/api/digitaloffer/offers", "")
	if rr.Code != http.StatusOK || strings.Contains(rr.Body.String(), "ghost") {
		t.Fatalf("audit failure must roll the domain write back, got %d %s", rr.Code, rr.Body.String())
	}
}

// TestDigitalOfferPurchasedEntitlementVisible proves the R2 money path lands
// rows the admin read surface can see (purchase via service API, §4.3).
func TestDigitalOfferPurchasedEntitlementVisible(t *testing.T) {
	env, svc, _ := newDigitalOfferEnv(t)
	token := env.login(t, testSeedUsername, testSeedPassword)

	now := time.Unix(1700000000, 0).UTC()
	offer, err := svc.CreateOffer(context.Background(), service.Actor{ID: "admin", Name: "Admin"}, service.CreateOfferInput{
		Name: "Pack 10", PriceAmount: 500, Currency: "CNY", EntitlementForm: store.FormCount,
		CountPerPurchase: 10, OnSale: true,
	}, now)
	if err != nil {
		t.Fatal(err)
	}
	sub, _, err := subject.NewStore(env.st).GetOrCreateSubject(context.Background(), "test", "buyer-1", now)
	if err != nil {
		t.Fatal(err)
	}
	repo := walletstore.NewRepository(env.st)
	acct, _, err := repo.GetOrCreateSubjectAccount(sub.ID, now)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := repo.Mutate(acct.ID, walletstore.LedgerEntryInput{
		EntryType: walletstore.EntryAdjust, AmountDelta: 1000, Memo: "fund", ActorID: "t", ActorName: "T",
	}, "fund-entry-1", now); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Purchase(context.Background(), sub.ID, offer.ID, "req-smoke", now); err != nil {
		t.Fatalf("purchase: %v", err)
	}

	rr := serveDigitalOffer(t, env, token, http.MethodGet, "/api/digitaloffer/entitlements?subjectId="+sub.ID, "")
	if rr.Code != http.StatusOK || !strings.Contains(rr.Body.String(), `"remainingCount":10`) {
		t.Fatalf("entitlement list = %d %s", rr.Code, rr.Body.String())
	}
	rr = serveDigitalOffer(t, env, token, http.MethodGet, "/api/digitaloffer/purchases", "")
	if rr.Code != http.StatusOK || !strings.Contains(rr.Body.String(), `"status":"fulfilled"`) {
		t.Fatalf("purchase list = %d %s", rr.Code, rr.Body.String())
	}
}

// TestDigitalOfferPublicListRateLimited proves §8: the C-end catalog bucket
// rejects with 429 + Retry-After when exhausted.
func TestDigitalOfferPublicListRateLimited(t *testing.T) {
	env, _, _ := newDigitalOfferEnv(t)
	for i := 0; i < bizOfferListRateLimiterMax; i++ {
		rr := serveDigitalOffer(t, env, "", http.MethodGet, "/api/biz/offers", "")
		if rr.Code != http.StatusOK {
			t.Fatalf("request %d = %d, want 200 within budget", i, rr.Code)
		}
	}
	rr := serveDigitalOffer(t, env, "", http.MethodGet, "/api/biz/offers", "")
	if rr.Code != http.StatusTooManyRequests || rr.Header().Get("Retry-After") == "" {
		t.Fatalf("exhausted bucket = %d, want 429 with Retry-After", rr.Code)
	}
}

func between(s, prefix, suffix string) string {
	i := strings.Index(s, prefix)
	if i < 0 {
		return ""
	}
	rest := s[i+len(prefix):]
	j := strings.Index(rest, suffix)
	if j < 0 {
		return ""
	}
	return rest[:j]
}
