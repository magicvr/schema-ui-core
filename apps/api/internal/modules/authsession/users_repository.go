package authsession

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

// ListUsers returns the filtered management projection and total count.
func (r *Repository) ListUsers(filter UserFilter) ([]User, int, error) {
	var items []User
	var total int
	err := r.withTx("list users", func(tx *sql.Tx) error {
		where, args := usersWhere(filter.Q)
		if err := tx.QueryRow(`SELECT COUNT(*) FROM users`+where, args...).Scan(&total); err != nil {
			return fmt.Errorf("count users: %w", err)
		}

		rows, err := tx.Query(
			`SELECT u.id, u.username, u.name, u.roles, u.password_hash, u.token_version, u.failed_login_count, u.locked_until, u.enabled, u.created_at, u.updated_at,
			        EXISTS(SELECT 1 FROM user_mfa um WHERE um.user_id = u.id AND um.status = 'active') AS mfa_enabled
			 FROM users u`+where+
				` ORDER BY `+usersSortSQL(filter.Sort, filter.Order)+`, u.id ASC`+
				` LIMIT ? OFFSET ?`,
			append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)...,
		)
		if err != nil {
			return fmt.Errorf("query users: %w", err)
		}
		items = make([]User, 0, filter.PageSize)
		for rows.Next() {
			user, err := scanUserListRow(rows)
			if err != nil {
				_ = rows.Close()
				return err
			}
			items = append(items, *user)
		}
		if err := rows.Err(); err != nil {
			_ = rows.Close()
			return fmt.Errorf("list users rows: %w", err)
		}
		if err := rows.Close(); err != nil {
			return fmt.Errorf("close users rows: %w", err)
		}
		for i := range items {
			roles, err := rolesForUser(tx, items[i].ID)
			if err != nil {
				return err
			}
			if !sameRoleSet(items[i].Roles, roles) {
				return fmt.Errorf("authsession: user %s role mismatch: legacy %v normalized %v", items[i].ID, items[i].Roles, roles)
			}
			items[i].Roles = roles
		}
		return nil
	})
	return items, total, err
}

// GetUser returns one managed user.
func (r *Repository) GetUser(id string) (*User, error) {
	return r.UserByID(id)
}

