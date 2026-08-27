package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
	"github.com/magicvr/schema-ui-core/apps/api/internal/requestid"
)

func TestOperationalGateModesAndAllowlist(t *testing.T) {
	for _, tc := range []struct {
		mode config.RuntimeMode
		want int
		code string
	}{
		{mode: config.RuntimeModeNormal, want: http.StatusCreated},
		{mode: config.RuntimeModeMaintenance, want: http.StatusServiceUnavailable, code: "SERVICE_MAINTENANCE"},
		{mode: config.RuntimeModeDegraded, want: http.StatusServiceUnavailable, code: "SERVICE_DEGRADED"},
		{mode: config.RuntimeModeReadOnly, want: http.StatusServiceUnavailable, code: "SERVICE_READ_ONLY"},
	} {
		t.Run(string(tc.mode), func(t *testing.T) {
			mux := http.NewServeMux()
			mux.Handle("GET /api/resource", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
			mux.Handle("POST /api/resource", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusCreated) }))
			mux.Handle("POST /api/auth/login", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) }))
			cfg := &config.Config{RuntimeMode: tc.mode}
			h := requestid.Middleware(WithOperationalGate(cfg, mux, WithJSONRouteErrors(mux)))

			get := httptest.NewRecorder()
			h.ServeHTTP(get, httptest.NewRequest(http.MethodGet, "/api/resource", nil))
			if get.Code != http.StatusOK {
				t.Fatalf("GET status = %d, want 200", get.Code)
			}

			post := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/resource", nil)
			req.Header.Set("X-Request-ID", "r5-gate-test")
			h.ServeHTTP(post, req)
			if post.Code != tc.want {
				t.Fatalf("POST status = %d, want %d; body=%s", post.Code, tc.want, post.Body.String())
			}
			if tc.code != "" {
				var body map[string]any
				if err := json.Unmarshal(post.Body.Bytes(), &body); err != nil {
					t.Fatal(err)
				}
				if body["error"] != tc.code || body[requestid.BodyName] != "r5-gate-test" {
					t.Fatalf("gate body = %v", body)
				}
			}

			login := httptest.NewRecorder()
			h.ServeHTTP(login, httptest.NewRequest(http.MethodPost, "/api/auth/login", nil))
			if login.Code != http.StatusNoContent {
				t.Fatalf("allowlisted login status = %d, want 204", login.Code)
			}
		})
	}
}

func TestOperationalGatePreservesUnknownAndMethodMismatch(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("GET /api/resource", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	cfg := &config.Config{RuntimeMode: config.RuntimeModeMaintenance}
	h := WithOperationalGate(cfg, mux, WithJSONRouteErrors(mux))

	unknown := httptest.NewRecorder()
	h.ServeHTTP(unknown, httptest.NewRequest(http.MethodPost, "/api/unknown", nil))
	if unknown.Code != http.StatusNotFound {
		t.Fatalf("unknown POST status = %d, want 404", unknown.Code)
	}

	mismatch := httptest.NewRecorder()
	h.ServeHTTP(mismatch, httptest.NewRequest(http.MethodPost, "/api/resource", nil))
	if mismatch.Code != http.StatusMethodNotAllowed {
		t.Fatalf("method mismatch status = %d, want 405", mismatch.Code)
	}
}

func TestOperationalGateCoversAllMutationMethods(t *testing.T) {
	methods := []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}
	mux := http.NewServeMux()
	for _, method := range methods {
		mux.Handle(method+" /api/mutation", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))
	}
	h := WithOperationalGate(&config.Config{RuntimeMode: config.RuntimeModeReadOnly}, mux, WithJSONRouteErrors(mux))
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			h.ServeHTTP(recorder, httptest.NewRequest(method, "/api/mutation", nil))
			if recorder.Code != http.StatusServiceUnavailable {
				t.Fatalf("%s status = %d, want 503", method, recorder.Code)
			}
		})
	}
}

func TestOperationalGateAllowsEveryRecoveryPath(t *testing.T) {
	paths := []string{"/api/auth/login", "/api/auth/refresh", "/api/auth/logout", "/api/auth/mfa/verify", "/api/account/password", "/api/mfa/enroll", "/api/mfa/confirm", "/api/mfa/disable", "/api/mfa/recovery/rotate",
		// workspace-019 R2: self-recovery stays reachable in maintenance mode.
		"/api/auth/recovery/start", "/api/auth/recovery/complete",
		"/api/auth/invite/accept"}
	mux := http.NewServeMux()
	for _, path := range paths {
		mux.Handle("POST "+path, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))
	}
	h := WithOperationalGate(&config.Config{RuntimeMode: config.RuntimeModeMaintenance}, mux, WithJSONRouteErrors(mux))
	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			h.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, path, nil))
			if recorder.Code != http.StatusNoContent {
				t.Fatalf("%s status = %d, want 204", path, recorder.Code)
			}
		})
	}
}
