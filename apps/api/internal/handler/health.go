package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/pkg/version"
)

type healthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version,omitempty"`
	Commit    string    `json:"commit,omitempty"`
}

// Register mounts R1 routes and the R4 account session route.
func Register(mux *http.ServeMux) {
	mux.Handle("GET /healthz", healthz())
	accountsHandler(mux)
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
