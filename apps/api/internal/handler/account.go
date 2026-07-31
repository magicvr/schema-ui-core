package handler

import (
	"net/http"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
)

// accountHandler serves the R4 minimal session contract.
type accountHandler struct {
	// sessionProvider is injectable for tests; nil is fail-closed and produces
	// no session (the /me endpoint responds Unauthorized).
	sessionProvider func() (account.Session, bool)
}

func sessionProvider() (account.Session, bool) {
	return account.StaticDevSession(), true
}

func accountsHandler(mux *http.ServeMux) {
	h := &accountHandler{sessionProvider: sessionProvider}
	mux.Handle("GET /api/accounts/me", h.me())
}

// me returns the current account session and $context snapshot: { user, features }.
func (h *accountHandler) me() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		provider := h.sessionProvider
		if provider == nil {
			writeError(w, http.StatusUnauthorized, "UNAUTHENTICATED", "no active session")
			return
		}
		session, ok := provider()
		if !ok {
			writeError(w, http.StatusUnauthorized, "UNAUTHENTICATED", "no active session")
			return
		}
		writeJSON(w, http.StatusOK, session)
	})
}
