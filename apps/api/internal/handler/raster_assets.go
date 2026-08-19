// Generic server-re-encoded raster asset store shared by the branding
// assets surface (W9 / GOAL-010) and the account avatar upload
// (W13 T-05 / GOAL-014). Brand icons and user avatars follow the same
// security model:
//   - every image is re-encoded server-side (PNG/JPEG/GIF/WebP -> PNG or
//     JPEG, dimension-limited); raw user bytes are never stored or served
//   - the public GET only ever serves server-produced raster output with
//     nosniff + sandbox + immutable caching
//   - stores live in dedicated directories (never the shared upload store)
// The permission gate differs per surface (branding: settings.write;
// avatar: any authenticated user) and is applied by the route wrappers,
// not by this store.
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
)

// RasterAssetStore persists processed raster assets (brand logos/favicons,
// account avatars) in a dedicated directory. Object ids are 16 random bytes
// (hex), same shape as the generic upload store; each object carries a
// {type, kind} meta file.
type RasterAssetStore struct {
	dir       string
	opts      BrandingAssetsOptions
	urlPrefix string
	// kinds maps a kind name to its longest-edge output target (from opts).
	kinds map[string]func(BrandingAssetsOptions) int
}

// NewRasterAssetStore constructs a generic store; zero/out-of-range options
// fall back to the documented defaults (config validation never passes zeros).
func NewRasterAssetStore(dir string, opts BrandingAssetsOptions, urlPrefix string, kinds map[string]func(BrandingAssetsOptions) int) *RasterAssetStore {
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
	if opts.AvatarDim > 0 {
		d.AvatarDim = opts.AvatarDim
	}
	if opts.JPEGQuality > 0 && opts.JPEGQuality <= 100 {
		d.JPEGQuality = opts.JPEGQuality
	}
	return &RasterAssetStore{dir: dir, opts: d, urlPrefix: urlPrefix, kinds: kinds}
}

// BrandAssetURLPrefix is the public URL prefix served by GET
// /api/branding/assets/{id}.
const BrandAssetURLPrefix = "/api/branding/assets/"

// AvatarAssetURLPrefix is the public URL prefix served by GET
// /api/account/avatars/{id}.
const AvatarAssetURLPrefix = "/api/account/avatars/"

// NewBrandingAssetStore constructs the branding store (kinds logo/favicon,
// targets from BrandingAssetsOptions).
func NewBrandingAssetStore(dir string, opts BrandingAssetsOptions) *RasterAssetStore {
	return NewRasterAssetStore(dir, opts, BrandAssetURLPrefix, map[string]func(BrandingAssetsOptions) int{
		"logo":    func(o BrandingAssetsOptions) int { return o.LogoMaxDim },
		"favicon": func(o BrandingAssetsOptions) int { return o.FaviconDim },
	})
}

// NewAvatarAssetStore constructs the avatar store (kind avatar, target
// BrandingAssetsOptions.AvatarDim, default 256).
func NewAvatarAssetStore(dir string, opts BrandingAssetsOptions) *RasterAssetStore {
	return NewRasterAssetStore(dir, opts, AvatarAssetURLPrefix, map[string]func(BrandingAssetsOptions) int{
		"avatar": func(o BrandingAssetsOptions) int { return o.AvatarDim },
	})
}

// BrandingAssetStore is the historical branding-store name; the branding
// surface now shares the generic raster store implementation.
type BrandingAssetStore = RasterAssetStore

// AssetIDFromURL extracts the asset id from a stored URL of this store (e.g.
// the value committed to site_settings.logo_url or users.avatar_url).
// Non-asset URLs (legacy http(s) or static same-origin paths) return ok=false.
func (s *RasterAssetStore) AssetIDFromURL(raw string) (string, bool) {
	if !strings.HasPrefix(raw, s.urlPrefix) {
		return "", false
	}
	id := strings.TrimPrefix(raw, s.urlPrefix)
	if !uploadFileIDPattern.MatchString(id) {
		return "", false
	}
	return id, true
}

// BrandAssetIDFromURL extracts a brand-asset id from a stored branding URL.
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

