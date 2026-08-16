// GOAL-020 A-003 F-002 regression: an invalid by-owner adjust still leaves
// the account created AND audited with wallet.account-create (audit precedes
// the Mutate call).
package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWalletByOwnerAdjustFailureStillAuditsCreate(t *testing.T) {
	env, _ := newWalletEnv(t)
	adminToken := env.login(t, testSeedUsername, testSeedPassword)

	// Invalid payload: amountDelta 0 (apply table rejects) + no memo.
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodPost, "/api/wallet/by-owner/u300/adjust", `{"amountDelta":0,"memo":""}`))
	if rr.Code != http.StatusBadRequest && rr.Code != http.StatusConflict {
		t.Fatalf("invalid adjust = %d, want 4xx", rr.Code)
	}

	// The account row exists (auto-open happened before the failed adjust).
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodGet, "/api/wallet/by-owner/u300", ""))
	if rr.Code != http.StatusOK {
		t.Fatalf("by-owner read after failed adjust = %d %s", rr.Code, rr.Body.String())
	}

	// wallet.account-create (auto) is on record for the opened account.
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodGet, "/api/operations?q=wallet", ""))
	if rr.Code != http.StatusOK {
		t.Fatal("operations failed")
	}
	var ops struct {
		Items []struct {
			Event string `json:"event"`
		} `json:"items"`
	}
	_ = json.Unmarshal(rr.Body.Bytes(), &ops)
	creates := 0
	for _, item := range ops.Items {
		if item.Event == "wallet.account-create" {
			creates++
		}
	}
	if creates < 1 {
		t.Fatalf("wallet.account-create events = %d, want >= 1", creates)
	}
}
