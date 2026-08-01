package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
)

// authHandler serves the R2 auth endpoints (GOAL-005): login, refresh and
// logout. Login/refresh are public; the access/refresh pair returned is consumed
// by the Web client (access in memory, refresh in localStorage, D-002).
type authHandler struct {
	a   *auth.Authenticator
	now func() time.Time
}

func authsHandler(mux *http.ServeMux, a *auth.Authenticator) {
	h := &authHandler{a: a, now: time.Now}
	mux.HandleFunc("POST /api/auth/login", h.login())
	mux.HandleFunc("POST /api/auth/refresh", h.refresh())
	mux.HandleFunc("POST /api/auth/logout", h.logout())
}

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type tokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type tokenResponse struct {
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	User         account.User `json:"user"`
}

// login authenticates a username/password and issues an access/refresh pair.
func (h *authHandler) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds credentials
		r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_LOGIN_BODY", "body must be JSON with username and password")
			return
		}
		if strings.TrimSpace(creds.Username) == "" || creds.Password == "" {
			writeError(w, http.StatusBadRequest, "INVALID_LOGIN_BODY", "username and password are required")
			return
		}
		access, refresh, user, err := h.a.Login(creds.Username, creds.Password, h.now().UTC())
		if errors.Is(err, auth.ErrInvalidCredentials) {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid username or password")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "LOGIN_FAILED", "authentication unavailable")
			return
		}
		writeJSON(w, http.StatusOK, tokenResponse{AccessToken: access, RefreshToken: refresh, User: user})
	}
}

// refresh rotates a valid refresh token into a new access/refresh pair.
func (h *authHandler) refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body tokenRequest
		r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_REFRESH_BODY", "body must be JSON with refreshToken")
			return
		}
		access, refresh, user, err := h.a.Refresh(body.RefreshToken, h.now().UTC())
		if errors.Is(err, auth.ErrInvalidToken) || errors.Is(err, auth.ErrExpiredToken) || errors.Is(err, auth.ErrTokenRevoked) {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid, expired or revoked refresh token")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "REFRESH_FAILED", "refresh unavailable")
			return
		}
		writeJSON(w, http.StatusOK, tokenResponse{AccessToken: access, RefreshToken: refresh, User: user})
	}
}

// logout revokes the presented refresh token (idempotent).
func (h *authHandler) logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body tokenRequest
		r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_LOGOUT_BODY", "body must be JSON with refreshToken")
			return
		}
		if err := h.a.Logout(body.RefreshToken, h.now().UTC()); err != nil {
			writeError(w, http.StatusInternalServerError, "LOGOUT_FAILED", "logout unavailable")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