// CreateUserManagement creates a user without implicitly creating role keys.
func (r *Repository) CreateUserManagement(user User) (*User, error) {
	err := r.withTx("create managed user", func(tx *sql.Tx) error {
		roles := dedupeKeys(user.Roles)
		for _, key := range roles {
			var count int
			if err := tx.QueryRow(`SELECT COUNT(*) FROM roles WHERE key = ?`, key).Scan(&count); err != nil {
				return fmt.Errorf("validate role %s: %w", key, err)
			}
			if count == 0 {
				return ErrInvalidRole
			}
		}
		var exists int
		if err := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)`, user.Username).Scan(&exists); err != nil {
			return fmt.Errorf("check username: %w", err)
		}
		if exists == 1 {
			return ErrUsernameTaken
		}
		rolesJSON, err := json.Marshal(roles)
		if err != nil {
			return fmt.Errorf("marshal roles: %w", err)
		}
		if _, err := tx.Exec(
			`INSERT INTO users (id, username, name, roles, password_hash, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`,
			user.ID, user.Username, user.Name, string(rolesJSON), user.PasswordHash,
			user.CreatedAt.Unix(), user.UpdatedAt.Unix(),
		); err != nil {
			return fmt.Errorf("insert user: %w", err)
		}
		for _, key := range roles {
			if _, err := tx.Exec(
				`INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)`, user.ID, "role-"+key,
			); err != nil {
				return fmt.Errorf("link user %s role %s: %w", user.ID, key, err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.GetUser(user.ID)
}

// UpdateUser applies management changes and enforces self/last-admin guards.
func (r *Repository) UpdateUser(id string, patch UserPatch, actorID string, now time.Time) (*User, error) {
	err := r.withTx("update managed user", func(tx *sql.Tx) error {
		var current User
		var rolesJSON string
		var createdAt, updatedAt int64
		err := tx.QueryRow(
			`SELECT id, username, name, roles, password_hash, token_version, failed_login_count, locked_until, enabled, created_at, updated_at FROM users WHERE id = ?`, id,
		).Scan(&current.ID, &current.Username, &current.Name, &rolesJSON, &current.PasswordHash, &current.TokenVersion, &current.FailedLoginCount, &current.LockedUntil, &current.Enabled, &createdAt, &updatedAt)
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("get user for update: %w", err)
		}
		if err := json.Unmarshal([]byte(rolesJSON), &current.Roles); err != nil {
			return fmt.Errorf("unmarshal roles: %w", err)
		}
		current.CreatedAt = time.Unix(createdAt, 0).UTC()
		current.UpdatedAt = time.Unix(updatedAt, 0).UTC()

		newRoles := current.Roles
		if patch.Roles != nil {
			roles := dedupeKeys(*patch.Roles)
			for _, key := range roles {
				var count int
				if err := tx.QueryRow(`SELECT COUNT(*) FROM roles WHERE key = ?`, key).Scan(&count); err != nil {
					return fmt.Errorf("validate role %s: %w", key, err)
				}
				if count == 0 {
					return ErrInvalidRole
				}
			}
			hadAdmin := slices.Contains(current.Roles, "admin")
			hasAdmin := slices.Contains(roles, "admin")
			if hadAdmin && !hasAdmin {
				if id == actorID {
					return ErrSelfOperation
				}
				other, err := countAdminUsersExcluding(tx, id)
				if err != nil {
					return err
				}
				if other == 0 {
					return ErrLastAdmin
				}
			}
			newRoles = roles
		}

		name := current.Name
		if patch.Name != nil {
			name = strings.TrimSpace(*patch.Name)
		}
		passwordHash := current.PasswordHash
		if patch.PasswordHash != nil {
			passwordHash = *patch.PasswordHash
		}
		rolesBytes, err := json.Marshal(newRoles)
		if err != nil {
			return fmt.Errorf("marshal roles: %w", err)
		}
		nextUpdatedAt := now.Unix()
		if nextUpdatedAt <= current.UpdatedAt.Unix() {
			nextUpdatedAt = current.UpdatedAt.Unix() + 1
		}
		// W4 P0-3: a password change bumps the user's token_version so every
		// already-issued access token (which carries the older version) is
		// rejected by the auth middleware immediately. Atomic with the same
		// transaction that persists the new hash and revokes refresh tokens.
		nextTokenVersion := current.TokenVersion
		if patch.PasswordHash != nil {
			nextTokenVersion = current.TokenVersion + 1
		}
		if _, err := tx.Exec(
			`UPDATE users SET name = ?, roles = ?, password_hash = ?, token_version = ?, updated_at = ? WHERE id = ?`,
			name, string(rolesBytes), passwordHash, nextTokenVersion, nextUpdatedAt, id,
		); err != nil {
			return fmt.Errorf("update user: %w", err)
		}
		if patch.PasswordHash != nil {
			if _, err := tx.Exec(
				`UPDATE refresh_tokens SET revoked_at = COALESCE(revoked_at, ?) WHERE user_id = ?`, now.Unix(), id,
			); err != nil {
				return fmt.Errorf("revoke refresh tokens after password change: %w", err)
			}
		}
		if _, err := tx.Exec(`DELETE FROM user_roles WHERE user_id = ?`, id); err != nil {
			return fmt.Errorf("clear user roles: %w", err)
		}
		for _, key := range newRoles {
			if _, err := tx.Exec(
				`INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)`, id, "role-"+key,
			); err != nil {
				return fmt.Errorf("relink user %s role %s: %w", id, key, err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.GetUser(id)
}

// DeleteUser removes an account and its refresh tokens atomically.
func (r *Repository) DeleteUser(id, actorID string) error {
	return r.withTx("delete managed user", func(tx *sql.Tx) error {
		var exists int
		if err := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`, id).Scan(&exists); err != nil {
			return fmt.Errorf("check delete user: %w", err)
		}
		if exists == 0 {
			return ErrNotFound
		}
		if id == actorID {
			return ErrSelfOperation
		}
		var isAdmin int
		if err := tx.QueryRow(
			`SELECT EXISTS(
				SELECT 1 FROM user_roles ur
				JOIN roles r ON r.id = ur.role_id
				WHERE ur.user_id = ? AND r.key = 'admin'
			)`, id,
		).Scan(&isAdmin); err != nil {
			return fmt.Errorf("check delete user admin role: %w", err)
		}
		if isAdmin == 1 {
			other, err := countAdminUsersExcluding(tx, id)
			if err != nil {
				return err
			}
			if other == 0 {
				return ErrLastAdmin
			}
		}
		if _, err := tx.Exec(`DELETE FROM refresh_tokens WHERE user_id = ?`, id); err != nil {
			return fmt.Errorf("revoke user tokens: %w", err)
		}
		result, err := tx.Exec(`DELETE FROM users WHERE id = ?`, id)
		if err != nil {
			return fmt.Errorf("delete user: %w", err)
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("read delete user affected rows: %w", err)
		}
		if affected != 1 {
			return ErrNotFound
		}
		return nil
	})
}

