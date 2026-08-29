// File library surface tests (S-02 · GOAL-007 D-002 §4/§7): list/detail
// read gates, download headers + audit, hard delete + audit, the upload
// confirmation endpoint (owner check, id validation), and the viewer/anonymous
// fail-closed paths.
package handler

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
	"path/filepath"
	"strings"
	"testing"
)

// uploadLibraryFile stores one file through the central upload endpoint and
// returns its id.
func uploadLibraryFile(t *testing.T, env *authTestEnv, token, filename string, content []byte) string {
	t.Helper()
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write(content)
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("upload status = %d: %s", rr.Code, rr.Body.String())
	}
	var out map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	id, _ := out["id"].(string)
	if id == "" {
		t.Fatal("upload response missing id")
	}
	return id
}

func bearerJSON(t *testing.T, env *authTestEnv, token, method, path, body string) (int, map[string]any) {
	t.Helper()
	req := bearer(t, token, method, path, body)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	var out map[string]any
	if rr.Body.Len() > 0 {
		if err := json.NewDecoder(rr.Body).Decode(&out); err != nil {
			t.Fatalf("decode %q: %v", rr.Body.String(), err)
		}
	}
	return rr.Code, out
}

// getResourceAs fetches a path with an explicit token (unlike getResource
// which always uses the seeded admin).
func getResourceAs(t *testing.T, env *authTestEnv, token, path string) (int, map[string]any) {
	t.Helper()
	return bearerJSON(t, env, token, http.MethodGet, path, "")
}

func fileEventRows(t *testing.T, env *authTestEnv) []struct{ Event, RecordID string } {
	t.Helper()
	// The operation-log Q search covers the event column; "files." selects
	// exactly the three file-library events.
	operations, _, err := env.operations.ListOperationsFiltered(operationlog.OperationFilter{Q: "files.", Sort: "created", Order: "asc", Page: 1, PageSize: 100})
	if err != nil {
		t.Fatal(err)
	}
	rows := make([]struct{ Event, RecordID string }, 0, len(operations))
	for _, op := range operations {
		recordID := ""
		if op.RecordID != nil {
			recordID = *op.RecordID
		}
		rows = append(rows, struct{ Event, RecordID string }{Event: op.Event, RecordID: recordID})
	}
	return rows
}

