// Account email identity binding/verification flow (workspace-018 R3 ·
// GOAL-004 D-001): the state machine over the R2 schema (migration 0054
// columns + 0055 challenge table), consuming ONLY the kernel.MailSender port.
//
// Frozen semantics (GOAL-002 D-001 §2/§4/§5/§6):
//   - bind reserves the unique slot immediately (pending occupies; the
//     lower(email) unique index backstops races);
//   - rebind = overwrite to the new address and reset to pending (the old
//     slot is released by the overwrite itself);
//   - verification is a 6-digit code, TTL 10 min, resend cooldown 60 s,
//     constant-time compare, 5 failed attempts void the challenge;
//   - there is NO path to verified without a delivered code (admin prefill
//     lands as pending too).
package authsession

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// Frozen I-005 numbers (user adjudication 2026-08-24 · GOAL-004 D-001).
const (
	emailCodeTTL          = 10 * time.Minute
	emailResendCooldown   = 60 * time.Second
	emailMaxFailedAttempts = 5
)

// Sentinel errors mapped to catalog codes by the HTTP surface.
var (
	ErrEmailInvalid         = errors.New("authsession: invalid email address")
	ErrEmailTaken           = errors.New("authsession: email already bound or pending on another account")
	ErrEmailNotPending      = errors.New("authsession: account has no pending email verification")
	ErrEmailCodeInvalid     = errors.New("authsession: verification code is invalid")
	ErrEmailCodeExpired     = errors.New("authsession: verification code expired")
	ErrEmailResendCooldown  = errors.New("authsession: resend cooldown active")
	ErrEmailSendFailed      = errors.New("authsession: verification email could not be sent")
)

// normalizeEmailInput trims and shape-checks a candidate address. Full RFC
// 5322 validation belongs to MailMessage.Validate at the send boundary; this
// guard only rejects obviously malformed input before it touches storage.
func normalizeEmailInput(raw string) (string, error) {
	email := strings.TrimSpace(raw)
	if email == "" || len(email) > 254 || strings.ContainsAny(email, " \t\r\n") {
		return "", ErrEmailInvalid
	}
	at := strings.Index(email, "@")
	if at <= 0 || at == len(email)-1 || strings.Count(email, "@") != 1 {
		return "", ErrEmailInvalid
	}
	return email, nil
}

func stringPtr(s string) *string { return &s }

// nullIfNil adapts an optional string for SQL storage: nil → NULL.
func nullIfNil(s *string) any {
	if s == nil {
		return nil
	}
	return *s
}

// EmailIdentityState reads back the managed identity triple for surfaces that
// render it (self profile). nil/nil = unbound.
func (r *Repository) EmailIdentityState(userID string) (*string, *string, error) {
	var email, status *string
	err := r.withTx("read email identity", func(tx kernel.Tx) error {
		return tx.QueryRow(context.Background(),
			`SELECT email, email_status FROM users WHERE id = ?`, userID,
		).Scan(&email, &status)
	})
	if err != nil {
		return nil, nil, err
	}
	return email, status, nil
}

func hashCode(emailCode string) string {
	sum := sha256.Sum256([]byte(emailCode))
	return hex.EncodeToString(sum[:])
}

func generateEmailCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", fmt.Errorf("generate email code: %w", err)
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

const emailVerificationMailSubject = "账号邮箱验证码 · Email verification code"

func emailVerificationMailBody(code string, expires time.Time) string {
	return fmt.Sprintf(
		"您的邮箱验证码：%s\n有效期至 %s（10 分钟）。\n\n"+
			"Your email verification code: %s\nValid for 10 minutes (until %s).\n\n"+
			"若非本人操作，请忽略本邮件。/ If you did not request this, ignore this message.\n",
		code, expires.Format("2006-01-02 15:04:05 MST"), code, expires.Format(time.RFC3339),
	)
}

func sendVerificationMail(ctx context.Context, sender kernel.MailSender, to, code string, expires time.Time) error {
	msg := kernel.MailMessage{
		To:       to,
		Subject:  emailVerificationMailSubject,
		TextBody: emailVerificationMailBody(code, expires),
	}
	if err := msg.Validate(); err != nil {
		return ErrEmailInvalid
	}
	if err := sender.Send(ctx, msg); err != nil {
		return fmt.Errorf("%w: %v", ErrEmailSendFailed, err)
	}
	return nil
}

