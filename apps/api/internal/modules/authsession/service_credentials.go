package authsession

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"slices"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/pagination"
)

// ServiceCredential is the metadata projection for a machine credential.
// The raw secret is intentionally not represented and is never persisted.
type ServiceCredential struct {
	ID          string
	Name        string
	TokenPrefix string
	TokenHash   string
	Scopes      []string
	ExpiresAt   time.Time
	RevokedAt   *time.Time
	LastUsedAt  *time.Time
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ServiceCredentialAudit is executed inside the credential mutation
// transaction. Returning an error rolls the mutation back.
type ServiceCredentialAudit func(kernel.Tx) error

// ServiceCredentialRevokeAudit is the revoke-specific audit callback.
type ServiceCredentialRevokeAudit func(kernel.Tx, ServiceCredential) error

var ErrCredentialNameTaken = errors.New("authsession: service credential name already taken")

// CreateServiceCredential inserts one credential and records its audit event
// in the same transaction when audit is non-nil.
func (r *Repository) CreateServiceCredential(credential ServiceCredential, audit ServiceCredentialAudit) error {
	credential.Name = strings.TrimSpace(credential.Name)
	credential.Scopes = sortedScopes(credential.Scopes)
	scopesJSON, err := json.Marshal(credential.Scopes)
	if err != nil {
		return fmt.Errorf("marshal service credential scopes: %w", err)
	}
	err = r.withTx("create service credential", func(tx kernel.Tx) error {
		if _, err := tx.Exec(context.Background(), `
INSERT INTO service_credentials
  (id, name, token_prefix, token_hash, scopes, expires_at, created_by, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			credential.ID, credential.Name, credential.TokenPrefix, credential.TokenHash,
			string(scopesJSON), credential.ExpiresAt.Unix(), credential.CreatedBy,
			credential.CreatedAt.Unix(), credential.UpdatedAt.Unix()); err != nil {
			// W9 F-011: the unique-violation predicate is dialect-agnostic, and
			// the name constraint is identified by both dialects' names (sqlite
			// auto-index "service_credentials.name"; postgres UNIQUE DDL name
			// "service_credentials_name_key").
			lowered := strings.ToLower(err.Error())
			if kernel.IsUniqueViolation(err) &&
				(strings.Contains(lowered, "service_credentials.name") || strings.Contains(lowered, "service_credentials_name_key")) {
				return ErrCredentialNameTaken
			}
			return fmt.Errorf("insert service credential: %w", err)
		}
		if audit != nil {
			if err := audit(tx); err != nil {
				return fmt.Errorf("record service credential create audit: %w", err)
			}
		}
		return nil
	})
	return err
}

// ListServiceCredentials returns metadata ordered newest first by contract.
func (r *Repository) ListServiceCredentials(page, pageSize int) ([]ServiceCredential, int, error) {
	var items []ServiceCredential
	var total int
	err := r.withTx("list service credentials", func(tx kernel.Tx) error {
		if err := tx.QueryRow(context.Background(), `SELECT COUNT(*) FROM service_credentials`).Scan(&total); err != nil {
			return fmt.Errorf("count service credentials: %w", err)
		}
		rows, err := tx.Query(context.Background(), `
SELECT id, name, token_prefix, token_hash, scopes, expires_at, revoked_at,
       last_used_at, created_by, created_at, updated_at
FROM service_credentials ORDER BY created_at DESC, id DESC LIMIT ? OFFSET ?`,
			pageSize, pagination.Offset(page, pageSize, total))
		if err != nil {
			return fmt.Errorf("query service credentials: %w", err)
		}
		defer rows.Close()
		items = make([]ServiceCredential, 0, pageSize)
		for rows.Next() {
			credential, err := scanServiceCredential(rows)
			if err != nil {
				return err
			}
			items = append(items, *credential)
		}
		if err := rows.Err(); err != nil {
			return fmt.Errorf("list service credentials rows: %w", err)
		}
		return nil
	})
	return items, total, err
}

// ServiceCredentialByID fetches one credential metadata projection.
func (r *Repository) ServiceCredentialByID(id string) (*ServiceCredential, error) {
	return r.serviceCredentialBy(`id = ?`, id)
}

// ServiceCredentialByHash fetches one credential by its persisted SHA-256.
func (r *Repository) ServiceCredentialByHash(tokenHash string) (*ServiceCredential, error) {
	return r.serviceCredentialBy(`token_hash = ?`, tokenHash)
}

func (r *Repository) serviceCredentialBy(where string, arg string) (*ServiceCredential, error) {
	var credential ServiceCredential
	err := r.withTx("get service credential", func(tx kernel.Tx) error {
		got, err := scanServiceCredential(tx.QueryRow(context.Background(), `
SELECT id, name, token_prefix, token_hash, scopes, expires_at, revoked_at,
       last_used_at, created_by, created_at, updated_at
FROM service_credentials WHERE `+where, arg))
		if err != nil {
			return err
		}
		credential = *got
		return nil
	})
	if err != nil {
		if errors.Is(err, kernel.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &credential, nil
}

// RevokeServiceCredential revokes an active credential. Repeated revocation
// is idempotent and returns changed=false without writing another audit row.
func (r *Repository) RevokeServiceCredential(id string, now time.Time, audit ServiceCredentialRevokeAudit) (*ServiceCredential, bool, error) {
	var credential ServiceCredential
	changed := false
	err := r.withTx("revoke service credential", func(tx kernel.Tx) error {
		got, err := scanServiceCredential(tx.QueryRow(context.Background(), `
SELECT id, name, token_prefix, token_hash, scopes, expires_at, revoked_at,
       last_used_at, created_by, created_at, updated_at
FROM service_credentials WHERE id = ?`, id))
		if errors.Is(err, kernel.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		credential = *got
		if credential.RevokedAt != nil {
			return nil
		}
		if _, err := tx.Exec(context.Background(), `UPDATE service_credentials SET revoked_at = ?, updated_at = ? WHERE id = ? AND revoked_at IS NULL`, now.Unix(), now.Unix(), id); err != nil {
			return fmt.Errorf("revoke service credential: %w", err)
		}
		revokedAt := now.UTC()
		credential.RevokedAt = &revokedAt
		credential.UpdatedAt = now.UTC()
		changed = true
		if audit != nil {
			if err := audit(tx, credential); err != nil {
				return fmt.Errorf("record service credential revoke audit: %w", err)
			}
		}
		return nil
	})
	return &credential, changed, err
}

// MarkServiceCredentialUsed is best effort metadata bookkeeping for requests.
func (r *Repository) MarkServiceCredentialUsed(id string, now time.Time) error {
	return r.withTx("mark service credential used", func(tx kernel.Tx) error {
		_, err := tx.Exec(context.Background(), `UPDATE service_credentials SET last_used_at = ?, updated_at = ? WHERE id = ? AND revoked_at IS NULL`, now.Unix(), now.Unix(), id)
		return err
	})
}

// MarkServiceCredentialUsedWithAudit updates usage metadata and appends the
// corresponding audit row on one transaction. Any audit error rolls back the
// metadata update, keeping credential usage and its audit evidence aligned.
func (r *Repository) MarkServiceCredentialUsedWithAudit(id string, now time.Time, audit ServiceCredentialAudit) error {
	if audit == nil {
		return errors.New("authsession: service credential use audit is required")
	}
	return r.withTx("mark service credential used with audit", func(tx kernel.Tx) error {
		if _, err := tx.Exec(context.Background(), `UPDATE service_credentials SET last_used_at = ?, updated_at = ? WHERE id = ? AND revoked_at IS NULL`, now.Unix(), now.Unix(), id); err != nil {
			return err
		}
		if err := audit(tx); err != nil {
			return fmt.Errorf("record service credential use audit: %w", err)
		}
		return nil
	})
}

func sortedScopes(scopes []string) []string {
	seen := make(map[string]struct{}, len(scopes))
	result := make([]string, 0, len(scopes))
	for _, scope := range scopes {
		scope = strings.TrimSpace(scope)
		if scope == "" {
			continue
		}
		if _, ok := seen[scope]; ok {
			continue
		}
		seen[scope] = struct{}{}
		result = append(result, scope)
	}
	slices.Sort(result)
	return result
}

func scanServiceCredential(row interface{ Scan(...any) error }) (*ServiceCredential, error) {
	var credential ServiceCredential
	var scopesJSON string
	var expiresAt, createdAt, updatedAt int64
	var revokedAt, lastUsedAt sql.NullInt64
	if err := row.Scan(&credential.ID, &credential.Name, &credential.TokenPrefix, &credential.TokenHash,
		&scopesJSON, &expiresAt, &revokedAt, &lastUsedAt, &credential.CreatedBy, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(scopesJSON), &credential.Scopes); err != nil {
		return nil, fmt.Errorf("decode service credential scopes: %w", err)
	}
	credential.Scopes = sortedScopes(credential.Scopes)
	credential.ExpiresAt = time.Unix(expiresAt, 0).UTC()
	credential.CreatedAt = time.Unix(createdAt, 0).UTC()
	credential.UpdatedAt = time.Unix(updatedAt, 0).UTC()
	if revokedAt.Valid {
		value := time.Unix(revokedAt.Int64, 0).UTC()
		credential.RevokedAt = &value
	}
	if lastUsedAt.Valid {
		value := time.Unix(lastUsedAt.Int64, 0).UTC()
		credential.LastUsedAt = &value
	}
	return &credential, nil
}
