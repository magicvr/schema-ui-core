// Wallet HTTP surface tests (S-14 · GOAL-019 D-002 §3/§6): permission gates,
// domain error mapping and the account/adjust/freeze/reconcile flow.
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/errorcatalog"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	walletstore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/wallet/store"
)

// walletServiceStub adapts the wallet store to WalletService (the module
// package imports handler, so handler tests use a store-backed double).
type walletServiceStub struct {
	repo *walletstore.Repository
}

func (s *walletServiceStub) ListAccounts(q, ownerType string, page, pageSize int) ([]walletstore.Account, int, error) {
	return s.repo.ListAccounts(walletstore.ListFilter{Q: q, OwnerType: ownerType, Page: page, PageSize: pageSize})
}
func (s *walletServiceStub) GetAccount(id string) (*walletstore.Account, error) {
	return s.repo.GetAccount(id)
}
func (s *walletServiceStub) CreateAccount(ownerType, ownerID, currency string, now time.Time) (*walletstore.Account, error) {
	if ownerType == "" || ownerID == "" {
		return nil, walletstore.ErrInvalidEntry
	}
	if currency == "" {
		currency = walletstore.DefaultCurrency
	}
	acct := walletstore.Account{ID: "acct-" + ownerID, OwnerType: ownerType, OwnerID: ownerID, Currency: currency, Status: walletstore.StatusActive, CreatedAt: now, UpdatedAt: now}
	if err := s.repo.CreateAccount(acct); err != nil {
		return nil, err
	}
	return &acct, nil
}
func (s *walletServiceStub) UpdateStatus(id, status string, version int64, now time.Time) (*walletstore.Account, error) {
	return s.repo.UpdateStatus(id, status, version, now)
}
func (s *walletServiceStub) ListEntries(accountID, entryType, q string, page, pageSize int) ([]walletstore.LedgerEntry, int, error) {
	if _, err := s.repo.GetAccount(accountID); err != nil {
		return nil, 0, err
	}
	return s.repo.ListEntries(accountID, entryType, q, page, pageSize)
}
func (s *walletServiceStub) Mutate(id string, in walletstore.LedgerEntryInput, now time.Time) (*walletstore.Account, *walletstore.LedgerEntry, bool, error) {
	if in.Memo == "" {
		return nil, nil, false, walletstore.ErrInvalidEntry
	}
	entryID := fmt.Sprintf("%016x", now.UnixMilli()) + newOperationID()
	account, entry, err := s.repo.Mutate(id, in, entryID, now)
	if err != nil {
		return nil, nil, false, err
	}
	return account, entry, entry.ID != entryID, nil
}
func (s *walletServiceStub) Reconcile(accountID, actorID string, now time.Time) (*walletstore.ReconciliationRun, error) {
	return s.repo.ReconcileRun(accountID, fmt.Sprintf("%016x", now.UnixMilli())+newOperationID(), actorID, now)
}
func (s *walletServiceStub) ListReconcileRuns(page, pageSize int) ([]walletstore.ReconciliationRun, int, error) {
	return s.repo.ListReconcileRuns(page, pageSize)
}

func mountWalletRoutes(t *testing.T, env *authTestEnv, service WalletService) {
	t.Helper()
	for _, r := range WalletRoutes(env.a, service, env.operations, "admin.wallet") {
		env.mux.Handle(r.Method+" "+r.Pattern, r.Handler)
	}
}

func newWalletEnv(t *testing.T) (*authTestEnv, *walletServiceStub) {
	t.Helper()
	env := newAuthTestEnv(t)
	stub := &walletServiceStub{repo: walletstore.NewRepository(env.st)}
	mountWalletRoutes(t, env, stub)
	return env, stub
}

