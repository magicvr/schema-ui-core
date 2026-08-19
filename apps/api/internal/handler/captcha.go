// Login captcha surface (S-11 · GOAL-011 D-002 §2/§3): the public preflight
// (GET /api/auth/captcha) and the admin settings surface
// (GET/PATCH /api/captcha/settings). The login gate itself is the
// CaptchaVerifier interface consumed by the auth handler; this file only
// owns the module's HTTP contributions.
package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

// CaptchaService is the challenge surface the captcha routes consume. It is
// satisfied structurally by the admin.login-captcha module Service (no handler
// import of the module package — the direction is module → handler).
type CaptchaService interface {
	// Generate issues a new one-time challenge.
	Generate() (id, question string, expiresInSeconds int64, err error)
	// Required reports whether login must present a captcha.
	Required() bool
	// SetEnabled flips the login captcha gate (admin settings).
	SetEnabled(enabled bool, now time.Time) error
}

// captchaGenerateLimiter bounds anonymous challenge generation per client IP
// (W7 F-006): even when the gate is enabled, a machine cannot flood
// /api/auth/captcha to obtain unlimited solvable questions or fill the
// captcha_challenges table. The same trusted-client-IP logic as login
// rate limiting is reused.
var captchaGenerateLimiter = newLoginRateLimiter(time.Minute, 10, 1<<16)

// CaptchaRoutes returns the admin.login-captcha HTTP surface.
func CaptchaRoutes(a *auth.Authenticator, service CaptchaService, operations operationlog.Recorder, moduleID string) []kernel.RouteContribution {
	var routes []kernel.RouteContribution

	// Public preflight: the login page asks BEFORE authentication, so this
	// route must not require a session (D-002 §2). When the gate is disabled
	// the client skips the challenge entirely; when enabled the client renders
	// the arithmetic question from the challenge payload.
	routes = append(routes, kernel.RouteContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("GET", "/api/auth/captcha")},
		Method:               "GET",
		Pattern:              "/api/auth/captcha",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			enabled := service.Required()
			body := map[string]any{"enabled": enabled}
			if enabled {
				// W7 F-006: rate-limit challenge generation (per real client IP)
				// so the public preflight cannot be used as an unlimited
				// solve-and-retry oracle or a table-filling pump.
				if !captchaGenerateLimiter.allow(loginClientIP(r), time.Now().UTC()) {
					writeLocalizedError(w, r, http.StatusTooManyRequests, "RATE_LIMITED", "too many captcha requests; try again later")
					return
				}
				// W7 F-006: record this generation attempt so the sliding window
				// actually counts requests (allow() only checks; record() creates
				// the entry). This bounds an anonymous client to 10 challenges per
				// minute.
				captchaGenerateLimiter.record(loginClientIP(r), time.Now().UTC())
				id, question, expiresInSeconds, err := service.Generate()
				if err != nil {
					writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not generate captcha")
					return
				}
				body["challenge"] = map[string]any{
					"id":               id,
					"question":         question,
					"expiresInSeconds": expiresInSeconds,
				}
			}
			writeJSON(w, http.StatusOK, body)
		}),
	})

	// Settings: read current gate state.
	routes = append(routes, kernel.RouteContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("GET", "/api/captcha/settings")},
		Method:               "GET",
		Pattern:              "/api/captcha/settings",
		Handler: a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := requirePermission(w, r, "captcha.read"); !ok {
				return
			}
			// Form-facing string value ("true"/"false") so the schema select
			// control can prefill (F-003, notifications F-001 pattern).
			writeJSON(w, http.StatusOK, map[string]any{"enabled": boolStringValue(service.Required())})
		})),
	})

	// Settings: flip the gate (audited).
	routes = append(routes, kernel.RouteContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("PATCH", "/api/captcha/settings")},
		Method:               "PATCH",
		Pattern:              "/api/captcha/settings",
		Handler: a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := requirePermission(w, r, "captcha.write")
			if !ok {
				return
			}
			var body struct {
				Enabled json.RawMessage `json:"enabled"`
			}
			r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.Enabled) == 0 {
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SETTINGS_BODY", "body must be JSON with enabled")
				return
			}
			// Accept a JSON bool or the "true"/"false" strings the schema select
			// control submits (F-003, notifications F-001 pattern).
			enabled, err := parseBoolValue(body.Enabled)
			if err != nil {
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SETTINGS_BODY", "enabled must be a boolean or "+`"true"/"false"`)
				return
			}
			now := time.Now().UTC()
			if err := service.SetEnabled(enabled, now); err != nil {
				writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not update captcha settings")
				return
			}
			recordCaptchaSettingsEvent(operations, user, enabled, now)
			writeJSON(w, http.StatusOK, map[string]any{"enabled": boolStringValue(enabled)})
		})),
	})
	return routes
}

// recordCaptchaSettingsEvent writes the captcha.settings-update audit row.
func recordCaptchaSettingsEvent(operations operationlog.Recorder, user account.User, enabled bool, now time.Time) {
	if operations == nil {
		return
	}
	recordAudit(operations, user, operationlog.EventCaptchaSettingsUpdate, "", auditDetail("settings-update", map[string]any{"enabled": enabled}), now, nil)
}
