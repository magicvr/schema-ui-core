package handler

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
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

// C1 hardening: HTML-bearing content is rejected server-side (never stored),
// and every download is forced to attachment so stored bytes can never render
// inline in the API's same-origin context.
func TestUploadRejectsHtmlAndForcesAttachment(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)

	uploadPart := func(name, declaredType, content string) *httptest.ResponseRecorder {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", `form-data; name="file"; filename="`+name+`"`)
		h.Set("Content-Type", declaredType)
		part, err := writer.CreatePart(h)
		if err != nil {
			t.Fatal(err)
		}
		_, _ = part.Write([]byte(content))
		if err := writer.Close(); err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		return rr
	}

	// A part declared text/plain whose real content is HTML is detected as
	// text/html and rejected — the client declaration is never trusted.
	html := uploadPart("evil.txt", "text/plain", "<!DOCTYPE html><html><script>alert(1)</script></html>")
	if html.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("html upload status = %d, want 415", html.Code)
	}
	var apiError map[string]string
	if err := json.Unmarshal(html.Body.Bytes(), &apiError); err != nil {
		t.Fatal(err)
	}
	if apiError["error"] != "UNSUPPORTED_FILE_TYPE" {
		t.Fatalf("html upload code = %q, want UNSUPPORTED_FILE_TYPE", apiError["error"])
	}

	// A plain-text file still uploads, and its download is forced to attachment
	// with a sandbox CSP — never inline-renderable from the same origin.
	plain := uploadPart("notes.txt", "text/plain", "just some text")
	if plain.Code != http.StatusOK {
		t.Fatalf("plain upload status = %d: %s", plain.Code, plain.Body.String())
	}
	var uploaded map[string]any
	if err := json.Unmarshal(plain.Body.Bytes(), &uploaded); err != nil {
		t.Fatal(err)
	}
	url, _ := uploaded["url"].(string)
	fileReq := httptest.NewRequest(http.MethodGet, url, nil)
	fileReq.Header.Set("Authorization", "Bearer "+token)
	fileResp := httptest.NewRecorder()
	env.mux.ServeHTTP(fileResp, fileReq)
	if fileResp.Code != http.StatusOK {
		t.Fatalf("file status = %d", fileResp.Code)
	}
	if got := fileResp.Header().Get("Content-Disposition"); got != `attachment; filename="download"` {
		t.Fatalf("Content-Disposition = %q, want attachment", got)
	}
	if got := fileResp.Header().Get("Content-Security-Policy"); got != "sandbox" {
		t.Fatalf("Content-Security-Policy = %q, want sandbox", got)
	}
	if got := fileResp.Header().Get("Content-Type"); got != "text/plain; charset=utf-8" {
		t.Fatalf("served Content-Type = %q, want detected text/plain", got)
	}

	// A-002 F-001: SVG, XML-prefixed SVG, header-smuggled HTML, and GIF-multiplexed
	// HTML are all rejected by the active-content marker gate even though
	// http.DetectContentType would sniff them as text/plain / text/xml / image/gif.
	smuggled := [][2]string{
		{"svg-plain.txt", "<svg xmlns=\"http://www.w3.org/2000/svg\"><script>alert(1)</script></svg>"},
		{"svg-xml.xml", "<?xml version=\"1.0\"?><svg xmlns=\"http://www.w3.org/2000/svg\"><script>alert(1)</script></svg>"},
		{"html-padded.txt", strings.Repeat("A", 600) + "<!DOCTYPE html><script>alert(1)</script>"},
		{"gif-html.gif", "GIF89a" + "<html><script>alert(1)</script></html>"},
	}
	for _, tc := range smuggled {
		resp := uploadPart(tc[0], "text/plain", tc[1])
		if resp.Code != http.StatusUnsupportedMediaType {
			t.Fatalf("smuggled %q status = %d, want 415", tc[0], resp.Code)
		}
	}

	// A-003 N-001: tag names are case-insensitive — mixed-case SVG/Script must
	// also be rejected.
	mixedCase := [][2]string{
		{"svg-mixed.txt", "<Svg xmlns=\"http://www.w3.org/2000/svg\"><Script>alert(1)</Script></Svg>"},
		{"xml-mixed.xml", "<?XML version=\"1.0\"?><Svg xmlns=\"http://www.w3.org/2000/svg\"><script>alert(1)</script></Svg>"},
	}
	for _, tc := range mixedCase {
		resp := uploadPart(tc[0], "text/plain", tc[1])
		if resp.Code != http.StatusUnsupportedMediaType {
			t.Fatalf("mixed-case %q status = %d, want 415", tc[0], resp.Code)
		}
	}

	_ = os.RemoveAll(env.uploadDir)
}

