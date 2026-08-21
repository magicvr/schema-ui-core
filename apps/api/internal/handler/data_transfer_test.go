// F-02 export/import tests (GOAL-004 S3): CSV shape/escaping/cap/permissions/
// audit for export; validation/partial-apply/error-report/owner/audit for import.
package handler

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

func uploadCSV(t *testing.T, env *authTestEnv, token, filename, content string) string {
	t.Helper()
	var body strings.Builder
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/upload", strings.NewReader(body.String()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("upload = %d: %s", rr.Code, rr.Body.String())
	}
	var out struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	return out.ID
}

func TestExportUsersCSV(t *testing.T) {
	env := newAuthTestEnv(t)
	env.addUser(t, "bob", "bob-password", []string{"editor"})
	req := bearer(t, adminToken(t, env), http.MethodGet, "/api/export/users", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("export = %d: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	if !strings.HasPrefix(body, "\uFEFF") {
		t.Fatal("export missing UTF-8 BOM")
	}
	lines := strings.Split(strings.TrimPrefix(body, "\uFEFF"), "\n")
	if len(lines) < 2 {
		t.Fatalf("export lines = %d, want >= 2", len(lines))
	}
	header := lines[0]
	for _, col := range []string{"id", "username", "name", "roles", "enabled", "locked", "createdAt", "updatedAt"} {
		if !strings.Contains(header, col) {
			t.Fatalf("export header %q missing %s", header, col)
		}
	}
	joined := strings.Join(lines, "\n")
	if !strings.Contains(joined, testSeedUsername) || !strings.Contains(joined, "bob") {
		t.Fatalf("export rows missing users: %s", joined)
	}
	ct := rr.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "text/csv") {
		t.Fatalf("export content-type = %q, want text/csv", ct)
	}
	if rr.Header().Get("Content-Disposition") == "" {
		t.Fatal("export missing Content-Disposition")
	}
}

func TestExportRolesCSV(t *testing.T) {
	env := newAuthTestEnv(t)
	req := bearer(t, adminToken(t, env), http.MethodGet, "/api/export/roles", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("export roles = %d: %s", rr.Code, rr.Body.String())
	}
	body := strings.TrimPrefix(rr.Body.String(), "\uFEFF")
	if !strings.Contains(body, "admin") || !strings.Contains(body, "editor") {
		t.Fatalf("roles export missing system roles: %s", body)
	}
}

func TestExportUnknownResource404(t *testing.T) {
	env := newAuthTestEnv(t)
	req := bearer(t, adminToken(t, env), http.MethodGet, "/api/export/orders", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("unknown export = %d, want 404", rr.Code)
	}
}

func TestExportPermissionGated(t *testing.T) {
	env := newAuthTestEnv(t)
	env.addUser(t, "viewer1", "viewer-password", []string{"viewer"})
	token := env.login(t, "viewer1", "viewer-password")
	req := bearer(t, token, http.MethodGet, "/api/export/users", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("viewer export = %d, want 403", rr.Code)
	}
}

