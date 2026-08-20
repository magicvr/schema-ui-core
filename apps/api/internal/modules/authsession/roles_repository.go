package authsession

import (
	"context"
	"errors"
	"fmt"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/pagination"
)

// ListRoles returns the filtered RBAC role projection and total count.
func (r *Repository) ListRoles(filter RoleFilter) ([]Role, int, error) {
	var items []Role
	var total int
	err := r.withTx("list roles", func(tx kernel.Tx) error {
		where, args := rolesWhere(filter.Q, filter.System)
		if err := tx.QueryRow(context.Background(), `SELECT COUNT(*) FROM roles`+where, args...).Scan(&total); err != nil {
			return fmt.Errorf("count roles: %w", err)
		}
		rows, err := tx.Query(context.Background(),
			`SELECT id, key, name, system, created_at, updated_at FROM roles`+where+
				` ORDER BY `+rolesSortSQL(filter.Sort, filter.Order)+`, id ASC`+
				` LIMIT ? OFFSET ?`,
			append(args, filter.PageSize, pagination.Offset(filter.Page, filter.PageSize, total))...,
		)
		if err != nil {
			return fmt.Errorf("query roles: %w", err)
		}
		items = make([]Role, 0, filter.PageSize)
		for rows.Next() {
			var role Role
			if err := scanRoleRow(rows, &role); err != nil {
				_ = rows.Close()
				return err
			}
			items = append(items, role)
		}
		if err := rows.Err(); err != nil {
			_ = rows.Close()
			return fmt.Errorf("list roles rows: %w", err)
		}
		if err := rows.Close(); err != nil {
			return fmt.Errorf("close roles rows: %w", err)
		}
		for i := range items {
			if err := hydrateRole(tx, &items[i]); err != nil {
				return err
			}
		}
		return nil
	})
	return items, total, err
}

