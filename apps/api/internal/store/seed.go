// Incremental R3 seed (GOAL-006 D-002.6 / I-006-001 §6): ensures the stable
// roles, permissions, menu item and grants exist idempotently. Each entity and
// relation is ensured independently, so an existing user never causes the whole
// relation seed to be skipped, and repeated startups create no duplicates. It
// never overwrites non-seed user fields and never upgrades editor to a write
// role.
package store

import (
	"database/sql"
	"fmt"
	"time"
)

// seedRBAC is called from Open when seedAdmin is enabled. It runs after
// seedAdmin and upgrades the derived roles (system=0, created by 0002/双写) to
// the stable seed (system=1) and links permissions and the representative menu
// grant. It is transactional and idempotent.
func (s *Store) seedRBAC() error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin rbac seed: %w", err)
	}
	defer tx.Rollback()
	now := time.Now().UTC().Unix()

	// Stable system roles: admin/editor/viewer (system=1).
	for _, key := range []string{"admin", "editor", "viewer"} {
		if _, err := tx.Exec(
			`INSERT INTO roles (id, key, name, system, created_at, updated_at)
			 VALUES (?, ?, ?, 1, ?, ?)
			 ON CONFLICT(id) DO UPDATE SET name = excluded.name, system = 1, updated_at = excluded.updated_at`,
			"role-"+key, key, key, now, now,
		); err != nil {
			return fmt.Errorf("seed role %s: %w", key, err)
		}
	}

	// GOAL-011 (I-011-001 §4): users/roles permissions + menus. records
	// permissions/menu were removed by 0006 records_retire (GOAL-011 S3); the
	// seed no longer creates them for fresh installs.
	for _, p := range []struct {
		id, key, desc string
	}{
		{"perm-users-read", "users.read", "users GET gate"},
		{"perm-users-write", "users.write", "users write gate"},
		{"perm-roles-read", "roles.read", "roles GET gate"},
		{"perm-roles-write", "roles.write", "roles write gate"},
		{"perm-roles-assign", "roles.assign", "user role assignment gate"},
	} {
		if _, err := tx.Exec(
			`INSERT INTO permissions (id, key, description, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?)
			 ON CONFLICT(id) DO UPDATE SET description = excluded.description, updated_at = excluded.updated_at`,
			p.id, p.key, p.desc, now, now,
		); err != nil {
			return fmt.Errorf("seed permission %s: %w", p.key, err)
		}
	}
	for _, m := range []struct {
		id, pageRef, featureKey string
	}{
		{"menu-users", "users", "menu_users"},
		{"menu-roles", "roles", "menu_roles"},
	} {
		if _, err := tx.Exec(
			`INSERT INTO menu_items (id, page_ref, feature_key, sort_order, enabled, created_at, updated_at)
			 VALUES (?, ?, ?, 0, 1, ?, ?)
			 ON CONFLICT(id) DO UPDATE SET page_ref = excluded.page_ref, feature_key = excluded.feature_key, enabled = 1, updated_at = excluded.updated_at`,
			m.id, m.pageRef, m.featureKey, now, now,
		); err != nil {
			return fmt.Errorf("seed menu item %s: %w", m.id, err)
		}
	}

	// Grants: admin read+write (+menus), editor read, viewer read (read-only).
	// GOAL-011: admin holds users/roles read+write + both menus; editor/viewer read.
	for _, p := range []string{"perm-users-read", "perm-users-write", "perm-roles-read", "perm-roles-write", "perm-roles-assign"} {
		if err := linkPermission(tx, "admin", p); err != nil {
			return err
		}
	}
	for _, m := range []string{"menu-users", "menu-roles"} {
		if err := linkMenu(tx, "admin", m); err != nil {
			return err
		}
	}
	for _, key := range []string{"editor", "viewer"} {
		if err := linkPermission(tx, key, "perm-users-read"); err != nil {
			return err
		}
		if err := linkPermission(tx, key, "perm-roles-read"); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func linkPermission(tx *sql.Tx, roleKey, permissionID string) error {
	if _, err := tx.Exec(
		`INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)
		 ON CONFLICT(role_id, permission_id) DO NOTHING`,
		"role-"+roleKey, permissionID,
	); err != nil {
		return fmt.Errorf("seed grant role %s permission %s: %w", roleKey, permissionID, err)
	}
	return nil
}

func linkMenu(tx *sql.Tx, roleKey, menuItemID string) error {
	if _, err := tx.Exec(
		`INSERT INTO role_menu_items (role_id, menu_item_id) VALUES (?, ?)
		 ON CONFLICT(role_id, menu_item_id) DO NOTHING`,
		"role-"+roleKey, menuItemID,
	); err != nil {
		return fmt.Errorf("seed grant role %s menu %s: %w", roleKey, menuItemID, err)
	}
	return nil
}
