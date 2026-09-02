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
	walletschema "github.com/magicvr/schema-ui-core/apps/api/modules/wallet/schema"
	walletstore "github.com/magicvr/schema-ui-core/apps/api/modules/wallet/store"
	"github.com/magicvr/schema-ui-core/apps/api/modules/wallet/subject"
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
	body := `{"batchId":"batch-viewer","count":5,"amount":"10.00","currency":"CNY"}`
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
	adminBody := `{"batchId":"batch-admin","count":3,"amount":"20.00","currency":"CNY"}`
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
	body := `{"batchId":"batch-list-test","count":2,"amount":"5.00","currency":"CNY"}`
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
		if row["voidable"] != true {
			t.Fatalf("unused voucher voidable = %v, want true", row["voidable"])
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
	if single["voidable"] != true {
		t.Fatalf("unused voucher get voidable = %v, want true", single["voidable"])
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
	body := `{"batchId":"batch-void-test","count":2,"amount":"3.00","currency":"CNY"}`
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
	if single["voidable"] != false {
		t.Fatalf("voided voucher voidable = %v, want false", single["voidable"])
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

	req = httptest.NewRequest("GET", "/api/wallet/vouchers/"+redeemID, nil)
	req.Header.Set("Authorization", "Bearer "+adminTok)
	w = httptest.NewRecorder()
	env.mux.ServeHTTP(w, req)
	var redeemed map[string]any
	_ = json.NewDecoder(w.Body).Decode(&redeemed)
	if redeemed["status"] != "redeemed" {
		t.Fatalf("redeemed status = %v, want redeemed", redeemed["status"])
	}
	if redeemed["voidable"] != false {
		t.Fatalf("redeemed voucher voidable = %v, want false", redeemed["voidable"])
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
	req = httptest.NewRequest("POST", "/api/wallet/vouchers/batches", strings.NewReader(`{"batchId":"b","count":-1,"amount":"1.00"}`))
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

	// 3. F-002: Invalid Currency USD -> 400 INVALID_VOUCHER_PARAMS
	req = httptest.NewRequest("POST", "/api/wallet/vouchers/batches", strings.NewReader(`{"batchId":"b","count":1,"amount":"1.00","currency":"USD"}`))
	req.Header.Set("Authorization", "Bearer "+adminTok)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	env.mux.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("USD currency status = %d, want 400", w.Code)
	}
	var currErr map[string]any
	_ = json.NewDecoder(w.Body).Decode(&currErr)
	if currErr["error"] != "INVALID_VOUCHER_PARAMS" {
		t.Fatalf("error = %v, want INVALID_VOUCHER_PARAMS", currErr["error"])
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
					Actions []map[string]any `json:"actions"`
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
	var voidAction map[string]any
	for _, child := range doc.Body.Children {
		for _, action := range child.Props.Actions {
			if action["key"] == "void" {
				voidAction = action
			}
		}
	}
	if voidAction == nil {
		t.Fatal("missing void row action in wallet-vouchers table")
	}
	if voidAction["visibleField"] != "voidable" {
		t.Fatalf("void visibleField = %v, want voidable", voidAction["visibleField"])
	}
}

// A-005 F-004 (A-008): the 0065 batch registry turns a repeated batchId into a
// 409 VOUCHER_BATCH_EXISTS conflict and never mixes a second list into the
// existing batch.
func TestVouchersBatchDuplicateConflict(t *testing.T) {
	env := newAuthTestEnv(t)
	repo := walletstore.NewRepository(env.st)
	service := newWalletStub(env, repo)
	mountWalletRoutes(t, env, service)
	adminTok := adminToken(t, env)

	post := func(body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest("POST", "/api/wallet/vouchers/batches", strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+adminTok)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		env.mux.ServeHTTP(w, req)
		return w
	}
	listTotal := func() int {
		req := httptest.NewRequest("GET", "/api/wallet/vouchers?batchId=batch-dup-h", nil)
		req.Header.Set("Authorization", "Bearer "+adminTok)
		w := httptest.NewRecorder()
		env.mux.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("list status = %d", w.Code)
		}
		var list struct {
			Total int `json:"total"`
		}
		if err := json.NewDecoder(w.Body).Decode(&list); err != nil {
			t.Fatal(err)
		}
		return list.Total
	}

	if w := post(`{"batchId":"batch-dup-h","count":2,"amount":"5.00","currency":"CNY"}`); w.Code != http.StatusCreated {
		t.Fatalf("first generate = %d %s", w.Code, w.Body.String())
	}

	w := post(`{"batchId":"batch-dup-h","count":1,"amount":"5.00","currency":"CNY"}`)
	if w.Code != http.StatusConflict {
		t.Fatalf("duplicate generate status = %d, want 409; body=%s", w.Code, w.Body.String())
	}
	var errBody map[string]any
	_ = json.NewDecoder(w.Body).Decode(&errBody)
	if errBody["error"] != "VOUCHER_BATCH_EXISTS" {
		t.Fatalf("error = %v, want VOUCHER_BATCH_EXISTS", errBody["error"])
	}
	if total := listTotal(); total != 2 {
		t.Fatalf("list total after rejected duplicate = %d, want 2 (no mixed rows)", total)
	}
}

// A-005 F-005 (A-008): expiresAt must be Unix SECONDS in 2001-09-09..2100 —
// millisecond timestamps and out-of-range values fail closed with 400 instead
// of silently minting immediate/never expiry.
func TestVouchersGenerateExpiresAtValidation(t *testing.T) {
	env := newAuthTestEnv(t)
	repo := walletstore.NewRepository(env.st)
	service := newWalletStub(env, repo)
	mountWalletRoutes(t, env, service)
	adminTok := adminToken(t, env)

	post := func(body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest("POST", "/api/wallet/vouchers/batches", strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+adminTok)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		env.mux.ServeHTTP(w, req)
		return w
	}

	// 1. Millisecond timestamp (2027-01-15T10:00:00 in ms) → 400.
	w := post(`{"batchId":"b-exp-ms","count":1,"amount":"1.00","currency":"CNY","expiresAt":1800000000000}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("milliseconds expiresAt status = %d, want 400; body=%s", w.Code, w.Body.String())
	}
	var msErr map[string]any
	_ = json.NewDecoder(w.Body).Decode(&msErr)
	if msErr["error"] != "INVALID_VOUCHER_PARAMS" {
		t.Fatalf("error = %v, want INVALID_VOUCHER_PARAMS", msErr["error"])
	}

	// 2. Pre-epoch-ish / ancient seconds value → 400.
	w = post(`{"batchId":"b-exp-old","count":1,"amount":"1.00","currency":"CNY","expiresAt":999999999}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("out-of-range expiresAt status = %d, want 400; body=%s", w.Code, w.Body.String())
	}

	// 3. Valid future seconds (2027-01-15T10:00:00Z) → 201.
	w = post(`{"batchId":"b-exp-ok","count":1,"amount":"1.00","currency":"CNY","expiresAt":1800000000}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("valid expiresAt status = %d, want 201; body=%s", w.Code, w.Body.String())
	}
}

// E-008: batchId is optional — the server generates a unique VB-… id; two
// submissions without batchId produce distinct batches.
func TestVouchersBatchGenerateAutoBatchID(t *testing.T) {
	env := newAuthTestEnv(t)
	repo := walletstore.NewRepository(env.st)
	service := newWalletStub(env, repo)
	mountWalletRoutes(t, env, service)
	adminTok := adminToken(t, env)

	post := func(body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest("POST", "/api/wallet/vouchers/batches", strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+adminTok)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		env.mux.ServeHTTP(w, req)
		return w
	}
	decodeBatchID := func(w *httptest.ResponseRecorder) string {
		t.Helper()
		if w.Code != http.StatusCreated {
			t.Fatalf("generate status = %d, body=%s", w.Code, w.Body.String())
		}
		var res struct {
			Items []struct {
				BatchID string `json:"batchId"`
			} `json:"items"`
		}
		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}
		if len(res.Items) == 0 || res.Items[0].BatchID == "" {
			t.Fatalf("no batchId in auto-generated response: %+v", res)
		}
		for _, item := range res.Items {
			if !strings.HasPrefix(item.BatchID, "VB-") {
				t.Fatalf("auto batchId %q does not start with VB-", item.BatchID)
			}
		}
		return res.Items[0].BatchID
	}

	// 1. No batchId → server-generated VB-… id.
	first := decodeBatchID(post(`{"count":2,"amount":"5.00","currency":"CNY"}`))

	// 2. A second submission gets a DIFFERENT auto id.
	second := decodeBatchID(post(`{"count":1,"amount":"2.50","currency":"CNY"}`))
	if first == second {
		t.Fatalf("auto batch ids collide: %s", first)
	}

	// 3. The generated batch is listable under its auto id.
	req := httptest.NewRequest("GET", "/api/wallet/vouchers?batchId="+second, nil)
	req.Header.Set("Authorization", "Bearer "+adminTok)
	list := httptest.NewRecorder()
	env.mux.ServeHTTP(list, req)
	if list.Code != http.StatusOK {
		t.Fatalf("list by auto batchId status = %d", list.Code)
	}
	var listRes struct {
		Total int `json:"total"`
	}
	if err := json.NewDecoder(list.Body).Decode(&listRes); err != nil {
		t.Fatal(err)
	}
	if listRes.Total != 1 {
		t.Fatalf("auto batch list total = %d, want 1", listRes.Total)
	}
}

// E-008: amount is CNY yuan with up to 2 decimal places (JSON number or
// string); it is converted to min units (分) internally — the returned voucher
// amounts stay cents (the table column formats them as currency).
func TestVouchersGenerateAmountYuanDecimals(t *testing.T) {
	env := newAuthTestEnv(t)
	repo := walletstore.NewRepository(env.st)
	service := newWalletStub(env, repo)
	mountWalletRoutes(t, env, service)
	adminTok := adminToken(t, env)

	post := func(body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest("POST", "/api/wallet/vouchers/batches", strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+adminTok)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		env.mux.ServeHTTP(w, req)
		return w
	}
	firstAmount := func(w *httptest.ResponseRecorder) float64 {
		t.Helper()
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
		}
		var res struct {
			Items []struct {
				Amount int64 `json:"amount"`
			} `json:"items"`
		}
		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}
		return float64(res.Items[0].Amount)
	}

	cases := []struct {
		body string
		want float64 // cents
	}{
		{`{"count":1,"amount":12.5}`, 1250},    // JSON number, one decimal
		{`{"count":1,"amount":"12.55"}`, 1255}, // string, two decimals
		{`{"count":1,"amount":0.5}`, 50},       // fractional yuan still > 0
		{`{"count":1,"amount":"88.00"}`, 8800}, // trailing zeros
		{`{"count":1,"amount":12.50}`, 1250},   // number with trailing zero digit
		{`{"count":1,"amount":"0.01"}`, 1},     // smallest unit
		{`{"count":1,"amount":1e2}`, 10000},    // JSON exponent notation → 100 yuan
	}
	for _, tc := range cases {
		if got := firstAmount(post(tc.body)); got != tc.want {
			t.Fatalf("amount %s = %v cents, want %v", tc.body, got, tc.want)
		}
	}

	rejects := []string{
		`{"count":1}`,                  // amount missing
		`{"count":1,"amount":0}`,       // zero
		`{"count":1,"amount":"0"}`,     // zero string
		`{"count":1,"amount":"1.234"}`, // > 2 decimals
		`{"count":1,"amount":"abc"}`,   // non-numeric
		`{"count":1,"amount":"-1.00"}`, // negative
	}
	for _, body := range rejects {
		w := post(body)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("amount %s status = %d, want 400; body=%s", body, w.Code, w.Body.String())
		}
		var errBody map[string]any
		_ = json.NewDecoder(w.Body).Decode(&errBody)
		if errBody["error"] != "INVALID_VOUCHER_PARAMS" {
			t.Fatalf("amount %s error = %v, want INVALID_VOUCHER_PARAMS", body, errBody["error"])
		}
	}
}

// E-009: expiresAt accepts a UTC date (YYYY-MM-DD) from the admin datePicker;
// the server converts it to end-of-day Unix seconds (23:59:59 UTC of the
// chosen day), so the whole day stays redeemable. Legacy Unix seconds input
// keeps working.
func TestVouchersGenerateExpiresAtDatePicker(t *testing.T) {
	env := newAuthTestEnv(t)
	repo := walletstore.NewRepository(env.st)
	service := newWalletStub(env, repo)
	mountWalletRoutes(t, env, service)
	adminTok := adminToken(t, env)

	post := func(body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest("POST", "/api/wallet/vouchers/batches", strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+adminTok)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		env.mux.ServeHTTP(w, req)
		return w
	}
	firstExpiry := func(w *httptest.ResponseRecorder) (string, bool) {
		t.Helper()
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, want 201; body=%s", w.Code, w.Body.String())
		}
		var res struct {
			Items []map[string]any `json:"items"`
		}
		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}
		exp, ok := res.Items[0]["expiresAt"]
		if !ok {
			return "", false
		}
		return exp.(string), true
	}

	// 1. Chosen date → expiry at 23:59:59 UTC of the SAME day.
	exp, ok := firstExpiry(post(`{"count":1,"amount":"1.00","expiresAt":"2027-01-15"}`))
	if !ok || exp != "2027-01-15T23:59:59.000Z" {
		t.Fatalf("date expiry = %q (ok=%v), want 2027-01-15T23:59:59.000Z", exp, ok)
	}

	// 2. Upper-edge date (2099-12-31) is accepted; the day after the window is not.
	if exp, ok := firstExpiry(post(`{"count":1,"amount":"1.00","expiresAt":"2099-12-31"}`)); !ok || exp != "2099-12-31T23:59:59.000Z" {
		t.Fatalf("2099-12-31 expiry = %q (ok=%v)", exp, ok)
	}
	rejects := []string{
		`{"count":1,"amount":"1.00","expiresAt":"2100-01-01"}`,
		`{"count":1,"amount":"1.00","expiresAt":"2001-09-08"}`,
		`{"count":1,"amount":"1.00","expiresAt":"2027-1-5"}`,             // not zero padded
		`{"count":1,"amount":"1.00","expiresAt":"2027-13-01"}`,           // invalid month
		`{"count":1,"amount":"1.00","expiresAt":"2027-01-15T10:00:00Z"}`, // time not accepted in date mode
	}
	for _, body := range rejects {
		w := post(body)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expiry %s status = %d, want 400; body=%s", body, w.Code, w.Body.String())
		}
		var errBody map[string]any
		_ = json.NewDecoder(w.Body).Decode(&errBody)
		if errBody["error"] != "INVALID_VOUCHER_PARAMS" {
			t.Fatalf("expiry %s error = %v, want INVALID_VOUCHER_PARAMS", body, errBody["error"])
		}
	}

	// 3. Empty date string = omitted (no expiry).
	if _, ok := firstExpiry(post(`{"count":1,"amount":"1.00","expiresAt":""}`)); ok {
		t.Fatal("empty expiry date must behave like omitted (no expiresAt in response)")
	}
}
