package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
	walletstore "github.com/magicvr/schema-ui-core/apps/api/modules/wallet/store"
	"github.com/magicvr/schema-ui-core/apps/api/modules/wallet/subject"
	walletschema "github.com/magicvr/schema-ui-core/apps/api/modules/wallet/schema"
)

func TestVouchersBatchGeneratePermissionAndAudit(t *testing.T) {
	env := newAuthTestEnv(t)
	repo := walletstore.NewRepository(env.st)
	service := newWalletStub(env, repo)
	mountWalletRoutes(t, env, service)

	env.addUser(t, "viewer1", "viewer-password", []string{"viewer"})
	adminTok := adminToken(t, env)
	viewerTok := env.login(t, "viewer1", "viewer-password")

	// F-003: Create a user with wallet.adjust ONLY (no wallet.voucher.issue).
	_ = env.st.Run(context.Background(), func(tx kernel.Tx) error {
		_, _ = tx.Exec(context.Background(), `INSERT INTO roles (id, name, description, is_system, created_at, updated_at) VALUES ('role-adjuster', 'Adjuster', 'Can adjust only', 0, 1000, 1000)`)
		_, _ = tx.Exec(context.Background(), `INSERT INTO role_permissions (role_id, permission_id) VALUES ('role-adjuster', 'wallet.adjust')`)
		return nil
	})
	env.addUser(t, "adjuster1", "adjuster-password", []string{"role-adjuster"})
	adjusterTok := env.login(t, "adjuster1", "adjuster-password")

	// 1. Viewer without permission should get 403.
	body := `{"batchId":"batch-viewer","count":5,"amount":1000,"currency":"CNY"}`
	req := httptest.NewRequest("POST", "/api/wallet/vouchers/batches", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+viewerTok)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	env.mux.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("viewer POST /api/wallet/vouchers/batches status = %d, want 403", w.Code)
	}

	// 2. F-003: User with wallet.adjust ONLY must also get 403 (strict permission isolation).
	req = httptest.NewRequest("POST", "/api/wallet/vouchers/batches", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+adjusterTok)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	env.mux.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("adjuster without voucher.issue POST /api/wallet/vouchers/batches status = %d, want 403", w.Code)
	}

	// 3. Admin with wallet.voucher.issue should succeed (201 Created).
	adminBody := `{"batchId":"batch-admin","count":3,"amount":2000,"currency":"CNY"}`
	req = httptest.NewRequest("POST", "/api/wallet/vouchers/batches", strings.NewReader(adminBody))
	req.Header.Set("Authorization", "Bearer "+adminTok)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	env.mux.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("admin POST /api/wallet/vouchers/batches status = %d, body: %s", w.Code, w.Body.String())
	}

	var res struct {
		Items []struct {
			ID         string `json:"id"`
			BatchID    string `json:"batchId"`
			CodePrefix string `json:"codePrefix"`
			Amount     int64  `json:"amount"`
			Code       string `json:"code"`
			Status     string `json:"status"`
		} `json:"items"`
		Total int `json:"total"`
	}
	if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
		t.Fatal(err)
	}
	if res.Total != 3 || len(res.Items) != 3 {
		t.Fatalf("unexpected items count: %d", res.Total)
	}
	for i, it := range res.Items {
		if len(it.Code) != 24 {
			t.Fatalf("item[%d] code len = %d, want 24", i, len(it.Code))
		}
		if it.Amount != 2000 || it.Status != "unused" {
			t.Fatalf("unexpected item fields: %+v", it)
		}
	}

	// 4. Verify operationlog: audit entry exists and DOES NOT contain plaintext codes!
	ops, _, err := env.operations.ListOperationsFiltered(operationlog.OperationFilter{
		Event: "records.create", Sort: "createdAt", Order: "desc", Page: 1, PageSize: 50,
	})
	if err != nil {
		t.Fatal(err)
	}
	var foundAudit bool
	for _, op := range ops {
		detail := ""
		if op.Detail != nil {
			detail = *op.Detail
		}
		if strings.Contains(detail, "batch-admin") {
			foundAudit = true
			// Assert that none of the generated plaintext codes appear in the audit detail
			for _, it := range res.Items {
				if strings.Contains(detail, it.Code) {
					t.Fatalf("security violation: plaintext code %s leaked into audit log!", it.Code)
				}
			}
		}
	}
	if !foundAudit {
		t.Fatal("expected records.create audit log entry for batch-admin")
	}
}

