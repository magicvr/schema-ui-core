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

// These paths preserve sign-in, session recovery and forced password change
// while business writes are paused. Matching is intentionally exact.
func operationalAllowlisted(r *http.Request) bool {
	if r.Method != http.MethodPost {
		return false
	}
	switch r.URL.Path {
	case "/api/auth/login", "/api/auth/refresh", "/api/auth/logout",
		"/api/auth/mfa/verify", "/api/account/password":
		return true
	default:
		return false
	}
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
