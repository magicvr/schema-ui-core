package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

// memEntity is an in-memory ResourceEntity used to prove the generic factory
// serves a resource with no hand-written handler (GOAL-010 S2 genericity). The
// SQLite-backed catalog instance is S4's job.
type memEntity struct {
	rows map[string]map[string]any
}

func newMemEntity(seed ...map[string]any) *memEntity {
	e := &memEntity{rows: map[string]map[string]any{}}
	for _, row := range seed {
		if id, _ := row["id"].(string); id != "" {
			e.rows[id] = row
		}
	}
	return e
}

func (e *memEntity) List(_ resourceFilter) ([]map[string]any, int, error) {
	items := make([]map[string]any, 0, len(e.rows))
	for _, row := range e.rows {
		items = append(items, row)
	}
	return items, len(items), nil
}

func (e *memEntity) Get(id string) (map[string]any, error) {
	row, ok := e.rows[id]
	if !ok {
		return nil, store.ErrNotFound
	}
	return row, nil
}

func (e *memEntity) Create(body map[string]any, id string, now time.Time, _ account.User) (map[string]any, error) {
	row := map[string]any{
		"id":        id,
		"sku":       stringField(body, "sku"),
		"title":     stringField(body, "title"),
		"updatedAt": now.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
	}
	e.rows[id] = row
	return row, nil
}

func (e *memEntity) Update(id string, body map[string]any, now time.Time, _ account.User) (map[string]any, error) {
	row, ok := e.rows[id]
	if !ok {
		return nil, store.ErrNotFound
	}
	if v, ok := body["title"]; ok {
		row["title"] = v
	}
	row["updatedAt"] = now.UTC().Format("2006-01-02T15:04:05.000Z07:00")
	return row, nil
}

func (e *memEntity) Delete(id string, _ account.User) error {
	if _, ok := e.rows[id]; !ok {
		return store.ErrNotFound
	}
	delete(e.rows, id)
	return nil
}

// catalogResource is a synthetic second resource (different path, fields, sort
// whitelist and id format) reusing the seeded admin grants, so the genericity
// test exercises the factory with no hand-written handler and no new grants.
func catalogResource(entity ResourceEntity) Resource {
	next := 0
	return Resource{
		ID:              "catalog",
		Path:            "/api/catalog",
		Listable:        true,
		SortFields:      []string{"sku", "title"},
		QSearch:         false,
		Entity:          entity,
		CreateFields:    []string{"sku", "title"},
		PatchFields:     []string{"title"},
		PermissionRead:  "records.read",
		PermissionWrite: "records.write",
		NewID: func() (string, error) {
			next++
			return fmt.Sprintf("cat-%d", next), nil
		},
	}
}

func expectError(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantCode string) {
	t.Helper()
	if rr.Code != wantStatus {
		t.Fatalf("status = %d, want %d: %s", rr.Code, wantStatus, rr.Body.String())
	}
	var body map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&body)
	if body["error"] != wantCode {
		t.Fatalf("error = %v, want %s", body["error"], wantCode)
	}
}

