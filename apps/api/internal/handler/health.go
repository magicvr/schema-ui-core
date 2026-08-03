package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
	"github.com/magicvr/schema-ui-core/apps/api/pkg/version"
)

type healthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version,omitempty"`
	Commit    string    `json:"commit,omitempty"`
}

// Register mounts the health endpoints, the R2 auth endpoints, the R4 account
// session route, the schema-driven resource CRUD API (SQLite-backed; users/roles
// are the GOAL-011 semantic resources), and the R1 schema document endpoint.
// Protected routes are wrapped in the request-identity middleware. The store is
// injected so the resource handlers read and write the same SQLite database that
// backs identity (GOAL-007 S3).
func Register(mux *http.ServeMux, a *auth.Authenticator, st *store.Store) {
	mux.Handle("GET /healthz", healthz())
	mux.Handle("GET /readyz", readyz(st))
	authsHandler(mux, a, st)
	accountsHandler(mux, a)
	registerResource(mux, a, usersResource(st))
	registerResource(mux, a, rolesResource(st))
	schemasHandler(mux)
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
// (A-002 F-002-006; compose uses this as service_healthy).
func readyz(st *store.Store) http.Handler {
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
