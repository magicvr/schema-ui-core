// F-03 self-service account surface (GOAL-005 D-002 `4): profile, own
// password change, and session listing/revocation. Identity-only endpoints
// (no permission key — self-service is available to every authenticated
// user). Password change reuses the management password semantics
// (8-72 bytes) and the W4 P0-3 token_version/refresh-revocation side effects.
package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/errorcatalog"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

// AccountRepository is the persistence surface consumed by the self-service
// account endpoints (admin.account module).
type AccountRepository interface {
	GetUser(string) (*authsession.User, error)
	UpdateUser(string, authsession.UserPatch, string, time.Time) (*authsession.User, error)
	ListRefreshTokensForUser(string) ([]authsession.RefreshToken, error)
	RevokeRefreshTokenIfOwned(string, string, time.Time) error
	BumpTokenVersionAndRevokeAll(string, time.Time) error
}

// AccountSelfRoutes returns the self-service route contributions (admin.account).
// avatarStore is the account avatar asset store (W13 T-05); nil disables the
// avatarUrl profile surface (used by bare test environments).
func AccountSelfRoutes(a *auth.Authenticator, repository AccountRepository, operations operationlog.Recorder, avatarStore *RasterAssetStore, moduleID string, notifier ...NotifyRepository) []kernel.RouteContribution {
	h := &accountSelfHandler{
		auth:            a,
		repository:      repository,
		operations:      operations,
		avatarStore:     avatarStore,
		now:             time.Now,
		passwordLimiter: newLoginRateLimiter(15*time.Minute, 5, 1<<16),
	}
	if len(notifier) > 0 {
		h.notifier = notifier[0]
	}
	var routes []kernel.RouteContribution
	add := func(method, pattern string, handler http.Handler) {
		routes = append(routes, kernel.RouteContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey(method, pattern)},
			Method:               method,
			Pattern:              pattern,
			Handler:              handler,
		})
	}
	add("GET", "/api/account/profile", a.Middleware(h.profile()))
	add("PATCH", "/api/account/profile", a.Middleware(h.updateProfile()))
	add("POST", "/api/account/password", a.Middleware(h.changePassword()))
	add("GET", "/api/account/sessions", a.Middleware(h.sessions()))
	add("POST", "/api/account/sessions/{id}/revoke", a.Middleware(h.revokeSession()))
	add("POST", "/api/account/sessions/revoke-others", a.Middleware(h.revokeOthers()))
	return routes
}

type accountSelfHandler struct {
	auth       *auth.Authenticator
	repository AccountRepository
	operations operationlog.Recorder
	// avatarStore owns the user avatar assets (W13 T-05); nil = no avatar surface.
	avatarStore *RasterAssetStore
	now         func() time.Time
	notifier    NotifyRepository
	// F-003 (A-003 recommended): wrong currentPassword attempts are brute-force
	// surface with a live access token — an in-memory sliding-window limiter per
	// client identity (same model as login rate limiting) brakes online
	// guessing. Best-effort and process-local, like the login limiter.
	passwordLimiter *loginRateLimiter
}

func (h *accountSelfHandler) identity(w http.ResponseWriter, r *http.Request) (account.User, bool) {
	user, ok := auth.UserIdentityFrom(r.Context())
	if !ok {
		writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "no active session")
		return account.User{}, false
	}
	return user, true
}

func accountProfileRow(u *authsession.User) map[string]any {
	return map[string]any{
		"id":        u.ID,
		"username":  u.Username,
		"name":      u.Name,
		"avatarUrl": u.AvatarURL,
		"enabled":   u.Enabled,
		"createdAt": u.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
		"updatedAt": u.UpdatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
	}
}

func (h *accountSelfHandler) profile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := h.identity(w, r)
		if !ok {
			return
		}
		u, err := h.repository.GetUser(user.ID)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not load profile")
			return
		}
		writeJSON(w, http.StatusOK, accountProfileRow(u))
	})
}

