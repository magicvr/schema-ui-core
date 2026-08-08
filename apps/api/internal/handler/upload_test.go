package handler

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// I-PROTO-FULL-001 · D-UPLOAD server-side contract (07-actions-contract.md §7.2):
// multipart upload → {url,id,name,size}; GET /api/files/{id} serves bytes.
func TestUploadEndpointContract(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "contract.pdf")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write([]byte("pdf-bytes"))
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp := httptest.NewRecorder()
	env.mux.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("upload status = %d: %s", resp.Code, resp.Body.String())
	}
	var uploaded map[string]any
	if err := json.Unmarshal(resp.Body.Bytes(), &uploaded); err != nil {
		t.Fatal(err)
	}
	id, _ := uploaded["id"].(string)
	url, _ := uploaded["url"].(string)
	if id == "" || url == "" {
		t.Fatalf("upload response missing id/url: %v", uploaded)
	}
	if uploaded["name"] != "contract.pdf" || uploaded["size"] != float64(9) {
		t.Fatalf("upload meta = %v", uploaded)
	}

	// The URL is fetchable and returns the exact bytes.
	fileReq := httptest.NewRequest(http.MethodGet, url, nil)
	fileReq.Header.Set("Authorization", "Bearer "+token)
	fileResp := httptest.NewRecorder()
	env.mux.ServeHTTP(fileResp, fileReq)
	if fileResp.Code != http.StatusOK {
		t.Fatalf("file status = %d: %s", fileResp.Code, fileResp.Body.String())
	}
	if fileResp.Body.String() != "pdf-bytes" {
		t.Fatalf("file body = %q, want pdf-bytes", fileResp.Body.String())
	}

	// Missing file part → INVALID_UPLOAD; empty file → INVALID_FILE.
	empty := new(bytes.Buffer)
	writer2 := multipart.NewWriter(empty)
	part2, err := writer2.CreateFormFile("file", "empty.txt")
	if err != nil {
		t.Fatal(err)
	}
	_ = part2
	if err := writer2.Close(); err != nil {
		t.Fatal(err)
	}
	req2 := httptest.NewRequest(http.MethodPost, "/api/upload", empty)
	req2.Header.Set("Authorization", "Bearer "+token)
	req2.Header.Set("Content-Type", writer2.FormDataContentType())
	resp2 := httptest.NewRecorder()
	env.mux.ServeHTTP(resp2, req2)
	if resp2.Code != http.StatusBadRequest {
		t.Fatalf("empty upload status = %d: %s", resp2.Code, resp2.Body.String())
	}
	var apiError map[string]string
	if err := json.Unmarshal(resp2.Body.Bytes(), &apiError); err != nil {
		t.Fatal(err)
	}
	if apiError["error"] != "INVALID_FILE" {
		t.Fatalf("empty upload code = %q, want INVALID_FILE", apiError["error"])
	}

	// Anonymous upload → 401.
	anon := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	anon.Header.Set("Content-Type", writer.FormDataContentType())
	anonResp := httptest.NewRecorder()
	env.mux.ServeHTTP(anonResp, anon)
	if anonResp.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous upload status = %d, want 401", anonResp.Code)
	}

	// Unknown file id → FILE_NOT_FOUND.
	missing := httptest.NewRequest(http.MethodGet, "/api/files/00000000000000000000000000000000", nil)
	missing.Header.Set("Authorization", "Bearer "+token)
	missingResp := httptest.NewRecorder()
	env.mux.ServeHTTP(missingResp, missing)
	if missingResp.Code != http.StatusNotFound {
		t.Fatalf("missing file status = %d, want 404", missingResp.Code)
	}

	// Cleanup the test upload store (avoid leaking files into the repo).
	_ = os.RemoveAll(env.uploadDir)
}
