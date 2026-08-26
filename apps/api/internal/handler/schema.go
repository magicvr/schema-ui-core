package handler

import (
	"net/http"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// schemaHandler serves page schema documents by manifest pageId. Documents are
// read-only, so the map is built once and shared across requests.
type schemaHandler struct {
	documents map[string][]byte // pageId -> raw JSON document
}

// RegisterSchemas publishes only finalized PageContribution documents. Profile
// filtering, owner validation, duplicate checks and JSON identity validation
// have already completed in kernel.RegisterContributions.
// W13 F-010 (GOAL-013 A-001): the surface rides the standard auth middleware —
// page documents describe authenticated admin surfaces, and an anonymous
// visitor must not be able to enumerate them.
func RegisterSchemas(mux routeRegistrar, a *auth.Authenticator, pages []kernel.PageContribution) {
	documents := make(map[string][]byte, len(pages))
	for _, page := range pages {
		documents[page.PageID] = append([]byte(nil), page.Document...)
	}
	h := &schemaHandler{documents: documents}
	mux.Handle("GET /api/schema/{pageId}", a.Middleware(h.schema()))
}

// schema serves GET /api/schema/{pageId}: the raw page document, or 404 when
// the pageId has no seeded document.
func (h *schemaHandler) schema() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageID := r.PathValue("pageId")
		raw, ok := h.documents[pageID]
		if !ok {
			writeLocalizedError(w, r, http.StatusNotFound, "SCHEMA_NOT_FOUND", "no page document for that pageId")
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(raw)
	})
}
