package handler

import (
	"embed"
	"net/http"
	"strings"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	activityschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/activity/schema"
	rolesschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/roles/schema"
	settingsschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/settings/schema"
	usersschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/users/schema"
)

// Core page schema documents backing GET /api/schema/{pageId}. Module-owned
// Settings and Activity documents are merged from their package-local embed
// sources below; the app manifest declares the same schemaUrl contract.
//
//go:embed fixtures/schema/*.json
var schemaPageFixtures embed.FS

// schemaHandler serves page schema documents by manifest pageId. Documents are
// read-only, so the map is built once and shared across requests.
type schemaHandler struct {
	documents map[string][]byte // pageId -> raw JSON document
}

func schemasHandler(mux *http.ServeMux, plan kernel.Plan) {
	h := &schemaHandler{documents: schemaDocumentsForPlan(plan)}
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
	for pageID, raw := range settingsschema.SchemaDocuments() {
		documents[pageID] = raw
	}
	for pageID, raw := range activityschema.SchemaDocuments() {
		documents[pageID] = raw
	}
	// R4 C3.3: users/roles schema documents are module-owned (content migrated
	// out of the central fixture embed).
	for pageID, raw := range usersschema.SchemaDocuments() {
		documents[pageID] = raw
	}
	for pageID, raw := range rolesschema.SchemaDocuments() {
		documents[pageID] = raw
	}
	return documents
}

func schemaDocumentsForPlan(plan kernel.Plan) map[string][]byte {
	documents := staticSchemaDocuments()
	// R4 C4.3: page→module ownership is derived from the module schema
	// contributors (no hardcoded owner map); plan gating stays by module.
	owners := map[string]string{}
	for _, contributor := range []struct {
		moduleID string
		pageIDs  []string
	}{
		{usersschema.ModuleID, usersschema.PageIDs()},
		{rolesschema.ModuleID, rolesschema.PageIDs()},
		{settingsschema.ModuleID, settingsschema.PageIDs()},
		{activityschema.ModuleID, activityschema.PageIDs()},
	} {
		for _, pageID := range contributor.pageIDs {
			owners[pageID] = contributor.moduleID
		}
	}
	for pageID, moduleID := range owners {
		if !plan.HasModule(moduleID) {
			delete(documents, pageID)
		}
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