// DeleteUsersBatch removes accounts and their refresh tokens in one
// transaction (ADR-0022 D5d whole-batch semantics, D-001 P0): every target runs
// the same existence/self guards first, plus a BATCH-LEVEL last-admin guard, so
// any failure rolls the whole batch back and nothing is partially committed.
func (r *Repository) DeleteUsersBatch(ids []string, actorID string) (int, error) {
	keys := dedupeKeys(ids)
	if len(keys) == 0 {
		return 0, nil
	}
	deleted := 0
	err := r.withTx("delete managed users batch", func(tx *sql.Tx) error {
		adminsInBatch := 0
		for _, id := range keys {
			var exists int
			if err := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`, id).Scan(&exists); err != nil {
				return fmt.Errorf("check delete user %s: %w", id, err)
			}
			if exists == 0 {
				return ErrNotFound
			}
			if id == actorID {
				return ErrSelfOperation
			}
			var isAdmin int
			if err := tx.QueryRow(
				`SELECT EXISTS(
					SELECT 1 FROM user_roles ur
					JOIN roles r ON r.id = ur.role_id
					WHERE ur.user_id = ? AND r.key = 'admin'
				)`, id,
			).Scan(&isAdmin); err != nil {
				return fmt.Errorf("check delete user %s admin role: %w", id, err)
			}
			if isAdmin == 1 {
				adminsInBatch++
			}
		}
		// Batch-level last-admin guard (A-002 F-001): per-id counting is
		// unsound for a batch — when every admin is in the selection, each id's
		// "other admins still exist" check passes because the others are also
		// still present, then the whole set is deleted, leaving zero admins.
		// Count the admins NOT in the selection: zero means the batch would
		// remove every admin, which must fail closed.
		if adminsInBatch > 0 {
			other, err := countAdminUsersExcludingBatch(tx, keys)
			if err != nil {
				return err
			}
			if other == 0 {
				return ErrLastAdmin
			}
		}
		for _, id := range keys {
			if _, err := tx.Exec(`DELETE FROM refresh_tokens WHERE user_id = ?`, id); err != nil {
				return fmt.Errorf("revoke user %s tokens: %w", id, err)
			}
			if _, err := tx.Exec(`DELETE FROM users WHERE id = ?`, id); err != nil {
				return fmt.Errorf("delete user %s: %w", id, err)
			}
		}
		deleted = len(keys)
		return nil
	})
	return deleted, err
}

// countAdminUsersExcludingBatch counts admins whose id is not in the batch
// selection. Used by the batch-level last-admin guard in DeleteUsersBatch.
func countAdminUsersExcludingBatch(tx *sql.Tx, ids []string) (int, error) {
	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	var count int
	if err := tx.QueryRow(
		`SELECT COUNT(*) FROM user_roles ur JOIN roles r ON r.id = ur.role_id
		 WHERE r.key = 'admin' AND ur.user_id NOT IN (`+placeholders+`)`, args...,
	).Scan(&count); err != nil {
		return 0, fmt.Errorf("count batch admin users: %w", err)
	}
	return count, nil
}

func countAdminUsersExcluding(tx *sql.Tx, id string) (int, error) {
	var count int
	if err := tx.QueryRow(
		`SELECT COUNT(*) FROM user_roles ur JOIN roles r ON r.id = ur.role_id
		 WHERE r.key = 'admin' AND ur.user_id != ?`, id,
	).Scan(&count); err != nil {
		return 0, fmt.Errorf("count admin users: %w", err)
	}
	return count, nil
}


// scanUserListRow scans the ListUsers projection: the 11 user columns plus the
// cross-module MFA flag (S-10 · GOAL-017 D-002 §4, user_mfa table contract
// from migration 0029). One Scan call — sql.Rows requires dest count ==
// column count.
func scanUserListRow(row interface{ Scan(...any) error }) (*User, error) {
	var user User
	var roles string
	var createdAt, updatedAt int64
	var mfaEnabled int
	err := row.Scan(&user.ID, &user.Username, &user.Name, &roles, &user.PasswordHash,
		&user.TokenVersion, &user.FailedLoginCount, &user.LockedUntil, &user.Enabled,
		&createdAt, &updatedAt, &mfaEnabled)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan user list row: %w", err)
	}
	if err := json.Unmarshal([]byte(roles), &user.Roles); err != nil {
		return nil, fmt.Errorf("unmarshal roles: %w", err)
	}
	user.MFAEnabled = mfaEnabled != 0
	user.CreatedAt = time.Unix(createdAt, 0).UTC()
	user.UpdatedAt = time.Unix(updatedAt, 0).UTC()
	return &user, nil
}

func usersWhere(query string) (string, []any) {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return "", nil
	}
	return ` WHERE (instr(lower(username), ?) > 0 OR instr(lower(name), ?) > 0)`, []any{query, query}
}

func usersSortSQL(sort, order string) string {
	column, collate := "username", " COLLATE NOCASE"
	switch sort {
	case "name":
		column, collate = "name", " COLLATE NOCASE"
	case "updatedAt":
		column, collate = "updated_at", ""
	}
	direction := "ASC"
	if order == "desc" {
		direction = "DESC"
	}
	return column + collate + " " + direction
}