// GetRole fetches one role by primary key.
func (r *Repository) GetRole(id string) (*Role, error) {
	var role Role
	err := r.withTx("get role", func(tx kernel.Tx) error {
		if err := scanRoleRow(tx.QueryRow(context.Background(),
			`SELECT id, key, name, system, created_at, updated_at FROM roles WHERE id = ?`, id,
		), &role); err != nil {
			return err
		}
		return hydrateRole(tx, &role)
	})
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// CreateRole creates a user-managed role without grants.
func (r *Repository) CreateRole(key, name string, now time.Time) (*Role, error) {
	return r.CreateRoleWithGrants(key, name, nil, nil, now)
}

// CreateRoleWithGrants creates a custom role and its grants atomically.
func (r *Repository) CreateRoleWithGrants(key, name string, permissions, menuItems []string, now time.Time) (*Role, error) {
	if !roleKeyRe.MatchString(key) {
		return nil, ErrInvalidKey
	}
	err := r.withTx("create role", func(tx kernel.Tx) error {
		var exists int
		if err := tx.QueryRow(context.Background(), `SELECT EXISTS(SELECT 1 FROM roles WHERE key = ?)`, key).Scan(&exists); err != nil {
			return fmt.Errorf("check role key: %w", err)
		}
		if exists == 1 {
			return ErrRoleTaken
		}
		nowUnix := now.Unix()
		if _, err := tx.Exec(context.Background(),
			`INSERT INTO roles (id, key, name, system, created_at, updated_at) VALUES (?, ?, ?, 0, ?, ?)`,
			"role-"+key, key, name, nowUnix, nowUnix,
		); err != nil {
			return fmt.Errorf("insert role: %w", err)
		}
		if err := replaceRolePermissions(tx, "role-"+key, permissions); err != nil {
			return err
		}
		return replaceRoleMenuItems(tx, "role-"+key, menuItems)
	})
	if err != nil {
		return nil, err
	}
	return r.GetRole("role-" + key)
}

// UpdateRole renames a custom role.
func (r *Repository) UpdateRole(id, name string, now time.Time) (*Role, error) {
	return r.UpdateRoleWithGrants(id, RolePatch{Name: &name}, now)
}

// UpdateRoleWithGrants updates a custom role and optional grant sets.
func (r *Repository) UpdateRoleWithGrants(id string, patch RolePatch, now time.Time) (*Role, error) {
	err := r.withTx("update role", func(tx kernel.Tx) error {
		var current Role
		if err := scanRoleRow(tx.QueryRow(context.Background(),
			`SELECT id, key, name, system, created_at, updated_at FROM roles WHERE id = ?`, id,
		), &current); err != nil {
			return err
		}
		if current.System {
			return ErrRoleSystem
		}
		name := current.Name
		if patch.Name != nil {
			name = strings.TrimSpace(*patch.Name)
		}
		updatedAt := now.Unix()
		if updatedAt <= current.UpdatedAt.Unix() {
			updatedAt = current.UpdatedAt.Unix() + 1
		}
		if _, err := tx.Exec(context.Background(), `UPDATE roles SET name = ?, updated_at = ? WHERE id = ?`, name, updatedAt, id); err != nil {
			return fmt.Errorf("update role: %w", err)
		}
		if patch.Permissions != nil {
			if err := replaceRolePermissions(tx, id, *patch.Permissions); err != nil {
				return err
			}
		}
		if patch.MenuItems != nil {
			if err := replaceRoleMenuItems(tx, id, *patch.MenuItems); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.GetRole(id)
}

// DeleteRole removes a custom role when it is not assigned.
func (r *Repository) DeleteRole(id string) error {
	return r.withTx("delete role", func(tx kernel.Tx) error {
		var current Role
		if err := scanRoleRow(tx.QueryRow(context.Background(),
			`SELECT id, key, name, system, created_at, updated_at FROM roles WHERE id = ?`, id,
		), &current); err != nil {
			return err
		}
		if current.System {
			return ErrRoleSystem
		}
		var used int
		if err := tx.QueryRow(context.Background(), `SELECT COUNT(*) FROM user_roles WHERE role_id = ?`, id).Scan(&used); err != nil {
			return fmt.Errorf("count role users: %w", err)
		}
		if used > 0 {
			return ErrRoleInUse
		}
		if _, err := tx.Exec(context.Background(), `DELETE FROM roles WHERE id = ?`, id); err != nil {
			return fmt.Errorf("delete role: %w", err)
		}
		return nil
	})
}

// DeleteRolesBatch removes custom roles in one transaction (ADR-0022 D5d
// whole-batch semantics, D-001 P0): every target runs the same
// existence/system/in-use guards first, so any failure rolls the whole batch
// back and nothing is partially committed.
func (r *Repository) DeleteRolesBatch(ids []string) (int, error) {
	keys := dedupeKeys(ids)
	if len(keys) == 0 {
		return 0, nil
	}
	deleted := 0
	err := r.withTx("delete roles batch", func(tx kernel.Tx) error {
		for _, id := range keys {
			var current Role
			if err := scanRoleRow(tx.QueryRow(context.Background(),
				`SELECT id, key, name, system, created_at, updated_at FROM roles WHERE id = ?`, id,
			), &current); err != nil {
				return err
			}
			if current.System {
				return ErrRoleSystem
			}
			var used int
			if err := tx.QueryRow(context.Background(), `SELECT COUNT(*) FROM user_roles WHERE role_id = ?`, id).Scan(&used); err != nil {
				return fmt.Errorf("count role %s users: %w", id, err)
			}
			if used > 0 {
				return ErrRoleInUse
			}
		}
		for _, id := range keys {
			if _, err := tx.Exec(context.Background(), `DELETE FROM roles WHERE id = ?`, id); err != nil {
				return fmt.Errorf("delete role %s: %w", id, err)
			}
		}
		deleted = len(keys)
		return nil
	})
	return deleted, err
}

// PermissionsForRoles returns the union of permissions for existing role keys.
func (r *Repository) PermissionsForRoles(roleKeys []string) ([]string, error) {
	keys := dedupeKeys(roleKeys)
	if len(keys) == 0 {
		return []string{}, nil
	}
	permissions := []string{}
	err := r.withTx("permissions for roles", func(tx kernel.Tx) error {
		placeholders := strings.TrimSuffix(strings.Repeat("?,", len(keys)), ",")
		args := make([]any, len(keys))
		for i, key := range keys {
			args[i] = key
		}
		var count int
		if err := tx.QueryRow(context.Background(), `SELECT COUNT(*) FROM roles WHERE key IN (`+placeholders+`)`, args...).Scan(&count); err != nil {
			return fmt.Errorf("count assignment roles: %w", err)
		}
		if count != len(keys) {
			return ErrInvalidRole
		}
		rows, err := tx.Query(context.Background(),
			`SELECT DISTINCT p.key
			 FROM roles r
			 JOIN role_permissions rp ON rp.role_id = r.id
			 JOIN permissions p ON p.id = rp.permission_id
			 WHERE r.key IN (`+placeholders+`) ORDER BY p.key`, args...,
		)
		if err != nil {
			return fmt.Errorf("query assignment role permissions: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var key string
			if err := rows.Scan(&key); err != nil {
				return fmt.Errorf("scan assignment role permission: %w", err)
			}
			permissions = append(permissions, key)
		}
		if err := rows.Err(); err != nil {
			return fmt.Errorf("query assignment role permissions: %w", err)
		}
		return nil
	})
	return permissions, err
}

// ValidatePermissionKeys rejects unknown permission references.
func (r *Repository) ValidatePermissionKeys(keys []string) error {
	return r.withTx("validate permission keys", func(tx kernel.Tx) error {
		for _, key := range dedupeKeys(keys) {
			var exists int
			if err := tx.QueryRow(context.Background(), `SELECT EXISTS(SELECT 1 FROM permissions WHERE key = ?)`, key).Scan(&exists); err != nil {
				return fmt.Errorf("validate permission %s: %w", key, err)
			}
			if exists == 0 {
				return ErrInvalidPermission
			}
		}
		return nil
	})
}

// ValidateMenuItemIDs rejects unknown navigation references.
func (r *Repository) ValidateMenuItemIDs(ids []string) error {
	return r.withTx("validate menu item ids", func(tx kernel.Tx) error {
		for _, id := range dedupeKeys(ids) {
			var exists int
			if err := tx.QueryRow(context.Background(), `SELECT EXISTS(SELECT 1 FROM menu_items WHERE id = ?)`, id).Scan(&exists); err != nil {
				return fmt.Errorf("validate menu item %s: %w", id, err)
			}
			if exists == 0 {
				return ErrInvalidMenuItem
			}
		}
		return nil
	})
}

func hydrateRole(tx kernel.Tx, role *Role) error {
	role.Permissions = []string{}
	rows, err := tx.Query(context.Background(),
		`SELECT p.key FROM role_permissions rp
		 JOIN permissions p ON p.id = rp.permission_id
		 WHERE rp.role_id = ? ORDER BY p.key`, role.ID,
	)
	if err != nil {
		return fmt.Errorf("query role permissions: %w", err)
	}
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			_ = rows.Close()
			return fmt.Errorf("scan role permission: %w", err)
		}
		role.Permissions = append(role.Permissions, key)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return fmt.Errorf("query role permissions rows: %w", err)
	}
	if err := rows.Close(); err != nil {
		return fmt.Errorf("close role permissions rows: %w", err)
	}

	role.MenuItems = []string{}
	rows, err = tx.Query(context.Background(),
		`SELECT m.id FROM role_menu_items rmi
		 JOIN menu_items m ON m.id = rmi.menu_item_id
		 WHERE rmi.role_id = ? ORDER BY m.id`, role.ID,
	)
	if err != nil {
		return fmt.Errorf("query role menu items: %w", err)
	}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			_ = rows.Close()
			return fmt.Errorf("scan role menu item: %w", err)
		}
		role.MenuItems = append(role.MenuItems, id)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return fmt.Errorf("query role menu item rows: %w", err)
	}
	if err := rows.Close(); err != nil {
		return fmt.Errorf("close role menu item rows: %w", err)
	}
	if err := tx.QueryRow(context.Background(), `SELECT COUNT(*) FROM user_roles WHERE role_id = ?`, role.ID).Scan(&role.AssignedUsers); err != nil {
		return fmt.Errorf("count assigned role users: %w", err)
	}
	return nil
}

func replaceRolePermissions(tx kernel.Tx, roleID string, keys []string) error {
	keys = dedupeKeys(keys)
	ids := make([]string, 0, len(keys))
	for _, key := range keys {
		var id string
		err := tx.QueryRow(context.Background(), `SELECT id FROM permissions WHERE key = ?`, key).Scan(&id)
		if errors.Is(err, kernel.ErrNoRows) {
			return ErrInvalidPermission
		}
		if err != nil {
			return fmt.Errorf("validate permission %s: %w", key, err)
		}
		ids = append(ids, id)
	}
	if _, err := tx.Exec(context.Background(), `DELETE FROM role_permissions WHERE role_id = ?`, roleID); err != nil {
		return fmt.Errorf("clear role permissions: %w", err)
	}
	for _, id := range ids {
		if _, err := tx.Exec(context.Background(), `INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)`, roleID, id); err != nil {
			return fmt.Errorf("grant role permission %s: %w", id, err)
		}
	}
	return nil
}

func replaceRoleMenuItems(tx kernel.Tx, roleID string, ids []string) error {
	ids = dedupeKeys(ids)
	for _, id := range ids {
		var exists int
		if err := tx.QueryRow(context.Background(), `SELECT EXISTS(SELECT 1 FROM menu_items WHERE id = ?)`, id).Scan(&exists); err != nil {
			return fmt.Errorf("validate menu item %s: %w", id, err)
		}
		if exists == 0 {
			return ErrInvalidMenuItem
		}
	}
	if _, err := tx.Exec(context.Background(), `DELETE FROM role_menu_items WHERE role_id = ?`, roleID); err != nil {
		return fmt.Errorf("clear role menu items: %w", err)
	}
	for _, id := range ids {
		if _, err := tx.Exec(context.Background(), `INSERT INTO role_menu_items (role_id, menu_item_id) VALUES (?, ?)`, roleID, id); err != nil {
			return fmt.Errorf("grant role menu item %s: %w", id, err)
		}
	}
	return nil
}

func scanRoleRow(row interface{ Scan(...any) error }, role *Role) error {
	var system int
	var createdAt, updatedAt int64
	err := row.Scan(&role.ID, &role.Key, &role.Name, &system, &createdAt, &updatedAt)
	if errors.Is(err, kernel.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("scan role: %w", err)
	}
	role.System = system == 1
	role.CreatedAt = time.Unix(createdAt, 0).UTC()
	role.UpdatedAt = time.Unix(updatedAt, 0).UTC()
	return nil
}

func rolesWhere(query string, system *bool) (string, []any) {
	clauses := []string{}
	args := []any{}
	if q := strings.ToLower(strings.TrimSpace(query)); q != "" {
		clauses = append(clauses, `(instr(lower(key), ?) > 0 OR instr(lower(name), ?) > 0)`)
		args = append(args, q, q)
	}
	if system != nil {
		clauses = append(clauses, `system = ?`)
		args = append(args, boolInt(*system))
	}
	if len(clauses) == 0 {
		return "", nil
	}
	return ` WHERE ` + strings.Join(clauses, ` AND `), args
}

func rolesSortSQL(sort, order string) string {
	// R4 S2 (R1 v1.4 F-002): portable case-insensitive collation via LOWER —
	// sqlite COLLATE NOCASE has no postgres equivalent.
	expr := "LOWER(key)"
	switch sort {
	case "name":
		expr = "LOWER(name)"
	case "updatedAt":
		expr = "updated_at"
	}
	direction := "ASC"
	if order == "desc" {
		direction = "DESC"
	}
	return expr + " " + direction
}

// ListPermissionCatalog returns every registered permission key (W11 · U-02).
// The catalog is the reconciled permissions table — the same source the
// roles resource validates grants against, so the UI options can never offer
// a key the backend would reject.
func (r *Repository) ListPermissionCatalog() ([]PermissionCatalogEntry, error) {
	var out []PermissionCatalogEntry
	err := r.withTx("list permission catalog", func(tx kernel.Tx) error {
		rows, err := tx.Query(context.Background(), `SELECT key, description FROM permissions ORDER BY key`)
		if err != nil {
			return fmt.Errorf("query permission catalog: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var entry PermissionCatalogEntry
			if err := rows.Scan(&entry.Key, &entry.Description); err != nil {
				return fmt.Errorf("scan permission catalog: %w", err)
			}
			out = append(out, entry)
		}
		return rows.Err()
	})
	return out, err
}

// ListMenuItemCatalog returns every enabled navigation node (W11 · U-02).
// menu_items.id is the same value the roles resource validates menu grants
// against ("menu-users" shape); the display label is derived from page_ref.
func (r *Repository) ListMenuItemCatalog() ([]MenuItemCatalogEntry, error) {
	var out []MenuItemCatalogEntry
	err := r.withTx("list menu item catalog", func(tx kernel.Tx) error {
		rows, err := tx.Query(context.Background(),
			`SELECT id, page_ref FROM menu_items WHERE enabled = 1 ORDER BY sort_order, id`,
		)
		if err != nil {
			return fmt.Errorf("query menu item catalog: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var entry MenuItemCatalogEntry
			if err := rows.Scan(&entry.ID, &entry.PageRef); err != nil {
				return fmt.Errorf("scan menu item catalog: %w", err)
			}
			entry.Label = menuItemLabel(entry.PageRef)
			out = append(out, entry)
		}
		return rows.Err()
	})
	return out, err
}

// menuItemLabel derives a display label from a page id ("data-dictionary" →
// "Data dictionary"), matching the admin console's English page titles.
func menuItemLabel(pageRef string) string {
	parts := strings.Split(pageRef, "-")
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, " ")
}
