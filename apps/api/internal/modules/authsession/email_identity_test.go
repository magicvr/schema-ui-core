// workspace-018 R3 · GOAL-004: binding/verification flow tests over the R2
// schema + 0055 challenge table (sqlite; the PG dialect is covered by the
// full-catalog bootstrap integration in internal/store).
package authsession

import (
	"context"
	"errors"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

type recordingSender struct {
	messages []kernel.MailMessage
}

func (s *recordingSender) Send(_ context.Context, msg kernel.MailMessage) error {
	s.messages = append(s.messages, msg)
	return nil
}

var sixDigits = regexp.MustCompile(`\b\d{6}\b`)

func (s *recordingSender) lastCode() string {
	if len(s.messages) == 0 {
		return ""
	}
	return sixDigits.FindString(s.messages[len(s.messages)-1].TextBody)
}

func openEmailIdentityFixture(t *testing.T) (*Repository, *recordingSender) {
	t.Helper()
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "email-identity.db"), "admin", "hash-v1", true)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	repository := NewRepository(st)
	now := time.Now().UTC().Truncate(time.Second)
	if _, err := repository.CreateUserManagement(User{
		ID: "u-bob", Username: "bob", Name: "Bob", Roles: []string{"viewer"},
		PasswordHash: "hash", CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("create second user: %v", err)
	}
	return repository, &recordingSender{}
}

func emailIdentityRow(t *testing.T, repo *Repository, userID string) (email, status *string) {
	t.Helper()
	email, status, err := repo.EmailIdentityState(userID)
	if err != nil {
		t.Fatalf("email identity state: %v", err)
	}
	return email, status
}

func TestBindEmailReservesSlotAndSends(t *testing.T) {
	repo, sender := openEmailIdentityFixture(t)
	base := time.Now().UTC().Truncate(time.Second)

	if err := repo.BindEmail("user-admin", "Alice@Example.com ", sender, base); err != nil {
		t.Fatalf("bind: %v", err)
	}
	email, status := emailIdentityRow(t, repo, "user-admin")
	if email == nil || *email != "Alice@Example.com" {
		t.Fatalf("email = %v, want Alice@Example.com (trimmed original case)", email)
	}
	if status == nil || *status != "pending" {
		t.Fatalf("status = %v, want pending", status)
	}
	if len(sender.messages) != 1 || sender.messages[0].To != "Alice@Example.com" {
		t.Fatalf("sent = %+v, want one message to the trimmed address", sender.messages)
	}

	// Another account binding the same address case-insensitively is rejected.
	if err := repo.BindEmail("u-bob", "alice@example.com", sender, base); !errors.Is(err, ErrEmailTaken) {
		t.Fatalf("cross-account bind err = %v, want ErrEmailTaken", err)
	}
}

func TestVerifyEmailHappyPathAndIdempotentVerified(t *testing.T) {
	repo, sender := openEmailIdentityFixture(t)
	base := time.Now().UTC().Truncate(time.Second)

	if err := repo.BindEmail("user-admin", "alice@example.com", sender, base); err != nil {
		t.Fatalf("bind: %v", err)
	}
	code := sender.lastCode()
	if len(code) != 6 {
		t.Fatalf("code = %q, want 6 digits", code)
	}
	// Wrong code first: invalid, still pending.
	if err := repo.VerifyEmail("user-admin", "000000", base.Add(time.Minute)); !errors.Is(err, ErrEmailCodeInvalid) {
		t.Fatalf("wrong code err = %v, want ErrEmailCodeInvalid", err)
	}
	if err := repo.VerifyEmail("user-admin", code, base.Add(2*time.Minute)); err != nil {
		t.Fatalf("verify: %v", err)
	}
	if _, status := emailIdentityRow(t, repo, "user-admin"); status == nil || *status != "verified" {
		t.Fatalf("status after verify, want verified")
	}

	// Re-binding the SAME verified address is an idempotent no-op (no mail).
	before := len(sender.messages)
	if err := repo.BindEmail("user-admin", "ALICE@example.com", sender, base.Add(3*time.Minute)); err != nil {
		t.Fatalf("idempotent rebind: %v", err)
	}
	if len(sender.messages) != before {
		t.Fatal("idempotent rebind must not send a new message")
	}
}

