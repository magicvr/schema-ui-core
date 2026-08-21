package obs

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"

	"net/http/httptest"
)

func newTestObserver() *Observer {
	return NewObserver(BuildInfo{Version: "t", Commit: "t", GoVersion: "t", Profile: "test"})
}

// TestServerDisabledIsNoOp pins the R1 default: disabled means nothing
// listens and lifecycle calls are safe no-ops.
func TestServerDisabledIsNoOp(t *testing.T) {
	s := NewServer(ServerOptions{Enabled: false, Addr: "127.0.0.1:25081"}, newTestObserver(), slog.Default())
	if s.Enabled() {
		t.Fatal("disabled server reports Enabled")
	}
	if err := s.Start(context.Background()); err != nil {
		t.Fatalf("disabled Start must be a no-op, got %v", err)
	}
	if s.Addr() != "" {
		t.Errorf("disabled Addr = %q, want empty", s.Addr())
	}
	if err := s.Stop(context.Background()); err != nil {
		t.Fatalf("disabled Stop must be a no-op, got %v", err)
	}
}

// TestServerEnabledServesMetrics starts a real loopback listener on an
// ephemeral port and scrapes it.
func TestServerEnabledServesMetrics(t *testing.T) {
	o := newTestObserver()
	o.RegisterModule("admin.users")
	s := NewServer(ServerOptions{Enabled: true, Addr: "127.0.0.1:0"}, o, slog.Default())
	if !s.Enabled() {
		t.Fatal("enabled server must report Enabled")
	}
	if err := s.Start(context.Background()); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer func() { _ = s.Stop(context.Background()) }()
	addr := s.Addr()
	if addr == "" {
		t.Fatal("Addr empty after Start")
	}

	resp, err := http.Get("http://" + addr + MetricsPath)
	if err != nil {
		t.Fatalf("scrape: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	if !strings.Contains(string(body), `suc_kernel_modules_enabled{module_id="admin.users"} 1`) {
		t.Errorf("module gauge missing in scrape:\n%s", body)
	}

	// Wrong path is not served on the dedicated listener.
	resp2, err := http.Get("http://" + addr + "/nope")
	if err != nil {
		t.Fatalf("GET /nope: %v", err)
	}
	resp2.Body.Close()
	if resp2.StatusCode != http.StatusNotFound {
		t.Errorf("GET /nope status = %d, want 404", resp2.StatusCode)
	}
}

// TestServerEnabledBindFailureFailsClosed pins the R2 D-001 §3 rule: an
// explicitly enabled listener that cannot bind fails startup.
func TestServerEnabledBindFailureFailsClosed(t *testing.T) {
	o := newTestObserver()
	blocker := NewServer(ServerOptions{Enabled: true, Addr: "127.0.0.1:0"}, o, slog.Default())
	if err := blocker.Start(context.Background()); err != nil {
		t.Fatalf("blocker Start: %v", err)
	}
	defer func() { _ = blocker.Stop(context.Background()) }()

	dup := NewServer(ServerOptions{Enabled: true, Addr: blocker.Addr()}, o, slog.Default())
	err := dup.Start(context.Background())
	if err == nil {
		t.Fatal("bind conflict must fail closed")
	}
	_ = dup.Stop(context.Background())
}

// TestServerNilObserverDisables keeps composition wiring safe when the
// observer could not be built.
func TestServerNilObserverDisables(t *testing.T) {
	s := NewServer(ServerOptions{Enabled: true, Addr: "127.0.0.1:0"}, nil, slog.Default())
	if s.Enabled() {
		t.Fatal("nil observer must disable the server")
	}
}

// TestInstrumentedMuxOwnership verifies Own + Handle/HandleFunc interception.
func TestInstrumentedMuxOwnership(t *testing.T) {
	o := newTestObserver()
	mux := NewInstrumentedMux(o)
	mux.Own("GET /api/users", "admin.users")
	mux.Handle("GET /api/users", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	mux.HandleFunc("POST /api/auth/login", func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	mux.ServeHTTP(httptest.NewRecorder(), req)
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	mux.ServeHTTP(httptest.NewRecorder(), req)

	body := scrape(t, o, "")
	for _, want := range []string{
		`suc_http_requests_total{method="GET",module_id="admin.users",route="/api/users",status="200"} 1`,
		`suc_http_requests_total{method="POST",module_id="core",route="/api/auth/login",status="200"} 1`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("missing %s in:\n%s", want, body)
		}
	}
}
