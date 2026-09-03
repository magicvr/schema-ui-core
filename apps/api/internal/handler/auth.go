package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/requestid"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
)

// MFAVerifier is the optional second-factor login gate (S-10 · GOAL-017
// D-002 §3): nil disables the gate entirely — the login contract stays
// byte-identical. Required()==true means the user must complete a second
// factor before tokens are issued; BeginChallenge issues the one-time proof
// after the password factor succeeded; Verify completes the login and returns
// the verified user id.
type MFAVerifier interface {
	Required(userID string) bool
	BeginChallenge(userID string, now time.Time) (proof string, err error)
	Verify(proof, code, recoveryCode string, now time.Time) (userID string, err error)
}

// CaptchaVerifier is the optional login captcha gate (S-11 · GOAL-011 D-002
// `2): nil disables the gate entirely (module not enabled) — the login
// contract is byte-identical to the pre-captcha behavior; a verifier with
// Required()==true demands a valid one-time challenge before credentials.
type CaptchaVerifier interface {
	Required() bool
	Verify(captchaID, answer string, now time.Time) error
}

// authHandler serves the R2 auth endpoints (GOAL-005): login, refresh and
// logout. Login/refresh are public; the access/refresh pair returned is consumed
// by the Web client (access in memory, refresh in localStorage, D-002). The
// store is used only for the R5 S6 operation log (auth events); identity data
// is resolved through the module-owned auth-session repository.
type authHandler struct {
	a           *auth.Authenticator
	operations  operationlog.Recorder
	now         func() time.Time
	rateLimiter kernel.RateLimiter
	captcha     CaptchaVerifier
	// mfa is the optional second-factor gate (S-10 · GOAL-017 D-002 §3);
	// nil keeps the login contract byte-identical.
	mfa MFAVerifier
}

func authsHandler(mux routeRegistrar, a *auth.Authenticator, operations operationlog.Recorder, limiters kernel.RateLimiterProvider, captcha CaptchaVerifier, mfa ...MFAVerifier) {
	h := &authHandler{
		a:           a,
		operations:  operations,
		now:         time.Now,
		rateLimiter: limiters.NewRateLimiter(15*time.Minute, 20, 1<<16),
		captcha:     captcha,
	}
	if len(mfa) > 0 && mfa[0] != nil {
		h.mfa = mfa[0]
		a.SetMFAEnforcer(mfa[0])
	}
	mux.HandleFunc("POST /api/auth/login", h.login())
	mux.HandleFunc("POST /api/auth/refresh", h.refresh())
	mux.HandleFunc("POST /api/auth/logout", h.logout())
}

type credentials struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	CaptchaID     string `json:"captchaId,omitempty"`
	CaptchaAnswer string `json:"captchaAnswer,omitempty"`
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
		// VP-032: atomic AllowRecord (optimistic slot reservation) eliminates
		// TOCTOU between check and record. Success clears the bucket.
		limiterKey := loginClientIP(r) + "|" + strings.ToLower(strings.TrimSpace(creds.Username))
		if h.rateLimiter != nil && !h.rateLimiter.AllowRecord(limiterKey, h.now().UTC()) {
			if sec := h.rateLimiter.RetryAfterSeconds(limiterKey, h.now().UTC()); sec > 0 {
				w.Header().Set("Retry-After", strconv.Itoa(sec))
			}
			writeLocalizedError(w, r, http.StatusTooManyRequests, "RATE_LIMITED", "too many failed login attempts; try again later")
			return
		}
		// S-11 (GOAL-011 D-002): when the captcha gate is enabled it runs after
		// the rate limiter (which bounds challenge exhaustion) and before
		// credential validation. Any failure maps to the single INVALID_CAPTCHA
		// code; failures do not count against the lockout budget.
		if h.captcha != nil && h.captcha.Required() {
			if err := h.captcha.Verify(creds.CaptchaID, creds.CaptchaAnswer, h.now().UTC()); err != nil {
				if h.rateLimiter != nil {
					h.rateLimiter.Clear(limiterKey)
				}
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_CAPTCHA", "captcha verification failed")
				return
			}
		}
		// GOAL-014 D-002: the real client identity feeds the per-(account|
		// source) lockout bucket — third-party failures against a known
		// username no longer lock the legitimate user's own logins.
		access, refresh, user, err := h.a.Login(creds.Username, creds.Password, h.now().UTC(), loginClientIP(r))
		// W7 F-009: locked / disabled accounts surface the SAME 401
		// UNAUTHORIZED envelope as an unknown user / wrong password. The
		// distinct 423/403 previously let an attacker distinguish "exists and
		// locked/disabled" from "does not exist" (account-enumeration oracle).
		// Fail-closed: no lock/disable status leak on the login surface.
		// W11 F-007 (D2 residual): these terminal states now also count
		// against the client's failure bucket. With entrance AllowRecord,
		// the slot is already recorded.
		if errors.Is(err, auth.ErrAccountLocked) || errors.Is(err, auth.ErrAccountDisabled) {
			writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "invalid username or password")
			return
		}
		if errors.Is(err, auth.ErrInvalidCredentials) {
			writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "invalid username or password")
			return
		}
		// S-10 (GOAL-017 D-002 §3): the password factor succeeded but the
		// account requires a second factor — no token is issued yet; the
		// client completes /api/auth/mfa/verify with the one-time proof.
		var mfaReq *auth.MFARequiredError
		if errors.As(err, &mfaReq) {
			if h.mfa == nil {
				// The gate vanished between Login and here — fail closed.
				writeLocalizedError(w, r, http.StatusInternalServerError, "LOGIN_FAILED", "authentication unavailable")
				return
			}
			// W11 F-003: the "password passed, awaiting second factor" state
			// is rate-limited like the password factor itself. With entrance
			// AllowRecord, this attempt was already counted in limiterKey.
			proof, perr := h.mfa.BeginChallenge(mfaReq.UserID, h.now().UTC())
			if perr != nil {
				writeLocalizedError(w, r, http.StatusInternalServerError, "LOGIN_FAILED", "authentication unavailable")
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{"mfaRequired": true, "mfaProof": proof})
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
			h.rateLimiter.Clear(limiterKey)
		}
		h.logOperation(operationlog.EventAuthLogin, user.ID, user.Name, newAuthDetail("login", creds.Username), requestid.FromContext(r.Context()), user.SessionID)
		// W17 GOAL-018 D-001 I-003: set the refresh token as an httpOnly cookie
		// (priority 1: browser SPA defense against XSS exfiltration). The JSON
		// response still contains refreshToken for non-browser client compatibility
		// (mobile SDKs, CLI tools) per I-002 three-layer fallback.
		setRefreshCookie(w, refresh, !isDevMode(r))
		writeJSON(w, http.StatusOK, tokenResponse{AccessToken: access, RefreshToken: refresh, User: user})
	}
}

