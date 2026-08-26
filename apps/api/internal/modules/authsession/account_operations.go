// Account-enable and session operations (F-03 · GOAL-005 D-002).
//
// These repository methods back the admin.account module: self-service session
// listing/revocation and admin enable/disable/unlock. The users table hosts the
// enabled column (migration 0013); the product-state semantics are owned by
// admin.account while persistence stays in core.auth-session.
package authsession

import (
	"context"
	"errors"
	"fmt"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"time"
)

// ErrSessionNotFound is returned when a session id does not exist or does not
// belong to the requesting identity (fail closed: no cross-user oracle).
var ErrSessionNotFound = errors.New("authsession: session not found")

// ListRefreshTokensForUser returns every refresh token of a user (active and
// revoked), newest first. Access tokens are short-lived and never listed.
func (r *Repository) ListRefreshTokensForUser(userID string) ([]RefreshToken, error) {
	var tokens []RefreshToken
	err := r.withTx("list user refresh tokens", func(tx kernel.Tx) error {
		rows, err := tx.Query(context.Background(),
			`SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
			 FROM refresh_tokens WHERE user_id = ?
			 ORDER BY created_at DESC, id DESC`, userID)
		if err != nil {
			return fmt.Errorf("query refresh tokens: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var token RefreshToken
			var expiresAt int64
			var revokedAt *int64
			var createdAt int64
			if err := rows.Scan(&token.ID, &token.UserID, &token.TokenHash, &expiresAt, &revokedAt, &createdAt); err != nil {
				return fmt.Errorf("scan refresh token: %w", err)
			}
			token.ExpiresAt = time.Unix(expiresAt, 0).UTC()
			if revokedAt != nil {
				value := time.Unix(*revokedAt, 0).UTC()
				token.RevokedAt = &value
			}
			token.CreatedAt = time.Unix(createdAt, 0).UTC()
			tokens = append(tokens, token)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

// RevokeRefreshTokenIfOwned revokes a refresh token only when it belongs to
// userID. Idempotent: an already-revoked owned token is a no-op success;
// unknown or foreign ids fail closed with ErrSessionNotFound.
func (r *Repository) RevokeRefreshTokenIfOwned(id, userID string, now time.Time) error {
	return r.withTx("revoke owned refresh token", func(tx kernel.Tx) error {
		res, err := tx.Exec(context.Background(),
			`UPDATE refresh_tokens SET revoked_at = ? WHERE id = ? AND user_id = ? AND revoked_at IS NULL`,
			now.Unix(), id, userID)
		if err != nil {
			return fmt.Errorf("revoke owned refresh token: %w", err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("rows affected: %w", err)
		}
		if affected == 0 {
			var owned int
			if err := tx.QueryRow(context.Background(),
				`SELECT COUNT(*) FROM refresh_tokens WHERE id = ? AND user_id = ?`, id, userID,
			).Scan(&owned); err != nil {
				return fmt.Errorf("check owned refresh token: %w", err)
			}
			if owned == 0 {
				return ErrSessionNotFound
			}
			// owned but already revoked: idempotent success
		}
		return nil
	})
}

// SetUserEnabled flips the product-state enabled flag. Disabling (enabled=false)
// is a security action: it bumps token_version (all issued access tokens die
// immediately) and revokes every live refresh token. Guards mirror the
// management surface: cannot disable self, cannot disable the last admin.
// Enabling only flips the flag (lock state is cleared by UnlockUser).
func (r *Repository) SetUserEnabled(id string, enabled bool, actorID string, now time.Time) (*User, error) {
	var updated *User
	err := r.withTx("set user enabled", func(tx kernel.Tx) error {
		var current User
		var rolesJSON string
		var createdAt, updatedAt int64
		var mustChangePassword int
		err := tx.QueryRow(context.Background(),
			`SELECT id, username, name, roles, password_hash, token_version, failed_login_count, locked_until, enabled, avatar_url, must_change_password, created_at, updated_at FROM users WHERE id = ?`, id,
		).Scan(&current.ID, &current.Username, &current.Name, &rolesJSON, &current.PasswordHash, &current.TokenVersion, &current.FailedLoginCount, &current.LockedUntil, &current.Enabled, &current.AvatarURL, &mustChangePassword, &createdAt, &updatedAt)
		if errors.Is(err, kernel.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("get user for enable toggle: %w", err)
		}
		if !enabled {
			if id == actorID {
				return ErrSelfOperation
			}
			isAdmin, err := userHasRoleKey(tx, id, "admin")
			if err != nil {
				return err
			}
			if isAdmin {
				// F-001 (A-003 independent): the last-admin invariant must hold
				// for ENABLED admins. Counting all admins lets a delegated
				// users.disable holder disable the last enabled admin while a
				// disabled admin still exists (recovery would then require
				// direct DB access, since users.enable is admin-only).
				other, err := countEnabledAdminUsersExcluding(tx, id)
				if err != nil {
					return err
				}
				if other == 0 {
					return ErrLastAdmin
				}
			}
		}
		nextTokenVersion := current.TokenVersion
		if !enabled && current.Enabled {
			nextTokenVersion = current.TokenVersion + 1
			if _, err := tx.Exec(context.Background(),
				`UPDATE refresh_tokens SET revoked_at = COALESCE(revoked_at, ?) WHERE user_id = ?`, now.Unix(), id,
			); err != nil {
				return fmt.Errorf("revoke refresh tokens on disable: %w", err)
			}
		}
		if _, err := tx.Exec(context.Background(),
			`UPDATE users SET enabled = ?, token_version = ?, updated_at = ? WHERE id = ?`,
			boolInt(enabled), nextTokenVersion, now.Unix(), id,
		); err != nil {
			return fmt.Errorf("update user enabled: %w", err)
		}
		// Post-update invariant (F-001 hardening): after a disable lands, at
		// least one ENABLED admin must remain. SQLite serializes writers, so a
		// concurrent disable either fails closed (busy) or hits this check;
		// either way zero-enabled-admin cannot commit.
		if !enabled {
			remaining, err := countEnabledAdminUsers(tx)
			if err != nil {
				return err
			}
			if remaining == 0 {
				return ErrLastAdmin
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	updated, err = r.GetUser(id)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

// UnlockUser manually clears the account-lock window and the consecutive
// failure counter (C-11 lock state). It does not touch the enabled flag.
// GOAL-014 D-002: an admin unlock also clears every per-(account|source)
// pair row — otherwise a source locked before the unlock would keep denying
// that client after the global state was wiped.
func (r *Repository) UnlockUser(id string, now time.Time) (*User, error) {
	err := r.withTx("unlock user", func(tx kernel.Tx) error {
		res, err := tx.Exec(context.Background(),
			`UPDATE users SET locked_until = 0, failed_login_count = 0, last_login_failure_at = 0, updated_at = ? WHERE id = ?`,
			now.Unix(), id)
		if err != nil {
			return fmt.Errorf("unlock user: %w", err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("rows affected: %w", err)
		}
		if affected == 0 {
			var exists int
			if err := tx.QueryRow(context.Background(), `SELECT COUNT(*) FROM users WHERE id = ?`, id).Scan(&exists); err != nil {
				return fmt.Errorf("check unlock user: %w", err)
			}
			if exists == 0 {
				return ErrNotFound
			}
		}
		if _, err := tx.Exec(context.Background(),
			`DELETE FROM login_failures WHERE user_id = ?`, id,
		); err != nil {
			return fmt.Errorf("clear source lock rows: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.GetUser(id)
}

// userHasRoleKey reports whether a user holds a role key (via user_roles).
func userHasRoleKey(tx kernel.Tx, userID, key string) (bool, error) {
	var count int
	if err := tx.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM user_roles ur JOIN roles r ON r.id = ur.role_id
		 WHERE ur.user_id = ? AND r.key = ?`, userID, key,
	).Scan(&count); err != nil {
		return false, fmt.Errorf("check user role %s: %w", key, err)
	}
	return count > 0, nil
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

// countEnabledAdminUsersExcluding counts ENABLED admin users other than id
// (F-001): the last-admin guard for enable/disable considers only accounts
// that can actually administer (enabled=1), so disabling the last enabled
// admin fails closed even when disabled admins still exist.
func countEnabledAdminUsersExcluding(tx kernel.Tx, id string) (int, error) {
	var count int
	if err := tx.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM user_roles ur
		 JOIN roles r ON r.id = ur.role_id
		 JOIN users u ON u.id = ur.user_id
		 WHERE r.key = 'admin' AND u.enabled = 1 AND ur.user_id != ?`, id,
	).Scan(&count); err != nil {
		return 0, fmt.Errorf("count enabled admin users: %w", err)
	}
	return count, nil
}

// countEnabledAdminUsers counts every ENABLED admin user (post-update
// invariant check for disable).
func countEnabledAdminUsers(tx kernel.Tx) (int, error) {
	var count int
	if err := tx.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM user_roles ur
		 JOIN roles r ON r.id = ur.role_id
		 JOIN users u ON u.id = ur.user_id
		 WHERE r.key = 'admin' AND u.enabled = 1`,
	).Scan(&count); err != nil {
		return 0, fmt.Errorf("count enabled admin users: %w", err)
	}
	return count, nil
}
