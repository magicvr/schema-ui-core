// Package store owns the admin.mfa persistence (S-10 · GOAL-017 D-002 §2):
// per-user TOTP state (pending/active) and one-time login proofs. It lives in
// a sub-package so the handler can consume row types without an import cycle
// with the module provider.
package store

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"time"
)

// TxRunner is the platform persistence boundary consumed by the repository.
type TxRunner interface {
	Run(context.Context, func(kernel.Tx) error) error
}

// Repository owns the MFA domain queries.
type Repository struct {
	runner TxRunner
}

// NewRepository constructs the MFA repository over a platform transaction
// runner.
func NewRepository(runner TxRunner) *Repository {
	return &Repository{runner: runner}
}

// Domain sentinels mapped by the handler to frozen error codes.
var (
	ErrNotFound       = errors.New("mfa row not found")
	ErrActiveConflict = errors.New("mfa enrollment is active; disable it first")
)

// State is one user_mfa row.
type State struct {
	UserID            string
	Status            string // pending | active
	SecretCiphertext  string
	RecoveryCodesHash string
	LastUsedStep      int64
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// Proof is one mfa_proofs row.
type Proof struct {
	ID        string
	UserID    string
	FailCount int
	ExpiresAt time.Time
	CreatedAt time.Time
}

// GetState returns the MFA state for one user.
func (r *Repository) GetState(userID string) (*State, error) {
	var s State
	var created, updated int64
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		err := tx.QueryRow(context.Background(),
			`SELECT user_id, status, totp_secret_ciphertext, recovery_codes_hash, last_used_step, created_at, updated_at
			 FROM user_mfa WHERE user_id = ?`, userID,
		).Scan(&s.UserID, &s.Status, &s.SecretCiphertext, &s.RecoveryCodesHash, &s.LastUsedStep, &created, &updated)
		if errors.Is(err, kernel.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("get mfa state: %w", err)
		}
		s.CreatedAt = time.Unix(created, 0)
		s.UpdatedAt = time.Unix(updated, 0)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// UpsertPending creates (or resets) a pending enrollment. Re-enrolling
// overwrites a previous PENDING row (A-005 recommended); an ACTIVE row is
// rejected with ErrActiveConflict — tearing down an active enrollment
// requires the second factor (A-007 F-002).
func (r *Repository) UpsertPending(userID, secretCiphertext, recoveryCodesHash string, now time.Time) error {
	return r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		res, err := tx.Exec(context.Background(),
			`INSERT INTO user_mfa (user_id, status, totp_secret_ciphertext, recovery_codes_hash, last_used_step, created_at, updated_at)
			 VALUES (?, 'pending', ?, ?, 0, ?, ?)
			 ON CONFLICT(user_id) DO UPDATE SET
			   status = 'pending',
			   totp_secret_ciphertext = excluded.totp_secret_ciphertext,
			   recovery_codes_hash = excluded.recovery_codes_hash,
			   last_used_step = 0,
			   updated_at = excluded.updated_at
			 WHERE user_mfa.status = 'pending'`,
			userID, secretCiphertext, recoveryCodesHash, now.Unix(), now.Unix(),
		)
		if err != nil {
			return fmt.Errorf("upsert pending mfa: %w", err)
		}
		// A filtered conflict-update (0 rows) means the row exists and is
		// ACTIVE — enrollment must not tear down an active enrollment without
		// the second factor (A-007 F-002).
		if n, _ := res.RowsAffected(); n == 0 {
			return ErrActiveConflict
		}
		return nil
	})
}

// Activate flips a pending enrollment to active (only pending rows may
// activate; the login gate requires active — A-004 F-001 response).
func (r *Repository) Activate(userID string, now time.Time) error {
	return r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		res, err := tx.Exec(context.Background(),
			`UPDATE user_mfa SET status = 'active', updated_at = ? WHERE user_id = ? AND status = 'pending'`,
			now.Unix(), userID,
		)
		if err != nil {
			return fmt.Errorf("activate mfa: %w", err)
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return ErrNotFound
		}
		return nil
	})
}

