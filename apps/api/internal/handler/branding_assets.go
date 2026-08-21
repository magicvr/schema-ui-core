// Brand asset upload endpoints (W9 / GOAL-010).
//
// Brand icons (logo / light / dark / favicon) are configured on the Settings
// page. The legacy textarea-URL input is replaced by uploads: images are
// stored in a dedicated brand-assets directory — NOT the generic /api/upload
// owner store and NOT the admin.file-library module — and served publicly,
// because the login page and shell load branding before authentication.
//
// The store implementation is shared with the account avatar upload (W13
// T-05): see raster_assets.go for the generic RasterAssetStore.
//
// Security model:
//   - uploads require settings.write (module permission gate)
//   - every image is re-encoded server-side (PNG/JPEG/GIF/WebP -> PNG or
//     JPEG, dimension-limited); raw user bytes are never stored or served
//   - the public GET only ever serves server-produced raster output with
//     nosniff + sandbox + immutable caching
package handler

import (
	"net/http"

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
	// AvatarDim is the longest-edge limit for account avatars (default 256).
	AvatarDim int
	// JPEGQuality (1..100) for opaque output (default 82).
	JPEGQuality int
}

// DefaultBrandingAssetsOptions returns the documented defaults.
func DefaultBrandingAssetsOptions() BrandingAssetsOptions {
	return BrandingAssetsOptions{
		MaxBytes:    4 << 20,
		LogoMaxDim:  512,
		FaviconDim:  64,
		AvatarDim:   256,
		JPEGQuality: 82,
	}
}

// BrandingAssetRoutes returns the admin.settings brand asset surface:
// authenticated upload (settings.write gate) + public GET.
func BrandingAssetRoutes(a authMiddleware, store *BrandingAssetStore, moduleID string) []kernel.RouteContribution {
	identity := func(method, pattern string) kernel.ContributionIdentity {
		return kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey(method, pattern)}
	}
	return []kernel.RouteContribution{
		{ContributionIdentity: identity("POST", "/api/branding/assets"), Method: "POST", Pattern: "/api/branding/assets", Handler: a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := requirePermission(w, r, "settings.write"); !ok {
				return
			}
			store.upload().ServeHTTP(w, r)
		}))},
		{ContributionIdentity: identity("GET", "/api/branding/assets/{id}"), Method: "GET", Pattern: "/api/branding/assets/{id}", Handler: store.file(), Public: true},
	}
}

// RegisterPublicBrandingAssets mounts the public GET for profiles without
// admin.settings (mvp): /api/branding may reference previously uploaded
// assets, so the read surface must stay available. Do not double-register
// when the settings module contributes the same route.
func RegisterPublicBrandingAssets(mux routeRegistrar, store *BrandingAssetStore) {
	mux.Handle("GET /api/branding/assets/{id}", store.file())
}
