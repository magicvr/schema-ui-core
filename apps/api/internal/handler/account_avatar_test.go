// Account avatar upload tests (W13 T-05 · GOAL-014).
package handler

import (
	"bytes"
	"encoding/json"
	"image/color"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func uploadAvatar(t *testing.T, env *authTestEnv, token string, body []byte, name string) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", name)
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := part.Write(body); err != nil {
		t.Fatalf("write part: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/account/avatar", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	return rr
}

func avatarIDOk(url string) bool {
	_, ok := NewAvatarAssetStore("", DefaultBrandingAssetsOptions()).AssetIDFromURL(url)
	return ok
}

// avatarURL validates and returns the url from an already-decoded upload body.
func avatarURL(t *testing.T, body map[string]any) string {
	t.Helper()
	url, _ := body["url"].(string)
	if url == "" || !avatarIDOk(url) {
		t.Fatalf("upload response url = %q, want /api/account/avatars/{id}", url)
	}
	return url
}

// decodeUpload decodes an upload response body exactly once.
func decodeUpload(t *testing.T, rr *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("decode upload response %q: %v", rr.Body.String(), err)
	}
	return body
}

func patchProfile(t *testing.T, env *authTestEnv, token, body string) *httptest.ResponseRecorder {
	t.Helper()
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, token, http.MethodPatch, "/api/account/profile", body))
	return rr
}

func TestAccountAvatarUploadAndPublicServe(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)

	// Opaque 1024x1024 avatar -> server re-encodes to <=256 JPEG.
	rr := uploadAvatar(t, env, token, makePNG(t, 1024, 1024, color.RGBA{255, 0, 0, 255}), "avatar.png")
	if rr.Code != http.StatusOK {
		t.Fatalf("upload = %d: %s", rr.Code, rr.Body.String())
	}
	body := decodeUpload(t, rr)
	if body["type"] != "image/jpeg" {
	t.Fatalf("opaque output type = %v, want image/jpeg", body["type"])
	}
	url := avatarURL(t, body)

	// Public GET (no auth): 200 + pinned headers; body is a <=256 JPEG.
	req := httptest.NewRequest(http.MethodGet, url, nil)
	serve := httptest.NewRecorder()
	env.mux.ServeHTTP(serve, req)
	if serve.Code != http.StatusOK {
		t.Fatalf("public get = %d", serve.Code)
	}
	if ct := serve.Header().Get("Content-Type"); ct != "image/jpeg" {
		t.Fatalf("content-type = %q", ct)
	}
	if serve.Header().Get("X-Content-Type-Options") != "nosniff" {
	t.Fatal("missing nosniff")
	}
	if cc := serve.Header().Get("Cache-Control"); cc != "public, max-age=31536000, immutable" {
	t.Fatalf("cache-control = %q", cc)
	}
	if csp := serve.Header().Get("Content-Security-Policy"); csp != "sandbox" {
	t.Fatalf("content-security-policy = %q, want sandbox", csp)
	}
	decoded, err := jpeg.Decode(bytes.NewReader(serve.Body.Bytes()))
	if err != nil {
	t.Fatalf("response is not a decodable jpeg: %v", err)
	}
	b := decoded.Bounds()
	if b.Dx() > 256 || b.Dy() > 256 {
	t.Fatalf("avatar dims = %dx%d, want <=256", b.Dx(), b.Dy())
	}
}

func TestAccountAvatarTransparencyStaysPNG(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	rr := uploadAvatar(t, env, token, makePNG(t, 400, 400, color.RGBA{0, 128, 0, 128}), "avatar.png")
	if rr.Code != http.StatusOK {
		t.Fatalf("upload = %d: %s", rr.Code, rr.Body.String())
	}
	body := decodeUpload(t, rr)
	if body["type"] != "image/png" {
	t.Fatalf("alpha output type = %v, want image/png", body["type"])
	}
	url := avatarURL(t, body)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	serve := httptest.NewRecorder()
	env.mux.ServeHTTP(serve, req)
	if _, err := png.Decode(bytes.NewReader(serve.Body.Bytes())); err != nil {
	t.Fatalf("response is not a decodable png: %v", err)
	}
}

