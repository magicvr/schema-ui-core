package handler

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"hash/crc32"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// makePNG encodes a solid-color WxH PNG (alpha 0 = fully transparent).
func makePNG(t *testing.T, w, h int, c color.Color) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, c)
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}
	return buf.Bytes()
}

func uploadBrandAsset(t *testing.T, env *authTestEnv, token, kind string, body []byte, name string) *httptest.ResponseRecorder {
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
	req := httptest.NewRequest(http.MethodPost, "/api/branding/assets?kind="+kind, &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	return rr
}

func assetURL(t *testing.T, rr *httptest.ResponseRecorder) string {
	t.Helper()
	var body map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("decode upload response %q: %v", rr.Body.String(), err)
	}
	url, _ := body["url"].(string)
	if url == "" {
		t.Fatalf("upload response missing url: %v", body)
	}
	return url
}

func patchSettings(t *testing.T, env *authTestEnv, token, body string) *httptest.ResponseRecorder {
	t.Helper()
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, token, http.MethodPatch, "/api/settings/default", body))
	return rr
}

func TestBrandingAssetUploadAndPublicServe(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)

	// Opaque 1024x1024 logo -> server re-encodes to <=512 JPEG.
	rr := uploadBrandAsset(t, env, token, "logo", makePNG(t, 1024, 1024, color.RGBA{255, 0, 0, 255}), "logo.png")
	if rr.Code != http.StatusOK {
		t.Fatalf("upload = %d: %s", rr.Code, rr.Body.String())
	}
	var body map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&body)
	if body["type"] != "image/jpeg" {
		t.Fatalf("opaque output type = %v, want image/jpeg", body["type"])
	}
	url, _ := body["url"].(string)
	if url == "" || !assetIDOk(url) {
		t.Fatalf("url = %q, want /api/branding/assets/{id}", url)
	}

	// Public GET (no auth): 200 + pinned headers; body is a <=512 JPEG.
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
	if b.Dx() > 512 || b.Dy() > 512 {
		t.Fatalf("logo dims = %dx%d, want <=512", b.Dx(), b.Dy())
	}
}

func TestBrandingAssetFaviconKeepsTransparencyAsPNG(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)

	// 400x400 semi-transparent favicon -> 64px PNG (alpha preserved).
	rr := uploadBrandAsset(t, env, token, "favicon", makePNG(t, 400, 400, color.RGBA{0, 128, 0, 128}), "favicon.png")
	if rr.Code != http.StatusOK {
		t.Fatalf("upload = %d: %s", rr.Code, rr.Body.String())
	}
	var body map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&body)
	if body["type"] != "image/png" {
		t.Fatalf("alpha output type = %v, want image/png", body["type"])
	}
	url, _ := body["url"].(string)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	serve := httptest.NewRecorder()
	env.mux.ServeHTTP(serve, req)
	decoded, err := png.Decode(bytes.NewReader(serve.Body.Bytes()))
	if err != nil {
		t.Fatalf("response is not a decodable png: %v", err)
	}
	b := decoded.Bounds()
	if b.Dx() > 64 || b.Dy() > 64 {
		t.Fatalf("favicon dims = %dx%d, want <=64", b.Dx(), b.Dy())
	}
}

