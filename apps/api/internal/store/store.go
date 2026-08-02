// Package store provides the R2 SQLite persistence for authentication
// (GOAL-005 D-003): a minimal users table (with roles as a JSON string until R3
// normalizes the model) and a refresh_tokens table storing opaque refresh
// tokens as SHA-256 hashes so a DB leak does not expose usable tokens.
package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// User is the persisted identity row backing account.Session.User.
type User struct {
	ID           string
	Username     string
	Name         string
	Roles        []string // R2 minimal; R3 replaces with normalized relations
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// RefreshToken is a stored opaque refresh token (only its hash is kept).
type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

// Sentinels for store lookups and mutations.
var (
	ErrNotFound       = errors.New("store: not found")
	ErrAlreadyRevoked = errors.New("store: refresh token already revoked")
)

// Store wraps the SQLite auth store. Concurrency is guarded by a single
// connection (modernc/sqlite is a single-writer backend).
type Store struct {
	db   *sql.DB
	path string // file path used for pre-upgrade recovery snapshots
}

// Open opens (creating if needed) the SQLite DB at path, applies versioned
// migrations (see migrate.go), and — when seedAdmin is true and the users table
// is empty — seeds the admin user with the given bcrypt password hash. The
// caller is responsible for enforcing a non-empty seed password in production
// (fail-closed on startup).
func Open(path, adminUsername, adminPasswordHash string, seedAdmin bool) (*Store, error) {
	if dir := filepath.Dir(path); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create db dir: %w", err)
		}
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	// modernc/sqlite handles a single writer; one conn avoids lock contention
	// and mirrors the R1 in-process single-writer posture.
	db.SetMaxOpenConns(1)

	s := &Store{db: db, path: path}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	if seedAdmin {
		if err := s.seedAdmin(adminUsername, adminPasswordHash); err != nil {
			db.Close()
			return nil, err
		}
		// Incremental R3 seed: stable roles/permissions/menu/grants (S3). Runs
		// whenever seeding is enabled so existing users never skip relation repair.
		if err := s.seedRBAC(); err != nil {
			db.Close()
			return nil, err
		}
	}
	return s, nil
}

// Close releases the underlying database handle.
func (s *Store) Close() error { return s.db.Close() }

