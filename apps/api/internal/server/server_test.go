package server

import (
	"log/slog"
	"net/http"
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
