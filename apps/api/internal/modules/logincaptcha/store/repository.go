// Package store owns the admin.login-captcha persistence (S-11 · GOAL-011
// D-002 `1/`3): one-time challenge rows (hashed answers) and the single-row
// enable flag. Lives in a sub-package so the handler can consume the types
// without an import cycle with the module provider.
package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// TxRunner is the platform persistence boundary consumed by the repository.
type TxRunner interface {
	WithTx(context.Context, func(*sql.Tx) error) error
}

// Repository owns the captcha domain queries.
type Repository struct {
	runner TxRunner
}

// NewRepository constructs the captcha repository over a platform transaction
// runner.
func NewRepository(runner TxRunner) *Repository {
	return &Repository{runner: runner}
}

// ErrChallengeNotFound is returned when a challenge id is unknown.
var ErrChallengeNotFound = errors.New("captcha challenge not found")

// CreateChallenge stores one challenge.
func (r *Repository) CreateChallenge(id, answerHash string, expiresAt, now time.Time) error {
	return r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		// Best-effort lazy purge of expired rows (D-002 `1).
		_, _ = tx.Exec(`DELETE FROM captcha_challenges WHERE expires_at <= ?`, now.Unix())
		_, err := tx.Exec(
			`INSERT INTO captcha_challenges (id, answer_hash, expires_at, created_at) VALUES (?, ?, ?, ?)`,
			id, answerHash, expiresAt.Unix(), now.Unix(),
		)
		if err != nil {
			return fmt.Errorf("insert captcha challenge: %w", err)
		}
		return nil
	})
}

// ConsumeChallenge atomically verifies-and-deletes one challenge: the row is
	// removed on ANY attempt (success or failure) so a challenge cannot be
	// brute-forced, and expiry is enforced inside the same transaction
	// (S-11 · GOAL-011 D-002 §1; grok A-003 F-001/F-004). Returns true only
	// when the challenge existed, was unexpired and the answer hash matched.
	func (r *Repository) ConsumeChallenge(id, answerHash string, now time.Time) (bool, error) {
		matched := false
		err := r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
			var stored string
			var expiresAt int64
			row := tx.QueryRow(`SELECT answer_hash, expires_at FROM captcha_challenges WHERE id = ?`, id)
			if err := row.Scan(&stored, &expiresAt); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil // unknown challenge: nothing to consume
				}
				return fmt.Errorf("get captcha challenge: %w", err)
			}
			// Consume on ANY attempt — success or failure (D-002 §1). A delete
			// failure fails the whole verify (fail-closed, A-003 F-004).
			if _, err := tx.Exec(`DELETE FROM captcha_challenges WHERE id = ?`, id); err != nil {
				return fmt.Errorf("delete captcha challenge: %w", err)
			}
			if expiresAt <= now.Unix() || stored != answerHash {
				return nil // expired or wrong answer: consumed, not matched
			}
			matched = true
			return nil
		})
		if err != nil {
			return false, err
		}
		return matched, nil
	}

// Enabled reports the config switch.
func (r *Repository) Enabled() (bool, error) {
	var enabled int
	err := r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		row := tx.QueryRow(`SELECT enabled FROM captcha_config WHERE id = 1`)
		if err := row.Scan(&enabled); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				enabled = 0 // default disabled
				return nil
			}
			return fmt.Errorf("get captcha config: %w", err)
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return enabled == 1, nil
}

// SetEnabled flips the config switch.
func (r *Repository) SetEnabled(enabled bool, now time.Time) error {
	return r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		_, err := tx.Exec(
			`INSERT INTO captcha_config (id, enabled, created_at, updated_at) VALUES (1, ?, ?, ?)
			 ON CONFLICT(id) DO UPDATE SET enabled = excluded.enabled, updated_at = excluded.updated_at`,
			boolInt(enabled), now.Unix(), now.Unix(),
		)
		if err != nil {
			return fmt.Errorf("set captcha config: %w", err)
		}
		return nil
	})
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