// S-14 (GOAL-019 D-002 §3): gates — anonymous 401; editor without the keys
// 403 on every wallet route.
func TestWalletRoutesGates(t *testing.T) {
	env, _ := newWalletEnv(t)

	for _, tc := range []struct{ method, path string }{
		{"GET", "/api/wallet/accounts"},
		{"POST", "/api/wallet/accounts"},
		{"PATCH", "/api/wallet/accounts/acct-1"},
		{"GET", "/api/wallet/accounts/acct-1/entries"},
		{"GET", "/api/wallet/entries?accountId=acct-1"},
		{"GET", "/api/wallet/by-owner/u1"},
		{"POST", "/api/wallet/by-owner/u1/adjust"},
		{"POST", "/api/wallet/accounts/acct-1/adjust"},
		{"POST", "/api/wallet/accounts/acct-1/freeze"},
		{"POST", "/api/wallet/accounts/acct-1/unfreeze"},
		{"POST", "/api/wallet/accounts/acct-1/deduct-frozen"},
		{"POST", "/api/wallet/reconcile"},
		{"GET", "/api/wallet/reconcile/runs"},
	} {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("%s %s = %d, want 401", tc.method, tc.path, rr.Code)
		}
	}

	env.addUser(t, "editor1", "editor-password", []string{"editor"})
	editorToken := env.login(t, "editor1", "editor-password")
	for _, tc := range []struct{ method, path string }{
		{"GET", "/api/wallet/accounts"},
		{"POST", "/api/wallet/accounts"},
		{"PATCH", "/api/wallet/accounts/acct-1"},
		{"GET", "/api/wallet/accounts/acct-1/entries"},
		{"GET", "/api/wallet/entries?accountId=acct-1"},
		{"GET", "/api/wallet/by-owner/u1"},
		{"POST", "/api/wallet/by-owner/u1/adjust"},
		{"POST", "/api/wallet/accounts/acct-1/adjust"},
		{"POST", "/api/wallet/accounts/acct-1/freeze"},
		{"POST", "/api/wallet/accounts/acct-1/unfreeze"},
		{"POST", "/api/wallet/accounts/acct-1/deduct-frozen"},
		{"POST", "/api/wallet/reconcile"},
		{"GET", "/api/wallet/reconcile/runs"},
	} {
		req := bearer(t, editorToken, tc.method, tc.path, "")
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("%s %s = %d, want 403", tc.method, tc.path, rr.Code)
		}
	}
}

