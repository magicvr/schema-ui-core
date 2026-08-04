package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// pageDocument is the minimal slice of a page document asserted by the schema
// endpoint tests (meta identity + a body node with a type).
type pageDocument struct {
	Meta struct {
		PageID          string `json:"pageId"`
		Title           string `json:"title"`
		ProtocolVersion string `json:"protocolVersion"`
	} `json:"meta"`
	Body map[string]any `json:"body"`
}

func TestSchemaEndpoint(t *testing.T) {
	mux := http.NewServeMux()
	schemasHandler(mux)

	t.Run("serves a seeded page document", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/schema/overview", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
		if got := rec.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
			t.Fatalf("content-type = %q, want application/json; charset=utf-8", got)
		}
		var doc pageDocument
		if err := json.NewDecoder(rec.Body).Decode(&doc); err != nil {
			t.Fatalf("body is not valid JSON: %v", err)
		}
		if doc.Meta.PageID != "overview" {
			t.Fatalf("meta.pageId = %q, want overview", doc.Meta.PageID)
		}
		if doc.Meta.Title == "" {
			t.Fatalf("meta.title is empty")
		}
		if doc.Meta.ProtocolVersion != "2.7" {
			t.Fatalf("meta.protocolVersion = %q, want 2.7", doc.Meta.ProtocolVersion)
		}
		if doc.Body["type"] == "" {
			t.Fatalf("body.type is missing")
		}
	})

	t.Run("every seeded pageId serves valid page JSON", func(t *testing.T) {
		documents := staticSchemaDocuments()
		if len(documents) == 0 {
			t.Fatalf("no seeded schema documents")
		}
		for pageID := range documents {
			req := httptest.NewRequest(http.MethodGet, "/api/schema/"+pageID, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
			if rec.Code != http.StatusOK {
				t.Fatalf("%s: status = %d, want 200", pageID, rec.Code)
			}
			var doc pageDocument
			if err := json.NewDecoder(rec.Body).Decode(&doc); err != nil {
				t.Fatalf("%s: body is not valid JSON: %v", pageID, err)
			}
			if doc.Meta.PageID != pageID {
				t.Fatalf("%s: meta.pageId = %q, want %q", pageID, doc.Meta.PageID, pageID)
			}
		}
	})

	t.Run("serves the R1 representative page set (GOAL-004)", func(t *testing.T) {
		// The 5 hand-written examples are migrated to Schema documents that the
		// manifest already routes to (D-003); these must be embedded and served.
		representative := []string{
			"data-table",
			"search-form-table",
			"form-controls",
			"form-with-reactions",
			"users",
			"roles",
			"settings",
			"activity",
		}
		documents := staticSchemaDocuments()
		for _, pageID := range representative {
			raw, ok := documents[pageID]
			if !ok {
				t.Fatalf("%s: representative page fixture is missing", pageID)
			}
			var doc pageDocument
			if err := json.Unmarshal(raw, &doc); err != nil {
				t.Fatalf("%s: fixture is not valid JSON: %v", pageID, err)
			}
			if doc.Meta.PageID != pageID {
				t.Fatalf("%s: meta.pageId = %q, want %q", pageID, doc.Meta.PageID, pageID)
			}
			if doc.Body["type"] == "" {
				t.Fatalf("%s: body.type is missing", pageID)
			}
		}
	})

	t.Run("unknown pageId returns 404 SCHEMA_NOT_FOUND", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/schema/does-not-exist", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", rec.Code)
		}
		var body map[string]string
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatalf("error body is not valid JSON: %v", err)
		}
		if body["error"] != "SCHEMA_NOT_FOUND" {
			t.Fatalf("error = %q, want SCHEMA_NOT_FOUND", body["error"])
		}
	})

	t.Run("non-GET method is not matched", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/schema/overview", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code == http.StatusOK {
			t.Fatalf("POST matched GET route; status = %d", rec.Code)
		}
	})

	// GOAL-012 / A-005 F-001: every pageId declared in the checked-in web
	// app-manifest must have an embed fixture so default Shell navigation never
	// advertises a SCHEMA_NOT_FOUND route.
	t.Run("checked-in web manifest pageIds all have embed fixtures", func(t *testing.T) {
		pageIDs := loadCheckedInManifestPageIDs(t)
		documents := staticSchemaDocuments()
		for _, pageID := range pageIDs {
			if _, ok := documents[pageID]; !ok {
				t.Fatalf("manifest pageId %q has no embed fixture under fixtures/schema/", pageID)
			}
			req := httptest.NewRequest(http.MethodGet, "/api/schema/"+pageID, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
			if rec.Code != http.StatusOK {
				t.Fatalf("%s: status = %d, want 200 (manifest must not reference dead schema routes)", pageID, rec.Code)
			}
		}
	})
}

// loadCheckedInManifestPageIDs reads the production web app-manifest and returns
// each page.pageId. Path is resolved from this test file so the check works
// regardless of the process working directory.
func loadCheckedInManifestPageIDs(t *testing.T) []string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	// apps/api/internal/handler → apps/web/public/.well-known/schema-ui/app-manifest.json
	manifestPath := filepath.Clean(filepath.Join(filepath.Dir(file),
		"..", "..", "..", "web", "public", ".well-known", "schema-ui", "app-manifest.json"))
	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read checked-in app-manifest: %v", err)
	}
	var manifest struct {
		Pages []struct {
			PageID    string `json:"pageId"`
			SchemaURL string `json:"schemaUrl"`
		} `json:"pages"`
	}
	if err := json.Unmarshal(raw, &manifest); err != nil {
		t.Fatalf("parse checked-in app-manifest: %v", err)
	}
	if len(manifest.Pages) == 0 {
		t.Fatal("checked-in app-manifest has no pages")
	}
	ids := make([]string, 0, len(manifest.Pages))
	for _, page := range manifest.Pages {
		if page.PageID == "" {
			t.Fatalf("manifest page missing pageId: %+v", page)
		}
		// Schema documents are served at /api/schema/{pageId}; enforce the
		// production convention so schemaUrl cannot silently diverge.
		wantURL := "/api/schema/" + page.PageID
		if page.SchemaURL != wantURL {
			t.Fatalf("page %q schemaUrl = %q, want %q", page.PageID, page.SchemaURL, wantURL)
		}
		if strings.Contains(page.PageID, "/") {
			t.Fatalf("pageId must not contain path separators: %q", page.PageID)
		}
		ids = append(ids, page.PageID)
	}
	return ids
}
