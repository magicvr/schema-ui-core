package handler

import (
	"net/http"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
)

// accountsHandler registers GET /api/accounts/me behind the request-identity
// middleware: the identity is resolved from the Bearer access token (or the
// explicit dev-session fallback), never from a process-injected static session.
func accountsHandler(mux *http.ServeMux, a *auth.Authenticator) {
	mux.Handle("GET /api/accounts/me", a.Middleware(http.HandlerFunc(me)))
}

// me returns the current account session and $context snapshot for the
// authenticated request identity.
func me(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.IdentityFrom(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHENTICATED", "no active session")
		return
	}
	writeJSON(w, http.StatusOK, account.Session{
		User: user,
		// R2 returns no server-side feature flags; the renderer treats absence
		// fail-closed.
		Features: map[string]bool{},
	})
}
