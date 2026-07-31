package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
)

func TestAccountsMe(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/accounts/me", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var body map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	user, ok := body["user"].(map[string]any)
	if !ok {
		t.Fatalf("user = %v, want object", body["user"])
	}
	if user["id"] != "dev-001" {
		t.Fatalf("user.id = %v, want dev-001", user["id"])
	}
	if _, ok := body["features"]; !ok {
		t.Fatalf("features missing from %v", body)
	}
}

func TestAccountsMeNoSessionFailsClosed(t *testing.T) {
	mux := http.NewServeMux()
	h := &accountHandler{
		sessionProvider: func() (account.Session, bool) {
			return account.Session{}, false
		},
	}
	mux.Handle("GET /api/accounts/me", h.me())

	req := httptest.NewRequest(http.MethodGet, "/api/accounts/me", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}
	if got := rr.Body.String(); got != "" && !containsString(got, "UNAUTHENTICATED") {
		t.Fatalf("body = %q, want UNAUTHENTICATED error", got)
	}
}

func TestAccountsMeNilProviderFailsClosed(t *testing.T) {
	mux := http.NewServeMux()
	// Nil sessionProvider must not panic and must fail closed.
	h := &accountHandler{}
	mux.Handle("GET /api/accounts/me", h.me())

	req := httptest.NewRequest(http.MethodGet, "/api/accounts/me", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}
	if got := rr.Body.String(); got != "" && !containsString(got, "UNAUTHENTICATED") {
		t.Fatalf("body = %q, want UNAUTHENTICATED error", got)
	}
}

func containsString(haystack, needle string) bool {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