// seedAdmin ensures the bootstrap admin user (with roles admin + editor) and its
// normalized role relations exist. It is idempotent: it never overwrites the
// password or other user fields of an existing admin, and the double-write
// closes the S1 fresh-DB intermediate state where 0002's backfill ran before the
// seed user was inserted (GOAL-006 D-004).
func (s *Store) seedAdmin(username, passwordHash string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin seed admin: %w", err)
	}
	defer tx.Rollback()
	now := time.Now().UTC().Unix()

	var id string
	err = tx.QueryRow(`SELECT id FROM users WHERE username = ?`, username).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		roles, err := json.Marshal([]string{"admin", "editor"})
		if err != nil {
			return fmt.Errorf("marshal seed roles: %w", err)
		}
		if _, err := tx.Exec(
			`INSERT INTO users (id, username, name, roles, password_hash, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`,
			"user-admin", username, "Admin", string(roles), passwordHash, now, now,
		); err != nil {
			return fmt.Errorf("seed admin: %w", err)
		}
		id = "user-admin"
	} else if err != nil {
		return fmt.Errorf("lookup seed admin: %w", err)
	}

	for _, key := range []string{"admin", "editor"} {
		if err := linkUserRole(tx, id, key, now); err != nil {
			return fmt.Errorf("seed admin role %s: %w", key, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit seed admin: %w", err)
	}
	return nil
}

// CreateUser inserts a new user. Since GOAL-006 阶段 B the legacy roles JSON and
// the normalized user_roles relation are written together in one transaction so
// the double-read compare in UserByID/UserByUsername never observes a drift.
// Roles are deduped before the JSON is stored so the two sources stay set-equal.
func (s *Store) CreateUser(u User) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin create user: %w", err)
	}
	defer tx.Rollback()
	now := time.Now().UTC().Unix()

	roles := dedupeKeys(u.Roles)
	rolesJSON, err := json.Marshal(roles)
	if err != nil {
		return fmt.Errorf("marshal roles: %w", err)
	}
	if _, err := tx.Exec(
		`INSERT INTO users (id, username, name, roles, password_hash, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		u.ID, u.Username, u.Name, string(rolesJSON), u.PasswordHash, u.CreatedAt.Unix(), u.UpdatedAt.Unix(),
	); err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	for _, key := range roles {
		if err := linkUserRole(tx, u.ID, key, now); err != nil {
			return fmt.Errorf("create user %s role %s: %w", u.ID, key, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit create user: %w", err)
	}
	return nil
}

// UserByUsername fetches a user by unique username.
func (s *Store) UserByUsername(username string) (*User, error) {
	return s.userWithRoles(s.db.QueryRow(
		`SELECT id, username, name, roles, password_hash, created_at, updated_at
		 FROM users WHERE username = ?`, username))
}

// UserByID fetches a user by primary key.
func (s *Store) UserByID(id string) (*User, error) {
	return s.userWithRoles(s.db.QueryRow(
		`SELECT id, username, name, roles, password_hash, created_at, updated_at
		 FROM users WHERE id = ?`, id))
}

// userWithRoles reads the persisted row and reconciles the two role sources
// (GOAL-006 阶段 B): the legacy roles JSON must be set-equal to the normalized
// user_roles relation, otherwise a diagnostic error is returned instead of a
// silently inconsistent identity. On agreement the normalized roles (ordered by
// role key ascending) are the authoritative read value. The users.roles column
// is retained until an explicit later migration removes it.
//
// Comparison follows the frozen set semantics (I-006-001 §5, A-002 F-004):
// duplicates in the legacy JSON are historical artifacts that 0002's backfill
// dedupes in the relations, so the two sources are compared as sets. A genuine
// set difference still fails closed.
func (s *Store) userWithRoles(row *sql.Row) (*User, error) {
	u, err := s.scanUser(row)
	if err != nil {
		return nil, err
	}
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

// rolesForUser returns a user's normalized roles ordered by role key ascending.
func (s *Store) rolesForUser(userID string) ([]string, error) {
	rows, err := s.db.Query(
		`SELECT r.key FROM user_roles ur JOIN roles r ON r.id = ur.role_id
		 WHERE ur.user_id = ? ORDER BY r.key`, userID)
	if err != nil {
		return nil, fmt.Errorf("query normalized roles: %w", err)
	}
	defer rows.Close()
	var keys []string
	for rows.Next() {
		var k string
		if err := rows.Scan(&k); err != nil {
			return nil, fmt.Errorf("scan normalized role: %w", err)
		}
		keys = append(keys, k)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("query normalized roles: %w", err)
	}
	return keys, nil
}

// sameRoleSet reports whether a and b hold the same set of role keys, ignoring
// order and duplicates (set semantics per I-006-001 §5). A duplicate role key
// in the legacy JSON is treated as the same role as the deduped relation.
func sameRoleSet(a, b []string) bool {
	as := make(map[string]bool, len(a))
	for _, v := range a {
		as[v] = true
	}
	bs := make(map[string]bool, len(b))
	for _, v := range b {
		bs[v] = true
	}
	if len(as) != len(bs) {
		return false
	}
	for v := range as {
		if !bs[v] {
			return false
		}
	}
	return true
}

// dedupeKeys returns keys with duplicates removed, preserving first-occurrence
// order.
func dedupeKeys(keys []string) []string {
	seen := make(map[string]bool, len(keys))
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		if !seen[k] {
			seen[k] = true
			out = append(out, k)
		}
	}
	return out
}

func (s *Store) scanUser(row *sql.Row) (*User, error) {
	var (
		u         User
		roles     string
		createdAt int64
		updatedAt int64
	)
	err := row.Scan(&u.ID, &u.Username, &u.Name, &roles, &u.PasswordHash, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan user: %w", err)
	}
	if err := json.Unmarshal([]byte(roles), &u.Roles); err != nil {
		return nil, fmt.Errorf("unmarshal roles: %w", err)
	}
	u.CreatedAt = time.Unix(createdAt, 0).UTC()
	u.UpdatedAt = time.Unix(updatedAt, 0).UTC()
	return &u, nil
}

// CreateRefreshToken persists a hashed opaque refresh token.
func (s *Store) CreateRefreshToken(rt RefreshToken) error {
	var revokedAt any
	if rt.RevokedAt != nil {
		revokedAt = rt.RevokedAt.Unix()
	}
	_, err := s.db.Exec(
		`INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		rt.ID, rt.UserID, rt.TokenHash, rt.ExpiresAt.Unix(), revokedAt, rt.CreatedAt.Unix(),
	)
	if err != nil {
		return fmt.Errorf("insert refresh token: %w", err)
	}
	return nil
}

// RefreshTokenByHash fetches a refresh token row by its stored hash.
func (s *Store) RefreshTokenByHash(hash string) (*RefreshToken, error) {
	var (
		rt        RefreshToken
		expiresAt int64
		revokedAt *int64
		createdAt int64
	)
	err := s.db.QueryRow(
		`SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
		 FROM refresh_tokens WHERE token_hash = ?`, hash,
	).Scan(&rt.ID, &rt.UserID, &rt.TokenHash, &expiresAt, &revokedAt, &createdAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan refresh token: %w", err)
	}
	rt.ExpiresAt = time.Unix(expiresAt, 0).UTC()
	if revokedAt != nil {
		t := time.Unix(*revokedAt, 0).UTC()
		rt.RevokedAt = &t
	}
	rt.CreatedAt = time.Unix(createdAt, 0).UTC()
	return &rt, nil
}

// RevokeRefreshToken marks a refresh token revoked. It returns ErrNotFound when
// the token does not exist and ErrAlreadyRevoked when it was already revoked.
func (s *Store) RevokeRefreshToken(id string, now time.Time) error {
	var current *int64
	err := s.db.QueryRow(`SELECT revoked_at FROM refresh_tokens WHERE id = ?`, id).Scan(&current)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("select revoked_at: %w", err)
	}
	if current != nil {
		return ErrAlreadyRevoked
	}
	if _, err := s.db.Exec(`UPDATE refresh_tokens SET revoked_at = ? WHERE id = ?`, now.Unix(), id); err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	return nil
}
