package handler

import (
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

// Register mounts the health endpoint, the R2 auth endpoints, the R4 account
// session route, the R4 records CRUD API (SQLite-backed), and the R1 schema
// document endpoint. Protected routes are wrapped in the request-identity
// middleware. The store is injected so the records handler reads and writes the
// same SQLite database that backs identity (GOAL-007 S3).
func Register(mux *http.ServeMux, a *auth.Authenticator, st *store.Store) {
	mux.Handle("GET /healthz", healthz())
	authsHandler(mux, a, st)
	accountsHandler(mux, a)
	recordsHandler(mux, a, st)
	schemasHandler(mux)
}

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

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]string{"error": code, "message": message})
}
