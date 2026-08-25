// Self-service password recovery domain tests (workspace-019 R2 · GOAL-003
// D-001): resolution, start (cooldown / anti-enumeration silence), evaluate
// (mismatch bookkeeping), completion semantics (challenge consumed first,
// UpdateUser lands token_version bump + refresh revocation + must-change
// clear). The e2e twin lives in recovery_e2e_test.go over the REAL channel.
package authsession

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/mail"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)


func openRecoveryFixture(t *testing.T) (*Repository, *mail.OutboxSink) {
	t.Helper()
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "recovery.db"), "admin", "hash-v1", true)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	repo := NewRepository(st)
	now := time.Now().UTC().Truncate(time.Second)
	if _, err := repo.CreateUserManagement(User{
		ID: "u-bob", Username: "bob", Name: "Bob", Roles: []string{"viewer"},
		PasswordHash: "old-hash", CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("create bob: %v", err)
	}
	return repo, mail.NewOutboxSink(st, 0)
}

func bindVerifiedEmail(t *testing.T, repo *Repository, sink *mail.OutboxSink, userID, address string, now time.Time) {
	t.Helper()
	if err := repo.BindEmail(userID, address, sink, now); err != nil {
		t.Fatalf("bind email %s: %v", userID, err)
	}
	code := codeFromOutbox(t, sink, address)
	if err := repo.VerifyEmail(userID, code, now.Add(time.Minute)); err != nil {
		t.Fatalf("verify email %s: %v", userID, err)
	}
}

func TestResolveRecoveryTarget(t *testing.T) {
	repo, sink := openRecoveryFixture(t)
	base := time.Now().UTC().Truncate(time.Second)
	bindVerifiedEmail(t, repo, sink, "user-admin", "Admin@Example.com", base)

	target, err := repo.ResolveRecoveryTarget("admin")
	if err != nil || target.UserID != "user-admin" || !target.Enabled {
		t.Fatalf("resolve by username = (%+v, %v)", target, err)
	}
	target, err = repo.ResolveRecoveryTarget("ADMIN@example.com")
	if err != nil || target.UserID != "user-admin" {
		t.Fatalf("resolve by email = (%+v, %v)", target, err)
	}
	if _, err := repo.ResolveRecoveryTarget("nobody"); !errors.Is(err, ErrRecoveryNotAvailable) {
		t.Fatalf("unknown err = %v, want ErrRecoveryNotAvailable", err)
	}
	// bound-but-unverified email must NOT be a recovery path (I-006)
	if err := repo.BindEmail("u-bob", "bob@example.com", sink, base); err != nil {
		t.Fatalf("bind bob pending email: %v", err)
	}
	if _, err := repo.ResolveRecoveryTarget("bob@example.com"); !errors.Is(err, ErrRecoveryNotAvailable) {
		t.Fatalf("pending-email resolve err = %v, want ErrRecoveryNotAvailable", err)
	}
}

func TestStartRecoveryDispatchesAndHonoursCooldown(t *testing.T) {
	repo, sink := openRecoveryFixture(t)
	base := time.Now().UTC().Truncate(time.Second)
	bindVerifiedEmail(t, repo, sink, "user-admin", "alice@example.com", base)

	sender := &recordingSender{}
	if err := repo.StartRecovery("admin", sender, base); err != nil {
		t.Fatalf("start: %v", err)
	}
	if len(sender.messages) != 1 || sender.messages[0].To != "alice@example.com" {
		t.Fatalf("dispatched = %+v, want one mail to alice@example.com", sender.messages)
	}
	code := outboxCodePattern.FindString(sender.messages[0].TextBody)
	if len(code) != 6 {
		t.Fatalf("body carries no 6-digit code: %q", sender.messages[0].TextBody)
	}

	// Frozen I-002 cooldown: an immediate second start is rejected.
	if err := repo.StartRecovery("admin", sender, base.Add(30*time.Second)); !errors.Is(err, ErrRecoveryCooldown) {
		t.Fatalf("cooldown err = %v, want ErrRecoveryCooldown", err)
	}
	// After the window a fresh challenge dispatches again.
	if err := repo.StartRecovery("admin", sender, base.Add(61*time.Second)); err != nil {
		t.Fatalf("second start: %v", err)
	}
	if len(sender.messages) != 2 {
		t.Fatalf("second dispatch missing: %+v", sender.messages)
	}
	fresh := outboxCodePattern.FindString(sender.messages[1].TextBody)
	if fresh == "" || fresh == code {
		t.Fatalf("fresh challenge did not rotate the code (%q vs old %q)", fresh, code)
	}
	// The rotated challenge invalidates the FIRST code.
	outcome, err := repo.EvaluateRecoveryCode("user-admin", code, base.Add(62*time.Second))
	if err != nil || outcome != RecoveryMismatch {
		t.Fatalf("stale code evaluate = (%v, %v), want mismatch", outcome, err)
	}
}

func TestStartRecoverySilentOnNoPath(t *testing.T) {
	repo, _ := openRecoveryFixture(t)
	base := time.Now().UTC().Truncate(time.Second)
	sender := &recordingSender{}
	if err := repo.StartRecovery("u-bob", sender, base); !errors.Is(err, ErrRecoveryNotAvailable) {
		t.Fatalf("unbound err = %v, want ErrRecoveryNotAvailable", err)
	}
	if err := repo.StartRecovery("ghost", sender, base); !errors.Is(err, ErrRecoveryNotAvailable) {
		t.Fatalf("unknown err = %v, want ErrRecoveryNotAvailable", err)
	}
	if len(sender.messages) != 0 {
		t.Fatalf("no-path starts dispatched %+v, want silence", sender.messages)
	}
}

