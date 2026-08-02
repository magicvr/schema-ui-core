package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

// msRFC3339 matches the frozen fixed-3-digit-millisecond updatedAt shape
// (GOAL-007 D-004): 2006-01-02T15:04:05.123Z.
var msRFC3339 = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$`)

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
	// T-API-08 · anonymous POST is 401 too.
	req = httptest.NewRequest(http.MethodPost, "/api/records", strings.NewReader(`{"name":"x","status":"a","owner":"b"}`))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("POST status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}

func TestRecordsWriteDeniedWithoutAdminRole(t *testing.T) {
	env := newAuthTestEnv(t)
	env.addUser(t, "editor", "pw", []string{"editor"})
	token := env.login(t, "editor", "pw")
	for _, tc := range []struct{ method, path string }{
		{http.MethodPatch, "/api/records/rec-3"},
		{http.MethodDelete, "/api/records/rec-3"},
		{http.MethodPost, "/api/records"},
	} {
		req := bearer(t, token, tc.method, tc.path, `{"name":"x"}`)
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("%s status = %d, want %d", tc.method, rr.Code, http.StatusForbidden)
		}
		var body map[string]any
		_ = jsonDecode(rr, &body)
		if body["error"] != "FORBIDDEN" {
			t.Fatalf("%s error = %v, want FORBIDDEN", tc.method, body["error"])
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
	for _, tc := range []struct{ method, path string }{
		{http.MethodPatch, "/api/records/rec-3"},
		{http.MethodDelete, "/api/records/rec-3"},
		{http.MethodPost, "/api/records"},
	} {
		req := bearer(t, token, tc.method, tc.path, `{"name":"x"}`)
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("%s status = %d, want 403", tc.method, rr.Code)
		}
		var body map[string]any
		_ = jsonDecode(rr, &body)
		if body["error"] != "FORBIDDEN" {
			t.Fatalf("%s error = %v, want FORBIDDEN", tc.method, body["error"])
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

// T-API-10 · POST with a valid body returns 201 + the full record (server id and
// updatedAt), and the create is visible through list/detail.
func TestRecordsCreate(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	req := bearer(t, token, http.MethodPost, "/api/records", `{"name":"Zephyr Labs","status":"active","owner":"nadia"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201: %s", rr.Code, rr.Body.String())
	}
	var rec map[string]any
	if err := jsonDecode(rr, &rec); err != nil {
		t.Fatalf("decode create: %v", err)
	}
	id, _ := rec["id"].(string)
	if !strings.HasPrefix(id, "rec-") || len(id) != len("rec-")+16 {
		t.Fatalf("id = %q, want rec- + 16 hex chars", id)
	}
	if rec["name"] != "Zephyr Labs" || rec["status"] != "active" || rec["owner"] != "nadia" {
		t.Fatalf("created = %v, want name/status/owner set", rec)
	}
	ua, _ := rec["updatedAt"].(string)
	if !msRFC3339.MatchString(ua) {
		t.Fatalf("created updatedAt %q, want fixed-3-ms RFC3339", ua)
	}
	// The new record is served by detail and counted by list.
	code, detail := getRecords(t, env, "/api/records/"+id)
	if code != http.StatusOK || detail["name"] != "Zephyr Labs" {
		t.Fatalf("detail = %d %v, want 200 Zephyr Labs", code, detail)
	}
	_, list := getRecords(t, env, "/api/records")
	if list["total"] != float64(9) {
		t.Fatalf("total = %v, want 9 after create", list["total"])
	}
}

// T-API-11 · POST with missing, blank or non-string required fields → 400
// INVALID_CREATE_FIELD (per-field message names the field).
func TestRecordsCreateInvalidField(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	for label, body := range map[string]string{
		"missing-name":     `{"status":"active","owner":"a"}`,
		"blank-name":       `{"name":"   ","status":"active","owner":"a"}`,
		"missing-status":   `{"name":"x","owner":"a"}`,
		"blank-status":     `{"name":"x","status":"","owner":"a"}`,
		"missing-owner":    `{"name":"x","status":"active"}`,
		"blank-owner":      `{"name":"x","status":"active","owner":" "}`,
		"non-string-name":  `{"name":123,"status":"active","owner":"a"}`,
		"null-name":        `{"name":null,"status":"active","owner":"a"}`,
	} {
		req := bearer(t, token, http.MethodPost, "/api/records", body)
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("%s: status = %d, want 400", label, rr.Code)
		}
		var out map[string]any
		_ = jsonDecode(rr, &out)
		if out["error"] != "INVALID_CREATE_FIELD" {
			t.Fatalf("%s: error = %v, want INVALID_CREATE_FIELD", label, out["error"])
		}
	}
}

