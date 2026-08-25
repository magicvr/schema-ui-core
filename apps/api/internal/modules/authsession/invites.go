// Admin invitation domain (workspace-019 R3 · GOAL-004 D-001 §3): the
// invite-creates-account flow with roles fixed AT ISSUANCE (user adjudication
// 2026-08-25 — 受邀账号角色以发布邀请时指定为准). The raw token is shown once
// (creation/resend); only its SHA-256 hash is stored. Lifecycle per frozen
// I-005: 7-day default validity, one-time consumption, instant revoke,
// resend replaces the token (60 s cooldown).
package authsession

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

const (
	defaultInviteTTL = 7 * 24 * time.Hour
	inviteResendCool = 60 * time.Second
	minInviteRoles   = 1 // an invite always carries its future roles
	maxInviteRoles   = 16
)

// Sentinels mapped to catalog codes by the HTTP surface.
var (
	ErrInviteNotFound = errors.New("authsession: invite not found")
	ErrInviteInvalid  = errors.New("authsession: invite is unknown, expired, consumed or revoked")
	ErrInviteCooldown = errors.New("authsession: invite resend cooldown active")
	ErrInviteRoleGone = errors.New("authsession: invited roles no longer exist; reissue the invite")
)

// Invite is one user_invites row (token material never re-read).
type Invite struct {
	ID         string
	Roles      []string
	InvitedBy  string
	Email      *string
	ExpiresAt  time.Time
	ConsumedAt *time.Time
	RevokedAt  *time.Time
	CreatedAt  time.Time
}

// newInviteToken mints an opaque invitation token: a 256-bit random URL-safe
// raw value shown exactly once, stored only as its SHA-256 hex hash. (The auth
// package cannot be imported here — it imports this module for account.User.)
func newInviteToken() (raw, hash string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", fmt.Errorf("invite token: %w", err)
	}
	raw = base64.RawURLEncoding.EncodeToString(b)
	hash = hashInviteToken(raw)
	return raw, hash, nil
}

func hashInviteToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func (i *Invite) live(now time.Time) bool {
	return i.ConsumedAt == nil && i.RevokedAt == nil && now.Before(i.ExpiresAt)
}

func scanInvite(row kernel.Row) (*Invite, error) {
	var inv Invite
	var rolesJSON string
	var email *string
	var consumed, revoked *int64
	var expires, sent, created int64
	err := row.Scan(&inv.ID, &rolesJSON, &inv.InvitedBy, &email, &expires, &consumed, &revoked, &sent, &created)
	if errors.Is(err, kernel.ErrNoRows) {
		return nil, ErrInviteNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan invite: %w", err)
	}
	if err := json.Unmarshal([]byte(rolesJSON), &inv.Roles); err != nil {
		return nil, fmt.Errorf("invite roles: %w", err)
	}
	inv.Email = email
	inv.ExpiresAt = time.Unix(expires, 0)
	inv.CreatedAt = time.Unix(created, 0)
	if consumed != nil {
		t := time.Unix(*consumed, 0)
		inv.ConsumedAt = &t
	}
	if revoked != nil {
		t := time.Unix(*revoked, 0)
		inv.RevokedAt = &t
	}
	return &inv, nil
}

const inviteColumns = `id, roles, invited_by, email, expires_at, consumed_at, revoked_at, last_sent_at, created_at`

