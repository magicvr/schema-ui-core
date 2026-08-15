package mfa

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/mfa/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func newService(t *testing.T) *Service {
	t.Helper()
	hash, err := auth.HashPassword("test-password", 4)
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "test.db"), "admin", hash, true)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return NewService(store.NewRepository(st), []byte("test-secret"))
}

func totpForSecret(t *testing.T, s *Service, secret string, now time.Time) string {
	step := now.Unix() / totpPeriodSeconds
	code, err := totpCode(secret, step)
	if err != nil {
		t.Fatal(err)
	}
	return code
}

// Full lifecycle: no enrollment → Required false; enroll (pending) → still
// false (no self-lock, A-004 F-001); confirm → active; login gate requires;
// verify with proof succeeds once; replay of the same proof fails; recovery
// code consumption is one-time.
func TestServiceLifecycle(t *testing.T) {
	s := newService(t)
	now := time.Now().UTC()

	if s.Required("user-admin") {
		t.Fatalf("Required before enrollment must be false")
	}
	if _, err := s.BeginChallenge("user-admin", now); err != ErrNotEnrolled {
		t.Fatalf("BeginChallenge before enrollment err = %v, want ErrNotEnrolled", err)
	}

	secret, _, codes, err := s.Enroll("user-admin", "User One", now)
	if err != nil || len(codes) != 10 {
		t.Fatalf("enroll: %v, codes=%d", err, len(codes))
	}
	if s.Required("user-admin") {
		t.Fatalf("Required while pending must be false (A-004 F-001: no self-lock)")
	}
	// Wrong code on confirm is rejected.
	if err := s.Confirm("user-admin", "000000", now); err != ErrMFAInvalid {
		t.Fatalf("confirm wrong code err = %v", err)
	}
	if err := s.Confirm("user-admin", totpForSecret(t, s, secret, now), now); err != nil {
		t.Fatalf("confirm: %v", err)
	}
	if !s.Required("user-admin") {
		t.Fatalf("Required after confirm must be true")
	}

	// Login gate: BeginChallenge issues a one-time proof.
	at := now.Add(90 * time.Second) // advance past the confirm step
	proof, err := s.BeginChallenge("user-admin", at)
	if err != nil {
		t.Fatalf("begin challenge: %v", err)
	}
	if _, err := s.Verify(proof, totpForSecret(t, s, secret, at), "", at); err != nil {
		t.Fatalf("verify: %v", err)
	}
	// One-shot: the same proof is now unknown → expired-equivalent.
	if _, err := s.Verify(proof, totpForSecret(t, s, secret, at), "", at); err != ErrProofExpired {
		t.Fatalf("replayed proof err = %v, want ErrProofExpired", err)
	}

	// Same-window TOTP replay: a new proof but the same time step is rejected.
	proof2, err := s.BeginChallenge("user-admin", at)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Verify(proof2, totpForSecret(t, s, secret, at), "", at); err != ErrMFAInvalid {
		t.Fatalf("same-step replay err = %v, want ErrMFAInvalid", err)
	}

	// Failure accounting: 5 failures exhaust the proof.
	proof3, err := s.BeginChallenge("user-admin", at)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 4; i++ {
		if _, err := s.Verify(proof3, "000000", "", at); err != ErrMFAInvalid {
			t.Fatalf("attempt %d err = %v, want ErrMFAInvalid", i, err)
		}
	}
	// The 5th consecutive failure exhausts and destroys the proof.
	if _, err := s.Verify(proof3, "000000", "", at); err != ErrProofExhausted {
		t.Fatalf("exhausted proof err = %v, want ErrProofExhausted", err)
	}
	// A destroyed proof is indistinguishable from a consumed one.
	if _, err := s.Verify(proof3, "000000", "", at); err != ErrProofExpired {
		t.Fatalf("destroyed proof err = %v, want ErrProofExpired", err)
	}

	// Recovery code: one-time consumption via verify.
	proof4, err := s.BeginChallenge("user-admin", at)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Verify(proof4, "", codes[0], at); err != nil {
		t.Fatalf("recovery verify: %v", err)
	}
	proof5, err := s.BeginChallenge("user-admin", at)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Verify(proof5, "", codes[0], at); err != ErrMFAInvalid {
		t.Fatalf("reused recovery code err = %v, want ErrMFAInvalid", err)
	}

	// Rotate: requires a valid code; returns a fresh set.
	at2 := at.Add(60 * time.Second)
	newCodes, err := s.RotateRecovery("user-admin", totpForSecret(t, s, secret, at2), "", at2)
	if err != nil || len(newCodes) != 10 {
		t.Fatalf("rotate: %v", err)
	}
	if strings.Join(newCodes, ",") == strings.Join(codes, ",") {
		t.Fatalf("rotated codes must differ")
	}

	// Disable: requires second factor; removes the enrollment.
	at3 := at2.Add(60 * time.Second)
	if err := s.Disable("user-admin", totpForSecret(t, s, secret, at3), "", at3); err != nil {
		t.Fatalf("disable: %v", err)
	}
	if s.Required("user-admin") {
		t.Fatalf("Required after disable must be false")
	}
}

func TestServiceExpiredProof(t *testing.T) {
	s := newService(t)
	now := time.Now().UTC()
	secret, _, _, err := s.Enroll("user-admin", "U", now)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Confirm("user-admin", totpForSecret(t, s, secret, now), now); err != nil {
		t.Fatal(err)
	}
	proof, err := s.BeginChallenge("user-admin", now)
	if err != nil {
		t.Fatal(err)
	}
	// Expired proof: rejected and deleted.
	if _, err := s.Verify(proof, totpForSecret(t, s, secret, now.Add(10*time.Minute)), "", now.Add(10*time.Minute)); err != ErrProofExpired {
		t.Fatalf("expired proof err = %v, want ErrProofExpired", err)
	}
}

func TestServiceAdminReset(t *testing.T) {
	s := newService(t)
	now := time.Now().UTC()
	secret, _, _, err := s.Enroll("user-admin", "U", now)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Confirm("user-admin", totpForSecret(t, s, secret, now), now); err != nil {
		t.Fatal(err)
	}
	if err := s.AdminReset("user-admin"); err != nil {
		t.Fatalf("admin reset: %v", err)
	}
	if s.Required("user-admin") {
		t.Fatalf("Required after admin reset must be false")
	}
}