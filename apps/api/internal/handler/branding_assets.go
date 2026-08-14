// Brand asset upload endpoints (W9 / GOAL-010).
//
// Brand icons (logo / light / dark / favicon) are configured on the Settings
// page. The legacy textarea-URL input is replaced by uploads: images are
// stored in a dedicated brand-assets directory — NOT the generic /api/upload
// owner store and NOT the admin.file-library module — and served publicly,
// because the login page and shell load branding before authentication.
//
// Security model:
//   - uploads require settings.write (module permission gate)
//   - every image is re-encoded server-side (PNG/JPEG/GIF/WebP -> PNG or
//     JPEG, dimension-limited); raw user bytes are never stored or served
//   - the public GET only ever serves server-produced raster output with
//     nosniff + sandbox + immutable caching
package handler

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/webp"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// BrandingAssetsOptions is the W9 processing policy; values are config-driven
// (config.yaml branding section, env-overridable).
type BrandingAssetsOptions struct {
	// MaxBytes bounds a single upload (default 4 MiB).
	MaxBytes int
	// LogoMaxDim is the longest-edge limit for logo assets (default 512).
	LogoMaxDim int
	// FaviconDim is the longest-edge limit for favicon assets (default 64).
	FaviconDim int
	// JPEGQuality (1..100) for opaque output (default 82).
	JPEGQuality int
}

// DefaultBrandingAssetsOptions returns the documented W9 defaults.
func DefaultBrandingAssetsOptions() BrandingAssetsOptions {
	return BrandingAssetsOptions{
		MaxBytes:    4 << 20,
		LogoMaxDim:  512,
		FaviconDim:  64,
		JPEGQuality: 82,
	}
}

// BrandingAssetStore persists processed brand images in a dedicated
// directory. Object ids are 16 random bytes (hex), same shape as the generic
// upload store; each object carries a {type, kind} meta file.
type BrandingAssetStore struct {
	dir  string
	opts BrandingAssetsOptions
}

// NewBrandingAssetStore constructs the store; zero/out-of-range options fall
// back to the documented defaults (config validation never passes zeros).
func NewBrandingAssetStore(dir string, opts BrandingAssetsOptions) *BrandingAssetStore {
	d := DefaultBrandingAssetsOptions()
	if opts.MaxBytes > 0 {
		d.MaxBytes = opts.MaxBytes
	}
	if opts.LogoMaxDim > 0 {
		d.LogoMaxDim = opts.LogoMaxDim
	}
	if opts.FaviconDim > 0 {
		d.FaviconDim = opts.FaviconDim
	}
	if opts.JPEGQuality > 0 && opts.JPEGQuality <= 100 {
		d.JPEGQuality = opts.JPEGQuality
	}
	return &BrandingAssetStore{dir: dir, opts: d}
}

// BrandAssetURLPrefix is the public URL prefix served by GET
// /api/branding/assets/{id}.
const BrandAssetURLPrefix = "/api/branding/assets/"

// BrandAssetIDFromURL extracts the asset id from a stored branding URL (e.g.
// the value committed to site_settings.logo_url). Non-asset URLs (legacy
// http(s) or static same-origin paths) return ok=false.
func BrandAssetIDFromURL(raw string) (string, bool) {
	if !strings.HasPrefix(raw, BrandAssetURLPrefix) {
		return "", false
	}
	id := strings.TrimPrefix(raw, BrandAssetURLPrefix)
	if !uploadFileIDPattern.MatchString(id) {
		return "", false
	}
	return id, true
}

// BrandingAssetRoutes returns the admin.settings brand asset surface:
// authenticated upload (settings.write gate) + public GET.
func BrandingAssetRoutes(a authMiddleware, store *BrandingAssetStore, moduleID string) []kernel.RouteContribution {
	identity := func(method, pattern string) kernel.ContributionIdentity {
		return kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey(method, pattern)}
	}
	return []kernel.RouteContribution{
		{ContributionIdentity: identity("POST", "/api/branding/assets"), Method: "POST", Pattern: "/api/branding/assets", Handler: a.Middleware(store.upload())},
		{ContributionIdentity: identity("GET", "/api/branding/assets/{id}"), Method: "GET", Pattern: "/api/branding/assets/{id}", Handler: store.file(), Public: true},
	}
}