func TestVouchersListAndGet(t *testing.T) {
	env := newAuthTestEnv(t)
	repo := walletstore.NewRepository(env.st)
	service := newWalletStub(env, repo)
	mountWalletRoutes(t, env, service)

	adminTok := adminToken(t, env)

	// 1. Generate a batch.
	body := `{"batchId":"batch-list-test","count":2,"amount":500,"currency":"CNY"}`
	req := httptest.NewRequest("POST", "/api/wallet/vouchers/batches", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+adminTok)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	env.mux.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("generate batch failed: %d", w.Code)
	}
	var genRes struct {
		Items []struct {
			ID   string `json:"id"`
			Code string `json:"code"`
		} `json:"items"`
	}
	_ = json.NewDecoder(w.Body).Decode(&genRes)
	firstID := genRes.Items[0].ID
	firstCode := genRes.Items[0].Code

	// 2. Query list: GET /api/wallet/vouchers
	req = httptest.NewRequest("GET", "/api/wallet/vouchers?batchId=batch-list-test", nil)
	req.Header.Set("Authorization", "Bearer "+adminTok)
	w = httptest.NewRecorder()
	env.mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list vouchers status = %d", w.Code)
	}
	var listRes struct {
		Items []map[string]any `json:"items"`
		Total int              `json:"total"`
	}
	if err := json.NewDecoder(w.Body).Decode(&listRes); err != nil {
		t.Fatal(err)
	}
	if listRes.Total != 2 {
		t.Fatalf("total = %d, want 2", listRes.Total)
	}
	for _, row := range listRes.Items {
		if _, ok := row["code"]; ok {
			t.Fatal("security violation: list vouchers must never return plaintext code")
		}
		if _, ok := row["codeHash"]; ok {
			t.Fatal("security violation: list vouchers must never return codeHash")
		}
		if _, ok := row["codePrefix"]; !ok {
			t.Fatal("expected codePrefix in list vouchers")
		}
	}

	// 3. Query single voucher: GET /api/wallet/vouchers/{id}
	req = httptest.NewRequest("GET", "/api/wallet/vouchers/"+firstID, nil)
	req.Header.Set("Authorization", "Bearer "+adminTok)
	w = httptest.NewRecorder()
	env.mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("get voucher status = %d", w.Code)
	}
	var single map[string]any
	if err := json.NewDecoder(w.Body).Decode(&single); err != nil {
		t.Fatal(err)
	}
	if single["id"] != firstID || single["codePrefix"] != firstCode[:6] {
		t.Fatalf("unexpected single voucher: %+v", single)
	}
	if _, ok := single["code"]; ok {
		t.Fatal("security violation: single voucher get must never return plaintext code")
	}

	// F-004: Query non-existent voucher -> 404 VOUCHER_NOT_FOUND
	req = httptest.NewRequest("GET", "/api/wallet/vouchers/non-existent-id", nil)
	req.Header.Set("Authorization", "Bearer "+adminTok)
	w = httptest.NewRecorder()
	env.mux.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("get non-existent voucher status = %d, want 404", w.Code)
	}
	var notFoundErr map[string]any
	_ = json.NewDecoder(w.Body).Decode(&notFoundErr)
	if notFoundErr["error"] != "VOUCHER_NOT_FOUND" {
		t.Fatalf("error = %v, want VOUCHER_NOT_FOUND", notFoundErr["error"])
	}
}

