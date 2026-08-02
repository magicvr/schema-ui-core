package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// adminToken logs in as the seeded admin and returns the access token.
func adminToken(t *testing.T, env *authTestEnv) string {
	t.Helper()
	return env.login(t, testSeedUsername, testSeedPassword)
}

func TestRecordsListDefault(t *testing.T) {
	env := newAuthTestEnv(t)
	code, body := getRecords(t, env, "/api/records")
	if code != http.StatusOK {
		t.Fatalf("status = %d, want %d", code, http.StatusOK)
	}
	items, ok := body["items"].([]any)
	if !ok {
		t.Fatalf("items = %v, want array", body["items"])
	}
	if len(items) != 8 {
		t.Fatalf("len(items) = %d, want 8", len(items))
	}
	if body["total"] != float64(8) {
		t.Fatalf("total = %v, want 8", body["total"])
	}
	if body["pageSize"] != float64(10) {
		t.Fatalf("pageSize = %v, want 10", body["pageSize"])
	}
	// Default sort is name asc: first item is "Acme Console".
	first, ok := items[0].(map[string]any)
	if !ok {
		t.Fatalf("first item = %v, want object", items[0])
	}
	if first["name"] != "Acme Console" {
		t.Fatalf("first name = %v, want Acme Console", first["name"])
	}
}

func TestRecordsListSearch(t *testing.T) {
	env := newAuthTestEnv(t)
	code, body := getRecords(t, env, "/api/records?q=alice")
	if code != http.StatusOK {
		t.Fatalf("status = %d, want %d", code, http.StatusOK)
	}
	items, ok := body["items"].([]any)
	if !ok {
		t.Fatalf("items = %v, want array", body["items"])
	}
	for _, item := range items {
		rec := item.(map[string]any)
		if rec["owner"] != "alice" && !strings.Contains(rec["name"].(string), "alice") {
			t.Fatalf("item %v does not match q=alice", rec)
		}
	}
}

func TestRecordsListSortDesc(t *testing.T) {
	env := newAuthTestEnv(t)
	code, body := getRecords(t, env, "/api/records?sort=updatedAt&order=desc")
	if code != http.StatusOK {
		t.Fatalf("status = %d, want %d", code, http.StatusOK)
	}
	items, ok := body["items"].([]any)
	if !ok {
		t.Fatalf("items = %v, want array", body["items"])
	}
	last, _ := items[0].(map[string]any)
	if last["name"] != "Globex Admin" {
		t.Fatalf("first (desc) name = %v, want Globex Admin", last["name"])
	}
}

func TestRecordsListPagination(t *testing.T) {
	env := newAuthTestEnv(t)
	code, body := getRecords(t, env, "/api/records?page=2&pageSize=3")
	if code != http.StatusOK {
		t.Fatalf("status = %d, want %d", code, http.StatusOK)
	}
	items, ok := body["items"].([]any)
	if !ok {
		t.Fatalf("items = %v, want array", body["items"])
	}
	if len(items) != 3 {
		t.Fatalf("len(items) = %d, want 3", len(items))
	}
	if body["total"] != float64(8) {
		t.Fatalf("total = %v, want 8", body["total"])
	}
	first, _ := items[0].(map[string]any)
	if first["name"] != "Initech Reports" {
		t.Fatalf("page 2 first name = %v, want Initech Reports", first["name"])
	}
}

func TestRecordsListInvalidParamsFailClosed(t *testing.T) {
	env := newAuthTestEnv(t)
	for _, path := range []string{
		"/api/records?sort=unknown",
		"/api/records?order=up",
		"/api/records?page=0",
		"/api/records?pageSize=abc",
	} {
		code, body := getRecords(t, env, path)
		if code != http.StatusBadRequest {
			t.Fatalf("%s: status = %d, want %d", path, code, http.StatusBadRequest)
		}
		if _, ok := body["error"]; !ok {
			t.Fatalf("%s: error code missing in %v", path, body)
		}
	}
}

