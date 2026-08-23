// MFA surface (S-10 · GOAL-017 D-002 §3/§4): the public second-factor
// verification endpoint (/api/auth/mfa/verify) plus the identity-scoped
// self-service endpoints (status/enroll/confirm/disable/recovery rotate) and
// the admin reset. The login gate itself is the MFAVerifier interface consumed
// by the auth handler; this file owns the module's HTTP contributions.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

// MFA domain error sentinels — returned by the admin.mfa module Service
// (module → handler direction, captcha precedent) and mapped here to the
// frozen wire codes.
var (
	ErrMFAProofExpired   = errors.New("mfa: proof expired")
	ErrMFAProofExhausted = errors.New("mfa: proof exhausted")
	ErrMFAInvalid        = errors.New("mfa: code invalid")
	ErrMFANotEnrolled    = errors.New("mfa: not enrolled")
	ErrMFAPendingOnly    = errors.New("mfa: pending activation")
	ErrMFAActive         = errors.New("mfa: already active")
)

// mfaVerifyRateLimiterWindow and mfaVerifyRateLimiterMax define the independent
// HTTP rate-limit budget for POST /api/auth/mfa/verify (A6 · security audit):
// 10 attempts per 15-minute window per client IP, matching the window of the
// login limiter (loginRateLimiter in auth.go: 20 attempts / 15 min). A tighter
// cap is appropriate because each proof is one-shot and 5-failure capped in the
// service; the HTTP limiter provides an outer bound independent of the login
// bucket so a replay-spray cannot exhaust the server's probe budget without
// being gated at the transport level.
const (
	mfaVerifyRateLimiterWindow   = 15 * time.Minute
	mfaVerifyRateLimiterMax      = 10
	mfaVerifyRateLimiterCapacity = 1 << 16
)

// MFASelfService is the identity-scoped self-service surface consumed by the
// /api/mfa/* routes (satisfied structurally by the admin.mfa module Service).
type MFASelfService interface {
	MFAVerifier
	Status(userID string) (enabled bool, enrolledAt time.Time, err error)
	Enroll(userID, name string, now time.Time) (secretBase32, otpauth string, recoveryCodes []string, err error)
	Confirm(userID, code string, now time.Time) error
	Disable(userID, code, recoveryCode string, now time.Time) error
	RotateRecovery(userID, code, recoveryCode string, now time.Time) ([]string, error)
	AdminReset(userID string) (removedActive bool, err error)
}

// SessionRevoker is the session-invalidation surface MFA disable / admin
// reset require (A-004 F-002: same strength as the self-service disable —
// token_version bump + full refresh revocation). Satisfied structurally by
// the auth-session repository.
type SessionRevoker interface {
	BumpTokenVersionAndRevokeAll(userID string, now time.Time) error
}

