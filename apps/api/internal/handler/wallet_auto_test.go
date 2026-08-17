// GOAL-020 · 钱包账户自动开户与用户绑定 — by-owner get-or-create 测试。
// 自动开户（读端点幂等 + auto 审计标记）、by-owner 调账自动开户、
// user 类型手动创建拒绝（WALLET_USER_AUTO_ONLY）。
package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	walletstore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/wallet/store"
)

// GetOrCreateUserAccount adapter lives here (walletServiceStub is declared in
// wallet_test.go; methods may live in any file of the same package).
func (s *walletServiceStub) GetOrCreateUserAccount(ownerID string, now time.Time) (*walletstore.Account, bool, error) {
	return s.repo.GetOrCreateUserAccount(ownerID, now)
}

func (s *walletServiceStub) GetUserAccountByOwner(ownerID string) (*walletstore.Account, error) {
	return s.repo.GetUserAccountByOwner(ownerID)
}

// D-001 §1: by-owner get-or-create — auto-created zero-balance account
// (audited with auto marker); repeated reads return the same row.
func TestWalletByOwnerAutoCreate(t *testing.T) {
	env, _ := newWalletEnv(t)
	adminToken := env.login(t, testSeedUsername, testSeedPassword)

	// GET is read-only (W15-F11): missing owner is 404 until POST creates.
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodGet, "/api/wallet/by-owner/u100", ""))
	if rr.Code != http.StatusNotFound {
		t.Fatalf("by-owner missing GET = %d %s, want 404", rr.Code, rr.Body.String())
	}
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodPost, "/api/wallet/by-owner/u100", ""))
	if rr.Code != http.StatusOK {
		t.Fatalf("by-owner create = %d %s", rr.Code, rr.Body.String())
	}
	var first map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &first); err != nil {
		t.Fatal(err)
	}
	if first["ownerType"] != "user" || first["ownerId"] != "u100" || first["balanceTotal"].(float64) != 0 {
		t.Fatalf("auto account = %v", first)
	}
	accountID := first["id"].(string)

	// Second read returns the same account (no duplicate).
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodGet, "/api/wallet/by-owner/u100", ""))
	if rr.Code != http.StatusOK {
		t.Fatal("second by-owner read failed")
	}
	var second map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &second)
	if second["id"] != accountID {
		t.Fatalf("by-owner read returned different account %v vs %v", second["id"], accountID)
	}

	// Only ONE wallet.account-create event (the auto marker).
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodGet, "/api/operations?q=wallet", ""))
	if rr.Code != http.StatusOK {
		t.Fatal("operations failed")
	}
	var ops struct {
		Items []struct {
			Event  string  `json:"event"`
			Detail *string `json:"detail"`
		} `json:"items"`
	}
	_ = json.Unmarshal(rr.Body.Bytes(), &ops)
	creates := 0
	for _, item := range ops.Items {
		if item.Event == "wallet.account-create" {
			creates++
			if item.Detail == nil || !strings.Contains(*item.Detail, `"auto":true`) {
				t.Fatalf("auto create detail missing marker: %v", item.Detail)
			}
		}
	}
	if creates != 1 {
		t.Fatalf("wallet.account-create events = %d, want 1", creates)
	}
}

// D-001 §1: by-owner adjust auto-creates the account and applies the entry.
func TestWalletByOwnerAdjustAutoCreate(t *testing.T) {
	env, _ := newWalletEnv(t)
	adminToken := env.login(t, testSeedUsername, testSeedPassword)

	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodPost, "/api/wallet/by-owner/u200/adjust", `{"amountDelta":2500,"memo":"initial grant"}`))
	if rr.Code != http.StatusOK {
		t.Fatalf("by-owner adjust = %d %s", rr.Code, rr.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	acct, _ := resp["account"].(map[string]any)
	if acct["ownerId"] != "u200" || acct["balanceTotal"].(float64) != 2500 {
		t.Fatalf("auto account after adjust = %v", acct)
	}
	entry, _ := resp["entry"].(map[string]any)
	if entry["entryType"] != "adjust" || entry["balanceAfterTotal"].(float64) != 2500 {
		t.Fatalf("entry = %v", entry)
	}
}

// D-001 §2: manual creation of user accounts is rejected.
func TestWalletManualUserAccountRejected(t *testing.T) {
	env, _ := newWalletEnv(t)
	adminToken := env.login(t, testSeedUsername, testSeedPassword)

	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodPost, "/api/wallet/accounts", `{"ownerType":"user","ownerId":"u9"}`))
	if rr.Code != http.StatusConflict {
		t.Fatalf("manual user create = %d, want 409", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "WALLET_USER_AUTO_ONLY") {
		t.Fatalf("body = %s, want WALLET_USER_AUTO_ONLY", rr.Body.String())
	}

	// business/system manual creation still works.
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodPost, "/api/wallet/accounts", `{"ownerType":"business","ownerId":"b1"}`))
	if rr.Code != http.StatusCreated {
		t.Fatalf("manual business create = %d %s", rr.Code, rr.Body.String())
	}
}