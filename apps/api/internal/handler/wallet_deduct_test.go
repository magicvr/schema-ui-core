// GOAL-021 D-001 §1: deduct-frozen endpoint — consumes from the frozen
// bucket (available untouched), audited with wallet.deduct-frozen.
package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWalletDeductFrozenEndpoint(t *testing.T) {
	env, _ := newWalletEnv(t)
	adminToken := env.login(t, testSeedUsername, testSeedPassword)

	// Auto-open + grant + freeze.
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodGet, "/api/wallet/by-owner/ufz", ""))
	if rr.Code != http.StatusOK {
		t.Fatalf("open = %d %s", rr.Code, rr.Body.String())
	}
	var acct map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &acct)
	accountID := acct["id"].(string)

	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodPost, "/api/wallet/accounts/"+accountID+"/adjust", "{\"amountDelta\":1000,\"memo\":\"grant\"}"))
	if rr.Code != http.StatusOK {
		t.Fatal("adjust failed")
	}
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodPost, "/api/wallet/accounts/"+accountID+"/freeze", "{\"amount\":400,\"memo\":\"hold\"}"))
	if rr.Code != http.StatusOK {
		t.Fatal("freeze failed")
	}

	// Deduct 250 from frozen.
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodPost, "/api/wallet/accounts/"+accountID+"/deduct-frozen", "{\"amount\":250,\"memo\":\"settle\"}"))
	if rr.Code != http.StatusOK {
		t.Fatalf("deduct-frozen = %d %s", rr.Code, rr.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	account, _ := resp["account"].(map[string]any)
	if account["balanceTotal"].(float64) != 750 || account["balanceAvailable"].(float64) != 600 || account["balanceFrozen"].(float64) != 150 {
		t.Fatalf("after deduct = %v", account)
	}

	// Over-deduct → 409 INSUFFICIENT_BALANCE (frozen bucket short).
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodPost, "/api/wallet/accounts/"+accountID+"/deduct-frozen", "{\"amount\":99999,\"memo\":\"too much\"}"))
	if rr.Code != http.StatusConflict {
		t.Fatalf("over-deduct = %d, want 409", rr.Code)
	}
	var errBody struct {
		Code string `json:"error"`
	}
	_ = json.Unmarshal(rr.Body.Bytes(), &errBody)
	if errBody.Code != "INSUFFICIENT_BALANCE" {
		t.Fatalf("over-deduct code = %s, want INSUFFICIENT_BALANCE", errBody.Code)
	}

	// wallet.deduct-frozen audit event on record.
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
	found := false
	for _, item := range ops.Items {
		if item.Event == "wallet.deduct-frozen" {
			found = true
		}
	}
	if !found {
		t.Fatal("missing wallet.deduct-frozen audit event")
	}
}