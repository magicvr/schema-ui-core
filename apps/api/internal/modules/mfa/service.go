// Package mfa provides the admin.mfa module surface (S-10 · GOAL-017 D-002):
// the TOTP second-factor login gate (handler.MFAVerifier), the self-service
// enrollment/confirm/disable/recovery surface and the admin reset. The TOTP
// algorithm is self-implemented (totp.go, RFC 6238 vectors) and secrets are
// encrypted at rest with AES-256-GCM under an HKDF-derived key from the
// server secret.
package mfa

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/hkdf"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/mfa/store"
)

// Proof TTL and the exhaustion threshold (D-002 §3): a proof dies after five
// consecutive second-factor failures.
const (
	proofTTL         = 5 * time.Minute
	proofFailLimit   = 5
	recoveryCodeLen  = 8
	recoveryCodeAlphabet = "23456789ABCDEFGHJKMNPQRSTUVWXYZ"
)

// Domain sentinels are the handler-package MFA errors (module → handler
// direction, captcha precedent): the handler maps them to the frozen wire
// codes without an import cycle.
var (
	ErrProofExpired   = handler.ErrMFAProofExpired
	ErrProofExhausted = handler.ErrMFAProofExhausted
	ErrMFAInvalid     = handler.ErrMFAInvalid
	ErrNotEnrolled    = handler.ErrMFANotEnrolled
	ErrPendingOnly    = handler.ErrMFAPendingOnly
	ErrActive         = handler.ErrMFAActive
)

// Service is the MFA domain service. It satisfies handler.MFAVerifier
// structurally (the direction is module → handler).
type Service struct {
	repo *store.Repository
	key  []byte // AES-GCM key derived from the server secret (HKDF)
}

// NewService constructs the MFA service. serverSecret is the same value used
// to sign JWTs (config secret); the TOTP key is derived with HKDF context
// "mfa/totp" (D-002 §2 — no dedicated key-management surface exists yet).
func NewService(repo *store.Repository, serverSecret []byte) *Service {
	key := make([]byte, 32)
	hk := hkdf.New(sha256.New, serverSecret, nil, []byte("mfa/totp"))
	if _, err := io.ReadFull(hk, key); err != nil {
		panic("mfa: derive totp key: " + err.Error())
	}
	return &Service{repo: repo, key: key}
}

// --- handler.MFAVerifier ---

// Required reports whether the user must complete a second factor before token
// issuance (only active enrollments gate login; pending ones do not — A-004
// F-001: no self-lock after enrollment). Fail-closed (W7 F-001): a storage
// error is NOT treated as "no MFA" — only a definitive ErrNotFound means the
// user has no enrollment. Any other read error forces the login into the
// second-factor branch, which then fails closed until the store can be read.
func (s *Service) Required(userID string) bool {
	// Defensive: a typed-nil *Service must never panic the login path (the
	// composition root passes a true nil interface when admin.mfa is off).
	if s == nil || s.repo == nil {
		return false
	}
	st, err := s.repo.GetState(userID)
	if errors.Is(err, store.ErrNotFound) {
		return false
	}
	if err != nil {
		return true
	}
	return st.Status == "active"
}

// BeginChallenge issues a one-time login proof for a pending login.
func (s *Service) BeginChallenge(userID string, now time.Time) (string, error) {
	st, err := s.repo.GetState(userID)
	if err != nil || st.Status != "active" {
		return "", ErrNotEnrolled
	}
	p, err := s.repo.CreateProof(userID, now.Add(proofTTL), now)
	if err != nil {
		return "", err
	}
	return p.ID, nil
}

