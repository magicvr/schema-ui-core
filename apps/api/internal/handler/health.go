package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
	"github.com/magicvr/schema-ui-core/apps/api/pkg/version"
)

type healthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version,omitempty"`
	Commit    string    `json:"commit,omitempty"`
}

// Register mounts core routes and the selected module contributions. The
// composition root passes the already-resolved plan so HTTP and Schema
// surfaces cannot silently diverge from the published Manifest.
//
// R4 C3.3 / R6 C6.1: admin.users/admin.roles HTTP routes are mounted by the
// composition root from their kernel.Provider surfaces (RegisterContributions),
// not by this central Register. This function keeps core auth/accounts/health
// registration only.
func Register(mux *http.ServeMux, a *auth.Authenticator, st kernel.Store, operations operationlog.Recorder, plan kernel.Plan, limiters kernel.RateLimiterProvider) {
	RegisterWithReadiness(mux, a, st, operations, plan, nil, limiters)
}

// RegisterWithReadiness is Register plus an optional module-graph readiness
// probe (R5). ready, when non-nil, gates /readyz on Start+Ready success.
// Schema pages are registered separately via RegisterSchemas so composition can
// pass runtime contribution ownership (R5 C5.1).
func RegisterWithReadiness(mux *http.ServeMux, a *auth.Authenticator, st kernel.Store, operations operationlog.Recorder, plan kernel.Plan, ready func() bool, limiters kernel.RateLimiterProvider, captcha ...CaptchaVerifier) {
	RegisterWithMFA(mux, a, st, operations, plan, ready, limiters, captcha, nil)
}

// RegisterWithMFA is RegisterWithReadiness plus the optional second-factor
// login gate (S-10 · GOAL-017 D-002 §3): nil keeps the login contract
// byte-identical.
func RegisterWithMFA(mux *http.ServeMux, a *auth.Authenticator, st kernel.Store, operations operationlog.Recorder, plan kernel.Plan, ready func() bool, limiters kernel.RateLimiterProvider, captcha []CaptchaVerifier, mfa MFAVerifier) {
	RegisterWithMFAProbes(mux, a, st, operations, plan, ready, limiters, captcha, mfa)
}

// RegisterWithMFAProbes is RegisterWithMFA plus optional readiness probes
// beyond the store ping (VP-014 GOAL-003 D-001): when an S3-compatible object
// backend is explicitly configured, composition passes a HeadBucket probe so
// readyz covers the backend too. Nil entries are ignored.
func RegisterWithMFAProbes(mux routeRegistrar, a *auth.Authenticator, st kernel.Store, operations operationlog.Recorder, plan kernel.Plan, ready func() bool, limiters kernel.RateLimiterProvider, captcha []CaptchaVerifier, mfa MFAVerifier, probes ...func(context.Context) error) {
	mux.Handle("GET /healthz", healthz())
	mux.Handle("GET /readyz", readyz(st, ready, probes...))
	if plan.HasModule("core.auth-session") {
		var verifier CaptchaVerifier
		if len(captcha) > 0 {
			verifier = captcha[0]
		}
		authsHandler(mux, a, operations, limiters, verifier, mfa)
		accountsHandler(mux, a)
	}
}

// healthz is the liveness probe: the process is up and serving. It never
// touches the database (A-002 F-002-006 separates liveness from readiness).
func healthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, healthResponse{
			Status:    "ok",
			Timestamp: time.Now().UTC(),
			Version:   version.Version,
			Commit:    version.Commit,
		})
	})
}

// readyz is the readiness probe: liveness plus a trivial SQLite read, so a
// dead, read-only or unmigrated database flips the container health gate
// (A-002 F-002-006; compose uses this as service_healthy). R5: when ready is
// non-nil it must also report true (module graph Start+Ready succeeded),
// otherwise the probe stays unavailable (freeze §3 — readyz is real module-graph
// readiness, not just store ping).
func readyz(st kernel.Store, ready func() bool, extra ...func(context.Context) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second)
		defer cancel()
		if err := st.Ping(ctx); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, healthResponse{
				Status:    "unavailable",
				Timestamp: time.Now().UTC(),
				Version:   version.Version,
				Commit:    version.Commit,
			})
			return
		}
		if ready != nil && !ready() {
			writeJSON(w, http.StatusServiceUnavailable, healthResponse{
				Status:    "not-ready",
				Timestamp: time.Now().UTC(),
				Version:   version.Version,
				Commit:    version.Commit,
			})
			return
		}
		// VP-014 GOAL-003: explicit object-backend probes share the readyz
		// deadline; any failure keeps the whole probe unavailable.
		for _, probe := range extra {
			if probe == nil {
				continue
			}
			if err := probe(ctx); err != nil {
				writeJSON(w, http.StatusServiceUnavailable, healthResponse{
					Status:    "unavailable",
					Timestamp: time.Now().UTC(),
					Version:   version.Version,
					Commit:    version.Commit,
				})
				return
			}
		}
		writeJSON(w, http.StatusOK, healthResponse{
			Status:    "ok",
			Timestamp: time.Now().UTC(),
			Version:   version.Version,
			Commit:    version.Commit,
		})
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]string{"error": code, "message": message})
}
