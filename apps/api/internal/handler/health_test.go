package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	mux := newAuthTestEnv(t).mux

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var body healthResponse
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Status != "ok" {
		t.Fatalf("status = %q, want ok", body.Status)
	}
}

// TestReadyz covers A-002 F-002-006: readiness reports ok while SQLite answers
// and 503 (unavailable) after the store is closed (fault injection).
func TestReadyz(t *testing.T) {
	t.Run("reports ok when SQLite answers", func(t *testing.T) {
		mux := newAuthTestEnv(t).mux

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/readyz", nil))
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}
		var body healthResponse
		if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body.Status != "ok" {
			t.Fatalf("status = %q, want ok", body.Status)
		}
	})

	t.Run("reports unavailable when SQLite is dead", func(t *testing.T) {
		env := newAuthTestEnv(t)
		// Fault injection: close the underlying store, then probe readiness.
		if err := env.st.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/readyz", nil))
		if rr.Code != http.StatusServiceUnavailable {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusServiceUnavailable)
		}
		var body healthResponse
		if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body.Status != "unavailable" {
			t.Fatalf("status = %q, want unavailable", body.Status)
		}
	})
}