func TestBrandingAssetUploadRejections(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)

	// SVG (active content) -> 415.
	svg := []byte("<svg xmlns='http://www.w3.org/2000/svg'><script>alert(1)</script></svg>")
	if rr := uploadBrandAsset(t, env, token, "logo", svg, "logo.svg"); rr.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("svg upload = %d, want 415", rr.Code)
	}

	// Empty file -> 400.
	if rr := uploadBrandAsset(t, env, token, "logo", nil, "empty.png"); rr.Code != http.StatusBadRequest {
		t.Fatalf("empty upload = %d, want 400", rr.Code)
	}

	// Unknown kind -> 400.
	if rr := uploadBrandAsset(t, env, token, "banner", makePNG(t, 8, 8, color.White), "b.png"); rr.Code != http.StatusBadRequest {
		t.Fatalf("bad kind = %d, want 400", rr.Code)
	}

	// Anonymous -> 401.
	if rr := uploadBrandAsset(t, env, "", "logo", makePNG(t, 8, 8, color.White), "a.png"); rr.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous upload = %d, want 401", rr.Code)
	}

	// Editor without settings.write -> 403.
	env.addUser(t, "ed", "editor-pass-1", []string{"editor"})
	editorToken := loginAs(t, env, "ed", "editor-pass-1")
	if rr := uploadBrandAsset(t, env, editorToken, "logo", makePNG(t, 8, 8, color.White), "e.png"); rr.Code != http.StatusForbidden {
		t.Fatalf("editor upload = %d, want 403", rr.Code)
	}

	// Non-decodable bytes with a benign sniff (text/plain) -> 415.
	if rr := uploadBrandAsset(t, env, token, "logo", []byte("just some text"), "x.txt"); rr.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("text upload = %d, want 415", rr.Code)
	}
}

func TestBrandingAssetOversizeRejected(t *testing.T) {
	testBrandOpts = &BrandingAssetsOptions{MaxBytes: 1024}
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	// Noisy PNG (poor compression) exceeds the 1 KiB policy cap.
	noisy := image.NewRGBA(image.Rect(0, 0, 128, 128))
	for y := 0; y < 128; y++ {
		for x := 0; x < 128; x++ {
			noisy.SetRGBA(x, y, color.RGBA{uint8(x * 7), uint8(y * 13), uint8(x * y), 255})
		}
	}
	var noisyBuf bytes.Buffer
	if err := png.Encode(&noisyBuf, noisy); err != nil {
		t.Fatalf("encode noisy png: %v", err)
	}
	big := noisyBuf.Bytes()
	if len(big) <= 1024 {
		t.Fatalf("fixture png = %d bytes, need > 1024", len(big))
	}
	if rr := uploadBrandAsset(t, env, token, "logo", big, "big.png"); rr.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("oversize upload = %d, want 413", rr.Code)
	}
}

func TestBrandingAssetGetMissingAndInvalidID(t *testing.T) {
	env := newAuthTestEnv(t)
	for _, id := range []string{"missing", "00000000000000000000000000000000"} {
		req := httptest.NewRequest(http.MethodGet, "/api/branding/assets/"+id, nil)
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("GET %q = %d, want 404", id, rr.Code)
		}
	}
	// Path traversal: the router normalizes (307 redirect) before the handler;
	// whatever the exact response, it must never serve content (never 200).
	req := httptest.NewRequest(http.MethodGet, "/api/branding/assets/../etc/passwd", nil)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code == http.StatusOK {
		t.Fatalf("traversal GET served content: %d", rr.Code)
	}
}

func TestBrandingAssetCleanupOnReplaceAndClear(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)

	// Upload A and commit it as the logo.
	rr := uploadBrandAsset(t, env, token, "logo", makePNG(t, 64, 64, color.RGBA{255, 0, 0, 255}), "a.png")
	urlA := assetURL(t, rr)
	idA, _ := BrandAssetIDFromURL(urlA)
	out := patchSettings(t, env, token, "{\"logoUrl\":\""+urlA+"\"}")
	if out.Code != http.StatusOK {
		t.Fatalf("patch A = %d: %s", out.Code, out.Body.String())
	}
	if _, err := os.Stat(filepath.Join(env.brandAssets.dir, idA)); err != nil {
		t.Fatalf("asset A missing after commit: %v", err)
	}

	// Replace with B -> A is deleted, B survives.
	rr = uploadBrandAsset(t, env, token, "logo", makePNG(t, 64, 64, color.RGBA{0, 255, 0, 255}), "b.png")
	urlB := assetURL(t, rr)
	idB, _ := BrandAssetIDFromURL(urlB)
	out = patchSettings(t, env, token, "{\"logoUrl\":\""+urlB+"\"}")
	if out.Code != http.StatusOK {
		t.Fatalf("patch B = %d: %s", out.Code, out.Body.String())
	}
	if _, err := os.Stat(filepath.Join(env.brandAssets.dir, idA)); !os.IsNotExist(err) {
		t.Fatalf("asset A still present after replace (err=%v)", err)
	}
	if _, err := os.Stat(filepath.Join(env.brandAssets.dir, idB)); err != nil {
		t.Fatalf("asset B missing after commit: %v", err)
	}

	// Clear the field -> B is deleted too.
	out = patchSettings(t, env, token, "{\"logoUrl\":\"\"}")
	if out.Code != http.StatusOK {
	t.Fatalf("clear = %d: %s", out.Code, out.Body.String())
	}
	if _, err := os.Stat(filepath.Join(env.brandAssets.dir, idB)); !os.IsNotExist(err) {
		t.Fatalf("asset B still present after clear (err=%v)", err)
	}
}