// S-14 flow: create → adjust → freeze → entries → over-freeze 409 →
// reconcile consistent → audit events.
func TestWalletLifecycleAndAdjustFlow(t *testing.T) {
	env, _ := newWalletEnv(t)
	adminToken := env.login(t, testSeedUsername, testSeedPassword)

	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodPost, "/api/wallet/by-owner/u1", ""))
	if rr.Code != http.StatusOK {
		t.Fatalf("by-owner open = %d %s", rr.Code, rr.Body.String())
	}
	var created map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	accountID, _ := created["id"].(string)
	if accountID == "" {
		t.Fatal("no account id")
	}

	// Adjust (wallet.adjust).
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodPost, "/api/wallet/accounts/"+accountID+"/adjust", `{"amountDelta":1000,"memo":"grant"}`))
	if rr.Code != http.StatusOK {
		t.Fatalf("adjust = %d %s", rr.Code, rr.Body.String())
	}
	var adjustResp map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &adjustResp); err != nil {
		t.Fatal(err)
	}
	acct, _ := adjustResp["account"].(map[string]any)
	if acct["balanceTotal"].(float64) != 1000 {
		t.Fatalf("balanceTotal = %v", acct["balanceTotal"])
	}

	// Freeze 300.
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodPost, "/api/wallet/accounts/"+accountID+"/freeze", `{"amount":300,"memo":"hold"}`))
	if rr.Code != http.StatusOK {
		t.Fatalf("freeze = %d %s", rr.Code, rr.Body.String())
	}

	// Entries via the canonical path variant.
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodGet, "/api/wallet/accounts/"+accountID+"/entries", ""))
	if rr.Code != http.StatusOK {
		t.Fatalf("entries = %d %s", rr.Code, rr.Body.String())
	}
	var listResp struct {
		Items []map[string]any `json:"items"`
		Total int              `json:"total"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &listResp); err != nil {
		t.Fatal(err)
	}
	if listResp.Total != 2 {
		t.Fatalf("entries total = %d, want 2", listResp.Total)
	}

	// Over-freeze → 409 INSUFFICIENT_BALANCE.
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodPost, "/api/wallet/accounts/"+accountID+"/freeze", `{"amount":99999,"memo":"too much"}`))
	if rr.Code != http.StatusConflict {
		t.Fatalf("over-freeze = %d, want 409", rr.Code)
	}

	// Unfreeze 100 (audited) then disable the account (audited).
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodPost, "/api/wallet/accounts/"+accountID+"/unfreeze", `{"amount":100,"memo":"release"}`))
	if rr.Code != http.StatusOK {
		t.Fatalf("unfreeze = %d %s", rr.Code, rr.Body.String())
	}
	var afterUnfreeze map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &afterUnfreeze)
	acctAfter, _ := afterUnfreeze["account"].(map[string]any)
	ver, _ := json.Marshal(acctAfter["version"])
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodPatch, "/api/wallet/accounts/"+accountID, `{"status":"disabled","version":`+string(ver)+`}`))
	if rr.Code != http.StatusOK {
		t.Fatalf("status update = %d %s", rr.Code, rr.Body.String())
	}

	// Reconcile → consistent.
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodPost, "/api/wallet/reconcile", ""))
	if rr.Code != http.StatusOK {
		t.Fatalf("reconcile = %d %s", rr.Code, rr.Body.String())
	}
	var run map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &run); err != nil {
		t.Fatal(err)
	}
	if run["result"] != "consistent" {
		t.Fatalf("reconcile result = %v", run["result"])
	}

	// A-007 F-001: operationlog actually received the six wallet events.
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodGet, "/api/operations?q=wallet", ""))
	if rr.Code != http.StatusOK {
		t.Fatalf("operations = %d %s", rr.Code, rr.Body.String())
	}
	var opsResp struct {
		Items []struct {
			Event string `json:"event"`
		} `json:"items"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &opsResp); err != nil {
		t.Fatal(err)
	}
	events := map[string]bool{}
	for _, item := range opsResp.Items {
		events[item.Event] = true
	}
	for _, want := range []string{"wallet.account-create", "wallet.account-update", "wallet.adjust", "wallet.freeze", "wallet.unfreeze", "wallet.reconcile"} {
		if !events[want] {
			t.Fatalf("missing audit event %s (got %v)", want, events)
		}
	}
}

