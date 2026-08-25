package mfa

// workspace-019 R2 (GOAL-003 D-001 §3 · A-001 F-002): VerifySecondFactor is
// the self-recovery completion gate — real TOTP validation and one-time
// recovery-code consumption WITHOUT a login proof. This test walks the REAL
// service (no fakes) so the MFA branch of success criterion 2 is same-chain.
import (
	"testing"
	"time"
)

func TestVerifySecondFactorRecoveryGate(t *testing.T) {
	s := newService(t)
	base := time.Now().UTC().Truncate(time.Second)

	// Not enrolled → fail closed (the handler treats this as an invalid factor).
	if err := s.VerifySecondFactor("user-admin", "000000", "", base); err != ErrNotEnrolled {
		t.Fatalf("un-enrolled err = %v, want ErrNotEnrolled", err)
	}

	secret, _, codes, err := s.Enroll("user-admin", "User One", base)
	if err != nil {
		t.Fatalf("enroll: %v", err)
	}
	at := base.Add(90 * time.Second)
	if err := s.Confirm("user-admin", totpForSecret(t, s, secret, at), at); err != nil {
		t.Fatalf("confirm: %v", err)
	}
	if !s.Required("user-admin") {
		t.Fatal("active enrollment must gate recovery completion")
	}

	// Wrong TOTP → ErrMFAInvalid, enrollment untouched.
	if err := s.VerifySecondFactor("user-admin", "000000", "", at.Add(90*time.Second)); err != ErrMFAInvalid {
		t.Fatalf("wrong TOTP err = %v, want ErrMFAInvalid", err)
	}
	// Correct TOTP at a fresh step passes.
	good := at.Add(180 * time.Second)
	if err := s.VerifySecondFactor("user-admin", totpForSecret(t, s, secret, good), "", good); err != nil {
		t.Fatalf("correct TOTP err = %v, want nil", err)
	}
	// Recovery-code alternative: first consumption succeeds, replay fails.
	recAt := good.Add(90 * time.Second)
	if err := s.VerifySecondFactor("user-admin", "", codes[0], recAt); err != nil {
		t.Fatalf("recovery code err = %v, want nil", err)
	}
	if err := s.VerifySecondFactor("user-admin", "", codes[0], recAt); err != ErrMFAInvalid {
		t.Fatalf("replayed recovery code err = %v, want ErrMFAInvalid", err)
	}
}
