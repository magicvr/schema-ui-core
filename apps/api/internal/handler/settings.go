// Site branding settings API (GOAL-013): public GET for shell/login branding,
// authenticated list + PATCH for the Schema-driven Settings page.
package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

// RegisterSettings exposes the Settings module registration adapter to the
// composition root. The module package owns when this contribution is added.
func RegisterSettings(mux *http.ServeMux, a *auth.Authenticator, st *store.Store) {
	settingsHandler(mux, a, st)
}

// SettingsRoutes returns the Settings module HTTP route contributions (public
// branding + authenticated settings list/detail/patch). R4 C4.1: module
// providers reuse it so the provider surface matches the central adapter.
func SettingsRoutes(a *auth.Authenticator, st *store.Store, moduleID string) []kernel.RouteContribution {
	return []kernel.RouteContribution{
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("GET", "/api/branding")}, Method: "GET", Pattern: "/api/branding", Handler: brandingGET(st)},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("GET", "/api/settings")}, Method: "GET", Pattern: "/api/settings", Handler: a.Middleware(settingsList(st))},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("GET", "/api/settings/{id}")}, Method: "GET", Pattern: "/api/settings/{id}", Handler: a.Middleware(settingsDetail(st))},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("PATCH", "/api/settings/{id}")}, Method: "PATCH", Pattern: "/api/settings/{id}", Handler: a.Middleware(settingsPatch(st))},
	}
}

func settingsHandler(mux *http.ServeMux, a *auth.Authenticator, st *store.Store) {
	// Public branding for login shell and document title (no secrets).
	mux.HandleFunc("GET /api/branding", brandingGET(st))
	// Schema table list envelope (one row) — requires settings.read.
	mux.Handle("GET /api/settings", a.Middleware(settingsList(st)))
	mux.Handle("GET /api/settings/{id}", a.Middleware(settingsDetail(st)))
	mux.Handle("PATCH /api/settings/{id}", a.Middleware(settingsPatch(st)))
}

type brandingResponse struct {
	SiteTitle string `json:"siteTitle"`
	LogoURL   string `json:"logoUrl"`
}

const configChangedHeader = "X-Schema-UI-Config-Changed"

func brandingGET(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		s, err := st.GetSiteSettings()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL", "could not load branding")
			return
		}
		writeJSON(w, http.StatusOK, brandingResponse{SiteTitle: s.SiteTitle, LogoURL: s.LogoURL})
	}
}

func settingsRow(s *store.SiteSettings) map[string]any {
	return map[string]any{
		"id":        s.ID,
		"siteTitle": s.SiteTitle,
		"logoUrl":   s.LogoURL,
		"updatedAt": s.UpdatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
	}
}

func settingsList(st *store.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r.Context(), "settings.read"); !ok {
			return
		}
		s, err := st.GetSiteSettings()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL", "could not list settings")
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

func settingsDetail(st *store.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r.Context(), "settings.read"); !ok {
			return
		}
		id := r.PathValue("id")
		if id != "default" {
			writeError(w, http.StatusNotFound, "SETTINGS_NOT_FOUND", "no settings with that id")
			return
		}
		s, err := st.GetSiteSettings()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL", "could not load settings")
			return
		}
		writeJSON(w, http.StatusOK, settingsRow(s))
	})
}

func settingsPatch(st *store.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r.Context(), "settings.write")
		if !ok {
			return
		}
		id := r.PathValue("id")
		if id != "default" {
			writeError(w, http.StatusNotFound, "SETTINGS_NOT_FOUND", "no settings with that id")
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxResourceBodyBytes)
		var body struct {
			SiteTitle *string `json:"siteTitle"`
			LogoURL   *string `json:"logoUrl"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_PATCH_BODY", "body must be JSON")
			return
		}
		cur, err := st.GetSiteSettings()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL", "could not load settings")
			return
		}
		title := cur.SiteTitle
		logo := cur.LogoURL
		if body.SiteTitle != nil {
			title = *body.SiteTitle
		}
		if body.LogoURL != nil {
			logo = *body.LogoURL
		}
		now := time.Now().UTC()
		updated, err := st.UpdateSiteSettings(title, logo, now)
		if errors.Is(err, store.ErrInvalidSiteTitle) {
			writeError(w, http.StatusBadRequest, "INVALID_SITE_TITLE", "siteTitle must not be empty")
			return
		}
		if errors.Is(err, store.ErrInvalidLogoURL) {
			writeError(w, http.StatusBadRequest, "INVALID_LOGO_URL", "logoUrl must be empty, a same-origin path, or an http(s) URL")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL", "could not update settings")
			return
		}
		// Best-effort operation log (A-006 R-003); never fails the write.
		recordID := "default"
		detail := `{"siteTitle":` + jsonQuote(updated.SiteTitle) + `}`
		op := store.Operation{
			ID:        newOperationID(),
			Event:     store.EventSettingsUpdate,
			ActorID:   user.ID,
			ActorName: user.Name,
			RecordID:  &recordID,
			Detail:    &detail,
			CreatedAt: now,
		}
		if err := st.RecordOperation(op); err != nil {
			slog.Error("operation log write failed", "event", store.EventSettingsUpdate, "err", err)
		}
		w.Header().Set(configChangedHeader, "settings.branding")
		writeJSON(w, http.StatusOK, settingsRow(updated))
	})
}