// BindEmail binds rawEmail to userID and starts verification: the address
// occupies its unique slot immediately (pending), any previous challenge is
// replaced, and the code goes out through the given MailSender. Rebinding to
// a new address is the same operation (overwrite releases the old slot).
// Idempotent when the SAME address is already verified. The Send happens
// inside the transaction so a delivery failure rolls the whole bind back
// (no pending row without a dispatched mail); the write window held during
// the synchronous Send is accepted at admin-scale traffic.
func (r *Repository) BindEmail(userID, rawEmail string, sender kernel.MailSender, now time.Time) error {
	email, err := normalizeEmailInput(rawEmail)
	if err != nil {
		return err
	}
	code, err := generateEmailCode()
	if err != nil {
		return err
	}
	expires := now.Add(emailCodeTTL)
	return r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		var currentEmail *string
		var currentStatus *string
		err := tx.QueryRow(context.Background(),
			`SELECT email, email_status FROM users WHERE id = ?`, userID,
		).Scan(&currentEmail, &currentStatus)
		if err != nil {
			return fmt.Errorf("load user %s: %w", userID, err)
		}
		if currentStatus != nil && *currentStatus == "verified" && currentEmail != nil &&
			strings.EqualFold(strings.TrimSpace(*currentEmail), email) {
			return nil // same address already verified: idempotent no-op
		}
		// A-001 F-002: re-binding the SAME pending address is resend semantics
		// (a fresh code to the same inbox) — honor the frozen cooldown. A
		// DIFFERENT address is a rebind and dispatches immediately.
		if currentStatus != nil && *currentStatus == "pending" && currentEmail != nil &&
			strings.EqualFold(strings.TrimSpace(*currentEmail), email) {
			var sentAt int64
			if err := tx.QueryRow(context.Background(),
				`SELECT sent_at FROM email_verification_challenges WHERE user_id = ?`, userID,
			).Scan(&sentAt); err == nil && now.Unix()-sentAt < int64(emailResendCooldown/time.Second) {
				return ErrEmailResendCooldown
			}
		}
		var taken int
		if err := tx.QueryRow(context.Background(),
			`SELECT COUNT(*) FROM users WHERE id <> ? AND email IS NOT NULL AND lower(email) = lower(?)`,
			userID, email,
		).Scan(&taken); err != nil {
			return fmt.Errorf("check email slot: %w", err)
		}
		if taken > 0 {
			return ErrEmailTaken
		}
		if _, err := tx.Exec(context.Background(),
			`UPDATE users SET email = ?, email_status = 'pending', updated_at = ? WHERE id = ?`,
			email, now.Unix(), userID,
		); err != nil {
			return fmt.Errorf("set pending email: %w", err)
		}
		if _, err := tx.Exec(context.Background(), `DELETE FROM email_verification_challenges WHERE user_id = ?`, userID); err != nil {
			return fmt.Errorf("clear old challenge: %w", err)
		}
		if _, err := tx.Exec(context.Background(),
			`INSERT INTO email_verification_challenges (user_id, code_hash, expires_at, sent_at, attempt_count)
			 VALUES (?, ?, ?, ?, 0)`,
			userID, hashCode(code), expires.Unix(), now.Unix(),
		); err != nil {
			return fmt.Errorf("store challenge: %w", err)
		}
		return sendVerificationMail(context.Background(), sender, email, code, expires)
	})
}

// VerifyEmail consumes the code: match → verified (challenge deleted);
// mismatch → attempt bookkeeping in ITS OWN transaction (the runner rolls a
// failed transaction back, so failure-path bookkeeping cannot share the main
// one); ≥5 failures void the challenge; stale → expired (challenge dropped).
func (r *Repository) VerifyEmail(userID, rawCode string, now time.Time) error {
	code := strings.TrimSpace(rawCode)
	if len(code) != 6 {
		return ErrEmailCodeInvalid
	}
	outcome, matched, err := r.evaluateVerification(userID, code, now)
	if err != nil {
		return fmt.Errorf("verify email: %w", err)
	}
	if matched {
		return nil
	}
	switch outcome {
	case verificationExpired:
		_ = r.runner.Run(context.Background(), func(tx kernel.Tx) error {
			_, err := tx.Exec(context.Background(),
				`DELETE FROM email_verification_challenges WHERE user_id = ? AND expires_at <= ?`,
				userID, now.Unix())
			return err
		})
		return ErrEmailCodeExpired
	case verificationMismatch:
		if r.registerFailedAttempt(userID, emailMaxFailedAttempts) {
			return ErrEmailCodeExpired // voided: a resend is required
		}
		return ErrEmailCodeInvalid
	default:
		return ErrEmailNotPending
	}
}

type verificationOutcome int

const (
	verificationMatch verificationOutcome = iota
	verificationMismatch
	verificationExpired
	verificationNotPending
)

