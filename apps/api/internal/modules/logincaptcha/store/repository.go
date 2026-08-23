// Package store owns the admin.login-captcha persistence (S-11 · GOAL-011
// D-002 `1/`3): one-time challenge rows (hashed answers) and the single-row
// enable flag. Lives in a sub-package so the handler can consume the types
// without an import cycle with the module provider.
package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// TxRunner is the platform persistence boundary consumed by the repository
// (kernel port; R4 — the public surface no longer exposes *sql.Tx).
type TxRunner interface {
	Run(context.Context, func(kernel.Tx) error) error
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
	return r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		// Best-effort lazy purge of expired rows (D-002 `1).
		_, _ = tx.Exec(context.Background(), `DELETE FROM captcha_challenges WHERE expires_at <= ?`, now.Unix())
		_, err := tx.Exec(context.Background(),
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
//
// W11 F-005: the guarded DELETE is the single statement that decides
// success — there is no read-then-delete window. Under READ COMMITTED two
// concurrent transactions cannot both claim the same row: the guarded delete
// wins for exactly one of them (RowsAffected == 1); the loser affects 0 rows
// and is treated as unmatched. The second best-effort delete preserves the
// consume-on-ANY-attempt contract (a wrong/expired first attempt removes the
// row so it cannot be brute-forced afterwards).
func (r *Repository) ConsumeChallenge(id, answerHash string, now time.Time) (bool, error) {
	matched := false
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		res, err := tx.Exec(context.Background(),
			`DELETE FROM captcha_challenges WHERE id = ? AND expires_at > ? AND answer_hash = ?`,
			id, now.Unix(), answerHash,
		)
		if err != nil {
			return fmt.Errorf("consume captcha challenge: %w", err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("consume captcha challenge rows: %w", err)
		}
		matched = affected == 1
		// Consume on ANY attempt — success or failure (D-002 §1). The
		// guarded statement above is the atomicity gate; this cleanup is a
		// same-transaction no-op for the winner and removes the row for a
		// wrong-answer or expired first attempt (fail-closed retry safety).
		if _, err := tx.Exec(context.Background(), `DELETE FROM captcha_challenges WHERE id = ?`, id); err != nil {
			return fmt.Errorf("consume captcha challenge cleanup: %w", err)
		}
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
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		row := tx.QueryRow(context.Background(), `SELECT enabled FROM captcha_config WHERE id = 1`)
		if err := row.Scan(&enabled); err != nil {
			if errors.Is(err, kernel.ErrNoRows) {
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
	return r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		_, err := tx.Exec(context.Background(),
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