// GOAL-003 · upload ownership: the uploader may download; another authenticated
// user receives 403; legacy meta without owner fails closed.
func TestUploadOwnerOnlyDownload(t *testing.T) {
	env := newAuthTestEnv(t)
	adminTok := adminToken(t, env)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "owned.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write([]byte("owner-only-bytes"))
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Authorization", "Bearer "+adminTok)
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
	url, _ := uploaded["url"].(string)
	id, _ := uploaded["id"].(string)
	if url == "" || id == "" {
		t.Fatalf("upload response missing url/id: %v", uploaded)
	}

	// Owner can download.
	ownReq := httptest.NewRequest(http.MethodGet, url, nil)
	ownReq.Header.Set("Authorization", "Bearer "+adminTok)
	ownResp := httptest.NewRecorder()
	env.mux.ServeHTTP(ownResp, ownReq)
	if ownResp.Code != http.StatusOK {
		t.Fatalf("owner download status = %d: %s", ownResp.Code, ownResp.Body.String())
	}
	if ownResp.Body.String() != "owner-only-bytes" {
		t.Fatalf("owner body = %q", ownResp.Body.String())
	}

	// A different authenticated user cannot download (IDOR closed).
	env.addUser(t, "viewer-upload", "viewer-pass-1", []string{"viewer"})
	viewerTok := env.login(t, "viewer-upload", "viewer-pass-1")
	otherReq := httptest.NewRequest(http.MethodGet, url, nil)
	otherReq.Header.Set("Authorization", "Bearer "+viewerTok)
	otherResp := httptest.NewRecorder()
	env.mux.ServeHTTP(otherResp, otherReq)
	if otherResp.Code != http.StatusForbidden {
		t.Fatalf("non-owner download status = %d, want 403: %s", otherResp.Code, otherResp.Body.String())
	}
	var apiError map[string]string
	if err := json.Unmarshal(otherResp.Body.Bytes(), &apiError); err != nil {
		t.Fatal(err)
	}
	if apiError["error"] != "FORBIDDEN" {
		t.Fatalf("non-owner error = %q, want FORBIDDEN", apiError["error"])
	}

	// Legacy/corrupt meta without owner fails closed (403), not world-readable.
	legacyID := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	if err := os.WriteFile(filepath.Join(env.uploadDir, legacyID), []byte("legacy"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(env.uploadDir, legacyID+".meta.json"),
		[]byte(`{"name":"legacy.txt","type":"text/plain"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	legacyReq := httptest.NewRequest(http.MethodGet, "/api/files/"+legacyID, nil)
	legacyReq.Header.Set("Authorization", "Bearer "+adminTok)
	legacyResp := httptest.NewRecorder()
	env.mux.ServeHTTP(legacyResp, legacyReq)
	if legacyResp.Code != http.StatusForbidden {
		t.Fatalf("legacy no-owner download status = %d, want 403: %s", legacyResp.Code, legacyResp.Body.String())
	}

	_ = os.RemoveAll(env.uploadDir)
}

// W4 P0-2: POST /api/upload requires files.write (admin-only by default). A
// low-privilege viewer account must be rejected with 403 before any storage
// happens, so a compromised viewer cannot fill the disk. The admin uploader
// retains the full contract (covered by TestUploadEndpointContract).
func TestUploadRequiresFilesWritePermission(t *testing.T) {
	env := newAuthTestEnv(t)

	// A viewer holds users.read only (no users.write/files.write) → 403.
	env.addUser(t, "viewer2", "pw", []string{"viewer"})
	viewerToken := env.login(t, "viewer2", "pw")

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "attack.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write([]byte("fill-the-disk"))
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Authorization", "Bearer "+viewerToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("viewer upload status = %d, want 403: %s", rr.Code, rr.Body.String())
	}
	var apiError map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &apiError); err != nil {
		t.Fatal(err)
	}
	if apiError["error"] != "FORBIDDEN" {
		t.Fatalf("viewer upload error = %q, want FORBIDDEN", apiError["error"])
	}
	// No file may have been stored (the upload dir is only created on a
	// successful save, so an absent dir also means zero stored files).
	entries, err := os.ReadDir(env.uploadDir)
	if err != nil {
		if os.IsNotExist(err) {
			entries = nil
		} else {
			t.Fatal(err)
		}
	}
	if len(entries) != 0 {
		t.Fatalf("viewer upload stored %d files, want 0", len(entries))
	}

	// The admin uploader still succeeds (files.write held by admin).
	adminTok := adminToken(t, env)
	body2 := new(bytes.Buffer)
	writer2 := multipart.NewWriter(body2)
	part2, err := writer2.CreateFormFile("file", "ok.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part2.Write([]byte("admin-file"))
	if err := writer2.Close(); err != nil {
		t.Fatal(err)
	}
	req2 := httptest.NewRequest(http.MethodPost, "/api/upload", body2)
	req2.Header.Set("Authorization", "Bearer "+adminTok)
	req2.Header.Set("Content-Type", writer2.FormDataContentType())
	rr2 := httptest.NewRecorder()
	env.mux.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusOK {
		t.Fatalf("admin upload status = %d, want 200: %s", rr2.Code, rr2.Body.String())
	}
	_ = os.RemoveAll(env.uploadDir)
}

// W4 P0-2: the per-user quota bounds a files.write holder's stored objects —
// a compromised admin cannot pump the disk past the configured cap. A tiny
// maxUserFiles makes the gate observable without thousands of uploads.
func TestUploadPerUserQuota(t *testing.T) {
	prevFiles := maxUserFiles
	prevBytes := maxUserBytes
	maxUserFiles = 2
	maxUserBytes = 0 // disable byte quota for a focused count check
	t.Cleanup(func() {
		maxUserFiles = prevFiles
		maxUserBytes = prevBytes
	})

	env := newAuthTestEnv(t)
	tok := adminToken(t, env)

	uploadOne := func() int {
		body := new(bytes.Buffer)
		w := multipart.NewWriter(body)
		part, err := w.CreateFormFile("file", "q.txt")
		if err != nil {
			t.Fatal(err)
		}
		_, _ = part.Write([]byte("quota"))
		if err := w.Close(); err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
		req.Header.Set("Authorization", "Bearer "+tok)
		req.Header.Set("Content-Type", w.FormDataContentType())
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		return rr.Code
	}

	if code := uploadOne(); code != http.StatusOK {
		t.Fatalf("first upload = %d, want 200", code)
	}
	if code := uploadOne(); code != http.StatusOK {
		t.Fatalf("second upload = %d, want 200", code)
	}
	if code := uploadOne(); code != http.StatusTooManyRequests {
		t.Fatalf("third upload = %d, want 429 (quota)", code)
	}
	_ = os.RemoveAll(env.uploadDir)
}
