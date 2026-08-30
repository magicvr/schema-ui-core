package mfa

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/modules/mfa/store"
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
	return NewService(store.NewRepository(st), []byte("test-secret"), nil)
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

// W13 F-004 (GOAL-013 A-001): Confirm must persist the MATCHED TOTP step,
// not the wall-clock step. The previous behavior wrote now.Unix()/period
// even when validation matched the ±1 neighbor step — a confirm with the
// PREVIOUS-step code (slow submit / device clock behind) set a watermark at
// the current step, so the first login's current-step code lost the replay
// check and was rejected for 30–60s. After the fix the first login right
// after confirm succeeds immediately.
func TestServiceConfirmPersistsMatchedStep(t *testing.T) {
	s := newService(t)
	now := time.Now().UTC()
	secret, _, _, err := s.Enroll("user-admin", "U", now)
	if err != nil {
		t.Fatal(err)
	}
	prevStep := now.Unix()/totpPeriodSeconds - 1
	if prevStep < 0 {
		t.Skip("process started within the first TOTP period")
	}
	prevCode, err := totpCode(secret, prevStep)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Confirm("user-admin", prevCode, now); err != nil {
		t.Fatalf("confirm with previous-step code: %v", err)
	}
	st, err := s.repo.GetState("user-admin")
	if err != nil {
		t.Fatal(err)
	}
	if st.LastUsedStep != prevStep {
		t.Fatalf("last_used_step = %d, want matched step %d", st.LastUsedStep, prevStep)
	}
	// Behavioral half: the first login in the SAME period as confirm accepts
	// the current-step code (the old wall-clock watermark rejected it).
	proof, err := s.BeginChallenge("user-admin", now)
	if err != nil {
		t.Fatal(err)
	}
	got, err := s.Verify(proof, totpForSecret(t, s, secret, now), "", now)
	if err != nil || got != "user-admin" {
		t.Fatalf("first login after confirm = %q %v, want immediate success", got, err)
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

// W15 F-004 (GOAL-016 A-001): the self-service step-up path (disable /
// recovery rotation / self-recovery gate) must advance the TOTP replay
// watermark like the login path. The same code resubmitted inside the window
// is rejected on the second use instead of succeeding twice.
func TestServiceStepUpTotpReplayRejected(t *testing.T) {
	s := newService(t)
	now := time.Now().UTC()
	secret, _, _, err := s.Enroll("user-admin", "U", now)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Confirm("user-admin", totpForSecret(t, s, secret, now), now); err != nil {
		t.Fatal(err)
	}
	at := now.Add(90 * time.Second) // advance past the confirm step
	code := totpForSecret(t, s, secret, at)

	// First use succeeds and advances the watermark (non-destructive op so the
	// enrollment survives for the replay probe).
	if _, err := s.RotateRecovery("user-admin", code, "", at); err != nil {
		t.Fatalf("first rotate: %v", err)
	}
	// Same-window replay of the same code must fail closed.
	if _, err := s.RotateRecovery("user-admin", code, "", at); err != ErrMFAInvalid {
		t.Fatalf("same-window replay err = %v, want ErrMFAInvalid", err)
	}
	// The self-recovery gate shares the same step-up semantics.
	if err := s.VerifySecondFactor("user-admin", code, "", at); err != ErrMFAInvalid {
		t.Fatalf("VerifySecondFactor replay err = %v, want ErrMFAInvalid", err)
	}
	// A FRESH code from the next step still works (window not poisoned).
	next := at.Add(time.Duration(totpPeriodSeconds) * time.Second)
	if _, err := s.RotateRecovery("user-admin", totpForSecret(t, s, secret, next), "", next); err != nil {
		t.Fatalf("next-step rotate: %v", err)
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
	removed, err := s.AdminReset("user-admin")
	if err != nil {
		t.Fatalf("admin reset: %v", err)
	}
	if !removed {
		t.Fatalf("admin reset of active enrollment must report removed=true")
	}
	if s.Required("user-admin") {
		t.Fatalf("Required after admin reset must be false")
	}
	// Resetting a user with no enrollment is a no-op and must NOT report an
	// active removal (W7 F-002: no generic forced-logout from mfa-reset).
	removed, err = s.AdminReset("user-none")
	if err != nil {
		t.Fatalf("admin reset no-enrollment: %v", err)
	}
	if removed {
		t.Fatalf("no-enrollment reset must report removed=false")
	}
}
// A-007 F-002: an active enrollment cannot be overwritten by re-enrolling
// (that would tear down MFA without the second factor); pending ones can.
func TestServiceEnrollCannotOverwriteActive(t *testing.T) {
	s := newService(t)
	now := time.Now().UTC()
	secret, _, _, err := s.Enroll("user-admin", "Admin", now)
	if err != nil {
		t.Fatal(err)
	}
	// A pending enrollment may be re-enrolled (A-005 recommended).
	secret2, _, _, err := s.Enroll("user-admin", "Admin", now)
	if err != nil {
		t.Fatalf("re-enroll over pending: %v", err)
	}
	if secret2 == secret {
		t.Fatalf("pending re-enroll must rotate the secret")
	}
	if err := s.Confirm("user-admin", totpForSecret(t, s, secret2, now), now); err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := s.Enroll("user-admin", "Admin", now); err != ErrActive {
		t.Fatalf("re-enroll over active err = %v, want ErrActive", err)
	}
	// After disable, enrollment works again.
	at := now.Add(90 * time.Second)
	if err := s.Disable("user-admin", totpForSecret(t, s, secret2, at), "", at); err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := s.Enroll("user-admin", "Admin", at); err != nil {
		t.Fatalf("enroll after disable: %v", err)
	}
}

// W11 F-004: a JWT-secret rotation must not lock MFA users out. A TOTP
// secret sealed under the PREVIOUS secret stays decryptable through the
// rotation-window fallback, and a successful second-factor verification
// re-wraps the ciphertext under the CURRENT key (after which the previous
// key is no longer needed).
func TestServiceRotationWindow(t *testing.T) {
	hash, err := auth.HashPassword("test-password", 4)
	if err != nil {
		t.Fatal(err)
	}
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "test.db"), "admin", hash, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	repo := store.NewRepository(st)
	now := time.Now().UTC()

	secretA := []byte("jwt-secret-a")
	secretB := []byte("jwt-secret-b")

	// Enrollment happens under secret A.
	first := NewService(repo, secretA, nil)
	secret, _, _, err := first.Enroll("user-admin", "Rotation", now)
	if err != nil {
		t.Fatal(err)
	}
	if err := first.Confirm("user-admin", totpForSecret(t, first, secret, now), now); err != nil {
		t.Fatal(err)
	}

	// Rotation: current = B, previous = A. Login must still complete.
	mid := NewService(repo, secretB, secretA)
	at := now.Add(90 * time.Second)
	proof, err := mid.BeginChallenge("user-admin", at)
	if err != nil {
		t.Fatalf("begin challenge after rotation: %v", err)
	}
	if got, err := mid.Verify(proof, totpForSecret(t, mid, secret, at), "", at); err != nil || got != "user-admin" {
		t.Fatalf("verify via previous key = %q %v, want user-admin/nil", got, err)
	}

	// The successful verify re-wrapped under current key: a service with
	// ONLY the current key (previous dropped) must still verify.
	post := NewService(repo, secretB, nil)
	at2 := at.Add(90 * time.Second)
	proof2, err := post.BeginChallenge("user-admin", at2)
	if err != nil {
		t.Fatalf("begin challenge after rewrap: %v", err)
	}
	if got, err := post.Verify(proof2, totpForSecret(t, post, secret, at2), "", at2); err != nil || got != "user-admin" {
		t.Fatalf("verify after rewrap = %q %v, want user-admin/nil", got, err)
	}

	// Sanity: the OLD current key alone can no longer decrypt the re-wrapped
	// ciphertext — the secret now lives under B only.
	stale := NewService(repo, secretA, nil)
	at3 := at2.Add(90 * time.Second)
	proof3, err := stale.BeginChallenge("user-admin", at3)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := stale.Verify(proof3, totpForSecret(t, stale, secret, at3), "", at3); err != ErrMFAInvalid {
		t.Fatalf("stale-key verify err = %v, want ErrMFAInvalid (decrypt failure surface)", err)
	}
}