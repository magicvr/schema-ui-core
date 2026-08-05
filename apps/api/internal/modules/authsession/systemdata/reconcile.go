package systemdata

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

type reconcileEntry struct {
	moduleID string
	kind     string
	key      string
	version  int
	checksum string
}

// Reconcile consumes only the finalized authorization/navigation contributions
// supplied by composition. It is idempotent, profile-aware, and never deletes
// rows that are not explicitly tracked as system-managed grants.
func Reconcile(ctx context.Context, runner TxRunner, permissions []kernel.PermissionContribution, navigation []kernel.NavigationContribution) error {
	perms := append([]kernel.PermissionContribution(nil), permissions...)
	navs := append([]kernel.NavigationContribution(nil), navigation...)
	sort.Slice(perms, func(i, j int) bool {
		if perms[i].ModuleID != perms[j].ModuleID {
			return perms[i].ModuleID < perms[j].ModuleID
		}
		return perms[i].Key < perms[j].Key
	})
	sort.Slice(navs, func(i, j int) bool {
		if navs[i].ModuleID != navs[j].ModuleID {
			return navs[i].ModuleID < navs[j].ModuleID
		}
		return navs[i].Key < navs[j].Key
	})
	if err := validateInputs(perms, navs); err != nil {
		return err
	}
	return runner.WithTx(ctx, func(tx *sql.Tx) error {
		now := time.Now().UTC().Unix()
		if err := ensureSystemRoles(tx, now); err != nil {
			return err
		}
		base := reconcileEntry{
			moduleID: "core.auth-session",
			kind:     "base",
			key:      "system_roles",
			version:  SystemDataVersion,
			checksum: checksumValue(struct {
				Roles []string `json:"roles"`
			}{Roles: []string{"admin", "editor", "viewer"}}),
		}
		if err := checkLedger(tx, base); err != nil {
			return err
		}
		if err := writeLedger(tx, base, now); err != nil {
			return err
		}

		for _, permission := range perms {
			entry := reconcileEntry{
				moduleID: permission.ModuleID,
				kind:     "authorization",
				key:      permission.Key,
				version:  permission.SystemDataVersion,
				checksum: permissionChecksum(permission),
			}
			if err := checkLedger(tx, entry); err != nil {
				return err
			}
			id := permissionID(permission.Permission)
			if err := ensurePermission(tx, id, permission); err != nil {
				return err
			}
			roles, _ := rolesForPolicy(permission.PolicyID)
			if err := syncGrant(tx, entry, id, roles, true); err != nil {
				return err
			}
			if err := writeLedger(tx, entry, now); err != nil {
				return err
			}
		}

		for _, navigation := range navs {
			entry := reconcileEntry{
				moduleID: navigation.ModuleID,
				kind:     "navigation",
				key:      navigation.Key,
				version:  navigation.SystemDataVersion,
				checksum: navigationChecksum(navigation),
			}
			if err := checkLedger(tx, entry); err != nil {
				return err
			}
			id := navigationID(navigation.NodeID)
			if err := ensureNavigation(tx, id, navigation); err != nil {
				return err
			}
			roles, _ := rolesForPolicy(navigation.Visibility)
			if err := syncGrant(tx, entry, id, roles, false); err != nil {
				return err
			}
			if err := writeLedger(tx, entry, now); err != nil {
				return err
			}
		}
		return nil
	})
}

func validateInputs(perms []kernel.PermissionContribution, navs []kernel.NavigationContribution) error {
	seen := map[string]bool{}
	for _, p := range perms {
		if p.ModuleID == "" || p.Key == "" || p.SystemDataVersion <= 0 {
			return fmt.Errorf("system-data permission %q has incomplete identity/version", p.Key)
		}
		if _, ok := rolesForPolicy(p.PolicyID); !ok {
			return fmt.Errorf("system-data permission %q has unknown policy %q", p.Key, p.PolicyID)
		}
		id := "authorization\x00" + p.ModuleID + "\x00" + p.Key
		if seen[id] {
			return fmt.Errorf("duplicate system-data permission %s/%s", p.ModuleID, p.Key)
		}
		seen[id] = true
	}
	for _, n := range navs {
		if n.ModuleID == "" || n.Key == "" || n.PageID == "" || n.SystemDataVersion <= 0 {
			return fmt.Errorf("system-data navigation %q has incomplete identity/version", n.Key)
		}
		if _, ok := rolesForPolicy(n.Visibility); !ok {
			return fmt.Errorf("system-data navigation %q has unknown visibility policy %q", n.Key, n.Visibility)
		}
		id := "navigation\x00" + n.ModuleID + "\x00" + n.Key
		if seen[id] {
			return fmt.Errorf("duplicate system-data navigation %s/%s", n.ModuleID, n.Key)
		}
		seen[id] = true
	}
	return nil
}

