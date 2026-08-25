// Self-service password recovery HTTP surface (workspace-019 R2 · GOAL-003
// D-001 §2): the pre-auth start/complete pair behind the login surface.
// Delivery rides the ONE composed kernel.MailSender; MFA-enrolled accounts
// face a second-factor gate between code match and the password write; every
// failure path consumes a recovery attempt so the total guess budget per
// challenge lifecycle stays ≤5 (Root D-002 frozen numbers).
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"log/slog"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

// RecoverySecondFactor is the optional admin.mfa gate for recovery
// completion. nil interface = module disabled → no second factor is demanded
// (GOAL-002 D-001 §1: accounts without MFA are unaffected).
type RecoverySecondFactor interface {
	Required(userID string) bool
	// VerifySecondFactor validates a TOTP or consumes a one-time recovery
	// code directly (no login proof); implemented by *mfa.Service.
	VerifySecondFactor(userID, code, recoveryCode string, now time.Time) error
}

// RecoveryRepository is the persistence surface consumed by the recovery
// endpoints.
type RecoveryRepository interface {
	ResolveRecoveryTarget(identifier string) (*authsession.RecoveryTarget, error)
	StartRecovery(identifier string, sender kernel.MailSender, now time.Time) error
	EvaluateRecoveryCode(userID, code string, now time.Time) (authsession.RecoveryOutcome, error)
	ConsumeRecoveryAttempt(userID string) bool
	DropStaleRecoveryChallenge(userID string, now time.Time)
	CompleteRecovery(userID, passwordHash, actorID string, now time.Time) error
	UserByID(id string) (*authsession.User, error)
}

// RegisterRecovery mounts the two public routes on the central mux (same
// layer as login/refresh/logout — NOT a module provider: self-recovery must
// exist on every profile that has core.auth-session).
func RegisterRecovery(mux routeRegistrar, operations operationlog.Recorder, repo RecoveryRepository, notifier NotifyRepository, sender kernel.MailSender, gate RecoverySecondFactor) {
	h := &recoveryHandler{
		repo:       repo,
		notifier:   notifier,
		sender:     sender,
		gate:       gate,
		operations: operations,
		now:        time.Now,
		rateLimiter: newLoginRateLimiter(15*time.Minute, 20, 1<<16),
	}
	mux.HandleFunc("POST /api/auth/recovery/start", h.start())
	mux.HandleFunc("POST /api/auth/recovery/complete", h.complete())
}

type recoveryHandler struct {
	repo        RecoveryRepository
	notifier    NotifyRepository
	sender      kernel.MailSender
	gate        RecoverySecondFactor
	operations  operationlog.Recorder
	now         func() time.Time
	rateLimiter *loginRateLimiter
}

func (h *recoveryHandler) limiterKey(r *http.Request, account string) string {
	return loginClientIP(r) + "|" + strings.ToLower(strings.TrimSpace(account))
}