func TestVerifyEmailExpiredDropsChallenge(t *testing.T) {
	repo, sender := openEmailIdentityFixture(t)
	base := time.Now().UTC().Truncate(time.Second)

	if err := repo.BindEmail("user-admin", "alice@example.com", sender, base); err != nil {
		t.Fatalf("bind: %v", err)
	}
	err := repo.VerifyEmail("user-admin", "123456", base.Add(emailCodeTTL+time.Minute))
	if !errors.Is(err, ErrEmailCodeExpired) {
		t.Fatalf("expired verify err = %v, want ErrEmailCodeExpired", err)
	}
	// Challenge dropped: further verifies report not-pending until a resend.
	if err := repo.VerifyEmail("user-admin", "123456", base.Add(emailCodeTTL+2*time.Minute)); !errors.Is(err, ErrEmailNotPending) {
		t.Fatalf("post-expiry verify err = %v, want ErrEmailNotPending", err)
	}
}

func TestFailedAttemptsVoidChallenge(t *testing.T) {
	repo, sender := openEmailIdentityFixture(t)
	base := time.Now().UTC().Truncate(time.Second)

	if err := repo.BindEmail("user-admin", "alice@example.com", sender, base); err != nil {
		t.Fatalf("bind: %v", err)
	}
	for i := 0; i < emailMaxFailedAttempts; i++ {
		err := repo.VerifyEmail("user-admin", "000000", base.Add(time.Duration(i)*time.Second))
		if i < emailMaxFailedAttempts-1 && !errors.Is(err, ErrEmailCodeInvalid) {
			t.Fatalf("attempt %d err = %v, want ErrEmailCodeInvalid", i+1, err)
		}
		if i == emailMaxFailedAttempts-1 && !errors.Is(err, ErrEmailCodeExpired) {
			t.Fatalf("final attempt err = %v, want ErrEmailCodeExpired (challenge voided)", err)
		}
	}
}

func TestBindSamePendingAddressHonorsCooldown(t *testing.T) {
	repo, sender := openEmailIdentityFixture(t)
	base := time.Now().UTC().Truncate(time.Second)

	if err := repo.BindEmail("user-admin", "alice@example.com", sender, base); err != nil {
		t.Fatalf("bind: %v", err)
	}
	// Same pending address = resend semantics: cooldown applies (A-001 F-002).
	if err := repo.BindEmail("user-admin", "ALICE@example.com", sender, base.Add(30*time.Second)); !errors.Is(err, ErrEmailResendCooldown) {
		t.Fatalf("same-address rebind err = %v, want ErrEmailResendCooldown", err)
	}
	// W13 F-009 (GOAL-013 A-001): a DIFFERENT address honors the same
	// cooldown — unlimited immediate dispatch to arbitrary addresses was a
	// mail-bomb primitive. Within the window it must be rejected…
	if err := repo.BindEmail("user-admin", "other@example.com", sender, base.Add(30*time.Second)); !errors.Is(err, ErrEmailResendCooldown) {
		t.Fatalf("rebind to new address within cooldown err = %v, want ErrEmailResendCooldown", err)
	}
	// …and after the cooldown window the rebind dispatches normally.
	if err := repo.BindEmail("user-admin", "other@example.com", sender, base.Add(61*time.Second)); err != nil {
		t.Fatalf("rebind to new address after cooldown: %v", err)
	}
}

func TestResendCooldown(t *testing.T) {
	repo, sender := openEmailIdentityFixture(t)
	base := time.Now().UTC().Truncate(time.Second)

	if err := repo.BindEmail("user-admin", "alice@example.com", sender, base); err != nil {
		t.Fatalf("bind: %v", err)
	}
	if err := repo.ResendEmailCode("user-admin", sender, base.Add(30*time.Second)); !errors.Is(err, ErrEmailResendCooldown) {
		t.Fatalf("early resend err = %v, want ErrEmailResendCooldown", err)
	}
	if err := repo.ResendEmailCode("user-admin", sender, base.Add(61*time.Second)); err != nil {
		t.Fatalf("resend after cooldown: %v", err)
	}
	code := sender.lastCode()
	if len(code) != 6 {
		t.Fatalf("resent code = %q, want 6 digits", code)
	}
	if err := repo.VerifyEmail("user-admin", code, base.Add(62*time.Second)); err != nil {
		t.Fatalf("verify resent code: %v", err)
	}
}

