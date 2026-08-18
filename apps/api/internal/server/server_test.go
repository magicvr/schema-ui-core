package server

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
)

// GOAL-003 · ReadHeaderTimeout must always be set (defaults to ReadTimeout).
func TestNewSetsReadHeaderTimeout(t *testing.T) {
	cfg := &config.Config{
		HTTPAddr:     ":0",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	srv := New(cfg, http.NotFoundHandler(), logger)
	if srv.ReadHeaderTimeout != 5*time.Second {
		t.Fatalf("ReadHeaderTimeout = %v, want 5s", srv.ReadHeaderTimeout)
	}
	if srv.ReadTimeout != 5*time.Second || srv.WriteTimeout != 10*time.Second {
		t.Fatalf("timeouts not preserved: read=%v write=%v", srv.ReadTimeout, srv.WriteTimeout)
	}
}

func TestWrapSecurityCORSAndNosniff(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	cfg := &config.Config{HTTPCORSOrigins: []string{"https://app.example"}}
	h := WrapSecurity(cfg, inner)

	t.Run("allowed origin preflight", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/api/users", nil)
		req.Header.Set("Origin", "https://app.example")
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)
		if rr.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want 204", rr.Code)
		}
		if rr.Header().Get("Access-Control-Allow-Origin") != "https://app.example" {
			t.Fatalf("allow-origin = %q", rr.Header().Get("Access-Control-Allow-Origin"))
		}
		if rr.Header().Get("X-Content-Type-Options") != "nosniff" {
			t.Fatal("missing nosniff")
		}
	})

	t.Run("unknown origin has nosniff but no CORS", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
		req.Header.Set("Origin", "https://evil.example")
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rr.Code)
		}
		if rr.Header().Get("Access-Control-Allow-Origin") != "" {
			t.Fatalf("unexpected CORS for foreign origin")
		}
		if rr.Header().Get("X-Content-Type-Options") != "nosniff" {
			t.Fatal("missing nosniff")
		}
	})
}

func TestNewEmitsCorrelationIDForDownstreamErrors(t *testing.T) {
	cfg := &config.Config{HTTPAddr: ":0"}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	srv := New(cfg, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}), logger)
	// The downstream body is intentionally irrelevant; the boundary guarantee is
	// the stable response header and request propagation, exercised by the
	// requestid package. Keep this test focused on the server composition.
	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set("X-Request-ID", "server-test-1")
	rr := httptest.NewRecorder()
	srv.Handler.ServeHTTP(rr, req)
	if got := rr.Header().Get("X-Request-ID"); got != "server-test-1" {
		t.Fatalf("response X-Request-ID = %q", got)
	}
}
