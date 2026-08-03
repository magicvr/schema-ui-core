// Users management repository (GOAL-011 S2 · I-011-001 §2): the management-API
// CRUD path over the persisted users table with the account-domain invariants —
// username uniqueness, sensitive-field isolation (password_hash never read into
// the API row), role assignment double-write (legacy roles JSON + user_roles)
// WITHOUT implicit role creation, self/last-admin protection, and atomic
// refresh-token cleanup on delete.
package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

// Account-domain sentinels surfaced to the handler for mapping to
// resource-specific error codes (I-011-001 §6).
var (
	ErrUsernameTaken = errors.New("store: username already taken")
	ErrLastAdmin     = errors.New("store: cannot remove the last admin")
	ErrSelfOperation = errors.New("store: self operation is not allowed")
	ErrInvalidRole   = errors.New("store: invalid role reference")
)

// UserFilter carries the list query parameters already validated by the handler
// (sort field, order, page, pageSize and search text). Sorting and filtering
// happen in SQL so the default path never loads the whole table.
type UserFilter struct {
	Q        string
	Sort     string // username | name | updatedAt (handler-validated)
	Order    string // asc | desc
	Page     int
	PageSize int
}

// UserPatch is the set of editable fields for the management PATCH; a nil
// pointer means "leave unchanged".
type UserPatch struct {
	Name         *string
	Roles        *[]string
	PasswordHash *string
}

