// workspace-018 R4 · GOAL-005 C1/C2: end-to-end evidence for the account
// email identity flow over THE VP-017 default-channel adapter — the mock
// OutboxSink publishing to the admin-visible mail_outbox record table
// (workspace-017 GOAL-007 contract). The verification code is recovered from
// the CHANNEL RECORD (not from any test stub), proving VP-018 exit criterion
// #2: "无生产渠道时须能从 017 默认渠道取出校验信".
package authsession

import (
	"context"
	"errors"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/mail"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

var outboxCodePattern = regexp.MustCompile(`\b\d{6}\b`)

func codeFromOutbox(t *testing.T, sink *mail.OutboxSink, to string) string {
	t.Helper()
	records, err := sink.List(context.Background(), 50, 0)
	if err != nil {
		t.Fatalf("list outbox records: %v", err)
	}
	for _, rec := range records {
		if rec.To != to {
			continue
		}
		full, err := sink.Get(context.Background(), rec.ID)
		if err != nil {
			t.Fatalf("get outbox record %s: %v", rec.ID, err)
		}
		code := outboxCodePattern.FindString(full.Body)
		if code == "" {
			t.Fatalf("outbox body for %s carries no 6-digit code: %q", to, full.Body)
		}
		return code
	}
	t.Fatalf("no outbox record addressed to %s", to)
	return ""
}

func openR4Fixture(t *testing.T) (*Repository, *mail.OutboxSink) {
	t.Helper()
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "r4-e2e.db"), "admin", "hash-v1", true)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	repo := NewRepository(st)
	now := time.Now().UTC().Truncate(time.Second)
	if _, err := repo.CreateUserManagement(User{
		ID: "u-bob", Username: "bob", Name: "Bob", Roles: []string{"viewer"},
		PasswordHash: "hash", CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("create second user: %v", err)
	}
	return repo, mail.NewOutboxSink(st, 0)
}

// TestR4EndToEndBindVerifyThroughMockChannel walks the frozen contract end to
// end through the real default-channel adapter.
func TestR4EndToEndBindVerifyThroughMockChannel(t *testing.T) {
	repo, sink := openR4Fixture(t)
	base := time.Now().UTC().Truncate(time.Second)

	// §2 bind reserves the slot; §4 the code rides the default channel.
	if err := repo.BindEmail("user-admin", "Alice@Example.com", sink, base); err != nil {
		t.Fatalf("bind: %v", err)
	}
	code := codeFromOutbox(t, sink, "Alice@Example.com")

	// Uniqueness fail-closed is observable at the same seam: another account
	// cannot take the address case-insensitively while it is pending.
	if err := repo.BindEmail("u-bob", "alice@example.com", sink, base); !errors.Is(err, ErrEmailTaken) {
		t.Fatalf("cross-account bind err = %v, want ErrEmailTaken", err)
	}

	// §6 complete the state machine with the code FROM the channel record.
	if err := repo.VerifyEmail("user-admin", code, base.Add(time.Minute)); err != nil {
		t.Fatalf("verify with channel-record code: %v", err)
	}
	email, status, err := repo.EmailIdentityState("user-admin")
	if err != nil || email == nil || *email != "Alice@Example.com" || status == nil || *status != "verified" {
		t.Fatalf("final state = (%v, %v) err %v, want verified identity", email, status, err)
	}

	// §1 no-email accounts are unaffected throughout.
	bobEmail, bobStatus, err := repo.EmailIdentityState("u-bob")
	if err != nil || bobEmail != nil || bobStatus != nil {
		t.Fatalf("bob state = (%v, %v) err %v, want unbound (nil, nil)", bobEmail, bobStatus, err)
	}
	if _, err := repo.UserByID("u-bob"); err != nil {
		t.Fatalf("no-email account lookup broken: %v", err)
	}
}
