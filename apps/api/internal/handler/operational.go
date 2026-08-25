package handler

import (
	"net/http"
	"strings"

	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
)

// WithOperationalGate applies the R5 write boundary after request-id/CORS and
// before the route handler. Only a currently registered mutation is gated;
// unknown paths and method mismatches remain the responsibility of the JSON
// 404/405 envelope.
func WithOperationalGate(cfg *config.Config, mux *http.ServeMux, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mode := cfg.RuntimeMode
		if mode == "" || mode == config.RuntimeModeNormal || !isRegisteredCurrentMethod(mux, r) || !isBusinessMutation(r) || operationalAllowlisted(r) {
			next.ServeHTTP(w, r)
			return
		}

		code := "SERVICE_DEGRADED"
		switch mode {
		case config.RuntimeModeMaintenance:
			code = "SERVICE_MAINTENANCE"
		case config.RuntimeModeReadOnly:
			code = "SERVICE_READ_ONLY"
		case config.RuntimeModeDegraded:
			// Keep the degraded code selected above.
		default:
			// Invalid values should have failed configuration loading. If a
			// zero/invalid Config is assembled directly, fail closed as degraded.
			code = "SERVICE_DEGRADED"
		}
		writeLocalizedError(w, r, http.StatusServiceUnavailable, code, operationalMessage(code))
	})
}

func isRegisteredCurrentMethod(mux *http.ServeMux, r *http.Request) bool {
	_, pattern := mux.Handler(r)
	if pattern == "" {
		return false
	}
	if method, _, ok := strings.Cut(pattern, " "); ok && method != r.Method {
		return false
	}
	return true
}

func isBusinessMutation(r *http.Request) bool {
	switch r.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

// These paths preserve sign-in, session recovery and self-service MFA/password
// recovery while business writes are paused. Keep this explicit registry in
// sync with the registered recovery routes; ordinary auth-adjacent mutations
// must remain gated.
var operationalRecoveryPaths = map[string]struct{}{
	"/api/auth/login":                {},
	"/api/auth/refresh":              {},
	"/api/auth/logout":               {},
	"/api/auth/mfa/verify":           {},
	"/api/account/password":          {},
	"/api/mfa/enroll":                {},
	"/api/mfa/confirm":               {},
	"/api/mfa/disable":               {},
	"/api/mfa/recovery/rotate":       {},
	// workspace-019 R2: self-recovery IS the session-recovery path — it must
	// stay reachable in maintenance/read-only modes (GOAL-003 D-001 §2).
	"/api/auth/recovery/start":     {},
	"/api/auth/recovery/complete":  {},
}

func operationalAllowlisted(r *http.Request) bool {
	if r.Method != http.MethodPost {
		return false
	}
	_, ok := operationalRecoveryPaths[r.URL.Path]
	return ok
}

func operationalMessage(code string) string {
	switch code {
	case "SERVICE_MAINTENANCE":
		return "service is under maintenance"
	case "SERVICE_READ_ONLY":
		return "service is read-only"
	default:
		return "service is operating in degraded mode"
	}
}
