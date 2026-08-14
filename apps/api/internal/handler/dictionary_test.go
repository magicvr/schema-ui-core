// Data dictionary surface tests (S-01 · GOAL-008 D-002 §6): two-level CRUD,
// permission gates, cascade delete, dict-key validation, audit events.
package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

func dictionaryCreateType(t *testing.T, env *authTestEnv, token, key, name string) string {
	t.Helper()
	code, body := bearerJSON(t, env, token, http.MethodPost, "/api/data-dictionary/types",
		"{\"key\":\""+key+"\",\"name\":\""+name+"\",\"enabled\":true,\"sort\":0}")
	if code != http.StatusCreated {
		t.Fatalf("create type = %d: %v", code, body)
	}
	id, _ := body["id"].(string)
	if id == "" {
		t.Fatal("create type missing id")
	}
	return id
}

// Two-level CRUD lifecycle: create types, create entries, q filter, patch,
// cascade delete on the type, and the three audit events.
func TestDictionaryLifecycle(t *testing.T) {
	env := newAuthTestEnv(t)
	admin := adminToken(t, env)

	typeID := dictionaryCreateType(t, env, admin, "order_status", "Order status")
	// duplicate type key → 409 DICT_TYPE_KEY_TAKEN
	code, errBody := bearerJSON(t, env, admin, http.MethodPost, "/api/data-dictionary/types",
		"{\"key\":\"order_status\",\"name\":\"Duplicate\"}")
	if code != http.StatusConflict || errBody["error"] != "DICT_TYPE_KEY_TAKEN" {
		t.Fatalf("duplicate type = %d %v", code, errBody)
	}

	// list types
	code, list := getResource(t, env, "/api/data-dictionary/types?pageSize=100")
	if code != http.StatusOK || list["total"] != float64(1) {
		t.Fatalf("list types = %d %v", code, list)
	}

	// create entries (valid dict key)
	code, entryBody := bearerJSON(t, env, admin, http.MethodPost, "/api/data-dictionary/entries",
		"{\"dictKey\":\"order_status\",\"entryKey\":\"pending\",\"label\":\"Pending\",\"enabled\":true,\"sort\":1}")
	if code != http.StatusCreated {
		t.Fatalf("create entry = %d %v", code, entryBody)
	}
	entryID, _ := entryBody["id"].(string)
	// unknown dict key → 400 DICT_KEY_NOT_FOUND
	code, errBody = bearerJSON(t, env, admin, http.MethodPost, "/api/data-dictionary/entries",
		"{\"dictKey\":\"nope\",\"entryKey\":\"x\",\"label\":\"X\"}")
	if code != http.StatusBadRequest || errBody["error"] != "DICT_KEY_NOT_FOUND" {
		t.Fatalf("unknown dict key = %d %v", code, errBody)
	}
	// duplicate (dict_key, entry_key) → 409 DICT_ENTRY_KEY_TAKEN
	code, errBody = bearerJSON(t, env, admin, http.MethodPost, "/api/data-dictionary/entries",
		"{\"dictKey\":\"order_status\",\"entryKey\":\"pending\",\"label\":\"Again\"}")
	if code != http.StatusConflict || errBody["error"] != "DICT_ENTRY_KEY_TAKEN" {
		t.Fatalf("duplicate entry = %d %v", code, errBody)
	}

	// q filter on entries
	code, q := getResource(t, env, "/api/data-dictionary/entries?q=pending")
	if code != http.StatusOK || q["total"] != float64(1) {
		t.Fatalf("entry q = %d %v", code, q)
	}
	// GOAL-015: dictKey exact filter (inner page) — a second type's entries
	// must not leak into the first type's filtered list.
	code, otherType := bearerJSON(t, env, admin, http.MethodPost, "/api/data-dictionary/types",
		"{\"key\":\"other_status\",\"name\":\"Other Status\"}")
	if code != http.StatusCreated {
		t.Fatalf("create other type = %d: %v", code, otherType)
	}
	code, otherEntry := bearerJSON(t, env, admin, http.MethodPost, "/api/data-dictionary/entries",
		"{\"dictKey\":\"other_status\",\"entryKey\":\"done\",\"label\":\"Done\"}")
	if code != http.StatusCreated {
		t.Fatalf("create other entry = %d: %v", code, otherEntry)
	}
	code, filtered := getResource(t, env, "/api/data-dictionary/entries?dictKey=order_status&pageSize=100")
	if code != http.StatusOK {
		t.Fatalf("dictKey filter = %d %v", code, filtered)
	}
	items, _ := filtered["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("dictKey filter items = %d, want 1 (only order_status)", len(items))
	}
	first := items[0].(map[string]any)
	if first["dictKey"] != "order_status" {
		t.Fatalf("filtered entry dictKey = %v, want order_status", first["dictKey"])
	}
	if first["dictTypeName"] != "Order status" {
		t.Fatalf("filtered entry dictTypeName = %v, want Order status", first["dictTypeName"])
	}
	code, otherFiltered := getResource(t, env, "/api/data-dictionary/entries?dictKey=other_status&pageSize=100")
	if code != http.StatusOK {
		t.Fatalf("other dictKey filter = %d %v", code, otherFiltered)
	}
	otherItems, _ := otherFiltered["items"].([]any)
	if len(otherItems) != 1 || otherItems[0].(map[string]any)["entryKey"] != "done" {
		t.Fatalf("other dictKey filter items = %v, want [done]", otherItems)
	}
	// dictKey composes with q.
	code, composed := getResource(t, env, "/api/data-dictionary/entries?dictKey=order_status&q=pending")
	if code != http.StatusOK || composed["total"] != float64(1) {
		t.Fatalf("dictKey+q = %d %v", code, composed)
	}
	code, composedMiss := getResource(t, env, "/api/data-dictionary/entries?dictKey=other_status&q=pending")
	if code != http.StatusOK || composedMiss["total"] != float64(0) {
		t.Fatalf("dictKey+q miss = %d %v", code, composedMiss)
	}

	// GOAL-015 F-002 (grok audit): after the LEFT JOIN dict_types the ORDER BY
	// columns must be qualified — dict_types also has sort/updated_at, so an
	// unqualified sort previously 500'd with "ambiguous column name". These
	// combos used to fail and now must succeed with stable ordering.
	code, sortedBySort := getResource(t, env, "/api/data-dictionary/entries?dictKey=order_status&sort=sort&order=asc&pageSize=100")
	if code != http.StatusOK {
		t.Fatalf("dictKey+sort=sort = %d %v", code, sortedBySort)
	}
	code, sortedByUpdated := getResource(t, env, "/api/data-dictionary/entries?sort=updatedAt&order=desc&pageSize=100")
	if code != http.StatusOK {
		t.Fatalf("sort=updatedAt = %d %v", code, sortedByUpdated)
	}
	// GOAL-015 F-003: the sortable dictTypeName column maps to dt.name.
	code, sortedByName := getResource(t, env, "/api/data-dictionary/entries?sort=dictTypeName&order=asc&pageSize=100")
	if code != http.StatusOK {
		t.Fatalf("sort=dictTypeName = %d %v", code, sortedByName)
	}
	nameItems, _ := sortedByName["items"].([]any)
	if len(nameItems) != 2 || nameItems[0].(map[string]any)["dictTypeName"] != "Order status" {
		t.Fatalf("sort=dictTypeName order = %v, want [Order status, Other Status]", nameItems)
	}
	// dictKey + sort + page composes (F-002 page slice path).
	code, paged := getResource(t, env, "/api/data-dictionary/entries?dictKey=order_status&sort=sort&page=1&pageSize=1")
	if code != http.StatusOK {
		t.Fatalf("dictKey+sort+page = %d %v", code, paged)
	}
	pagedItems, _ := paged["items"].([]any)
	if len(pagedItems) != 1 || paged["total"] != float64(1) {
		t.Fatalf("dictKey+sort+page items = %v total=%v", pagedItems, paged["total"])
	}

	// patch entry
	code, patchBody := bearerJSON(t, env, admin, http.MethodPatch, "/api/data-dictionary/entries/"+entryID,
		"{\"label\":\"Pending (updated)\",\"enabled\":false}")
	if code != http.StatusOK {
		t.Fatalf("patch entry = %d: %v", code, patchBody)
	}

	// description round-trip on create + patch (fixed: create dropped it).
	code, descBody := bearerJSON(t, env, admin, http.MethodPost, "/api/data-dictionary/types",
		"{\"key\":\"desc_type\",\"name\":\"Desc\",\"description\":\"hello world\"}")
	if code != http.StatusCreated {
		t.Fatalf("create type with description = %d: %v", code, descBody)
	}
	if got, _ := descBody["description"].(string); got != "hello world" {
		t.Fatalf("create description = %q, want hello world", got)
	}
	descID, _ := descBody["id"].(string)
	code, descBody = bearerJSON(t, env, admin, http.MethodPatch, "/api/data-dictionary/types/"+descID,
		"{\"description\":\"updated desc\"}")
	if code != http.StatusOK {
		t.Fatalf("patch type description = %d: %v", code, descBody)
	}
	if got, _ := descBody["description"].(string); got != "updated desc" {
		t.Fatalf("patched description = %q, want updated desc", got)
	}
	// clearing to empty is allowed (JSONFields optional semantics).
	code, descBody = bearerJSON(t, env, admin, http.MethodPatch, "/api/data-dictionary/types/"+descID,
		"{\"description\":\"\"}")
	if code != http.StatusOK {
		t.Fatalf("patch type description to empty = %d: %v", code, descBody)
	}
	if got, _ := descBody["description"].(string); got != "" {
		t.Fatalf("cleared description = %q, want empty", got)
	}

	// cascade: delete the type removes its entries
	code, _ = bearerJSON(t, env, admin, http.MethodDelete, "/api/data-dictionary/types/"+typeID, "")
	if code != http.StatusNoContent {
		t.Fatalf("delete type = %d", code)
	}
	code, _ = getResource(t, env, "/api/data-dictionary/entries?q=pending")
	if code != http.StatusOK {
		t.Fatalf("entries after cascade = %d", code)
	}
	if code == http.StatusOK {
		// GOAL-015: cascade removes order_status entries; the other_status
		// type's entry (added for the dictKey filter test) survives.
		code2, q2 := getResource(t, env, "/api/data-dictionary/entries?dictKey=order_status&pageSize=100")
		if code2 != http.StatusOK || q2["total"] != float64(0) {
			t.Fatalf("cascade left order_status entries: %d %v", code2, q2)
		}
		code3, q3 := getResource(t, env, "/api/data-dictionary/entries?dictKey=other_status&pageSize=100")
		if code3 != http.StatusOK || q3["total"] != float64(1) {
			t.Fatalf("other_status entries missing after cascade: %d %v", code3, q3)
		}
	}

	// audit events recorded
	ops, _, err := env.operations.ListOperationsFiltered(operationlog.OperationFilter{Q: "dictionary.", Sort: "created", Order: "asc", Page: 1, PageSize: 100})
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]bool{}
	for _, op := range ops {
		got[op.Event] = true
	}
	for _, want := range []string{"dictionary.create", "dictionary.update", "dictionary.delete"} {
		if !got[want] {
			t.Fatalf("missing audit event %s (got %v)", want, got)
		}
	}
}

// Permission gates fail closed: viewer without dictionary.read/write → 403;
// anonymous → 401.
func TestDictionaryPermissionGates(t *testing.T) {
	env := newAuthTestEnv(t)
	_ = adminToken(t, env)
	env.addUser(t, "dict-viewer", "pw", []string{"viewer"})
	viewer := env.login(t, "dict-viewer", "pw")

	code, _ := getResourceAs(t, env, viewer, "/api/data-dictionary/types")
	if code != http.StatusForbidden {
		t.Fatalf("viewer list types = %d, want 403", code)
	}
	code, _ = bearerJSON(t, env, viewer, http.MethodPost, "/api/data-dictionary/types",
		"{\"key\":\"x\",\"name\":\"X\"}")
	if code != http.StatusForbidden {
		t.Fatalf("viewer create type = %d, want 403", code)
	}
	code, _ = getResourceAs(t, env, viewer, "/api/data-dictionary/entries")
	if code != http.StatusForbidden {
		t.Fatalf("viewer list entries = %d, want 403", code)
	}

	// anonymous → 401
	for _, path := range []string{"/api/data-dictionary/types", "/api/data-dictionary/entries"} {
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, path, nil))
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("anonymous GET %s = %d, want 401", path, rr.Code)
		}
	}
}