// The library surface serves the stored files with metadata, honors the
// schema-driven list contract, streams downloads with the frozen hardening
// headers, deletes hard, and records the three audit events.
func TestFileLibraryLifecycle(t *testing.T) {
	env := newAuthTestEnv(t)
	admin := adminToken(t, env)

	reportID := uploadLibraryFile(t, env, admin, "quarterly-report.pdf", []byte("report-bytes"))
	photoID := uploadLibraryFile(t, env, admin, "photo.png", []byte("photo-bytes"))

	// List: both files, metadata complete, name sort default.
	code, list := getResource(t, env, "/api/library/files?pageSize=100")
	if code != http.StatusOK {
		t.Fatalf("list status = %d: %v", code, list)
	}
	if list["total"] != float64(2) {
		t.Fatalf("list total = %v, want 2", list["total"])
	}
	items, _ := list["items"].([]any)
	if len(items) != 2 {
		t.Fatalf("list items = %d, want 2", len(items))
	}
	names := map[string]bool{}
	for _, raw := range items {
		item, _ := raw.(map[string]any)
		names[item["name"].(string)] = true
		if item["owner"] != "user-admin" {
			t.Fatalf("row owner = %v, want user-admin", item["owner"])
		}
		if _, ok := item["size"].(float64); !ok {
			t.Fatalf("row size missing: %v", item)
		}
	}
	if !names["quarterly-report.pdf"] || !names["photo.png"] {
		t.Fatalf("list names = %v", names)
	}

	// Q search narrows to one row.
	code, q := getResource(t, env, "/api/library/files?q=photo")
	if code != http.StatusOK || q["total"] != float64(1) {
		t.Fatalf("q search = %d %v, want 1 row", code, q)
	}

	// Detail row.
	code, detail := getResource(t, env, "/api/library/files/"+photoID)
	if code != http.StatusOK || detail["name"] != "photo.png" {
		t.Fatalf("detail = %d %v", code, detail)
	}

	// Download: exact bytes + frozen hardening headers.
	req := bearer(t, admin, http.MethodGet, "/api/library/files/"+reportID+"/download", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK || rr.Body.String() != "report-bytes" {
		t.Fatalf("download = %d %q", rr.Code, rr.Body.String())
	}
	if rr.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("download missing nosniff")
	}
	if !strings.Contains(rr.Header().Get("Content-Disposition"), "attachment") {
		t.Fatalf("download must force attachment: %q", rr.Header().Get("Content-Disposition"))
	}

	// Upload confirmation: bare id and url form both accepted.
	code, ack := bearerJSON(t, env, admin, http.MethodPost, "/api/library/files/upload", "{\"file\":\""+reportID+"\"}")
	if code != http.StatusOK || ack["id"] != reportID {
		t.Fatalf("upload ack = %d %v", code, ack)
	}
	code, ack = bearerJSON(t, env, admin, http.MethodPost, "/api/library/files/upload", "{\"file\":\"/api/files/"+photoID+"\"}")
	if code != http.StatusOK {
		t.Fatalf("upload ack url form = %d", code)
	}

	// Hard delete + repeat 404.
	code, _ = bearerJSON(t, env, admin, http.MethodDelete, "/api/library/files/"+reportID, "")
	if code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want 204", code)
	}
	if _, err := os.Stat(filepath.Join(env.uploadDir, reportID)); !os.IsNotExist(err) {
		t.Fatalf("stored object still present after delete: %v", err)
	}
	code, delErr := bearerJSON(t, env, admin, http.MethodDelete, "/api/library/files/"+reportID, "")
	if code != http.StatusNotFound || delErr["error"] != "FILE_NOT_FOUND" {
		t.Fatalf("repeat delete = %d %v, want 404 FILE_NOT_FOUND", code, delErr)
	}

	// The three audit events are recorded (upload ack, download, delete).
	events := fileEventRows(t, env)
	got := map[string]bool{}
	for _, e := range events {
		got[e.Event] = true
		if e.Event == "files.upload" && e.RecordID != reportID && e.RecordID != photoID {
			t.Fatalf("upload event record_id = %q", e.RecordID)
		}
	}
	for _, want := range []string{"files.upload", "files.download", "files.delete"} {
		if !got[want] {
			t.Fatalf("missing audit event %s (got %v)", want, got)
		}
	}

	_ = os.RemoveAll(env.uploadDir)
}

// Permission gates fail closed: viewer (files.read/files.delete absent) → 403
// on every library surface; anonymous → 401.
func TestFileLibraryPermissionGates(t *testing.T) {
	env := newAuthTestEnv(t)
	admin := adminToken(t, env)
	env.addUser(t, "file-viewer", "pw", []string{"viewer"})
	viewer := env.login(t, "file-viewer", "pw")

	id := uploadLibraryFile(t, env, admin, "secret.txt", []byte("s"))

	for _, path := range []string{
		"/api/library/files",
		"/api/library/files/" + id,
		"/api/library/files/" + id + "/download",
	} {
		code, _ := getResourceAs(t, env, viewer, path)
		if code != http.StatusForbidden {
			t.Fatalf("viewer GET %s = %d, want 403", path, code)
		}
	}
	code, _ := bearerJSON(t, env, viewer, http.MethodDelete, "/api/library/files/"+id, "")
	if code != http.StatusForbidden {
		t.Fatalf("viewer DELETE = %d, want 403", code)
	}
	code, _ = bearerJSON(t, env, viewer, http.MethodPost, "/api/library/files/upload", "{\"file\":\""+id+"\"}")
	if code != http.StatusForbidden {
		t.Fatalf("viewer upload ack = %d, want 403", code)
	}

	// Anonymous → 401 on every surface.
	for _, req := range []*http.Request{
		httptest.NewRequest(http.MethodGet, "/api/library/files", nil),
		httptest.NewRequest(http.MethodDelete, "/api/library/files/"+id, nil),
		httptest.NewRequest(http.MethodGet, "/api/library/files/"+id+"/download", nil),
		httptest.NewRequest(http.MethodPost, "/api/library/files/upload", strings.NewReader("{}")),
	} {
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("anonymous %s %s = %d, want 401", req.Method, req.URL.Path, rr.Code)
		}
	}

	_ = os.RemoveAll(env.uploadDir)
}