func TestBrandingAssetCleanupOnReset(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)

	// One referenced + one orphan asset; reset deletes both.
	rr := uploadBrandAsset(t, env, token, "logo", makePNG(t, 32, 32, color.RGBA{255, 0, 0, 255}), "ref.png")
	urlRef := assetURL(t, rr)
	rr = uploadBrandAsset(t, env, token, "favicon", makePNG(t, 32, 32, color.RGBA{0, 0, 255, 255}), "orphan.png")
	urlOrphan := assetURL(t, rr)
	idRef, _ := BrandAssetIDFromURL(urlRef)
	idOrphan, _ := BrandAssetIDFromURL(urlOrphan)

	out := patchSettings(t, env, token, "{\"logoUrl\":\""+urlRef+"\",\"faviconUrl\":\""+urlOrphan+"\"}")
	if out.Code != http.StatusOK {
		t.Fatalf("patch = %d: %s", out.Code, out.Body.String())
	}

	reset := httptest.NewRecorder()
	env.mux.ServeHTTP(reset, bearer(t, token, http.MethodPost, "/api/settings/default/reset", ""))
	if reset.Code != http.StatusOK {
		t.Fatalf("reset = %d: %s", reset.Code, reset.Body.String())
	}
	for _, id := range []string{idRef, idOrphan} {
		if _, err := os.Stat(filepath.Join(env.brandAssets.dir, id)); !os.IsNotExist(err) {
			t.Fatalf("asset %s still present after reset (err=%v)", id, err)
		}
	}
}

func TestBrandingAssetSharedReferenceSurvivesReplace(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)

	// The same asset id is referenced by logoUrl AND faviconUrl; replacing
	// logoUrl alone must not delete the asset still used by faviconUrl.
	rr := uploadBrandAsset(t, env, token, "logo", makePNG(t, 32, 32, color.RGBA{9, 9, 9, 255}), "shared.png")
	urlShared := assetURL(t, rr)
	idShared, _ := BrandAssetIDFromURL(urlShared)
	out := patchSettings(t, env, token, "{\"logoUrl\":\""+urlShared+"\",\"faviconUrl\":\""+urlShared+"\"}")
	if out.Code != http.StatusOK {
		t.Fatalf("patch shared = %d: %s", out.Code, out.Body.String())
	}
	rr = uploadBrandAsset(t, env, token, "logo", makePNG(t, 32, 32, color.RGBA{8, 8, 8, 255}), "replacement.png")
	urlNew := assetURL(t, rr)
	out = patchSettings(t, env, token, "{\"logoUrl\":\""+urlNew+"\"}")
	if out.Code != http.StatusOK {
		t.Fatalf("patch replace = %d: %s", out.Code, out.Body.String())
	}
	if _, err := os.Stat(filepath.Join(env.brandAssets.dir, idShared)); err != nil {
		t.Fatalf("shared asset deleted although faviconUrl still references it: %v", err)
	}
	// Clearing faviconUrl now releases it.
	out = patchSettings(t, env, token, "{\"faviconUrl\":\"\"}")
	if out.Code != http.StatusOK {
		t.Fatalf("patch clear favicon = %d: %s", out.Code, out.Body.String())
	}
	if _, err := os.Stat(filepath.Join(env.brandAssets.dir, idShared)); !os.IsNotExist(err) {
		t.Fatalf("shared asset still present after last reference cleared (err=%v)", err)
	}
}