func (h *accountSelfHandler) updateProfile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := h.identity(w, r)
		if !ok {
			return
		}
		var body map[string]json.RawMessage
		r.Body = http.MaxBytesReader(w, r.Body, maxResourceBodyBytes)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PATCH_BODY", "body must be JSON")
			return
		}
		raw, present := body["name"]
		if !present {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PATCH_BODY", "name is required")
			return
		}
		var name string
		if err := json.Unmarshal(raw, &name); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PATCH_BODY", "name must be a string")
			return
		}
		name = strings.TrimSpace(name)
		if name == "" {
			writeLocalizedFieldError(w, r, http.StatusBadRequest, "INVALID_PATCH_FIELD", "name must not be empty", []errorcatalog.FieldError{{Field: "name", Reason: "must not be empty"}})
			return
		}
		patch := authsession.UserPatch{Name: &name}
		// W13 T-05: optional avatarUrl ("" clears the avatar; a value must be a
		// URL served by the account avatar store — nothing else may be
		// committed to the profile). The previous avatar file is deleted
		// best-effort once the new value is persisted.
		if rawAvatar, hasAvatar := body["avatarUrl"]; hasAvatar {
			if h.avatarStore == nil {
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PATCH_BODY", "avatarUrl is not supported")
				return
			}
			var avatarURL string
			if err := json.Unmarshal(rawAvatar, &avatarURL); err != nil {
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PATCH_BODY", "avatarUrl must be a string")
				return
			}
			avatarURL = strings.TrimSpace(avatarURL)
			if avatarURL != "" {
				if _, ok := h.avatarStore.AssetIDFromURL(avatarURL); !ok {
					writeLocalizedFieldError(w, r, http.StatusBadRequest, "INVALID_PATCH_FIELD", "avatarUrl is not a valid avatar asset", []errorcatalog.FieldError{{Field: "avatarUrl", Reason: "must be a URL of the account avatar store"}})
					return
				}
			}
			patch.AvatarURL = &avatarURL
		}
		// Capture the previous avatar before the update so the replaced or
		// cleared file can be dropped afterwards (best-effort).
		oldAvatar := ""
		if patch.AvatarURL != nil && h.avatarStore != nil {
			if current, err := h.repository.GetUser(user.ID); err == nil {
				oldAvatar = current.AvatarURL
			}
		}
		u, err := h.repository.UpdateUser(user.ID, patch, user.ID, h.now().UTC())
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not update profile")
			return
		}
		if patch.AvatarURL != nil && h.avatarStore != nil && oldAvatar != *patch.AvatarURL {
			// Only avatar-store assets are ever deleted here (the new value
			// was validated above; the old value is either empty, the same
			// asset, or an avatar-store asset from an earlier commit).
			if err := h.avatarStore.DeleteOrphan(oldAvatar); err != nil {
				slog.Error("avatar clear cleanup failed", "user", user.ID, "err", err)
			}
		}
		// W13 T-05 follow-up (user 2026-08-16): the shell session (user-menu
		// avatar + display name) must refresh immediately after a profile save
		// — same config-change channel the settings branding surface uses. The
		// host reacts by re-resolving /me; no page reload required.
		w.Header().Set(configChangedHeader, "account.profile")
		writeJSON(w, http.StatusOK, accountProfileRow(u))
	})
}

type changePasswordBody struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