// storeUpload runs the shared multipart upload pipeline: multipart walk with
// hard size limits, dangerous-type sniffing, server-side re-encode and
// persistence. It writes error responses itself and returns ok=false on
// failure; on success it returns the response payload (id/name/type/size/url)
// and the caller writes it (so surface-specific side effects, e.g. avatar
// replacement cleanup, can run first). Permission gating is the caller's
// responsibility.
func (s *RasterAssetStore) storeUpload(w http.ResponseWriter, r *http.Request) (map[string]any, bool) {
	return s.storeUploadForOwner(w, r, "")
}

// storeUploadForOwner is storeUpload with an explicit asset owner (used by the
// avatar surface so ownership is recorded in the meta file and later enforced
// by profile PATCH and cleanup). An empty owner means a shared/global asset
// (branding).
func (s *RasterAssetStore) storeUploadForOwner(w http.ResponseWriter, r *http.Request, owner string) (map[string]any, bool) {
	kind := strings.TrimSpace(r.URL.Query().Get("kind"))
	if kind == "" {
		// Single-kind stores (avatar) default; multi-kind stores require an
		// explicit kind (branding defaults to logo for back-compat).
		if len(s.kinds) == 1 {
			for k := range s.kinds {
				kind = k
			}
		} else {
			kind = "logo"
		}
	}
	target, known := s.kinds[kind]
	if !known {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_KIND", "kind is not supported by this store")
		return nil, false
	}
	// Manual multipart walk: NextPart returns the part header WITHOUT
	// consuming the body, so an oversized part is rejected (413) before
	// any payload is read; the payload itself is read through a
	// LimitReader, so a lying declared size cannot bypass the cap.
	reader, err := r.MultipartReader()
	if err != nil {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_UPLOAD", "expected a multipart form")
		return nil, false
	}
	part, err := reader.NextPart()
	if err != nil {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_UPLOAD", "expected a multipart file part named file")
		return nil, false
	}
	defer part.Close()
	if part.FormName() != "file" {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_UPLOAD", "expected a multipart file part named file")
		return nil, false
	}
	if part.FileName() == "" {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_FILE", "empty files are rejected")
		return nil, false
	}
	body, err := io.ReadAll(io.LimitReader(part, int64(s.opts.MaxBytes)+1))
	if err != nil {
		writeLocalizedError(w, r, http.StatusInternalServerError, "STORAGE_UNAVAILABLE", "could not read upload")
		return nil, false
	}
	if len(body) > s.opts.MaxBytes {
		writeLocalizedError(w, r, http.StatusRequestEntityTooLarge, "FILE_TOO_LARGE", "file exceeds the server size limit")
		return nil, false
	}
	if len(body) == 0 {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_FILE", "empty files are rejected")
		return nil, false
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
		return nil, false
	}
	processed, contentType, err := processRasterImage(body, target(s.opts), s.opts)
	if err != nil {
		writeLocalizedError(w, r, http.StatusUnsupportedMediaType, "UNSUPPORTED_FILE_TYPE", "image must be a decodable PNG, JPEG, GIF or WebP")
		return nil, false
	}
	id, err := s.save(contentType, kind, owner, processed)
	if err != nil {
		writeLocalizedError(w, r, http.StatusInternalServerError, "STORAGE_UNAVAILABLE", "could not store asset")
		return nil, false
	}
	return map[string]any{
		"id":   id,
		"name": part.FileName(),
		"type": contentType,
		"size": len(processed),
		"url":  s.urlPrefix + id,
	}, true
}

// upload serves a store upload without a permission gate (the route
// wrapper adds the surface gate).
func (s *RasterAssetStore) upload() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload, ok := s.storeUpload(w, r)
		if !ok {
			return
		}
		writeJSON(w, http.StatusOK, payload)
	})
}