// S-14 idempotency + optimistic lock + disabled rejection.
func TestWalletIdempotencyAndStatus(t *testing.T) {
	env, _ := newWalletEnv(t)
	adminToken := env.login(t, testSeedUsername, testSeedPassword)

	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodPost, "/api/wallet/by-owner/u9", ""))
	if rr.Code != http.StatusOK {
		t.Fatal("by-owner open failed")
	}
	if got := rr.Header().Get("ETag"); got != `"v0"` {
		t.Fatalf("initial ETag = %q, want %q", got, `"v0"`)
	}
	var created map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &created)
	accountID := created["id"].(string)
	version := int64(created["version"].(float64))

	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodGet, "/api/wallet/by-owner/u9", ""))
	if rr.Code != http.StatusOK || rr.Header().Get("ETag") != `"v0"` {
		t.Fatalf("by-owner GET = %d ETag=%q", rr.Code, rr.Header().Get("ETag"))
	}
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodGet, "/api/wallet/accounts", ""))
	if rr.Code != http.StatusOK || rr.Header().Get("ETag") != "" {
		t.Fatalf("account list = %d ETag=%q, want no account ETag", rr.Code, rr.Header().Get("ETag"))
	}

	body := `{"amountDelta":500,"memo":"grant","idempotencyKey":"k9"}`
	var lastAccount map[string]any
	var operationID string
	for i := 0; i < 2; i++ {
		rr = httptest.NewRecorder()
		env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodPost, "/api/wallet/accounts/"+accountID+"/adjust", body))
		if rr.Code != http.StatusOK {
			t.Fatalf("idempotent adjust #%d = %d %s", i, rr.Code, rr.Body.String())
		}
		var ar map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &ar)
		lastAccount, _ = ar["account"].(map[string]any)
		operation, _ := ar["operation"].(map[string]any)
		if operation["state"] != "succeeded" || operation["idempotencyKey"] != "k9" || int64(operation["resourceVersion"].(float64)) != 1 {
			t.Fatalf("operation #%d = %#v", i, operation)
		}
		if got, _ := operation["replayed"].(bool); got != (i == 1) {
			t.Fatalf("operation #%d replayed = %v", i, got)
		}
		if i == 0 {
			operationID, _ = operation["operationId"].(string)
		} else if operation["operationId"] != operationID {
			t.Fatalf("replay operationId = %v, want %s", operation["operationId"], operationID)
		}
		if got := rr.Header().Get("ETag"); got != `"v1"` {
			t.Fatalf("adjust #%d ETag = %q, want %q", i, got, `"v1"`)
		}
	}
	if lastAccount != nil {
		version = int64(lastAccount["version"].(float64))
	}
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodGet, "/api/wallet/accounts/"+accountID+"/entries", ""))
	if rr.Code != http.StatusOK {
		t.Fatal("entries failed")
	}
	var listResp struct {
		Total int `json:"total"`
	}
	_ = json.Unmarshal(rr.Body.Bytes(), &listResp)
	if listResp.Total != 1 {
		t.Fatalf("idempotent replay wrote %d entries, want 1", listResp.Total)
	}
	ops, opTotal, err := env.operations.ListOperationsFiltered(operationlog.OperationFilter{
		Event: operationlog.EventWalletAdjust, Sort: "createdAt", Order: "desc", Page: 1, PageSize: 20,
	})
	if err != nil || opTotal != 1 || len(ops) != 1 {
		t.Fatalf("wallet.adjust audit rows = %d/%d err=%v, want 1", len(ops), opTotal, err)
	}

	// Same key with a different payload is a 409 and does not mint an operation.
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodPost, "/api/wallet/accounts/"+accountID+"/adjust", `{"amountDelta":501,"memo":"grant","idempotencyKey":"k9"}`))
	if rr.Code != http.StatusConflict || !bodyHasCode(rr, "LEDGER_IDEMPOTENCY_CONFLICT") {
		t.Fatalf("different payload replay = %d %s", rr.Code, rr.Body.String())
	}

	// Missing, invalid, or contradictory preconditions fail closed.
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodPatch, "/api/wallet/accounts/"+accountID, `{"status":"disabled"}`))
	if rr.Code != http.StatusPreconditionRequired || !bodyHasCode(rr, "PRECONDITION_REQUIRED") {
		t.Fatalf("missing precondition = %d %s", rr.Code, rr.Body.String())
	}
	req := bearer(t, adminToken, http.MethodPatch, "/api/wallet/accounts/"+accountID, `{"status":"disabled","expectedVersion":1}`)
	req.Header.Set("If-Match", `"v0"`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest || !bodyHasCode(rr, "INVALID_PRECONDITION") {
		t.Fatalf("contradictory precondition = %d %s", rr.Code, rr.Body.String())
	}

	// Stale legacy version → 409; fresh legacy version remains compatible.
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodPatch, "/api/wallet/accounts/"+accountID, `{"status":"disabled","version":0}`))
	if rr.Code != http.StatusConflict || !bodyHasCode(rr, "LEDGER_VERSION_CONFLICT") {
		t.Fatalf("stale status = %d %s, want 409", rr.Code, rr.Body.String())
	}
	ver, _ := json.Marshal(version)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodPatch, "/api/wallet/accounts/"+accountID, `{"status":"disabled","version":`+string(ver)+`}`))
	if rr.Code != http.StatusOK {
		t.Fatalf("fresh status = %d, want 200", rr.Code)
	}
	if got := rr.Header().Get("ETag"); got != `"v2"` {
		t.Fatalf("legacy update ETag = %q, want %q", got, `"v2"`)
	}

	// Disabled account rejects mutations.
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodPost, "/api/wallet/accounts/"+accountID+"/adjust", `{"amountDelta":1,"memo":"m"}`))
	if rr.Code != http.StatusConflict {
		t.Fatalf("disabled adjust = %d, want 409", rr.Code)
	}

	// expectedVersion and header-only If-Match both drive successful updates.
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodPatch, "/api/wallet/accounts/"+accountID, `{"status":"active","expectedVersion":2}`))
	if rr.Code != http.StatusOK || rr.Header().Get("ETag") != `"v3"` {
		t.Fatalf("expectedVersion update = %d ETag=%q", rr.Code, rr.Header().Get("ETag"))
	}
	req = bearer(t, adminToken, http.MethodPatch, "/api/wallet/accounts/"+accountID, `{"status":"disabled"}`)
	req.Header.Set("If-Match", `"v3"`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK || rr.Header().Get("ETag") != `"v4"` {
		t.Fatalf("If-Match update = %d ETag=%q body=%s", rr.Code, rr.Header().Get("ETag"), rr.Body.String())
	}
}

