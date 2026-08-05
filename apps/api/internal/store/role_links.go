package store

import (
	"database/sql"
	"fmt"
	"regexp"
)

// roleKeyRe matches stable role keys used by current user and role writes.
var roleKeyRe = regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)

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
		 ON CONFLICT(user_id, role_id) DO NOTHING`,
		userID, "role-"+key,
	); err != nil {
		return fmt.Errorf("link user %s role %s: %w", userID, key, err)
	}
	return nil
}
