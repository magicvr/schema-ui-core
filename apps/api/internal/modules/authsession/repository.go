// Package authsession owns account, refresh-token, user-management, and RBAC
// persistence. The platform store supplies only the transaction boundary.
package authsession

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// TxRunner is the platform persistence boundary consumed by the repository.
// Domain code never imports the concrete store implementation.
type TxRunner interface {
	WithTx(context.Context, func(*sql.Tx) error) error
}

// Repository owns the auth-session and RBAC domain queries.
type Repository struct {
	runner TxRunner
}

// NewRepository constructs the module-owned repository over a platform
// transaction runner.
func NewRepository(runner TxRunner) *Repository {
	return &Repository{runner: runner}
}

// User is the persisted identity row backing account.Session.User.
type User struct {
	ID           string
	Username     string
	Name         string
	Roles        []string
	PasswordHash string
	// TokenVersion is a per-user monotonic counter (W4 P0-3): incremented on
	// password change; the auth middleware rejects access-token JWTs issued at
	// an older version, revoking already-signed tokens immediately.
	TokenVersion int
	// FailedLoginCount counts consecutive failed password attempts (GOAL-004
	// S4-6 account lock); reset on successful login.
	FailedLoginCount int
	// LockedUntil is the unix-second lock window end (0 = not locked). Locks
	// expire automatically once now passes the window.
	LockedUntil int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// RefreshToken is a stored opaque refresh token; only its hash is persisted.
type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

// UserFilter carries handler-validated user list parameters.
type UserFilter struct {
	Q        string
	Sort     string
	Order    string
	Page     int
	PageSize int
}

// UserPatch is the editable subset of a managed user.
type UserPatch struct {
	Name         *string
	Roles        *[]string
	PasswordHash *string
}

// Role is the persisted RBAC role projection.
type Role struct {
	ID            string
	Key           string
	Name          string
	System        bool
	Permissions   []string
	MenuItems     []string
	AssignedUsers int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// RolePatch applies partial changes to a user-created role.
type RolePatch struct {
	Name        *string
	Permissions *[]string
	MenuItems   *[]string
}

// RoleFilter carries handler-validated role list parameters.
type RoleFilter struct {
	Q        string
	Sort     string
	Order    string
	Page     int
	PageSize int
}

// Domain sentinels are owned by core.auth-session rather than the platform
// persistence package.
var (
	ErrNotFound          = errors.New("authsession: not found")
	ErrAlreadyRevoked    = errors.New("authsession: refresh token already revoked")
	ErrUsernameTaken     = errors.New("authsession: username already taken")
	ErrLastAdmin         = errors.New("authsession: cannot remove the last admin")
	ErrSelfOperation     = errors.New("authsession: self operation is not allowed")
	ErrInvalidRole       = errors.New("authsession: invalid role reference")
	ErrRoleTaken         = errors.New("authsession: role key already taken")
	ErrRoleInUse         = errors.New("authsession: role is in use by users")
	ErrRoleSystem        = errors.New("authsession: role is a system role")
	ErrInvalidKey        = errors.New("authsession: invalid role key")
	ErrInvalidPermission = errors.New("authsession: invalid permission reference")
	ErrInvalidMenuItem   = errors.New("authsession: invalid menu item reference")
)

func (r *Repository) withTx(operation string, fn func(*sql.Tx) error) error {
	if r == nil || r.runner == nil {
		return fmt.Errorf("%s: authsession repository is not configured", operation)
	}
	if err := r.runner.WithTx(context.Background(), fn); err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	return nil
}

func dedupeKeys(keys []string) []string {
	seen := make(map[string]bool, len(keys))
	out := make([]string, 0, len(keys))
	for _, key := range keys {
		if !seen[key] {
			seen[key] = true
			out = append(out, key)
		}
	}
	return out
}

func sameRoleSet(a, b []string) bool {
	as := make(map[string]bool, len(a))
	for _, value := range a {
		as[value] = true
	}
	bs := make(map[string]bool, len(b))
	for _, value := range b {
		bs[value] = true
	}
	if len(as) != len(bs) {
		return false
	}
	for value := range as {
		if !bs[value] {
			return false
		}
	}
	return true
}