func TestRebindOverwriteReleasesOldSlot(t *testing.T) {
	repo, sender := openEmailIdentityFixture(t)
	base := time.Now().UTC().Truncate(time.Second)

	if err := repo.BindEmail("user-admin", "old@example.com", sender, base); err != nil {
		t.Fatalf("bind old: %v", err)
	}
	// Rebind overwrites to the new address (contract §5).
	if err := repo.BindEmail("user-admin", "new@example.com", sender, base.Add(time.Minute)); err != nil {
		t.Fatalf("rebind: %v", err)
	}
	// The old slot is released by the overwrite: another account can take it.
	if err := repo.BindEmail("u-bob", "OLD@example.com", sender, base.Add(2*time.Minute)); err != nil {
		t.Fatalf("bind released slot: %v", err)
	}
}

func TestAdminPrefillPendingCannotVerifyDirectly(t *testing.T) {
	repo, _ := openEmailIdentityFixture(t)
	now := time.Now().UTC()

	if _, err := repo.UpdateUser("u-bob", UserPatch{Email: strPtr("Bob@Example.com")}, "user-admin", now); err != nil {
		t.Fatalf("admin prefill: %v", err)
	}
	email, status := emailIdentityRow(t, repo, "u-bob")
	if email == nil || *email != "Bob@Example.com" || status == nil || *status != "pending" {
		t.Fatalf("prefill state = (%v, %v), want (Bob@Example.com, pending)", email, status)
	}
	// No challenge exists yet: verification is not possible without delivery.
	if err := repo.VerifyEmail("u-bob", "123456", now); !errors.Is(err, ErrEmailNotPending) {
		t.Fatalf("verify without challenge err = %v, want ErrEmailNotPending", err)
	}
	// Clearing resets to unbound.
	if _, err := repo.UpdateUser("u-bob", UserPatch{Email: strPtr("")}, "user-admin", now.Add(time.Minute)); err != nil {
		t.Fatalf("clear email: %v", err)
	}
	email, status = emailIdentityRow(t, repo, "u-bob")
	if email != nil || status != nil {
		t.Fatalf("cleared state = (%v, %v), want (nil, nil)", email, status)
	}
}

func TestVerifyUnboundAccountIsControlledNotPending(t *testing.T) {
	// A-003 F-001: an account that never bound an email must get the
	// controlled ErrEmailNotPending (→ HTTP 409), not an internal error.
	repo, _ := openEmailIdentityFixture(t)
	err := repo.VerifyEmail("u-bob", "123456", time.Now().UTC())
	if !errors.Is(err, ErrEmailNotPending) {
		t.Fatalf("unbound verify err = %v, want ErrEmailNotPending", err)
	}
}

func TestSendFailureRollsBindBack(t *testing.T) {
	repo := func() *Repository {
		st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "send-fail.db"), "admin", "hash-v1", true)
		if err != nil {
			t.Fatalf("open store: %v", err)
		}
		t.Cleanup(func() { _ = st.Close() })
		return NewRepository(st)
	}()
	failing := &failingSender{}
	err := repo.BindEmail("user-admin", "alice@example.com", failing, time.Now().UTC())
	if !errors.Is(err, ErrEmailSendFailed) {
		t.Fatalf("bind with failing sender err = %v, want ErrEmailSendFailed", err)
	}
	email, status := emailIdentityRow(t, repo, "user-admin")
	if email != nil || status != nil {
		t.Fatalf("state after failed send = (%v, %v), want rolled back (nil, nil)", email, status)
	}
}

type failingSender struct{}

func (failingSender) Send(context.Context, kernel.MailMessage) error {
	return errors.New("smtp dial failed")
}

func strPtr(s string) *string { return &s }
