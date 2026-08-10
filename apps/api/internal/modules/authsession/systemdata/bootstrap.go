package systemdata

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// TxRunner is the small platform boundary needed by auth/session persistence.
// The module owns the SQL; the platform owns transaction lifecycle.
type TxRunner interface {
	WithTx(context.Context, func(*sql.Tx) error) error
}

// NeedsBootstrap reports whether no user rows exist yet — the bootstrap admin
// is missing and must be (re)created. Used instead of a one-shot fresh-database
// gate so a failed first bootstrap cannot permanently lock the instance with
// zero users and no login path (C4).
func NeedsBootstrap(ctx context.Context, runner TxRunner) (bool, error) {
	var count int64
	err := runner.WithTx(ctx, func(tx *sql.Tx) error {
		return tx.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count)
	})
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// Bootstrap creates the initial administrator and base system roles for a
// genuinely fresh database. Existing user fields are never updated.
func Bootstrap(ctx context.Context, runner TxRunner, username, passwordHash string) error {
	if username == "" {
		return errors.New("authsession bootstrap: username is required")
	}
	return runner.WithTx(ctx, func(tx *sql.Tx) error {
		now := time.Now().UTC().Unix()
		if err := ensureSystemRoles(tx, now); err != nil {
			return err
		}

		var id string
		err := tx.QueryRow(`SELECT id FROM users WHERE username = ?`, username).Scan(&id)
		if errors.Is(err, sql.ErrNoRows) {
			roles, marshalErr := json.Marshal([]string{"admin", "editor"})
			if marshalErr != nil {
				return fmt.Errorf("marshal bootstrap roles: %w", marshalErr)
			}
			if _, err := tx.Exec(
				`INSERT INTO users (id, username, name, roles, password_hash, created_at, updated_at)
				 VALUES (?, ?, ?, ?, ?, ?, ?)`,
				"user-admin", username, "Admin", string(roles), passwordHash, now, now,
			); err != nil {
				return fmt.Errorf("insert bootstrap admin: %w", err)
			}
			id = "user-admin"
		} else if err != nil {
			return fmt.Errorf("lookup bootstrap admin: %w", err)
		}

		for _, key := range []string{"admin", "editor"} {
			if err := linkUserRole(tx, id, key); err != nil {
				return fmt.Errorf("link bootstrap role %s: %w", key, err)
			}
		}
		return nil
	})
}

func ensureSystemRoles(tx *sql.Tx, now int64) error {
	for _, key := range []string{"admin", "editor", "viewer"} {
		if _, err := tx.Exec(
			`INSERT INTO roles (id, key, name, system, created_at, updated_at)
			 VALUES (?, ?, ?, 1, ?, ?)
			 ON CONFLICT(id) DO UPDATE SET system = 1, updated_at = excluded.updated_at
			 WHERE roles.system <> 1`,
			"role-"+key, key, key, now, now,
		); err != nil {
			return fmt.Errorf("ensure system role %s: %w", key, err)
		}
		var storedKey string
		if err := tx.QueryRow(`SELECT key FROM roles WHERE id = ?`, "role-"+key).Scan(&storedKey); err != nil {
			return fmt.Errorf("verify system role %s: %w", key, err)
		}
		if storedKey != key {
			return fmt.Errorf("system role id role-%s has key %q", key, storedKey)
		}
	}
	return nil
}

func linkUserRole(tx *sql.Tx, userID, roleKey string) error {
	if _, err := tx.Exec(
		`INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)
		 ON CONFLICT(user_id, role_id) DO NOTHING`, userID, "role-"+roleKey,
	); err != nil {
		return err
	}
	return nil
}