func TestBrandingAssetStartupGC(t *testing.T) {
	dir := t.TempDir()
	store := NewBrandingAssetStore(dir, DefaultBrandingAssetsOptions())
	keepID, err := store.save("image/png", "logo", "", makePNG(t, 16, 16, color.RGBA{255, 0, 0, 255}))
	if err != nil {
		t.Fatalf("save keep: %v", err)
	}
	dropID, err := store.save("image/png", "favicon", "", makePNG(t, 16, 16, color.RGBA{0, 0, 255, 255}))
	if err != nil {
		t.Fatalf("save drop: %v", err)
	}
	// Legacy URL references and empty strings are ignored (no-op).
	if err := store.GC([]string{BrandAssetURLPrefix + keepID, "https://cdn.example/logo.png", ""}); err != nil {
		t.Fatalf("gc: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, keepID)); err != nil {
		t.Fatalf("referenced asset removed by gc: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, dropID)); !os.IsNotExist(err) {
		t.Fatalf("orphan asset survived gc (err=%v)", err)
	}
}


// oversizedPNG builds a PNG header declaring 30000x30000 with a valid IHDR
// CRC (decompression-bomb fixture, A-002 F-001).
func oversizedPNG(t *testing.T) []byte {
	t.Helper()
	var buf bytes.Buffer
	buf.Write([]byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a})
	var length [4]byte
	binary.BigEndian.PutUint32(length[:], 13)
	buf.Write(length[:])
	ihdr := make([]byte, 13)
	binary.BigEndian.PutUint32(ihdr[0:4], 30000)
	binary.BigEndian.PutUint32(ihdr[4:8], 30000)
	ihdr[8] = 8 // bit depth
	ihdr[9] = 2 // color type RGB
	chunk := append([]byte("IHDR"), ihdr...)
	buf.Write(chunk)
	var crc [4]byte
	binary.BigEndian.PutUint32(crc[:], crc32.ChecksumIEEE(chunk))
	buf.Write(crc[:])
	return buf.Bytes()
}

func TestBrandingAssetRejectsOversizedDimensions(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	// Tiny file (header only) declaring 30000x30000 -> rejected before decode.
	if rr := uploadBrandAsset(t, env, token, "logo", oversizedPNG(t), "bomb.png"); rr.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("oversized-dimension upload = %d, want 415", rr.Code)
	}
}

func TestBrandingAssetJpegAndGifInputs(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)

	// JPEG input: opaque photo-ish source re-encoded (JPEG out).
	photo := image.NewRGBA(image.Rect(0, 0, 320, 240))
	for y := 0; y < 240; y++ {
		for x := 0; x < 320; x++ {
			photo.SetRGBA(x, y, color.RGBA{uint8(x), uint8(y), 128, 255})
		}
	}
	var jpgBuf bytes.Buffer
	if err := jpeg.Encode(&jpgBuf, photo, &jpeg.Options{Quality: 90}); err != nil {
		t.Fatalf("encode jpeg: %v", err)
	}
	if rr := uploadBrandAsset(t, env, token, "logo", jpgBuf.Bytes(), "photo.jpg"); rr.Code != http.StatusOK {
		t.Fatalf("jpeg upload = %d: %s", rr.Code, rr.Body.String())
	}

	// GIF input: first frame accepted (animated input collapses to a frame).
	var gifBuf bytes.Buffer
	if err := gif.Encode(&gifBuf, image.NewRGBA(image.Rect(0, 0, 16, 16)), nil); err != nil {
		t.Fatalf("encode gif: %v", err)
	}
	if rr := uploadBrandAsset(t, env, token, "logo", gifBuf.Bytes(), "frame.gif"); rr.Code != http.StatusOK {
		t.Fatalf("gif upload = %d: %s", rr.Code, rr.Body.String())
	}
}

// assetIDOk reports whether url matches the public asset URL shape.
func assetIDOk(url string) bool {
	_, ok := BrandAssetIDFromURL(url)
	return ok
}