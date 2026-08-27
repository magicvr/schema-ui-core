// Self-service password recovery domain (workspace-019 R2 · GOAL-003 D-001):
// the challenge state machine over migration 0056, consuming ONLY the
// kernel.MailSender port. Frozen semantics (Root D-002 + GOAL-002 D-001 §1):
//   - proof form is the VP-018-isomorphic 6-digit code (sha256 at rest,
//     constant-time compare); delivery rides the ONE composed MailSender;
//   - TTL 10 min, resend cooldown 60 s, 5 failed attempts void the challenge;
//   - accounts without a bound+verified email have NO self-recovery path
//     (I-006 · controlled negative, admin reset stays the privilege path);
//   - completion hashes at the handler layer and lands through UpdateUser,
//     which bumps token_version and revokes refresh tokens atomically (§4).
package authsession

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// Frozen I-002 numbers (user adjudication 2026-08-25 · Root D-002): identical
// to the VP-018 verification constants on purpose — one mental model, shared
// bookkeeping shape.
const (
	recoveryCodeTTL           = 10 * time.Minute
	recoveryResendCooldown    = 60 * time.Second
	recoveryMaxFailedAttempts = 5
)

// Sentinel errors mapped to catalog codes by the HTTP surface. Unknown
// account / unbound email / disabled collapse into ErrRecoveryNotAvailable so
// the pre-auth start surface can answer with one enumeration-neutral shape.
var (
	ErrRecoveryNotAvailable = errors.New("authsession: account has no self-recovery path")
	ErrRecoveryCooldown     = errors.New("authsession: recovery resend cooldown active")
	ErrRecoverySendFailed   = errors.New("authsession: recovery email could not be sent")
	ErrRecoveryCodeInvalid  = errors.New("authsession: recovery code is invalid")
	ErrRecoveryCodeExpired  = errors.New("authsession: recovery code expired or voided")
	ErrRecoveryNotPending   = errors.New("authsession: no live recovery challenge")
)

const recoveryMailSubject = "密码重置验证码 · Password reset code"

func recoveryMailBody(code string, expires time.Time) string {
	return fmt.Sprintf(
		"您的密码重置验证码：%s\n有效期至 %s（10 分钟）。\n\n"+
			"Your password reset code: %s\nValid for 10 minutes (until %s).\n\n"+
			"若非本人操作，请忽略本邮件并考虑修改密码。/ If you did not request this, ignore this message.\n",
		code, expires.Format("2006-01-02 15:04:05 MST"), code, expires.Format(time.RFC3339),
	)
}

func sendRecoveryMail(ctx context.Context, sender kernel.MailSender, to, code string, expires time.Time) error {
	msg := kernel.MailMessage{
		To:       to,
		Subject:  recoveryMailSubject,
		TextBody: recoveryMailBody(code, expires),
	}
	if err := msg.Validate(); err != nil {
		return ErrRecoverySendFailed
	}
	if err := sender.Send(ctx, msg); err != nil {
		return fmt.Errorf("%w: %v", ErrRecoverySendFailed, err)
	}
	return nil
}

// RecoveryTarget is the resolved self-recovery subject.
type RecoveryTarget struct {
	UserID  string
	Enabled bool
}

