// Site branding settings API (GOAL-013 + VP-007 S3): public GET for shell/login
// startup configuration, authenticated list + PATCH for the Schema-driven
// Settings page, plus the authenticated reset-to-defaults endpoint.
package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	settingsrepository "github.com/magicvr/schema-ui-core/apps/api/internal/modules/settings/repository"
)

// RegisterSettings is removed in R6 C6.1: the Settings module mounts its HTTP
// surface via the module provider (SettingsRoutes + kernel.RegisterContributions);
// the handler test environment mounts SettingsRoutes directly. No central
// adapter remains.

// SettingsRoutes returns the Settings module HTTP route contributions (public
// startup config + authenticated settings list/detail/patch/reset). R4 C4.1:
// module providers reuse it so the provider surface matches the central adapter.
type SettingsRepository interface {
	GetSiteSettings() (*settingsrepository.SiteSettings, error)
	PatchSiteSettings(*string, *string, *string, *string, *string, *string, *string, *string, time.Time) (*settingsrepository.SiteSettings, error)
	ResetSiteSettings(time.Time) (*settingsrepository.SiteSettings, error)
}

func SettingsRoutes(a *auth.Authenticator, repository SettingsRepository, operations operationlog.Recorder, moduleID, configNamespace string, assets *BrandingAssetStore) []kernel.RouteContribution {
	return []kernel.RouteContribution{
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("GET", "/api/branding")}, Method: "GET", Pattern: "/api/branding", Handler: brandingGET(repository), Public: true},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("GET", "/api/settings")}, Method: "GET", Pattern: "/api/settings", Handler: a.Middleware(settingsList(repository))},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("GET", "/api/settings/{id}")}, Method: "GET", Pattern: "/api/settings/{id}", Handler: a.Middleware(settingsDetail(repository))},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("PATCH", "/api/settings/{id}")}, Method: "PATCH", Pattern: "/api/settings/{id}", Handler: a.Middleware(settingsPatch(repository, operations, configNamespace, assets))},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("POST", "/api/settings/{id}/reset")}, Method: "POST", Pattern: "/api/settings/{id}/reset", Handler: a.Middleware(settingsReset(repository, operations, configNamespace, assets))},
	}
}

// RegisterPublicBranding mounts GET /api/branding for profiles that do not load
// admin.settings (VP-007 exit 4 / I-L10N-003: public startup config must not
// require the settings edit module). When admin.settings is present, branding
// is contributed by SettingsRoutes instead — do not double-register.
func RegisterPublicBranding(mux *http.ServeMux, repository SettingsRepository) {
	mux.Handle("GET /api/branding", brandingGET(repository))
}

// brandingResponse is the public startup configuration (I-L10N-003 compatible
// extension of the legacy {siteTitle, logoUrl} contract; additive fields only).
type brandingResponse struct {
	SiteTitle        string   `json:"siteTitle"`
	LogoURL          string   `json:"logoUrl"`
	LogoURLLight     string   `json:"logoUrlLight"`
	LogoURLDark      string   `json:"logoUrlDark"`
	FaviconURL       string   `json:"faviconUrl"`
	DefaultLocale    string   `json:"defaultLocale"`
	SupportedLocales []string `json:"supportedLocales"`
	SiteTimezone     string   `json:"siteTimezone"`
	DefaultTheme     string   `json:"defaultTheme"`
}

const configChangedHeader = "X-Schema-UI-Config-Changed"

func brandingGET(repository SettingsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, err := repository.GetSiteSettings()
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not load branding")
			return
		}
		writeJSON(w, http.StatusOK, brandingRow(s))
	}
}

func brandingRow(s *settingsrepository.SiteSettings) brandingResponse {
	locale := s.DefaultLocale
	if locale == "" {
		locale = "auto"
	}
	timezone := s.SiteTimezone
	if timezone == "" {
		timezone = "auto"
	}
	theme := s.DefaultTheme
	if theme == "" {
		theme = "auto"
	}
	return brandingResponse{
		SiteTitle:        s.SiteTitle,
		LogoURL:          s.LogoURL,
		LogoURLLight:     s.LogoURLLight,
		LogoURLDark:      s.LogoURLDark,
		FaviconURL:       s.FaviconURL,
		DefaultLocale:    locale,
		SupportedLocales: settingsrepository.SupportedLocales,
		SiteTimezone:     timezone,
		DefaultTheme:     theme,
	}
}

func settingsRow(s *settingsrepository.SiteSettings) map[string]any {
	return map[string]any{
		"id":            s.ID,
		"siteTitle":     s.SiteTitle,
		"logoUrl":       s.LogoURL,
		"logoUrlLight":  s.LogoURLLight,
		"logoUrlDark":   s.LogoURLDark,
		"faviconUrl":    s.FaviconURL,
		"defaultLocale": s.DefaultLocale,
		"siteTimezone":  s.SiteTimezone,
		"defaultTheme":  s.DefaultTheme,
		"updatedAt":     s.UpdatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
	}
}

func settingsList(repository SettingsRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r, "settings.read"); !ok {
			return
		}
		s, err := repository.GetSiteSettings()
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not list settings")
			return
		}
		writeJSON(w, http.StatusOK, resourceList{
			Items:    []map[string]any{settingsRow(s)},
			Total:    1,
			Page:     1,
			PageSize: 10,
		})
	})
}