func checkLedger(tx *sql.Tx, entry reconcileEntry) error {
	var version int
	var checksum string
	err := tx.QueryRow(
		`SELECT version, checksum FROM system_data_reconcile WHERE module_id = ? AND kind = ? AND contribution_key = ?`,
		entry.moduleID, entry.kind, entry.key,
	).Scan(&version, &checksum)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read system-data ledger %s/%s: %w", entry.kind, entry.key, err)
	}
	if version > entry.version {
		return fmt.Errorf("system-data %s/%s version %d is newer than code %d", entry.kind, entry.key, version, entry.version)
	}
	if version == entry.version && checksum != entry.checksum {
		return fmt.Errorf("system-data %s/%s checksum drift (ledger %s, code %s)", entry.kind, entry.key, checksum, entry.checksum)
	}
	return nil
}

func writeLedger(tx *sql.Tx, entry reconcileEntry, now int64) error {
	if _, err := tx.Exec(
		`INSERT INTO system_data_reconcile (module_id, kind, contribution_key, version, checksum, applied_at)
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON CONFLICT(module_id, kind, contribution_key) DO UPDATE SET version = excluded.version, checksum = excluded.checksum, applied_at = excluded.applied_at
		 WHERE system_data_reconcile.version <> excluded.version OR system_data_reconcile.checksum <> excluded.checksum`,
		entry.moduleID, entry.kind, entry.key, entry.version, entry.checksum, now,
	); err != nil {
		return fmt.Errorf("write system-data ledger %s/%s: %w", entry.kind, entry.key, err)
	}
	return nil
}

func ensurePermission(tx *sql.Tx, id string, p kernel.PermissionContribution) error {
	description := strings.TrimSpace(p.Resource + " " + p.Action + " gate")
	if _, err := tx.Exec(
		`INSERT INTO permissions (id, key, description, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT(id) DO NOTHING`, id, p.Permission, description, time.Now().UTC().Unix(), time.Now().UTC().Unix(),
	); err != nil {
		return fmt.Errorf("ensure permission %s: %w", p.Permission, err)
	}
	var storedID, storedKey string
	if err := tx.QueryRow(`SELECT id, key FROM permissions WHERE key = ?`, p.Permission).Scan(&storedID, &storedKey); err != nil {
		return fmt.Errorf("verify permission %s: %w", p.Permission, err)
	}
	if storedID != id || storedKey != p.Permission {
		return fmt.Errorf("permission %s identity conflict (id %s, want %s)", p.Permission, storedID, id)
	}
	return nil
}

func ensureNavigation(tx *sql.Tx, id string, n kernel.NavigationContribution) error {
	if _, err := tx.Exec(
		`INSERT INTO menu_items (id, page_ref, feature_key, sort_order, enabled, created_at, updated_at)
		 VALUES (?, ?, ?, ?, 1, ?, ?)
		 ON CONFLICT(id) DO NOTHING`, id, n.PageID, n.NodeID, n.Order, time.Now().UTC().Unix(), time.Now().UTC().Unix(),
	); err != nil {
		return fmt.Errorf("ensure navigation %s: %w", n.NodeID, err)
	}
	var storedID, pageRef, featureKey string
	if err := tx.QueryRow(`SELECT id, page_ref, feature_key FROM menu_items WHERE feature_key = ?`, n.NodeID).Scan(&storedID, &pageRef, &featureKey); err != nil {
		return fmt.Errorf("verify navigation %s: %w", n.NodeID, err)
	}
	if storedID != id || pageRef != n.PageID || featureKey != n.NodeID {
		return fmt.Errorf("navigation %s identity conflict (id=%s page=%s feature=%s)", n.NodeID, storedID, pageRef, featureKey)
	}
	return nil
}

