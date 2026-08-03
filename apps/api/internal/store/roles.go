// Roles management repository (GOAL-011 S2 · I-011-001 §3): the management-API
// CRUD path over the roles table with the system-role and in-use protections.
// User-created roles are system=0; the seeded admin/editor/viewer are system=1
// and read-only from the API. User-created roles manage their permission/menu
// grants atomically as the bounded IAM extension frozen in I-011-004 §4.
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
	ErrRoleTaken         = errors.New("store: role key already taken")
	ErrRoleInUse         = errors.New("store: role is in use by users")
	ErrRoleSystem        = errors.New("store: role is a system role")
	ErrInvalidKey        = errors.New("store: invalid role key")
	ErrInvalidPermission = errors.New("store: invalid permission reference")
	ErrInvalidMenuItem   = errors.New("store: invalid menu item reference")
)

// Role is the persisted roles row. System marks the seed-managed roles
// (admin/editor/viewer) which the API cannot modify or delete.
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

// RolePatch applies partial updates to a user-created role. Nil grant pointers
// preserve the corresponding relation; non-nil values replace the whole set.
type RolePatch struct {
	Name        *string
	Permissions *[]string
	MenuItems   *[]string
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
	items := make([]Role, 0, f.PageSize)
	for rows.Next() {
		var r Role
		if err := scanRoleRow(rows, &r); err != nil {
			return nil, 0, err
		}
		items = append(items, r)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, 0, fmt.Errorf("list roles rows: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, 0, fmt.Errorf("close roles rows: %w", err)
	}
	for i := range items {
		if err := s.hydrateRole(&items[i]); err != nil {
			return nil, 0, err
		}
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
	if err := s.hydrateRole(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

// CreateRole inserts a user-created role (system=0) with id role-<key>. A
// malformed key (roleKeyRe) fails closed as ErrInvalidKey; a duplicate key is
// reported deterministically as ErrRoleTaken (single-connection pre-check).
func (s *Store) CreateRole(key, name string, now time.Time) (*Role, error) {
	return s.CreateRoleWithGrants(key, name, nil, nil, now)
}

// CreateRoleWithGrants creates a user role and its permission/menu grants in
// one transaction. Unknown catalog references fail closed without a partial
// role row.
func (s *Store) CreateRoleWithGrants(key, name string, permissions, menuItems []string, now time.Time) (*Role, error) {
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
	if err := replaceRolePermissions(tx, "role-"+key, permissions); err != nil {
		return nil, err
	}
	if err := replaceRoleMenuItems(tx, "role-"+key, menuItems); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit create role: %w", err)
	}
	return s.GetRole("role-" + key)
}

// UpdateRole renames a role (key immutable). System roles fail closed
// (ErrRoleSystem). Returns ErrNotFound when the role does not exist.
func (s *Store) UpdateRole(id, name string, now time.Time) (*Role, error) {
	return s.UpdateRoleWithGrants(id, RolePatch{Name: &name}, now)
}

// UpdateRoleWithGrants updates a custom role and optionally replaces either
// grant set atomically. System roles remain immutable.
func (s *Store) UpdateRoleWithGrants(id string, patch RolePatch, now time.Time) (*Role, error) {
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
	name := cur.Name
	if patch.Name != nil {
		name = strings.TrimSpace(*patch.Name)
	}
	updatedAt := now.Unix()
	if updatedAt <= cur.UpdatedAt.Unix() {
		updatedAt = cur.UpdatedAt.Unix() + 1
	}
	if _, err := tx.Exec(`UPDATE roles SET name = ?, updated_at = ? WHERE id = ?`, name, updatedAt, id); err != nil {
		return nil, fmt.Errorf("update role: %w", err)
	}
	if patch.Permissions != nil {
		if err := replaceRolePermissions(tx, id, *patch.Permissions); err != nil {
			return nil, err
		}
	}
	if patch.MenuItems != nil {
		if err := replaceRoleMenuItems(tx, id, *patch.MenuItems); err != nil {
			return nil, err
		}
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

// PermissionsForRoles returns the union of effective permission keys for role
// keys. Every role must already exist; callers use this to enforce delegated
// assignment as a subset of the actor's current permissions.
func (s *Store) PermissionsForRoles(roleKeys []string) ([]string, error) {
	keys := dedupeKeys(roleKeys)
	if len(keys) == 0 {
		return []string{}, nil
	}
	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(keys)), ",")
	args := make([]any, len(keys))
	for i, key := range keys {
		args[i] = key
	}
	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM roles WHERE key IN (`+placeholders+`)`, args...).Scan(&count); err != nil {
		return nil, fmt.Errorf("count assignment roles: %w", err)
	}
	if count != len(keys) {
		return nil, ErrInvalidRole
	}
	rows, err := s.db.Query(
		`SELECT DISTINCT p.key
		 FROM roles r
		 JOIN role_permissions rp ON rp.role_id = r.id
		 JOIN permissions p ON p.id = rp.permission_id
		 WHERE r.key IN (`+placeholders+`) ORDER BY p.key`, args...,
	)
	if err != nil {
		return nil, fmt.Errorf("query assignment role permissions: %w", err)
	}
	defer rows.Close()
	permissions := []string{}
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, fmt.Errorf("scan assignment role permission: %w", err)
		}
		permissions = append(permissions, key)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("query assignment role permissions: %w", err)
	}
	return permissions, nil
}

// ValidatePermissionKeys and ValidateMenuItemIDs check catalog references
// without mutating grants. Handlers use them to distinguish invalid input from
// a valid-but-not-delegable permission set.
func (s *Store) ValidatePermissionKeys(keys []string) error {
	for _, key := range dedupeKeys(keys) {
		var exists int
		if err := s.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM permissions WHERE key = ?)`, key).Scan(&exists); err != nil {
			return fmt.Errorf("validate permission %s: %w", key, err)
		}
		if exists == 0 {
			return ErrInvalidPermission
		}
	}
	return nil
}

func (s *Store) ValidateMenuItemIDs(ids []string) error {
	for _, id := range dedupeKeys(ids) {
		var exists int
		if err := s.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM menu_items WHERE id = ?)`, id).Scan(&exists); err != nil {
			return fmt.Errorf("validate menu item %s: %w", id, err)
		}
		if exists == 0 {
			return ErrInvalidMenuItem
		}
	}
	return nil
}

func (s *Store) hydrateRole(r *Role) error {
	r.Permissions = []string{}
	rows, err := s.db.Query(
		`SELECT p.key FROM role_permissions rp
		 JOIN permissions p ON p.id = rp.permission_id
		 WHERE rp.role_id = ? ORDER BY p.key`, r.ID,
	)
	if err != nil {
		return fmt.Errorf("query role permissions: %w", err)
	}
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			rows.Close()
			return fmt.Errorf("scan role permission: %w", err)
		}
		r.Permissions = append(r.Permissions, key)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("query role permissions rows: %w", err)
	}
	if err := rows.Close(); err != nil {
		return fmt.Errorf("close role permissions rows: %w", err)
	}

	r.MenuItems = []string{}
	rows, err = s.db.Query(
		`SELECT m.id FROM role_menu_items rmi
		 JOIN menu_items m ON m.id = rmi.menu_item_id
		 WHERE rmi.role_id = ? ORDER BY m.id`, r.ID,
	)
	if err != nil {
		return fmt.Errorf("query role menu items: %w", err)
	}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return fmt.Errorf("scan role menu item: %w", err)
		}
		r.MenuItems = append(r.MenuItems, id)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("query role menu item rows: %w", err)
	}
	if err := rows.Close(); err != nil {
		return fmt.Errorf("close role menu item rows: %w", err)
	}
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM user_roles WHERE role_id = ?`, r.ID).Scan(&r.AssignedUsers); err != nil {
		return fmt.Errorf("count assigned role users: %w", err)
	}
	return nil
}

