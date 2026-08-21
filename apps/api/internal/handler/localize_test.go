package handler

// Server-side locale negotiation tests (VP-007 S4 · C2/C3):
// cataloged codes negotiate zh-CN/en-US messages + Content-Language;
// uncataloged codes stay English with no messageKey and no diagnostics leak.

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func requestWithLanguage(lang string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/api/branding", nil)
	if lang != "" {
		req.Header.Set("Accept-Language", lang)
	}
	return req
}

func decodeErrorBody(t *testing.T, rr *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return body
}

func TestWriteLocalizedErrorNegotiatesLanguage(t *testing.T) {
	// zh-CN
	rr := httptest.NewRecorder()
	writeLocalizedError(rr, requestWithLanguage("zh-CN"), http.StatusBadRequest, "INVALID_SITE_TITLE", "siteTitle must not be empty")
	if rr.Header().Get("Content-Language") != "zh-CN" {
		t.Fatalf("Content-Language = %q, want zh-CN", rr.Header().Get("Content-Language"))
	}
	body := decodeErrorBody(t, rr)
	if body["error"] != "INVALID_SITE_TITLE" {
		t.Fatalf("error code = %v", body["error"])
	}
	if body["message"] != "站点标题不能为空" {
		t.Fatalf("zh message = %v", body["message"])
	}
	if body["messageKey"] != "error.invalidSiteTitle" {
		t.Fatalf("messageKey = %v", body["messageKey"])
	}

	// en-US
	rr = httptest.NewRecorder()
	writeLocalizedError(rr, requestWithLanguage("en-US"), http.StatusBadRequest, "INVALID_SITE_TITLE", "siteTitle must not be empty")
	if rr.Header().Get("Content-Language") != "en-US" {
		t.Fatalf("Content-Language = %q", rr.Header().Get("Content-Language"))
	}
	if body := decodeErrorBody(t, rr); body["message"] != "siteTitle must not be empty" {
		t.Fatalf("en message = %v", body["message"])
	}

	// Accept-Language prefix matching (zh / en-GB) + fallback.
	rr = httptest.NewRecorder()
	writeLocalizedError(rr, requestWithLanguage("zh"), http.StatusBadRequest, "INVALID_TIMEZONE", "tz")
	if body := decodeErrorBody(t, rr); body["message"] != "默认时区必须是 auto 或有效的 IANA 时区" {
		t.Fatalf("zh prefix message = %v", body["message"])
	}
	rr = httptest.NewRecorder()
	writeLocalizedError(rr, requestWithLanguage("en-GB"), http.StatusBadRequest, "INVALID_TIMEZONE", "tz")
	if body := decodeErrorBody(t, rr); body["message"] != "siteTimezone must be auto or a valid IANA timezone" {
		t.Fatalf("en-GB message = %v", body["message"])
	}
	rr = httptest.NewRecorder()
	writeLocalizedError(rr, requestWithLanguage("fr-FR, de-DE"), http.StatusBadRequest, "INVALID_TIMEZONE", "tz")
	if body := decodeErrorBody(t, rr); body["message"] != "siteTimezone must be auto or a valid IANA timezone" {
		t.Fatalf("unsupported fallback = %v", body["message"])
	}
	rr = httptest.NewRecorder()
	writeLocalizedError(rr, requestWithLanguage(""), http.StatusBadRequest, "INVALID_TIMEZONE", "tz")
	if body := decodeErrorBody(t, rr); body["message"] != "siteTimezone must be auto or a valid IANA timezone" {
		t.Fatalf("empty header fallback = %v", body["message"])
	}
}

func TestWriteLocalizedErrorUncatalogedStaysEnglishWithoutKey(t *testing.T) {
	for _, lang := range []string{"zh-CN", "en-US", ""} {
		rr := httptest.NewRecorder()
		writeLocalizedError(rr, requestWithLanguage(lang), http.StatusInternalServerError, "INTERNAL", "could not load branding")
		body := decodeErrorBody(t, rr)
		if body["error"] != "INTERNAL" {
			t.Fatalf("code = %v", body["error"])
		}
		if body["message"] != "could not load branding" {
			t.Fatalf("INTERNAL message must stay English generic, got %v", body["message"])
		}
		if _, ok := body["messageKey"]; ok {
			t.Fatalf("INTERNAL must not carry messageKey: %v", body)
		}
		if _, ok := body["params"]; ok {
			t.Fatalf("INTERNAL must not carry params: %v", body)
		}
	}
}

func TestLoginFailureLocalizedEndToEnd(t *testing.T) {
	env := newAuthTestEnv(t)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"username":"admin","password":"wrong"}`))
	req.Header.Set("Accept-Language", "zh-CN")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("login status = %d", rr.Code)
	}
	body := decodeErrorBody(t, rr)
	if body["error"] != "UNAUTHORIZED" {
		t.Fatalf("code = %v", body["error"])
	}
	if body["message"] != "用户名或密码错误" {
		t.Fatalf("zh login message = %v", body["message"])
	}
	if rr.Header().Get("Content-Language") != "zh-CN" {
		t.Fatalf("Content-Language = %q", rr.Header().Get("Content-Language"))
	}
}

func TestSettingsValidationLocalizedEndToEnd(t *testing.T) {
	env := newAuthTestEnv(t)
	token := env.login(t, testSeedUsername, testSeedPassword)
	req := bearer(t, token, http.MethodPatch, "/api/settings/default", `{"siteTimezone":"Foo/Bar"}`)
	req.Header.Set("Accept-Language", "zh-CN")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", rr.Code)
	}
	body := decodeErrorBody(t, rr)
	if body["error"] != "INVALID_TIMEZONE" {
		t.Fatalf("code = %v", body["error"])
	}
	if body["message"] != "默认时区必须是 auto 或有效的 IANA 时区" {
		t.Fatalf("zh message = %v", body["message"])
	}
	if body["messageKey"] != "error.invalidTimezone" {
		t.Fatalf("messageKey = %v", body["messageKey"])
	}
}

func TestAuthMiddlewareLocalized(t *testing.T) {
	env := newAuthTestEnv(t)
	req := httptest.NewRequest(http.MethodGet, "/api/settings", nil)
	req.Header.Set("Accept-Language", "zh-CN")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", rr.Code)
	}
	body := decodeErrorBody(t, rr)
	if body["error"] != "UNAUTHENTICATED" {
		t.Fatalf("code = %v", body["error"])
	}
	if body["message"] != "未登录或会话已失效" {
		t.Fatalf("zh message = %v", body["message"])
	}
}
