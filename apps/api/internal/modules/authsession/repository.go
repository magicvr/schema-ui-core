// Package authsession owns account, refresh-token, user-management, and RBAC
// persistence. The platform store supplies only the transaction boundary.
package authsession

import (
	"context"
	"errors"
	"fmt"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"time"
)

// TxRunner is the platform persistence boundary consumed by the repository.
// Domain code never imports the concrete store implementation.
type TxRunner interface {
	Run(context.Context, func(kernel.Tx) error) error
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
	// MFAEnabled is the product-state second-factor flag (S-10 · GOAL-017
	// D-002 §4, migration 0029): true when user_mfa has an active enrollment.
	// Cross-module read of the admin.mfa table (contract = migration 0029).
	MFAEnabled bool
	// Enabled is the product-state account flag (F-03 · migration 0013):
	// false = disabled by an admin (login/refresh/middleware fail closed;
	// disable also bumps TokenVersion and revokes all refresh tokens).
	Enabled bool
	// FailedLoginCount counts consecutive failed password attempts (GOAL-004
	// S4-6 account lock); reset on successful login.
	FailedLoginCount int
	// LockedUntil is the unix-second lock window end (0 = not locked). Locks
	// expire automatically once now passes the window.
	LockedUntil int64
	// AvatarURL is the self-service avatar asset URL (W13 T-05 · migration
	// 0035, account module): "" = no avatar. Values are committed by the
	// account profile PATCH and validated against the avatar store.
	AvatarURL string
	// MustChangePassword is true when the user has not yet replaced the
	// initial/reset password (W16-F01).
	MustChangePassword bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
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
	// T-02 (GOAL-013 D-003): optional management-list filters (nil = no
	// constraint). Enabled matches the enabled flag; Locked matches rows
	// currently lock-expired (locked_until in the future).
	Enabled *bool
	Locked  *bool
}

// UserPatch is the editable subset of a managed user.
type UserPatch struct {
	Name         *string
	Roles        *[]string
	PasswordHash *string
	// AvatarURL is the self-service avatar value (W13 T-05); nil = untouched.
	AvatarURL *string
	// MustChangePassword is the forced-password-change flag (W16-F01); nil =
	// untouched.
	MustChangePassword *bool
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
	// T-02 (GOAL-013 D-003): optional system-flag filter (nil = no constraint).
	System *bool
}

// PermissionCatalogEntry is one row of the RBAC permission catalog (W11 ·
// U-02 dynamic permission options).
type PermissionCatalogEntry struct {
	Key         string
	Description string
}

// MenuItemCatalogEntry is one row of the RBAC navigation catalog (W11 · U-02
// dynamic menu-access options). Label is a deterministic display derivation
// from the page reference (admin console page ids are English tokens).
type MenuItemCatalogEntry struct {
	ID      string
	PageRef string
	Label   string
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

func (r *Repository) withTx(operation string, fn func(kernel.Tx) error) error {
	if r == nil || r.runner == nil {
		return fmt.Errorf("%s: authsession repository is not configured", operation)
	}
	if err := r.runner.Run(context.Background(), fn); err != nil {
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