// CreateInvite validates the requested roles, persists the invite and returns
// the ONE-TIME raw token. Mail dispatch belongs to the handler layer (it owns
// base-URL composition); a dispatch failure there compensates by revoking.
func (r *Repository) CreateInvite(invitedBy string, roles []string, email string, ttl time.Duration, now time.Time) (string, *Invite, error) {
	roles = dedupeKeys(roles)
	if len(roles) < minInviteRoles || len(roles) > maxInviteRoles {
		return "", nil, ErrInviteRoleGone
	}
	for _, key := range roles {
		var count int
		err := r.withTx("check invite role", func(tx kernel.Tx) error {
			return tx.QueryRow(context.Background(),
				`SELECT COUNT(*) FROM roles WHERE key = ?`, key).Scan(&count)
		})
		if err != nil || count == 0 {
			return "", nil, ErrInviteRoleGone
		}
	}
	raw, hash, err := newInviteToken()
	if err != nil {
		return "", nil, err
	}
	rolesJSON, err := json.Marshal(roles)
	if err != nil {
		return "", nil, err
	}
	id := randomHexID("inv")
	var emailArg any
	if strings.TrimSpace(email) != "" {
		emailArg = strings.TrimSpace(email)
	}
	err = r.withTx("create invite", func(tx kernel.Tx) error {
		_, err := tx.Exec(context.Background(),
			`INSERT INTO user_invites (id, token_hash, roles, invited_by, email, expires_at, last_sent_at, created_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			id, hash, string(rolesJSON), invitedBy, emailArg,
			now.Add(ttl).Unix(), now.Unix(), now.Unix())
		return err
	})
	if err != nil {
		return "", nil, err
	}
	inv, err := r.GetInvite(id)
	if err != nil {
		return "", nil, err
	}
	return raw, inv, nil
}

// GetInvite loads one row by id.
func (r *Repository) GetInvite(id string) (*Invite, error) {
	var inv *Invite
	err := r.withTx("get invite", func(tx kernel.Tx) error {
		row, err := getInviteRow(tx, id)
		if err != nil {
			return err
		}
		inv, err = scanInvite(row)
		return err
	})
	return inv, err
}

func getInviteRow(tx kernel.Tx, id string) (kernel.Row, error) {
	return tx.QueryRow(context.Background(),
		`SELECT `+inviteColumns+` FROM user_invites WHERE id = ?`, id), nil
}

// ListInvites returns newest-first paged invites for the admin surface.
func (r *Repository) ListInvites(page, pageSize int) ([]Invite, int, error) {
	var out []Invite
	var total int
	err := r.withTx("list invites", func(tx kernel.Tx) error {
		if err := tx.QueryRow(context.Background(),
			`SELECT COUNT(*) FROM user_invites`).Scan(&total); err != nil {
			return err
		}
		rows, err := tx.Query(context.Background(),
			`SELECT `+inviteColumns+` FROM user_invites ORDER BY created_at DESC LIMIT ? OFFSET ?`,
			pageSize, (page-1)*pageSize)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			inv, err := scanInvite(rows)
			if err != nil {
				return err
			}
			out = append(out, *inv)
		}
		return rows.Err()
	})
	return out, total, err
}

// RevokeInvite invalidates a pending invite instantly (I-005 撤销即时失效).
// Consumed/already-revoked invites revoke as a no-op (idempotent admin UX).
func (r *Repository) RevokeInvite(id string, now time.Time) error {
	return r.withTx("revoke invite", func(tx kernel.Tx) error {
		// Existence first (unknown id must surface as ErrInviteNotFound).
		var exists int
		if err := tx.QueryRow(context.Background(),
			`SELECT COUNT(*) FROM user_invites WHERE id = ?`, id).Scan(&exists); err != nil {
			return fmt.Errorf("exists: %w", err)
		}
		if exists == 0 {
			return ErrInviteNotFound
		}
		result, err := tx.Exec(context.Background(),
			`UPDATE user_invites SET revoked_at = ? WHERE id = ? AND consumed_at IS NULL AND revoked_at IS NULL`,
			now.Unix(), id)
		if err != nil {
			return fmt.Errorf("update: %w", err)
		}
		if _, rerr := result.RowsAffected(); rerr != nil {
			return rerr
		}
		return nil
	})
}

// ResendInvite rotates the token on a still-live invite (撤旧发新：the old link
// dies because the stored hash is replaced) and refreshes the expiry window,
// honouring the 60 s resend cooldown. Expired-but-unconsumed invites may be
// resent (fresh window); consumed ones are dead forever.
func (r *Repository) ResendInvite(id string, ttl time.Duration, now time.Time) (string, *Invite, error) {
	raw, hash, err := newInviteToken()
	if err != nil {
		return "", nil, err
	}
	err = r.withTx("resend invite", func(tx kernel.Tx) error {
		var consumed, revoked *int64
		var lastSent int64
		err := tx.QueryRow(context.Background(),
			`SELECT consumed_at, revoked_at, last_sent_at FROM user_invites WHERE id = ?`,
			id).Scan(&consumed, &revoked, &lastSent)
		if errors.Is(err, kernel.ErrNoRows) {
			return ErrInviteNotFound
		}
		if err != nil {
			return err
		}
		if consumed != nil {
			return ErrInviteInvalid
		}
		if revoked != nil {
			return ErrInviteInvalid
		}
		if now.Unix()-lastSent < int64(inviteResendCool/time.Second) {
			return ErrInviteCooldown
		}
		_, uerr := tx.Exec(context.Background(),
			`UPDATE user_invites SET token_hash = ?, expires_at = ?, last_sent_at = ? WHERE id = ?`,
			hash, now.Add(ttl).Unix(), now.Unix(), id)
		return uerr
	})
	if err != nil {
		return "", nil, err
	}
	inv, err := r.GetInvite(id)
	if err != nil {
		return "", nil, err
	}
	return raw, inv, nil
}

// AcceptInvite redeems a raw token into a real account: live checks → role
// re-validation against the CURRENT role table (deleted roles fail closed) →
// username uniqueness → account creation WITH the invite's roles and a
// must_change_password=false flag (the invitee chose their own password) →
// one-time consumption. Everything lands in ONE transaction so a consumed
// invite can never yield a half-created account. passwordHash arrives already
// policy-checked and hashed (handler layer owns hashing).
func (r *Repository) AcceptInvite(rawToken, username, name, passwordHash string, now time.Time) (*User, error) {
	hash := hashInviteToken(rawToken)
	cleanName := strings.TrimSpace(name)
	if cleanName == "" {
		cleanName = strings.TrimSpace(username)
	}
	var created *User
	err := r.withTx("accept invite", func(tx kernel.Tx) error {
		inv, serr := scanInvite(tx.QueryRow(context.Background(),
			`SELECT `+inviteColumns+` FROM user_invites WHERE token_hash = ?`, hash))
		if errors.Is(serr, ErrInviteNotFound) {
			// Unknown hash ≡ expired/consumed/revoked on the pre-auth surface:
			// one uniform invalid answer (no oracle).
			return ErrInviteInvalid
		}
		if serr != nil {
			return serr
		}
		if !inv.live(now) {
			return ErrInviteInvalid
		}
		// Role re-validation: an admin-deleted role makes the invite stale.
		roleSet := make(map[string]bool, len(inv.Roles))
		for _, key := range inv.Roles {
			var count int
			if err := tx.QueryRow(context.Background(),
				`SELECT COUNT(*) FROM roles WHERE key = ?`, key).Scan(&count); err != nil {
				return fmt.Errorf("invite role %s: %w", key, err)
			}
			if count == 0 {
				return ErrInviteRoleGone
			}
			roleSet[key] = true
		}
		// Username uniqueness inside the same transaction (fail-closed).
		var taken int
		if err := tx.QueryRow(context.Background(),
			`SELECT COUNT(*) FROM users WHERE lower(username) = lower(?)`,
			strings.TrimSpace(username)).Scan(&taken); err != nil {
			return fmt.Errorf("check username: %w", err)
		}
		if taken > 0 {
			return ErrUsernameTaken
		}
		u := User{
			ID:           randomHexID("user"),
			Username:     strings.TrimSpace(username),
			Name:         cleanName,
			Roles:        inv.Roles,
			PasswordHash: passwordHash,
			Enabled:      true,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		rolesBytes, merr := json.Marshal(u.Roles)
		if merr != nil {
			return merr
		}
		if _, err := tx.Exec(context.Background(),
			`INSERT INTO users (id, username, name, roles, password_hash, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`,
			u.ID, u.Username, u.Name, string(rolesBytes), u.PasswordHash, now.Unix(), now.Unix()); err != nil {
			return fmt.Errorf("insert invited user: %w", err)
		}
		for key := range roleSet {
			if _, err := tx.Exec(context.Background(),
				`INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)`, u.ID, "role-"+key); err != nil {
				return fmt.Errorf("link invited role %s: %w", key, err)
			}
		}
		if _, err := tx.Exec(context.Background(),
			`UPDATE user_invites SET consumed_at = ? WHERE id = ? AND consumed_at IS NULL`, now.Unix(), inv.ID); err != nil {
			return fmt.Errorf("consume invite: %w", err)
		}
		created = &u
		return nil
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}
