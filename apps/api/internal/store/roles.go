// Roles management repository (GOAL-011 S2 · I-011-001 §3): the management-API
// CRUD path over the roles table with the system-role and in-use protections.
// User-created roles are system=0; the seeded admin/editor/viewer are system=1
// and read-only from the API. Grants (role_permissions / role_menu_items) are NOT
// managed here (full IAM is out of scope, I-011-001 §3.4).
package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Roles-domain sentinels surfaced to the handler for mapping (I-011-001 §6).
var (
	ErrRoleTaken   = errors.New("store: role key already taken")
	ErrRoleInUse   = errors.New("store: role is in use by users")
	ErrRoleSystem  = errors.New("store: role is a system role")
	ErrInvalidKey  = errors.New("store: invalid role key")
)

// Role is the persisted roles row. System marks the seed-managed roles
// (admin/editor/viewer) which the API cannot modify or delete.
type Role struct {
	ID        string
	Key       string
	Name      string
	System    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// RoleFilter carries the validated list query parameters.
type RoleFilter struct {
	Q        string
	Sort     string // key | name | updatedAt (handler-validated)
	Order    string // asc | desc
	Page     int
	PageSize int
}

// ListRoles returns the rows matching f.Q (case-insensitive substring across
// key/name), ordered by the validated sort column, plus the total before
// pagination.
func (s *Store) ListRoles(f RoleFilter) ([]Role, int, error) {
	where, args := rolesWhere(f.Q)
	var total int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM roles`+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count roles: %w", err)
	}

	rows, err := s.db.Query(
		`SELECT id, key, name, system, created_at, updated_at FROM roles`+where+
			` ORDER BY `+rolesSortSQL(f.Sort, f.Order)+`, id ASC`+
			` LIMIT ? OFFSET ?`,
		append(args, f.PageSize, (f.Page-1)*f.PageSize)...,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list roles: %w", err)
	}
	defer rows.Close()

	items := make([]Role, 0, f.PageSize)
	for rows.Next() {
		var r Role
		if err := scanRoleRow(rows, &r); err != nil {
			return nil, 0, err
		}
		items = append(items, r)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("list roles rows: %w", err)
	}
	return items, total, nil
}

// GetRole fetches one role by primary key (role-<key>).
func (s *Store) GetRole(id string) (*Role, error) {
	var r Role
	if err := scanRoleRow(s.db.QueryRow(
		`SELECT id, key, name, system, created_at, updated_at FROM roles WHERE id = ?`, id), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// CreateRole inserts a user-created role (system=0) with id role-<key>. A
// malformed key (roleKeyRe) fails closed as ErrInvalidKey; a duplicate key is
// reported deterministically as ErrRoleTaken (single-connection pre-check).
func (s *Store) CreateRole(key, name string, now time.Time) (*Role, error) {
	if !roleKeyRe.MatchString(key) {
		return nil, ErrInvalidKey
	}
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin create role: %w", err)
	}
	defer tx.Rollback()
	var exists int
	if err := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM roles WHERE key = ?)`, key).Scan(&exists); err != nil {
		return nil, fmt.Errorf("check role key: %w", err)
	}
	if exists == 1 {
		return nil, ErrRoleTaken
	}
	nowUnix := now.Unix()
	if _, err := tx.Exec(
		`INSERT INTO roles (id, key, name, system, created_at, updated_at) VALUES (?, ?, ?, 0, ?, ?)`,
		"role-"+key, key, name, nowUnix, nowUnix,
	); err != nil {
		return nil, fmt.Errorf("insert role: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit create role: %w", err)
	}
	return s.GetRole("role-" + key)
}

// UpdateRole renames a role (key immutable). System roles fail closed
// (ErrRoleSystem). Returns ErrNotFound when the role does not exist.
func (s *Store) UpdateRole(id, name string, now time.Time) (*Role, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin update role: %w", err)
	}
	defer tx.Rollback()
	var cur Role
	if err := scanRoleRow(tx.QueryRow(
		`SELECT id, key, name, system, created_at, updated_at FROM roles WHERE id = ?`, id), &cur); err != nil {
		return nil, err
	}
	if cur.System {
		return nil, ErrRoleSystem
	}
	updatedAt := now.Unix()
	if updatedAt <= cur.UpdatedAt.Unix() {
		updatedAt = cur.UpdatedAt.Unix() + 1
	}
	if _, err := tx.Exec(`UPDATE roles SET name = ?, updated_at = ? WHERE id = ?`, name, updatedAt, id); err != nil {
		return nil, fmt.Errorf("update role: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit update role: %w", err)
	}
	return s.GetRole(id)
}

// DeleteRole removes a user-created role. System roles fail closed
// (ErrRoleSystem); a role still assigned to users fails closed (ErrRoleInUse,
// DB ON DELETE RESTRICT is the backstop). Returns ErrNotFound when absent.
func (s *Store) DeleteRole(id string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin delete role: %w", err)
	}
	defer tx.Rollback()
	var cur Role
	if err := scanRoleRow(tx.QueryRow(
		`SELECT id, key, name, system, created_at, updated_at FROM roles WHERE id = ?`, id), &cur); err != nil {
		return err
	}
	if cur.System {
		return ErrRoleSystem
	}
	var used int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM user_roles WHERE role_id = ?`, id).Scan(&used); err != nil {
		return fmt.Errorf("count role users: %w", err)
	}
	if used > 0 {
		return ErrRoleInUse
	}
	if _, err := tx.Exec(`DELETE FROM roles WHERE id = ?`, id); err != nil {
		return fmt.Errorf("delete role: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit delete role: %w", err)
	}
	return nil
}

// scanRoleRow scans one row (sql.Row or sql.Rows) into a Role.
func scanRoleRow(row interface{ Scan(...any) error }, r *Role) error {
	var system int
	var createdAt, updatedAt int64
	err := row.Scan(&r.ID, &r.Key, &r.Name, &system, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("scan role: %w", err)
	}
	r.System = system == 1
	r.CreatedAt = time.Unix(createdAt, 0).UTC()
	r.UpdatedAt = time.Unix(updatedAt, 0).UTC()
	return nil
}

// rolesWhere builds the optional WHERE clause for q (case-insensitive substring
// across key/name).
func rolesWhere(q string) (string, []any) {
	q = strings.ToLower(strings.TrimSpace(q))
	if q == "" {
		return "", nil
	}
	return ` WHERE (instr(lower(key), ?) > 0 OR instr(lower(name), ?) > 0)`,
		[]any{q, q}
}

// rolesSortSQL maps a handler-validated sort field to its column and ORDER BY
// clause (defaults to key).
func rolesSortSQL(sort, order string) string {
	col, collate := "key", " COLLATE NOCASE"
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
