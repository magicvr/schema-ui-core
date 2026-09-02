// Package subject provides channel-agnostic external subject seam
// (issuer, external_id) -> subject_id (VP-029 · GOAL-002).
// External subjects are NOT stored in admin.users.
package subject

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// Domain sentinels.
var (
	ErrNotFound       = errors.New("subject not found")
	ErrInvalidSubject = errors.New("invalid subject issuer or external_id")
)

// Subject represents an external subject mapping.
type Subject struct {
	ID         string    `json:"id"`
	Issuer     string    `json:"issuer"`
	ExternalID string    `json:"externalId"`
	CreatedAt  time.Time `json:"createdAt"`
}

// TxRunner is the persistence boundary.
type TxRunner interface {
	Run(context.Context, func(kernel.Tx) error) error
}

// Store provides persistence for external subjects.
type Store struct {
	runner TxRunner
}

// NewStore constructs a new subject Store over a platform transaction runner.
func NewStore(runner TxRunner) *Store {
	return &Store{runner: runner}
}

// newID generates a time-ordered millisecond prefix + random hex ID.
func newID(now time.Time) (string, error) {
	randBytes := make([]byte, 12)
	if _, err := rand.Read(randBytes); err != nil {
		return "", fmt.Errorf("subject: random bytes: %w", err)
	}
	return fmt.Sprintf("%016x%s", now.UnixMilli(), hex.EncodeToString(randBytes)), nil
}

// isUniqueViolation detects unique constraint violations portably.
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique") || strings.Contains(msg, "23505")
}

// GetOrCreateSubject retrieves an existing subject or creates a new one idempotently.
func (s *Store) GetOrCreateSubject(ctx context.Context, issuer, externalID string, now time.Time) (*Subject, bool, error) {
	issuer = strings.TrimSpace(issuer)
	externalID = strings.TrimSpace(externalID)
	if issuer == "" || externalID == "" {
		return nil, false, ErrInvalidSubject
	}

	// 1. Fast path: check if existing.
	if existing, err := s.GetSubjectByExternalID(ctx, issuer, externalID); err == nil {
		return existing, false, nil
	} else if !errors.Is(err, ErrNotFound) {
		return nil, false, fmt.Errorf("get subject by external id: %w", err)
	}

	// 2. Create new row.
	id, err := newID(now)
	if err != nil {
		return nil, false, err
	}

	insertErr := s.runner.Run(ctx, func(tx kernel.Tx) error {
		_, err := tx.Exec(ctx,
			`INSERT INTO subjects (id, issuer, external_id, created_at) VALUES (?, ?, ?, ?)`,
			id, issuer, externalID, now.Unix(),
		)
		return err
	})

	if insertErr != nil {
		if !isUniqueViolation(insertErr) {
			return nil, false, fmt.Errorf("insert subject: %w", insertErr)
		}
		// Concurrent create won: reload in a fresh transaction.
		existing, err := s.GetSubjectByExternalID(ctx, issuer, externalID)
		if err != nil {
			return nil, false, fmt.Errorf("reload subject after conflict: %w", err)
		}
		return existing, false, nil
	}

	return &Subject{
		ID:         id,
		Issuer:     issuer,
		ExternalID: externalID,
		CreatedAt:  now.Truncate(time.Second),
	}, true, nil
}

// GetSubject retrieves a subject by its internal subject_id.
func (s *Store) GetSubject(ctx context.Context, id string) (*Subject, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrNotFound
	}
	var sub Subject
	var created int64
	err := s.runner.Run(ctx, func(tx kernel.Tx) error {
		row := tx.QueryRow(ctx,
			`SELECT id, issuer, external_id, created_at FROM subjects WHERE id = ?`,
			id,
		)
		if err := row.Scan(&sub.ID, &sub.Issuer, &sub.ExternalID, &created); err != nil {
			if errors.Is(err, kernel.ErrNoRows) {
				return ErrNotFound
			}
			return fmt.Errorf("query subject: %w", err)
		}
		sub.CreatedAt = time.Unix(created, 0)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

// SubjectExists checks if a subject with the given ID exists.
func (s *Store) SubjectExists(ctx context.Context, id string) (bool, error) {
	sub, err := s.GetSubject(ctx, id)
	if errors.Is(err, ErrNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return sub != nil, nil
}

// GetSubjectByExternalID retrieves a subject by (issuer, external_id).
func (s *Store) GetSubjectByExternalID(ctx context.Context, issuer, externalID string) (*Subject, error) {
	issuer = strings.TrimSpace(issuer)
	externalID = strings.TrimSpace(externalID)
	if issuer == "" || externalID == "" {
		return nil, ErrNotFound
	}
	var sub Subject
	var created int64
	err := s.runner.Run(ctx, func(tx kernel.Tx) error {
		row := tx.QueryRow(ctx,
			`SELECT id, issuer, external_id, created_at FROM subjects WHERE issuer = ? AND external_id = ?`,
			issuer, externalID,
		)
		if err := row.Scan(&sub.ID, &sub.Issuer, &sub.ExternalID, &created); err != nil {
			if errors.Is(err, kernel.ErrNoRows) {
				return ErrNotFound
			}
			return fmt.Errorf("query subject by external id: %w", err)
		}
		sub.CreatedAt = time.Unix(created, 0)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &sub, nil
}
