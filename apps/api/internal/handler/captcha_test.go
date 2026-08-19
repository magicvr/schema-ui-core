// Login captcha handler tests (S-11 · GOAL-011 D-002 §2/§3): the public
// preflight, the admin settings surface, the login gate (default off; enabled
// → INVALID_CAPTCHA without a valid challenge; correct challenge passes) and
// the captcha.settings-update audit event.
package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCaptchaPreflightDisabledByDefault(t *testing.T) {
	env := newAuthTestEnv(t)
	rec := httptest.NewRecorder()
	env.mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/auth/captcha", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("preflight = %d: %s", rec.Code, rec.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["enabled"] != false {
		t.Fatalf("enabled = %v, want false (default off, D-001 §5)", body["enabled"])
	}
	if _, has := body["challenge"]; has {
		t.Fatalf("challenge must be absent while the gate is off: %v", body)
	}
}

func TestCaptchaPreflightEnabledServesChallenge(t *testing.T) {
	env := newAuthTestEnv(t)
	if err := env.captcha.SetEnabled(true, time.Now().UTC()); err != nil {
		t.Fatalf("enable: %v", err)
	}
	rec := httptest.NewRecorder()
	env.mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/auth/captcha", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("preflight = %d: %s", rec.Code, rec.Body.String())
	}
	var body struct {
		Enabled   bool `json:"enabled"`
		Challenge *struct {
			ID               string `json:"id"`
			Question         string `json:"question"`
			ExpiresInSeconds int64  `json:"expiresInSeconds"`
		} `json:"challenge"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !body.Enabled || body.Challenge == nil {
		t.Fatalf("body = %+v, want enabled challenge", body)
	}
	if body.Challenge.ID == "" || body.Challenge.Question == "" || body.Challenge.ExpiresInSeconds <= 0 {
		t.Fatalf("challenge = %+v, want id/question/expiry", body.Challenge)
	}
}

func TestCaptchaPreflightRateLimited(t *testing.T) {
	env := newAuthTestEnv(t)
	if err := env.captcha.SetEnabled(true, time.Now().UTC()); err != nil {
		t.Fatalf("enable: %v", err)
	}
	old := captchaGenerateLimiter
	captchaGenerateLimiter = newLoginRateLimiter(time.Minute, 10, 1<<16)
	defer func() { captchaGenerateLimiter = old }()

	for i := 0; i < 10; i++ {
		rec := httptest.NewRecorder()
		env.mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/auth/captcha", nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("preflight %d = %d: %s", i, rec.Code, rec.Body.String())
		}
	}
	// The 11th request within the window is rejected (W7 F-006: generation is
	// actually counted, not just checked).
	rec := httptest.NewRecorder()
	env.mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/auth/captcha", nil))
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("11th preflight = %d, want 429: %s", rec.Code, rec.Body.String())
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["error"] != "RATE_LIMITED" {
		t.Fatalf("error = %q, want RATE_LIMITED", body["error"])
	}
}

func TestCaptchaSettingsReadWrite(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)

	get := func() map[string]any {
		t.Helper()
		req := bearer(t, token, http.MethodGet, "/api/captcha/settings", "")
		rec := httptest.NewRecorder()
		env.mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET settings = %d: %s", rec.Code, rec.Body.String())
		}
		var body map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		return body
	}
	if got := get(); got["enabled"] != "false" {
		t.Fatalf("settings = %v, want disabled string form (F-003)", got)
	}
	// The schema select submits "true"/"false" strings (F-003).
	req := bearer(t, token, http.MethodPatch, "/api/captcha/settings", `{"enabled":"true"}`)
	rec := httptest.NewRecorder()
	env.mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("PATCH settings = %d: %s", rec.Code, rec.Body.String())
	}
	var patched map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &patched); err != nil {
		t.Fatalf("decode patched: %v", err)
	}
	if patched["enabled"] != "true" {
		t.Fatalf("patched = %v, want enabled string", patched)
	}
	if got := get(); got["enabled"] != "true" {
		t.Fatalf("settings after patch = %v, want enabled", got)
	}
	// JSON bool stays accepted for API clients.
	reqBool := bearer(t, token, http.MethodPatch, "/api/captcha/settings", `{"enabled":false}`)
	recBool := httptest.NewRecorder()
	env.mux.ServeHTTP(recBool, reqBool)
	if recBool.Code != http.StatusOK {
		t.Fatalf("PATCH bool = %d: %s", recBool.Code, recBool.Body.String())
	}
}

func TestCaptchaSettingsPermissionAndBody(t *testing.T) {
	env := newAuthTestEnv(t)
	anon := httptest.NewRecorder()
	env.mux.ServeHTTP(anon, httptest.NewRequest(http.MethodGet, "/api/captcha/settings", nil))
	if anon.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous settings = %d, want 401", anon.Code)
	}
	token := adminToken(t, env)
	for _, payload := range []string{"not json", `{}`, `{"enabled":"yes"}`} {
		req := bearer(t, token, http.MethodPatch, "/api/captcha/settings", payload)
		rec := httptest.NewRecorder()
		env.mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("PATCH %q = %d, want 400", payload, rec.Code)
		}
	}
}

func TestCaptchaLoginGateDefaultOff(t *testing.T) {
	env := newAuthTestEnv(t)
	code, body := loginBody(t, env, testSeedUsername, testSeedPassword)
	if code != http.StatusOK {
		t.Fatalf("login = %d: %v", code, body)
	}
}

func TestCaptchaLoginGateEnabledWrong(t *testing.T) {
	env := newAuthTestEnv(t)
	if err := env.captcha.SetEnabled(true, time.Now().UTC()); err != nil {
		t.Fatalf("enable: %v", err)
	}
	code, body := loginBody(t, env, testSeedUsername, testSeedPassword)
	if code != http.StatusBadRequest {
		t.Fatalf("login without captcha = %d: %v", code, body)
	}
	if body["error"] != "INVALID_CAPTCHA" {
		t.Fatalf("error = %v, want INVALID_CAPTCHA", body["error"])
	}
	payload := strings.NewReader("{\"username\":\"" + testSeedUsername + "\",\"password\":\"" + testSeedPassword + "\",\"captchaId\":\"cap-test-1\",\"captchaAnswer\":\"" + "99" + "\"}")
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", payload)
	rec := httptest.NewRecorder()
	env.mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("login with wrong captcha = %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCaptchaLoginGateEnabledCorrect(t *testing.T) {
	env := newAuthTestEnv(t)
	if err := env.captcha.SetEnabled(true, time.Now().UTC()); err != nil {
		t.Fatalf("enable: %v", err)
	}
	payload := strings.NewReader("{\"username\":\"" + testSeedUsername + "\",\"password\":\"" + testSeedPassword + "\",\"captchaId\":\"cap-test-1\",\"captchaAnswer\":\"" + "2" + "\"}")
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", payload)
	rec := httptest.NewRecorder()
	env.mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("login with correct captcha = %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCaptchaSettingsUpdateAudited(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	req := bearer(t, token, http.MethodPatch, "/api/captcha/settings", `{"enabled":"true"}`)
	rec := httptest.NewRecorder()
	env.mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("PATCH settings = %d: %s", rec.Code, rec.Body.String())
	}
	opsReq := bearer(t, token, http.MethodGet, "/api/operations?pageSize=100", "")
	opsRec := httptest.NewRecorder()
	env.mux.ServeHTTP(opsRec, opsReq)
	if opsRec.Code != http.StatusOK {
		t.Fatalf("operations = %d: %s", opsRec.Code, opsRec.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(opsRec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode operations: %v", err)
	}
	items, _ := body["items"].([]any)
	found := false
	for _, item := range items {
		row, _ := item.(map[string]any)
		if row["event"] == "captcha.settings-update" {
			found = true
		}
	}
	if !found {
		t.Fatalf("operations = %v, want captcha.settings-update row", body)
	}
}