// RegisterPublicBrandingAssets mounts the public GET for profiles without
// admin.settings (mvp): /api/branding may reference previously uploaded
// assets, so the read surface must stay available. Do not double-register
// when the settings module contributes the same route.
func RegisterPublicBrandingAssets(mux *http.ServeMux, store *BrandingAssetStore) {
	mux.Handle("GET /api/branding/assets/{id}", store.file())
}

// upload handles POST /api/branding/assets?kind=logo|favicon.
func (s *BrandingAssetStore) upload() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r, "settings.write"); !ok {
			return
		}
		kind := strings.TrimSpace(r.URL.Query().Get("kind"))
		if kind == "" {
			kind = "logo"
		}
		if kind != "logo" && kind != "favicon" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_KIND", "kind must be logo or favicon")
			return
		}
		// Manual multipart walk: NextPart returns the part header WITHOUT
		// consuming the body, so an oversized part is rejected (413) before
		// any payload is read; the payload itself is read through a
		// LimitReader, so a lying declared size cannot bypass the cap.
		reader, err := r.MultipartReader()
		if err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_UPLOAD", "expected a multipart form")
			return
		}
		part, err := reader.NextPart()
		if err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_UPLOAD", "expected a multipart file part named file")
			return
		}
		defer part.Close()
		if part.FormName() != "file" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_UPLOAD", "expected a multipart file part named file")
			return
		}
		if part.FileName() == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_FILE", "empty files are rejected")
			return
		}
		body, err := io.ReadAll(io.LimitReader(part, int64(s.opts.MaxBytes)+1))
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "STORAGE_UNAVAILABLE", "could not read upload")
			return
		}
		if len(body) > s.opts.MaxBytes {
			writeLocalizedError(w, r, http.StatusRequestEntityTooLarge, "FILE_TOO_LARGE", "file exceeds the server size limit")
			return
		}
		if len(body) == 0 {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_FILE", "empty files are rejected")
			return
		}
		// Same hard rejections as the generic upload store: sniffed dangerous
		// types and active-content markers (SVG/HTML/script) never reach the
		// processor; every accepted image is re-encoded anyway.
		detected := http.DetectContentType(body)
		base := detected
		if i := strings.IndexByte(base, ';'); i >= 0 {
			base = strings.TrimSpace(base[:i])
		}
		if dangerousInlineTypes[base] || containsActiveContent(body) {
			writeLocalizedError(w, r, http.StatusUnsupportedMediaType, "UNSUPPORTED_FILE_TYPE", "file type is not allowed")
			return
		}
		processed, contentType, err := processBrandingImage(body, kind, s.opts)
		if err != nil {
			writeLocalizedError(w, r, http.StatusUnsupportedMediaType, "UNSUPPORTED_FILE_TYPE", "image must be a decodable PNG, JPEG, GIF or WebP")
			return
		}
		id, err := s.save(contentType, kind, processed)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "STORAGE_UNAVAILABLE", "could not store asset")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"id":   id,
			"name": part.FileName(),
			"type": contentType,
			"size": len(processed),
			"url":  BrandAssetURLPrefix + id,
		})
	})
}

// file serves GET /api/branding/assets/{id} publicly. The store only ever
// contains server-produced raster output, so inline rendering is safe; the
// headers still pin type, nosniff and immutable caching (content-addressed).
func (s *BrandingAssetStore) file() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		body, meta, err := s.load(id)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				writeLocalizedError(w, r, http.StatusNotFound, "ASSET_NOT_FOUND", "no asset with that id")
				return
			}
			writeLocalizedError(w, r, http.StatusInternalServerError, "STORAGE_UNAVAILABLE", "could not read asset")
			return
		}
		contentType := meta["type"]
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Security-Policy", "sandbox")
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	})
}

// save persists a processed asset + its meta file.
func (s *BrandingAssetStore) save(contentType, kind string, body []byte) (string, error) {
	idBytes := make([]byte, 16)
	if _, err := rand.Read(idBytes); err != nil {
		return "", err
	}
	id := hex.EncodeToString(idBytes)
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(s.dir, id), body, 0o644); err != nil {
		return "", err
	}
	meta := map[string]string{"type": contentType, "kind": kind}
	raw, err := json.Marshal(meta)
	if err == nil {
		_ = os.WriteFile(filepath.Join(s.dir, id+".meta.json"), raw, 0o644)
	}
	return id, nil
}

