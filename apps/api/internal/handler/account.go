package handler

import (
	"net/http"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
)

// accountsHandler registers GET /api/accounts/me behind the request-identity
// middleware: the identity is resolved from the Bearer access token (or the
// explicit dev-session fallback), never from a process-injected static session.
func accountsHandler(mux routeRegistrar, a *auth.Authenticator) {
	mux.Handle("GET /api/accounts/me", a.Middleware(meHandler(a)))
}

// meHandler returns the current account session ($context snapshot) for the
// authenticated request identity, including the boolean features projection from
// the persisted menu grants (GOAL-006 S5 / I-006-002).
func meHandler(a *auth.Authenticator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := auth.UserIdentityFrom(r.Context())
		if !ok {
			writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "no active session")
			return
		}
		features, err := a.Features(user)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not resolve features")
			return
		}
		writeJSON(w, http.StatusOK, account.Session{User: user, Features: features})
	}
}