// MFARoutes returns the admin.mfa HTTP surface.
func MFARoutes(a *auth.Authenticator, service MFASelfService, operations operationlog.Recorder, revoker SessionRevoker, moduleID string) []kernel.RouteContribution {
	var routes []kernel.RouteContribution
	add := func(method, pattern string, h http.Handler) {
		routes = append(routes, kernel.RouteContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey(method, pattern)},
			Method:               method,
			Pattern:              pattern,
			Handler:              h,
		})
	}

	// A6 · independent MFA verify rate limiter (security audit finding):
	// POST /api/auth/mfa/verify gets its own counting bucket, never shared with
	// the login limiter. Key = client IP (loginClientIP, same proxy trust rules
	// as the login limiter). Window/threshold mirror the login limiter (15 min /
	// 20 max) with a tighter cap (10) because each proof is one-shot, so 10
	// HTTP-level attempts already cover all legitimate retry scenarios.
	mfaVerifyLimiter := newLoginRateLimiter(
		mfaVerifyRateLimiterWindow,
		mfaVerifyRateLimiterMax,
		mfaVerifyRateLimiterCapacity,
	)

	// Second-factor completion: public (holds the one-time proof). Issues the
	// real token pair on success (D-002 §3) and records mfa.login.
	add("POST", "/api/auth/mfa/verify", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// A6: independent per-IP rate limit on the verify surface. The key is
		// the client IP only (no username, since the proof already encodes the
		// user; an IP key still stops a replay-spray from a single host).
		limiterKey := loginClientIP(r)
		now := time.Now().UTC()
		if !mfaVerifyLimiter.allow(limiterKey, now) {
			if sec := mfaVerifyLimiter.retryAfterSeconds(limiterKey, now); sec > 0 {
				w.Header().Set("Retry-After", strconv.Itoa(sec))
			}
			writeLocalizedError(w, r, http.StatusTooManyRequests, "RATE_LIMITED", "too many MFA verify attempts; try again later")
			return
		}
		var body struct {
			Proof        string `json:"proof"`
			Code         string `json:"code"`
			RecoveryCode string `json:"recoveryCode,omitempty"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.Proof) == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_MFA_BODY", "body must be JSON with proof and code")
			return
		}
		userID, err := service.Verify(body.Proof, body.Code, body.RecoveryCode, time.Now().UTC())
		if err != nil {
			// Failed verifications count against the limiter bucket so that a
			// spray of invalid codes from one IP is bounded at the HTTP layer
			// independently of the service-level proof exhaustion counter.
			mfaVerifyLimiter.record(limiterKey, now)
			writeMFAError(w, r, err)
			return
		}
		access, refresh, user, err := a.IssueTokensFor(userID, time.Now().UTC())
		if err != nil {
			writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "authentication unavailable")
			return
		}
		if operations != nil {
			recordAudit(operations, user, operationlog.EventMFALogin, userID, auditDetail("login", map[string]any{"userId": userID}), time.Now().UTC(), r.Context())
		}
		writeJSON(w, http.StatusOK, tokenResponse{AccessToken: access, RefreshToken: refresh, User: user})
	}))

	// Status: self-service enrollment state.
	add("GET", "/api/mfa/status", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requireIdentity(w, r)
		if !ok {
			return
		}
		enabled, enrolledAt, err := service.Status(user.ID)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not read MFA status")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"enabled": enabled, "enrolledAt": enrolledAt})
	})))

	// Enroll: one-time secret + recovery codes (pending until confirm).
	add("POST", "/api/mfa/enroll", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requireIdentity(w, r)
		if !ok {
			return
		}
		var body struct {
			CurrentPassword string `json:"currentPassword"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.CurrentPassword) == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "MFA_CURRENT_PASSWORD_REQUIRED", "currentPassword is required to start MFA enrollment")
			return
		}
		// W7 F-007 step-up: a Bearer session alone must not be enough to bind a
		// new TOTP secret to the account. Confirm the current password first; an
		// account that already has active MFA is rejected by Enroll with
		// ErrActive (and must use disable/rotate instead).
		current, err := a.UserByID(user.ID)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not load account")
			return
		}
		if !auth.VerifyPassword(current.PasswordHash, body.CurrentPassword) {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PASSWORD", "current password is incorrect")
			return
		}
		secret, otpauth, codes, err := service.Enroll(user.ID, user.Name, time.Now().UTC())
		if err != nil {
			// A-008 recommended: an active enrollment maps to 400
			// MFA_ALREADY_ACTIVE (disable first), not a generic 500.
			writeMFAError(w, r, err)
			return
		}
		recordAudit(operations, user, operationlog.EventMFAEnroll, user.ID, auditDetail("enroll", map[string]any{"userId": user.ID}), time.Now().UTC(), r.Context())
		writeJSON(w, http.StatusOK, map[string]any{"secretBase32": secret, "otpauthURL": otpauth, "recoveryCodes": codes})
	})))

	// Confirm: activate the pending enrollment with a correct TOTP code.
	add("POST", "/api/mfa/confirm", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requireIdentity(w, r)
		if !ok {
			return
		}
		var body struct {
			Code string `json:"code"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.Code) == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_MFA_BODY", "body must be JSON with code")
			return
		}
		if err := service.Confirm(user.ID, body.Code, time.Now().UTC()); err != nil {
			writeSelfServiceMFAError(w, r, err)
			return
		}
		recordAudit(operations, user, operationlog.EventMFAConfirm, user.ID, auditDetail("confirm", map[string]any{"userId": user.ID}), time.Now().UTC(), r.Context())
		w.WriteHeader(http.StatusNoContent)
	})))

	// Disable: remove the enrollment after a valid code/recovery, then
	// invalidate all sessions (A-004 F-002 parity with self-service disable).
	add("POST", "/api/mfa/disable", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requireIdentity(w, r)
		if !ok {
			return
		}
		var body struct {
			Code         string `json:"code"`
			RecoveryCode string `json:"recoveryCode,omitempty"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_MFA_BODY", "body must be JSON")
			return
		}
		now := time.Now().UTC()
		if err := service.Disable(user.ID, body.Code, body.RecoveryCode, now); err != nil {
			writeSelfServiceMFAError(w, r, err)
			return
		}
		if err := revoker.BumpTokenVersionAndRevokeAll(user.ID, now); err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not invalidate sessions")
			return
		}
		recordAudit(operations, user, operationlog.EventMFADisable, user.ID, auditDetail("disable", map[string]any{"userId": user.ID}), time.Now().UTC(), r.Context())
		w.WriteHeader(http.StatusNoContent)
	})))

	// Recovery rotate: replace the one-time recovery set after validation.
	add("POST", "/api/mfa/recovery/rotate", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requireIdentity(w, r)
		if !ok {
			return
		}
		var body struct {
			Code         string `json:"code"`
			RecoveryCode string `json:"recoveryCode,omitempty"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_MFA_BODY", "body must be JSON")
			return
		}
		codes, err := service.RotateRecovery(user.ID, body.Code, body.RecoveryCode, time.Now().UTC())
		if err != nil {
			writeSelfServiceMFAError(w, r, err)
			return
		}
		recordAudit(operations, user, operationlog.EventMFARecoveryRotate, user.ID, auditDetail("recovery-rotate", map[string]any{"userId": user.ID}), time.Now().UTC(), r.Context())
		writeJSON(w, http.StatusOK, map[string]any{"recoveryCodes": codes})
	})))

	// Admin reset: remove another user's enrollment + invalidate sessions.
	add("POST", "/api/users/{id}/mfa/reset", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, "users.mfa-reset")
		if !ok {
			return
		}
		targetID := r.PathValue("id")
		now := time.Now().UTC()
		// W7 F-002: mirror the users-resource admin target boundary (users.go
		// authorizeAdminTargetBoundary) — a delegated users.mfa-reset actor
		// must not strip 2FA from an admin account. The store's SELF_OPERATION
		// guard does not apply here; this endpoint is management of another
		// user, so the actor is allowed to reset themselves only when they hold
		// the management permission (same posture as users.write).
		target, err := a.UserByID(targetID)
		if err != nil {
			writeLocalizedError(w, r, http.StatusNotFound, "USER_NOT_FOUND", "no user with that id")
			return
		}
		if !slices.Contains(user.Roles, "admin") && slices.Contains(target.Roles, "admin") {
			writeLocalizedError(w, r, http.StatusForbidden, "ADMIN_ACCOUNT_FORBIDDEN", "only an admin may reset an admin's MFA")
			return
		}
		removedActive, err := service.AdminReset(targetID)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not reset MFA")
			return
		}
		// W7 F-002: only an ACTIVE enrollment removal should invalidate the
		// target's sessions. Resetting a user who has no active MFA must not be
		// treated as a generic force-logout primitive.
		if removedActive {
			if err := revoker.BumpTokenVersionAndRevokeAll(targetID, now); err != nil {
				writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not invalidate sessions")
				return
			}
		}
		recordAudit(operations, user, operationlog.EventMFAAdminReset, targetID, auditDetail("admin-reset", map[string]any{"userId": targetID, "revokedSessions": removedActive}), time.Now().UTC(), r.Context())
		w.WriteHeader(http.StatusNoContent)
	})))

	return routes
}

// requireIdentity resolves the request identity for self-service endpoints
// (no permission key — every authenticated user can manage their own MFA).
func requireIdentity(w http.ResponseWriter, r *http.Request) (account.User, bool) {
	user, ok := auth.UserIdentityFrom(r.Context())
	if !ok {
		writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "no active session")
		return account.User{}, false
	}
	return user, true
}

// writeMFAError maps MFA domain errors to the frozen wire codes. Used by the
// login second-factor endpoint (/api/auth/mfa/verify): an invalid code there
// is an authentication failure (401) handled directly by the login surface.
func writeMFAError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ErrMFAProofExpired):
		writeLocalizedError(w, r, http.StatusUnauthorized, "MFA_PROOF_EXPIRED", "second-factor proof expired; sign in again")
	case errors.Is(err, ErrMFAProofExhausted):
		writeLocalizedError(w, r, http.StatusUnauthorized, "MFA_PROOF_EXHAUSTED", "too many failed attempts; sign in again")
	case errors.Is(err, ErrMFAInvalid):
		writeLocalizedError(w, r, http.StatusUnauthorized, "MFA_INVALID", "invalid second-factor code")
	case errors.Is(err, ErrMFANotEnrolled):
		writeLocalizedError(w, r, http.StatusBadRequest, "MFA_NOT_ENROLLED", "no MFA enrollment for this account")
	case errors.Is(err, ErrMFAPendingOnly):
		writeLocalizedError(w, r, http.StatusBadRequest, "MFA_PENDING_ONLY", "MFA is not activated yet")
	case errors.Is(err, ErrMFAActive):
		writeLocalizedError(w, r, http.StatusBadRequest, "MFA_ALREADY_ACTIVE", "MFA is already active")
	default:
		writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not complete MFA operation")
	}
}

// writeSelfServiceMFAError maps MFA domain errors for the identity-scoped
// self-service endpoints (confirm/disable/recovery rotate). An invalid code
// here is a CLIENT VALIDATION failure (400), not a lost session (401): the
// web auth wrapper (authFetch) must not treat it as session expiry and force
// a sign-out (W11 · M-02/M-03). Proof errors never occur on this surface.
func writeSelfServiceMFAError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ErrMFAInvalid):
		writeLocalizedError(w, r, http.StatusBadRequest, "MFA_INVALID", "invalid second-factor code")
	case errors.Is(err, ErrMFANotEnrolled):
		writeLocalizedError(w, r, http.StatusBadRequest, "MFA_NOT_ENROLLED", "no MFA enrollment for this account")
	case errors.Is(err, ErrMFAPendingOnly):
		writeLocalizedError(w, r, http.StatusBadRequest, "MFA_PENDING_ONLY", "MFA is not activated yet")
	case errors.Is(err, ErrMFAActive):
		writeLocalizedError(w, r, http.StatusBadRequest, "MFA_ALREADY_ACTIVE", "MFA is already active")
	default:
		writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not complete MFA operation")
	}
}