// DeleteState removes the MFA enrollment (disable / admin reset).
func (r *Repository) DeleteState(userID string) error {
	return r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		if _, err := tx.Exec(context.Background(), `DELETE FROM user_mfa WHERE user_id = ?`, userID); err != nil {
			return fmt.Errorf("delete mfa state: %w", err)
		}
		return nil
	})
}

// UpdateRecoveryCodes replaces the recovery-code hash set.
func (r *Repository) UpdateRecoveryCodes(userID, recoveryCodesHash string, now time.Time) error {
	return r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		res, err := tx.Exec(context.Background(),
			`UPDATE user_mfa SET recovery_codes_hash = ?, updated_at = ? WHERE user_id = ?`,
			recoveryCodesHash, now.Unix(), userID,
		)
		if err != nil {
			return fmt.Errorf("update recovery codes: %w", err)
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return ErrNotFound
		}
		return nil
	})
}

// SetLastUsedStep records the last consumed TOTP time step (replay window).
func (r *Repository) SetLastUsedStep(userID string, step int64, now time.Time) error {
	return r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		if _, err := tx.Exec(context.Background(),
			`UPDATE user_mfa SET last_used_step = ?, updated_at = ? WHERE user_id = ?`,
			step, now.Unix(), userID,
		); err != nil {
			return fmt.Errorf("set last used step: %w", err)
		}
		return nil
	})
}

// CreateProof inserts a one-time login proof (5-minute TTL, GOAL-017 D-002 §3).
func (r *Repository) CreateProof(userID string, expiresAt time.Time, now time.Time) (*Proof, error) {
	var idBytes [16]byte
	if _, err := rand.Read(idBytes[:]); err != nil {
		return nil, fmt.Errorf("mfa proof id: %w", err)
	}
	id := hex.EncodeToString(idBytes[:])
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		_, err := tx.Exec(context.Background(),
			`INSERT INTO mfa_proofs (id, user_id, fail_count, expires_at, created_at) VALUES (?, ?, 0, ?, ?)`,
			id, userID, expiresAt.Unix(), now.Unix(),
		)
		if err != nil {
			return fmt.Errorf("create mfa proof: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &Proof{ID: id, UserID: userID, FailCount: 0, ExpiresAt: expiresAt, CreatedAt: now}, nil
}

// GetProof returns one proof row (used by verify before consuming).
func (r *Repository) GetProof(id string) (*Proof, error) {
	var p Proof
	var expires, created int64
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		err := tx.QueryRow(context.Background(),
			`SELECT id, user_id, fail_count, expires_at, created_at FROM mfa_proofs WHERE id = ?`, id,
		).Scan(&p.ID, &p.UserID, &p.FailCount, &expires, &created)
		if errors.Is(err, kernel.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("get mfa proof: %w", err)
		}
		p.ExpiresAt = time.Unix(expires, 0)
		p.CreatedAt = time.Unix(created, 0)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// IncrementProofFailures counts one failed second-factor attempt.
func (r *Repository) IncrementProofFailures(id string, now time.Time) error {
	return r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		if _, err := tx.Exec(context.Background(),
			`UPDATE mfa_proofs SET fail_count = fail_count + 1 WHERE id = ?`, id,
		); err != nil {
			return fmt.Errorf("increment mfa proof failures: %w", err)
		}
		return nil
	})
}

// DeleteProof consumes a proof (successful verify or exhaustion).
func (r *Repository) DeleteProof(id string) error {
	return r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		if _, err := tx.Exec(context.Background(), `DELETE FROM mfa_proofs WHERE id = ?`, id); err != nil {
			return fmt.Errorf("delete mfa proof: %w", err)
		}
		return nil
	})
}