// file serves GET {urlPrefix}{id} publicly. The store only ever contains
// server-produced raster output, so inline rendering is safe; the headers
// still pin type, nosniff and immutable caching (content-addressed).
func (s *RasterAssetStore) file() http.Handler {
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
func (s *RasterAssetStore) save(contentType, kind, owner string, body []byte) (string, error) {
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
	meta := map[string]string{"type": contentType, "kind": kind, "owner": owner}
	raw, err := json.Marshal(meta)
	if err == nil {
		_ = os.WriteFile(filepath.Join(s.dir, id+".meta.json"), raw, 0o644)
	}
	return id, nil
}

// load reads a stored asset by id (id shape validated — same guard as the
// generic upload store, so a crafted PathValue cannot escape the directory).
func (s *RasterAssetStore) load(id string) ([]byte, map[string]string, error) {
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
func (s *RasterAssetStore) Delete(id string) error {
	if !uploadFileIDPattern.MatchString(id) {
		return nil
	}
	if err := os.Remove(filepath.Join(s.dir, id)); err != nil && !os.IsNotExist(err) {
		return err
	}
	_ = os.Remove(filepath.Join(s.dir, id+".meta.json"))
	return nil
}

// DeleteOrphan deletes the asset referenced by a raw URL when (and only when)
// the URL points into this store (no-op otherwise). Best-effort cleanup used
// when a stored reference is replaced or cleared; callers log failures.
func (s *RasterAssetStore) DeleteOrphan(raw string) error {
	if id, ok := s.AssetIDFromURL(raw); ok {
		return s.Delete(id)
	}
	return nil
}


// CountOwner returns the number of stored assets whose meta marks them as owned
// by owner. Corrupt/unreadable meta files count conservatively toward every
// caller so a failed read cannot be used to bypass a per-user avatar quota.
func (s *RasterAssetStore) CountOwner(owner string) (int, error) {
	if owner == "" {
		return 0, nil
	}
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	count := 0
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".meta.json") {
			continue
		}
		id := strings.TrimSuffix(entry.Name(), ".meta.json")
		if !uploadFileIDPattern.MatchString(id) {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(s.dir, entry.Name()))
		if err != nil {
			count++ // conservative: unreadable meta may still be an owned asset
			continue
		}
		meta := map[string]string{}
		if err := json.Unmarshal(raw, &meta); err != nil {
			count++ // conservative
			continue
		}
		if meta["owner"] == owner {
			count++
		}
	}
	return count, nil
}

// DeleteAll removes every stored asset of this store.
func (s *RasterAssetStore) DeleteAll() error {
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

// GC deletes stored assets not referenced by the given URL list (startup
// housekeeping: crashed uploads / cancelled edits).
func (s *RasterAssetStore) GC(referenced []string) error {
	keep := map[string]bool{}
	for _, ref := range referenced {
		if id, ok := s.AssetIDFromURL(ref); ok {
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

// processRasterImage re-encodes a decoded raster: dimension-limited (never
// upscaled), PNG when the result has transparency, JPEG otherwise. Output is
// always a fresh server-produced raster — never the input bytes.
func processRasterImage(body []byte, targetDim int, opts BrandingAssetsOptions) ([]byte, string, error) {
	src, err := decodeRasterImage(body)
	if err != nil {
		return nil, "", err
	}
	bounds := src.Bounds()
	maxSide := bounds.Dx()
	if bounds.Dy() > maxSide {
		maxSide = bounds.Dy()
	}
	scale := 1.0
	if maxSide > targetDim {
		scale = float64(targetDim) / float64(maxSide)
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

// maxRasterInputDimension is the hard longest-edge bound for DECODED input
// (decompression-bomb guard): the output is at most 512/256/64px, so larger
// input only wastes memory; bounding the decoded allocation keeps the
// endpoint's per-request memory bounded.
const maxRasterInputDimension = 2048 // W7 F-005: 2048^2*4 ≈ 16 MiB worst decode allocation

// decodeRasterImage decodes PNG, JPEG, GIF (first frame) or WebP after a
// header-only DecodeConfig pre-check (tiny file + huge declared dimensions is
// rejected before any allocation). Any other content (including SVG/HTML
// smuggled past the sniff checks) fails decode.
func decodeRasterImage(body []byte) (image.Image, error) {
	cfg, err := rasterImageConfig(body)
	if err != nil {
		return nil, err
	}
	if cfg.Width > maxRasterInputDimension || cfg.Height > maxRasterInputDimension {
		return nil, errors.New("image dimensions exceed the server limit")
	}
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

// rasterImageConfig reads the header-only config for PNG/JPEG/GIF/WebP.
func rasterImageConfig(body []byte) (image.Config, error) {
	reader := func() *bytes.Reader { return bytes.NewReader(body) }
	if cfg, err := png.DecodeConfig(reader()); err == nil {
		return cfg, nil
	}
	if cfg, err := jpeg.DecodeConfig(reader()); err == nil {
		return cfg, nil
	}
	if cfg, err := gif.DecodeConfig(reader()); err == nil {
		return cfg, nil
	}
	if cfg, err := webp.DecodeConfig(reader()); err == nil {
		return cfg, nil
	}
	return image.Config{}, errors.New("unsupported image format")
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