func TestWalletMutationWithoutIdempotencyKeyReturnsOperation(t *testing.T) {
	env, _ := newWalletEnv(t)
	token := env.login(t, testSeedUsername, testSeedPassword)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, token, http.MethodPost, "/api/wallet/by-owner/no-key", ""))
	if rr.Code != http.StatusOK {
		t.Fatalf("open = %d %s", rr.Code, rr.Body.String())
	}
	var account map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &account)
	accountID, _ := account["id"].(string)

	req := bearer(t, token, http.MethodPost, "/api/wallet/accounts/"+accountID+"/adjust", `{"amountDelta":10,"memo":"no-key"}`)
	req.Header.Set("If-Match", `W/"v999"`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("no-key adjust = %d %s", rr.Code, rr.Body.String())
	}
	var response map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &response)
	operation, _ := response["operation"].(map[string]any)
	operationID, _ := operation["operationId"].(string)
	if operationID == "" || operation["state"] != "succeeded" || operation["replayed"] != false || operation["resourceVersion"] != float64(1) {
		t.Fatalf("operation = %#v", operation)
	}
	if _, exists := operation["idempotencyKey"]; exists {
		t.Fatalf("no-key operation unexpectedly contains idempotencyKey: %#v", operation)
	}
}

// S-14 error codes reach the frozen wire contract (localized catalog).
func TestWalletErrorCodesCataloged(t *testing.T) {
	for _, code := range []string{"WALLET_NOT_FOUND", "WALLET_OWNER_TAKEN", "WALLET_DISABLED", "INSUFFICIENT_BALANCE", "LEDGER_VERSION_CONFLICT", "LEDGER_IDEMPOTENCY_CONFLICT", "INVALID_LEDGER_ENTRY", "INVALID_WALLET_BODY", "WALLET_USER_AUTO_ONLY", "PRECONDITION_REQUIRED", "INVALID_PRECONDITION"} {
		entry, ok := errorcatalog.Catalog[code]
		if !ok || entry.En == "" || entry.Zh == "" || entry.MessageKey == "" {
			t.Errorf("wallet code %s not cataloged: %+v", code, entry)
		}
	}
}

// Keep kernel import used (gates test builds route keys through the provider
// surface in the module tests).
var _ = kernel.RouteKey
