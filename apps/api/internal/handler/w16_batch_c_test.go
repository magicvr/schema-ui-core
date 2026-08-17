package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCronPreviewEndpoint(t *testing.T) {
	env := newAuthTestEnv(t)
	req := bearer(t, adminToken(t, env), http.MethodPost, "/api/scheduled-tasks/cron/preview", `{"cron":"0 2 * * *"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("cron preview status = %d: %s", rr.Code, rr.Body.String())
	}
	var body struct {
		Description string   `json:"description"`
		NextRuns    []string `json:"nextRuns"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if len(body.NextRuns) != 3 {
		t.Fatalf("nextRuns = %d, want 3", len(body.NextRuns))
	}
	for _, run := range body.NextRuns {
		if !strings.Contains(run, "T") || !strings.HasSuffix(run, "Z") {
			t.Fatalf("nextRuns entry %q not RFC3339 UTC", run)
		}
	}
}

func TestDictEntryBadgeStylePersists(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	// Create a dictionary type.
	typeReq := bearer(t, token, http.MethodPost, "/api/data-dictionary/types", `{"key":"status","name":"Status"}`)
	typeRR := httptest.NewRecorder()
	env.mux.ServeHTTP(typeRR, typeReq)
	if typeRR.Code != http.StatusCreated {
		t.Fatalf("create dict type = %d: %s", typeRR.Code, typeRR.Body.String())
	}
	// Create an entry with a badge style.
	entryReq := bearer(t, token, http.MethodPost, "/api/data-dictionary/entries", `{"dictKey":"status","entryKey":"ok","label":"OK","badgeStyle":"success"}`)
	entryRR := httptest.NewRecorder()
	env.mux.ServeHTTP(entryRR, entryReq)
	if entryRR.Code != http.StatusCreated {
		t.Fatalf("create dict entry = %d: %s", entryRR.Code, entryRR.Body.String())
	}
	var created map[string]any
	_ = json.NewDecoder(entryRR.Body).Decode(&created)
	if created["badgeStyle"] != "success" {
		t.Fatalf("created badgeStyle = %v, want success", created["badgeStyle"])
	}
	id, _ := created["id"].(string)
	if id == "" {
		t.Fatal("created entry missing id")
	}
	// Fetch it back.
	getReq := bearer(t, token, http.MethodGet, "/api/data-dictionary/entries/"+id, "")
	getRR := httptest.NewRecorder()
	env.mux.ServeHTTP(getRR, getReq)
	if getRR.Code != http.StatusOK {
		t.Fatalf("get dict entry = %d: %s", getRR.Code, getRR.Body.String())
	}
	var fetched map[string]any
	_ = json.NewDecoder(getRR.Body).Decode(&fetched)
	if fetched["badgeStyle"] != "success" {
		t.Fatalf("fetched badgeStyle = %v, want success", fetched["badgeStyle"])
	}
}

func TestSettingsFooterFieldsPersist(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	patch := bearer(t, token, http.MethodPatch, "/api/settings/default", `{"copyrightText":"© 2026 Example","icpNumber":"ICP-123"}`)
	pr := httptest.NewRecorder()
	env.mux.ServeHTTP(pr, patch)
	if pr.Code != http.StatusOK {
		t.Fatalf("patch settings = %d: %s", pr.Code, pr.Body.String())
	}
	var patched map[string]any
	_ = json.NewDecoder(pr.Body).Decode(&patched)
	if patched["copyrightText"] != "© 2026 Example" || patched["icpNumber"] != "ICP-123" {
		t.Fatalf("patched footer = %v", patched)
	}
	get := bearer(t, token, http.MethodGet, "/api/settings/default", "")
	gr := httptest.NewRecorder()
	env.mux.ServeHTTP(gr, get)
	if gr.Code != http.StatusOK {
		t.Fatalf("get settings = %d", gr.Code)
	}
	var got map[string]any
	_ = json.NewDecoder(gr.Body).Decode(&got)
	if got["copyrightText"] != "© 2026 Example" || got["icpNumber"] != "ICP-123" {
		t.Fatalf("fetched footer = %v", got)
	}
}
