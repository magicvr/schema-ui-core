package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
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
// R4 C3.3: admin.users/admin.roles HTTP routes are mounted by the composition
// root from their kernel.Provider surfaces (MountProviderRoutes /
// RegisterContributions), not by this central Register. This function keeps
// core auth/accounts/health/schema registration only.
func Register(mux *http.ServeMux, a *auth.Authenticator, st *store.Store, plan kernel.Plan) {
	RegisterWithReadiness(mux, a, st, plan, nil)
}

// RegisterWithReadiness is Register plus an optional module-graph readiness
// probe (R5). ready, when non-nil, gates /readyz on Start+Ready success.
// Schema pages are registered separately via RegisterSchemas so composition can
// pass runtime contribution ownership (R5 C5.1).
func RegisterWithReadiness(mux *http.ServeMux, a *auth.Authenticator, st *store.Store, plan kernel.Plan, ready func() bool) {
	mux.Handle("GET /healthz", healthz())
	mux.Handle("GET /readyz", readyz(st, ready))
	if plan.HasModule("core.auth-session") {
		authsHandler(mux, a, st)
		accountsHandler(mux, a)
	}
}

// MountProviderRoutes mounts the admin.users/admin.roles HTTP routes generated
// by the generic resource factory for enabled modules. TEST-ONLY: the handler
// test environment uses it to mirror the production surface; the composition
// root MUST consume kernel.RegisterContributions (provider finalize) instead,
// never this function. Keeping it here avoids duplicating the route factory in
// the test package while production has a single mounting chain.
func MountProviderRoutes(mux *http.ServeMux, a *auth.Authenticator, st *store.Store, plan kernel.Plan) {
	if plan.HasModule("admin.users") {
		for _, route := range resourceRoutes(a, usersResource(st), "admin.users") {
			mux.Handle(route.Method+" "+route.Pattern, route.Handler)
		}
	}
	if plan.HasModule("admin.roles") {
		for _, route := range resourceRoutes(a, rolesResource(st), "admin.roles") {
			mux.Handle(route.Method+" "+route.Pattern, route.Handler)
		}
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
func readyz(st *store.Store, ready func() bool) http.Handler {
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
