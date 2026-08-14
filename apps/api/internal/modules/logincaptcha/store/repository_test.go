// Package store tests for the admin.login-captcha persistence (S-11 ·
// GOAL-011 D-002 §1/§3): one-time challenge rows and the single-row switch.
package store

import (
	"context"
	"errors"
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

func TestChallengeLifecycle(t *testing.T) {
	r := newTestEnv(t)
	now := time.Now().UTC()
	if err := r.CreateChallenge("cap-1", "hash-1", now.Add(5*time.Minute), now); err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := r.GetChallenge("cap-1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got == nil || *got != "hash-1" {
		t.Fatalf("answer hash = %v, want hash-1", got)
	}
	if err := r.DeleteChallenge("cap-1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := r.GetChallenge("cap-1"); !errors.Is(err, ErrChallengeNotFound) {
		t.Fatalf("get after delete = %v, want ErrChallengeNotFound", err)
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
	if _, err := r.GetChallenge("cap-old"); !errors.Is(err, ErrChallengeNotFound) {
		t.Fatalf("expired challenge readable: %v", err)
	}
	if got, err := r.GetChallenge("cap-new"); err != nil || got == nil || *got != "hash-new" {
		t.Fatalf("new challenge = %v, %v", got, err)
	}
}

var _ = context.Background
