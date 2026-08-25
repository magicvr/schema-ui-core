// R4 end-to-end HTTP evidence (workspace-019 · GOAL-005): the three IAM
// chains walked over the REAL handler mux with the REAL mock channel —
// recovery (start → outbox code → complete → login), invitation (admin
// issue → outbox link → accept → login with issued roles), and policy
// enforcement at the HTTP seam (tightened minLength blocks creation).
// Codes/links come from CHANNEL RECORDS, never from test stubs.
package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

var evidenceCodePattern = regexp.MustCompile(`\b\d{6}\b`)
var evidenceTokenPattern = regexp.MustCompile(`token=([A-Za-z0-9_-]+)`)

func testOutboxCurrent(t *testing.T, env *authTestEnv, to, mustContain string) (body string) {
	t.Helper()
	records, err := env.recoverySender.List(context.Background(), 50, 0)
	if err != nil {
		t.Fatalf("list outbox: %v", err)
	}
	for i := 0; i < len(records); i++ {
		rec := records[i]
		if rec.To != to {
			continue
		}
		full, err := env.recoverySender.Get(context.Background(), rec.ID)
		if err != nil {
			t.Fatalf("get outbox %s: %v", rec.ID, err)
		}
		if mustContain != "" && !strings.Contains(full.Body, mustContain) {
			continue
		}
		return full.Body
	}
	t.Fatalf("no outbox record addressed to %s matching %q (all: %+v)", to, mustContain, func() []string {
		var out []string
		for _, rec := range records {
			full, _ := env.recoverySender.Get(context.Background(), rec.ID)
			out = append(out, rec.To+" => "+full.Body[:min(len(full.Body), 40)])
		}
		return out
	}())
	return ""
}

func serveBearer(t *testing.T, mux *http.ServeMux, accessToken, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, bearer(t, accessToken, method, path, body))
	return rr
}

func TestR4RecoveryChainOverHTTP(t *testing.T) {
	env := newAuthTestEnv(t)
	token := env.login(t, testSeedUsername, testSeedPassword)
	if rec := serveBearer(t, env.mux, token, http.MethodPost, "/api/account/email/bind", `{"email":"r4admin@example.com"}`); rec.Code != http.StatusOK {
		t.Fatalf("bind status = %d %s", rec.Code, rec.Body.String())
	}
	code := evidenceCodePattern.FindString(testOutboxCurrent(t, env, "r4admin@example.com", "邮箱验证码"))
	if rec := serveBearer(t, env.mux, token, http.MethodPost, "/api/account/email/verify", `{"code":"`+code+`"}`); rec.Code != http.StatusOK {
		t.Fatalf("verify status = %d %s", rec.Code, rec.Body.String())
	}
	if code2, _ := sendJSON(t, env.mux, http.MethodPost, "/api/auth/recovery/start", `{"account":"admin"}`); code2 != http.StatusAccepted {
		t.Fatalf("start status = %d", code2)
	}
	rCode := evidenceCodePattern.FindString(testOutboxCurrent(t, env, "r4admin@example.com", "password reset code"))
	if code3, _ := sendJSON(t, env.mux, http.MethodPost, "/api/auth/recovery/complete",
		`{"account":"admin","code":"`+rCode+`","newPassword":"r4-recovered-pass-1"}`); code3 != http.StatusNoContent {
		t.Fatalf("complete status = %d", code3)
	}
	_, out := sendJSON(t, env.mux, http.MethodPost, "/api/auth/login", `{"username":"admin","password":"r4-recovered-pass-1"}`)
	if out["accessToken"] == nil {
		t.Fatal("login with recovered password failed — recovery chain not closed over HTTP")
	}
}

func TestR4InviteChainOverHTTP(t *testing.T) {
	env := newAuthTestEnv(t)
	admin := env.login(t, testSeedUsername, testSeedPassword)
	if rec := serveBearer(t, env.mux, admin, http.MethodPost, "/api/users/invites",
		`{"email":"r4inv@example.com","roles":["viewer"],"expiresInDays":7}`); rec.Code != http.StatusCreated {
		t.Fatalf("create invite status = %d %s", rec.Code, rec.Body.String())
	}
	m := evidenceTokenPattern.FindStringSubmatch(testOutboxCurrent(t, env, "r4inv@example.com", "账号邀请"))
	if len(m) != 2 {
		t.Fatal("invitation mail carries no token")
	}
	tok := m[1]
	if code, _ := sendJSON(t, env.mux, http.MethodPost, "/api/auth/invite/accept",
		`{"token":"`+tok+`","username":"r4invitee","name":"R4 Invitee","password":"r4-invited-pass-1"}`); code != http.StatusNoContent {
		t.Fatalf("accept status = %d", code)
	}
	_, body := sendJSON(t, env.mux, http.MethodPost, "/api/auth/login", `{"username":"r4invitee","password":"r4-invited-pass-1"}`)
	if body["accessToken"] == nil {
		t.Fatal("invitee login failed — invitation chain not closed over HTTP")
	}
	if roles := roleKeysOf(body); len(roles) != 1 || roles[0] != "viewer" {
		t.Fatalf("invitee roles = %v, want [viewer]", roles)
	}
	_, replay := sendJSON(t, env.mux, http.MethodPost, "/api/auth/invite/accept",
		`{"token":"`+tok+`","username":"r4again","password":"r4-invited-pass-1"}`)
	if !strings.Contains(stringify(replay), "INVITE_INVALID") {
		t.Fatalf("replay body = %v, want INVITE_INVALID", replay)
	}
}

func TestR4PolicyEnforcementOverHTTP(t *testing.T) {
	env := newAuthTestEnv(t)
	admin := env.login(t, testSeedUsername, testSeedPassword)
	if rec := serveBearer(t, env.mux, admin, http.MethodPatch, "/api/settings/password-policy",
		`{"minLength":12,"minCategories":0,"historyDepth":0}`); rec.Code != http.StatusOK {
		t.Fatalf("patch policy status = %d %s", rec.Code, rec.Body.String())
	}
	if rec := serveBearer(t, env.mux, admin, http.MethodPost, "/api/users",
		`{"username":"weakpw","name":"Weak","password":"eightbte"}`); !strings.Contains(rec.Body.String(), "INVALID_PASSWORD") {
		t.Fatalf("weak create body = %s, want INVALID_PASSWORD", rec.Body.String())
	}
	if rec := serveBearer(t, env.mux, admin, http.MethodPost, "/api/users",
		`{"username":"strongpw","name":"Strong","password":"strong-pass-123"}`); rec.Code != http.StatusOK && rec.Code != http.StatusCreated {
		t.Fatalf("strong create status = %d %s", rec.Code, rec.Body.String())
	}
}

func roleKeysOf(body map[string]any) []string {
	user, ok := body["user"].(map[string]any)
	if !ok {
		return nil
	}
	raw, ok := user["roles"].([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(raw))
	for _, r := range raw {
		if s, ok := r.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

func stringify(v map[string]any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}