// refresh rotates a valid refresh token into a new access/refresh pair.
func (h *authHandler) refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// W17 GOAL-018 D-001 I-002: three-layer fallback (Cookie → Header → Body)
		refreshToken := extractRefreshToken(r)
		if refreshToken == "" {
			// Priority 3: try JSON body if Cookie and Header both empty
			var body tokenRequest
			r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_REFRESH_BODY", "body must be JSON with refreshToken")
				return
			}
			refreshToken = body.RefreshToken
		}
		if refreshToken == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "MISSING_REFRESH_TOKEN", "refresh token required in cookie, header, or body")
			return
		}
		access, refresh, user, err := h.a.Refresh(refreshToken, h.now().UTC())
		if errors.Is(err, auth.ErrInvalidToken) || errors.Is(err, auth.ErrExpiredToken) || errors.Is(err, auth.ErrTokenRevoked) {
			writeLocalizedError(w, r, http.StatusUnauthorized, "REFRESH_TOKEN_EXPIRED", "invalid, expired or revoked refresh token")
			return
		}
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "REFRESH_FAILED", "refresh unavailable")
			return
		}
		h.authEvent(operationlog.EventAuthRefresh, user.ID, requestid.FromContext(r.Context()), user.SessionID)
		// W17 GOAL-018 D-001 I-003: update the httpOnly cookie with the new token
		// (cookie rotation on every refresh). The JSON response still contains
		// refreshToken for non-browser client compatibility.
		setRefreshCookie(w, refresh, !isDevMode(r))
		writeJSON(w, http.StatusOK, tokenResponse{AccessToken: access, RefreshToken: refresh, User: user})
	}
}

// logout revokes the presented refresh token (idempotent) and records the
// auth.logout operation for the token's owner (I-008-003 §2/§5).
func (h *authHandler) logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// W17 GOAL-018 D-001 I-002: three-layer fallback (Cookie → Header → Body)
		refreshToken := extractRefreshToken(r)
		if refreshToken == "" {
			// Priority 3: try JSON body if Cookie and Header both empty
			var body tokenRequest
			r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_LOGOUT_BODY", "body must be JSON with refreshToken")
				return
			}
			refreshToken = body.RefreshToken
		}
		if refreshToken == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "MISSING_REFRESH_TOKEN", "refresh token required in cookie, header, or body")
			return
		}
		userID, sessionID, err := h.a.Logout(refreshToken, h.now().UTC())
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "LOGOUT_FAILED", "logout unavailable")
			return
		}
		if userID != "" {
			h.authEvent(operationlog.EventAuthLogout, userID, requestid.FromContext(r.Context()), sessionID)
		}
		// W17 GOAL-018 D-001: clear the httpOnly cookie on logout
		clearRefreshCookie(w)
		w.WriteHeader(http.StatusNoContent)
	}
}

// authEvent records an auth operation-log event with the R2 versioned username
// detail. The account snapshot from auth does not carry the login username, so
// it is resolved from the store;
// best-effort: an unresolvable actor still logs with the actor id and a
// service-log error, never blocking the business response (§5).
func (h *authHandler) authEvent(event, userID, correlationID, sessionID string) {
	u, err := h.a.UserByID(userID)
	if err != nil {
		slog.Error("operation log auth event: resolve user", "event", event, "user_id", userID, "err", err)
		h.logOperation(event, userID, "", "", correlationID, sessionID)
		return
	}
	h.logOperation(event, u.ID, u.Name, newAuthDetail(strings.TrimPrefix(event, "auth."), u.Username), correlationID, sessionID)
}

func newAuthDetail(action, username string) string {
	detail, err := operationlog.NewDetail(action, nil, map[string]any{"username": username})
	if err != nil {
		slog.Error("operation log auth detail: build", "action", action, "err", err)
		return ""
	}
	return detail
}

// logOperation appends one operation-log row (R5 S6 optional bonus checkpoint,
// I-008-003 §5). Best-effort: a logging failure is logged to the service log
// and never changes the business response.
func (h *authHandler) logOperation(event, actorID, actorName, detail, correlationID, sessionID string) {
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
	op.CorrelationID = correlationID
	op.SessionID = sessionID
	if h.operations == nil {
		return
	}
	if err := h.operations.RecordOperation(op); err != nil {
		slog.Error("operation log write failed", "event", event, "err", err)
	}
}