// load reads a stored asset by id (id shape validated — same guard as the
// generic upload store, so a crafted PathValue cannot escape the directory).
func (s *BrandingAssetStore) load(id string) ([]byte, map[string]string, error) {
	if !uploadFileIDPattern.MatchString(id) {
		return nil, nil, os.ErrNotExist
	}
	body, err := os.ReadFile(filepath.Join(s.dir, id))
	if err != nil {
		return nil, nil, err
	}
	meta := map[string]string{}
	raw, err := os.ReadFile(filepath.Join(s.dir, id+".meta.json"))
	if err == nil {
		_ = json.Unmarshal(raw, &meta)
	}
	return body, meta, nil
}

// Delete removes an asset (idempotent; a missing file is a no-op).
func (s *BrandingAssetStore) Delete(id string) error {
	if !uploadFileIDPattern.MatchString(id) {
		return nil
	}
	if err := os.Remove(filepath.Join(s.dir, id)); err != nil && !os.IsNotExist(err) {
		return err
	}
	_ = os.Remove(filepath.Join(s.dir, id+".meta.json"))
	return nil
}

// DeleteAll removes every stored brand asset (reset-to-defaults).
func (s *BrandingAssetStore) DeleteAll() error {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if uploadFileIDPattern.MatchString(strings.TrimSuffix(entry.Name(), ".meta.json")) {
			if err := os.Remove(filepath.Join(s.dir, entry.Name())); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	}
	return nil
}

// GC deletes stored assets not referenced by the current settings singleton
// (startup housekeeping: crashed uploads / cancelled form edits). Referenced
// values are the site_settings branding columns (URL strings).
func (s *BrandingAssetStore) GC(referenced []string) error {
	keep := map[string]bool{}
	for _, ref := range referenced {
		if id, ok := BrandAssetIDFromURL(ref); ok {
			keep[id] = true
		}
	}
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		id := strings.TrimSuffix(entry.Name(), ".meta.json")
		if !uploadFileIDPattern.MatchString(id) {
			continue
		}
		if keep[id] {
			continue
		}
		if err := os.Remove(filepath.Join(s.dir, entry.Name())); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

// processBrandingImage re-encodes a decoded raster: dimension-limited
// (never upscaled), PNG when the result has transparency, JPEG otherwise.
// Output is always a fresh server-produced raster — never the input bytes.
func processBrandingImage(body []byte, kind string, opts BrandingAssetsOptions) ([]byte, string, error) {
	target := opts.LogoMaxDim
	if kind == "favicon" {
		target = opts.FaviconDim
	}
	src, err := decodeBrandingImage(body)
	if err != nil {
		return nil, "", err
	}
	bounds := src.Bounds()
	maxSide := bounds.Dx()
	if bounds.Dy() > maxSide {
		maxSide = bounds.Dy()
	}
	scale := 1.0
	if maxSide > target {
		scale = float64(target) / float64(maxSide)
	}
	w := int(math.Round(float64(bounds.Dx()) * scale))
	h := int(math.Round(float64(bounds.Dy()) * scale))
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	xdraw.CatmullRom.Scale(rgba, rgba.Bounds(), src, bounds, xdraw.Over, nil)
	var out bytes.Buffer
	if imageIsOpaque(rgba) {
		if err := jpeg.Encode(&out, rgba, &jpeg.Options{Quality: opts.JPEGQuality}); err != nil {
			return nil, "", fmt.Errorf("encode jpeg: %w", err)
		}
		return out.Bytes(), "image/jpeg", nil
	}
	if err := png.Encode(&out, rgba); err != nil {
		return nil, "", fmt.Errorf("encode png: %w", err)
	}
	return out.Bytes(), "image/png", nil
}

// decodeBrandingImage decodes PNG, JPEG, GIF (first frame) or WebP. Any other
// content (including SVG/HTML smuggled past the sniff checks) fails decode.
func decodeBrandingImage(body []byte) (image.Image, error) {
	reader := func() *bytes.Reader { return bytes.NewReader(body) }
	if img, err := png.Decode(reader()); err == nil {
		return img, nil
	}
	if img, err := jpeg.Decode(reader()); err == nil {
		return img, nil
	}
	if img, err := gif.Decode(reader()); err == nil {
		return img, nil
	}
	if img, err := webp.Decode(reader()); err == nil {
		return img, nil
	}
	return nil, errors.New("unsupported image format")
}

// imageIsOpaque reports whether every pixel has full alpha (JPEG output) or
// transparency is present (PNG output keeps it).
func imageIsOpaque(img *image.RGBA) bool {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if img.RGBAAt(x, y).A != 255 {
				return false
			}
		}
	}
	return true
}
