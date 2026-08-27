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
	"log/slog"
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
	// prevKey is the HKDF key derived from the previous JWT secret during a
	// VP-016 rotation window (W11 F-004); nil/empty when no previous secret
	// is configured keeps exact single-key behavior.
	prevKey []byte
}

// NewService constructs the MFA service. serverSecret is the same value used
// to sign JWTs (config secret); the TOTP key is derived with HKDF context
// "mfa/totp" (D-002 §2 — no dedicated key-management surface exists yet).
// previousSecret, when non-empty, is the VP-016 rotation window fallback: a
// ciphertext sealed under the previous secret stays decryptable after
// AUTH_JWT_SECRET rotation, and successful second-factor verifications
// re-wrap the secret under the current key (W11 F-004 — the previous
// check: only the current secret was derived, so a rotation locked every MFA
// user into an unsatisfiable second factor forever).
func NewService(repo *store.Repository, serverSecret, previousSecret []byte) *Service {
	key := make([]byte, 32)
	hk := hkdf.New(sha256.New, serverSecret, nil, []byte("mfa/totp"))
	if _, err := io.ReadFull(hk, key); err != nil {
		panic("mfa: derive totp key: " + err.Error())
	}
	s := &Service{repo: repo, key: key}
	if len(previousSecret) > 0 {
		prev := make([]byte, 32)
		pk := hkdf.New(sha256.New, previousSecret, nil, []byte("mfa/totp"))
		if _, err := io.ReadFull(pk, prev); err != nil {
			panic("mfa: derive previous totp key: " + err.Error())
		}
		s.prevKey = prev
	}
	return s
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
		plain, fromPrevious := s.decryptSecret(st.SecretCiphertext)
		step, valid := ValidateTotp(plain, code, now, totpWindow, st.LastUsedStep)
		if valid {
			// W9 F-005: the guarded advance IS the replay gate — a concurrent
			// verification of the same code loses the CAS (0 rows affected) and
			// is rejected below instead of being accepted twice.
			advanced, advErr := s.repo.AdvanceLastUsedStep(proof.UserID, step, now)
			if advErr != nil {
				return "", advErr
			}
			ok = advanced
			if ok {
				// W11 F-004: only the CAS winner re-wraps; a concurrent
				// loser must not rotate the ciphertext under a stale state.
				s.maybeRewrap(st, plain, fromPrevious, now)
			}
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
	// W11 F-004: a pending enrollment sealed under the previous secret is
	// re-wrapped on first successful confirmation.
	plain, fromPrevious := s.decryptSecret(st.SecretCiphertext)
	// W13 F-004 (GOAL-013 A-001): persist the MATCHED time step, not the
	// wall-clock step. ValidateTotp accepts the ±1 neighbor steps, so the
	// previous `now.Unix()/period` watermark could exceed the actually
	// consumed code's step (e.g. a confirm with the next-step code wrote the
	// current step); the first login then presented its current-step code,
	// which lost the `candidate <= lastUsedStep` replay check and was
	// rejected for up to two periods (~30–60s). Persisting the matched step
	// mirrors the login path (AdvanceLastUsedStep with the validated step).
	matchedStep, ok := ValidateTotp(plain, code, now, totpWindow, 0)
	if !ok {
		return ErrMFAInvalid
	}
	s.maybeRewrap(st, plain, fromPrevious, now)
	if err := s.repo.SetLastUsedStep(userID, matchedStep, now); err != nil {
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
	plain, fromPrevious := s.decryptSecret(st.SecretCiphertext)
	if _, ok := ValidateTotp(plain, code, now, totpWindow, st.LastUsedStep); !ok {
		return ErrMFAInvalid
	}
	// W11 F-004: a rotation-window decrypt re-wraps on any successful
	// second-factor proof (disable / recovery rotation included).
	s.maybeRewrap(st, plain, fromPrevious, now)
	return nil
}

// VerifySecondFactor validates a TOTP code or consumes a one-time recovery
// code against an active enrollment WITHOUT a login proof (workspace-019 R2 ·
// GOAL-003 D-001 §3): the self-recovery completion gate. Unlike the login
// path there is no mfa_proofs row — the caller's recovery challenge carries
// the guess budget (≤5 failed attempts void it), so this check stays a thin
// export of the self-service step-up semantics. A typed-nil receiver fails
// closed (the composition root passes a true nil interface when admin.mfa is
// off, but the guard mirrors Service.Required against future misuse).
func (s *Service) VerifySecondFactor(userID, code, recoveryCode string, now time.Time) error {
	if s == nil || s.repo == nil {
		return ErrNotEnrolled
	}
	return s.requireActiveSecondFactor(userID, code, recoveryCode, now)
}

// consumeRecoveryCode bcrypt-compares the code against the stored hash set and
// removes the matched hash (one-time consumption). W9 F-006: the rewrite is a
// guarded compare-and-set on updated_at — a concurrent redemption of another
// code moves the set first, so this caller's guarded write fails, it re-reads
// and retries; a concurrently re-presented SAME code loses the race and is
// rejected. The previous read-list/rewrite-whole-list pair resurrected consumed
// codes and allowed double use under concurrency.
func (s *Service) consumeRecoveryCode(st *store.State, code string, now time.Time) bool {
	trimmed := strings.TrimSpace(code)
	for attempt := 0; attempt < 4; attempt++ {
		var hashes []string
		if err := json.Unmarshal([]byte(st.RecoveryCodesHash), &hashes); err != nil {
			return false
		}
		matched := -1
		for i, h := range hashes {
			if bcrypt.CompareHashAndPassword([]byte(h), []byte(trimmed)) == nil {
				matched = i
				break
			}
		}
		if matched < 0 {
			return false
		}
		remaining := append(append([]string{}, hashes[:matched]...), hashes[matched+1:]...)
		next, err := json.Marshal(remaining)
		if err != nil {
			return false
		}
		consumed, err := s.repo.UpdateRecoveryCodesIfUnchanged(st.UserID, string(next), st.RecoveryCodesHash, now)
		if err != nil || consumed {
			return err == nil
		}
		// Lost the optimistic race: re-read the current set and retry so the
		// one-time semantics hold under concurrent redemption.
		fresh, err := s.repo.GetState(st.UserID)
		if err != nil {
			return false
		}
		st = fresh
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
// failed second factor, never a panic). W11 F-004: with a previous key
// configured the previous-secret key is tried as the rotation window
// fallback; the second result tells the caller the ciphertext was sealed
// under the PREVIOUS secret so it can re-wrap under the current key.
func (s *Service) decryptSecret(encoded string) (plain string, fromPrevious bool) {
	if p := decryptWithKey(s.key, encoded); p != "" {
		return p, false
	}
	if len(s.prevKey) > 0 {
		if p := decryptWithKey(s.prevKey, encoded); p != "" {
			return p, true
		}
	}
	return "", false
}

// maybeRewrap re-encrypts a TOTP secret that was just decrypted with the
// previous-secret key under the CURRENT key (W11 F-004). Best-effort: a
// failure keeps the previous-key window open for the next successful login —
// it must never fail an already-verified second factor.
func (s *Service) maybeRewrap(state *store.State, plain string, fromPrevious bool, now time.Time) {
	if !fromPrevious {
		return
	}
	if err := s.repo.UpdateSecretCiphertext(state.UserID, encryptSecret(s.key, []byte(plain)), now); err != nil {
		slog.Warn("mfa: rewrap totp secret under current key failed", "userID", state.UserID, "err", err)
	}
}

// decryptWithKey is the raw AES-256-GCM unwrap under exactly one key.
func decryptWithKey(key []byte, encoded string) string {
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
