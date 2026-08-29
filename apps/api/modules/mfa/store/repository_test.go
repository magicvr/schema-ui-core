package store

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func newRepo(t *testing.T) *Repository {
	t.Helper()
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "test.db"), "admin", "hash", true)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return NewRepository(st)
}

// W11 F-003: CreateProof lazily purges the user's EXPIRED proof rows in the
// same transaction (captcha precedent) — proof issuances cannot grow
// mfa_proofs unboundedly across challenges.
func TestCreateProofPurgesExpired(t *testing.T) {
	repo := newRepo(t)
	now := time.Now().UTC()

	oldExpired, err := repo.CreateProof("user-admin", now.Add(-time.Minute), now.Add(-5*time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	live, err := repo.CreateProof("user-admin", now.Add(5*time.Minute), now)
	if err != nil {
		t.Fatal(err)
	}
	// A new issuance for the same user drops the expired row only.
	if _, err := repo.CreateProof("user-admin", now.Add(5*time.Minute), now); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.GetProof(oldExpired.ID); err == nil {
		t.Fatalf("expired proof %s survived the purge", oldExpired.ID)
	}
	if _, err := repo.GetProof(live.ID); err != nil {
		t.Fatalf("live proof %s was purged: %v", live.ID, err)
	}
}

// W11 F-003: the failure counter is guarded in SQL (fail_count < 5) so
// concurrent wrong guesses cannot exceed the exhaustion budget via
// check-then-act.
func TestIncrementProofFailuresCapped(t *testing.T) {
	repo := newRepo(t)
	now := time.Now().UTC()
	proof, err := repo.CreateProof("user-admin", now.Add(5*time.Minute), now)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 8; i++ {
		if err := repo.IncrementProofFailures(proof.ID, now); err != nil {
			t.Fatalf("increment %d: %v", i, err)
		}
	}
	got, err := repo.GetProof(proof.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.FailCount != proofFailLimit {
		t.Fatalf("fail_count = %d, want cap %d", got.FailCount, proofFailLimit)
	}
}

// W9 A-005 R-F-003: regression locks for the F-005/F-006 guarded primitives —
// the TOTP watermark only moves forward (CAS), and the recovery-code rewrite
// only lands when the caller saw the exact previous set value.
func TestAdvanceLastUsedStepGuardedCAS(t *testing.T) {
	repo := newRepo(t)
	now := time.Now().UTC()
	if err := repo.UpsertPending("user-admin", "cipher", "[]", now); err != nil {
		t.Fatal(err)
	}

	advanced, err := repo.AdvanceLastUsedStep("user-admin", 100, now)
	if err != nil || !advanced {
		t.Fatalf("first advance = %v %v, want true/nil", advanced, err)
	}
	// Same step again: replay gate must reject.
	advanced, err = repo.AdvanceLastUsedStep("user-admin", 100, now)
	if err != nil || advanced {
		t.Fatalf("same-step advance = %v %v, want false/nil", advanced, err)
	}
	// Older step: rejected.
	advanced, err = repo.AdvanceLastUsedStep("user-admin", 99, now)
	if err != nil || advanced {
		t.Fatalf("older-step advance = %v %v, want false/nil", advanced, err)
	}
	// Newer step: accepted.
	advanced, err = repo.AdvanceLastUsedStep("user-admin", 101, now)
	if err != nil || !advanced {
		t.Fatalf("newer-step advance = %v %v, want true/nil", advanced, err)
	}
	st, err := repo.GetState("user-admin")
	if err != nil || st.LastUsedStep != 101 {
		t.Fatalf("last_used_step = %d %v, want 101", st.LastUsedStep, err)
	}
}

func TestUpdateRecoveryCodesIfUnchangedValueCAS(t *testing.T) {
	repo := newRepo(t)
	now := time.Now().UTC()
	const initial = `["h1","h2"]`
	if err := repo.UpsertPending("user-admin", "cipher", initial, now); err != nil {
		t.Fatal(err)
	}

	// Correct previous value: the swap lands.
	consumed, err := repo.UpdateRecoveryCodesIfUnchanged("user-admin", `["h2"]`, initial, now)
	if err != nil || !consumed {
		t.Fatalf("first CAS = %v %v, want true/nil", consumed, err)
	}
	// Stale previous value (the same-second window the W9 A-005 R-F-002 lock
	// closes): rejected even when updated_at is identical.
	consumed, err = repo.UpdateRecoveryCodesIfUnchanged("user-admin", `["h1"]`, initial, now)
	if err != nil || consumed {
		t.Fatalf("stale-token CAS = %v %v, want false/nil", consumed, err)
	}
	// Current value: lands again.
	consumed, err = repo.UpdateRecoveryCodesIfUnchanged("user-admin", `[]`, `["h2"]`, now)
	if err != nil || !consumed {
		t.Fatalf("second CAS = %v %v, want true/nil", consumed, err)
	}
	st, err := repo.GetState("user-admin")
	if err != nil || st.RecoveryCodesHash != `[]` {
		t.Fatalf("recovery_codes_hash = %q %v, want []", st.RecoveryCodesHash, err)
	}
}