func (h *accountSelfHandler) changePassword() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := h.identity(w, r)
		if !ok {
			return
		}
		var body changePasswordBody
		r.Body = http.MaxBytesReader(w, r.Body, maxResourceBodyBytes)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PASSWORD_BODY", "body must be JSON with currentPassword and newPassword")
			return
		}
		if body.CurrentPassword == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PASSWORD", "currentPassword is required")
			return
		}
		length := len([]byte(body.NewPassword))
		if length < minPasswordBytes || length > maxPasswordBytes || strings.TrimSpace(body.NewPassword) == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PASSWORD", "new password must be a non-whitespace string of 8 to 72 bytes")
			return
		}
		if body.NewPassword == body.CurrentPassword {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PASSWORD", "new password must differ from the current password")
			return
		}
		u, err := h.repository.GetUser(user.ID)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not load account")
			return
		}
		wasForced := u.MustChangePassword
		now := h.now().UTC()
		// F-003: wrong-current-password attempts share the login limiter model
		// (per client identity, 5 failures / 15 min window). The limiter is
		// checked before the bcrypt work; a blocked client gets 429.
		limiterKey := loginClientIP(r) + "|" + user.ID
		if h.passwordLimiter != nil && !h.passwordLimiter.allow(limiterKey, h.now().UTC()) {
			writeLocalizedError(w, r, http.StatusTooManyRequests, "RATE_LIMITED", "too many failed password attempts; try again later")
			return
		}
		if !auth.VerifyPassword(u.PasswordHash, body.CurrentPassword) {
			if h.passwordLimiter != nil {
				h.passwordLimiter.record(limiterKey, h.now().UTC())
			}
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PASSWORD", "current password is incorrect")
			return
		}
		if h.passwordLimiter != nil {
			h.passwordLimiter.clear(limiterKey)
		}
		hash, err := auth.HashPassword(body.NewPassword, passwordHashCost)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not hash password")
			return
		}
		// Reuses the management update path: bumps token_version and revokes
		// every live refresh token atomically (W4 P0-3; GOAL-005 D-002 `2).
		// W16-F01: a successful self-service password change clears the forced
		// must-change flag.
		mustChange := false
		if _, err := h.repository.UpdateUser(user.ID, authsession.UserPatch{PasswordHash: &hash, MustChangePassword: &mustChange}, user.ID, now); err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not change password")
			return
		}
		h.record(operationlog.EventAccountPasswordChange, user, user.ID, "")
		NotifyAccountEvent(h.notifier, user.ID, "account.password-changed", now)
		// W16-F01: when the change was a forced initial-password replacement,
		// reissue a fresh token pair so the user stays signed in and enters the
		// app immediately. Normal password changes keep the historical 204 +
		// re-login contract.
		if wasForced {
			access, refresh, acct, err := h.auth.IssueTokensFor(user.ID, now)
			if err != nil {
				writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not refresh current session")
				return
			}
			writeJSON(w, http.StatusOK, tokenResponse{AccessToken: access, RefreshToken: refresh, User: acct})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

func (h *accountSelfHandler) sessions() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := h.identity(w, r)
		if !ok {
			return
		}
		page, ok := intParam(r.URL.Query().Get("page"), 1)
		if !ok {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PAGE", "page must be a positive integer")
			return
		}
		pageSize, ok := intParam(r.URL.Query().Get("pageSize"), DefaultPageSize)
		if !ok {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PAGE_SIZE", "pageSize must be a positive integer")
			return
		}
		if pageSize > maxPageSize {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PAGE_SIZE", "pageSize must not exceed 100")
			return
		}
		all, err := h.repository.ListRefreshTokensForUser(user.ID)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not list sessions")
			return
		}
		// Optional status filter (schema-driven: sessions-table filters.status).
		switch statusFilter := strings.TrimSpace(r.URL.Query().Get("status")); statusFilter {
		case "":
			// no filter
		case "active", "revoked":
			filtered := make([]authsession.RefreshToken, 0, len(all))
			for _, token := range all {
				revoked := token.RevokedAt != nil
				if (statusFilter == "active" && !revoked) || (statusFilter == "revoked" && revoked) {
					filtered = append(filtered, token)
				}
			}
			all = filtered
		default:
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_STATUS_FILTER", "status must be active, revoked, or empty")
			return
		}
		start := (page - 1) * pageSize
		end := start + pageSize
		if start > len(all) {
			start = len(all)
		}
		if end > len(all) {
			end = len(all)
		}
		currentHash := ""
		if raw := strings.TrimSpace(r.Header.Get("X-Refresh-Token")); raw != "" {
			currentHash = auth.HashToken(raw)
		}
		items := make([]map[string]any, 0, end-start)
		for _, token := range all[start:end] {
			status := "active"
			if token.RevokedAt != nil {
				status = "revoked"
			}
			current := currentHash != "" && token.TokenHash == currentHash
			row := map[string]any{
				"id":        token.ID,
				"createdAt": token.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
				"expiresAt": token.ExpiresAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
				"status":    status,
				"current":   current,
			}
			if current {
				row["userAgent"] = r.UserAgent()
				row["ip"] = loginClientIP(r)
			}
			if token.RevokedAt != nil {
				row["revokedAt"] = token.RevokedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00")
			}
			items = append(items, row)
		}
		writeJSON(w, http.StatusOK, resourceList{Items: items, Total: len(all), Page: page, PageSize: pageSize})
	})
}

func (h *accountSelfHandler) revokeSession() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := h.identity(w, r)
		if !ok {
			return
		}
		id := r.PathValue("id")
		if err := h.repository.RevokeRefreshTokenIfOwned(id, user.ID, h.now().UTC()); err != nil {
			if errors.Is(err, authsession.ErrSessionNotFound) {
				writeLocalizedError(w, r, http.StatusNotFound, "SESSION_NOT_FOUND", "no session with that id")
				return
			}
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not revoke session")
			return
		}
		h.record(operationlog.EventAccountSessionRevoke, user, id, "")
		w.WriteHeader(http.StatusNoContent)
	})
}

func (h *accountSelfHandler) revokeOthers() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := h.identity(w, r)
		if !ok {
			return
		}
		now := h.now().UTC()
		// W16-F07: revoke every other device's refresh token and invalidate all
		// previously issued access tokens by bumping token_version, then mint a
		// fresh pair for the current caller so this device stays signed in.
		if err := h.repository.BumpTokenVersionAndRevokeAll(user.ID, now); err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not revoke other sessions")
			return
		}
		access, refresh, acct, err := h.auth.IssueTokensFor(user.ID, now)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not refresh current session")
			return
		}
		h.record(operationlog.EventAccountSessionRevoke, user, "others", "")
		writeJSON(w, http.StatusOK, tokenResponse{AccessToken: access, RefreshToken: refresh, User: acct})
	})
}

func (h *accountSelfHandler) record(event string, user account.User, recordID, detail string) {
	var detailPtr *string
	if strings.TrimSpace(detail) != "" {
		detailPtr = auditDetail(strings.TrimPrefix(event, "account."), map[string]any{"detail": detail})
	}
	recordAudit(h.operations, user, event, recordID, detailPtr, h.now().UTC(), nil)
}
