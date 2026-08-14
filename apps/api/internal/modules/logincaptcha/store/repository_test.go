// Package store tests for the admin.login-captcha persistence (S-11 ·
// GOAL-011 D-002 §1/§3): one-time challenge rows (atomic consume with expiry)
// and the single-row switch.
package store

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func newTestEnv(t *testing.T) *Repository {
	t.Helper()
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "test.db"), "admin", "hash", true)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return NewRepository(st)
}

func TestConsumeChallengeLifecycle(t *testing.T) {
	r := newTestEnv(t)
	now := time.Now().UTC()
	if err := r.CreateChallenge("cap-1", "hash-1", now.Add(5*time.Minute), now); err != nil {
		t.Fatalf("create: %v", err)
	}
	ok, err := r.ConsumeChallenge("cap-1", "hash-1", now)
	if err != nil || !ok {
		t.Fatalf("consume = %v, %v; want true", ok, err)
	}
	// consumed: second attempt fails
	ok, err = r.ConsumeChallenge("cap-1", "hash-1", now)
	if err != nil || ok {
		t.Fatalf("consume after consume = %v, %v; want false", ok, err)
	}
	// wrong answer consumes too
	if err := r.CreateChallenge("cap-2", "hash-2", now.Add(5*time.Minute), now); err != nil {
		t.Fatalf("create 2: %v", err)
	}
	ok, err = r.ConsumeChallenge("cap-2", "wrong", now)
	if err != nil || ok {
		t.Fatalf("wrong answer consume = %v, %v; want false", ok, err)
	}
	// unknown id: no-op false
	ok, err = r.ConsumeChallenge("cap-unknown", "x", now)
	if err != nil || ok {
		t.Fatalf("unknown consume = %v, %v; want false", ok, err)
	}
}

func TestConsumeChallengeExpiredFails(t *testing.T) {
	r := newTestEnv(t)
	now := time.Now().UTC()
	if err := r.CreateChallenge("cap-old", "hash-old", now.Add(-time.Minute), now.Add(-time.Hour)); err != nil {
		t.Fatalf("create expired: %v", err)
	}
	// expired challenge must not verify even with the right answer (F-001)
	ok, err := r.ConsumeChallenge("cap-old", "hash-old", now)
	if err != nil || ok {
		t.Fatalf("expired consume = %v, %v; want false (F-001)", ok, err)
	}
	// and it is consumed
	ok, err = r.ConsumeChallenge("cap-old", "hash-old", now)
	if err != nil || ok {
		t.Fatalf("expired re-consume = %v, %v; want false", ok, err)
	}
}

func TestEnabledDefaultsFalseAndToggles(t *testing.T) {
	r := newTestEnv(t)
	now := time.Now().UTC()
	enabled, err := r.Enabled()
	if err != nil {
		t.Fatalf("enabled: %v", err)
	}
	if enabled {
		t.Fatal("captcha must default to disabled (D-001 §5)")
	}
	if err := r.SetEnabled(true, now); err != nil {
		t.Fatalf("set enabled: %v", err)
	}
	enabled, err = r.Enabled()
	if err != nil || !enabled {
		t.Fatalf("enabled after set = %v, %v; want true", enabled, err)
	}
	if err := r.SetEnabled(false, now.Add(time.Minute)); err != nil {
		t.Fatalf("set disabled: %v", err)
	}
	enabled, err = r.Enabled()
	if err != nil || enabled {
		t.Fatalf("enabled after unset = %v, %v; want false", enabled, err)
	}
}

func TestCreateChallengeLazyPurgesExpired(t *testing.T) {
	r := newTestEnv(t)
	now := time.Now().UTC()
	if err := r.CreateChallenge("cap-old", "hash-old", now.Add(-time.Minute), now.Add(-time.Hour)); err != nil {
		t.Fatalf("create expired: %v", err)
	}
	if err := r.CreateChallenge("cap-new", "hash-new", now.Add(5*time.Minute), now); err != nil {
		t.Fatalf("create new: %v", err)
	}
	ok, err := r.ConsumeChallenge("cap-old", "hash-old", now)
	if err != nil || ok {
		t.Fatalf("expired challenge verifies: %v %v", ok, err)
	}
	ok, err = r.ConsumeChallenge("cap-new", "hash-new", now)
	if err != nil || !ok {
		t.Fatalf("new challenge = %v, %v; want true", ok, err)
	}
}
