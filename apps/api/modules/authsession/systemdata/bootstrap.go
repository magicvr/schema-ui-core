package systemdata

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"time"
)

// TxRunner is the small platform boundary needed by auth/session persistence.
// The module owns the SQL; the platform owns transaction lifecycle.
type TxRunner interface {
	Run(context.Context, func(kernel.Tx) error) error
}

// NeedsBootstrap reports whether no user rows exist yet — the bootstrap admin
// is missing and must be (re)created. Used instead of a one-shot fresh-database
// gate so a failed first bootstrap cannot permanently lock the instance with
// zero users and no login path (C4).
func NeedsBootstrap(ctx context.Context, runner TxRunner) (bool, error) {
	var count int64
	err := runner.Run(ctx, func(tx kernel.Tx) error {
		return tx.QueryRow(context.Background(), `SELECT COUNT(*) FROM users`).Scan(&count)
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
	return runner.Run(ctx, func(tx kernel.Tx) error {
		now := time.Now().UTC().Unix()
		if err := ensureSystemRoles(tx, now); err != nil {
			return err
		}

		var id string
		err := tx.QueryRow(context.Background(), `SELECT id FROM users WHERE username = ?`, username).Scan(&id)
		if errors.Is(err, kernel.ErrNoRows) {
			roles, marshalErr := json.Marshal([]string{"admin", "editor"})
			if marshalErr != nil {
				return fmt.Errorf("marshal bootstrap roles: %w", marshalErr)
			}
			if _, err := tx.Exec(context.Background(),
				`INSERT INTO users (id, username, name, roles, password_hash, must_change_password, created_at, updated_at)
				 VALUES (?, ?, ?, ?, ?, 1, ?, ?)`,
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

func ensureSystemRoles(tx kernel.Tx, now int64) error {
	for _, key := range []string{"admin", "editor", "viewer"} {
		if _, err := tx.Exec(context.Background(),
			`INSERT INTO roles (id, key, name, system, created_at, updated_at)
			 VALUES (?, ?, ?, 1, ?, ?)
			 ON CONFLICT(id) DO UPDATE SET system = 1, updated_at = excluded.updated_at
			 WHERE roles.system <> 1`,
			"role-"+key, key, key, now, now,
		); err != nil {
			return fmt.Errorf("ensure system role %s: %w", key, err)
		}
		var storedKey string
		if err := tx.QueryRow(context.Background(), `SELECT key FROM roles WHERE id = ?`, "role-"+key).Scan(&storedKey); err != nil {
			return fmt.Errorf("verify system role %s: %w", key, err)
		}
		if storedKey != key {
			return fmt.Errorf("system role id role-%s has key %q", key, storedKey)
		}
	}
	return nil
}

func linkUserRole(tx kernel.Tx, userID, roleKey string) error {
	if _, err := tx.Exec(context.Background(),
		`INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)
		 ON CONFLICT(user_id, role_id) DO NOTHING`, userID, "role-"+roleKey,
	); err != nil {
		return err
	}
	return nil
}