// start issues a reset code to the account's bound+verified address. Every
// "no self-recovery path" outcome answers the SAME 202 dispatched shape with
// no mail sent (anti-enumeration fail-closed, W7 F-009口径): an attacker cannot
// distinguish unknown usernames from accounts without a verified email.
func (h *recoveryHandler) start() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Account string `json:"account"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.Account) == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_RECOVERY_BODY", "body must be JSON with an account field")
			return
		}
		key := h.limiterKey(r, body.Account)
		if !h.rateLimiter.allow(key, h.now().UTC()) {
			if sec := h.rateLimiter.retryAfterSeconds(key, h.now().UTC()); sec > 0 {
				w.Header().Set("Retry-After", strconv.Itoa(sec))
			}
			writeLocalizedError(w, r, http.StatusTooManyRequests, "RATE_LIMITED", "too many recovery requests; try again later")
			return
		}
		if err := h.repo.StartRecovery(body.Account, h.sender, h.now().UTC()); err != nil {
			switch {
			case errors.Is(err, authsession.ErrRecoveryNotAvailable):
				// Enumeration-neutral: identical shape, nothing was sent.
				h.rateLimiter.record(key, h.now().UTC())
			case errors.Is(err, authsession.ErrRecoveryCooldown):
				writeLocalizedError(w, r, http.StatusTooManyRequests, "EMAIL_RESEND_COOLDOWN", "please wait before requesting another code")
				return
			case errors.Is(err, authsession.ErrRecoverySendFailed):
				writeLocalizedError(w, r, http.StatusBadGateway, "EMAIL_SEND_FAILED", "the reset email could not be sent")
				return
			default:
				writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not process the recovery request")
				return
			}
		}
		writeJSON(w, http.StatusAccepted, map[string]string{"status": "dispatched"})
	}
}

// complete validates the emailed code (+ second factor for MFA accounts) and
// swaps the password in one shot. Failure paths consume attempts; success
// lands through UpdateUser semantics (token_version bump revokes every live
// session) and returns to the sign-in page WITHOUT issuing tokens (D-001 §4).
func (h *recoveryHandler) complete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Account          string `json:"account"`
			Code             string `json:"code"`
			NewPassword      string `json:"newPassword"`
			SecondFactorCode string `json:"secondFactorCode"`
			RecoveryCode     string `json:"recoveryCode"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_RECOVERY_BODY", "body must be JSON with account, code and newPassword")
			return
		}
		if strings.TrimSpace(body.Account) == "" || strings.TrimSpace(body.Code) == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_RECOVERY_BODY", "account and code are required")
			return
		}
		key := h.limiterKey(r, body.Account)
		if !h.rateLimiter.allow(key, h.now().UTC()) {
			if sec := h.rateLimiter.retryAfterSeconds(key, h.now().UTC()); sec > 0 {
				w.Header().Set("Retry-After", strconv.Itoa(sec))
			}
			writeLocalizedError(w, r, http.StatusTooManyRequests, "RATE_LIMITED", "too many recovery attempts; try again later")
			return
		}
		target, err := h.repo.ResolveRecoveryTarget(body.Account)
		if err != nil {
			// Unknown identifier: uniform invalid-code answer (no oracle).
			writeLocalizedError(w, r, http.StatusBadRequest, "RECOVERY_CODE_INVALID", "recovery code is invalid")
			return
		}

		outcome, err := h.repo.EvaluateRecoveryCode(target.UserID, body.Code, h.now().UTC())
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not process the recovery request")
			return
		}
		switch outcome {
		case authsession.RecoveryMatch:
			// continue below
		case authsession.RecoveryExpired:
			h.repo.DropStaleRecoveryChallenge(target.UserID, h.now().UTC())
			writeLocalizedError(w, r, http.StatusBadRequest, "RECOVERY_CODE_EXPIRED", "recovery code expired; request a new one")
			return
		case authsession.RecoveryNotPending:
			writeLocalizedError(w, r, http.StatusBadRequest, "RECOVERY_CODE_INVALID", "recovery code is invalid")
			return
		default: // mismatch
			h.failAttempt(w, r, target.UserID)
			return
		}

		// Second factor for enrolled accounts (V-F076 anti-bypass): demanded
		// AFTER the email code matched, BEFORE any password mutation.
		if h.gate != nil && h.gate.Required(target.UserID) {
			code2, rec2 := strings.TrimSpace(body.SecondFactorCode), strings.TrimSpace(body.RecoveryCode)
			if code2 == "" && rec2 == "" {
				writeLocalizedError(w, r, http.StatusBadRequest, "RECOVERY_SECOND_FACTOR_REQUIRED", "your second factor is required to finish recovery")
				return
			}
			if verr := h.gate.VerifySecondFactor(target.UserID, code2, rec2, h.now().UTC()); verr != nil {
				// A wrong second factor burns a challenge attempt too: the
				// 6-digit code alone must never unlock unlimited TOTP guesses.
				h.failAttemptMFA(w, r, target.UserID, verr)
				return
			}
		}

		newPassword := body.NewPassword
		length := len([]byte(newPassword))
		if length < minPasswordBytes || length > maxPasswordBytes || strings.TrimSpace(newPassword) == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PASSWORD", "new password must be a non-whitespace string of 8 to 72 bytes")
			return
		}
		hash, herr := auth.HashPassword(newPassword, passwordHashCost)
		if herr != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not hash password")
			return
		}
		if cerr := h.repo.CompleteRecovery(target.UserID, hash, target.UserID, h.now().UTC()); cerr != nil {
			if errors.Is(cerr, authsession.ErrRecoveryNotPending) {
				writeLocalizedError(w, r, http.StatusBadRequest, "RECOVERY_CODE_INVALID", "recovery code is invalid")
				return
			}
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not process the recovery request")
			return
		}
		if h.operations != nil {
			h.recordRecoveryEvent(target.UserID)
		}
		if h.notifier != nil {
			NotifyAccountEvent(h.notifier, target.UserID, "account.password-changed", h.now().UTC())
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// failAttempt records the failed code guess and answers uniformly.
func (h *recoveryHandler) failAttempt(w http.ResponseWriter, r *http.Request, userID string) {
	if h.repo.ConsumeRecoveryAttempt(userID) {
		writeLocalizedError(w, r, http.StatusBadRequest, "RECOVERY_CODE_EXPIRED", "recovery code expired; request a new one")
		return
	}
	writeLocalizedError(w, r, http.StatusBadRequest, "RECOVERY_CODE_INVALID", "recovery code is invalid")
}

// failAttemptMFA records the rejected second factor as a challenge attempt
// and maps the underlying mfa error when the budget survives.
func (h *recoveryHandler) failAttemptMFA(w http.ResponseWriter, r *http.Request, userID string, verr error) {
	if h.repo.ConsumeRecoveryAttempt(userID) {
		writeLocalizedError(w, r, http.StatusBadRequest, "RECOVERY_CODE_EXPIRED", "recovery code expired; request a new one")
		return
	}
	switch {
	case errors.Is(verr, ErrMFAInvalid), errors.Is(verr, ErrMFANotEnrolled), errors.Is(verr, ErrMFAPendingOnly):
		writeLocalizedError(w, r, http.StatusBadRequest, "MFA_INVALID", "invalid second-factor code")
	default:
		writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not process the recovery request")
	}
}

// recordRecoveryEvent appends the audit row under the existing
// account.password-change event family (D-001 §5: no CHECK-rebuild migration);
// the detail action distinguishes the recovery path. Actor name resolves
// best-effort and never blocks the 204.
func (h *recoveryHandler) recordRecoveryEvent(userID string) {
	if h.operations == nil {
		return
	}
	actorName := ""
	if u, err := h.repo.UserByID(userID); err == nil && u != nil {
		actorName = u.Name
	}
	detail, derr := operationlog.NewDetail("password-recovery", nil, map[string]any{"userId": userID})
	if derr != nil {
		return
	}
	op := operationlog.Operation{
		ID:        newOperationID(),
		Event:     operationlog.EventAccountPasswordChange,
		ActorID:   userID,
		ActorName: actorName,
		Detail:    &detail,
		CreatedAt: time.Now().UTC(),
	}
	if err := h.operations.RecordOperation(op); err != nil {
		slog.Error("operation log recovery event write failed", "event", op.Event, "err", err)
	}
}