// T-API-12 · POST with non-JSON, truncated or oversized bodies → 400
// INVALID_CREATE_BODY.
func TestRecordsCreateInvalidBody(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	for label, body := range map[string]string{
		"not-json":  `not json`,
		"array":     `[1,2,3]`,
		"truncated": `{"name":"x"`,
		"oversized": `{"name":"` + strings.Repeat("x", maxRecordBodyBytes) + `","status":"a","owner":"b"}`,
	} {
		req := bearer(t, token, http.MethodPost, "/api/records", body)
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("%s: status = %d, want 400", label, rr.Code)
		}
		var out map[string]any
		_ = jsonDecode(rr, &out)
		if out["error"] != "INVALID_CREATE_BODY" {
			t.Fatalf("%s: error = %v, want INVALID_CREATE_BODY", label, out["error"])
		}
	}
}

// T-API-13 · admin (records.read + records.write) walks the whole lifecycle:
// create → list → detail → patch → delete → 404.
func TestRecordsAdminLifecycle(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)

	req := bearer(t, token, http.MethodPost, "/api/records", `{"name":"Lifecycle Co","status":"pending","owner":"omar"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201", rr.Code)
	}
	var created map[string]any
	_ = jsonDecode(rr, &created)
	id, _ := created["id"].(string)

	if code, list := getRecords(t, env, "/api/records?q=lifecycle"); code != http.StatusOK || list["total"] != float64(1) {
		t.Fatalf("search = %d total=%v, want 200/1", code, list["total"])
	}
	req = bearer(t, token, http.MethodPatch, "/api/records/"+id, `{"name":"Lifecycle Rebranded"}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("patch status = %d, want 200", rr.Code)
	}
	req = bearer(t, token, http.MethodDelete, "/api/records/"+id, "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want 204", rr.Code)
	}
	if code, _ := getRecords(t, env, "/api/records/"+id); code != http.StatusNotFound {
		t.Fatalf("detail after delete = %d, want 404", code)
	}
}

// T-API-05 (R-001) · consecutive updates strictly increase updatedAt without
// sleeps: millisecond precision plus the monotonic clamp make the assertion
// deterministic even for same-millisecond back-to-back writes.
func TestRecordsUpdateStrictlyIncreasesAcrossRapidPatches(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	var prev time.Time
	for i := 0; i < 4; i++ {
		req := bearer(t, token, http.MethodPatch, "/api/records/rec-3", fmt.Sprintf(`{"name":"Patch %d"}`, i))
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("patch %d status = %d, want 200: %s", i, rr.Code, rr.Body.String())
		}
		var body map[string]any
		_ = jsonDecode(rr, &body)
		ua, _ := body["updatedAt"].(string)
		if !msRFC3339.MatchString(ua) {
			t.Fatalf("patch %d updatedAt %q, want fixed-3-ms RFC3339", i, ua)
		}
		now, err := time.Parse(time.RFC3339, ua)
		if err != nil {
			t.Fatalf("parse patch %d updatedAt %q: %v", i, ua, err)
		}
		if !prev.IsZero() && !now.After(prev) {
			t.Fatalf("patch %d updatedAt %v not strictly after %v", i, now, prev)
		}
		prev = now
	}
}

// T-DB-09 · the handler is backed by SQLite: a record written directly through
// the store is served by the HTTP handler — there is no in-process slice
// fallback in the production path.
func TestRecordsHandlerReadsFromStore(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	if _, err := env.st.CreateRecord(store.Record{
		ID: "rec-direct", Name: "Direct Write", Status: "active", Owner: "zed",
		UpdatedAt: time.Date(2026, 8, 2, 12, 0, 0, 0, time.UTC),
	}); err != nil {
		t.Fatalf("store create: %v", err)
	}
	code, detail := getRecords(t, env, "/api/records/rec-direct")
	if code != http.StatusOK || detail["name"] != "Direct Write" {
		t.Fatalf("handler detail = %d %v, want 200 Direct Write", code, detail)
	}
	// And a record the handler created is present in the store's table (round
	// trip through the shared SQLite database, not a per-handler copy).
	req := bearer(t, token, http.MethodPost, "/api/records", `{"name":"Via Handler","status":"active","owner":"yara"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create via handler = %d, want 201", rr.Code)
	}
	var created map[string]any
	_ = jsonDecode(rr, &created)
	id, _ := created["id"].(string)
	got, err := env.st.GetRecord(id)
	if err != nil || got.Name != "Via Handler" {
		t.Fatalf("store read of handler create = %v (err %v), want Via Handler", got, err)
	}
}

// A-003 R-002 · PATCH stores trimmed field values (consistent with create):
// leading/trailing whitespace on an accepted value is normalized.
func TestRecordsUpdateTrimsValues(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	req := bearer(t, token, http.MethodPatch, "/api/records/rec-3", `{"name":"  Hooli Rebrand  ","owner":" carol "}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	var body map[string]any
	_ = jsonDecode(rr, &body)
	if body["name"] != "Hooli Rebrand" || body["owner"] != "carol" {
		t.Fatalf("updated = %v, want trimmed Hooli Rebrand / carol", body)
	}
}

func jsonDecode(rr *httptest.ResponseRecorder, out *map[string]any) error {
	if rr.Body.Len() == 0 {
		return nil
	}
	return json.NewDecoder(rr.Body).Decode(out)
}
