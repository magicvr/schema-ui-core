package authsession

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"
)

var roleKeyRe = regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)

// CreateUser inserts an internal account and keeps the legacy roles JSON and
// normalized user_roles relation set-equal in one transaction.
func (r *Repository) CreateUser(user User) error {
	return r.withTx("create user", func(tx *sql.Tx) error {
		roles := dedupeKeys(user.Roles)
		rolesJSON, err := json.Marshal(roles)
		if err != nil {
			return fmt.Errorf("marshal roles: %w", err)
		}
		if _, err := tx.Exec(
			`INSERT INTO users (id, username, name, roles, password_hash, must_change_password, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			user.ID, user.Username, user.Name, string(rolesJSON), user.PasswordHash,
			boolInt(user.MustChangePassword), user.CreatedAt.Unix(), user.UpdatedAt.Unix(),
		); err != nil {
			return fmt.Errorf("insert user: %w", err)
		}
		now := time.Now().UTC().Unix()
		for _, key := range roles {
			if err := linkUserRole(tx, user.ID, key, now); err != nil {
				return fmt.Errorf("create user %s role %s: %w", user.ID, key, err)
			}
		}
		return nil
	})
}

// UserByUsername fetches an identity by its unique username.
func (r *Repository) UserByUsername(username string) (*User, error) {
	return r.userBy("get user by username", func(tx *sql.Tx) *sql.Row {
		return tx.QueryRow(
			`SELECT id, username, name, roles, password_hash, token_version, failed_login_count, locked_until, enabled, avatar_url, must_change_password, created_at, updated_at
			 FROM users WHERE username = ?`, username)
	})
}

// UserByID fetches an identity by primary key.
func (r *Repository) UserByID(id string) (*User, error) {
	return r.userBy("get user by id", func(tx *sql.Tx) *sql.Row {
		return tx.QueryRow(
			`SELECT id, username, name, roles, password_hash, token_version, failed_login_count, locked_until, enabled, avatar_url, must_change_password, created_at, updated_at
			 FROM users WHERE id = ?`, id)
	})
}

func (r *Repository) userBy(operation string, query func(*sql.Tx) *sql.Row) (*User, error) {
	var user *User
	err := r.withTx(operation, func(tx *sql.Tx) error {
		var err error
		user, err = userWithRoles(tx, query(tx))
		return err
	})
	return user, err
}

func userWithRoles(tx *sql.Tx, row *sql.Row) (*User, error) {
	user, err := scanUser(row)
	if err != nil {
		return nil, err
	}
	roles, err := rolesForUser(tx, user.ID)
	if err != nil {
		return nil, err
	}
	if !sameRoleSet(user.Roles, roles) {
		return nil, fmt.Errorf("authsession: user %s role mismatch: legacy %v normalized %v", user.ID, user.Roles, roles)
	}
	user.Roles = roles
	return user, nil
}

func rolesForUser(tx *sql.Tx, userID string) ([]string, error) {
	rows, err := tx.Query(
		`SELECT r.key FROM user_roles ur JOIN roles r ON r.id = ur.role_id
		 WHERE ur.user_id = ? ORDER BY r.key`, userID)
	if err != nil {
		return nil, fmt.Errorf("query normalized roles: %w", err)
	}
	defer rows.Close()
	var keys []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, fmt.Errorf("scan normalized role: %w", err)
		}
		keys = append(keys, key)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("query normalized roles: %w", err)
	}
	return keys, nil
}

// PermissionsForUser returns the normalized permission projection.
func (r *Repository) PermissionsForUser(userID string) ([]string, error) {
	var keys []string
	err := r.withTx("permissions for user", func(tx *sql.Tx) error {
		rows, err := tx.Query(
			`SELECT DISTINCT p.key
			 FROM user_roles ur
			 JOIN role_permissions rp ON rp.role_id = ur.role_id
			 JOIN permissions p ON p.id = rp.permission_id
			 WHERE ur.user_id = ? ORDER BY p.key`, userID)
		if err != nil {
			return fmt.Errorf("query permissions: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var key string
			if err := rows.Scan(&key); err != nil {
				return fmt.Errorf("scan permission: %w", err)
			}
			keys = append(keys, key)
		}
		if err := rows.Err(); err != nil {
			return fmt.Errorf("query permissions: %w", err)
		}
		return nil
	})
	return keys, err
}

// FeaturesForUser returns the enabled menu feature projection.
func (r *Repository) FeaturesForUser(userID string) (map[string]bool, error) {
	features := make(map[string]bool)
	err := r.withTx("features for user", func(tx *sql.Tx) error {
		rows, err := tx.Query(
			`SELECT m.feature_key, EXISTS(
				SELECT 1 FROM user_roles ur
				JOIN role_menu_items rmi ON rmi.role_id = ur.role_id
				WHERE ur.user_id = ? AND rmi.menu_item_id = m.id
			 )
			 FROM menu_items m
			 WHERE m.enabled = 1
			 ORDER BY m.feature_key`, userID)
		if err != nil {
			return fmt.Errorf("query features: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var key string
			var granted bool
			if err := rows.Scan(&key, &granted); err != nil {
				return fmt.Errorf("scan feature: %w", err)
			}
			features[key] = granted
		}
		if err := rows.Err(); err != nil {
			return fmt.Errorf("query features: %w", err)
		}
		return nil
	})
	return features, err
}

// RecordLoginFailure bumps the user's consecutive-failure counter and, past
// the threshold, opens a lock window ending at lockedUntil (GOAL-004 S4-6).
// It reports whether the lock opened on this failure.
func (r *Repository) RecordLoginFailure(userID string, threshold int, lockedUntil time.Time, now time.Time) (bool, error) {
	locked := false
	err := r.withTx("record login failure", func(tx *sql.Tx) error {
		var count int
		if err := tx.QueryRow(`SELECT failed_login_count FROM users WHERE id = ?`, userID).Scan(&count); err != nil {
			return fmt.Errorf("read failed login count: %w", err)
		}
		next := count + 1
		lockedUntilUnix := int64(0)
		if next >= threshold {
			lockedUntilUnix = lockedUntil.Unix()
			next = 0 // counter resets when the lock opens; it counts consecutive failures
			locked = true
		}
		if _, err := tx.Exec(
			`UPDATE users SET failed_login_count = ?, locked_until = ?, updated_at = ? WHERE id = ?`,
			next, lockedUntilUnix, now.Unix(), userID,
		); err != nil {
			return fmt.Errorf("record login failure: %w", err)
		}
		return nil
	})
	return locked, err
}

// ResetLoginFailures clears the consecutive-failure counter and lock window on
// a successful login.
func (r *Repository) ResetLoginFailures(userID string, now time.Time) error {
	return r.withTx("reset login failures", func(tx *sql.Tx) error {
		if _, err := tx.Exec(
			`UPDATE users SET failed_login_count = 0, locked_until = 0, updated_at = ? WHERE id = ?`,
			now.Unix(), userID,
		); err != nil {
			return fmt.Errorf("reset login failures: %w", err)
		}
		return nil
	})
}

// BumpTokenVersionAndRevokeAll bumps the user's token_version (invalidating
// every already-signed access token, W4 P0-3) and revokes every live refresh
// token in one transaction (S-10 · GOAL-017 D-002 §4: MFA disable / admin
// reset force re-authentication at the same strength as password change).
func (r *Repository) BumpTokenVersionAndRevokeAll(userID string, now time.Time) error {
	return r.withTx("bump token version and revoke refresh tokens", func(tx *sql.Tx) error {
		if _, err := tx.Exec(
			`UPDATE users SET token_version = token_version + 1, updated_at = ? WHERE id = ?`,
			now.Unix(), userID,
		); err != nil {
			return fmt.Errorf("bump token version: %w", err)
		}
		if _, err := tx.Exec(
			`UPDATE refresh_tokens SET revoked_at = ? WHERE user_id = ? AND revoked_at IS NULL`,
			now.Unix(), userID,
		); err != nil {
			return fmt.Errorf("revoke user refresh tokens: %w", err)
		}
		return nil
	})
}

// RevokeAllRefreshTokensForUser revokes every live refresh token of the user
// (account-lock hardening: a locked account's sessions must not keep rotating).
func (r *Repository) RevokeAllRefreshTokensForUser(userID string, now time.Time) error {
	return r.withTx("revoke user refresh tokens", func(tx *sql.Tx) error {
		if _, err := tx.Exec(
			`UPDATE refresh_tokens SET revoked_at = ? WHERE user_id = ? AND revoked_at IS NULL`,
			now.Unix(), userID,
		); err != nil {
			return fmt.Errorf("revoke user refresh tokens: %w", err)
		}
		return nil
	})
}

func scanUser(row interface{ Scan(...any) error }) (*User, error) {
	var user User
	var roles string
	var createdAt, updatedAt int64
	var mustChangePassword int
	err := row.Scan(&user.ID, &user.Username, &user.Name, &roles, &user.PasswordHash, &user.TokenVersion, &user.FailedLoginCount, &user.LockedUntil, &user.Enabled, &user.AvatarURL, &mustChangePassword, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan user: %w", err)
	}
	if err := json.Unmarshal([]byte(roles), &user.Roles); err != nil {
		return nil, fmt.Errorf("unmarshal roles: %w", err)
	}
	user.MustChangePassword = mustChangePassword != 0
	user.CreatedAt = time.Unix(createdAt, 0).UTC()
	user.UpdatedAt = time.Unix(updatedAt, 0).UTC()
	return &user, nil
}

// CreateRefreshToken persists a hashed opaque refresh token.
func (r *Repository) CreateRefreshToken(token RefreshToken) error {
	return r.withTx("create refresh token", func(tx *sql.Tx) error {
		var revokedAt any
		if token.RevokedAt != nil {
			revokedAt = token.RevokedAt.Unix()
		}
		if _, err := tx.Exec(
			`INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at)
			 VALUES (?, ?, ?, ?, ?, ?)`,
			token.ID, token.UserID, token.TokenHash, token.ExpiresAt.Unix(), revokedAt, token.CreatedAt.Unix(),
		); err != nil {
			return fmt.Errorf("insert refresh token: %w", err)
		}
		return nil
	})
}

// RefreshTokenByHash fetches a refresh token by its stored hash.
func (r *Repository) RefreshTokenByHash(hash string) (*RefreshToken, error) {
	var token RefreshToken
	err := r.withTx("get refresh token", func(tx *sql.Tx) error {
		var expiresAt int64
		var revokedAt *int64
		var createdAt int64
		err := tx.QueryRow(
			`SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
			 FROM refresh_tokens WHERE token_hash = ?`, hash,
		).Scan(&token.ID, &token.UserID, &token.TokenHash, &expiresAt, &revokedAt, &createdAt)
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("scan refresh token: %w", err)
		}
		token.ExpiresAt = time.Unix(expiresAt, 0).UTC()
		if revokedAt != nil {
			value := time.Unix(*revokedAt, 0).UTC()
			token.RevokedAt = &value
		}
		token.CreatedAt = time.Unix(createdAt, 0).UTC()
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// RevokeRefreshToken marks a refresh token revoked atomically. The guarded
// UPDATE (revoked_at IS NULL) makes concurrent revocations race-safe: exactly
// one caller wins, every later caller gets ErrAlreadyRevoked instead of a
// check-then-act window that would let two rotations both proceed.
func (r *Repository) RevokeRefreshToken(id string, now time.Time) error {
	return r.withTx("revoke refresh token", func(tx *sql.Tx) error {
		res, err := tx.Exec(`UPDATE refresh_tokens SET revoked_at = ? WHERE id = ? AND revoked_at IS NULL`, now.Unix(), id)
		if err != nil {
			return fmt.Errorf("update refresh token: %w", err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("rows affected: %w", err)
		}
		if affected == 0 {
			var exists int64
			if err := tx.QueryRow(`SELECT COUNT(*) FROM refresh_tokens WHERE id = ?`, id).Scan(&exists); err != nil {
				return fmt.Errorf("check refresh token: %w", err)
			}
			if exists == 0 {
				return ErrNotFound
			}
			return ErrAlreadyRevoked
		}
		return nil
	})
}

func ensureRole(tx *sql.Tx, key string, now int64) error {
	if _, err := tx.Exec(
		`INSERT INTO roles (id, key, name, system, created_at, updated_at)
		 VALUES (?, ?, ?, 0, ?, ?)
		 ON CONFLICT(id) DO NOTHING`,
		"role-"+key, key, key, now, now,
	); err != nil {
		return fmt.Errorf("ensure role %s: %w", key, err)
	}
	return nil
}

func linkUserRole(tx *sql.Tx, userID, key string, now int64) error {
	if !roleKeyRe.MatchString(key) {
		return fmt.Errorf("invalid role key %q", key)
	}
	if err := ensureRole(tx, key, now); err != nil {
		return err
	}
	if _, err := tx.Exec(
		`INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)
		 ON CONFLICT(user_id, role_id) DO NOTHING`, userID, "role-"+key,
	); err != nil {
		return fmt.Errorf("link user %s role %s: %w", userID, key, err)
	}
	return nil
}