// Upload confirmation validates the id shape and ownership: foreign files are
// never confirmable, unknown ids 404, malformed ids 400.
func TestFileLibraryUploadAckValidation(t *testing.T) {
	env := newAuthTestEnv(t)
	admin := adminToken(t, env)
	env.addUser(t, "second-admin", "pw", []string{"admin"})
	second := env.login(t, "second-admin", "pw")

	adminID := uploadLibraryFile(t, env, admin, "mine.txt", []byte("a"))
	foreignID := uploadLibraryFile(t, env, second, "theirs.txt", []byte("b"))

	// Malformed id → INVALID_FILE_ID.
	code, errBody := bearerJSON(t, env, admin, http.MethodPost, "/api/library/files/upload", "{\"file\":\"../etc/passwd\"}")
	if code != http.StatusBadRequest || errBody["error"] != "INVALID_FILE_ID" {
		t.Fatalf("malformed id = %d %v", code, errBody)
	}
	// Unknown well-formed id → FILE_NOT_FOUND.
	code, errBody = bearerJSON(t, env, admin, http.MethodPost, "/api/library/files/upload", "{\"file\":\"00000000000000000000000000000000\"}")
	if code != http.StatusNotFound || errBody["error"] != "FILE_NOT_FOUND" {
		t.Fatalf("unknown id = %d %v", code, errBody)
	}
	// Foreign owner → 403 (never confirmable into the library).
	code, errBody = bearerJSON(t, env, admin, http.MethodPost, "/api/library/files/upload", "{\"file\":\""+foreignID+"\"}")
	if code != http.StatusForbidden || errBody["error"] != "FORBIDDEN" {
		t.Fatalf("foreign ack = %d %v", code, errBody)
	}
	// Owner's own file confirms fine.
	code, _ = bearerJSON(t, env, admin, http.MethodPost, "/api/library/files/upload", "{\"file\":\""+adminID+"\"}")
	if code != http.StatusOK {
		t.Fatalf("own ack = %d", code)
	}

	_ = os.RemoveAll(env.uploadDir)
}

// A-002 R-002: a ghost sidecar (body-less .meta.json leftover from the
// pre-port era) must NOT delete with 204 — the response is 404 and the
// best-effort cleanup still removes the sidecar so it cannot ghost-list.
func TestFileLibraryGhostSidecarDelete(t *testing.T) {
	env := newAuthTestEnv(t)
	admin := adminToken(t, env)

	id := "cccccccccccccccccccccccccccccccc"
	sidecar := `{"name":"ghost.csv","type":"text/csv","owner":"` + admin + `"}`

	// Simulate pre-port residue: only the meta sidecar, no body.
	if err := os.MkdirAll(env.uploadDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Resolve the real admin user id first via a normal upload ack flow is
	// overkill; owner does not matter for the ghost path — write it directly.
	_ = sidecar
	if err := os.WriteFile(filepath.Join(env.uploadDir, id+".meta.json"), []byte(`{"name":"ghost.csv","type":"text/csv","owner":"someone"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	code, _ := bearerJSON(t, env, admin, http.MethodDelete, "/api/library/files/"+id, "")
	if code != http.StatusNotFound {
		t.Fatalf("ghost sidecar delete = %d, want 404 (GOAL-004 D-001 section 5)", code)
	}
	if _, err := os.Stat(filepath.Join(env.uploadDir, id+".meta.json")); !os.IsNotExist(err) {
		t.Fatalf("best-effort cleanup did not remove the ghost sidecar: %v", err)
	}
}
