package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/requestid"
)

func TestJSONRouteErrors404And405(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"ok": "true"})
	})
	h := WithJSONRouteErrors(mux)

	t.Run("unknown path is JSON 404", func(t *testing.T) {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/does-not-exist", nil))
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", rr.Code)
		}
		if ct := rr.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
			t.Fatalf("content-type = %q, want application/json", ct)
		}
		var body map[string]any
		if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
			t.Fatalf("json: %v body=%s", err, rr.Body.String())
		}
		if body["error"] != "NOT_FOUND" {
			t.Fatalf("error = %v, want NOT_FOUND", body["error"])
		}
	})

	t.Run("wrong method is JSON 405", func(t *testing.T) {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/health", nil))
		if rr.Code != http.StatusMethodNotAllowed {
			t.Fatalf("status = %d, want 405 body=%s", rr.Code, rr.Body.String())
		}
		if ct := rr.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
			t.Fatalf("content-type = %q, want application/json", ct)
		}
		var body map[string]any
		if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
			t.Fatalf("json: %v body=%s", err, rr.Body.String())
		}
		if body["error"] != "METHOD_NOT_ALLOWED" {
			t.Fatalf("error = %v, want METHOD_NOT_ALLOWED", body["error"])
		}
	})

	t.Run("registered GET still 200", func(t *testing.T) {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/health", nil))
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rr.Code)
		}
	})

	t.Run("HEAD on GET route is not JSON 405", func(t *testing.T) {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodHead, "/api/health", nil))
		if rr.Code == http.StatusMethodNotAllowed {
			t.Fatalf("HEAD mapped to 405: %s", rr.Body.String())
		}
	})
}

func TestR1CorrelationIDAppearsInRouteErrorEnvelope(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"ok": "true"})
	})
	h := requestid.Middleware(WithJSONRouteErrors(mux))
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/missing", nil)
	req.Header.Set(requestid.HeaderName, "r1-route-001")
	h.ServeHTTP(rr, req)
	if rr.Header().Get(requestid.HeaderName) != "r1-route-001" {
		t.Fatalf("response request id = %q", rr.Header().Get(requestid.HeaderName))
	}
	var body map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("json: %v", err)
	}
	if body[requestid.BodyName] != "r1-route-001" {
		t.Fatalf("correlation_id = %v", body[requestid.BodyName])
	}
}
