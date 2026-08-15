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
	ErrMFAProofExpired  = errors.New("mfa: proof expired")
	ErrMFAProofExhausted = errors.New("mfa: proof exhausted")
	ErrMFAInvalid       = errors.New("mfa: code invalid")
	ErrMFANotEnrolled   = errors.New("mfa: not enrolled")
	ErrMFAPendingOnly   = errors.New("mfa: pending activation")
	ErrMFAActive        = errors.New("mfa: already active")
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
	AdminReset(userID string) error
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

	// Second-factor completion: public (holds the one-time proof). Issues the
	// real token pair on success (D-002 §3) and records mfa.login.
	add("POST", "/api/auth/mfa/verify", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			writeMFAError(w, r, err)
			return
		}
		access, refresh, user, err := a.IssueTokensFor(userID, time.Now().UTC())
		if err != nil {
			writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "authentication unavailable")
			return
		}
		if operations != nil {
			now := time.Now().UTC()
			detail := `{"userId":` + jsonQuote(userID) + `}`
			_ = operations.RecordOperation(operationlog.Operation{
				ID: newOperationID(), Event: operationlog.EventMFALogin,
				ActorID: user.ID, ActorName: user.Name, Detail: &detail, CreatedAt: now,
			})
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
		secret, otpauth, codes, err := service.Enroll(user.ID, user.Name, time.Now().UTC())
		if err != nil {
			// A-008 recommended: an active enrollment maps to 400
			// MFA_ALREADY_ACTIVE (disable first), not a generic 500.
			writeMFAError(w, r, err)
			return
		}
		recordMFAEvent(operations, user, operationlog.EventMFAEnroll, `{"userId":`+jsonQuote(user.ID)+`}`)
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
		recordMFAEvent(operations, user, operationlog.EventMFAConfirm, `{"userId":`+jsonQuote(user.ID)+`}`)
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
		recordMFAEvent(operations, user, operationlog.EventMFADisable, `{"userId":`+jsonQuote(user.ID)+`}`)
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
		recordMFAEvent(operations, user, operationlog.EventMFARecoveryRotate, `{"userId":`+jsonQuote(user.ID)+`}`)
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
		if err := service.AdminReset(targetID); err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not reset MFA")
			return
		}
		if err := revoker.BumpTokenVersionAndRevokeAll(targetID, now); err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not invalidate sessions")
			return
		}
		recordMFAEvent(operations, user, operationlog.EventMFAAdminReset, `{"userId":`+jsonQuote(targetID)+`}`)
		w.WriteHeader(http.StatusNoContent)
	})))

	return routes
}

// requireIdentity resolves the request identity for self-service endpoints
// (no permission key — every authenticated user can manage their own MFA).
func requireIdentity(w http.ResponseWriter, r *http.Request) (account.User, bool) {
	user, ok := auth.IdentityFrom(r.Context())
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

// recordMFAEvent writes an MFA audit row.
func recordMFAEvent(operations operationlog.Recorder, user account.User, event, detail string) {
	if operations == nil {
		return
	}
	now := time.Now().UTC()
	_ = operations.RecordOperation(operationlog.Operation{
		ID: newOperationID(), Event: event,
		ActorID: user.ID, ActorName: user.Name, Detail: &detail, CreatedAt: now,
	})
}
