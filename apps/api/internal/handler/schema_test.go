package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
}
