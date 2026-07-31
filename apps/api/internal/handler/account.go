package handler

import (
	"net/http"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
)

// accountHandler serves the R4 minimal session contract.
type accountHandler struct {
	// sessionProvider is injectable for tests; nil means the static dev
	// session. Fail-closed: an empty provider never produces a session.
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
		session, ok := h.sessionProvider()
		if !ok {
			writeError(w, http.StatusUnauthorized, "UNAUTHENTICATED", "no active session")
			return
		}
		writeJSON(w, http.StatusOK, session)
	})
}
