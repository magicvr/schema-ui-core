package handler

import (
	"embed"
	"net/http"
	"strings"
)

// Page schema documents backing GET /api/schema/{pageId}. The app manifest
// declares each page's schemaUrl as /api/schema/<pageId>; these fixtures make
// that contract real so the Web loader has a runtime source. GOAL-002 seeds a
// minimal subset (overview/catalog); the representative page set is authored
// under GOAL-004.
//
//go:embed fixtures/schema/*.json
var schemaPageFixtures embed.FS

// schemaHandler serves page schema documents by manifest pageId. Documents are
// read-only, so the map is built once and shared across requests.
type schemaHandler struct {
	documents map[string][]byte // pageId -> raw JSON document
}

func schemasHandler(mux *http.ServeMux) {
	h := &schemaHandler{documents: staticSchemaDocuments()}
	mux.Handle("GET /api/schema/{pageId}", h.schema())
}

// staticSchemaDocuments loads the embedded page fixtures into a pageId -> raw
// JSON map. The embed set is a build-time invariant; a missing file is a
// programming error surfaced at startup.
func staticSchemaDocuments() map[string][]byte {
	entries, err := schemaPageFixtures.ReadDir("fixtures/schema")
	if err != nil {
		panic(err)
	}
	documents := make(map[string][]byte, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		pageID := strings.TrimSuffix(entry.Name(), ".json")
		if pageID == "" || pageID == entry.Name() {
			continue
		}
		raw, err := schemaPageFixtures.ReadFile("fixtures/schema/" + entry.Name())
		if err != nil {
			panic(err)
		}
		documents[pageID] = raw
	}
	return documents
}

// schema serves GET /api/schema/{pageId}: the raw page document, or 404 when
// the pageId has no seeded document.
func (h *schemaHandler) schema() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageID := r.PathValue("pageId")
		raw, ok := h.documents[pageID]
		if !ok {
			writeError(w, http.StatusNotFound, "SCHEMA_NOT_FOUND", "no page document for that pageId")
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(raw)
	})
}