// GOAL-010 S2 · a resource registered through the generic factory gets the full
// five-route CRUD surface with no hand-written handler.
func TestResourceFactoryServesNewResourceWithoutHandwrittenHandler(t *testing.T) {
	env := newAuthTestEnv(t)
	entity := newMemEntity(
		map[string]any{"id": "cat-1", "sku": "S-1", "title": "Widget", "updatedAt": "2026-08-03T00:00:00.000Z"},
	)
	mux := http.NewServeMux()
	registerResource(mux, env.a, catalogResource(entity))
	token := adminToken(t, env)

	// list
	req := bearer(t, token, http.MethodGet, "/api/catalog", "")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("list status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	var list map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&list)
	items, ok := list["items"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("list items = %v, want 1", list["items"])
	}
	if list["total"] != float64(1) {
		t.Fatalf("list total = %v, want 1", list["total"])
	}

	// create → 201 with a generated cat-N id
	req = bearer(t, token, http.MethodPost, "/api/catalog", `{"sku":"S-2","title":"Gadget"}`)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201: %s", rr.Code, rr.Body.String())
	}
	var created map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&created)
	id, _ := created["id"].(string)
	if !strings.HasPrefix(id, "cat-") {
		t.Fatalf("created id = %q, want cat- prefix", id)
	}

	// detail + patch
	req = bearer(t, token, http.MethodGet, "/api/catalog/"+id, "")
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("detail status = %d, want 200", rr.Code)
	}
	req = bearer(t, token, http.MethodPatch, "/api/catalog/"+id, `{"title":"Gadget Pro"}`)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("patch status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	var updated map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&updated)
	if updated["title"] != "Gadget Pro" {
		t.Fatalf("patched title = %v, want Gadget Pro", updated["title"])
	}

	// delete → 204, then detail 404 with the default {ID}_NOT_FOUND code
	req = bearer(t, token, http.MethodDelete, "/api/catalog/"+id, "")
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want 204", rr.Code)
	}
	req = bearer(t, token, http.MethodGet, "/api/catalog/"+id, "")
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("detail after delete status = %d, want 404", rr.Code)
	}
	var errBody map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&errBody)
	if errBody["error"] != "CATALOG_NOT_FOUND" {
		t.Fatalf("not found error = %v, want CATALOG_NOT_FOUND", errBody["error"])
	}
}

// GOAL-010 S2 · the shared gates (sort/page validation, create/patch field
// validation, body error, permission) apply uniformly to every registered
// resource.
func TestResourceFactorySharedGates(t *testing.T) {
	env := newAuthTestEnv(t)
	entity := newMemEntity()
	mux := http.NewServeMux()
	registerResource(mux, env.a, catalogResource(entity))
	token := adminToken(t, env)

	req := bearer(t, token, http.MethodGet, "/api/catalog?sort=unknown", "")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	expectError(t, rr, http.StatusBadRequest, "INVALID_SORT_FIELD")

	req = bearer(t, token, http.MethodPost, "/api/catalog", `{"title":"no sku"}`)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	expectError(t, rr, http.StatusBadRequest, "INVALID_CREATE_FIELD")

	req = bearer(t, token, http.MethodPatch, "/api/catalog/cat-1", `{"title":"  "}`)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	expectError(t, rr, http.StatusBadRequest, "INVALID_PATCH_FIELD")

	req = bearer(t, token, http.MethodPost, "/api/catalog", `not json`)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	expectError(t, rr, http.StatusBadRequest, "INVALID_CREATE_BODY")

	// QSearch=false: the q param is ignored, not rejected.
	req = bearer(t, token, http.MethodGet, "/api/catalog?q=nope", "")
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("q on a non-search resource status = %d, want 200", rr.Code)
	}

	// permission gate: viewer (records.read, no records.write) lists but not writes.
	env.addUser(t, "viewer2", "pw", []string{"viewer"})
	vToken := env.login(t, "viewer2", "pw")
	req = bearer(t, vToken, http.MethodGet, "/api/catalog", "")
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("viewer list status = %d, want 200", rr.Code)
	}
	req = bearer(t, vToken, http.MethodPost, "/api/catalog", `{"sku":"S-9","title":"Denied"}`)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	expectError(t, rr, http.StatusForbidden, "FORBIDDEN")
}

// GOAL-010 S2 · permission keys default to "{id}.read" / "{id}.write" when a
// resource declares none, and the factory enforces the derived key.
func TestResourceFactoryDefaultPermissionDerivation(t *testing.T) {
	env := newAuthTestEnv(t)
	mux := http.NewServeMux()
	registerResource(mux, env.a, Resource{
		ID:           "widget",
		Path:         "/api/widget",
		Listable:     true,
		SortFields:   []string{"name"},
		Entity:       newMemEntity(),
		CreateFields: []string{"name"},
		PatchFields:  []string{"name"},
		NewID:        func() (string, error) { return "w-1", nil },
	})
	token := adminToken(t, env)

	// admin holds records.read/write, not the derived widget.read → 403.
	req := bearer(t, token, http.MethodGet, "/api/widget", "")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	expectError(t, rr, http.StatusForbidden, "FORBIDDEN")
}