func replaceRolePermissions(tx *sql.Tx, roleID string, keys []string) error {
	keys = dedupeKeys(keys)
	ids := make([]string, 0, len(keys))
	for _, key := range keys {
		var id string
		err := tx.QueryRow(`SELECT id FROM permissions WHERE key = ?`, key).Scan(&id)
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidPermission
		}
		if err != nil {
			return fmt.Errorf("validate permission %s: %w", key, err)
		}
		ids = append(ids, id)
	}
	if _, err := tx.Exec(`DELETE FROM role_permissions WHERE role_id = ?`, roleID); err != nil {
		return fmt.Errorf("clear role permissions: %w", err)
	}
	for _, id := range ids {
		if _, err := tx.Exec(`INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)`, roleID, id); err != nil {
			return fmt.Errorf("grant role permission %s: %w", id, err)
		}
	}
	return nil
}

func replaceRoleMenuItems(tx *sql.Tx, roleID string, ids []string) error {
	ids = dedupeKeys(ids)
	for _, id := range ids {
		var exists int
		if err := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM menu_items WHERE id = ?)`, id).Scan(&exists); err != nil {
			return fmt.Errorf("validate menu item %s: %w", id, err)
		}
		if exists == 0 {
			return ErrInvalidMenuItem
		}
	}
	if _, err := tx.Exec(`DELETE FROM role_menu_items WHERE role_id = ?`, roleID); err != nil {
		return fmt.Errorf("clear role menu items: %w", err)
	}
	for _, id := range ids {
		if _, err := tx.Exec(`INSERT INTO role_menu_items (role_id, menu_item_id) VALUES (?, ?)`, roleID, id); err != nil {
			return fmt.Errorf("grant role menu item %s: %w", id, err)
		}
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