// ResolveRecoveryTarget locates the unique account behind a start/complete
// identifier: exact username first (the login-page key), then the account's
// bound+VERIFIED address folded case-insensitively. Unresolved / unbound /
// unverified surfaces as ErrRecoveryNotAvailable; callers decide whether that
// means silence (start, anti-enumeration) or a uniform invalid code (complete).
func (r *Repository) ResolveRecoveryTarget(identifier string) (*RecoveryTarget, error) {
	id := trimIdentifier(identifier)
	if id == "" {
		return nil, ErrRecoveryNotAvailable
	}
	var out RecoveryTarget
	err := r.withTx("resolve recovery target", func(tx kernel.Tx) error {
		var email *string
		var status *string
		var enabled int
		err := tx.QueryRow(context.Background(),
			`SELECT id, enabled, email, email_status FROM users WHERE username = ?`, id,
		).Scan(&out.UserID, &enabled, &email, &status)
		switch {
		case err == nil:
			out.Enabled = enabled != 0
			return nil
		case !errors.Is(err, kernel.ErrNoRows):
			return fmt.Errorf("lookup user %q: %w", id, err)
		}
		// Not a username: try the verified-email projection (lower fold).
		err = tx.QueryRow(context.Background(),
			`SELECT id, enabled FROM users WHERE email IS NOT NULL AND email_status = 'verified' AND lower(email) = lower(?)`, id,
		).Scan(&out.UserID, &enabled)
		if errors.Is(err, kernel.ErrNoRows) {
			return ErrRecoveryNotAvailable
		}
		if err != nil {
			return fmt.Errorf("lookup recovery email %q: %w", id, err)
		}
		out.Enabled = enabled != 0
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// StartRecovery resolves the account, honours the frozen cooldown, replaces
// any prior challenge and dispatches the code through the given MailSender.
// TWO-PHASE like BindEmail: the staged challenge is compensated away when the
// dispatch fails, leaving no live challenge without a dispatched mail.
// Accounts without a self-recovery path (unknown / unbound / unverified /
// disabled) fail with ErrRecoveryNotAvailable BEFORE any mail leaves.
func (r *Repository) StartRecovery(identifier string, sender kernel.MailSender, now time.Time) error {
	target, err := r.ResolveRecoveryTarget(identifier)
	if err != nil {
		return err
	}
	if !target.Enabled {
		return ErrRecoveryNotAvailable
	}
	var address string
	if err := r.withTx("read recovery address", func(tx kernel.Tx) error {
		return tx.QueryRow(context.Background(),
			`SELECT email FROM users WHERE id = ? AND email IS NOT NULL AND email_status = 'verified'`,
			target.UserID,
		).Scan(&address)
	}); err != nil {
		if errors.Is(err, kernel.ErrNoRows) {
			return ErrRecoveryNotAvailable
		}
		return fmt.Errorf("read verified email: %w", err)
	}

	code, err := generateEmailCode()
	if err != nil {
		return err
	}
	expires := now.Add(recoveryCodeTTL)

	// Cooldown + idempotent replace in ONE transaction (0056 PK on user_id).
	if err := r.withTx("stage recovery challenge", func(tx kernel.Tx) error {
		var sentAt int64
		scanErr := tx.QueryRow(context.Background(),
			`SELECT sent_at FROM password_recovery_challenges WHERE user_id = ?`, target.UserID,
		).Scan(&sentAt)
		if scanErr == nil && now.Unix()-sentAt < int64(recoveryResendCooldown/time.Second) {
			return ErrRecoveryCooldown
		}
		if scanErr != nil && !errors.Is(scanErr, kernel.ErrNoRows) {
			return fmt.Errorf("read prior challenge: %w", scanErr)
		}
		if _, err := tx.Exec(context.Background(), `DELETE FROM password_recovery_challenges WHERE user_id = ?`, target.UserID); err != nil {
			return fmt.Errorf("clear old challenge: %w", err)
		}
		if _, err := tx.Exec(context.Background(),
			`INSERT INTO password_recovery_challenges (user_id, code_hash, expires_at, sent_at, attempt_count)
			 VALUES (?, ?, ?, ?, 0)`,
			target.UserID, hashCode(code), expires.Unix(), now.Unix(),
		); err != nil {
			return fmt.Errorf("store recovery challenge: %w", err)
		}
		return nil
	}); err != nil {
		return err
	}

	if serr := sendRecoveryMail(context.Background(), sender, address, code, expires); serr != nil {
		if cerr := r.withTx("compensate recovery challenge", func(tx kernel.Tx) error {
			_, err := tx.Exec(context.Background(), `DELETE FROM password_recovery_challenges WHERE user_id = ?`, target.UserID)
			return err
		}); cerr != nil {
			return fmt.Errorf("%w: %v (compensation failed: %v)", ErrRecoverySendFailed, serr, cerr)
		}
		return serr
	}
	return nil
}

// RecoveryOutcome classifies an evaluate pass (mirrors verificationOutcome;
// exported because the completion gate orchestrates across the handler seam).
type RecoveryOutcome int

const (
	RecoveryMatch RecoveryOutcome = iota
	RecoveryMismatch
	RecoveryExpired
	RecoveryNotPending
)

// EvaluateRecoveryCode runs the read-only classification pass for the
// completion gate. Only classification happens here; every mutation path
// (attempt bumping, voiding, consumption) is explicit at the call site so a
// rollback cannot erase failure-path bookkeeping (VerifyEmail precedent).
func (r *Repository) EvaluateRecoveryCode(userID, rawCode string, now time.Time) (RecoveryOutcome, error) {
	code := trimCode(rawCode)
	outcome := RecoveryNotPending
	err := r.withTx("evaluate recovery code", func(tx kernel.Tx) error {
		var codeHash string
		var expiresAt int64
		err := tx.QueryRow(context.Background(),
			`SELECT code_hash, expires_at FROM password_recovery_challenges WHERE user_id = ?`,
			userID,
		).Scan(&codeHash, &expiresAt)
		if errors.Is(err, kernel.ErrNoRows) {
			outcome = RecoveryNotPending
			return nil
		}
		if err != nil {
			return fmt.Errorf("load recovery challenge: %w", err)
		}
		// W13 F-008 (GOAL-013 A-001): expiry alone must not classify the
		// response. The previous shape answered RECOVERY_CODE_EXPIRED for ANY
		// submitted value once a challenge had aged out, letting anyone who
		// knows an account identifier probe whether that account recently
		// requested a reset. Expiry is only surfaced when the code hash
		// MATCHES (the legitimate holder gets the helpful message); every
		// other value falls through to the uniform invalid classification.
		if now.Unix() >= expiresAt {
			if subtle.ConstantTimeCompare([]byte(hashCode(code)), []byte(codeHash)) == 1 {
				outcome = RecoveryExpired
			} else {
				outcome = RecoveryMismatch
			}
			return nil
		}
		if subtle.ConstantTimeCompare([]byte(hashCode(code)), []byte(codeHash)) == 1 {
			outcome = RecoveryMatch
		} else {
			outcome = RecoveryMismatch
		}
		return nil
	})
	if err != nil {
		return RecoveryNotPending, err
	}
	return outcome, nil
}

// ConsumeRecoveryAttempt records one failed completion attempt in its own
// transaction and reports whether the challenge reached the void threshold
// (deleted — a resend is required). Called for wrong codes AND rejected
// second factors alike: the total guess budget per challenge lifecycle stays ≤5.
func (r *Repository) ConsumeRecoveryAttempt(userID string) bool {
	voided := false
	_ = r.withTx("record recovery failure", func(tx kernel.Tx) error {
		if _, err := tx.Exec(context.Background(),
			`UPDATE password_recovery_challenges SET attempt_count = attempt_count + 1 WHERE user_id = ?`,
			userID); err != nil {
			return err
		}
		var attempts int
		if err := tx.QueryRow(context.Background(),
			`SELECT attempt_count FROM password_recovery_challenges WHERE user_id = ?`, userID,
		).Scan(&attempts); err != nil {
			return err
		}
		if attempts >= recoveryMaxFailedAttempts {
			voided = true
			_, err := tx.Exec(context.Background(),
				`DELETE FROM password_recovery_challenges WHERE user_id = ?`, userID)
			return err
		}
		return nil
	})
	return voided
}

// DropStaleRecoveryChallenge removes an expired challenge (best-effort, same
// as VerifyEmail's stale path) so the next start starts clean.
func (r *Repository) DropStaleRecoveryChallenge(userID string, now time.Time) {
	_ = r.withTx("drop stale recovery challenge", func(tx kernel.Tx) error {
		_, err := tx.Exec(context.Background(),
			`DELETE FROM password_recovery_challenges WHERE user_id = ? AND expires_at <= ?`,
			userID, now.Unix())
		return err
	})
}

// CompleteRecovery consumes the live challenge (guarded: exactly one row must
// go — replaying a consumed code cannot re-enter) and then hands the hashed
// replacement password to UpdateUser, whose transaction bumps token_version,
// revokes every refresh token and clears must_change_password. Consumption
// happens FIRST by design: a failure afterwards forces a fresh start instead
// of leaving a live challenge against a half-mutated account.
func (r *Repository) CompleteRecovery(userID, passwordHash string, actorID string, now time.Time) error {
	consumed := false
	err := r.withTx("consume recovery challenge", func(tx kernel.Tx) error {
		result, err := tx.Exec(context.Background(),
			`DELETE FROM password_recovery_challenges WHERE user_id = ?`, userID)
		if err != nil {
			return fmt.Errorf("consume recovery challenge: %w", err)
		}
		rows, rerr := result.RowsAffected()
		if rerr != nil {
			return fmt.Errorf("consume recovery challenge rows: %w", rerr)
		}
		consumed = rows == 1
		return nil
	})
	if err != nil {
		return err
	}
	if !consumed {
		return ErrRecoveryNotPending
	}
	mustChange := false
	_, err = r.UpdateUser(userID, UserPatch{PasswordHash: &passwordHash, MustChangePassword: &mustChange}, actorID, now)
	return err
}

func trimIdentifier(raw string) string {
	return strings.TrimSpace(raw)
}

func trimCode(raw string) string { return strings.TrimSpace(raw) }
