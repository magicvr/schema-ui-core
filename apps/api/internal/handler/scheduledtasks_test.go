// Scheduled tasks surface tests (S-04 · GOAL-010 D-002 §6): CRUD lifecycle
// with cron validation, manual trigger + run history, permission gates, audit
// events.
package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

func taskCreate(t *testing.T, env *authTestEnv, token, key, cron, name string) string {
	t.Helper()
	code, body := bearerJSON(t, env, token, http.MethodPost, "/api/scheduled-tasks",
		"{\"key\":\""+key+"\",\"cron\":\""+cron+"\",\"name\":\""+name+"\",\"enabled\":true}")
	if code != http.StatusCreated {
		t.Fatalf("create task = %d: %v", code, body)
	}
	id, _ := body["id"].(string)
	if id == "" {
		t.Fatal("create task missing id")
	}
	return id
}

func TestScheduledTasksLifecycle(t *testing.T) {
	env := newAuthTestEnv(t)
	admin := adminToken(t, env)

	id := taskCreate(t, env, admin, "cleanup", "*/15 * * * *", "Cleanup")

	// Invalid cron → 400 INVALID_CRON.
	code, errBody := bearerJSON(t, env, admin, http.MethodPost, "/api/scheduled-tasks",
		"{\"key\":\"bad\",\"cron\":\"99 99 99 99 99\",\"name\":\"Bad\"}")
	if code != http.StatusBadRequest || errBody["error"] != "INVALID_CRON" {
		t.Fatalf("invalid cron = %d %v", code, errBody)
	}
	// Duplicate key → 409 TASK_KEY_TAKEN.
	code, errBody = bearerJSON(t, env, admin, http.MethodPost, "/api/scheduled-tasks",
		"{\"key\":\"cleanup\",\"cron\":\"* * * * *\",\"name\":\"Dup\"}")
	if code != http.StatusConflict || errBody["error"] != "TASK_KEY_TAKEN" {
		t.Fatalf("duplicate key = %d %v", code, errBody)
	}

	// List contains the task.
	code, list := getResource(t, env, "/api/scheduled-tasks?pageSize=100")
	if code != http.StatusOK || list["total"] != float64(1) {
		t.Fatalf("list = %d %v", code, list)
	}

	// Manual trigger records a run row.
	code, _ = bearerJSON(t, env, admin, http.MethodPost, "/api/scheduled-tasks/"+id+"/run", "")
	if code != http.StatusNoContent {
		t.Fatalf("run = %d", code)
	}
	code, runs := getResource(t, env, "/api/scheduled-tasks/"+id+"/runs")
	if code != http.StatusOK || runs["total"] != float64(1) {
		t.Fatalf("task runs = %d %v", code, runs)
	}
	code, allRuns := getResource(t, env, "/api/task-runs?pageSize=100")
	if code != http.StatusOK || allRuns["total"] != float64(1) {
		t.Fatalf("all runs = %d %v", code, allRuns)
	}

	// Patch (cron update + disable) keeps merge semantics.
	code, _ = bearerJSON(t, env, admin, http.MethodPatch, "/api/scheduled-tasks/"+id,
		"{\"name\":\"Cleanup renamed\",\"enabled\":false}")
	if code != http.StatusOK {
		t.Fatalf("patch = %d", code)
	}
	code, detail := getResource(t, env, "/api/scheduled-tasks/"+id)
	if code != http.StatusOK || detail["name"] != "Cleanup renamed" || detail["enabled"] != false {
		t.Fatalf("detail after patch = %d %v", code, detail)
	}

	// Delete cascades run rows.
	code, _ = bearerJSON(t, env, admin, http.MethodDelete, "/api/scheduled-tasks/"+id, "")
	if code != http.StatusNoContent {
		t.Fatalf("delete = %d", code)
	}
	code, _ = getResource(t, env, "/api/task-runs?pageSize=100")
	if code != http.StatusOK {
		t.Fatalf("runs after delete = %d", code)
	}

	// Audit events.
	ops, _, err := env.operations.ListOperationsFiltered(operationlog.OperationFilter{Q: "scheduled-tasks.", Sort: "created", Order: "asc", Page: 1, PageSize: 100})
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]bool{}
	for _, op := range ops {
		got[op.Event] = true
	}
	for _, want := range []string{"scheduled-tasks.create", "scheduled-tasks.update", "scheduled-tasks.delete"} {
		if !got[want] {
			t.Fatalf("missing audit event %s (got %v)", want, got)
		}
	}
}

func TestScheduledTasksPermissionGates(t *testing.T) {
	env := newAuthTestEnv(t)
	env.addUser(t, "task-viewer", "pw", []string{"viewer"})
	viewer := env.login(t, "task-viewer", "pw")

	code, _ := getResourceAs(t, env, viewer, "/api/scheduled-tasks")
	if code != http.StatusForbidden {
		t.Fatalf("viewer list = %d, want 403", code)
	}
	code, _ = bearerJSON(t, env, viewer, http.MethodPost, "/api/scheduled-tasks",
		"{\"key\":\"x\",\"cron\":\"* * * * *\",\"name\":\"X\"}")
	if code != http.StatusForbidden {
		t.Fatalf("viewer create = %d, want 403", code)
	}
	// Anonymous → 401.
	anon := httptest.NewRecorder()
	env.mux.ServeHTTP(anon, httptest.NewRequest(http.MethodGet, "/api/scheduled-tasks", nil))
	if anon.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous list = %d, want 401", anon.Code)
	}
}

// F-01 (W14 D-003): the handler directory endpoint lists registered keys.
func TestScheduledTasksHandlersDirectory(t *testing.T) {
	env := newAuthTestEnv(t)
	admin := adminToken(t, env)
	code, body := getResourceAs(t, env, admin, "/api/scheduled-tasks/handlers")
	if code != http.StatusOK {
		t.Fatalf("handlers = %d: %v", code, body)
	}
	items, ok := body["items"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("handlers items = %#v, want one system.noop", body["items"])
	}
	first, ok := items[0].(map[string]any)
	if !ok || first["key"] != "system.noop" {
		t.Fatalf("first handler = %#v, want system.noop", items[0])
	}
}

// Unknown handler keys are rejected at write time (A-003 F-003).
func TestScheduledTasksInvalidHandler(t *testing.T) {
	env := newAuthTestEnv(t)
	admin := adminToken(t, env)
	code, body := bearerJSON(t, env, admin, http.MethodPost, "/api/scheduled-tasks",
		"{\"key\":\"h\",\"cron\":\"* * * * *\",\"name\":\"H\",\"handler\":\"typo.handler\"}")
	if code != http.StatusBadRequest || body["error"] != "INVALID_HANDLER" {
		t.Fatalf("unknown handler = %d %v, want 400 INVALID_HANDLER", code, body)
	}
}