// ListUsers returns the rows matching f.Q (case-insensitive substring across
// username/name), ordered by the validated sort column, plus the total number of
// matching rows before pagination. Each row is role-reconciled (userWithRoles
// semantics) so the API never emits a drift between the legacy JSON and the
// normalized relation.
func (s *Store) ListUsers(f UserFilter) ([]User, int, error) {
	where, args := usersWhere(f.Q)
	var total int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM users`+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	rows, err := s.db.Query(
		`SELECT id, username, name, roles, password_hash, created_at, updated_at FROM users`+where+
			` ORDER BY `+usersSortSQL(f.Sort, f.Order)+`, id ASC`+
			` LIMIT ? OFFSET ?`,
		append(args, f.PageSize, (f.Page-1)*f.PageSize)...,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}

	items := make([]User, 0, f.PageSize)
	for rows.Next() {
		var u User
		if err := scanUserRow(rows, &u); err != nil {
			rows.Close()
			return nil, 0, err
		}
		items = append(items, u)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, 0, fmt.Errorf("list users rows: %w", err)
	}
	// Close the list rows BEFORE role reconciliation: the store holds a single
	// connection, so a nested query while rows are open would deadlock (the
	// same constraint dbHasRows documents in migrate.go).
	rows.Close()
	for i := range items {
		reconciled, err := s.reconcileRoles(&items[i])
		if err != nil {
			return nil, 0, err
		}
		items[i] = *reconciled
	}
	return items, total, nil
}

// GetUser fetches one user by primary key with normalized, reconciled roles.
func (s *Store) GetUser(id string) (*User, error) {
	return s.userWithRoles(s.db.QueryRow(
		`SELECT id, username, name, roles, password_hash, created_at, updated_at FROM users WHERE id = ?`, id))
}

// reconcileRoles compares the legacy roles JSON with the normalized user_roles
// relation (set semantics) and, on agreement, makes the normalized (ascending)
// list authoritative.
func (s *Store) reconcileRoles(u *User) (*User, error) {
	norm, err := s.rolesForUser(u.ID)
	if err != nil {
		return nil, err
	}
	if !sameRoleSet(u.Roles, norm) {
		return nil, fmt.Errorf("store: user %s role mismatch: legacy %v normalized %v", u.ID, u.Roles, norm)
	}
	u.Roles = norm
	return u, nil
}

// scanUserRow scans one *sql.Rows row into a User (created/updated as unix UTC).
func scanUserRow(row *sql.Rows, u *User) error {
	var roles string
	var createdAt, updatedAt int64
	if err := row.Scan(&u.ID, &u.Username, &u.Name, &roles, &u.PasswordHash, &createdAt, &updatedAt); err != nil {
		return fmt.Errorf("scan user: %w", err)
	}
	if err := json.Unmarshal([]byte(roles), &u.Roles); err != nil {
		return fmt.Errorf("unmarshal roles: %w", err)
	}
	u.CreatedAt = time.Unix(createdAt, 0).UTC()
	u.UpdatedAt = time.Unix(updatedAt, 0).UTC()
	return nil
}

// CreateUserManagement inserts a new user via the management API (I-011-001
// §2/§2.3). Unlike the seed/internal CreateUser, it validates that every role
// key already exists (NO implicit role creation — the API never calls ensureRole)
// and reports a duplicate username as ErrUsernameTaken. The double-write
// (legacy roles JSON + user_roles) commits atomically.
func (s *Store) CreateUserManagement(u User) (*User, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin create user: %w", err)
	}
	defer tx.Rollback()

	roles := dedupeKeys(u.Roles)
	for _, key := range roles {
		var n int
		if err := tx.QueryRow(`SELECT COUNT(*) FROM roles WHERE key = ?`, key).Scan(&n); err != nil {
			return nil, fmt.Errorf("validate role %s: %w", key, err)
		}
		if n == 0 {
			return nil, ErrInvalidRole
		}
	}
	var exists int
	if err := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)`, u.Username).Scan(&exists); err != nil {
		return nil, fmt.Errorf("check username: %w", err)
	}
	if exists == 1 {
		return nil, ErrUsernameTaken
	}

	rolesJSON, err := json.Marshal(roles)
	if err != nil {
		return nil, fmt.Errorf("marshal roles: %w", err)
	}
	if _, err := tx.Exec(
		`INSERT INTO users (id, username, name, roles, password_hash, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		u.ID, u.Username, u.Name, string(rolesJSON), u.PasswordHash, u.CreatedAt.Unix(), u.UpdatedAt.Unix(),
	); err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}
	// Link roles WITHOUT ensureRole: existence was validated above, so the
	// relation insert can never dangle (I-011-001 §2.3 F-004).
	for _, key := range roles {
		if _, err := tx.Exec(
			`INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)`,
			u.ID, "role-"+key,
		); err != nil {
			return nil, fmt.Errorf("link user %s role %s: %w", u.ID, key, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit create user: %w", err)
	}
	return s.GetUser(u.ID)
}

// UpdateUser applies a management PATCH (name / roles / password) to an existing
// user and refreshes updated_at. Role changes enforce the account-domain
// invariants: a user cannot remove their own admin role (ErrSelfOperation) and
// no change may leave zero admin users (ErrLastAdmin). Unknown role keys fail
// closed (ErrInvalidRole, no implicit creation). Returns ErrNotFound when the
// user does not exist.
func (s *Store) UpdateUser(id string, patch UserPatch, actorID string, now time.Time) (*User, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin update user: %w", err)
	}
	defer tx.Rollback()

	var cur User
	var rolesJSON string
	var createdAt, curUpdatedAt int64
	err = tx.QueryRow(
		`SELECT id, username, name, roles, password_hash, created_at, updated_at FROM users WHERE id = ?`, id,
	).Scan(&cur.ID, &cur.Username, &cur.Name, &rolesJSON, &cur.PasswordHash, &createdAt, &curUpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user for update: %w", err)
	}
	if err := json.Unmarshal([]byte(rolesJSON), &cur.Roles); err != nil {
		return nil, fmt.Errorf("unmarshal roles: %w", err)
	}
	cur.CreatedAt = time.Unix(createdAt, 0).UTC()
	cur.UpdatedAt = time.Unix(curUpdatedAt, 0).UTC()

	newRoles := cur.Roles
	if patch.Roles != nil {
		roles := dedupeKeys(*patch.Roles)
		for _, key := range roles {
			var n int
			if err := tx.QueryRow(`SELECT COUNT(*) FROM roles WHERE key = ?`, key).Scan(&n); err != nil {
				return nil, fmt.Errorf("validate role %s: %w", key, err)
			}
			if n == 0 {
				return nil, ErrInvalidRole
			}
		}
		hadAdmin := slices.Contains(cur.Roles, "admin")
		hasAdmin := slices.Contains(roles, "admin")
		if hadAdmin && !hasAdmin {
			if id == actorID {
				return nil, ErrSelfOperation
			}
			other, err := countAdminUsersExcluding(tx, id)
			if err != nil {
				return nil, err
			}
			if other == 0 {
				return nil, ErrLastAdmin
			}
		}
		newRoles = roles
	}

	name := cur.Name
	if patch.Name != nil {
		name = strings.TrimSpace(*patch.Name)
	}
	passwordHash := cur.PasswordHash
	if patch.PasswordHash != nil {
		passwordHash = *patch.PasswordHash
	}

	rolesBytes, err := json.Marshal(newRoles)
	if err != nil {
		return nil, fmt.Errorf("marshal roles: %w", err)
	}
	rolesJSON = string(rolesBytes)
	updatedAt := now.Unix()
	if updatedAt <= cur.UpdatedAt.Unix() {
		updatedAt = cur.UpdatedAt.Unix() + 1 // monotonic clamp (records D-004 parity)
	}
	if _, err := tx.Exec(
		`UPDATE users SET name = ?, roles = ?, password_hash = ?, updated_at = ? WHERE id = ?`,
		name, string(rolesJSON), passwordHash, updatedAt, id,
	); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	if patch.PasswordHash != nil {
		if _, err := tx.Exec(
			`UPDATE refresh_tokens SET revoked_at = COALESCE(revoked_at, ?) WHERE user_id = ?`,
			now.Unix(), id,
		); err != nil {
			return nil, fmt.Errorf("revoke refresh tokens after password change: %w", err)
		}
	}
	// Rewrite the user_roles relation to match the new role set (double-write).
	if _, err := tx.Exec(`DELETE FROM user_roles WHERE user_id = ?`, id); err != nil {
		return nil, fmt.Errorf("clear user roles: %w", err)
	}
	for _, key := range newRoles {
		if _, err := tx.Exec(
			`INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)`,
			id, "role-"+key,
		); err != nil {
			return nil, fmt.Errorf("relink user %s role %s: %w", id, key, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit update user: %w", err)
	}
	return s.GetUser(id)
}

// DeleteUser removes a user and, atomically, their refresh tokens (revoking any
// live session) — user_roles rows cascade (I-011-001 §2.5). Enforces
// self-protection (ErrSelfOperation) and last-admin protection (ErrLastAdmin).
// Returns ErrNotFound when the user does not exist.
func (s *Store) DeleteUser(id, actorID string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin delete user: %w", err)
	}
	defer tx.Rollback()

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
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit delete user: %w", err)
	}
	return nil
}

// countAdminUsersExcluding counts users (other than id) holding the admin role
// via the normalized user_roles relation.
func countAdminUsersExcluding(tx *sql.Tx, id string) (int, error) {
	var n int
	if err := tx.QueryRow(
		`SELECT COUNT(*) FROM user_roles ur JOIN roles r ON r.id = ur.role_id
		 WHERE r.key = 'admin' AND ur.user_id != ?`, id,
	).Scan(&n); err != nil {
		return 0, fmt.Errorf("count admin users: %w", err)
	}
	return n, nil
}

// usersWhere builds the optional WHERE clause for q (case-insensitive substring
// across username/name).
func usersWhere(q string) (string, []any) {
	q = strings.ToLower(strings.TrimSpace(q))
	if q == "" {
		return "", nil
	}
	return ` WHERE (instr(lower(username), ?) > 0 OR instr(lower(name), ?) > 0)`,
		[]any{q, q}
}

// usersSortSQL maps a handler-validated sort field to its column and ORDER BY
// clause. Only whitelisted fields are mapped (defaults to username).
func usersSortSQL(sort, order string) string {
	col, collate := "username", " COLLATE NOCASE"
	switch sort {
	case "name":
		col, collate = "name", " COLLATE NOCASE"
	case "updatedAt":
		col, collate = "updated_at", ""
	}
	dir := "ASC"
	if order == "desc" {
		dir = "DESC"
	}
	return col + collate + " " + dir
}
