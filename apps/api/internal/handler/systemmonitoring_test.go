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
	for _, field := range []string{"status", "ready", "version", "commit", "uptimeSeconds", "modules", "dbSizeBytes"} {
		if _, ok := body[field]; !ok {
			t.Fatalf("status missing field %q: %v", field, body)
		}
	}
	if body["status"] != "ok" || body["ready"] != true {
		t.Fatalf("status = %v, want ok/true", body)
	}
	modules, _ := body["modules"].([]any)
	if len(modules) == 0 {
		t.Fatalf("status modules empty: %v", body)
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
