package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	"github.com/magicvr/schema-ui-core/apps/api/internal/requestid"
)

func TestBrandingPublicAndSettingsPatch(t *testing.T) {
	env := newAuthTestEnv(t)
	withRequestID := requestid.Middleware(env.mux)

	// Public branding without auth.
	req := httptest.NewRequest(http.MethodGet, "/api/branding", nil)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("branding status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	var branding map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&branding); err != nil {
		t.Fatalf("branding json: %v", err)
	}
	if branding["siteTitle"] != "Schema UI Core" {
		t.Fatalf("default siteTitle = %v", branding["siteTitle"])
	}
	if branding["logoUrl"] != "" {
		t.Fatalf("default logoUrl = %v, want empty", branding["logoUrl"])
	}

	token := env.login(t, testSeedUsername, testSeedPassword)

	// List requires settings.read
	req = bearer(t, token, http.MethodGet, "/api/settings", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("settings list status = %d: %s", rr.Code, rr.Body.String())
	}

	// Empty title rejected
	req = bearer(t, token, http.MethodPatch, "/api/settings/default", `{"siteTitle":"   "}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("empty title status = %d, want 400: %s", rr.Code, rr.Body.String())
	}

	// Valid patch
	req = bearer(t, token, http.MethodPatch, "/api/settings/default",
		`{"siteTitle":"Acme Admin","logoUrl":"https://example.com/logo.png"}`)
	req.Header.Set(requestid.HeaderName, "r2-settings-001")
	rr = httptest.NewRecorder()
	withRequestID.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("patch status = %d: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("X-Schema-UI-Config-Changed"); got != "settings.branding" {
		t.Fatalf("config change namespace = %q, want settings.branding", got)
	}
	var row map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&row)
	if row["siteTitle"] != "Acme Admin" || row["logoUrl"] != "https://example.com/logo.png" {
		t.Fatalf("patched row = %v", row)
	}

	// Branding reflects update
	req = httptest.NewRequest(http.MethodGet, "/api/branding", nil)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	_ = json.NewDecoder(rr.Body).Decode(&branding)
	if branding["siteTitle"] != "Acme Admin" {
		t.Fatalf("branding after patch = %v", branding)
	}

	// Clear logo
	req = bearer(t, token, http.MethodPatch, "/api/settings/default", `{"logoUrl":""}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("clear logo status = %d: %s", rr.Code, rr.Body.String())
	}

	// A-006 R-003: settings PATCH appends settings.update to the operation log.
	ops, total, err := env.operations.ListOperationsFiltered(operationlog.OperationFilter{
		Q: "settings.update", Sort: "createdAt", Order: "desc", Page: 1, PageSize: 20,
	})
	if err != nil {
		t.Fatalf("list operations: %v", err)
	}
	if total < 1 {
		t.Fatal("expected at least one settings.update operation after patch")
	}
	found := false
	for _, op := range ops {
		if op.Event == operationlog.EventSettingsUpdate && op.CorrelationID == "r2-settings-001" {
			found = true
			if op.RecordID == nil || *op.RecordID != "default" {
				t.Fatalf("settings.update record_id = %v, want default", op.RecordID)
			}
			if op.Detail == nil {
				t.Fatal("settings.update missing structured detail")
			}
			detail, parseErr := operationlog.ParseDetail(*op.Detail)
			if parseErr != nil || detail.SchemaVersion != operationlog.DetailSchemaVersion {
				t.Fatalf("settings.update detail = %+v, err=%v", detail, parseErr)
			}
			if got := detail.After["logoUrl"]; got != operationlog.RedactedValue {
				t.Fatalf("settings logoUrl = %v, want redacted", got)
			}
			break
		}
	}
	if !found {
		t.Fatalf("settings.update with correlation r2-settings-001 not found in operations: %+v", ops)
	}
}

func TestOperationsListAndDetailReadOnly(t *testing.T) {
	env := newAuthTestEnv(t)
	token := env.login(t, testSeedUsername, testSeedPassword)

	// Login already wrote auth.login operation
	req := bearer(t, token, http.MethodGet, "/api/operations?pageSize=20", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("operations list status = %d: %s", rr.Code, rr.Body.String())
	}
	var list map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&list)
	items, _ := list["items"].([]any)
	if len(items) == 0 {
		t.Fatal("expected at least one operation from login")
	}
	first, _ := items[0].(map[string]any)
	id, _ := first["id"].(string)
	if id == "" {
		t.Fatalf("first op missing id: %v", first)
	}

	req = bearer(t, token, http.MethodGet, "/api/operations/"+id, "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("detail status = %d: %s", rr.Code, rr.Body.String())
	}

	// Write routes must not be mounted
	req = bearer(t, token, http.MethodPost, "/api/operations", `{}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code == http.StatusOK || rr.Code == http.StatusCreated {
		t.Fatalf("POST operations should not succeed, status = %d", rr.Code)
	}

	req = bearer(t, token, http.MethodDelete, "/api/operations/"+id, "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code == http.StatusNoContent || rr.Code == http.StatusOK {
		t.Fatalf("DELETE operations should not succeed, status = %d", rr.Code)
	}
}

func TestSettingsWriteRequiresPermission(t *testing.T) {
	env := newAuthTestEnv(t)
	// editor seed has no settings.write
	env.addUser(t, "ed", "editor-pass-1", []string{"editor"})
	// need password set - check addUser
	tok := loginAs(t, env, "ed", "editor-pass-1")
	req := bearer(t, tok, http.MethodPatch, "/api/settings/default", `{"siteTitle":"Nope"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("editor patch status = %d, want 403: %s", rr.Code, rr.Body.String())
	}
}