func TestRecordsListPageSizeCap(t *testing.T) {
	env := newAuthTestEnv(t)
	code, body := getRecords(t, env, "/api/records?pageSize=1000")
	if code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", code, http.StatusBadRequest)
	}
	if body["error"] != "INVALID_PAGE_SIZE" {
		t.Fatalf("error = %v, want INVALID_PAGE_SIZE", body["error"])
	}
}

func TestRecordsWriteRequiresAuth(t *testing.T) {
	env := newAuthTestEnv(t)
	code, body := sendJSON(t, env.mux, http.MethodPatch, "/api/records/rec-3", `{"name":"x"}`)
	if code != http.StatusUnauthorized {
		t.Fatalf("PATCH status = %d, want %d", code, http.StatusUnauthorized)
	}
	if body["error"] != "UNAUTHENTICATED" {
		t.Fatalf("PATCH error = %v, want UNAUTHENTICATED", body["error"])
	}
	req := httptest.NewRequest(http.MethodDelete, "/api/records/rec-3", nil)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("DELETE status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}

func TestRecordsWriteDeniedWithoutAdminRole(t *testing.T) {
	env := newAuthTestEnv(t)
	env.addUser(t, "editor", "pw", []string{"editor"})
	token := env.login(t, "editor", "pw")
	for _, method := range []string{http.MethodPatch, http.MethodDelete} {
		req := bearer(t, token, method, "/api/records/rec-3", `{"name":"x"}`)
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("%s status = %d, want %d", method, rr.Code, http.StatusForbidden)
		}
		var body map[string]any
		_ = jsonDecode(rr, &body)
		if body["error"] != "FORBIDDEN" {
			t.Fatalf("%s error = %v, want FORBIDDEN", method, body["error"])
		}
	}
}

func TestRecordsUpdateBodyTooLarge(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	huge := `{"name":"` + strings.Repeat("x", maxRecordBodyBytes) + `"}`
	req := bearer(t, token, http.MethodPatch, "/api/records/rec-3", huge)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestRecordsDetail(t *testing.T) {
	env := newAuthTestEnv(t)
	code, body := getRecords(t, env, "/api/records/rec-3")
	if code != http.StatusOK {
		t.Fatalf("status = %d, want %d", code, http.StatusOK)
	}
	if body["id"] != "rec-3" {
		t.Fatalf("id = %v, want rec-3", body["id"])
	}
}

func TestRecordsDetailNotFound(t *testing.T) {
	env := newAuthTestEnv(t)
	code, body := getRecords(t, env, "/api/records/rec-999")
	if code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", code, http.StatusNotFound)
	}
	if body["error"] != "RECORD_NOT_FOUND" {
		t.Fatalf("error = %v, want RECORD_NOT_FOUND", body["error"])
	}
}

func TestRecordsUpdate(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	req := bearer(t, token, http.MethodPatch, "/api/records/rec-3", `{"name":"Hooli Rebrand","status":"archived"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", rr.Code, http.StatusOK, rr.Body.String())
	}
	var body map[string]any
	_ = jsonDecode(rr, &body)
	if body["name"] != "Hooli Rebrand" || body["status"] != "archived" {
		t.Fatalf("updated = %v, want name=Hooli Rebrand status=archived", body)
	}
	// Persisted: a subsequent GET reflects the patch.
	_, detail := getRecords(t, env, "/api/records/rec-3")
	if detail["name"] != "Hooli Rebrand" {
		t.Fatalf("detail name = %v, want Hooli Rebrand", detail["name"])
	}
}

func TestRecordsUpdateRefreshesUpdatedAt(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)

	_, before := getRecords(t, env, "/api/records/rec-3")
	beforeValue, ok := before["updatedAt"].(string)
	if !ok {
		t.Fatalf("updatedAt = %v, want string", before["updatedAt"])
	}
	beforeTime, err := time.Parse(time.RFC3339, beforeValue)
	if err != nil {
		t.Fatalf("parse before updatedAt %q: %v", beforeValue, err)
	}

	req := bearer(t, token, http.MethodPatch, "/api/records/rec-3", `{"name":"Refreshed"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	var body map[string]any
	_ = jsonDecode(rr, &body)
	afterValue, ok := body["updatedAt"].(string)
	if !ok {
		t.Fatalf("updatedAt = %v, want string", body["updatedAt"])
	}
	afterTime, err := time.Parse(time.RFC3339, afterValue)
	if err != nil {
		t.Fatalf("parse after updatedAt %q: %v", afterValue, err)
	}
	if !afterTime.After(beforeTime) {
		t.Fatalf("updatedAt = %v, want after %v", afterValue, beforeValue)
	}
}

func TestRecordsUpdateInvalidFailClosed(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	for _, body := range []string{
		`{"name":""}`,
		`{"status":""}`,
		`not json`,
	} {
		req := bearer(t, token, http.MethodPatch, "/api/records/rec-3", body)
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("%s: status = %d, want %d", body, rr.Code, http.StatusBadRequest)
		}
	}
}

func TestRecordsUpdateNotFound(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	req := bearer(t, token, http.MethodPatch, "/api/records/rec-999", `{"name":"x"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
	var body map[string]any
	_ = jsonDecode(rr, &body)
	if body["error"] != "RECORD_NOT_FOUND" {
		t.Fatalf("error = %v, want RECORD_NOT_FOUND", body["error"])
	}
}

func TestRecordsDelete(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	req := bearer(t, token, http.MethodDelete, "/api/records/rec-3", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNoContent)
	}
	// Removed from the list and detail now 404s.
	_, list := getRecords(t, env, "/api/records")
	if list["total"] != float64(7) {
		t.Fatalf("total = %v, want 7 after delete", list["total"])
	}
	code, _ := getRecords(t, env, "/api/records/rec-3")
	if code != http.StatusNotFound {
		t.Fatalf("detail status = %d, want %d after delete", code, http.StatusNotFound)
	}
}

func TestRecordsDeleteNotFound(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	req := bearer(t, token, http.MethodDelete, "/api/records/rec-999", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}

// GOAL-006 S4 · anonymous read is 401 (reads are now gated too).
func TestRecordsReadRequiresAuth(t *testing.T) {
	env := newAuthTestEnv(t)
	for _, path := range []string{"/api/records", "/api/records/rec-3"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("GET %s status = %d, want 401", path, rr.Code)
		}
	}
}

// GOAL-006 S4 · viewer can read (records.read) but cannot write (no
// records.write): the gate checks persisted permission keys, not roles.
func TestRecordsViewerCanReadNotWrite(t *testing.T) {
	env := newAuthTestEnv(t)
	env.addUser(t, "viewer", "pw", []string{"viewer"})
	token := env.login(t, "viewer", "pw")

	req := bearer(t, token, http.MethodGet, "/api/records", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("GET status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	for _, method := range []string{http.MethodPatch, http.MethodDelete} {
		req := bearer(t, token, method, "/api/records/rec-3", `{"name":"x"}`)
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("%s status = %d, want 403", method, rr.Code)
		}
		var body map[string]any
		_ = jsonDecode(rr, &body)
		if body["error"] != "FORBIDDEN" {
			t.Fatalf("%s error = %v, want FORBIDDEN", method, body["error"])
		}
	}
}

// GOAL-006 S4 · an authenticated user whose roles grant no permission is denied
// on read (records.read missing).
func TestRecordsReadDeniedWithoutPermission(t *testing.T) {
	env := newAuthTestEnv(t)
	// A derived role with no grants: carries no permissions.
	env.addUser(t, "auditor", "pw", []string{"custom"})
	token := env.login(t, "auditor", "pw")
	req := bearer(t, token, http.MethodGet, "/api/records", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("GET status = %d, want 403 (no records.read)", rr.Code)
	}
}

func jsonDecode(rr *httptest.ResponseRecorder, out *map[string]any) error {
	if rr.Body.Len() == 0 {
		return nil
	}
	return json.NewDecoder(rr.Body).Decode(out)
}