func syncGrant(tx *sql.Tx, entry reconcileEntry, targetID string, roles []string, permission bool) error {
	rows, err := tx.Query(
		`SELECT role_key, target_id FROM system_data_grants WHERE module_id = ? AND kind = ? AND contribution_key = ?`,
		entry.moduleID, entry.kind, entry.key,
	)
	if err != nil {
		return fmt.Errorf("read managed grants %s/%s: %w", entry.kind, entry.key, err)
	}
	type grant struct{ roleKey, targetID string }
	var previous []grant
	for rows.Next() {
		var g grant
		if err := rows.Scan(&g.roleKey, &g.targetID); err != nil {
			rows.Close()
			return err
		}
		previous = append(previous, g)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("iterate managed grants %s/%s: %w", entry.kind, entry.key, err)
	}
	if err := rows.Close(); err != nil {
		return err
	}
	desired := make(map[string]bool, len(roles))
	for _, role := range roles {
		desired[role] = true
	}
	for _, old := range previous {
		if desired[old.roleKey] {
			continue
		}
		roleID := "role-" + old.roleKey
		if permission {
			_, err = tx.Exec(`DELETE FROM role_permissions WHERE role_id = ? AND permission_id = ?`, roleID, old.targetID)
		} else {
			_, err = tx.Exec(`DELETE FROM role_menu_items WHERE role_id = ? AND menu_item_id = ?`, roleID, old.targetID)
		}
		if err != nil {
			return fmt.Errorf("remove managed grant %s/%s: %w", entry.kind, entry.key, err)
		}
	}
	if _, err := tx.Exec(`DELETE FROM system_data_grants WHERE module_id = ? AND kind = ? AND contribution_key = ?`, entry.moduleID, entry.kind, entry.key); err != nil {
		return fmt.Errorf("reset managed grants %s/%s: %w", entry.kind, entry.key, err)
	}
	for _, role := range roles {
		roleID := "role-" + role
		if permission {
			if _, err := tx.Exec(`INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?) ON CONFLICT(role_id, permission_id) DO NOTHING`, roleID, targetID); err != nil {
				return fmt.Errorf("grant permission %s to %s: %w", targetID, role, err)
			}
		} else if _, err := tx.Exec(`INSERT INTO role_menu_items (role_id, menu_item_id) VALUES (?, ?) ON CONFLICT(role_id, menu_item_id) DO NOTHING`, roleID, targetID); err != nil {
			return fmt.Errorf("grant menu %s to %s: %w", targetID, role, err)
		}
		if _, err := tx.Exec(
			`INSERT INTO system_data_grants (module_id, kind, contribution_key, role_key, target_id) VALUES (?, ?, ?, ?, ?)`,
			entry.moduleID, entry.kind, entry.key, role, targetID,
		); err != nil {
			return fmt.Errorf("record managed grant %s/%s: %w", entry.kind, entry.key, err)
		}
	}
	return nil
}

func permissionID(key string) string {
	return "perm-" + strings.NewReplacer(".", "-", "_", "-").Replace(key)
}

func navigationID(nodeID string) string {
	return strings.NewReplacer("_", "-").Replace(nodeID)
}

func permissionChecksum(p kernel.PermissionContribution) string {
	return checksumValue(struct {
		ModuleID          string `json:"module_id"`
		Key               string `json:"key"`
		Version           int    `json:"version"`
		Permission        string `json:"permission"`
		Resource          string `json:"resource"`
		Action            string `json:"action"`
		PolicyID          string `json:"policy_id"`
		SecretSensitivity string `json:"secret_sensitivity"`
	}{p.ModuleID, p.Key, p.SystemDataVersion, p.Permission, p.Resource, p.Action, p.PolicyID, p.SecretSensitivity})
}

func navigationChecksum(n kernel.NavigationContribution) string {
	return checksumValue(struct {
		ModuleID   string `json:"module_id"`
		Key        string `json:"key"`
		Version    int    `json:"version"`
		NodeID     string `json:"node_id"`
		PageID     string `json:"page_id"`
		Parent     string `json:"parent"`
		Order      int    `json:"order"`
		Label      string `json:"label"`
		Visibility string `json:"visibility"`
		Permission string `json:"permission"`
	}{n.ModuleID, n.Key, n.SystemDataVersion, n.NodeID, n.PageID, n.Parent, n.Order, n.Label, n.Visibility, n.Permission})
}

func checksumValue(value any) string {
	data, _ := json.Marshal(value)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
