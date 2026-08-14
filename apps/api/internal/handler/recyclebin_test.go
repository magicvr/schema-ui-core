// Recycle bin handler tests (S-12 · GOAL-012 D-002 §3): list/detail,
// restore, purge, permission gates and the audit events. The fake service is
// a test double (the module package imports handler, so handler tests cannot
// import the real service); the module tests cover the store-backed service.
package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	datadictionarystore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/datadictionary/store"
)

func recycleItem(id, resource, resourceID string) RecycleItem {
	return RecycleItem{
		ID:         id,
		Resource:   resource,
		ResourceID: resourceID,
		Payload:    map[string]any{"id": resourceID, "key": resourceID + "-key"},
		ActorID:    "user-admin",
		ActorName:  "Admin",
		DeletedAt:  time.Now().UTC(),
	}
}

func TestRecycleBinListAndDetail(t *testing.T) {
	env := newAuthTestEnv(t)
	env.recycle.add(recycleItem("recycle-1", "dict-types", "t1"))
	env.recycle.add(recycleItem("recycle-2", "scheduled-tasks", "task-a"))
	token := adminToken(t, env)

	req := bearer(t, token, http.MethodGet, "/api/recycle-bin", "")
	rec := httptest.NewRecorder()
	env.mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("list = %d: %s", rec.Code, rec.Body.String())
	}
	var body struct {
		Items []map[string]any `json:"items"`
		Total int              `json:"total"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Total != 2 || len(body.Items) != 2 {
		t.Fatalf("items = %d/%d, want 2", len(body.Items), body.Total)
	}

	req = bearer(t, token, http.MethodGet, "/api/recycle-bin/recycle-1", "")
	rec = httptest.NewRecorder()
	env.mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("detail = %d: %s", rec.Code, rec.Body.String())
	}
	var detail map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &detail); err != nil {
		t.Fatalf("decode detail: %v", err)
	}
	if detail["resourceId"] != "t1" || detail["restored"] != false {
		t.Fatalf("detail = %v", detail)
	}
}

func TestRecycleBinRestoreAndPurge(t *testing.T) {
	env := newAuthTestEnv(t)
	env.recycle.add(recycleItem("recycle-1", "dict-types", "t1"))
	token := adminToken(t, env)

	req := bearer(t, token, http.MethodPost, "/api/recycle-bin/recycle-1/restore", "")
	rec := httptest.NewRecorder()
	env.mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("restore = %d: %s", rec.Code, rec.Body.String())
	}
	var row map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &row); err != nil {
		t.Fatalf("decode restore: %v", err)
	}
	if row["key"] != "t1-key" {
		t.Fatalf("restored row = %v", row)
	}
	// Item is now restored: detail reports it.
	req = bearer(t, token, http.MethodGet, "/api/recycle-bin/recycle-1", "")
	rec = httptest.NewRecorder()
	env.mux.ServeHTTP(rec, req)
	var detail map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &detail)
	if detail["restored"] != true {
		t.Fatalf("detail after restore = %v, want restored true", detail)
	}

	// Purge the second item.
	env.recycle.add(recycleItem("recycle-2", "dict-entries", "e1"))
	req = bearer(t, token, http.MethodDelete, "/api/recycle-bin/recycle-2", "")
	rec = httptest.NewRecorder()
	env.mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("purge = %d: %s", rec.Code, rec.Body.String())
	}
	req = bearer(t, token, http.MethodGet, "/api/recycle-bin/recycle-2", "")
	rec = httptest.NewRecorder()
	env.mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("detail after purge = %d, want 404", rec.Code)
	}
}

func TestRecycleBinPermissionGates(t *testing.T) {
	env := newAuthTestEnv(t)
	env.recycle.add(recycleItem("recycle-1", "dict-types", "t1"))
	token := adminToken(t, env)

	anon := httptest.NewRecorder()
	env.mux.ServeHTTP(anon, httptest.NewRequest(http.MethodGet, "/api/recycle-bin", nil))
	if anon.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous list = %d, want 401", anon.Code)
	}
	// editor has no recycle.write → restore 403.
	env.addUser(t, "ed", "editor-pass-1", []string{"editor"})
	edToken := loginAs(t, env, "ed", "editor-pass-1")
	req := bearer(t, edToken, http.MethodPost, "/api/recycle-bin/recycle-1/restore", "")
	rec := httptest.NewRecorder()
	env.mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("editor restore = %d, want 403", rec.Code)
	}
	_ = token
}

func TestRecycleBinNotFound(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	req := bearer(t, token, http.MethodGet, "/api/recycle-bin/nope", "")
	rec := httptest.NewRecorder()
	env.mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("unknown detail = %d, want 404", rec.Code)
	}
}

func TestRecycleBinAuditEvents(t *testing.T) {
	env := newAuthTestEnv(t)
	env.recycle.add(recycleItem("recycle-1", "dict-types", "t1"))
	env.recycle.add(recycleItem("recycle-2", "dict-types", "t2"))
	token := adminToken(t, env)

	req := bearer(t, token, http.MethodPost, "/api/recycle-bin/recycle-1/restore", "")
	rec := httptest.NewRecorder()
	env.mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("restore = %d", rec.Code)
	}
	req = bearer(t, token, http.MethodDelete, "/api/recycle-bin/recycle-2", "")
	rec = httptest.NewRecorder()
	env.mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("purge = %d", rec.Code)
	}

	opsReq := bearer(t, token, http.MethodGet, "/api/operations?pageSize=100", "")
	opsRec := httptest.NewRecorder()
	env.mux.ServeHTTP(opsRec, opsReq)
	if opsRec.Code != http.StatusOK {
		t.Fatalf("operations = %d", opsRec.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(opsRec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode operations: %v", err)
	}
	items, _ := body["items"].([]any)
	events := map[string]bool{}
	for _, item := range items {
		row, _ := item.(map[string]any)
		ev, _ := row["event"].(string)
		events[ev] = true
	}
	if !events["recycle.restore"] {
		t.Fatalf("operations missing recycle.restore: %v", events)
	}
	if !events["recycle.purge"] {
		t.Fatalf("operations missing recycle.purge: %v", events)
	}
}

// recordingTrash is a TrashRecorder test double capturing calls.
type recordingTrash struct {
	calls []struct {
		resource string
		id       string
		row      map[string]any
	}
}

func (r *recordingTrash) Record(_ context.Context, resource, id string, row map[string]any, _ account.User, _ time.Time) error {
	r.calls = append(r.calls, struct {
		resource string
		id       string
		row      map[string]any
	}{resource: resource, id: id, row: row})
	return nil
}

// TestRecycleFactoryHookSnapshotsOnDelete proves the S-12 delete hook: a
// successful delete records a snapshot with the pre-delete row; a failed
// delete records nothing (no orphan snapshots).
func TestRecycleFactoryHookSnapshotsOnDelete(t *testing.T) {
	env := newAuthTestEnv(t)
	trash := &recordingTrash{}
	mux := http.NewServeMux()
	for _, route := range DictionaryRoutes(env.a, datadictionarystore.NewRepository(env.st), env.operations, "admin.data-dictionary", trash) {
		mux.Handle(route.Method+" "+route.Pattern, route.Handler)
	}
	token := adminToken(t, env)
	// Create a dict type via the API (bearer-authenticated).
	createReq := bearer(t, token, http.MethodPost, "/api/data-dictionary/types", `{"key":"status","name":"Status"}`)
	createRec := httptest.NewRecorder()
	mux.ServeHTTP(createRec, createReq)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("create = %d: %s", createRec.Code, createRec.Body.String())
	}
	var createdBody struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(createRec.Body.Bytes(), &createdBody); err != nil || createdBody.ID == "" {
		t.Fatalf("create body missing id: %v", err)
	}
	id := createdBody.ID
	// Delete through the factory (with the trash-wired routes).
	req := bearer(t, token, http.MethodDelete, "/api/data-dictionary/types/"+id, "")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete = %d: %s", rec.Code, rec.Body.String())
	}
	if len(trash.calls) != 1 {
		t.Fatalf("trash calls = %d, want 1", len(trash.calls))
	}
	if trash.calls[0].resource != "dict-types" || trash.calls[0].id != id {
		t.Fatalf("trash call = %+v", trash.calls[0])
	}
	if trash.calls[0].row["key"] != "status" {
		t.Fatalf("snapshot row = %v, want key status", trash.calls[0].row)
	}
	// Deleting an unknown id fails → no snapshot recorded.
	req = bearer(t, token, http.MethodDelete, "/api/data-dictionary/types/nope", "")
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("unknown delete = %d, want 404", rec.Code)
	}
	if len(trash.calls) != 1 {
		t.Fatalf("trash calls after failed delete = %d, want 1 (no orphan)", len(trash.calls))
	}
}

// TestRecycleFactoryHookNilKeepsLegacySemantics proves Trash nil leaves the
// delete path byte-identical (the env-mounted dictionary routes have no trash).
func TestRecycleFactoryHookNilKeepsLegacySemantics(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	createReq := bearer(t, token, http.MethodPost, "/api/data-dictionary/types", `{"key":"legacy","name":"Legacy"}`)
	createRec := httptest.NewRecorder()
	env.mux.ServeHTTP(createRec, createReq)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("create = %d: %s", createRec.Code, createRec.Body.String())
	}
	var createdBody struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(createRec.Body.Bytes(), &createdBody); err != nil || createdBody.ID == "" {
		t.Fatalf("create body missing id: %v", err)
	}
	id := createdBody.ID
	req := bearer(t, token, http.MethodDelete, "/api/data-dictionary/types/"+id, "")
	rec := httptest.NewRecorder()
	env.mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete = %d: %s", rec.Code, rec.Body.String())
	}
}
