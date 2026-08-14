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

	// patch entry
	code, patchBody := bearerJSON(t, env, admin, http.MethodPatch, "/api/data-dictionary/entries/"+entryID,
		"{\"label\":\"Pending (updated)\",\"enabled\":false}")
	if code != http.StatusOK {
		t.Fatalf("patch entry = %d: %v", code, patchBody)
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
		code2, q2 := getResource(t, env, "/api/data-dictionary/entries?pageSize=100")
		if code2 != http.StatusOK || q2["total"] != float64(0) {
			t.Fatalf("cascade left entries: %d %v", code2, q2)
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
