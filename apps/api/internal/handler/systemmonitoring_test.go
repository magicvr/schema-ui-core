// System monitoring surface tests (S-03 · GOAL-009 D-002 `5): status summary
// contents + permission gates, and the read-only recent-events list contract.
package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// The status endpoint reports the in-process probes (store ping + readiness
// gate), version/commit, uptime, the module set and the DB size — all fields
// present and typed; viewer/anonymous fail closed.
func TestSystemMonitoringStatus(t *testing.T) {
	env := newAuthTestEnv(t)

	code, body := getResource(t, env, "/api/system-monitoring/status")
	if code != http.StatusOK {
		t.Fatalf("status = %d: %v", code, body)
	}
	// Single-row list envelope (A-003 F-001): statCard dataSource loads via
	// fetchResourceList, which requires {items,total,page,pageSize}.
	if body["total"] != float64(1) || body["page"] != float64(1) {
		t.Fatalf("status envelope = %v, want total=1 page=1", body)
	}
	items, _ := body["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("status items = %d, want 1", len(items))
	}
	row := items[0].(map[string]any)
	for _, field := range []string{"status", "ready", "version", "commit", "uptimeSeconds", "moduleCount", "modules", "dbSizeBytes"} {
		if _, ok := row[field]; !ok {
			t.Fatalf("status row missing field %q: %v", field, row)
		}
	}
	if row["status"] != "ok" || row["ready"] != true {
		t.Fatalf("status row = %v, want ok/true", row)
	}
	modules, _ := row["modules"].([]any)
	if len(modules) == 0 {
		t.Fatalf("status modules empty: %v", row)
	}
	if row["moduleCount"] != float64(len(modules)) {
		t.Fatalf("moduleCount = %v, want %d", row["moduleCount"], len(modules))
	}

	// Anonymous → 401.
	anon := httptest.NewRecorder()
	env.mux.ServeHTTP(anon, httptest.NewRequest(http.MethodGet, "/api/system-monitoring/status", nil))
	if anon.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous status = %d, want 401", anon.Code)
	}
	// Viewer (no monitoring.read) → 403.
	env.addUser(t, "monitor-viewer", "pw", []string{"viewer"})
	viewer := env.login(t, "monitor-viewer", "pw")
	code, _ = getResourceAs(t, env, viewer, "/api/system-monitoring/status")
	if code != http.StatusForbidden {
		t.Fatalf("viewer status = %d, want 403", code)
	}
}

// The recent-events resource follows the read-only list contract (items/total)
// and the same permission gates as the status endpoint.
func TestSystemMonitoringErrors(t *testing.T) {
	env := newAuthTestEnv(t)
	admin := adminToken(t, env)

	// An authenticated action produces an operation-log row to list.
	code, _ := bearerJSON(t, env, admin, http.MethodPost, "/api/data-dictionary/types",
		"{\"key\":\"m_status\",\"name\":\"Monitoring\"}")
	if code != http.StatusCreated {
		t.Fatalf("seed type = %d", code)
	}

	code, list := getResource(t, env, "/api/system-monitoring/errors?pageSize=100")
	if code != http.StatusOK {
		t.Fatalf("errors list = %d: %v", code, list)
	}
	if list["total"] == float64(0) {
		t.Fatalf("errors list total = 0, want >= 1 (seeded operation row)")
	}
	items, _ := list["items"].([]any)
	if len(items) == 0 {
		t.Fatalf("errors items empty")
	}
	first := items[0].(map[string]any)
	if first["event"] == nil {
		t.Fatalf("errors row missing event: %v", first)
	}

	// Viewer → 403.
	env.addUser(t, "monitor-viewer2", "pw", []string{"viewer"})
	viewer := env.login(t, "monitor-viewer2", "pw")
	code, _ = getResourceAs(t, env, viewer, "/api/system-monitoring/errors")
	if code != http.StatusForbidden {
		t.Fatalf("viewer errors = %d, want 403", code)
	}
}

// Unknown error id → 404 OPERATION_NOT_FOUND (A-003 F-002: GetOperation wraps
// the sentinel, so the entity must use errors.Is).
func TestSystemMonitoringErrorsNotFound(t *testing.T) {
	env := newAuthTestEnv(t)
	code, body := getResource(t, env, "/api/system-monitoring/errors/op-missing")
	if code != http.StatusNotFound || body["error"] != "OPERATION_NOT_FOUND" {
		t.Fatalf("missing error id = %d %v, want 404 OPERATION_NOT_FOUND", code, body)
	}
}