// evaluateVerification runs the read/consume pass. Only the SUCCESS path
// mutates (mark verified + consume challenge); every failure path returns
// without writes so its rollback is harmless. Controlled negative outcomes
// return (outcome, false, nil); hard storage failures return a non-nil error.
func (r *Repository) evaluateVerification(userID, code string, now time.Time) (verificationOutcome, bool, error) {
	outcome := verificationNotPending
	matched := false
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		var status *string
		if err := tx.QueryRow(context.Background(),
			`SELECT email_status FROM users WHERE id = ? AND email IS NOT NULL`, userID,
		).Scan(&status); err != nil {
			outcome = verificationNotPending
			return err
		}
		if status == nil || *status != "pending" {
			outcome = verificationNotPending
			return errNotSentinel
		}
		var codeHash string
		var expiresAt int64
		err := tx.QueryRow(context.Background(),
			`SELECT code_hash, expires_at FROM email_verification_challenges WHERE user_id = ?`,
			userID,
		).Scan(&codeHash, &expiresAt)
		if err != nil {
			// Missing row (never sent / already consumed): controlled negative.
			outcome = verificationNotPending
			return errNotSentinel
		}
		if now.Unix() >= expiresAt {
			outcome = verificationExpired
			return errNotSentinel
		}
		if subtle.ConstantTimeCompare([]byte(hashCode(code)), []byte(codeHash)) != 1 {
			outcome = verificationMismatch
			return errNotSentinel
		}
		result, err := tx.Exec(context.Background(),
			`UPDATE users SET email_status = 'verified', updated_at = ? WHERE id = ? AND email IS NOT NULL AND email_status = 'pending'`,
			now.Unix(), userID,
		)
		if err != nil {
			return fmt.Errorf("mark verified: %w", err)
		}
		if rows, _ := result.RowsAffected(); rows != 1 {
			outcome = verificationNotPending
			return errNotSentinel
		}
		if _, err := tx.Exec(context.Background(),
			`DELETE FROM email_verification_challenges WHERE user_id = ?`, userID,
		); err != nil {
			return fmt.Errorf("consume challenge: %w", err)
		}
		outcome = verificationMatch
		matched = true
		return nil
	})
	if err != nil && !errors.Is(err, errNotSentinel) {
		return verificationNotPending, false, fmt.Errorf("evaluate verification: %w", err)
	}
	return outcome, matched, nil
}

// registerFailedAttempt bumps the failure counter in its own transaction and
// reports whether the challenge reached the void threshold (deleted).
func (r *Repository) registerFailedAttempt(userID string, maxAttempts int) bool {
	_ = r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		if _, err := tx.Exec(context.Background(),
			`UPDATE email_verification_challenges SET attempt_count = attempt_count + 1 WHERE user_id = ?`,
			userID); err != nil {
			return err
		}
		var attempts int
		if err := tx.QueryRow(context.Background(),
			`SELECT attempt_count FROM email_verification_challenges WHERE user_id = ?`, userID,
		).Scan(&attempts); err != nil {
			return err
		}
		if attempts >= maxAttempts {
			_, err := tx.Exec(context.Background(),
				`DELETE FROM email_verification_challenges WHERE user_id = ?`, userID)
			return err
		}
		return nil
	})
	// Read-back is best-effort; a voided challenge surfaces as expired on the
	// next evaluate pass either way.
	var attempts int
	_ = r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		return tx.QueryRow(context.Background(),
			`SELECT COALESCE(attempt_count, 0) FROM email_verification_challenges WHERE user_id = ?`, userID,
		).Scan(&attempts)
	})
	return attempts == 0
}

// errNotSentinel is a control-flow marker: abort a transaction for a
// controlled negative outcome without mapping to a domain error.
var errNotSentinel = errors.New("authsession: controlled verification outcome")

// ResendEmailCode issues a fresh code for an existing pending address,
// honoring the frozen resend cooldown.
func (r *Repository) ResendEmailCode(userID string, sender kernel.MailSender, now time.Time) error {
	code, err := generateEmailCode()
	if err != nil {
		return err
	}
	expires := now.Add(emailCodeTTL)
	return r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		var status *string
		var email *string
		if err := tx.QueryRow(context.Background(),
			`SELECT email_status, email FROM users WHERE id = ? AND email IS NOT NULL`, userID,
		).Scan(&status, &email); err != nil || status == nil || *status != "pending" || email == nil {
			return ErrEmailNotPending
		}
		var sentAt int64
		err := tx.QueryRow(context.Background(),
			`SELECT sent_at FROM email_verification_challenges WHERE user_id = ?`, userID,
		).Scan(&sentAt)
		switch {
		case err == nil:
			if now.Unix()-sentAt < int64(emailResendCooldown/time.Second) {
				return ErrEmailResendCooldown
			}
		default:
			// No live challenge (expired/consumed): issuing a fresh one is fine.
		}
		if _, err := tx.Exec(context.Background(), `DELETE FROM email_verification_challenges WHERE user_id = ?`, userID); err != nil {
			return fmt.Errorf("clear old challenge: %w", err)
		}
		if _, err := tx.Exec(context.Background(),
			`INSERT INTO email_verification_challenges (user_id, code_hash, expires_at, sent_at, attempt_count)
			 VALUES (?, ?, ?, ?, 0)`,
			userID, hashCode(code), expires.Unix(), now.Unix(),
		); err != nil {
			return fmt.Errorf("store refreshed challenge: %w", err)
		}
		return sendVerificationMail(context.Background(), sender, *email, code, expires)
	})
}
