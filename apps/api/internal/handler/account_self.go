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
}

// AccountSelfRoutes returns the self-service route contributions (admin.account).
func AccountSelfRoutes(a *auth.Authenticator, repository AccountRepository, operations operationlog.Recorder, moduleID string) []kernel.RouteContribution {
	h := &accountSelfHandler{
		repository:      repository,
		operations:      operations,
		now:             time.Now,
		passwordLimiter: newLoginRateLimiter(15*time.Minute, 5, 1<<16),
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
	return routes
}

type accountSelfHandler struct {
	repository AccountRepository
	operations operationlog.Recorder
	now        func() time.Time
	// F-003 (A-003 recommended): wrong currentPassword attempts are brute-force
	// surface with a live access token — an in-memory sliding-window limiter per
	// client identity (same model as login rate limiting) brakes online
	// guessing. Best-effort and process-local, like the login limiter.
	passwordLimiter *loginRateLimiter
}

func (h *accountSelfHandler) identity(w http.ResponseWriter, r *http.Request) (account.User, bool) {
	user, ok := auth.IdentityFrom(r.Context())
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
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PATCH_FIELD", "name must not be empty")
			return
		}
		u, err := h.repository.UpdateUser(user.ID, authsession.UserPatch{Name: &name}, user.ID, h.now().UTC())
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not update profile")
			return
		}
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
		u, err := h.repository.GetUser(user.ID)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not load account")
			return
		}
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
		if _, err := h.repository.UpdateUser(user.ID, authsession.UserPatch{PasswordHash: &hash}, user.ID, h.now().UTC()); err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not change password")
			return
		}
		h.record(operationlog.EventAccountPasswordChange, user.ID, user.Name, user.ID, "")
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
		pageSize, ok := intParam(r.URL.Query().Get("pageSize"), 10)
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
		start := (page - 1) * pageSize
		end := start + pageSize
		if start > len(all) {
			start = len(all)
		}
		if end > len(all) {
			end = len(all)
		}
		items := make([]map[string]any, 0, end-start)
		for _, token := range all[start:end] {
			status := "active"
			if token.RevokedAt != nil {
				status = "revoked"
			}
			row := map[string]any{
				"id":        token.ID,
				"createdAt": token.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
				"expiresAt": token.ExpiresAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
				"status":    status,
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
		h.record(operationlog.EventAccountSessionRevoke, user.ID, user.Name, id, "")
		w.WriteHeader(http.StatusNoContent)
	})
}

func (h *accountSelfHandler) record(event, actorID, actorName, recordID, detail string) {
	op := operationlog.Operation{
		ID:        newOperationID(),
		Event:     event,
		ActorID:   actorID,
		ActorName: actorName,
		CreatedAt: h.now().UTC(),
	}
	if recordID != "" {
		op.RecordID = &recordID
	}
	if detail != "" {
		op.Detail = &detail
	}
	if h.operations == nil {
		return
	}
	if err := h.operations.RecordOperation(op); err != nil {
		slog.Error("operation log write failed", "event", event, "err", err)
	}
}