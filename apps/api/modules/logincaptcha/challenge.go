// Arithmetic captcha challenge service (S-11 · GOAL-011 D-002 §1): generates
// small add/sub questions, stores hashed answers, verifies one-time with
// expiry, and exposes the Required()/Verify() surface the login handler
// consumes (handler.CaptchaVerifier) plus the Generate()/SetEnabled() surface
// the captcha routes consume (handler.CaptchaService).
package logincaptcha

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/modules/logincaptcha/store"
)

// challengeTTL bounds a challenge lifetime (D-002 §1).
const challengeTTL = 5 * time.Minute

// ErrInvalidCaptcha is returned when verification fails for any reason
// (unknown, expired, consumed, or wrong answer — the caller maps it to
// INVALID_CAPTCHA without leaking which).
var ErrInvalidCaptcha = errors.New("captcha verification failed")

// Service implements the login-facing captcha surface.
type Service struct {
	repository *store.Repository
	now        func() time.Time
}

// NewService constructs the captcha service.
func NewService(repository *store.Repository) *Service {
	return &Service{repository: repository, now: time.Now}
}

// Generate issues a new arithmetic challenge and persists it, returning the
// public projection (id, question, expiresInSeconds) the preflight route
// serves (D-002 §2).
func (s *Service) Generate() (id, question string, expiresInSeconds int64, err error) {
	now := s.now().UTC()
	a := randInt(1, 50)
	b := randInt(1, 50)
	op := "+"
	if randInt(0, 1) == 1 {
		op = "-"
		if b > a {
			a, b = b, a
		}
	}
	var answer int
	if op == "+" {
		answer = a + b
	} else {
		answer = a - b
	}
	id = newChallengeID()
	expiresAt := now.Add(challengeTTL)
	hash := answerHash(id, fmt.Sprint(answer))
	if err := s.repository.CreateChallenge(id, hash, expiresAt, now); err != nil {
		return "", "", 0, err
	}
	return id, fmt.Sprintf("%d %s %d = ?", a, op, b), int64(challengeTTL.Seconds()), nil
}

// Required reports whether login must present a captcha. Fail-closed on
// config read errors (grok A-003 F-006): if the switch cannot be read the gate
// is treated as ON so an operator must fix the config before logins resume.
func (s *Service) Required() bool {
	enabled, err := s.repository.Enabled()
	if err != nil {
		return true
	}
	return enabled
}

// SetEnabled flips the login captcha gate (D-002 §3).
func (s *Service) SetEnabled(enabled bool, now time.Time) error {
	return s.repository.SetEnabled(enabled, now)
}

// Verify checks a submitted challenge one-time: the challenge is consumed
// (deleted) on ANY attempt — success or failure — so a challenge cannot be
// brute-forced; expiry is enforced inside the same transaction (D-002 §1;
// grok A-003 F-001/F-004). Any failure — unknown, expired, consumed, wrong
// answer or store error — maps to ErrInvalidCaptcha without leaking which.
func (s *Service) Verify(captchaID, answer string, now time.Time) error {
	if captchaID == "" {
		return ErrInvalidCaptcha
	}
	ok, err := s.repository.ConsumeChallenge(captchaID, answerHash(captchaID, answer), now)
	if err != nil {
		return ErrInvalidCaptcha
	}
	if !ok {
		return ErrInvalidCaptcha
	}
	return nil
}

func answerHash(id, answer string) string {
	sum := sha256.Sum256([]byte(id + ":" + answer))
	return hex.EncodeToString(sum[:])
}

func randInt(min, max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		return min
	}
	return min + int(n.Int64())
}

func newChallengeID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "cap-" + fmt.Sprint(time.Now().UnixNano())
	}
	return hex.EncodeToString(b[:])
}