func loginAs(t *testing.T, env *authTestEnv, username, password string) string {
	t.Helper()
	return env.login(t, username, password)
}

func TestBrandingVp007StartupFieldsAndPatch(t *testing.T) {
	env := newAuthTestEnv(t)
	token := env.login(t, testSeedUsername, testSeedPassword)

	// Patch the VP-007 fields across all four categories.
	req := bearer(t, token, http.MethodPatch, "/api/settings/default",
		`{"logoUrlLight":"/assets/logo-light.svg","logoUrlDark":"/assets/logo-dark.svg","faviconUrl":"/favicon.ico","defaultLocale":"zh-CN","siteTimezone":"Asia/Shanghai","defaultTheme":"dark"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("vp007 patch status = %d: %s", rr.Code, rr.Body.String())
	}
	var row map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&row)
	if row["defaultLocale"] != "zh-CN" || row["defaultTheme"] != "dark" || row["siteTimezone"] != "Asia/Shanghai" {
		t.Fatalf("patched row = %v", row)
	}
	if row["logoUrlLight"] != "/assets/logo-light.svg" || row["logoUrlDark"] != "/assets/logo-dark.svg" || row["faviconUrl"] != "/favicon.ico" {
		t.Fatalf("branding row = %v", row)
	}

	// Public startup configuration reflects the fields + supported locales.
	req = httptest.NewRequest(http.MethodGet, "/api/branding", nil)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("branding status = %d", rr.Code)
	}
	var branding map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&branding)
	if branding["defaultLocale"] != "zh-CN" || branding["defaultTheme"] != "dark" || branding["siteTimezone"] != "Asia/Shanghai" {
		t.Fatalf("branding startup = %v", branding)
	}
	if branding["logoUrlLight"] != "/assets/logo-light.svg" || branding["logoUrlDark"] != "/assets/logo-dark.svg" || branding["faviconUrl"] != "/favicon.ico" {
		t.Fatalf("branding assets = %v", branding)
	}
	locales, ok := branding["supportedLocales"].([]any)
	if !ok || len(locales) != 2 {
		t.Fatalf("supportedLocales = %v", branding["supportedLocales"])
	}
}

func TestSettingsValidationAndReset(t *testing.T) {
	env := newAuthTestEnv(t)
	token := env.login(t, testSeedUsername, testSeedPassword)

	// Invalid IANA timezone → 400 INVALID_TIMEZONE, previous value untouched.
	req := bearer(t, token, http.MethodPatch, "/api/settings/default", `{"siteTimezone":"Foo/Bar"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("bad timezone status = %d, want 400", rr.Code)
	}
	var errBody map[string]string
	_ = json.NewDecoder(rr.Body).Decode(&errBody)
	if errBody["error"] != "INVALID_TIMEZONE" {
		t.Fatalf("bad timezone code = %v, want INVALID_TIMEZONE", errBody["error"])
	}

	req = bearer(t, token, http.MethodPatch, "/api/settings/default", `{"defaultLocale":"fr-FR"}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest || !bodyHasCode(rr, "INVALID_DEFAULT_LOCALE") {
		t.Fatalf("bad locale = %d %s", rr.Code, rr.Body.String())
	}

	req = bearer(t, token, http.MethodPatch, "/api/settings/default", `{"defaultTheme":"neon"}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest || !bodyHasCode(rr, "INVALID_DEFAULT_THEME") {
		t.Fatalf("bad theme = %d %s", rr.Code, rr.Body.String())
	}

	// Reset restores defaults and fires the config-change header.
	req = bearer(t, token, http.MethodPost, "/api/settings/default/reset", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("reset status = %d: %s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("X-Schema-UI-Config-Changed"); got != "settings.branding" {
		t.Fatalf("reset config namespace = %q", got)
	}
	var row map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&row)
	if row["siteTitle"] != "Schema UI Core" || row["defaultLocale"] != "auto" || row["defaultTheme"] != "auto" || row["siteTimezone"] != "auto" {
		t.Fatalf("reset row = %v", row)
	}
	if row["logoUrlLight"] != "" || row["logoUrlDark"] != "" || row["faviconUrl"] != "" {
		t.Fatalf("reset assets = %v", row)
	}

	// Reset requires settings.write.
	env.addUser(t, "ed", "editor-pass-1", []string{"editor"})
	editorToken := env.login(t, "ed", "editor-pass-1")
	req = bearer(t, editorToken, http.MethodPost, "/api/settings/default/reset", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("editor reset status = %d, want 403", rr.Code)
	}
}

func bodyHasCode(rr *httptest.ResponseRecorder, want string) bool {
	var body map[string]string
	_ = json.NewDecoder(rr.Body).Decode(&body)
	return body["error"] == want
}
