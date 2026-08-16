// GOAL-022 · 我的钱包自服务面（D-002 §1/§2）— identity-only 测试：
// 任意已认证用户只读自己的钱包；惰性 get-or-create（auto 审计标记）；
// 幂等（不重复开户/审计）；流水只返回本人账户；无权限键门禁。
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	walletstore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/wallet/store"
)

func newWalletSelfEnv(t *testing.T) (*authTestEnv, *walletServiceStub) {
	t.Helper()
	env := newAuthTestEnv(t)
	stub := &walletServiceStub{repo: walletstore.NewRepository(env.st)}
	mountWalletRoutes(t, env, stub)
	for _, r := range WalletSelfRoutes(env.a, stub, env.operations, "admin.wallet") {
		env.mux.Handle(r.Method+" "+r.Pattern, r.Handler)
	}
	return env, stub
}

// D-002 §2: identity-only — anonymous 401；editor（无任何 wallet.* 权限键）
// 也能读自己的钱包（与 /api/account/profile 同款，无权限键门禁）。
func TestWalletSelfIdentityOnly(t *testing.T) {
	env, _ := newWalletSelfEnv(t)
	for _, path := range []string{"/api/wallet/me", "/api/wallet/me/entries"} {
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, path, nil))
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("GET %s = %d, want 401", path, rr.Code)
		}
	}

	env.addUser(t, "editor1", "editor-password", []string{"editor"})
	editorToken := env.login(t, "editor1", "editor-password")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, editorToken, http.MethodGet, "/api/wallet/me", ""))
	if rr.Code != http.StatusOK {
		t.Fatalf("editor /me = %d %s, want 200", rr.Code, rr.Body.String())
	}
}

