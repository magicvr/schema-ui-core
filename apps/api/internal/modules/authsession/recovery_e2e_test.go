// workspace-019 R2 · GOAL-003 C4: end-to-end evidence for self-recovery over
// THE VP-017 default-channel adapter — the mock OutboxSink publishing to the
// admin-visible mail_outbox record table. The reset code is recovered from
// the CHANNEL RECORD (never from a test stub), proving the VP-019 exit path:
// "无生产渠道时须能从 017 默认渠道取出恢复信并完成设密".
package authsession

import (
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func TestR2EndToEndRecoveryThroughMockChannel(t *testing.T) {
	repo, sink := openRecoveryFixture(t)
	base := time.Now().UTC().Truncate(time.Second)
	bindVerifiedEmail(t, repo, sink, "user-admin", "alice@example.com", base)

	// start: the code rides the default channel.
	if err := repo.StartRecovery("admin", sink, base); err != nil {
		t.Fatalf("start: %v", err)
	}
	code := codeFromOutbox(t, sink, "alice@example.com")

	// complete: match + consume + rotate in one call. The replacement hash is
	// real bcrypt (A-001 F-002) so the stored credential is verifiable with
	// the same compare the login surface performs (auth.VerifyPassword).
	hash, herr := bcrypt.GenerateFromPassword([]byte("brand-new-pass-1"), 4)
	if herr != nil {
		t.Fatalf("hash new password: %v", herr)
	}
	outcome, err := repo.EvaluateRecoveryCode("user-admin", code, base.Add(time.Minute))
	if err != nil || outcome != RecoveryMatch {
		t.Fatalf("evaluate channel-record code = (%v, %v), want match", outcome, err)
	}
	if err := repo.CompleteRecovery("user-admin", string(hash), "user-admin", base.Add(2*time.Minute)); err != nil {
		t.Fatalf("complete: %v", err)
	}
	u, err := repo.UserByID("user-admin")
	if err != nil {
		t.Fatalf("reload user: %v", err)
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte("brand-new-pass-1")) != nil {
		t.Fatal("stored hash does not verify the new password (login would fail)")
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte("old-credential")) == nil {
		t.Fatal("stored hash unexpectedly verifies a wrong password")
	}

	// The consumed code cannot re-enter (one-shot challenge).
	if err := repo.CompleteRecovery("user-admin", string(hash), "user-admin", base.Add(3*time.Minute)); !errors.Is(err, ErrRecoveryNotPending) {
		t.Fatalf("replay err = %v, want ErrRecoveryNotPending", err)
	}
}
