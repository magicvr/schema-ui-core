package composition

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/ratelimit"
	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
	authsession "github.com/magicvr/schema-ui-core/apps/api/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
	settingsrepository "github.com/magicvr/schema-ui-core/apps/api/modules/settings/repository"
)

func operationalGateServer(t *testing.T, mode config.RuntimeMode) *http.Server {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "r5-gate.db")
	st, err := testsupport.OpenStore(dbPath, "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	cfg := &config.Config{ProfileName: "admin", DBPath: dbPath, RuntimeMode: mode}
	plan, err := ResolvePlan(cfg)
	if err != nil {
		t.Fatal(err)
	}
	a := auth.New([]byte("test-secret"), 0, 0, st, false)
	jobs, err := newJobRuntime(st)
	if err != nil {
		t.Fatal(err)
	}
	cachePort, err := newCache(cfg)
	if err != nil {
		t.Fatal(err)
	}
	mux, err := newMux(
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
		nil,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		cachePort,
		ratelimit.NewProvider(),
	)
	if err != nil {
		t.Fatal(err)
	}
	return newServer(cfg, mux, slog.New(slog.NewTextHandler(io.Discard, nil)))
}

func TestOperationalGateCompositionCoversCoreAndProviderRoutes(t *testing.T) {
	for _, tc := range []struct {
		mode config.RuntimeMode
		code string
	}{
		{mode: config.RuntimeModeMaintenance, code: "SERVICE_MAINTENANCE"},
		{mode: config.RuntimeModeDegraded, code: "SERVICE_DEGRADED"},
		{mode: config.RuntimeModeReadOnly, code: "SERVICE_READ_ONLY"},
	} {
		t.Run(string(tc.mode), func(t *testing.T) {
			srv := operationalGateServer(t, tc.mode)

			// Provider contribution: dictionary create is denied before auth.
			provider := httptest.NewRecorder()
			providerReq := httptest.NewRequest(http.MethodPost, "/api/data-dictionary/types", nil)
			providerReq.Header.Set("X-Request-ID", "r5-provider")
			srv.Handler.ServeHTTP(provider, providerReq)
			assertOperationalError(t, provider, tc.code, "r5-provider")

			// R6 central management mutation is covered by the same operational
			// boundary before auth or credential creation runs.
			credential := httptest.NewRecorder()
			credentialReq := httptest.NewRequest(http.MethodPost, "/api/service-credentials", nil)
			credentialReq.Header.Set("X-Request-ID", "r6-service-credential")
			srv.Handler.ServeHTTP(credential, credentialReq)
			assertOperationalError(t, credential, tc.code, "r6-service-credential")

			// Core route: login remains available for session recovery and keeps its
			// own validation error rather than becoming an operational denial.
			login := httptest.NewRecorder()
			srv.Handler.ServeHTTP(login, httptest.NewRequest(http.MethodPost, "/api/auth/login", nil))
			if login.Code != http.StatusBadRequest {
				t.Fatalf("login status = %d, want 400; body=%s", login.Code, login.Body.String())
			}

			health := httptest.NewRecorder()
			srv.Handler.ServeHTTP(health, httptest.NewRequest(http.MethodGet, "/healthz", nil))
			if health.Code != http.StatusOK {
				t.Fatalf("health status = %d, want 200", health.Code)
			}
		})
	}
}

func TestOperationalGateCompositionPreservesNormalAuthAndEnvelopeRouting(t *testing.T) {
	srv := operationalGateServer(t, config.RuntimeModeNormal)

	provider := httptest.NewRecorder()
	srv.Handler.ServeHTTP(provider, httptest.NewRequest(http.MethodPost, "/api/data-dictionary/types", nil))
	if provider.Code != http.StatusUnauthorized {
		t.Fatalf("normal provider write status = %d, want 401; body=%s", provider.Code, provider.Body.String())
	}

	unknown := httptest.NewRecorder()
	srv.Handler.ServeHTTP(unknown, httptest.NewRequest(http.MethodPost, "/api/unknown", nil))
	if unknown.Code != http.StatusNotFound {
		t.Fatalf("unknown POST status = %d, want 404; body=%s", unknown.Code, unknown.Body.String())
	}

	mismatch := httptest.NewRecorder()
	srv.Handler.ServeHTTP(mismatch, httptest.NewRequest(http.MethodPost, "/healthz", nil))
	if mismatch.Code != http.StatusMethodNotAllowed {
		t.Fatalf("health POST status = %d, want 405; body=%s", mismatch.Code, mismatch.Body.String())
	}
}

func assertOperationalError(t *testing.T, response *httptest.ResponseRecorder, code, requestID string) {
	t.Helper()
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503; body=%s", response.Code, response.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["error"] != code || body["correlation_id"] != requestID {
		t.Fatalf("operational body = %v, want code=%s correlation_id=%s", body, code, requestID)
	}
}