func settingsDetail(repository SettingsRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r, "settings.read"); !ok {
			return
		}
		id := r.PathValue("id")
		if id != "default" {
			writeLocalizedError(w, r, http.StatusNotFound, "SETTINGS_NOT_FOUND", "no settings with that id")
			return
		}
		s, err := repository.GetSiteSettings()
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not load settings")
			return
		}
		writeJSON(w, http.StatusOK, settingsRow(s))
	})
}

func settingsPatch(repository SettingsRepository, operations operationlog.Recorder, configNamespace string, assets *BrandingAssetStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, "settings.write")
		if !ok {
			return
		}
		id := r.PathValue("id")
		if id != "default" {
			writeLocalizedError(w, r, http.StatusNotFound, "SETTINGS_NOT_FOUND", "no settings with that id")
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxResourceBodyBytes)
		var body struct {
			SiteTitle     *string `json:"siteTitle"`
			LogoURL       *string `json:"logoUrl"`
			LogoURLLight  *string `json:"logoUrlLight"`
			LogoURLDark   *string `json:"logoUrlDark"`
			FaviconURL    *string `json:"faviconUrl"`
			DefaultLocale *string `json:"defaultLocale"`
			SiteTimezone  *string `json:"siteTimezone"`
			DefaultTheme  *string `json:"defaultTheme"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PATCH_BODY", "body must be JSON")
			return
		}
		// W9 (GOAL-010): capture the pre-patch values so replaced brand assets
		// can be deleted right after a successful patch (best-effort; the
		// startup GC catches leftovers).
		current, _ := repository.GetSiteSettings()
		now := time.Now().UTC()
		updated, err := repository.PatchSiteSettings(
			body.SiteTitle, body.LogoURL, body.LogoURLLight, body.LogoURLDark, body.FaviconURL,
			body.DefaultLocale, body.SiteTimezone, body.DefaultTheme, now,
		)
		if err != nil {
			writeSettingsError(w, r, err)
			return
		}
		if current != nil {
			cleanupReplacedBrandAssets(current, updated, assets)
		}
		recordSettingsOperation(operations, user, "updated", updated, now)
		w.Header().Set(configChangedHeader, configNamespace)
		writeJSON(w, http.StatusOK, settingsRow(updated))
	})
}

// cleanupReplacedBrandAssets deletes brand assets that the patch replaced or
// cleared (I-004: replace/clear deletes immediately). Only assets previously
// referenced by a branding column are touched; legacy URL values are skipped.
func cleanupReplacedBrandAssets(current, updated *settingsrepository.SiteSettings, assets *BrandingAssetStore) {
	if assets == nil {
		return
	}
	fields := []struct{ before, after string }{
		{current.LogoURL, updated.LogoURL},
		{current.LogoURLLight, updated.LogoURLLight},
		{current.LogoURLDark, updated.LogoURLDark},
		{current.FaviconURL, updated.FaviconURL},
	}
	for _, f := range fields {
		if f.before == f.after {
			continue
		}
		if id, ok := BrandAssetIDFromURL(f.before); ok {
			_ = assets.Delete(id)
		}
	}
}

func settingsReset(repository SettingsRepository, operations operationlog.Recorder, configNamespace string, assets *BrandingAssetStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, "settings.write")
		if !ok {
			return
		}
		id := r.PathValue("id")
		if id != "default" {
			writeLocalizedError(w, r, http.StatusNotFound, "SETTINGS_NOT_FOUND", "no settings with that id")
			return
		}
		now := time.Now().UTC()
		updated, err := repository.ResetSiteSettings(now)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not reset settings")
			return
		}
		// I-004: reset clears every stored brand asset.
		if assets != nil {
			_ = assets.DeleteAll()
		}
		recordSettingsOperation(operations, user, "reset", updated, now)
		w.Header().Set(configChangedHeader, configNamespace)
		writeJSON(w, http.StatusOK, settingsRow(updated))
	})
}

// writeSettingsError maps repository validation errors onto the stable error
// code contract (S3; codes are part of the D-002 appendix A enumeration).
func writeSettingsError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, settingsrepository.ErrInvalidSiteTitle):
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SITE_TITLE", "siteTitle must not be empty")
	case errors.Is(err, settingsrepository.ErrInvalidLogoURL):
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_LOGO_URL", "logoUrl fields must be empty, a same-origin path, or an http(s) URL")
	case errors.Is(err, settingsrepository.ErrInvalidDefaultLocale):
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_DEFAULT_LOCALE", "defaultLocale must be auto, zh-CN or en-US")
	case errors.Is(err, settingsrepository.ErrInvalidDefaultTheme):
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_DEFAULT_THEME", "defaultTheme must be auto, light or dark")
	case errors.Is(err, settingsrepository.ErrInvalidSiteTimezone):
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_TIMEZONE", "siteTimezone must be auto or a valid IANA timezone")
	default:
		writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not update settings")
	}
}

func recordSettingsOperation(operations operationlog.Recorder, user account.User, action string, updated *settingsrepository.SiteSettings, now time.Time) {
	if operations == nil {
		return
	}
	recordID := "default"
	detail := `{"siteTitle":` + jsonQuote(updated.SiteTitle) + `,"action":"` + action + `"}`
	op := operationlog.Operation{
		ID:        newOperationID(),
		Event:     operationlog.EventSettingsUpdate,
		ActorID:   user.ID,
		ActorName: user.Name,
		RecordID:  &recordID,
		Detail:    &detail,
		CreatedAt: now,
	}
	if err := operations.RecordOperation(op); err != nil {
		slog.Error("operation log write failed", "event", operationlog.EventSettingsUpdate, "err", err)
	}
}
