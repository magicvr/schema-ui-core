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

// EnsureTestAdmin upserts an OPTIONAL test-only bootstrap admin (TEST_ADMIN_
// USERNAME / TEST_ADMIN_PASSWORD). When the password is non-empty it (re)creates
// the user with the given bcrypt hash and resets must_change_password to 0 on
// every boot, so local/CI verification has a stable credential without touching
// the existing "admin" bootstrap user. Empty password = no-op. It never runs on
// its own: the composition root calls it only when the env pair is configured.
func EnsureTestAdmin(ctx context.Context, runner TxRunner, username, passwordHash string) error {
	if username == "" || passwordHash == "" {
		return nil
	}
	return runner.Run(ctx, func(tx kernel.Tx) error {
		now := time.Now().UTC().Unix()
		if err := ensureSystemRoles(tx, now); err != nil {
			return err
		}
		var id string
		err := tx.QueryRow(context.Background(), `SELECT id FROM users WHERE username = ?`, username).Scan(&id)
		switch {
		case errors.Is(err, kernel.ErrNoRows):
			roles, marshalErr := json.Marshal([]string{"admin", "editor"})
			if marshalErr != nil {
				return fmt.Errorf("marshal test-admin roles: %w", marshalErr)
			}
			id = "user-" + username
			if _, execErr := tx.Exec(context.Background(),
				`INSERT INTO users (id, username, name, roles, password_hash, must_change_password, created_at, updated_at)
				 VALUES (?, ?, ?, ?, ?, 0, ?, ?)`,
				id, username, username, string(roles), passwordHash, now, now,
			); execErr != nil {
				return fmt.Errorf("insert test admin: %w", execErr)
			}
		case err != nil:
			return fmt.Errorf("lookup test admin: %w", err)
		default:
			// User exists: reset password to the env value and clear the
			// must-change flag so the configured credential always works.
			if _, execErr := tx.Exec(context.Background(),
				`UPDATE users SET password_hash = ?, must_change_password = 0, updated_at = ? WHERE id = ?`,
				passwordHash, now, id,
			); execErr != nil {
				return fmt.Errorf("update test admin: %w", execErr)
			}
		}

		for _, key := range []string{"admin", "editor"} {
			if err := linkUserRole(tx, id, key); err != nil {
				return fmt.Errorf("link test-admin role %s: %w", key, err)
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
