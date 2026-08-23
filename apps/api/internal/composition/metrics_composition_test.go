package composition

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	settingsrepository "github.com/magicvr/schema-ui-core/apps/api/internal/modules/settings/repository"
	"github.com/magicvr/schema-ui-core/apps/api/internal/obs"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

// metricsDrillMux assembles the composition mux with a live observer and the
// test-only probe module enabled (same host pattern as the S2 drill).
func metricsDrillMux(t *testing.T, cfg *config.Config) (*http.ServeMux, *obs.Observer) {
	t.Helper()
	all := append([]kernel.Module(nil), kernel.BuiltinModules()...)
	all = append(all, probeModule())
	registry, err := kernel.NewRegistry(all)
	if err != nil {
		t.Fatalf("registry with probe: %v", err)
	}
	resolution, err := kernel.ResolveProfile("admin", nil)
	if err != nil {
		t.Fatalf("resolve admin: %v", err)
	}
	plan, err := registry.Resolve(append(resolution.Modules, probeModule().ID))
	if err != nil {
		t.Fatalf("plan with probe: %v", err)
	}
	st, err := testsupport.OpenStore(":memory:", "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	a := auth.New([]byte("test-secret"), 0, 0, st, false)
	jobs, err := newJobRuntime(st)
	if err != nil {
		t.Fatal(err)
	}
	observer := newObserver(cfg, plan, obs.NewTracing(obs.TracingOptions{Enabled: false}, slog.Default()))
	mux, err := newMuxWithExtraProviders(
		cfg,
		a,
		st,
		authsession.NewRepository(st),
		operationlog.NewRepository(st),
		settingsrepository.New(st),
		plan,
		&readinessGate{},
		jwtSecret("test-secret"),
		jobs,
		[]kernel.Provider{&probeProvider{desc: probeModule()}},
		observer,
		slog.Default(),
	)
	if err != nil {
		t.Fatalf("instrumented composition: %v", err)
	}
	return mux, observer
}

func scrapeObserver(t *testing.T, o *obs.Observer) string {
	t.Helper()
	srv := httptest.NewServer(o.Handler(""))
	t.Cleanup(srv.Close)
	resp, err := http.Get(srv.URL + "/metrics")
	if err != nil {
		t.Fatalf("scrape: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

// TestMetricsCompositionTagsModuleOwnership pins the R2 acceptance core:
// contributed routes are measured under their owning module_id, central
// registrations under core, and the build/module series are present.
func TestMetricsCompositionTagsModuleOwnership(t *testing.T) {
	cfg := &config.Config{ProfileName: "admin", DBPath: "unused.db"}
	mux, observer := metricsDrillMux(t, cfg)

	// Contributed probe route (200) and central health route (200).
	mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/api/probe-items", nil))
	mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/healthz", nil))

	body := scrapeObserver(t, observer)
	for _, want := range []string{
		`suc_http_requests_total{method="GET",module_id="admin.probe",route="/api/probe-items",status="200"} 1`,
		`suc_http_requests_total{method="GET",module_id="core",route="/healthz",status="200"} 1`,
		`suc_kernel_modules_enabled{module_id="admin.probe"} 1`,
		`suc_build_info{`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("missing %q in exposition:\n%s", want, body)
		}
	}
}

// TestMetricsServerWiringEnabledScrapes proves the composition mapping
// (config -> obs.Server) serves the exposition face end to end.
func TestMetricsServerWiringEnabledScrapes(t *testing.T) {
	cfg := &config.Config{ProfileName: "admin", MetricsEnabled: true, MetricsAddr: "127.0.0.1:0"}
	observer := newObserver(cfg, kernel.Plan{}, obs.NewTracing(obs.TracingOptions{Enabled: false}, slog.Default()))
	srv := newMetricsServer(cfg, observer, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if !srv.Enabled() {
		t.Fatal("enabled config must produce an enabled server")
	}
	if err := srv.Start(context.Background()); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer func() { _ = srv.Stop(context.Background()) }()

	resp, err := http.Get("http://" + srv.Addr() + "/metrics")
	if err != nil {
		t.Fatalf("scrape: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), `suc_build_info{`) {
		t.Errorf("suc_build_info missing:\n%s", body)
	}
}

// TestTracingWiringNoopByDefault pins GOAL-004 D-001 §3 at the wiring layer:
// the default config produces an inert tracer path.
func TestTracingWiringNoopByDefault(t *testing.T) {
	cfg := &config.Config{ProfileName: "admin"}
	tr := newTracing(cfg)
	if tr.Enabled() {
		t.Fatal("default config must keep tracing a no-op")
	}
	if err := tr.Shutdown(context.Background()); err != nil {
		t.Fatalf("no-op Shutdown must return nil, got %v", err)
	}
}

// TestTracingWiringEnabledBuildsProvider proves the config mapping flips the
// tracer path on with an explicit endpoint.
func TestTracingWiringEnabledBuildsProvider(t *testing.T) {
	cfg := &config.Config{
		ProfileName:       "admin",
		TracesEnabled:     true,
		TracesEndpoint:    "http://127.0.0.1:4318",
		TracesSampleRatio: 1.0,
		AppName:           "schema-ui-core-api",
	}
	tr := newTracing(cfg)
	if !tr.Enabled() {
		t.Fatal("enabled config must produce an enabled tracer path")
	}
	// Shutdown without any spans must still succeed (provider exists).
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := tr.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
}

// TestMetricsServerWiringDisabledIsInert pins the default: disabled config
// yields an inert server even with a live observer.
func TestMetricsServerWiringDisabledIsInert(t *testing.T) {
	cfg := &config.Config{ProfileName: "admin"}
	observer := newObserver(cfg, kernel.Plan{}, obs.NewTracing(obs.TracingOptions{Enabled: false}, slog.Default()))
	srv := newMetricsServer(cfg, observer, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if srv.Enabled() {
		t.Fatal("default config must keep metrics disabled")
	}
	if err := srv.Start(context.Background()); err != nil {
		t.Fatalf("disabled Start must be no-op: %v", err)
	}
	if srv.Addr() != "" {
		t.Errorf("disabled server bound %q", srv.Addr())
	}
}
