// Password-policy configuration surface (workspace-019 R3 · GOAL-004
// D-001 §2): the admin.settings tab extension — GET/PATCH the singleton row,
// gated by settings.write. Range validation lives here so the stored policy
// always stays inside the bcrypt-safe window.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
)

// PolicySettingsRepository is what the configuration endpoints need.
type PolicySettingsRepository interface {
	GetPasswordPolicy() (authsession.PasswordPolicy, error)
	UpdatePasswordPolicy(p authsession.PasswordPolicy) error
	PermissionsForUser(userID string) ([]string, error)
}

// PasswordPolicyRoutes returns the GET/PATCH pair for the admin.settings
// module (D-001 §2: UI 仅 admin.settings tab 扩展；不进 mvp 默认集).
func PasswordPolicyRoutes(a *auth.Authenticator, repo PolicySettingsRepository, moduleID string) []kernel.RouteContribution {
	h := &policySettingsHandler{a: a, repo: repo, now: time.Now}
	var routes []kernel.RouteContribution
	add := func(method, pattern string, handler http.Handler) {
		routes = append(routes, kernel.RouteContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey(method, pattern)},
			Method:               method,
			Pattern:              pattern,
			Handler:              handler,
		})
	}
	add("GET", "/api/settings/password-policy", a.Middleware(h.requirePermission("settings.read", h.get())))
	add("PATCH", "/api/settings/password-policy", a.Middleware(h.requirePermission("settings.write", h.patch())))
	return routes
}

type policySettingsHandler struct {
	a    *auth.Authenticator
	repo PolicySettingsRepository
	now  func() time.Time
}

// requirePermission gates on ONE resolved permission key; the read face takes
// settings.read while only the write face demands settings.write (A-001 F-004,
// mirroring the mail-tab read/write split).
func (h *policySettingsHandler) requirePermission(permission string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actor, ok := auth.UserIdentityFrom(r.Context())
		if !ok {
			writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "no active session")
			return
		}
		granted, err := h.repo.PermissionsForUser(actor.ID)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not resolve permissions")
			return
		}
		for _, p := range granted {
			if p == permission {
				next.ServeHTTP(w, r)
				return
			}
		}
		writeLocalizedError(w, r, http.StatusForbidden, "FORBIDDEN", "you do not have permission for this action")
	})
}

func (h *policySettingsHandler) get() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		policy, err := h.repo.GetPasswordPolicy()
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not read the password policy")
			return
		}
		writeJSON(w, http.StatusOK, map[string]int{
			"minLength":     policy.MinLength,
			"minCategories": policy.MinCategories,
			"historyDepth":  policy.HistoryDepth,
		})
	})
}

func (h *policySettingsHandler) patch() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			MinLength     *int `json:"minLength"`
			MinCategories *int `json:"minCategories"`
			HistoryDepth  *int `json:"historyDepth"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxResourceBodyBytes)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_BODY", "body must be JSON")
			return
		}
		current, err := h.repo.GetPasswordPolicy()
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not read the password policy")
			return
		}
		apply := func(v *int, fallback int) int {
			if v != nil {
				return *v
			}
			return fallback
		}
		next := authsession.PasswordPolicy{
			MinLength:     apply(body.MinLength, current.MinLength),
			MinCategories: apply(body.MinCategories, current.MinCategories),
			HistoryDepth:  apply(body.HistoryDepth, current.HistoryDepth),
		}
		if next.MinLength < 8 || next.MinLength > 72 ||
			next.MinCategories < 0 || next.MinCategories > 4 ||
			next.HistoryDepth < 0 || next.HistoryDepth > 10 {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_BODY",
				"minLength must be 8-72, minCategories 0-4, historyDepth 0-10")
			return
		}
		if uerr := h.repo.UpdatePasswordPolicy(next); uerr != nil {
			// A-001 F-001 sentinel split: a missing singleton row (legacy
			// pre-0057 store) is a 404 on the frozen SETTINGS_NOT_FOUND code,
			// not a blanket 500; anything else stays a storage failure.
			if errors.Is(uerr, authsession.ErrPasswordPolicyNotSeeded) {
				writeLocalizedError(w, r, http.StatusNotFound, "SETTINGS_NOT_FOUND", "password policy row is not seeded")
				return
			}
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not save the password policy")
			return
		}
		writeJSON(w, http.StatusOK, map[string]int{
			"minLength":     next.MinLength,
			"minCategories": next.MinCategories,
			"historyDepth":  next.HistoryDepth,
		})
	})
}