func TestVoucherVoidAndConflict(t *testing.T) {
	env := newAuthTestEnv(t)
	repo := walletstore.NewRepository(env.st)
	service := newWalletStub(env, repo)
	mountWalletRoutes(t, env, service)

	env.addUser(t, "viewer1", "viewer-password", []string{"viewer"})
	adminTok := adminToken(t, env)
	viewerTok := env.login(t, "viewer1", "viewer-password")

	// 1. Generate 2 vouchers.
	body := `{"batchId":"batch-void-test","count":2,"amount":300,"currency":"CNY"}`
	req := httptest.NewRequest("POST", "/api/wallet/vouchers/batches", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+adminTok)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	env.mux.ServeHTTP(w, req)
	var genRes struct {
		Items []struct {
			ID   string `json:"id"`
			Code string `json:"code"`
		} `json:"items"`
	}
	_ = json.NewDecoder(w.Body).Decode(&genRes)
	targetID := genRes.Items[0].ID
	redeemID := genRes.Items[1].ID
	redeemCode := genRes.Items[1].Code

	// 2. Viewer trying to void -> 403 Forbidden.
	req = httptest.NewRequest("POST", "/api/wallet/vouchers/"+targetID+"/void", nil)
	req.Header.Set("Authorization", "Bearer "+viewerTok)
	w = httptest.NewRecorder()
	env.mux.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("viewer void status = %d, want 403", w.Code)
	}

	// 3. F-004: Void non-existent voucher -> 404 VOUCHER_NOT_FOUND
	req = httptest.NewRequest("POST", "/api/wallet/vouchers/non-existent-id/void", nil)
	req.Header.Set("Authorization", "Bearer "+adminTok)
	w = httptest.NewRecorder()
	env.mux.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("void non-existent status = %d, want 404", w.Code)
	}

	// 4. Admin void voucher: POST /api/wallet/vouchers/{id}/void -> 200 OK.
	req = httptest.NewRequest("POST", "/api/wallet/vouchers/"+targetID+"/void", nil)
	req.Header.Set("Authorization", "Bearer "+adminTok)
	w = httptest.NewRecorder()
	env.mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("void voucher status = %d, body = %s", w.Code, w.Body.String())
	}

	// Verify status changed to void.
	req = httptest.NewRequest("GET", "/api/wallet/vouchers/"+targetID, nil)
	req.Header.Set("Authorization", "Bearer "+adminTok)
	w = httptest.NewRecorder()
	env.mux.ServeHTTP(w, req)
	var single map[string]any
	_ = json.NewDecoder(w.Body).Decode(&single)
	if single["status"] != "void" {
		t.Fatalf("status = %v, want void", single["status"])
	}

	// 5. F-004: Redeem the second voucher, then try to void it -> 409 VOUCHER_ALREADY_REDEEMED
	subStore := subject.NewStore(env.st)
	sub, _, err := subStore.GetOrCreateSubject(context.Background(), "test", "sub-void", time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.vouchers.Redeem(context.Background(), sub.ID, redeemCode, time.Now().UTC()); err != nil {
		t.Fatalf("redeem: %v", err)
	}

	req = httptest.NewRequest("POST", "/api/wallet/vouchers/"+redeemID+"/void", nil)
	req.Header.Set("Authorization", "Bearer "+adminTok)
	w = httptest.NewRecorder()
	env.mux.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Fatalf("void redeemed voucher status = %d, want 409", w.Code)
	}
	var conflictErr map[string]any
	_ = json.NewDecoder(w.Body).Decode(&conflictErr)
	if conflictErr["error"] != "VOUCHER_ALREADY_REDEEMED" {
		t.Fatalf("error = %v, want VOUCHER_ALREADY_REDEEMED", conflictErr["error"])
	}
}

func TestVoucherInvalidBodyAndParams(t *testing.T) {
	env := newAuthTestEnv(t)
	repo := walletstore.NewRepository(env.st)
	service := newWalletStub(env, repo)
	mountWalletRoutes(t, env, service)

	adminTok := adminToken(t, env)

	// 1. Invalid JSON body -> 400 INVALID_VOUCHER_BODY
	req := httptest.NewRequest("POST", "/api/wallet/vouchers/batches", strings.NewReader(`not-json`))
	req.Header.Set("Authorization", "Bearer "+adminTok)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	env.mux.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	var bodyErr map[string]any
	_ = json.NewDecoder(w.Body).Decode(&bodyErr)
	if bodyErr["error"] != "INVALID_VOUCHER_BODY" {
		t.Fatalf("error = %v, want INVALID_VOUCHER_BODY", bodyErr["error"])
	}

	// 2. Invalid count / amount -> 400 INVALID_VOUCHER_PARAMS
	req = httptest.NewRequest("POST", "/api/wallet/vouchers/batches", strings.NewReader(`{"batchId":"b","count":-1,"amount":100}`))
	req.Header.Set("Authorization", "Bearer "+adminTok)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	env.mux.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	var paramErr map[string]any
	_ = json.NewDecoder(w.Body).Decode(&paramErr)
	if paramErr["error"] != "INVALID_VOUCHER_PARAMS" {
		t.Fatalf("error = %v, want INVALID_VOUCHER_PARAMS", paramErr["error"])
	}
}

func TestVoucherSchemaRegistration(t *testing.T) {
	docBytes := walletschema.SchemaDocuments()["wallet-vouchers"]
	if len(docBytes) == 0 {
		t.Fatal("wallet-vouchers schema document not found")
	}
	var doc struct {
		Meta struct {
			PageID string `json:"pageId"`
		} `json:"meta"`
		Actions map[string]any `json:"actions"`
		Body    struct {
			Children []struct {
				Props struct {
					Toolbar []map[string]any `json:"toolbar"`
				} `json:"props"`
			} `json:"children"`
		} `json:"body"`
	}
	if err := json.Unmarshal(docBytes, &doc); err != nil {
		t.Fatalf("unmarshal wallet-vouchers schema: %v", err)
	}
	if doc.Meta.PageID != "wallet-vouchers" {
		t.Fatalf("pageId = %s, want wallet-vouchers", doc.Meta.PageID)
	}
	if _, ok := doc.Actions["generateBatch"]; !ok {
		t.Fatal("missing generateBatch action in schema")
	}
	if _, ok := doc.Actions["openGenerate"]; !ok {
		t.Fatal("missing openGenerate action in schema")
	}
	if _, ok := doc.Actions["voidVoucher"]; !ok {
		t.Fatal("missing voidVoucher action in schema")
	}
}