// Verify completes the second factor: proof validity (existence/expiry/one-shot)
// plus a TOTP code or a one-time recovery code. Returns the verified user id.
// Same-window TOTP replay is rejected via last_used_step (D-002 §2/§6).
func (s *Service) Verify(proofID, code, recoveryCode string, now time.Time) (string, error) {
	proof, err := s.repo.GetProof(proofID)
	if err != nil {
		return "", ErrProofExpired // unknown proofs are indistinguishable from consumed ones
	}
	if now.After(proof.ExpiresAt) {
		_ = s.repo.DeleteProof(proofID)
		return "", ErrProofExpired
	}
	if proof.FailCount >= proofFailLimit {
		_ = s.repo.DeleteProof(proofID)
		return "", ErrProofExhausted
	}
	st, err := s.repo.GetState(proof.UserID)
	if err != nil || st.Status != "active" {
		return "", ErrNotEnrolled
	}
	ok := false
	if strings.TrimSpace(recoveryCode) != "" {
		ok = s.consumeRecoveryCode(st, recoveryCode, now)
	} else {
		step, valid := ValidateTotp(decryptSecret(s.key, st.SecretCiphertext), code, now, totpWindow, st.LastUsedStep)
		if valid {
			if err := s.repo.SetLastUsedStep(proof.UserID, step, now); err != nil {
				return "", err
			}
			ok = true
		}
	}
	if !ok {
		_ = s.repo.IncrementProofFailures(proofID, now)
		if proof.FailCount+1 >= proofFailLimit {
			_ = s.repo.DeleteProof(proofID)
			return "", ErrProofExhausted
		}
		return "", ErrMFAInvalid
	}
	_ = s.repo.DeleteProof(proofID) // one-shot
	return proof.UserID, nil
}

// --- self-service / admin surface ---

// Status returns the enrollment state for the current user.
func (s *Service) Status(userID string) (enabled bool, enrolledAt time.Time, err error) {
	st, err := s.repo.GetState(userID)
	if errors.Is(err, store.ErrNotFound) {
		return false, time.Time{}, nil
	}
	if err != nil {
		return false, time.Time{}, err
	}
	return st.Status == "active", st.CreatedAt, nil
}

// Enroll creates a pending enrollment and returns the one-time secret payload.
func (s *Service) Enroll(userID, name string, now time.Time) (secretBase32, otpauth string, recoveryCodes []string, err error) {
	// A-007 F-002: an ACTIVE enrollment must not be overwritten without the
	// second factor (disabling first requires a valid code/recovery). Pending
	// enrollments may be re-enrolled (A-005 recommended).
	if existing, err := s.repo.GetState(userID); err == nil {
		if existing.Status == "active" {
			return "", "", nil, ErrActive
		}
	} else if !errors.Is(err, store.ErrNotFound) {
		return "", "", nil, err
	}
	secret, err := GenerateSecret()
	if err != nil {
		return "", "", nil, err
	}
	codes := make([]string, 10)
	hashes := make([]string, 10)
	for i := range codes {
		codes[i] = randomRecoveryCode()
		h, err := bcrypt.GenerateFromPassword([]byte(codes[i]), bcrypt.DefaultCost)
		if err != nil {
			return "", "", nil, fmt.Errorf("mfa: hash recovery code: %w", err)
		}
		hashes[i] = string(h)
	}
	hashJSON, err := json.Marshal(hashes)
	if err != nil {
		return "", "", nil, err
	}
	ciphertext := encryptSecret(s.key, []byte(secret))
	if err := s.repo.UpsertPending(userID, ciphertext, string(hashJSON), now); err != nil {
		return "", "", nil, err
	}
	return secret, otpauthURL("Schema UI Core", name, secret), codes, nil
}

// Confirm activates a pending enrollment after a correct TOTP code.
func (s *Service) Confirm(userID, code string, now time.Time) error {
	st, err := s.repo.GetState(userID)
	if errors.Is(err, store.ErrNotFound) {
		return ErrNotEnrolled
	}
	if err != nil {
		return err
	}
	if st.Status == "active" {
		return ErrActive
	}
	if _, ok := ValidateTotp(decryptSecret(s.key, st.SecretCiphertext), code, now, totpWindow, 0); !ok {
		return ErrMFAInvalid
	}
	if err := s.repo.SetLastUsedStep(userID, now.Unix()/totpPeriodSeconds, now); err != nil {
		return err
	}
	return s.repo.Activate(userID, now)
}

// Disable removes the enrollment after a valid code/recovery (caller then
// bumps token_version and revokes sessions — A-004 F-002 parity).
func (s *Service) Disable(userID, code, recoveryCode string, now time.Time) error {
	if err := s.requireActiveSecondFactor(userID, code, recoveryCode, now); err != nil {
		return err
	}
	return s.repo.DeleteState(userID)
}