func TestAccountAvatarRejections(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)

	// Anonymous -> 401.
	if rr := uploadAvatar(t, env, "", makePNG(t, 8, 8, color.White), "a.png"); rr.Code != http.StatusUnauthorized {
	t.Fatalf("anonymous upload = %d, want 401", rr.Code)
	}

	// SVG (active content) -> 415.
	svg := []byte("<svg xmlns='http://www.w3.org/2000/svg'><script>alert(1)</script></svg>")
	if rr := uploadAvatar(t, env, token, svg, "logo.svg"); rr.Code != http.StatusUnsupportedMediaType {
	t.Fatalf("svg upload = %d, want 415", rr.Code)
	}

	// Empty file -> 400.
	if rr := uploadAvatar(t, env, token, nil, "empty.png"); rr.Code != http.StatusBadRequest {
	t.Fatalf("empty upload = %d, want 400", rr.Code)
	}

	// Non-decodable bytes -> 415.
	if rr := uploadAvatar(t, env, token, []byte("just some text"), "x.txt"); rr.Code != http.StatusUnsupportedMediaType {
	t.Fatalf("text upload = %d, want 415", rr.Code)
	}
}

func TestAccountAvatarProfileCommitReplaceClear(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)

	// Upload A, commit it via the profile PATCH, and confirm the profile GET.
	rr := uploadAvatar(t, env, token, makePNG(t, 64, 64, color.RGBA{255, 0, 0, 255}), "a.png")
	if rr.Code != http.StatusOK {
	t.Fatalf("upload A = %d: %s", rr.Code, rr.Body.String())
	}
	urlA := avatarURL(t, decodeUpload(t, rr))
	idA, _ := NewAvatarAssetStore("", DefaultBrandingAssetsOptions()).AssetIDFromURL(urlA)
	out := patchProfile(t, env, token, "{\"name\":\"Admin\",\"avatarUrl\":\""+urlA+"\"}")
	if out.Code != http.StatusOK {
	t.Fatalf("patch A = %d: %s", out.Code, out.Body.String())
	}
	// Profile GET returns avatarUrl.
	get := httptest.NewRecorder()
	env.mux.ServeHTTP(get, bearer(t, token, http.MethodGet, "/api/account/profile", ""))
	var row map[string]any
	_ = json.NewDecoder(get.Body).Decode(&row)
	if row["avatarUrl"] != urlA {
	t.Fatalf("profile avatarUrl = %v, want %s", row["avatarUrl"], urlA)
	}
	if _, err := os.Stat(filepath.Join(env.avatarAssets.dir, idA)); err != nil {
	t.Fatalf("asset A missing after commit: %v", err)
	}

	// Replace with B -> the profile commit deletes A and keeps B.
	rr = uploadAvatar(t, env, token, makePNG(t, 64, 64, color.RGBA{0, 255, 0, 255}), "b.png")
	if rr.Code != http.StatusOK {
	t.Fatalf("upload B = %d: %s", rr.Code, rr.Body.String())
	}
	urlB := avatarURL(t, decodeUpload(t, rr))
	idB, _ := NewAvatarAssetStore("", DefaultBrandingAssetsOptions()).AssetIDFromURL(urlB)
	out = patchProfile(t, env, token, "{\"name\":\"Admin\",\"avatarUrl\":\""+urlB+"\"}")
	if out.Code != http.StatusOK {
	t.Fatalf("patch B = %d: %s", out.Code, out.Body.String())
	}
	if _, err := os.Stat(filepath.Join(env.avatarAssets.dir, idA)); !os.IsNotExist(err) {
	t.Fatalf("asset A still present after replace (err=%v)", err)
	}
	if _, err := os.Stat(filepath.Join(env.avatarAssets.dir, idB)); err != nil {
	t.Fatalf("asset B missing after commit: %v", err)
	}

	// Clearing avatarUrl deletes B too.
	out = patchProfile(t, env, token, "{\"name\":\"Admin\",\"avatarUrl\":\"\"}")
	if out.Code != http.StatusOK {
	t.Fatalf("clear = %d: %s", out.Code, out.Body.String())
	}
	if _, err := os.Stat(filepath.Join(env.avatarAssets.dir, idB)); !os.IsNotExist(err) {
	t.Fatalf("asset B still present after clear (err=%v)", err)
	}
}

func TestAccountAvatarProfileRejectsForeignURL(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	// A brand-asset URL (or any other origin) must never be committed as an avatar.
	out := patchProfile(t, env, token, "{\"name\":\"Admin\",\"avatarUrl\":\"/api/branding/assets/00000000000000000000000000000000\"}")
	if out.Code != http.StatusBadRequest {
	t.Fatalf("foreign avatarUrl = %d, want 400: %s", out.Code, out.Body.String())
	}
}

func TestAccountAvatarMissingAsset404(t *testing.T) {
	env := newAuthTestEnv(t)
	req := httptest.NewRequest(http.MethodGet, "/api/account/avatars/00000000000000000000000000000000", nil)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
	t.Fatalf("missing avatar = %d, want 404", rr.Code)
	}
}
