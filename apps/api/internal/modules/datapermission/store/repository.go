// Package store owns the admin.data-permission persistence (S-09 · GOAL-016
// D-002 §4): per-resource scope policies and user × resource assignments. It
// lives in a sub-package so the handler can consume the row types without an
// import cycle with the module provider.
package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"time"
)

// TxRunner is the platform persistence boundary consumed by the repository.
type TxRunner interface {
	Run(context.Context, func(kernel.Tx) error) error
}

// Repository owns the data-permission domain queries.
type Repository struct {
	runner TxRunner
}

// NewRepository constructs the data-permission repository over a platform
// transaction runner.
func NewRepository(runner TxRunner) *Repository {
	return &Repository{runner: runner}
}

// Domain sentinels mapped by the handler to frozen error codes.
var (
	ErrNotFound       = errors.New("data permission row not found")
	ErrInvalidScope   = errors.New("invalid scope type")
	ErrNotEnforceable = errors.New("resource is not enforceable")
)

// Scope types (v1; org deferred to B-10, I-011-001 §5).
const (
	ScopeAll  = "all"
	ScopeSelf = "self"
)

// ValidScope reports whether s is a known scope type.
func ValidScope(s string) bool { return s == ScopeAll || s == ScopeSelf }

// Policy is one data_scope_policies row.
type Policy struct {
	Resource     string
	OwnerColumn  string
	DefaultScope string
	Enabled      bool
	UpdatedAt    time.Time
}

// Assignment is one user_data_scopes row.
type Assignment struct {
	UserID    string
	Resource  string
	ScopeType string
	UpdatedAt time.Time
}

// ListPolicies returns all registered scope policies ordered by resource.
func (r *Repository) ListPolicies() ([]Policy, error) {
	policies := []Policy{}
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		rows, err := tx.Query(context.Background(),
			`SELECT resource, owner_column, default_scope, enabled, updated_at FROM data_scope_policies ORDER BY resource`)
		if err != nil {
			return fmt.Errorf("list scope policies: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var p Policy
			var updated int64
			var enabled int
			if err := rows.Scan(&p.Resource, &p.OwnerColumn, &p.DefaultScope, &enabled, &updated); err != nil {
				return fmt.Errorf("scan scope policy: %w", err)
			}
			p.Enabled = enabled != 0
			p.UpdatedAt = time.Unix(updated, 0)
			policies = append(policies, p)
		}
		return rows.Err()
	})
	return policies, err
}

// GetPolicy returns one policy by resource id.
func (r *Repository) GetPolicy(resource string) (*Policy, error) {
	var p Policy
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		var updated int64
		var enabled int
		err := tx.QueryRow(context.Background(),
			`SELECT resource, owner_column, default_scope, enabled, updated_at FROM data_scope_policies WHERE resource = ?`,
			resource,
		).Scan(&p.Resource, &p.OwnerColumn, &p.DefaultScope, &enabled, &updated)
		if errors.Is(err, kernel.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("get scope policy: %w", err)
		}
		p.Enabled = enabled != 0
		p.UpdatedAt = time.Unix(updated, 0)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// UpsertPolicy registers or updates one resource scope policy. default_scope
// is required (A-004 F-001: no implicit default) and validated.
func (r *Repository) UpsertPolicy(resource, ownerColumn, defaultScope string, enabled bool, now time.Time) error {
	if !ValidScope(defaultScope) {
		return ErrInvalidScope
	}
	enabledInt := 0
	if enabled {
		enabledInt = 1
	}
	return r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		_, err := tx.Exec(context.Background(),
			`INSERT INTO data_scope_policies (resource, owner_column, default_scope, enabled, updated_at)
			 VALUES (?, ?, ?, ?, ?)
			 ON CONFLICT(resource) DO UPDATE SET
			   owner_column = excluded.owner_column,
			   default_scope = excluded.default_scope,
			   enabled = excluded.enabled,
			   updated_at = excluded.updated_at`,
			resource, ownerColumn, defaultScope, enabledInt, now.Unix(),
		)
		if err != nil {
			return fmt.Errorf("upsert scope policy: %w", err)
		}
		return nil
	})
}

// ListAssignments returns the scope assignments for one user.
func (r *Repository) ListAssignments(userID string) ([]Assignment, error) {
	assignments := []Assignment{}
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		rows, err := tx.Query(context.Background(),
			`SELECT user_id, resource, scope_type, updated_at FROM user_data_scopes WHERE user_id = ? ORDER BY resource`,
			userID,
		)
		if err != nil {
			return fmt.Errorf("list scope assignments: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var a Assignment
			var updated int64
			if err := rows.Scan(&a.UserID, &a.Resource, &a.ScopeType, &updated); err != nil {
				return fmt.Errorf("scan scope assignment: %w", err)
			}
			a.UpdatedAt = time.Unix(updated, 0)
			assignments = append(assignments, a)
		}
		return rows.Err()
	})
	return assignments, err
}

// GetAssignment returns one user × resource assignment.
func (r *Repository) GetAssignment(userID, resource string) (*Assignment, error) {
	var a Assignment
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		var updated int64
		err := tx.QueryRow(context.Background(),
			`SELECT user_id, resource, scope_type, updated_at FROM user_data_scopes WHERE user_id = ? AND resource = ?`,
			userID, resource,
		).Scan(&a.UserID, &a.Resource, &a.ScopeType, &updated)
		if errors.Is(err, kernel.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("get scope assignment: %w", err)
		}
		a.UpdatedAt = time.Unix(updated, 0)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// UpsertAssignments upserts one user's assignments (map resource → scope_type).
// A scope of "all" is stored explicitly (overrides a self default).
func (r *Repository) UpsertAssignments(userID string, scopes map[string]string, now time.Time) error {
	for _, s := range scopes {
		if !ValidScope(s) {
			return ErrInvalidScope
		}
	}
	return r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		for resource, scopeType := range scopes {
			if _, err := tx.Exec(context.Background(),
				`INSERT INTO user_data_scopes (user_id, resource, scope_type, updated_at)
				 VALUES (?, ?, ?, ?)
				 ON CONFLICT(user_id, resource) DO UPDATE SET
				   scope_type = excluded.scope_type,
				   updated_at = excluded.updated_at`,
				userID, resource, scopeType, now.Unix(),
			); err != nil {
				return fmt.Errorf("upsert scope assignment: %w", err)
			}
		}
		return nil
	})
}