func TestCompleteRecoveryConsumesChallengeAndRotatesCredentials(t *testing.T) {
	repo, sink := openRecoveryFixture(t)
	base := time.Now().UTC().Truncate(time.Second)
	bindVerifiedEmail(t, repo, sink, "user-admin", "alice@example.com", base)
	if err := repo.StartRecovery("admin", sink, base); err != nil {
		t.Fatalf("start: %v", err)
	}
	code := codeFromOutbox(t, sink, "alice@example.com")

	wrong := "000000"
	if wrong == code {
		wrong = "000001"
	}
	outcome, err := repo.EvaluateRecoveryCode("user-admin", wrong, base.Add(time.Minute))
	if err != nil || outcome != RecoveryMismatch {
		t.Fatalf("wrong-code evaluate = (%v, %v), want mismatch", outcome, err)
	}
	if repo.ConsumeRecoveryAttempt("user-admin") {
		t.Fatal("first failure must not void the challenge")
	}

	outcome, err = repo.EvaluateRecoveryCode("user-admin", code, base.Add(time.Minute))
	if err != nil || outcome != RecoveryMatch {
		t.Fatalf("right-code evaluate = (%v, %v), want match", outcome, err)
	}
	if err := repo.CompleteRecovery("user-admin", "new-hash", "user-admin", base.Add(2*time.Minute)); err != nil {
		t.Fatalf("complete: %v", err)
	}
	// Challenge is one-shot: replay cannot re-enter.
	if err := repo.CompleteRecovery("user-admin", "another-hash", "user-admin", base.Add(3*time.Minute)); !errors.Is(err, ErrRecoveryNotPending) {
		t.Fatalf("replay complete err = %v, want ErrRecoveryNotPending", err)
	}
	u, err := repo.UserByID("user-admin")
	if err != nil || u.PasswordHash != "new-hash" {
		t.Fatalf("password after recovery = (%v, %v), want new-hash", u, err)
	}
	// §4 session semantics: token_version bumped + refresh tokens revoked +
	// forced-initial-password flag cleared (W16-F01 parity).
	var version, revoked int
	err = repo.withTx("recovery session-semantics readback", func(tx kernel.Tx) error {
		if err := tx.QueryRow(context.Background(), `SELECT token_version FROM users WHERE id = 'user-admin'`).Scan(&version); err != nil {
			return err
		}
		return tx.QueryRow(context.Background(), `SELECT COUNT(*) FROM refresh_tokens WHERE user_id = 'user-admin' AND revoked_at IS NULL`).Scan(&revoked)
	})
	if err != nil {
		t.Fatalf("session-semantics readback: %v", err)
	}
	if version == 0 {
		t.Fatal("token_version was not bumped by recovery completion")
	}
	if revoked != 0 {
		t.Fatalf("live refresh tokens after recovery = %d, want 0", revoked)
	}
	if u.MustChangePassword {
		t.Fatal("must_change_password survived recovery completion")
	}
}

func TestCompleteRecoveryVoidAfterFiveFailures(t *testing.T) {
	repo, sink := openRecoveryFixture(t)
	base := time.Now().UTC().Truncate(time.Second)
	bindVerifiedEmail(t, repo, sink, "user-admin", "alice@example.com", base)
	if err := repo.StartRecovery("admin", sink, base); err != nil {
		t.Fatalf("start: %v", err)
	}
	for i := 0; i < 4; i++ {
		if voided := repo.ConsumeRecoveryAttempt("user-admin"); voided {
			t.Fatalf("attempt %d prematurely voided the challenge", i+1)
		}
	}
	if voided := repo.ConsumeRecoveryAttempt("user-admin"); !voided {
		t.Fatal("5th failure must void the challenge")
	}
	outcome, err := repo.EvaluateRecoveryCode("user-admin", "123456", base.Add(time.Minute))
	if err != nil || outcome != RecoveryNotPending {
		t.Fatalf("voided evaluate = (%v, %v), want notPending", outcome, err)
	}
}

// TestCompleteRecoveryExpiredCodeClassifiesStale walks the expiry branch:
// past-TTL challenges classify expired and the stale row drops on demand.
func TestCompleteRecoveryExpiredCodeClassifiesStale(t *testing.T) {
	repo, sink := openRecoveryFixture(t)
	base := time.Now().UTC().Truncate(time.Second)
	bindVerifiedEmail(t, repo, sink, "user-admin", "alice@example.com", base)
	if err := repo.StartRecovery("admin", sink, base); err != nil {
		t.Fatalf("start: %v", err)
	}
	code := codeFromOutbox(t, sink, "alice@example.com")
	outcome, err := repo.EvaluateRecoveryCode("user-admin", code, base.Add(11*time.Minute))
	if err != nil || outcome != RecoveryExpired {
		t.Fatalf("expired evaluate = (%v, %v), want expired", outcome, err)
	}
	repo.DropStaleRecoveryChallenge("user-admin", base.Add(11*time.Minute))
	outcome, err = repo.EvaluateRecoveryCode("user-admin", code, base.Add(12*time.Minute))
	if err != nil || outcome != RecoveryNotPending {
		t.Fatalf("post-drop evaluate = (%v, %v), want notPending", outcome, err)
	}
}