// RotateRecovery replaces the recovery-code set after a valid code/recovery.
func (s *Service) RotateRecovery(userID, code, recoveryCode string, now time.Time) ([]string, error) {
	if err := s.requireActiveSecondFactor(userID, code, recoveryCode, now); err != nil {
		return nil, err
	}
	codes := make([]string, 10)
	hashes := make([]string, 10)
	for i := range codes {
		codes[i] = randomRecoveryCode()
		h, err := bcrypt.GenerateFromPassword([]byte(codes[i]), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("mfa: hash recovery code: %w", err)
		}
		hashes[i] = string(h)
	}
	hashJSON, err := json.Marshal(hashes)
	if err != nil {
		return nil, err
	}
	if err := s.repo.UpdateRecoveryCodes(userID, string(hashJSON), now); err != nil {
		return nil, err
	}
	return codes, nil
}

// AdminReset removes another user's enrollment (caller then bumps
// token_version and revokes sessions only when an ACTIVE enrollment existed —
// W7 F-002: resetting a user without active MFA must not become a generic
// forced-logout primitive). The returned bool reports whether an active
// enrollment was removed; a pending state is removed but reported false.
func (s *Service) AdminReset(userID string) (bool, error) {
	st, err := s.repo.GetState(userID)
	if errors.Is(err, store.ErrNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if err := s.repo.DeleteState(userID); err != nil {
		return false, err
	}
	return st.Status == "active", nil
}

// requireActiveSecondFactor validates a code or recovery code against an
// active enrollment (used by disable / recovery rotation).
func (s *Service) requireActiveSecondFactor(userID, code, recoveryCode string, now time.Time) error {
	st, err := s.repo.GetState(userID)
	if err != nil {
		return ErrNotEnrolled
	}
	if st.Status != "active" {
		return ErrPendingOnly
	}
	if strings.TrimSpace(recoveryCode) != "" {
		if !s.consumeRecoveryCode(st, recoveryCode, now) {
			return ErrMFAInvalid
		}
		return nil
	}
	if _, ok := ValidateTotp(decryptSecret(s.key, st.SecretCiphertext), code, now, totpWindow, st.LastUsedStep); !ok {
		return ErrMFAInvalid
	}
	return nil
}

// consumeRecoveryCode bcrypt-compares the code against the stored hash set and
// removes the matched hash (one-time consumption).
func (s *Service) consumeRecoveryCode(st *store.State, code string, now time.Time) bool {
	var hashes []string
	if err := json.Unmarshal([]byte(st.RecoveryCodesHash), &hashes); err != nil {
		return false
	}
	for i, h := range hashes {
		if bcrypt.CompareHashAndPassword([]byte(h), []byte(strings.TrimSpace(code))) == nil {
			remaining := append(append([]string{}, hashes[:i]...), hashes[i+1:]...)
			next, err := json.Marshal(remaining)
			if err != nil {
				return false
			}
			if err := s.repo.UpdateRecoveryCodes(st.UserID, string(next), now); err != nil {
				return false
			}
			return true
		}
	}
	return false
}

func randomRecoveryCode() string {
	raw := make([]byte, recoveryCodeLen)
	_, _ = rand.Read(raw)
	var b strings.Builder
	for _, v := range raw {
		b.WriteByte(recoveryCodeAlphabet[int(v)%len(recoveryCodeAlphabet)])
	}
	return b.String()
}

// encryptSecret seals a plaintext secret with AES-256-GCM: nonce || ciphertext,
// base64 encoded.
func encryptSecret(key, plain []byte) string {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic("mfa: aes cipher: " + err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic("mfa: gcm: " + err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		panic("mfa: nonce: " + err.Error())
	}
	sealed := gcm.Seal(nonce, nonce, plain, nil)
	return base64.StdEncoding.EncodeToString(sealed)
}

// decryptSecret reverses encryptSecret; an undecryptable blob yields "" (a
// failed second factor, never a panic).
func decryptSecret(key []byte, encoded string) string {
	sealed, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return ""
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return ""
	}
	if len(sealed) < gcm.NonceSize() {
		return ""
	}
	plain, err := gcm.Open(nil, sealed[:gcm.NonceSize()], sealed[gcm.NonceSize():], nil)
	if err != nil {
		return ""
	}
	return string(plain)
}

var _ auth.MFAEnforcer = (*Service)(nil)