// D-002 §2: 惰性 get-or-create — 首次 /me 自动开户（resourceList 信封 +
// auto 审计标记）；重复读取返回同一账户且只产生一次 create 事件；
// ownerId 恒等于会话用户。
func TestWalletSelfAutoCreateAndIdempotency(t *testing.T) {
	env, _ := newWalletSelfEnv(t)
	env.addUser(t, "alice", "alice-password", []string{"editor"})
	token := env.login(t, "alice", "alice-password")

	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, token, http.MethodGet, "/api/wallet/me", ""))
	if rr.Code != http.StatusOK {
		t.Fatalf("first /me = %d %s", rr.Code, rr.Body.String())
	}
	var first struct {
		Items []map[string]any `json:"items"`
		Total int              `json:"total"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &first); err != nil {
		t.Fatal(err)
	}
	if first.Total != 1 || len(first.Items) != 1 {
		t.Fatalf("envelope = %+v", first)
	}
	acct := first.Items[0]
	if acct["ownerType"] != "user" || acct["ownerId"] != "user-alice" || acct["balanceTotal"].(float64) != 0 {
		t.Fatalf("auto account = %v", acct)
	}
	accountID := acct["id"].(string)

	// 第二次读取：同一账户，仍只有一次 create 事件。
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, token, http.MethodGet, "/api/wallet/me", ""))
	if rr.Code != http.StatusOK {
		t.Fatal("second /me failed")
	}
	var second struct {
		Items []map[string]any `json:"items"`
	}
	_ = json.Unmarshal(rr.Body.Bytes(), &second)
	if second.Items[0]["id"] != accountID {
		t.Fatalf("second /me id = %v, want %s", second.Items[0]["id"], accountID)
	}

	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, token, http.MethodGet, "/api/operations?q=wallet", ""))
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
				t.Fatalf("auto marker missing: %v", item.Detail)
			}
		}
	}
	if creates != 1 {
		t.Fatalf("wallet.account-create events = %d, want 1", creates)
	}
}

// D-002 §1: 自服务面只暴露会话用户自己的钱包 —— 流水按身份推导的账户
// 过滤，不接收任何客户端 ownerId；两个用户互不可见。
func TestWalletSelfEntriesOwnScope(t *testing.T) {
	env, _ := newWalletSelfEnv(t)
	adminToken := env.login(t, testSeedUsername, testSeedPassword)

	env.addUser(t, "alice", "alice-password", []string{"editor"})
	env.addUser(t, "bob", "bob-password", []string{"viewer"})
	// 管理员分别给 alice 2500、bob 1000（不同账户）。
	for _, tc := range []struct{ owner, memo string; delta int }{
		{"user-alice", "grant alice", 2500},
		{"user-bob", "grant bob", 1000},
	} {
		body := fmt.Sprintf(`{"amountDelta":%d,"memo":"%s"}`, tc.delta, tc.memo)
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, bearer(t, adminToken, http.MethodPost, "/api/wallet/by-owner/"+tc.owner+"/adjust", body))
		if rr.Code != http.StatusOK {
			t.Fatalf("grant %s = %d %s", tc.owner, rr.Code, rr.Body.String())
		}
	}

	aliceToken := env.login(t, "alice", "alice-password")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, aliceToken, http.MethodGet, "/api/wallet/me", ""))
	if rr.Code != http.StatusOK {
		t.Fatalf("alice /me = %d %s", rr.Code, rr.Body.String())
	}
	var me struct {
		Items []map[string]any `json:"items"`
	}
	_ = json.Unmarshal(rr.Body.Bytes(), &me)
	if me.Items[0]["ownerId"] != "user-alice" || me.Items[0]["balanceTotal"].(float64) != 2500 {
		t.Fatalf("alice self account = %v", me.Items[0])
	}

	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, aliceToken, http.MethodGet, "/api/wallet/me/entries", ""))
	if rr.Code != http.StatusOK {
		t.Fatalf("alice /me/entries = %d %s", rr.Code, rr.Body.String())
	}
	var entries struct {
		Items []map[string]any `json:"items"`
		Total int              `json:"total"`
	}
	_ = json.Unmarshal(rr.Body.Bytes(), &entries)
	if entries.Total != 1 || len(entries.Items) != 1 {
		t.Fatalf("alice entries = %+v", entries)
	}
	if entries.Items[0]["memo"] != "grant alice" || entries.Items[0]["amountDelta"].(float64) != 2500 {
		t.Fatalf("alice entry = %v", entries.Items[0])
	}

	// bob 只看到自己的账户（1000）与自己的流水，与 alice 完全隔离。
	bobToken := env.login(t, "bob", "bob-password")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, bobToken, http.MethodGet, "/api/wallet/me", ""))
	if rr.Code != http.StatusOK {
		t.Fatalf("bob /me = %d %s", rr.Code, rr.Body.String())
	}
	_ = json.Unmarshal(rr.Body.Bytes(), &me)
	if me.Items[0]["ownerId"] != "user-bob" || me.Items[0]["balanceTotal"].(float64) != 1000 {
		t.Fatalf("bob self account = %v", me.Items[0])
	}
	bobAccountID := me.Items[0]["id"].(string)

	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, bobToken, http.MethodGet, "/api/wallet/me/entries", ""))
	if rr.Code != http.StatusOK {
		t.Fatalf("bob /me/entries = %d %s", rr.Code, rr.Body.String())
	}
	var bobEntries struct {
		Items []map[string]any `json:"items"`
		Total int              `json:"total"`
	}
	_ = json.Unmarshal(rr.Body.Bytes(), &bobEntries)
	if bobEntries.Total != 1 || len(bobEntries.Items) != 1 {
		t.Fatalf("bob entries = %+v", bobEntries)
	}
	if bobEntries.Items[0]["memo"] != "grant bob" || bobEntries.Items[0]["amountDelta"].(float64) != 1000 {
		t.Fatalf("bob entry = %v", bobEntries.Items[0])
	}

	// A-002 F-001：查询参数注入被忽略 —— /me 不接受 ?ownerId=，
	// /me/entries 不接受 ?accountId=（账户恒由会话推导）。
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, aliceToken, http.MethodGet, "/api/wallet/me?ownerId=user-bob", ""))
	if rr.Code != http.StatusOK {
		t.Fatalf("alice /me?ownerId= = %d %s", rr.Code, rr.Body.String())
	}
	_ = json.Unmarshal(rr.Body.Bytes(), &me)
	if me.Items[0]["ownerId"] != "user-alice" {
		t.Fatalf("ownerId query ignored expected alice, got %v", me.Items[0]["ownerId"])
	}

	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, aliceToken, http.MethodGet, "/api/wallet/me/entries?accountId="+bobAccountID, ""))
	if rr.Code != http.StatusOK {
		t.Fatalf("alice /me/entries?accountId= = %d %s", rr.Code, rr.Body.String())
	}
	_ = json.Unmarshal(rr.Body.Bytes(), &entries)
	if entries.Total != 1 || entries.Items[0]["memo"] != "grant alice" {
		t.Fatalf("accountId query ignored expected alice entry, got %+v", entries)
	}
}