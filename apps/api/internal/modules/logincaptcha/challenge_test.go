// Service tests for the login captcha (S-11 · GOAL-011 D-002 §1/§3): the
// arithmetic challenge lifecycle, the default-off gate and the one-time
// verify contract.
package logincaptcha

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/logincaptcha/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func newServiceEnv(t *testing.T) *Service {
	t.Helper()
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "test.db"), "admin", "hash", true)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return NewService(store.NewRepository(st))
}

func TestGenerateIssuesPersistedChallenge(t *testing.T) {
	s := newServiceEnv(t)
	id, question, expiresInSeconds, err := s.Generate()
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if id == "" || question == "" {
		t.Fatalf("challenge = %q %q, want non-empty id/question", id, question)
	}
	if expiresInSeconds != 300 {
		t.Fatalf("expiresInSeconds = %d, want 300", expiresInSeconds)
	}
}

func TestGateDefaultsOff(t *testing.T) {
	s := newServiceEnv(t)
	if s.Required() {
		t.Fatal("gate must default to off (D-001 §5)")
	}
}

func TestSetEnabledFlipsGate(t *testing.T) {
	s := newServiceEnv(t)
	now := time.Now().UTC()
	if err := s.SetEnabled(true, now); err != nil {
		t.Fatalf("enable: %v", err)
	}
	if !s.Required() {
		t.Fatal("gate must be on after enable")
	}
}

func TestVerifyConsumesChallengeOneTime(t *testing.T) {
	s := newServiceEnv(t)
	if err := s.SetEnabled(true, time.Now().UTC()); err != nil {
		t.Fatalf("enable: %v", err)
	}
	id, question, _, err := s.Generate()
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	// parse "a op b = ?" from the question
	var a, b int
	var op string
	if _, err := fmt.Sscanf(question, "%d %s %d = ?", &a, &op, &b); err != nil {
		t.Fatalf("parse question %q: %v", question, err)
	}
	answer := ""
	if op == "+" {
		answer = strconv.Itoa(a + b)
	} else {
		answer = strconv.Itoa(a - b)
	}
	if err := s.Verify(id, answer, time.Now().UTC()); err != nil {
		t.Fatalf("verify correct answer: %v", err)
	}
	// second attempt must fail: challenge consumed
	if err := s.Verify(id, answer, time.Now().UTC()); !errors.Is(err, ErrInvalidCaptcha) {
		t.Fatalf("verify after consume = %v, want ErrInvalidCaptcha", err)
	}
}

func TestVerifyWrongAnswerFailsAndConsumes(t *testing.T) {
	s := newServiceEnv(t)
	if err := s.SetEnabled(true, time.Now().UTC()); err != nil {
		t.Fatalf("enable: %v", err)
	}
	id, _, _, err := s.Generate()
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if err := s.Verify(id, "999999", time.Now().UTC()); !errors.Is(err, ErrInvalidCaptcha) {
		t.Fatalf("verify wrong = %v, want ErrInvalidCaptcha", err)
	}
	// consumed: correct answer also fails now
	if err := s.Verify(id, "2", time.Now().UTC()); !errors.Is(err, ErrInvalidCaptcha) {
		t.Fatalf("verify after wrong attempt = %v, want ErrInvalidCaptcha", err)
	}
}

func TestVerifyExpiredChallengeFails(t *testing.T) {
	s := newServiceEnv(t)
	now := time.Now().UTC()
	// Create an expired challenge directly through the repository (the service
	// Generate always uses a fresh TTL). The stored hash matches the scheme
	// (answerHash(id, answer)) so the test isolates the EXPIRY rejection from
	// a wrong-answer rejection (grok A-004 F-010).
	expired := now.Add(-time.Minute)
	correctHash := answerHash("cap-expired", "42")
	if err := s.repository.CreateChallenge("cap-expired", correctHash, expired, now.Add(-time.Hour)); err != nil {
		t.Fatalf("create expired: %v", err)
	}
	if err := s.Verify("cap-expired", "42", now); !errors.Is(err, ErrInvalidCaptcha) {
		t.Fatalf("verify expired = %v, want ErrInvalidCaptcha (F-001)", err)
	}
}

func TestVerifyEmptyOrUnknownChallenge(t *testing.T) {
	s := newServiceEnv(t)
	if err := s.Verify("", "2", time.Now().UTC()); !errors.Is(err, ErrInvalidCaptcha) {
		t.Fatalf("verify empty id = %v, want ErrInvalidCaptcha", err)
	}
	if err := s.Verify("cap-unknown", "2", time.Now().UTC()); !errors.Is(err, ErrInvalidCaptcha) {
		t.Fatalf("verify unknown id = %v, want ErrInvalidCaptcha", err)
	}
}

// F-011 (grok A-004): Required() must fail CLOSED on config read errors —
// an unreadable switch treats the gate as ON.
func TestRequiredFailsClosedOnConfigError(t *testing.T) {
	failing := &failingRunner{err: errors.New("store unavailable")}
	s := &Service{repository: store.NewRepository(failing), now: time.Now}
	if !s.Required() {
		t.Fatal("Required() must return true when the config cannot be read (F-011)")
	}
}

// failingRunner is a TxRunner that always fails, simulating a store outage.
type failingRunner struct{ err error }

func (r *failingRunner) Run(_ context.Context, _ func(kernel.Tx) error) error {
	return r.err
}
