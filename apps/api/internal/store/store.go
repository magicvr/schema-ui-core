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
	db *sql.DB
}

const schema = `
CREATE TABLE IF NOT EXISTS users (
  id            TEXT PRIMARY KEY,
  username      TEXT NOT NULL UNIQUE,
  name          TEXT NOT NULL,
  roles         TEXT NOT NULL, -- JSON array; R3 normalizes
  password_hash TEXT NOT NULL,
  created_at    INTEGER NOT NULL,
  updated_at    INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS refresh_tokens (
  id         TEXT PRIMARY KEY,
  user_id    TEXT NOT NULL REFERENCES users(id),
  token_hash TEXT NOT NULL UNIQUE,
  expires_at INTEGER NOT NULL,
  revoked_at INTEGER,
  created_at INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
`

// Open opens (creating if needed) the SQLite DB at path, applies the idempotent
// schema, and — when seedAdmin is true and the users table is empty — seeds the
// admin user with the given bcrypt password hash. The caller is responsible for
// enforcing a non-empty seed password in production (fail-closed on startup).
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

	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	if seedAdmin {
		if err := s.seedAdmin(adminUsername, adminPasswordHash); err != nil {
			db.Close()
			return nil, err
		}
	}
	return s, nil
}

// Close releases the underlying database handle.
func (s *Store) Close() error { return s.db.Close() }

func (s *Store) migrate() error {
	if _, err := s.db.Exec(schema); err != nil {
		return fmt.Errorf("apply schema: %w", err)
	}
	return nil
}

// SeedAdmin inserts the bootstrap admin user when no users exist. It is
// idempotent: once any user exists, it is a no-op.
func (s *Store) seedAdmin(username, passwordHash string) error {
	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count); err != nil {
		return fmt.Errorf("count users: %w", err)
	}
	if count > 0 {
		return nil
	}
	now := time.Now().UTC()
	roles, err := json.Marshal([]string{"admin", "editor"})
	if err != nil {
		return fmt.Errorf("marshal seed roles: %w", err)
	}
	_, err = s.db.Exec(
		`INSERT INTO users (id, username, name, roles, password_hash, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		"user-admin", username, "Admin", string(roles), passwordHash, now.Unix(), now.Unix(),
	)
	if err != nil {
		return fmt.Errorf("seed admin: %w", err)
	}
	return nil
}

// CreateUser inserts a new user.
func (s *Store) CreateUser(u User) error {
	roles, err := json.Marshal(u.Roles)
	if err != nil {
		return fmt.Errorf("marshal roles: %w", err)
	}
	_, err = s.db.Exec(
		`INSERT INTO users (id, username, name, roles, password_hash, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		u.ID, u.Username, u.Name, string(roles), u.PasswordHash, u.CreatedAt.Unix(), u.UpdatedAt.Unix(),
	)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

// UserByUsername fetches a user by unique username.
func (s *Store) UserByUsername(username string) (*User, error) {
	return s.scanUser(s.db.QueryRow(
		`SELECT id, username, name, roles, password_hash, created_at, updated_at
		 FROM users WHERE username = ?`, username))
}

// UserByID fetches a user by primary key.
func (s *Store) UserByID(id string) (*User, error) {
	return s.scanUser(s.db.QueryRow(
		`SELECT id, username, name, roles, password_hash, created_at, updated_at
		 FROM users WHERE id = ?`, id))
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
