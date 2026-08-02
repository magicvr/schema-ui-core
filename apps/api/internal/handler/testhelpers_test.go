package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

// authTestEnv wires a full handler mux backed by a temp SQLite store and an
// Authenticator. Tests that only exercise public routes can use env.mux
// directly; write-route tests log in first to obtain a Bearer access token.
type authTestEnv struct {
	mux *http.ServeMux
	a   *auth.Authenticator
	st  *store.Store
}

const (
	testSeedUsername = "admin"
	testSeedPassword = "test-password"
	testBcryptCost   = 4 // cheap for tests
	testJWTSecret    = "test-secret"
)

// newAuthTestEnv seeds the admin user and mounts all routes (no dev-session
// fallback).
func newAuthTestEnv(t *testing.T) *authTestEnv {
	t.Helper()
	return newAuthTestEnvWith(t, false)
}

// newDevSessionTestEnv mounts all routes with the explicit dev-session fallback
// enabled (acceptance M9: production must not enable this).
func newDevSessionTestEnv(t *testing.T) *authTestEnv {
	t.Helper()
	return newAuthTestEnvWith(t, true)
}

func newAuthTestEnvWith(t *testing.T, devSession bool) *authTestEnv {
	t.Helper()
	hash, err := auth.HashPassword(testSeedPassword, testBcryptCost)
	if err != nil {
		t.Fatalf("hash seed password: %v", err)
	}
	st, err := store.Open(filepath.Join(t.TempDir(), "test.db"), testSeedUsername, hash, true)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	a := auth.New([]byte(testJWTSecret), 15*time.Minute, 30*24*time.Hour, st, devSession)
	mux := http.NewServeMux()
	Register(mux, a)
	return &authTestEnv{mux: mux, a: a, st: st}
}

// login performs POST /api/auth/login and returns the access token.
func (e *authTestEnv) login(t *testing.T, username, password string) string {
	t.Helper()
	body := `{"username":` + quote(username) + `,"password":` + quote(password) + `}`
	code, out := sendJSON(t, e.mux, http.MethodPost, "/api/auth/login", body)
	if code != http.StatusOK {
		t.Fatalf("login status = %d, want 200: %v", code, out)
	}
	tok, _ := out["accessToken"].(string)
	if tok == "" {
		t.Fatalf("accessToken missing in %v", out)
	}
	return tok
}

// addUser inserts a user directly into the store for permission tests.
func (e *authTestEnv) addUser(t *testing.T, username, password string, roles []string) {
	t.Helper()
	hash, err := auth.HashPassword(password, testBcryptCost)
	if err != nil {
		t.Fatalf("hash %s password: %v", username, err)
	}
	now := time.Now().UTC()
	if err := e.st.CreateUser(store.User{
		ID:           "user-" + username,
		Username:     username,
		Name:         username,
		Roles:        roles,
		PasswordHash: hash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}); err != nil {
		t.Fatalf("create user %s: %v", username, err)
	}
}

// bearer returns a request with the given access token attached.
func bearer(t *testing.T, token, method, path string, body string) *http.Request {
	t.Helper()
	var reader *strings.Reader
	if body == "" {
		reader = strings.NewReader("")
	} else {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Authorization", "Bearer "+token)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}

func quote(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func getJSON(t *testing.T, mux *http.ServeMux, path string) (int, map[string]any) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	var body map[string]any
	if rr.Body.Len() > 0 {
		if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
			t.Fatalf("decode %q: %v", rr.Body.String(), err)
		}
	}
	return rr.Code, body
}

// getRecords fetches a records path as the seeded admin (GOAL-006 S4: reads are
// authenticated and permission-gated, so a Bearer token is required).
func getRecords(t *testing.T, env *authTestEnv, path string) (int, map[string]any) {
	t.Helper()
	req := bearer(t, adminToken(t, env), http.MethodGet, path, "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	var body map[string]any
	if rr.Body.Len() > 0 {
		if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
			t.Fatalf("decode %q: %v", rr.Body.String(), err)
		}
	}
	return rr.Code, body
}

func sendJSON(t *testing.T, mux *http.ServeMux, method, path, body string) (int, map[string]any) {
	t.Helper()
	var reader *strings.Reader
	if body == "" {
		reader = strings.NewReader("")
	} else {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	var out map[string]any
	if rr.Body.Len() > 0 && rr.Header().Get("Content-Type") != "" {
		_ = json.NewDecoder(rr.Body).Decode(&out)
	}
	return rr.Code, out
}

func containsString(haystack, needle string) bool {
	return strings.Contains(haystack, needle)
}
