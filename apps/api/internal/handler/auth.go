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
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

// authHandler serves the R2 auth endpoints (GOAL-005): login, refresh and
// logout. Login/refresh are public; the access/refresh pair returned is consumed
// by the Web client (access in memory, refresh in localStorage, D-002). The
// store is used only for the R5 S6 operation log (auth events); identity data
// is resolved through the module-owned auth-session repository.
type authHandler struct {
	a           *auth.Authenticator
	operations  operationlog.Recorder
	now         func() time.Time
	rateLimiter *loginRateLimiter
}

func authsHandler(mux *http.ServeMux, a *auth.Authenticator, operations operationlog.Recorder) {
	h := &authHandler{
		a:           a,
		operations:  operations,
		now:         time.Now,
		rateLimiter: newLoginRateLimiter(15*time.Minute, 20, 1<<16),
	}
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
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_LOGIN_BODY", "body must be JSON with username and password")
			return
		}
		if strings.TrimSpace(creds.Username) == "" || creds.Password == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_LOGIN_BODY", "username and password are required")
			return
		}
		// D-001 P1: the bucket key is the real client IP (trusted reverse proxy
		// X-Real-IP) plus the attempted username, so one attacker spraying many
		// usernames cannot lock out unrelated clients behind the same proxy.
		limiterKey := loginClientIP(r) + "|" + strings.ToLower(strings.TrimSpace(creds.Username))
		if h.rateLimiter != nil && !h.rateLimiter.allow(limiterKey, h.now().UTC()) {
			writeLocalizedError(w, r, http.StatusTooManyRequests, "RATE_LIMITED", "too many failed login attempts; try again later")
			return
		}
		access, refresh, user, err := h.a.Login(creds.Username, creds.Password, h.now().UTC())
		if errors.Is(err, auth.ErrInvalidCredentials) {
			if h.rateLimiter != nil {
				h.rateLimiter.record(limiterKey, h.now().UTC())
			}
			writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "invalid username or password")
			return
		}
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "LOGIN_FAILED", "authentication unavailable")
			return
		}
		// A successful login clears the client's failure bucket (D-001 P1): a
		// legitimate user who mis-typed the password a few times must not be
		// locked out by their own earlier failures.
		if h.rateLimiter != nil {
			h.rateLimiter.clear(limiterKey)
		}
		h.logOperation(operationlog.EventAuthLogin, user.ID, user.Name, `{"username":`+jsonQuote(creds.Username)+`}`)
		writeJSON(w, http.StatusOK, tokenResponse{AccessToken: access, RefreshToken: refresh, User: user})
	}
}

// refresh rotates a valid refresh token into a new access/refresh pair.
func (h *authHandler) refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body tokenRequest
		r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_REFRESH_BODY", "body must be JSON with refreshToken")
			return
		}
		access, refresh, user, err := h.a.Refresh(body.RefreshToken, h.now().UTC())
		if errors.Is(err, auth.ErrInvalidToken) || errors.Is(err, auth.ErrExpiredToken) || errors.Is(err, auth.ErrTokenRevoked) {
			writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "invalid, expired or revoked refresh token")
			return
		}
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "REFRESH_FAILED", "refresh unavailable")
			return
		}
		h.authEvent(operationlog.EventAuthRefresh, user.ID)
		writeJSON(w, http.StatusOK, tokenResponse{AccessToken: access, RefreshToken: refresh, User: user})
	}
}

// logout revokes the presented refresh token (idempotent) and records the
// auth.logout operation for the token's owner (I-008-003 §2/§5).
func (h *authHandler) logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body tokenRequest
		r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_LOGOUT_BODY", "body must be JSON with refreshToken")
			return
		}
		userID, err := h.a.Logout(body.RefreshToken, h.now().UTC())
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "LOGOUT_FAILED", "logout unavailable")
			return
		}
		if userID != "" {
			h.authEvent(operationlog.EventAuthLogout, userID)
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// authEvent records an auth operation-log event with the frozen username detail
// (I-008-003 §3: detail = {"username":"<用户名>"}). The account snapshot from
// auth does not carry the login username, so it is resolved from the store;
// best-effort: an unresolvable actor still logs with the actor id and a
// service-log error, never blocking the business response (§5).
func (h *authHandler) authEvent(event, userID string) {
	u, err := h.a.UserByID(userID)
	if err != nil {
		slog.Error("operation log auth event: resolve user", "event", event, "user_id", userID, "err", err)
		h.logOperation(event, userID, "", "")
		return
	}
	h.logOperation(event, u.ID, u.Name, `{"username":`+jsonQuote(u.Username)+`}`)
}

// logOperation appends one operation-log row (R5 S6 optional bonus checkpoint,
// I-008-003 §5). Best-effort: a logging failure is logged to the service log
// and never changes the business response.
func (h *authHandler) logOperation(event, actorID, actorName, detail string) {
	op := operationlog.Operation{
		ID:        newOperationID(),
		Event:     event,
		ActorID:   actorID,
		ActorName: actorName,
		CreatedAt: time.Now().UTC(),
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
