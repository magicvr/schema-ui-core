package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// W16-F03: the import template route returns a downloadable users CSV header.
func TestImportTemplateUsers(t *testing.T) {
	env := newAuthTestEnv(t)
	req := bearer(t, adminToken(t, env), http.MethodGet, "/api/import/users/template", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("template status = %d, want 200", rr.Code)
	}
	if !strings.HasPrefix(rr.Header().Get("Content-Type"), "text/csv") {
		t.Fatalf("template content-type = %q, want text/csv", rr.Header().Get("Content-Type"))
	}
	if !strings.Contains(rr.Body.String(), "username,name,roles,password") {
		t.Fatalf("template body = %q, want header row", rr.Body.String())
	}
}

// W16-F03: import failures include rowNumber/field/reason in fieldErrors while
// keeping the legacy errors list.
func TestImportFieldErrors(t *testing.T) {
	env := newAuthTestEnv(t)
	// Row 1 is valid; row 2 misses password; row 3 has a bad role reference.
	csv := "username,name,roles,password\nok,OK,editor,import-pass-1\nbad,Missing,,short\nrole,R,no-such-role,import-pass-1\n"
	fileID := uploadCSV(t, env, adminToken(t, env), "users.csv", csv)
	req := bearer(t, adminToken(t, env), http.MethodPost, "/api/import/users", `{"fileId":"`+fileID+`"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("import status = %d: %s", rr.Code, rr.Body.String())
	}
	var body struct {
		Applied     int `json:"applied"`
		Failed      int `json:"failed"`
		FieldErrors []struct {
			RowNumber int    `json:"rowNumber"`
			Field     string `json:"field"`
			Reason    string `json:"reason"`
		} `json:"fieldErrors"`
		Errors []struct {
			Row int `json:"row"`
		} `json:"errors"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("decode import response: %v", err)
	}
	if body.Applied != 1 || body.Failed != 2 {
		t.Fatalf("applied/failed = %d/%d, want 1/2", body.Applied, body.Failed)
	}
	if len(body.FieldErrors) != 2 {
		t.Fatalf("fieldErrors = %d, want 2: %+v", len(body.FieldErrors), body.FieldErrors)
	}
	if len(body.Errors) != 2 {
		t.Fatalf("legacy errors = %d, want 2", len(body.Errors))
	}
	seen := map[string]bool{}
	for _, fe := range body.FieldErrors {
		seen[fmt.Sprintf("%d:%s", fe.RowNumber, fe.Field)] = true
	}
	if !seen["3:password"] {
		t.Fatalf("fieldErrors missing row 3 password: %+v", body.FieldErrors)
	}
	if !seen["4:roles"] {
		t.Fatalf("fieldErrors missing row 4 roles: %+v", body.FieldErrors)
	}
}

// W16-F02: file library rows expose downloadUrl for preview/copy actions.
func TestFileLibraryRowsExposeDownloadUrl(t *testing.T) {
	env := newAuthTestEnv(t)
	fileID := uploadCSV(t, env, adminToken(t, env), "hello.csv", "a,b\n1,2\n")
	req := bearer(t, adminToken(t, env), http.MethodGet, "/api/library/files", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("list status = %d: %s", rr.Code, rr.Body.String())
	}
	var body struct {
		Items []map[string]any `json:"items"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, item := range body.Items {
		if item["id"] == fileID {
			url, _ := item["downloadUrl"].(string)
			if url != "/api/library/files/"+fileID+"/download" {
				t.Fatalf("downloadUrl = %q, want /api/library/files/%s/download", url, fileID)
			}
			found = true
		}
	}
	if !found {
		t.Fatalf("uploaded file %s not in list", fileID)
	}
}