func TestImportUsersPartialApply(t *testing.T) {
	env := newAuthTestEnv(t)
	csv := "username,name,roles,password\nimp1,Import One,editor,import-pass-1\nimp2,Import Two,,import-pass-2\ndup,Duplicate,," // dup missing password → row error
	fileID := uploadCSV(t, env, adminToken(t, env), "users.csv", csv)
	req := bearer(t, adminToken(t, env), http.MethodPost, "/api/import/users", `{"fileId":"`+fileID+`"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("import = %d: %s", rr.Code, rr.Body.String())
	}
	var result struct {
		Applied int              `json:"applied"`
		Failed  int              `json:"failed"`
		Total   int              `json:"total"`
		Errors  []importRowError `json:"errors"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result.Applied != 2 || result.Failed != 1 || result.Total != 3 {
		t.Fatalf("import result = %+v, want applied=2 failed=1 total=3", result)
	}
	if len(result.Errors) != 1 || result.Errors[0].Row != 4 || !strings.Contains(result.Errors[0].Message, "password") {
		t.Fatalf("import errors = %+v", result.Errors)
	}
	if code, out := getResource(t, env, "/api/users?q=imp1"); code != http.StatusOK || out["total"].(float64) != 1 {
		t.Fatalf("imp1 not applied: %v", out)
	}
	ops, _, err := env.operations.ListOperationsFiltered(operationlog.OperationFilter{Sort: "createdAt", Order: "desc", Page: 1, PageSize: 20})
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, op := range ops {
		if op.Event == operationlog.EventDataImport && op.Detail != nil && strings.Contains(*op.Detail, `"applied":2`) {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("data.import operation log entry missing")
	}
}

func TestImportUsersValidationErrors(t *testing.T) {
	env := newAuthTestEnv(t)
	env.addUser(t, "existing", "existing-pass", []string{"editor"})
	csv := "username,name,roles,password\nexisting,Duplicate Username,,import-pass-1\nx1,NoSuchRole,nope-role,import-pass-1\nx2,ShortPw,,short\n"
	fileID := uploadCSV(t, env, adminToken(t, env), "users.csv", csv)
	req := bearer(t, adminToken(t, env), http.MethodPost, "/api/import/users", `{"fileId":"`+fileID+`"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("import = %d: %s", rr.Code, rr.Body.String())
	}
	var result struct {
		Applied int              `json:"applied"`
		Failed  int              `json:"failed"`
		Errors  []importRowError `json:"errors"`
	}
	_ = json.NewDecoder(rr.Body).Decode(&result)
	if result.Applied != 0 || result.Failed != 3 {
		t.Fatalf("import result = %+v, want applied=0 failed=3", result)
	}
	messages := ""
	for _, e := range result.Errors {
		messages += e.Message + " | "
	}
	for _, want := range []string{"username already exists", "unknown role key", "8 to 72"} {
		if !strings.Contains(messages, want) {
			t.Fatalf("import errors missing %q: %s", want, messages)
		}
	}
}

func TestImportForeignFileForbidden(t *testing.T) {
	env := newAuthTestEnv(t)
	// files.write is admin-only, so a foreign-file import needs a second admin:
	// admin2 uploads, admin1 tries to import admin2's file → 403 owner mismatch.
	env.addUser(t, "admin2", "admin2-password", []string{"admin"})
	admin2Token := env.login(t, "admin2", "admin2-password")
	csv := "username,name,roles,password\nimpX,Import X,,import-pass-1\n"
	fileID := uploadCSV(t, env, admin2Token, "users.csv", csv)
	req := bearer(t, adminToken(t, env), http.MethodPost, "/api/import/users", `{"fileId":"`+fileID+`"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("foreign file import = %d, want 403: %s", rr.Code, rr.Body.String())
	}
}

func TestImportMissingFile404(t *testing.T) {
	env := newAuthTestEnv(t)
	req := bearer(t, adminToken(t, env), http.MethodPost, "/api/import/users", `{"fileId":"00000000000000000000000000000000"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("missing file import = %d, want 404", rr.Code)
	}
}

func TestImportPermissionGated(t *testing.T) {
	env := newAuthTestEnv(t)
	env.addUser(t, "editor1", "editor-password", []string{"editor"})
	token := env.login(t, "editor1", "editor-password")
	req := bearer(t, token, http.MethodPost, "/api/import/users", `{"fileId":"00000000000000000000000000000000"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("editor import = %d, want 403", rr.Code)
	}
}


// --- A-003 F-005 补充用例 ---

func TestTransferEndpointsAnonymous401(t *testing.T) {
	env := newAuthTestEnv(t)
	for _, tc := range []struct{ method, path string }{
		{http.MethodGet, "/api/export/users"},
		{http.MethodGet, "/api/export/roles"},
		{http.MethodPost, "/api/import/users"},
	} {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("%s %s anonymous = %d, want 401", tc.method, tc.path, rr.Code)
		}
	}
}

func TestEditorCanExport(t *testing.T) {
	env := newAuthTestEnv(t)
	env.addUser(t, "editor1", "editor-password", []string{"editor"})
	token := env.login(t, "editor1", "editor-password")
	req := bearer(t, token, http.MethodGet, "/api/export/users", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("editor export = %d, want 200 (data.export is PolicyAdminEditor)", rr.Code)
	}
}

func TestExportAuditLogged(t *testing.T) {
	env := newAuthTestEnv(t)
	req := bearer(t, adminToken(t, env), http.MethodGet, "/api/export/users", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("export = %d", rr.Code)
	}
	ops, _, err := env.operations.ListOperationsFiltered(operationlog.OperationFilter{Sort: "createdAt", Order: "desc", Page: 1, PageSize: 20})
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, op := range ops {
		if op.Event == operationlog.EventDataExport && op.Detail != nil && strings.Contains(*op.Detail, `"rows":`) {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("data.export operation log entry missing")
	}
}

func TestExportCSVEscaping(t *testing.T) {
	env := newAuthTestEnv(t)
	env.addUser(t, "comma,user", "comma-password", []string{"editor"})
	req := bearer(t, adminToken(t, env), http.MethodGet, "/api/export/users?q=comma", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("export = %d", rr.Code)
	}
	body := strings.TrimPrefix(rr.Body.String(), "\uFEFF")
	if !strings.Contains(body, "\"comma,user\"") {
		t.Fatalf("comma field not quoted: %s", body)
	}
}

func TestExportFormulaInjectionGuarded(t *testing.T) {
	env := newAuthTestEnv(t)
	env.addUser(t, "formula1", "formula-password", []string{"editor"})
	// rename via repository to a formula-like name
	formulaName := "=HYPERLINK(\"http://evil\")"
	now := time.Now().UTC()
	if _, err := env.authRepository.UpdateUser("user-formula1", authsession.UserPatch{Name: &formulaName}, "user-admin", now); err != nil {
		t.Fatalf("rename: %v", err)
	}
	req := bearer(t, adminToken(t, env), http.MethodGet, "/api/export/users?q=formula", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("export = %d", rr.Code)
	}
	body := strings.TrimPrefix(rr.Body.String(), "\uFEFF")
	if !strings.Contains(body, "'=HYPERLINK") {
		t.Fatalf("formula cell not neutralized: %s", body)
	}
}

func TestImportSizeLimit(t *testing.T) {
	env := newAuthTestEnv(t)
	// 3 MiB body (upload allows up to 8 MiB; import rejects > 2 MiB)
	big := strings.Repeat("a", 3<<20)
	csv := "username,name,roles,password\n" + big
	fileID := uploadCSV(t, env, adminToken(t, env), "big.csv", csv)
	req := bearer(t, adminToken(t, env), http.MethodPost, "/api/import/users", `{"fileId":"`+fileID+`"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("oversize import = %d, want 413", rr.Code)
	}
}

func TestImportMissingHeaderInvalidCsv(t *testing.T) {
	env := newAuthTestEnv(t)
	fileID := uploadCSV(t, env, adminToken(t, env), "bad.csv", "\n")
	req := bearer(t, adminToken(t, env), http.MethodPost, "/api/import/users", `{"fileId":"`+fileID+`"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("missing header = %d, want 400 INVALID_CSV", rr.Code)
	}
}

func TestExportLimitExceeded(t *testing.T) {
	env := newAuthTestEnv(t)
	req := bearer(t, adminToken(t, env), http.MethodGet, "/api/export/users?pageSize=20000", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("oversize export = %d, want 400", rr.Code)
	}
}

func TestImportRoleAssignmentBoundary(t *testing.T) {
	env := newAuthTestEnv(t)
	now := time.Now().UTC()
	// A delegated data-import holder WITHOUT admin / roles.assign: importing
	// roles=admin must fail per-row (F-001), not 403 the whole request.
	// The import modal's upload control requires files.write (upload is a
	// separate shared capability); a delegated importer needs both keys.
	if _, err := env.authRepository.CreateRoleWithGrants(
		"import-only", "Import only",
		[]string{"data.import", "files.write"}, nil, now,
	); err != nil {
		t.Fatalf("create import-only role: %v", err)
	}
	env.addUser(t, "importer", "importer-password", []string{"import-only"})
	token := env.login(t, "importer", "importer-password")
	csv := "username,name,roles,password\nimpA,Import A,admin,import-pass-1\nimpB,Import B,,import-pass-1\n"
	fileID := uploadCSV(t, env, token, "users.csv", csv)
	req := bearer(t, token, http.MethodPost, "/api/import/users", `{"fileId":"`+fileID+`"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("import = %d, want 200 (per-row failures): %s", rr.Code, rr.Body.String())
	}
	var result struct {
		Applied int              `json:"applied"`
		Failed  int              `json:"failed"`
		Errors  []importRowError `json:"errors"`
	}
	_ = json.NewDecoder(rr.Body).Decode(&result)
	if result.Applied != 1 || result.Failed != 1 {
		t.Fatalf("import result = %+v, want applied=1 failed=1", result)
	}
	if !strings.Contains(result.Errors[0].Message, "roles.assign") && !strings.Contains(result.Errors[0].Message, "admin") {
		t.Fatalf("delegation error message unexpected: %s", result.Errors[0].Message)
	}
}

func TestImportUnknownResource404(t *testing.T) {
	env := newAuthTestEnv(t)
	req := bearer(t, adminToken(t, env), http.MethodPost, "/api/import/orders", `{"fileId":"00000000000000000000000000000000"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("unknown import = %d, want 404", rr.Code)
	}
}